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
	"io/ioutil"
	"testing"
)

func testLager(level Level) error {
	buf := new(bytes.Buffer)

	lgr := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(level),
		Drinker: NewLogDrinker(buf),
	})

	newError := func(level Level) error {
		return fmt.Errorf("wasn't expecting log of type %s", level)
	}

	lgr.Tracef("This is %s", Trace)
	if actual, _ := ioutil.ReadAll(buf); Trace == level && len(actual) == 0 {
		return newError(Trace)
	}
	lgr.Debugf("This is %s", Debug)
	if actual, _ := ioutil.ReadAll(buf); Debug == level && len(actual) == 0 {
		return newError(Debug)
	}
	lgr.Infof("This is %s", Info)
	if actual, _ := ioutil.ReadAll(buf); Info == level && len(actual) == 0 {
		return newError(Info)
	}
	lgr.Warnf("This is %s", Warn)
	if actual, _ := ioutil.ReadAll(buf); Warn == level && len(actual) == 0 {
		return newError(Warn)
	}
	lgr.Errorf("This is %s", Error)
	if actual, _ := ioutil.ReadAll(buf); Error == level && len(actual) == 0 {
		return newError(Error)
	}

	return nil
}

func TestLagerTrace(t *testing.T) {
	if err := testLager(Trace); err != nil {
		t.Fatal(err)
	}
}

func TestLagerDebug(t *testing.T) {
	if err := testLager(Debug); err != nil {
		t.Fatal(err)
	}
}

func TestLagerInfo(t *testing.T) {
	if err := testLager(Info); err != nil {
		t.Fatal(err)
	}
}

func TestLagerWarn(t *testing.T) {
	if err := testLager(Warn); err != nil {
		t.Fatal(err)
	}
}

func TestLagerError(t *testing.T) {
	if err := testLager(Error); err != nil {
		t.Fatal(err)
	}
}
