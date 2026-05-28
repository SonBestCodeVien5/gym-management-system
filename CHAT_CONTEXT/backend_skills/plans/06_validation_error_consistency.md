# Cycle 06 - Validation Hardening & Error Consistency

## Status

- Status: completed
- Priority: medium
- Depends on: employee management completed and pushed
- Scope: backend HTTP error contract and validation hardening only

## Goal

Chuẩn hóa response lỗi và validation HTTP toàn backend để FE/manual clients có thể xử lý lỗi theo
`error.code` thay vì phải parse `message` tự do.

Cycle này không đổi success response shape nếu không cần thiết. Các response thành công hiện tại vẫn
giữ dạng:

```json
{
  "message": "operation successful",
  "data": {}
}
```

hoặc chỉ:

```json
{
  "message": "operation successful"
}
```

## Current baseline

Hiện status code nhìn chung đã hợp lý, nhưng body lỗi chưa đồng nhất:

- Bind JSON lỗi có nơi trả `{"message":"invalid request body","error":"<raw bind error>"}`.
- Validation/service lỗi thường trả `{"message":"..."}`.
- Auth middleware trả `{"message":"missing access token"}` / `{"message":"invalid access token"}`.
- Một số path/query ObjectID, RFC3339 date, nearby query parse lỗi đang tự viết từng handler.
- Employee duplicate đã được map qua `repository.ErrDuplicate`, nhưng cycle này cần giữ nguyên nguyên
  tắc không trả raw Mongo duplicate-key/driver text ra API.

Các file handler chính cần migrate:

- `internal/handlers/auth_handler.go`
- `internal/handlers/auth_middleware.go`
- `internal/handlers/member_handler.go`
- `internal/handlers/subscription_handler.go`
- `internal/handlers/attendance_handler.go`
- `internal/handlers/session_handler.go`
- `internal/handlers/course_handler.go`
- `internal/handlers/branch_handler.go`
- `internal/handlers/employee_handler.go`

## Error response contract

Tất cả response lỗi mới dùng shape:

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "invalid subscription input",
    "details": {}
  }
}
```

Rules:

- `error.code` là enum ổn định cho client.
- `error.message` là thông điệp ngắn, đã sanitize, có thể hiển thị hoặc log client-side.
- `error.details` luôn là object; mặc định `{}`.
- Không expose `err.Error()` từ JSON binder, Mongo driver, JWT parser, bcrypt, hoặc storage internals.
- Có thể dùng `details` cho field-level lỗi đơn giản như `{ "field": "start_date" }`, nhưng không bắt
  buộc phải enrich toàn bộ validation trong cycle này.

## Error codes and HTTP mapping

| Code | HTTP | Dùng khi |
|---|---:|---|
| `INVALID_INPUT` | 400 | Body sai, thiếu required field, enum sai, money/count âm, business validation input sai |
| `INVALID_ID` | 400 | Path/query/body ObjectID không parse được |
| `INVALID_DATE` | 400 | RFC3339 date/datetime không parse được |
| `UNAUTHORIZED` | 401 | Missing/malformed/expired/inactive access token, invalid credentials, invalid refresh token |
| `FORBIDDEN` | 403 | Đã auth nhưng role không đủ quyền |
| `NOT_FOUND` | 404 | Resource/reference không tồn tại |
| `CONFLICT` | 409 | Duplicate unique field, trạng thái không hợp lệ cho action, self-lockout, weekly limit, double-use |
| `INTERNAL_ERROR` | 500 | Storage/token/internal failure không mong muốn |

Ghi chú:

- Invalid JSON/body cũng dùng `INVALID_INPUT`, không thêm code riêng để tránh nở enum sớm.
- Service error input hiện tại như `ErrInvalidSubscriptionInput`, `ErrInvalidEmployeeInput`,
  `ErrInvalidDiscount` map sang `INVALID_INPUT`.
- Service conflict như duplicate, already suspended, not active, weekly limit, self-deactivation map
  sang `CONFLICT`.

## API contract impact

Update `docs/api_contract.md` với section global `Error response`.

Ví dụ `400` invalid ObjectID:

```json
{
  "error": {
    "code": "INVALID_ID",
    "message": "invalid subscription id",
    "details": {}
  }
}
```

Ví dụ `401`:

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid access token",
    "details": {}
  }
}
```

