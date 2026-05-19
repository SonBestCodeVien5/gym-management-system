# Cycle 06 — Indexes & Data Integrity

## Status

- Status: planned
- Priority: medium
- Depends on: final set of collections/features known

## Goal

Tạo index cần thiết và siết data integrity để MongoDB hỗ trợ business rules.

## Index plan

### Members

- unique `ccid`

Already noted:
- `members.ccid` unique index exists.

### Branches

- unique `branch_code`
- `location` 2dsphere

### Subscriptions

- `member_id`
- `course_id`
- `home_branch_id`
- `status`
- compound `{member_id, status}`

### Attendance

- `subscription_id`
- `member_id`
- `session_id`
- `reported_missed_ref_id`
- optional compound to prevent duplicate check-in per business day/session

### Sessions

- `branch_id`
- `start_time`
- `tags`
- compound `{branch_id, start_time}`

### Refunds

- unique `subscription_id`

### Employees/Auth

- unique `email` or `username`
- refresh token hash unique
- refresh token expiry TTL if using token collection

## Implementation plan

Create central bootstrap function or keep repo init constructors consistent.

Option A:
- `pkg/database/indexes.go`
- `EnsureIndexes(ctx, db) error`

Option B:
- each repository constructor creates its own indexes

Preferred:
- central `EnsureIndexes` for visibility.

Call from `cmd/server/main.go` after DB selection.

## Data integrity rules

- Unique fields return conflict, not raw Mongo error.
- Geo query index must exist before nearby endpoint.
- Refund double-submit guarded by unique index and atomic update.
- Attendance makeup reuse guarded by query/index if possible.
- Auth refresh tokens expire via TTL index if persisted.

## Docs/test plan

Update:
- `docs/local_dev_guide.md` if index bootstrapping matters.
- `CHAT_CONTEXT/README.md`
- `worklog.md`

Run:
```bash
go build ./...
go test ./...
```

Manual:
- start app, verify no index creation error.
- create duplicate ccid/branch_code/refund to confirm conflict behavior.

## Risks

- Creating unique indexes on dirty existing data can fail.
- Need migration/cleanup if duplicate data already exists.
- TTL index behavior not immediate.