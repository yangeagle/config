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
	data  section   // all key-value
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

/*
read string line by line

special symbol definition:
#			comment
=			simple assignment			key=value
[]			`section`					`struct`
[[]]		[]`section`					the array of `struct`

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
		//none line or comment (line start with "#")
		if len(rowNospace) == 0 || strings.HasPrefix(rowNospace, "#") {
			continue
		}

		currentLevel, err := c.getCurrentLevel(row)
		if err != nil {
			fmt.Printf("getCurrentLevel failed in line %d:%s\n", line, err)
			return fmt.Errorf("unexpected error:%s", err)
		}

		levelCount := len(c.order)

		currentSection, err := c.getCurrentSection(currentLevel)
		if err != nil {
			fmt.Printf("getCurrentSection failed in line %d:%s\n", line, err)
			return fmt.Errorf("unexpected error:%s", err)
		}

		row = rowNospace
		rowLen := len(row)

		//section[[]]
		if strings.HasPrefix(row, "[[") && strings.HasSuffix(row, "]]") {
			// todo

		} else if strings.HasPrefix(row, "[") && strings.HasSuffix(row, "]") { //section []

			key := row[1 : rowLen-1]
			key = strings.TrimSpace(key)

			value, ok := currentSection[key]
			if !ok {
				newSection := section{}
				currentSection[key] = newSection

				if currentLevel < levelCount {
					// update current section
					c.order[currentLevel] = newSection
				} else if currentLevel == levelCount {
					// new level
					c.order = append(c.order, newSection)
				}

			} else {
				return fmt.Errorf("%s:%s duplicate definition in line %d", key, value, line)
			}

		} else if index := strings.Index(row, "="); index >= 0 { //spliter: "="

			key := row[:index]
			key = strings.TrimSpace(key)

			value, ok := currentSection[key]
			if !ok {
				valueTmp := row[index+1:]
				valueTmp = strings.TrimSpace(valueTmp)

				currentSection[key] = valueTmp

				//for test
				fmt.Println("--->", currentSection)
			} else {
				return fmt.Errorf("%s:%s already exist", key, value)

			}

		} else {
			return fmt.Errorf("missing symbol in line %d", line)
		}

	} //for

	return nil
}
