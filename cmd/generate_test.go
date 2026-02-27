package cmd

import (
	"strings"
	"testing"

	gitcollector "github.com/jeversonmisael/ez-gocommit/internal/git"
)

func TestInferScope_DeepPath(t *testing.T) {
	files := []string{"internal/ai/client.go"}
	if got := inferScope(files); got != "ai" {
		t.Errorf("inferScope() = %q, want %q", got, "ai")
	}
}

func TestInferScope_TopLevelFile(t *testing.T) {
	files := []string{"main.go"}
	if got := inferScope(files); got != "main" {
		t.Errorf("inferScope() = %q, want %q", got, "main")
	}
}

func TestInferScope_Empty(t *testing.T) {
	if got := inferScope(nil); got != "core" {
		t.Errorf("inferScope(nil) = %q, want %q", got, "core")
	}
}

func TestDescribeChanges_Single(t *testing.T) {
	got := describeChanges([]string{"main.go"})
	if !strings.Contains(got, "main.go") {
		t.Errorf("describeChanges(single) = %q, should mention filename", got)
	}
}

func TestDescribeChanges_Multiple(t *testing.T) {
	got := describeChanges([]string{"a.go", "b.go", "c.go"})
	if !strings.Contains(got, "3") {
		t.Errorf("describeChanges(multiple) = %q, should mention count", got)
	}
}

func TestStyleVerbs_Conventional(t *testing.T) {
	verbs, prefix := styleVerbs("conventional")
	if prefix != "" {
		t.Errorf("conventional prefix = %q, want empty", prefix)
	}
	if verbs[0] != "feat" {
		t.Errorf("verbs[0] = %q, want feat", verbs[0])
	}
}

func TestStyleVerbs_Gitmoji(t *testing.T) {
	_, prefix := styleVerbs("gitmoji")
	if prefix == "" {
		t.Error("gitmoji prefix should not be empty")
	}
}

func TestStyleVerbs_Free(t *testing.T) {
	verbs, prefix := styleVerbs("free")
	if prefix != "" {
		t.Errorf("free prefix = %q, want empty", prefix)
	}
	if verbs[0] != "add" {
		t.Errorf("free verbs[0] = %q, want add", verbs[0])
	}
}

func TestMockSuggestions_ReturnThree(t *testing.T) {
	ctx := &gitcollector.Context{
		BranchName:   "feat/login",
		ChangedFiles: []string{"internal/auth/handler.go"},
	}
	suggestions := mockSuggestions(ctx, "conventional")
	if len(suggestions) != 3 {
		t.Fatalf("mockSuggestions() returned %d suggestions, want 3", len(suggestions))
	}
}

func TestMockSuggestions_RankedInOrder(t *testing.T) {
	ctx := &gitcollector.Context{
		BranchName:   "main",
		ChangedFiles: []string{"cmd/root.go"},
	}
	suggestions := mockSuggestions(ctx, "conventional")
	for i, s := range suggestions {
		if s.Rank != i+1 {
			t.Errorf("suggestions[%d].Rank = %d, want %d", i, s.Rank, i+1)
		}
	}
}

func TestMockSuggestions_ConfidenceLevels(t *testing.T) {
	ctx := &gitcollector.Context{
		BranchName:   "main",
		ChangedFiles: []string{"main.go"},
	}
	suggestions := mockSuggestions(ctx, "conventional")
	expected := []string{"high", "medium", "low"}
	for i, want := range expected {
		if suggestions[i].Confidence != want {
			t.Errorf("suggestions[%d].Confidence = %q, want %q", i, suggestions[i].Confidence, want)
		}
	}
}

func TestMockSuggestions_MessagesNotEmpty(t *testing.T) {
	ctx := &gitcollector.Context{
		BranchName:   "fix/bug",
		ChangedFiles: []string{"pkg/server/server.go", "pkg/server/handler.go"},
	}
	suggestions := mockSuggestions(ctx, "conventional")
	for i, s := range suggestions {
		if strings.TrimSpace(s.Message) == "" {
			t.Errorf("suggestions[%d].Message is empty", i)
		}
		if strings.TrimSpace(s.Reasoning) == "" {
			t.Errorf("suggestions[%d].Reasoning is empty", i)
		}
	}
}

func TestMockSuggestions_GitmojHasPrefix(t *testing.T) {
	ctx := &gitcollector.Context{
		BranchName:   "main",
		ChangedFiles: []string{"main.go"},
	}
	suggestions := mockSuggestions(ctx, "gitmoji")
	for _, s := range suggestions {
		if !strings.HasPrefix(s.Message, "✨") {
			t.Errorf("gitmoji suggestion should start with ✨, got: %q", s.Message)
		}
	}
}

func TestMockSuggestions_ReasoningMentionsBranch(t *testing.T) {
	ctx := &gitcollector.Context{
		BranchName:   "feat/payment",
		ChangedFiles: []string{"payment.go"},
	}
	suggestions := mockSuggestions(ctx, "conventional")
	if !strings.Contains(suggestions[0].Reasoning, "feat/payment") {
		t.Errorf("top suggestion reasoning should mention branch name, got: %q", suggestions[0].Reasoning)
	}
}
