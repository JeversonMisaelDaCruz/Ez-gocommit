package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagStyle  string
	flagModel  string
	flagConfig string
)

var rootCmd = &cobra.Command{
	Use:   "ezgocommit",
	Short: "AI-powered Git commit message generator",
	Long: `Ez-gocommit analyzes your staged changes and generates
semantic commit messages using Claude AI.

Run it inside a Git repository after staging your changes:
  git add .
  ezgocommit`,
	RunE: runGenerate,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagStyle, "style", "", "commit style: conventional, gitmoji, free, custom")
	rootCmd.PersistentFlags().StringVar(&flagModel, "model", "", "Claude model to use (default: claude-sonnet-4-6)")
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "path to config file")

	rootCmd.AddCommand(versionCmd)
}
