---
description: Generate semantic Git commit messages using Claude AI. Use when the user wants to create a commit, generate a commit message, or asks "commit my changes".
argument-hint: "[optional extra context]"
allowed-tools: Bash(git *)
---

Generate a semantic commit message for the current staged changes using the Conventional Commits standard.

## Steps

1. Run `git diff --cached --stat` to check for staged changes.

2. If there are **no staged changes**, inform the user:
   > "No staged changes found. Stage your files first with `git add <files>` and try again."

3. If there **are staged changes**, collect context:
   - `git diff --cached` — full staged diff
   - `git branch --show-current` — current branch name
   - `git log --oneline -10` — last 10 commits (style reference)

4. Analyze the diff and generate **3 commit message suggestions** following Conventional Commits:
   - Each suggestion has a **subject** (type(scope): description) and an optional **body**
   - Keep subjects under 72 characters
   - Use types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`, `perf`, `style`
   - Consider the branch name and recent commits to match the project's style
   - If the user passed `$ARGUMENTS`, use them as extra context for the suggestions

5. Present the suggestions clearly, numbered 1–3, and ask the user which one to use (or if they want to edit).

6. Once the user confirms a message, run:
   ```
   git commit -m "<subject>" -m "<body>"
   ```
   (omit `-m "<body>"` if there is no body)

## Notes

- Do NOT commit automatically without user confirmation.
- If the user asks to edit a suggestion, apply the edit and confirm before committing.
- No API key or external binary required — Claude Code handles everything.
