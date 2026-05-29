# Test - 02 Dashboard Reference

## Status

- Status: tested with local limitations
- Feature: Reference-inspired operational dashboard
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/02_dashboard_reference.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/02_dashboard_reference.md`
- Review file: `CHAT_CONTEXT/frontend_skills/reviews/02_dashboard_reference.md`
- Tested at: 2026-05-29

## Commands

```bash
cd frontend
npm run build
```

```bash
cd frontend
npm run dev -- --host 127.0.0.1 --port 5173
```

```bash
curl -sS -i http://127.0.0.1:5173/
curl -sS -i http://127.0.0.1:5173/login
curl -sS -i http://127.0.0.1:5173/app
```

```bash
cd frontend
node --input-type=module -e "<DashboardHome SSR smoke check>"
```

```bash
curl -sS -i http://127.0.0.1:18080/api/v1/auth/me
git diff --check
```

## Command Results

| Command | Result | Notes |
|---|---|---|
| `npm run build` | pass | Vite built 29 modules and emitted production assets. |
| `npm run dev -- --host 127.0.0.1 --port 5173` | pass with escalation | Sandbox blocked local listen with `EPERM`; running with approved local network permission started Vite at `http://127.0.0.1:5173/`. |
| `curl` `/`, `/login`, `/app` against Vite | pass | All returned `HTTP/1.1 200 OK` and the SPA HTML with `/src/main.jsx`. |
| `DashboardHome` SSR smoke check | pass with sandbox warning | Rendered dashboard sample content, KPI labels, latest member table, schedule section, staff context, and donut accessible label. Vite emitted a sandbox `WebSocket server` listen warning, but the render assertion passed. |
| `rg` static checks for topbar review fix | pass | Confirmed `Notifications coming soon` and `Search coming soon` disabled controls in `AppShell.jsx`. |
| `rg` static checks for responsive/table CSS | pass | Confirmed disabled topbar styling, table horizontal overflow, and mobile single-column CSS rules exist. |
| `curl` backend `/api/v1/auth/me` on `127.0.0.1:18080` | skipped/blocked | Backend was not listening locally, so real login/session/logout flow could not be verified in this test pass. |
| `git diff --check` | pass | No whitespace errors. |

## Manual UI/API Checks

- [x] Build: production build passes.
- [x] Vite dev server: starts locally after sandbox network permission.
- [x] Route smoke: `/`, `/login`, and `/app` are served by Vite as SPA routes.
- [x] Dashboard render smoke: dashboard sections and sample data render through React SSR.
- [x] Review fix smoke: notification/search topbar placeholders are disabled and labeled as coming soon.
- [x] Static responsive checks: CSS has mobile single-column rules and table-contained horizontal scroll.
- [ ] Desktop browser viewport: not run; no browser automation package is installed in the frontend.
- [ ] Mobile browser viewport: not run; no browser automation package is installed in the frontend.
- [ ] API success path login/session restore/logout: not run; backend was not listening on `127.0.0.1:18080`.
- [ ] API error path: not run; backend was not listening on `127.0.0.1:18080`.

## Issues Found

| Severity | Case | Issue | Fix |
|---|---|---|---|
| none | build/render smoke | No new issue found in feasible checks. | |

## Final Result

- Result: complete with local limitations
- Ready for `$gym-fe-complete`: yes

The frontend implementation passes build, route smoke, component render smoke, review-fix static
checks, and whitespace checks. Browser visual checks at 320px and desktop width, plus authenticated
API flow verification, are still recorded as local limitations because browser automation and a
seeded backend were not available during this pass.

## Retest - 2026-05-29

Rechecked the FE02 slice after the latest handoff.

### Commands

```bash
cd frontend
npm run build
```

```bash
cd frontend
npm run dev -- --host 127.0.0.1 --port 5173
```

```bash
curl -sS -i http://127.0.0.1:5173/
curl -sS -i http://127.0.0.1:5173/login
curl -sS -i http://127.0.0.1:5173/app
```

### Results

- `npm run build` passed.
- Vite dev server started on `http://127.0.0.1:5173/`.
- `/`, `/login`, and `/app` all returned `HTTP/1.1 200 OK` and the SPA HTML shell.
- No new frontend code issue was found in this recheck.

### Still open

- Desktop and mobile browser viewport checks are still not run because local browser automation is not installed.
- Authenticated API flow verification still depends on a running backend with seeded login data.
- These gaps are recorded as known limitations for the FE02 completion handoff.

## Completion Note

FE02 is complete as a frontend-only dashboard cycle: the protected `/app` dashboard now uses the
reference-inspired static dashboard composition, live staff identity still comes from auth state,
and the remaining browser/backend checks are captured as follow-up limitations rather than blocking
issues.

## Manual Checklist

Use this when the backend is running locally with `PORT=8080` and the frontend runs on `127.0.0.1:5173`.

### Prerequisites

- Backend started with the local `.env` values.
- Bootstrap admin account available from `BOOTSTRAP_ADMIN_EMAIL`.
- Frontend dev server running with `npm run dev -- --host 127.0.0.1 --port 5173`.

### Auth Flow

- [ ] Open `http://127.0.0.1:5173/login`.
- [ ] Log in with the bootstrap admin email from `.env` and the matching password.
- [ ] Confirm the app redirects to `/app` after login.
- [ ] Refresh the page and confirm the session is restored through `GET /api/v1/auth/me`.
- [ ] Log out and confirm the app returns to `/login`.
- [ ] Enter a wrong password and confirm a visible login error appears.

### Dashboard Layout

- [ ] Open `/app` and confirm the hero, KPI grid, revenue panel, donut panel, latest members table, and today schedule are visible.
- [ ] Confirm the notification and search topbar controls are disabled and labeled as coming soon.
- [ ] Confirm the staff context panel shows the live employee identity from the auth session.
- [ ] Confirm the dashboard copy makes it clear the metrics are sample data, not live API data.

### Desktop View

- [ ] Check the layout at a wide viewport around `1280px`.
- [ ] Confirm the KPI cards stay in a multi-column grid.
- [ ] Confirm the chart row sits side by side without overflow.
- [ ] Confirm the members table remains readable and does not break the page width.
- [ ] Confirm the schedule panel remains aligned with the rest of the dashboard.

### Mobile View

- [ ] Check the layout at a narrow viewport around `320px`.
- [ ] Confirm the dashboard collapses to a single-column stack.
- [ ] Confirm the members table can scroll horizontally instead of breaking the layout.
- [ ] Confirm the sidebar/topbar remain usable and do not overlap the content.
- [ ] Confirm text still fits without clipping in the hero, KPI cards, and schedule list.

### Pass Criteria

- [ ] Build passes with `npm run build`.
- [ ] Login, restore, and logout all work against the local backend.
- [ ] No visible layout breakage at desktop or mobile width.
- [ ] No misleading live-data claim appears in the dashboard.
