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

	"github.com/yangeagle/config"
)

type ConfigOption struct {
	Hostname string    `config:"host"`
	Addr     string    `config:"ipaddr"`
	PortNum  int       `config:"port"`
	Height   []float32 `config:"height"`
	Active   bool      `config:"active"`
	Clusters []string  `config:"cluster"`
	Dist     int       `config:"distance"`
	Temp     float64   `config:"temprature"`
	TopLevel *int      `config:"top_level"`
	NumConn  int       `config:"max_conn"`
	Order    []int     `config:"order"`
}

const configFile = "simple.conf"

func main() {

	confParser := config.NewConfig()

	err := confParser.ParseFile(configFile)
	if err != nil {
		fmt.Println("ParseFile failed:", err)
		return
	}

	confOption := new(ConfigOption)

	err = confParser.Unmarshal(confOption)
	if err != nil {
		fmt.Println("Unmarshal failed:", err)
		return
	}

	fmt.Println("Hostname:", confOption.Hostname)
	fmt.Println("Addr:", confOption.Addr)
	fmt.Println("Port:", confOption.PortNum)
	fmt.Println("Height:", confOption.Height)
	fmt.Println("Active:", confOption.Active)
	fmt.Println("Clusters:", confOption.Clusters)
	fmt.Println("Dist:", confOption.Dist)
	fmt.Println("Temp:", confOption.Temp)
	fmt.Println("TopLevel:", *confOption.TopLevel)
	fmt.Println("NumConn:", confOption.NumConn)
	fmt.Println("Order:", confOption.Order)
}
