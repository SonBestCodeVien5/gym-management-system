---
name: gym-implement
description: Implement planned backend features for this gym management system. Use when Codex is asked to code from an existing backend plan, continue implementation notes, or respond to an explicit `$gym-implement` or `/gym-implement` request in this repository.
---

# Gym Implement

## Read First

1. Read `docs/README.md`.
2. Read `CHAT_CONTEXT/README.md`.
3. Read `CHAT_CONTEXT/backend_skills/README.md`.
4. Read `../gym-project-maintainer/references/backend-delivery.md`.
5. Read `CHAT_CONTEXT/backend_skills/plans/<feature>.md`.
6. Read `CHAT_CONTEXT/backend_skills/implementations/<feature>.md` if it exists; otherwise use
   `CHAT_CONTEXT/backend_skills/implementations/_template.md`.
7. Read `docs/api_contract.md` and `api_test.http` when HTTP behavior changes.
8. Inspect only source files required by the plan.

## Focus

- Implement by existing repo layers: model, repository, service, handler, route wiring, docs/tests.
- Keep business rules in services, MongoDB work in repositories, and HTTP parsing/error mapping in
  handlers.
- Reject client-controlled money, status, roles, and computed counters.
- Use atomic updates or indexes where double-submit or race risk matters.

## Output Rules

- Update `CHAT_CONTEXT/backend_skills/implementations/<feature>.md` with changed files, decisions,
  commands, limitations, and review handoff.
- Keep `docs/api_contract.md` and `api_test.http` aligned when the implementation changes HTTP
  behavior.
- Run `gofmt` on changed Go files.
- Run `go build ./...`; run `go test ./...` when feasible and record skipped reasons when not.
- Do not write the review result for `$gym-review`.
