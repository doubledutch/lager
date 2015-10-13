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
	"bytes"
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
func (drkr *LogDrinker) Drink(v map[string]interface{}) error {
	b := new(bytes.Buffer)

	for _, key := range []string{"time", "level", "msg"} {
		appendKeyValue(b, key, v[key])
		delete(v, key)
	}

	for key, value := range v {
		appendKeyValue(b, key, value)
	}

	b.WriteByte('\n')

	fmt.Fprint(drkr.output, string(b.Bytes()))
	return nil
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return false
		}
	}
	return true
}

func appendKeyValue(b *bytes.Buffer, key string, value interface{}) {

	b.WriteString(key)
	b.WriteByte('=')

	switch value := value.(type) {
	case string:
		if needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}

	b.WriteByte(' ')
}
