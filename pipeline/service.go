package pipeline

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"

	git "gopkg.in/libgit2/git2go.v22"
)

// Service is a struct for services in the pipeline
type Service struct {
	conf        config.Service
	Diff        bool
	State       string
	creds       *credentials.RawCredentials
	index       *git.Index
	shouldBuild bool
}

// NewService returns an instance of a pipeline service
func NewService(srv config.Service, creds *credentials.RawCredentials) *Service {
	return &Service{
		conf:        srv,
		Diff:        false,
		State:       "Pending",
		creds:       creds,
		index:       nil,
		shouldBuild: false,
	}
}

// ShouldBuild diffs a service's path and determs whether or not it needs to run the pipeline on it
func (s *Service) ShouldBuild() (bool, error) {
	s.shouldBuild = true
	return true, nil
}

func (s *Service) execBuild() error {
	cmdStr := strings.Split(s.conf.BuildCMD, " ")
	log.Println("executing build for ", s.conf.Name)
	cmdPath, lookupErr := exec.LookPath(cmdStr[0])
	if lookupErr != nil {
		return lookupErr
	}
	cmd := exec.Command(cmdPath)
	cmdLen := len(cmdStr)
	for i := 1; i < cmdLen; i++ {
		if strings.Contains(cmdStr[i], ".") {
			cmdStr[i] = strings.Replace(cmdStr[i], ".", s.conf.Path, 1)
		}
		if strings.Contains(cmdStr[i], "~") {
			currUser, userErr := user.Current()
			if userErr != nil {
				return userErr
			}
			cmdStr[i] = strings.Replace(cmdStr[i], "~", currUser.HomeDir, 1)
		}
		cmd.Args = append(cmd.Args, cmdStr[i])
	}
	log.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return nil
}
