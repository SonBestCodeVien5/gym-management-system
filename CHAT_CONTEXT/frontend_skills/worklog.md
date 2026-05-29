# Frontend Worklog

Use this file for short frontend roadmap and completion summaries.

## Current frontend roadmap

- [x] React/Vite scaffold with Iron Forge style
- [x] Auth UI and API client integration
- [x] Dashboard layout and protected routing
- [ ] Reference-inspired operational dashboard
- [ ] Core member/subscription workflows

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
