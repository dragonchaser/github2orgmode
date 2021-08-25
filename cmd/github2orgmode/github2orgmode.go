package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func main() {
	reposlist := strings.Fields(os.Getenv("GH_REPOS"))
	token := os.Getenv("GH_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	o := []string{
		fmt.Sprintf("#+TITLE: Github issues for %s", strings.Join(reposlist, ", ")),
		"#+CATEGORY: Github",
		"",
	}
	// list all repositories for the authenticated user for the selected repos
	for _, r := range reposlist {
		o = append(o, fmt.Sprintf("* %s", r))
		data := strings.Split(r, "/")
		opts := github.IssueListByRepoOptions{
			State: "open",
		}
		issues, _, err := client.Issues.ListByRepo(ctx, data[0], data[1], &opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, i := range issues {
			st := ""
			if *i.State == "open" {
				st = "TODO"
			} else {
				st = "DONE"
			}
			o = append(o, fmt.Sprintf("** %s - %s", st, *i.Title))
			// the timestamps have been created without < & > intentionally
			// I do not want them to show up in the daily agenda
			o = append(o, fmt.Sprintf("    Created: %s", *i.CreatedAt))
			o = append(o, fmt.Sprintf("    Updated: %s", *i.UpdatedAt))
			o = append(o, fmt.Sprintf("    [%s]", *i.URL))
			o = append(o, "\n")
		}
	}
	fmt.Printf(strings.Join(o, "\n") + "\n")
}
