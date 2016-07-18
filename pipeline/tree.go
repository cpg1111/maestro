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
	"errors"
	"log"

	git "gopkg.in/libgit2/git2go.v22"
)

// DepTree is a dependency tree to determine whether or not to build services
type DepTree struct {
	CurrNode *DepService
}

// TraverseTree traverses a dependency tree
func TraverseTree(depSrv *DepService, repo *git.Repository, lastBuildCommit, currBuildCommit *string) error {
	if depSrv == nil {
		return errors.New("Service is nil in tree")
	}
	shouldBuild, buildErr := depSrv.build.ShouldBuild(repo, lastBuildCommit, currBuildCommit)
	if buildErr != nil {
		return buildErr
	}
	for i := range depSrv.Children {
		depSrv.Children[i].build.shouldBuild = shouldBuild
		travErr := TraverseTree(depSrv.Children[i], repo, lastBuildCommit, currBuildCommit)
		if travErr != nil {
			return travErr
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

func dependsOnChild(child, parent *DepService) int {
	for j := range child.build.conf.DependsOn {
		if parent.Children[child.build.conf.DependsOn[j]] != nil {
			return j
		}
	}
	return -1
}

func getDependencies(depSrv *DepService, tree *DepTree, created map[string]*DepService, proj *Project) {
	for j := range depSrv.build.conf.DependsOn {
		if created[depSrv.build.conf.DependsOn[j]] != nil {
			youngerParentIndex := dependsOnChild(depSrv, created[depSrv.build.conf.DependsOn[j]])
			if youngerParentIndex > -1 {
				depSrv.Parent = created[depSrv.build.conf.DependsOn[j]].Children[depSrv.build.conf.DependsOn[youngerParentIndex]]
				depSrv.Children[depSrv.build.conf.Name] = depSrv
			} else {
				depSrv.Parent = created[depSrv.build.conf.DependsOn[j]]
				if depSrv.Parent.Children[depSrv.build.conf.Name] == nil {
					depSrv.Parent.Children[depSrv.build.conf.Name] = depSrv
				}
			}
		} else {
			parent := &DepService{
				build:    proj.Services[depSrv.build.conf.DependsOn[j]],
				Children: make(map[string]*DepService),
			}
			depSrv.Parent = parent
			created[depSrv.build.conf.DependsOn[j]] = parent
			parent.Children[depSrv.build.conf.Name] = depSrv
			created[depSrv.build.conf.Name] = depSrv
			if len(parent.build.conf.DependsOn) > 0 {
				getDependencies(parent, tree, created, proj)
			} else if tree.CurrNode == nil {
				tree.CurrNode = parent
			} else if tree.CurrNode.Children[parent.build.conf.Name] == nil {
				tree.CurrNode.Children[parent.build.conf.Name] = parent
			}
		}
	}
}

// NewTreeList returns a list of dependency trees
func NewTreeList(proj *Project) (newTree *DepTree) {
	created := make(map[string]*DepService)
	log.Println("Creating a dependency tree")
	newTree = &DepTree{CurrNode: nil}
	for i := range proj.Services {
		log.Println("Finding a spot in the tree for ", proj.Services[i].conf.Name)
		if created[proj.Services[i].conf.Name] != nil {
			continue
		}
		depSrv := &DepService{build: proj.Services[i], Children: make(map[string]*DepService)}
		if len(depSrv.build.conf.DependsOn) != 0 {
			getDependencies(depSrv, newTree, created, proj)
		} else if newTree.CurrNode == nil {
			newTree.CurrNode = depSrv
		} else {
			newTree.CurrNode.Children[depSrv.build.conf.Name] = depSrv
		}
	}
	return
}
