# FE Plan - 05 Members

Status: Planned

Created: 2026-06-01

## Goal

Build the first backend-backed resource workspace for member operations in the React/Vite staff
console.

FE05 should replace the Members placeholder with a live workflow for creating member profiles,
opening a member by ID, viewing member details, listing that member's subscriptions, and confirming
offline payment for a pending subscription. It should use the existing FE03 route/API foundation and
FE04 brand shell without changing backend behavior.

Because the backend currently has no `GET /api/v1/members` list/search endpoint, FE05 should not
fake a member directory. The first pass should make the limitation explicit in the UI and support a
direct ObjectID lookup MVP.

## Current Baseline

- Stack: React 18 + Vite 8, custom dependency-free route matcher.
- Auth/session state is owned by `AuthContext`.
- Low-level API behavior is in `frontend/src/lib/api.js`:
  - `apiRequest(path, { method, body, accessToken })`
  - normalized backend error envelope handling
  - auth helpers only
- Current member routes are placeholders from FE03:
  - `/app/members`
  - `/app/members/:id`
- Current shared UI primitives:
  - `PageHeader`
  - `DataPanel`
  - `StateBlock`
  - route-level `RouteGuard`
  - role helpers in `lib/permissions.js`
- FE04 brand assets are already integrated into login, shell, loading, and not-found surfaces.

## Screens And Routes

| Route | Access | FE05 behavior |
|---|---|---|
| `/app/members` | `admin`, `manager`, `receptionist` | Members command center with direct member ID lookup, create action, backend-gap note for directory search, and post-create shortcut state. |
| `/app/members/new` | `admin`, `manager`, `receptionist` | Create member form. On success, navigate to `/app/members/:id`. |
| `/app/members/:id` | `admin`, `manager`, `receptionist` | Live member profile, status/contact/training summary, subscriptions panel, and offline payment confirmation action. |

Route config changes:

- Add a static `/app/members/new` route before `/app/members/:id`.
- Change `/app/members` and `/app/members/:id` from `planned` to `ready` after implementation.
- Keep Members navigation active under the existing `Tong quan` group.
- Keep FE route role gating as UX only; backend remains the security source of truth.

## Component Plan

Add or update these frontend files:

| Path | Responsibility |
|---|---|
| `src/lib/membersApi.js` | Small endpoint helpers wrapping `apiRequest` for create/get/activate/list-subscriptions. |
| `src/components/members/MembersPage.jsx` | `/app/members` command center with lookup form, create navigation, and no-directory backend-gap panel. |
| `src/components/members/MemberCreateView.jsx` | Create member form and submit state. |
| `src/components/members/MemberDetailView.jsx` | Fetch member + subscriptions by route `id`, render profile and action panels. |
| `src/components/members/MemberLookupPanel.jsx` | ObjectID lookup form reusable on command center and optionally detail error states. |
| `src/components/members/MemberProfilePanel.jsx` | Member identity, contact, status, attendance count, created/updated timestamps. |
| `src/components/members/MemberSubscriptionsPanel.jsx` | Member-scoped subscription table/list with status badges and payment action affordance. |
| `src/components/members/OfflinePaymentPanel.jsx` | Select/manual `subscription_id`, confirm action, conflict/error/success feedback. |
| `src/components/members/memberFormatters.js` | Local formatting for dates, boolean statuses, money, and subscription labels if needed. |
| `src/routes/routeConfig.js` | Add `/app/members/new`; mark implemented Members routes ready. |
| `src/App.jsx` | Render Members route components instead of `ModulePlaceholder` for member route keys. |
| `src/index.css` | Add scoped member form/table/action styles using existing tokens and compact panel patterns. |

Keep member components local under `components/members/` until another feature needs the same form or
table abstraction.

## State And API Plan

API helpers:

```js
createMember(accessToken, payload)
getMember(accessToken, memberId)
activateMember(accessToken, memberId, subscriptionId)
listMemberSubscriptions(accessToken, memberId)
```

Backend endpoints:

| Action | Endpoint | Notes |
|---|---|---|
| Create member | `POST /api/v1/members` | Body: `ccid`, `full_name`, optional `email`, `phone`, `gender`, `level`. Response returns the created member. |
| Get member | `GET /api/v1/members/:id` | Requires ObjectID path param. |
| Confirm offline payment | `PATCH /api/v1/members/:id/activate` | Body: `{ "subscription_id": "..." }`; activates the pending subscription and sets `is_registered=true`. |
| List member subscriptions | `GET /api/v1/members/:id/subscriptions` | Returns an array; empty array is valid. |

Local state:

- `MembersPage`:
  - lookup ID input
  - lookup validation message
  - optional last created member shortcut, kept only in component state
- `MemberCreateView`:
  - controlled form values
  - touched/submitted state
  - submit status: `idle | submitting | success | error`
  - normalized API error
