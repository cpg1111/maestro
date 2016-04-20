package pipeline

import (
	//"os/exec"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
)

// Service is a struct for services in the pipeline
type Service struct {
	conf  config.Service
	Diff  bool
	State string
	creds *credentials.RawCredentials
	index *git.Index
}

// NewService returns an instance of a pipeline service
func NewService(srv config.Service, creds *credentials.RawCredentials) *Service {
	return &Service{
		conf:  srv,
		Diff:  false,
		State: "Pending",
		creds: creds,
		index: nil,
	}
}

// ShouldRunPipeLine diffs a service's path and determs whether or not it needs to run the pipeline on it
func (s *Service) ShouldRunPipeLine() (bool, error) {
	return true, nil
}
