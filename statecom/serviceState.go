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
