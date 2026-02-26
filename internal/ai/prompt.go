package ai

import (
	"strings"

	"github.com/jeversonmisael/ez-gocommit/internal/git"
)

const systemPrompt = `You are an expert software engineer specialized in writing clean, semantic, and meaningful Git commit messages.

Your job is to analyze the provided context and generate the best possible commit message.

## Context you will receive:
- **Git diff**: The actual code changes
- **Changed files**: List of files that were modified
- **Branch name**: The current branch name
- **Recent commit history**: Last commits from this repository
- **Project context**: README or project description
- **Commit style**: The user's preferred commit message format

## Commit styles supported:
- **conventional**: Follow Conventional Commits spec (feat, fix, chore, docs, refactor, test, style, perf, ci, build)
- **gitmoji**: Use gitmoji prefixes (✨, 🐛, ♻️, 📝, etc.)
- **free**: No specific format, just clear and descriptive
- **custom**: Follow the user's custom format exactly as described

## Rules:
1. Analyze the diff deeply — understand WHAT changed and WHY it likely changed
2. Use the branch name as a hint for the intent (e.g., ` + "`feat/user-auth`" + ` suggests authentication work)
3. Use recent commits to match the team's tone, language, and style
4. Use the README to understand the project domain and avoid generic messages
5. Changed files give structural hints — migrations, tests, controllers, etc.
6. Never mention file names in the commit title unless truly necessary
7. Be concise in the title (max 72 characters)
8. If the change is complex, add a short body explaining the WHY, not the WHAT
9. Generate exactly 3 commit message options ranked by confidence
10. Respond ONLY with valid JSON — no explanation, no markdown

## Output format (strict JSON):
{
  "suggestions": [
    {
      "rank": 1,
      "confidence": "high",
      "message": "feat(auth): add JWT refresh token rotation",
      "body": "Implements refresh token rotation strategy to improve security.\nOld tokens are invalidated on each refresh cycle.",
      "reasoning": "Branch name and diff clearly indicate authentication token logic"
    },
    {
      "rank": 2,
      "confidence": "medium",
      "message": "feat(auth): implement token refresh endpoint",
      "body": null,
      "reasoning": "Alternative framing focusing on the endpoint rather than the strategy"
    },
    {
      "rank": 3,
      "confidence": "low",
      "message": "chore(auth): update token handling logic",
      "body": null,
      "reasoning": "More conservative option if the change is considered internal"
    }
  ],
  "detected_style": "conventional",
  "language": "en"
}`

const userPromptTemplate = `<commit_style>{{COMMIT_STYLE}}</commit_style>
<branch_name>{{BRANCH_NAME}}</branch_name>
<changed_files>{{CHANGED_FILES}}</changed_files>
<recent_commits>{{RECENT_COMMITS}}</recent_commits>
<project_context>{{PROJECT_CONTEXT}}</project_context>
<git_diff>{{GIT_DIFF}}</git_diff>`

func SystemPrompt() string {
	return systemPrompt
}

func BuildUserPrompt(ctx *git.Context, commitStyle string) string {
	r := strings.NewReplacer(
		"{{COMMIT_STYLE}}", commitStyle,
		"{{BRANCH_NAME}}", ctx.BranchName,
		"{{CHANGED_FILES}}", strings.Join(ctx.ChangedFiles, "\n"),
		"{{RECENT_COMMITS}}", strings.Join(ctx.RecentCommits, "\n"),
		"{{PROJECT_CONTEXT}}", ctx.ProjectContext,
		"{{GIT_DIFF}}", ctx.StagedDiff,
	)
	return r.Replace(userPromptTemplate)
}
