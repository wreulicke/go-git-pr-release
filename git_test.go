package main

import (
	"testing"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func TestFindPullRequestReference(t *testing.T) {
	// like git ls-remote origin refs/pull/*/head
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

func TestFilterMergedPullRequest(t *testing.T) {
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

	hash, err := r.ResolveRevision(plumbing.Revision("3704d81253329f6abff59a3ed6a542030ff1cabc"))
	if err != nil {
		t.Error(err)
		return
	}
	base := plumbing.NewHashReference("base", *hash)
	results, err := FilterMergedPullRequest(r, base, refs)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(results)
	var contains bool
	for _, v := range results {
		if v.Name() == "refs/pull/1/head" {
			contains = true
		} else if v.Name() == "refs/pull/2/head" {
			t.Error("refs/pull/2/head should not be contained")
		}
	}
	if !contains {
		t.Logf("Results: %v", results)
		t.Error("refs/pull/1/head is not found")
	}
}
