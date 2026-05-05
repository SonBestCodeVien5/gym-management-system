# PHASE 2 - BỔ SUNG & CẢI THIỆN THIẾT KẾ HỆ THỐNG

> Tài liệu này bổ sung cho `project_2_NguyenKhanhSon_20235212.md`, tập trung vào các phần còn thiếu hoặc cần làm rõ thêm.

---

## I. AUTHENTICATION & AUTHORIZATION

### 1. Cơ chế xác thực (JWT)

Hệ thống sử dụng **JWT (JSON Web Token)** stateless với cặp Access Token / Refresh Token:

- **Access Token:** TTL = 15 phút. Gửi kèm trong header `Authorization: Bearer <token>` mỗi request.
- **Refresh Token:** TTL = 7 ngày. Lưu trong HttpOnly cookie, dùng để cấp lại Access Token mà không cần đăng nhập lại.

**Payload của JWT:**
```json
{
  "sub": "<employeeId hoặc memberId>",
  "role": ["Receptionist"],
  "branchId": ["<id1>", "<id2>"],
  "level": "Advanced",
  "iat": 1700000000,
  "exp": 1700000900
}
```

### 2. Bảng phân quyền theo Role

| Endpoint | Member | Receptionist | Trainer | Manager |
|----------|--------|--------------|---------|---------|
| POST /members | ✗ | ✓ | ✗ | ✓ |
| POST /subscriptions | ✗ | ✓ | ✗ | ✓ |
| PATCH /subscriptions/:id/suspend | ✗ | ✓ | ✗ | ✓ |
| POST /attendance/checkin | ✗ | ✓ | ✗ | ✓ |
| POST /attendance/report | ✓ | ✓ | ✗ | ✓ |
| GET /attendance/history | ✓ (bản thân) | ✓ | ✓ | ✓ |
| GET /subscriptions/:id | ✓ (bản thân) | ✓ | ✗ | ✓ |
| GET /branches/nearby | ✓ | ✓ | ✓ | ✓ |
| POST /employees | ✗ | ✗ | ✗ | ✓ |
| GET /reports/* | ✗ | ✗ | ✗ | ✓ |

> **Quy tắc bổ sung:**
> - Receptionist chỉ thao tác được với dữ liệu thuộc `branchId` của mình (kiểm tra trong middleware).
> - Manager chỉ quản lý được branch mà mình là `managerId`.

### 3. Middleware xác thực (Go)

```
Request → [AuthMiddleware] → [RoleGuard] → [BranchScopeGuard] → Handler
```

- **AuthMiddleware:** Verify JWT signature, kiểm tra expiry.
- **RoleGuard:** So sánh `role` trong token với role được phép của endpoint.
- **BranchScopeGuard:** Với các role giới hạn chi nhánh, kiểm tra `branchId` trong token có chứa `branchId` của resource đang truy cập không.

---

## II. API ENDPOINTS - BỔ SUNG ĐẦY ĐỦ

### Auth

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/auth/login | Đăng nhập, trả về Access Token + Refresh Token. |
| POST | /api/v1/auth/refresh | Cấp lại Access Token bằng Refresh Token. |
| POST | /api/v1/auth/logout | Hủy Refresh Token. |

### Members

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/members | Đăng ký hồ sơ học viên mới (Check CCID). |
| GET | /api/v1/members/:id | Xem hồ sơ học viên. |
| PATCH | /api/v1/members/:id | Cập nhật thông tin học viên. |

### Subscriptions

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/subscriptions | Tạo thẻ tập & kích hoạt thanh toán. |
| GET | /api/v1/subscriptions/:id | Xem chi tiết thẻ tập (trạng thái, số buổi còn lại). |
| GET | /api/v1/members/:id/subscriptions | Lịch sử toàn bộ thẻ tập của học viên. |
| PATCH | /api/v1/subscriptions/:id/suspend | Gửi yêu cầu bảo lưu. |
| PATCH | /api/v1/subscriptions/:id/unsuspend | Kích hoạt lại thẻ sau bảo lưu (thủ công). |
| POST | /api/v1/subscriptions/:id/refund | Yêu cầu hoàn tiền (áp dụng rule 72h/50%). |

### Attendance

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/attendance/checkin | Điểm danh (trừ remainingSessions, kiểm tra quota tuần). |
| POST | /api/v1/attendance/report | Báo nghỉ có phép (kiểm tra luật 30 ngày). |
| POST | /api/v1/attendance/makeup | Đăng ký buổi tập bù (kiểm tra Reported_Missed 7 ngày). |
| GET | /api/v1/subscriptions/:id/attendance | Lịch sử điểm danh của một thẻ tập. |

### Branches

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | /api/v1/branches/nearby | Tìm chi nhánh gần nhất theo tọa độ GPS. |
| GET | /api/v1/branches/:id | Xem thông tin chi nhánh. |
| GET | /api/v1/branches/:id/stats | Thống kê check-in, học viên active (Manager only). |

### Employees

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/employees | Tạo hồ sơ nhân viên mới. |
| GET | /api/v1/employees/:id | Xem thông tin nhân viên. |
| PATCH | /api/v1/employees/:id | Cập nhật vai trò / chi nhánh của nhân viên. |

### Courses

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | /api/v1/courses | Danh sách gói tập (có filter theo level). |
| POST | /api/v1/courses | Tạo gói tập mới (Manager only). |
| PATCH | /api/v1/courses/:id | Cập nhật gói tập. |

---

## III. CƠ CHẾ ENFORCE `sessionPerWeek`

Field `sessionPerWeek` trong Subscriptions cần được kiểm tra tại mỗi lần check-in. Logic bổ sung cho endpoint `POST /attendance/checkin`:

```
1. Xác định tuần hiện tại: [Thứ Hai ~ Chủ Nhật] chứa ngày check-in.
2. Đếm số bản ghi Attendance của subId này trong tuần đó
   với status IN (Attended, Makeup).
3. Nếu count >= sessionPerWeek → Trả về lỗi 409:
   "Đã đạt giới hạn buổi tập trong tuần này."
4. Nếu count < sessionPerWeek → Cho phép check-in, tạo bản ghi Attended.
```

> **Lưu ý:** Buổi tập bù (Makeup) vẫn tính vào quota tuần để tránh lạm dụng.

---

## IV. CRON JOBS - BỔ SUNG ĐẦY ĐỦ

| Job | Lịch chạy | Mô tả |
|-----|-----------|-------|
| `expire-suspensions` | Hàng ngày 00:00 | Nếu `currentDate > suspension.endDate` → status = `Expired`, `remainingSessions = 0`. |
| `expire-subscriptions` | Hàng ngày 00:00 | Nếu `remainingSessions == 0` và status = `Active` → status = `Expired`. |
| `notify-low-sessions` | Hàng ngày 08:00 | Nếu `remainingSessions <= 3` → Gửi thông báo nhắc gia hạn (push/email). |
| `weekly-quota-reset` | Mỗi Thứ Hai 00:00 | (Tùy chọn) Log hoặc reset bộ đếm tuần nếu dùng cache thay vì đếm trực tiếp từ DB. |
| `geo-index-health` | Hàng tuần | Kiểm tra index `2dsphere` trên collection Branches còn hoạt động bình thường. |

> Với backend Go, các Cron Job nên được triển khai bằng thư viện `robfig/cron` hoặc tách thành một worker service riêng để không block main server.

---

## V. SƠ ĐỒ THỰC THỂ QUAN HỆ (ERD) - MÔ TẢ CHI TIẾT

Do tài liệu gốc chỉ mô tả bằng text, dưới đây là mô tả đầy đủ để vẽ ERD:

```
┌──────────┐       ┌───────────────┐       ┌─────────┐
│ Branches │ 1───N │ Subscriptions │ N───1 │ Members │
│          │       │               │       │         │
│          │ 1───N │  Attendance   │       └─────────┘
└──────────┘       └───────────────┘
     │ 1                  │ 1
     │                    │ N
     N                    │
┌───────────┐      ┌──────────────┐
│ Employees │      │  Attendance  │
│           │      │ (chi tiết)   │
└───────────┘      └──────────────┘
                          │ N───1
                    ┌─────────┐
                    │ Courses │
                    └─────────┘
```

**Quan hệ đầy đủ:**

- `Members` (1) → (N) `Subscriptions`: Một học viên có nhiều thẻ tập theo thời gian.
- `Courses` (1) → (N) `Subscriptions`: Một gói tập mẫu sinh ra nhiều thẻ tập thực tế.
- `Branches` (1) → (N) `Subscriptions`: Chi nhánh gốc (home branch) của thẻ tập.
- `Subscriptions` (1) → (N) `Attendance`: Mỗi thẻ tập có nhiều bản ghi điểm danh.
- `Branches` (1) → (N) `Attendance`: Ghi nhận học viên tập ở chi nhánh nào (roaming).
- `Employees` (N) ↔ (N) `Branches`: Một nhân viên làm tại nhiều chi nhánh (lưu qua `branchId[]` trong Employees).

---

## VI. BỔ SUNG COLLECTION: Sessions (Tùy chọn nâng cao)

Nếu hệ thống cần quản lý **lịch tập theo slot giờ cụ thể** (không chỉ điểm danh tự do), cân nhắc thêm collection này:

- **_id**: ObjectId.
- **branchId**: Ref (Branches).
- **trainerId**: Ref (Employees).
- **courseLevel**: Enum (Basic, Advanced, Professional).
- **scheduledAt**: Date (ngày + giờ bắt đầu).
- **durationMinutes**: Int.
- **capacity**: Int (sĩ số tối đa).
- **enrolledCount**: Int (số học viên đã đăng ký).

Khi đó `Attendance` sẽ có thêm field `sessionId`: Ref (Sessions) để liên kết buổi tập cụ thể.

---

## VII. CÁC RÀNG BUỘC BỔ SUNG

- **Không tạo 2 thẻ Active cùng lúc:** Trước khi tạo Subscription mới, kiểm tra member đã có thẻ nào `status = Active` chưa. Nếu có → yêu cầu kết thúc hoặc hoàn tiền thẻ cũ trước.
- **Refund chỉ do Manager thực hiện:** Endpoint `/refund` yêu cầu role Manager và ghi log lý do.
- **Audit log:** Các thao tác quan trọng (tạo thẻ, hoàn tiền, bảo lưu) nên ghi vào một collection `AuditLogs` riêng với `{ actorId, action, targetId, timestamp, metadata }` để phục vụ tra cứu khi có tranh chấp.
