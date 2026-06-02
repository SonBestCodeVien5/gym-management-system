# Review - 05 Members

## Status

- Status: reviewed
- Feature: FE05 Members
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/05_members.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/05_members.md`
- Reviewed at: 2026-06-02

## Review summary

- Result: needs fixes
- Build status: `npm run build` passed
- Test status: Browser review pass was not run because MCP Playwright is not available in this turn;
  implementation note records prior mocked browser smoke.

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
| medium | `frontend/src/components/members/OfflinePaymentPanel.jsx:46`, `frontend/src/components/members/MemberDetailView.jsx:30` | Successful offline-payment feedback is effectively lost. `OfflinePaymentPanel` sets success state/message, then immediately awaits `onActivated()`. That calls `loadMember()`, which sets member state to `loading`; `MemberDetailView` then returns the full loading branch and unmounts the payment panel, resetting its local success message before users can reliably see it. This misses the FE05 plan's success-feedback requirement. | Keep the detail page mounted during post-activation refresh, or move activation success notice into `MemberDetailView` state so it survives refresh. Use a separate background refresh state for subscriptions/member data after activation. |
| low | `frontend/src/components/members/OfflinePaymentPanel.jsx:35`, `frontend/src/components/members/OfflinePaymentPanel.jsx:119` | Invalid manual `subscription_id` has no reachable field-level explanation while the button is disabled. The only invalid-ID message is set inside submit handling, but submit cannot run when `!isObjectId(subscriptionId)` disables the button. This leaves keyboard/screen-reader users with a disabled action and no immediate reason beyond placeholder text. | Show validation text when a non-empty subscription ID is invalid, or keep submit enabled and show the existing validation error on submit. |

## Fixes applied during review

- None. Review only.

## Handoff to test

- Fix the two review findings, then run `npm run build`.
- Run a browser pass for:
  - `/app/members`
  - `/app/members/new`
  - `/app/members/:id`
  - successful activation with a pending subscription
  - invalid manual subscription ID feedback
  - one narrow viewport around 320-375px
