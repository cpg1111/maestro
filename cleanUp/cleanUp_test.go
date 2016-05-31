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

package cleanUp

import (
	"testing"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/util"
)

var conf = config.CleanUp{
	AdditionalCMDs: []string{"docker ps -a"},
	InDaemon:       false,
	Artifacts: []config.Artifact{
		config.Artifact{
			RuntimeFilePath: "./dist/maestro",
			SaveFilePath:    "/tmp/",
		},
	},
}

func TestHandleCMDs(t *testing.T) {
	cmdErr := handleCMDs(conf.AdditionalCMDs)
	if cmdErr != nil {
		t.Error(cmdErr)
	}
}

func TestSaveArtifacts(t *testing.T) {
	artifactErr := saveArtifacts(conf.Artifacts)
	if artifactErr != nil {
		t.Error(artifactErr)
	}
	for i := range conf.Artifacts {
		checkErr := util.CheckForFile(conf.Artifacts[i].SaveFilePath)
		if checkErr != nil {
			t.Error(checkErr)
		}
	}
}
