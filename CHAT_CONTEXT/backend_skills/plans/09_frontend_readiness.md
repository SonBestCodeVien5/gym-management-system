# Feature Plan - Frontend Readiness Mini-Cycle

## Status
- Planned: yes
- Implemented: no
- Reviewed: no
- Tested: no
- Docs updated: no

## Goal
Prepare the backend for the first frontend integration pass without expanding the business domain.
This mini-cycle should make browser calls from a local FE app predictable and give FE a stable way
to restore the current logged-in employee after page reload.

## Current context
- Current public auth endpoints are `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`, and
  `POST /api/v1/auth/logout`.
- Protected business routes already use `handlers.AuthRequired` and role guards in
  `internal/app/router.go`.
- `AuthRequired` validates the access token, reloads the employee, rejects inactive or missing
  employees, and stores `auth_employee_id` plus `auth_roles` in Gin context.
- Login returns a compact employee shape through `service.AuthEmployeeResponse`, but refresh does
  not return employee data.
- There is no CORS middleware yet, so a browser FE on Vite or another dev origin will be blocked
  even when the API works through REST Client or curl.

## API shape

### Add current employee endpoint

`GET /api/v1/auth/me`

Auth:
- Requires `Authorization: Bearer <access_token>`.
- Any active authenticated employee may call it.
- No role guard beyond `AuthRequired`.

Response `200`:

```json
{
  "message": "current employee fetched successfully",
  "data": {
    "id": "69f20c000c4cd4cdf5768500",
    "employee_id": "ADMIN001",
    "email": "admin@gym.test",
    "full_name": "Gym Admin",
    "role": ["admin"],
    "branch_id": []
  }
}
```

Status codes:
- `200`: access token is valid and maps to an active employee.
- `401`: missing, malformed, expired, invalid, deleted-employee, or inactive-employee token.
- `500`: unexpected storage/internal failure.

### Add CORS support

Configuration:
- Add `CORS_ALLOWED_ORIGINS` as a comma-separated env value.
- Recommended local value:
  `http://localhost:5173,http://127.0.0.1:5173`
- Empty value means no CORS headers are emitted.

Behavior:
- Apply CORS globally before route handling so browser preflight requests do not hit auth guards.
- Reflect the request `Origin` only when it exactly matches the configured allow-list.
- Allow methods: `GET`, `POST`, `PATCH`, `DELETE`, `OPTIONS`.
- Allow headers: `Authorization`, `Content-Type`.
- Keep credentials disabled for now because current auth uses bearer tokens, not cookies.
- Do not use wildcard `*` with credentials.

Preflight:
- `OPTIONS` from an allowed origin should return a successful empty response with CORS headers.
- `OPTIONS` from a disallowed origin should not receive `Access-Control-Allow-Origin`.

## Business rules
- `/auth/me` must not trust any client-sent employee payload or role. It must load the employee from
  storage through the existing authenticated employee id.
- Inactive employees and removed employees should behave as invalid access-token cases because the
  token no longer represents a usable staff session.
- Refresh-token rotation remains unchanged. FE can call `/auth/refresh` for a new token pair, then
  `/auth/me` to reload employee context.
- No member, subscription, attendance, course, branch, session, employee-management, refund, or
  pricing behavior changes in this cycle.

## Data model and indexes
- No collection schema changes.
- No new MongoDB indexes.
- No data migration.
- No seed command in this mini-cycle. Keep demo/bootstrap guidance in docs and REST samples only.

## Layer plan

### Config/startup
- Extend `app.Config` with CORS allow-list configuration.
- Parse `CORS_ALLOWED_ORIGINS` in `cmd/server/main.go` into trimmed, non-empty origins.
- Add `.env.example` and local-dev docs entries for FE dev origin.

### App/router
- Register CORS middleware on the Gin engine before `RegisterRoutes`.
- Add `protected.GET("/auth/me", h.Auth.Me)` after creating the protected API group and before
  role-specific route groups.
- Keep public auth route behavior unchanged.

### Handler
- Add `AuthHandler.Me`.
- Read `handlers.AuthEmployeeIDKey` from Gin context.
- If the key is missing or malformed, return shared `401 UNAUTHORIZED`.
- Call the auth service to fetch the current employee response.
- Return the existing success envelope: `{"message":"current employee fetched successfully","data":...}`.

### Service
- Extend `service.AuthService` with a `CurrentEmployee(ctx, employeeID string)` style method.
- Reuse the same compact employee response helper used by login.
- Convert repository not-found to `ErrInvalidToken` and inactive status to `ErrInactiveEmployee`.
- Preserve sanitized error mapping. Do not leak storage or token internals.

### Tests
- Add or extend integration tests for:
  - Login then `GET /api/v1/auth/me` returns the same compact employee identity.
  - Missing access token on `/auth/me` returns `401` with shared error envelope.
  - Allowed-origin preflight returns CORS headers without requiring auth.
  - Disallowed origin does not receive `Access-Control-Allow-Origin`.
- Keep `go test ./...` behavior developer-friendly. Integration tests may continue to skip when
  MongoDB is not reachable.

## Docs and samples
- Update `docs/api_contract.md`:
  - Add `GET /api/v1/auth/me` to the Auth table.
  - Add endpoint details and note that browser FE should send bearer tokens in `Authorization`.
  - Add a short CORS local-dev note.
- Update `api_test.http` with a current-user request after login.
- Update `.env.example` with `CORS_ALLOWED_ORIGINS`.
- Update `docs/local_dev_guide.md` route list and local FE setup guidance.
- Update backend phase notes during implement/review/test/complete phases as usual.

## Risks and decisions
- CORS is browser-only protection. It does not replace auth, role guards, or server-side validation.
- A permissive wildcard origin would make local testing easy but weakens production posture, so the
  plan uses an explicit allow-list.
- Preflight must run before auth. Otherwise browsers will fail before the actual authenticated API
  request is sent.
- `/auth/me` intentionally returns the compact login employee shape, not the full admin employee
  management shape. FE can use employee-management endpoints later when it needs full staff detail.
- A proper seed/demo-data command is useful later, but adding it now would create a broader data
  lifecycle concern. This mini-cycle keeps that as follow-up.

## Out of scope
- CI automation.
- Session-specific integration coverage expansion.
- Full seed/demo-data command.
- Cookie-based auth, CSRF, OAuth, password reset, or refresh-token cookies.
- Frontend implementation.

## Verification target
- `git diff --check`
- `env GOCACHE=/tmp/gocache go build ./...`
- `env GOCACHE=/tmp/gocache go test ./...`
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1` when MongoDB is reachable

## Next action
Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/09_frontend_readiness.md`.
