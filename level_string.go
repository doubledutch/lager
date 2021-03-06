/*
Copyright 2015 Doubledutch
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

// generated by stringer -type=Level; DO NOT EDIT

package lager

import "fmt"

const (
	_Level_name_0 = "Error"
	_Level_name_1 = "Warn"
	_Level_name_2 = "Info"
	_Level_name_3 = "Debug"
	_Level_name_4 = "Trace"
)

var (
	_Level_index_0 = [...]uint8{0, 5}
	_Level_index_1 = [...]uint8{0, 4}
	_Level_index_2 = [...]uint8{0, 4}
	_Level_index_3 = [...]uint8{0, 5}
	_Level_index_4 = [...]uint8{0, 5}
)

func (i Level) String() string {
	switch {
	case i == 1:
		return _Level_name_0
	case i == 4:
		return _Level_name_1
	case i == 8:
		return _Level_name_2
	case i == 16:
		return _Level_name_3
	case i == 32:
		return _Level_name_4
	default:
		return fmt.Sprintf("Level(%d)", i)
	}
}
