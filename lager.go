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
}

// PaleLager represents a logger
type paleLager interface {
	Logf(log Level, msg string, v ...interface{})
}

type lager struct {
	pale paleLager
}

func newLager(lgr paleLager) Lager {
	return &lager{
		pale: lgr,
	}
}

// Tracef logs with level Trace
func (lgr *lager) Tracef(msg string, v ...interface{}) {
	lgr.pale.Logf(Trace, msg, v...)
}

// Debugf logs with level Debug
func (lgr *lager) Debugf(msg string, v ...interface{}) {
	lgr.pale.Logf(Debug, msg, v...)
}

// Infof logs with level Info
func (lgr *lager) Infof(msg string, v ...interface{}) {
	lgr.pale.Logf(Info, msg, v...)
}

// Warnf logs with level Warn
func (lgr *lager) Warnf(msg string, v ...interface{}) {
	lgr.pale.Logf(Warn, msg, v...)
}

// Errorf logs with level Error
func (lgr *lager) Errorf(msg string, v ...interface{}) {
	lgr.pale.Logf(Error, msg, v...)
}
