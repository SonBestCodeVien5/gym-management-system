# Implementation — refund pricing

## Status

- Status: implemented
- Feature: refund pricing
- Plan file: `CHAT_CONTEXT/backend_skills/plans/01_refund_pricing.md`
- Started at: 2026-05-19 21:40 +07
- Finished at: 2026-05-19 21:47 +07

## Scope implemented

- [x] Model changes
- [x] Repository changes
- [x] Service changes
- [x] Handler changes
- [x] Route changes
- [x] Docs/API sample changes

## Files changed

- `internal/models/subscription.go`
- `internal/models/refund.go`
- `internal/repository/subscription_repo.go`
- `internal/repository/refund_repo.go`
- `internal/service/subscription_service.go`
- `internal/handlers/subscription_handler.go`
- `cmd/server/main.go`
- `docs/api_contract.md`
- `api_test.http`
- `CHAT_CONTEXT/backend_skills/implementations/01_refund_pricing.md`

## Key decisions

- Server calculates subscription pricing from course snapshot.
- Discount supports `none`, empty, `percent`, and `fixed`.
- Refund only allows `active` subscriptions with remaining sessions.
- Refund rejects `suspended` subscriptions because bảo lưu thẻ không áp dụng hoàn tiền.
- Refund amount uses integer formula: `total_amount_paid * remaining_sessions / total_sessions`.
- Subscription refund uses atomic update from active with remaining sessions to `refunded` and `remaining_sessions = 0`.
- Refund audit record saved in `refunds` collection.

## Implementation notes

- `Subscription.Total_Amount_Paid` refactored to `TotalAmountPaid`.
- Added pricing fields: `subtotal_amount`, `discount_type`, `discount_value`, `discount_amount`, `promo_code`, `total_amount_paid`.
- Added `Refund` model with processed status.
- Added `RefundRepository` with create and lookup by subscription ID.
- Added `POST /api/v1/subscriptions/:id/refund`.
- Updated refund rule to match latest plan: `suspended` cannot refund; only `active` can refund.
- No transaction wrapper exists; refund flow can still become partial if refund insert fails after subscription update.

## Commands run

- `gofmt -w internal/models/subscription.go internal/models/refund.go internal/repository/subscription_repo.go internal/repository/refund_repo.go internal/service/subscription_service.go internal/handlers/subscription_handler.go cmd/server/main.go`
- `go build ./...`
- `go test ./...`

## Known limitations

- Unique index for `refunds.subscription_id` not bootstrapped.
- No Mongo transaction wrapper exists around subscription update + refund insert.
- Duplicate refund prevention does pre-check by subscription ID plus atomic subscription status transition.

## Handoff to review

- Review service refund partial failure behavior.
- Review whether `refunds.subscription_id` unique index should be created during app startup or migration.
- Review handler body binding for refund if empty body should be accepted.