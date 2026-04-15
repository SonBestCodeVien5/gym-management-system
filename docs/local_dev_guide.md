# Huong Dan Local Dev - Phase 2

## 1. Muc tieu tai lieu
Tai lieu nay tong hop toan bo quy trinh da lam de ban co the:
- Chay backend Go + MongoDB local
- Test API (`/ping`, `/api/v1/registration`)
- Xem record trong MongoDB nhanh
- Ket noi MongoDB Compass dung cau hinh
- Xu ly cac loi thuong gap

## 2. Kien truc dang dung (ban rut gon)
Luong xu ly hien tai:
1. `cmd/server/main.go` khoi dong app, doc `.env`, ket noi MongoDB, mo route.
2. `pkg/database/mongodb.go` quan ly ket noi DB (`ConnectMongoDB`).
3. `internal/repository/member_repo.go` query vao collection `members`.
4. `internal/service/member_service.go` xu ly logic dang ky member.
5. `internal/handlers/member_handler.go` nhan request HTTP va tra JSON.

Luot request dang ky:
HTTP Request -> Handler -> Service -> Repository -> MongoDB -> JSON Response.

## 3. Chay he thong local
### 3.1 Chay MongoDB bang Docker
```bash
docker compose up -d
```

Kiem tra container:
```bash
docker ps --filter name=gym_mongodb --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
```

### 3.2 Kiem tra file `.env`
Noi dung can dung:
```env
MONGODB_URI=mongodb://admin:password123@localhost:27017/?authSource=admin
PORT=8080
```

### 3.3 Build va run backend
```bash
go build ./... && echo "BUILD OK"
go run cmd/server/main.go
```

Neu thanh cong, log se co:
- `Connected to MongoDB successfully`
- `Listening and serving HTTP on :8080`

## 4. Test API
### 4.1 File `api_test.http` dung format dung
```http
# Health check
GET http://localhost:8080/ping

###

# Member registration
POST http://localhost:8080/api/v1/registration
Content-Type: application/json

{
  "ccid": "012345678901",
  "full_name": "Nguyen Van A",
  "email": "a@example.com",
  "phone": "0900000000",
  "gender": "male",
  "level": "basic"
}
```

Luu y:
- Bat buoc co `###` de REST Client nhan la 2 request rieng.
- Neu khong co `###` thi co the khong hien `Send Request` cho request thu 2.

### 4.2 Test nhanh bang curl
```bash
curl -s http://localhost:8080/ping
```

```bash
curl -s -X POST http://localhost:8080/api/v1/registration \
  -H 'Content-Type: application/json' \
  -d '{"ccid":"012345678901","full_name":"Nguyen Van A","email":"a@example.com","phone":"0900000000","gender":"male","level":"basic"}'
```

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

## 8. Checkpoint da hoan thanh
1. Ket noi MongoDB thanh cong.
2. Route `GET /ping` hoat dong.
3. Route `POST /api/v1/registration` hoat dong.
4. Member duoc insert vao `gym_management.members`.
5. Da giam coupling nhe: service khong check `mongo.ErrNoDocuments` truc tiep, su dung loi trung lap tu repository.
