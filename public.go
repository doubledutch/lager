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

var defaultLager = NewContextLager(nil)

// SetLevels sets the levels of the package lager
func SetLevels(levels *Levels) ContextLager {
	defaultLager.SetLevels(levels)
	return defaultLager
}

// SetDrinker sets the drinker of the package lager
func SetDrinker(drinker Drinker) ContextLager {
	defaultLager.(*contextLager).drinker = drinker
	return defaultLager
}

// Tracef logs with level Trace using the package lager
func Tracef(msg string, v ...interface{}) {
	defaultLager.Tracef(msg, v...)
}

// Debugf logs with level Debug using the package lager
func Debugf(msg string, v ...interface{}) {
	defaultLager.Debugf(msg, v...)
}

// Infof logs with level Info using the package lager
func Infof(msg string, v ...interface{}) {
	defaultLager.Infof(msg, v...)
}

// Warnf logs with level Warn using the package lager
func Warnf(msg string, v ...interface{}) {
	defaultLager.Warnf(msg, v...)
}

// Errorf logs with level Error using the package lager
func Errorf(msg string, v ...interface{}) {
	defaultLager.Errorf(msg, v...)
}

// Set sets a key to value in the lager map  using the package lager
func Set(key, value string) ContextLager {
	return defaultLager.Set(key, value)
}

// Child creates a child ContextLager using the package lager as the parent.
// The child inherits all the parent values.
func Child() ContextLager {
	return defaultLager.Child()
}
