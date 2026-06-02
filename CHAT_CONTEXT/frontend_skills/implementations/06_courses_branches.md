# Implementation - FE 06 Courses And Branches

## Status

- Status: implemented
- Feature: Courses and branches settings workspace
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/06_courses_branches.md`
- Started at: 2026-06-02
- Finished at: 2026-06-02

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/lib/coursesApi.js` - Course CRUD API helpers.
- `frontend/src/lib/branchesApi.js` - Branch CRUD and nearby search helpers.
- `frontend/src/lib/featureHelpers.js` - Shared ObjectID, date, tag, query, and formatting helpers.
- `frontend/src/components/settings/` - Course/branch list, detail, forms, nearby search, and formatters.
- `frontend/src/App.jsx` - Renders courses and branches settings routes.
- `frontend/src/routes/routeConfig.js` - Marks settings routes ready and adds detail routes.
- `frontend/src/index.css` - Shared resource workspace, form, list, detail, and mobile styles.

## Key decisions

- Kept manager selection as optional manual ObjectID because FE10 employee list is admin-only.
- Used numeric longitude/latitude nearby search without map dependencies.
- Course and branch update forms send the full backend shape.

## Commands run

```bash
npm run build
```

## Known limitations

- No browser/manual backend CRUD smoke was run in this implementation turn.
- Branch deletion/reference rejection is surfaced from the backend without extra client rules.

## Handoff to review

- Review course/branch validation, delete confirmation flows, nearby search states, and mobile stacked records.

## Review fixes - 2026-06-02

- Fixed blank branch/nearby coordinate validation so empty longitude/latitude no longer passes as
  numeric `0`.
- Connected branch and nearby field errors with `aria-invalid` and `aria-describedby`.
- Build passed with `npm run build`.
