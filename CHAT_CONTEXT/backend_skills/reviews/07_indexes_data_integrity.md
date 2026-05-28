# Code Review - Indexes and Data Integrity

## Status

- Status: reviewed
- Feature: indexes and data-integrity hardening
- Plan file: `CHAT_CONTEXT/backend_skills/plans/07_indexes_data_integrity.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/07_indexes_data_integrity.md`
- Reviewed at: 2026-05-28

## Review summary

- Result: pass; no blocking findings.
- Build status: pass.
- Test status: pass.

## Checklist

- [x] Code compiles.
- [x] Handler only handles HTTP parse/response.
- [x] Service owns business rules.
- [x] Repository only handles DB.
- [x] Model tags match API/DB contract.
- [x] Errors map to correct HTTP status.
- [x] Atomic updates used where needed.
- [x] No client-controlled money/status/role introduced.
- [x] Routes have correct order.
- [x] Docs/API samples match behavior.

## Passed

- `pkg/database/indexes.go` centralizes index definitions for all current collections and includes
  the planned unique, query, partial unique, and TTL indexes.
- `cmd/server/main.go` runs `database.EnsureIndexes` after selecting the database and before
  repository construction, so startup fails before serving requests if index creation fails.
- Repository constructors now bind collections without hidden index side effects.
- Member, branch, refund, and attendance repositories normalize Mongo duplicate-key errors to
  `repository.ErrDuplicate`.
- Service mapping keeps business meaning out of handlers:
  - duplicate member CCID -> `ErrMemberCCIDAlreadyExists`
  - duplicate branch code -> `ErrBranchCodeAlreadyExists`
  - duplicate refund audit -> `ErrRefundAlreadyExists`
  - duplicate makeup/session attendance -> existing conflict errors
- Handler mapping preserves the shared error envelope and returns `409` for new duplicate conflicts.
- The implementation matches the plan's deferral of Mongo transactions and reference-hardening
  checks.
- `docs/api_contract.md`, `docs/local_dev_guide.md`, and `api_test.http` reflect the visible index
  and duplicate-conflict behavior.

## Issues found

- None blocking.

## Fixes applied during review

- None.

## Remaining risks

- Manual API/DB verification has not been run yet; startup index creation should be checked against
  the local MongoDB dataset because dirty duplicate data can block startup.
- There is no automated integration test that asserts actual MongoDB index creation or duplicate-key
  behavior; current automated coverage is build/unit level.
- Refund still has a known multi-write limitation: subscription status can change before refund
  audit insert if an unexpected storage failure happens after the status update.
- Attendance still performs insert, subscription remaining-session update, and member counter update
  as separate writes. Partial unique indexes cover duplicate session check-in/makeup reuse but not
  full transactional consistency.
- Branch manager and session branch/trainer reference hardening remains deferred.

## Commands run

- `env GOCACHE=/tmp/gocache go build ./...` - pass; Go printed the existing read-only module
  stat-cache warning but exited `0`.
- `env GOCACHE=/tmp/gocache go test ./...` - pass.
- `git diff --check` - pass.
- `rg -n "EnsureIndexes|ErrDuplicate|ErrBranchCodeAlreadyExists|ErrRefundAlreadyExists|ErrMakeupAlreadyUsed|ErrSessionCheckInClosed|MongoDB indexes|branch_code_unique|expires_at_ttl" . --glob '*.go' --glob '*.md' --glob '*.http'`

## Handoff to test

- Use `$gym-test` with this review note.
- Test startup against local MongoDB and verify key indexes exist.
- Manually verify duplicate member CCID, duplicate branch code, duplicate refund, duplicate session
  check-in, and duplicate makeup reuse return `409` with `error.code = CONFLICT`.
- Verify existing nearby branch search, auth login/refresh/logout, employee list filters, session
  list filters, and subscription list by member still work after index bootstrap.
