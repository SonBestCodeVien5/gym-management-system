# Cycle 08 - Integration Tests & Fixtures

## Status

- Status: planned
- Priority: medium
- Depends on: Cycle 07 indexes/data integrity complete
- Planned at: 2026-05-28

## Goal

Add a reusable integration-test harness and fixture layer so important backend behavior can be
verified with Go tests instead of relying only on `api_test.http` and manual local checks.

This cycle should not add new API endpoints. It should make the existing HTTP surface testable
through `httptest` with real repositories, real services, real handlers, and a real MongoDB test
database.

## Current Baseline

- Existing automated tests are mostly service-level unit tests with stub repositories:
  - `internal/service/subscription_service_test.go`
  - `internal/service/auth_service_test.go`
  - `internal/service/employee_service_test.go`
  - `internal/handlers/auth_middleware_test.go`
- Manual API coverage lives in `api_test.http`.
- Route wiring currently lives inline in `cmd/server/main.go`.
- MongoDB index bootstrap lives in `pkg/database.EnsureIndexes`.
- Local MongoDB is provided by `docker-compose.yml` on `localhost:27017`.

## API Contract

No public API contract changes are planned.

Integration tests should assert existing contract behavior:

- Protected business routes require `Authorization: Bearer <access_token>`.
- Success responses keep the current `message`/`data` shape.
- Error responses use:

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "branch code already exists",
    "details": {}
  }
}
```

- Important status mappings remain:
  - invalid input -> `400`
  - unauthorized -> `401`
  - forbidden -> `403`
  - missing resource/reference -> `404`
  - business conflict or duplicate unique value -> `409`
  - unexpected storage/internal failure -> `500`

## Test Harness Plan

### Router Construction

Prefer extracting route/dependency construction out of `cmd/server/main.go` into a reusable package,
for example:

- `internal/app/router.go`
- `internal/app/dependencies.go`

Planned responsibilities:

- Build repositories from a supplied `*mongo.Database`.
- Build services from repositories.
- Build handlers from services.
- Register all current routes in the same order as production.
- Return `*gin.Engine` for both production and tests.

`cmd/server/main.go` should keep environment loading, MongoDB connection, index bootstrap,
bootstrap admin config, and `r.Run`. The new app/router layer should avoid reading process env
directly except through explicit config passed from `main`.

Reason: integration tests must cover real route wiring without copying route registration into
test-only code.

### MongoDB Test Database

Create test utilities under `internal/testutil`, for example:

- `mongo.go`
- `fixtures.go`
- `http.go`

Planned Mongo behavior:

- Use env `GYM_TEST_MONGODB_URI` when present.
- If it is absent, use the local Docker default:
  `mongodb://admin:password123@localhost:27017/?authSource=admin&directConnection=true`
- If MongoDB cannot be reached quickly, skip integration tests with `t.Skip`, not fail the full
  unit-test suite.
- Create a unique database per test or per test package, such as
  `gym_test_<unix_nano>_<short_random>`.
- Run `pkg/database.EnsureIndexes(ctx, db)` before seeding fixtures.
- Drop the test database in `t.Cleanup`.
- Never clean or drop the development database `gym_management`.

### Fixture Layer

Fixtures should create data through HTTP when the test is about API behavior, and through
repositories only when setup would distract from the behavior being tested.

Planned fixture helpers:

- create admin employee/bootstrap credentials
- login and return access/refresh tokens
- create authenticated request helpers
- create branch with unique `branch_code`
- create course with deterministic pricing/session count/tags
- create member with unique `ccid`
- create subscription and activate it through the existing member activation flow
- create session with matching branch/course level/tags
- parse JSON responses and assert shared error shape

Use deterministic timestamps where business rules depend on time. Prefer explicit RFC3339 dates in
requests rather than `time.Now()` inside tests.

## Integration Test Scope

### Must Cover In This Cycle

1. **Startup and router smoke**
   - Build router against test DB.
   - Ensure indexes are created.
   - `GET /ping` returns `200`.

2. **Auth and role guard**
   - Bootstrap or seed admin.
   - Login returns access + refresh token.
   - Protected route without token returns `401`.
   - Admin token can access an admin-only route.
   - Non-admin token cannot access employee management and returns `403`.
   - Refresh rotates token.
   - Logout revokes refresh token.

