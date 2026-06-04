# Huong Dan Local Dev Va Demo

## 1. Muc tieu tai lieu
Tai lieu nay tong hop toan bo quy trinh da lam de ban co the:
- Chay full stack bang Docker Compose
- Chay backend Go + MongoDB local
- Chay frontend React/Vite local
- Nap demo data bang seed command
- Test API (`/ping`, auth, members, subscriptions, attendance, sessions)
- Xem record trong MongoDB nhanh
- Ket noi MongoDB Compass dung cau hinh
- Xu ly cac loi thuong gap

Tai lieu API contract hien tai: [api_contract.md](api_contract.md)

## 2. Kien truc dang dung (ban rut gon)
Luong xu ly hien tai:
1. `cmd/server/main.go` khoi dong app, doc `.env`, ket noi MongoDB, chay index bootstrap.
2. `pkg/database/mongodb.go` quan ly ket noi DB (`ConnectMongoDB`).
3. `internal/app/router.go` dung chung de khoi tao repository/service/handler va dang ky route.
4. `internal/repository/member_repo.go` query vao collection `members`.
5. `internal/service/member_service.go` xu ly logic dang ky member.
6. `internal/handlers/member_handler.go` nhan request HTTP va tra JSON.
7. `internal/repository/course_repo.go`, `branch_repo.go` query theo `_id` cho subscription.
8. `internal/repository/subscription_repo.go`, `internal/service/subscription_service.go`, `internal/handlers/subscription_handler.go` phuc vu luong subscription.
9. `internal/handlers/course_handler.go`, `branch_handler.go` xu ly CRUD course/branch.
10. `internal/handlers/attendance_handler.go` xu ly check-in, report missed, makeup va history.
11. `internal/handlers/session_handler.go` xu ly session create/list/get/enroll/check-in.
12. `internal/handlers/auth_handler.go`, `auth_middleware.go`, `internal/service/auth_service.go`
    xu ly login, current employee, refresh, logout, access token va role guard.
13. `internal/app/cors.go` xu ly CORS allow-list va preflight cho browser FE local.

Luot request dang ky:
HTTP Request -> Handler -> Service -> Repository -> MongoDB -> JSON Response.

Luot request subscription:
HTTP Request -> Subscription Handler -> Subscription Service -> Member/Course/Branch Repo + Subscription Repo -> MongoDB -> JSON Response.

## 3. Chay he thong local
### 3.1 Chay full stack bang Docker
```bash
docker compose up -d --build
```

Neu Docker bao thieu plugin `docker-buildx`, build bang legacy builder truoc:

```bash
DOCKER_BUILDKIT=0 docker compose build
docker compose up -d
```

Nap demo data:
```bash
docker compose --profile seed run --rm seed
```

Neu may thieu `docker-buildx`, dung cung fallback cho seed profile:
```bash
DOCKER_BUILDKIT=0 docker compose --profile seed run --rm seed
```

Mo ung dung:
- Frontend: `http://localhost:5173`
- API: `http://localhost:8080/ping`
- MongoDB: `localhost:27017`

Tai khoan demo:

| Role | Email | Password |
|---|---|---|
| Admin | `admin@gym.test` | `demo123456` |
| Manager | `manager@gym.test` | `demo123456` |
| Receptionist | `receptionist@gym.test` | `demo123456` |
| Trainer | `trainer@gym.test` | `demo123456` |

Reset Docker database chi khi that su muon xoa data local:
```bash
docker compose down -v
```

Neu MongoDB bao loi featureCompatibilityVersion do volume cu duoc tao boi ban MongoDB moi hon,
hay dung lenh reset tren chi khi ban dong y xoa database Docker local.

### 3.2 Chay rieng MongoDB bang Docker de dev backend
```bash
docker compose up -d mongodb
```

Kiem tra container:
```bash
docker ps --filter name=gym_mongodb --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
```

