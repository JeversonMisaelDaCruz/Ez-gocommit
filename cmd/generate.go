package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jeversonmisael/ez-gocommit/internal/ai"
	"github.com/jeversonmisael/ez-gocommit/internal/config"
	gitcollector "github.com/jeversonmisael/ez-gocommit/internal/git"
	"github.com/jeversonmisael/ez-gocommit/internal/ui"
	"github.com/spf13/cobra"
)

func runGenerate(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadWithOverrides(flagStyle, flagModel, flagLanguage)
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	ctx, err := gitcollector.Collect(cwd, cfg.MaxDiffLines)
	if err != nil {
		return err
	}

	userPrompt := ai.BuildUserPrompt(ctx, cfg.CommitStyle)

	stopSpinner := startSpinner("Analyzing your changes with Claude...")
	suggestions, err := ai.GenerateSuggestions(userPrompt, cfg.APIKey, cfg.Model)
	stopSpinner()

	if err != nil {
		return err
	}

	fmt.Println()

	result, err := ui.Run(suggestions)
	if err != nil {
		return err
	}

	if result.Cancelled {
		color.Yellow("\nAborted.")
		return nil
	}

	if err := doCommit(result.Message, result.Body); err != nil {
		return err
	}

	color.Green("\n✔ Committed: %s\n", result.Message)
	return nil
}

func doCommit(message, body string) error {
	full := message
	if strings.TrimSpace(body) != "" {
		full = message + "\n\n" + body
	}

	gitCmd := exec.Command("git", "commit", "-m", full)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	return gitCmd.Run()
}

func startSpinner(label string) func() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	stop := make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case <-stop:
				fmt.Printf("\r%s\r", strings.Repeat(" ", len(label)+4))
				return
			default:
				fmt.Printf("\r%s %s", frames[i%len(frames)], label)
				time.Sleep(80 * time.Millisecond)
				i++
			}
		}
	}()

	return func() {
		close(stop)
		time.Sleep(100 * time.Millisecond)
	}
}
