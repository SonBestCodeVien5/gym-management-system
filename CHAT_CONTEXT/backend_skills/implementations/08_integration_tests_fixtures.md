# Implementation - integration tests fixtures

## Status

- Status: implemented
- Feature: integration tests fixtures
- Plan file: `CHAT_CONTEXT/backend_skills/plans/08_integration_tests_fixtures.md`
- Started at: 2026-05-28
- Finished at: 2026-05-28

## Scope implemented

- [ ] Model changes
- [ ] Repository changes
- [ ] Service changes
- [ ] Handler changes
- [x] Route/dependency wiring changes
- [x] Test utility changes
- [x] Integration tests
- [x] Docs/context changes

## Files changed

- `cmd/server/main.go`
- `internal/app/router.go`
- `internal/testutil/mongo.go`
- `internal/testutil/http.go`
- `internal/testutil/fixtures.go`
- `internal/integration/integration_test.go`
- `docs/code_reading_guide.md`
- `docs/local_dev_guide.md`
- `README.md`
- `CHAT_CONTEXT/backend_skills/worklog.md`

## Key decisions

- Extracted production router/dependency wiring to `internal/app.NewRouter` and
  `internal/app.RegisterRoutes`.
- Kept `cmd/server/main.go` responsible for env loading, MongoDB connection, index bootstrap,
  app config, and `r.Run`.
- Integration tests use the same `internal/app` route registration as production.
- MongoDB integration tests create an isolated `gym_test_<id>` database and drop only that database
  in cleanup.
- Test harness runs `pkg/database.EnsureIndexes` on the test DB before building the router.
- If MongoDB is not reachable, integration tests skip instead of failing the whole unit-test suite.

## Implementation notes

- `internal/testutil.NewTestApp` builds the router with deterministic test auth secrets and a
  bootstrapped admin employee.
- HTTP helpers assert status, shared error response shape, and response `data` payloads.
- Fixture helpers create branch, course, member, pending subscription, and active subscription
  through real HTTP endpoints.
- Integration coverage currently includes:
  - `/ping` smoke
  - missing-token `401`
  - admin access to employee management
  - non-admin `403` on employee management
  - login, refresh, logout, and revoked refresh-token rejection
  - branch/course/member/subscription create and member activation
  - duplicate branch code, duplicate member CCID, and duplicate refund `409`
  - refund audit document count
  - reported missed attendance, valid makeup, and reused makeup `409`
  - branch nearby success and invalid query `400`

## Commands run

- `gofmt -w cmd/server/main.go internal/app/router.go`
- `gofmt -w internal/testutil/mongo.go internal/testutil/http.go internal/testutil/fixtures.go internal/integration/integration_test.go`
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1` - pass.

## Known limitations

- No GitHub Actions/CI service setup in this cycle.
- Integration tests depend on a reachable MongoDB instance and skip when it is unavailable.
- Testcontainers were intentionally not added.
- Session integration coverage is not included yet.
- Refund and attendance multi-write flows are verified by final state only; this does not add Mongo
  transactions.

## Handoff to review

Review route parity between `internal/app/router.go` and the previous `cmd/server/main.go` wiring,
test DB cleanup guard, integration test assertions, and docs alignment.

Next file for `$gym-review`:
`CHAT_CONTEXT/backend_skills/implementations/08_integration_tests_fixtures.md`
