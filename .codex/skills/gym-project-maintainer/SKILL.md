---
name: gym-project-maintainer
description: Maintain cross-surface docs, report material, chat context, and backend memory boundaries for this gym management system. Use when Codex must coordinate more than one gym documentation/context surface or decide which project material to load; use the phase-specific gym skills for focused plan, implement, review, test, completion, docs, or report work.
---

# Gym Project Maintainer

## Start

1. Read `docs/README.md`.
2. Read `CHAT_CONTEXT/README.md`.
3. Read only the source and context files required by the active task.

## Choose The Lane

Use the focused skill when one lane is clear:

| Task | Skill |
|---|---|
| Plan backend feature | `$gym-plan` |
| Implement backend feature | `$gym-implement` |
| Review backend feature | `$gym-review` |
| Test backend feature | `$gym-test` |
| Finalize feature docs/context | `$gym-complete` |
| Resume backend work | `$gym-resume` |
| Summarize backend status | `$gym-status` |
| Maintain durable docs | `$gym-docs` |
| Prepare report material | `$gym-report` |
| Git/GitHub and version control | `$gym-git` |

### Documentation

- Keep durable project docs in `docs/`.
- Update `docs/api_contract.md` and `api_test.http` when HTTP behavior changes.
- Check handlers, services, routes, and models before documenting exact behavior.
- Update `docs/README.md` when a new durable document becomes an entrypoint.

### Report Material

- Keep report chapter inputs in `docs/report-materials/`.
- Reconcile planned report text with code and `docs/api_contract.md` before claiming a behavior is
  implemented.
- Use `docs/system_analysis_design_guide.md` and `docs/faq_why.md` when report structure or
  rationale is needed.

### Chat Context

- Keep `CHAT_CONTEXT/README.md` short enough to resume a new chat.
- Store status, decisions, risks, and next files to load. Do not paste code, secrets, long logs, or
  report drafts into chat context.
- Move durable explanations into `docs/` and report drafts into `docs/report-materials/`.

### Backend Delivery

Use the backend phase skill that matches the active task. Those skills read
[references/backend-delivery.md](references/backend-delivery.md) for shared rules and feature
memory layout.

## Quality Rules

- Prefer one clear source of truth over duplicate summaries.
- Treat source code and route wiring as truth for implemented behavior, then align API docs and
  context summaries.
- Preserve report drafts when they represent target scope, but label planned behavior clearly.
- Keep context loading narrow; do not read every feature log for one task.
- Verify docs structure after moving files by searching for stale paths.
- Validate this skill after changing its metadata or bundled references.
