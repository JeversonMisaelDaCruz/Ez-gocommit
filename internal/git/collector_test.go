package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func initTestRepo(t *testing.T) (string, *gogit.Repository) {
	t.Helper()
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("failed to init test repo: %v", err)
	}
	return dir, repo
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

func stageFile(t *testing.T, repo *gogit.Repository, name string) {
	t.Helper()
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("worktree error: %v", err)
	}
	if _, err := wt.Add(name); err != nil {
		t.Fatalf("failed to stage %s: %v", name, err)
	}
}

func makeCommit(t *testing.T, repo *gogit.Repository, msg string) {
	t.Helper()
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("worktree error: %v", err)
	}
	_, err = wt.Commit(msg, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@test.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
}

func TestCollect_NoStagedChanges(t *testing.T) {
	dir, _ := initTestRepo(t)

	_, err := Collect(dir, 500)
	if err == nil {
		t.Fatal("expected error for empty repo, got nil")
	}
	if err != ErrNoStagedChanges {
		t.Errorf("expected ErrNoStagedChanges, got: %v", err)
	}
}

func TestCollect_InitialCommitStagedFiles(t *testing.T) {
	dir, repo := initTestRepo(t)

	writeFile(t, dir, "main.go", `package main

func main() {}
`)
	stageFile(t, repo, "main.go")

	ctx, err := Collect(dir, 500)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if !contains(ctx.ChangedFiles, "main.go") {
		t.Errorf("ChangedFiles = %v, want to contain main.go", ctx.ChangedFiles)
	}
	if ctx.StagedDiff == "" {
		t.Error("StagedDiff should not be empty")
	}
	if !strings.Contains(ctx.StagedDiff, "main.go") {
		t.Errorf("StagedDiff should mention main.go, got: %s", ctx.StagedDiff)
	}
}

func TestCollect_BranchName(t *testing.T) {
	dir, repo := initTestRepo(t)

	writeFile(t, dir, "file.go", "package main")
	stageFile(t, repo, "file.go")
	makeCommit(t, repo, "initial commit")

	writeFile(t, dir, "file.go", "package main\n\nfunc Foo() {}")
	stageFile(t, repo, "file.go")

	ctx, err := Collect(dir, 500)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if ctx.BranchName == "" {
		t.Error("BranchName should not be empty")
	}
}

func TestCollect_StagedDiffAfterCommit(t *testing.T) {
	dir, repo := initTestRepo(t)

	writeFile(t, dir, "service.go", "package main\n\nfunc Run() {}")
	stageFile(t, repo, "service.go")
	makeCommit(t, repo, "feat: initial service")

	writeFile(t, dir, "service.go", "package main\n\nfunc Run() {}\n\nfunc Stop() {}")
	stageFile(t, repo, "service.go")

	ctx, err := Collect(dir, 500)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if ctx.StagedDiff == "" {
		t.Error("StagedDiff should not be empty after modifying a tracked file")
	}
	if len(ctx.ChangedFiles) == 0 {
		t.Error("ChangedFiles should not be empty")
	}
}

func TestCollect_RecentCommits(t *testing.T) {
	dir, repo := initTestRepo(t)

	messages := []string{
		"feat: first feature",
		"fix: correct bug",
		"chore: update deps",
	}

	for i, msg := range messages {
		fname := filepath.Base(dir) + "_file" + string(rune('a'+i)) + ".go"
		writeFile(t, dir, fname, "package main")
		stageFile(t, repo, fname)
		makeCommit(t, repo, msg)
	}

	writeFile(t, dir, "new.go", "package main\nfunc New() {}")
	stageFile(t, repo, "new.go")

	ctx, err := Collect(dir, 500)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if len(ctx.RecentCommits) == 0 {
		t.Error("RecentCommits should not be empty")
	}
	if !contains(ctx.RecentCommits, "chore: update deps") {
		t.Errorf("RecentCommits = %v, expected to contain recent commits", ctx.RecentCommits)
	}
}

func TestCollect_ProjectContext_README(t *testing.T) {
	dir, repo := initTestRepo(t)

	readmeContent := "# TestProject\nThis is a test project for Ez-gocommit.\n"
	writeFile(t, dir, "README.md", readmeContent)
	writeFile(t, dir, "main.go", "package main")
	stageFile(t, repo, "main.go")

	ctx, err := Collect(dir, 500)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if !strings.Contains(ctx.ProjectContext, "TestProject") {
		t.Errorf("ProjectContext should contain README content, got: %q", ctx.ProjectContext)
	}
}

func TestCollect_ProjectContext_NoREADME(t *testing.T) {
	dir, repo := initTestRepo(t)

	writeFile(t, dir, "main.go", "package main")
	stageFile(t, repo, "main.go")

	ctx, err := Collect(dir, 500)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if ctx.ProjectContext != "" {
		t.Errorf("ProjectContext should be empty when no README exists, got: %q", ctx.ProjectContext)
	}
}

func TestCollect_DiffTruncation(t *testing.T) {
	dir, repo := initTestRepo(t)

	writeFile(t, dir, "big.go", "package main\n")
	stageFile(t, repo, "big.go")
	makeCommit(t, repo, "chore: initial file")

	var sb strings.Builder
	sb.WriteString("package main\n\nvar data = []string{\n")
	for i := 0; i < 600; i++ {
		sb.WriteString(`  "line",` + "\n")
	}
	sb.WriteString("}\n")
	writeFile(t, dir, "big.go", sb.String())
	stageFile(t, repo, "big.go")

	const maxLines = 50
	ctx, err := Collect(dir, maxLines)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	lineCount := strings.Count(ctx.StagedDiff, "\n")
	if lineCount > maxLines+5 {
		t.Errorf("StagedDiff should be truncated: got %d lines, max was %d", lineCount, maxLines)
	}
	if !strings.Contains(ctx.StagedDiff, "truncated") {
		t.Error("Truncated diff should contain a truncation notice")
	}
}

func TestCollect_MultipleChangedFiles(t *testing.T) {
	dir, repo := initTestRepo(t)

	files := []string{"a.go", "b.go", "c.go"}
	for _, f := range files {
		writeFile(t, dir, f, "package main")
		stageFile(t, repo, f)
	}

	ctx, err := Collect(dir, 500)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if len(ctx.ChangedFiles) != len(files) {
		t.Errorf("ChangedFiles count = %d, want %d", len(ctx.ChangedFiles), len(files))
	}
}

func TestCollect_NotARepository(t *testing.T) {
	dir := t.TempDir()

	_, err := Collect(dir, 500)
	if err == nil {
		t.Error("Collect() should return error for non-git directory")
	}
}

func TestTruncateLines_BelowLimit(t *testing.T) {
	input := "line1\nline2\nline3"
	result := truncateLines(input, 10)
	if result != input {
		t.Errorf("truncateLines() should not modify input below limit")
	}
}

func TestTruncateLines_AboveLimit(t *testing.T) {
	lines := make([]string, 20)
	for i := range lines {
		lines[i] = "line"
	}
	input := strings.Join(lines, "\n")
	result := truncateLines(input, 5)
	if !strings.Contains(result, "truncated") {
		t.Error("truncateLines() should add truncation notice")
	}
	if strings.Count(result, "\n") >= 20 {
		t.Error("truncateLines() should reduce line count")
	}
}

func TestTruncateLines_ZeroLimit(t *testing.T) {
	input := "line1\nline2\nline3"
	result := truncateLines(input, 0)
	if result != input {
		t.Error("truncateLines() with limit=0 should return input unchanged")
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
