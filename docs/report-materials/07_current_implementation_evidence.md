# Current Implementation Evidence for Formal Report

Tai lieu nay gom bang chung hien trang de viet bao cao chinh thuc. Muc tieu la tach ro:

- **Implemented**: da co trong code va da duoc test.
- **Partial/MVP**: da co ban chay duoc nhung con gioi han.
- **Planned/Future work**: chua co trong backend hien tai.

Nguon doi chieu chinh:

- `docs/api_contract.md`
- `api_test.http`
- `README.md`
- `docs/local_dev_guide.md`
- source code trong `cmd/server`, `internal/handlers`, `internal/service`, `internal/repository`,
  `internal/models`

---

## 1. Executive Summary

He thong hien tai la backend quan ly gym da chi nhanh duoc xay dung bang Go, Gin va MongoDB. Kien
truc duoc tach theo luong `Handler -> Service -> Repository -> MongoDB`, giup phan tach HTTP
parsing, business rules va thao tac du lieu.

Phien ban hien tai da hoan thanh cac nhom chuc nang backend chinh:

- Quan ly hoc vien.
- Quan ly goi tap, chi nhanh va tim chi nhanh gan nhat.
- Tao, kich hoat, bao luu, tiep tuc, het han va hoan tien the tap.
- Diem danh, bao nghi va tap bu.
- Lich lop/session: tao, liet ke, xem chi tiet, dang ky, check-in theo session.
- Dang nhap nhan vien, refresh/logout token, middleware xac thuc va role guard.
- Quan ly nhan vien admin-only: tao, liet ke, xem chi tiet, cap nhat, reset password/deactivate.
- Dashboard/report aggregates cho admin/manager: summary KPI, revenue, plan distribution, recent
  members, today sessions.
- Frontend Staff Portal MVP bang React/Vite tich hop live API.
- Docker Compose full-stack va seed demo data phuc vu demo/cham bai.

Nhung phan van la dinh huong:

- Branch-scope authorization theo tung chi nhanh.
- Cron jobs va thong bao tu dong.
- Member App rieng cho hoc vien.
- Report export PDF/CSV va thong ke tai chinh nang cao.
- Thanh toan online, face recognition, audit log day du.

---

## 2. Current Backend Scope

| Area | Current implemented behavior | Report note |
|---|---|---|
| Members | Register, get by ID, activate offline payment, list member subscriptions | Implemented |
| Courses | CRUD course, price/session metadata for subscription snapshot | Implemented |
| Branches | CRUD branch, GeoJSON location, nearby query with `2dsphere` index | Implemented |
| Subscriptions | Create pending subscription, activate through member payment confirmation, suspend, unsuspend, expire, refund | Implemented |
| Attendance | Check-in, report missed, makeup, history by subscription | Implemented |
| Sessions | Create/list/get/enroll/check-in using existing attendance rules | Implemented |
| Auth | Employee login, access token, refresh token rotation, logout revoke, role guard | Implemented |
| Employees | Admin-only create/list/get/update/password reset/deactivate | Implemented |
| Reports/dashboard | Summary KPIs, revenue buckets, plan distribution, recent members, today sessions | Implemented as MVP |
| Frontend | React/Vite staff portal consuming live backend APIs | Implemented as MVP |
| Demo/package | Docker Compose full stack and deterministic seed command | Implemented as MVP |
| Cron/notifications | Not implemented | Planned |

---

## 3. Requirement Traceability Matrix

