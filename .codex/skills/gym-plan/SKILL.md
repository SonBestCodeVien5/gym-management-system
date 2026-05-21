---
name: gym-plan
description: Plan backend feature work for this gym management system. Use when Codex is asked to design a backend cycle, define API/business/data/layer steps, prepare a feature plan before coding, or respond to an explicit `$gym-plan` or `/gym-plan` request in this repository.
---

# Gym Plan

## Read First

1. Read `docs/README.md`.
2. Read `CHAT_CONTEXT/README.md`.
3. Read `CHAT_CONTEXT/backend_skills/README.md`.
4. Read `../gym-project-maintainer/references/backend-delivery.md`.
5. Read `docs/api_contract.md`.
6. Read the target plan in `CHAT_CONTEXT/backend_skills/plans/<feature>.md` if it exists.
7. Inspect only relevant models, repositories, services, handlers, route wiring, tests, and docs.

## Focus

- Establish the current code and API baseline before proposing changes.
- Define API shape, business rules, data model/query/index changes, layer responsibilities, docs,
  verification, and risks.
- Surface security, atomicity, data-integrity, and route-order concerns while they are still design
  decisions.
- Keep current behavior and planned behavior distinct.

## Output Rules

- Update or create `CHAT_CONTEXT/backend_skills/plans/<feature>.md`.
- Include status, goal, API contract, business rules, data changes, layer plan, docs/test plan, and
  risks for a new plan.
- Update `CHAT_CONTEXT/backend_skills/worklog.md` only when roadmap status or next action changes.
- Do not implement production code in this phase unless the user explicitly changes scope.
- End with the file to read next for `$gym-implement`.
