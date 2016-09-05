package statecom

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cpg1111/maestro/config"
)

type State struct {
	Project    string
	Branch     string
	StateLabel string
	TimeStamp  time.Time
}

type ServiceState struct {
	Name  string
	State *State
}

type StateCom struct {
	Project  string
	Branch   string
	Services map[string]*ServiceStateMgr
	Global   *State
	client   *http.Client
}

func New(conf config.Config, maestrodEndpoint, branch string) *StateCom {
	var client *http.Client
	if len(maestrodEndpoint) > 0 {
		client = &http.Client{}
	}
	stateCom := &StateCom{
		Project: conf.Project.RepoURL,
		Branch:  branch,
		Global: &State{
			Project:    conf.Project.RepoURL,
			Branch:     branch,
			StateLabel: "pending",
			TimeStamp:  time.Now(),
		},
		client: client,
	}
	stateCom.Services = make(map[string]*ServiceStateMgr)
	for i := 0; i < len(conf.Services); i++ {
		stateCom.Services[conf.Services[i].Name] = NewServiceStateMgr(conf.Services[i], stateCom)
	}
	return stateCom
}

func (s *StateCom) Send(state interface{}) {
	if s.client != nil {
		go func() {
			payload, marshErr := json.Marshal(state)
			if marshErr != nil {
				log.Println("WARNING", marshErr)
			}
			payloadRdr := bytes.NewReader(payload)
			resp, postErr := s.client.Post("/state", "application/json", payloadRdr)
			if postErr != nil {
				log.Println("WARNING", postErr.Error())
			}
			if resp.StatusCode != 201 {
				log.Println("WARNING STATEUPDATE NOT SENT")
			}
		}()
	}
}

func (s *StateCom) setState(state *State) {
	state.Project = s.Project
	state.Branch = s.Branch
	s.Send(state)
	s.Global = state
}

func (s *StateCom) Start() {
	startState := &State{
		StateLabel: "started",
		TimeStamp:  time.Now(),
	}
	s.setState(startState)
}

func (s *StateCom) Env() {
	envState := &State{
		StateLabel: "creating env",
		TimeStamp:  time.Now(),
	}
	s.setState(envState)
}

func (s *StateCom) Cloning() {
	cloneState := &State{
		StateLabel: "cloning repo",
		TimeStamp:  time.Now(),
	}
	s.setState(cloneState)
}

func (s *StateCom) CleanUp() {
	cleanUpState := &State{
		StateLabel: "clean up",
		TimeStamp:  time.Now(),
	}
	s.setState(cleanUpState)
}

func (s *StateCom) Done() {
	doneState := &State{
		StateLabel: "done",
		TimeStamp:  time.Now(),
	}
	s.setState(doneState)
}

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
	s.Send(srvState)
}
