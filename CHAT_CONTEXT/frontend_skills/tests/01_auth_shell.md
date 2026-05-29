# Test - 01 Auth Shell

## Status

- Status: retested
- Feature: Staff auth shell
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/01_auth_shell.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/01_auth_shell.md`
- Review file: `CHAT_CONTEXT/frontend_skills/reviews/01_auth_shell.md`
- Tested at: 2026-05-29

## Commands

```bash
cd frontend
npm run build
rg -n "password\\.length|toi thieu 8|min-length|Mat khau can" frontend/src/components/LoginView.jsx
curl -sS -i http://127.0.0.1:5173/login
curl -sS -i http://127.0.0.1:5173/app
curl -sS -i http://127.0.0.1:5173/
curl -sS -i http://127.0.0.1:8080/ping
node -e '<auth API smoke script without printing tokens>'
```

## Command results

| Command | Result | Notes |
|---|---|---|
| `npm run build` | pass | Vite built 39 modules and emitted `dist/index.html`, CSS, and JS assets after the login validation fix. |
| `npm run build` after favicon/wordmark update | pass | Vite built with the new `frontend/public/favicon.svg` and reference-style wordmark changes. |
| `rg -n "password\\.length\|toi thieu 8\|min-length\|Mat khau can" frontend/src/components/LoginView.jsx` | pass | No matches; login form no longer blocks passwords shorter than 8 characters. |
| `curl -sS -i http://127.0.0.1:5173/login` | pass | Vite dev server returned `200 OK` and the React app HTML shell. |
| `curl -sS -i http://127.0.0.1:5173/app` | pass | Vite dev server returned `200 OK` and the React app HTML shell for the protected route. |
| `curl -sS -i http://127.0.0.1:5173/` | pass | Vite dev server returned `200 OK` and the React app HTML shell for the redirect entry route. |
| `curl -sS -i http://127.0.0.1:8080/ping` | pass | Backend started successfully against local `gym_mongodb` after running with escalated local-network access. |
| `node -e '<auth API smoke script without printing tokens>'` | pass | Login `200`, token presence true, `/auth/me` `200`, refresh `200`, logout `200`, wrong password `401 UNAUTHORIZED`, CORS preflight `204` with `allow-origin=http://127.0.0.1:5173`. |

## Manual UI/API checks

- [x] Review finding fixed: login validation no longer enforces password length before API submit.
- [x] Dev server route shell: `/`, `/login`, and `/app` return the Vite React shell.
- [ ] Desktop viewport: not run; no browser automation dependency is present in `frontend/package.json`.
- [ ] Mobile viewport: not run; no browser automation dependency is present in `frontend/package.json`.
- [x] API success path: login, `/auth/me`, refresh, logout.
- [x] API error path: wrong password returns `401` with `UNAUTHORIZED`.
- [x] Mismatched favicon replaced with a new Iron Forge `IF` favicon and restored in `frontend/index.html`.

## Issues found

| Severity | Case | Issue | Fix |
|---|---|---|---|
| none | login validation | Previously open review finding is fixed. | No test-phase code fix needed. |

## Skipped checks

- Full browser login interaction through the React form: not run because no browser automation/headless
  browser is installed in the local frontend toolchain.
- Reload `/app` and observe React route guard visually: not run because no browser automation/headless
  browser is installed in the local frontend toolchain.
- Refresh-token retry through the React boot flow: API refresh was tested directly, but the browser
  boot retry path was not exercised in a real browser.
- Logout button interaction in the React shell: API logout was tested directly, but the browser button
  flow was not exercised in a real browser.
- Backend-down UI error in the React form: not run in a real browser.
- 320px/1280px visual checks: skipped because no browser automation tool is installed in the frontend.

## Final result

- Result: pass for build, static route shell, review-finding regression, backend auth API flow, and CORS
  preflight. Browser visual/interaction checks remain unverified due missing browser automation.
- Ready for `$gym-fe-complete`: yes, with the residual visual/browser-interaction risk noted above.

## Next action

Use `$gym-fe-complete` to finalize FE 01 Auth Shell context/docs, or manually open
`http://127.0.0.1:5173` for a human viewport pass before completion if visual sign-off is required.
