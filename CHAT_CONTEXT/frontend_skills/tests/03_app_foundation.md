# Test - 03 App Routing And API Foundation

## Status

- Status: tested
- Feature: 03 App Routing And API Foundation
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/03_app_foundation.md`
- Review file: `CHAT_CONTEXT/frontend_skills/reviews/03_app_foundation.md`
- Tested at: 2026-05-31

## Commands

```bash
cd frontend
npm run build
```

```bash
cd frontend
node --input-type=module -e "<route matcher smoke check>"
```

```bash
cd frontend
node --input-type=module -e "<permission smoke check>"
```

```bash
git diff --check
```

```bash
cd frontend
npm run dev -- --host 127.0.0.1 --port 5174
```

```bash
curl -sS -o /tmp/fe03-*.html -w "%{http_code} %{content_type}\n" <route>
```

## Command results

| Command | Result | Notes |
|---|---|---|
| `npm run build` | pass | Vite built 36 modules and emitted production assets. |
| route matcher smoke check | pass | Covered `/`, `/login`, `/app`, `/app/dashboard`, module routes, detail params, unknown `/app/*`, non-app redirect, and malformed encoded detail URL. |
| permission smoke check | pass | Covered dashboard, members, sessions, reports, employees, courses, branches role outcomes. |
| `git diff --check` | pass | No whitespace errors. |
| `npm run dev -- --host 127.0.0.1 --port 5173` | blocked | Sandbox bind failed first with `EPERM`; escalated retry found port `5173` already in use. |
| `npm run dev -- --host 127.0.0.1 --port 5174` | pass | Vite dev server started at `http://127.0.0.1:5174/`; stopped after checks. |

## Dev Server Route Checks

| Route | Result | Notes |
|---|---|---|
| `/` | `200 text/html` | SPA entry served. |
| `/login` | `200 text/html` | SPA entry served. |
| `/app` | `200 text/html` | SPA entry served for protected alias route. |
| `/app/dashboard` | `200 text/html` | SPA entry served for dashboard route. |
| `/app/members` | `200 text/html` | SPA entry served for placeholder route. |
| `/app/subscriptions/sub-1` | `200 text/html` | SPA entry served for detail placeholder route. |
| `/app/unknown` | `200 text/html` | SPA entry served so app-level not-found can render. |
| `/app/members/%E0%A4%A` | `404` | Vite rejected the raw malformed request URL before React ran; the matcher-level malformed URL behavior is covered by the Node smoke check. |

## Manual UI/API checks

- [x] Desktop route coverage: verified by matcher smoke and dev-server HTTP checks.
- [x] Mobile route coverage: route behavior is viewport-independent; no browser viewport screenshot was run.
- [x] API success path: not applicable for FE03 business modules because this cycle adds no new business API calls.
- [x] API error path: not applicable for FE03 business modules; auth/backend error UI was not exercised.
- [x] Review finding regression: malformed encoded detail URL no longer throws in `matchRoute()`.

## Skipped

- Browser DOM/manual auth flow was not run because no browser automation dependency or local Chromium
  binary is present in the frontend environment.
- Backend login/current-user/refresh/logout flow was not run because FE03 does not add new auth API
  behavior and no local backend credentials were provided for this test pass.

## Issues found

| Severity | Case | Issue | Fix |
|---|---|---|---|
| none | FE03 verification | No blocking issue found in build, matcher, permission, whitespace, or dev-server route checks. | Not applicable. |

## Final result

- Result: pass with skipped browser/auth manual checks noted.
- Ready for `$gym-fe-complete`: yes.
