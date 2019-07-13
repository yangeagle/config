/*
   Copyright 2019 Yang

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
limitations under the License.
*/

package kitty

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

const (
	tagName = "kitty"
)

type section map[string]interface{}

/*
Config type declaration
*/
type Config struct {
	allConfigOptions section   // the top level section, including all key-value
	order            []section // current section in every level
}

/*
* NewConfig create a instance of Config
 */
func NewConfig() *Config {
	conf := Config{}

	conf.allConfigOptions = section{}
	conf.order = make([]section, 0)
	conf.order = append(conf.order, conf.allConfigOptions)

	return &conf
}

/*
open config file and parse
*/
func (c *Config) ParseFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	err = c.parse(file)
	if err != nil {
		return err
	}

	return nil
}

// parse config string
func (c *Config) ParseBytes(conf []byte) error {

	buf := bytes.NewBuffer(conf)

	if err := c.parse(buf); err != nil {
		return err
	}

	return nil
}

/*
* return all key-value in map
 */
func (c *Config) GetAllConfigOptions() section {
	return c.allConfigOptions
}

func (c *Config) getCurrentLevel(s string) (int, error) {
	tabCount, err := getTabCount(s)
	if err != nil {
		return 0, err
	}

	return tabCount + 1, nil

}

//get current section
func (c *Config) getCurrentSection(currentLevel int) (section, error) {

	levelCount := len(c.order)
	if currentLevel > levelCount {
		return nil, fmt.Errorf("syntax error: please check indentation.")
	}

	return c.order[currentLevel-1], nil
}

// update the level if already exist, or add a new level
func (c *Config) updateOrAddLevel(currentLevel int) (section, error) {
	levelCount := len(c.order)
	if currentLevel > levelCount {
		return nil, fmt.Errorf("syntax error: please check indentation.")
	}

	newSection := section{}

	if currentLevel < levelCount {
		//update
		c.order[currentLevel] = newSection
	} else if currentLevel == levelCount {
		//new level
		c.order = append(c.order, newSection)

	}

	return newSection, nil
}

/*
read string line by line

special symbol definition:
#			comment
=			simple assignment, key=value
[]			`section`
[[]]			[]`section`, the array of `section`

*/
func (c *Config) parse(reader io.Reader) error {
	var (
		err  error
		row  string
		rd   = bufio.NewReader(reader)
		line int
	)

	for {
		line++

		row, err = rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(row) == 0 {
					// file end
					break
				}
			} else {
				// error
				return err
			}
		}

		rowNospace := strings.TrimSpace(row)
		// comment (line start with "#") or none line
		if strings.HasPrefix(rowNospace, "#") || len(rowNospace) == 0 {
			continue
		}

		currentLevel, err := c.getCurrentLevel(row)
		if err != nil {
			return fmt.Errorf("getCurrentLevel failed in line %d:%s\n", line, err)
		}

		currentSection, err := c.getCurrentSection(currentLevel)
		if err != nil {
			return fmt.Errorf("getCurrentSection failed in line %d:%s\n", line, err)
		}

		row = rowNospace
		rowLen := len(row)

		if strings.HasPrefix(row, "[[") && strings.HasSuffix(row, "]]") { //section[[]]
			newSection, err := c.updateOrAddLevel(currentLevel)
			if err != nil {
				return fmt.Errorf("updateOrAddLevel failed in line %d:%s\n", line, err)
			}

			key := row[2 : rowLen-2]
			key = strings.TrimSpace(key)

			var newSlice []section
			value, ok := currentSection[key]
			if !ok {
				// new slice
				newSlice = []section{newSection}
			} else {
				// slice already exist and then just append
				sliceTmp, ok := value.([]section)
				if !ok {
					return fmt.Errorf("%v not slice type in line %d", value, line)
				}

				newSlice = append(sliceTmp, newSection)
			}

			// save
			currentSection[key] = newSlice
		} else if strings.HasPrefix(row, "[") && strings.HasSuffix(row, "]") { //section []
			key := row[1 : rowLen-1]
			key = strings.TrimSpace(key)

			value, ok := currentSection[key]
			if ok {
				return fmt.Errorf("%s:%s duplicate definition in line %d", key, value, line)
			}

			newSection, err := c.updateOrAddLevel(currentLevel)
			if err != nil {
				return fmt.Errorf("updateOrAddLevel failed in line %d:%s\n", line, err)
			}

			// save
			currentSection[key] = newSection
		} else if index := strings.Index(row, "="); index >= 0 { //"="
			key := row[:index]
			key = strings.TrimSpace(key)

			value, ok := currentSection[key]
			if ok {
				return fmt.Errorf("%s:%s already exist", key, value)
			}

			valueTmp := row[index+1:]
			valueTmp = strings.TrimSpace(valueTmp)
			// save
			currentSection[key] = valueTmp

		} else {
			return fmt.Errorf("missing separator in line %d", line)
		}

	} //for

	return nil
}

