# FE Plan - 01 Auth Shell

Status: Planned

Created: 2026-05-29

## Goal

Build the first operational frontend slice: staff login, token-backed session restore, logout, and a
protected dashboard shell. This gives later member, subscription, attendance, session, and employee
screens a shared API client, auth state, route guard, and error handling baseline.

This cycle should replace the current static marketing-first `App.jsx` experience with a staff portal
entry point. The existing Iron Forge visual language can be reused, but the first screen should be the
usable staff login/workspace, not a landing page.

## Current Baseline

- Frontend stack is React 18 + Vite with no router, state, form, test, or UI dependencies.
- Current source files are only `frontend/src/main.jsx`, `frontend/src/App.jsx`, and
  `frontend/src/index.css`.
- Current `App.jsx` is a static Iron Forge landing page.
- Backend browser readiness is done:
  - CORS supports configured Vite origins through `CORS_ALLOWED_ORIGINS`.
  - Auth uses bearer access tokens in the `Authorization` header.
  - `GET /api/v1/auth/me` can restore the current staff identity from an access token.
- Backend error shape is stable:
  `{"error":{"code":"...","message":"...","details":{}}}`.

## Screens And Routes

Use simple browser routes or hash routes without adding a router dependency in this first pass. If the
implementation chooses React Router later, keep the same public paths.

| Route | Access | Purpose |
|---|---|---|
| `/login` | Public, redirect authenticated users to `/app` | Staff login form |
| `/app` | Protected | Dashboard shell with current employee context |
| `/` | Redirect | Send unauthenticated users to `/login`, authenticated users to `/app` |

Initial protected shell content:

- Top bar with product name, current employee name/email, role badges, and logout action.
- Left navigation for upcoming areas: Dashboard, Members, Subscriptions, Attendance, Sessions,
  Employees. Links can be disabled or point to placeholder panels in this cycle.
- Dashboard landing area with compact status panels:
  - current staff identity
  - assigned branch count
  - allowed role-derived modules
  - API connection/session status
- No business CRUD forms in this cycle.

## Component Plan

Create a small feature-oriented structure under `frontend/src/`:

| Path | Responsibility |
|---|---|
| `App.jsx` | App composition, route switching, auth provider boundary |
| `lib/api.js` | Fetch wrapper, base URL, auth headers, error-envelope parsing |
| `lib/authStorage.js` | Token persistence helpers |
| `context/AuthContext.jsx` | Auth state, login, restore, refresh, logout actions |
| `components/LoginView.jsx` | Login form and unauthenticated layout |
| `components/AppShell.jsx` | Protected staff portal layout |
| `components/DashboardHome.jsx` | First dashboard content |
| `components/RouteGuard.jsx` | Protected/public redirect decisions |
| `components/StatusMessage.jsx` | Reusable loading/error/empty message block |

Keep components narrow and local. Do not introduce a design-system abstraction until more screens share
the same controls.

## State Plan

Auth state:

```js
{
  status: 'checking' | 'anonymous' | 'authenticated',
  employee: null | {
    id: string,
    employee_id: string,
    email: string,
    full_name: string,
    role: string[],
    branch_id: string[]
  },
  accessToken: string | null,
  refreshToken: string | null,
  error: null | {
    code: string,
    message: string,
    details: object
  }
}
```

Storage:

- Store `access_token` and `refresh_token` in `localStorage` for this MVP because the backend currently
  exposes bearer tokens only and does not use HttpOnly cookies.
- Namespace keys, for example `gym.accessToken` and `gym.refreshToken`.
- Never log tokens or render token strings.
- Clear both tokens on logout, failed refresh, inactive account, missing token, or malformed token.

Session restore flow:

1. On app boot, read stored access token.
2. If no access token, show anonymous state.
3. If access token exists, call `GET /api/v1/auth/me`.
4. If `/auth/me` returns `200`, set authenticated employee state.
5. If `/auth/me` returns `401` and refresh token exists, call `POST /api/v1/auth/refresh`.
6. If refresh succeeds, store replacement tokens and retry `/auth/me` once.
7. If refresh fails, clear storage and show login with a concise session-expired message.

## API Plan

Base URL:

- Add support for `VITE_API_BASE_URL`.
- Default to `http://localhost:8080` when the env var is absent.
- Local backend must set `CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173` for
  browser integration.

