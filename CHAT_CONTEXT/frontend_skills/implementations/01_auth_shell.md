# Implementation - 01 Auth Shell

## Status

- Status: implemented
- Feature: Staff auth shell
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/01_auth_shell.md`
- Started at: 2026-05-29
- Finished at: 2026-05-29

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/App.jsx` - replaced static landing page with `/`, `/login`, and `/app` route
  switching plus auth provider boundary.
- `frontend/src/lib/api.js` - added Vite-configured API base URL, JSON fetch wrapper, bearer header
  support, auth endpoints, and backend error-envelope parsing.
- `frontend/src/lib/authStorage.js` - added namespaced localStorage token helpers.
- `frontend/src/context/AuthContext.jsx` - added login, session restore, one refresh retry,
  authenticated employee state, and logout cleanup.
- `frontend/src/components/LoginView.jsx` - added staff login form, validation, loading, API error, and
  session-expired states.
- `frontend/src/components/LoginView.jsx` - follow-up fix: removed login-time password min-length
  validation so FE matches the backend login contract, which requires only a non-empty password.
- `frontend/src/components/RouteGuard.jsx` - added public/protected route redirect behavior.
- `frontend/src/components/AppShell.jsx` - added protected staff workspace shell, role-aware
  navigation placeholders, current employee summary, and logout action.
- `frontend/src/components/DashboardHome.jsx` - added first dashboard content for identity, roles,
  module availability, and API session status.
- `frontend/src/components/StatusMessage.jsx` - added reusable loading/status surface.
- `frontend/src/index.css` - replaced marketing-page styling with responsive staff portal styling while
  keeping the Iron Forge dark/orange visual direction.
- `frontend/vite.config.js` - pinned the local Vite dev URL to `http://127.0.0.1:5173/`.
- `frontend/index.html` - restored the favicon link after adding a new Iron Forge favicon asset.
- `frontend/public/favicon.svg` - added a new `IF` favicon using the Iron Forge dark/orange palette.
- `frontend/src/components/AppShell.jsx` - aligned sidebar wordmark with the reference text-logo style:
  `IRON` text color and `FORGE` accent color.
- `frontend/src/components/LoginView.jsx` - aligned login wordmark with the same reference text-logo
  style.
- `CHAT_CONTEXT/frontend_skills/worklog.md` - updated frontend roadmap status and review handoff.

## Key decisions

- Kept routing dependency-free for this cycle by using small `window.history` route switching.
- Used `VITE_API_BASE_URL` with fallback `http://localhost:8080`.
- Stored `gym.accessToken` and `gym.refreshToken` in localStorage for the bearer-token MVP.
- Automatic refresh is limited to app boot/session restore and retries `/auth/me` once.
- Logout clears local state/storage even if the logout API request fails.
- Business modules beyond Dashboard are visible by role but disabled as placeholders.

## Commands run

```bash
cd frontend
npm run build
```

Result: pass.

Re-run after adding the Vite dev-server config: pass.

Re-run after fixing review finding for login password validation: pass.

Re-run after replacing the mismatched favicon with a new Iron Forge favicon and applying the reference
wordmark style: pass.

## Known limitations

- Manual browser checks against a live backend were not run in this implementation phase.
- Token storage is localStorage-based because the backend does not yet provide HttpOnly cookie auth.
- Placeholder business modules do not contain CRUD screens yet.
- Branch-scope authorization is displayed as out of scope and is not enforced in the frontend.

## Handoff to review

- Review finding in `CHAT_CONTEXT/frontend_skills/reviews/01_auth_shell.md` has been addressed.
- Check route guard behavior for `/`, `/login`, `/app`, and browser back/forward.
- Check token storage clearing on login failure, refresh failure, and logout.
- Check the one-refresh retry path in `AuthContext`.
- Check error parsing for backend error envelopes and network failures.
- Check mobile layout at 320px and desktop layout around 1280px.
