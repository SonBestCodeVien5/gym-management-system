# PHASE 4: KẾ HOẠCH TRIỂN KHAI & TECH STACK

---

> Ghi chú đối chiếu hiện trạng: tài liệu này mô tả kiến trúc mục tiêu và lộ trình triển khai để đưa
> vào báo cáo. Hiện tại repo đã hoàn thành backend Go + Gin + MongoDB, Staff Portal React/Vite,
> Docker Compose full stack, nginx frontend container và `cmd/seed` demo data. Redis, cron worker,
> CI/CD, Member App, report export và một số thư viện đề xuất bên dưới là định hướng/future work.
> Xem thêm [07_current_implementation_evidence.md](07_current_implementation_evidence.md).

## I. KIẾN TRÚC HỆ THỐNG TỔNG THỂ

Hệ thống được thiết kế theo mô hình **Monolithic Modular** — đủ đơn giản để một nhóm nhỏ phát triển và bảo trì, nhưng có ranh giới module rõ ràng để dễ tách thành microservice sau này nếu cần.

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                         │
│                                                             │
│   ┌──────────────────┐          ┌──────────────────┐        │
│   │   Staff Portal   │          │   Member App     │        │
│   │  (React / Web)   │          │ (React / Mobile) │        │
│   └────────┬─────────┘          └────────┬─────────┘        │
└────────────┼───────────────────────────────────────────────┘
             │ HTTPS / REST API
┌────────────▼───────────────────────────────────────────────┐
│                       API GATEWAY                           │
│              (Rate Limiting, Auth Middleware)               │
└────────────┬───────────────────────────────────────────────┘
             │
┌────────────▼───────────────────────────────────────────────┐
│                    BACKEND LAYER (Go)                       │
│                                                             │
│  ┌──────────┐ ┌────────────┐ ┌──────────┐ ┌────────────┐  │
│  │  Member  │ │Subscription│ │Attendance│ │  Branch &  │  │
│  │  Module  │ │  Module    │ │  Module  │ │  Employee  │  │
│  └──────────┘ └────────────┘ └──────────┘ └────────────┘  │
│                                                             │
│  ┌──────────────────────┐    ┌──────────────────────────┐  │
│  │    Auth Module       │    │    Cron Job Worker       │  │
│  │  (JWT + Middleware)  │    │  (robfig/cron)           │  │
│  └──────────────────────┘    └──────────────────────────┘  │
└────────────┬───────────────────────────────────────────────┘
             │
