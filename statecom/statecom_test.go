package statecom

import (
	"testing"

	"github.com/cpg1111/maestro/config"
)

var (
	conf, confErr = config.Load("../test_conf.toml", "/tmp/test")
	stateCom      = New(conf, "")
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
