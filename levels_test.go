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

import "testing"

func TestLevels(t *testing.T) {
	levels := new(Levels)

	levels.Set(Error)

	if !levels.Contains(Error) {
		t.Fatal("levels should contain error")
	}

	levels.Set(Trace)

	if !levels.Contains(Error) || !levels.Contains(Trace) {
		t.Fatal("levels should contain error and trace")
	}
}

func TestAllLevels(t *testing.T) {
	levels := new(Levels)

	testLevels := []Level{Error, Warn, Info, Debug, Trace}

	for _, v := range testLevels {
		levels.Set(v)
	}

	for _, v := range testLevels {
		if !levels.Contains(v) {
			t.Fatalf("levels didn't contain %s", v)
		}
	}
}

func TestLevelsFromString(t *testing.T) {
	verifyLevelsString(t, "E", []Level{Error}, []Level{Warn, Info, Trace, Debug})

	verifyLevelsString(t, "EI", []Level{Error, Info}, []Level{Warn, Trace, Debug})

	verifyLevelsString(t, "", []Level{}, []Level{Warn, Info, Trace, Debug, Error})

	verifyLevelsString(t, "EWITD", []Level{Error, Warn, Info, Trace, Debug}, []Level{})

}

func verifyLevelsString(t *testing.T, sLevels string, contains, excludes []Level) {
	levels := LevelsFromString(sLevels)

	for _, contain := range contains {
		if !levels.Contains(contain) {
			t.Fatalf("levels contains %s", contain)
		}
	}

	for _, exclude := range excludes {
		if levels.Contains(exclude) {
			t.Fatalf("levels includes %s", exclude)
		}
	}
}

func TestLevelsAll(t *testing.T) {
	levels := new(Levels).All()

	includes := []Level{Debug, Trace, Info, Warn, Error}
	for _, include := range includes {
		if !levels.Contains(include) {
			t.Fatalf("levels doesn't contain %s", include)
		}
	}
}
