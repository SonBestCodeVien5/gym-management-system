# Cycle 07 — Integration Tests & Fixtures

## Status

- Status: planned
- Priority: medium
- Depends on: core backend features stable

## Goal

Thêm integration tests và fixtures để backend không phụ thuộc test tay bằng `api_test.http`.

## Test strategy

Use Go tests with:
- `testing`
- `net/http/httptest`
- Gin router setup
- Mongo test database

## Test utilities plan

Create:
- `internal/testutil/router.go`
- `internal/testutil/mongo.go`
- `internal/testutil/fixtures.go`

Responsibilities:
- setup test DB with unique name
- clean collections before/after tests
- build router with real repositories/services/handlers
- seed branch/course/member/subscription/session data

## Test flows

### Core subscription flow

1. Create branch.
2. Create course.
3. Create member.
4. Create subscription pending.
5. Activate subscription.
6. Get subscription.
7. List member subscriptions.

### Attendance flow

1. Active subscription.
2. Check-in.
3. Report missed.
4. Makeup from valid report.
5. Reject reused report.

### Sessions flow

1. Create session.
2. Enroll subscription.
3. Check-in session.
4. Reject over-capacity or invalid tag.

### Refund flow

1. Active subscription.
2. Use one session via attendance.
3. Refund subscription.
4. Verify status/refund amount/remaining sessions.

### Branch nearby flow

1. Seed branches with coordinates.
2. Query nearby.
3. Verify sorted/distance if available.

### Auth flow

1. Seed employee.
2. Login.
3. Access protected route.
4. Refresh.
5. Logout.

## Commands

```bash
go test ./...
go build ./...
```

## Fixture data

Need stable sample:
- branch HCM with GeoJSON Point
- course with base price, session count, allowed tags
- member with unique ccid
- subscription active
- session with tags matching subscription allowed tags
- employee admin/trainer/receptionist

## CI-ready later

Future:
- add GitHub Actions
- start MongoDB service
- run `go test ./...`
- run `go build ./...`

## Docs update

Update:
- `docs/local_dev_guide.md`
- `docs/code_reading_guide.md`
- `CHAT_CONTEXT/README.md`
- `worklog.md`

## Risks

- Tests need MongoDB running.
- Test DB cleanup must not touch dev DB.
- Time-based rules need deterministic timestamps.