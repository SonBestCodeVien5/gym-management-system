# PHASE 2: THIẾT KẾ HỆ THỐNG

# I. KIẾN TRÚC DỮ LIỆU (DATABASE SCHEMA - MONGODB)

Hệ thống sử dụng MongoDB để lưu trữ dữ liệu linh hoạt, tối ưu hóa hiệu năng bằng cách kết hợp giữa Tham chiếu (Reference) và Nhúng (Embedding).

### 1. Collection: Branches (Cơ sở/Chi nhánh)

Quản lý mạng lưới đa điểm và phục vụ tính năng tìm kiếm vị trí (Roaming).

- **_id**: ObjectId (Primary Key).
- **branchCode**: String (Unique Index - ví dụ: "HN01", "HCM02").
- **name**: String (Tên cơ sở).
- **address**: String (Địa chỉ chi tiết).
- **province**: String (Tỉnh/Thành phố).
- **location**: `{ type: "Point", coordinates: [Longitude, Latitude] }` (Chuẩn GeoJSON để tìm kiếm chi nhánh gần nhất).
- **managerId**: Ref (Employees) (ID của quản lý chi nhánh).

### 2. Collection: Members (Hồ sơ học viên)

- **_id**: ObjectId.
- **ccid**: String (Unique Index - Định danh bằng CCCD để tránh gian lận tài khoản ảo).
- **fullName**: String.
- **email**: String.
- **phone**: String.
- **gender**: String.
- **level**: String.
- **isRegistered**: Boolean (True sau khi hoàn thành khóa đầu để hưởng ưu đãi 10% cho lần sau).
- **isSuspended**: Boolean (True sau khi học viên đã thực hiện bảo lưu).

### 3. Collection: Employees (Nhân sự đa năng)

- **_id**: ObjectId.
- **employeeId**: String (Mã nhân viên duy nhất).
- **fullName**: String.
- **role**: Array[String] (Manager, Trainer, Receptionist - hỗ trợ đa vai trò).
- **level**: Enum (Basic, Advanced, Professional - dùng để khớp với trình độ gói tập).
- **phone**: String.
- **email**: String.
- **branchId**: Array[ObjectId] (Danh sách các chi nhánh nhân viên đang làm việc).

### 4. Collection: Courses (Gói tập mẫu)

- **_id**: ObjectId.
- **title**: String (Tên gói tập).
- **level**: Enum (Basic, Advanced, Professional).
- **basePrice**: Int64 (Giá gốc chưa chiết khấu - đơn vị VNĐ).
- **sessionCount**: Int (Số buổi định danh của gói).
- **Description**: String.

### 5. Collection: Subscriptions (Thẻ tập - Trọng tâm logic)

- **ID**: String.
- **memberId**: Ref (Members).
- **courseId**: Ref (Courses).
- **homeBranchId**: Ref (Branches).
- **status**: Enum (Pending, Active, Suspended, Expired, Refunded).

**Tài chính (Độ chính xác tuyệt đối):**

- **paymentDate**: Date (nullable, chi co sau khi confirm offline).
- **totalAmountPaid**: Int64 (Tiền thực thu sau chiết khấu).
- **unitPrice**: Int64 (Giá trị 01 buổi tập = totalAmountPaid / totalSessions - dùng để hoàn tiền).

**Vận hành:**

- **totalSessions / remainingSessions**: Int.
- **sessionPerWeek**: Int (Hạn mức tập luyện hàng tuần).
- **suspension**: `{ startDate, endDate, frozenSession, reason }` (Thông tin bảo lưu hiện tại).

### 6. Collection: Attendance (Nhật ký điểm danh)

- **ID**: String.
- **subId**: Ref (Subscriptions).
- **branchId**: Ref (Branches).
- **date**: Date.
- **status**: Enum (Attended, Absent, Reported_Missed, Makeup).
- **isMakeupFor**: Date (Nếu status là "Makeup", trỏ về ngày đã báo nghỉ hợp lệ).

---

## II. SƠ ĐỒ THỰC THỂ QUAN HỆ (ERD) - CẤU TRÚC LIÊN KẾT

- **Member (1) - Subscription (N):** Lịch sử tham gia của một học viên.
- **Course (1) - Subscription (N):** Các thẻ tập được tạo từ khung chương trình mẫu.
- **Employee (1) - Branch (N):** Một nhân viên có thể hỗ trợ/quản lý tại nhiều cơ sở.
- **Branch (1) - Subscription/Attendance (N):** Quản lý học viên đăng ký gốc và lưu lượng check-in thực tế.
- **Subscription (1) - Attendance (N):** Theo dõi lịch sử thực thi của thẻ tập.

*Hình 1: Sơ đồ usecase về các chức năng nghiệp vụ*

---

## III. CHI TIẾT CÁC THUẬT TOÁN NGHIỆP VỤ (BUSINESS LOGIC)

### 1. Thuật toán Hoàn tiền (Financial Precision)

