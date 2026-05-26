# Implementation - Employee management

## Status

- Status: implemented
- Feature: employee management
- Plan file: `CHAT_CONTEXT/backend_skills/plans/05_employee_management.md`
- Started at: 2026-05-26
- Finished at: 2026-05-26

## Scope implemented

- [x] Model changes
- [x] Repository changes
- [x] Service changes
- [x] Handler changes
- [x] Route changes
- [x] Docs/API sample changes

## Files changed

- `cmd/server/main.go` - wired employee service, handler, and admin-only routes.
- `internal/repository/errors.go` - added storage-agnostic duplicate error.
- `internal/repository/employee_repo.go` - added list, partial update, password update, duplicate
  key normalization.
- `internal/repository/refresh_token_repo.go` - added active refresh-token revoke by employee ID.
- `internal/service/employee_service.go` - added employee management business rules, password hash,
  branch validation, self-lockout guard, and refresh-token revocation.
- `internal/handlers/employee_handler.go` - added employee HTTP request parsing and error mapping.
- `internal/service/employee_service_test.go` - added focused employee service tests.
- `internal/service/auth_service_test.go` - updated auth test stubs for expanded repository
  interfaces.
- `docs/api_contract.md` - documented employee endpoints, role guard, request/response/status
  behavior.
- `api_test.http` - added employee management API samples.
- `docs/code_reading_guide.md` - added employee management code-reading path.
- `docs/local_dev_guide.md` - noted employee management checkpoint.

## Key decisions

- Employee management routes are admin-only.
- No hard delete endpoint was added; offboarding uses `status = inactive`.
- Create/update normalize employee email and store it as both `email` and `normalized_email` to
  match existing auth lookup behavior.
- Responses use a service DTO and do not expose `password_hash` or `normalized_email`.
- Password reset and active-to-inactive status transition revoke active refresh tokens for the
  target employee.
- Admin self-deactivation and self-removal of `admin` role return conflict.
- Branch IDs are validated against existing branch documents before persistence.

## Implementation notes

### Models

- Kept `models.Employee` as the persistence model.
- Added safe response and input DTOs in `employee_service.go` and request DTOs in
  `employee_handler.go`.

### Repository

- `EmployeeRepository.List` supports optional `role`, `status`, and `branch_id` filters.
- `EmployeeRepository.UpdateByID` uses `FindOneAndUpdate` with return-after to return the updated
  document.
- Mongo duplicate-key errors are mapped to `repository.ErrDuplicate`.
- `RefreshTokenRepository.RevokeActiveByEmployeeID` uses `UpdateMany` and succeeds even when no
  active refresh token exists.

### Service

- `EmployeeService` owns role/status/level/password validation and email normalization.
- Passwords are hashed with bcrypt and never accepted as a direct hash.
- `trainer` role requires a valid level.
- Update preserves omitted fields.

### Handler

- Handler parses JSON, ObjectIDs, query filters, and actor ID from auth middleware context.
- Error mapping follows current project contract: invalid input `400`, not found `404`, conflict
  `409`, unexpected `500`.

### Routes

- Added:
  - `POST /api/v1/employees`
  - `GET /api/v1/employees`
  - `GET /api/v1/employees/:id`
  - `PATCH /api/v1/employees/:id/password`
  - `PATCH /api/v1/employees/:id`

### Docs/API samples

- API contract and REST samples are aligned with implemented endpoint behavior.

## Commands run

```bash
gofmt -w cmd/server/main.go internal/repository/errors.go internal/repository/employee_repo.go internal/repository/refresh_token_repo.go internal/service/employee_service.go internal/service/employee_service_test.go internal/service/auth_service_test.go internal/handlers/employee_handler.go
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
```

## Known limitations

- Password reset/profile update and refresh-token revocation are not wrapped in a MongoDB
  transaction.
- Existing access tokens are rejected only on the next protected request because auth middleware
  reloads employee state then.
- The service prevents self-lockout, but it does not yet enforce "there must always be at least one
  active admin".
- Session creation still does not validate that `trainer_id` references an active trainer.

## Handoff to review

- Check employee route authorization and route matching around `/employees/:id/password`.
- Check service validation for role/status/level/password and branch references.
- Check partial update semantics and self-lockout behavior.
- Check ordering and residual risk around password/status update plus refresh-token revocation.
