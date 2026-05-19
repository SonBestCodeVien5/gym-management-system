# Code Review — refund pricing

## Status

- Status: reviewed
- Feature: refund pricing
- Plan file: `CHAT_CONTEXT/backend_skills/plans/01_refund_pricing.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/01_refund_pricing.md`
- Reviewed at: 2026-05-19 21:56 +07

## Review summary

- Result: pass with recorded risks
- Build status: pass (`go build ./...`)
- Test status: pass (`go test ./...`, no test files)

## Checklist

- [x] Code compiles.
- [x] Handler only handles HTTP parse/response.
- [x] Service owns business rules.
- [x] Repository only handles DB.
- [x] Model tags match API/DB contract.
- [x] Errors map to correct HTTP status.
- [x] Atomic updates used where needed.
- [x] Routes have correct order.
- [x] Docs/API samples match behavior.

## Passed

- Subscription pricing is calculated server-side from course snapshot.
- Client cannot set `subtotal_amount`, `discount_amount`, `total_amount_paid`, `unit_price`, `total_sessions`, or `remaining_sessions` through request DTO.
- Discount validation matches plan:
  - empty/`none` => zero discount
  - `percent` requires `0..100`
  - `fixed` requires `0..subtotal`
- `Subscription.Total_Amount_Paid` was refactored to idiomatic `TotalAmountPaid`.
- Refund service blocks `pending`, `suspended`, `expired`, and `refunded` via status validation.
- Refund service blocks no remaining sessions and invalid session data.
- Refund amount uses planned formula:
  - `total_amount_paid * remaining_sessions / total_sessions`
- Repository refund update is atomic for status/session transition:
  - status must be `active`
  - `remaining_sessions > 0`
  - update sets `status = refunded`, `remaining_sessions = 0`
- Refund audit record includes subscription/member/session/money/reason/status/timestamps.
- Route `POST /api/v1/subscriptions/:id/refund` is registered before `GET /api/v1/subscriptions/:id`, so no route conflict.
- `docs/api_contract.md` and `api_test.http` include refund/pricing behavior.
- Build and test commands pass.

## Issues found

- [medium] `refunds.subscription_id` unique index is not bootstrapped. Atomic subscription status update prevents normal double refund, but DB-level uniqueness is still missing from plan requirement and remains important for data integrity.
- [medium] Refund flow has partial failure risk: subscription can be updated to `refunded` before refund audit insert fails. No transaction wrapper exists, so this is accepted as recorded limitation for now.
- [low] `RefundSubscription` maps atomic update `ErrNotFound` to `ErrSubscriptionCannotRefund`. If subscription is deleted between read and update, response becomes `409` instead of `404`. Race is unlikely, but status is not perfectly precise.
- [low] Handler comment says refund body is optional, but `ShouldBindJSON` rejects empty body. Current API plan shows a JSON body with `reason`, so behavior is acceptable; comment should be corrected later if empty-body refund is desired.
- [low] API contract records endpoint and behavior summary, but does not include detailed request/response JSON schema for pricing/refund. Existing project contract is high-level, so this is not blocking.

## Fixes applied during review

- Review after latest status-rule update confirmed refund eligibility now matches plan:
  - service allows only `active`
  - repository atomic filter allows only `active`
  - service test expects `suspended` to return `ErrSubscriptionCannotRefund`

## Remaining risks

| Risk | Severity | Impact | Fix timing |
|---|---:|---|---|
| Missing unique index on `refunds.subscription_id` | medium | DB does not enforce one refund per subscription. Current service pre-check + atomic subscription update prevents normal double refund, but data integrity still depends on app logic. Manual DB writes or future code paths could create duplicate refund records. | Soon, before production or before more refund flows. Good fit for `06_indexes_data_integrity`. |
| No Mongo transaction around subscription update + refund insert | medium | If subscription update succeeds but refund insert fails, subscription becomes `refunded` with no audit record. Money/audit history becomes incomplete and support/accounting cannot prove refund details. | Later for MVP if accepted limitation, but soon before real payment/accounting integration. Needs Mongo session/transaction design. |
| Atomic update `ErrNotFound` maps to `ErrSubscriptionCannotRefund` | low | Rare delete/race case can return HTTP `409` instead of `404`. Does not corrupt data. | Later. Clean up during validation/error consistency cycle. |
| Refund handler comment says body optional but `ShouldBindJSON` rejects empty body | low | Developer/API confusion. Client must send `{}` or `{ "reason": "..." }`; empty body gets `400`. | Later, unless API should allow empty body. Then fix soon by tolerating `io.EOF`. |
| API contract lacks detailed JSON schema/examples for pricing/refund | low | FE may need to inspect `api_test.http` or code for exact request/response fields. | Later or during API contract cleanup. Not blocking backend correctness. |
| No automated refund/pricing tests | medium | Regressions in discount math, refund amount, and double refund behavior may slip through manual testing. | Soon. Add in `/backend-test 01_refund_pricing` or integration test cycle. |

### Fix priority

1. Soon: add unique index for `refunds.subscription_id`.
2. Soon: add automated tests for pricing/refund edge cases.
3. Later unless production payments start: add Mongo transaction for update + audit insert.
4. Later: polish error mapping race edge case.
5. Later: align refund empty-body behavior/comment.
6. Later: expand API contract examples.

## Handoff to test

- Test create subscription with:
  - no discount
  - `discount_type = percent`, valid and invalid values
  - `discount_type = fixed`, valid and invalid values
- Test refund rejects:
  - invalid ObjectID
  - missing subscription
  - `pending`
  - `suspended`
  - `expired`
  - already `refunded`
  - `remaining_sessions = 0`
- Test refund success after activation:
  - verify response `refund.refund_amount`
  - verify subscription status becomes `refunded`
  - verify `remaining_sessions = 0`
  - verify refund audit document exists
- Test double refund returns conflict.
- Record DB verification in `CHAT_CONTEXT/backend_skills/tests/01_refund_pricing.md`.