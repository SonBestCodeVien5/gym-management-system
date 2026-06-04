# FE Plan - 11 Live Dashboard APIs

Status: Completed

Created: 2026-06-02

## Goal

Replace the FE02 static dashboard sample data with live operational metrics once backend
dashboard/report endpoints exist.

FE11 should preserve the current staff-console dashboard shape where it still works: KPI cards,
revenue trend, membership mix, latest members, and today's sessions. The implementation should make
data provenance clear, remove sample labels only after live responses are wired, and leave the
existing auth/session shell unchanged.

The backend dashboard/report aggregate contract is now available through the Cycle 10 backend
implementation, so FE11 can use the live dashboard endpoints documented in `docs/api_contract.md`.

## Current Baseline

- Stack: React 18 + Vite 8.
- `/app/dashboard` is implemented and route status is `ready`.
- `DashboardHome` renders live staff identity from `AuthContext`.
- Dashboard metrics currently come from static `frontend/src/components/dashboardData.js`.
- `/app/reports` exists in `routeConfig.js` but is `blocked`.
- Current resource APIs cover auth, employees, members, subscriptions, courses, branches,
  attendance, and sessions.
- Current backend API docs and code expose live dashboard aggregate endpoints for admin/manager
  roles.

## Screens And Routes

| Route | Access | FE11 behavior |
|---|---|---|
| `/app/dashboard` | `admin`, `manager`, `trainer`, `receptionist` | Replace static dashboard sections with live dashboard API data, with loading/empty/error/stale states. |
| `/app/reports` | `admin`, `manager` | Optional first live report landing after backend contract exists; otherwise keep route blocked and do not fake reports. |

Route config changes:

- Keep `/app/dashboard` ready for all staff roles.
- Update dashboard route description from sample metrics to live metrics only after backend data is
  connected.
- Keep `/app/reports` blocked until a concrete backend route set exists.
- If backend adds report endpoints that support a useful report landing, change `/app/reports` to
  `ready` and add a real `ReportsPage`; otherwise leave reports out of FE11 implementation.

## Component Plan

Target frontend files after the backend contract is available:

| Path | Responsibility |
|---|---|
| `src/lib/dashboardApi.js` | Endpoint helpers for summary, revenue, plan/member distribution, recent members, and today's sessions. |
| `src/components/DashboardHome.jsx` | Replace static import with API-backed query state and section-level fallback states. |
| `src/components/dashboardData.js` | Remove or demote to local mock/test fixture once live data is wired. |
| `src/components/KpiCard.jsx` | Accept numeric/label/delta values from API without sample assumptions. |
| `src/components/RevenueBars.jsx` | Render empty/loading/error-friendly chart state from live series. |
| `src/components/PlanDonut.jsx` | Render empty/loading/error-friendly distribution state from live data. |
| `src/components/MemberTable.jsx` | Render recent members from API and support empty response. |
| `src/components/ScheduleList.jsx` | Render today's sessions from API and support empty response. |
| `src/components/reports/ReportsPage.jsx` | Optional report landing only if backend report endpoints exist beyond dashboard data. |
| `src/routes/routeConfig.js` | Update route metadata and optionally mark reports ready. |
| `src/index.css` | Add only scoped dashboard live-state styles needed for stable loading/empty/error panels. |

Keep dashboard-specific transforms local. Do not add a global charting dependency unless the existing
CSS bars/donut become insufficient for the backend response.

## State And API Plan

Preferred API helpers after backend implementation:

```js
getDashboardSummary(accessToken, params)
getDashboardRevenue(accessToken, params)
getDashboardPlanDistribution(accessToken, params)
getDashboardRecentMembers(accessToken, params)
getDashboardTodaySessions(accessToken, params)
```

Suggested backend contract to confirm in a backend plan before FE implementation:

| Action | Suggested endpoint | Response purpose |
|---|---|---|
| Dashboard summary | `GET /api/v1/dashboard/summary?branch_id=&from=&to=` | KPI values and deltas: active members, revenue, check-ins, classes. |
| Revenue trend | `GET /api/v1/dashboard/revenue?from=&to=&bucket=day` | Ordered time-series buckets in VND. |
| Plan distribution | `GET /api/v1/dashboard/plans?branch_id=&from=&to=` | Course/plan counts for donut/list display. |
| Recent members | `GET /api/v1/dashboard/members/recent?limit=5` | Latest member rows for dashboard table. |
| Today's sessions | `GET /api/v1/dashboard/sessions/today?branch_id=` | Session cards with trainer, time, capacity, and status. |

