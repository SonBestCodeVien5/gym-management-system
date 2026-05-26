# Cycle 05 - Employee Management

## Status

- Status: completed
- Ngày plan: 2026-05-26
- Ngày complete: 2026-05-26
- Priority: medium
- Depends on: auth/login + role guard
- Phase tiếp theo: none; cycle complete

## Goal

Thêm API admin-only để tạo, xem, cập nhật, vô hiệu hóa và reset mật khẩu cho tài khoản
employee/staff sau khi Cycle 04 đã có bootstrap admin và auth.

Cycle này giúp admin quản lý staff account thật, không còn phụ thuộc hoàn toàn vào bootstrap env,
đồng thời không expose password hash hoặc cho client tự điều khiển các field auth do server tính.

## Baseline hiện tại

- `models.Employee` đã có:
  - `_id`, `employee_id`, `full_name`, `email`, `normalized_email`, `password_hash`, `status`
  - `role []string`, `level`, `phone`, `branch_id []ObjectID`
  - `created_at`, `updated_at`
- `EmployeeRepository` hiện có create, lookup theo ID, lookup theo normalized email, và bootstrap
  admin.
- Startup đã tạo unique sparse index cho:
  - `employees.normalized_email`
  - `employees.employee_id`
- Auth đã reload employee state khi validate access token và refresh token. Employee inactive sẽ bị
  reject.
- Chưa có employee management handler, service, route, API docs, hoặc REST sample.

## API contract

Tất cả employee management route yêu cầu access token và role `admin`.

| Method | Endpoint | Purpose |
|---|---|---|
| `POST` | `/api/v1/employees` | Tạo staff account với mật khẩu ban đầu. |
| `GET` | `/api/v1/employees` | Danh sách staff account. |
| `GET` | `/api/v1/employees/:id` | Xem một staff account. |
| `PATCH` | `/api/v1/employees/:id` | Cập nhật một phần profile, role, branch, level, email, employee ID, hoặc status. |
| `PATCH` | `/api/v1/employees/:id/password` | Admin reset mật khẩu employee. |

Không thêm `DELETE /employees/:id` trong cycle này. Offboarding dùng `status = inactive` để giữ ổn
định audit history, sessions, và auth records.

### Employee response

Mọi response employee phải bỏ `password_hash` và `normalized_email`.

```json
{
  "id": "69f20c000c4cd4cdf5768500",
  "employee_id": "EMP001",
  "full_name": "Tran Van Trainer",
  "email": "trainer@gym.test",
  "status": "active",
  "role": ["trainer"],
  "level": "advanced",
  "phone": "0900000002",
  "branch_id": ["69f20a180c4cd4cdf57684fe"],
  "created_at": "2026-05-26T08:00:00Z",
  "updated_at": "2026-05-26T08:00:00Z"
}
```

### Create employee

`POST /api/v1/employees`

Request:

```json
{
  "employee_id": "EMP001",
  "full_name": "Tran Van Trainer",
  "email": "trainer@gym.test",
  "password": "strong-password-123",
  "role": ["trainer"],
  "level": "advanced",
  "phone": "0900000002",
  "branch_id": ["69f20a180c4cd4cdf57684fe"],
  "status": "active"
}
```

Response `201`:

```json
{
  "message": "employee created successfully",
  "data": {
    "id": "69f20c000c4cd4cdf5768500",
    "employee_id": "EMP001",
    "full_name": "Tran Van Trainer",
    "email": "trainer@gym.test",
    "status": "active",
    "role": ["trainer"],
    "level": "advanced",
    "phone": "0900000002",
    "branch_id": ["69f20a180c4cd4cdf57684fe"],
    "created_at": "2026-05-26T08:00:00Z",
    "updated_at": "2026-05-26T08:00:00Z"
  }
}
```

Status codes:

