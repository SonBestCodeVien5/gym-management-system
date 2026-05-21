---
name: gym-review
description: Review backend feature changes for this gym management system. Use when Codex is asked for a code review, backend feature review, phase review note, or an explicit `$gym-review` or `/gym-review` request in this repository.
---

# Gym Review

## Read First

1. Read `CHAT_CONTEXT/README.md`.
2. Read `CHAT_CONTEXT/backend_skills/README.md`.
3. Read `../gym-project-maintainer/references/backend-delivery.md`.
4. Read the target plan and implementation note for `<feature>`.
5. Read `CHAT_CONTEXT/backend_skills/reviews/<feature>.md` if it exists; otherwise use
   `CHAT_CONTEXT/backend_skills/reviews/_template.md`.
6. Inspect the changed source files, relevant tests, `docs/api_contract.md`, and `api_test.http`.

## Focus

- Lead with bugs, regressions, data-integrity risks, security gaps, wrong layer ownership, contract
  drift, and missing tests.
- Check handler/service/repository boundaries, error mapping, atomic updates, route order, model
  tags, and client-controlled fields.
- Verify the implementation matches the plan or call out the mismatch.
- Keep review evidence grounded in file and line references.

## Output Rules

- Update `CHAT_CONTEXT/backend_skills/reviews/<feature>.md` with findings, fixes if applied,
  remaining risks, and test handoff.
- Present findings first when reporting to the user.
- Do not turn review into a broad refactor unless the user asks for fixes.
- Do not claim verification that was not run.
