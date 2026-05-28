# Test - Validation Error Consistency

## Status

- Status: passed
- Feature: validation error consistency
- Plan file: `CHAT_CONTEXT/backend_skills/plans/06_validation_error_consistency.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/06_validation_error_consistency.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/06_validation_error_consistency.md`
- Tested at: 2026-05-26

## Commands

- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed a read-only module stat-cache warning
  but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.
- `docker ps --filter name=gym_mongodb --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'` -
  pass; local `gym_mongodb` was running.
- `env GOCACHE=/tmp/gocache PORT=18082 go run cmd/server/main.go` - pass after running outside the
  sandbox so the server could connect to local MongoDB.

## Manual API tests

Base URL: `http://127.0.0.1:18082`

### Happy path

- [x] `GET /ping` returned `200`.
- [x] `POST /api/v1/auth/login` with bootstrap admin returned `200` and an access token.

### Invalid input

- [x] `POST /api/v1/employees` with invalid JSON field type returned `400` +
  `error.code = INVALID_INPUT`.
- [x] `PATCH /api/v1/subscriptions/:id/suspend` with invalid `start_date` returned `400` +
  `error.code = INVALID_DATE`.

### Invalid ID

- [x] `GET /api/v1/employees/not-an-object-id` returned `400` + `error.code = INVALID_ID`.

### Auth and authorization

- [x] `GET /api/v1/courses` without token returned `401` + `error.code = UNAUTHORIZED`.
- [x] `GET /api/v1/courses` with malformed token returned `401` + `error.code = UNAUTHORIZED`.
- [x] Created a temporary receptionist employee, logged in as that employee, and confirmed
  `POST /api/v1/courses` returned `403` + `error.code = FORBIDDEN`.

### Not found

- [x] `GET /api/v1/employees/000000000000000000000000` returned `404` +
  `error.code = NOT_FOUND`.

### Conflict/business rule

- [x] Duplicate admin employee create request returned `409` + `error.code = CONFLICT`.

## Response-shape checks

- [x] Every manual error response checked above contained `error.details` as an object.
- [x] Invalid request body returned the sanitized `INVALID_INPUT` contract instead of the old raw
  Gin bind error shape.

## DB state verification

- [x] Temporary receptionist employee was deactivated through
  `PATCH /api/v1/employees/:id` after the forbidden-role check returned `200`.
- Direct MongoDB read was not required for this cycle because the feature only changes HTTP error
  response shape; no schema or persistence behavior was introduced.

## Issues found

- None.

## Final result

- Result: pass.
- Ready for `$gym-complete`: yes.

## Notes

- Initial attempt to run the server inside the sandbox failed because the sandbox could not open the
  local MongoDB socket. The server was then run with elevated local permissions for manual API
  verification.
- Manual test created and then deactivated temporary employee ID `6a15b633ac178aaaab1f83fe`.
