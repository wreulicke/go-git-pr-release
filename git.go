package main

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
)

func FindMergedPullRequest() {
	var r *git.Repository
	re, _ := r.Remote("origin")
	rfs, _ := re.List(&git.ListOptions{})
	for _, rf := range rfs {
		fmt.Println(rf.Name())
	}
	return
}
