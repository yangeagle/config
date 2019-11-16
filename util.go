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

package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/*
* get the count of tab
 */
func getTabCount(s string) (int, error) {
	var count int

	if len(s) == 0 {
		return -1, fmt.Errorf("empty string")
	}

	originString := s
	newString := strings.TrimSpace(originString)

	if index := strings.IndexRune(originString, rune(newString[0])); index > 0 {
		spacePart := originString[:index]
		tabLen := len("\t")
		count = len(spacePart) / tabLen
	}

	return count, nil
}

/*
description:
	1.parse simple type
	2.convert string to reflect.Value
intput:
	1. typ: type in string format
	2. str: string to convert

output:
	1. reflect.Value
	2. error
*/
func convert2Value(typ, str string) (reflect.Value, error) {
	var value reflect.Value

	switch typ {
	case "bool":
		boolTmp, err := strconv.ParseBool(str)
		if err != nil {
			return value, err
		}

		value = reflect.ValueOf(boolTmp)
	case "int":
		intTmp, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return value, err
		}

		value = reflect.ValueOf(int(intTmp))
	case "float32":
		float32Tmp, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return value, nil
		}

		value = reflect.ValueOf(float32(float32Tmp))
	case "float64":
		float64Tmp, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return value, err
		}

		value = reflect.ValueOf(float64(float64Tmp))
	case "string":
		value = reflect.ValueOf(str)
	default:
		return value, fmt.Errorf("unknown type: %s", typ)
	}

	return value, nil
}