┌────────────▼───────────────────────────────────────────────┐
│                      DATA LAYER                             │
│                                                             │
│   ┌─────────────────┐         ┌────────────────────────┐   │
│   │    MongoDB      │         │    Redis (Cache)        │   │
│   │  (Primary DB)   │         │  (Session, Rate Limit)  │   │
│   └─────────────────┘         └────────────────────────┘   │
└────────────────────────────────────────────────────────────┘
```

---

## II. TECH STACK CHI TIẾT

### Backend

| Thành phần | Công nghệ | Lý do chọn |
|-----------|-----------|------------|
| Ngôn ngữ | **Go** | Hiệu năng cao, xử lý concurrent tốt, phù hợp API server |
| Web Framework | **Gin** | Nhẹ, nhanh, middleware ecosystem phong phú |
| MongoDB Driver | **mongo-go-driver** | Driver chính thức, hỗ trợ đầy đủ aggregation pipeline |
| Auth | **HS256 JWT-compatible token bằng Go stdlib crypto + bcrypt** | Đã triển khai access/refresh token, refresh hash storage và role guard; có thể thay bằng thư viện JWT chuẩn khi hardening |
| Cron Jobs | **robfig/cron v3** | Planned/future work |
| Config | **godotenv + environment variables** | Hiện đang đọc `.env`; `viper` là hướng nâng cấp nếu config phức tạp hơn |
| Logging | **Go log + Gin logger** | Đủ cho MVP; `zap` là hướng nâng cấp structured logging |
| Validation | **Manual validation in handlers/services** | Đang dùng parse/validation thủ công; validator tags là hướng nâng cấp |

### Frontend

| Thành phần | Công nghệ | Lý do chọn |
|-----------|-----------|------------|
| Framework | **React 18** | Đã triển khai Staff Portal MVP với component/routing/state trong repo |
| Language | **JavaScript / JSX** | Phù hợp phạm vi MVP; TypeScript là hướng nâng cấp nếu cần type-safety sâu hơn |
| State/API | **React state/context + fetch client nội bộ** | Đủ cho auth, dashboard và các module staff portal hiện tại |
| Styling | **CSS thuần trong `frontend/src/index.css`** | Dễ kiểm soát giao diện MVP và không phụ thuộc UI framework ngoài |
| Build Tool | **Vite** | Build nhanh, dễ đóng gói vào nginx container |

### Database & Infrastructure

| Thành phần | Công nghệ | Ghi chú |
|-----------|-----------|---------|
| Database | **MongoDB 7.0** | Index `2dsphere` cho Geo-query |
| Cache | **Redis 7** | Planned; refresh token hiện lưu dạng hash trong MongoDB |
| Deployment | **Docker + Docker Compose** | Đóng gói nhất quán, dễ chạy local và production |
| Static frontend server | **Nginx** | Đã dùng trong `frontend/Dockerfile` để serve Vite build |
| CI/CD | **GitHub Actions** | Planned |

---

## III. CẤU TRÚC THƯ MỤC DỰ ÁN (Backend Go)

```
gym-system/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── auth/                    # JWT, middleware
│   ├── member/                  # Handler, service, repository
│   ├── subscription/
│   ├── attendance/
│   ├── branch/
│   ├── employee/
│   ├── course/
│   └── cron/                    # Cron job definitions
├── pkg/
│   ├── database/                # MongoDB connection
│   ├── cache/                   # Redis client
│   ├── config/                  # Viper config loader
│   └── logger/                  # Zap logger setup
├── docker-compose.yml
├── Dockerfile
└── .env.example
```

---

## IV. KẾ HOẠCH PHÁT TRIỂN (TIMELINE)

Tổng thời gian dự kiến ban đầu: **10 tuần**. Trạng thái dưới đây ghi rõ phần đã hoàn thành trong MVP
và phần còn là future work.

### Giai đoạn 1 — Nền tảng (Tuần 1–2)
- [x] Thiết lập project Go + Gin, kết nối MongoDB.
- [x] Tạo Docker Compose full stack: MongoDB, API, frontend và seed profile.
- [x] Implement Auth module: đăng nhập, current user, refresh/logout, middleware phân quyền.
- [x] Thiết lập cấu trúc thư mục backend theo `cmd/`, `internal/`, `pkg/`.
- [ ] Redis cache/rate limit: future work.

### Giai đoạn 2 — Core Business (Tuần 3–5)
- [x] Member registration/get/activate/list subscriptions với `ccid` unique.
- [x] CRUD Courses.
- [x] Tạo và quản lý Subscriptions, snapshot pricing/session/allowed tags từ course.
- [x] Điểm danh, kiểm tra quota tuần, trừ buổi.
- [x] Logic báo nghỉ rolling 30 ngày và tập bù 7 ngày.

### Giai đoạn 3 — Tính năng nâng cao (Tuần 6–7)
- [x] Geo-query tìm chi nhánh gần nhất.
- [x] Hoàn tiền MVP theo số buổi còn lại và ghi collection `refunds`.
- [x] Bảo lưu/thôi bảo lưu/hết hạn thẻ tập.
- [x] Sessions: create/list/get/enroll/check-in.
- [x] Employee management admin-only.
- [x] Dashboard/report aggregate APIs MVP.
- [ ] Cron Jobs, notification và audit log nâng cao: future work.

### Giai đoạn 4 — Frontend (Tuần 7–9)
- [x] Staff Portal MVP: dashboard, members, subscriptions, attendance, sessions, employees,
  courses, branches.
- [x] Tích hợp live API qua bearer token và auth restore.
- [x] Docker frontend image serve Vite build bằng nginx.
- [ ] Member App riêng cho học viên: future work.

### Giai đoạn 5 — Kiểm thử & Hoàn thiện (Tuần 10)
- [x] Unit/service tests cho business logic quan trọng.
- [x] Integration tests dùng router thật và MongoDB test DB riêng.
- [x] Index bootstrap tập trung cho unique/query/partial unique/TTL indexes.
- [x] Tài liệu API bằng `docs/api_contract.md` và `api_test.http`.
- [x] Demo package: README, Docker Compose, seed data, final report assembly.
- [ ] Swagger/OpenAPI và CI/CD: future work.

---

## V. CHIẾN LƯỢC TESTING

### Unit Test (Go)
Tập trung vào các hàm business logic không phụ thuộc database:

```
✅ Tính refund amount (3 trường hợp)
✅ Kiểm tra rolling window 30 ngày báo nghỉ
✅ Kiểm tra window 7 ngày tập bù
✅ Kiểm tra quota tuần (sessionPerWeek)
✅ Kiểm tra phân cấp Trainer level
```

### Integration Test
Chạy trên môi trường Docker với MongoDB test instance:
```
✅ Luồng đăng ký → tạo thẻ → check-in
✅ Luồng báo nghỉ → tập bù
✅ Luồng bảo lưu → thôi bảo lưu / hết hạn thủ công
✅ Hoàn tiền MVP + kiểm tra refund record
```

### Manual Test
- Kiểm tra Geo-query với tọa độ thực tế.
- Kiểm tra phân quyền từng role trên Staff Portal.
- Kiểm tra Docker Compose, seed data, login demo admin và dashboard frontend.
- Geo-fencing trên thiết bị di động là future work.

---

## VI. RỦI RO & PHƯƠNG ÁN DỰ PHÒNG

| Rủi ro | Mức độ | Phương án xử lý |
|--------|--------|-----------------|
| Nhận diện khuôn mặt phức tạp, tốn thời gian | Cao | Fallback sang CCCD / mã thẻ thủ công cho MVP |
| MongoDB Geo-query chậm khi dữ liệu lớn | Trung bình | Đảm bảo index `2dsphere` được tạo từ đầu |
| Đồng bộ Cron Job chạy nhiều instance | Trung bình | Dùng Redis lock để đảm bảo mỗi job chỉ chạy 1 lần |
| Member App chưa triển khai | Trung bình | Staff Portal MVP đã hoàn thành; Member App làm sau |
