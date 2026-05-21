# Cycle 04 - Auth/Login + Role Guard

## Status

- Status: planned
- Priority: medium
- Roadmap position: next backend cycle after attendance report/makeup completion
- Depends on: existing `Employee` model being extended into an authenticatable staff record
- Endpoints:
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/refresh`
  - `POST /api/v1/auth/logout`

## Goal

Thêm đăng nhập cho employee/staff, access token + refresh token có thể revoke, và middleware role
guard cho các API nghiệp vụ hiện đang public.

## Current baseline

- `docs/api_contract.md` đã liệt kê ba auth endpoint ở trạng thái `Planned`.
- `cmd/server/main.go` hiện chỉ có `/ping` và `/api/v1/*`; toàn bộ member, course, branch,
  subscription, attendance, session route đều chưa có auth middleware.
- `internal/models/employee.go` đã tồn tại với:
  - `employee_id`, `full_name`, `email`, `phone`, `level`
  - `role []string`
  - `branch_id []ObjectID`
- Chưa có employee repository/service/handler cho login.
- Chưa có `password_hash`, employee status, auth timestamps, refresh-token collection, token
  config, hoặc JWT dependency trong route/service wiring hiện tại.
- Handler hiện tại thường trả JSON theo shape `message` + `data`; auth response nên giữ cùng style
  trong cycle này thay vì đổi toàn backend trước cycle 05.

## Scope decisions

- Giữ role assignment dạng danh sách để khớp model hiện tại; guard pass nếu employee có ít nhất
  một role được phép.
- Giữ branch assignment dạng danh sách trong employee model; branch-scope authorization chưa áp
  dụng ở cycle này vì feature này chỉ chốt identity + role guard.
- Không thêm endpoint quản trị employee trong cycle này.
- Bootstrap admin đầu tiên bằng env-based bootstrap cho local/dev MVP:
  - đọc bootstrap values từ env khi server start
  - chỉ tạo admin khi account bootstrap chưa tồn tại
  - hash password bằng bcrypt trước khi persist
  - không log bootstrap password hoặc password hash
- Chỉ `/ping`, login, và refresh là public trong cycle này. Logout revoke bằng refresh-token body
  không cần access token; toàn bộ business routes còn lại phải protected theo role matrix bên dưới.
- Nếu product/FE sau này cần public catalog trước login, mở riêng `GET /api/v1/courses`,
  `GET /api/v1/branches`, hoặc `GET /api/v1/branches/nearby` bằng contract update có chủ đích,
  không để route public do wiring ngẫu nhiên.

## API contract plan

### Login

```http
POST /api/v1/auth/login
```

Request:

```json
{
  "email": "admin@gym.test",
  "password": "secret"
}
```

Response `200`:

```json
{
  "message": "login successful",
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "employee": {
      "id": "ObjectID",
      "employee_id": "EMP001",
      "email": "admin@gym.test",
      "full_name": "Gym Admin",
      "role": ["admin"],
      "branch_id": []
    }
  }
}
```

Status codes:

- `200`: credentials accepted and token pair issued.
- `400`: invalid body or missing email/password.
- `401`: invalid credentials or inactive employee.
- `500`: token/storage/internal failure.

### Refresh

```http
POST /api/v1/auth/refresh
```

Request:

```json
{
  "refresh_token": "..."
}
```

Response `200`:

```json
{
  "message": "token refreshed successfully",
  "data": {
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

Status codes:

- `200`: refresh token valid; issue a new access token and rotate refresh token.
- `400`: invalid body or missing refresh token.
- `401`: invalid, expired, revoked, or unknown refresh token.
- `500`: token/storage/internal failure.

### Logout

```http
POST /api/v1/auth/logout
```

Request:

```json
{
  "refresh_token": "..."
}
```

Response `200`:

```json
{
  "message": "logout successful"
}
```

Status codes:

- `200`: known refresh token revoked; repeated logout should stay idempotent.
- `400`: invalid body or missing refresh token.
- `401`: malformed or unverifiable refresh token.
- `500`: token/storage/internal failure.

## Protected-route plan

Public:

- `GET /ping`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`

Authenticated route with no extra role decision:

- `POST /api/v1/auth/logout` uses the refresh token body for revoke and does not require an access
  token.

Role guard matrix:

| Surface | Roles |
|---|---|
| Course CRUD | `admin`, `manager` |
| Branch CRUD | `admin`, `manager` |
| Member create/get/list-subscriptions/activate | `admin`, `manager`, `receptionist` |
| Subscription create/get/refund/suspend/unsuspend/expire | `admin`, `manager`, `receptionist` |
| Attendance checkin/report/makeup/history | `admin`, `manager`, `receptionist` |
| Session create | `admin`, `manager`, `trainer` |
| Session list/get/enroll/checkin | `admin`, `manager`, `trainer` |

Middleware behavior:

- Missing/invalid/expired access token returns `401`.
- Authenticated employee without an allowed role returns `403`.
- Middleware stores trusted employee ID and role list in Gin context; handlers must not trust
  client-sent role fields.
- Register static auth routes before protected API groups and preserve existing static route order
  such as `/branches/nearby` before `/branches/:id`.

## Business and security rules

- Normalize login email before lookup; do not reveal whether email or password failed.
- Never store or return plain passwords.
- Hash employee passwords with bcrypt before persistence.
- Access token is short lived and signed from `JWT_ACCESS_SECRET`.
- Refresh token is longer lived, signed from `JWT_REFRESH_SECRET`, and persisted only as a hash for
  revocation/rotation checks.
- Refresh rotates on success: conditionally revoke the presented non-revoked refresh token and only
  persist/return the replacement token when that revoke succeeds.
- Logout revokes the matching non-revoked refresh token hash and should not leak whether a token was
  already revoked.
- Access-token claims should carry only trusted auth data needed by middleware, at least employee ID
  and roles; repository remains source of truth when status must be checked.
- Reject login/refresh for inactive staff once employee status is added.
- Never log passwords, raw tokens, token hashes, or JWT secrets.

## Data and query plan

### Employee

Extend the existing employee model instead of replacing it:

- Keep current `employee_id`, `full_name`, `role []string`, `level`, `phone`, `email`,
  `branch_id []ObjectID`.
- Add `password_hash`.
- Add status field, defaulting to an enabled value such as `active`.
- Add `created_at` and `updated_at` if employee persistence is introduced in this cycle.
- Validate roles against:
  - `admin`
  - `manager`
  - `trainer`
  - `receptionist`

Repository needs:

- Lookup employee by normalized email for login.
- Get employee by ID when middleware/service must confirm staff state.
- Bootstrap one admin from env when no bootstrap admin exists, without exposing plain-password seed
  data in source.

Index/query decisions:

- Create the unique normalized-email index in cycle 04 so login identity stays deterministic.
- Keep employee ID unique for admin bootstrap/ops flows if repository bootstrap uses it as a lookup
  key.

### Refresh token

Add refresh-token storage:

- `id`
- `employee_id`
- `token_hash`
- `expires_at`
- `revoked_at` optional
- `created_at`
- `updated_at`

Repository needs:

- Create refresh token hash at login and after rotation.
- Find non-revoked token by hash.
- Conditionally revoke one non-revoked token for refresh so replayed/double-submit refresh calls do
  not both issue replacements.
- Revoke one token for logout.
- Reject expired records even if TTL cleanup is not yet present.

Index/query decisions:

- Create unique `token_hash` index in cycle 04 to keep refresh-token identity deterministic.
- TTL cleanup on `expires_at` is optional in cycle 04 and may stay tracked under cycle 07; service
  expiry validation is mandatory either way.

## Layer plan

### Models

- Extend `internal/models/employee.go` with auth fields/status/timestamps while preserving current
  role and branch shapes.
- Add refresh-token model under `internal/models/`.

### Repository

- Add employee repository for email/ID lookup and the env bootstrap persistence helper for the first
  admin account.
- Add refresh-token repository for create/find/revoke/rotation support.
- Keep Mongo collection access and index bootstrap in repository layer.

### Service

- Add auth service for credential verification, token pair issue, refresh rotation, logout revoke,
  role/status validation inputs, and service-level auth errors.
- Keep secret/TTL parsing close to auth config/service wiring instead of handlers.
- Define explicit errors for invalid credentials, invalid token, inactive employee, and token storage
  failures so handler status mapping stays predictable.

### Handler and middleware

- Add auth handler for request binding and HTTP status mapping.
- Add `AuthRequired` middleware for access-token validation and trusted Gin context values.
- Add `RequireRoles` middleware accepting allowed roles and matching against the role list from auth
  context.
- Wire public auth routes and protected route groups in `cmd/server/main.go`.

## Docs and verification plan

Update after implementation changes HTTP behavior:

- `.env.example` with JWT secret, TTL, and bootstrap admin config names selected in implementation.
- `docs/api_contract.md` with auth details, protected-route notes, `401`, and `403` behavior.
- `api_test.http` with login, refresh, logout, authorized call, missing-token call, and forbidden
  role call samples.
- `CHAT_CONTEXT/README.md` when the implemented backend surface and next resume point change.
- `CHAT_CONTEXT/backend_skills/worklog.md` when feature status advances beyond planned.

Automated checks:

```bash
go build ./...
go test ./...
```

Test coverage to add:

- Login success, wrong password, missing fields, inactive employee.
- Refresh success with rotation, expired token, revoked token, unknown token.
- Logout revoke and repeated logout behavior.
- Auth middleware `401` cases for missing, malformed, and expired access token.
- Role middleware `403` case plus allowed-role pass case.
- Route-level coverage proving selected CRUD/action routes are protected and public auth routes still
  work.

Manual/API checks:

- Bootstrap one admin from env and prepare one restricted-role employee for role checks.
- Login and call one allowed route with `Authorization: Bearer <access_token>`.
- Call the same protected route without access token and with a role that is not allowed.
- Refresh once, confirm old refresh token no longer works, then logout the replacement token.

## Risks and follow-up decisions

- Bootstrap admin env values must be configured carefully; the bootstrap path must not overwrite an
  existing account or expose credentials in logs.
- Employee model already uses multi-role and multi-branch fields; changing them to single values in
  implementation would broaden migration cost and needs an explicit decision.
- Refresh rotation remains double-submit-sensitive; implementation must prove the conditional revoke
  path with tests.
- JWT token/session tests need stable TTL and time control to avoid flaky expiry assertions.
- Protecting every current route will break existing unauthenticated `api_test.http` flows until auth
  samples and headers are added.
- Error response shape should follow current handler style in this cycle; broad response
  normalization remains cycle 05 work.
