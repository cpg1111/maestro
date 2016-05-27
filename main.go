package main

import (
	"flag"
	"log"
	"os"

	"github.com/cpg1111/maestro/cleanUp"
	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/environment"
	"github.com/cpg1111/maestro/pipeline"
)

var (
	confPath        = flag.String("config", "./conf.toml", "Path to the config for maestro to use")
	clonePath       = flag.String("clone-path", "./", "Local path to clone repo to defaults to PWD")
	checkoutBranch  = flag.String("branch", "master", "Git branch to checkout for project")
	lastBuildCommit = flag.String("prev-commit", "", "Previous commit to compare to")
	deploy          = flag.Bool("deploy", false, "Whether or not to deploy this build")
)

func main() {
	flag.Parse()
	if *lastBuildCommit == "" {
		log.Println("Maestro requires a previous commit to build from.")
		os.Exit(1)
	}
	log.Println("Loading Configuration...")
	conf, confErr := config.Load(*confPath, *clonePath)
	if confErr != nil {
		log.Fatal(confErr)
	}
	log.Println("Loading Credentials...")
	creds, credErr := credentials.NewCreds(&conf.Project)
	if credErr != nil {
		log.Fatal(credErr)
	}
	if len(conf.Environment.Exec) > 0 || len(conf.Environment.ExecSync) > 0 {
		log.Println("Loading Environment...")
		envErr := environment.Load(&conf.Environment)
		if envErr != nil {
			log.Fatal(envErr)
		}
	}
	log.Println("Creating Pipeline...")
	pipe := pipeline.New(&conf, creds, *clonePath, *checkoutBranch)
	repo, cloneErr := pipe.Clone(pipe.CloneOpts)
	if cloneErr != nil {
		log.Fatal(cloneErr)
	}
	log.Println("Building Dependency Tree...")
	depTrees := pipeline.NewTreeList(pipe)
	log.Println("Building Serivces...", *deploy)
	buildErr := pipeline.Run(depTrees, repo, lastBuildCommit, deploy)
	if buildErr != nil {
		os.RemoveAll(*clonePath)
		log.Fatal(buildErr)
	}
	log.Println("Cleaning Up Build...")
	cleanUp.Run(&conf.CleanUp, clonePath)
}
