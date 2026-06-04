# Backend Worklog

Dùng file này để giữ roadmap và completion summary ngắn cho feature backend.

## Current backend roadmap

- [x] Refund flow & pricing rules
- [x] Branch nearby geo query
- [x] Attendance report/makeup endpoints nếu route còn thiếu
- [x] Auth/login + role guard
- [x] Employee management
- [x] Validation hardening & error consistency
- [x] Indexes and data integrity
- [x] Integration tests & fixtures
- [x] Frontend readiness mini-cycle
- [x] Dashboard/report aggregate APIs
- [x] Final project package: seed data, Docker, README, report alignment

---

# Feature - Final project package

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary - 2026-06-04

### Goal
Finish the project as a runnable, demonstrable, and report-ready submission package: deterministic
demo data, full-stack Docker, root README/project overview, local run guidance, and report material
alignment.

### Key decisions
- Do not add new public HTTP endpoints in this phase unless implementation finds a small support gap
  that must be explicitly documented.
- Add a seed/demo command instead of reusing integration-test fixtures for local data.
- Make seed data idempotent by stable unique keys and require any destructive reset to be explicit.
- Add configurable `DB_NAME` support so local, Docker, and seed flows can target the same intended
  database.
- Package MongoDB, backend API, and frontend in Docker with documented local placeholder secrets.
- Reconcile stale report material so employee management and dashboard/report APIs are shown as
  implemented.

