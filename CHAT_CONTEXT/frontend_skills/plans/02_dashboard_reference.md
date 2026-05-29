# FE Plan - 02 Dashboard Reference

Status: Completed

Created: 2026-05-29

## Goal

Upgrade the protected `/app` dashboard from the current auth/session placeholder into an operational
staff dashboard based on `fe-tham-khao/iron_forge_gym_dashboard.html`.

This cycle is frontend-only. It should improve the first authenticated screen with realistic dashboard
structure, role-aware navigation polish, KPI cards, charts, member table, and class schedule while
keeping API integration limited to the existing auth state. Backend aggregate/report endpoints do not
exist yet, so this cycle must not pretend to fetch live dashboard metrics.

## Current Baseline

- FE 01 Auth Shell exists:
  - `/login` public route.
  - `/app` protected route.
  - `AuthContext` restores current staff via `GET /api/v1/auth/me`.
  - `AppShell` owns sidebar/topbar/logout.
  - `DashboardHome` currently shows employee identity, role badges, module count, and API session
    status.
- Reference file:
  - `fe-tham-khao/iron_forge_gym_dashboard.html`
  - Static HTML/CSS dashboard with sidebar sections, topbar, KPI cards, bar chart, donut chart,
    member table, and schedule list.
- Reference file has no real assets, no favicon, no React components, and uses CSS variables that do
  not match current app tokens.
- Current backend has business CRUD/action APIs but no dashboard/report aggregate endpoints.

## Screens And Routes

No new route is needed.

| Route | Access | Change |
|---|---|---|
| `/app` | Protected | Replace placeholder dashboard content with the reference-inspired dashboard home. |

Keep `/login` and `/` behavior unchanged.

## Reference Elements To Adapt

From `iron_forge_gym_dashboard.html`, adapt these ideas:

- Sidebar grouping:
  - `Tong quan`: Dashboard, Members, Sessions/Attendance, Reports
  - `Quan ly`: Trainers/Employees, Courses, Branches, Payments
  - bottom settings/logout area remains app-owned
- Topbar:
  - dashboard title
  - current date
  - current employee avatar/initials
  - optional notification/search buttons as non-functional visual controls
- KPI cards:
  - active members
  - monthly revenue
  - today's check-ins
  - classes this week
- Chart row:
  - revenue bar chart
  - subscription/course distribution donut
- Bottom row:
  - latest member registrations table
  - today's class schedule list
- Status badges:
  - paid, pending, waiting payment
  - class capacity ok/warn/full

Do not copy raw CSS wholesale. Translate the reference into the app's current token names and
component structure.

## Component Plan

Update or add files under `frontend/src/`:

| Path | Responsibility |
|---|---|
| `components/DashboardHome.jsx` | Compose the new dashboard home from smaller sections. |
| `components/dashboardData.js` | Static dashboard sample data for this frontend-only cycle. |
| `components/KpiCard.jsx` | Reusable KPI card with label, value, delta, and tone. |
| `components/RevenueBars.jsx` | CSS bar chart using static values. |
| `components/PlanDonut.jsx` | SVG donut chart based on static distribution values. |
| `components/MemberTable.jsx` | Latest members table with status badges. |
| `components/ScheduleList.jsx` | Today's class schedule list with capacity states. |
| `components/AppShell.jsx` | Refine sidebar groups and topbar if needed, keeping auth/logout behavior unchanged. |
| `index.css` | Add dashboard-specific layout classes, responsive rules, badges, charts, and compact table styles. |

Keep components simple and local. Do not add chart/icon dependencies for this cycle.

## State And Data Plan

State sources:

- Auth employee from `AuthContext`.
- Static dashboard sample data from `components/dashboardData.js`.

Suggested data shape:

```js
{
  kpis: [
    { key, label, value, delta, tone, accent }
  ],
  revenue: [
    { day, value, tone }
  ],
  planDistribution: [
    { label, value, color }
  ],
  latestMembers: [
    { name, phone, plan, trainer, status }
  ],
  todaySessions: [
    { time, name, trainer, room, capacity, maxCapacity, tone }
  ]
}
```

Important behavior:

- Keep sample data in one clearly named module so later API-backed implementation can replace it.
- Do not store or display tokens.
- Keep current employee identity live from auth state, not static data.
- Keep role filtering for navigation from FE 01.

## API Plan

No new backend API calls in this cycle.

Existing dependency remains:

| Purpose | Endpoint | Owner |
|---|---|---|
| Restore staff session | `GET /api/v1/auth/me` | Existing AuthContext |

Backend-contract gap to document:

- Dashboard metrics such as active member count, monthly revenue, today check-ins, class counts,
  recent member registrations, and plan distribution are not available as aggregate dashboard APIs.
- If live metrics are required later, plan a backend report/dashboard cycle first.

## UX States

Dashboard `/app`:

- Auth checking state remains owned by `RouteGuard`.
- Authenticated state shows dashboard content immediately.
- Static dashboard data means no dashboard loading spinner is needed.
- If `employee.branch_id` is empty, show branch count as `0` in the staff summary area.
- Disabled navigation items should remain visibly disabled until their module screens are implemented.

Empty/fallback handling:

- If `latestMembers` is empty, show a compact empty row.
- If `todaySessions` is empty, show a compact empty schedule item.
- If revenue/distribution arrays are empty, render an empty chart panel with stable dimensions.

## Responsive And Accessibility Notes

- Desktop: keep dense dashboard layout with KPI grid, two-column chart row, and table/schedule bottom
  row.
- Tablet: collapse dashboard sections to two columns.
- Mobile around 320px:
  - single-column cards
  - horizontally scrollable table only inside table container
  - no page-level horizontal overflow
  - topbar staff summary must not overlap logout button
- Use semantic headings for dashboard sections.
- Use real `<table>` for latest members.
- SVG donut should include `role="img"` and an accessible label.
- Chart bars should have text labels or `aria-label` values so information is not color-only.
- Badges/capacity states should not rely only on color; include text like `Paid`, `Pending`, `Full`.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/02_dashboard_reference.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/02_dashboard_reference.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/02_dashboard_reference.md`

Verification:

```sh
cd frontend
npm run build
npm run dev
```

Manual checks:

- `/login` still renders and can log in when backend is running.
- `/app` renders the new dashboard after login.
- Refresh `/app` still restores session through `GET /api/v1/auth/me`.
- Logout still clears tokens and returns to `/login`.
- Dashboard content does not imply metrics are live API data.
- Check 320px mobile width and 1280px desktop width for overflow/overlap.

## Risks And Boundaries

- This cycle intentionally uses static dashboard data because backend dashboard/report endpoints do not
  exist.
- Copying the reference raw CSS would conflict with current app tokens; implementation should translate
  the design into existing CSS.
- The reference uses `ti ti-*` icon classes, but current FE has no icon library. Use text labels,
  simple CSS marks, or existing inline SVG only if needed; do not add an icon dependency just for this
  cycle.
- Tables and charts can make mobile overflow likely; responsive layout must be verified manually.
- Live member/subscription workflows remain a later cycle.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/02_dashboard_reference.md`.
