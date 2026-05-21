---
name: gym-test
description: Verify backend features for this gym management system. Use when Codex is asked to test a backend cycle, run build/test/manual API verification, update a feature test note, or respond to an explicit `$gym-test` or `/gym-test` request in this repository.
---

# Gym Test

## Read First

1. Read `CHAT_CONTEXT/README.md`.
2. Read `CHAT_CONTEXT/backend_skills/README.md`.
3. Read `../gym-project-maintainer/references/backend-delivery.md`.
4. Read the target plan, implementation note, and review note for `<feature>`.
5. Read `CHAT_CONTEXT/backend_skills/tests/<feature>.md` if it exists; otherwise use
   `CHAT_CONTEXT/backend_skills/tests/_template.md`.
6. Read `docs/api_contract.md`, `api_test.http`, relevant automated tests, and source needed to
   understand expected behavior.

## Focus

- Verify build, automated tests, HTTP behavior, business-rule failures, and DB state effects that
  matter for the feature.
- Cover happy path, invalid input, not found, conflict/business-rule behavior, and race-sensitive
  cases when practical.
- Record what was actually run, what was skipped, and why.
- Prefer real command/API evidence over assumptions.

## Output Rules

- Update `CHAT_CONTEXT/backend_skills/tests/<feature>.md`.
- Run `go build ./...` and `go test ./...` when feasible.
- Use manual API checks only when the local environment and data prerequisites are available; record
  skipped prerequisites instead of inventing results.
- State whether the feature is ready for `$gym-complete`.
