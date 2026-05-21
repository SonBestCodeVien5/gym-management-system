# Backend Delivery Memory

This folder keeps backend feature memory. Workflow rules live in repo-scoped Codex skills, not in
prompt alias documents here.

| Phase | Skill |
|---|---|
| Plan | `$gym-plan` |
| Implement | `$gym-implement` |
| Review | `$gym-review` |
| Test | `$gym-test` |
| Complete | `$gym-complete` |
| Resume | `$gym-resume` |
| Status | `$gym-status` |

## Keep

| Path | Purpose |
|---|---|
| `plans/` | Feature plans and API/business-rule decisions |
| `implementations/` | Implementation notes and handoff state |
| `reviews/` | Review findings and fixes |
| `tests/` | Verification notes |
| `worklog.md` | Short chronology and current roadmap |

## Read Narrowly

1. Read `CHAT_CONTEXT/README.md`.
2. Use the focused backend phase skill.
3. Read the current feature plan and only the phase notes needed for the task.
4. Inspect relevant source files and docs before editing.

## Update Rules

- Keep feature notes concise: decisions, files changed, commands, risks, and next action.
- Do not paste full source files, huge diffs, secrets, or local scratch output.
- Update `docs/api_contract.md` and `api_test.http` when HTTP behavior changes.
- Update `CHAT_CONTEXT/README.md` when project-level state or the next resume point changes.
- Keep scratch/session-only notes out of tracked memory.
