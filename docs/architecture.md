# Architecture

## Overview

Ez-gocommit is a single-binary Go CLI tool. It has no daemon, no server, and no persistent state beyond config files. Each invocation runs the full pipeline and exits.

## Project structure

```
ez-gocommit/
├── main.go                      # Entry point — calls cmd.Execute()
├── go.mod
├── go.sum
│
├── cmd/
│   ├── root.go                  # Cobra root command + flag definitions
│   ├── generate.go              # Main pipeline: collect → AI → TUI → commit
│   └── version.go               # `ezgocommit version` subcommand
│
└── internal/
    ├── config/
    │   └── config.go            # Load and validate configuration
    │
    ├── git/
    │   └── collector.go         # Collect git context from the repository
    │
    ├── ai/
    │   ├── types.go             # Suggestion and AIResponse structs
    │   ├── prompt.go            # System prompt + BuildUserPrompt()
    │   └── client.go            # Anthropic API call + JSON parsing
    │
    └── ui/
        └── selector.go          # Bubbletea interactive TUI
```

## Data flow

```
┌─────────────────────────────────────────────────────────────┐
│                        cmd/generate.go                       │
│                                                              │
│  1. config.LoadWithOverrides()                               │
│       ↓                                                      │
│  2. git.Collect()  ──────────────────────────────────────┐  │
│       reads: staged diff, changed files,                  │  │
│              branch name, 10 recent commits, README.md    │  │
│       ↓                                                   │  │
│  3. ai.BuildUserPrompt()  ←──────────────────────────────┘  │
│       fills {{PLACEHOLDERS}} in the user prompt template     │
│       ↓                                                      │
│  4. ai.GenerateSuggestions()                                 │
│       sends system prompt + user prompt to Claude API        │
│       parses JSON → []Suggestion                             │
│       ↓                                                      │
│  5. ui.Run()                                                 │
│       Bubbletea TUI: navigate, pick, optionally edit         │
│       returns (message, body, cancelled)                     │
│       ↓                                                      │
│  6. git commit -m "message\n\nbody"                          │
└─────────────────────────────────────────────────────────────┘
```

## Package responsibilities

### `internal/config`

Loads configuration using Viper. Reads from:
- Global file: `~/.config/ezgocommit/config.toml`
- Local file: `.ezgocommit.toml` (current directory)
- Environment variable: `ANTHROPIC_API_KEY`

CLI flags applied by `cmd` after loading override the file values.

### `internal/git`

Wraps `go-git/v5` to extract everything needed for a meaningful prompt:

| Function | What it collects |
|----------|-----------------|
| `getBranchName` | Current branch short name |
| `getStagedDiff` | Unified diff of HEAD vs index (truncated to `max_diff_lines`) |
| `getRecentCommits` | Last 10 commit subject lines |
| `getProjectContext` | First 100 lines of `README.md` |

For repositories with no commits yet (initial commit), `getStagedDiff` falls back to a file list rather than a real patch.

### `internal/ai`

**`prompt.go`** holds two string constants:
- `systemPrompt` — the full instruction set sent as the Claude system turn. Defines rules, supported styles, and the strict JSON output format.
- `userPromptTemplate` — the user turn template with `{{PLACEHOLDER}}` tokens replaced at runtime.

**`client.go`** calls `client.Messages.New()` from the official `anthropic-sdk-go`. It strips any accidental markdown fences from the response before JSON parsing.

**`types.go`** defines `Suggestion` (one option) and `AIResponse` (the full parsed response).

### `internal/ui`

A self-contained [Bubbletea](https://github.com/charmbracelet/bubbletea) program with two modes:

- **`modeSelect`** — arrow key / vim key navigation, number shortcuts (1-3), `e` to enter edit
- **`modeEdit`** — inline text editing with left/right movement, `Enter` to confirm, `Esc` to cancel

Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss). Confidence badges are color-coded:
- `●● HIGH` → green
- `●○ MED` → yellow
- `○○ LOW` → red

### `cmd`

Thin orchestration layer built with Cobra. `root.go` registers flags and wires the default command to `runGenerate`. `generate.go` owns the full pipeline in sequence. No business logic lives here.

## Dependencies

| Package | Version | Role |
|---------|---------|------|
| `github.com/spf13/cobra` | v1.10.2 | CLI command structure |
| `github.com/spf13/viper` | v1.21.0 | Config file + env var loading |
| `github.com/go-git/go-git/v5` | v5.17.0 | Git operations (pure Go) |
| `github.com/anthropics/anthropic-sdk-go` | v1.26.0 | Anthropic API client |
| `github.com/charmbracelet/bubbletea` | v1.3.10 | TUI framework |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | TUI styling |
| `github.com/fatih/color` | v1.18.0 | Colored terminal output |

## AI prompt design

The system prompt instructs Claude to:
- Analyze the diff deeply (not just filenames)
- Use the branch name as an intent hint
- Mirror the tone and style of recent commits
- Understand the project domain from the README
- Produce exactly 3 JSON-serialized suggestions ranked by confidence
- Never output anything outside the JSON object

The user prompt wraps the runtime context in XML-like tags (`<git_diff>`, `<branch_name>`, etc.) to give Claude clear boundaries between each piece of data.

## Error handling strategy

- Missing API key → clear error with setup instructions, exit 1
- No staged changes → clear error prompting `git add`, exit 1
- Not a git repository → error from go-git, exit 1
- API error → wrapped error with original message, exit 1
- Malformed JSON from AI → error with raw response for debugging, exit 1
- User aborts TUI → prints "Aborted.", exit 0
