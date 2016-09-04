package statecom

import (
	"github.com/cpg1111/maestro/config"
)

type ServiceStateMgr struct {
	Name      string
	State     string
	HasFailed bool
	Parent    *StateCom
}

func NewServiceStateMgr(conf config.Service, parent *StateCom) *ServiceStateMgr {
	return &ServiceStateMgr{
		Name:      conf.Name,
		State:     "pending",
		HasFailed: false,
		Parent:    parent,
	}
}

func (s *ServiceStateMgr) SetState(state string, success bool) {
	s.State = state
	s.HasFailed = !success
	s.Parent.SetServiceState(s)
}
