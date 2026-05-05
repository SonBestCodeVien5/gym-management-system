# PHASE 4: KẾ HOẠCH TRIỂN KHAI & TECH STACK

---

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
| Ngôn ngữ | **Go 1.22** | Hiệu năng cao, xử lý concurrent tốt, phù hợp API server |
| Web Framework | **Gin** | Nhẹ, nhanh, middleware ecosystem phong phú |
| MongoDB Driver | **mongo-go-driver v1.15** | Driver chính thức, hỗ trợ đầy đủ aggregation pipeline |
| Auth | **golang-jwt/jwt v5** | JWT chuẩn, dễ tích hợp middleware |
| Cron Jobs | **robfig/cron v3** | Cú pháp cron quen thuộc, hỗ trợ timezone |
| Config | **viper** | Đọc config từ file `.env` và biến môi trường |
| Logging | **zap (Uber)** | Structured logging, hiệu năng cao |
| Validation | **go-playground/validator** | Validation tag trực tiếp trên struct |

### Frontend

| Thành phần | Công nghệ | Lý do chọn |
|-----------|-----------|------------|
| Framework | **React 18 + TypeScript** | Type-safe, ecosystem lớn |
| State Management | **Zustand** | Nhẹ hơn Redux, đủ dùng cho quy mô dự án này |
| UI Library | **shadcn/ui + Tailwind CSS** | Đẹp, dễ tùy chỉnh, không lock-in |
| HTTP Client | **Axios** | Interceptor tiện cho việc gắn JWT và refresh token |
| Map | **Leaflet.js** | Mã nguồn mở, không tốn phí API như Google Maps |
| Build Tool | **Vite** | Nhanh hơn CRA nhiều lần |

### Database & Infrastructure

| Thành phần | Công nghệ | Ghi chú |
|-----------|-----------|---------|
| Database | **MongoDB 7.0** | Index `2dsphere` cho Geo-query |
| Cache | **Redis 7** | Lưu Refresh Token, rate limiting |
| Deployment | **Docker + Docker Compose** | Đóng gói nhất quán, dễ chạy local và production |
| Reverse Proxy | **Nginx** | SSL termination, load balancing |
| CI/CD | **GitHub Actions** | Tự động test và deploy khi push lên main |

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

Tổng thời gian dự kiến: **10 tuần**

### Giai đoạn 1 — Nền tảng (Tuần 1–2)
- [ ] Thiết lập project Go + Gin, kết nối MongoDB, Redis.
- [ ] Tạo Docker Compose (Go + MongoDB + Redis).
- [ ] Implement Auth module: đăng nhập, JWT, middleware phân quyền.
- [ ] Thiết lập cấu trúc thư mục và coding conventions.

### Giai đoạn 2 — Core Business (Tuần 3–5)
- [ ] CRUD Members (có kiểm tra CCCD unique).
- [ ] CRUD Courses.
- [ ] Tạo & quản lý Subscriptions (tính tiền, ưu đãi, unitPrice).
- [ ] Điểm danh (kiểm tra quota tuần, trừ buổi).
- [ ] Logic báo nghỉ (rolling 30 ngày) và tập bù (7 ngày).

### Giai đoạn 3 — Tính năng nâng cao (Tuần 6–7)
- [ ] Geo-query tìm chi nhánh gần nhất.
- [ ] Logic hoàn tiền (3 trường hợp + audit log).
- [ ] Bảo lưu thẻ tập (suspend/unsuspend).
- [ ] Cron Jobs: expire subscription, expire suspension, thông báo sắp hết buổi.

### Giai đoạn 4 — Frontend (Tuần 7–9)
- [ ] Staff Portal: Dashboard, Check-in, Quản lý thẻ tập, Quản lý học viên.
- [ ] Member App: Trang chủ, Lịch sử tập, Tìm chi nhánh, Báo nghỉ.
- [ ] Tích hợp API + xử lý JWT refresh token tự động.

### Giai đoạn 5 — Kiểm thử & Hoàn thiện (Tuần 10)
- [ ] Unit test các module business logic quan trọng (hoàn tiền, báo nghỉ).
- [ ] Integration test các luồng chính (đăng ký → check-in → hoàn tiền).
- [ ] Fix bug, tối ưu query MongoDB (kiểm tra index).
- [ ] Viết tài liệu API (Swagger/OpenAPI).
- [ ] Demo và bàn giao.

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
✅ Luồng bảo lưu → hết hạn bảo lưu (Cron)
✅ Hoàn tiền + kiểm tra AuditLog
```

### Manual Test
- Kiểm tra Geo-query với tọa độ thực tế.
- Kiểm tra Geo-fencing trên thiết bị di động.
- Kiểm tra phân quyền từng role trên Staff Portal.

---

## VI. RỦI RO & PHƯƠNG ÁN DỰ PHÒNG

| Rủi ro | Mức độ | Phương án xử lý |
|--------|--------|-----------------|
| Nhận diện khuôn mặt phức tạp, tốn thời gian | Cao | Fallback sang CCCD / mã thẻ thủ công cho MVP |
| MongoDB Geo-query chậm khi dữ liệu lớn | Trung bình | Đảm bảo index `2dsphere` được tạo từ đầu |
| Đồng bộ Cron Job chạy nhiều instance | Trung bình | Dùng Redis lock để đảm bảo mỗi job chỉ chạy 1 lần |
| Thiếu thời gian hoàn thành Frontend | Cao | Ưu tiên Staff Portal, Member App làm sau |