| Req ID | Requirement | Backend/API evidence | Data evidence | Test evidence | Status |
|---|---|---|---|---|---|
| FR-01 | Dang ky hoc vien | `POST /api/v1/members` | `members`, unique `ccid` index | Build/test + manual API history in phase notes | Implemented |
| FR-02 | Xem hoc vien va the tap | `GET /api/v1/members/:id`, `GET /api/v1/members/:id/subscriptions` | `members`, `subscriptions.member_id` | Automated build/test; API samples | Implemented |
| FR-03 | Tao the tap tu course/branch/member | `POST /api/v1/subscriptions` | `subscriptions`, refs to member/course/branch | Service tests for pricing/discount | Implemented |
| FR-04 | Thanh toan offline/kich hoat | `PATCH /api/v1/members/:id/activate` | `subscriptions.status`, `payment_date`, `members.is_registered` | Manual API phase evidence | Implemented |
| FR-05 | Bao luu/tiep tuc/het han | `PATCH /subscriptions/:id/suspend`, `/unsuspend`, `/expire` | `subscriptions.status`, `suspension` | Build/test + API samples | Implemented |
| FR-06 | Hoan tien | `POST /api/v1/subscriptions/:id/refund` | `refunds`, `subscriptions.status=refunded`, `remaining_sessions=0` | Service tests + manual API phase evidence | Implemented as MVP |
| FR-07 | Diem danh | `POST /api/v1/attendance/checkin` | `attendance`, `subscriptions.remaining_sessions`, `members.total_sessions_attended` | Manual API + service flow evidence | Implemented |
| FR-08 | Bao nghi/tap bu | `POST /attendance/report`, `POST /attendance/makeup` | `attendance.status`, `is_makeup_for` | Manual API + DB verification in phase notes | Implemented |
| FR-09 | Tim chi nhanh gan nhat | `GET /api/v1/branches/nearby` | `branches.location` GeoJSON + `2dsphere` index | Manual API phase evidence | Implemented |
| FR-10 | Lich lop/session | `/api/v1/sessions*` endpoints | `sessions`, enrolled subscription IDs | Build/test + API samples | Implemented |
| FR-11 | Dang nhap va phan quyen | `/api/v1/auth/login`, `/refresh`, `/logout`; protected business routes | `employees`, `refresh_tokens` | Unit tests, middleware tests, manual API + DB verification | Implemented |
| FR-12 | Quan ly nhan vien | `/api/v1/employees*` admin-only endpoints | `employees`, `refresh_tokens` revoke on reset/deactivate | Unit/service/integration and build evidence in phase notes | Implemented |
| FR-13 | Bao cao/thong ke | `/api/v1/dashboard/*` admin/manager endpoints | Aggregates from members/subscriptions/refunds/attendance/sessions | Build/test and dashboard phase evidence | Implemented as MVP |
| FR-14 | Frontend staff portal | React/Vite routes under `frontend/src` | Live API integration through bearer tokens | Frontend build/browser smoke phase evidence | Implemented as MVP |
| FR-15 | Demo/packaging | Docker Compose and `cmd/seed` | Deterministic demo records in MongoDB | Build/config/seed smoke evidence | Implemented as MVP |

---

## 4. Architecture Material

### 4.1 Layered Backend

Bao cao co the mo ta kien truc backend nhu sau:

```text
HTTP Request
  -> Gin route
  -> Handler: parse body/path/query, map error to HTTP response
  -> Service: validate business rules, orchestrate repositories
  -> Repository: MongoDB query/update/index
  -> MongoDB
```

Ly do phu hop:

- Handler khong chua business rules phuc tap.
- Service khong phu thuoc truc tiep vao `mongo.ErrNoDocuments`; repository map loi thanh
  `repository.ErrNotFound`.
- Repository giu chi tiet MongoDB va index, giup tang kha nang test va bao tri.

### 4.2 Security and Authorization

He thong auth hien tai danh cho staff/employee:

- Login bang email + password.
- Password duoc hash bang bcrypt.
- Access token ngan han, refresh token dai han hon.
- Refresh token duoc luu trong DB duoi dang SHA-256 hash, khong luu raw token.
- Refresh token duoc rotate: refresh thanh cong thi revoke token cu va cap token moi.
- Business routes bat buoc co `Authorization: Bearer <access_token>`.
- Role guard kiem tra cac role `admin`, `manager`, `trainer`, `receptionist`.

