package main

import (
	"log"
	"time"

	"github.com/dbunt1tled/github-stargazers/internal/config"
	"github.com/dbunt1tled/github-stargazers/internal/db"
	"github.com/dbunt1tled/github-stargazers/internal/github"
)

func main() {
	var (
		cfg                               *config.Config
		githubManager                     *github.GitHubManager
		storage                           *db.Storage
		err                               error
		repos, stargazers, added, removed []string
		today, prev                       string
	)

	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}
	githubManager = github.NewGitHubManager(cfg.GitHubToken, cfg.GitHubUsername)
	storage, err = db.New(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	repos, err = githubManager.GetRepositories(cfg.GitHubUsername)
	if err != nil {
		log.Fatal(err)
	}
	today = time.Now().Format("2006-01-02")
	for _, repo := range repos {
		log.Printf("Processing %s\n", repo)
		stargazers, err = githubManager.GetStargazers(cfg.GitHubUsername, repo)
		if err != nil {
			log.Fatal(err)
		}
		err = storage.Add(repo, today, stargazers)
		if err != nil {
			log.Fatal(err)
		}
		prev, err = storage.GetPreviousDate(repo, today)
		if err != nil {
			log.Fatal(err)
		}
		if prev == "" {
			prev = today
		}
		added, removed, err = storage.Diff(repo, today, prev)
		if err != nil {
			log.Fatal(err)
		}
		if len(added) > 0 {
			log.Printf(" ➕Added: %v\n", added)
		}
		if len(removed) > 0 {
			log.Printf(" ➖Removed: %v\n", removed)
		}
		if len(added) == 0 && len(removed) == 0 {
			log.Println("   No changes since", prev)
		}
	}
}