### 3.3 Kiem tra file `.env`
Noi dung can dung:
```env
MONGODB_URI=mongodb://admin:password123@127.0.0.1:27017/?authSource=admin&directConnection=true
DB_NAME=gym_management
PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
JWT_ACCESS_SECRET=replace-with-a-long-random-access-secret
JWT_REFRESH_SECRET=replace-with-a-long-random-refresh-secret
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
BOOTSTRAP_ADMIN_EMPLOYEE_ID=ADMIN001
BOOTSTRAP_ADMIN_FULL_NAME=Gym Admin
BOOTSTRAP_ADMIN_EMAIL=admin@gym.test
BOOTSTRAP_ADMIN_PASSWORD=demo123456
```

`JWT_ACCESS_SECRET` va `JWT_REFRESH_SECRET` bat buoc phai co. Neu thieu, server se dung loi khi
khoi tao auth service. `BOOTSTRAP_ADMIN_EMAIL` + `BOOTSTRAP_ADMIN_PASSWORD` tao admin dau tien neu
email chua ton tai.

`DB_NAME` mac dinh la `gym_management` neu khong khai bao. Nen khai bao ro bien nay khi chay seed,
Docker, hoac khi muon tach database demo/test/local.

`CORS_ALLOWED_ORIGINS` dung cho browser FE local. Neu chay Vite o `localhost:5173`, giu gia tri mau
o tren. Neu bien nay rong, backend khong tra CORS header.

### 3.4 Build va run backend
```bash
go build ./... && echo "BUILD OK"
go run ./cmd/server
```

Neu thanh cong, log se co:
- `Connected to MongoDB successfully`
- `MongoDB indexes ensured successfully`
- `Listening and serving HTTP on :8080`

### 3.5 Nap demo data local
Seed command doc `.env`, tao index, va upsert demo data theo ID co dinh nen co the chay lai nhieu
lan ma khong tao ban ghi trung:

```bash
go run ./cmd/seed
```

Seed tao demo employees, branches, courses, members, subscriptions, attendances, sessions va refund
de dashboard/frontend co du lieu mau.

### 3.6 Chay frontend local
```bash
cp frontend/.env.example frontend/.env
npm --prefix frontend install
npm --prefix frontend run dev
```

`frontend/.env` dung:

```env
VITE_API_BASE_URL=http://localhost:8080
```

### 3.7 Route hien co
- `GET /ping`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `POST /api/v1/employees`
- `GET /api/v1/employees`
- `GET /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id/password`
- `GET /api/v1/dashboard/summary`
- `GET /api/v1/dashboard/revenue`
- `GET /api/v1/dashboard/plans`
- `GET /api/v1/dashboard/members/recent`
- `GET /api/v1/dashboard/sessions/today`
- `POST /api/v1/members`
- `GET /api/v1/members/:id`
- `GET /api/v1/members/:id/subscriptions`
- `PATCH /api/v1/members/:id/activate`
- `POST /api/v1/courses`
- `GET /api/v1/courses`
- `GET /api/v1/courses/:id`
- `PATCH /api/v1/courses/:id`
- `DELETE /api/v1/courses/:id`
- `POST /api/v1/branches`
- `GET /api/v1/branches`
- `GET /api/v1/branches/nearby`
- `GET /api/v1/branches/:id`
- `PATCH /api/v1/branches/:id`
- `DELETE /api/v1/branches/:id`
- `POST /api/v1/subscriptions`
- `POST /api/v1/subscriptions/:id/refund`
- `GET /api/v1/subscriptions/:id`
- `PATCH /api/v1/subscriptions/:id/suspend`
- `PATCH /api/v1/subscriptions/:id/unsuspend`
- `PATCH /api/v1/subscriptions/:id/expire`
- `POST /api/v1/attendance/checkin`
- `POST /api/v1/attendance/report`
- `POST /api/v1/attendance/makeup`
- `GET /api/v1/subscriptions/:id/attendance`
- `POST /api/v1/sessions`
- `GET /api/v1/sessions`
- `GET /api/v1/sessions/:id`
- `POST /api/v1/sessions/:id/enroll`
- `POST /api/v1/sessions/:id/checkin`