Gioi han can ghi ro:

- Branch-scope authorization chua ap dung trong cycle nay.
- Employee CRUD endpoint da co nhung van la admin-only; branch-scope permission theo tung chi nhanh
  chua duoc ap dung.
- JWT duoc ky HS256 bang stdlib crypto trong code hien tai, khong dung thu vien `golang-jwt/jwt`.

### 4.3 Data Model Highlights

| Collection | Vai tro trong bao cao |
|---|---|
| `members` | Dinh danh hoc vien, CCID unique, trang thai dang ky |
| `courses` | Mau goi tap, gia, so buoi, level/tag |
| `branches` | Chi nhanh, dia chi, GeoJSON location |
| `subscriptions` | The tap, trang thai, payment, discount, remaining sessions |
| `attendance` | Lich su check-in, report missed, makeup |
| `sessions` | Lich lop co trainer, capacity, enrolled subscriptions |
| `employees` | Staff auth identity, roles, branch assignment |
| `refresh_tokens` | Hash refresh token, expiry, revoke timestamp |
| `refunds` | Audit record cho hoan tien MVP |

Dashboard/report endpoints doc aggregate tu `members`, `subscriptions`, `refunds`, `attendance`,
va `sessions`; hien chua co collection report rieng.

---

## 5. Business Rules Material

### Subscription and Pricing

- Client khong duoc tu set tien thanh toan.
- Service lay `base_price`, `session_count`, `allowed_tags` tu course de snapshot vao subscription.
- Subscription moi o trang thai `pending`.
- Payment offline duoc confirm qua `PATCH /api/v1/members/:id/activate` voi `subscription_id`.
- Sau confirm, subscription thanh `active`, member thanh registered.

### Refund MVP

Hien tai refund implementation khac voi target rule 72h/50%/20% trong requirement ban dau:

- Chi cho refund subscription dang `active`.
- Tu choi `pending`, `suspended`, `expired`, `refunded`.
- Yeu cau `remaining_sessions > 0`.
- Cong thuc hien tai: `refund_amount = total_amount_paid * remaining_sessions / total_sessions`.
- Sau refund: subscription thanh `refunded`, `remaining_sessions = 0`, va ghi refund record.

Khi viet bao cao, nen ghi target refund rule 72h/50%/20% la **planned/business target**, con rule
tren la **MVP implemented rule**.

### Attendance, Report and Makeup

- Check-in thanh cong tao attendance, tru remaining sessions va tang total attended count.
- Weekly limit dua tren `session_per_week`.
- Report missed chi duoc ghi trong cua so 30 ngay.
- Makeup phai tham chieu ngay reported-missed trong 7 ngay va tieu thu mot session.
- Handler khong cho client tu dieu khien `reported_missed`/`makeup`; endpoint rieng set status.

### Branch Nearby

- Chi nhanh luu `location` theo GeoJSON: `type = Point`, `coordinates = [lng, lat]`.
- Query nearby dung `lng`, `lat`, optional `max_distance`, `limit`.
- Repository tao index `branches.location` loai `2dsphere`.

---

## 6. Testing Evidence Summary

Bang chung kiem thu co the dua vao chuong testing:

| Test layer | Evidence |
|---|---|
| Build | `env GOCACHE=/tmp/gocache go build ./...` pass |
| Automated tests | `env GOCACHE=/tmp/gocache go test ./...` pass |
| Service unit tests | Subscription pricing/refund; auth login/refresh/logout |
| Middleware tests | Missing token `401`, invalid token `401`, allowed role `200`, forbidden role `403`, unexpected service error `500` |
| Manual API | Auth login, protected route with/without token, refresh rotation, reused old refresh token, logout idempotency |
| DB verification | Admin bootstrap, bcrypt password hash, refresh-token hash storage, no raw token, `revoked_at` after refresh/logout |
| Integration tests | `internal/integration` uses real router wiring and isolated MongoDB test DBs |
| Frontend build/smoke | `npm --prefix frontend run build` and frontend phase browser smokes |
| Docker/demo package | `docker compose config/build` and deterministic `cmd/seed` demo data flow |

