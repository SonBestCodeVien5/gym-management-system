# Backend Worklog

Dùng file này để ghi lại plan, implement, code review, test cho từng feature backend.

## Current backend roadmap

- [ ] Refund flow & pricing rules
- [ ] Branch nearby geo query
- [ ] Attendance report/makeup endpoints nếu route còn thiếu
- [ ] Auth/login + role guard
- [ ] Validation hardening & error consistency
- [ ] Indexes and data integrity
- [ ] Integration tests & fixtures

---

# Template

Copy block này cho mỗi feature.

```md
# Feature — <feature_name>

## Status
- Planned:
- Implemented:
- Reviewed:
- Tested:
- Docs updated:

## Plan — <date>

### Goal
...

### API
- Method:
- Path:
- Request:
- Response:
- Status codes:

### Business rules
- ...

### Data changes
- Model:
- Collection:
- Index:

### Implementation steps
1. ...
2. ...
3. ...

### Risks
- ...

## Implementation — <date>

### Files changed
- ...

### Notes
- ...

### Commands
- ...

## Code Review — <date>

### Passed
- ...

### Issues found
- ...

### Fixes applied
- ...

### Remaining risks
- ...

## Test — <date>

### Commands
- `go build ./...` — pass/fail
- `go test ./...` — pass/fail/skipped reason

### Manual API
- ...

### Results
- ...

### Issues
- ...

### Fixed
- ...

## Docs updated
- [ ] `docs/api_contract.md`
- [ ] `api_test.http`
- [ ] `CHAT_CONTEXT/README.md`
```

---

# Feature — Refund flow & pricing rules

## Status
- Planned: yes
- Implemented: no
- Reviewed: no
- Tested: no
- Docs updated: no

## Plan summary

### Goal
Implement `POST /api/v1/subscriptions/:id/refund` and pricing discount rules for subscription creation.

### API
- `POST /api/v1/subscriptions/:id/refund`
- Request:
```json
{
  "reason": "member requested cancellation"
}
```
- Response should include refund record or refund summary.

### Business rules
- Only `active` or `suspended` subscription can be refunded.
- Cannot refund `pending`, `expired`, `refunded`.
- Cannot refund if `remaining_sessions <= 0`.
- `used_sessions = total_sessions - remaining_sessions`.
- `refund_amount = total_amount_paid * remaining_sessions / total_sessions`.
- After refund:
  - subscription `status = refunded`
  - `remaining_sessions = 0`
  - refund record inserted.
- Prevent double refund via atomic update and/or unique index.

### Pricing rules
- Server calculates money from course snapshot.
- Optional discount:
  - `none`
  - `percent`
  - `fixed`
- Percent must be `0 <= value <= 100`.
- Fixed must be `0 <= value <= subtotal`.
- `total_amount_paid = subtotal - discount_amount`.

### Files expected
- `internal/models/subscription.go`
- `internal/models/refund.go`
- `internal/repository/subscription_repo.go`
- `internal/repository/refund_repo.go`
- `internal/service/subscription_service.go`
- `internal/handlers/subscription_handler.go`
- `cmd/server/main.go`
- `docs/api_contract.md`
- `api_test.http`

---

# Feature — Branch nearby geo query

## Status
- Planned: yes
- Implemented: no
- Reviewed: no
- Tested: no
- Docs updated: no

## Plan summary

### Goal
Implement `GET /api/v1/branches/nearby`.

### API
- `GET /api/v1/branches/nearby?lng=106.7&lat=10.8&max_distance=5000&limit=10`

### Business rules
- Validate lng/lat range.
- Default `max_distance = 5000`.
- Default `limit = 10`, max `100`.
- GeoJSON coordinate order is `[lng, lat]`.
- Route must be before `/branches/:id`.

### Data/index
- Mongo index: `branches.location` 2dsphere.

### Files expected
- `internal/repository/branch_repo.go`
- `internal/service/branch_service.go`
- `internal/handlers/branch_handler.go`
- `cmd/server/main.go`
- Mongo index bootstrap location
- `docs/api_contract.md`
- `api_test.http`