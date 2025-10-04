package github

import (
	"context"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

type GitHubManager struct {
	client   *github.Client
	username string
}

func NewGitHubManager(token string, username string) *GitHubManager {
	var ctx = context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubManager{
		client:   client,
		username: username,
	}
}

func (gm *GitHubManager) GetRepositories(owner string) ([]string, error) {
	var allRepositories []string
	ctx := context.Background()
	opts := &github.RepositoryListByUserOptions{
		Type:      "owner",
		Sort:      "created",
		Direction: "asc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	for {
		repositories, resp, err := gm.client.Repositories.ListByUser(ctx, owner, opts)
		if err != nil {
			return nil, err
		}
		for _, repository := range repositories {
			allRepositories = append(allRepositories, repository.GetName())
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allRepositories, nil
}

func (gm *GitHubManager) GetStargazers(owner string, repo string) ([]string, error) {
	var allStargazers []string
	ctx := context.Background()
	opt := &github.ListOptions{
		PerPage: 100,
	}
	for {
		stars, resp, err := gm.client.Activity.ListStargazers(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}
		for _, star := range stars {
			if star.User != nil && star.User.Login != nil {
				allStargazers = append(allStargazers, *star.User.Login)
			}

		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allStargazers, nil

}
