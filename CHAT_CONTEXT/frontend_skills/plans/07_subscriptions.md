# FE Plan - 07 Subscriptions

Status: Planned

Created: 2026-06-02

## Goal

Build the subscription workspace for creating pending subscriptions and managing subscription
lifecycle actions.

FE07 should use live subscription APIs, FE05 member direct lookup, and FE06 course/branch reference
data. Because the backend has no global `GET /api/v1/subscriptions` list/search endpoint, the first
UI should use direct ObjectID lookup plus create flow rather than fake a subscription directory.

## Current Baseline

- Existing routes:
  - `/app/subscriptions`
  - `/app/subscriptions/:id`
- Route registry does not yet include `/app/subscriptions/new`.
- FE05 has member creation/detail and member-scoped subscriptions.
- FE06 is expected to provide course and branch list helpers for selects.
- Backend subscription create validates member/course/branch references and snapshots course pricing.

## Screens And Routes

| Route | Access | Behavior |
|---|---|---|
| `/app/subscriptions` | `admin`, `manager`, `receptionist` | Command center with direct subscription ID lookup, create action, backend-gap note for missing global list, and optional member-scoped lookup link. |
| `/app/subscriptions/new` | `admin`, `manager`, `receptionist` | Create pending subscription form. |
| `/app/subscriptions/:id` | `admin`, `manager`, `receptionist` | Subscription detail, lifecycle actions, refund panel, and attendance shortcut. |

Route config changes:

- Add `/app/subscriptions/new` before `/app/subscriptions/:id`.
- Mark subscriptions routes ready after implementation.
- Keep `/app/subscriptions/:id/attendance` reserved for FE08 and make sure it is registered before
  `/app/subscriptions/:id` when FE08 adds it.

## Component Plan

Add or update:

| Path | Responsibility |
|---|---|
| `src/lib/subscriptionsApi.js` | Create/get/suspend/unsuspend/expire/refund helpers. |
| `src/components/subscriptions/SubscriptionsPage.jsx` | Direct lookup command center and missing-list note. |
| `src/components/subscriptions/SubscriptionCreateView.jsx` | Create pending subscription form. |
| `src/components/subscriptions/SubscriptionDetailView.jsx` | Fetch detail and render lifecycle panels. |
| `src/components/subscriptions/SubscriptionLookupPanel.jsx` | Reusable direct ObjectID lookup. |
| `src/components/subscriptions/SubscriptionSummaryPanel.jsx` | Status, pricing, sessions, dates, member/course/branch IDs. |
| `src/components/subscriptions/SubscriptionLifecyclePanel.jsx` | Suspend/unsuspend/expire actions. |
| `src/components/subscriptions/RefundPanel.jsx` | Refund action with reason and backend refund response. |
| `src/components/subscriptions/subscriptionFormatters.js` | Money, date, status, ObjectID, discount formatting. |
| `src/App.jsx` | Render FE07 route components. |
| `src/routes/routeConfig.js` | Add `/app/subscriptions/new`; mark routes ready. |
| `src/index.css` | Scoped subscription form/action/status layout. |

## State And API Plan

API helpers:

```js
createSubscription(accessToken, payload)
getSubscription(accessToken, subscriptionId)
suspendSubscription(accessToken, subscriptionId, payload)
unsuspendSubscription(accessToken, subscriptionId)
expireSubscription(accessToken, subscriptionId)
refundSubscription(accessToken, subscriptionId, reason)
```

Create payload:

- `member_id`
- `course_id`
- `home_branch_id`
- `start_date` RFC3339
- `end_date` RFC3339
- `session_per_week`
- `discount_type`
- `discount_value`
- `promo_code`

Suspend payload:

- `start_date` RFC3339
- `end_date` RFC3339
- `frozen_session`
- `reason`

Refund payload:

- `reason`

Reference data:

- Use `GET /api/v1/courses` and `GET /api/v1/branches` from FE06 for selects when available.
- Member remains direct ObjectID lookup because no member list/search endpoint exists.

State:

- Create page:
  - form state, reference-list query state, submit state, API error
  - navigate to detail after successful create
- Detail page:
  - subscription query state
  - lifecycle mutation state
  - refund mutation state
  - refresh detail after lifecycle/refund action

Validation:

- ObjectID validation for member, course, branch, subscription IDs.
- Dates must be valid and `start_date <= end_date`.
- `session_per_week` positive integer.
- `discount_type` optional but should be constrained to backend-supported values when confirmed.
- `discount_value` non-negative.
- Suspend frozen sessions non-negative/positive according to backend behavior; surface backend `400`
  if stricter.

## UX States

- Missing global subscription list gap on command center.
- Direct lookup empty/invalid ID.
- Create reference loading and reference API error.
- Create invalid reference `404`.
- Detail `404` and invalid ID.
- Status-aware action availability:
  - pending: no suspend/refund until activated
  - active: suspend, expire, refund
  - suspended: unsuspend
  - expired/refunded: actions disabled
- `409` lifecycle conflicts.
- Refund success displays refund amount/status from response `refund`.
- Attendance shortcut links to `/app/subscriptions/:id/attendance` once FE08 exists.

## Responsive And Accessibility Notes

- Keep forms dense but readable; two-column desktop, one-column mobile.
- Lifecycle buttons should wrap and stay readable at 320px.
- Destructive actions need confirmation and clear status text.
- Use visible labels and `aria-describedby` for validation.
- Mutation result messages use `role="status"` or `role="alert"`.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/07_subscriptions.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/07_subscriptions.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/07_subscriptions.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Checks:

- Direct lookup invalid, not-found, and success.
- Create pending subscription using known member/course/branch IDs.
- Lifecycle action state matrix.
- Refund response handling.
- Mobile and desktop layout.

## Backend Contract Gaps

- No `GET /api/v1/subscriptions` global list/search endpoint.
- No member list/search endpoint for member selection.
- Confirm exact accepted `discount_type` values from service/tests during implementation.

## Risks And Boundaries

- Do not implement payment activation here; FE05 handles member activation through pending
  subscription ID.
- Do not implement attendance history here; FE08 owns it.
- Do not fake subscription lists.
- Reference selects depend on FE06 or must gracefully degrade to manual ObjectID fields.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/07_subscriptions.md`.
