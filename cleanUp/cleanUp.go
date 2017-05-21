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
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if len(cmds) > 0 {
		for i := range cmds {
			cmd, err := util.FmtCommand(cmds[i], pwd)
			if err != nil {
				return err
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
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
		errChan := make(chan error)
		for i := range artifacts {
			go func() {
				buildFile, err := ioutil.ReadFile(artifacts[i].RuntimeFilePath)
				if err != nil {
					errChan <- err
				} else {
					errChan <- ioutil.WriteFile(artifacts[i].SaveFilePath, buildFile, 0644)
				}
			}()
		}
		total := 0
		for {
			err := <-errChan
			if err != nil {
				return err
			}
			total++
			if total == artifactsLen {
				break
			}
		}
	}
	return nil
}

// Run runs the clean tasks
func Run(conf *config.CleanUp, clonePath *string) error {
	err := handleCMDs(conf.AdditionalCMDs)
	if err != nil {
		return err
	}
	err = saveArtifacts(conf.Artifacts)
	if err != nil {
		return err
	}
	if conf.InDaemon {
		// TODO send status to daemon
		return nil
	}
	os.RemoveAll(*clonePath)
	return nil
}
