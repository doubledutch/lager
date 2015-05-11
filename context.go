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
	"os"
	"runtime/debug"
	"time"
)

// ContextLager is a Lager that adds context to logs with key value pairs
type ContextLager interface {
	Lager
	Set(key, value string) ContextLager
	Child() ContextLager
}

// ContextConfig is defines the configuration for ContextLager.
// ContextConfig wraps Config.
type ContextConfig struct {
	Levels      *Levels
	Drinker     Drinker
	Values      map[string]string
	Stacktraces bool
}

// DefaultContextConfig creates a default ContextConfig
func DefaultContextConfig() *ContextConfig {
	return &ContextConfig{
		Levels:  new(Levels).Set(Error),
		Drinker: NewJSONDrinker(os.Stdout),
	}
}

type contextLager struct {
	Lager

	levels  *Levels
	drinker Drinker

	values      map[string]string
	stacktraces bool
}

// NewContextLager creates a JSONLager
func NewContextLager(config *ContextConfig) ContextLager {
	values := make(map[string]string)

	if config == nil {
		config = DefaultContextConfig()
	}

	//copy all keys and values into allValues
	for k, v := range config.Values {
		if _, ok := values[k]; !ok {
			values[k] = v
		}
	}

	lgr := &contextLager{
		levels:      config.Levels,
		drinker:     config.Drinker,
		values:      values,
		stacktraces: config.Stacktraces,
	}

	lgr.Lager = newLager(lgr)
	return lgr
}

// Set sets a key in the MapLager
func (lgr *contextLager) Set(key, value string) ContextLager {
	lgr.values[key] = value
	return lgr
}

//Logf writes a log to the standard output
func (lgr *contextLager) Logf(lvl Level, message string, v ...interface{}) {
	if !lgr.levels.Contains(lvl) {
		return
	}

	allValues := make(map[string]string)
	for k, v := range lgr.values {
		allValues[k] = v
	}

	if lvl == Error && lgr.stacktraces {
		allValues["stacktrace"] = string(debug.Stack())
	}

	//add all standard values
	allValues["@timestamp"] = time.Now().UTC().Format("Mon Jan 2 15:04:05 MST 2006")
	allValues["message"] = fmt.Sprintf(message, v...)
	allValues["level"] = lvl.String()

	//not sure what to do if the logger fails here
	lgr.drinker.Drink(allValues)
}

// Child creates a child ContextLager from this, the parent.
// The child inherits all the parent values.
func (lgr *contextLager) Child() ContextLager {
	return NewContextLager(&ContextConfig{
		Levels:  lgr.levels,
		Drinker: lgr.drinker,
		Values:  lgr.values,
	})
}