### Next action
Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/11_final_project_package.md`.

## Implementation summary - 2026-06-04

### Result
- Added configurable `DB_NAME` support with `gym_management` fallback.
- Added `cmd/seed` deterministic demo data command and seeded the local `gym_management` database.
- Added full-stack Docker packaging for MongoDB, API, frontend, and a `seed` compose profile.
- Rewrote root README as final evaluator quickstart with Docker, seed, demo accounts, local dev,
  verification commands, and docs/report links.
- Updated `.env`, `.env.example`, `frontend/.env.example`, local dev docs, docs hub, API samples,
  report evidence, and a final report assembly draft.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `npm --prefix frontend run build` - pass.
- `docker compose config` - pass.
- `docker compose --profile seed config` - pass.
- `env DOCKER_BUILDKIT=0 docker compose build` outside sandbox - pass. Plain `docker compose build`
  could not run because the local Docker CLI is missing the `docker-buildx` plugin.
- `env GOCACHE=/tmp/gocache go run ./cmd/seed` outside sandbox - pass and seeded demo data.

### Next action
Use `$gym-review` with `CHAT_CONTEXT/backend_skills/implementations/11_final_project_package.md`.

## Review summary - 2026-06-04

### Result
- Review requested fixes before final test/complete.
- Main finding: seed natural-key upserts can leave relationship records pointing at fixed demo IDs
  when matching natural-key records already exist with different `_id` values.
- Additional findings: frontend Docker context needs its own `.dockerignore`, Docker buildx fallback
  should be documented, and `api_test.http` still contains stale sample IDs.

### Verification
- `git diff --check` - pass.
- Build/test/Docker/seed pass evidence remains in the implementation note; review did not rerun the
  full command matrix.

### Next action
Use `$gym-implement` to fix findings from
`CHAT_CONTEXT/backend_skills/reviews/11_final_project_package.md`.

## Review fix implementation - 2026-06-04

### Result
- Fixed seed relationship safety by capturing actual `_id` values after natural-key upserts and using
  them for branches, subscriptions, attendances, sessions, and refunds.
- Added `frontend/.dockerignore` for reproducible frontend Docker context.
- Documented the `DOCKER_BUILDKIT=0 docker compose build` fallback in README/local dev guide.
- Replaced stale `api_test.http` ObjectIDs with seeded demo variables and avoided sample employee
  email conflict with the seeded trainer.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `env GOCACHE=/tmp/gocache go run ./cmd/seed` outside sandbox - pass on an already-seeded DB.
- `npm --prefix frontend run build` - pass.
- `docker compose config` - pass.
- `docker compose --profile seed config` - pass.
- `env DOCKER_BUILDKIT=0 docker compose build` outside sandbox - pass.

### Next action
Use `$gym-review` with `CHAT_CONTEXT/backend_skills/implementations/11_final_project_package.md`.

## Re-review summary - 2026-06-04

### Result
- Review passed after the fix pass.
- Previous high seed relationship finding is resolved by actual-ID capture after natural-key upserts.
- Previous Docker/docs/API sample findings are resolved by `frontend/.dockerignore`,
  buildx fallback docs, and seeded demo ID variables in `api_test.http`.

### Verification
- `git diff --check` - pass.
- Build/test/Docker/seed pass evidence remains in the implementation note; re-review did not rerun
  the full command matrix.

### Next action
Use `$gym-test` with `CHAT_CONTEXT/backend_skills/reviews/11_final_project_package.md`.

## Test summary - 2026-06-04

### Result
- Final package test pass with one local-volume caveat.
- Go build/test, frontend build, compose config/profile config, Docker legacy-builder build, clean
  Compose stack, seed idempotency, protected API smoke, frontend login/dashboard smoke, and DB count
  checks passed.
- Default `docker compose up -d` hit an existing local `gym-management-system_mongo_data` volume with
  Mongo featureCompatibilityVersion `8.2`, which `mongo:7` cannot open. The volume was preserved;
  clean-volume verification used `docker compose -p gym-final-test ...`.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass with existing read-only stat-cache warning.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `npm --prefix frontend run build` - pass.
- `docker compose config` and `docker compose --profile seed config` - pass.
- `env DOCKER_BUILDKIT=0 docker compose build` - pass.
- `env DOCKER_BUILDKIT=0 docker compose -p gym-final-test up -d --build` - pass.
- `env DOCKER_BUILDKIT=0 docker compose -p gym-final-test --profile seed run --rm seed` twice -
  pass.
- Manual API/frontend smoke - pass.
- `git diff --check` - pass.

### Next action
Use `$gym-complete` with `CHAT_CONTEXT/backend_skills/tests/11_final_project_package.md`.

## Completion - 2026-06-04

### Result
- Final project package cycle completed.
- Durable docs, API samples, report material, backend memory, and chat snapshot are aligned with the
  shipped Docker/seed/README package.
- No public HTTP contract changed in this phase; `docs/api_contract.md` remains the current API
  contract.

### Docs updated
- [x] `README.md`
- [x] `docs/README.md`
- [x] `docs/local_dev_guide.md`
- [x] `docs/report-materials/07_current_implementation_evidence.md`
- [x] `docs/report-materials/README.md`
- [x] `docs/report-materials/final_report.md`
- [x] `api_test.http`
- [x] `CHAT_CONTEXT/README.md`
- [x] `CHAT_CONTEXT/backend_skills/worklog.md`

### Remaining risks
- Existing local Docker volumes created by newer MongoDB versions may need an intentional
  `docker compose down -v` reset before using the default `mongo:7` service.
- Frontend Docker API base URL is build-time config.
- Branch-scope authorization, report export, online payment, notifications, and Member App remain
  future work.

### Next action
Use `$gym-git` to inspect the full diff and prepare the final commit, or hand off the project for
submission review.

---

# Feature - Dashboard/report aggregate APIs

## Status
- Planned: yes
- Implemented: yes
- Reviewed: no
- Tested: yes
- Docs updated: yes

## Plan summary - 2026-06-02

### Goal
Add read-only dashboard/report aggregate endpoints for FE11 live dashboard data:
summary KPIs, revenue buckets, plan distribution, recent members, and today's sessions.

### Key decisions
- Use protected `GET /api/v1/dashboard/*` endpoints for `admin` and `manager`.
- Keep endpoints read-only and server-computed; clients cannot provide money totals.
- Net revenue means paid subscription totals in range minus processed refund totals in range.
- Add focused read indexes for dashboard date-range queries; no schema changes.
- Keep recent members unscoped by branch because members do not currently store a branch field.

### Next action
Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/10_dashboard_reports.md`.

## Implementation summary - 2026-06-02

### Result
- Added admin/manager-only dashboard aggregate endpoints:
  - `GET /api/v1/dashboard/summary`
  - `GET /api/v1/dashboard/revenue`
  - `GET /api/v1/dashboard/plans`
  - `GET /api/v1/dashboard/members/recent`
  - `GET /api/v1/dashboard/sessions/today`
- Added dashboard DTOs, repository aggregate/read queries, service defaults/deltas/validation, handler
  query parsing, route wiring, and dashboard-supporting read indexes.
- Revenue is net revenue: paid subscription totals minus processed refunds.
- Updated `docs/api_contract.md` and `api_test.http`.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass, including `internal/integration`.
- `git diff --check` - pass.

### Next action
Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/11_live_dashboard_apis.md`.

## Completion - 2026-06-02

### Result
- Dashboard/report aggregate API cycle completed with implementation and test evidence.
- Backend now provides the FE11 live dashboard contract.
- Review phase was not run separately in this chained pass; implementation and test notes record the
  main review focus areas and residual risks.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `api_test.http`
- [x] `CHAT_CONTEXT/backend_skills/worklog.md`
- [x] `CHAT_CONTEXT/README.md`

### Follow-up risks
- Recent members are not branch-scoped because member documents have no branch field.
- Revenue semantics are limited to subscription payments minus processed refunds.
- Dashboard report export and broader `/api/v1/reports/*` APIs remain future work.

### Next action
Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/11_live_dashboard_apis.md`.

---

# Feature - Frontend readiness mini-cycle

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary - 2026-05-28

### Goal
Prepare the backend for initial frontend integration without adding new business-domain behavior.

### Key decisions
- Add explicit allow-list CORS support through `CORS_ALLOWED_ORIGINS` for browser FE dev origins.
- Apply CORS globally before auth guards so preflight requests do not fail on missing bearer tokens.
- Add protected `GET /api/v1/auth/me` so FE can restore current staff context after reload or token
  refresh.
- Return the same compact employee shape as login and keep refresh-token behavior unchanged.
- Do not add schema changes, indexes, migrations, or a seed command in this mini-cycle.
- Update API/local-dev docs, `.env.example`, REST samples, and focused integration tests.

### Next action
Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/09_frontend_readiness.md`.

## Implementation summary - 2026-05-28

### Result
- Added allow-list CORS support through `CORS_ALLOWED_ORIGINS`.
- Added global preflight handling before auth guards.
- Added protected `GET /api/v1/auth/me` using existing access-token validation and compact login
  employee response shape.
- Added service/unit and integration coverage for current employee lookup and CORS preflight.
- Updated API contract, local-dev guide, `.env.example`, REST sample, and chat snapshot.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1` - pass.

### Next action
Use `$gym-review` with `CHAT_CONTEXT/backend_skills/implementations/09_frontend_readiness.md`.

## Review summary - 2026-05-28

### Result
- Review passed with no blocking findings.
- Checked CORS middleware ordering, explicit origin allow-list, no wildcard/credentials, `/auth/me`
  route placement, handler/service boundaries, error mapping, docs/API sample alignment, and focused
  test coverage.

### Verification
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass; integration package used cached result in that
  full-suite run.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1` - pass.

### Next action
Use `$gym-test` with `CHAT_CONTEXT/backend_skills/reviews/09_frontend_readiness.md`.

## Test summary - 2026-05-28

### Result
- Automated build/tests passed.
- Integration tests passed against real local MongoDB outside sandbox.
- Manual local-server checks passed for allowed CORS preflight, disallowed-origin preflight, login
  then `/auth/me`, and missing-token `/auth/me`.
- Integration DB cleanup check passed: no leftover `gym_test_*` databases.

### Verification
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass; integration package used cached result in that
  full-suite run.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1 -v` in sandbox - pass with all
  integration cases skipped because sandbox could not open `localhost:27017`.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1 -v` outside sandbox - pass
  against real MongoDB.
- Manual API checks against temporary server on `PORT=18084` - pass.
- Direct MongoDB `gym_test_*` cleanup check - pass.

### Next action
Use `$gym-complete` with `CHAT_CONTEXT/backend_skills/tests/09_frontend_readiness.md`.

## Completion - 2026-05-28

### Result
- Frontend readiness mini-cycle completed end-to-end.
- Backend now supports explicit allow-list CORS through `CORS_ALLOWED_ORIGINS`.
- Browser preflight requests are handled globally before auth guards.
- Protected `GET /api/v1/auth/me` returns the current active employee from the bearer access token.
- API contract, REST sample, `.env.example`, local-dev guide, code-reading guide, README, backend
  phase notes, and chat snapshot are aligned.
- No schema, index, migration, seed command, cookie auth, or FE implementation was added.

### Verification
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass; integration package used cached result in that
  full-suite run.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1 -v` outside sandbox - pass
  against real MongoDB.
- Manual API checks against temporary server on `PORT=18084` - pass.
- Direct MongoDB `gym_test_*` cleanup check - pass.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `docs/local_dev_guide.md`
- [x] `docs/code_reading_guide.md`
- [x] `README.md`
- [x] `.env.example`
- [x] `api_test.http`
- [x] `CHAT_CONTEXT/README.md`

### Follow-up risks
- No browser UI has consumed the API yet; first FE pass should still verify real browser requests.
- Seed/demo data remains a follow-up if frontend demos need predictable sample records.
- CI automation is still not configured.
- Expanded session/not-found integration coverage remains a good backlog item.

### Next action
Use `$gym-git` to review/commit/push Cycle 09, or `$gym-plan` for the next backlog item.

---

# Feature - Integration tests & fixtures

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary - 2026-05-28

### Goal
Add a reusable integration-test harness and fixtures so current backend behavior can be verified
through Go tests with real HTTP routing, real services/repositories, and a real MongoDB test DB.

### Key decisions
- Do not add or change public API endpoints in this cycle.
- Extract reusable router/dependency construction from `cmd/server/main.go` so production and tests
  use the same route registration.
- Add MongoDB test helpers that create a unique test database, run `pkg/database.EnsureIndexes`, and
  drop only that test database in cleanup.
- Keep `go test ./...` developer-friendly by skipping integration tests when MongoDB is not
  reachable.
- Cover auth/role guard, member-subscription activation, duplicate unique conflicts, attendance
  makeup reuse, and branch nearby as the must-have integration flows.
- Keep CI service setup, Testcontainers, exhaustive invalid-body matrices, and transaction work out
  of this cycle.

### Next action
Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/08_integration_tests_fixtures.md`.

## Implementation summary - 2026-05-28

### Result
- Extracted reusable route/dependency wiring into `internal/app`.
- Kept production startup in `cmd/server/main.go` focused on env, MongoDB connection, index
  bootstrap, app config, and server run.
- Added `internal/testutil` for isolated MongoDB test DB setup, index bootstrap, HTTP helpers, auth
  login, and fixtures.
- Added `internal/integration` tests for auth/role guard, member-subscription activation, duplicate
  conflicts, attendance makeup reuse, branch nearby, and refund audit count.
- Updated local-dev/code-reading docs and README with the integration test surface.

### Verification
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1` - pass.

### Next action
Use `$gym-review` with `CHAT_CONTEXT/backend_skills/implementations/08_integration_tests_fixtures.md`.

## Review summary - 2026-05-28

### Result
- Review passed with no blocking findings.
- Checked production route extraction, route order, role matrix, test DB cleanup guard,
  integration-test coverage, docs alignment, and remaining coverage gaps.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./internal/integration -count=1` - pass.
- `git diff --check` - pass.

### Next action
Use `$gym-test` with `CHAT_CONTEXT/backend_skills/reviews/08_integration_tests_fixtures.md`.

## Test summary - 2026-05-28

### Result
- Automated build/tests passed.
- Integration tests passed against real local MongoDB on `localhost:27017`.
- Integration suite verified auth/role guard, token refresh/logout, member-subscription activation,
  duplicate branch/member/refund conflicts, attendance makeup reuse, branch nearby, and invalid
  nearby query behavior.
- DB cleanup verification passed: no `gym_test_*` databases remained after tests.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass in sandbox; integration package used cached
  result in that run.
- `go test ./internal/integration -count=1 -v` outside sandbox - pass with real MongoDB.
- `go test ./...` outside sandbox - pass with real MongoDB.
- `git diff --check` - pass.
- `docker exec gym_mongodb mongosh ... listDatabases` - pass; no leftover `gym_test_*` databases.

### Next action
Use `$gym-complete` with `CHAT_CONTEXT/backend_skills/tests/08_integration_tests_fixtures.md`.

## Completion - 2026-05-28

### Result
- Integration tests and fixtures cycle completed end-to-end.
- Production route/dependency wiring now lives in `internal/app` and is reused by tests.
- `cmd/server/main.go` remains the production entry point for env loading, MongoDB connection,
  central index bootstrap, app config, and server run.
- `internal/testutil` provides isolated MongoDB test DB setup, index bootstrap, auth login, HTTP
  helpers, and core fixtures.
- `internal/integration` verifies key HTTP behavior through real router/service/repository/MongoDB
  paths.
- Durable local-dev/code-reading docs, README, backend phase notes, and chat snapshot are aligned.
- No public API contract or REST sample behavior changed, so `docs/api_contract.md` and
  `api_test.http` did not need Cycle 08 edits.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass in sandbox; integration package used cached
  result in that run.
- `go test ./internal/integration -count=1 -v` outside sandbox - pass with real MongoDB.
- `go test ./...` outside sandbox - pass with real MongoDB.
- `git diff --check` - pass.
- `docker exec gym_mongodb mongosh ... listDatabases` - pass; no leftover `gym_test_*` databases.

### Docs updated
- [x] `docs/local_dev_guide.md`
- [x] `docs/code_reading_guide.md`
- [x] `README.md`
- [x] `CHAT_CONTEXT/README.md`

### Follow-up risks
- Session create/enroll/check-in integration coverage remains a good follow-up.
- HTTP not-found integration coverage is still limited.
- CI automation is not configured yet.
- Refund and attendance multi-write flows remain without MongoDB transactions.

### Next action
Use `$gym-git` to review/commit/push Cycle 08, or `$gym-plan` for the next backlog item.

---

# Feature - Indexes and data integrity

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary - 2026-05-28

### Goal
Add central MongoDB index bootstrap and harden DB-backed integrity for current API surfaces without
adding new endpoints.

### Key decisions
- Prefer `pkg/database.EnsureIndexes(ctx, db)` called from startup before repository construction.
- Keep startup idempotent and fail fast if unique indexes cannot be created because local data is
  dirty.
- Preserve current HTTP success shapes and shared error envelope.
- Map duplicate-key errors to domain conflicts instead of returning raw Mongo errors.
- Add unique indexes for `branches.branch_code` and `refunds.subscription_id`.
- Add query indexes for subscriptions, attendances, sessions, employees, and refresh tokens.
- Add partial unique attendance indexes for duplicate session check-in and makeup reuse.
- Defer full Mongo transactions for refund/attendance multi-write flows unless implementation finds
  a small, safe pattern.

### Next action
Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/07_indexes_data_integrity.md`.

## Implementation summary - 2026-05-28

### Result
- Added central MongoDB index bootstrap at startup.
- Moved existing repository-constructor index ownership into `pkg/database.EnsureIndexes`.
- Added unique/query/partial unique/TTL indexes for current API surfaces.
- Added duplicate-key normalization and conflict mapping for member, branch, refund, and attendance
  paths.
- Updated API docs, local dev guide, and REST samples for visible integrity behavior.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.

### Next action
Use `$gym-review` with `CHAT_CONTEXT/backend_skills/implementations/07_indexes_data_integrity.md`.

## Review summary - 2026-05-28

### Result
- Review passed with no blocking findings.
- Checked index definitions, startup ordering, repository/service/handler error mapping, docs/API
  sample alignment, and known transaction/reference-hardening deferrals.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.

### Next action
Use `$gym-test` with `CHAT_CONTEXT/backend_skills/reviews/07_indexes_data_integrity.md`.

## Test summary - 2026-05-28

### Result
- Automated build/tests passed.
- Server startup against local MongoDB passed and logged central index bootstrap success.
- Direct MongoDB index inspection passed for members, branches, subscriptions, attendances, sessions,
  refunds, employees, and refresh tokens.
- Manual API verification passed for auth, branch nearby, employee/session filters, subscription
  list, duplicate member CCID, duplicate branch code, duplicate refund, duplicate session check-in,
  and duplicate makeup reuse.
- Refund DB state verification passed: duplicate refund left exactly one refund audit row and the
  subscription stayed `refunded` with `remaining_sessions = 0`.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.
- Manual API script against local server on `PORT=18083` - pass.
- Direct MongoDB index/refund checks through `mongosh` - pass.

### Next action
Use `$gym-complete` with `CHAT_CONTEXT/backend_skills/tests/07_indexes_data_integrity.md`.

## Completion - 2026-05-28

### Result
- Indexes and data-integrity hardening cycle completed end-to-end.
- Centralized MongoDB index bootstrap in `pkg/database.EnsureIndexes`.
- Startup now ensures indexes before repository construction and fails fast on dirty duplicate data.
- Added unique, query, partial unique, and TTL indexes for current API surfaces.
- Mapped duplicate-key races to public `409 CONFLICT` responses through repository/service/handler
  layers.
- Updated durable API/local-dev docs, REST sample, backend phase notes, and chat snapshot.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.
- Manual API script against local server on `PORT=18083` - pass.
- Direct MongoDB index/refund checks through `mongosh` - pass.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `docs/local_dev_guide.md`
- [x] `docs/code_reading_guide.md`
- [x] `README.md`
- [x] `api_test.http`
- [x] `CHAT_CONTEXT/README.md`

### Follow-up risks
- Refund and attendance side effects remain multi-write flows without Mongo transactions.
- Branch manager and session branch/trainer reference hardening remains deferred.
- Integration tests should cover startup index creation and duplicate-key API behavior.

### Next action
Use `$gym-plan` or `$gym-implement` for the next backend cycle:
`CHAT_CONTEXT/backend_skills/plans/08_integration_tests_fixtures.md`.

---

# Feature - Validation hardening & error consistency

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary - 2026-05-26

### Goal
Chuẩn hóa toàn bộ backend error response sang contract ổn định:

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "invalid input",
    "details": {}
  }
}
```

### Key decisions
- Giữ nguyên success response shape hiện tại để giảm tác động lên FE/manual clients.
- Dùng enum lỗi chung: `INVALID_INPUT`, `INVALID_ID`, `INVALID_DATE`, `UNAUTHORIZED`,
  `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `INTERNAL_ERROR`.
- Không trả raw bind/Mongo/JWT/bcrypt/storage errors ra API.
- Handler tiếp tục chịu trách nhiệm parse HTTP input và map service errors sang status + code.
- Service/repository không biết HTTP response shape.

### Next action
Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/06_validation_error_consistency.md`.

## Implementation summary - 2026-05-26

### Result
- Added shared handler error response helpers with stable codes and nested `error` payloads.
- Migrated auth middleware and all current handlers to the shared error contract.
- Sanitized invalid request body responses so raw Gin bind errors are no longer returned.
- Kept success response shapes unchanged.
- Updated API contract and REST samples with the new error response contract.

### Verification
- `env GOCACHE=/tmp/gocache go test ./internal/handlers -count=1` - pass.
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.

### Next action
Use `$gym-review` with `CHAT_CONTEXT/backend_skills/implementations/06_validation_error_consistency.md`.

## Review summary - 2026-05-26

### Result
- Review passed with no blocking findings.
- Checked shared error helper, auth middleware, handler mapping, service/repository boundaries,
  old error-shape sweep, docs/API sample alignment, and focused middleware body assertions.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./internal/handlers -count=1` - pass.
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.

### Next action
Use `$gym-test` with `CHAT_CONTEXT/backend_skills/reviews/06_validation_error_consistency.md`.

## Test summary - 2026-05-26

### Result
- Automated build/tests passed.
- Manual API verification passed for shared error response contract:
  `UNAUTHORIZED`, `FORBIDDEN`, `INVALID_ID`, `INVALID_DATE`, `INVALID_INPUT`, `NOT_FOUND`, and
  `CONFLICT`.
- Every checked error response included nested `error.code` and object `error.details`.
- Temporary receptionist employee used for `403` verification was deactivated after the check.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.
- Manual API script against local server on `PORT=18082` - pass.

### Next action
Use `$gym-complete` with `CHAT_CONTEXT/backend_skills/tests/06_validation_error_consistency.md`.

## Completion - 2026-05-26

### Result
- Validation/error consistency cycle completed end-to-end.
- Added shared backend HTTP error contract:
  `{"error":{"code":"...","message":"...","details":{}}}`.
- Migrated all current handler error paths and auth middleware to stable public error codes.
- Sanitized invalid body responses so raw Gin binding errors are not returned.
- Preserved existing success response shapes.
- Updated API contract and REST samples with representative error checks.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.
- Manual API on `PORT=18082` - pass for `UNAUTHORIZED`, `FORBIDDEN`, `INVALID_ID`, `INVALID_DATE`,
  `INVALID_INPUT`, `NOT_FOUND`, and `CONFLICT`.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `api_test.http`
- [x] `README.md`
- [x] `docs/code_reading_guide.md`
- [x] `CHAT_CONTEXT/README.md`

### Follow-up risks
- Handler-wide automated body assertions are still limited; broader integration tests can cover the
  full API contract later.
- Existing MVP data-integrity limitations remain outside this cycle: last-active-admin enforcement,
  trainer reference validation for sessions, and transactional refresh-token revocation.
- Manual test created and deactivated temporary employee `6a15b633ac178aaaab1f83fe`.

### Next action
Use `$gym-plan` or `$gym-implement` for the next backend cycle:
`CHAT_CONTEXT/backend_skills/plans/07_indexes_data_integrity.md`.

---

# Feature - Employee management

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary - 2026-05-26

### Goal
Thêm API admin-only để tạo, list, xem chi tiết, cập nhật, vô hiệu hóa, và reset mật khẩu cho staff
account sau cycle bootstrap admin/auth.

### Planned API
- `POST /api/v1/employees`
- `GET /api/v1/employees`
- `GET /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id/password`

### Key decisions
- Không hard delete trong cycle này; offboarding dùng `status = inactive`.
- Response không được expose `password_hash` hoặc `normalized_email`.
- Employee management là admin-only.
- Password reset và update từ active sang inactive nên revoke refresh token active của employee đó.
- Admin tự deactivate hoặc tự remove role `admin` của mình nên bị conflict để giảm rủi ro tự khóa hệ
  thống.

### Next action
Dùng `$gym-review` với `CHAT_CONTEXT/backend_skills/implementations/05_employee_management.md`.

## Implementation summary - 2026-05-26

### Result
- Added admin-only employee create/list/get/update/password reset endpoints.
- Added employee service validation for role/status/level/password, email normalization, branch
  references, and self-lockout prevention.
- Added employee repository list/update/password-update operations and duplicate-key mapping.
- Added refresh-token revoke by employee ID for password reset and deactivation.
- Updated API contract, REST samples, code-reading guide, and local dev checkpoint.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.

### Review handoff
- Review route authorization, partial update semantics, refresh-token revocation ordering, and
  employee response safety.

## Review summary - 2026-05-26

### Result
- Review passed with no blocking findings.
- Checked route authorization/order, handler/service/repository ownership, error mapping, response
  safety, docs/API sample alignment, and focused employee service tests.

### Verification
- `env GOCACHE=/tmp/gocache go test ./internal/service -run TestEmployeeService -count=1` - pass.
- `env GOCACHE=/tmp/gocache go build ./...` - pass.
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.

### Next action
Use `$gym-test` with `CHAT_CONTEXT/backend_skills/reviews/05_employee_management.md`.

## Test summary - 2026-05-26

### Result
- Automated build/tests passed.
- Manual API verification passed for admin create/list/get/update/reset/deactivate, invalid input,
  not found, duplicate conflict, self-deactivation conflict, missing-token `401`, non-admin `403`,
  refresh-token revoke after reset, inactive login rejection, and inactive access-token rejection.
- Direct MongoDB verification passed for normalized email, password hash presence, inactive final
  status, and revoked refresh tokens.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./internal/service -run TestEmployeeService -count=1` - pass.
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- Manual API script against local server on `PORT=18081` - pass.
- Temporary Go DB check script - pass.

### Next action
Use `$gym-complete` with `CHAT_CONTEXT/backend_skills/tests/05_employee_management.md`.

## Completion - 2026-05-26

### Result
- Employee management cycle completed end-to-end.
- Added admin-only APIs for employee create/list/get/update/password reset.
- Employee create/update normalizes email, validates role/status/level/password and branch
  references, hashes passwords with bcrypt, and returns safe employee responses.
- Password reset and active-to-inactive deactivation revoke active refresh tokens for the employee.
- Admin self-deactivation and self-removal of `admin` role return conflict.
- Durable docs and API samples were aligned with implemented behavior.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./internal/service -run TestEmployeeService -count=1` - pass.
- `git diff --check` - pass.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- Manual API - pass for admin login, employee create/list/get/update/reset/deactivate, invalid
  input, not found, duplicate conflict, self-deactivation conflict, missing-token `401`, non-admin
  `403`, refresh-token revoke after reset, inactive login rejection, and inactive access-token
  rejection.
- Direct Mongo verification - pass for normalized email, bcrypt-like password hash, inactive final
  status, and revoked refresh tokens.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `api_test.http`
- [x] `README.md`
- [x] `docs/local_dev_guide.md`
- [x] `docs/code_reading_guide.md`
- [x] `CHAT_CONTEXT/README.md`

### Follow-up risks
- Password/profile update and refresh-token revocation are not transactional.
- Last-active-admin invariant is not enforced beyond self-lockout prevention.
- Session create still does not validate `trainer_id` as an active trainer.
- Manual test data remains as inactive employee `codex.employee.1779803637@gym.test` with revoked
  refresh tokens.

### Next action
Use `$gym-plan` or `$gym-implement` for `CHAT_CONTEXT/backend_skills/plans/06_validation_error_consistency.md`.

---

# Feature - Auth/login + role guard

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Completion - 2026-05-25

### Result
- Added employee login with bcrypt password verification.
- Added access + refresh token issue, refresh rotation, logout revoke, and refresh-token hash
  persistence.
- Added env bootstrap for first admin account.
- Added auth middleware and role guard for current business routes.
- Updated API contract and REST samples for auth and protected routes.
- Updated local development and code-reading docs for auth/env/role guard.

### Verification
- `env GOCACHE=/tmp/gocache go build ./...` - pass.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- Manual API - pass for health, admin login, protected route with/without token, wrong password,
  missing refresh token, refresh rotation, reused old refresh token, logout idempotency, inactive
  employee login rejection, and receptionist role forbidden check.
- Direct Mongo verification - pass for admin bootstrap, bcrypt password hash presence, refresh-token
  hash storage, absence of raw token fields, and `revoked_at` after refresh/logout.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `api_test.http`
- [x] `README.md`
- [x] `docs/local_dev_guide.md`
- [x] `docs/code_reading_guide.md`
- [x] `CHAT_CONTEXT/README.md`

### Follow-up risks
- Refresh-token TTL cleanup remains for the index/data-integrity cycle.
- Refresh rotation can invalidate the old token before replacement persistence succeeds; accepted as
  residual availability risk for MVP.
- Employee management is now complete; validation/error consistency hardening is the next backend
  cycle.

---

# Feature — Refund flow & pricing rules

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary

### Goal
Implement `POST /api/v1/subscriptions/:id/refund` and pricing discount rules for subscription creation.

### API
- `POST /api/v1/subscriptions/:id/refund`
- Request:
```json
{
  "reason": "member requested cancellation"
}
```
- Response should include refund record or refund summary.

### Business rules
- Only `active` subscription can be refunded.
- Cannot refund `pending`, `suspended`, `expired`, `refunded`.
- Cannot refund if `remaining_sessions <= 0`.
- `used_sessions = total_sessions - remaining_sessions`.
- `refund_amount = total_amount_paid * remaining_sessions / total_sessions`.
- After refund:
  - subscription `status = refunded`
  - `remaining_sessions = 0`
  - refund record inserted.
- Prevent double refund via atomic update and/or unique index.

### Pricing rules
- Server calculates money from course snapshot.
- Optional discount:
  - `none`
  - `percent`
  - `fixed`
- Percent must be `0 <= value <= 100`.
- Fixed must be `0 <= value <= subtotal`.
- `total_amount_paid = subtotal - discount_amount`.

### Files expected
- `internal/models/subscription.go`
- `internal/models/refund.go`
- `internal/repository/subscription_repo.go`
- `internal/repository/refund_repo.go`
- `internal/service/subscription_service.go`
- `internal/handlers/subscription_handler.go`
- `cmd/server/main.go`
- `docs/api_contract.md`
- `api_test.http`

## Completion — 2026-05-20

### Result
- `POST /api/v1/subscriptions` pricing/discount implemented, reviewed, tested, and documented.
- `POST /api/v1/subscriptions/:id/refund` implemented, reviewed, tested, and documented.
- Pricing is server-calculated from course snapshot:
  - `subtotal_amount`
  - `discount_amount`
  - `total_amount_paid`
- Refund allows only `active` subscriptions with valid remaining sessions.
- Refund atomically changes subscription to `refunded` and `remaining_sessions = 0`, then inserts refund audit record.

### Verification
- `go build ./...` — pass in test phase.
- `go test ./...` — pass in test phase.
- Automated service tests — pass for pricing, invalid discount inputs, refund conflict cases, duplicate prevention, and success.
- Manual API — pass for create subscription pricing, activation, refund success, and post-refund subscription state.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `api_test.http`
- [x] `CHAT_CONTEXT/README.md`

### Remaining risks
- `refunds.subscription_id` unique index is not bootstrapped yet; track under `07_indexes_data_integrity`.
- No Mongo transaction around subscription update + refund audit insert; partial failure risk remains accepted for MVP.
- Rare delete/race case may return `409` instead of `404`.
- Refund handler requires JSON body; empty body returns `400`.

---

# Feature — Branch nearby geo query

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary

### Goal
Implement `GET /api/v1/branches/nearby`.

### API
- `GET /api/v1/branches/nearby?lng=106.7&lat=10.8&max_distance=5000&limit=10`

### Business rules
- Validate lng/lat range.
- Default `max_distance = 5000`.
- Default `limit = 10`, max `100`.
- GeoJSON coordinate order is `[lng, lat]`.
- Route must be before `/branches/:id`.

### Data/index
- Mongo index: `branches.location` 2dsphere.

### Files expected
- `internal/repository/branch_repo.go`
- `internal/service/branch_service.go`
- `internal/handlers/branch_handler.go`
- `cmd/server/main.go`
- Mongo index bootstrap location
- `docs/api_contract.md`
- `api_test.http`

## Completion — 2026-05-20

### Result
- `GET /api/v1/branches/nearby` implemented, reviewed, tested, and documented.
- Query uses required `lng`, `lat`; optional `max_distance`, `limit`.
- Response includes `distance_meters`.
- MongoDB `branches.location` 2dsphere index created at repository init.
- Route order verified: `/branches/nearby` before `/branches/:id`.

### Verification
- `go build ./...` — pass in test phase.
- `go test ./...` — pass in test phase.
- Manual API — pass for happy path, default query, invalid inputs, route order.
- Manual DB cleanup — done.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `api_test.http`
- [x] `CHAT_CONTEXT/README.md`

### Remaining risks
- Existing malformed `branches.location` documents can fail index creation or be excluded from geo results.

---

# Feature - Attendance report/makeup endpoints

## Status
- Planned: yes
- Implemented: yes
- Reviewed: yes
- Tested: yes
- Docs updated: yes

## Plan summary

### Goal
Expose dedicated attendance report/makeup routes without exposing client-controlled attendance status.

### API
- `POST /api/v1/attendance/report`
- `POST /api/v1/attendance/makeup`
- Report request uses `subscription_id`, `branch_id`, optional `date`.
- Makeup request uses `subscription_id`, `branch_id`, optional `date`, required `is_makeup_for`.

### Business rules
- Report stores `reported_missed`, keeps remaining sessions unchanged, and enforces one report in the 30-day window.
- Makeup stores `makeup`, must reference a reported-missed date within 7 days, cannot reuse the same reference, respects weekly limits, and consumes one remaining session.

## Completion - 2026-05-21

### Result
- Dedicated report and makeup handlers/routes are implemented.
- `AttendanceService.CheckIn` remains the shared rule path.
- `docs/api_contract.md` documents request, response, and status behavior.
- `api_test.http` has report and makeup request samples.

### Verification
- `go build ./...` - pass in test and re-review phases.
- `go test ./...` - pass in test and re-review phases.
- Manual API - pass for happy path, invalid input, not found, subscription-state conflict, report window conflict, missing/overdue makeup reference, and duplicate makeup.
- Direct Mongo verification - pass for attendance records, remaining-session decrement, member attended counter, and rejected overdue makeup non-insert.

### Docs updated
- [x] `docs/api_contract.md`
- [x] `api_test.http`
- [x] `CHAT_CONTEXT/README.md`

### Remaining risks
- Duplicate makeup protection is not DB-enforced yet; track under `07_indexes_data_integrity`.
- Attendance insert, subscription decrement, and member attended-count increment are not atomic as one unit.
- Makeup still references the exact reported-missed RFC3339 instant instead of a stable report ID.
- Feature-specific integration coverage remains for `08_integration_tests_fixtures`.
