package main

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

// FindMergeBase provides like git merge-base
func FindMergeBase(r *git.Repository, c1 *object.Commit, c2 *object.Commit) (*object.Commit, error) {
	iter1 := object.NewCommitIterCTime(c1, nil, nil)
	iter2 := object.NewCommitIterCTime(c2, nil, nil)
	history := map[plumbing.Hash]bool{}

	for {
		var hasErrorOne bool
		if c, err := iter1.Next(); err == nil {
			if _, ok := history[c.Hash]; ok {
				return c, nil
			}
			history[c.Hash] = true
		} else {
			hasErrorOne = true
		}
		if c, err := iter2.Next(); err == nil {
			if _, ok := history[c.Hash]; ok {
				return c, nil
			}
			history[c.Hash] = true
		} else {
			if hasErrorOne {
				return nil, errors.New("merge-base is not found")
			}
		}
	}
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
		base, err := FindMergeBase(r, productionHead, c)
		if err != nil {
			if err.Error() == "object not found" {
				// ignore
				continue
			}
			return nil, errors.Wrapf(err, "Cannot find merge base between %s and %s", productionHead.Hash.String(), ref.Name())
		}
		if base.Hash != c.Hash {
			results = append(results, ref)
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
