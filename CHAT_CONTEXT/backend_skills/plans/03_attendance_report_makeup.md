# Cycle 03 — Attendance Report/Makeup Endpoints

## Status

- Status: completed on 2026-05-21
- Priority: medium-high
- Depends on: existing attendance service rules already present in `AttendanceService.CheckIn`
- Endpoints:
  - `POST /api/v1/attendance/report`
  - `POST /api/v1/attendance/makeup`

## Goal

Expose dedicated attendance report/makeup APIs while keeping business rules centralized in `internal/service/attendance_service.go`.

Current service already accepts statuses:
- `attended`
- `absent`
- `reported_missed`
- `makeup`

Dedicated endpoints should avoid exposing free-form `status` from client and should internally set correct status.

## Current source findings

Read files:
- `internal/models/attendance.go`
- `internal/repository/attendance_repo.go`
- `internal/service/attendance_service.go`
- `internal/handlers/attendance_handler.go`
- `cmd/server/main.go`
- `docs/api_contract.md`

Findings:
- Model has no `member_id`, no `reason`, no `reported_missed_ref_id`.
- `Attendance` fields:
  - `sub_id`
  - `branch_id`
  - optional `session_id`
  - `date`
  - `status`
  - `is_makeup_for`
- Repository supports:
  - `Create`
  - `ListBySubscriptionID`
- Service method `CheckIn` already contains report/makeup business rules:
  - `reported_missed` 30-day sliding window via `validateReportedMissedWindow`
  - `makeup` validates `IsMakeupFor` date within 7 days via `validateMakeupRequest`
  - `makeup` rejects reused `IsMakeupFor` date
  - `attended` and `makeup` enforce weekly session limit
- Handler currently exposes only:
  - `POST /api/v1/attendance/checkin`
  - `GET /api/v1/subscriptions/:id/attendance`
- Route wiring does not register:
  - `POST /api/v1/attendance/report`
  - `POST /api/v1/attendance/makeup`
- `docs/api_contract.md` still marks report/makeup endpoints as `Planned`.

## API plan

### Report missed

```http
POST /api/v1/attendance/report
```

Request:

```json
{
  "subscription_id": "ObjectID",
  "branch_id": "ObjectID",
  "date": "2026-05-19T10:00:00Z"
}
```

Notes:
- Do not accept `status`; handler sets `status = "reported_missed"`.
- Do not accept `member_id`; service validates subscription existence/status. Current model does not store member on attendance.
- Do not accept `reason` in this cycle unless model change is explicitly added.

Success response `201`:

```json
{
  "message": "attendance report recorded successfully",
  "data": {
    "id": "ObjectID",
    "sub_id": "ObjectID",
    "branch_id": "ObjectID",
    "date": "2026-05-19T10:00:00Z",
    "status": "reported_missed",
    "is_makeup_for": null
  }
}
```

Status codes:
- `201`: report created.
- `400`: invalid JSON, invalid `subscription_id`, invalid `branch_id`, invalid `date`.
- `404`: subscription not found.
- `409`: subscription not active/expired, report limit reached in 30-day window.
- `500`: internal error.

### Makeup

```http
POST /api/v1/attendance/makeup
```

Request:

```json
{
  "subscription_id": "ObjectID",
  "branch_id": "ObjectID",
  "date": "2026-05-22T10:00:00Z",
  "is_makeup_for": "2026-05-19T10:00:00Z"
}
```

Notes:
- Do not accept `status`; handler sets `status = "makeup"`.
- Current service references reported missed record by exact missed date (`is_makeup_for`), not by attendance ID.
- `reported_missed_ref_id` should not be introduced in this cycle unless model/repository/service are expanded to support ID lookup and storage.

Success response `201`:

```json
{
  "message": "attendance makeup recorded successfully",
  "data": {
    "id": "ObjectID",
    "sub_id": "ObjectID",
    "branch_id": "ObjectID",
    "date": "2026-05-22T10:00:00Z",
    "status": "makeup",
    "is_makeup_for": "2026-05-19T10:00:00Z"
  }
}
```

Status codes:
- `201`: makeup created.
- `400`: invalid JSON, invalid `subscription_id`, invalid `branch_id`, invalid `date`, invalid `is_makeup_for`.
- `404`: subscription not found.
- `409`: subscription not active/expired, weekly session limit reached, invalid makeup reference, makeup reference not found, makeup already used, no remaining sessions.
- `500`: internal error.

## Business rules

Report:
- Subscription must exist.
- Subscription must be `active`.
- Subscription must not be expired at report date.
- Report date defaults to current server time if omitted.
- Report date must be within existing 30-day sliding window rule: one `reported_missed` per 30 days for same subscription.
- No remaining-session decrement for `reported_missed`.

Makeup:
- Subscription must exist.
- Subscription must be `active`.
- Subscription must not be expired at makeup date.
- Makeup date defaults to current server time if omitted.
- `is_makeup_for` is required.
- `is_makeup_for` must match existing `reported_missed` attendance date for same subscription.
- Makeup date must be after or equal to missed date.
- Makeup date must be within 7 days after missed date.
- Same `is_makeup_for` date cannot be reused.
- Must respect `session_per_week`.
- Must consume one remaining session and increment member attended count, same as existing `CheckIn` behavior.

## Data changes

### Model

No required model change for minimal scope.

Current `Attendance` supports needed fields through:
- `status = "reported_missed"`
- `status = "makeup"`
- `is_makeup_for` as missed-date reference for makeup

