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

package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strings"
	"text/template"
)

// Commits is a struct for templating commits
type Commits struct {
	Prev string
	Curr string
}

// TemplateCommits templates the commits hashes into commans
func TemplateCommits(strCMD, lastCommit, currCommit string) (string, error) {
	buff := &bytes.Buffer{}
	commits := &Commits{Prev: lastCommit, Curr: currCommit}
	tmpl, err := template.New("cmd").Parse(strCMD)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(buff, commits)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

// FmtCommand formats a string to an exec'able command
func FmtCommand(strCMD, path string) (*exec.Cmd, error) {
	var cmdArr []string
	if len(strCMD) >= 8 && strings.Contains(strCMD[0:8], "bash -c") {
		cmdArr = []string{strCMD[0:4], strCMD[5:7], strCMD[8:]}
	} else {
		cmdArr = strings.Split(strCMD, " ")
	}
	cmdPath, err := exec.LookPath(cmdArr[0])
	if err != nil {
		return nil, err
	}
	cmdLen := len(cmdArr)
	for i := 1; i < cmdLen; i++ {
		if strings.Contains(cmdArr[i], "./") {
			if path[len(path)-1] != '/' {
				path = fmt.Sprintf("%s/", path)
			}
			cmdArr[i] = strings.Replace(cmdArr[i], "./", path, -1)
		}
		if strings.Contains(cmdArr[i], "~") {
			currUser, err := user.Current()
			if err != nil {
				return &exec.Cmd{}, err
			}
			cmdArr[i] = strings.Replace(cmdArr[i], "~", currUser.HomeDir, -1)
		}
	}
	cmd := exec.Command(cmdPath, cmdArr[1:]...)
	return cmd, nil
}
