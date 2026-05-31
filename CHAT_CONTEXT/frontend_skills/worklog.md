# Frontend Worklog

Use this file for short frontend roadmap and completion summaries.

## Current frontend roadmap

- [x] FE 01 Auth Shell - staff login, API client, token restore, logout, protected shell
- [x] FE 02 Dashboard Reference - adapt `fe-tham-khao/iron_forge_gym_dashboard.html` into `/app`
- [x] FE 02.1 Dashboard Responsive Repair - fix shell/dashboard mobile and tablet layout before FE03
- [x] FE 03 App Routing And API Foundation - scalable routes, resource API helpers, shared UI states
- [ ] FE 04 Brand Asset Integration - official logo/favicon/color/web assets from `iron-forge-brand-assets`
- [ ] FE 05 Members - create/search/detail, activate offline payment, member subscriptions
- [ ] FE 06 Courses And Branches - CRUD settings and selectable reference data for later forms
- [ ] FE 07 Subscriptions - create pending subscription and lifecycle actions
- [ ] FE 08 Attendance - check-in, report missed, makeup, subscription attendance history
- [ ] FE 09 Sessions - session list/calendar, create, enroll, session check-in
- [ ] FE 10 Employees - admin-only staff management and password reset
- [ ] FE 11 Live Dashboard APIs - replace static dashboard metrics after backend report APIs exist
- [ ] FE 12 UX/Test Hardening - browser automation, accessibility, mobile/desktop visual checks

Roadmap source: `CHAT_CONTEXT/frontend_skills/plans/00_frontend_roadmap.md`.

## Planned - 2026-05-29 - FE 00 Frontend Roadmap

Created a top-level frontend roadmap covering dashboard reference work, brand asset integration,
shared routing/API foundation, each backend-backed resource module, dashboard API gaps, and app-wide
test/UX hardening.

Feature-specific plans should reference `CHAT_CONTEXT/frontend_skills/plans/00_frontend_roadmap.md`
before narrowing their own scope.

Recommended next action: use `$gym-fe-implement` with
`CHAT_CONTEXT/frontend_skills/plans/02_dashboard_reference.md`, or use `$gym-fe-plan` to expand FE 03
or FE 04 first if the next cycle should prepare routing/assets before dashboard implementation.

## Planned - 2026-05-29 - FE 01 Auth Shell

Create the first operational frontend slice: staff login, bearer-token API client, session restore via
`GET /api/v1/auth/me`, logout, and a protected dashboard shell.

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/01_auth_shell.md`.

## Implemented - 2026-05-29 - FE 01 Auth Shell

Implemented staff login, dependency-free route switching, auth state, token storage helpers,
session restore with one refresh retry, logout, protected dashboard shell, role-aware disabled module
navigation, and responsive staff portal styling.

Build passed with `npm run build`.

Use `$gym-fe-review` with `CHAT_CONTEXT/frontend_skills/implementations/01_auth_shell.md`.

## Reviewed - 2026-05-29 - FE 01 Auth Shell

Review found one medium contract-drift issue: the login form blocks passwords shorter than 8
characters even though backend login accepts any non-empty password for bcrypt comparison.

Use `$gym-fe-implement` to fix `CHAT_CONTEXT/frontend_skills/reviews/01_auth_shell.md`, then
`$gym-fe-test`.

## Tested - 2026-05-29 - FE 01 Auth Shell

Build passed and Vite served `/login` plus `/app` as React app shell routes. Backend auth flow checks
were blocked because `127.0.0.1:8080` was not running. The review finding about login password
min-length validation remains open.

Not ready for `$gym-fe-complete`. Fix the review finding, then rerun `$gym-fe-test` with backend
available.

## Implemented Fix - 2026-05-29 - FE 01 Auth Shell

Removed login-time password min-length validation so the frontend matches backend auth behavior:
login now requires only a non-empty password before calling `/api/v1/auth/login`.

Build passed with `npm run build`.

Use `$gym-fe-test` with `CHAT_CONTEXT/frontend_skills/tests/01_auth_shell.md`.

## Retested - 2026-05-29 - FE 01 Auth Shell

Retest passed for `npm run build`, removed login min-length regression, Vite route shell responses for
`/`, `/login`, and `/app`, backend auth login/current-user/refresh/logout/wrong-password behavior, and
CORS preflight from `http://127.0.0.1:5173`.

Browser visual and React form interaction checks were not run because no local browser automation is
installed. Ready for `$gym-fe-complete` with that residual risk noted in
`tests/01_auth_shell.md`.

## UI Brand Cleanup - 2026-05-29 - FE 01 Auth Shell

Applied the usable brand direction from `fe-tham-khao`: text wordmark styling now uses `IRON` in the
main text color and `FORGE` in the orange accent. Replaced the mismatched tab icon with a new
Iron Forge `IF` favicon asset in `frontend/public/favicon.svg`.

Build passed with `npm run build`.

## Planned - 2026-05-29 - FE 02 Dashboard Reference

