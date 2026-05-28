---
name: gym-fe-implement
description: Implement planned React/Vite frontend features for this gym management system. Use when Codex is asked to code from an existing frontend plan, continue FE implementation notes, build UI/components/API client work, or respond to an explicit `$gym-fe-implement` or `/gym-fe-implement` request in this repository.
---

# Gym FE Implement

## Read First

1. Read `docs/README.md`.
2. Read `CHAT_CONTEXT/README.md`.
3. Read `CHAT_CONTEXT/frontend_skills/README.md`.
4. Read `CHAT_CONTEXT/frontend_skills/plans/<feature>.md`.
5. Read `CHAT_CONTEXT/frontend_skills/implementations/<feature>.md` if it exists; otherwise use
   `CHAT_CONTEXT/frontend_skills/implementations/_template.md`.
6. Read `docs/api_contract.md` and `api_test.http` when API behavior or FE API calls are involved.
7. Inspect only source files required by the plan.

## Focus

- Implement with existing React/Vite structure under `frontend/`.
- Keep UI components, state, API client code, and styling scoped to the planned feature.
- Use `VITE_*` env variables for browser-visible config; do not hardcode temporary backend ports.
- Do not store secrets in frontend code, env samples, localStorage examples, docs, or screenshots.
- Build responsive, usable states: loading, error, empty, disabled, and success where relevant.
- Preserve the app's Iron Forge design language unless the plan intentionally changes it.

## Output Rules

- Update `CHAT_CONTEXT/frontend_skills/implementations/<feature>.md` with changed files, decisions,
  commands, limitations, and review handoff.
- Keep docs/API contract notes aligned if FE work exposes or depends on backend behavior.
- Run `npm run build` when feasible and record skipped reasons when not.
- Do not write the review result for `$gym-fe-review`.