- `MemberDetailView`:
  - member query state: `loading | success | error`
  - subscriptions query state: `loading | success | error`
  - refresh counter or callback after activation
- `OfflinePaymentPanel`:
  - selected/manual subscription ID
  - activation status
  - success message
  - normalized API error

Client-side validation:

- Required before create:
  - `ccid`
  - `full_name`
- Optional trimmed fields:
  - `email`
  - `phone`
  - `gender`
  - `level`
- Direct lookup and activation IDs:
  - validate likely Mongo ObjectID shape (`24` hex characters) before calling the API.
- Do not enforce email or phone policies that the backend does not enforce yet.

Request behavior:

- Always pass `employee` access token from `useAuth()`.
- On `401`, rely on existing auth/session handling where available; show a session-expired API error
  if a resource request fails without automatic recovery.
- On successful create, navigate to `/app/members/:id`.
- On successful activation, refetch member and subscriptions so the profile and subscription status
  reflect the backend result.

## UX States

Members command center:

- No live directory/search because backend support is absent.
- Direct ID lookup empty/invalid state.
- Create member primary action.
- Backend-gap message for list/search without blocking implemented workflows.

Create member:

- Initial empty form.
- Required field errors for `ccid` and `full_name`.
- Submit loading state.
- `409 CONFLICT` duplicate CCID shown next to the CCID field or as a form-level conflict.
- Backend `400 INVALID_INPUT` shown as sanitized API message.
- Success should route to detail instead of leaving a stale form.

Member detail:

- Initial loading.
- Invalid ID before fetch.
- `404 NOT_FOUND` member not found.
- Profile loaded but no subscriptions.
- Subscriptions loaded with pending/active/suspended/expired/refunded statuses.
- Subscription list API error independent from member profile error when possible.

Offline payment confirmation:

- Pending subscription can be selected from the member subscriptions list when present.
- Manual `subscription_id` input remains available because FE07 subscription creation is not built
  yet and the source subscription may be known externally.
- Disable submit when ID is invalid or request is submitting.
- `404` subscription/member not found and `409` invalid lifecycle/member mismatch shown clearly.
- Success message confirms payment activation and refreshes data.

## Responsive And Accessibility Notes

- Keep the staff-console layout dense and operational; no hero or marketing sections.
- At 320px:
  - form fields stack vertically
  - action bars wrap without horizontal overflow
  - subscription rows can become compact stacked records instead of forcing a full table
- On desktop:
  - use two-column panels only where both columns remain readable
  - keep profile and payment actions near each other but avoid nested cards
- Forms need visible labels, field-level errors, and submit buttons with loading text.
- Status and error feedback after submit should use existing state/message patterns and `aria-live`
  where the text changes after user action.
- The lookup form must be keyboard-submittable.
- Do not make disabled actions the only explanation; pair disabled payment actions with concise
  state text.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/05_members.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/05_members.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/05_members.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Static/manual checks:

- `/app/members` renders the live member command center for allowed roles.
- `/app/members/new` renders the create form.
- `/app/members/:id` renders invalid-ID, loading, not-found, and success states.
- Browser back/forward works after create navigation and detail navigation.
- Forbidden roles still show the route-level access-denied state.
- Mobile 320px and desktop widths have no page-level horizontal overflow.

Backend checks when local API and credentials are available:

- Login as `admin`, `manager`, or `receptionist`.
- Create member with unique `ccid`; verify `201` result appears on detail route.
- Create duplicate `ccid`; verify `409` is shown without clearing the form.
- Open member by valid ID; verify profile fields render.
- Open unknown valid ObjectID; verify `404` state.
- List member subscriptions; verify empty array state and non-empty status rows when fixtures exist.
- Confirm offline payment with a pending subscription for that member; verify success and data refresh.
- Confirm with wrong-member or already-active subscription; verify `409` state.

## Backend Contract Gaps

- No `GET /api/v1/members` list/search endpoint exists.
- No `GET /api/v1/members?ccid=...` lookup exists.
- No endpoint returns member + subscriptions in one request.
- FE05 should not change backend docs unless implementation discovers contract drift. A later backend
  cycle can add list/search if the product needs a real member directory.

## Risks And Boundaries

- Do not implement FE07 subscription creation in FE05. The payment panel may use existing member
  subscriptions or manual subscription ID input only.
- Do not implement attendance history in FE05; leave that for FE08.
- Do not add broad table/form abstractions before the Courses/Branches and Subscriptions cycles prove
  the shared shape.
- Do not add a router dependency in this cycle unless the current route matcher blocks `/new` and
  `/:id` behavior.
- Backend validation for member optional fields is minimal; keep frontend validation aligned instead
  of inventing stricter rules.
- FE-created member data may expose that there is no current edit-member endpoint. Treat editing as
  out of scope unless a backend cycle adds it.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/05_members.md`.
