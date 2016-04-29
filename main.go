package main

import (
	"flag"
	"log"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"
	"github.com/cpg1111/maestro/pipeline"
)

var confPath = flag.String("config", "./conf.toml", "Path to the config for maestro to use")
var clonePath = flag.String("clone-path", "./", "Local path to clone repo to defaults to PWD")
var checkoutBranch = flag.String("branch", "master", "Git branch to checkout for project")

func main() {
	flag.Parse()
	log.Println("Running")
	conf, confErr := config.Load(*confPath)
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
	log.Println(repo)
	log.Println(depTrees)
}
