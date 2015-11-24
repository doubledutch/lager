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

// Lager is a Lager that explicitly defines all of the log levels
// as log methods.
type Lager interface {
	Tracef(msg string, v ...interface{})
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	SetLevels(levels *Levels)
	Levels() *Levels
}

// palgeLager represents a logger
type paleLager interface {
	Logf(log Level, msg string, v ...interface{})
	SetLevels(levels *Levels)
	Levels() *Levels
}

type lager struct {
	pale paleLager

	levels *Levels
}

func newLager(lgr paleLager, levels *Levels) Lager {
	return &lager{
		pale:   lgr,
		levels: levels,
	}
}

// Tracef logs with level Trace
func (lgr *lager) Tracef(msg string, v ...interface{}) {
	lgr.logf(Trace, msg, v...)
}

// Debugf logs with level Debug
func (lgr *lager) Debugf(msg string, v ...interface{}) {
	lgr.logf(Debug, msg, v...)
}

// Infof logs with level Info
func (lgr *lager) Infof(msg string, v ...interface{}) {
	lgr.logf(Info, msg, v...)
}

// Warnf logs with level Warn
func (lgr *lager) Warnf(msg string, v ...interface{}) {
	lgr.logf(Warn, msg, v...)
}

// Errorf logs with level Error
func (lgr *lager) Errorf(msg string, v ...interface{}) {
	lgr.logf(Error, msg, v...)
}

func (lgr *lager) logf(lvl Level, msg string, v ...interface{}) {
	if lgr.levels == nil {
		return
	}

	if !lgr.levels.Contains(lvl) {
		return
	}

	lgr.pale.Logf(lvl, msg, v...)
}

func (lgr *lager) SetLevels(levels *Levels) {
	lgr.levels.Replace(levels)
}

func (lgr *lager) Levels() *Levels {
	return lgr.levels
}