// parse content, and storage to v
func (c *Config) Unmarshal(v interface{}) error {
	refValue := reflect.ValueOf(v)

	//pointer check
	if refValue.Kind() != reflect.Ptr || refValue.IsNil() {
		return fmt.Errorf("pointer type needed and not nil.")
	}

	//derefer pointer
	elemValue := refValue.Elem()

	if err := section2Struct(c.allConfigOptions, elemValue); err != nil {
		return err
	}

	return nil
}

// section to struct
func section2Struct(sec section, refValue reflect.Value) error {
	//get type
	valueType := refValue.Type()

	//pointer check
	if refValue.Kind() == reflect.Ptr {
		if refValue.IsNil() {
			// new
			// element type, eg. *struct --> struct
			elemType := valueType.Elem()
			// *struct
			newValue := reflect.New(elemType)
			refValue.Set(newValue)
		}

		//derefer pointer
		refValue = refValue.Elem()
		// element type
		valueType = refValue.Type()
	}

	// fields
	n := refValue.NumField()
	for i := 0; i < n; i++ {
		fieldType := valueType.Field(i)
		fieldValue := refValue.Field(i)

		//get tag
		fieldTag := fieldType.Tag.Get(tagName)

		optionValue, ok := sec[fieldTag]
		if !ok {
			//this config option not exist in config file
			continue
		}

		err := setOptionValue2RefValue(fieldValue, optionValue)
		if err != nil {
			return err
		}

	} // for n

	return nil
}

/*
description:
	set optionValue to refValue

input:
	1. refValue
	2. optionValue
output:
	error
*/
func setOptionValue2RefValue(refValue reflect.Value, optionValue interface{}) error {

	basicType := refValue.Kind()

	switch basicType {
	case reflect.String:
		fallthrough
	case reflect.Bool:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		optionValueStr, ok := optionValue.(string)
		if !ok {
			return fmt.Errorf("unknown string:%v", optionValue)
		}

		newValue, err := convert2Value(basicType.String(), optionValueStr)
		if err != nil {
			return err
		}

		// set value
		refValue.Set(newValue)
	case reflect.Ptr:
		if refValue.IsNil() {
			// new
			//element type, eg. *int --> int
			elemType := refValue.Type().Elem()
			// pointer
			newValue := reflect.New(elemType)
			refValue.Set(newValue)
		}

		//derefer pointer
		elemValue := refValue.Elem()

		err := setOptionValue2RefValue(elemValue, optionValue)
		if err != nil {
			return err
		}
	case reflect.Struct:
		sec, ok := optionValue.(section)
		if !ok {
			return fmt.Errorf("unknown struct")
		}

		err := section2Struct(sec, refValue)
		if err != nil {
			return err
		}
	case reflect.Slice:
		// value type
		valueType := refValue.Type()
		// element type, eg. []string --> string
		elemType := valueType.Elem()

		// the slice of string
		if sliceStr, ok := optionValue.(string); ok { //simple type
			sep := ","
			strs := strings.Split(sliceStr, sep)
			// new slice
			newSlice := reflect.MakeSlice(valueType, 0, len(strs))

			// every element
			for _, str := range strs {
				str = strings.TrimSpace(str)

				// pointer
				newValue := reflect.New(elemType)
				newElemValue := newValue.Elem()

				err := setOptionValue2RefValue(newElemValue, str)
				if err != nil {
					return err
				}

				// append
				newSlice = reflect.Append(newSlice, newElemValue)
			}

			refValue.Set(newSlice)
		} else if sliceSection, ok := optionValue.([]section); ok { //composite type
			// new slice
			newSlice := reflect.MakeSlice(valueType, 0, len(sliceSection))

			// every element
			for _, sec := range sliceSection {
				//pointer
				newValue := reflect.New(elemType)
				newElemValue := newValue.Elem()

				err := section2Struct(sec, newElemValue)
				if err != nil {
					return err
				}

				newSlice = reflect.Append(newSlice, newElemValue)
			}

			refValue.Set(newSlice)
		} else {
			return fmt.Errorf("unknown string:%v", optionValue)
		}

	default:
		return fmt.Errorf("unknown basic type:%v", basicType)
	}

	return nil
}
