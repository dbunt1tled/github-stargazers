package cli

import (
	"slices"

	"github.com/dbunt1tled/github-stargazers/internal/config"
	"github.com/dbunt1tled/github-stargazers/internal/db"
	"github.com/dbunt1tled/github-stargazers/internal/github"
	"github.com/spf13/cobra"
)

func NewUnStargazerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unstar",
		Short: "List not your's stargazers",
		Long:  "List your's stargazers who removed star from your repositories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				cfg                         *config.Config
				githubManager               *github.GitHubManager
				storage                     *db.Storage
				err                         error
				userStarred, userStargazers []string
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

			userStarred, err = githubManager.GetStarredUsers(cfg.GitHubUsername)
			if err != nil {
				return err
			}
			userStargazers, err = storage.GetStargazers(cmd.Context())
			if err != nil {
				return err
			}
			seen := make(map[string]bool)
			for _, u := range userStarred {
				if !slices.Contains(userStargazers, u) && !seen[u] {
					seen[u] = true
					cmd.Println(u)
				}
			}
			return nil
		},
	}
}
