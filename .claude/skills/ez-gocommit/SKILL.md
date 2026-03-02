---
name: ez-gocommit
description: Generate semantic Git commit messages using Claude AI via ez-gocommit. Use when the user wants to create a commit, generate a commit message, or asks "commit my changes".
disable-model-invocation: true
argument-hint: "[optional extra context]"
allowed-tools: Bash(git *), Bash(ezgocommit *)
---

Generate a semantic commit message for the current staged changes using ez-gocommit.

## Steps

1. Check for staged changes:
   ```
   git diff --cached --stat
   ```

2. If there are **no staged changes**, inform the user:
   > "No staged changes found. Stage your files first with `git add <files>` and try again."

3. If there **are staged changes**, run:
   ```
   ezgocommit $ARGUMENTS
   ```

## How the tool works

`ezgocommit` analyzes the staged diff, branch name, recent commits, and README to generate 3 commit message suggestions following the Conventional Commits standard.

- Use **↑ ↓** to navigate suggestions
- Press **Enter** to confirm and commit
- Press **e** to edit the message inline
- Press **q** to abort without committing

## Requirements

- `ANTHROPIC_API_KEY` env var must be set, or `api_key` configured in the config file.
- The binary `ezgocommit` must be installed and available in PATH.
