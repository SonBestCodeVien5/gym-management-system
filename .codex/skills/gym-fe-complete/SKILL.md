---
name: gym-fe-complete
description: Finalize frontend feature documentation and context for this gym management system. Use when Codex is asked to complete a React/Vite frontend cycle, align UI/API docs and samples after implementation/test, update frontend worklog/chat context, or respond to an explicit `$gym-fe-complete` or `/gym-fe-complete` request.
---

# Gym FE Complete

## Read First

1. Read `docs/README.md`.
2. Read `CHAT_CONTEXT/README.md`.
3. Read `CHAT_CONTEXT/frontend_skills/README.md`.
4. Read the plan, implementation note, review note, and test note for `<feature>`.
5. Read `docs/api_contract.md`, relevant frontend docs, and actual changed code that defines shipped
   behavior.

## Focus

- Align frontend docs, API assumptions, frontend memory, and chat resume state with code that exists.
- Distinguish complete, skipped, planned, and remaining-risk work.
- Keep chat context short and keep detailed phase evidence in feature memory notes.
- Catch UI/API/docs drift before the feature leaves the cycle.
- For UI-bearing features, confirm the test note includes MCP Playwright/browser evidence or a clear
  skip reason before marking the cycle complete.

## Output Rules

- Update `CHAT_CONTEXT/frontend_skills/worklog.md` and `CHAT_CONTEXT/README.md` when project state
  or next resume point changed.
- Update durable docs touched by the feature, but do not rewrite report drafts unless asked.
- Do not mark a feature tested or complete without matching evidence or a recorded skip reason.
- State the next suggested skill: `$gym-git`, `$gym-fe-plan`, or another focused skill.
