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
	"fmt"
	"io"
)

// LogDrinker is a Drinker that uses log.Logger
type LogDrinker struct {
	output io.Writer
}

// NewLogDrinker creates a new Log Drinker
func NewLogDrinker(output io.Writer) Drinker {
	return &LogDrinker{
		output: output,
	}
}

// Drink drinks logs
func (drkr *LogDrinker) Drink(v interface{}) error {
	fmt.Fprintf(drkr.output, "%v\n", v)
	return nil
}