- `201`: tạo employee thành công.
- `400`: invalid JSON, ObjectID, thiếu field bắt buộc, role/status/level/password không hợp lệ.
- `409`: normalized email hoặc employee ID đã tồn tại.
- `500`: storage/internal error.

### List employees

`GET /api/v1/employees`

Query filter optional:

| Field | Type | Rule |
|---|---|---|
| `role` | string | Nếu gửi thì phải là `admin`, `manager`, `trainer`, hoặc `receptionist`. |
| `status` | string | Nếu gửi thì phải là `active` hoặc `inactive`. |
| `branch_id` | ObjectID string | Nếu gửi thì phải là ObjectID hợp lệ. |

Response `200`:

```json
{
  "message": "employees fetched successfully",
  "data": []
}
```

List được phép trả `data` rỗng. Sort theo `created_at` giảm dần, hoặc `_id` giảm dần nếu thiếu
timestamp cũ.

### Get employee

`GET /api/v1/employees/:id`

Status codes:

- `200`: tìm thấy employee.
- `400`: employee ID không hợp lệ.
- `404`: không tìm thấy employee.
- `500`: storage/internal error.

### Update employee

`PATCH /api/v1/employees/:id`

Request hỗ trợ partial update. Field nào không gửi thì giữ nguyên.

```json
{
  "employee_id": "EMP001",
  "full_name": "Tran Van Trainer Updated",
  "email": "trainer.updated@gym.test",
  "role": ["trainer", "manager"],
  "level": "professional",
  "phone": "0900000003",
  "branch_id": ["69f20a180c4cd4cdf57684fe"],
  "status": "inactive"
}
```

Response `200`:

```json
{
  "message": "employee updated successfully",
  "data": {
    "id": "69f20c000c4cd4cdf5768500",
    "employee_id": "EMP001",
    "full_name": "Tran Van Trainer Updated",
    "email": "trainer.updated@gym.test",
    "status": "inactive",
    "role": ["trainer", "manager"],
    "level": "professional",
    "phone": "0900000003",
    "branch_id": ["69f20a180c4cd4cdf57684fe"],
    "created_at": "2026-05-26T08:00:00Z",
    "updated_at": "2026-05-26T09:00:00Z"
  }
}
```

Status codes:

- `200`: cập nhật thành công.
- `400`: invalid JSON, ObjectID, role/status/level không hợp lệ, hoặc không gửi field mutable nào.
- `404`: không tìm thấy employee.
- `409`: trùng normalized email/employee ID hoặc thay đổi tự khóa tài khoản.
- `500`: storage/internal error.

### Reset employee password

`PATCH /api/v1/employees/:id/password`

Request:

```json
{
  "password": "new-strong-password-123"
}
```

Response `200`:

```json
{
  "message": "employee password updated successfully"
}
```

Status codes:

- `200`: thay password hash và revoke refresh token active của employee đó.
- `400`: invalid body hoặc password không đạt policy.
- `404`: không tìm thấy employee.
- `500`: storage/internal error.

## Business rules

- Employee management là admin-only trong cycle này.
- Role hợp lệ:
  - `admin`
  - `manager`
  - `trainer`
  - `receptionist`
- Role input được lowercase, trim, deduplicate, và không được rỗng.
- Status hợp lệ:
  - `active`
  - `inactive`
- Create default `status` thành `active` nếu client không gửi.
- Create bắt buộc có `employee_id`, `full_name`, `email`, `password`, và ít nhất một role.
- Email dùng normalization hiện tại của auth trước khi persist và lookup.
- Chỉ lưu bcrypt password hash; không nhận hoặc trả `password_hash`.
- Password policy cho cycle này:
  - tối thiểu 8 ký tự
  - chỉ trim để check rỗng, còn bytes password gửi lên được giữ nguyên để hash
- `branch_id` optional, nhưng mọi ObjectID gửi lên phải reference branch đang tồn tại.
- `level` optional với staff không phải trainer. Nếu gửi thì phải là `basic`, `advanced`, hoặc
  `professional`.
