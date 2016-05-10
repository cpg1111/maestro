package pipeline

import "log"

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
func RunBuild(depTrees []*DepTree) error {
	log.Println("run")
	for i := range depTrees {
		TraverseTree(depTrees[i].CurrNode)
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