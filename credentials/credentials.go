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

package credentials

import (
	"os/user"
	"strings"

	"github.com/cpg1111/maestro/config"

	prompt "github.com/segmentio/go-prompt"
	git "gopkg.in/libgit2/git2go.v22"
)

// RawCredentials is a struct for any credentials
type RawCredentials struct {
	project    config.Project
	SSHPrivKey string
	SSHPubKey  string
	Username   string
	Password   string
}

// ToGitCredentials converts RawCredentials to Git credentials
func (rc *RawCredentials) ToGitCredentials() git.Cred {
	var num int
	var creds git.Cred
	switch rc.project.AuthType {
	case "SSH":
		num, creds = git.NewCredSshKey(rc.Username, rc.SSHPubKey, rc.SSHPrivKey, rc.Password)
		if num != 0 {
			panic("GIT ERROR WHEN LOADING CREDENTIALS")
		}
		return creds
	case "HTTP":
		num, creds = git.NewCredUserpassPlaintext(rc.Username, rc.Password)
		if num != 0 {
			panic("GIT ERROR WHEN LOADING CREDENTIALS")
		}
		return creds
	}
	return creds
}

func readKey(path string) (string, error) {
	fullPath := path
	if strings.Contains(path, "~") {
		currUser, uErr := user.Current()
		if uErr != nil {
			return "", uErr
		}
		fullPath = strings.Replace(path, "~", currUser.HomeDir, 1)
	}
	return fullPath, nil
}

// NewCreds returns a pointer to a new instance of RawCredentials
func NewCreds(project *config.Project) (*RawCredentials, error) {
	var privKey, pubKey string
	var privErr, pubErr error
	if strings.Contains(project.AuthType, "SSH") {
		privKey, privErr = readKey(project.SSHPrivKeyPath)
		if privErr != nil {
			return nil, privErr
		}
		pubKey, pubErr = readKey(project.SSHPubKeyPath)
		if pubErr != nil {
			return nil, pubErr
		}
	}
	var pwd string
	if project.PromptForPWD {
		pwd = prompt.PasswordMasked("Please Enter Your Password")
	} else {
		pwd = project.Password
	}
	creds := &RawCredentials{
		project:    *project,
		SSHPrivKey: privKey,
		SSHPubKey:  pubKey,
		Username:   project.Username,
		Password:   pwd,
	}
	return creds, nil
}
