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

type WebConf struct {
	IP  string `config:"ip"`
	Mac string `config:"mac"`
}

type ClusterConf struct {
	Ipaddr string `config:"addr"`
	Weight int    `config:"wgh"`
}

type PortalConf struct {
	Enabled     bool           `config:"enabled"`
	Ip          string         `config:"ip"`
	Port        int            `config:"port"`
	Clusters    []ClusterConf  `config:"cluster"`
	Clusters1   []*ClusterConf `config:"cluster"`
	ConnTimeout int            `config:"connecttimeout"`
	CallTimeout int            `config:"calltimeout"`
	Web         WebConf        `config:"web"`
	Web1        *WebConf       `config:"web1"`
}

type MacConf struct {
	Mac1 *string `config:"mac1"`
	Mac2 string  `config:"mac2"`
}

type MonitorConf struct {
	Enabled bool      `config:"enabled"`
	IP      string    `config:"ip"`
	Macs    MacConf   `config:"MAC"`
	Port    int       `config:"port"`
	Cluster []*string `config:"clust"`
}

type ConfigOption struct {
	Hostname string      `config:"host"`
	Addr     string      `config:"ipaddr"`
	PortNum  int         `config:"port"`
	Status   string      `config:"compression"`
	Height   []float32   `config:"height"`
	Active   bool        `config:"active"`
	Clusters []string    `config:"cluster"`
	Dist     int         `config:"distance"`
	Temp     float64     `config:"temprature"`
	Level    *int        `config:"top_level"`
	NumConn  int         `config:"max_conn"`
	PortSW   bool        `config:"port_enable"`
	Order    []int       `config:"order"`
	Monitor  MonitorConf `config:"monitor"`
	Portal   PortalConf  `config:"portal"`
}

const configFile = "test.conf"

func main() {

	parser := config.NewConfig()

	err := parser.ParseFile(configFile)
	if err != nil {
		fmt.Println("ParseFile failed:", err)
		return
	}

	confOption := new(ConfigOption)

	err = parser.Unmarshal(confOption)
	if err != nil {
		fmt.Println("Unmarshal failed:", err)
		return
	}

	fmt.Println("Hostname:", confOption.Hostname)
	fmt.Println("Addr:", confOption.Addr)
	fmt.Println("Port:", confOption.PortNum)
	fmt.Println("Status:", confOption.Status)
	fmt.Println("Height:", confOption.Height)
	fmt.Println("Active:", confOption.Active)
	fmt.Println("Clusters:", confOption.Clusters)
	fmt.Println("Dist:", confOption.Dist)
	fmt.Println("Temp:", confOption.Temp)
	fmt.Println("TopLevel:", *confOption.Level)
	fmt.Println("NumConn:", confOption.NumConn)
	fmt.Println("PortSW:", confOption.PortSW)
	fmt.Println("Order:", confOption.Order)
	fmt.Println("Monitor:", confOption.Monitor)
	fmt.Println("Portal:", confOption.Portal)
}
