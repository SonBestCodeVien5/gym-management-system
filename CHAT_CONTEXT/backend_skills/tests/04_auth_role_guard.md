# Test — auth role guard

## Status

- Status: completed
- Feature: auth role guard
- Plan file: `CHAT_CONTEXT/backend_skills/plans/04_auth_role_guard.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/04_auth_role_guard.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/04_auth_role_guard.md`
- Tested at: 2026-05-25

## Commands

- `env GOCACHE=/tmp/gocache go build ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `env GOCACHE=/tmp/gocache PORT=18080 MONGODB_URI=... JWT_ACCESS_SECRET=... JWT_REFRESH_SECRET=... BOOTSTRAP_ADMIN_* go run ./cmd/server` - pass after running with local-network escalation; sandboxed run could not connect to `127.0.0.1:27017`.

## Manual API tests

### Happy path

- [x] `GET /ping` returned `200`.
- [x] Admin login returned `200` with access and refresh tokens present.
- [x] `GET /api/v1/courses` with admin access token returned `200`.
- [x] Refresh with current refresh token returned `200` and a replacement refresh token.
- [x] Logout with replacement refresh token returned `200`.
- [x] Repeated logout with the same refresh token returned `200`.

### Invalid input

- [x] `GET /api/v1/courses` without access token returned `401`.
- [x] Admin login with wrong password returned `401`.
- [x] Refresh with missing `refresh_token` returned `400`.
- [x] Login for an inactive seeded employee returned `401`.

### Not found

- [x] Reusing the old rotated refresh token returned `401`, covering revoked/unknown refresh-token behavior.

### Conflict/business rule

- [x] Seeded receptionist login returned `200`, then `POST /api/v1/courses` with receptionist token returned `403`.
- [x] Added automated middleware tests for missing token `401`, invalid token `401`, allowed role `200`, forbidden role `403`, and unexpected auth-service error `500`.

## DB state verification

- [x] Bootstrap admin exists with role `admin`.
- [x] Admin has a bcrypt-looking `password_hash`.
- [x] Refresh token documents have `token_hash`.
- [x] Refresh token documents do not contain raw `refresh_token` or `token` fields.
- [x] Refresh/logout test left admin refresh tokens revoked with `revoked_at`.
- [x] Temporary seeded employees `restricted@gym.test` and `inactive@gym.test` were deleted after test.
- [x] One orphan refresh token from the deleted restricted employee was deleted after test cleanup.

## Issues found

- Review test gap was addressed by adding `internal/handlers/auth_middleware_test.go`.
- No runtime/API blocker found.

## Final result

- Result: pass
- Ready to update docs/context: yes; ready for `$gym-complete`.
