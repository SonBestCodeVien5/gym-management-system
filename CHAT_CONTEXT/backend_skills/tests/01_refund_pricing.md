# Test — refund pricing

## Status

- Status: tested
- Feature: refund pricing
- Plan file: `CHAT_CONTEXT/backend_skills/plans/01_refund_pricing.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/01_refund_pricing.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/01_refund_pricing.md`
- Tested at: 2026-05-19 22:26 +07

## Commands

- `gofmt -w internal/service/subscription_service_test.go` — pass
- `go test ./...` — pass
  - `internal/service` tests pass.
  - other packages have no test files.
- `go build ./...` — pass

## Automated tests added

File:
- `internal/service/subscription_service_test.go`

Coverage added:
- Create subscription pricing:
  - no discount / empty discount type
  - percent discount
  - fixed discount
  - invalid percent value
  - invalid fixed value
  - invalid discount type
  - server-calculated `subtotal_amount`, `discount_amount`, `total_amount_paid`
  - server-owned `unit_price`, `total_sessions`, `remaining_sessions`, `status`
- Refund flow:
  - invalid ObjectID
  - missing subscription
  - pending cannot refund
  - expired cannot refund
  - existing audit record blocks duplicate refund
  - no remaining sessions
  - remaining sessions greater than total sessions data conflict
  - atomic update conflict maps to cannot refund
  - active refund success
  - suspended refund success
  - refund amount formula
  - reason trimming
  - refund audit creation

## Manual API tests

### Happy path

- [x] Live Docker MongoDB/API test passed on `gym_mongodb` via `mongodb://admin:password123@127.0.0.1:27017/?authSource=admin&directConnection=true`.
- [x] Temporary API server ran on `PORT=18080`.
- [x] Created member/course/branch/subscription.
- [x] Verified create subscription pricing:
  - `subtotal_amount = 1000000`
  - `discount_amount = 200000`
  - `total_amount_paid = 800000`
  - `remaining_sessions = 10`
- [x] Activated subscription through member activation endpoint.
- [x] Refunded subscription.
- [x] Verified refund response:
  - `refund_amount = 800000`
  - `remaining_sessions = 10`
  - `used_sessions = 0`
  - `status = processed`
- [x] Verified `GET /api/v1/subscriptions/:id` after refund:
  - `status = refunded`
  - `remaining_sessions = 0`

### Invalid input

- [x] Covered by automated service tests:
  - invalid ObjectID
  - invalid discount type
  - invalid discount values

### Not found

- [x] Covered by automated service test:
  - missing subscription returns `ErrSubscriptionNotFound`

### Conflict/business rule

- [x] Covered by automated service tests:
  - pending refund rejected
  - expired refund rejected
  - duplicate refund audit rejected
  - no remaining sessions rejected
  - invalid session data rejected
  - atomic update conflict rejected

## DB state verification

- [x] Live Docker MongoDB/API test used running `gym_mongodb` container.
- [x] Refund audit was inserted and returned with id `6a0c80617f52ee27b625c986`.
- [x] Refunded subscription id `6a0c80617f52ee27b625c985`.
- [x] Repository atomic filter was previously reviewed.
- [x] Service tests verify refund audit creation call and subscription refund update call.

## Issues found

- Initial automated test expectation mutated stub subscription state before assertion; fixed by making test stub not mutate loaded subscription during `RefundSubscription`.
- Existing recorded risks remain:
  - `refunds.subscription_id` unique index not bootstrapped.
  - no Mongo transaction around subscription update + refund audit insert.

## Final result

- Result: pass for automated backend test phase and live Docker MongoDB/API verification.
- Ready to update docs/context: yes.
