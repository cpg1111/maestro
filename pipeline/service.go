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
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/util"

	git "gopkg.in/libgit2/git2go.v22"
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
	currCommitObject, parseErr := repo.RevparseSingle(currCommit)
	if parseErr != nil {
		return nil, parseErr
	}
	currCommitID := currCommitObject.Id()
	currCommitRef, lookupErr := repo.LookupCommit(currCommitID)
	if lookupErr != nil {
		return nil, lookupErr
	}
	currTree, treeErr := currCommitRef.Tree()
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
	prevCommitObject, parseErr := repo.RevparseSingle(*lastBuildCommit)
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
	return os.OpenFile(s.conf.BuildLogFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

func logToFile(in *bufio.Scanner, out *os.File, errChan chan error) {
	for in.Scan() {
		text := fmt.Sprintf("%s STDOUT: %s\n", time.Now(), in.Text())
		log.Println("writing ", text)
		fd, writeErr := out.WriteString(text)
		log.Println(fd)
		if writeErr != nil {
			errChan <- writeErr
		}
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
	logFile, fileErr := s.getLogFile()
	if fileErr != nil {
		log.Fatal("fileErr: ", fileErr)
		panic(fileErr)
	}
	defer logFile.Close()
	errChan := make(chan error)
	go logToFile(in1, logFile, errChan)
	go logToFile(in2, logFile, errChan)
	writeErr := <-errChan
	if writeErr != nil {
		return writeErr
	}
	syncErr := logFile.Sync()
	if syncErr != nil {
		panic(syncErr)
	}
	in1Err := in1.Err()
	if in1Err != nil {
		log.Println(in1Err)
		return in1Err
	}
	in2Err := in2.Err()
	if in2 != nil {
		log.Println(in2Err)
		return in2Err
	}
	return nil
}

func (s *Service) execSrvCmd(cmdStr, path string) (*exec.Cmd, error) {
	cmdStr, tmplErr := util.TemplateCommits(cmdStr, s.lastCommit, s.currCommit)
	if tmplErr != nil {
		log.Println(tmplErr)
		return nil, tmplErr
	}
	log.Println("executing", cmdStr)
	cmd, cmdErr := util.FmtCommand(cmdStr, path)
	if cmdErr != nil {
		log.Println(cmdErr)
		return cmd, cmdErr
	}
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
	if s.conf.CheckCMD == "" {
		return true, nil
	}
	cmdStr, tmplErr := util.TemplateCommits(s.conf.CheckCMD, s.lastCommit, s.currCommit)
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
