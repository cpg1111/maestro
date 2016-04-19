package pipeline

import (
	"os"
	"os/exec"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"

	git "gopkg.in/libgit2/git2go.v22"
)

type Service struct {
	config.Service
	Diff     bool
	State    string
	creds    credentials.RawCredentials
	gitCreds git.Cred
}

func NewService(srv config.Service, creds credentials.RawCredentials) *Service {
	return &Service{
		srv,
		Diff:     false,
		State:    "Pending",
		creds:    creds,
		gitCreds: creds.ToGitCredentials(),
	}
}

func (s *Service) ShouldRun() (bool, error) {

}

func (s *Service) Build() error {

}

type Project struct {
	config.Project
	State    string
	ABSPath  string
	Services map[string]Service
}
