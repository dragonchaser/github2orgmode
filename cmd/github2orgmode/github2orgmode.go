package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
	"os"
)

func main() {
	token := os.Getenv("GH_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(repos)
}