Hệ thống sử dụng kiểu dữ liệu Int64 để triệt tiêu sai số dấu phẩy động:

- **Implemented MVP:** `refund_amount = total_amount_paid * remaining_sessions / total_sessions`.
  Sau khi hoàn tiền, subscription chuyển sang `refunded` và `remaining_sessions = 0`.
- **Business target/planned extension:** Hoàn 50% trong 72h đầu nếu chưa tập; hoàn 20% giá trị còn
  lại nếu đã tập dưới hoặc bằng 50% số buổi; không hoàn nếu tập > 50%.

### 2. Logic Báo nghỉ & Tập bù (Sliding Window)

- **Báo nghỉ (Reported_Missed):** Chỉ được chấp nhận nếu trong 30 ngày gần nhất (rolling window) chưa có bản ghi báo nghỉ nào khác.
- **Tập bù (Makeup):** Chỉ được kích hoạt nếu tồn tại bản ghi `Reported_Missed` trong vòng 07 ngày trước đó. Khi tập bù thành công, status bản ghi mới là `Makeup` và `isMakeupFor` trỏ về ngày nghỉ gốc.

### 3. Logic Phân cấp chuyên môn (Level-based Access)

- Khi gán một Employee vào Attendance với vai trò Trainer, hệ thống kiểm tra: `Employee.level >= Course.level`. HLV cấp Professional có thể dạy các gói Basic/Advanced, nhưng không có chiều ngược lại.

### 4. Logic Tự động hóa (Cron Jobs)

- **Quét hết hạn bảo lưu:** Nếu `currentDate > endDate` => Chuyển trạng thái thẻ sang `Expired` và reset `remainingSessions = 0`.

---

## IV. DANH SÁCH API ENDPOINTS (DÀNH CHO BACKEND GO)

| Method | Endpoint | Mô tả nghiệp vụ |
|--------|----------|-----------------|
| POST | /api/v1/members | Đăng ký hồ sơ học viên mới (Check CCID). |
| GET | /api/v1/members/:id | Xem hồ sơ học viên. |
| PATCH | /api/v1/members/:id/activate | Confirm offline payment -> subscription pending to active. |
| POST | /api/v1/subscriptions | Tạo thẻ tập (status = pending). |
| GET | /api/v1/subscriptions/:id | Xem chi tiết thẻ tập. |
| PATCH | /api/v1/subscriptions/:id/suspend | Gửi yêu cầu bảo lưu (Cập nhật frozenSession). |
| PATCH | /api/v1/subscriptions/:id/unsuspend | Kích hoạt lại thẻ sau bảo lưu. |
| PATCH | /api/v1/subscriptions/:id/expire | Hết hạn thủ công. |
| POST | /api/v1/attendance/checkin | Điểm danh (Tự động trừ remainingSessions). |
| GET | /api/v1/subscriptions/:id/attendance | Lịch sử điểm danh theo thẻ tập. |
| POST | /api/v1/attendance/report | Báo nghỉ có phép (Kiểm tra luật 30 ngày). |
| POST | /api/v1/attendance/makeup | Đăng ký buổi tập bù (Kiểm tra 7 ngày). |
| GET | /api/v1/branches/nearby | Tìm chi nhánh theo tọa độ GPS. |
| POST | /api/v1/subscriptions/:id/refund | Yêu cầu hoàn tiền theo số buổi còn lại. |
| POST | /api/v1/sessions | Tạo lịch lớp. |
| GET | /api/v1/sessions | Lọc/lấy danh sách lịch lớp. |
| GET | /api/v1/sessions/:id | Xem chi tiết lịch lớp. |
| POST | /api/v1/sessions/:id/enroll | Đăng ký subscription vào lịch lớp. |
| POST | /api/v1/sessions/:id/checkin | Check-in theo lịch lớp. |
| POST | /api/v1/auth/login | Đăng nhập nhân viên. |
| POST | /api/v1/auth/refresh | Rotate refresh token và cấp lại access token. |
| POST | /api/v1/auth/logout | Hủy refresh token. |

---

## V. RÀNG BUỘC HỆ THỐNG (SYSTEM CONSTRAINTS)

- **Unique CCID:** Ngăn chặn việc tạo nhiều tài khoản để hưởng ưu đãi "Người mới" liên tục và cho phép tập liên chi nhánh.
- **Frozen Status:** Khi `Subscription.status` là `"Suspended"`, mọi API điểm danh/tập bù cho thẻ đó đều bị khóa.
- **Geo-Fencing:** Giao diện Frontend có thể dùng GeoLocation để chỉ cho phép điểm danh khi học viên đang ở bán kính < 500m so với tọa độ chi nhánh.
- **Lịch tập:** Chọn cố định 3 đến 6 buổi trên 1 tuần.
- **Ưu đãi:** Biến `isRegistered` phải được kiểm tra trước khi tính `totalPaid` cho Subscription mới.
