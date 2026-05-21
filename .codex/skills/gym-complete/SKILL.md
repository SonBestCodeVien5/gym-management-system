---
name: gym-complete
description: Finalize backend feature documentation and context for this gym management system. Use when Codex is asked to complete a backend cycle, align API docs and samples after implementation/test, update worklog/chat context, or respond to an explicit `$gym-complete` or `/gym-complete` request.
---

# Gym Complete

## Read First

1. Read `docs/README.md`.
2. Read `CHAT_CONTEXT/README.md`.
3. Read `CHAT_CONTEXT/backend_skills/README.md`.
4. Read `../gym-project-maintainer/references/backend-delivery.md`.
5. Read the plan, implementation note, review note, and test note for `<feature>`.
6. Read `docs/api_contract.md`, `api_test.http`, and the actual changed code that defines the
   shipped behavior.

## Focus

- Align durable docs, API samples, backend memory, and chat resume state with code that exists.
- Distinguish complete, skipped, planned, and remaining-risk work.
- Keep chat context short and keep detailed phase evidence in feature memory notes.
- Catch contract/docs drift before the feature leaves the cycle.

## Output Rules

- Update `docs/api_contract.md` and `api_test.http` when HTTP behavior changed.
- Update `CHAT_CONTEXT/backend_skills/worklog.md` and `CHAT_CONTEXT/README.md` when project state
  or next resume point changed.
- Update durable docs touched by the feature, but do not rewrite report drafts unless asked.
- Do not mark a feature tested or complete without the matching evidence or a recorded skip reason.
