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
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	cfg, err := config.LoadWithOverrides(flagStyle, flagModel, flagLanguage)
	if err != nil {
		return err
	}

	if !flagDryRun {
		if err := cfg.Validate(); err != nil {
			return err
		}
	}

	ctx, err := gitcollector.Collect(cwd, cfg.MaxDiffLines)
	if err != nil {
		return err
	}

	var suggestions []ai.Suggestion

	if flagDryRun {
		suggestions = mockSuggestions(ctx, cfg.CommitStyle)
		color.Yellow("\n[dry-run] skipping API call — using mock suggestions\n")
	} else {
		userPrompt := ai.BuildUserPrompt(ctx, cfg.CommitStyle)
		stopSpinner := startSpinner("Analyzing your changes with Claude...")
		suggestions, err = ai.GenerateSuggestions(userPrompt, cfg.APIKey, cfg.Model)
		stopSpinner()
		if err != nil {
			return err
		}
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

	if flagDryRun {
		color.Yellow("\n[dry-run] Would commit: %q\n", result.Message)
		if strings.TrimSpace(result.Body) != "" {
			color.Yellow("[dry-run] With body:\n%s\n", result.Body)
		}
		return nil
	}

	if err := doCommit(result.Message, result.Body); err != nil {
		return err
	}

	color.Green("\n✔ Committed: %s\n", result.Message)
	return nil
}

func mockSuggestions(ctx *gitcollector.Context, style string) []ai.Suggestion {
	scope := inferScope(ctx.ChangedFiles)
	verb, prefix := styleVerbs(style)

	return []ai.Suggestion{
		{
			Rank:       1,
			Confidence: "high",
			Message:    fmt.Sprintf("%s%s(%s): %s staged changes", prefix, verb[0], scope, describeChanges(ctx.ChangedFiles)),
			Body:       "",
			Reasoning:  fmt.Sprintf("Based on %d staged file(s) on branch %q", len(ctx.ChangedFiles), ctx.BranchName),
		},
		{
			Rank:       2,
			Confidence: "medium",
			Message:    fmt.Sprintf("%s%s(%s): update %s implementation", prefix, verb[1], scope, scope),
			Body:       "",
			Reasoning:  "Alternative framing focused on the updated component",
		},
		{
			Rank:       3,
			Confidence: "low",
			Message:    fmt.Sprintf("%s%s: apply changes to %s", prefix, verb[2], scope),
			Body:       "",
			Reasoning:  "Conservative option without scope qualifier",
		},
	}
}

func inferScope(files []string) string {
	if len(files) == 0 {
		return "core"
	}
	parts := strings.Split(files[0], "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	name := parts[0]
	name = strings.TrimSuffix(name, ".go")
	name = strings.TrimSuffix(name, ".ts")
	name = strings.TrimSuffix(name, ".js")
	return name
}

func describeChanges(files []string) string {
	if len(files) == 1 {
		return "in " + files[0]
	}
	return fmt.Sprintf("across %d files", len(files))
}

func styleVerbs(style string) (verbs [3]string, prefix string) {
	switch style {
	case config.StyleGitmoji:
		return [3]string{"feat", "refactor", "chore"}, "✨ "
	case config.StyleFree:
		return [3]string{"add", "update", "adjust"}, ""
	default:
		return [3]string{"feat", "refactor", "chore"}, ""
	}
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
