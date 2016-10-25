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
