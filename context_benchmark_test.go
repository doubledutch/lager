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
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	stout "github.com/pivotal-golang/lager"
)

func BenchmarkJSONContextLagerOneLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := NewContextLager(&ContextConfig{
		Levels:      new(Levels).Set(Error),
		Drinker:     NewJSONDrinker(f),
		Stacktraces: false,
	})

	lgr.Set("app", "benchmark")
	lgr.Set("type", "lager")

	msg := "test"

	for i := 0; i < b.N; i++ {
		lgr.Tracef("This is a %s", msg)
		lgr.Debugf("This is a %s", msg)
		lgr.Infof("This is a %s", msg)
		lgr.Warnf("This is a %s", msg)
		lgr.Errorf("This is a %s", msg)
	}
}

func BenchmarkContextPivotalLagerOneLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := stout.NewLogger("my-app")

	lgr.RegisterSink(stout.NewWriterSink(f, stout.INFO))

	msg := "pivotal lager"
	boomErr := errors.New("boom")

	ctxt := stout.Data{
		"app":  "benchmark",
		"type": "stout",
	}

	for i := 0; i < b.N; i++ {
		// No trace, treat as debug
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), ctxt)
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), ctxt)
		lgr.Info(fmt.Sprintf("benchmark %s", msg), ctxt)
		// No Warn, treat as info
		lgr.Info(fmt.Sprintf("benchmark %s", msg), ctxt)
		lgr.Error(fmt.Sprintf("benchmark %s", msg), boomErr, ctxt)
	}
}

func BenchmarkContextLogrusOneLevel(b *testing.B) {
	// Doesn't work with /dev/null
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	logger := logrus.New()
	logger.Out = f
	logger.Level = logrus.ErrorLevel

	msg := "logrus"
	data := logrus.Fields{
		"app":  "benchmark",
		"type": "logrus",
	}

	for i := 0; i < b.N; i++ {
		logger.WithFields(data).Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Info(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Warn(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Error(fmt.Sprintf("The %s breaks!", msg))
	}
}

func BenchmarkJSONContextLagerAllLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := NewContextLager(&ContextConfig{
		Levels:      new(Levels).Set(Error | Warn | Info | Debug | Trace),
		Drinker:     NewJSONDrinker(f),
		Stacktraces: false,
	})

	lgr.Set("app", "benchmark")
	lgr.Set("type", "yes")

	msg := "test"

	for i := 0; i < b.N; i++ {
		lgr.Tracef("This is a %s", msg)
		lgr.Debugf("This is a %s", msg)
		lgr.Infof("This is a %s", msg)
		lgr.Warnf("This is a %s", msg)
		lgr.Errorf("This is a %s", msg)
	}
}

func BenchmarkContextPivotalLagerAllLevels(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := stout.NewLogger("my-app")

	// Uses minLevel
	lgr.RegisterSink(stout.NewWriterSink(f, stout.DEBUG))

	msg := "pivotal lager"
	boomErr := errors.New("boom")

	ctxt := stout.Data{
		"app":  "benchmark",
		"type": "stout",
	}

	for i := 0; i < b.N; i++ {
		// No trace, treat as debug
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), ctxt)
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), ctxt)
		lgr.Info(fmt.Sprintf("benchmark %s", msg), ctxt)
		// No Warn, treat as info
		lgr.Info(fmt.Sprintf("benchmark %s", msg), ctxt)
		lgr.Error(fmt.Sprintf("benchmark %s", msg), boomErr, ctxt)
	}
}

func BenchmarkContextLogrusAllLevels(b *testing.B) {
	// Doesn't work with /dev/null
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	logger := logrus.New()
	logger.Out = f
	logger.Level = logrus.DebugLevel

	msg := "logrus"

	data := logrus.Fields{
		"app":  "benchmark",
		"type": "logrus",
	}

	for i := 0; i < b.N; i++ {
		logger.Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Info(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Warn(fmt.Sprintf("The %s breaks!", msg))
		logger.WithFields(data).Error(fmt.Sprintf("The %s breaks!", msg))
	}
}
