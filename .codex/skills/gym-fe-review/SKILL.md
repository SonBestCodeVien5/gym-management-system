---
name: gym-fe-review
description: Review React/Vite frontend feature changes for this gym management system. Use when Codex is asked for a frontend code review, UI/UX review, API integration review, phase review note, or an explicit `$gym-fe-review` or `/gym-fe-review` request in this repository.
---

# Gym FE Review

## Read First

1. Read `CHAT_CONTEXT/README.md`.
2. Read `CHAT_CONTEXT/frontend_skills/README.md`.
3. Read the target plan and implementation note for `<feature>`.
4. Read `CHAT_CONTEXT/frontend_skills/reviews/<feature>.md` if it exists; otherwise use
   `CHAT_CONTEXT/frontend_skills/reviews/_template.md`.
5. Inspect changed frontend files, relevant CSS, docs, `docs/api_contract.md`, and API samples when
   integration is involved.

## Focus

- Lead with bugs, UX regressions, responsive breakage, accessibility gaps, auth/session risks,
  hardcoded endpoints, contract drift, and missing states/tests.
- Check component boundaries, state ownership, event handlers, form validation, error handling,
  loading/empty states, and API response assumptions.
- Verify implementation matches the plan or call out the mismatch.
- Keep review evidence grounded in file and line references.

## Browser Checks

- When a change affects UI, routing, responsive behavior, auth flow, or form interaction, use MCP
  Playwright when available for at least one real browser pass.
- Prefer DOM/accessibility snapshots for interaction evidence; use screenshots when layout, visual
  overlap, or responsive behavior matters.
- Cover the smallest route set that proves the reviewed behavior: happy path, empty/error/forbidden
  states when applicable, and one narrow viewport for responsive-sensitive work.
- If Playwright or a dev server/backend prerequisite is unavailable, record the skip reason and keep
  the risk in the review handoff.

## Output Rules

- Update `CHAT_CONTEXT/frontend_skills/reviews/<feature>.md` with findings, fixes if applied,
  remaining risks, and test handoff.
- Present findings first when reporting to the user.
- Do not turn review into a broad redesign unless the user asks for fixes.
- Do not claim verification that was not run.
