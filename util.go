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
	"fmt"
	"strings"
)

/*
* get the count of tab
 */
func getTabCount(s string) (int, error) {

	if len(s) == 0 {
		return -1, fmt.Errorf("empty string")
	}

	var count int

	originString := s
	newString := strings.TrimSpace(originString)

	if index := strings.IndexRune(originString, rune(newString[0])); index > 0 {
		spacePart := originString[:index]
		tabLen := len("\t")
		count = len(spacePart) / tabLen
	}

	return count, nil
}
