# FE Plan - 00 Frontend Roadmap

Status: Planned

Created: 2026-05-29

## Goal

Create the top-level frontend roadmap for the gym management React/Vite app. Feature-specific plans
should reference this file first, then narrow their own scope, routes, components, API calls, test
cases, and risks.

The roadmap converts the current staff auth shell into a full operations console for members,
subscriptions, attendance, sessions, employees, courses, branches, and dashboard reporting while
keeping backend-contract gaps explicit.

## Current Baseline

- Stack: React 18 + Vite 8.
- Current app routes:
  - `/` redirects by auth state.
  - `/login` is the public staff login screen.
  - `/app` is the protected staff workspace shell.
- Current frontend capabilities:
  - bearer-token API client
  - localStorage token persistence
  - session restore through `GET /api/v1/auth/me`
  - one refresh retry through `POST /api/v1/auth/refresh`
  - logout through `POST /api/v1/auth/logout`
  - protected shell, role-aware disabled navigation, and the FE02 static dashboard
  - FE02.1 responsive repair completed as temporary containment because manual viewport review found
    the dashboard layout visually poor
- Current reference material:
  - `fe-tham-khao/iron_forge_gym_dashboard.html` - operational dashboard mockup.
  - `fe-tham-khao/iron-forge-gym/iron-forge-gym/` - marketing/brand React reference.
  - `frontend/iron-forge-brand-assets/` - logo, color, social, print, and website asset kit.
- Current backend APIs exist for auth, employees, members, subscriptions, courses, branches,
  attendance, and sessions.
- Current backend APIs do not include dashboard aggregate/report endpoints for live KPI/revenue/chart
  data.

## Roadmap Summary

| FE | Feature | Primary goal | Data mode | Backend gap |
|---|---|---|---|---|
| 01 | Auth Shell | Staff login, restore, logout, protected shell | Live auth API | None |
| 02 | Dashboard Reference | Replace placeholder dashboard with reference-inspired static ops dashboard | Static sample data + live staff identity | Live metrics/report APIs |
| 02.1 | Dashboard Responsive Repair | Fix shell/dashboard mobile and tablet layout before adding more pages | Static sample data + live staff identity | None |
| 03 | App Routing And API Foundation | Scale routes, API hooks, shared table/form/status patterns | Live auth + shared client | None |
| 04 | Brand Asset Integration | Adopt official favicon/logo/color/web assets consistently | Static assets | None |
| 05 | Members | Create/search/view members, activate offline payment, member subscriptions | Live member APIs | List/search endpoint may be needed if absent |
| 06 | Courses And Branches | Manage selectable courses and branches for later forms | Live CRUD APIs | None |
| 07 | Subscriptions | Create pending subscriptions and lifecycle actions | Live subscription APIs | List endpoint may be needed if absent |
| 08 | Attendance | Free check-in, report missed, makeup, subscription attendance history | Live attendance APIs | Member/subscription lookup UX may need list/search |
| 09 | Sessions | Session calendar/list, detail, enroll, check-in | Live session APIs | None |
| 10 | Employees | Admin staff management and password reset | Live employee APIs | None |
| 11 | Live Dashboard APIs | Replace static dashboard metrics with backend data | Live dashboard/report APIs | New backend cycle required |
| 12 | UX/Test Hardening | Browser automation, accessibility pass, mobile/desktop visual checks, responsive cleanup from manual review | App-wide | Optional test tooling |

## Screens And Routes Overview

Target protected route groups:

- `/app/dashboard` - dashboard home.
- `/app/members`, `/app/members/new`, `/app/members/:id`.
- `/app/subscriptions`, `/app/subscriptions/new`, `/app/subscriptions/:id`.
- `/app/attendance`, `/app/subscriptions/:id/attendance`.
- `/app/sessions`, `/app/sessions/new`, `/app/sessions/:id`.
- `/app/employees`, `/app/employees/new`, `/app/employees/:id`.
- `/app/settings/courses`, `/app/settings/courses/:id`.
- `/app/settings/branches`, `/app/settings/branches/:id`.

