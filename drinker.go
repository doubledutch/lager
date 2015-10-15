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

package lager

import (
	"errors"
	"io"
)

// ErrNoDrinker is used when a drinker cannot be returned, primarly DrinkerFromString
var ErrNoDrinker = errors.New("No Drinker")

// Drinker will drink logs and output them
type Drinker interface {
	Drink(v map[string]interface{}) error
}

// NewDrinkerFunc creates a new drinker
type NewDrinkerFunc func(output io.Writer) Drinker

// DrinkerFromString provides a way to get a NewDrinkerFunc
// using a string. Useful for deciding on a Drinker using the environment.
func DrinkerFromString(str string) (NewDrinkerFunc, error) {
	switch str {
	case "JSON":
		return NewJSONDrinker, nil
	case "LOG":
		return NewLogDrinker, nil
	default:
		return nil, ErrNoDrinker
	}
}
