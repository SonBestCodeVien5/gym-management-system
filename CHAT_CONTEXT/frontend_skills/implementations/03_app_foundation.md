# Implementation - 03 App Routing And API Foundation

## Status

- Status: implemented
- Feature: App Routing And API Foundation
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md`
- Started at: 2026-05-31
- Finished at: 2026-05-31

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/routes/routeConfig.js` - added the app route registry, nav groups, roles, module
  metadata, planned API scope, and placeholder route definitions.
- `frontend/src/routes/matchRoute.js` - added dependency-free route matching with simple `:id`
  params, `/app` redirect handling, workspace not-found detection, and safe malformed param
  handling.
- `frontend/src/lib/permissions.js` - added shared role labels and route role helpers.
- `frontend/src/lib/resourceState.js` - added small resource-state constants/helpers for future
  resource screens.
- `frontend/src/components/AppShell.jsx` - refactored shell navigation to come from route config and
  active URL route instead of local active-item state.
- `frontend/src/App.jsx` - replaced the hardcoded `KNOWN_ROUTES` flow with route matching, `/app`
  to `/app/dashboard` redirect behavior, route-level forbidden/not-found handling, dashboard render,
  and module placeholder render.
- `frontend/src/components/ModulePlaceholder.jsx` - added shared planned-module page.
- `frontend/src/components/PageHeader.jsx` - added compact page header primitive.
- `frontend/src/components/DataPanel.jsx` - added shared content panel primitive.
- `frontend/src/components/StateBlock.jsx` - added reusable loading/empty/error/forbidden/not-found
  state block.
- `frontend/src/index.css` - added page header, data panel, state block, placeholder page, feature
  list, and API list styles.
- `CHAT_CONTEXT/frontend_skills/plans/03_app_foundation.md` - marked implemented and aligned the route
  table with Reports/Payments placeholders.
- `CHAT_CONTEXT/frontend_skills/worklog.md` - updated FE03 implementation handoff.

## Key decisions

- Kept routing dependency-free for this cycle. The matcher supports static routes and simple `:id`
  params, which is enough for placeholder/detail shells.
- Malformed encoded params now fail the current route pattern and fall through to workspace
  not-found instead of throwing during render.
- `/app` now redirects to `/app/dashboard`; `/app/dashboard` renders the existing FE02 dashboard.
- Business module routes render placeholders only. No member/course/branch/subscription/session/etc.
  API calls were added.
- Route-level role gating is implemented as frontend UX only; backend remains the real security
  boundary.
- Existing FE02.1 responsive behavior was preserved. FE12 remains the place for broader responsive
  redesign.

## Commands run

```bash
cd frontend
npm run build
```

Result: pass. Vite built 36 modules and emitted production assets.

```bash
git diff --check
```

Result: pass.

```bash
cd frontend
node --input-type=module -e "<route matcher smoke check>"
```

Result: pass for `/`, `/login`, `/app`, `/app/dashboard`, `/app/members/:id`, malformed encoded
detail URL, unknown `/app/*`, and non-app redirect behavior.

```bash
cd frontend
node --input-type=module -e "<permission smoke check>"
```

Result: pass for admin/manager/trainer/receptionist route access cases.

## Known limitations

- No browser route/manual auth flow was run in this implementation phase.
- No business API calls are implemented yet; placeholder pages list planned APIs only.
- The dependency-free matcher is intentionally small. Revisit React Router if nested workflows become
  hard to maintain.
- Dashboard responsive UX remains deferred to FE12.

## Handoff to review

- Review `App.jsx`, `routes/routeConfig.js`, and `routes/matchRoute.js` for route correctness.
- Confirm role-gated forbidden states are shown for inaccessible routes.
- Confirm `/app` redirects to `/app/dashboard` and unknown `/app/*` routes render not-found.
- Confirm placeholders do not imply live business data.
- Use `$gym-fe-review` with this implementation note.
