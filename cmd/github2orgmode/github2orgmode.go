package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
	"regexp"
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
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

			// add org- and repo-name to tags
			labels := fmt.Sprintf(":%s:%s:", reg.ReplaceAllString(data[0], ""), reg.ReplaceAllLiteralString(data[1], ""))
			for _, l := range i.Labels {
				labels = labels + reg.ReplaceAllString(l.GetName(), "") + ":"
			}
			o = append(o, fmt.Sprintf("** TODO %s \t\t %s", i.GetTitle(), labels))
			// the timestamps have been created without < & > intentionally
			// I do not want them to show up in the daily agenda
			o = append(o, fmt.Sprintf("\tCreated  : %s", i.GetCreatedAt()))
			o = append(o, fmt.Sprintf("\tUpdated  : %s", i.GetUpdatedAt()))
			o = append(o, fmt.Sprintf("\tCreator  : %s", i.GetUser().GetLogin()))
			for _, a := range i.Assignees {
				o = append(o, fmt.Sprintf("\tAssignee : %s", a.GetLogin()))
			}
			o = append(o, fmt.Sprintf("\t[%s]", i.GetURL()))
			o = append(o, "\n")
			o = append(o, formatBody(i.GetBody()))
			o = append(o, "\n")
		}
	}
	fmt.Printf(strings.Join(o, "\n") + "\n")
}

func formatBody(body string) string {
	body = strings.ReplaceAll(body, "\r\n", "\n\t")
	return "\t" + body
}
