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
		CheckCMD:  []string{"touch check_file"},
		CreateCMD: []string{"touch create_file"},
		UpdateCMD: []string{"touch update_file"},
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
	existErr := util.CheckForFile("check_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("check_file", "")
	if err != nil {
		t.Error(err)
	}
}

func TestCreate(t *testing.T) {
	create(testDeploySrv)
	existErr := util.CheckForFile("create_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("create_file", "")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	update(testDeploySrv)
	existErr := util.CheckForFile("update_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("update_file", "")
	if err != nil {
		t.Error(err)
	}
}
