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

package statecom

import (
	"testing"

	"github.com/cpg1111/maestro/config"
)

var (
	conf, confErr = config.Load("../test_conf.toml", "/tmp/test")
	stateCom      = New(conf, "", "")
)

func TestStart(t *testing.T) {
	if confErr != nil {
		t.Error(confErr)
	}
	stateCom.Start()
	if stateCom.Global.StateLabel != "started" {
		t.Error("state is not 'started")
	}
}

func TestEnv(t *testing.T) {
	if confErr != nil {
		t.Error(confErr)
	}
	stateCom.Env()
	if stateCom.Global.StateLabel != "creating env" {
		t.Error("state is not 'creating env'")
	}
}

func TestCloning(t *testing.T) {
	if confErr != nil {
		t.Error(confErr)
	}
	stateCom.Cloning()
	if stateCom.Global.StateLabel != "cloning repo" {
		t.Error("state is not 'cloning repo'")
	}
}

func TestCleanUp(t *testing.T) {
	if confErr != nil {
		t.Error(confErr)
	}
	stateCom.CleanUp()
	if stateCom.Global.StateLabel != "clean up" {
		t.Error("state is not 'clean up'")
	}
}

func TestDone(t *testing.T) {
	if confErr != nil {
		t.Error(confErr)
	}
	stateCom.Done()
	if stateCom.Global.StateLabel != "done" {
		t.Error("state is not 'done'")
	}
}
