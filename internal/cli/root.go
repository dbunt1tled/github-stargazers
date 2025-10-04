package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd *cobra.Command
}

func NewCLI() *CLI {
	cli := &CLI{
		rootCmd: &cobra.Command{
			Use:   "github-stargazer",
			Short: "GitHub Stargazer",
			Long:  "Get Statistics about your stargazers",
		},
	}

	cli.rootCmd.AddCommand(NewStatCommand())
	cli.rootCmd.AddCommand(NewUnStargazerCommand())

	return cli
}

func (c *CLI) Execute() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Println(err) //nolint:forbidigo // print error
		os.Exit(1)
	}
}
