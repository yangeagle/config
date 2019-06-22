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
	data  section   // the top level section, including all key-value
	order []section // current section in every level
}

/*
* NewConfig create a instance of Config
 */
func NewConfig() *Config {
	conf := Config{}

	conf.data = section{}
	conf.order = make([]section, 0)
	conf.order = append(conf.order, conf.data)

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
func (c *Config) ReturnAll() section {
	return c.data
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

			//for test
			fmt.Println("--->", currentSection)
		} else {
			return fmt.Errorf("missing separator in line %d", line)
		}

	} //for

	return nil
}
