---
name: gym-fe-plan
description: Plan frontend feature work for this gym management system. Use when Codex is asked to design a React/Vite frontend cycle, define pages/routes/components/state/API integration, prepare UI implementation steps before coding, or respond to an explicit `$gym-fe-plan` or `/gym-fe-plan` request in this repository.
---

# Gym FE Plan

## Read First

1. Read `docs/README.md`.
2. Read `CHAT_CONTEXT/README.md`.
3. Read `CHAT_CONTEXT/frontend_skills/README.md`.
4. Read `docs/api_contract.md` when the feature calls backend APIs.
5. Read `frontend/package.json`, `frontend/vite.config.js`, `frontend/src/App.jsx`, and
   `frontend/src/index.css` when planning inside the current React app.
6. Read the target plan in `CHAT_CONTEXT/frontend_skills/plans/<feature>.md` if it exists.
7. Inspect only relevant frontend source, docs, and API samples.

## Focus

- Establish current UI, routes, app structure, API contract, and design baseline before proposing
  changes.
- Define pages/routes, components, state shape, API calls, env variables, loading/empty/error states,
  validation, responsive behavior, accessibility basics, docs, tests, and risks.
- Keep backend-contract needs separate from frontend-only decisions.
- Prefer existing app style, design tokens, and local component patterns.

## Output Rules

- Update or create `CHAT_CONTEXT/frontend_skills/plans/<feature>.md`.
- Include status, goal, screens/routes, component plan, state/API plan, UX states, responsive/accessibility
  notes, docs/test plan, and risks.
- Update `CHAT_CONTEXT/frontend_skills/worklog.md` only when roadmap status or next action changes.
- Do not implement production code in this phase unless the user explicitly changes scope.
- End with the file to read next for `$gym-fe-implement`.
