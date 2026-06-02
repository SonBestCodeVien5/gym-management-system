# Test - FE 06 Courses And Branches

## Status

- Status: build and mocked browser route verified; live backend checks skipped with reason
- Feature: Courses and branches settings workspace
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/06_courses_branches.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/06_courses_branches.md`
- Review file: `CHAT_CONTEXT/frontend_skills/reviews/06_courses_branches.md`
- Tested at: 2026-06-02

## Verification summary

- Result: production build passed after review fixes.
- Browser route smoke: passed with Playwright mocked auth/API for `/app/settings/courses`,
  `/app/settings/courses/:id`, `/app/settings/branches`, and `/app/settings/branches/:id`.
- Live backend CRUD smoke: skipped because no seeded backend credentials/session were available in
  this pass.

## Commands run

```bash
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

## Checks covered

- Production bundle compiles after the blank coordinate validation and ARIA wiring fixes.
- FE06 review findings are recorded as fixed in the implementation note.
- Mocked-auth browser route rendering for courses/branches list and detail routes.

## Checks not covered

- Course create/update/delete against a live backend.
- Branch create/update/delete against a live backend, including duplicate branch code `409`.
- Nearby branch search with live GeoJSON data.
- Full desktop/mobile browser interaction verification for stacked records and field error focus.
