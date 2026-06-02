# FE Plan - 12 UX/Test Hardening

Status: Planned

Created: 2026-06-02

## Goal

Run a final frontend hardening cycle across the React/Vite staff console after FE05-FE10 resource
screens are implemented.

FE12 should convert the current build/mocked-route evidence into a repeatable browser verification
workflow, fix responsive and accessibility issues found during that pass, and leave a clear
acceptance matrix for future frontend work. This cycle is frontend-focused and should not change
backend API behavior.

## Current Baseline

- Stack: React 18 + Vite 8 with no current frontend test dependency.
- Production build has passed across FE05-FE10.
- Mocked Playwright/browser checks have covered:
  - FE05 members route rendering
  - FE06 courses/branches route rendering
  - FE07 subscription mutation success plus stale refresh-failure alert
  - FE08 attendance route rendering
  - FE09 session mutation success plus stale refresh-failure alert on `390x844`
  - FE10 employees route rendering
- Live backend CRUD/action smokes remain skipped because no seeded credentials/session data were
  available.
- FE02.1 was a temporary responsive containment pass; final shell/dashboard/resource responsive
  acceptance is still deferred to FE12.

## Screens And Routes

FE12 should verify and harden these existing routes rather than add a new product route:

| Route group | Routes |
|---|---|
| Auth/shell | `/`, `/login`, `/app/dashboard`, app not-found route |
| Members | `/app/members`, `/app/members/new`, `/app/members/:id` |
| Courses | `/app/settings/courses`, `/app/settings/courses/:id` |
| Branches | `/app/settings/branches`, `/app/settings/branches/:id` |
| Subscriptions | `/app/subscriptions`, `/app/subscriptions/new`, `/app/subscriptions/:id` |
| Attendance | `/app/attendance`, `/app/subscriptions/:id/attendance` |
| Sessions | `/app/sessions`, `/app/sessions/new`, `/app/sessions/:id` |
| Employees | `/app/employees`, `/app/employees/new`, `/app/employees/:id` |

Route config changes:

- No new route is expected.
- Keep blocked routes (`reports`, `payments`) visually clear and keyboard reachable if they remain in
  navigation.
- If a route is discovered as unreachable because of route ordering or role metadata, fix route config
  and add it to the test matrix.

## Component Plan

Potential frontend/test files:

| Path | Responsibility |
|---|---|
| `package.json` | Add repeatable browser-test scripts if the project accepts a test dependency. |
| `playwright.config.js` | Preferred browser automation config for Vite app checks. |
| `tests/e2e/` or `src/tests/e2e/` | Mocked-auth route, interaction, responsive, and accessibility smoke specs. |
| `tests/fixtures/` | Shared mocked API responses for auth and resource endpoints. |
| `src/components/AppShell.jsx` | Fix sidebar/topbar/mobile navigation issues discovered by viewport checks. |
| `src/components/*` resource screens | Fix field labels, error associations, button states, stale alerts, and empty states found by tests. |
| `src/index.css` | Central responsive cleanup for forms, tables, panels, action bars, and dashboard layout. |
| `CHAT_CONTEXT/frontend_skills/tests/12_ux_test_hardening.md` | Record final verification matrix and skipped live-backend checks, if any. |

Preferred test dependency:

- Add `@playwright/test` as a dev dependency if network/dependency installation is allowed.
- If a dependency cannot be added, use the available MCP Playwright/manual browser tooling and record
  the limitation clearly in the FE12 test note.

Keep production changes scoped to issues found by the hardening pass. Do not redesign the whole app
or introduce a component library.

## State And API Plan

No new backend API endpoints are expected.

Mocked API coverage should include:

- auth login/current-user/refresh/logout
- resource success responses for every route group
- representative `400`, `401`, `403`, `404`, and `409` error responses using the shared backend error
  envelope
- background refresh failure after successful mutation for subscription/session detail flows

Live API coverage, when seeded backend credentials exist:

- login/current employee/logout
- at least one success and one expected error path for members, courses, branches, subscriptions,
  attendance, sessions, and employees
- role checks for admin-only employees and manager-only settings

Frontend state checks:

- loading states do not block navigation permanently
- success messages survive background refreshes when the mutation already succeeded
- field-level validation errors are reachable and announced
- disabled buttons have a visible reason when the user can reasonably be confused
- stale-data alerts are visually distinct from hard failure states

## UX States

Hardening acceptance should cover:

- login empty fields and invalid credentials
- protected route while auth restore is loading
- forbidden route for a mocked role
- app not-found route
- each resource command/list page empty and success state
- create forms with required field errors
- detail pages with invalid ID, loading, `404`, success, and mutation failure
- conflict states for duplicate unique values or invalid lifecycle transitions
- blocked reports/payments navigation state
- mobile sidebar open/close behavior
- browser back/forward across command, create, and detail routes

## Responsive And Accessibility Notes

Viewport matrix:

- `320x720`
- `390x844`
- `768x1024`
- `1080x800`
- `1280x800`

Responsive acceptance:

- no page-level horizontal overflow
- topbar, sidebar trigger, staff context, and action bars wrap without covering content
- tables either become stacked records or scroll inside their own wrapper
- form grids stack cleanly with labels and errors visible
- dashboard charts have compact alternatives on small screens
- buttons keep stable dimensions and readable labels
- blocked/forbidden/not-found states fit on mobile

Accessibility acceptance:

- interactive controls are keyboard reachable in logical order
- visible labels or accessible names exist for inputs, selects, and icon-like buttons
- `aria-describedby` connects field errors where practical
- submit/status changes use `aria-live="polite"` for action feedback
- route-level headings are unique and present
- table/list semantics remain understandable after responsive stacking
- focus outlines remain visible on dark UI
- color is not the only signal for status badges, errors, or lifecycle states

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/12_ux_test_hardening.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/12_ux_test_hardening.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/12_ux_test_hardening.md`

Baseline commands:

```sh
cd frontend
npm run build
```

If Playwright is added:

```sh
cd frontend
npm run test:e2e
```

Manual/browser checks:

- Run the viewport matrix above for auth, dashboard, one resource create form, one resource detail
  action flow, and one admin-only route.
- Verify browser back/forward on custom route matcher routes.
- Verify no route redirects to `/login` during mocked-auth browser tests.

Live backend checks when credentials are available:

- Use a seeded admin and at least one non-admin role.
- Smoke live success/error paths listed in FE05-FE10 test notes.
- Record any skipped live check with the missing credential/fixture reason.

Docs updates:

- Update only frontend phase notes unless a backend contract drift is found.
- If a reusable browser test workflow is added, document the command in this plan's test note and
  consider adding it to the frontend worklog summary.

## Risks And Boundaries

- Browser automation may require a new dev dependency and browser install; if network access is not
  available, FE12 may need an MCP/manual pass instead.
- Live API verification still depends on seeded backend credentials and stable fixture data.
- Broad CSS cleanup can accidentally change many screens; prefer targeted fixes backed by viewport
  evidence.
- Do not use FE12 to add new business functionality, dashboard report APIs, payment screens, or a new
  router library unless a test exposes a blocking routing bug.
- Do not replace existing component patterns with a design system during hardening.

## Next Action

Use `$gym-fe-implement` with
`CHAT_CONTEXT/frontend_skills/plans/12_ux_test_hardening.md` after FE11 is either implemented or
explicitly deferred until backend dashboard/report APIs exist.
