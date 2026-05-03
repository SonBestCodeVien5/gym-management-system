# Chat Context Snapshot

Read this first when continuing the project in a new chat.

## Current state
- Stack: Go + Gin + MongoDB, Clean Architecture.
- Member flow implemented:
  - `POST /api/v1/registration`
  - `GET /api/v1/members/:id`
  - `PATCH /api/v1/members/:id/activate` for offline payment confirmation.
- Subscription flow implemented:
  - `POST /api/v1/subscriptions`
  - `GET /api/v1/subscriptions/:id`
- Repositories exist for:
  - member
  - course
  - branch
  - subscription
- `members.ccid` has a unique index.
- Subscription input parses `start_date` and `end_date` using RFC3339.
- Subscription creation currently:
  - validates member/course/branch references
  - sets `status = active`
  - copies price and session count from `Course`
  - does NOT apply discounts/refunds yet.
- Offline payment is handled by a separate member activation endpoint, not inside subscription creation.

## Testing notes
- `api_test.http` contains sample requests for ping, member registration, member activation, and subscription.
- Subscription testing needs real seeded `course_id` and `branch_id` because there is no create API for those yet.
- `go build ./...` was last verified to pass.

## Recommended next step
- Continue with the payment/offline-confirm flow or move into subscription status/refund/suspension logic.

## Quick resume prompt
- "Read `CHAT_CONTEXT/README.md` first, then continue from the current project state."