3. **Core member/subscription flow**
   - Create branch.
   - Create course.
   - Create member.
   - Create pending subscription.
   - Activate through offline payment endpoint.
   - Get subscription.
   - List member subscriptions.

4. **Data-integrity conflicts**
   - Duplicate member `ccid` returns `409 CONFLICT`.
   - Duplicate branch `branch_code` returns `409 CONFLICT`.
   - Duplicate refund for the same subscription returns `409 CONFLICT`.

5. **Attendance makeup guard**
   - Active subscription can report missed attendance.
   - Valid makeup consumes one session.
   - Reusing the same reported-missed date returns `409 CONFLICT`.

6. **Branch nearby**
   - Seed at least two branches with GeoJSON points.
   - `GET /branches/nearby` returns `200`.
   - Results are non-empty and include `distance_meters`.
   - Invalid coordinate/query returns `400`.

### Nice To Cover If Scope Stays Small

- Session create/enroll/check-in happy path.
- Duplicate session check-in returns `409 CONFLICT`.
- Session invalid tag returns `409 CONFLICT`.
- Employee create/list/update/password reset HTTP flow.
- Refund amount calculation through end-to-end HTTP response plus direct DB audit verification.

### Explicitly Out Of Scope For This Cycle

- GitHub Actions or CI service setup.
- Testcontainers or new Docker orchestration dependencies.
- Full exhaustive matrix of every invalid body for every endpoint.
- New business rules or API response shape changes.
- MongoDB transactions for multi-write refund/attendance flows.

## Layer Plan

### `internal/app`

- Add an app/router builder that production `main` and integration tests can both call.
- Keep dependency construction explicit and boring.
- Keep auth config and bootstrap admin config passed in from callers.
- Preserve route order, especially `/branches/nearby` before `/branches/:id`.

### `internal/testutil`

- Add Mongo test DB setup and cleanup helpers.
- Add fixture creation helpers.
- Add HTTP helpers for JSON requests, auth headers, response decoding, and error assertions.
- Keep helpers small enough that failed tests still show the API behavior clearly.

### Integration Tests

Recommended location:

- `internal/integration/...`

Alternative:

- `internal/handlers/..._integration_test.go`

Use a build tag only if the default test run becomes too environment-sensitive. Preferred first
approach: tests auto-skip when MongoDB is unavailable, so `go test ./...` remains developer-friendly.

### Docs

Update:

- `docs/local_dev_guide.md` with integration-test command and MongoDB requirement.
- `docs/code_reading_guide.md` with the new app/router entry point if introduced.
- `CHAT_CONTEXT/README.md` only after implementation/completion changes project state.
- `CHAT_CONTEXT/backend_skills/worklog.md` as the cycle moves through phases.

## Verification Plan

During implementation/review/test phases run:

```bash
env GOCACHE=/tmp/gocache go test ./...
env GOCACHE=/tmp/gocache go build ./...
git diff --check
```

For integration coverage with local MongoDB:

```bash
docker compose up -d mongodb
env GOCACHE=/tmp/gocache go test ./internal/integration -count=1
```

If tests are implemented in another package, use that package path instead.

## Risks And Decisions

- Integration tests need MongoDB. They should skip cleanly when MongoDB is unavailable unless an
  explicit env flag later says "fail if integration dependencies are missing".
- Test DB cleanup must be guarded so it never drops `gym_management`.
- Extracting router setup touches production startup path; review must compare route order and auth
  role matrices against `docs/api_contract.md`.
- Fixture helpers can become too large. Keep them API-focused and avoid building a second fake app.
- Time-based attendance rules need explicit timestamps to avoid flaky tests.
- Refund and attendance remain multi-write flows without transactions; tests should verify final
  observable state but not pretend atomicity is solved.
- Index creation can fail if duplicate fixture data is accidentally seeded before `EnsureIndexes`;
  generate unique IDs/codes per test.

## Handoff To Implementation

Start with:

1. `internal/app` router/dependency extraction.
2. `internal/testutil` Mongo + HTTP helpers.
3. A small integration smoke test for `/ping`, auth login, and one protected route.
4. Expand into the must-cover flows above.

Next file for `$gym-implement`:
`CHAT_CONTEXT/backend_skills/plans/08_integration_tests_fixtures.md`
