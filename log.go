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
	"io"
	"log"
	"os"
)

// LogConfig represents for a Lager
type LogConfig struct {
	Levels *Levels
	Output io.Writer
}

// DefaultLogConfig is the default config
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Levels: new(Levels).Set(Error),
		Output: os.Stdout,
	}
}

// LogLager implements Logger using *log.Logger
type LogLager struct {
	Lager
	logger *log.Logger
}

// NewLogLager creates a new LogLager
func NewLogLager(config *LogConfig) Lager {
	if config == nil {
		config = DefaultLogConfig()
	}

	logger := &LogLager{
		logger: log.New(config.Output, "", log.LstdFlags),
	}

	logger.Lager = newLager(logger, config.Levels)

	return logger
}

// Logf will log the given msg formatted with v if min is greater than or equal
// to the log level of LogLager
func (lgr *LogLager) Logf(level Level, msg string, v ...interface{}) {
	lgr.logger.Printf(msg, v...)
}
