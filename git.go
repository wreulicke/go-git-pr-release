package main

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

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

// FilterMergedPullRequest returns merged reference from refs since productionRef.
func FilterMergedPullRequest(r *git.Repository, productionRef *plumbing.Reference, refs []*plumbing.Reference) (results []*plumbing.Reference, err error) {
	productionHead, err := r.CommitObject(productionRef.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "Cannot find commit object")
	}

	for _, ref := range refs {
		c, err := r.CommitObject(ref.Hash())
		if err != nil {
			if err.Error() == "object not found" {
				// ignore
				continue
			}
			return nil, errors.Wrapf(err, "Cannot find commit object for %s", ref.Name())
		}
		base, err := productionHead.MergeBase(c)
		if err != nil {
			if err.Error() == "object not found" {
				// ignore
				continue
			}
			return nil, errors.Wrapf(err, "Cannot find merge base between %s and %s", productionHead.Hash.String(), ref.Name())
		}
		for _, b := range base {
			if b.Hash != c.Hash {
				results = append(results, ref)
				break
			}
		}
	}
	return results, nil
}

func FindPullRequestNumbers(r *git.Repository, productionRef *plumbing.Reference) (nums []int, err error) {
	refs, err := FindPullRequestReference(r)
	if err != nil {
		return nil, errors.Wrap(err, "Cannnot find pull request reference")
	}
	refs, err = FilterMergedPullRequest(r, productionRef, refs)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to filter merged pull request")
	}
	for _, ref := range refs {
		name := string(ref.Name())
		if strings.HasPrefix(name, "refs/pull/") && strings.HasSuffix(name, "/head") {
			name = strings.TrimPrefix(name, "refs/pull/")
			name = strings.TrimSuffix(name, "/head")
			if num, err := strconv.Atoi(name); err == nil {
				nums = append(nums, num)
			}
		}
	}
	return nums, nil
}
