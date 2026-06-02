# Implementation - 11 Live Dashboard APIs

## Status

- Status: implemented
- Feature: FE11 Live Dashboard APIs
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/11_live_dashboard_apis.md`
- Started at: 2026-06-02
- Finished at: 2026-06-02

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/lib/dashboardApi.js` - added dashboard endpoint helpers.
- `frontend/src/components/DashboardHome.jsx` - replaced static sample data source with live dashboard
  API loading/error/section-error state and response mapping.
- `frontend/src/components/RevenueBars.jsx` - removed sample wording and rendered live formatted values.
- `frontend/src/components/PlanDonut.jsx` - removed sample wording and changed total label to
  subscriptions.
- `frontend/src/components/MemberTable.jsx` - removed sample wording and keyed rows by member ID.
- `frontend/src/components/ScheduleList.jsx` - removed sample wording and keyed rows by session ID.
- `frontend/src/routes/routeConfig.js` - updated dashboard metadata and planned API list.
- `CHAT_CONTEXT/frontend_skills/implementations/11_live_dashboard_apis.md` - implementation handoff.
- `CHAT_CONTEXT/frontend_skills/worklog.md` - implementation chronology.

## Key decisions

- Kept the existing dashboard layout and visual components, mapping backend responses into their local
  display shape.
- Dashboard metrics now require admin/manager access. Trainer/receptionist roles see a dashboard
  access denied state instead of static sample metrics.
- Requested dashboard sections independently with `Promise.allSettled` so revenue/plan/member/session
  failures can show section alerts while summary remains usable.
- Kept `/app/reports` blocked because backend added dashboard endpoints, not broader report/export
  endpoints.
- FE12 removed `dashboardData.js` after confirming it was no longer imported by production source.

## Review Fixes

- `frontend/src/components/DashboardHome.jsx` now keeps the dashboard in stale state when a later
  summary refresh fails after a successful snapshot, instead of dropping straight to full error.
- `frontend/src/components/DashboardHome.jsx` and `frontend/src/components/MemberTable.jsx` now
  render recent members with real backend fields only: name, level, joined time, and registration
  status. The fake trainer column is gone.
- `frontend/src/components/DashboardHome.jsx`, `frontend/src/components/RevenueBars.jsx`, and
  `frontend/src/index.css` now preserve negative revenue values and render them below a zero
  baseline instead of clamping them to empty-state behavior.

## Commands run

```bash
npm run build
git diff --check
```

## Command results

- `npm run build` passed.
- `git diff --check` passed.
- Mocked Playwright dashboard smoke passed for live admin data, recent member labels, and negative
  revenue rendering. Stale refresh state could not be forced cleanly from the browser because the
  dashboard currently has no user-triggered refresh action.

## Known limitations

- Live browser/API smoke was not run in this FE11 implementation pass.
- Dashboard requests currently use global backend defaults and do not expose branch/date filters in
  the UI.
- No report landing page was added.
- The FE11 stale-state branch was fixed in code and retested at build level, but not forced through
  a user-visible refresh action because the current UI has no explicit refresh control.

## Handoff to review

- FE11 review findings are fixed in code.
- Re-run `$gym-fe-test` for a final browser/build pass.