Public routes stay small:

- `/login`
- `/`

## Component Plan

Shared components to introduce before or during FE 03:

- route layout and route registry
- module sidebar groups
- page header/action bar
- data table and mobile table wrapper
- empty/loading/error state components
- form field, field error, and submit bar
- status badge and role badge
- confirm dialog for destructive or irreversible actions
- API-backed select/autocomplete for courses, branches, members, and subscriptions

Feature components should stay local until two features need the same pattern.

## State And API Plan

State ownership:

- Auth state remains in `AuthContext`.
- Resource screens start with local component state.
- Add reusable hooks for repeated query/mutation shapes only after FE 03 identifies the common
  loading/error/retry contract.

API ownership:

- `frontend/src/lib/api.js` stays the low-level request/error parser.
- Resource modules should add small endpoint helpers, for example `membersApi`, `subscriptionsApi`,
  `attendanceApi`, instead of calling raw paths from views.
- Token refresh remains boot/session-restore scoped until repeated protected API screens justify a
  broader retry strategy.

## UX States

Every resource screen should plan these states:

- initial loading
- empty data
- validation error
- backend error-envelope message
- `401` session expired path
- `403` forbidden role path
- `404` not found
- `409` business conflict
- success feedback after create/update/action
- disabled actions for unsupported role or invalid lifecycle state

## Responsive And Accessibility Notes

- Mobile target starts at 320px with no page-level horizontal overflow.
- Tables may scroll inside their own wrapper, not the entire app viewport.
- Topbar, staff summary, and action bars must wrap without covering content.
- Forms use visible labels, field-level errors, and submit loading states.
- Error/status messages use `aria-live="polite"` where they update after user action.
- Modals and confirmations must be keyboard reachable before they are used for destructive actions.
- Dense dashboard panels should keep stable dimensions so loading, empty, and hover states do not
  shift the layout.

## Feedback Loop

- Use implementation notes, review notes, test notes, and screenshots as the source for plan updates
  when a feature reveals a layout, accessibility, or contract mismatch.
- If desktop or mobile screenshots show overflow, clipping, or broken stacking, record the exact
  breakpoint and affected component in the feature plan before fixing it.
- Keep the feature plan and test note aligned: the plan states what should happen, and the test note
  records what was actually verified in browser/manual checks.
- For FE 02 and later dashboard-like screens, treat mobile responsive behavior as a required list
  item, not an optional polish item.

## Feature Plans

### FE 01 Auth Shell

Status: implemented and retested.

Plan file: `CHAT_CONTEXT/frontend_skills/plans/01_auth_shell.md`

Scope:

- `/login`, `/app`, `/`
- `AuthContext`, `api.js`, token storage, route guard, app shell, dashboard placeholder
- Backend auth endpoints only

Follow-up:

- Use `CHAT_CONTEXT/frontend_skills/tests/01_auth_shell.md` before final FE completion if a human
  browser pass is required.

### FE 02 Dashboard Reference

Status: completed with responsive follow-up.

Plan file: `CHAT_CONTEXT/frontend_skills/plans/02_dashboard_reference.md`

Scope:

- Adapt `fe-tham-khao/iron_forge_gym_dashboard.html` into the protected `/app` dashboard.
- Add KPI cards, revenue bars, plan donut, latest members table, and today's class schedule.
- Use static `dashboardData` until backend aggregate/report endpoints exist.
- Keep current live staff identity from auth state.

Implementation notes:

- Do not copy the raw reference CSS wholesale.
- Keep static data visibly internal/sample; do not imply metrics are live.
- Use the current Iron Forge dark/orange token direction.
- Manual viewport review found the current responsive layout visually poor, so FE 02.1 should run
  before FE 03.

### FE 02.1 Dashboard Responsive Repair

Status: implemented and quick reviewed/tested.

Plan file: `CHAT_CONTEXT/frontend_skills/plans/02_1_dashboard_responsive_repair.md`

Scope:

