# Backend Plan - 10 Dashboard Reports

Status: Planned

Created: 2026-06-02

## Goal

Add backend dashboard/report aggregate endpoints so FE11 can replace static dashboard sample data
with live operational data.

This cycle should expose a small, read-only dashboard API surface for management reporting without
changing existing member, subscription, attendance, session, refund, or employee mutation behavior.

## Current Baseline

- The backend currently has no `/api/v1/dashboard/*` or `/api/v1/reports/*` routes.
- Existing implemented collections can support the first dashboard pass:
  - `members`: active/registered counts, recent registrations.
  - `subscriptions`: paid subscription revenue, active subscriptions, plan/course distribution.
  - `refunds`: processed refund amounts for net revenue.
  - `attendances`: today check-in counts.
  - `sessions`: today's and this week's class counts, today's session list.
- Existing indexes cover some query shapes but not all dashboard date-range reads:
  - `sessions` has `branch_id + scheduled_at` and `course_level + scheduled_at`.
  - `attendances` has `sub_id + date`, but not date-only.
  - `subscriptions` has `status`, `course_id`, `home_branch_id`, but not payment-date range.
  - `refunds` has subscription/member indexes, but not processed-date range.

## API Contract

All endpoints are protected and require `admin` or `manager`.

### Summary

`GET /api/v1/dashboard/summary?branch_id=&from=&to=`

Query params:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `branch_id` | ObjectID string | no | Limits branch-scoped subscription/session/attendance metrics when supplied. |
| `from` | RFC3339 string | no | Range start. Defaults to the start of the current month for revenue/member deltas. |
| `to` | RFC3339 string | no | Range end. Defaults to now. Must be after `from`. |

Response `200`:

```json
{
  "message": "dashboard summary fetched successfully",
  "data": {
    "active_members": 128,
    "active_members_delta": 12,
    "monthly_revenue": 142000000,
    "monthly_revenue_delta": 12000000,
    "today_checkins": 34,
    "today_checkins_delta": -4,
    "classes_this_week": 18,
    "classes_this_week_delta": 3,
    "range": {
      "from": "2026-06-01T00:00:00Z",
      "to": "2026-06-02T08:00:00Z"
    }
  }
}
```

Revenue is net revenue: paid subscription totals whose `payment_date` falls in range, minus processed
refund amounts whose `processed_at` falls in range.

### Revenue

`GET /api/v1/dashboard/revenue?branch_id=&from=&to=&bucket=day`

Query params:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `branch_id` | ObjectID string | no | Limits subscription payments by `home_branch_id`. |
| `from` | RFC3339 string | no | Defaults to six days before `to` at day start. |
| `to` | RFC3339 string | no | Defaults to now. |
| `bucket` | string | no | Only `day` is supported in this cycle. |

Response `200`:

```json
{
  "message": "dashboard revenue fetched successfully",
  "data": {
    "bucket": "day",
    "from": "2026-05-27T00:00:00Z",
    "to": "2026-06-02T08:00:00Z",
    "items": [
      {
        "label": "2026-06-01",
        "gross_amount": 15000000,
        "refund_amount": 1000000,
        "net_amount": 14000000
      }
    ]
  }
}
```

### Plan distribution

`GET /api/v1/dashboard/plans?branch_id=&from=&to=`

Response `200`:

```json
{
  "message": "dashboard plan distribution fetched successfully",
  "data": {
    "items": [
      {
        "course_id": "69f20a180c4cd4cdf57684fe",
        "label": "Advanced Strength",
        "count": 42
      }
    ]
  }
}
```

Distribution counts subscriptions in `active`, `suspended`, and `expired` states by `course_id`.
Pending subscriptions are excluded because they are not paid/registered business yet. Refunded
subscriptions are excluded because they no longer represent active plan mix.

### Recent members

`GET /api/v1/dashboard/members/recent?limit=5`

Response `200`:

```json
{
  "message": "dashboard recent members fetched successfully",
  "data": {
    "items": [
      {
        "id": "69f20c000c4cd4cdf5768500",
        "full_name": "Nguyen Minh Khoa",
        "phone": "0912345678",
        "level": "advanced",
        "is_registered": true,
        "created_at": "2026-06-01T08:00:00Z"
      }
    ]
  }
}
```

`limit` defaults to `5`, max `20`.

### Today's sessions

`GET /api/v1/dashboard/sessions/today?branch_id=&date=`

Query params:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `branch_id` | ObjectID string | no | Limits sessions by branch. |
| `date` | RFC3339 string | no | Defaults to today. Uses the date portion to build a day range. |

Response `200`:

```json
{
  "message": "dashboard today sessions fetched successfully",
  "data": {
    "items": [
      {
        "id": "69f20c000c4cd4cdf5768500",
        "branch_id": "69f20a180c4cd4cdf57684fe",
        "trainer_id": "69f20c000c4cd4cdf5768501",
        "course_level": "advanced",
        "scheduled_at": "2026-06-02T09:00:00Z",
        "duration_min": 60,
        "capacity": 15,
        "enrolled_count": 12,
        "tags": ["strength"]
      }
    ]
  }
}
```

Status codes:

- `200`: aggregate fetched; empty arrays/counts are valid.
- `400`: invalid ObjectID, invalid date, invalid range, invalid bucket, or invalid limit.
- `401`: missing/malformed/expired token.
- `403`: authenticated role is not `admin` or `manager`.
- `500`: storage/internal failure.

## Business Rules

