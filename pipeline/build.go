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

	git "gopkg.in/libgit2/git2go.v22"
)

func build(srv *DepService, index string, done chan string, errChan chan error, shouldDeploy *bool) {
	err := srv.build.execBuild()
	if err != nil {
		errChan <- err
		return
	}
	log.Println("Run tests")
	err = RunTests(srv.build)
	if err != nil {
		errChan <- err
		return
	}
	log.Println("Tests done")
	if !*shouldDeploy {
		log.Println("No Deploy")
		done <- index
		return
	}
	log.Println("Checking deployment...")
	err = check(srv.build)
	if err != nil {
		errChan <- err
		return
	}
	srv.build.shouldBuild = false
	done <- index
	return
}

func runServiceBuild(srvs map[string]*DepService, testAll, shouldDeploy *bool) error {
	doneChan := make(chan string)
	errChan := make(chan error)
	buildTotal := 0
	for i := range srvs {
		if srvs[i].build.shouldBuild || *testAll {
			buildTotal++
			log.Println("Building", srvs[i].build.conf.Name)
			go build(srvs[i], i, doneChan, errChan, shouldDeploy)
		} else {
			runErr := runServiceBuild(srvs[i].Children, testAll, shouldDeploy)
			if runErr != nil {
				return runErr
			}
		}
	}
	total := 0
	if buildTotal > 0 {
		for {
			select {
			case index := <-doneChan:
				total++
				if len(srvs[index].Children) > 0 {
					runErr := runServiceBuild(srvs[index].Children, testAll, shouldDeploy)
					if runErr != nil {
						return runErr
					}
				}
				if total == buildTotal {
					return nil
				}
			case errMsg := <-errChan:
				if errMsg != nil {
					log.Println(errMsg)
					return errMsg
				}
			}
		}
	}
	return nil
}

// Run runs the build for all changed services
func Run(depTree *DepTree, repo *git.Repository, lastBuildCommit, currBuildCommit *string, testAll, shouldDeploy *bool) error {
	currNode := depTree.CurrNode
	travErr := TraverseTree(currNode, repo, lastBuildCommit, currBuildCommit)
	if travErr != nil {
		return travErr
	}
	rootMap := make(map[string]*DepService)
	rootMap["root"] = currNode
	err := runServiceBuild(rootMap, testAll, shouldDeploy)
	if err != nil {
		return err
	}
	return nil
}
