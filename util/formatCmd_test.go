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
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFmtCommand(t *testing.T) {
	path, pathErr := os.Getwd()
	if pathErr != nil {
		t.Error(pathErr)
	}
	cmd, err := FmtCommand("bash -c 'NODE_ENV=production gulp --gulpfile=./client/js/gulpfile.js scripts'", path)
	if err != nil {
		t.Error(err)
	}
	expected := []string{BASH, "-c", fmt.Sprintf("'NODE_ENV=production gulp --gulpfile=%s/client/js/gulpfile.js scripts'", path)}
	for i := range cmd.Args {
		if strings.Compare(cmd.Args[i], expected[i]) != 0 {
			t.Errorf("expected %s, found %s", expected[i], cmd.Args[i])
		}
	}
}
