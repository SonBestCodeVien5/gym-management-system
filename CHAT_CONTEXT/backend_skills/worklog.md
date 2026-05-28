# Backend Worklog

Dùng file này để giữ roadmap và completion summary ngắn cho feature backend.

## Current backend roadmap

- [x] Refund flow & pricing rules
- [x] Branch nearby geo query
- [x] Attendance report/makeup endpoints nếu route còn thiếu
- [x] Auth/login + role guard
- [x] Employee management
- [x] Validation hardening & error consistency
- [ ] Indexes and data integrity
- [ ] Integration tests & fixtures

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
