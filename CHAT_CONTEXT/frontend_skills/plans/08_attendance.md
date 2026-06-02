# FE Plan - 08 Attendance

Status: Planned

Created: 2026-06-02

## Goal

Build the attendance command center for free check-in, missed-session reporting, makeup attendance,
and subscription attendance history.

FE08 should use the existing attendance APIs and keep lookup workflows explicit where backend search
APIs are absent. The first pass should support direct subscription/member context and branch
selection from FE06 when available.

## Current Baseline

- Existing route:
  - `/app/attendance`
- Existing route registry does not yet include:
  - `/app/subscriptions/:id/attendance`
- FE05 provides member detail and member-scoped subscriptions.
- FE07 is expected to provide subscription detail lookup.
- Branch selection should use FE06 branch list helpers when available.

## Screens And Routes

| Route | Access | Behavior |
|---|---|---|
| `/app/attendance` | `admin`, `manager`, `receptionist` | Command center with check-in, report missed, makeup, and history lookup panels. |
| `/app/subscriptions/:id/attendance` | `admin`, `manager`, `receptionist` | Attendance history for one subscription, plus quick report/makeup/check-in actions scoped to that subscription. |

Route config changes:

- Add `/app/subscriptions/:id/attendance` before `/app/subscriptions/:id`.
- Keep `/app/attendance` in the existing `Tong quan` group.
- Mark attendance route ready after implementation.

## Component Plan

Add or update:

| Path | Responsibility |
|---|---|
| `src/lib/attendanceApi.js` | Check-in/report/makeup/history helpers. |
| `src/components/attendance/AttendancePage.jsx` | Command center with panels/tabs. |
| `src/components/attendance/AttendanceHistoryView.jsx` | Subscription-scoped history route. |
| `src/components/attendance/CheckInPanel.jsx` | Free check-in form. |
| `src/components/attendance/ReportMissedPanel.jsx` | Report missed form. |
| `src/components/attendance/MakeupPanel.jsx` | Makeup form referencing reported missed date. |
| `src/components/attendance/AttendanceHistoryPanel.jsx` | Fetch/render attendance rows. |
| `src/components/attendance/SubscriptionHistoryLookup.jsx` | Direct subscription ID lookup. |
| `src/components/attendance/attendanceFormatters.js` | Status/date/ObjectID helpers. |
| `src/App.jsx` | Render FE08 route components. |
| `src/routes/routeConfig.js` | Add subscription attendance route and mark attendance route ready. |
| `src/index.css` | Scoped attendance command/history responsive styles. |

## State And API Plan

API helpers:

```js
checkInAttendance(accessToken, payload)
reportMissedAttendance(accessToken, payload)
createMakeupAttendance(accessToken, payload)
listSubscriptionAttendance(accessToken, subscriptionId)
```

Check-in payload:

- `subscription_id`
- `branch_id`
- `date` optional RFC3339
- `status` should be `attended` for free check-in
- `session_id` optional only when context requires it; FE09 owns session check-in
- `is_makeup_for` optional only for special cases; dedicated makeup endpoint preferred

Report payload:

- `subscription_id`
- `branch_id`
- `date` optional RFC3339

Makeup payload:

- `subscription_id`
- `branch_id`
- `date` optional RFC3339
- `is_makeup_for` required RFC3339

History:

- `GET /api/v1/subscriptions/:id/attendance`

State:

- Each command panel owns form and mutation state.
- History owns query state and refetch callback.
- Command success can optionally trigger history refetch when scoped to a subscription detail route.
- Branch options use FE06 branch API when available; fallback to manual branch ID if branch list fails.

Validation:

- Subscription and branch IDs must be ObjectIDs.
- Dates must be valid RFC3339 when supplied.
- Makeup `is_makeup_for` required.
- Do not client-enforce weekly limits or makeup windows beyond obvious date shape; backend is source
  of truth.

## UX States

- Command center starts with empty forms and direct subscription lookup.
- Branch list loading/error/fallback.
- Check-in success and `409` weekly limit/no remaining sessions.
- Report missed success and `409` reported-missed 30-day limit.
- Makeup success and `409` invalid/reused/not-found makeup reference.
- History loading, empty, invalid ID, and backend error.
- Unknown subscription history route should show `404`/error state.

## Responsive And Accessibility Notes

- Use compact segmented controls or panels for command modes; no marketing layout.
- At 320px, command panels stack and form actions become full-width.
- History rows should become stacked records on mobile.
- Labels visible for all date/ID fields.
- Success/error messages use `aria-live`.
- Date input helper text should state that backend expects RFC3339; implementation may use
  `datetime-local` and convert to RFC3339.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/08_attendance.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/08_attendance.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/08_attendance.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Checks:

- Command center render and role guard.
- Check-in/report/makeup validation.
- History empty/success/error.
- `409` conflict display from backend or mocked API.
- Mobile and desktop route layout.

## Backend Contract Gaps

- No global subscription/member search endpoint; direct subscription ID remains required.
- Branch list exists, but branch selection depends on FE06 implementation.
- FE09 owns session-specific enroll/check-in flow; FE08 should not duplicate that workflow.

## Risks And Boundaries

- Do not create sessions or enroll subscriptions in FE08.
- Do not assume `reported_missed` source lookup exists beyond attendance history.
- Avoid complex calendar UI; use simple date fields until UX hardening.
- Backend conflict messages must be surfaced without inventing client-side business rules.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/08_attendance.md`.
