# Cycle 05 — Validation Hardening & Error Consistency

## Status

- Status: planned
- Priority: medium
- Depends on: core feature endpoints stabilized

## Goal

Chuẩn hóa validation và response lỗi toàn backend.

## Problem

Hiện handler/service có nhiều cách trả lỗi khác nhau. Cần một pattern thống nhất để FE dễ xử lý.

## Error response contract

Proposed shape:

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "invalid subscription input",
    "details": {}
  }
}
```

## Error code plan

- `INVALID_INPUT`
- `INVALID_ID`
- `INVALID_DATE`
- `NOT_FOUND`
- `CONFLICT`
- `UNAUTHORIZED`
- `FORBIDDEN`
- `INTERNAL_ERROR`

## HTTP status mapping

- invalid input/date/id → 400
- unauthorized/missing token → 401
- forbidden/role denied → 403
- not found → 404
- conflict/business rule → 409
- unknown → 500

## Implementation plan

Create:
- `internal/handlers/response.go`

Functions:
- `RespondOK`
- `RespondCreated`
- `RespondError`
- `RespondValidationError`

Optional:
- central error mapper per domain or shared helper.

## Validation hardening plan

- Validate ObjectID before repo query where possible.
- Validate RFC3339 date strings in handler.
- Validate GeoJSON branch type and coordinate range.
- Validate money/count not negative.
- Validate enum values:
  - subscription status
  - discount type
  - attendance type/status
  - employee role
- Hide raw Mongo duplicate key errors.

## Migration plan

Do per handler:
1. subscription
2. attendance
3. branch
4. course
5. member
6. session
7. auth

## Docs/test plan

Update:
- `docs/api_contract.md`
- `api_test.http`
- `CHAT_CONTEXT/README.md`
- `worklog.md`

Run:
```bash
go build ./...
go test ./...
```

## Risks

- Broad change can touch many handlers.
- Response shape changes may break existing manual clients.
- Do after feature endpoints stable.