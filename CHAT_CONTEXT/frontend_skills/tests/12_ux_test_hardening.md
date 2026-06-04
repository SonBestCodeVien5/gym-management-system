# Test - 12 UX/Test Hardening

## Status

- Status: tested with limitations
- Feature: FE12 UX/Test Hardening MVP completion
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/12_ux_test_hardening.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/12_ux_test_hardening.md`
- Review file: `CHAT_CONTEXT/frontend_skills/reviews/12_ux_test_hardening.md`
- Tested at: 2026-06-04

## Commands

```bash
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
npm run dev -- --host 127.0.0.1 --port 5174
npm run dev -- --host 127.0.0.1 --port 5178
npm run dev -- --host 127.0.0.1 --port 5184
MCP Playwright browser tab listing
Chromium/CDP mocked dashboard smoke on 127.0.0.1:5184
git diff --check
curl -sS -i http://127.0.0.1:5174/app/dashboard
```

## Command results

| Command | Result | Notes |
|---|---|---|
| `npm run build` | passed | Ran after FE11 and after deleting `dashboardData.js`. |
| `npm run dev -- --host 127.0.0.1 --port 5173` | failed in sandbox | `EPERM` binding localhost. |
| escalated `npm run dev -- --host 127.0.0.1 --port 5173` | failed | Port `5173` was already in use. |
| escalated `npm run dev -- --host 127.0.0.1 --port 5174` | passed | Vite served on `http://127.0.0.1:5174/`. |
| escalated `npm run dev -- --host 127.0.0.1 --port 5178` | passed | Vite served the FE12 follow-up smoke. |
| `npm run dev -- --host 127.0.0.1 --port 5184` | passed after escalation | Sandbox listen failed with `EPERM`; escalated Vite served the final FE12 test smoke. |
| MCP Playwright browser tab listing | blocked | MCP browser profile was locked before navigation (`mcp-chrome-for-testing-64b1e2a` already in use). |
| Chromium/CDP mocked dashboard smoke | passed | Desktop dashboard, stale refresh, mobile no-overflow/members/sessions controls, receptionist forbidden state, and console/exception checks passed. |
| `git diff --check` | passed | No whitespace errors. |
| final `curl` to `127.0.0.1:5174` | failed to connect as expected | Confirmed scoped Vite shutdown. |

## Browser checks

- [x] Desktop mocked admin dashboard rendered live API data and stayed on `/app/dashboard`.
- [x] Desktop dashboard refresh button was visible and exercised a stale refresh failure path.
- [x] Refresh summary failure after a successful snapshot kept the prior revenue visible and rendered the stale-data alert.
- [x] Mobile `390x844` mocked admin dashboard had no page-level horizontal overflow.
- [x] Mobile `390x844` measured `scrollWidth=390` and `clientWidth=390`.
- [x] Mobile members expandable content rendered.
- [x] Mobile sessions expandable content rendered.
- [x] Mobile members and sessions expand buttons had distinct accessible names.
- [x] Mocked receptionist dashboard rendered `Dashboard access denied`.
- [x] Console warning/error check returned no messages.
- [ ] Full resource route viewport matrix was not run in this pass.
- [ ] Live backend browser/API checks were not run with seeded credentials.
- [ ] Permanent Playwright automation was not added.

## Issues found

| Severity | Case | Issue | Fix |
|---|---|---|---|
| low | stale source cleanup | `dashboardData.js` was no longer imported after FE11. | Deleted the unused static sample fixture. |
| low | stale refresh path | The live dashboard could refresh stale, but the branch was not visible from the UI. | Added a manual refresh button to the dashboard header. |
| low | mobile panel naming | The Members and Sessions expand buttons were not reliably distinguishable to automation. | Added distinct accessible labels and controls. |

## Final FE12 smoke - 2026-06-04

- `npm run build` passed.
- `git diff --check` passed.
- MCP Playwright was attempted first but remained blocked by a locked browser profile.
- Temporary Chromium/CDP fallback passed against Vite on `127.0.0.1:5184` with mocked auth/API:
  - desktop admin dashboard stayed on `/app/dashboard`, rendered live revenue, and exposed Refresh;
  - manual refresh failed the summary endpoint after prior success and kept the old snapshot visible;
  - mobile `390x844` had no page-level horizontal overflow;
  - mobile Members and Sessions panels were independently targetable by accessible name;
  - mocked receptionist role rendered `Dashboard access denied`;
  - console warnings/errors and runtime exceptions were empty.

## Final result

- Result: build, targeted mocked dashboard browser smoke, stale refresh smoke, no-overflow check,
  forbidden role check, and whitespace check passed.
- Ready for `$gym-fe-complete`: yes, as an MVP hardening completion with documented deferred
  follow-ups.
