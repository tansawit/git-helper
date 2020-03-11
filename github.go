package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v29/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
)

// getGitHubClient returns a github API client
func getGitHubClient(ghToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

// GetUserRepos returns a slice of the user's public GitHub Repositories
// Implementation Credit: https://github.com/lox/alfred-github-jump/repos.go
func githubGetUserRepos(ghToken string) ([]*github.Repository, map[string]github.Repository, error) {
	client := getGitHubClient(ghToken)
	ctx := context.Background()
	var repoMap = map[string]github.Repository{}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 45},
		Sort:        "pushed",
	}

	repos := []*github.Repository{}

	for {
		result, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return repos, repoMap, err
		}
		repos = append(repos, result...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	for _, repo := range repos {
		repoMap[*repo.Name] = *repo
	}

	return repos, repoMap, nil
}

// repoList prints a list of all user repositories
func repoList() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Name\tDescription\tLanguage\tURL\t# of Open Issues")
	for _, repo := range repos {
		fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t%v\t%v", *repo.Name, nilableString(repo.Description, 50), nilableString(repo.Language, 0), *repo.HTMLURL, *repo.OpenIssuesCount))
	}
	fmt.Fprintln(w)
	w.Flush()
}

// repoInfo prints info about a specific repository
func repoInfo() {
	if !validRepoName(*infoRepoName) {
		fmt.Println(`Please enter a valid reponame, without the username in front (i.e. git-helper instead of tansawit/git-helper)`)
	} else {
		repo := repoMap[*infoRepoName]
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 2, 8, 2, '\t', 0)
		fmt.Fprintln(w, "Name\tDescription\tLanguage\tURL\t# of Open Issues")
		fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t%v\t%v", *repo.Name, nilableString(repo.Description, 50), nilableString(repo.Language, 0), *repo.HTMLURL, *repo.OpenIssuesCount))
		fmt.Fprintln(w)
		w.Flush()
	}
}

// repoOpen openspen repo in browser
func repoOpen() {
	if !validRepoName(*infoRepoName) {
		fmt.Println(`Please enter a valid reponame, without the username in front (i.e. git-helper instead of tansawit/git-helper)`)
	} else {
		repo := repoMap[*infoRepoName]
		openbrowser(*repo.HTMLURL)
	}
}

// validRepoName checks if <name> is a valid user repository
func validRepoName(name string) bool {
	if _, ok := repoMap[name]; ok {
		return true
	}
	return false
}