Luu y auth:
- Public: `/ping`, `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/logout`.
- Protected current-user route: `/api/v1/auth/me`.
- Cac route protected con lai can header `Authorization: Bearer <access_token>`.
- Role guard kiem quyen theo ma tran trong [api_contract.md](api_contract.md).

## 4. Test API
### 4.1 File [api_test.http](../api_test.http) dung format dung
```http
# Health check
GET http://localhost:8080/ping

###

# Login
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@gym.test",
  "password": "demo123456"
}
```

Luu y:
- Bat buoc co `###` de REST Client nhan la 2 request rieng.
- Neu khong co `###` thi co the khong hien `Send Request` cho request thu 2.
- Sau khi login, copy `access_token` vao bien `@access_token` trong `api_test.http`.
- Goi `GET /api/v1/auth/me` de kiem tra token va lay lai employee hien tai cho FE.

### 4.2 Test nhanh bang curl
```bash
curl -s http://localhost:8080/ping
```

Login:
```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@gym.test","password":"demo123456"}'
```

Gan token de goi route protected:
```bash
ACCESS_TOKEN='PASTE_ACCESS_TOKEN_HERE'
```

```bash
curl -s -X POST http://localhost:8080/api/v1/members \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"ccid":"012345678901","full_name":"Nguyen Van A","email":"a@example.com","phone":"0900000000","gender":"male","level":"basic"}'
```

Subscription test mau (RFC3339):
```bash
curl -s -X POST http://localhost:8080/api/v1/subscriptions \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"member_id":"PUT_MEMBER_OBJECT_ID","course_id":"PUT_COURSE_OBJECT_ID","home_branch_id":"PUT_BRANCH_OBJECT_ID","start_date":"2026-04-28T10:00:00Z","end_date":"2026-05-28T10:00:00Z","session_per_week":3}'
```

Activate member (offline payment confirm):
```bash
curl -s -X PATCH http://localhost:8080/api/v1/members/PUT_MEMBER_OBJECT_ID/activate \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"subscription_id":"PUT_SUBSCRIPTION_OBJECT_ID"}'
```

Attendance check-in (attended/makeup):
```bash
curl -s -X POST http://localhost:8080/api/v1/attendance/checkin \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"subscription_id":"PUT_SUBSCRIPTION_OBJECT_ID","branch_id":"PUT_BRANCH_OBJECT_ID","date":"2026-05-10T08:00:00Z","status":"attended"}'
```

### 4.3 Chay automated tests
Unit tests khong can server dang chay:
```bash
go test ./...
```

Integration tests dung `httptest` + MongoDB that qua `internal/app` va `internal/testutil`.
Neu MongoDB local khong reachable, cac integration tests se skip de `go test ./...` van chay duoc
unit test.

Chay MongoDB roi chay rieng integration tests:
```bash
docker compose up -d mongodb
GYM_TEST_MONGODB_URI='mongodb://admin:password123@localhost:27017/?authSource=admin&directConnection=true' \
  go test ./internal/integration -count=1
```

Luu y:
- Test DB co ten dang `gym_test_<id>` va duoc drop trong cleanup.
- Integration tests goi `pkg/database.EnsureIndexes` tren DB test truoc khi seed fixture.
- Khong chay integration tests tren DB dev `gym_management`.

## 5. Xem record trong MongoDB
### 5.1 Cach 1 - one-shot, in ket qua ngay
```bash
docker exec gym_mongodb mongosh "mongodb://admin:password123@localhost:27017/admin" --quiet --eval 'db.getSiblingDB("gym_management").members.find({}, {_id:0, ccid:1, full_name:1, email:1}).limit(20).forEach(d => print(d.ccid + " | " + d.full_name + " | " + d.email));'
```

### 5.2 Cach 2 - vao interactive shell roi query tay
```bash
docker exec -it gym_mongodb mongosh "mongodb://admin:password123@localhost:27017/admin"
```

Trong shell:
```javascript
show dbs
use gym_management
show collections
db.members.find().pretty()
db.employees.find({}, {password_hash:0}).pretty()
db.refresh_tokens.find({}, {token_hash:1, employee_id:1, expires_at:1, revoked_at:1}).pretty()
exit
```

