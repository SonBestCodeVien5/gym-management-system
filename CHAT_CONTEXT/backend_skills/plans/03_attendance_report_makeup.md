# Cycle 03 — Attendance Report/Makeup Endpoints

## Status

- Status: planned
- Priority: medium-high
- Depends on: attendance service rules already existing
- Endpoints:
  - `POST /api/v1/attendance/report`
  - `POST /api/v1/attendance/makeup`

## Goal

Expose planned attendance report/makeup APIs if business rules already exist in service but no routes/handlers exist.

## Current context

README says:
- `reported_missed` enforces 30-day sliding window.
- `makeup` requires valid `reported_missed` reference within 7 days.
- Same report cannot be reused twice.

API contract still marks:
- `POST /api/v1/attendance/report` as Planned
- `POST /api/v1/attendance/makeup` as Planned

## API plan

### Report missed

```http
POST /api/v1/attendance/report
```

Request:
```json
{
  "subscription_id": "ObjectID",
  "member_id": "ObjectID",
  "branch_id": "ObjectID",
  "date": "2026-05-19T10:00:00Z",
  "reason": "sick"
}
```

Creates attendance record with type/status `reported_missed`.

### Makeup

```http
POST /api/v1/attendance/makeup
```

Request:
```json
{
  "subscription_id": "ObjectID",
  "member_id": "ObjectID",
  "branch_id": "ObjectID",
  "date": "2026-05-22T10:00:00Z",
  "reported_missed_ref_id": "ObjectID"
}
```

Creates attendance record with type/status `makeup`.

## Business rules

Report:
- Subscription must exist.
- Member must match subscription.
- Subscription must be active.
- Report date must be within 30-day sliding window.

Makeup:
- Must reference valid reported_missed record.
- Reference must belong to same subscription/member.
- Makeup date must be within 7 days from reported_missed.
- Report cannot be reused.
- Must respect session/week rule if existing service enforces it.

## Implementation plan

- Inspect `attendance_service.go` for existing methods.
- If methods exist, add handler methods + routes.
- If methods missing, add service methods around existing attendance creation rules.
- Add repo query to find attendance by ID if missing.
- Add route wiring in `cmd/server/main.go`.

## Docs/test plan

Update:
- `docs/api_contract.md`
- `api_test.http`
- `CHAT_CONTEXT/README.md`
- `worklog.md`

Run:
```bash
go build ./...
go test ./...
```

## Risks

- Existing service may use same `CheckIn` method with type field instead of separate report/makeup methods.
- Need avoid duplicate logic and keep rules centralized.