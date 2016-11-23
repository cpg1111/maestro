package pipeline

import (
	"testing"

	"github.com/cpg1111/maestro/config"
)

var testRunnerService = &Service{
	conf: config.Service{
		Name:             "testRunnerService",
		Tag:              "0.1.0",
		TagType:          "git",
		Path:             "./",
		BuildLogFilePath: "/var/log/maestro/test.log",
		BuildCMD:         []string{"echo build"},
		TestCMD:          []string{"echo test"},
		CheckCMD:         []string{"echo check"},
		CreateCMD:        []string{"echo create"},
		UpdateCMD:        []string{"echo update"},
	},
	State:       "created",
	shouldBuild: true,
	path:        "./",
	diffPath:    "./",
	HasFailed:   false,
}

func TestRunTests(t *testing.T) {
	err := RunTests(testRunnerService)
	if err != nil {
		t.Error(err)
	}
}
