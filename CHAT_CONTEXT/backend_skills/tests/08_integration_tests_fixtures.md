# Test - integration tests fixtures

## Status

- Status: tested
- Feature: integration tests fixtures
- Plan file: `CHAT_CONTEXT/backend_skills/plans/08_integration_tests_fixtures.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/08_integration_tests_fixtures.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/08_integration_tests_fixtures.md`
- Tested at: 2026-05-28

## Commands

- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass in sandbox; integration package used cached
  result in that run.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1 -v` - sandbox run passed by
  skipping all integration tests because socket access to `localhost:27017` was not permitted.
- `go test ./internal/integration -count=1 -v` outside sandbox - pass with real MongoDB on
  `localhost:27017`.
- `go test ./...` outside sandbox - pass with real MongoDB integration tests.
- `git diff --check` - pass.
- `docker ps --filter name=gym_mongodb --format ...` - MongoDB container running and exposing
  `27017`.
- `docker exec gym_mongodb mongosh ... listDatabases` - pass; no leftover `gym_test_*` databases
  after test cleanup.

## Automated HTTP Tests

The integration suite uses `httptest` with the real router from `internal/app`, real services, real
repositories, and a real MongoDB test database.

### Happy path

- [x] `GET /ping` returns `200`.
- [x] Admin login returns access and refresh tokens.
- [x] Admin can access employee management.
- [x] Refresh token rotates successfully.
- [x] Logout revokes refresh token.
- [x] Create branch, course, member, pending subscription, and activate subscription.
- [x] Get subscription and list member subscriptions.
- [x] Report missed attendance and create a valid makeup attendance.
- [x] Branch nearby query returns non-empty results with `distance_meters`.

### Invalid input

- [x] Branch nearby invalid longitude returns `400 INVALID_INPUT`.

### Not found

- [ ] Not covered by Cycle 08 integration tests. Existing service tests cover some not-found
  mappings; broader HTTP not-found integration coverage remains follow-up.

### Conflict/business rule

- [x] Non-admin token on employee management returns `403 FORBIDDEN`.
- [x] Missing access token returns `401 UNAUTHORIZED`.
- [x] Reusing logged-out refresh token returns `401 UNAUTHORIZED`.
- [x] Duplicate branch code returns `409 CONFLICT`.
- [x] Duplicate member CCID returns `409 CONFLICT`.
- [x] Duplicate refund returns `409 CONFLICT`.
- [x] Reusing the same makeup reference returns `409 CONFLICT`.

## DB State Verification

- [x] Integration tests create isolated databases named `gym_test_<ObjectID>`.
- [x] Test setup runs `pkg/database.EnsureIndexes` before seeding fixtures.
- [x] Duplicate refund test verifies exactly one refund audit document remains in the test DB.
- [x] Post-test MongoDB inspection found no leftover `gym_test_*` databases.

## Manual API Tests

- Manual server/curl checks were not run for this phase because the feature is specifically an
  automated integration-test harness and the `httptest` suite exercises the HTTP surface without a
  live server port.

## Issues found

- None.

## Final result

- Result: pass.
- Ready to update docs/context: yes.
- Ready for `$gym-complete`: yes.