Plan added to adapt `fe-tham-khao/iron_forge_gym_dashboard.html` into the protected `/app` dashboard:
KPI cards, revenue bars, plan donut, latest members table, and today's class schedule using static
frontend data until backend dashboard/report endpoints exist.

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/02_dashboard_reference.md`.

## Implemented - 2026-05-29 - FE 02 Dashboard Reference

Implemented the reference-inspired `/app` dashboard with static sample KPI data, revenue bars, plan
donut, latest member table, today's class schedule, grouped sidebar placeholders, dashboard topbar
tools, and staff context sourced from the live auth state.

Build passed with `npm run build`.

Use `$gym-fe-review` with `CHAT_CONTEXT/frontend_skills/implementations/02_dashboard_reference.md`.

## Reviewed - 2026-05-29 - FE 02 Dashboard Reference

Review found one low-priority accessibility issue: the notification and search topbar controls were
focusable placeholders without a real action path or explicit coming-soon state.

## Tested - 2026-05-29 - FE 02 Dashboard Reference

Build, route smoke, dashboard render smoke, review-fix static checks, and whitespace checks passed.
Browser viewport verification and authenticated backend flow verification were not run because
local browser automation and a seeded backend were unavailable in this pass.

## Completed - 2026-05-29 - FE 02 Dashboard Reference

FE02 is complete as a frontend-only dashboard cycle.

- `/app` now renders the reference-inspired operational dashboard with static sample metrics.
- Live staff identity remains sourced from the auth session.
- Remaining viewport/backend checks are documented as local limitations in the FE02 test note.

Manual viewport feedback later found the FE02 responsive layout visually poor, so FE02.1 is planned
below before FE03.

## Planned - 2026-05-31 - FE 02.1 Dashboard Responsive Repair

Manual viewport review found the FE02 dashboard responsive layout visually poor. Added a narrow
repair cycle before FE03 so the shared shell, topbar, mobile navigation, dashboard cards, charts,
member table, schedule list, and staff context are fixed before more resource pages reuse the same
layout foundation.

Plan file: `CHAT_CONTEXT/frontend_skills/plans/02_1_dashboard_responsive_repair.md`.

Recommended next action: use `$gym-fe-implement` with
`CHAT_CONTEXT/frontend_skills/plans/02_1_dashboard_responsive_repair.md`.

## Implemented - 2026-05-31 - FE 02.1 Dashboard Responsive Repair

Implemented the responsive repair for the FE02 `/app` dashboard: raised the primary shell breakpoint
to `1080px`, added a responsive Menu button for the grouped sidebar, moved the staff context card to
the top on the normal dashboard, converted mobile KPIs into compact rows, replaced mobile chart panels
with number summaries, and made Members/Classes expandable on narrow screens. Auth/API behavior stayed
unchanged.

Build passed with `npm run build`.

Use `$gym-fe-review` with
`CHAT_CONTEXT/frontend_skills/implementations/02_1_dashboard_responsive_repair.md`.

## Reviewed/Tested - 2026-05-31 - FE 02.1 Dashboard Responsive Repair

Quick review/test passed for build, whitespace, and SSR render smoke. No blocking code issue was
found, but manual feedback says responsive still is not final. Broader responsive UX should move to
FE12 instead of continuing to polish FE02.1.

Review note: `CHAT_CONTEXT/frontend_skills/reviews/02_1_dashboard_responsive_repair.md`.
Test note: `CHAT_CONTEXT/frontend_skills/tests/02_1_dashboard_responsive_repair.md`.

## Planned - 2026-05-31 - FE 03 App Routing And API Foundation

Planned the next foundation feature after FE02.1: route registry, simple dependency-free route
matching, role-aware module navigation, protected placeholder pages, shared page/data/state UI
primitives, and API-state conventions for later resource screens.

Plan file: `CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md`.

Recommended next action: use `$gym-fe-implement` with
`CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md`.

## Implemented - 2026-05-31 - FE 03 App Routing And API Foundation

Implemented the route/API foundation: route registry, dependency-free route matcher, role helpers,
route-driven `AppShell`, `/app` to `/app/dashboard` redirect, forbidden/not-found states, placeholder
module pages, shared page/data/state UI primitives, and route matcher smoke checks.

No business API calls or CRUD screens were added.

Build passed with `npm run build`.

Use `$gym-fe-review` with `CHAT_CONTEXT/frontend_skills/implementations/03_app_foundation.md`.

## Reviewed - 2026-05-31 - FE 03 App Routing And API Foundation

Review found one medium route matcher issue: malformed encoded detail params could throw
`URIError` during render instead of falling through to the app not-found state.

The finding was fixed in `frontend/src/routes/matchRoute.js` by safely decoding dynamic param
segments and treating malformed params as an unmatched pattern.

Review note: `CHAT_CONTEXT/frontend_skills/reviews/03_app_foundation.md`.

## Tested - 2026-05-31 - FE 03 App Routing And API Foundation

Build, route matcher smoke, permission smoke, whitespace check, Vite dev-server startup on
`127.0.0.1:5174`, and HTTP SPA route checks passed. The malformed encoded detail URL regression is
covered by the matcher smoke check; raw HTTP malformed URL is rejected by Vite before React runs.

Browser DOM/manual auth flow was not run because no browser automation or local Chromium is present,
and no backend credentials were provided. FE03 adds no new business API calls.

Test note: `CHAT_CONTEXT/frontend_skills/tests/03_app_foundation.md`.

## Completed - 2026-05-31 - FE 03 App Routing And API Foundation

FE03 is complete as a frontend foundation cycle.

- App routes now come from a shared route registry and matcher.
- `/app` aliases to `/app/dashboard`; `/app/dashboard` renders the existing dashboard.
- Members, subscriptions, attendance, sessions, reports, employees, courses, branches, and payments
  routes render role-aware placeholders only.
- Shared page header, data panel, state block, role helper, and resource-state primitives are in
  place for later modules.
- No backend API contract changed; planned API lists in placeholders are metadata only.

Recommended next action: use `$gym-git` to review/commit the current frontend changes, or
`$gym-fe-plan` to start FE04 Brand Asset Integration / FE05 Members.
