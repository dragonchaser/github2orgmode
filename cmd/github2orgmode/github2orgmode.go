package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
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
	// list all repositories for the authenticated user
	for _, r := range reposlist {
		data := strings.Split(r, "/")
		opts := github.IssueListByRepoOptions{
			State: "all",
		}
		issues, _, err := client.Issues.ListByRepo(ctx, data[0], data[1], &opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, i := range issues {
			fmt.Printf("[%s] %s - %s\n", *i.State, r, *i.Title)
		}
	}
}
