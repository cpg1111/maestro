package pipeline

import (
	"os"
	"os/exec"

	"github.com/cpg1111/maestro/config"

	git "github.com/libgit2/git2go"
)

type Service struct {
	config.Service
	Diff  bool
	State string
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
