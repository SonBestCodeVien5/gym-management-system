# Test - Indexes and Data Integrity

## Status

- Status: passed
- Feature: indexes and data-integrity hardening
- Plan file: `CHAT_CONTEXT/backend_skills/plans/07_indexes_data_integrity.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/07_indexes_data_integrity.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/07_indexes_data_integrity.md`
- Tested at: 2026-05-28

## Commands

| Command | Result | Notes |
|---|---|---|
| `env GOCACHE=/tmp/gocache go build ./...` | pass | Go printed the existing read-only module stat-cache warning but exited `0`. |
| `env GOCACHE=/tmp/gocache go test ./...` | pass | All packages passed. |
| `git diff --check` | pass | No whitespace errors. |
| `docker ps --filter name=gym_mongodb --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'` | pass | `gym_mongodb` was running and exposed `27017`. |
| `env GOCACHE=/tmp/gocache PORT=18083 go run cmd/server/main.go` | pass after elevated run | Sandbox could not open the local MongoDB socket; running outside the sandbox connected successfully. |

## Startup / Index Bootstrap

- [x] Server started against local MongoDB on `PORT=18083`.
- [x] Startup logged `Connected to MongoDB successfully`.
- [x] Startup logged `MongoDB indexes ensured successfully`.
- [x] Server registered routes and listened on `:18083`.

## DB Index Verification

Checked via `mongosh` inside `gym_mongodb`.

- [x] `members`: `_id_`, `ccid_1`
- [x] `branches`: `_id_`, `location_2dsphere`, `branch_code_unique`
- [x] `subscriptions`: `_id_`, `member_id_idx`, `status_idx`, `member_status_idx`, `course_id_idx`, `home_branch_id_idx`
- [x] `attendances`: `_id_`, `sub_id_date_desc_idx`, `session_id_idx`, `session_sub_unique`, `makeup_sub_ref_unique`
- [x] `sessions`: `_id_`, `branch_scheduled_at_idx`, `level_scheduled_at_idx`, `tags_idx`
- [x] `refunds`: `_id_`, `subscription_id_unique`, `member_id_idx`
- [x] `employees`: `_id_`, `normalized_email_unique`, `employee_id_unique`, `role_status_created_idx`, `branch_status_idx`
- [x] `refresh_tokens`: `_id_`, `token_hash_unique`, `employee_revoked_idx`, `expires_at_ttl`

## Manual API Tests

Base URL: `http://127.0.0.1:18083`

### Happy path

- [x] `GET /ping` returned `200`.
- [x] `POST /api/v1/auth/login` returned `200`.
- [x] `POST /api/v1/auth/refresh` returned `200`.
- [x] `POST /api/v1/auth/logout` returned `200`.
- [x] `POST /api/v1/branches` created a branch and returned `201`.
- [x] `GET /api/v1/branches/nearby` returned `200`.
- [x] `POST /api/v1/members` created a member and returned `201`.
- [x] `POST /api/v1/courses` created a course and returned `201`.
- [x] `POST /api/v1/subscriptions` created subscriptions and returned `201`.
- [x] `PATCH /api/v1/members/:id/activate` activated subscriptions and returned `200`.
- [x] `GET /api/v1/members/:id/subscriptions` returned `200`.
- [x] `POST /api/v1/subscriptions/:id/refund` returned `200`.
- [x] `POST /api/v1/sessions` created a session and returned `201`.
- [x] `POST /api/v1/sessions/:id/enroll` returned `200`.
- [x] `POST /api/v1/sessions/:id/checkin` returned `201`.
- [x] `POST /api/v1/attendance/report` returned `201`.
- [x] `POST /api/v1/attendance/makeup` returned `201`.
- [x] `GET /api/v1/employees?role=admin&status=active` returned `200`.
- [x] `GET /api/v1/sessions?level=basic&date=2026-06-13T00:00:00Z` returned `200`.

### Invalid input

- [x] `GET /api/v1/branches/not-an-object-id` returned `400` + `error.code = INVALID_ID`.

### Not found

- [x] `POST /api/v1/subscriptions/000000000000000000000000/refund` returned `404` +
  `error.code = NOT_FOUND`.

### Conflict / business rule

- [x] Duplicate `branch_code` returned `409` + `error.code = CONFLICT`.
- [x] Duplicate member `ccid` returned `409` + `error.code = CONFLICT`.
- [x] Duplicate refund request returned `409` + `error.code = CONFLICT`.
- [x] Duplicate session check-in returned `409` + `error.code = CONFLICT`.
- [x] Duplicate makeup reuse returned `409` + `error.code = CONFLICT`.

## DB State Verification

- [x] Refund audit uniqueness checked for test subscription `6a17dead392346c3dfb24332`.
- [x] `refunds.countDocuments({subscription_id})` returned `1` after first refund and duplicate
  refund attempt.
- [x] Subscription status for the refunded subscription was `refunded`.
- [x] Refunded subscription `remaining_sessions` was `0`.

## Issues found

- None.

## Final result

- Result: pass.
- Ready for `$gym-complete`: yes.

## Notes

- Initial server run inside the sandbox failed because it could not connect to
  `127.0.0.1:27017`; running the same command outside the sandbox succeeded.
- Curl from inside the sandbox could not reach the elevated server namespace, so manual API checks
  were run outside the sandbox.
- Manual API verification created persistent local test data with unique suffixes.
