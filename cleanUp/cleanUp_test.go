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
