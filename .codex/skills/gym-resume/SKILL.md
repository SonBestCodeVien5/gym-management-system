---
name: gym-resume
description: Resume backend work for this gym management system from current chat context and backend feature memory. Use when Codex is asked to continue an interrupted backend task, infer the active phase from notes and workspace state, or respond to an explicit `$gym-resume` or `/gym-resume` request.
---

# Gym Resume

## Read First

1. Read `CHAT_CONTEXT/README.md`.
2. Read `CHAT_CONTEXT/backend_skills/README.md`.
3. Read `../gym-project-maintainer/references/backend-delivery.md`.
4. Read `CHAT_CONTEXT/backend_skills/worklog.md`.
5. Read the current feature memory files named by the snapshot, worklog, or user.
6. Inspect `git status --short` and changed files before continuing edits.

## Focus

- Reconstruct the active feature, active phase, last verified result, and next concrete step.
- Load only the memory and source files needed to continue safely.
- Work with existing workspace changes; do not revert or overwrite unrelated edits.
- Hand off to the focused skill mentally: plan, implement, review, test, or complete.

## Output Rules

- Continue the requested work when the next action is clear.
- Update the active feature note when resume changes status, decisions, commands, or blockers.
- State the active feature and phase before major edits when the resume state is ambiguous.
- Ask only when the repo and memory do not reveal a safe next step.
