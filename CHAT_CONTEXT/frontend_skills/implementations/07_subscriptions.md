# Implementation - FE 07 Subscriptions

## Status

- Status: implemented
- Feature: Subscription command center, create flow, detail lifecycle, and refund
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/07_subscriptions.md`
- Started at: 2026-06-02
- Finished at: 2026-06-02

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/lib/subscriptionsApi.js` - Create, get, suspend, unsuspend, expire, and refund helpers.
- `frontend/src/components/subscriptions/` - Lookup, create, detail, summary, lifecycle, refund, and formatters.
- `frontend/src/App.jsx` - Renders subscription routes.
- `frontend/src/routes/routeConfig.js` - Adds `/app/subscriptions/new`, marks subscription routes ready, and registers attendance before detail.
- `frontend/src/index.css` - Shared resource workspace styles used by subscription screens.

## Key decisions

- Kept subscription directory as direct ObjectID lookup because the backend has no global list endpoint.
- Course and branch reference options load from FE06 helpers, with manual ObjectID entry still available if reference loading fails.
- Discount type uses backend-supported `none`, `percent`, and `fixed`.

## Commands run

```bash
npm run build
```

## Known limitations

- No live subscription create/lifecycle/refund smoke was run in this implementation turn.
- Member selection remains manual ObjectID because there is no member list/search endpoint.

## Handoff to review

- Review status-aware lifecycle availability, refund response handling, invalid ID states, and create form date conversion.

## Review fixes - 2026-06-02

- Changed subscription detail refresh after lifecycle/refund mutations to background refresh so
  success notices and refund amount feedback remain mounted.
- Connected subscription create validation errors with `aria-invalid` and `aria-describedby`.
- Build passed with `npm run build`.

## Post-push review fix - 2026-06-02

- Render a visible refresh-failure alert in the subscription detail success branch when a background
  refresh after lifecycle/refund succeeds at mutation time but fails to reload detail data.
