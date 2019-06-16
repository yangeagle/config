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
	"bytes"
	"io"
	"os"
)

const (
	tagName = "kitty"
)

type section map[string]interface{}

/*
Config type declaration
*/
type Config struct {
	data  section    // all key-value
	order []*section // level
}

/*
* NewConfig create a instance of Config
 */
func NewConfig() *Config {

	conf := Config{}

	conf.data = section{}
	conf.order = make([]*section, 0)
	conf.order = append(conf.order, &(conf.data))

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

/*
read string line by line

special symbol definition:
#			comment
=			simple assignment			key=value
[]			`section`					`struct`
[[]]		[]`section`					the array of `struct`

*/
func (c *Config) parse(reader io.Reader) error {
	return nil
}
