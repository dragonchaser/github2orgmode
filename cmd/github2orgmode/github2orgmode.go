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
	handlelist := strings.Fields(os.Getenv("GH_HANDLES"))
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
		project := fmt.Sprintf("* %s", r)
		todo := []string{}
		prog := []string{}
		done := []string{}
		data := strings.Split(r, "/")
		opts := github.IssueListByRepoOptions{
			State: "all",
			ListOptions: github.ListOptions{
				PerPage: 100000,
			},
		}
		issues, _, err := client.Issues.ListByRepo(ctx, data[0], data[1], &opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, i := range issues {
			issue := []string{}
			// add org- and repo-name to tags
			labels := fmt.Sprintf(":%s:%s:", reg.ReplaceAllString(data[0], ""), reg.ReplaceAllLiteralString(data[1], ""))
			for _, l := range i.Labels {
				labels = labels + reg.ReplaceAllString(l.GetName(), "") + ":"
			}
			status := "TODO"
			for _, a := range i.Assignees {
				for _, handle := range handlelist {
					if a.GetLogin() == handle {
						status = "IN PROGRESS"
					}
				}
			}
			if i.GetState() == "closed" {
				status = "DONE"
			}
			issue = append(issue, fmt.Sprintf("** %s %s \t\t %s", status, i.GetTitle(), labels))
			issue = append(issue, fmt.Sprintf("\tState     : %s", i.GetState()))
			issue = append(issue, fmt.Sprintf("\tCreator   : %s", i.GetUser().GetLogin()))
			// the timestamps have been created without < & > intentionally
			// I do not want them to show up in the daily agenda
			issue = append(issue, fmt.Sprintf("\tCreated   : %s", i.GetCreatedAt()))
			issue = append(issue, fmt.Sprintf("\tUpdated   : %s", i.GetUpdatedAt()))
			for _, a := range i.Assignees {
				issue = append(issue, fmt.Sprintf("\tAssignee : %s", a.GetLogin()))
			}
			if i.GetState() == "closed" {
				issue = append(issue, fmt.Sprintf("\tClosed at : %s", i.GetClosedAt()))
			}
			issue = append(issue, fmt.Sprintf("\t[%s]", i.GetURL()))
			issue = append(issue, "\n")
			issue = append(issue, formatBody(i.GetBody()))
			issue = append(issue, "\n")
			if status == "TODO" {
				todo = append(todo, strings.Join(issue, "\n"))
			} else if status == "DONE" {
				done = append(done, strings.Join(issue, "\n"))
			} else {
				prog = append(prog, strings.Join(issue, "\n"))
			}
		}
		o = append(o, project, strings.Join(prog, "\n"), strings.Join(todo, "\n"), strings.Join(done, "\n"))
	}
	fmt.Printf(strings.Join(o, "\n") + "\n")
}

func formatBody(body string) string {
	body = strings.ReplaceAll(body, "\r\n", "\n\t")
	return "\t" + body
}
