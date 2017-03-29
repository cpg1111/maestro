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
	currTree, err := util.CommitToTree(repo, currCommit)
	if err != nil {
		return nil, err
	}
	return repo.DiffTreeToTree(prev, currTree, opts)
}

// ShouldBuild diffs a service's path and determs whether or not it needs to run the pipeline on it
func (s *Service) ShouldBuild(repo *git.Repository, lastBuildCommit, currBuildCommit *string) (bool, error) {
	if s.shouldBuild {
		return s.shouldBuild, nil
	}
	prevTree, err := util.CommitToTree(repo, *lastBuildCommit)
	if err != nil {
		return false, err
	}
	diffOpts, err := git.DefaultDiffOptions()
	if err != nil {
		return false, err
	}
	diffOpts.Pathspec = []string{s.diffPath}
	var diff *git.Diff
	if *currBuildCommit == "" {
		diff, err = diffToWorkingDir(repo, prevTree, &diffOpts)
	} else {
		diff, err = diffToMostRecentCommit(repo, prevTree, &diffOpts, *currBuildCommit)
	}
	if err != nil {
		return false, diffErr
	}
	deltas, err := diff.NumDeltas()
	if err != nil {
		return false, err
	}
	if deltas == 0 {
		return false, nil
	}
	log.Println("Found", deltas, "deltas in", s.conf.Name, "beginning build")
	s.shouldBuild = true
	return true, nil
}

func (s *Service) getLogFile() (*os.File, error) {
	return os.OpenFile(s.conf.BuildLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
}

func (s *Service) logToFile(stream string, in *bufio.Scanner) {
	out, err := s.getLogFile()
	if err != nil {
		panic(err)
	}
	defer out.Close()
	for in.Scan() {
		text := fmt.Sprintf("%s %s: %s\n", time.Now(), stream, in.Text())
		_, err = out.WriteString(text)
		if err != nil {
			panic(err)
		}
	}
	err = out.Sync()
	if err != nil {
		panic(err)
	}
	err = in.Err()
	if err != nil && err.Error() != "read |0: bad file descriptor" {
		panic(err)
	}
}

func (s *Service) logStdoutToFile(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
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
	cmdStr, err := util.TemplateCommits(cmdStr, s.lastCommit, s.currCommit)
	if err != nil {
		log.Println(tmplErr)
		return nil, tmplErr
	}
	cmd, err := util.FmtCommand(cmdStr, path)
	if err != nil {
		log.Println(err)
		return cmd, err
	}
	log.Printf("Running %v\n", cmd.Args)
	if s.conf.BuildLogFilePath != "" {
		err = s.logStdoutToFile(cmd)
		if err != nil {
			return cmd, err
		}
		err = cmd.Run()
		if err != nil {
			log.Println(err)
			return cmd, err
		}
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return cmd, err
		}
	}
	return cmd, nil
}

func (s *Service) execCheck() (bool, error) {
	for i := range s.conf.CheckCMD {
		if s.conf.CheckCMD[i] == "" {
			continue
		}
		cmdStr, err := util.TemplateCommits(s.conf.CheckCMD[i], s.lastCommit, s.currCommit)
		if err != nil {
			return false, err
		}
		cmd, err := util.FmtCommand(cmdStr, s.conf.Path)
		if err != nil {
			return false, err
		}
		fmt.Printf("%d\n", len(cmd.Args))
		err = cmd.Run()
		if err != nil {
			return false, err
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
		if s.conf.HealthCheck.Type == PTAttach {
			passPid := HealthCheck(&s.conf).(func(pid int) error)
			return passPid(cmd.Process.Pid).(error)
		}
	}
	checkRes := HealthCheck(&s.conf)
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
		if s.conf.HealthCheck.Type == PTAttach {
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
