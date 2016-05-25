package pipeline

import (
	"log"

	git "gopkg.in/libgit2/git2go.v22"
)

func build(srv *DepService, index string, done chan string, errChan chan error) {
	err := srv.build.execBuild()
	if err != nil {
		errChan <- err
		return
	}
	err = RunTests(srv.build)
	if err != nil {
		errChan <- err
		return
	}
	err = check(srv.build)
	if err != nil {
		errChan <- err
		return
	}
	srv.build.shouldBuild = false
	done <- index
	return
}

func runServiceBuild(srvs map[string]*DepService) error {
	log.Println("building services")
	doneChan := make(chan string)
	errChan := make(chan error)
	for i := range srvs {
		log.Println("building ", srvs[i].build.conf.Name, srvs[i].build.shouldBuild)
		if srvs[i].build.shouldBuild {
			go build(srvs[i], i, doneChan, errChan)
		}
	}
	total := 0
	for {
		select {
		case index := <-doneChan:
			total++
			if len(srvs[index].Children) > 0 {
				runServiceBuild(srvs[index].Children)
			}
			if total == len(srvs) {
				return nil
			}
		case errMsg := <-errChan:
			if errMsg != nil {
				return errMsg
			}
		}
	}
}

// Run runs the build for all changed services
func Run(depTrees []*DepTree, repo *git.Repository, lastBuildCommit string) error {
	log.Println("run")
	for i := range depTrees {
		travErr := TraverseTree(depTrees[i].CurrNode, repo, lastBuildCommit)
		if travErr != nil {
			return travErr
		}
		log.Println(i+1, "tree")
		rootMap := make(map[string]*DepService)
		rootMap["root"] = depTrees[i].CurrNode
		err := runServiceBuild(rootMap)
		if err != nil {
			return err
		}
	}
	return nil
}
