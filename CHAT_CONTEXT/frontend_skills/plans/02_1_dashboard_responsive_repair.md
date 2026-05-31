# FE Plan - 02.1 Dashboard Responsive Repair

Status: Implemented

Created: 2026-05-31

## Goal

Repair the FE02 `/app` dashboard responsive layout after manual viewport review found the current
mobile experience visually poor.

This is a narrow frontend-only repair cycle before FE03. It should stabilize the shared shell,
topbar, mobile navigation, and dashboard section stacking so later resource pages do not inherit a
broken layout foundation.

## Current Baseline

- FE02 dashboard is implemented and pushed.
- `/app` still uses the current custom route shell from `App.jsx`.
- The dashboard uses static sample data from `dashboardData.js` plus live employee identity from
  `AuthContext`.
- Manual user review found the responsive UI unacceptable.
- Existing automated checks only covered build, route smoke, SSR render smoke, and static CSS checks;
  they did not visually verify browser layout.

## Screens And Routes

No new route is needed.

| Route | Access | Change |
|---|---|---|
| `/app` | Protected | Repair shell/dashboard responsive behavior while keeping FE02 content and auth behavior intact. |

Keep `/login`, `/`, and auth flow behavior unchanged.

## Problem Areas To Inspect

- Topbar:
  - title and staff summary stack poorly on narrow widths
  - avatar, employee text, tool placeholders, and logout compete for one row
- Mobile navigation:
  - too many disabled future modules are exposed at once
  - horizontal nav may feel noisy before those modules exist
- Dashboard density:
  - KPI cards and large display type feel oversized on mobile
  - chart panels and donut legend need better narrow-width spacing
  - bottom table/schedule section feels too dense for small screens
- Table:
  - current `min-width: 560px` keeps table semantics but may create a heavy mobile interaction
  - page-level overflow must remain zero; table overflow may stay inside its wrapper
- Staff context:
  - long employee values and role chips must wrap cleanly

## Component Plan

Update only the current FE02 surface unless implementation reveals a tiny helper is needed.

| Path | Responsibility |
|---|---|
| `frontend/src/components/AppShell.jsx` | Adjust mobile-visible navigation and topbar markup if CSS alone is not enough. Keep logout/auth behavior unchanged. |
| `frontend/src/components/DashboardHome.jsx` | Adjust section grouping or labels only if needed for responsive clarity. Do not change data source. |
| `frontend/src/components/MemberTable.jsx` | Consider mobile-friendly rendering while preserving desktop table semantics. |
| `frontend/src/components/ScheduleList.jsx` | Keep schedule rows readable at 320px without overlap. |
| `frontend/src/index.css` | Main repair target: breakpoints, spacing, font sizes, grids, topbar wrapping, mobile nav, chart/table/schedule/staff context layout. |

Avoid adding new dependencies in this repair cycle.

## State And API Plan

No state or API changes.

- Keep auth state from `AuthContext`.
- Keep static dashboard data from `dashboardData.js`.
- Do not add dashboard/report API calls.
- Do not alter token storage or refresh behavior.

## UX States

Retain FE02 states:

- authenticated `/app` renders dashboard immediately
- disabled module placeholders remain disabled
- empty fallbacks for table, chart, donut, and schedule remain available

Responsive repair behavior:

- At mobile widths, avoid showing every future module as equal priority if it makes the navigation
  noisy. Prefer a compact horizontally scrollable nav, reduced labels, or hiding disabled future
  modules on mobile.
- Staff identity should remain visible but may be compressed to name/email with ellipsis or stacked
  below the title.
- Dashboard panels should feel like operational mobile sections, not shrunken desktop cards.

## Responsive And Accessibility Requirements

Required viewport targets:

- 320px mobile
- 375px mobile
- 768px tablet
- 1080px responsive shell breakpoint
- 1280px desktop

Pass criteria:

- no page-level horizontal overflow at 320px
- no overlapping text/buttons in topbar or staff summary
- logout remains visible and reachable
- mobile nav remains usable and does not dominate the viewport
- KPI values fit inside cards without oversized display type
- revenue bars, donut chart, and legends do not overlap
- latest members table either scrolls only inside its wrapper or switches to a readable mobile layout
- schedule rows wrap without clipping names, room, trainer, or capacity
- staff context values and role chips wrap without breaking the panel
- disabled notification/search placeholders remain non-focusable through `disabled`

Accessibility:

- preserve semantic headings
- preserve table semantics on desktop
- preserve donut `role="img"` label and bar `aria-label` values
- do not remove visible text from badges/capacity states
- keep focusable controls meaningful; placeholders must remain disabled or non-interactive

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/02_1_dashboard_responsive_repair.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/02_1_dashboard_responsive_repair.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/02_1_dashboard_responsive_repair.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Manual browser checks:

- log in or use an existing local auth session to open `/app`
- check 320px, 375px, 768px, 1080px, and 1280px widths
- confirm no page-level horizontal overflow
- confirm topbar, mobile nav, KPI grid, charts, member table, schedule list, and staff context are
  readable and do not overlap
- capture or describe any remaining breakpoint-specific issue in the test note

Backend/API checks:

- only auth login/restore/logout if backend is running
- backend is not required for static dashboard responsive repair if a valid session is already
  available

## Risks And Boundaries

- This cycle should not become a full FE12 hardening pass.
- Do not introduce Playwright or other browser automation here unless explicitly requested.
- Do not start FE03 routing changes inside this repair cycle.
- If the dashboard design still feels too dense after CSS repair, document a smaller dashboard
  content redesign before implementing resource pages.

## Next Action

Use `$gym-fe-implement` with
`CHAT_CONTEXT/frontend_skills/plans/02_1_dashboard_responsive_repair.md`.
