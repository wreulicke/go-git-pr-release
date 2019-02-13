package main

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type PullRequestReference struct {
	plumbing.Reference
}

func FindPullRequestReference(r *git.Repository) ([]*plumbing.Reference, error) {
	re, err := r.Remote("origin")
	if err != nil {
		return nil, err
	}
	rfs, err := re.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []*plumbing.Reference
	for _, rf := range rfs {
		if strings.HasPrefix(string(rf.Name()), "refs/pull/") && strings.HasSuffix(string(rf.Name()), "/head") {
			result = append(result, rf)
		}
	}

	return result, nil
}
