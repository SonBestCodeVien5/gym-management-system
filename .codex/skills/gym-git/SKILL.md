---
name: gym-git
description: Manage Git, GitHub, and version-control workflow for this gym management system. Use when Codex is asked to inspect branch/worktree status, prepare commits, review diffs/history, create branch or PR handoff guidance, use `gh` for repository workflow, or respond to an explicit `$gym-git` or `/gym-git` request.
---

# Gym Git

## Read First

1. Run `git status --short`.
2. Check the current branch with `git branch --show-current`.
3. Inspect remotes with `git remote -v` when GitHub or push/pull/PR work is requested.
4. Read `git diff --stat` and targeted `git diff` output before staging, committing, or describing
   changes.
5. Read `docs/README.md` or affected docs only when the version-control task needs project context.

## Focus

- Keep Git operations explicit, reviewable, and scoped to the user's requested change.
- Preserve dirty worktree changes unless the user clearly asks to remove or rewrite them.
- Separate inspect, stage, commit, push, PR, and release actions instead of bundling risky steps.
- Use GitHub tooling only when GitHub state is needed; prefer local Git for local truth.

## Safety Rules

- Do not run destructive commands such as `git reset --hard`, `git checkout --`, branch deletion, or
  file removal unless the user clearly requests that operation.
- Do not revert edits that were not made for the current task.
- Do not stage unrelated files when preparing a commit.
- Do not commit or push unless the user asks for it.
- Before commit, summarize staged scope and use a message that matches the actual diff.
- Before push, confirm branch intent when on `main` or when remote effects are ambiguous.
- Never expose secrets from `.env`, credentials, tokens, or remote auth configuration.

## Task Playbook

### Inspect

- Use `git status --short`, `git diff --stat`, targeted `git diff`, `git log --oneline`, or
  `git show` as needed.
- Report important lines because the user does not see command output.

### Stage And Commit

1. Identify files owned by the requested change.
2. Review their diff.
3. Stage only that scope.
4. Re-check staged diff.
5. Commit only when requested and report the commit hash/message.

### GitHub

- Check local branch and remote before PR or push guidance.
- Use `gh` only when the requested GitHub operation needs live remote state and permissions allow it.
- For PR summaries, include scope, verification, docs/context changes, and residual risk.

## Output Rules

- State worktree/branch conditions that affect the requested operation.
- Distinguish local state from remote GitHub state.
- Record commands not run when blocked by permissions, network, or missing GitHub tooling.
