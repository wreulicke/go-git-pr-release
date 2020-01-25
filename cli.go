package main

import (
	prompt "github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func Run(args []string) error {
	c := ReleaseOption{}
	rootCmd := cobra.Command{
		Use: "go-git-pr-release",
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.AccessToken == "" {
				c.AccessToken = prompt.Input("access token: ", func(t prompt.Document) []prompt.Suggest {
					return []prompt.Suggest{}
				})
			}
			r, err := git.PlainOpen(".")
			if err != nil {
				return err
			}
			auth, err := ssh.NewSSHAgentAuth("git")
			if err != nil {
				return err
			}
			releaser := Releaser{
				Repository: r,
				Auth:       auth,
			}

			err = releaser.StartRelease(c)
			if err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.Flags().StringVarP(&c.AccessToken, "access-token", "t", "", "access-token")
	rootCmd.Flags().StringVarP(&c.Owner, "owner", "o", "", "owner")
	rootCmd.Flags().StringVarP(&c.RepositoryName, "repository", "r", "", "repository")
	rootCmd.Flags().StringVarP(&c.BaseBranch, "base", "b", "master", "base branch")
	rootCmd.Flags().StringVarP(&c.HeadBranch, "head", "d", "develop", "head branch")
	rootCmd.Flags().StringVar(&c.ReleaseBranch, "release", "release", "release branch")
	return rootCmd.Execute()
}
