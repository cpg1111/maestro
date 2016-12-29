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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cpg1111/maestro/config"
)

// State is the overall state
type State struct {
	Project    string
	Branch     string
	StateLabel string
	TimeStamp  time.Time
}

// ServiceState is the state of a service
type ServiceState struct {
	Name  string
	State *State
}

// StateCom is responsible for sending messages between Maestro and Maestrod for state
type StateCom struct {
	Project      string
	Branch       string
	Services     map[string]*ServiceStateMgr
	Global       *State
	client       *http.Client
	maestrodHost string
	maestrodPort int
}

func getMaestrodInfo(conf config.Project) (host string, port int) {
	if len(conf.MaestrodHostEnv) > 0 {
		host = os.Getenv(conf.MaestrodHostEnv)
	} else {
		host = conf.MaestrodHost
	}
	if len(conf.MaestrodPortEnv) > 0 {
		var pErr error
		port, pErr = strconv.Atoi(os.Getenv(conf.MaestrodPortEnv))
		if pErr != nil {
			panic(pErr)
		}
	} else {
		port = conf.MaestrodPort
	}
	return
}

// New returns a pointer to a StateCom struct
func New(conf config.Config, maestrodEndpoint, branch string) *StateCom {
	var client *http.Client
	if len(maestrodEndpoint) > 0 {
		client = &http.Client{}
	}
	mdHost, mdPort := getMaestrodInfo(conf.Project)
	stateCom := &StateCom{
		Project: conf.Project.RepoURL,
		Branch:  branch,
		Global: &State{
			Project:    conf.Project.RepoURL,
			Branch:     branch,
			StateLabel: "pending",
			TimeStamp:  time.Now(),
		},
		client:       client,
		maestrodHost: mdHost,
		maestrodPort: mdPort,
	}
	stateCom.Services = make(map[string]*ServiceStateMgr)
	for i := 0; i < len(conf.Services); i++ {
		stateCom.Services[conf.Services[i].Name] = NewServiceStateMgr(conf.Services[i], stateCom)
	}
	return stateCom
}

func (s *StateCom) send(state interface{}) {
	payload, marshErr := json.Marshal(state)
	if marshErr != nil {
		log.Println("WARNING", marshErr)
	}
	payloadRdr := bytes.NewReader(payload)
	resp, postErr := s.client.Post(
		fmt.Sprintf(
			"http://%s:%d/state",
			s.maestrodHost,
			s.maestrodPort,
		),
		"application/json",
		payloadRdr,
	)
	if postErr != nil {
		log.Println("WARNING", postErr.Error())
	}
	if resp != nil && resp.StatusCode != 201 {
		log.Println("WARNING STATEUPDATE NOT SENT")
	}
}

// Send sends the messages out to maestrod
func (s *StateCom) Send(state interface{}, bkg bool) {
	if s.client != nil {
		if bkg {
			go s.send(state)
		} else {
			s.send(state)
		}
	}
}

func (s *StateCom) setState(state *State) {
	state.Project = s.Project
	state.Branch = s.Branch
	s.Send(state, true)
	s.Global = state
}

func (s *StateCom) updateState(state string) {
	s.setState(&State{
		StateLabel: state,
		TimeStamp:  time.Now(),
	})
}

// Start sets the state of the build to started
func (s *StateCom) Start() {
	s.updateState("started")
}

// Env sets the state of the build to
// creating the environment
func (s *StateCom) Env() {
	s.updateState("creating env")
}

// Cloning sets the state of the build to
// cloning repo
func (s *StateCom) Cloning() {
	s.updateState("cloning repo")
}

// CleanUp sets the state of the build to
// clean up
func (s *StateCom) CleanUp() {
	s.updateState("clean up")
}

type successPayload struct {
	Proj   string `json:"project"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}

func (s *StateCom) success(commit string) {
	if s.client != nil {
		body := successPayload{
			Proj:   s.Project,
			Branch: s.Branch,
			Commit: commit,
		}
		s.Send(body, false)
	}
}

// Done sets the state of the build to
// done
func (s *StateCom) Done(success bool, commit string) {
	if success {
		s.success(commit)
	}
	s.updateState("done")
}

// SetServiceState sets the state of a specific service
func (s *StateCom) SetServiceState(srv *ServiceStateMgr) {
	srvState := &ServiceState{
		Name: srv.Name,
		State: &State{
			Project:    s.Project,
			Branch:     s.Project,
			StateLabel: srv.State,
			TimeStamp:  time.Now(),
		},
	}
	s.Send(srvState, true)
}
