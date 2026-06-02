# Review - FE 06 Courses And Branches

## Status

- Status: reviewed
- Feature: Courses and branches settings workspace
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/06_courses_branches.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/06_courses_branches.md`
- Reviewed at: 2026-06-02

## Review summary

- Result: issues found
- Build status: passed with `npm run build`
- Test status: browser route check attempted, but protected route redirected to `/login` because no backend/auth session was available.

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
| medium | `frontend/src/components/settings/BranchForm.jsx:29` and `frontend/src/components/settings/NearbyBranchesPanel.jsx:15` | Blank longitude/latitude fields pass validation because `Number('')` becomes `0`. Branch create/update can send `[0, 0]`, and nearby search can query `(0,0)` instead of showing the required-field error planned for FE06. | Validate the raw string before numeric conversion, e.g. require `values.lng.trim()` and `values.lat.trim()` before range checks. |
| low | `frontend/src/components/settings/BranchForm.jsx:86` | Most branch field errors are rendered without `aria-describedby`, while the plan required connected field errors. This is inconsistent with `CourseForm`. | Add stable error IDs and connect each invalid branch input with `aria-invalid` and `aria-describedby`. |

## Fixes applied during review

- None. Review only.

## Handoff to test

- After fixes, test branch create/update with blank coordinates, nearby search with blank coordinates, invalid coordinate ranges, duplicate branch code `409`, and 320px stacked layout.
