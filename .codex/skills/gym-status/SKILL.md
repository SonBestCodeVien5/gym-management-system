---
name: gym-status
description: Summarize backend status for this gym management system from current project memory. Use when Codex is asked for roadmap status, current cycle progress, next backend step, or an explicit `$gym-status` or `/gym-status` request without making implementation changes.
---

# Gym Status

## Read First

1. Read `CHAT_CONTEXT/README.md`.
2. Read `CHAT_CONTEXT/backend_skills/README.md`.
3. Read `CHAT_CONTEXT/backend_skills/worklog.md`.
4. Read only the current or named feature plan/status notes needed to answer accurately.

## Focus

- Report completed work, active/planned cycle, blockers, verification state, and next step.
- Prefer summary memory over broad source-code loading.
- Distinguish planned work from implemented and tested work.
- Keep the response short unless the user asks for a detailed status report.

## Output Rules

- Do not edit code or feature notes unless the user asks to update status memory.
- Do not infer completion from a plan alone.
- Mention stale or missing evidence when status notes disagree.
