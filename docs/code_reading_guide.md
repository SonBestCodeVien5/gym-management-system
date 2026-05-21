# Huong Dan Doc Code

Tai lieu nay chi giu ban do doc code ngan. Chi tiet endpoint nam trong
[api_contract.md](api_contract.md); cach chay va test local nam trong
[local_dev_guide.md](local_dev_guide.md).

## 1. Doc Tu Entry Point

Bat dau o `cmd/server/main.go`:

1. Xem `.env`, ket noi MongoDB va database name.
2. Xem repository -> service -> handler duoc khoi tao theo thu tu nao.
3. Xem route Gin de biet feature nao dang duoc wire that.

## 2. Lop Kien Truc

| Lop | Vai tro |
|---|---|
| `internal/handlers` | Parse HTTP input, goi service, map loi sang response |
| `internal/service` | Business rules va orchestration |
| `internal/repository` | MongoDB query/update |
| `internal/models` | Struct DB/JSON |

Khi doc mot feature, di theo luong:

`route -> handler -> service -> repository -> model`

## 3. Luong Nen Doc Truoc

| Feature | File chinh |
|---|---|
| Member registration | `member_handler.go`, `member_service.go`, `member_repo.go`, `member.go` |
| Subscription create/activate | `subscription_handler.go`, `subscription_service.go`, `subscription_repo.go`, `subscription.go` |
| Attendance rules | `attendance_handler.go`, `attendance_service.go`, `attendance_repo.go`, `attendance.go` |
| Sessions | `session_handler.go`, `session_service.go`, `session_repo.go`, `session.go` |
| Branch nearby | `branch_handler.go`, `branch_service.go`, `branch_repo.go`, `branch.go` |

## 4. Rule Dang Co

- Member registration can unique `ccid`.
- Subscription creation validate member/course/branch references and snapshot course data.
- Offline payment activation goes through member activation and subscription confirmation.
- Refund, enrollment, and other double-submit risks need atomic repository behavior.
- Attendance handles weekly quota, reported-missed window, makeup reference, and remaining session
  effects.
- Nearby branch search depends on GeoJSON coordinates and MongoDB geo index behavior.

Doc rule trong service truoc khi sua handler hay repository.

## 5. Cach Debug Theo Luong

1. Tim route trong `cmd/server/main.go`.
2. Kiem tra handler nhan body, param, query va parse date/ObjectID dung chua.
3. Kiem tra service tra domain error nao va business rule nao dang chan.
4. Kiem tra repository query field, filter status, atomic update, va MongoDB index lien quan.
5. Doi chieu `docs/api_contract.md`, `api_test.http`, va test hien co.

## 6. Thu Tu Goi Y

1. `cmd/server/main.go`
2. member handler/service/repository/model
3. course va branch handler/service/repository/model
4. subscription handler/service/repository/model
5. attendance handler/service/repository/model
6. session handler/service/repository/model
7. tests va API contract cho feature dang sua

## 7. Tai Lieu Lien Quan

- [api_contract.md](api_contract.md): endpoint current/planned va request/response quan trong.
- [local_dev_guide.md](local_dev_guide.md): run backend, test API, inspect MongoDB.
- [faq_why.md](faq_why.md): ly do kien truc va nghiep vu can giai thich.
- [../CHAT_CONTEXT/README.md](../CHAT_CONTEXT/README.md): snapshot de resume phien moi.
