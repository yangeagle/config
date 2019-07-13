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

package main

import (
	"fmt"

	"github.com/kitty"
)

var configBytes []byte = []byte(`
#comment like this
host = example.com

port = 43
compression = on

#comment like this

active = true
group = ["abc", 192.168.1.10, 172.10.156.23, what are you doing]

Order = [1, 2, 3, 110]

timeout = 100
temprature = 36.8

height = 8848.996

`)

func main() {

	parser := kitty.NewConfig()

	if err := parser.ParseBytes(configBytes); err != nil {
		fmt.Printf("parse str failed: %s", err)
		return
	}

	allConfigOptions := parser.GetAllConfigOptions()

	for option, value := range allConfigOptions {
		fmt.Println(option, ":", value)
	}
}
