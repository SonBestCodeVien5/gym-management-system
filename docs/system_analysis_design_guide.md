# Huong Dan Viet Bao Cao: Mo Hinh Phat Trien & Phan Tich Thiet Ke He Thong

## 1. Muc tieu tai lieu
Tai lieu nay tong hop cac noi dung can de viet phan:
1. Mo hinh phat trien phan mem (SDLC)
2. Phan tich nghiep vu
3. Thiet ke he thong
4. Truy vet yeu cau va kiem thu

Tai lieu wiki tham chieu nhanh:
- [docs/faq_why.md](docs/faq_why.md)
- [docs/development_journal.md](docs/development_journal.md)

## 2. De xuat mo hinh phat trien cho du an nay

### 2.1 Mo hinh de xuat: Hybrid (Iterative-Incremental + Scrum-lite)
Ly do phu hop:
1. Du an co nghiep vu thay doi trong qua trinh hoc va lam.
2. Co the chot tung luong chuc nang nho de demo som (registration truoc, payment sau).
3. De quan ly rui ro ky thuat (DB, API, coupling) theo tung dot.

### 2.2 Cach mo ta trong bao cao
Ban co the mo ta theo 3 giai doan:
1. **Foundation Sprint**:
   - Khoi tao kien truc (handler/service/repository)
   - Setup Docker + MongoDB + ket noi local
2. **Core Business Sprint**:
   - Hoan thanh registration flow va validate data
   - Hoan thien test API co ban
3. **Expansion Sprint**:
   - Payment, suspension/resume, attendance report
   - Hardening: index, test, response convention
4. **Subscription Sprint**:
   - Tao subscription dua tren member/course/branch
   - Kiem tra tham chieu va parse ngay gio RFC3339
   - Mo rong route API cho subscription

### 2.3 Definition of Done (DoD) de dua vao bao cao
Moi user story/use case duoc xem la xong khi:
1. Build pass (`go build ./...`)
2. API test pass (REST Client/curl)
3. Co validation + xu ly loi ro rang
4. Co cap nhat tai lieu (journal/faq)

## 3. Huong dan phan tich nghiep vu

### 3.1 Actors va pham vi
Actors chinh:
1. Member
2. Receptionist
3. Trainer
4. Manager

Pham vi he thong hien tai:
1. Dang ky hoc vien
2. Quan ly thong tin hoc vien
3. Tao subscription
4. Nen tang de mo rong thanh toan, bao luu, diem danh

### 3.2 Danh sach business rules uu tien cao
1. `ccid` phai unique
2. Khong cho dang ky neu input thieu truong bat buoc
3. Tu dong set field he thong khi tao member:
   - `is_registered = false`
   - `is_suspended = false`
   - `total_sessions_attended = 0`
   - `created_at`, `updated_at`

### 3.3 Functional vs Non-functional Requirements
Functional (FR):
1. Dang ky member
2. Lay thong tin member
3. Tao subscription
4. (Mo rong) payment/suspension/attendance

Non-functional (NFR):
1. Tinh dung dan du lieu: unique `ccid`
2. Kha nang bao tri: tach tang ro rang
3. Kha nang mo rong: bo sung use case moi khong pha vo flow cu
4. Tinh quan sat: co logging/co tai lieu debug

## 4. Huong dan thiet ke he thong

### 4.1 Kien truc logic (Layered/Clean-lite)
Luong xu ly:
`Handler -> Service -> Repository -> MongoDB`

Vai tro:
1. Handler: giao tiep HTTP, parse request, map response
2. Service: chua business logic
3. Repository: truy cap du lieu
4. Model: schema du lieu dung chung

### 4.2 Thiet ke du lieu
Collections da co trong phase hien tai:
1. members
2. branches
3. courses
4. subscriptions
5. attendance
6. employees

Luu y quan trong khi viet bao cao:
1. Tien te dung `int64` de tranh sai so float
2. Embedded vs reference can duoc giai trinh ro
3. Cac field optional dung pointer + `omitempty`

### 4.3 Thiet ke API
Nguyen tac trinh bay:
1. Endpoint theo resource
2. Ma HTTP ro nghia
3. Validation o handler, rule o service
4. Response co cau truc thong nhat

Mau endpoint hien co:
1. `GET /ping`
2. `POST /api/v1/registration`
3. `GET /api/v1/members/:id`
4. `POST /api/v1/subscriptions`
5. `GET /api/v1/subscriptions/:id`

## 5. Truy vet yeu cau (Requirement Traceability)

### 5.1 Mau bang truy vet de chen vao bao cao
| Req ID | Use Case | API | Du lieu lien quan | Test Case |
|---|---|---|---|---|
| FR-01 | Dang ky member | POST /api/v1/registration | members.ccid, full_name | TC-REG-001 |
| FR-02 | Kiem tra trung CCCD | POST /api/v1/registration | members.ccid | TC-REG-002 |
| FR-03 | Tao subscription | POST /api/v1/subscriptions | subscriptions.member_id, course_id, home_branch_id | TC-SUB-001 |

### 5.2 Loi ich
1. Hoi dong de doi chieu yeu cau voi code de dang
2. Nhom de kiem soat pham vi trien khai
3. De viet phan kiem thu ket qua

## 6. Chien luoc kiem thu de dua vao bao cao

### 6.1 Test theo tang
1. Unit test service: rule nghiep vu
2. Integration test: API + DB
3. Manual test: REST Client + Mongo query

### 6.2 Test theo nghiep vu
Vi du cho registration:
1. Input hop le -> 201
2. Thieu `ccid`/`full_name` -> 400
3. Trung `ccid` -> 409

Vi du cho subscription:
1. Input hop le + tham chieu ton tai -> 201
2. Sai ID / ngay gio -> 400
3. Member/course/branch khong ton tai -> 404

## 7. Rui ro va giam thieu
1. Coupling storage qua manh -> giam bang error abstraction
2. Du lieu trung do race -> them unique index + xu ly duplicate key
3. Sai cau hinh local -> chuan hoa `.env` va local guide
4. Vuot scope -> chia sprint, moi sprint co DoD

## 8. Ghi chu thiet ke subscription
1. Subscription luu tham chieu qua ID, khong nhung toan bo member/course/branch.
2. Price va session total lay tu Course de tranh sai du lieu tu client.
3. Ngay gio request dung RFC3339 de nhat quan giua Go, JSON, MongoDB.

## 9. Doan mo ta mau de chen thang vao bao cao
"Du an ap dung mo hinh phat trien Iterative-Incremental ket hop Scrum-lite de giam rui ro ky thuat va cho phep phan hoi som. Kien truc he thong duoc tach thanh cac tang Handler, Service, Repository va Database nham dam bao tinh bao tri, de kiem thu va de mo rong use case theo tung sprint. Qua trinh phan tich va thiet ke duoc truy vet thong qua ma yeu cau (FR), use case, endpoint API, thanh phan du lieu va test case tuong ung, giup dam bao tinh nhat quan giua tai lieu va hien thuc." 

---

Cap nhat lan cuoi: 2026-04-28