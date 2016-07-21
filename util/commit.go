package util

import (
	git "gopkg.in/libgit2/git2go.v22"
)

func CommitToTree(repo *git.Repository, hash string) (*git.Tree, error) {
	commitObj, commitErr := repo.RevparseSingle(hash)
	if commitErr != nil {
		return nil, commitErr
	}
	commitID := commitObj.Id()
	commit, lookupErr := repo.LookupCommit(commitID)
	if lookupErr != nil {
		return nil, lookupErr
	}
	return commit.Tree()
}
