# Code Review — auth role guard

## Status

- Status: completed
- Feature: auth role guard
- Plan file: `CHAT_CONTEXT/backend_skills/plans/04_auth_role_guard.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/04_auth_role_guard.md`
- Reviewed at: 2026-05-25

## Review summary

- Result: pass with test gaps and residual risks
- Build status: pass (`env GOCACHE=/tmp/gocache go build ./...`)
- Test status: pass (`env GOCACHE=/tmp/gocache go test ./...`)

## Checklist

- [x] Code compiles.
- [x] Handler only handles HTTP parse/response.
- [x] Service owns business rules.
- [x] Repository only handles DB.
- [x] Model tags match API/DB contract.
- [x] Errors map to correct HTTP status.
- [x] Atomic updates used where needed.
- [x] Routes have correct order.
- [x] Docs/API samples match behavior.

## Passed

- Auth endpoints are registered before the protected route group.
- `/branches/nearby` remains registered before `/branches/:id`.
- Role guard matrix in route wiring matches the cycle 04 plan.
- Access-token validation reloads employee state and does not trust client-sent roles.
- Refresh token storage uses hashed token values and conditional revoke.
- Login does not return password hash or normalized-email internals.
- API contract and REST samples document auth, protected routes, `401`, and `403`.

## Issues found

- Medium: Route/middleware behavior was not covered by automated tests at review time. The implementation had
  service-level tests for login, refresh rotation, reused refresh token rejection, logout
  idempotency, and wrong password, but there are no handler/router tests proving missing access
  token `401`, role guard `403`, public auth route access, or protected business route coverage.
  Evidence at review time: route wiring lives in `cmd/server/main.go`; tests lived only in
  `internal/service/auth_service_test.go`. Test phase added `internal/handlers/auth_middleware_test.go`
  to cover middleware/role guard status behavior.
- Low: Refresh rotation revokes the presented refresh token before inserting the replacement token.
  If replacement persistence fails after revoke, the client receives an error and the old refresh
  token is already invalidated. This does not create a privilege escalation, but it is an
  availability/session-continuity risk.

## Fixes applied during review

- None.

## Remaining risks

- Manual API/DB verification was still pending at review time; test phase completed it successfully.
- Refresh-token TTL cleanup remains deferred to the data-integrity/index cycle.
- Branch-scope authorization is intentionally out of scope for cycle 04.
- JWT secrets are required but strength/length is not enforced beyond non-empty config.
- Employee management is not implemented yet, so future non-bootstrap staff account management still
  needs the next backend cycle.

## Handoff to test

- Run manual API checks for bootstrap admin, login, authorized request, missing-token `401`,
  forbidden-role `403`, refresh rotation, reused old refresh token rejection, and logout
  idempotency.
- Verify MongoDB documents for `employees` and `refresh_tokens`, including password hash presence,
  absence of raw refresh tokens, and `revoked_at` after refresh/logout.
- Consider adding handler/router tests before marking cycle 04 fully tested.
