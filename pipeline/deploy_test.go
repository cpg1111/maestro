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

package pipeline

import (
	"testing"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/util"
)

var testDeploySrv = &Service{
	conf: config.Service{
		Name:      "test",
		Tag:       "0.1.0",
		TagType:   "git",
		Path:      ".",
		BuildCMD:  []string{"echo build"},
		CheckCMD:  []string{"echo check >> test_file"},
		CreateCMD: []string{"echo create >> test_file"},
		UpdateCMD: []string{"echo update >> test_file"},
	},
	shouldBuild:   true,
	logFileOffset: 1,
	lastCommit:    "fake-commit",
	currCommit:    "fake-commit",
	path:          ".",
	diffPath:      ".",
	HasFailed:     false,
}

func TestCheck(t *testing.T) {
	check(testDeploySrv)
	existErr := util.CheckForFile("test_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("test_file", "check\nupdate\n")
	if err != nil {
		t.Error(err)
	}
}

func TestCreate(t *testing.T) {
	create(testDeploySrv)
	existErr := util.CheckForFile("test_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("test_file", "create\n")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	update(testDeploySrv)
	existErr := util.CheckForFile("test_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("test_file", "update\n")
	if err != nil {
		t.Error(err)
	}
}
