# Cycle 07 - Indexes & Data Integrity

## Status

- Status: planned
- Feature: indexes and data-integrity hardening
- Planned at: 2026-05-28
- Priority: medium
- Depends on: current feature set through cycle 06
- Next phase: `$gym-implement`

## Goal

Make MongoDB enforce the data rules that should not rely only on read-before-write checks, and add
the query indexes needed by the current API surfaces.

This cycle is not an endpoint expansion. Existing HTTP routes and success response shapes should
stay unchanged unless an integrity failure needs clearer conflict mapping.

## Current Baseline

Already implemented:

- `members.ccid` unique index is created in `NewMemberRepository`.
- `branches.location` 2dsphere index is created in `NewBranchRepository`.
- `employees.normalized_email` and `employees.employee_id` unique sparse indexes are created in
  `NewEmployeeRepository`.
- `refresh_tokens.token_hash` unique index is created in `NewRefreshTokenRepository`.
- `sessionRepo.ReserveEnrollment` uses an atomic `FindOneAndUpdate` with capacity and duplicate
  enrollment guards.
- `subscriptionRepo.RefundSubscription` atomically changes active subscriptions with remaining
  sessions to `refunded`.

Gaps:

- Index creation is split across repository constructors and is incomplete for subscriptions,
  attendances, sessions, refunds, and several auth/list filters.
- `branches.branch_code` is not unique at the DB layer.
- `refunds.subscription_id` is not unique at the DB layer, so future direct writes or race paths can
  create duplicate refund audit rows.
- Duplicate key mapping is inconsistent: employee maps duplicate key to `repository.ErrDuplicate`,
  but member/branch/refund create paths currently return raw Mongo duplicate errors.
- Attendance duplicate prevention and remaining-session decrement still have race windows because
  attendance insert and subscription/member counter updates are separate operations.
- Refresh tokens have no TTL index on `expires_at`.

## API Contract

No new endpoints.

Expected visible behavior:

- Existing duplicate member CCID, branch code, employee email/employee ID, refresh token hash, and
  refund audit conflicts should map to the shared error contract with `409` and
  `error.code = CONFLICT`.
- Startup should fail fast with a clear log if index creation fails.
- Success response shapes remain unchanged.

Docs to update if implemented:

- `docs/api_contract.md`: add/adjust notes about DB-enforced uniqueness where behavior is visible.
- `docs/local_dev_guide.md`: mention startup index bootstrap and dirty-data duplicate-index failure.
- `api_test.http`: add representative duplicate branch/refund conflict samples if useful.

## Business Rules To Enforce

### Unique Identity Fields

- `members.ccid` remains unique.
- `branches.branch_code` becomes unique.
- `employees.normalized_email` and `employees.employee_id` remain unique sparse indexes.
- `refresh_tokens.token_hash` remains unique.
- Duplicate-key errors on API create/update flows must become service-level conflicts, not raw
  storage errors.

### Refund Audit Integrity

- `refunds.subscription_id` must be unique.
- `RefundRepository.Create` should map Mongo duplicate-key errors to `repository.ErrDuplicate`.
- `SubscriptionService.RefundSubscription` should map refund duplicate insert to
  `ErrRefundAlreadyExists`.
- Keep the current subscription atomic status update guard.
- Do not add a Mongo transaction in this cycle unless implementation discovers a simple local
  pattern. Record the remaining risk: subscription update can still succeed before audit insert
  fails if there is an unexpected write/storage failure after the status update.

### Query Performance Indexes

Add indexes matching current query/filter paths:

- `subscriptions.member_id`
- `subscriptions.status`
- compound `subscriptions.member_id + subscriptions.status`
- optional supporting indexes on `subscriptions.course_id` and `subscriptions.home_branch_id`
- `attendances.sub_id + date desc` for history and service scans
- `attendances.session_id` for session check-in duplicate scans when a session exists
- partial unique index for session check-in duplicate prevention:
  `attendances.session_id + sub_id` where `session_id` exists
- partial unique index for makeup reuse:
  `attendances.sub_id + is_makeup_for + status` where `status = "makeup"` and `is_makeup_for`
  exists
- `sessions.branch_id + scheduled_at`
- `sessions.course_level + scheduled_at`
- `sessions.tags`
- `employees.role + status + created_at desc`
- `employees.branch_id + status`
- `refresh_tokens.employee_id + revoked_at`
- TTL index on `refresh_tokens.expires_at`

### Attendance Race Windows

Minimum for this cycle:

- Add partial unique indexes that reject duplicate session check-ins and duplicate makeup reuse.
- Map duplicate-key errors from attendance insert to existing conflict errors where the service can
  identify the attempted operation:
  - duplicate session check-in -> `ErrSessionCheckInClosed`
  - duplicate makeup reuse -> `ErrMakeupAlreadyUsed`

Out of scope unless explicitly expanded:

- Full transaction for attendance insert + subscription remaining-session decrement + member
  attended counter increment.
- Atomic decrement of `remaining_sessions` based on status and session limit in one DB operation.

