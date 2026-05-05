# Chat Context Snapshot

Read this first when continuing the project in a new chat.

## Current state
- Stack: Go + Gin + MongoDB, Clean Architecture.
- Member flow implemented:
  - `POST /api/v1/members`
  - `GET /api/v1/members/:id`
  - `PATCH /api/v1/members/:id/activate` for offline payment confirmation.
- Subscription flow implemented:
  - `POST /api/v1/subscriptions`
  - `GET /api/v1/subscriptions/:id`
  - `PATCH /api/v1/subscriptions/:id/suspend`
  - `PATCH /api/v1/subscriptions/:id/unsuspend`
  - `PATCH /api/v1/subscriptions/:id/expire`
- Course flow implemented:
  - `POST /api/v1/courses`
  - `GET /api/v1/courses`
  - `GET /api/v1/courses/:id`
  - `PATCH /api/v1/courses/:id`
  - `DELETE /api/v1/courses/:id`
- Branch flow implemented:
  - `POST /api/v1/branches`
  - `GET /api/v1/branches`
  - `GET /api/v1/branches/:id`
  - `PATCH /api/v1/branches/:id`
  - `DELETE /api/v1/branches/:id`
- Attendance flow implemented:
  - `POST /api/v1/attendance/checkin`
  - `GET /api/v1/subscriptions/:id/attendance`
- Repositories exist for:
  - member
  - course
  - branch
  - subscription
  - attendance
- `members.ccid` has a unique index.
- Subscription input parses `start_date` and `end_date` using RFC3339.
- Subscription creation currently:
  - validates member/course/branch references
  - sets `status = pending`
  - copies price and session count from `Course`
  - does NOT apply discounts/refunds yet.
- Offline payment is handled by a separate member activation endpoint, not inside subscription creation.

## Testing notes
- `api_test.http` contains sample requests for ping, member registration, member activation, subscription, course/branch CRUD, and attendance.
- Subscription testing needs real `course_id` and `branch_id`, but now there are create APIs for both.
- `go build ./...` was last verified to pass.

## Recommended next step
- Implement remaining Phase 2 rules: `sessionPerWeek`, attendance report/makeup, refund rules, and branches/nearby.

## Quick resume prompt
- "Read `CHAT_CONTEXT/README.md` first, then continue from the current project state."
