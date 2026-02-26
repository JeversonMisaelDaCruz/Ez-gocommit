package git

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Context struct {
	BranchName     string
	StagedDiff     string
	ChangedFiles   []string
	RecentCommits  []string
	ProjectContext string
}

var ErrNoStagedChanges = errors.New("no staged changes found — run `git add` first")

func Collect(repoPath string, maxDiffLines int) (*Context, error) {
	repo, err := gogit.PlainOpenWithOptions(repoPath, &gogit.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	branch, err := getBranchName(repo)
	if err != nil {
		branch = "unknown"
	}

	diff, files, err := getStagedDiff(repo, maxDiffLines)
	if err != nil {
		return nil, err
	}
	if diff == "" {
		return nil, ErrNoStagedChanges
	}

	commits, err := getRecentCommits(repo, 10)
	if err != nil {
		commits = []string{}
	}

	projectCtx := getProjectContext(repoPath)

	return &Context{
		BranchName:     branch,
		StagedDiff:     diff,
		ChangedFiles:   files,
		RecentCommits:  commits,
		ProjectContext: projectCtx,
	}, nil
}

func getBranchName(repo *gogit.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	if head.Name().IsBranch() {
		return head.Name().Short(), nil
	}
	return head.Hash().String()[:8], nil
}

func getStagedDiff(repo *gogit.Repository, maxLines int) (string, []string, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return "", nil, fmt.Errorf("cannot open worktree: %w", err)
	}

	status, err := wt.Status()
	if err != nil {
		return "", nil, fmt.Errorf("cannot get status: %w", err)
	}

	var stagedFiles []string
	for path, s := range status {
		if s.Staging != gogit.Unmodified && s.Staging != gogit.Untracked {
			stagedFiles = append(stagedFiles, path)
		}
	}

	if len(stagedFiles) == 0 {
		return "", nil, nil
	}

	head, err := repo.Head()
	var headCommit *object.Commit
	if err == nil {
		headCommit, err = repo.CommitObject(head.Hash())
		if err != nil {
			headCommit = nil
		}
	}

	var diffBuf bytes.Buffer

	if headCommit != nil {
		headTree, err := headCommit.Tree()
		if err != nil {
			return buildSimpleDiff(stagedFiles), stagedFiles, nil
		}

		idx, err := repo.Storer.Index()
		if err != nil {
			return buildSimpleDiff(stagedFiles), stagedFiles, nil
		}

		idxTree := &object.Tree{}
		for _, entry := range idx.Entries {
			idxTree.Entries = append(idxTree.Entries, object.TreeEntry{
				Name: entry.Name,
				Mode: entry.Mode,
				Hash: entry.Hash,
			})
		}

		changes, err := headTree.Diff(idxTree)
		if err != nil {
			return buildSimpleDiff(stagedFiles), stagedFiles, nil
		}

		for _, change := range changes {
			patch, err := change.Patch()
			if err != nil {
				continue
			}
			diffBuf.WriteString(patch.String())
		}
	} else {
		return buildSimpleDiff(stagedFiles), stagedFiles, nil
	}

	diffStr := truncateLines(diffBuf.String(), maxLines)
	return diffStr, stagedFiles, nil
}

func buildSimpleDiff(files []string) string {
	var sb strings.Builder
	sb.WriteString("Initial commit — new files added:\n")
	for _, f := range files {
		sb.WriteString("  + ")
		sb.WriteString(f)
		sb.WriteString("\n")
	}
	return sb.String()
}

func getRecentCommits(repo *gogit.Repository, n int) ([]string, error) {
	head, err := repo.Head()
	if err != nil {
		return []string{}, nil
	}

	iter, err := repo.Log(&gogit.LogOptions{From: head.Hash()})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var messages []string
	count := 0
	err = iter.ForEach(func(c *object.Commit) error {
		if count >= n {
			return plumbing.ErrObjectNotFound
		}
		lines := strings.SplitN(c.Message, "\n", 2)
		messages = append(messages, lines[0])
		count++
		return nil
	})
	if err != nil && !errors.Is(err, plumbing.ErrObjectNotFound) {
		return messages, err
	}
	return messages, nil
}

func getProjectContext(repoPath string) string {
	candidates := []string{"README.md", "readme.md", "README.rst", "README"}
	for _, name := range candidates {
		path := repoPath + "/" + name
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		var sb strings.Builder
		scanner := bufio.NewScanner(f)
		lineCount := 0
		for scanner.Scan() && lineCount < 100 {
			sb.WriteString(scanner.Text())
			sb.WriteString("\n")
			lineCount++
		}
		return sb.String()
	}
	return ""
}

func truncateLines(s string, maxLines int) string {
	if maxLines <= 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	truncated := strings.Join(lines[:maxLines], "\n")
	return truncated + fmt.Sprintf("\n\n[... diff truncated at %d lines ...]", maxLines)
}