De viet bao cao, nen trinh bay testing theo 3 lop:

1. **Unit/service tests**: business rules.
2. **Middleware/router tests**: auth va role guard.
3. **Manual API + DB verification**: flow dau-cuoi voi MongoDB local.

---

## 7. Current Limitations and Future Work

### Can state as current limitations

- Chua co branch-scope guard.
- TTL index cleanup cho refresh tokens la eventual theo co che MongoDB, khong phai xoa tuc thi.
- Refresh rotation co residual availability risk: token cu bi revoke truoc khi replacement persist.
- Chua co transaction MongoDB cho mot so flow gom nhieu write nhu refund audit va attendance side
  effects.
- Dashboard/report hien la MVP read-only; chua co PDF/CSV export hoac report scheduling.
- Frontend Staff Portal la MVP cho nhan vien; chua co Member App rieng cho hoc vien.
- Chua co Cron jobs, notification, online payment, face recognition.

### Suitable future work list

1. Branch-scope authorization theo branch assignment.
2. Permanent Playwright regression suite cho frontend.
3. Report export PDF/CSV, scheduled reports, va deeper revenue/accounting views.
4. Payment gateway va audit log tai chinh.
5. Cron jobs, notification/reminder, va operational automation.
6. Member App rieng cho hoc vien.
7. MongoDB transaction hardening cho cac flow nhieu write neu he thong len production.

---

## 8. Suggested Final Report Structure

1. **Mo dau**
   - Boi canh gym da chi nhanh.
   - Van de: du lieu roi rac, check-in/hoan tien/bao luu thieu minh bach.
2. **Phan tich yeu cau**
   - Actors: Member, Receptionist, Trainer, Manager, Admin/System.
   - FR/NFR, phan tach current vs planned.
3. **Phan tich va thiet ke he thong**
   - Use case diagram.
   - ERD/logical data model.
   - Component/layer diagram.
   - Sequence cho: login, tao subscription + activate, check-in, report/makeup, refund.
4. **Thiet ke API va du lieu**
   - Khong copy toan bo API table; dan chieu `docs/api_contract.md`.
   - Neu can, chen bang rut gon cac endpoint implemented.
5. **Cai dat**
   - Go + Gin + MongoDB, layered architecture.
   - Auth/role guard.
   - Geo query, pricing/refund, attendance rules.
6. **Kiem thu**
   - Build/test commands.
   - Unit/middleware/manual API/DB verification.
   - Bang test case dai dien.
7. **Ket qua, gioi han, huong phat trien**
   - Implemented backend scope.
   - Limitations va future work.

---

## 9. Report-Ready Paragraph

"He thong Gym Management System duoc trien khai theo kien truc phan lop Handler - Service -
Repository - MongoDB, kem frontend Staff Portal bang React/Vite. Cach to chuc nay giup tach ro trach
nhiem: handler chi xu ly HTTP input/output, service dam nhan business rules, repository phu trach
truy cap MongoDB va index. Phien ban hien tai da hoan thanh cac luong nghiep vu cot loi gom quan ly
hoc vien, goi tap, chi nhanh, the tap, diem danh, bao nghi/tap bu, lich lop, tim chi nhanh gan nhat,
xac thuc/phan quyen nhan vien, quan ly staff va dashboard quan tri MVP. Qua trinh kiem thu ket hop
automated tests, middleware tests, integration tests, frontend build/browser smoke, manual API va xac
minh truc tiep tren MongoDB, dam bao cac rule quan trong nhu khong tin client ve tien, refresh token
chi luu dang hash, va role guard tra dung `401/403` duoc kiem chung."
