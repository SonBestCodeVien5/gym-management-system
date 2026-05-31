# FE Plan - 03 App Routing And API Foundation

Status: Completed

Created: 2026-05-31

## Goal

Prepare the React/Vite frontend for real resource pages without building those resource workflows yet.

FE03 should turn the current two-route auth shell into a scalable operations console foundation:
route registry, role-aware navigation, protected module placeholders, shared UI primitives, and a
repeatable API/state pattern that later features can reuse for Members, Courses, Branches,
Subscriptions, Attendance, Sessions, and Employees.

This cycle should not implement CRUD screens or call new business APIs. It should make FE05+ easier
and less duplicated.

## Current Baseline

- Stack: React 18 + Vite 8, no router dependency.
- Current app routes are custom and hardcoded in `frontend/src/App.jsx`:
  - `/`
  - `/login`
  - `/app`
- `AppShell` owns:
  - sidebar/mobile menu config
  - active item local state
  - topbar/staff summary/logout
  - dashboard rendering
- `AuthContext` owns login, session restore, one refresh retry, logout, and token storage.
- `frontend/src/lib/api.js` owns low-level API request parsing and auth endpoint helpers.
- FE02/FE02.1 dashboard is static sample data plus live staff identity.
- Responsive dashboard work is explicitly not final; broader responsive UX is deferred to FE12.

## Screens And Routes

Use FE03 to introduce route ownership, but keep first-pass resource screens as placeholders.

Preferred route behavior:

| Route | Access | FE03 behavior |
|---|---|---|
| `/` | Public redirect | Redirect by auth state as today. |
| `/login` | Anonymous only | Keep current login screen. |
| `/app` | Protected | Redirect/alias to `/app/dashboard`. |
| `/app/dashboard` | Protected | Render current FE02 dashboard. |
| `/app/members` | Protected, role-gated | Placeholder page for FE05. |
| `/app/members/:id` | Protected, role-gated | Placeholder/detail shell for FE05. |
| `/app/subscriptions` | Protected, role-gated | Placeholder page for FE07. |
| `/app/subscriptions/:id` | Protected, role-gated | Placeholder/detail shell for FE07/FE08. |
| `/app/attendance` | Protected, role-gated | Placeholder page for FE08. |
| `/app/sessions` | Protected, role-gated | Placeholder page for FE09. |
| `/app/sessions/:id` | Protected, role-gated | Placeholder/detail shell for FE09. |
| `/app/reports` | Protected, manager/admin | Blocked placeholder for future dashboard/report APIs. |
| `/app/employees` | Protected, admin only | Placeholder page for FE10. |
| `/app/settings/courses` | Protected, manager/admin | Placeholder page for FE06. |
| `/app/settings/branches` | Protected, manager/admin | Placeholder page for FE06. |
| `/app/payments` | Protected, role-gated | Blocked placeholder for later payment/refund workspace. |

Routing decision for FE03:

- Keep a dependency-free route registry for this cycle instead of adding `react-router-dom`.
- Implement simple static + `:id` path matching in local utilities.
- Revisit React Router before deeply nested workflows become complex.

Rationale:

- Current app has no router dependency.
- FE03 can solve immediate duplication with a route config and matcher.
- Avoid spending this cycle on package churn while the first business screens are still planned.

## Component Plan

Add or update these files under `frontend/src/`:

| Path | Responsibility |
|---|---|
| `routes/routeConfig.js` | Source of truth for public/protected routes, module labels, roles, groups, readiness, and placeholder metadata. |
| `routes/matchRoute.js` | Match current path to route config, including simple `:id` params and fallback route. |
| `components/AppShell.jsx` | Stop owning route config. Accept active route/module info and render children. Keep logout/staff summary. |
| `components/DashboardHome.jsx` | Render under `/app/dashboard`; no data/API change. |
| `components/ModulePlaceholder.jsx` | Shared placeholder for planned modules with title, description, required roles, and planned API scope. |
| `components/PageHeader.jsx` | Compact title/action header for future resource screens. |
| `components/DataPanel.jsx` | Shared panel wrapper for resource content, loading/empty/error blocks, and future tables/forms. |
| `components/StateBlock.jsx` | Reusable loading/empty/error/forbidden/not-found state block. |
| `lib/permissions.js` | Role helper functions such as `hasAnyRole`, `canAccessRoute`, and readable role labels. |
| `lib/resourceState.js` | Small helpers/constants for `idle/loading/success/error/empty` states if implementation needs them. |
| `App.jsx` | Replace `KNOWN_ROUTES` with route matching and render route element through shell/guard. |
| `index.css` | Add shared page header, data panel, placeholder, and state block styles. Keep FE02 responsive debt out of scope. |

