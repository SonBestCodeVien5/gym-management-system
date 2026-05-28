# Huong Dan Doc Code

Tai lieu nay chi giu ban do doc code ngan. Chi tiet endpoint nam trong
[api_contract.md](api_contract.md); cach chay va test local nam trong
[local_dev_guide.md](local_dev_guide.md).

## 1. Doc Tu Entry Point

Bat dau o `cmd/server/main.go`:

1. Xem `.env`, ket noi MongoDB va database name.
2. Xem `pkg/database.EnsureIndexes` duoc goi truoc khi khoi tao repository.
3. Xem repository -> service -> handler duoc khoi tao theo thu tu nao.
4. Xem route Gin de biet feature nao dang duoc wire that.

## 2. Lop Kien Truc

| Lop | Vai tro |
|---|---|
| `internal/handlers` | Parse HTTP input, goi service, map loi sang response |
| `internal/service` | Business rules va orchestration |
| `internal/repository` | MongoDB query/update |
| `internal/models` | Struct DB/JSON |
| `pkg/database` | MongoDB connect va startup index bootstrap |

Khi doc mot feature, di theo luong:

`route -> handler -> service -> repository -> model`

HTTP error response shape dung helper chung trong `internal/handlers/response.go`. Handler map
service error sang public `error.code`; service va repository khong biet HTTP response shape.

## 3. Luong Nen Doc Truoc

| Feature | File chinh |
|---|---|
| Member registration | `member_handler.go`, `member_service.go`, `member_repo.go`, `member.go` |
| Subscription create/activate | `subscription_handler.go`, `subscription_service.go`, `subscription_repo.go`, `subscription.go` |
| Attendance rules | `attendance_handler.go`, `attendance_service.go`, `attendance_repo.go`, `attendance.go` |
| Sessions | `session_handler.go`, `session_service.go`, `session_repo.go`, `session.go` |
| Branch nearby | `branch_handler.go`, `branch_service.go`, `branch_repo.go`, `branch.go` |
| Auth and role guard | `auth_handler.go`, `auth_middleware.go`, `auth_service.go`, `employee_repo.go`, `refresh_token_repo.go`, `employee.go`, `refresh_token.go` |
| Employee management | `employee_handler.go`, `employee_service.go`, `employee_repo.go`, `refresh_token_repo.go`, `employee.go` |
| Error response consistency | `response.go`, `auth_middleware.go`, cac `*_handler.go`, `api_contract.md` |

## 4. Rule Dang Co

- Member registration can unique `ccid`.
- Subscription creation validate member/course/branch references and snapshot course data.
- Offline payment activation goes through member activation and subscription confirmation.
- Refund, enrollment, and other double-submit risks need atomic repository behavior.
- Attendance handles weekly quota, reported-missed window, makeup reference, and remaining session
  effects.
- Nearby branch search depends on GeoJSON coordinates and MongoDB geo index behavior.
- Auth normalizes employee email, verifies bcrypt password hash, signs access/refresh tokens, stores
  refresh tokens as hashes, and reloads employee state during access-token validation.
- Role guard checks trusted roles from auth middleware; handlers must not trust client-sent role
  fields.
- Employee management is admin-only, hashes password in service, never returns password hash, and
  revokes active refresh tokens on password reset or deactivation.
- Error responses use `{"error":{"code":"...","message":"...","details":{}}}`. Success responses
  keep the current `message`/`data` shape.
- Startup creates MongoDB indexes centrally in `pkg/database.EnsureIndexes`; repositories should not
  hide index creation in constructors.
- Unique indexes enforce member CCID, branch code, employee email/ID, refresh-token hash, refund
  subscription, duplicate session check-in, and duplicate makeup reuse.

Doc rule trong service truoc khi sua handler hay repository.

## 5. Cach Debug Theo Luong

1. Tim route trong `cmd/server/main.go`.
2. Neu thay doi data integrity/query behavior, kiem tra `pkg/database/indexes.go`.
3. Kiem tra handler nhan body, param, query va parse date/ObjectID dung chua.
4. Kiem tra service tra domain error nao va business rule nao dang chan.
5. Kiem tra repository query field, filter status, atomic update, va MongoDB index lien quan.
6. Doi chieu `docs/api_contract.md`, `api_test.http`, va test hien co.

## 6. Thu Tu Goi Y

1. `cmd/server/main.go`
2. `pkg/database/mongodb.go` va `pkg/database/indexes.go`
3. member handler/service/repository/model
4. course va branch handler/service/repository/model
5. subscription handler/service/repository/model
6. attendance handler/service/repository/model
7. session handler/service/repository/model
8. auth va employee handler/middleware/service/repository/model
9. tests va API contract cho feature dang sua

## 7. Tai Lieu Lien Quan

- [api_contract.md](api_contract.md): endpoint current/planned va request/response quan trong.
- [local_dev_guide.md](local_dev_guide.md): run backend, test API, inspect MongoDB.
- [faq_why.md](faq_why.md): ly do kien truc va nghiep vu can giai thich.
- [../CHAT_CONTEXT/README.md](../CHAT_CONTEXT/README.md): snapshot de resume phien moi.
