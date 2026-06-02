# Review - 11 Live Dashboard APIs

## Status

- Status: reviewed with findings
- Feature: FE11 Live Dashboard APIs
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/11_live_dashboard_apis.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/11_live_dashboard_apis.md`
- Reviewed at: 2026-06-02

## Review summary

- Result: findings found
- Build status: `npm run build` passed
- Test status: mocked browser review passed for happy path, section-error alert, forbidden role, and
  mobile no-overflow; live backend browser smoke was not run

## Checklist

- [x] UI matches intended design/style.
- [x] Routes and components are scoped cleanly.
- [x] API base URL/env handling is correct.
- [x] Auth/session state is not hardcoded or leaked.
- [x] Loading, empty, and error states are handled where relevant.
- [x] Responsive layout works on the reviewed mobile dashboard viewport.
- [x] Accessibility basics are covered for reviewed states.
- [x] Docs/context are aligned.

## Issues found

| Severity | File | Issue | Fix |
|---|---|---|---|
| medium | `frontend/src/components/DashboardHome.jsx:151` and `frontend/src/components/MemberTable.jsx:11` | Recent member rows now come from `/dashboard/members/recent`, whose contract exposes member identity/status fields only. The mapper fills the table's `Plan` column from member `level` and the `Trainer` column from compact member ID, so the live dashboard labels backend data as plan/trainer data it does not have. This violates the FE11 boundary to avoid fake live metrics. | Either change `MemberTable` for dashboard use to columns like `Name`, `Level`, `Member ID`, `Status`, or extend the backend contract to return real plan/trainer fields before showing those labels. |
| medium | `frontend/src/components/DashboardHome.jsx:130` and `frontend/src/components/RevenueBars.jsx:6` | Revenue buckets clamp negative `net_amount` to zero for chart sizing, and all-negative/refund-only ranges render as "No net revenue recorded." Since backend defines net revenue as payments minus refunds, a negative period is valid data and should not be hidden as empty. | Preserve negative bucket state in the chart/summary copy, or render a distinct refund/net-negative state instead of treating it as no data. |
| low | `frontend/src/components/DashboardHome.jsx:233` | The stale-data branch is fragile: `loadDashboard` first changes `status` from `success` to `refreshing`, then summary failure checks `current.status === 'success'`, so a summary failure during a later refresh would fall to full `error` instead of keeping the existing dashboard as stale. There is no manual refresh yet, so impact is currently limited. | Check whether `current.data.summary` exists instead of checking only `current.status === 'success'`, and render a stale alert when summary refresh fails after prior data exists. |

## Browser checks

- Started Vite on `http://127.0.0.1:5175/` after port `5174` was busy.
- Mocked admin happy path:
  - `/app/dashboard` rendered Dashboard and live revenue value.
  - Mobile `390x844` had no page-level horizontal overflow (`scrollWidth=390`, `clientWidth=390`).
- Mocked section failure:
  - summary succeeded, revenue returned `500`.
  - dashboard still rendered summary KPIs and showed "Some dashboard sections could not refresh."
- Mocked receptionist:
  - `/app/dashboard` rendered `Dashboard access denied`.
- Dev server was stopped after review; final curl to `5175` failed to connect as expected.

## Fixes applied during review

- None. Review only.

## Remaining risks

- Live backend browser/API smoke with seeded credentials was not run.
- The full FE12 route/viewport matrix is still pending.
- Browser console logged expected errors during the mocked `500` section-failure case.

## Handoff to test

- Fix the medium findings first.
- Retest `npm run build`.
- Re-run mocked dashboard browser checks for:
  - live happy path
  - recent members table labels
  - refund-only or negative net revenue range
  - section-error alert
  - receptionist forbidden state
