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
	"log"
	"os"
	"testing"
)

func BenchmarkLoggerOneAndAllLevels(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := log.New(f, "", log.LstdFlags)

	msg := "test"
	for i := 0; i < b.N; i++ {
		lgr.Printf("[TRACE] This is a %s", msg)
		lgr.Printf("[DEBUG] This is a %s", msg)
		lgr.Printf("[INFO] This is a %s", msg)
		lgr.Printf("[WARN] This is a %s", msg)
		lgr.Printf("[ERROR] This is a %s", msg)
	}
}

func BenchmarkLogLagerOneLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := NewLogLager(&LogConfig{
		Levels: new(Levels).Set(Error),
		Output: f,
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

func BenchmarkLogLagerAllLevel(b *testing.B) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	lgr := NewLogLager(&LogConfig{
		Levels: new(Levels).Set(Error | Warn | Info | Debug | Trace),
		Output: f,
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
