package main

import (
	"flag"
	"log"
	"os"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/pipeline"
)

var (
	confPath        = flag.String("config", "./conf.toml", "Path to the config for maestro to use")
	clonePath       = flag.String("clone-path", "./", "Local path to clone repo to defaults to PWD")
	checkoutBranch  = flag.String("branch", "master", "Git branch to checkout for project")
	lastBuildCommit = flag.String("prev-commit", "", "Previous commit to compare to")
)

func main() {
	flag.Parse()
	if *lastBuildCommit == "" {
		log.Println("Maestro requires a previous commit to build from.")
		os.Exit(1)
	}
	log.Println("Running")
	conf, confErr := config.Load(*confPath, *clonePath)
	if confErr != nil {
		log.Fatal(confErr)
	}
	creds, credErr := credentials.NewCreds(conf.Project)
	if credErr != nil {
		log.Fatal(credErr)
	}
	pipe := pipeline.New(&conf, creds, *clonePath, *checkoutBranch)
	repo, cloneErr := pipe.Clone(pipe.CloneOpts)
	if cloneErr != nil {
		log.Fatal(cloneErr)
	}
	depTrees := pipeline.NewTreeList(pipe)
	buildErr := pipeline.RunBuild(depTrees, repo, *lastBuildCommit)
	if buildErr != nil {
		os.RemoveAll(*clonePath)
		log.Fatal(buildErr)
	}
	os.RemoveAll(*clonePath)
	log.Println(*repo)
	log.Println(*depTrees[0].CurrNode)
}
