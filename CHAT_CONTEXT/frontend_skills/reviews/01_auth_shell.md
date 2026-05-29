# Review - 01 Auth Shell

## Status

- Status: reviewed
- Feature: Staff auth shell
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/01_auth_shell.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/01_auth_shell.md`
- Reviewed at: 2026-05-29

## Review summary

- Result: changes requested
- Build status: pass, `cd frontend && npm run build`
- Test status: manual backend/browser auth flow not run in this review phase

## Checklist

- [x] UI matches intended design/style.
- [x] Routes and components are scoped cleanly.
- [x] API base URL/env handling is correct.
- [x] Auth/session state is not hardcoded or leaked.
- [x] Loading, empty, and error states are handled where relevant.
- [x] Responsive layout has mobile/desktop CSS breakpoints.
- [x] Accessibility basics are covered.
- [x] Docs/context are aligned.

## Issues found

| Severity | File | Issue | Fix |
|---|---|---|---|
| medium | `frontend/src/components/LoginView.jsx:15` | Login blocks any password shorter than 8 characters before calling the API. The backend login contract only requires a non-empty password, and `internal/service/auth_service.go:147` validates only `password == ""` before bcrypt comparison. Backend auth tests also log in with `"secret"` at `internal/service/auth_service_test.go:109`. This can reject valid existing/bootstrap credentials in the browser while the API would accept them. | Remove the login-time min-length check and keep only required-field validation. Keep the 8-character rule for employee create/reset screens when those screens are implemented. |

## Fixes applied during review

- None. Review note only; production code was not changed in this phase.

## Remaining risks

- Manual login, restore, refresh, logout, wrong-password, backend-down, and mobile viewport checks still
  need to run against a live backend.
- Token persistence remains localStorage-based by design for this MVP.
- Business modules beyond Dashboard are placeholders.

## Handoff to test

- Fix the login password validation issue before or during the next implementation pass.
- Then use `$gym-fe-test` with this review note.
- Required checks:
  - `cd frontend && npm run build`
  - Valid admin login through `POST /api/v1/auth/login`
  - Reload `/app` and confirm restore through `GET /api/v1/auth/me`
  - Expired/invalid access token with valid refresh token retries restore once
  - Logout clears local tokens and returns to `/login`
  - Wrong password shows sanitized backend error
  - Backend unavailable shows network error
  - Mobile around 320px and desktop around 1280px
