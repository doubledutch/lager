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
	"io/ioutil"
	"strings"
	"testing"
)

func TestLogLogf(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogLager(&LogConfig{
		Levels: new(Levels).Set(Error),
		Output: buf,
	})

	expected := "This is a test\n"
	logger.Errorf("This is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(actual), expected) {
		t.Fatalf("actual: '%s' to contain expected: '%s'", string(actual), expected)
	}
}

func TestLogNotLogf(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogLager(&LogConfig{
		Levels: new(Levels).Set(Error),
		Output: buf,
	})

	logger.Debugf("This is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) != 0 {
		t.Fatal("expected no logs")
	}
}

func TestDefaultLogConfig(t *testing.T) {
	if DefaultLogConfig() == nil {
		t.Fatal("DefaultLogConfig() == nil")
	}
}

func TestDefaultLog(t *testing.T) {
	NewLogLager(nil)
}
