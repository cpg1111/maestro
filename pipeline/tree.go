package pipeline

// DepTree is a dependency tree to determine whether or not to build services
type DepTree struct {
	CurrNode *DepService
}

// TraverseTree traverses a dependency tree
func TraverseTree(depSrv *DepService) error {
	shouldBuild, buildErr := depSrv.build.ShouldBuild()
	if buildErr != nil {
		return buildErr
	}
	if shouldBuild {
		for i := range depSrv.Children {
			depSrv.Children[i].build.shouldBuild = true
			travErr := TraverseTree(depSrv.Children[i])
			if travErr != nil {
				return travErr
			}
		}
	}
	return nil
}

// DepService represents a service in the tree
type DepService struct {
	build    *Service
	Parent   *DepService
	Children map[string]*DepService
}

func getDependencies(depSrv *DepService, created map[string]*DepService, proj *Project) {
	for j := range depSrv.build.conf.DependsOn {
		if created[depSrv.build.conf.DependsOn[j]] != nil {
			depSrv.Parent = created[depSrv.build.conf.DependsOn[j]]
			if created[depSrv.build.conf.DependsOn[j]].Children[depSrv.build.conf.Name] == nil {
				created[depSrv.build.conf.DependsOn[j]].Children[depSrv.build.conf.Name] = depSrv
			}
		} else {
			parent := &DepService{build: proj.Services[depSrv.build.conf.DependsOn[j]]}
			depSrv.Parent = parent
			if parent.Children[depSrv.build.conf.Name] == nil {
				parent.Children[depSrv.build.conf.Name] = depSrv
			}
		}
	}
}

// NewTreeList returns a list of dependency trees
func NewTreeList(proj *Project) []*DepTree {
	var newTrees []*DepTree
	created := make(map[string]*DepService)
	for i := range proj.Services {
		if created[proj.Services[i].conf.Name] != nil {
			depSrv := created[proj.Services[i].conf.Name]
			getDependencies(depSrv, created, proj)
		}
		depSrv := &DepService{build: proj.Services[i]}
		if len(depSrv.build.conf.DependsOn) == 0 {
			newTrees = append(newTrees, &DepTree{CurrNode: depSrv})
			created[depSrv.build.conf.Name] = depSrv
		} else {
			getDependencies(depSrv, created, proj)
		}
	}
	return newTrees
}
