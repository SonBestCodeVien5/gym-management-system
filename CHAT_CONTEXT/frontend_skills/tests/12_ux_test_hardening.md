# Test - 12 UX/Test Hardening

## Status

- Status: tested with limitations
- Feature: FE12 UX/Test Hardening
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/12_ux_test_hardening.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/12_ux_test_hardening.md`
- Review file: not run in this sequence
- Tested at: 2026-06-02

## Commands

```bash
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
npm run dev -- --host 127.0.0.1 --port 5174
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
| MCP Playwright mocked browser smoke | passed with limitation | Desktop dashboard, mobile no-overflow/classes, and receptionist forbidden state passed. |
| `git diff --check` | passed | No whitespace errors. |
| final `curl` to `127.0.0.1:5174` | failed to connect as expected | Confirmed scoped Vite shutdown. |

## Browser checks

- [x] Desktop mocked admin dashboard rendered live API data and stayed on `/app/dashboard`.
- [x] Mobile `390x844` mocked admin dashboard had no page-level horizontal overflow.
- [x] Mobile classes expandable content rendered.
- [x] Mocked receptionist dashboard rendered `Dashboard access denied`.
- [x] Console warning/error check returned no messages.
- [ ] Full resource route viewport matrix was not run in this pass.
- [ ] Live backend browser/API checks were not run with seeded credentials.
- [ ] Permanent Playwright automation was not added.

## Issues found

| Severity | Case | Issue | Fix |
|---|---|---|---|
| low | stale source cleanup | `dashboardData.js` was no longer imported after FE11. | Deleted the unused static sample fixture. |
| low | mobile member panel smoke | The Members expandable button was not conclusively targeted because the accessible name collided with other Members text. | Recorded as follow-up for a fuller FE12 automation pass. |

## Final result

- Result: build, targeted mocked dashboard browser smoke, no-overflow check, forbidden role check, and
  whitespace check passed.
- Ready for `$gym-fe-complete`: no. Full FE12 completion still needs broader route viewport matrix or
  a repeatable Playwright test suite.
