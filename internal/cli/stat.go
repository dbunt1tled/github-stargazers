package cli

import (
	"time"

	"github.com/dbunt1tled/github-stargazers/internal/config"
	"github.com/dbunt1tled/github-stargazers/internal/db"
	"github.com/dbunt1tled/github-stargazers/internal/github"
	"github.com/spf13/cobra"
)

func NewStatCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stargazers",
		Short: "Statistics about your stargazers",
		Long:  "Get Information who add/remove star to your repositories",
		RunE: func(cmd *cobra.Command, _ []string) error {
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
				return err
			}
			githubManager = github.NewGitHubManager(cfg.GitHubToken, cfg.GitHubUsername)
			storage, err = db.New(cmd.Context(), cfg.DatabasePath)
			if err != nil {
				return err
			}

			repos, err = githubManager.GetRepositories(cfg.GitHubUsername)
			if err != nil {
				return err
			}
			today = time.Now().Format("2006-01-02")
			for _, repo := range repos {
				cmd.Printf("Processing %s\n", repo)
				stargazers, err = githubManager.GetStargazers(cfg.GitHubUsername, repo)
				if err != nil {
					return err
				}
				err = storage.Add(cmd.Context(), repo, today, stargazers)
				if err != nil {
					return err
				}
				prev, err = storage.GetPreviousDate(cmd.Context(), repo, today)
				if err != nil {
					return err
				}
				if prev == "" {
					prev = today
				}
				added, removed, err = storage.Diff(cmd.Context(), repo, today, prev)
				if err != nil {
					return err
				}
				if len(added) > 0 {
					cmd.Printf(" ➕Added: %v\n", added)
				}
				if len(removed) > 0 {
					cmd.Printf(" ➖Removed: %v\n", removed)
				}
				if len(added) == 0 && len(removed) == 0 {
					cmd.Println("  No changes since", prev)
				}
			}
			return nil
		},
	}
}