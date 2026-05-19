# Cycle 01 — Refund Flow & Pricing Rules

## Status

- Status: planned
- Priority: highest
- Depends on: existing members/courses/branches/subscriptions/attendance
- Main endpoint:
  - `POST /api/v1/subscriptions/:id/refund`
- Related endpoint:
  - `POST /api/v1/subscriptions`

## Goal

Hoàn thiện pricing khi tạo subscription và thêm refund flow.

Scope:
- Subscription creation hỗ trợ discount server-side.
- Refund endpoint tính tiền hoàn dựa trên buổi còn lại.
- Refund tạo audit record.
- Subscription sau refund chuyển `status = refunded`, `remaining_sessions = 0`.

## Why now

Theo `CHAT_CONTEXT/README.md`, recommended next step là:

> Implement refund flow and pricing rules next.

Lý do:
- Subscription lifecycle hiện thiếu trạng thái kết thúc bằng refund.
- Pricing hiện chỉ copy base price và session count từ course, chưa có discount.
- Refund là rule tiền quan trọng, nên làm trước nearby/auth.

## Current code facts

- `Subscription` hiện có:
  - `Status`
  - `PaymentDate`
  - `Total_Amount_Paid`
  - `UnitPrice`
  - `TotalSessions`
  - `RemainingSessions`
- `CreateSubscription` hiện:
  - validate member/course/branch
  - set `status = pending`
  - snapshot `AllowedTags`
  - set `UnitPrice = course.BasePrice`
  - set `TotalSessions = course.SessionCount`
  - set `Total_Amount_Paid = course.BasePrice * course.SessionCount`
  - set `RemainingSessions = course.SessionCount`
- `SubscriptionRepository` đã có:
  - `UpdateStatus`
  - `UpdateRemainingSessions`
  - `UpdateRemainingSessionsAndStatus`

## API plan

### Create subscription with pricing rules

Endpoint:

```http
POST /api/v1/subscriptions
```

Add optional request fields:

```json
{
  "discount_type": "percent",
  "discount_value": 10,
  "promo_code": "SUMMER10"
}
```

Allowed `discount_type`:
- empty
- `none`
- `percent`
- `fixed`

Response should include:
- `subtotal_amount`
- `discount_type`
- `discount_value`
- `discount_amount`
- `promo_code`
- `total_amount_paid`

### Refund subscription

Endpoint:

```http
POST /api/v1/subscriptions/:id/refund
```

Request:

```json
{
  "reason": "member requested cancellation"
}
```

Success response:

```json
{
  "message": "subscription refunded successfully",
  "refund": {
    "id": "ObjectID",
    "subscription_id": "ObjectID",
    "member_id": "ObjectID",
    "used_sessions": 4,
    "remaining_sessions": 8,
    "refund_amount": 800000,
    "reason": "member requested cancellation",
    "status": "processed",
    "created_at": "2026-05-19T13:00:00Z",
    "processed_at": "2026-05-19T13:00:00Z"
  }
}
```

## Business rules

### Pricing rules

- Server calculates money from course snapshot.
- Client must not control total amount.
- `subtotal_amount = course.base_price * course.session_count`
- If `discount_type` empty or `none`:
  - `discount_amount = 0`
  - `discount_value = 0`
- If `discount_type = percent`:
  - `discount_value >= 0`
  - `discount_value <= 100`
  - `discount_amount = subtotal_amount * discount_value / 100`
- If `discount_type = fixed`:
  - `discount_value >= 0`
  - `discount_value <= subtotal_amount`
  - `discount_amount = discount_value`
- `total_amount_paid = subtotal_amount - discount_amount`
- `unit_price = course.base_price`
- `total_sessions = course.session_count`
- `remaining_sessions = course.session_count`
- `status = pending`

### Refund rules

- Only allow refund status:
  - `active`
  - `suspended`
- Reject refund status:
  - `pending`
  - `expired`
  - `refunded`
- `total_sessions > 0`
- `remaining_sessions > 0`
- `used_sessions = total_sessions - remaining_sessions`
- If `used_sessions < 0`, reject as data conflict.
- `refund_amount = total_amount_paid * remaining_sessions / total_sessions`
- After refund:
  - subscription `status = refunded`
  - subscription `remaining_sessions = 0`
  - refund record inserted