- Dashboard/report endpoints are read-only and must not mutate domain state.
- Revenue is server-computed only; clients cannot provide money totals.
- Net revenue = subscription payments in range minus processed refunds in range.
- Subscription payments count only documents with a non-null `payment_date`.
- Branch filter applies to:
  - subscription/revenue/plan metrics through `home_branch_id`
  - sessions through `branch_id`
  - attendance through `branch_id`
  - recent members do not filter by branch in this cycle because members have no branch field
- Date filters are inclusive start and exclusive end: `from <= value < to`.
- If `from`/`to` are absent, handlers build stable defaults:
  - summary: current month start through now
  - revenue: last 7 days through now
  - sessions today: selected/current day start through next day start
- Deltas compare the selected range with the immediately previous range of the same duration.
- Use UTC for truncating default/report buckets in this cycle.

## Data And Index Changes

No schema changes.

Add indexes through `pkg/database.EnsureIndexes`:

- `subscriptions`: `payment_date_idx`
- `subscriptions`: `home_branch_payment_date_idx`
- `members`: `created_at_desc_idx`
- `attendances`: `date_idx`
- `attendances`: `branch_date_idx`
- `sessions`: `scheduled_at_idx`
- `refunds`: `processed_at_idx`

These indexes support date-range aggregation and recent-member reads without changing existing unique
constraints.

## Layer Plan

### Models

Add `internal/models/dashboard.go` with response DTOs:

- `DashboardRange`
- `DashboardSummary`
- `DashboardRevenueResponse`
- `DashboardRevenueBucket`
- `DashboardPlanDistributionResponse`
- `DashboardPlanDistributionItem`
- `DashboardRecentMembersResponse`
- `DashboardRecentMember`
- `DashboardTodaySessionsResponse`

### Repository

Add `internal/repository/dashboard_repo.go`:

- Owns Mongo collection access for aggregate/read-only dashboard queries.
- Uses aggregation or focused `CountDocuments`/`Find` calls.
- Accepts typed filter structs, not raw HTTP query strings.
- Returns empty arrays/counts instead of nil where useful for frontend rendering.

Planned methods:

```go
GetSummary(ctx, filter DashboardRangeFilter) (*models.DashboardSummary, error)
GetRevenue(ctx, filter DashboardRevenueFilter) (*models.DashboardRevenueResponse, error)
GetPlanDistribution(ctx, filter DashboardRangeFilter) (*models.DashboardPlanDistributionResponse, error)
GetRecentMembers(ctx, limit int) (*models.DashboardRecentMembersResponse, error)
GetTodaySessions(ctx, filter DashboardTodaySessionsFilter) (*models.DashboardTodaySessionsResponse, error)
```

### Service

Add `internal/service/dashboard_service.go`:

- Validates filter defaults and range semantics.
- Enforces supported bucket values.
- Coordinates repository calls and delta calculations.
- Keeps money/count semantics out of handlers.

### Handler

Add `internal/handlers/dashboard_handler.go`:

- Parses query params.
- Maps invalid IDs/dates/ranges/buckets/limits to existing error helpers.
- Delegates to service and returns existing `message`/`data` success shape.

### Route Wiring

Update `internal/app/router.go`:

- Instantiate dashboard repository/service/handler.
- Add protected `admin`/`manager` routes before catch-all-like dynamic routes are not relevant here,
  but keep dashboard routes grouped near other management reads.

Routes:

```go
dashboardRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager))
dashboardRoutes.GET("/dashboard/summary", h.Dashboard.Summary)
dashboardRoutes.GET("/dashboard/revenue", h.Dashboard.Revenue)
dashboardRoutes.GET("/dashboard/plans", h.Dashboard.PlanDistribution)
dashboardRoutes.GET("/dashboard/members/recent", h.Dashboard.RecentMembers)
dashboardRoutes.GET("/dashboard/sessions/today", h.Dashboard.TodaySessions)
```

## Docs And Test Plan

Docs/API samples:

- Update `docs/api_contract.md` with dashboard endpoint tables and response examples.
- Add `api_test.http` samples for all dashboard endpoints.
- Update `CHAT_CONTEXT/README.md` and `backend_skills/worklog.md` when complete.

Automated verification:

```sh
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
```

Focused tests:

- Unit tests for dashboard service default range/bucket/limit validation if behavior grows complex.
- Integration tests with seeded member/subscription/payment/refund/attendance/session data:
  - admin can fetch dashboard endpoints
  - receptionist is forbidden
  - invalid branch/date/range/limit/bucket returns `400`
  - empty DB returns zero counts and empty arrays
  - revenue net amount subtracts processed refunds

Manual API checks when local backend/Mongo/credentials are available:

- Login as admin.
- Call all five dashboard endpoints.
- Call one endpoint with invalid date and invalid branch ID.
- Confirm receptionist or trainer token gets `403`.

## Risks And Boundaries

- Revenue semantics are intentionally narrow: paid subscription totals minus processed refund totals.
  They do not model cash drawer adjustments, online payment provider fees, or refund approval states.
- Recent members cannot be branch-scoped without adding branch ownership to members.
- Plan labels may need course names; current course model should be joined/loaded for labels, but
  `course_id` remains the stable identifier.
- Aggregation reads are eventually consistent with Mongo writes and do not require transactions.
- Do not add report export, CSV/PDF, scheduled reports, or payment-history screens in this cycle.
- Do not expose raw internal errors or database pipeline details in response messages.

## Next Action

Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/10_dashboard_reports.md`.
