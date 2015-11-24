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
	With(map[string]string) ContextLager
	WithError(error) ContextLager
	Set(key, value string) ContextLager
	Child() ContextLager
}

// ContextConfig is defines the configuration for ContextLager.
type ContextConfig struct {
	Levels      *Levels
	Drinker     Drinker
	Values      map[string]string
	Stacktraces bool
	FileType    FileType
}

// DefaultContextConfig creates a default ContextConfig
func DefaultContextConfig() *ContextConfig {
	return &ContextConfig{
		Levels:   new(Levels).Set(Error),
		Drinker:  NewJSONDrinker(os.Stdout),
		FileType: PackageFile,
	}
}

type contextLager struct {
	Lager

	drinker Drinker

	values      map[string]string
	stacktraces bool
	fileType    FileType
}

// NewContextLager creates a JSONLager
func NewContextLager(config *ContextConfig) ContextLager {
	values := make(map[string]string)

	if config == nil {
		config = DefaultContextConfig()
	}

	if config.Values == nil {
		config.Values = make(map[string]string)
	}

	//copy all keys and values into allValues
	for k, v := range config.Values {
		if _, ok := values[k]; !ok {
			values[k] = v
		}
	}

	logger := &contextLager{
		drinker:     config.Drinker,
		values:      values,
		stacktraces: config.Stacktraces,
		fileType:    config.FileType,
	}

	logger.Lager = newLager(logger, config.Levels)
	return logger
}

func (lgr *contextLager) With(fields map[string]string) ContextLager {
	if fields == nil {
		return lgr
	}

	clgr := lgr.Child()
	for key, value := range fields {
		clgr.Set(key, value)
	}
	return clgr
}

func (lgr *contextLager) WithError(err error) ContextLager {
	if err == nil {
		return lgr
	}

	clgr := lgr.Child()
	clgr.Set("error", err.Error())
	return clgr
}

// Set sets a key to value in the lager map
func (lgr *contextLager) Set(key, value string) ContextLager {
	lgr.values[key] = value
	return lgr
}

func (lgr *contextLager) Unset(key string) ContextLager {
	delete(lgr.values, key)
	return lgr
}

//Logf writes a log to the standard output
func (lgr *contextLager) Logf(lvl Level, message string, v ...interface{}) {
	allValues := make(map[string]interface{})
	for k, v := range lgr.values {
		allValues[k] = v
	}

	if lvl == Error && lgr.stacktraces {
		allValues["stacktrace"] = string(debug.Stack())
	}

	file := lgr.fileType.Caller(5)
	if file != "" {
		allValues["file"] = file
	}

	//add all standard values
	allValues["time"] = time.Now().UTC().Format(time.RFC3339)
	allValues["msg"] = fmt.Sprintf(message, v...)
	allValues["level"] = lvl.String()

	//not sure what to do if the logger fails here
	lgr.drinker.Drink(allValues)
}

// Child creates a child ContextLager from this, the parent.
// The child inherits all the parent values.
func (lgr *contextLager) Child() ContextLager {
	return NewContextLager(&ContextConfig{
		Levels:      lgr.Levels(),
		Drinker:     lgr.drinker,
		Values:      lgr.values,
		Stacktraces: lgr.stacktraces,
		FileType:    lgr.fileType,
	})
}