- Double refund prevention:
  - atomic update status active/suspended → refunded
  - unique index `refunds.subscription_id`

## Data model plan

### Update `Subscription`

Refactor current field:
- `Total_Amount_Paid` → `TotalAmountPaid`

Add fields:

```go
SubtotalAmount  int64  `bson:"subtotal_amount" json:"subtotal_amount"`
DiscountType    string `bson:"discount_type" json:"discount_type"`
DiscountValue   int64  `bson:"discount_value" json:"discount_value"`
DiscountAmount  int64  `bson:"discount_amount" json:"discount_amount"`
PromoCode       string `bson:"promo_code,omitempty" json:"promo_code,omitempty"`
TotalAmountPaid int64  `bson:"total_amount_paid" json:"total_amount_paid"`
```

### Add `Refund`

File:
- `internal/models/refund.go`

Fields:
- `ID`
- `SubscriptionID`
- `MemberID`
- `UsedSessions`
- `RemainingSessions`
- `RefundAmount`
- `Reason`
- `Status`
- `CreatedAt`
- `ProcessedAt`

## Repository plan

### New file

- `internal/repository/refund_repo.go`

Interface:
- `Create(ctx context.Context, refund *models.Refund) error`
- `GetBySubscriptionID(ctx context.Context, subscriptionID string) (*models.Refund, error)` optional

### Update subscription repo

Add:

```go
RefundSubscription(ctx context.Context, id string) error
```

Atomic filter:

```go
bson.M{
  "_id": objID,
  "status": bson.M{"$in": []string{"active", "suspended"}},
  "remaining_sessions": bson.M{"$gt": 0},
}
```

Update:

```go
bson.M{
  "$set": bson.M{
    "status": "refunded",
    "remaining_sessions": 0,
  },
}
```

## Service plan

Update constructor:
- inject `RefundRepository`

Add errors:
- `ErrInvalidDiscount`
- `ErrSubscriptionCannotRefund`
- `ErrSubscriptionNoRemainingSessions`
- `ErrRefundAlreadyExists`

Add method:

```go
RefundSubscription(ctx context.Context, id string, reason string) (*models.Refund, error)
```

Update `CreateSubscription`:
- validate pricing input.
- calculate subtotal/discount/total.
- set snapshot fields.

Refund flow:
1. Validate ObjectID by repo or service.
2. Load subscription.
3. Validate status and sessions.
4. Calculate refund.
5. Atomic update subscription to refunded.
6. Insert refund record.
7. Return refund.

## Handler plan

Update create request DTO:
- `DiscountType string`
- `DiscountValue int64`
- `PromoCode string`

Add refund request DTO:
```go
type refundSubscriptionRequest struct {
    Reason string `json:"reason"`
}
```

Add handler:
```go
func (h *SubscriptionHandler) Refund(c *gin.Context)
```

Error mapping:
- invalid input/discount/id → 400
- subscription not found → 404
- cannot refund/no remaining/already refunded → 409
- unknown → 500

## Route plan

In `cmd/server/main.go`:

```go
api.POST("/subscriptions/:id/refund", subscriptionHandler.Refund)
```

## Docs/test plan

Update:
- `docs/api_contract.md`
- `api_test.http`
- `CHAT_CONTEXT/README.md`
- `CHAT_CONTEXT/backend_skills/worklog.md`

Commands:
```bash
gofmt -w internal/models/subscription.go internal/models/refund.go internal/repository/subscription_repo.go internal/repository/refund_repo.go internal/service/subscription_service.go internal/handlers/subscription_handler.go cmd/server/main.go
go build ./...
go test ./...
```

Manual API flow:
1. Create branch.
2. Create course.
3. Create member.
4. Create subscription with discount.
5. Activate member/subscription.
6. Check-in once.
7. Refund subscription.
8. GET subscription, verify:
   - `status = refunded`
   - `remaining_sessions = 0`

## Risks

- No transaction wrapper exists. Atomic subscription update + refund insert can become partial if insert fails.
- Refactor `Total_Amount_Paid` can touch multiple files.
- Unique refund index may need bootstrap work.
- Existing API clients may ignore new pricing fields safely because they are additive.