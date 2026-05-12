# Chat Context Snapshot

Read this first when continuing the project in a new chat.

## Current state
- Stack: Go + Gin + MongoDB, Clean Architecture.
- Current focus: backend feature completion, with Sessions now implemented and the next work centered on refund/pricing and branch nearby search.
- Member flow implemented:
   - `POST /api/v1/members`
   - `GET /api/v1/members/:id`
   - `GET /api/v1/members/:id/subscriptions` to list member subscriptions.
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
- Sessions workflow implemented:
  - `POST /api/v1/sessions`
  - `GET /api/v1/sessions`
  - `GET /api/v1/sessions/:id`
  - `POST /api/v1/sessions/:id/enroll`
  - `POST /api/v1/sessions/:id/checkin`
- Repositories exist for:
  - member
  - course
  - branch
  - subscription
  - attendance
  - session
- `members.ccid` has a unique index.
- Subscription input parses `start_date` and `end_date` using RFC3339.
- Subscription creation currently:
  - validates member/course/branch references
  - sets `status = pending`
  - copies price and session count from `Course`
  - does NOT apply discounts/refunds yet.
- Offline payment is handled by a separate member activation endpoint, not inside subscription creation.
- Attendance check-in now enforces `sessionPerWeek` for `attended` and `makeup` records.
- `reported_missed` now enforces a 30-day sliding window.
- `makeup` now requires a valid `reported_missed` reference within 7 days and cannot reuse the same report twice.

## Session architecture notes
- `Course.allowed_tags` defines which tags are allowed for a course.
- `Subscription.allowed_tags` snapshots the course tags at purchase time.
- `Session` does not have `course_id` anymore.
- Session enrollment stores enrolled subscription IDs directly on the session document.
- Enrollment is atomic in MongoDB, so the last slot cannot be double-booked by concurrent requests.
- Enrollment tag validation is "any match is enough": if at least one `session.tags` item exists in `subscription.allowed_tags`, enroll is allowed.
- Session check-in reuses the existing attendance service, so the same attendance rules still apply.
- `attendance.session_id` is optional and supports both free check-in and session-based check-in.

## Docs alignment
- Current vs planned API contract snapshot: see [docs/api_contract.md](docs/api_contract.md).
- Phase 2 design docs aligned to the current contract (including auth/refund/nearby placeholders).

## Testing notes
- `api_test.http` contains sample requests for ping, member registration, member activation, subscription, course/branch CRUD, and attendance.
- Subscription testing needs real `course_id` and `branch_id`, but now there are create APIs for both.
- Sessions workflow currently covers create/list/get/enroll/checkin.
- Enroll uses atomic Mongo update plus tag allow-list validation.
- Check-in creates attendance with `session_id` and then applies the normal attendance rules.
- `go build ./...` was last verified to pass.

## Recommended next step
- Implement refund flow and pricing rules next.
- After that: branch nearby geo query, auth/login + role guard, validation/error consistency, indexes/data integrity, then integration tests and fixtures.

## Todo list (current)
- [x] Chuan hoa API contract & docs
- [x] Enforce sessionPerWeek rule
- [x] Report/Makeup attendance rules
- [x] Sessions enroll/checkin workflow
- [x] Subscriptions list by member
- [ ] Refund flow & pricing rules
- [ ] Branch nearby geo query
- [ ] Auth/login + role guard
- [ ] Validation hardening & error consistency
- [ ] Indexes and data integrity
- [ ] Integration tests & fixtures

## Quick resume prompt
- "Read `CHAT_CONTEXT/README.md` first, then continue from the current project state."