Optional future model changes if product needs richer report/makeup audit:
- add `member_id` denormalized field
- add `reason`
- add `reported_missed_ref_id`
- add unique/index support to prevent duplicate makeup by reference at DB level

### Collection

Use existing `attendances` collection.

### Index

No index required for endpoint exposure.

Recommended follow-up for Cycle 06:
- `{ sub_id: 1, date: -1 }` for history/window checks.
- Optional partial unique index for makeup reuse if moving from date reference to stable report ID.

## Repository plan

Minimal scope:
- No repository interface change required.
- Continue using `ListBySubscriptionID` for service validation.

Optional stronger scope:
- Add `GetByID(ctx, id string)` only if API uses `reported_missed_ref_id`.
- Add targeted query methods for report/makeup validation to avoid scanning full subscription history:
  - `CountReportedMissedInWindow`
  - `FindReportedMissedByDate`
  - `ExistsMakeupForDate`

## Service plan

Minimal scope:
- Keep business rules inside existing `CheckIn`.
- No new service method required if handlers call `CheckIn` with fixed status.

Cleaner optional API:
- Add `ReportMissed(ctx, attendance *models.Attendance) error`
  - force status to `reported_missed`
  - delegate to shared validation/create path
- Add `Makeup(ctx, attendance *models.Attendance) error`
  - force status to `makeup`
  - require `IsMakeupFor`
  - delegate to shared validation/create path

If adding methods, avoid duplicating rules. Refactor private helper only if needed.

## Handler plan

Add request DTOs:
- `attendanceReportRequest`
  - `subscription_id`
  - `branch_id`
  - optional `date`
- `attendanceMakeupRequest`
  - `subscription_id`
  - `branch_id`
  - optional `date`
  - required `is_makeup_for`

Add handler methods:
- `ReportMissed(c *gin.Context)`
- `Makeup(c *gin.Context)`

Handler responsibilities:
- parse JSON
- parse ObjectIDs
- parse RFC3339 dates
- set server-controlled status
- call service
- map errors to HTTP response
- return `201` with attendance payload

Do not place 30-day, 7-day, active subscription, or session/week logic in handler.

## Route plan

Add in `cmd/server/main.go`:

```go
api.POST("/attendance/report", attendanceHandler.ReportMissed)
api.POST("/attendance/makeup", attendanceHandler.Makeup)
```

Place near existing attendance routes:

```go
api.POST("/attendance/checkin", attendanceHandler.CheckIn)
api.POST("/attendance/report", attendanceHandler.ReportMissed)
api.POST("/attendance/makeup", attendanceHandler.Makeup)
api.GET("/subscriptions/:id/attendance", attendanceHandler.ListBySubscription)
```

No route conflict expected.

## Error mapping

Map service errors:
- `ErrInvalidAttendanceInput` -> `400`
- `ErrSubscriptionNotFound` -> `404`
- `ErrAttendanceCheckInNotAllowed` -> `409`
- `ErrSubscriptionExpired` -> `409`
- `ErrNoRemainingSessions` -> `409`
- `ErrWeeklySessionLimitReached` -> `409`
- `ErrReportedMissedLimitReached` -> `409`
- `ErrMakeupReferenceInvalid` -> `409`
- `ErrMakeupReferenceNotFound` -> `409`
- `ErrMakeupAlreadyUsed` -> `409`
- default -> `500`

Note:
- Current handler maps `ErrMakeupReferenceNotFound` to `409`. If API consistency later wants missing report as `404`, include in Cycle 05 validation/error consistency.

## Docs/test plan

Update after implementation:
- `docs/api_contract.md`
  - mark report/makeup as Implemented
  - add endpoint details
- `api_test.http`
  - sample `POST /attendance/report`
  - sample `POST /attendance/makeup`
- `CHAT_CONTEXT/README.md`
  - current state if feature completed
- `CHAT_CONTEXT/backend_skills/implementations/03_attendance_report_makeup.md`
  - implementation log during `/backend-implement`
- `CHAT_CONTEXT/backend_skills/worklog.md`
  - short status during completion phase

Verification during implement/test phases:
```bash
gofmt -w internal/handlers/attendance_handler.go cmd/server/main.go
go build ./...
go test ./...
```

Manual API flow:
1. Create member/course/branch/subscription.
2. Activate member/subscription flow.
3. `POST /api/v1/attendance/report` with active subscription.
4. Confirm second report in 30-day window returns `409`.
5. `POST /api/v1/attendance/makeup` with `is_makeup_for` equal report date.
6. Confirm repeated makeup returns `409`.
7. Confirm old `is_makeup_for` beyond 7 days returns `409`.

## Risks

- Current makeup reference uses exact `time.Time` equality. Client must send same RFC3339 instant as report date. This is fragile.
- No DB-level uniqueness prevents concurrent duplicate makeup reuse. Service scan can race under double-submit.
- `CheckIn` creates attendance before checking remaining sessions for `attended/makeup`; if remaining sessions is already 0, a makeup record can be inserted before `ErrNoRemainingSessions`. This pre-existing risk should be handled in Cycle 06 or fixed during implementation if approved.
- Current plan does not add `reason` or `reported_missed_ref_id`; adding them expands scope to model/repository/service/API docs.
- Error consistency should be revisited in Cycle 05.

## Implementation recommendation

Implement minimal endpoint exposure now:
1. Add dedicated handler methods that set fixed statuses.
2. Reuse existing `AttendanceService.CheckIn`.
3. Register two routes.
4. Update API docs and samples during completion phase.
5. Leave DB indexes/concurrency hardening for Cycle 06 unless user approves scope expansion.