- Repair `/app` dashboard responsive behavior at 320px, 375px, 768px, and 1280px.
- Keep FE02 dashboard content, static data, auth state, and route behavior unchanged.
- Fix shared shell/topbar/mobile-nav issues before resource pages reuse the same app frame.
- Make mobile dashboard sections readable without page-level horizontal overflow.

Backend APIs:

- No new backend API calls.
- Auth flow may be manually checked if the backend is running, but this cycle is primarily layout
  repair.

Risks:

- If this is skipped, FE03+ resource pages will inherit a poor shell/topbar/mobile foundation.
- This is not the full FE12 hardening cycle; keep it focused on the current dashboard and shared shell.

### FE 03 App Routing And API Foundation

Status: planned.

Plan file: `CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md`

Scope:

- Introduce scalable route definitions for business modules:
  - `/app/dashboard`
  - `/app/members`
  - `/app/members/:id`
  - `/app/subscriptions`
  - `/app/subscriptions/:id`
  - `/app/attendance`
  - `/app/sessions`
  - `/app/sessions/:id`
  - `/app/employees`
  - `/app/settings/courses`
  - `/app/settings/branches`
- Decide whether to keep local history routing or add `react-router-dom`. For many nested screens,
  prefer React Router unless the implementation can keep the local route map simple and testable.
- Extend the API client without adding a complex interceptor:
  - typed endpoint helpers by resource
  - backend error-envelope normalization
  - loading/error/retry state helpers
  - role/permission helpers
- Add shared UI primitives:
  - page header
  - data panel
  - table wrapper
  - empty state
  - form row/field/error
  - confirm action modal
  - status badge

Backend APIs:

- Existing auth endpoints only for foundation validation.

Risks:

- Adding all resource screens before shared patterns will duplicate loading/error/table code.
- Keeping custom routing too long will make browser back/forward and nested detail screens brittle.

### FE 04 Brand Asset Integration

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/04_brand_assets.md`

Scope:

- Use `frontend/iron-forge-brand-assets/` selectively:
  - `01_logo/` for official logo and app icon variants
  - `02_colors/` for token reconciliation
  - `06_website/` for favicon, OG image, loading icon, 404 illustration, service icons
- Decide whether to replace `frontend/public/favicon.svg` with
  `frontend/iron-forge-brand-assets/06_website/iron-forge-favicon-if.svg`.
- Add only runtime assets that the app actually imports or serves. Keep social/print/mockup assets as
  reference unless a feature needs them.
- Align app CSS variables with `iron-forge-css-variables.css` while preserving existing contrast and
  compact staff-console layout.

Backend APIs:

- None.

Risks:

- Some SVG logos use text and font fallback; verify rendering before using them as critical UI assets.
- Large/social/print assets should not be copied into `public/` unless needed at runtime.

### FE 05 Members

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/05_members.md`

Screens/routes:

- `/app/members` - member search/list workspace.
- `/app/members/new` - create member form.
- `/app/members/:id` - member profile, status, contact details, subscriptions panel.

Backend APIs:

- `POST /api/v1/members`
- `GET /api/v1/members/:id`
- `PATCH /api/v1/members/:id/activate`
- `GET /api/v1/members/:id/subscriptions`

Backend-contract gap:

- `GET /api/v1/members` is not listed in the current API contract. If the UI needs a real directory,
  plan a backend list/search endpoint or use a direct-ID lookup MVP.

UX states:

- create validation, duplicate CCID conflict, not found by ID, inactive/pending state, activate
  success/failure, member has no subscriptions.

### FE 06 Courses And Branches

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/06_courses_branches.md`

Screens/routes:

- `/app/settings/courses`
- `/app/settings/courses/:id`
- `/app/settings/branches`
- `/app/settings/branches/:id`

Backend APIs:

- Courses: `POST/GET/GET:id/PATCH/DELETE /api/v1/courses`
- Branches: `POST/GET/GET:id/PATCH/DELETE /api/v1/branches`
- Nearby: `GET /api/v1/branches/nearby`

Purpose:

- Provide selectable data for subscriptions, member home branch, sessions, and attendance.
- Avoid hardcoding course/branch options in later forms.

UX states:

- loading lists, empty lists, create/edit validation, delete confirmation, duplicate code/name
  conflicts, nearby search permission/error state.

### FE 07 Subscriptions

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/07_subscriptions.md`

