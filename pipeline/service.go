package pipeline

import (
	"log"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/util"

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
func (s *Service) ShouldBuild(repo *git.Repository, lastBuildCommit *string) (bool, error) {
	log.Println("diff")
	prevCommitObject, _, parseErr := repo.RevparseExt(*lastBuildCommit)
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

func execSrvCmd(cmdStr, path string) error {
	cmd, cmdErr := util.FormatCommand(cmdStr, path)
	if cmdErr != nil {
		return cmdErr
	}
	cmd.Run()
	return nil
}

func (s *Service) execCheck() (bool, error) {
	if s.conf.CheckCMD == "" {
		return true, nil
	}
	cmd, cmdErr := util.FormatCommand(s.conf.CheckCMD, s.conf.Path)
	if cmdErr != nil {
		return false, cmdErr
	}
	checkErr := cmd.Run()
	if checkErr != nil {
		return false, checkErr
	}
	return true, nil
}

func (s *Service) execBuild() error {
	return execSrvCmd(s.conf.BuildCMD, s.conf.Path)
}

func (s *Service) execTests() error {
	return execSrvCmd(s.conf.TestCMD, s.conf.Path)
}

func (s *Service) execCreate() error {
	if s.conf.CreateCMD == "" {
		return nil
	}
	return execSrvCmd(s.conf.CreateCMD, s.conf.Path)
}

func (s *Service) execUpdate() error {
	if s.conf.UpdateCMD == "" {
		return nil
	}
	return execSrvCmd(s.conf.UpdateCMD, s.conf.Path)
}
