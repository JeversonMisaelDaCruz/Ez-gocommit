# Getting Started

## Prerequisites

- **Go 1.22+** — [download](https://go.dev/dl/)
- **Git** — any recent version
- **Anthropic API key** — get one at [console.anthropic.com](https://console.anthropic.com/)

## Installation

### From source

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go build -o ezgocommit .
```

Move the binary to somewhere in your `$PATH`:

```bash
sudo mv ezgocommit /usr/local/bin/
```

### With go install

```bash
go install github.com/jeversonmisael/ez-gocommit@latest
```

### Build with version tag

```bash
go build -ldflags="-X github.com/jeversonmisael/ez-gocommit/cmd.Version=1.0.0" -o ezgocommit .
```

## Setting up the API key

The simplest way is an environment variable:

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

To make it permanent, add that line to your `~/.zshrc`, `~/.bashrc`, or equivalent.

Alternatively, create a config file at `~/.config/ezgocommit/config.toml`:

```toml
api_key = "sk-ant-..."
```

See [configuration.md](configuration.md) for details on all options.

## First use

Stage some changes and run the tool:

```bash
cd your-project
git add .
ezgocommit
```

The tool will:

1. Read your staged diff, branch name, and recent commit history
2. Send that context to Claude
3. Display 3 ranked commit message suggestions in an interactive UI
4. Commit the one you choose

## Verifying the install

```bash
ezgocommit version
```

## Running without an API key

If you want to test the binary without spending API credits, you can check that the tool detects the missing key correctly:

```bash
ANTHROPIC_API_KEY="" ezgocommit
# Error: Anthropic API key not found.
# ...
```

## Next steps

- [Configuration Reference](configuration.md) — customize model, style, diff limits
- [Architecture](architecture.md) — understand how the codebase is structured
- [Contributing](contributing.md) — how to add features or fix bugs
