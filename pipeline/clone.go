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

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"

	pb "gopkg.in/cheggaaa/pb.v1"
	git "gopkg.in/libgit2/git2go.v22"
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
	progbar  *pb.ProgressBar
	received uint
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
	return git.ErrOk
}

// New returns a new instance of a pipeline project
func New(conf *config.Config, creds *credentials.RawCredentials, clonePath, branch string) *Project {
	newServices := make(map[string]*Service)
	for i := range conf.Services {
		newServices[conf.Services[i].Name] = NewService(conf.Services[i], creds)
	}
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		panic(cwdErr)
	}
	gitCreds := creds.ToGitCredentials()
	cloneOpts := &git.CloneOptions{
		RemoteCallbacks: &git.RemoteCallbacks{
			CredentialsCallback:      credCB(&gitCreds),
			CertificateCheckCallback: certCheckCB,
			TransferProgressCallback: handleProgress,
		},
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafeCreate,
		},
		Bare:           false,
		CheckoutBranch: branch,
	}
	//conf.CloneOpts.RemoteCreateCallback = createRemote
	return &Project{
		conf:      conf.Project,
		State:     "Pending",
		ABSPath:   fmt.Sprintf("%s/%s", cwd, clonePath),
		Services:  newServices,
		creds:     creds,
		gitCreds:  &gitCreds,
		clonePath: clonePath,
		CloneOpts: cloneOpts,
	}
}

// Clone clones a git repo
func (p *Project) Clone(opts *git.CloneOptions) (*git.Repository, error) {
	log.Println("Cloning Repo...")
	return git.Clone(p.conf.RepoURL, fmt.Sprintf("%s/", p.clonePath), opts)
}
