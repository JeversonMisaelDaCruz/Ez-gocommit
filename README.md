# Ez-gocommit

A CLI tool written in Go that generates semantic Git commit messages using the Claude API (Anthropic). It analyzes your staged diff, branch name, recent commit history, and project README to produce 3 ranked suggestions — displayed in an interactive terminal UI where you can pick, edit, or abort.

## How it works

```
git add .  →  ezgocommit  →  [TUI with 3 ranked suggestions]  →  git commit
```

1. Reads your staged diff, changed files, branch name, recent commits, and `README.md`
2. Sends that context to Claude via the Anthropic API
3. Returns 3 ranked commit messages (high / medium / low confidence)
4. Lets you navigate, pick, or edit one inline — then commits automatically

## Install

**From source:**

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go build -o ezgocommit .
sudo mv ezgocommit /usr/local/bin/
```

**With `go install`:**

```bash
go install github.com/jeversonmisael/ez-gocommit@latest
```

## Requirements

- Go 1.22+
- An [Anthropic API key](https://console.anthropic.com/)
- A Git repository with staged changes

## Setup

The only required configuration is your API key.

**Option 1 — environment variable (recommended):**

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

Add it to your `~/.zshrc` or `~/.bashrc` to persist it.

**Option 2 — config file:**

Create `~/.config/ezgocommit/config.toml`:

```toml
api_key = "sk-ant-..."
```

See [docs/configuration.md](docs/configuration.md) for all available options.

## Usage

```bash
git add .
ezgocommit
```

```
⠸ Analyzing your changes with Claude...

╭──────────────────────────────────────────────────────────────────────╮
│  Ez-gocommit — Select a commit message                               │
│                                                                      │
│  ▶ [1] ●● HIGH   feat(auth): add JWT refresh token rotation          │
│    [2] ●○ MED    feat(auth): implement token refresh endpoint        │
│    [3] ○○ LOW    chore(auth): update token handling logic            │
│                                                                      │
│  💬 Branch name and diff clearly indicate authentication token logic │
│                                                                      │
│  ↑↓/jk navigate • 1-3 jump • Enter confirm • e edit • q abort       │
╰──────────────────────────────────────────────────────────────────────╯

✔ Committed: feat(auth): add JWT refresh token rotation
```

### TUI controls

| Key | Action |
|-----|--------|
| `↑` / `↓` or `j` / `k` | Navigate between suggestions |
| `1` / `2` / `3` | Jump directly to that suggestion |
| `Enter` | Confirm and commit |
| `e` | Edit the selected message inline |
| `q` / `Esc` / `Ctrl+C` | Abort |

**In edit mode:**

| Key | Action |
|-----|--------|
| `Enter` | Confirm edited message |
| `Esc` | Cancel edit, return to selection |
| `←` / `→` | Move cursor |
| `Ctrl+A` / `Home` | Go to start |
| `Ctrl+E` / `End` | Go to end |
| `Backspace` | Delete character |

### Flags

```bash
ezgocommit --style gitmoji       # use gitmoji instead of conventional commits
ezgocommit --style free          # no format constraints
ezgocommit --model claude-opus-4-6  # use a different Claude model
```

## Commit styles

| Style | Example |
|-------|---------|
| `conventional` (default) | `feat(auth): add login with OAuth` |
| `gitmoji` | `✨ add login with OAuth` |
| `free` | `Add OAuth login support` |
| `custom` | Defined by you in `custom_format` |

## Documentation

- [Getting Started](docs/getting-started.md)
- [Configuration Reference](docs/configuration.md)
- [Architecture](docs/architecture.md)
- [Contributing](docs/contributing.md)

## License

MIT
