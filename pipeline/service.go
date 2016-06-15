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

package pipeline

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/util"

	git "gopkg.in/libgit2/git2go.v22"
)

// Service is a struct for services in the pipeline
type Service struct {
	conf          config.Service
	Diff          bool
	State         string
	creds         *credentials.RawCredentials
	index         *git.Index
	shouldBuild   bool
	logFileOffset int64
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

func (s *Service) getLogFile() (*os.File, error) {
	stat := &syscall.Stat_t{}
	statErr := syscall.Stat(s.conf.BuildLogFilePath, stat)
	falseErrStr := "no such file or directory"
	if statErr != nil && statErr.Error() != falseErrStr {
		return nil, statErr
	}
	if stat.Size > 0 {
		return os.Open(s.conf.BuildLogFilePath)
	}
	return os.Create(s.conf.BuildLogFilePath)
}

func (s *Service) execSrvCmd(cmdStr, path string) (*exec.Cmd, error) {
	cmd, cmdErr := util.FormatCommand(cmdStr, path)
	if cmdErr != nil {
		return cmd, cmdErr
	}
	if s.conf.BuildLogFilePath != "" {
		logFile, openErr := s.getLogFile()
		if openErr != nil {
			return cmd, openErr
		}
		log.Println(logFile)
		output, outputErr := cmd.Output()
		offset, writeErr := logFile.WriteAt(output, s.logFileOffset)
		if outputErr != nil {
			return cmd, outputErr
		}
		if writeErr != nil {
			return cmd, writeErr
		}
		s.logFileOffset = (int64)(offset)
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		runErr := cmd.Run()
		if runErr != nil {
			return cmd, runErr
		}
	}
	return cmd, nil
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
	_, err := s.execSrvCmd(s.conf.BuildCMD, s.conf.Path)
	log.Println("Built")
	return err
}

func (s *Service) execTests() error {
	log.Println("Testing")
	_, err := s.execSrvCmd(s.conf.TestCMD, s.conf.Path)
	log.Println("Tested")
	return err
}

func (s *Service) execCreate() error {
	if s.conf.CreateCMD == "" {
		return nil
	}
	cmd, err := s.execSrvCmd(s.conf.CreateCMD, s.conf.Path)
	if err != nil {
		return err
	}
	if s.conf.HealthCheck.Type == "PTRACE_ATTACH" {
		passPid := HealthCheck(&s.conf).(func(pid int) error)
		return passPid(cmd.Process.Pid).(error)
	}
	return HealthCheck(&s.conf).(error)
}

func (s *Service) execUpdate() error {
	if s.conf.UpdateCMD == "" {
		return nil
	}
	cmd, err := s.execSrvCmd(s.conf.UpdateCMD, s.conf.Path)
	if err != nil {
		return err
	}
	if s.conf.HealthCheck.Type == "PTRACE_ATTACH" {
		passPid := HealthCheck(&s.conf).(func(pid int) error)
		passed := passPid(cmd.Process.Pid)
		if passed != nil {
			return passed.(error)
		}
	}
	checkRes := HealthCheck(&s.conf)
	if checkRes != nil {
		return checkRes.(error)
	}
	return nil
}