Screens/routes:

- `/app/subscriptions` - lookup/list workspace.
- `/app/subscriptions/new` - create pending subscription.
- `/app/subscriptions/:id` - lifecycle, remaining sessions, attendance shortcut, refund panel.

Backend APIs:

- `POST /api/v1/subscriptions`
- `GET /api/v1/subscriptions/:id`
- `PATCH /api/v1/subscriptions/:id/suspend`
- `PATCH /api/v1/subscriptions/:id/unsuspend`
- `PATCH /api/v1/subscriptions/:id/expire`
- `POST /api/v1/subscriptions/:id/refund`
- Member subscription context: `GET /api/v1/members/:id/subscriptions`

Backend-contract gap:

- `GET /api/v1/subscriptions` is not listed in the current API contract. If the UI needs global
  subscription search/list, plan backend support or keep the first pass member-scoped/direct-ID.

UX states:

- pending vs active vs suspended vs expired/refunded, invalid reference IDs, conflict on invalid
  lifecycle transition, refund amount preview, disabled actions by role/status.

### FE 08 Attendance

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/08_attendance.md`

Screens/routes:

- `/app/attendance` - check-in/report/makeup command center.
- `/app/subscriptions/:id/attendance` - subscription attendance history.

Backend APIs:

- `POST /api/v1/attendance/checkin`
- `GET /api/v1/subscriptions/:id/attendance`
- `POST /api/v1/attendance/report`
- `POST /api/v1/attendance/makeup`

UX states:

- check-in success, weekly limit conflict, expired/suspended subscription, invalid makeup reference,
  no attendance history, backend date validation errors, branch/session selection.

Dependencies:

- Courses/branches and member/subscription lookup UX should exist first, or this feature needs manual
  ObjectID entry fields for the first pass.

### FE 09 Sessions

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/09_sessions.md`

Screens/routes:

- `/app/sessions` - filterable session list/calendar.
- `/app/sessions/new` - create session.
- `/app/sessions/:id` - detail, enrollment list, enroll member subscription, check-in.

Backend APIs:

- `POST /api/v1/sessions`
- `GET /api/v1/sessions`
- `GET /api/v1/sessions/:id`
- `POST /api/v1/sessions/:id/enroll`
- `POST /api/v1/sessions/:id/checkin`

UX states:

- branch/date/level filters, empty day, capacity full, invalid subscription for enrollment, already
  enrolled, already checked in, trainer/manager role visibility.

