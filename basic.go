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
)

// BasicConfig represents for a Lager
type BasicConfig struct {
	Levels  *Levels
	Drinker Drinker
}

// DefaultBasicConfig is the default config
func DefaultBasicConfig() *BasicConfig {
	return &BasicConfig{
		Levels:  new(Levels).Set(Error),
		Drinker: NewLogDrinker(os.Stdout),
	}
}

// BasicLager implements Logger using *log.Logger
type BasicLager struct {
	Lager
	levels  *Levels
	drinker Drinker
}

// NewBasicLager creates a new BasicLager
func NewBasicLager(config *BasicConfig) Lager {
	if config == nil {
		config = DefaultBasicConfig()
	}

	logger := &BasicLager{
		levels:  config.Levels,
		drinker: config.Drinker,
	}

	logger.Lager = newLager(logger)

	return logger
}

// Logf will log the given msg formatted with v if min is greater than or equal
// to the log level of LogLager
func (lgr *BasicLager) Logf(level Level, msg string, v ...interface{}) {
	if !lgr.levels.Contains(level) {
		return
	}

	lgr.drinker.Drink(fmt.Sprintf(msg, v...))
}
