package main

import (
	"fmt"
	"testing"

	"gopkg.in/src-d/go-git.v4"
)

func TestList(t *testing.T) {
	r, _ := git.PlainOpen(".")
	re, _ := r.Remote("origin")
	rfs, _ := re.List(&git.ListOptions{})
	for _, rf := range rfs {
		fmt.Println(rf.Name())
	}
	t.Error("aaaaaaaaaa")
}
