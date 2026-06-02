# Review - FE 10 Employees

## Status

- Status: reviewed
- Feature: Admin-only employee management workspace
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/10_employees.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/10_employees.md`
- Reviewed at: 2026-06-02

## Review summary

- Result: minor issues found
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
| low | `frontend/src/components/employees/EmployeeForm.jsx:135` | Branch assignment uses a `textarea` with a `list` attribute, but datalist suggestions are not supported for `textarea`. Staff will not get the branch options that FE10 planned to provide from FE06 branch data. | Use an `input` for single branch entry plus add/remove chips, or keep a textarea but remove the non-functional datalist and show selectable branch IDs separately. |
| low | `frontend/src/components/employees/EmployeeForm.jsx:87` | Update submits the full mutable employee payload every time and does not enforce "at least one changed mutable field" from the plan. This is unlikely to break backend behavior but does not match the planned update validation. | Track original values in detail view and submit only changed fields, or add an unchanged-form guard before calling `PATCH`. |

## Fixes applied during review

- None. Review only.

## Handoff to test

- After fixes, test admin-only access, branch assignment UX, unchanged update behavior, deactivation confirmation/conflict display, and password reset field clearing.
