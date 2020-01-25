package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"

	"github.com/google/go-github/v24/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Releaser struct {
	Repository *git.Repository
	Auth       ssh.AuthMethod
}

type ReleaseOption struct {
	Owner            string
	RepositoryName   string
	DevelopBranch    string
	ReleaseBranch    string
	ProductionBranch string
	AccessToken      string
}

func (re *Releaser) StartRelease(o ReleaseOption) error {
	r := re.Repository
	err := r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	if err != nil {
		if err.Error() != "already up-to-date" {
			return errors.Wrap(err, "Cannot fetch from origin")
		}
	}
	headRefRame := plumbing.NewRemoteReferenceName("origin", o.ProductionBranch)
	headRef, err := r.Reference(headRefRame, true)
	if err != nil {
		return errors.Wrapf(err, "Reference is not found. ref: %s", headRefRame)
	}
	nums, err := FindPullRequestNumbers(r, headRef)
	if err != nil {
		return errors.Wrap(err, "Failed to find pull request have beem merged since head")
	}
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/remotes/origin/" + o.DevelopBranch + ":refs/heads/" + o.ReleaseBranch),
		},
	})
	if err != nil {
		if err.Error() != "already up-to-date" {
			return errors.Wrap(err, "Cannot push branch")
		}
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: o.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	var pullRequests []github.PullRequest
	for _, n := range nums {
		pr, _, err := client.PullRequests.Get(context.TODO(), o.Owner, o.RepositoryName, n)
		if err != nil {
			return err
		}
		if *pr.State == "closed" && *pr.Merged == true {
			pullRequests = append(pullRequests, *pr)
		}
	}
	tmpl, err := template.New("release template").
		Parse("{{range .}}* {{range .Labels}}[{{.Name}}]{{end -}} {{.Title}} by @{{.User.Login}} #{{.Number}}\n{{end}}")
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, pullRequests)
	if err != nil {
		return err
	}
	body := buf.String()
	fmt.Println(body)
	title := "Release Test"
	releaseBranch := o.ReleaseBranch
	baseBranch := o.ProductionBranch
	pr, _, err := client.PullRequests.Create(context.TODO(), o.Owner, o.RepositoryName, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Head:  &releaseBranch,
		Base:  &baseBranch,
	})
	if err != nil {
		return err
	}
	log.Printf("PR is created. here: %s", pr.GetHTMLURL())
	return nil
}
