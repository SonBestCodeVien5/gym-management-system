# Implementation — attendance report makeup

## Status

- Status: implemented
- Feature: attendance report makeup
- Plan file: `CHAT_CONTEXT/backend_skills/plans/03_attendance_report_makeup.md`
- Started at: 2026-05-21
- Finished at: 2026-05-21

## Scope implemented

- [x] Model scope reviewed; no model change required
- [x] Repository scope reviewed; no repository change required
- [x] Service changes
- [x] Handler changes
- [x] Route changes
- [x] Docs/API sample changes

## Files changed

- `internal/handlers/attendance_handler.go`
- `internal/service/attendance_service.go`
- `cmd/server/main.go`
- `docs/api_contract.md`
- `api_test.http`
- `CHAT_CONTEXT/backend_skills/implementations/03_attendance_report_makeup.md`

## Key decisions

- No model change needed for minimal scope.
- No repository change needed; existing `ListBySubscriptionID` still supports service validation.
- `AttendanceService.CheckIn` owns report/makeup rules and remaining-session guards.
- Dedicated endpoints set server-controlled statuses:
  - `POST /api/v1/attendance/report` sets `status = "reported_missed"`.
  - `POST /api/v1/attendance/makeup` sets `status = "makeup"`.
- Handlers do not accept free-form `status` for new dedicated report/makeup endpoints.

## Implementation notes

- Added report request DTO with:
  - `subscription_id`
  - `branch_id`
  - optional `date`
- Added makeup request DTO with:
  - `subscription_id`
  - `branch_id`
  - optional `date`
  - required `is_makeup_for`
- Added `AttendanceHandler.ReportMissed`.
- Added `AttendanceHandler.Makeup`.
- Added shared `handleAttendanceError` helper for new endpoints.
- Registered routes:
  - `POST /api/v1/attendance/report`
  - `POST /api/v1/attendance/makeup`
- Existing business rules remain centralized in `internal/service/attendance_service.go`:
  - active subscription required.
  - subscription must not be expired at attendance date.
  - `reported_missed` enforces 30-day window.
  - `makeup` requires existing reported missed reference by exact date.
  - `makeup` must be within 7 days and cannot reuse same reference.
  - `makeup` duplicate use is checked during validation and rechecked immediately before insert to narrow double-submit race scope.
  - `attended`/`makeup` reject `remaining_sessions <= 0` before inserting attendance.
  - `makeup` enforces weekly session limit and consumes remaining session.
- Updated `docs/api_contract.md` to mark report/makeup endpoints implemented.
- Added `api_test.http` samples for report and makeup endpoints.

## Commands run

- `gofmt -w internal/handlers/attendance_handler.go cmd/server/main.go`
- `go build ./...`
- `gofmt -w internal/service/attendance_service.go internal/handlers/attendance_handler.go cmd/server/main.go`
- `go build ./...`

## Known limitations

- Makeup still references missed report by exact `is_makeup_for` date because model has no `reported_missed_ref_id`.
- Makeup duplicate prevention remains service-level only; latest change narrows double-submit scope but does not provide DB-level uniqueness.

## Handoff to review

- Review route exposure and error mapping for:
  - `POST /api/v1/attendance/report`
  - `POST /api/v1/attendance/makeup`
- Confirm status remains server-controlled in dedicated endpoints.
- Confirm no business rule leaked from service into handler.
