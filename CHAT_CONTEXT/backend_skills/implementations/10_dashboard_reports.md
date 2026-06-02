# Implementation - 10 Dashboard Reports

## Status

- Status: implemented
- Feature: Dashboard/report aggregate APIs
- Plan file: `CHAT_CONTEXT/backend_skills/plans/10_dashboard_reports.md`
- Started at: 2026-06-02
- Finished at: 2026-06-02

## Scope implemented

- [x] Model changes
- [x] Repository changes
- [x] Service changes
- [x] Handler changes
- [x] Route changes
- [x] Docs/API sample changes

## Files changed

- `internal/models/dashboard.go` - added dashboard response DTOs.
- `internal/repository/dashboard_repo.go` - added Mongo read/aggregate queries for dashboard data.
- `internal/service/dashboard_service.go` - added dashboard defaults, validation, deltas, and bucket filling.
- `internal/handlers/dashboard_handler.go` - added query parsing, error mapping, and success responses.
- `internal/app/router.go` - wired dashboard repository/service/handler and admin/manager routes.
- `pkg/database/indexes.go` - added dashboard date-range/read indexes.
- `internal/integration/integration_test.go` - added dashboard integration coverage.
- `docs/api_contract.md` - documented dashboard endpoints, role matrix, responses, and status codes.
- `api_test.http` - added dashboard request samples.
- `CHAT_CONTEXT/backend_skills/plans/10_dashboard_reports.md` - added the feature plan.
- `CHAT_CONTEXT/backend_skills/worklog.md` - added roadmap/worklog entry.

## Key decisions

- Dashboard endpoints are admin/manager-only because they expose revenue/reporting data.
- Revenue is server-computed as paid subscription totals minus processed refund totals.
- Revenue buckets support `bucket=day` only in this cycle.
- Summary member deltas use newly registered members in the selected range versus the immediately
  previous range; total active members remains the current registered/not-suspended count.
- Branch filtering applies to subscription, revenue, attendance, and session metrics; recent members
  remain unscoped because member documents have no branch field.
- Date defaults are UTC-based.

## Implementation notes

### Models

- Added DTOs for summary, revenue buckets, plan distribution, recent members, and today's sessions.

### Repository

- Added a dashboard repository that owns read-only Mongo collection access.
- Used focused counts/finds and aggregation pipelines for revenue and plan distribution.
- Refunds are branch-filtered by looking up their source subscription.

### Service

- Added validation for branch ObjectID, date ranges, bucket, and recent-member limit.
- Added default ranges:
  - summary: current month start to now
  - revenue: last 7-day window to now
  - today sessions: selected/current day
- Added delta calculations for active member signups, revenue, today check-ins, and classes this week.

### Handler

- Added handlers for:
  - `GET /api/v1/dashboard/summary`
  - `GET /api/v1/dashboard/revenue`
  - `GET /api/v1/dashboard/plans`
  - `GET /api/v1/dashboard/members/recent`
  - `GET /api/v1/dashboard/sessions/today`
- Handler maps invalid IDs to `INVALID_ID`, invalid date/range to `INVALID_DATE`, invalid bucket/limit
  to `INVALID_INPUT`, and storage failures to `INTERNAL_ERROR`.

### Routes

- Added an admin/manager dashboard route group under the existing protected `/api/v1` router.

### Docs/API samples

- Updated the API contract collection table, role matrix, detailed endpoint docs, notes, and samples.
- Added REST samples for success, forbidden, and invalid-range dashboard checks.

## Commands run

```bash
gofmt -w internal/models/dashboard.go internal/repository/dashboard_repo.go internal/service/dashboard_service.go internal/handlers/dashboard_handler.go internal/app/router.go pkg/database/indexes.go internal/integration/integration_test.go
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
git diff --check
```

## Command results

- `go build ./...` passed. Go printed the existing read-only module stat-cache warning but exited `0`.
- `go test ./...` passed, including `internal/integration`.
- `git diff --check` passed.

## Known limitations

- No report export or `/api/v1/reports/*` endpoint was added.
- Revenue semantics do not include provider fees, cash adjustments, or future payment history beyond
  subscription/refund records.
- Recent members are not branch-scoped until member records store branch ownership.
- Bucket support is daily only.

## Handoff to review

- Check dashboard revenue/refund aggregation semantics, especially branch-filtered refunds through
  subscription lookup.
- Check role scope for dashboard endpoints.
- Check summary delta definitions and UTC default windows.
- Check route/docs/API sample alignment.
