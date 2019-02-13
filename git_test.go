package main

import (
	"testing"

	"gopkg.in/src-d/go-git.v4"
)

func TestFindPullRequestReference(t *testing.T) {
	// like ls-remote
	r, err := git.PlainOpen(".")
	if err != nil {
		t.Error(err)
		return
	}
	refs, err := FindPullRequestReference(r)
	if err != nil {
		t.Error(err)
		return
	}

	var contains bool
	for _, ref := range refs {
		if ref.Name() == "refs/pull/1/head" {
			contains = true
		}
	}
	if !contains {
		t.Error("Pull request reference is not found")
	}
}
