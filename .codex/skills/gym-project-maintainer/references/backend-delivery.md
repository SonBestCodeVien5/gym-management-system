# Backend Delivery

Use this reference when the task changes backend code or backend feature memory.

## Read Path

1. Read `CHAT_CONTEXT/README.md`.
2. Read `CHAT_CONTEXT/backend_skills/README.md`.
3. Read `docs/api_contract.md` and `api_test.http` when HTTP behavior is touched.
4. Read only the feature memory files needed for the current phase.
5. Inspect only relevant models, repositories, services, handlers, route wiring, tests, and docs.

## Memory Layout

| Need | Path |
|---|---|
| Plan | `CHAT_CONTEXT/backend_skills/plans/<feature>.md` |
| Implementation note | `CHAT_CONTEXT/backend_skills/implementations/<feature>.md` |
| Review note | `CHAT_CONTEXT/backend_skills/reviews/<feature>.md` |
| Test note | `CHAT_CONTEXT/backend_skills/tests/<feature>.md` |
| Short chronology | `CHAT_CONTEXT/backend_skills/worklog.md` |

Use `_template.md` files in the implementation, review, and test folders when a new feature note
does not exist yet.

## Phase Expectations

| Phase | Minimum update |
|---|---|
| Plan | Goal, API shape, business rules, data/query changes, layer steps, tests, risks |
| Implement | Files changed, decisions, commands, limitations, review handoff |
| Review | Findings, fixes, remaining risks, test handoff |
| Test | Commands, manual/API checks, DB checks when relevant, final result |
| Complete | API contract, API samples, worklog, and chat snapshot when project state changes |

## Backend Rules

- Keep HTTP parsing and error mapping in handlers.
- Keep business rules and orchestration in services.
- Keep MongoDB access in repositories.
- Do not trust client input for money, status, roles, or server-computed counters.
- Prefer atomic updates or indexes for refund, enrollment, payment confirmation, duplicate check-in,
  and other double-submit risks.
- Use the current error/status mapping unless the contract is intentionally changed: invalid input
  -> `400`, not found -> `404`, conflicts/business rules -> `409`, unexpected server/storage errors
  -> `500`.
- Keep feature memory concise. Store decisions and verification results, not large diffs or copied
  source files.
