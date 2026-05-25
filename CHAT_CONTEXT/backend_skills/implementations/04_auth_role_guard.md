# Implementation — auth role guard

## Status

- Status: implemented
- Feature: auth role guard
- Plan file: `CHAT_CONTEXT/backend_skills/plans/04_auth_role_guard.md`
- Started at: 2026-05-25
- Finished at: 2026-05-25

## Scope implemented

- [x] Model changes
- [x] Repository changes
- [x] Service changes
- [x] Handler changes
- [x] Route changes
- [x] Docs/API sample changes

## Files changed

- `.env.example`
- `api_test.http`
- `cmd/server/main.go`
- `docs/api_contract.md`
- `internal/handlers/auth_handler.go`
- `internal/handlers/auth_middleware.go`
- `internal/handlers/auth_middleware_test.go`
- `internal/models/employee.go`
- `internal/models/refresh_token.go`
- `internal/repository/employee_repo.go`
- `internal/repository/refresh_token_repo.go`
- `internal/service/auth_service.go`
- `internal/service/auth_service_test.go`

## Key decisions

- Implemented JWT-compatible HS256 tokens with Go stdlib crypto instead of adding a new JWT
  dependency.
- Added random `jti` to each token so refresh rotation cannot reissue the same token when claims
  are otherwise identical.
- Stored refresh tokens only as SHA-256 hashes.
- Stored normalized login identity in `employees.normalized_email`; repository creates a unique
  sparse index for it and a unique sparse index for `employee_id`.
- `POST /api/v1/auth/logout` remains public and uses the refresh token body for revoke.
- Business routes are protected by access-token middleware plus role guard groups.

## Implementation notes

- Extended `Employee` with `normalized_email`, `password_hash`, `status`, and timestamps while
  preserving multi-role and multi-branch fields.
- Added env-based bootstrap admin; startup only inserts when the normalized email does not already
  exist and never logs the password or hash.
- Added refresh-token collection with unique `token_hash` index and active-token revoke lookup.
- Added `AuthService` login, refresh rotation, logout, access-token validation, and bootstrap
  methods.
- Added `AuthHandler` for login/refresh/logout and auth middleware for `401`/`403` behavior.
- Protected course/branch/member/subscription/attendance/session routes according to the cycle 04
  role matrix.
- Updated API contract and REST samples for auth, protected calls, missing-token, and forbidden-role
  checks.

## Commands run

- `gofmt -w cmd/server/main.go internal/models/employee.go internal/models/refresh_token.go internal/repository/employee_repo.go internal/repository/refresh_token_repo.go internal/service/auth_service.go internal/handlers/auth_handler.go internal/handlers/auth_middleware.go`
- `env GOCACHE=/tmp/gocache go build ./...` — pass; emitted a sandbox-related stat-cache warning
  for the default module cache but exited `0`.
- `gofmt -w internal/service/auth_service_test.go`
- `env GOCACHE=/tmp/gocache go test ./...` — pass.
- Test phase added `internal/handlers/auth_middleware_test.go`; `env GOCACHE=/tmp/gocache go build ./...` and `env GOCACHE=/tmp/gocache go test ./...` still pass.

## Known limitations

- No branch-scope authorization yet; role guard only checks role membership.
- Employee management endpoints are still planned for cycle 05, so non-bootstrap staff accounts
  still need direct DB/seed setup until that cycle lands.
- Refresh-token TTL cleanup index is not added; service validates expiry and index cleanup remains
  suitable for cycle 07 data-integrity work.
- No manual API/DB verification has been run in this implementation phase.

## Handoff to review

- Review JWT parsing/signing, refresh rotation/revoke edge cases, role guard route coverage, and
  startup failure behavior when JWT/bootstrap env values are missing or invalid.
- Verify route order still preserves `/branches/nearby` before `/branches/:id`.