Backend needs to define:

- role access for dashboard and reports
- branch scoping behavior for staff assigned to one or more branches
- date range defaults and timezone rules
- currency units and whether revenue includes refunds
- whether check-in totals include makeup and session check-ins
- response shapes for empty periods

Frontend local state:

- dashboard query status: `loading | success | error | stale`
- section states for summary, revenue, plan distribution, recent members, and sessions
- filter state if backend supports branch/date range filters
- last successful payload for stale-data display after a background refresh fails

Request behavior:

- Fetch dashboard data after auth restore succeeds.
- Use the existing access token from `useAuth()`.
- Prefer independent section requests if backend supports separate endpoints; one failed section
  should not blank the entire dashboard.
- If backend provides one aggregate endpoint, keep section-level derived states from that payload.
- Keep sample data out of production render once live endpoints are connected.

## UX States

Dashboard:

- Initial loading with stable KPI/chart/member/session panel dimensions.
- Empty branch/date period with clear zero-state values.
- Section-level API error for revenue, plan distribution, recent members, or sessions.
- Full dashboard error only when the core summary cannot load.
- Stale-data alert if a manual/background refresh fails after live data was already shown.
- `401` session expired path through existing protected route/auth behavior.
- `403` forbidden state if backend dashboard access is narrower than the route roles.
- No dashboard endpoint available: keep FE02 sample dashboard and route metadata explicit; do not
  silently present sample numbers as live.

Reports route:

- If still blocked, continue using the existing module placeholder/blocked navigation state.
- If implemented, show loading, empty, filter validation, and backend error states before adding any
  export or advanced reporting controls.

## Responsive And Accessibility Notes

- Preserve the FE02.1 mobile containment work until FE12 does the full responsive redesign.
- Live loading and error copy must not resize KPI cards or chart panels enough to cause layout shift.
- At 320px and 390px, use compact KPI rows and summary numbers before dense charts.
- Chart values need text equivalents so the dashboard is understandable without visual bars/donut
  alone.
- Dashboard refresh/error messages should use `aria-live="polite"` when they change after user
  action.
- Filter controls, if added, need visible labels and keyboard-reachable reset/apply actions.
- Avoid adding a new card layer around existing dashboard panels.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/11_live_dashboard_apis.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/11_live_dashboard_apis.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/11_live_dashboard_apis.md`

Frontend verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Mocked browser checks:

- `/app/dashboard` renders loading, success, empty, section error, and stale-data states with mocked
  dashboard responses.
- Mobile `390x844` and desktop `1280x800` viewports have no page-level horizontal overflow.
- Route `/app/reports` remains blocked or renders the new report landing according to the confirmed
  backend contract.

Live backend checks when API and credentials are available:

- Login as an allowed staff role.
- Verify dashboard summary, revenue, plan distribution, recent members, and today's sessions render
  from live API responses.
- Verify at least one zero-data period.
- Verify branch/date filter behavior if the backend contract includes filters.
- Verify `403` with a role or account that should not access reports.

Docs updates:

- Do not update `docs/api_contract.md` until backend dashboard/report endpoints are implemented.
- After backend implementation, reconcile FE11 plan, `docs/api_contract.md`, `api_test.http`, and
  frontend tests with the final endpoint names and response shapes.

## Backend Contract Notes

- Implemented dashboard endpoints are listed in `docs/api_contract.md`.
- Current dashboard values should come from `/api/v1/dashboard/*`, not static sample fixtures.
- Broader `/api/v1/reports/*` export/reporting endpoints still do not exist.
- Member and subscription global list/search gaps may limit dashboard drill-down links.

## Risks And Boundaries

- Do not implement fake live metrics by aggregating incomplete frontend resource screens.
- Do not scrape existing list pages for dashboard data.
- Do not add dashboard export, report scheduling, or payment-history reporting in this cycle.
- Revenue semantics are product-sensitive; wait for backend definition before showing money totals as
  live.
- FE12 remains responsible for broad responsive/accessibility hardening; FE11 should only make
  dashboard live-state layout stable enough for its new data states.

## Next Action

Use `$gym-fe-review` with `CHAT_CONTEXT/frontend_skills/implementations/11_live_dashboard_apis.md`.
