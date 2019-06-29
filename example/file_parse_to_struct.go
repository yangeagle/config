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
	//	Mac1	string		`kitty:"mac1"`
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

type Conf struct {
	Hostname string    `kitty:"host"`
	Addr     string    `kitty:"ipaddr"`
	PortNum  int       `kitty:"port"`
	PortNum1 *int      `kitty:"port"`
	Status   string    `kitty:"compression"`
	Height   []float32 `kitty:"height"`
	Active   bool      `kitty:"active"`
	Clusters []string  `kitty:"cluster"`
	Dist     int       `kitty:"distance"`
	Temp     float64   `kitty:"temprature"`
	Level    int       `kitty:"top_level"`
	NumConn  int       `kitty:"max_conn"`
	PortSW   bool      `kitty:"port_enable"`
	Order    []int     `kitty:"order"`
	//	Monitor		*MonitorConf	`kitty:"monitor"`
	Monitor MonitorConf `kitty:"monitor"`
	Portal  PortalConf  `kitty:"portal"`
}

const configfile = "test.conf"

func main() {

	conf := kitty.NewConfig()

	err := conf.ParseFile(configfile)
	if err != nil {
		fmt.Println("parse eror:", err)
		return
	}

	conf_info := new(Conf)
	conf_info.Dist = 90

	var a int = 2
	conf_info.PortNum1 = &a

	conf_info.Portal.Web1 = new(WebConf)

	//no use
	/*for _, v := range conf_info.Portal.Clusters1 {
		v = new(ClusterConf)
	}*/

	err_u := conf.Unmarshal(conf_info)
	if err_u != nil {
		fmt.Println("err_u:", err_u)
		return
	}

	fmt.Println("conf_info:", conf_info)

}
