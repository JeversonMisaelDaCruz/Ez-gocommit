# Contributing

## Development setup

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go mod download
```

Build:

```bash
go build -o ezgocommit .
```

Run without installing:

```bash
ANTHROPIC_API_KEY=sk-ant-... go run .
```

## Running checks

```bash
go build ./...   # must compile with no errors
go vet ./...     # must pass with no warnings
```

## Project layout

```
cmd/            CLI commands (Cobra) — no business logic
internal/
  config/       configuration loading
  git/          git context collection (go-git)
  ai/           Anthropic API integration
  ui/           Bubbletea TUI
docs/           project documentation
```

See [architecture.md](architecture.md) for a detailed explanation of each package.

## Adding a new commit style

1. Add a constant in `internal/config/config.go`:
   ```go
   StyleMyStyle = "mystyle"
   ```

2. Document the style in `internal/ai/prompt.go` inside the `systemPrompt` constant under the `## Commit styles supported:` section.

3. Add the style to the configuration reference in `docs/configuration.md`.

## Modifying the AI prompt

The system prompt lives in `internal/ai/prompt.go` as the `systemPrompt` constant. The user prompt template is `userPromptTemplate` in the same file.

Rules for prompt changes:
- The output format must remain strict JSON matching `AIResponse` in `internal/ai/types.go`
- Do not change placeholder names (`{{GIT_DIFF}}`, etc.) without updating `BuildUserPrompt`
- Test manually with a real API key and a staged change before opening a PR

## Changing the TUI

The TUI lives entirely in `internal/ui/selector.go`. It follows the standard Bubbletea model/update/view pattern:

- `model` — state
- `Update()` — handles key events, dispatches to `updateSelect` or `updateEdit`
- `View()` — renders the current state as a string using Lipgloss

Colors are defined as package-level `lipgloss.Style` variables at the top of the file.

## Commit messages

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(scope): short description
fix: correct something
docs: update configuration reference
refactor(git): simplify diff collection
```

Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `style`, `perf`

## Opening a pull request

1. Fork and create a branch: `git checkout -b feat/my-feature`
2. Make your changes
3. Run `go build ./...` and `go vet ./...`
4. Open a PR with a clear description of the change and why it was made
