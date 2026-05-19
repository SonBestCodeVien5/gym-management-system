# Cycle 04 — Auth/Login + Role Guard

## Status

- Status: planned
- Priority: medium
- Depends on: employee model/repository design
- Endpoints:
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/refresh`
  - `POST /api/v1/auth/logout`

## Goal

Thêm authentication cho employee/staff và middleware phân quyền theo role.

## API plan

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

Response:
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "employee": {
    "id": "ObjectID",
    "email": "admin@gym.test",
    "role": "admin"
  }
}
```

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

## Roles

- `admin`
- `manager`
- `trainer`
- `receptionist`

## Permission plan

Public:
- `/ping`
- `/api/v1/auth/login`
- `/api/v1/auth/refresh`

Admin/manager:
- course CRUD
- branch CRUD

Receptionist/manager/admin:
- member create/get/activate
- subscription create/refund/suspend/unsuspend/expire
- attendance checkin/report/makeup

Trainer/manager/admin:
- sessions create/list/get/enroll/checkin

## Data model plan

Employee:
- ID
- email or username
- password_hash
- full_name
- role
- branch_id optional
- status
- created_at
- updated_at

Refresh token:
- token_hash
- employee_id
- expires_at
- revoked_at optional
- created_at

## Security rules

- Never store plain password.
- Hash password with bcrypt.
- JWT secrets from env:
  - `JWT_ACCESS_SECRET`
  - `JWT_REFRESH_SECRET`
- Access token short TTL.
- Refresh token longer TTL.
- Store hash of refresh token if persisted.
- Never log token/password.

## Implementation plan

- Add employee repository if missing.
- Add auth service.
- Add auth handler.
- Add middleware:
  - `AuthRequired`
  - `RequireRoles`
- Wire routes in `main.go`.
- Protect selected routes.
- Add seed/admin bootstrap plan if needed.

## Docs/test plan

Update:
- `.env.example`
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

- Need decide seed admin strategy.
- Adding auth can break existing manual tests unless public/protected route plan clear.
- Role guard should be introduced carefully after core APIs are stable.