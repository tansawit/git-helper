package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	repos, repoMap, _ = githubGetUserRepos("975bb0d7acde4a86dce7c87907007f150cf9117a")

	app  = kingpin.New("git-helper", "A command line tool for interacting with git repo information").Version("1.0.0")
	list = app.Command("list", "List all of your repositories")

	info         = app.Command("info", "Get information about repository <reponame>")
	infoRepoName = info.Arg("reponame", "Name of repository to use").Required().String()

	open         = app.Command("open", "Open repository in default browser")
	openRepoName = open.Arg("reponame", "Name of repository to use").Required().String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// List all repositories
	case list.FullCommand():
		repoList()
	// Get information about a repository
	case info.FullCommand():
		repoInfo()
	// Open a repository page in browser
	case open.FullCommand():
	}
}
