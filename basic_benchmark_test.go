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

	// Packages we'll be benchmarking
	"log"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/logutils"
	stout "github.com/pivotal-golang/lager"
)

func BenchmarkBasicLagerOneLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := NewBasicLager(&BasicConfig{
		Levels:  new(Levels).Set(Error),
		Drinker: NewLogDrinker(f),
	})

	msg := "test"

	for i := 0; i < b.N; i++ {
		lgr.Tracef("This is a %s", msg)
		lgr.Debugf("This is a %s", msg)
		lgr.Infof("This is a %s", msg)
		lgr.Warnf("This is a %s", msg)
		lgr.Errorf("This is a %s", msg)
	}
}

func BenchmarkLogUtilsOneLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "ERROR",
		Writer:   f,
	}
	log.SetOutput(filter)

	msg := "test"
	for i := 0; i < b.N; i++ {
		log.Printf("[TRACE] This is a %s", msg)
		log.Printf("[DEBUG] This is a %s", msg)
		log.Printf("[INFO] This is a %s", msg)
		log.Printf("[WARN] This is a %s", msg)
		log.Printf("[ERROR] This is a %s", msg)
	}
}

func BenchmarkPivotalLagerOneLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := stout.NewLogger("my-app")

	lgr.RegisterSink(stout.NewWriterSink(f, stout.INFO))

	msg := "pivotal lager"
	boomErr := errors.New("boom")

	for i := 0; i < b.N; i++ {
		// No trace, treat as debug
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), nil)
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), nil)
		lgr.Info(fmt.Sprintf("benchmark %s", msg), nil)
		// No Warn, treat as info
		lgr.Info(fmt.Sprintf("benchmark %s", msg), nil)
		lgr.Error(fmt.Sprintf("benchmark %s", msg), boomErr, nil)
	}
}

func BenchmarkLogrusOneLevel(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		logger.Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.Info(fmt.Sprintf("The %s breaks!", msg))
		logger.Warn(fmt.Sprintf("The %s breaks!", msg))
		logger.Error(fmt.Sprintf("The %s breaks!", msg))
	}
}

func BenchmarkBasicLagerAllLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := NewBasicLager(&BasicConfig{
		Levels:  new(Levels).Set(Error | Warn | Info | Debug | Trace),
		Drinker: NewLogDrinker(f),
	})

	msg := "test"

	for i := 0; i < b.N; i++ {
		lgr.Tracef("This is a %s", msg)
		lgr.Debugf("This is a %s", msg)
		lgr.Infof("This is a %s", msg)
		lgr.Warnf("This is a %s", msg)
		lgr.Errorf("This is a %s", msg)
	}
}

func BenchmarkLogUtilsAllLevels(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	// Uses minLevel
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "TRACE",
		Writer:   f,
	}
	log.SetOutput(filter)

	msg := "test"
	for i := 0; i < b.N; i++ {
		log.Printf("[TRACE] This is a %s", msg)
		log.Printf("[DEBUG] This is a %s", msg)
		log.Printf("[INFO] This is a %s", msg)
		log.Printf("[WARN] This is a %s", msg)
		log.Printf("[ERROR] This is a %s", msg)
	}
}

func BenchmarkPivotalLagerAllLevels(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		// No trace, treat as debug
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), nil)
		lgr.Debug(fmt.Sprintf("benchmark %s", msg), nil)
		lgr.Info(fmt.Sprintf("benchmark %s", msg), nil)
		// No Warn, treat as info
		lgr.Info(fmt.Sprintf("benchmark %s", msg), nil)
		lgr.Error(fmt.Sprintf("benchmark %s", msg), boomErr, nil)
	}
}

func BenchmarkLogrusAllLevels(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		logger.Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.Debug(fmt.Sprintf("The %s breaks!", msg))
		logger.Info(fmt.Sprintf("The %s breaks!", msg))
		logger.Warn(fmt.Sprintf("The %s breaks!", msg))
		logger.Error(fmt.Sprintf("The %s breaks!", msg))
	}
}
