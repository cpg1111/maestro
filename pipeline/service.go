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
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/util"

	git "gopkg.in/libgit2/git2go.v24"
)

// Service is a struct for services in the pipeline
type Service struct {
	conf          config.Service
	State         string
	creds         *credentials.RawCredentials
	index         *git.Index
	shouldBuild   bool
	logFileOffset int64
	lastCommit    string
	currCommit    string
	path          string
	diffPath      string
	HasFailed     bool
}

// NewService returns an instance of a pipeline service
func NewService(srv config.Service, creds *credentials.RawCredentials, clonePath, last, curr string) *Service {
	diffPath := util.FmtDiffPath(clonePath, srv.Path)
	return &Service{
		conf:        srv,
		State:       "Pending",
		creds:       creds,
		index:       nil,
		shouldBuild: false,
		lastCommit:  last,
		currCommit:  curr,
		path:        srv.Path,
		diffPath:    diffPath,
	}
}

func diffToWorkingDir(repo *git.Repository, prev *git.Tree, opts *git.DiffOptions) (*git.Diff, error) {
	return repo.DiffTreeToWorkdir(prev, opts)
}

func diffToMostRecentCommit(repo *git.Repository, prev *git.Tree, opts *git.DiffOptions, currCommit string) (*git.Diff, error) {
	currTree, treeErr := util.CommitToTree(repo, currCommit)
	if treeErr != nil {
		return nil, treeErr
	}
	return repo.DiffTreeToTree(prev, currTree, opts)
}

// ShouldBuild diffs a service's path and determs whether or not it needs to run the pipeline on it
func (s *Service) ShouldBuild(repo *git.Repository, lastBuildCommit, currBuildCommit *string) (bool, error) {
	if s.shouldBuild {
		return s.shouldBuild, nil
	}
	prevTree, treeErr := util.CommitToTree(repo, *lastBuildCommit)
	if treeErr != nil {
		return false, treeErr
	}
	diffOpts, optsErr := git.DefaultDiffOptions()
	if optsErr != nil {
		return false, optsErr
	}
	diffOpts.Pathspec = []string{s.diffPath}
	var diff *git.Diff
	var diffErr error
	if *currBuildCommit == "" {
		diff, diffErr = diffToWorkingDir(repo, prevTree, &diffOpts)
	} else {
		diff, diffErr = diffToMostRecentCommit(repo, prevTree, &diffOpts, *currBuildCommit)
	}
	if diffErr != nil {
		return false, diffErr
	}
	deltas, deltaErr := diff.NumDeltas()
	if deltaErr != nil {
		return false, deltaErr
	}
	if deltas == 0 {
		return false, nil
	}
	log.Println("Found", deltas, "deltas in", s.conf.Name, "beginning build")
	s.shouldBuild = true
	return true, nil
}

func (s *Service) getLogFile() (*os.File, error) {
	return os.OpenFile(s.conf.BuildLogFilePath, os.O_CREATE|os.O_WRONLY, 0664)
}

func (s *Service) logToFile(stream string, in *bufio.Scanner) {
	out, fileErr := s.getLogFile()
	if fileErr != nil {
		panic(fileErr)
	}
	defer out.Close()
	for in.Scan() {
		text := fmt.Sprintf("%s %s: %s\n", time.Now(), stream, in.Text())
		_, writeErr := out.WriteString(text)
		if writeErr != nil {
			panic(writeErr)
		}
	}
	syncErr := out.Sync()
	if syncErr != nil {
		panic(syncErr)
	}
	inErr := in.Err()
	if inErr != nil && inErr.Error() != "read |0: bad file descriptor" {
		panic(inErr)
	}
}

func (s *Service) logStdoutToFile(cmd *exec.Cmd) error {
	stdout, outErr := cmd.StdoutPipe()
	if outErr != nil {
		return outErr
	}
	stderr, errErr := cmd.StderrPipe()
	if errErr != nil {
		return errErr
	}
	in1 := bufio.NewScanner(stdout)
	in2 := bufio.NewScanner(stderr)
	go s.logToFile("STDOUT", in1)
	go s.logToFile("STDERR", in2)
	return nil
}

