/*
Copyright 2016 Christian Grabowski All rights reserved.
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

package util

import (
	"testing"
)

// TestFmtDiffPath tests FmtDiffPath
func TestFmtDiffPath(t *testing.T) {
	testVals := []string{
		"/usr/lib",
		"usr/lib/",
		"/usr/lib/",
		"usr/lib",
	}
	for i := range testVals {
		res := FmtDiffPath("/usr/", testVals[i])
		if res != "lib/" {
			t.Errorf("expected lib/ found %s", res)
		}
	}
}

// TestFmtClonePath tests FmtDiffPath
func TestFmtClonePath(t *testing.T) {
	testPath1 := "/tmp/build/"
	testPathPtr1 := &testPath1
	res1 := FmtClonePath(testPathPtr1)
	if *res1 != "/tmp/build" {
		t.Errorf("expected /tmp/build found %s", *res1)
	}
	testPath2 := "/tmp/build"
	testPathPtr2 := &testPath2
	res2 := FmtClonePath(testPathPtr2)
	if *res2 != "/tmp/build" {
		t.Errorf("expected /tmp/build found %s", *res2)
	}
}
