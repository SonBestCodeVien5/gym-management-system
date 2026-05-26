# Test - Employee management

## Status

- Status: tested
- Feature: employee management
- Plan file: `CHAT_CONTEXT/backend_skills/plans/05_employee_management.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/05_employee_management.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/05_employee_management.md`
- Tested at: 2026-05-26

## Commands

```bash
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./internal/service -run TestEmployeeService -count=1
git diff --check
env GOCACHE=/tmp/gocache go test ./...
env GOCACHE=/tmp/gocache PORT=18081 go run ./cmd/server
python3 <manual employee API script>
env GOCACHE=/tmp/gocache go run /tmp/gym_employee_db_check.go 6a15a5f5f6939d1a7c18fae4 codex.employee.1779803637@gym.test
```

## Command results

| Command | Result | Notes |
|---|---|---|
| `env GOCACHE=/tmp/gocache go build ./...` | pass | Go printed a stat-cache warning for read-only module cache outside workspace, but exited 0. |
| `env GOCACHE=/tmp/gocache go test ./internal/service -run TestEmployeeService -count=1` | pass | Focused employee service tests. |
| `git diff --check` | pass | No whitespace errors. |
| `env GOCACHE=/tmp/gocache go test ./...` | pass | Full automated test suite. |
| `env GOCACHE=/tmp/gocache PORT=18081 go run ./cmd/server` | pass | Required escalation because sandbox blocked local MongoDB socket access. Server stopped after manual tests. |
| Manual API script | pass | Required escalation because sandbox blocked local HTTP socket access. |
| DB check script | pass | Required escalation because direct MongoDB socket access is blocked in sandbox. |

## Manual API tests

Server:

- `PORT=18081`
- MongoDB local connection succeeded.
- Bootstrap admin login used `.env` values without printing secrets.

Test employee:

- ID: `6a15a5f5f6939d1a7c18fae4`
- Email: `codex.employee.1779803637@gym.test`
- Final status: `inactive`

### Happy path

- [x] Request: admin login.
- [x] Expected: `200` with access token.
- [x] Actual: `200`.
- [x] Result: pass.

- [x] Request: `POST /api/v1/employees`.
- [x] Expected: `201`, safe employee response, no `password_hash`, no `normalized_email`.
- [x] Actual: `201`; response safety predicate passed.
- [x] Result: pass.

- [x] Request: `GET /api/v1/employees?role=trainer&status=active`.
- [x] Expected: `200`, includes created employee.
- [x] Actual: `200`.
- [x] Result: pass.

- [x] Request: `GET /api/v1/employees/:id`.
- [x] Expected: `200`.
- [x] Actual: `200`.
- [x] Result: pass.

- [x] Request: `PATCH /api/v1/employees/:id` with `full_name` and `phone`.
- [x] Expected: `200`, updated safe response.
- [x] Actual: `200`.
- [x] Result: pass.

- [x] Request: created employee login with initial password.
- [x] Expected: `200`.
- [x] Actual: `200`.
- [x] Result: pass.

- [x] Request: `PATCH /api/v1/employees/:id/password` with valid new password.
- [x] Expected: `200`.
- [x] Actual: `200`.
- [x] Result: pass.

- [x] Request: created employee login with reset password.
- [x] Expected: `200`.
- [x] Actual: `200`.
- [x] Result: pass.

- [x] Request: deactivate created employee.
- [x] Expected: `200`.
- [x] Actual: `200`.
- [x] Result: pass.

### Invalid input

- [x] Request: create trainer without `level`.
- [x] Expected status: `400`.
- [x] Actual: `400`.
- [x] Result: pass.

- [x] Request: `GET /api/v1/employees/not-an-id`.
- [x] Expected status: `400`.
- [x] Actual: `400`.
- [x] Result: pass.

- [x] Request: reset password with short password.
- [x] Expected status: `400`.
- [x] Actual: `400`.
- [x] Result: pass.

### Auth and authorization

- [x] Request: `GET /api/v1/employees` without token.
- [x] Expected status: `401`.
- [x] Actual: `401`.
- [x] Result: pass.

- [x] Request: `GET /api/v1/employees` with trainer token.
- [x] Expected status: `403`.
- [x] Actual: `403`.
- [x] Result: pass.

- [x] Request: login inactive employee.
- [x] Expected status: `401`.
- [x] Actual: `401`.
- [x] Result: pass.

- [x] Request: use inactive employee access token on `GET /api/v1/sessions`.
- [x] Expected status: `401`.
- [x] Actual: `401`.
- [x] Result: pass.

### Not found

- [x] Request: `GET /api/v1/employees/000000000000000000000000`.
- [x] Expected status: `404`.
- [x] Actual: `404`.
- [x] Result: pass.

### Conflict/business rule

- [x] Request: duplicate employee create with same employee ID/email.
- [x] Expected status: `409`.
- [x] Actual: `409`.
- [x] Result: pass.

- [x] Request: refresh old employee refresh token after password reset.
- [x] Expected status: `401`.
- [x] Actual: `401`.
- [x] Result: pass.

- [x] Request: admin self-deactivation.
- [x] Expected status: `409`.
- [x] Actual: `409`.
- [x] Result: pass.

## DB state verification

- [x] Expected DB changes:
  - Employee exists.
  - `email` and `normalized_email` equal normalized lowercase email.
  - `password_hash` exists and looks bcrypt-like.
  - Final employee `status` is `inactive`.
  - Refresh tokens created during employee login flows are revoked after reset/deactivation.
- [x] Actual DB changes:
  - `employee_found=true`
  - `email_normalized=true`
  - `password_hash_present=true`
  - `status_inactive=true`
  - `refresh_tokens_total=2`
  - `refresh_tokens_revoked=2`

## Issues found

No test failures found.

| Severity | Case | Issue | Fix |
|---|---|---|---|
| none | - | No issue found. | No fix applied. |

## Skipped / not covered

- Concurrent duplicate create/update races were not manually stress-tested. Unique indexes and
  duplicate-key mapping were covered by manual duplicate create.
- Last-active-admin invariant is intentionally not implemented in this cycle.
- Session `trainer_id` active-trainer validation is outside this cycle.
- Docker-based Mongo inspection was skipped because Docker socket access was denied in this
  environment; direct MongoDB verification was done with a temporary Go script instead.

## Final result

- Result: pass
- Ready for `$gym-complete`: yes
