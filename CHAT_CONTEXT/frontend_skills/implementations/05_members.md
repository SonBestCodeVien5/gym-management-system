# Implementation - 05 Members

## Status

- Status: implemented
- Feature: FE05 Members
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/05_members.md`
- Started at: 2026-06-01
- Finished at: 2026-06-01

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/App.jsx` - renders the live Members command center, create page, and detail page
  instead of the FE03 placeholder for member routes.
- `frontend/src/routes/routeConfig.js` - marks Members routes ready and adds `/app/members/new`
  before `/app/members/:id`.
- `frontend/src/lib/membersApi.js` - adds member create/get/activate/list-subscriptions API helpers.
- `frontend/src/components/members/MembersPage.jsx` - direct ID lookup command center and backend-gap
  note for missing directory search.
- `frontend/src/components/members/MemberCreateView.jsx` - controlled create-member form, validation,
  submit state, and success navigation.
- `frontend/src/components/members/MemberDetailView.jsx` - member/detail fetch state, subscriptions
  fetch state, invalid/not-found/error states, and activation refresh flow.
- `frontend/src/components/members/MemberLookupPanel.jsx` - reusable ObjectID lookup form.
- `frontend/src/components/members/MemberProfilePanel.jsx` - profile/status/contact/training summary.
- `frontend/src/components/members/MemberSubscriptionsPanel.jsx` - member-scoped subscription
  table/list with pending subscription selection.
- `frontend/src/components/members/OfflinePaymentPanel.jsx` - manual/select subscription activation
  form and success/error feedback.
- `frontend/src/components/members/memberFormatters.js` - local ObjectID validation and member/
  subscription formatting helpers.
- `frontend/src/index.css` - scoped FE05 form, member detail, subscription table, success, select, and
  responsive styling.
- `CHAT_CONTEXT/frontend_skills/worklog.md` - FE05 implementation entry.

## Key decisions

- Kept the directory/search gap visible instead of rendering fake members because the backend has no
  `GET /api/v1/members` or CCID search endpoint.
- Added `/app/members/new` as a static route before `/app/members/:id` so the existing matcher does
  not treat `new` as an ID.
- Kept validation aligned with backend behavior: create requires only `ccid` and `full_name`;
  optional contact/training fields are trimmed but not over-validated.
- Activation supports both a pending subscription selected from member subscriptions and a manual
  subscription ID because FE07 subscription creation is not implemented yet.
- Kept member UI components local under `components/members/` instead of introducing shared
  table/form abstractions.

## Review fixes

2026-06-02:

- Fixed the offline-payment success state so `MemberDetailView` owns a persistent activation notice
  and refreshes member/subscription data in the background instead of unmounting the detail page.
- Fixed manual `subscription_id` validation so non-empty invalid IDs show field-level feedback while
  the confirm button remains disabled.

## Commands run

```bash
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
npm run dev -- --host 127.0.0.1 --port 5174
curl -sS -i http://127.0.0.1:5174/app/members
curl -sS -i http://127.0.0.1:5174/app/members/new
curl -sS -i http://127.0.0.1:5174/app/members/69e100da9359b4be784078df
Playwright mocked browser smoke for `/app/members`, `/app/members/new`, and `/app/members/:id`
npm run build
```

Result:

- Build passed.
- Review-fix build passed.
- Vite on `5173` was blocked by an existing process, so the dev server is running on
  `http://127.0.0.1:5174/`.
- HTTP SPA route smoke returned `200 OK` for `/app/members`, `/app/members/new`, and
  `/app/members/69e100da9359b4be784078df`.
- Mocked browser smoke passed for the Members command center, create page, and detail page with
  mocked `auth/me`, member detail, and member subscriptions responses.

## Known limitations

- Live backend member create/get/subscription/activation flows were not manually smoked in this pass;
  the implementation is wired to the documented endpoints and build-verified only.
- The member command center remains a direct ObjectID lookup MVP until a backend list/search endpoint
  exists.
- FE05 does not create subscriptions, edit members, or show attendance history; those remain later
  cycles.

## Handoff to test

Review findings were fixed on 2026-06-02. Test results are recorded in
`CHAT_CONTEXT/frontend_skills/tests/05_members.md`.
