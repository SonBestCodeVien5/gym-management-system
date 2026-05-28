# Implementation - Indexes and Data Integrity

## Status

- Status: implemented
- Feature: indexes and data-integrity hardening
- Plan file: `CHAT_CONTEXT/backend_skills/plans/07_indexes_data_integrity.md`
- Started at: 2026-05-28
- Finished at: 2026-05-28

## Scope implemented

- [ ] Model changes
- [x] Repository changes
- [x] Service changes
- [x] Handler changes
- [x] Route/startup changes
- [x] Docs/API sample changes

## Files changed

- `pkg/database/indexes.go` - central MongoDB index bootstrap.
- `cmd/server/main.go` - runs index bootstrap before repository construction.
- `internal/repository/member_repo.go` - duplicate-key normalization; index creation moved out.
- `internal/repository/branch_repo.go` - duplicate-key normalization; index creation moved out.
- `internal/repository/employee_repo.go` - index creation moved out.
- `internal/repository/refresh_token_repo.go` - index creation moved out.
- `internal/repository/refund_repo.go` - duplicate-key normalization.
- `internal/repository/attendance_repo.go` - duplicate-key normalization.
- `internal/service/member_service.go` - duplicate CCID race maps to member conflict.
- `internal/service/branch_service.go` - duplicate branch code maps to branch-code conflict.
- `internal/service/subscription_service.go` - duplicate refund audit maps to refund conflict.
- `internal/service/attendance_service.go` - duplicate attendance insert maps to makeup/session
  conflicts.
- `internal/handlers/branch_handler.go` - duplicate branch code returns `409`.
- `internal/handlers/attendance_handler.go` - duplicate session attendance conflict maps to `409`.
- `docs/api_contract.md` - documents central indexes and visible conflict behavior.
- `docs/local_dev_guide.md` - documents startup index bootstrap and dirty-data failure mode.
- `api_test.http` - adds duplicate branch-code sample.
- `CHAT_CONTEXT/backend_skills/worklog.md` - records implementation summary and next action.

## Key decisions

- Centralized index creation in `pkg/database.EnsureIndexes` and removed repository-constructor
  index creation to avoid split ownership.
- Kept `members.ccid` index name as `ccid_1` for compatibility with local DBs that previously
  created the default-named unique index.
- Added query indexes for current list/filter paths and unique indexes for branch code, refund audit,
  duplicate session check-in, and duplicate makeup reuse.
- Added `refresh_tokens.expires_at` TTL index with `ExpireAfterSeconds(0)` while keeping auth service
  expiry checks as the source of truth.
- Mapped duplicate-key errors to existing domain conflicts where possible and added
  `ErrBranchCodeAlreadyExists` for branch-code collisions.
- Left full MongoDB transaction work out of this cycle because the local setup is not guaranteed to
  run as a replica set.

## Implementation notes

### Database

- Startup now calls `EnsureIndexes` after selecting `gym_management` and before repository
  construction.
- `EnsureIndexes` creates indexes for members, branches, subscriptions, attendances, sessions,
  refunds, employees, and refresh tokens.
- Startup fails fast if dirty existing data prevents unique index creation.

### Repository

- Member, branch, refund, and attendance create/update paths now map Mongo duplicate key errors to
  `repository.ErrDuplicate`.
- Employee duplicate mapping remains in place.

### Service

- Member duplicate create races return `ErrMemberCCIDAlreadyExists`.
- Branch duplicate create/update returns `ErrBranchCodeAlreadyExists`.
- Refund duplicate audit insert returns `ErrRefundAlreadyExists`.
- Attendance duplicate insert returns `ErrMakeupAlreadyUsed` for makeup and `ErrSessionCheckInClosed`
  when `session_id` exists.

### Handler

- Branch duplicate code now returns the shared error envelope with `409`/`CONFLICT`.
- Attendance/session duplicate conflicts continue to return `409`.

### Docs/API samples

- API docs now mention central startup index bootstrap, unique indexes, partial unique attendance
  indexes, and TTL cleanup behavior.
- Local dev guide now documents duplicate-data startup failure.
- REST samples include duplicate branch-code check.

## Commands run

```bash
gofmt -w cmd/server/main.go pkg/database/indexes.go internal/repository/member_repo.go internal/repository/branch_repo.go internal/repository/employee_repo.go internal/repository/refresh_token_repo.go internal/repository/refund_repo.go internal/repository/attendance_repo.go internal/service/member_service.go internal/service/branch_service.go internal/service/subscription_service.go internal/service/attendance_service.go internal/handlers/branch_handler.go internal/handlers/attendance_handler.go
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
git diff --check
```

Results:

- `go build ./...` - pass; Go printed the existing read-only module stat-cache warning but exited
  `0`.
- `go test ./...` - pass.
- `git diff --check` - pass.

## Known limitations

- Manual API/DB verification was not run in this implementation phase.
- Existing dirty local data can block startup index creation; cleanup is intentionally manual.
- Refund still updates subscription status before writing the refund audit row; unexpected audit
  insert failures after the status update remain a transaction follow-up.
- Attendance insert, subscription remaining-session decrement, and member attended counter update
  remain separate writes. Partial unique indexes now cover duplicate session check-in and makeup
  reuse, but they do not make all attendance side effects transactional.
- Service-level reference hardening for branch manager and session trainer/branch remains follow-up
  work.

## Handoff to review

- Review `pkg/database/indexes.go` index names/options, especially partial unique attendance indexes
  and TTL behavior.
- Check startup ordering in `cmd/server/main.go`.
- Check duplicate-key mapping from repository to service to handler for member, branch, refund, and
  attendance paths.
- Decide whether branch/session reference hardening should stay deferred or be added before test.
