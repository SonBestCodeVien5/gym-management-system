# Review - FE 08 Attendance

## Status

- Status: reviewed
- Feature: Attendance command center and subscription attendance history
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/08_attendance.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/08_attendance.md`
- Reviewed at: 2026-06-02

## Review summary

- Result: issues found
- Build status: passed with `npm run build`
- Test status: browser route check attempted in this review batch, but protected routes redirected to `/login` because no backend/auth session was available.

## Checklist

- [x] UI matches intended design/style.
- [x] Routes and components are scoped cleanly.
- [x] API base URL/env handling is correct.
- [x] Auth/session state is not hardcoded or leaked.
- [x] Loading, empty, and error states are handled where relevant.
- [ ] Responsive layout works on mobile and desktop.
- [ ] Accessibility basics are covered.
- [x] Docs/context are aligned.

## Issues found

| Severity | File | Issue | Fix |
|---|---|---|---|
| medium | `frontend/src/components/attendance/AttendanceHistoryView.jsx:48`, `frontend/src/components/attendance/CheckInPanel.jsx:24`, `frontend/src/components/attendance/ReportMissedPanel.jsx:24`, `frontend/src/components/attendance/MakeupPanel.jsx:28` | Subscription-scoped quick action forms initialize `subscription_id` only once from props. If staff navigates from one `/app/subscriptions/:id/attendance` route to another in the same mounted app, history updates to the new ID but the check-in/report/makeup forms can keep the old subscription ID. | Sync form state when `subscriptionId` prop changes, or key each panel by `subscriptionId` in `AttendanceHistoryView`. |
| low | `frontend/src/components/attendance/CheckInPanel.jsx:59` | Field-level errors are rendered without `aria-describedby`; the plan required visible labels and connected errors. | Add stable error IDs plus `aria-invalid` and `aria-describedby` to invalid fields in check-in/report/makeup panels. |

## Fixes applied during review

- None. Review only.

## Handoff to test

- After fixes, test navigating between two subscription attendance IDs, command success triggering history refresh, invalid ObjectID validation, and backend conflict display.
