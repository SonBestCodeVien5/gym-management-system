# Implementation - 02 Dashboard Reference

## Status

- Status: implemented
- Feature: Reference-inspired operational dashboard
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/02_dashboard_reference.md`
- Started at: 2026-05-29
- Finished at: 2026-05-29

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/components/AppShell.jsx` - refined sidebar into reference-style groups, added
  disabled dashboard placeholder tool buttons, staff initials avatar, and expanded future module
  placeholders while keeping logout/auth behavior unchanged.
- `frontend/src/components/DashboardHome.jsx` - replaced the auth placeholder dashboard with the new
  operational dashboard composition.
- `frontend/src/components/dashboardData.js` - added static sample dashboard data for KPI cards,
  revenue bars, plan distribution, latest members, and today's sessions.
- `frontend/src/components/KpiCard.jsx` - added reusable KPI card component.
- `frontend/src/components/RevenueBars.jsx` - added accessible CSS bar chart using static sample data.
- `frontend/src/components/PlanDonut.jsx` - added accessible SVG donut chart using static sample data.
- `frontend/src/components/MemberTable.jsx` - added latest registrations table with status badges.
- `frontend/src/components/ScheduleList.jsx` - added today's class schedule list with capacity states.
- `frontend/src/index.css` - added grouped sidebar styling, topbar controls, dashboard grid,
  chart/table/schedule styles, status badges, staff context panel, and responsive layout rules.
- `CHAT_CONTEXT/frontend_skills/worklog.md` - updated FE roadmap and implementation handoff.

## Key decisions

- Kept this cycle frontend-only. No new dashboard/report API calls were added because backend
  aggregate endpoints do not exist yet.
- Put sample data in `dashboardData.js` so a future live-dashboard cycle can replace it without
  rewriting view components.
- Kept `/app` as the only protected route for this feature; deeper business routes remain planned for
  FE 03.
- Kept future module buttons visible by role but disabled, matching the existing auth-shell behavior.
- Used CSS/SVG for charts and small tool symbols instead of adding an icon or chart dependency.
- Kept live staff identity from `AuthContext`; only dashboard metrics are static sample data.

## Commands run

```bash
cd frontend
npm run build
```

Result: pass. Vite built 29 modules and emitted production assets.

```bash
git diff --check
```

Result: pass.

## Review fix - 2026-05-29

- Fixed the FE02 review finding in `AppShell.jsx` by marking the notification and search topbar
  placeholders as disabled controls with "coming soon" accessible labels.
- Added disabled styling for topbar placeholder buttons in `index.css` so the visual state matches
  the non-interactive behavior.
- Re-ran `npm run build` and `git diff --check`; both passed.

## Known limitations

- Dashboard KPIs, revenue, plan distribution, recent members, and class schedule are static sample
  data.
- Browser visual checks at 320px and 1280px were not run in this implementation phase.
- Tool buttons in the dashboard topbar are disabled visual placeholders until those features exist.
- Navigation still does not open business modules; route foundation remains FE 03.

## Handoff to review

- Verify the dashboard does not imply sample metrics are live API data.
- Check responsive layout for KPI grid, chart row, table overflow, schedule list, and staff topbar.
- Review accessibility for chart labels, donut `role="img"`, real table semantics, badges, and
  disabled navigation.
- Confirm auth restore/logout behavior in `AppShell` stayed unchanged.
- Use `$gym-fe-review` with this implementation note.
