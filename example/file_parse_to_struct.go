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

type WebConf struct {
	IP  string `kitty:"ip"`
	Mac string `kitty:"mac"`
}

type ClusterConf struct {
	Ipaddr string `kitty:"addr"`
	Weight int    `kitty:"wgh"`
}

type PortalConf struct {
	Enabled     bool           `kitty:"enabled"`
	Ip          string         `kitty:"ip"`
	Port        int            `kitty:"port"`
	Clusters    []ClusterConf  `kitty:"cluster"`
	Clusters1   []*ClusterConf `kitty:"cluster"`
	ConnTimeout int            `kitty:"connecttimeout"`
	CallTimeout int            `kitty:"calltimeout"`
	Web         WebConf        `kitty:"web"`
	Web1        *WebConf       `kitty:"web1"`
}

type MacConf struct {
	Mac1 *string `kitty:"mac1"`
	Mac2 string  `kitty:"mac2"`
}

type MonitorConf struct {
	Enabled bool      `kitty:"enabled"`
	IP      string    `kitty:"ip"`
	Macs    MacConf   `kitty:"MAC"`
	Port    int       `kitty:"port"`
	Cluster []*string `kitty:"clust"`
}

type ConfigOption struct {
	Hostname string      `kitty:"host"`
	Addr     string      `kitty:"ipaddr"`
	PortNum  int         `kitty:"port"`
	Status   string      `kitty:"compression"`
	Height   []float32   `kitty:"height"`
	Active   bool        `kitty:"active"`
	Clusters []string    `kitty:"cluster"`
	Dist     int         `kitty:"distance"`
	Temp     float64     `kitty:"temprature"`
	Level    *int        `kitty:"top_level"`
	NumConn  int         `kitty:"max_conn"`
	PortSW   bool        `kitty:"port_enable"`
	Order    []int       `kitty:"order"`
	Monitor  MonitorConf `kitty:"monitor"`
	Portal   PortalConf  `kitty:"portal"`
}

const configFile = "test.conf"

func main() {

	parser := kitty.NewConfig()

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
