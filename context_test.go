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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDefaultContextConfig(t *testing.T) {
	if DefaultContextConfig() == nil {
		t.Fatal("DefaultContextConfig() == nil")
	}
}

func TestContextJSONLogf(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(Trace),
		Drinker: NewJSONDrinker(buf),
	})

	logger.Set("hello", "world")

	expected := "this is a test"
	logger.Tracef("this is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err = json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	if logMap["msg"] != expected {
		t.Fatalf("actual: '%s' to contain expected: '%s'", logMap["message"], expected)
	}

	if logMap["hello"] != "world" {
		t.Fatalf("expected hello == equal")
	}
}

func TestContextJSONLogfErrorStacktrace(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:      new(Levels).Set(Error),
		Drinker:     NewJSONDrinker(buf),
		Stacktraces: true,
	})

	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("%s", r.(error))
			}
			done <- struct{}{}
		}()

		panic(errors.New("PANIC"))
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err := json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	if len(logMap["stacktrace"]) == 0 {
		t.Fatal("expected stacktrace")
	}
}

func TestContextJSONNotLogf(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(Error),
		Drinker: NewJSONDrinker(buf),
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

func TestContextChild(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(Error),
		Drinker: NewJSONDrinker(buf),
	})

	keys := map[string]string{
		"a": "one",
		"b": "two",
		"c": "three",
	}

	logger.Set("a", keys["a"])

	child1 := logger.Child()
	child1.Set("b", keys["b"])

	child2 := child1.Child()
	child2.Set("c", keys["c"])

	child2.Errorf("this is a test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err = json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	for k, v := range keys {
		if logMap[k] != v {
			t.Fatalf("expected key '%s' == '%s', got '%s'", k, v, logMap[k])
		}
	}
}

func TestContextJSONNoFile(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(Trace),
		Drinker: NewJSONDrinker(buf),
	})

	logger.Tracef("this is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err = json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	if logMap["file"] != "" {
		t.Fatalf("expected no file to be logged")
	}
}

func TestContextJSONShortFile(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:   new(Levels).Set(Trace),
		Drinker:  NewJSONDrinker(buf),
		FileType: ShortFile,
	})

	logger.Tracef("this is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err = json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	if logMap["file"] == "" {
		t.Fatalf("expected file to be logged")
	}

	parts := strings.Split(logMap["file"], ":")
	if parts[0] != "context_test.go" {
		t.Fatalf("expected %s, actual %s", "context_test.go", parts[0])
	}
}

func TestContextJSONPackageFile(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:   new(Levels).Set(Trace),
		Drinker:  NewJSONDrinker(buf),
		FileType: PackageFile,
	})

	logger.Tracef("this is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err = json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	if logMap["file"] == "" {
		t.Fatalf("expected file to be logged")
	}

	parts := strings.Split(logMap["file"], string(os.PathSeparator))
	if len(parts) != 4 {
		t.Fatalf("Expected 4 parts in the path, actual %d", len(parts))
	}
	expected := []string{"github.com", "doubledutch", "lager"}
	for i, expect := range expected {
		if parts[i] != expect {
			t.Fatalf("expected %s, actual %s", expect, parts[i])
		}
	}

	parts = strings.Split(parts[3], ":")
	if parts[0] != "context_test.go" {
		t.Fatalf("expected %s, actual %s", "context_test.go", parts[0])
	}
}

func TestContextJSONFullFile(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:   new(Levels).Set(Trace),
		Drinker:  NewJSONDrinker(buf),
		FileType: FullFile,
	})

	logger.Tracef("this is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err = json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	if logMap["file"] == "" {
		t.Fatalf("expected file to be logged")
	}

	parts := strings.Split(logMap["file"], string(os.PathSeparator))
	if len(parts) <= 4 {
		t.Fatalf("Expected more than 4 parts to the path, actual %d", len(parts))
	}
}

func TestContextJSONChildFile(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewContextLager(&ContextConfig{
		Levels:      new(Levels).Set(Trace),
		Drinker:     NewJSONDrinker(buf),
		FileType:    FullFile,
		Stacktraces: true,
	})

	child := logger.Child()
	child.Tracef("this is a %s", "test")

	actual, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	var logMap map[string]string
	if err = json.Unmarshal(actual, &logMap); err != nil {
		t.Fatal(err)
	}

	if logMap["file"] == "" {
		t.Fatalf("expected file to be logged")
	}
}

func TestContextLager(t *testing.T) {
	NewContextLager(nil)
}
