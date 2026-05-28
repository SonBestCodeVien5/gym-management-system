# Implementation - Frontend Readiness Mini-Cycle

## Status

- Status: implemented, pending review
- Feature: Frontend readiness mini-cycle
- Plan file: `CHAT_CONTEXT/backend_skills/plans/09_frontend_readiness.md`
- Started at: 2026-05-28
- Finished at: 2026-05-28

## Scope implemented

- [ ] Model changes
- [ ] Repository changes
- [x] Service changes
- [x] Handler changes
- [x] Route changes
- [x] Docs/API sample changes

## Files changed

- `.env.example` - added `CORS_ALLOWED_ORIGINS` sample for local FE dev.
- `cmd/server/main.go` - parses comma-separated CORS origins and passes them into app config.
- `internal/app/cors.go` - added allow-list CORS middleware with preflight handling.
- `internal/app/router.go` - applies CORS middleware and registers protected `GET /api/v1/auth/me`.
- `internal/handlers/auth_handler.go` - added current employee handler.
- `internal/handlers/auth_middleware_test.go` - updated auth service test stub for the new method.
- `internal/service/auth_service.go` - added `CurrentEmployee` to load the authenticated active employee.
- `internal/service/auth_service_test.go` - added unit coverage for current employee lookup.
- `internal/testutil/mongo.go` - configures test router with an allowed CORS origin.
- `internal/integration/integration_test.go` - added `/auth/me` and CORS preflight integration checks.
- `docs/api_contract.md` - documented `/auth/me` and CORS local-dev behavior.
- `docs/local_dev_guide.md` - documented CORS env and `/auth/me`.
- `api_test.http` - added current employee request.
- `CHAT_CONTEXT/README.md` - updated active backend snapshot and resume point.
- `CHAT_CONTEXT/backend_skills/worklog.md` - added implementation summary and review handoff.

## Key decisions

- CORS uses an explicit origin allow-list from `CORS_ALLOWED_ORIGINS`; empty config emits no CORS
  headers.
- Preflight requests are handled before auth so browsers can complete the CORS handshake without an
  access token.
- `/auth/me` is protected by the existing `AuthRequired` middleware and has no extra role guard.
- `/auth/me` reloads the employee from storage via the authenticated employee id and returns the same
  compact employee response shape used by login.
- Refresh-token rotation remains unchanged. FE can refresh tokens first, then call `/auth/me`.

## Implementation notes

### Models

- No model changes.

### Repository

- No repository changes.

### Service

- `AuthService.CurrentEmployee` maps missing employees to `ErrInvalidToken` and inactive employees
  to `ErrInactiveEmployee`.
- Unexpected repository errors are still returned for handler-level `500` mapping.

### Handler

- `AuthHandler.Me` reads `AuthEmployeeIDKey` from Gin context, rejects missing/malformed context as
  unauthorized, and maps auth errors to the shared error envelope.

### Routes

- `GET /api/v1/auth/me` is registered on the protected API group before role-specific groups.
- CORS middleware is registered on the Gin engine before route registration.

### Docs/API samples

- API contract, local-dev guide, `.env.example`, and REST Client sample now include the new frontend
  readiness behavior.

## Commands run

```bash
gofmt -w cmd/server/main.go internal/app/router.go internal/app/cors.go internal/handlers/auth_handler.go internal/handlers/auth_middleware_test.go internal/service/auth_service.go internal/service/auth_service_test.go internal/integration/integration_test.go internal/testutil/mongo.go
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
env GOCACHE=/tmp/gocache go test ./internal/integration -count=1
```

`go build ./...` passed with the existing Go module stat-cache read-only warning and exit code `0`.

## Known limitations

- CORS is configured for bearer-token browser requests only. Cookie auth and credentials are still
  out of scope.
- No seed/demo-data command was added in this mini-cycle.
- CI automation remains a follow-up.

## Handoff to review

- Check CORS behavior, especially global preflight handling and no wildcard origin.
- Check `/auth/me` route placement and error mapping.
- Check docs/sample alignment with the new protected auth endpoint.
