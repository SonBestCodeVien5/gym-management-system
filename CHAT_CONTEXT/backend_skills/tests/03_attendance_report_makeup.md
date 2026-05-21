# Test - attendance report makeup

## Status

- Status: tested with automated build/test, manual API checks, and direct MongoDB verification
- Feature: attendance report makeup
- Plan file: `CHAT_CONTEXT/backend_skills/plans/03_attendance_report_makeup.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/03_attendance_report_makeup.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/03_attendance_report_makeup.md`
- Tested at: 2026-05-21
- Final result: pass

## Commands

- `go build ./...` - pass
  - Exit code `0`.
  - Sandbox run emitted a read-only Go module stat-cache warning, but build completed successfully.
- `go test ./...` - pass

```txt
?   	github.com/SonBestCodeVien5/gym-management-system/cmd/server	[no test files]
?   	github.com/SonBestCodeVien5/gym-management-system/internal/handlers	[no test files]
?   	github.com/SonBestCodeVien5/gym-management-system/internal/models	[no test files]
?   	github.com/SonBestCodeVien5/gym-management-system/internal/repository	[no test files]
ok  	github.com/SonBestCodeVien5/gym-management-system/internal/service	(cached)
?   	github.com/SonBestCodeVien5/gym-management-system/pkg/database	[no test files]
```

- Manual server bring-up - pass
  - Started `PORT=18080 go run cmd/server/main.go` against local MongoDB.
  - `GET /ping` returned `200`.
  - Local API and Mongo calls required outside-sandbox execution because sandboxed socket access to `127.0.0.1:27017` was denied.

## Manual Setup

Fresh records were created through the API for this run:

- branch_id: `6a0edfb36de3c53cea94c431`
- course_id: `6a0edfb36de3c53cea94c430`
- member_id: `6a0edfb16de3c53cea94c42f`
- active main subscription_id: `6a0edfd46de3c53cea94c432`
- active overdue-window subscription_id: `6a0edfd56de3c53cea94c433`
- pending subscription_id: `6a0edfd66de3c53cea94c434`

The main and overdue-window subscriptions were activated through `PATCH /api/v1/members/:id/activate`. The pending subscription was intentionally left pending for conflict coverage.

## Manual API Results

| Case | Expected | Observed | Result |
|---|---|---|---|
| Report missed on active main subscription | `201` | `201`, `status = "reported_missed"` | pass |
| Makeup for exact report date | `201` | `201`, `status = "makeup"` and `is_makeup_for` preserved | pass |
| Report with invalid `subscription_id` | `400` | `400`, `invalid subscription id` | pass |
| Report with invalid `branch_id` | `400` | `400`, `invalid branch id` | pass |
| Report with invalid `date` | `400` | `400`, `invalid date format` | pass |
| Makeup missing `is_makeup_for` | `400` | `400`, `is_makeup_for is required` | pass |
| Report with unknown valid subscription ObjectID | `404` | `404`, `subscription not found` | pass |
| Report on pending subscription | `409` | `409`, subscription status conflict | pass |
| Second report inside 30-day window | `409` | `409`, `reported missed limit reached within 30 days` | pass |
| Makeup with no matching reported-missed date | `409` | `409`, `makeup reference not found` | pass |
| Makeup more than 7 days after source report | `409` | `409`, `invalid makeup reference` | pass |
| Reuse same report date for second makeup | `409` | `409`, `makeup reference already used` | pass |

Observed happy-path payloads:

```json
{"data":{"id":"6a0ee05f6de3c53cea94c436","sub_id":"6a0edfd46de3c53cea94c432","branch_id":"6a0edfb36de3c53cea94c431","date":"2026-05-12T08:00:00Z","status":"reported_missed","is_makeup_for":null},"message":"attendance report recorded successfully"}
```

```json
{"data":{"id":"6a0ee0b96de3c53cea94c437","sub_id":"6a0edfd46de3c53cea94c432","branch_id":"6a0edfb36de3c53cea94c431","date":"2026-05-14T08:00:00Z","status":"makeup","is_makeup_for":"2026-05-12T08:00:00Z"},"message":"attendance makeup recorded successfully"}
```

State reads after the main flow:

- `GET /api/v1/subscriptions/:id/attendance` returned the one `reported_missed` record and the one `makeup` record for the main subscription.
- `GET /api/v1/subscriptions/:id` returned `remaining_sessions = 11` from an initial `12`.
- `GET /api/v1/members/:id` returned `total_sessions_attended = 1`.

## DB Verification

Direct Mongo verification was run through `mongosh` in the local `gym_mongodb` container.

- [x] Main subscription stored exactly the successful report and makeup records.
- [x] Report record stored `status = "reported_missed"` and `is_makeup_for = null`.
- [x] Makeup record stored `status = "makeup"` and `is_makeup_for = ISODate("2026-05-12T08:00:00Z")`.
- [x] Main subscription stored `remaining_sessions = 11` after one successful makeup.
- [x] Member stored `total_sessions_attended = 1` after one successful makeup.
- [x] Overdue-window subscription stored only its source report; the rejected overdue makeup did not insert an attendance record.

Direct DB snapshot:

```txt
main_attendances:
- 2026-05-12T08:00:00Z reported_missed is_makeup_for=null
- 2026-05-14T08:00:00Z makeup is_makeup_for=2026-05-12T08:00:00Z
main_subscription:
- status=active total_sessions=12 remaining_sessions=11
member:
- total_sessions_attended=1
overdue_attendances:
- 2026-05-01T08:00:00Z reported_missed only
```

## Failures And Blockers

- Failures observed in this run: none.
- Blockers for the requested Cycle 03 verification: none.
- Previous partial blocker note is cleared in this run:
  - unknown valid subscription ObjectID returned `404`, not `400`.
  - the overdue makeup-window case was executed and returned `409`.

## Remaining Risks

- No feature-specific integration tests exist yet for report/makeup endpoints.
- Duplicate makeup prevention is still service-level and not protected by a DB-level unique/atomic guard for concurrent double-submit.
- Makeup lookup still depends on exact RFC3339 instant equality for `is_makeup_for`.
