# Code Review - Employee management

## Status

- Status: reviewed
- Feature: employee management
- Plan file: `CHAT_CONTEXT/backend_skills/plans/05_employee_management.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/05_employee_management.md`
- Reviewed at: 2026-05-26

## Review summary

- Result: pass, no blocking findings
- Reviewer: Codex
- Build status: pass
- Test status: pass

## Checklist

- [x] Code compiles.
- [x] Handler only handles HTTP parse/response.
- [x] Service owns business rules.
- [x] Repository only handles DB.
- [x] Model tags match API/DB contract.
- [x] Errors map to correct HTTP status.
- [x] Atomic updates used where needed.
- [x] No unsafe client-controlled password hash or auth-computed field.
- [x] Routes have correct order.
- [x] Docs/API samples match behavior.

## Passed

- Employee routes are protected by `AuthRequired` and `RequireRoles(service.RoleAdmin)`.
- Handler parses JSON, path/query ObjectIDs, and actor ID; business rules stay in service.
- Service validates roles, status, level, password policy, branch references, and self-lockout.
- Password hashes are generated with bcrypt; API responses use a safe DTO without `password_hash`
  or `normalized_email`.
- Employee repository owns MongoDB list/update/password operations and maps duplicate key errors.
- Refresh-token repository can revoke active tokens by employee ID.
- Route order has `/employees/:id/password` before `/employees/:id` for PATCH.
- API contract and `api_test.http` include employee management endpoints and role behavior.

## Issues found

No blocking issues found during review.

| Severity | File | Issue | Fix |
|---|---|---|---|
| none | - | No blocking issue found. | No fix applied. |

## Fixes applied during review

- None.

## Commands run

```bash
env GOCACHE=/tmp/gocache go test ./internal/service -run TestEmployeeService -count=1
env GOCACHE=/tmp/gocache go build ./...
git diff --check
env GOCACHE=/tmp/gocache go test ./...
```

Notes:
- `go build` passed. Go printed a stat-cache warning for the read-only module cache outside the
  workspace, but the command exited successfully.

## Remaining risks

- Password reset/profile update and refresh-token revocation are not transactional. A partial
  failure can leave profile/password changed while token revocation failed.
- The implementation prevents self-lockout but does not enforce "at least one active admin remains".
- Session creation still does not validate that `trainer_id` points to an active trainer.
- Manual API and Mongo verification have not been run yet in this review phase.

## Handoff to test

- Test admin create/list/get/update/reset-password endpoints with real tokens.
- Verify non-admin employee management returns `403` and missing token returns `401`.
- Verify created employee can login and inactive employee cannot login.
- Verify password reset and deactivation revoke active refresh tokens in MongoDB.
- Verify response bodies never include `password_hash` or `normalized_email`.