### Reference Integrity Hardening

Plan as service-level checks because MongoDB does not enforce foreign keys:

- Session creation should verify `branch_id` exists.
- Session creation should verify `trainer_id` exists, is active, and has trainer role.
- Branch manager assignment should verify `manager_id` exists when provided.
- Employee branch IDs are already validated by employee service; keep that behavior.

These checks can be implemented in this cycle if scope allows after index work. If not, keep them as
explicit follow-up notes in the implementation handoff.

## Data / Index Design

Preferred implementation:

- Add `pkg/database/indexes.go` with `EnsureIndexes(ctx context.Context, db *mongo.Database) error`.
- Call `database.EnsureIndexes` from `cmd/server/main.go` after selecting `gym_management` and
  before repository construction.
- Keep index definitions centralized with stable names.
- Either remove repository-constructor index creation or keep it only when definitions are identical
  and idempotent. Avoid creating the same key with different names/options.

Suggested index names:

| Collection | Index |
|---|---|
| `members` | `ccid_unique` |
| `branches` | `branch_code_unique`, `location_2dsphere` |
| `subscriptions` | `member_id_idx`, `status_idx`, `member_status_idx`, `course_id_idx`, `home_branch_id_idx` |
| `attendances` | `sub_id_date_desc_idx`, `session_id_idx`, `session_sub_unique`, `makeup_sub_ref_unique` |
| `sessions` | `branch_scheduled_at_idx`, `level_scheduled_at_idx`, `tags_idx` |
| `refunds` | `subscription_id_unique`, `member_id_idx` |
| `employees` | `normalized_email_unique`, `employee_id_unique`, `role_status_created_idx`, `branch_status_idx` |
| `refresh_tokens` | `token_hash_unique`, `employee_revoked_idx`, `expires_at_ttl` |

TTL detail:

- `refresh_tokens.expires_at` should use `ExpireAfterSeconds(0)`.
- TTL cleanup is eventual, so auth logic must continue checking `expires_at` directly.

## Layer Plan

### Database Bootstrap

- Create a central index bootstrap function.
- Use a bounded context, around 10 seconds, for startup index creation.
- Return errors to `main` and log fatal with collection/index context.
- Keep startup idempotent across repeated local runs.

### Repository

- Add duplicate-key normalization in create/update paths that can hit unique indexes:
  - member create
  - branch create/update
  - refund create
  - attendance create where partial unique indexes can fire
- Keep `repository.ErrDuplicate` as the storage-agnostic duplicate signal.
- Add focused repository tests if a Mongo test harness exists; otherwise cover behavior at service
  level and manual API/DB checks.

### Service

- Map `repository.ErrDuplicate` to existing service conflicts:
  - member -> `ErrMemberCCIDAlreadyExists`
  - branch -> likely introduce `ErrBranchConflict` or reuse `ErrInvalidBranchInput` only if the API
    should stay `400`; prefer new conflict error for duplicate branch code.
  - refund -> `ErrRefundAlreadyExists`
  - attendance makeup/session duplicates -> existing makeup/session conflict errors.
- Add reference checks for session branch/trainer and branch manager only if repositories/interfaces
  can support it without large refactors.

### Handler

- Preserve current error response shape.
- Map any new service conflict sentinel to `RespondConflict`.
- Do not expose index names or raw Mongo error text.

### Docs / Context

- Update API and local-dev docs if visible behavior changes.
- Update implementation/review/test notes in later phases.
- Update `CHAT_CONTEXT/README.md` and worklog at completion.

## Verification Plan

Automated:

```bash
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
git diff --check
```

Manual local API / DB checks:

- Start app against local MongoDB and confirm startup succeeds with index bootstrap.
- Verify indexes exist with `mongosh` / Compass for key collections.
- Duplicate member CCID returns `409 CONFLICT`.
- Duplicate branch code returns `409 CONFLICT`.
- Duplicate refund for a subscription returns `409 CONFLICT`, with only one refund audit row.
- Duplicate session check-in returns existing conflict response.
- Duplicate makeup for the same reported-missed reference returns existing conflict response.
- Existing nearby branch search still works with `location_2dsphere`.
- Refresh-token login/refresh/logout still works; expired-token TTL is documented as eventual.

Dirty-data safety check:

- If startup index creation fails because local data already contains duplicates, do not silently
  drop or rewrite data. Record the duplicate collection/index and require cleanup.

## Risks And Deferrals

- Unique indexes can fail on existing dirty local data.
- TTL deletion is not immediate and must not be treated as auth enforcement.
- MongoDB transactions may require replica set configuration; this project currently runs a simple
  local MongoDB container.
- Attendance remaining-session updates and member counters remain multi-write operations unless a
  future cycle adds transactions or atomic counter claims.
- Centralizing indexes touches startup and repository constructors; keep the change narrow and avoid
  unrelated repository refactors.

## Implementation Handoff

Start `$gym-implement` from this file:

`CHAT_CONTEXT/backend_skills/plans/07_indexes_data_integrity.md`