### FE 10 Employees

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/10_employees.md`

Screens/routes:

- `/app/employees` - admin-only staff list.
- `/app/employees/new` - create employee.
- `/app/employees/:id` - profile, roles, branches, level, status, password reset.

Backend APIs:

- `POST /api/v1/employees`
- `GET /api/v1/employees`
- `GET /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id/password`

UX states:

- forbidden for non-admin, inactive employee badges, duplicate email/employee ID, reset password
  confirmation, role/branch selectors, token revocation notice after reset/deactivate.

### FE 11 Live Dashboard APIs

Status: blocked by backend contract.

Plan file to create after backend planning: `CHAT_CONTEXT/frontend_skills/plans/11_live_dashboard.md`

Frontend target:

- Replace FE 02 static dashboard sample data with API-backed metrics.
- Keep the same visual sections if FE 02 proves usable:
  - active members
  - monthly revenue
  - today check-ins
  - classes this week
  - revenue chart
  - subscription/course distribution
  - recent members
  - today's sessions

Backend-contract gap:

- Need new backend dashboard/report endpoints. Suggested future API shape:
  - `GET /api/v1/dashboard/summary`
  - `GET /api/v1/dashboard/revenue`
  - `GET /api/v1/dashboard/members/recent`
  - `GET /api/v1/dashboard/sessions/today`

### FE 12 UX/Test Hardening

Status: planned.

Plan file to create: `CHAT_CONTEXT/frontend_skills/plans/12_ux_test_hardening.md`

Scope:

- Add browser automation if the project accepts a test dependency.
- Verify login, protected route, CRUD/action flows, and visual layout at mobile/desktop widths.
- Add accessibility checks for form labels, focus states, table semantics, modals, keyboard flows, and
  aria-live error messages.
- Add a repeatable frontend verification matrix for build, audit, route shell, backend API smoke, and
  browser interaction.
- Include responsive bugfixes and viewport cleanup that were discovered during manual review so
  later feature plans can route layout issues back into this cycle.
- Treat FE02.1 as a temporary containment pass, not the final responsive design. A future FE12 pass
  should revisit dashboard/shell responsive UX with browser evidence instead of continuing to polish
  FE02.1 ad hoc.
- Define a reusable responsive pattern for operational dashboards and resource pages:
  - when to use sidebar drawer vs persistent sidebar
  - when dense charts become summary numbers
  - when tables become cards, expandable sections, or horizontal scroll
  - viewport acceptance criteria for 320px, 375px, 768px, 1080px, and desktop

Risks:

- Without Playwright or another browser runner, frontend verification remains partly manual.
- Dense admin dashboards can regress on 320px mobile unless table overflow and topbar wrapping are
  tested every cycle.

## Cross-Cutting Architecture Plan

Routing:

- FE 02 can keep the current `/app` shell.
- FE 03 should introduce scalable route ownership before business modules are implemented.
- Route access should combine auth status and role visibility; backend remains source of truth.

API:

- Keep `Authorization: Bearer <accessToken>` in one API layer.
- Keep backend error-envelope parsing centralized.
- Avoid global refresh retry for every request until there are enough protected APIs to justify it.
- Keep resource-specific helpers close to the API layer, not inside view components.

State:

- Auth remains context-owned.
- Resource features should use local component state first.
- Add shared hooks only when at least two resource screens share the same loading/error/retry pattern.
- Do not introduce a global store until cross-screen resource caching becomes a real problem.

UI:

- Keep the app a dense staff operations console, not a marketing landing page.
- Use asset-kit logo/color material selectively.
- Prefer compact tables, forms, panels, badges, and command bars.
- Avoid nested cards and page-level decorative effects.

Backend contract:

- Existing APIs are enough to start core CRUD/action screens.
- Dashboard live metrics, member global list/search, and subscription global list/search are the main
  visible gaps for a polished admin workflow.

## Docs And Test Plan

Every feature plan should include:

- target routes
- components/files to add or update
- state shape
- exact backend endpoints
- loading/empty/error states
- role/permission behavior
- responsive/accessibility checks
- build/audit/API/browser verification
- known backend-contract gaps

Expected phase files:

- `CHAT_CONTEXT/frontend_skills/plans/<feature>.md`
- `CHAT_CONTEXT/frontend_skills/implementations/<feature>.md`
- `CHAT_CONTEXT/frontend_skills/reviews/<feature>.md`
- `CHAT_CONTEXT/frontend_skills/tests/<feature>.md`

Common verification:

```sh
cd frontend
npm audit
npm run build
```

When a feature calls backend APIs:

- Start backend with CORS for `http://127.0.0.1:5173`.
- Smoke valid auth login/current-user/refresh/logout if auth behavior is touched.
- Smoke at least one success and one expected error path for each new resource/action.

## Recommended Order

1. FE 02 Dashboard Reference.
2. FE 04 Brand Asset Integration if visual consistency should be fixed before more screens.
3. FE 03 App Routing And API Foundation.
4. FE 06 Courses And Branches.
5. FE 05 Members.
6. FE 07 Subscriptions.
7. FE 08 Attendance.
8. FE 09 Sessions.
9. FE 10 Employees.
10. Backend dashboard/report cycle.
11. FE 11 Live Dashboard APIs.
12. FE 12 UX/Test Hardening, then repeat as needed.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/02_dashboard_reference.md`, or use
`$gym-fe-plan` to expand the next chosen roadmap item into its own detailed plan before coding.
