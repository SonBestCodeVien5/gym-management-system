# Review - FE 07 Subscriptions

## Status

- Status: reviewed
- Feature: Subscription command center, create flow, detail lifecycle, and refund
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/07_subscriptions.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/07_subscriptions.md`
- Reviewed at: 2026-06-02

## Review summary

- Result: issues found
- Build status: passed with `npm run build`
- Test status: browser route check attempted in this review batch, but protected routes redirected to `/login` because no backend/auth session was available.

## Checklist

- [x] UI matches intended design/style.
- [x] Routes and components are scoped cleanly.
- [x] API base URL/env handling is correct.
- [x] Auth/session state is not hardcoded or leaked.
- [x] Loading, empty, and error states are handled where relevant.
- [ ] Responsive layout works on mobile and desktop.
- [ ] Accessibility basics are covered.
- [x] Docs/context are aligned.

## Issues found

| Severity | File | Issue | Fix |
|---|---|---|---|
| medium | `frontend/src/components/subscriptions/SubscriptionDetailView.jsx:27`, `frontend/src/components/subscriptions/SubscriptionLifecyclePanel.jsx:55`, `frontend/src/components/subscriptions/RefundPanel.jsx:22` | Lifecycle and refund success notices are lost or only flash briefly. Child panels set success state, then call `onChanged`; parent `loadSubscription` immediately sets detail state to `loading`, replacing the whole detail view and unmounting the child panel state. This breaks the planned success feedback, including refund amount display. | Keep refresh as a background load that does not replace the detail route, or lift mutation notices/refund result into `SubscriptionDetailView` so they survive the refresh. |
| low | `frontend/src/components/subscriptions/SubscriptionCreateView.jsx:145` | Field-level validation errors are not connected with `aria-describedby`, despite the plan requiring connected field errors. | Add stable error IDs and `aria-invalid`/`aria-describedby` for invalid create fields. |

## Fixes applied during review

- None. Review only.

## Handoff to test

- After fixes, test lifecycle success messages, refund amount display after refresh, invalid direct lookup, reference-list fallback, and create date conversion.

## Post-push review - 2026-06-02

| Severity | File | Issue | Fix |
|---|---|---|---|
| low | `frontend/src/components/subscriptions/SubscriptionDetailView.jsx:35` | Background detail refresh errors after lifecycle/refund mutations are stored in `subscriptionState.error`, but the success branch never renders that error. If the mutation succeeds and the follow-up detail fetch fails, staff keep the success notice but see stale subscription data with no visible refresh failure. | Render a small alert in the success branch when `subscriptionState.error` is present, or make background refresh return/throw a refresh error that the child panel can display separately. |
