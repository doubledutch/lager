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

import "sync"

// Level represents a logging level
type Level uint

//go:generate stringer -type=Level
const (
	// Trace = a log level for tracing
	Trace Level = 1 << 5
	// Debug = a log level for debugging
	Debug Level = 1 << 4
	// Info = a log level for informational messages
	Info Level = 1 << 3
	// Warn = a log level for warnings
	Warn Level = 1 << 2
	// Error = a log level for errors
	Error Level = 1
)

// Levels contains the log levels a logger will write to a Drinker
type Levels struct {
	bits Level
	lock sync.RWMutex
}

// LevelsFromString creates a levels object from a string
// Levels are specified using a capital letter corresponding
// to the first level of the desired level.
func LevelsFromString(sLevels string) *Levels {
	levels := new(Levels)
	for _, sLevel := range sLevels {
		switch sLevel {
		case 'E':
			levels.Set(Error)
		case 'W':
			levels.Set(Warn)
		case 'I':
			levels.Set(Info)
		case 'T':
			levels.Set(Trace)
		case 'D':
			levels.Set(Debug)
		}
	}

	return levels
}

// Set sets a log level
func (lvls *Levels) Set(level Level) *Levels {
	lvls.lock.Lock()
	lvls.bits |= level
	lvls.lock.Unlock()

	return lvls
}

// Unset removes a log level
func (lvls *Levels) Unset(level Level) *Levels {
	lvls.lock.Lock()
	lvls.bits &= ^level
	lvls.lock.Unlock()

	return lvls
}

// Contains checks to see if a log level is contained in a logger
func (lvls *Levels) Contains(level Level) bool {
	lvls.lock.RLock()
	defer lvls.lock.RUnlock()
	return lvls.bits&level == level
}

// All sets all levels
func (lvls *Levels) All() *Levels {
	lvls.Set(Trace | Debug | Info | Warn | Error)

	return lvls
}

// Replace changes it's value to match level
func (lvls *Levels) Replace(level *Levels) *Levels {
	lvls.lock.Lock()
	lvls.bits = level.bits
	lvls.lock.Unlock()

	return lvls
}
