# Review - 03 App Routing And API Foundation

## Status

- Status: reviewed
- Feature: 03 App Routing And API Foundation
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/03_app_foundation.md`
- Reviewed at: 2026-05-31

## Review summary

- Result: changes requested
- Build status: pass (`cd frontend && npm run build`)
- Test status: pass for route matcher and permission smoke checks; browser/manual auth flow not run

## Checklist

- [x] UI matches intended design/style.
- [x] Routes and components are scoped cleanly.
- [x] API base URL/env handling is correct.
- [x] Auth/session state is not hardcoded or leaked.
- [x] Loading, empty, and error states are handled where relevant.
- [x] Responsive layout works on mobile and desktop by code inspection; no browser screenshot run.
- [x] Accessibility basics are covered.
- [x] Docs/context are aligned.

## Issues found

| Severity | File | Issue | Fix |
|---|---|---|---|
| medium | `frontend/src/routes/matchRoute.js:27` | `decodeURIComponent(pathPart)` can throw `URIError` for malformed encoded route params, e.g. `/app/members/%E0%A4%A`. Because `matchRoute()` runs during render in `App.jsx`, this crashes the app instead of showing the planned not-found/invalid route state. | Decode params through a safe helper that catches `URIError` and treats the route as unmatched, or stores the raw segment and lets the target page validate the ID. Add a smoke check for malformed detail URLs. |

## Fixes applied during review

- None. This was a review-only pass.

## Verification run

```bash
cd frontend
npm run build
```

Result: pass. Vite built 36 modules and emitted production assets.

```bash
cd frontend
node --input-type=module -e "<route matcher smoke check>"
```

Result: pass for `/`, `/login`, `/app`, `/app/dashboard`, `/app/members/:id`, unknown `/app/*`, and non-app redirect behavior.

```bash
cd frontend
node --input-type=module -e "<permission smoke check>"
```

Result: pass for admin/manager/trainer/receptionist route access cases.

```bash
cd frontend
node --input-type=module -e "<malformed encoded param check>"
```

Result: fails with `URIError: URI malformed`.

## Remaining risks

- No browser route/manual auth flow was run in this review.
- No API integration was exercised because FE03 intentionally adds no new business API calls.
- Dashboard responsive behavior remains the FE02.1 temporary state, with broader polish deferred to FE12.

## Handoff to test

- Fix malformed route-param decoding before FE03 test signoff.
- After the fix, rerun build, route matcher smoke checks, and manual browser checks for `/app`,
  `/app/dashboard`, role-forbidden placeholders, unknown `/app/*`, browser back/forward, and logout.
