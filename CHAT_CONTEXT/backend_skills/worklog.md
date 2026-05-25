# Backend Worklog

Dùng file này để giữ roadmap và completion summary ngắn cho feature backend.

## Current backend roadmap

- [x] Refund flow & pricing rules
- [x] Branch nearby geo query
- [x] Attendance report/makeup endpoints nếu route còn thiếu
- [x] Auth/login + role guard
- [ ] Employee management
- [ ] Validation hardening & error consistency
- [ ] Indexes and data integrity
- [ ] Integration tests & fixtures

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
- Employee management remains the next backend cycle for creating and maintaining non-bootstrap
  staff accounts.

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