Ví dụ `409`:

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "weekly session limit reached",
    "details": {}
  }
}
```

## Layer plan

### Handler layer

Create `internal/handlers/response.go`:

- Error code constants.
- `RespondOK(c, message, data)` and `RespondCreated(c, message, data)` helpers may be added, but
  success-response migration is optional and should not create noisy diffs.
- `RespondMessage(c, status, message)` or equivalent helper for success without `data` if useful.
- `RespondError(c, status, code, message, details)`.
- Convenience wrappers:
  - `RespondInvalidInput`
  - `RespondInvalidID`
  - `RespondInvalidDate`
  - `RespondUnauthorized`
  - `RespondForbidden`
  - `RespondNotFound`
  - `RespondConflict`
  - `RespondInternal`
- Optional parse helpers if they reduce duplication without hiding handler intent:
  - `ParseObjectID(value, fieldName)`
  - `ParseObjectIDParam(c, name, label)`
  - `ParseRFC3339(value, fieldName)`

Handler responsibility:

- Parse request body/path/query.
- Validate syntactic HTTP input: ObjectID, RFC3339, numeric query parse.
- Map service errors to status + `error.code`.
- Never return raw binder/driver/token internals.

### Service layer

Keep existing business-rule validation in services:

- subscription pricing/discount/status transitions
- attendance weekly limit/report/makeup windows
- course required fields and allowed values
- branch GeoJSON/nearby constraints
- employee role/status/level/password/branch references/self-lockout
- auth credential/token business errors

Do not introduce HTTP response shapes into services.

### Repository layer

Keep repositories storage-only.

- Continue mapping duplicate-key errors to storage-agnostic `repository.ErrDuplicate` where needed.
- Do not return HTTP codes or API error structs from repositories.

## Validation hardening plan

- ObjectID:
  - path params: `:id` across members/subscriptions/courses/branches/sessions/employees
  - body fields: member/course/home branch/subscription/session/employee branch IDs
  - query fields: employee `branch_id`, session `branchId`, nearby query numeric fields where relevant
- RFC3339 date/datetime:
  - subscription `start_date`, `end_date`
  - attendance `date`, `is_makeup_for`
  - session `scheduled_at`, list `date`
- Branch/GeoJSON:
  - Keep service validation for type `Point`, exactly two coordinates, lng/lat ranges.
  - Handler nearby query parse errors should return `INVALID_INPUT` with sanitized messages.
- Numeric constraints:
  - Keep service validation for price, session counts, discount value, coordinates, capacity, weekly
    count.
  - Ensure parse errors and negative query values do not leak raw conversion text.
- Enum values:
  - subscription status/discount type
  - attendance type/status
  - employee role/status/level
  - branch status if present
- Duplicate/internal errors:
  - Duplicate unique constraints return `CONFLICT`.
  - Unknown Mongo/JWT/bcrypt/storage errors return `INTERNAL_ERROR`.

## Migration order

1. Add response/error helpers and focused tests if practical.
2. Migrate auth middleware first so protected-route failures use the new contract.
3. Migrate auth handler.
4. Migrate employee handler because it is the newest surface and has duplicate/role/security cases.
5. Migrate subscription and attendance handlers because they have the richest validation/status
   mapping.
6. Migrate session handler.
7. Migrate branch handler.
8. Migrate course and member handlers.
9. Run a final `rg` sweep for old error shapes:
   - `gin.H{"message":`
   - `"error": err.Error()`
   - `AbortWithStatusJSON`

## Docs and samples plan

Update:

- `docs/api_contract.md`
  - Add global error response section.
  - Clarify status codes refer to the new error body.
  - Keep endpoint success examples unchanged unless implementation actually changes them.
- `api_test.http`
  - Add or update representative error checks for invalid token, forbidden role, invalid ObjectID,
    invalid date, invalid body, not found, conflict.
- `CHAT_CONTEXT/backend_skills/implementations/06_validation_error_consistency.md` during
  implementation.
- `CHAT_CONTEXT/backend_skills/worklog.md` as phases advance.
- `CHAT_CONTEXT/README.md` only after completion or if the resume point changes.

## Verification plan

Automated:

```bash
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
git diff --check
```

Focused manual/API checks:

- Missing token on protected endpoint returns `401` + `UNAUTHORIZED`.
- Non-admin employee endpoint access returns `403` + `FORBIDDEN`.
- Invalid path ID returns `400` + `INVALID_ID`.
- Invalid RFC3339 date returns `400` + `INVALID_DATE`.
- Invalid JSON/body returns `400` + `INVALID_INPUT` without raw bind error.
- Missing reference/resource returns `404` + `NOT_FOUND`.
- Duplicate employee email/ID or known business conflict returns `409` + `CONFLICT`.
- Forced/observed unexpected storage error, if practical, returns `500` + `INTERNAL_ERROR`.

## Acceptance criteria

- All implemented backend error responses use the new `{"error":{...}}` contract.
- Success response shapes remain compatible with the current API contract.
- No raw binder, Mongo, JWT, bcrypt, or internal storage error text is returned to clients.
- Existing service/repository responsibilities remain intact.
- `go build ./...`, `go test ./...`, and `git diff --check` pass.
- API docs and REST samples match implemented behavior.

## Risks

- This is a broad cross-handler change; use the `rg` sweep to avoid leaving mixed contracts.
- FE/manual clients reading top-level `message` on errors will need to switch to `error.message`.
- Too much success-response helper migration could create noisy diffs without changing behavior; keep
  that optional.
- Sanitizing errors improves API safety but may reduce debugging detail; server logs can be expanded
  later if needed.
- This cycle does not solve transactional/data-integrity gaps such as refresh revoke atomicity,
  last-active-admin, or trainer validation for sessions.
