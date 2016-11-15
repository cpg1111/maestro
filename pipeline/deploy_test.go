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
		CheckCMD:  []string{"touch ./test_file"},
		CreateCMD: []string{"bash -c \"touch ./test_file && echo 'create' >> ./test_file\""},
		UpdateCMD: []string{"bash -c \"touch ./test_file && echo 'update' >> ./test_file\""},
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
	existErr := util.CheckForFile("./test_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("./test_file", "check\nupdate\n")
	if err != nil {
		t.Error(err)
	}
}

func TestCreate(t *testing.T) {
	create(testDeploySrv)
	existErr := util.CheckForFile("./test_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("./test_file", "create\n")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	update(testDeploySrv)
	existErr := util.CheckForFile("./test_file")
	if existErr != nil {
		t.Error(existErr)
	}
	err := util.CheckFileContents("./test_file", "update\n")
	if err != nil {
		t.Error(err)
	}
}
