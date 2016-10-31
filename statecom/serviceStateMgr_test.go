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
	srvConf, srvConfErr = config.Load("../test_conf.toml", "/tmp/test/")
	srvStateCom         = New(srvConf, "", "")
	serviceMgr          = NewServiceStateMgr(srvConf.Services[0], srvStateCom)
)

func TestSetState(t *testing.T) {
	if srvConfErr != nil {
		t.Error(srvConfErr)
	}
	serviceMgr.SetState("test", true)
	if serviceMgr.State != "test" {
		t.Error("serviceMgr's state is not 'test'")
	}
	if serviceMgr.HasFailed {
		t.Error("serviceMgr has failed")
	}
}