func (s *Service) execSrvCmd(cmdStr, path string) (*exec.Cmd, error) {
	if cmdStr == "" {
		return nil, errors.New("empty string is not a valid command")
	}
	cmdStr, tmplErr := util.TemplateCommits(cmdStr, s.lastCommit, s.currCommit)
	if tmplErr != nil {
		log.Println(tmplErr)
		return nil, tmplErr
	}
	cmd, cmdErr := util.FmtCommand(cmdStr, path)
	if cmdErr != nil {
		log.Println(cmdErr)
		return cmd, cmdErr
	}
	log.Printf("Running %v\n", cmd.Args)
	if s.conf.BuildLogFilePath != "" {
		stdoutErr := s.logStdoutToFile(cmd)
		if stdoutErr != nil {
			return cmd, stdoutErr
		}
		runErr := cmd.Run()
		if runErr != nil {
			log.Println(runErr)
			return cmd, runErr
		}
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		runErr := cmd.Run()
		if runErr != nil {
			log.Println(runErr)
			return cmd, runErr
		}
	}
	return cmd, nil
}

func (s *Service) execCheck() (bool, error) {
	for i := range s.conf.CheckCMD {
		if s.conf.CheckCMD[i] == "" {
			continue
		}
		cmdStr, tmplErr := util.TemplateCommits(s.conf.CheckCMD[i], s.lastCommit, s.currCommit)
		if tmplErr != nil {
			return false, tmplErr
		}
		cmd, cmdErr := util.FmtCommand(cmdStr, s.conf.Path)
		if cmdErr != nil {
			return false, cmdErr
		}
		checkErr := cmd.Run()
		if checkErr != nil {
			return false, checkErr
		}
	}
	return true, nil
}

func (s *Service) execBuild() error {
	for i := range s.conf.BuildCMD {
		_, err := s.execSrvCmd(s.conf.BuildCMD[i], s.conf.Path)
		if err != nil {
			s.HasFailed = true
			return err
		}
	}
	log.Println("Built")
	return nil
}

func (s *Service) execTests() error {
	log.Println("Testing")
	for i := range s.conf.TestCMD {
		_, err := s.execSrvCmd(s.conf.TestCMD[i], s.conf.Path)
		if err != nil {
			s.HasFailed = true
			return err
		}
	}
	log.Println("Tested")
	return nil
}

func (s *Service) execCreate() error {
	for i := range s.conf.CreateCMD {
		cmd, err := s.execSrvCmd(s.conf.CreateCMD[i], s.conf.Path)
		if err != nil {
			s.HasFailed = true
			return err
		}
		if s.conf.HealthCheck.Type == "PTRACE_ATTACH" {
			passPid := HealthCheck(&s.conf).(func(pid int) error)
			return passPid(cmd.Process.Pid).(error)
		}
	}
	checkRes := HealthCheck(&s.conf).(error)
	if checkRes != nil {
		s.HasFailed = true
		return checkRes.(error)
	}
	return nil
}

func (s *Service) execUpdate() error {
	for i := range s.conf.UpdateCMD {
		cmd, err := s.execSrvCmd(s.conf.UpdateCMD[i], s.conf.Path)
		if err != nil {
			s.HasFailed = true
			return err
		}
		if s.conf.HealthCheck.Type == "PTRACE_ATTACH" {
			passPid := HealthCheck(&s.conf).(func(pid int) error)
			passed := passPid(cmd.Process.Pid)
			if passed != nil {
				s.HasFailed = true
				return passed.(error)
			}
		}
	}
	checkRes := HealthCheck(&s.conf)
	if checkRes != nil {
		s.HasFailed = true
		return checkRes.(error)
	}
	return nil
}
