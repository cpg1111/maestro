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

package cleanUp

import (
	"io/ioutil"
	"os"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/util"
)

func handleCMDs(cmds []string) error {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		return pwdErr
	}
	if len(cmds) > 0 {
		for i := range cmds {
			cmd, cmdErr := util.FormatCommand(cmds[i], pwd)
			if cmdErr != nil {
				return cmdErr
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func saveArtifacts(artifacts []config.Artifact) error {
	artifactsLen := len(artifacts)
	if artifactsLen > 0 {
		finished := make(chan bool)
		err := make(chan error)
		for i := range artifacts {
			go func() {
				buildFile, readErr := ioutil.ReadFile(artifacts[i].RuntimeFilePath)
				if readErr != nil {
					err <- readErr
					finished <- false
				} else {
					writeErr := ioutil.WriteFile(artifacts[i].SaveFilePath, buildFile, 0644)
					if writeErr != nil {
						err <- writeErr
						finished <- false
					} else {
						finished <- true
					}
				}
			}()
		}
		total := 0
		for {
			select {
			case done := <-finished:
				if done {
					total++
				}
				if total == artifactsLen {
					return nil
				}
			case errMsg := <-err:
				if errMsg != nil {
					return errMsg
				}
			}
		}
	}
	return nil
}

// Run runs the clean tasks
func Run(conf *config.CleanUp, clonePath *string) error {
	cmdErr := handleCMDs(conf.AdditionalCMDs)
	if cmdErr != nil {
		return cmdErr
	}
	artifactErr := saveArtifacts(conf.Artifacts)
	if artifactErr != nil {
		return artifactErr
	}
	if conf.InDaemon {
		// TODO send status to daemon
		return nil
	}
	os.RemoveAll(*clonePath)
	return nil
}
