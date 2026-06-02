# FE Plan - 09 Sessions

Status: Planned

Created: 2026-06-02

## Goal

Build the session scheduling workspace for listing/filtering sessions, creating sessions, enrolling
subscriptions, and checking in enrolled subscriptions.

FE09 should use live session APIs, branch/course reference context from FE06, and direct
subscription lookup from FE07/FE08. It should not depend on FE10 employee management for trainer
selection because employee listing is admin-only while session routes are also available to managers
and trainers.

## Current Baseline

- Existing routes:
  - `/app/sessions`
  - `/app/sessions/:id`
- Route registry does not yet include `/app/sessions/new`.
- Session list API supports filters:
  - `branchId`
  - `level`
  - `date` RFC3339
- Session create requires `trainer_id`; no trainer search endpoint is available to non-admins.

## Screens And Routes

| Route | Access | Behavior |
|---|---|---|
| `/app/sessions` | `admin`, `manager`, `trainer` | Filterable session list, create action, empty/error states. |
| `/app/sessions/new` | `admin`, `manager`, `trainer` | Create session form. |
| `/app/sessions/:id` | `admin`, `manager`, `trainer` | Session detail, enrollment list, enroll subscription, and session check-in. |

Route config changes:

- Add `/app/sessions/new` before `/app/sessions/:id`.
- Mark sessions routes ready after implementation.

## Component Plan

Add or update:

| Path | Responsibility |
|---|---|
| `src/lib/sessionsApi.js` | Create/list/get/enroll/check-in helpers. |
| `src/components/sessions/SessionsPage.jsx` | Filterable list and command center. |
| `src/components/sessions/SessionCreateView.jsx` | Create session form. |
| `src/components/sessions/SessionDetailView.jsx` | Detail fetch and action panels. |
| `src/components/sessions/SessionFilters.jsx` | Branch/level/date filters. |
| `src/components/sessions/SessionList.jsx` | Table/list responsive rendering. |
| `src/components/sessions/EnrollmentPanel.jsx` | Enroll direct subscription ID. |
| `src/components/sessions/SessionCheckInPanel.jsx` | Check in enrolled subscription ID. |
| `src/components/sessions/sessionFormatters.js` | Date/duration/capacity/status/ObjectID helpers. |
| `src/App.jsx` | Render FE09 route components. |
| `src/routes/routeConfig.js` | Add `/app/sessions/new`; mark session routes ready. |
| `src/index.css` | Scoped session list/filter/action responsive styling. |

## State And API Plan

API helpers:

```js
listSessions(accessToken, { branchId, level, date })
createSession(accessToken, payload)
getSession(accessToken, sessionId)
enrollSubscription(accessToken, sessionId, subscriptionId)
checkInSessionSubscription(accessToken, sessionId, subscriptionId)
```

Create payload:

- `branch_id`
- `trainer_id`
- `course_level`
- `scheduled_at` RFC3339
- `duration_min`
- `capacity`
- `tags`

List query:

- `branchId` ObjectID
- `level`
- `date` RFC3339

State:

- Session list owns filter and query state.
- Create form owns reference data, form, and submit state.
- Detail owns session query state plus enroll/check-in mutation states.
- After enroll or check-in, refetch session detail.

Reference data:

- Branch list from FE06 when available; fallback to manual branch ID.
- Course levels/tags can be taken from course list when available, but create form must allow manual
  tags because sessions store tags directly.
- Trainer ID:
  - default to current employee ID when current role includes `trainer`.
  - allow manual trainer ObjectID for managers/admins.
  - do not rely on employee list unless current user is admin and FE10 is already implemented.

Validation:

- ObjectID validation for branch, trainer, session, subscription IDs.
- `scheduled_at` valid RFC3339.
- duration and capacity positive integers.
- course_level required.
- tags parsed into string array.
- Enroll/check-in subscription ID required ObjectID.

## UX States

- List loading, empty, error.
- Filter invalid date or branch ID.
- Create success navigates to detail.
- Detail invalid ID, not found, loading, success.
- Capacity display: enrolled count vs capacity.
- Enroll conflicts:
  - already enrolled
  - full session
  - tag not allowed
- Check-in conflicts:
  - not enrolled
  - check-in closed
  - subscription/attendance business conflicts
- Successful enroll/check-in notice and refreshed detail.

## Responsive And Accessibility Notes

- List filters wrap on tablet and stack on mobile.
- Session rows become stacked records at narrow widths.
- Create/detail action panels use full-width buttons at 320px.
- Enrollment/check-in forms must be keyboard-submittable.
- Use visible labels and connected field errors.
- Success/error messages use `aria-live`.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/09_sessions.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/09_sessions.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/09_sessions.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Checks:

- Session list filters.
- Create session.
- Detail fetch.
- Enroll success/conflict.
- Check-in success/conflict.
- Mobile and desktop layout.

## Backend Contract Gaps

- No trainer search/list endpoint available to manager/trainer roles.
- No course-level enum endpoint; levels/tags come from existing course list or manual entry.
- Session list query uses `branchId` camelCase, not `branch_id`.

## Risks And Boundaries

- Do not implement employee management in FE09.
- Do not duplicate free attendance check-in from FE08; session check-in uses session endpoint.
- Avoid calendar drag/drop; list/filter UI is enough for this cycle.
- Role guard in frontend is UX only; backend remains security source of truth.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/09_sessions.md`.
