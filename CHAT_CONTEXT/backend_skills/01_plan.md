# Skill 01 — Plan Backend Feature

Dùng skill này trước khi code feature backend.

## Input cần đọc

- `CHAT_CONTEXT/README.md`
- `docs/api_contract.md`
- Model liên quan trong `internal/models`
- Handler liên quan trong `internal/handlers`
- Service liên quan trong `internal/service`
- Repository liên quan trong `internal/repository`
- Route wiring trong `cmd/server/main.go`

## Output của bước plan

Tạo plan ngắn nhưng đủ các phần:

1. Mục tiêu feature
2. API contract
3. Business rules
4. Data model changes
5. Repository changes
6. Service changes
7. Handler changes
8. Route changes
9. Error mapping
10. Tests/API samples
11. Risk/cạnh tranh dữ liệu nếu có

## Checklist plan

- [ ] Xác định endpoint/method/path.
- [ ] Xác định request body/query params.
- [ ] Xác định response success.
- [ ] Xác định error cases: 400, 404, 409, 500.
- [ ] Xác định model cần thêm/sửa.
- [ ] Xác định Mongo query/index cần thiết.
- [ ] Xác định business rules thuộc service.
- [ ] Xác định atomic update nếu có race condition.
- [ ] Xác định route order nếu có path conflict.
- [ ] Xác định docs/test cần cập nhật.

## Template plan

```md
# Plan — <feature_name>

## Goal
...

## API
- Method:
- Path:
- Request:
- Response:
- Status codes:

## Business rules
- ...

## Data changes
- Model:
- Collection:
- Index:

## Implementation steps
1. Model
2. Repository
3. Service
4. Handler
5. Route
6. Docs/API test
7. Build/test

## Risks
- ...
```

## Rule quan trọng

- Không code khi chưa biết repository/service hiện tại.
- Không đặt business rule trong handler.
- Không tin tiền/số buổi/status từ client nếu server có thể tự tính.
- Với refund/enroll/payment, ưu tiên update atomic hoặc unique index để chặn double-submit.

---

# Current Plan Progress

File này chỉ giữ tổng hợp quan trọng. Plan con theo từng chu kỳ nằm trong:

- `CHAT_CONTEXT/backend_skills/plans/`

## Roadmap chu kỳ backend

| Cycle | Feature | Plan file | Status |
|---|---|---|---|
| 01 | Refund flow & pricing rules | `plans/01_refund_pricing.md` | completed |
| 02 | Branch nearby geo query | `plans/02_branch_nearby.md` | completed |
| 03 | Attendance report/makeup endpoints | `plans/03_attendance_report_makeup.md` | completed |
| 04 | Auth/login + role guard | `plans/04_auth_role_guard.md` | planned |
| 05 | Validation hardening & error consistency | `plans/05_validation_error_consistency.md` | planned |
| 06 | Indexes and data integrity | `plans/06_indexes_data_integrity.md` | planned |
| 07 | Integration tests & fixtures | `plans/07_integration_tests_fixtures.md` | planned |

## Next action

Implement Cycle 04: Auth/login + role guard.

## Rule cập nhật tiến độ

Sau mỗi chu kỳ:
- Update status trong bảng roadmap này.
- Update plan con trong `plans/<cycle>.md`.
- Update implementation log trong `implementations/<cycle>.md`.
- Update review result trong `reviews/<cycle>.md`.
- Update test report trong `tests/<cycle>.md`.
- Update `worklog.md`.
- Update `CHAT_CONTEXT/README.md`.