Endpoints:

| Action | Method + path | Request | Success |
|---|---|---|---|
| Login | `POST /api/v1/auth/login` | `{ "email": string, "password": string }` | stores `access_token`, `refresh_token`, `employee` |
| Restore current employee | `GET /api/v1/auth/me` | bearer access token | sets employee state |
| Refresh | `POST /api/v1/auth/refresh` | `{ "refresh_token": string }` | stores replacement token pair |
| Logout | `POST /api/v1/auth/logout` | `{ "refresh_token": string }` | clears local state/storage |

Fetch wrapper rules:

- Add `Authorization: Bearer <accessToken>` only when an access token is present.
- Always send `Content-Type: application/json` for JSON body requests.
- Parse success as `{ message, data }`.
- Parse errors from `error.code`, `error.message`, and `error.details`.
- Handle non-JSON/network failures as a local error object, for example
  `{ code: 'NETWORK_ERROR', message: 'Cannot reach API server', details: {} }`.
- Keep automatic token refresh scoped to boot/session restore in this first cycle. Avoid a complex
  global retry interceptor until there are multiple protected API screens.

## UX States

Login view:

- Fields: email, password.
- Client validation:
  - email required
  - password required
  - do not enforce create/reset password policy here; login should only reject an empty password and
    let the backend verify the submitted credential
- Submit loading state disables the form.
- Invalid credentials show the sanitized backend message.
- Network errors show a short API connection message.
- After successful login, navigate to `/app`.

App boot:

- Show a full-page checking state while `/auth/me` or refresh is running.
- If restore fails, route to `/login` and show session-expired state only once.

Protected shell:

- Logout action remains usable even if logout API fails; clear local tokens after the attempt.
- Show `403` as "Tai khoan khong co quyen truy cap khu vuc nay" when later route-level guards are
  added. For this cycle, all authenticated roles can open `/app`.
- Placeholder navigation items should communicate unavailable modules with disabled styling, not fake
  successful screens.

## Responsive And Accessibility Notes

- Login layout must work at 320px width without horizontal scroll.
- App shell should collapse navigation into a top/bottom compact menu on mobile instead of forcing a
  wide sidebar.
- Keep focus visible for inputs and buttons.
- Use semantic form labels, not placeholder-only labels.
- Use `aria-live="polite"` for login/session error messages.
- Button text must not overflow at Vietnamese copy lengths.
- Avoid nested cards. Use one shell surface plus panels/rows.
- Keep palette close to current Iron Forge dark/orange style, but reduce marketing-scale typography
  inside the app shell.

## Docs And Test Plan

Implementation note should be created at:

- `CHAT_CONTEXT/frontend_skills/implementations/01_auth_shell.md`

Review note should check:

- token storage and clearing behavior
- refresh retry limit
- backend error-envelope parsing
- route guard behavior
- mobile shell layout

Test note should be created at:

- `CHAT_CONTEXT/frontend_skills/tests/01_auth_shell.md`

Verification commands:

```sh
cd frontend
npm run build
```

Manual browser/API verification:

- Start backend with bootstrap admin, JWT secrets, MongoDB, and CORS allowed for Vite.
- Start Vite dev server.
- Login with valid admin credentials.
- Refresh the page on `/app`; current employee should restore through `/auth/me`.
- Invalidate access token path by clearing only access token or using an expired token if available;
  refresh should rotate tokens and restore once.
- Logout; tokens should be cleared and `/app` should redirect to `/login`.
- Try wrong password; show backend `401` message without exposing raw internals.
- Stop backend; login should show network error.
- Check mobile viewport around 320px and desktop viewport around 1280px.

## Risks And Boundaries

- `localStorage` token persistence is acceptable for this MVP because the backend currently has no
  cookie auth flow. It is not the final security posture for a production staff portal.
- No branch-scoped authorization exists yet. FE should display branch assignment but not assume it can
  enforce branch data access.
- No seed/demo data command exists yet. Manual verification depends on bootstrap admin or existing
  employee data.
- Do not add member/subscription CRUD in this cycle. The goal is shared auth and app shell plumbing.
- Do not change backend endpoints unless frontend integration reveals contract drift.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/01_auth_shell.md`.
