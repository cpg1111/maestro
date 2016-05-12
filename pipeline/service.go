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
func (s *Service) ShouldBuild(repo *git.Repository, lastBuildCommit string) (bool, error) {
	log.Println("diff")
	prevCommitObject, _, parseErr := repo.RevparseExt(lastBuildCommit)
	if parseErr != nil {
		return false, parseErr
	}
	prevCommitID := prevCommitObject.Id()
	prevCommit, lookupErr := repo.LookupCommit(prevCommitID)
	if lookupErr != nil {
		return false, lookupErr
	}
	prevTree, treeErr := prevCommit.Tree()
	if treeErr != nil {
		return false, treeErr
	}
	diffOpts := &git.DiffOptions{
		Flags:            git.DiffNormal,
		IgnoreSubmodules: git.SubmoduleIgnoreNone,
		Pathspec:         []string{s.conf.Path},
	}
	diff, diffErr := repo.DiffTreeToWorkdir(prevTree, diffOpts)
	if diffErr != nil {
		return false, diffErr
	}
	deltas, deltaErr := diff.NumDeltas()
	if deltaErr != nil {
		return false, deltaErr
	}
	if deltas > 0 {
		return false, nil
	}
	s.shouldBuild = true
	return true, nil
}

func formatCommand(strCMD, path, name string) (*exec.Cmd, error) {
	cmdStr := strings.Split(strCMD, " ")
	log.Println("executing build for ", name)
	cmdPath, lookupErr := exec.LookPath(cmdStr[0])
	if lookupErr != nil {
		return &exec.Cmd{}, lookupErr
	}
	cmd := exec.Command(cmdPath)
	cmdLen := len(cmdStr)
	for i := 1; i < cmdLen; i++ {
		if strings.Contains(cmdStr[i], ".") {
			cmdStr[i] = strings.Replace(cmdStr[i], ".", path, 1)
		}
		if strings.Contains(cmdStr[i], "~") {
			currUser, userErr := user.Current()
			if userErr != nil {
				return &exec.Cmd{}, userErr
			}
			cmdStr[i] = strings.Replace(cmdStr[i], "~", currUser.HomeDir, 1)
		}
		cmd.Args = append(cmd.Args, cmdStr[i])
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func execSrvCmd(cmdStr, path, name string) error {
	cmd, cmdErr := formatCommand(cmdStr, path, name)
	if cmdErr != nil {
		return cmdErr
	}
	cmd.Run()
	return nil
}

func (s *Service) execCheck() (bool, error) {
	cmd, cmdErr := formatCommand(s.conf.CheckCMD, s.conf.Path, s.conf.Name)
	if cmdErr != nil {
		return false, cmdErr
	}
	cmd.Run()
	return false, nil
}

func (s *Service) execBuild() error {
	return execSrvCmd(s.conf.BuildCMD, s.conf.Path, s.conf.Name)
}

func (s *Service) execTests() error {
	return execSrvCmd(s.conf.TestCMD, s.conf.Path, s.conf.Name)
}

func (s *Service) execCreate() error {
	return execSrvCmd(s.conf.CreateCMD, s.conf.Path, s.conf.Name)
}

func (s *Service) execUpdate() error {
	return execSrvCmd(s.conf.UpdateCMD, s.conf.Path, s.conf.Name)
}
