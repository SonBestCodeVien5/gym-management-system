---
name: gym-fe-test
description: Verify React/Vite frontend features for this gym management system. Use when Codex is asked to test a frontend cycle, run build/manual browser checks/API flow checks, update a frontend test note, or respond to an explicit `$gym-fe-test` or `/gym-fe-test` request in this repository.
---

# Gym FE Test

## Read First

1. Read `CHAT_CONTEXT/README.md`.
2. Read `CHAT_CONTEXT/frontend_skills/README.md`.
3. Read the target plan, implementation note, and review note for `<feature>`.
4. Read `CHAT_CONTEXT/frontend_skills/tests/<feature>.md` if it exists; otherwise use
   `CHAT_CONTEXT/frontend_skills/tests/_template.md`.
5. Read relevant frontend source, CSS, API contract, and docs needed to understand expected behavior.

## Focus

- Verify `npm run build`, route/page rendering, responsive layout, interaction behavior, forms, API
  success/error behavior, and auth/session flows that matter for the feature.
- Prefer real browser/dev-server evidence when available. Use screenshots or DOM checks when visual
  correctness matters.
- Record what was actually run, what was skipped, and why.
- Do not invent API/manual results when local backend or data prerequisites are unavailable.

## Output Rules

- Update `CHAT_CONTEXT/frontend_skills/tests/<feature>.md`.
- Run `npm run build` when feasible.
- Start the Vite dev server for manual checks when the app needs a browser/dev-server environment.
- State whether the feature is ready for `$gym-fe-complete`.
