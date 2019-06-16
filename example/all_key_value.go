package main

import (
	"fmt"

	"github.com/kitty"
)

func main() {

	conf := kitty.NewConfig()

	configFile := "./test.conf"

	if err := conf.ParseFile(configFile); err != nil {
		fmt.Printf("parse file %s failed", configFile)
		return
	}

	allKeyValue := conf.ReturnAll()
	fmt.Println("allKeyValue:", allKeyValue)
}