Luu y quan trong:
- Lenh khong co `--eval` chi mo shell, khong tu in du lieu.
- Muon in ngay thi phai dung `--eval`.

## 6. Ket noi MongoDB Compass
### 6.1 Connection string de dung ngay
```text
mongodb://admin:password123@localhost:27017/?authSource=admin&directConnection=true
```

### 6.2 Neu nhap tay trong UI Compass
1. Hostname: `localhost`
2. Port: `27017`
3. Auth method: Username/Password
4. Username: `admin`
5. Password: `password123`
6. Authentication DB: `admin`
7. TLS/SSL: Off

### 6.3 Luu y ve MongoDB indexes
- Startup goi `pkg/database.EnsureIndexes` de tao index tap trung cho cac collection hien co.
- Unique indexes quan trong: `members.ccid`, `branches.branch_code`, `refunds.subscription_id`,
  `employees.normalized_email`, `employees.employee_id`, va `refresh_tokens.token_hash`.
- Attendance co partial unique indexes de chan duplicate session check-in va duplicate makeup reuse.
- Refresh token co TTL index tren `expires_at`; MongoDB xoa document het han theo co che eventual.
- Neu DB local da co duplicate data, server co the dung khi tao unique index. Khi do can cleanup data
  trung trong collection bi log bao loi, khong drop/rewrite data tu dong.

## 7. Loi thuong gap va cach xu ly
### 7.1 `Authentication failed`
Nguyen nhan:
- Sai password (nham `adminpassword123` voi `password123`).

Cach sua:
- Dong bo dung password trong `docker-compose.yml` va `.env`.

### 7.2 Build pass nhung khong thay output
- `go build ./...` chi compile, khong chay app.
- Khong co loi thi terminal thuong im lang.

### 7.3 `Send Request` khong hien trong file `.http`
- Chua cai extension REST Client.
- Thieu `###` de tach request.
- File khong dung duoi `.http`.

### 7.4 `go run` exit code 1
- Thuong do DB auth sai, port dang ban, hoac env thieu.
- Kiem tra lai `.env`, `docker ps`, va log terminal.
- Neu log bao `Failed to initialize auth service`, kiem tra `JWT_ACCESS_SECRET` va
  `JWT_REFRESH_SECRET`.

### 7.5 Route business tra `401`
- Thieu header `Authorization: Bearer <access_token>`.
- Token het han, sai secret, malformed, hoac employee da inactive.
- Login lai bang `/api/v1/auth/login` hoac refresh token bang `/api/v1/auth/refresh`.

### 7.6 Route business tra `403`
- Token hop le nhung role khong du quyen cho route.
- Xem ma tran role trong [api_contract.md](api_contract.md).

### 7.7 `invalid start_date format` khi test subscription
- Do input date khong dung RFC3339.
- Dinh dang dung la `2026-04-28T10:00:00Z` hoac co timezone ro rang.

## 8. Checkpoint da hoan thanh
1. Ket noi MongoDB thanh cong.
2. Route `GET /ping` hoat dong.
3. Route `POST /api/v1/members` hoat dong.
4. Member duoc insert vao `gym_management.members`.
5. Da giam coupling nhe: service khong check `mongo.ErrNoDocuments` truc tiep, su dung loi trung lap tu repository.
6. Da co central MongoDB index bootstrap cho unique/query/TTL indexes.
7. Da co route subscription va parser RFC3339.
8. Da co auth login/refresh/logout, bootstrap admin, refresh-token hash storage, access-token
   middleware, va role guard cho route business.
9. Da co employee management admin-only de tao/list/get/update/reset password staff account.
10. Da co `internal/app` de dung chung route/dependency wiring cho server va integration tests.
11. Da co integration tests voi MongoDB test DB rieng cho auth, role guard, subscription,
    duplicate conflict, attendance makeup, va branch nearby.
12. Da co dashboard/report aggregate APIs cho admin/manager.
13. Da co frontend React/Vite staff portal ket noi live API.
14. Da co full-stack Docker Compose va seed demo data command.