Keep the component set minimal. Do not add form/table abstractions that are not used by at least the
placeholder/foundation screens.

## State And API Plan

State:

- Auth state remains owned by `AuthContext`.
- Route state remains local in `App.jsx` for this cycle.
- Active module should be derived from the matched route, not from local sidebar state.
- Placeholder pages should not fetch data.

API:

- No new backend API calls in FE03.
- Keep `frontend/src/lib/api.js` as the low-level request/error parser.
- Do not add token refresh for every resource request yet.
- Add only lightweight API foundation if needed:
  - consistent `accessToken` passing convention for future helpers
  - shared request-state shape documentation/helpers
  - endpoint metadata in placeholder pages, not live calls

Current API contracts to keep in mind for later features:

| Future feature | Relevant backend APIs |
|---|---|
| FE05 Members | `POST /api/v1/members`, `GET /api/v1/members/:id`, `PATCH /api/v1/members/:id/activate`, `GET /api/v1/members/:id/subscriptions` |
| FE06 Courses/Branches | `POST/GET/GET:id/PATCH/DELETE /api/v1/courses`, `POST/GET/GET:id/PATCH/DELETE /api/v1/branches`, `GET /api/v1/branches/nearby` |
| FE07 Subscriptions | `POST /api/v1/subscriptions`, `GET /api/v1/subscriptions/:id`, lifecycle/refund endpoints |
| FE08 Attendance | `POST /api/v1/attendance/checkin`, `POST /api/v1/attendance/report`, `POST /api/v1/attendance/makeup`, `GET /api/v1/subscriptions/:id/attendance` |
| FE09 Sessions | `POST /api/v1/sessions`, `GET /api/v1/sessions`, `GET /api/v1/sessions/:id`, enroll/checkin endpoints |
| FE10 Employees | `POST /api/v1/employees`, `GET /api/v1/employees`, `GET/PATCH /api/v1/employees/:id`, password reset endpoint |

Backend gaps remain unchanged:

- no live dashboard aggregate/report APIs
- no global members list/search endpoint
- no global subscriptions list/search endpoint

## UX States

FE03 placeholder/resource shell should support these states:

- protected route loading: existing `RouteGuard`
- forbidden route: user has auth but lacks required role
- not found route: unknown path under `/app`
- module planned/disabled: route exists but feature is not implemented yet
- empty placeholder: explain planned workflow without fake data
- API error display pattern: show normalized `error.code`, `message`, and optional field details

Do not create fake CRUD data in FE03.

## Responsive And Accessibility Notes

- Preserve FE02.1 behavior as a temporary layout state.
- Do not continue polishing dashboard responsive design in FE03.
- New placeholder pages must use compact operational layouts:
  - no marketing hero sections
  - page header + data panel
  - readable text at 320px
  - no page-level horizontal overflow
- Navigation:
  - desktop sidebar remains persistent above the responsive breakpoint
  - responsive menu button opens grouped module navigation
  - disabled/future modules must communicate planned state
- Accessibility:
  - active route should be clear through `aria-current="page"` where practical
  - route placeholders should use semantic headings
  - forbidden/not-found states should be discoverable and not blank
  - buttons that do not perform real actions should be disabled or clearly placeholder actions

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/03_app_foundation.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/03_app_foundation.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/03_app_foundation.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Manual checks:

- `/` still redirects by auth state.
- `/login` still renders and does not show protected shell when anonymous.
- `/app` aliases/redirects to `/app/dashboard`.
- `/app/dashboard` renders the existing dashboard.
- placeholder routes render correct title/roles/planned API summary.
- forbidden routes show a real forbidden state for a role that lacks access.
- unknown routes show a not-found state.
- browser back/forward updates the active route.
- logout still clears session and returns to `/login`.

Backend checks:

- If backend is running, verify login, current-user restore, refresh, and logout still work.
- Do not require business API calls in this cycle.

## Risks And Boundaries

- Do not build Members/Courses/Branches CRUD in FE03.
- Do not add `react-router-dom` unless implementation proves the local route matcher is already too
  brittle.
- Do not move FE12 responsive redesign work into this cycle.
- Avoid over-abstracting forms/tables before real resource screens need them.
- Route-level role gating in FE is only UX; backend remains the security source of truth.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md`.
