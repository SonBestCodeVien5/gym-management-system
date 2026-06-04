# Final Report Assembly - Gym Management System

Tai lieu nay la ban lap ghep bao cao cuoi cung tu cac report materials hien co. Khi can nop ban
chinh thuc, co the dung cau truc nay de chuyen sang `.docx` hoac PDF.

## 1. Mo dau

Gym Management System la he thong quan ly phong gym da chi nhanh. Muc tieu cua he thong la gom cac
nghiep vu dang roi rac nhu quan ly hoc vien, goi tap, chi nhanh, the tap, diem danh, lich lop, nhan
vien va dashboard quan tri vao mot ung dung co API va giao dien nhan vien ro rang.

Van de thuc te:

- Du lieu hoc vien va the tap de bi trung lap neu quan ly thu cong.
- Diem danh, bao nghi va tap bu can quy tac nhat quan.
- Nhan vien can dang nhap va duoc phan quyen theo vai tro.
- Quan ly can xem nhanh doanh thu, lop trong ngay va tinh hinh goi tap.

## 2. Phan tich yeu cau

Tac nhan chinh:

- Admin: quan ly nhan vien, cau hinh he thong, xem dashboard.
- Manager: quan ly course/branch/session va xem dashboard.
- Receptionist: tao hoc vien, the tap, kich hoat thanh toan offline va diem danh.
- Trainer: tao/xem session, enroll va check-in session.
- Member: doi tuong duoc quan ly trong he thong.

Chuc nang da implement:

- Auth va role guard.
- Employee management admin-only.
- Member/course/branch/subscription management.
- Attendance, report missed, makeup.
- Sessions va session check-in.
- Dashboard/report aggregate APIs.
- React/Vite staff portal.
- Docker Compose va seed demo data.

Chuc nang future work:

- Branch-scope authorization theo tung chi nhanh.
- Online payment, notification, cron jobs.
- Report export PDF/CSV va scheduled reports.
- Member App rieng cho hoc vien.
- Audit log tai chinh day du.

## 3. Thiet ke he thong

Backend dung kien truc phan lop:

```text
HTTP Request
  -> Gin Router
  -> Handler
  -> Service
  -> Repository
  -> MongoDB
```

Ly do chon kien truc nay:

- Handler chi xu ly request/response.
- Service tap trung business rules nhu pricing, refund, attendance window va role behavior.
- Repository tach MongoDB query/update/index khoi logic nghiep vu.
- `internal/app` gom route/dependency wiring dung chung cho server va integration tests.

Du lieu duoc luu trong cac collection:

- `members`
- `courses`
- `branches`
- `subscriptions`
- `attendance`
- `sessions`
- `employees`
- `refresh_tokens`
- `refunds`

MongoDB indexes duoc bootstrap tap trung trong `pkg/database.EnsureIndexes`.

## 4. Thiet ke API

API dung prefix `/api/v1`. Cac route business can:

```http
Authorization: Bearer <access_token>
```

Nhom endpoint chinh:

- `/api/v1/auth/*`
- `/api/v1/employees*`
- `/api/v1/members*`
- `/api/v1/courses*`
- `/api/v1/branches*`
- `/api/v1/subscriptions*`
- `/api/v1/attendance*`
- `/api/v1/sessions*`
- `/api/v1/dashboard*`

Loi HTTP dung contract chung:

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "invalid request body",
    "details": {}
  }
}
```

API chi tiet xem `docs/api_contract.md`.

## 5. Cai dat

Backend:

- Go + Gin.
- MongoDB official driver.
- JWT-like HMAC token implementation bang stdlib crypto.
- Refresh token chi luu dang SHA-256 hash.
- Password nhan vien hash bang bcrypt.

Frontend:

- React 18 + Vite.
- Bearer-token API client.
- Staff portal gom dashboard va cac module nghiep vu.

Packaging:

- `Dockerfile` build backend server va seed binary.
- `frontend/Dockerfile` build static Vite app va serve bang nginx.
- `docker-compose.yml` chay MongoDB, API, frontend.
- `cmd/seed` nap demo data idempotent.

## 6. Kiem thu

Lenh kiem thu backend:

```bash
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
```

Lenh kiem thu frontend:

```bash
npm --prefix frontend run build
```

Lenh kiem thu Docker:

```bash
docker compose config
docker compose build
```

Tren may thieu plugin `docker-buildx`, qua trinh kiem thu Docker da dung fallback:

```bash
DOCKER_BUILDKIT=0 docker compose build
```

Nhan xet:

- Unit/service tests tap trung vao auth va subscription rules.
- Middleware tests kiem `401/403`.
- Integration tests dung router that va MongoDB test DB rieng.
- Frontend da co build va browser smoke trong cac phase frontend.
- Seed data giup demo dashboard/module khong bi rong.
- Final package smoke da chay full stack tren Docker volume sach, seed hai lan de kiem idempotency,
  login demo admin va dashboard frontend bang Playwright.

## 7. Ket qua va gioi han

Ket qua:

- Backend API hoan thanh cac nghiep vu cot loi.
- Frontend staff portal co the login va thao tac cac module chinh.
- Docker Compose giup chay full stack bang mot lenh.
- Demo seed data giup nguoi cham co san tai khoan va du lieu mau.

Gioi han:

- Chua co branch-scope guard theo tung chi nhanh.
- Chua co online payment.
- Chua co cron/notification.
- Dashboard/report moi la MVP read-only, chua export PDF/CSV.
- Chua co Member App rieng.

## 8. Tai lieu tham chieu

- `README.md`
- `docs/api_contract.md`
- `docs/local_dev_guide.md`
- `docs/code_reading_guide.md`
- `docs/faq_why.md`
- `docs/report-materials/07_current_implementation_evidence.md`
