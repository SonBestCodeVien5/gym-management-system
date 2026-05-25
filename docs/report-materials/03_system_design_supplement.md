# PHASE 2 - BỔ SUNG & CẢI THIỆN THIẾT KẾ HỆ THỐNG

> Tài liệu này bổ sung cho `project_2_NguyenKhanhSon_20235212.md`, tập trung vào các phần còn thiếu hoặc cần làm rõ thêm.

---

## I. AUTHENTICATION & AUTHORIZATION

### 1. Cơ chế xác thực (JWT)

Hệ thống sử dụng **JWT (JSON Web Token)** stateless với cặp Access Token / Refresh Token:

- **Access Token:** TTL = 15 phút. Gửi kèm trong header `Authorization: Bearer <token>` mỗi request.
- **Refresh Token:** TTL = 7 ngày. Backend hiện nhận refresh token trong request body, lưu trong
  MongoDB dưới dạng hash để phục vụ revoke/rotation. HttpOnly cookie có thể là cải tiến frontend
  sau này.

**Payload JWT trong implementation hiện tại:**
```json
{
  "employee_id": "<employeeObjectId>",
  "role": ["receptionist"],
  "token_type": "access",
  "jti": "<random-token-id>",
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
> - Role guard đã implemented.
> - Branch-scope authorization theo `branchId` là planned/future work.

### 3. Middleware xác thực (Go)

```
Request → [AuthMiddleware] → [RoleGuard] → [BranchScopeGuard] → Handler
```

- **AuthMiddleware:** Verify JWT signature, kiểm tra expiry.
- **RoleGuard:** So sánh `role` trong token với role được phép của endpoint.
- **BranchScopeGuard:** Planned/future work. Với các role giới hạn chi nhánh, kiểm tra `branchId`
  trong token có chứa `branchId` của resource đang truy cập không.

---

## II. API ENDPOINTS - BỔ SUNG ĐẦY ĐỦ

Ghi chu:
- Implemented: da co trong code hien tai.
- Planned: dung cho dinh huong, chua co trong code.

### Auth

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/auth/login | Đăng nhập, trả về Access Token + Refresh Token. (implemented) |
| POST | /api/v1/auth/refresh | Rotate refresh token và cấp lại access token. (implemented) |
| POST | /api/v1/auth/logout | Hủy refresh token. (implemented) |

### Members

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/members | Đăng ký hồ sơ học viên mới (Check CCID). (implemented) |
| GET | /api/v1/members/:id | Xem hồ sơ học viên. (implemented) |
| PATCH | /api/v1/members/:id/activate | Confirm offline payment cho subscription. (implemented) |
| PATCH | /api/v1/members/:id | Cập nhật thông tin học viên. (planned) |

### Subscriptions

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/subscriptions | Tạo thẻ tập (status = pending). (implemented) |
| GET | /api/v1/subscriptions/:id | Xem chi tiết thẻ tập (trạng thái, số buổi còn lại). (implemented) |
| GET | /api/v1/members/:id/subscriptions | Lịch sử toàn bộ thẻ tập của học viên. (implemented) |
| PATCH | /api/v1/subscriptions/:id/suspend | Gửi yêu cầu bảo lưu. (implemented) |
| PATCH | /api/v1/subscriptions/:id/unsuspend | Kích hoạt lại thẻ sau bảo lưu (thủ công). (implemented) |
| PATCH | /api/v1/subscriptions/:id/expire | Hết hạn thu cong. (implemented) |
| POST | /api/v1/subscriptions/:id/refund | Yêu cầu hoàn tiền theo số buổi còn lại. (implemented MVP) |

### Attendance

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/attendance/checkin | Điểm danh (trừ remainingSessions, kiểm tra quota tuần). (implemented) |
| POST | /api/v1/attendance/report | Báo nghỉ có phép (kiểm tra luật 30 ngày). (implemented) |
| POST | /api/v1/attendance/makeup | Đăng ký buổi tập bù (kiểm tra Reported_Missed 7 ngày). (implemented) |
| GET | /api/v1/subscriptions/:id/attendance | Lịch sử điểm danh của một thẻ tập. (implemented) |

### Branches

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | /api/v1/branches/nearby | Tìm chi nhánh gần nhất theo tọa độ GPS. (implemented) |
| GET | /api/v1/branches | Danh sach chi nhanh. (implemented) |
| GET | /api/v1/branches/:id | Xem thông tin chi nhánh. (implemented) |
| POST | /api/v1/branches | Tao chi nhanh. (implemented) |
| PATCH | /api/v1/branches/:id | Cap nhat chi nhanh. (implemented) |
| DELETE | /api/v1/branches/:id | Xoa chi nhanh. (implemented) |
| GET | /api/v1/branches/:id/stats | Thống kê check-in, học viên active (Manager only). (planned) |

### Employees

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/employees | Tạo hồ sơ nhân viên mới. (planned) |
| GET | /api/v1/employees/:id | Xem thông tin nhân viên. (planned) |
| PATCH | /api/v1/employees/:id | Cập nhật vai trò / chi nhánh của nhân viên. (planned) |

### Courses

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | /api/v1/courses | Danh sách gói tập (có filter theo level). (implemented) |
| GET | /api/v1/courses/:id | Xem chi tiet goi tap. (implemented) |
| POST | /api/v1/courses | Tạo gói tập mới (Manager only). (implemented) |
| PATCH | /api/v1/courses/:id | Cập nhật gói tập. (implemented) |
| DELETE | /api/v1/courses/:id | Xoa goi tap. (implemented) |

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

## VI. COLLECTION: Sessions

Hệ thống hiện đã có collection `sessions` để quản lý lịch tập theo slot giờ cụ thể:

- **_id**: ObjectId.
- **branchId**: Ref (Branches).
- **trainerId**: Ref (Employees).
- **courseLevel**: Enum (Basic, Advanced, Professional).
- **scheduledAt**: Date (ngày + giờ bắt đầu).
- **durationMin**: Int.
- **capacity**: Int (sĩ số tối đa).
- **enrolledCount**: Int (số học viên đã đăng ký).
- **enrolledSubscriptionIds**: Array[ObjectId] (danh sách subscription đã đăng ký).
- **tags**: Array[String] (ràng buộc course/subscription phù hợp).

`Attendance` có thể gắn `sessionId` khi check-in theo lịch lớp.

---

## VII. CÁC RÀNG BUỘC BỔ SUNG

- **Không tạo 2 thẻ Active cùng lúc:** Planned/future hardening. Trước khi tạo Subscription mới,
  kiểm tra member đã có thẻ nào `status = Active` chưa. Nếu có → yêu cầu kết thúc hoặc hoàn tiền
  thẻ cũ trước.
- **Refund chỉ do Manager thực hiện:** Endpoint `/refund` yêu cầu role Manager và ghi log lý do.
- **Audit log:** Các thao tác quan trọng (tạo thẻ, hoàn tiền, bảo lưu) nên ghi vào một collection `AuditLogs` riêng với `{ actorId, action, targetId, timestamp, metadata }` để phục vụ tra cứu khi có tranh chấp.
