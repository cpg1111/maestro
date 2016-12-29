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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cpg1111/maestro/cleanUp"
	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/environment"
	"github.com/cpg1111/maestro/pipeline"
	"github.com/cpg1111/maestro/statecom"
	"github.com/cpg1111/maestro/util"
)

var (
	confPath        = flag.String("config", "/etc/maestro/conf.toml", "Path to the config for maestro to use")
	clonePath       = flag.String("clone-path", "/tmp/clone", "Local path to clone repo to defaults to PWD")
	checkoutBranch  = flag.String("branch", "master", "Git branch to checkout for project")
	lastBuildCommit = flag.String("prev-commit", "", "Previous commit to compare to")
	currBuildCommit = flag.String("curr-commit", "", "Current commit to compare to, if not specified, will diff HEAD of branch")
	testAll         = flag.Bool("test-all", false, "Whether to test all services or not")
	deploy          = flag.Bool("deploy", false, "Whether or not to deploy this build")
)

func setEnv(conf *config.Config) {
	os.Setenv("LAST_COMMIT", *lastBuildCommit)
	if *currBuildCommit != "" {
		os.Setenv("CURR_COMMIT", *currBuildCommit)
	} else {
		os.Setenv("CURR_COMMIT", "HEAD")
	}
	if len(conf.Environment.Exec) > 0 || len(conf.Environment.ExecSync) > 0 {
		log.Println("Loading Environment...")
		envErr := environment.Load(&conf.Environment)
		if envErr != nil {
			log.Fatal(envErr)
		}
	}
}

func checkNoPrevCommit() {
	if *lastBuildCommit == "" {
		log.Println("Maestro requires a previous commit to build from.")
		os.Exit(2)
	}
}

func getMaestrodEndpoint() string {
	endpoint := fmt.Sprintf("%s:%s", os.Getenv("MAESTROD_SERVICE_HOST"), os.Getenv("MAESTROD_SERVICE_PORT"))
	if len(endpoint) > 1 {
		return endpoint
	}
	return ""
}

func main() {
	flag.Parse()
	checkNoPrevCommit()
	clonePath := util.FmtClonePath(clonePath)
	log.Println("Loading Configuration...")
	conf, confErr := config.Load(*confPath, *clonePath)
	if confErr != nil {
		log.Fatal(confErr)
	}
	stateCom := statecom.New(conf, getMaestrodEndpoint(), *checkoutBranch)
	log.Println("Loading Credentials...")
	stateCom.Start()
	creds, credErr := credentials.NewCreds(&conf.Project)
	if credErr != nil {
		log.Fatal(credErr)
	}
	stateCom.Env()
	setEnv(&conf)
	log.Println("Creating Pipeline...")
	pipe := pipeline.New(&conf, creds, *clonePath, *checkoutBranch, *lastBuildCommit, *currBuildCommit)
	stateCom.Cloning()
	repo, cloneErr := pipe.Clone()
	if cloneErr != nil {
		log.Fatal(cloneErr)
	}
	if *currBuildCommit != "" {
		log.Println("Checking out current commit...")
		pipe.Checkout(repo, *currBuildCommit)
	}
	log.Println("Building Dependency Tree...")
	depTree := pipeline.NewTreeList(pipe)
	log.Println("Building Serivces...")
	buildErr := pipeline.Run(depTree, repo, stateCom, lastBuildCommit, currBuildCommit, testAll, deploy)
	if buildErr != nil {
		os.RemoveAll(*clonePath)
		log.Fatal(buildErr)
	}
	log.Println("Cleaning Up Build...")
	stateCom.CleanUp()
	cleanUp.Run(&conf.CleanUp, clonePath)
	stateCom.Done(true, *currBuildCommit)
}
