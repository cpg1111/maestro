package pipeline

import (
	"log"

	git "gopkg.in/libgit2/git2go.v22"
)

func runServiceBuild(srvs map[string]*DepService) error {
	log.Println("building services")
	for i := range srvs {
		log.Println("building ", srvs[i].build.conf.Name, srvs[i].build.shouldBuild)
		if srvs[i].build.shouldBuild {
			err := srvs[i].build.execBuild()
			if err != nil {
				return err
			}
			if len(srvs[i].Children) > 0 {
				runServiceBuild(srvs[i].Children)
			}
		}
	}
	return nil
}

// RunBuild runs the build for all changed services
func RunBuild(depTrees []*DepTree, repo *git.Repository, lastBuildCommit string) error {
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