- Nếu role có `trainer`, `level` bắt buộc và phải hợp lệ.
- Partial update phải giữ nguyên giá trị hiện tại với field không gửi.
- Update email phải tính lại `normalized_email`.
- Update status sang `inactive` làm access token hiện có fail ở request protected tiếp theo, vì auth
  đã reload employee state.
- Khi status đổi từ `active` sang `inactive`, revoke refresh token active của employee đó.
- Khi reset password, revoke refresh token active của employee đó.
- Chặn admin tự deactivate tài khoản của mình hoặc tự remove role `admin` của mình để giảm rủi ro
  tự khóa hệ thống. Rule "luôn còn ít nhất một active admin" đầy đủ có thể để cycle data-integrity
  nếu MVP rule này chưa đủ.

## Data và query changes

Không cần collection MongoDB mới.

Reuse collection `employees` và field hiện có:

| Field | Rule |
|---|---|
| `_id` | Service/repository tạo ObjectID mới khi create. |
| `employee_id` | Mã staff unique; conflict trả `409`. |
| `email` | Lưu dạng normalized/lowercase để đồng bộ với bootstrap/login hiện tại. |
| `normalized_email` | Login key do server tính, unique. |
| `password_hash` | Chỉ lưu bcrypt hash. |
| `status` | `active` hoặc `inactive`. |
| `role` | Danh sách role hợp lệ đã deduplicate. |
| `level` | Level course/trainer. |
| `branch_id` | Reference ObjectID tới branch hiện có. |
| `created_at` | Set khi create. |
| `updated_at` | Set khi create, update, và reset password. |

Repository cần thêm:

- `List(ctx, filter EmployeeListFilter) ([]models.Employee, error)`
- `UpdateByID(ctx, id string, update EmployeeUpdate) (*models.Employee, error)`
- `UpdatePasswordByID(ctx, id string, passwordHash string, now time.Time) error`
- Có thể thêm `ExistsByID`/`GetByIDs` cho branch repository, hoặc service validate bằng
  `BranchRepository.GetByID` hiện có.
- Thêm method refresh-token repository:
  - `RevokeActiveByEmployeeID(ctx, employeeID primitive.ObjectID, now time.Time) error`

Duplicate-key handling:

- Normalize Mongo duplicate key error từ employee create/update thành service-level conflict error
  để handler trả `409`.

## Layer plan

### Models

- Giữ `models.Employee` làm persistence shape.
- Thêm request/response DTO ở handler hoặc service; không dùng raw model cho request chứa password.
- Thêm list filter/update struct ở repository hoặc service nếu cần.

### Repository

- Mở rộng employee repository với list và partial update.
- Chỉ update field mutable; profile update không được overwrite password hash.
- Set `updated_at` khi profile update và password update.
- Return employee sau update để handler trả safe response mới nhất.
- Build list filter theo optional `role`, `status`, `branch_id`.
- Mở rộng refresh-token repository để revoke active token theo employee ID.

### Service

- Thêm `EmployeeService`.
- Service chịu trách nhiệm validation và normalization:
  - required fields
  - roles/status/level/password
  - email normalization
  - validate branch reference
  - map duplicate/conflict
  - self-lockout guard dựa trên authenticated actor ID từ handler context
- Hash password bằng bcrypt.
- Convert employee model sang safe response.
- Revoke refresh token khi reset password và khi active-to-inactive.

### Handlers và routes

- Thêm `internal/handlers/employee_handler.go`.
- Handler parse ObjectID và request body.
- Handler đọc authenticated employee ID từ `handlers.AuthEmployeeIDKey` để service check
  self-lockout.
- Error mapping:
  - invalid input -> `400`
  - not found -> `404`
  - duplicate/self-lockout conflict -> `409`
  - unexpected storage/internal -> `500`
- Wire routes trong `cmd/server/main.go`:

