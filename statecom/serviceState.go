package statecom

import (
	"github.com/cpg1111/maestro/config"
)

// ServiceStateMgr is responsible for managing the state of a service
type ServiceStateMgr struct {
	Name      string
	State     string
	HasFailed bool
	Parent    *StateCom
}

// NewServiceStateMgr returns a pointer to a ServiceStateMgr struct
func NewServiceStateMgr(conf config.Service, parent *StateCom) *ServiceStateMgr {
	return &ServiceStateMgr{
		Name:      conf.Name,
		State:     "pending",
		HasFailed: false,
		Parent:    parent,
	}
}

// SetState sets the state of a service
func (s *ServiceStateMgr) SetState(state string, success bool) {
	s.State = state
	s.HasFailed = !success
	s.Parent.SetServiceState(s)
}
