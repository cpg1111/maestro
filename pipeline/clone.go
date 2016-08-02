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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/util"

	pb "gopkg.in/cheggaaa/pb.v1"
	git "gopkg.in/libgit2/git2go.v24"
)

// Project is a struct for the Project in the pipeline
type Project struct {
	conf      config.Project
	State     string
	ABSPath   string
	Services  map[string]*Service
	creds     *credentials.RawCredentials
	gitCreds  *git.Cred
	clonePath string
	CloneOpts *git.CloneOptions
}

func credCB(gitCreds *git.Cred) git.CredentialsCallback {
	return func(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
		return 0, gitCreds
	}
}

func certCheckCB(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

var (
	progbar     *pb.ProgressBar
	received    uint
	hasFinished = false
	done        = make(chan bool)
)

func handleProgress(stats git.TransferProgress) git.ErrorCode {
	if progbar == nil {
		progbar = pb.StartNew((int)(stats.TotalObjects))
	}
	newObjs := stats.ReceivedObjects - received
	for i := 0; i < (int)(newObjs); i++ {
		progbar.Increment()
	}
	received = stats.ReceivedObjects
	if (int)(received) == (int)(stats.TotalObjects) && !hasFinished {
		hasFinished = true
		time.Sleep(time.Second)
		done <- true
		return git.ErrOk
	}
	return git.ErrOk
}

// New returns a new instance of a pipeline project
func New(conf *config.Config, creds *credentials.RawCredentials, clonePath, branch, last, curr string) *Project {
	var absPath string
	cwd, cwdErr := os.Getwd()
	if clonePath[0] == '.' {
		absPath = fmt.Sprintf("%s/%s", cwd, clonePath)
	} else {
		absPath = clonePath
	}
	newServices := make(map[string]*Service)
	for i := range conf.Services {
		newServices[conf.Services[i].Name] = NewService(conf.Services[i], creds, clonePath, last, curr)
	}
	if cwdErr != nil {
		panic(cwdErr)
	}
	gitCreds := creds.ToGitCredentials()
	fetchOpts := &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credCB(&gitCreds),
			CertificateCheckCallback: certCheckCB,
			TransferProgressCallback: handleProgress,
		},
	}
	cloneOpts := &git.CloneOptions{
		FetchOptions: fetchOpts,
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafe,
		},
		Bare:           false,
		CheckoutBranch: branch,
	}
	return &Project{
		conf:      conf.Project,
		State:     "Pending",
		ABSPath:   absPath,
		Services:  newServices,
		creds:     creds,
		gitCreds:  &gitCreds,
		clonePath: clonePath,
		CloneOpts: cloneOpts,
	}
}

// Clone clones a git repo
func (p *Project) Clone() (resRepo *git.Repository, resErr error) {
	log.Println("Cloning Repo...")
	repoChan := make(chan *git.Repository)
	errChan := make(chan error)
	go func() {
		repo, repoErr := git.Clone(p.conf.RepoURL, fmt.Sprintf("%s/", p.clonePath), p.CloneOpts)
		repoChan <- repo
		errChan <- repoErr
	}()
	var doneMsg bool
	for {
		select {
		case resRepo = <-repoChan:
			if resRepo != nil && doneMsg {
				cdErr := os.Chdir(p.ABSPath)
				if cdErr != nil {
					log.Fatal(cdErr)
				}
				return
			}
		case resErr = <-errChan:
			if resErr != nil {
				panic(resErr)
			}
		case doneMsg = <-done:
			if resRepo != nil && doneMsg {
				cdErr := os.Chdir(p.ABSPath)
				if cdErr != nil {
					log.Fatal(cdErr)
				}
				return
			}
		}
	}
}

// Checkout checks out the repo to the current commit or HEAD of the branch
func (p *Project) Checkout(repo *git.Repository, commit string) error {
	tree, treeErr := util.CommitToTree(repo, commit)
	if treeErr != nil {
		return treeErr
	}
	opts := &git.CheckoutOpts{}
	return repo.CheckoutTree(tree, opts)
}