```go
employeeRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin))
employeeRoutes.POST("/employees", employeeHandler.Create)
employeeRoutes.GET("/employees", employeeHandler.List)
employeeRoutes.GET("/employees/:id", employeeHandler.GetByID)
employeeRoutes.PATCH("/employees/:id", employeeHandler.Update)
employeeRoutes.PATCH("/employees/:id/password", employeeHandler.UpdatePassword)
```

Route order note: kiểm tra route `/employees/:id/password` khi implement. Nếu Gin conflict với
wildcard theo method/path trong thực tế, đăng ký route password trước route wildcard.

## Docs plan

Cập nhật trong phase implement/complete:

- `docs/api_contract.md`
  - Thêm Employees collection table.
  - Thêm role guard row cho employee management.
  - Thêm endpoint details, request/response examples, status codes.
- `api_test.http`
  - Thêm sample create/list/get/update/reset password employee bằng admin token.
  - Thêm sample non-admin bị forbidden.
  - Thêm sample login bằng employee vừa tạo.
- `docs/code_reading_guide.md`
  - Thêm flow employee management sau khi implement.
- `docs/local_dev_guide.md`
  - Thêm manual employee seeding/verification nếu cần.
- `CHAT_CONTEXT/README.md`
  - Chỉ chuyển employee management từ planned sang implemented ở phase complete.
- `CHAT_CONTEXT/backend_skills/worklog.md`
  - Cập nhật phase status khi implement/review/test/complete.

## Test và verification plan

Minimum automated checks:

- `go build ./...`
- `go test ./...`

Focused service tests:

- Create employee success hash password và trả safe response.
- Thiếu required fields trả invalid input.
- Role/status/level không hợp lệ trả invalid input.
- Trainer không có level trả invalid input.
- Duplicate normalized email/employee ID map thành conflict.
- Branch reference validation reject branch thiếu/không hợp lệ.
- Update giữ nguyên field không gửi.
- Update email tính lại normalized email.
- Self-deactivation và self-admin-role removal trả conflict.
- Password reset đổi hash và revoke active refresh tokens.
- Deactivation revoke active refresh tokens và làm auth reject employee.

Focused handler/manual API checks:

- Admin create/list/get/update/reset password thành công.
- Non-admin authenticated employee nhận `403`.
- Missing token nhận `401`.
- Employee vừa tạo login được bằng initial password.
- Refresh token cũ fail sau password reset.
- Inactive employee không login được và không dùng được access token hiện có ở protected routes.
- Response không chứa `password_hash` hoặc `normalized_email`.

Direct Mongo checks:

- `employees.password_hash` tồn tại và nhìn giống bcrypt hash.
- `employees.normalized_email` là lowercase.
- Duplicate `employee_id`/`normalized_email` bị reject.
- `refresh_tokens.revoked_at` được set sau password reset/deactivation.

## Risks

- Password reset và refresh-token revocation không nằm trong MongoDB transaction; nếu xử lý thứ tự
  không cẩn thận có thể update password nhưng chưa revoke được refresh token cũ.
- Access token hiện có vẫn chỉ bị reject ở protected request tiếp theo, khi middleware reload
  employee state. Với inactive status, hành vi này chấp nhận được vì middleware đã reload state.
- Self-lockout protection không đảm bảo chắc chắn luôn còn một active admin khác. Rule last-active
  admin đầy đủ có thể cần indexed query và nên thuộc cycle data-integrity nếu cần.
- Session create hiện lưu `trainer_id` nhưng chưa validate trainer đó là active trainer. Cycle
  employee này chuẩn hóa staff records, còn trainer enforcement nên là task session/data-integrity
  riêng nếu không mở rộng scope.
- Unique sparse indexes đã có, nhưng legacy employee document thiếu normalized fields có thể cần
  cleanup trước khi strict management semantics đáng tin cậy hoàn toàn.
