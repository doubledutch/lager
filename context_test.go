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
		t.Fatalf("Error unmarshalling log: %s", err)
	}

	if logMap["msg"] != expected {
		t.Fatalf("expected '%s', got '%s'", expected, logMap["msg"])
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
		t.Fatalf("Error unmarshalling json: %s", err)
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
	dec := json.NewDecoder(buf)

	logger := NewContextLager(&ContextConfig{
		Levels:      new(Levels).Set(Trace),
		Drinker:     NewJSONDrinker(buf),
		FileType:    FullFile,
		Stacktraces: true,
	})

	child := logger.Child()
	child.Tracef("this is a %s", "test")

	var logMap map[string]string
	if err := dec.Decode(&logMap); err != nil {
		t.Fatal(err)
	}

	if logMap["file"] == "" {
		t.Fatalf("expected file to be logged")
	}
}

func TestContextWith(t *testing.T) {
	buf := new(bytes.Buffer)
	dec := json.NewDecoder(buf)

	logger := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(Trace),
		Drinker: NewJSONDrinker(buf),
	})

	logger.Set("global", "peace")
	logger.Tracef("hello world")

	var first map[string]string
	err := dec.Decode(&first)
	if err != nil {
		t.Fatal(err)
	}

	if first["global"] != "peace" {
		t.Fatalf("expected global to be peace, got %s", first["global"])
	}

	unexpectedKey := func(key string, m map[string]string) {
		if _, ok := m[key]; ok {
			t.Fatalf("unexpected key %s", key)
		}
	}

	unexpectedKey("a", first)
	unexpectedKey("b", first)

	logger.With(map[string]string{
		"a": "one",
		"b": "two",
	}).Tracef("hello again, world")

	var second map[string]string
	err = dec.Decode(&second)
	if err != nil {
		t.Fatal(err)
	}

	if second["global"] != "peace" {
		t.Fatalf("expected global to be peace, got %s", second["global"])
	}

	if second["a"] != "one" {
		t.Fatalf("expected a to be one, got %s", second["a"])
	}

	if second["b"] != "two" {
		t.Fatalf("expected b to be two, got %s", second["b"])
	}

	logger.Tracef("good night, world")
	var third map[string]string
	err = dec.Decode(&third)
	if err != nil {
		t.Fatal(err)
	}

	if third["global"] != "peace" {
		t.Fatalf("expected global to be peace, got %s", third["global"])
	}

	unexpectedKey("a", third)
	unexpectedKey("b", third)
}

func TestContextWithErr(t *testing.T) {
	buf := new(bytes.Buffer)
	dec := json.NewDecoder(buf)

	logger := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(Trace),
		Drinker: NewJSONDrinker(buf),
	})

	err := errors.New("failure")
	logger.WithError(err).Tracef("hello world")

	var actual map[string]string
	if err := dec.Decode(&actual); err != nil {
		t.Fatal(err)
	}

	if actual["error"] != err.Error() {
		t.Fatalf("expected error to be failure, got %s", actual["error"])
	}
}

func TestContextWithErrNil(t *testing.T) {
	buf := new(bytes.Buffer)
	dec := json.NewDecoder(buf)

	logger := NewContextLager(&ContextConfig{
		Levels:  new(Levels).Set(Trace),
		Drinker: NewJSONDrinker(buf),
	})

	var err error
	logger.WithError(err).Tracef("hello world")

	var actual map[string]string
	if err := dec.Decode(&actual); err != nil {
		t.Fatal(err)
	}

	if _, ok := actual["error"]; ok {
		t.Fatalf("expected no error")
	}
}

func TestContextLager(t *testing.T) {
	NewContextLager(nil)
}
