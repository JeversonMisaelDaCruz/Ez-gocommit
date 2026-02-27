package ai

import (
	"strings"
	"testing"

	"github.com/jeversonmisael/ez-gocommit/internal/git"
)

func TestSystemPrompt_NotEmpty(t *testing.T) {
	p := SystemPrompt()
	if strings.TrimSpace(p) == "" {
		t.Error("SystemPrompt() returned empty string")
	}
}

func TestSystemPrompt_ContainsRequiredSections(t *testing.T) {
	p := SystemPrompt()

	required := []string{
		"conventional",
		"gitmoji",
		"free",
		"custom",
		`"suggestions"`,
		`"rank"`,
		`"confidence"`,
		`"message"`,
		`"reasoning"`,
	}

	for _, s := range required {
		if !strings.Contains(p, s) {
			t.Errorf("SystemPrompt() missing expected content: %q", s)
		}
	}
}

func TestBuildUserPrompt_AllPlaceholdersFilled(t *testing.T) {
	ctx := &git.Context{
		BranchName:     "feat/login",
		StagedDiff:     "diff --git a/main.go ...",
		ChangedFiles:   []string{"main.go", "auth/handler.go"},
		RecentCommits:  []string{"feat: add user model", "fix: correct typo"},
		ProjectContext: "# MyApp\nA web application.",
	}

	prompt := BuildUserPrompt(ctx, "conventional")

	checks := map[string]string{
		"commit_style":    "conventional",
		"branch_name":     "feat/login",
		"git_diff":        "diff --git a/main.go",
		"changed_files":   "main.go",
		"recent_commits":  "feat: add user model",
		"project_context": "MyApp",
	}

	for label, expected := range checks {
		if !strings.Contains(prompt, expected) {
			t.Errorf("BuildUserPrompt() missing %s content %q", label, expected)
		}
	}
}

func TestBuildUserPrompt_NoRemainingPlaceholders(t *testing.T) {
	ctx := &git.Context{
		BranchName:     "main",
		StagedDiff:     "+ added line",
		ChangedFiles:   []string{"file.go"},
		RecentCommits:  []string{"initial commit"},
		ProjectContext: "",
	}

	prompt := BuildUserPrompt(ctx, "free")

	if strings.Contains(prompt, "{{") || strings.Contains(prompt, "}}") {
		t.Errorf("BuildUserPrompt() left unreplaced placeholders in output:\n%s", prompt)
	}
}

func TestBuildUserPrompt_MultipleChangedFiles(t *testing.T) {
	ctx := &git.Context{
		BranchName:   "main",
		StagedDiff:   "some diff",
		ChangedFiles: []string{"a.go", "b.go", "c.go"},
	}

	prompt := BuildUserPrompt(ctx, "conventional")

	for _, f := range ctx.ChangedFiles {
		if !strings.Contains(prompt, f) {
			t.Errorf("BuildUserPrompt() missing changed file %q", f)
		}
	}
}

func TestBuildUserPrompt_GitmojStyle(t *testing.T) {
	ctx := &git.Context{
		BranchName:   "fix/crash",
		StagedDiff:   "- bad line\n+ good line",
		ChangedFiles: []string{"server.go"},
	}

	prompt := BuildUserPrompt(ctx, "gitmoji")

	if !strings.Contains(prompt, "gitmoji") {
		t.Error("BuildUserPrompt() should include gitmoji style in output")
	}
}
