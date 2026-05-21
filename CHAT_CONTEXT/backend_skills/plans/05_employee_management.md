# Cycle 05 - Employee Management

## Status

- Status: planned
- Priority: medium
- Depends on: auth/login + role guard

## Goal

Thêm API quản trị employee để admin tạo và quản lý staff account sau khi auth bootstrap đã có.

## API plan

| Method | Endpoint | Purpose |
|---|---|---|
| `POST` | `/api/v1/employees` | Tạo employee/staff account. |
| `GET` | `/api/v1/employees` | Danh sách employee. |
| `GET` | `/api/v1/employees/:id` | Xem chi tiết employee. |
| `PATCH` | `/api/v1/employees/:id` | Cập nhật profile, role, branch, status. |
| `PATCH` | `/api/v1/employees/:id/password` | Đặt lại hoặc đổi password theo policy đã chốt. |

## Business rules

- Scope đầu tiên là admin-only; manager branch scope để feature sau nếu cần.
- Không trả `password_hash` ra response.
- Validate employee role và branch references trước khi persist.
- Không tin client input cho auth-computed fields hoặc password hash.
- Employee inactive không được login sau khi status flow được nối với auth.

## Data plan

- Reuse employee auth fields từ cycle 04.
- Thêm repository queries cho create/list/get/update/password update nếu chưa có.
- Dùng unique email/employee identifier rules đã chốt cho auth và index hardening.

## Layer plan

- Model: chỉ mở rộng employee fields nếu API management cần.
- Repository: employee CRUD/password update queries.
- Service: validation, role/status/password rules.
- Handler/routes: admin-protected employee endpoints.

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

- Employee lifecycle và auth rules phải đồng bộ để inactive/password reset không để token flow mơ hồ.
- Role/branch assignment dễ mở rộng scope sang branch authorization nếu không giữ admin-only ban đầu.
- Password setup/reset policy cần chốt rõ trước khi endpoint password được implement.
