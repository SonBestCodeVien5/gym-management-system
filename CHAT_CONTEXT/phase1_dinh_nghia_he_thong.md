# PHASE 1: ĐỊNH NGHĨA & PHÂN TÍCH YÊU CẦU HỆ THỐNG

---

## I. BỐI CẢNH & ĐẶT VẤN ĐỀ

Thị trường phòng gym tư nhân tại Việt Nam đang phát triển mạnh theo mô hình đa chi nhánh. Tuy nhiên, phần lớn các hệ thống quản lý hiện tại vẫn hoạt động rời rạc giữa các cơ sở, gây ra nhiều bất cập:

- Học viên không thể tập tại chi nhánh khác ngoài nơi đăng ký gốc.
- Dữ liệu điểm danh, thẻ tập, và tài chính được lưu thủ công hoặc trên các phần mềm không đồng bộ.
- Không có cơ chế kiểm soát gian lận tài khoản (một người lập nhiều tài khoản để hưởng ưu đãi "Người mới").
- Quy trình hoàn tiền, bảo lưu thiếu minh bạch và dễ xảy ra tranh chấp.

**Mục tiêu dự án:** Xây dựng hệ thống quản lý gym đa chi nhánh tích hợp, cho phép học viên tập linh hoạt tại nhiều cơ sở, tự động hóa các nghiệp vụ tài chính và điểm danh, đồng thời cung cấp công cụ quản lý tập trung cho đội ngũ vận hành.

---

## II. MÔ HÌNH VẬN HÀNH & ĐỐI TƯỢNG SỬ DỤNG

### 1. Mô hình kinh doanh

- **Loại hình:** Private Gym đa cơ sở (Multi-branch).
- **Phạm vi:** Nhiều chi nhánh trên nhiều tỉnh/thành phố, dùng chung một hệ thống dữ liệu trung tâm.
- **Tính năng đặc trưng:** Roaming — học viên đăng ký tại chi nhánh A vẫn có thể check-in tại chi nhánh B ở tỉnh khác bằng cùng một thẻ tập.

### 2. Phân cấp học viên (Member Level)

| Level | Điều kiện | Đặc quyền |
|-------|-----------|-----------|
| Cơ bản (Basic) | Dưới 1 năm tập luyện | Tham gia các lớp Basic |
| Nâng cao (Advanced) | Trên 1 năm tập luyện | Tham gia lớp Basic + Advanced |
| Chuyên nghiệp (Professional) | Định hướng thi đấu | Tham gia tất cả các lớp |

### 3. Các Actor trong hệ thống

| Actor | Mô tả |
|-------|-------|
| **Học viên (Member)** | Người đăng ký và sử dụng dịch vụ tập luyện. |
| **Lễ tân (Receptionist)** | Nhân viên tiếp đón, tạo thẻ tập, điểm danh, xử lý báo nghỉ. |
| **Huấn luyện viên (Trainer)** | Phụ trách buổi tập, chỉ được dạy lớp tương ứng với cấp độ chuyên môn. |
| **Quản lý (Manager)** | Quản lý toàn bộ chi nhánh, xem báo cáo, xử lý hoàn tiền. |
| **Hệ thống tự động (Cron)** | Thực thi các tác vụ định kỳ (hết hạn thẻ, nhắc gia hạn...). |

---

## III. YÊU CẦU CHỨC NĂNG (FUNCTIONAL REQUIREMENTS)

### FR-01: Quản lý học viên
- FR-01.1: Đăng ký hồ sơ học viên mới với xác thực CCCD (duy nhất toàn hệ thống).
- FR-01.2: Định danh học viên bằng CCCD và nhận diện khuôn mặt để đảm bảo nguyên tắc "1 thẻ - 1 người".
- FR-01.3: Xem và cập nhật thông tin cá nhân.

### FR-02: Quản lý thẻ tập (Subscription)
- FR-02.1: Tạo thẻ tập mới từ danh sách gói tập có sẵn, ghi nhận số buổi và giá trị.
- FR-02.2: Áp dụng ưu đãi 10% cho học viên đã hoàn thành ít nhất một khóa trước đó (`isRegistered = true`).
- FR-02.3: Cho phép bảo lưu thẻ tập khi học viên có lý do hợp lệ, ghi đông băng số buổi còn lại.
- FR-02.4: Tự động chuyển trạng thái thẻ sang `Expired` khi hết buổi hoặc hết thời hạn bảo lưu.

### FR-03: Điểm danh & Lịch tập
- FR-03.1: Ghi nhận điểm danh theo thẻ tập, tự động trừ số buổi còn lại.
- FR-03.2: Giới hạn số buổi tập trong tuần theo `sessionPerWeek` (3–6 buổi, không tính buổi tập bù).
- FR-03.3: Hỗ trợ Roaming — học viên check-in tại bất kỳ chi nhánh nào trong hệ thống.
- FR-03.4: Tùy chọn Geo-fencing — chỉ cho phép check-in khi học viên ở trong bán kính 500m so với chi nhánh.

### FR-04: Báo nghỉ & Tập bù
- FR-04.1: Học viên được báo nghỉ tối đa 1 lần trong mỗi cửa sổ 30 ngày (rolling window, không cộng dồn).
- FR-04.2: Buổi tập bù phải được thực hiện trong vòng 7 ngày kể từ ngày báo nghỉ hợp lệ.
- FR-04.3: Buổi tập bù vẫn tính vào quota tuần.

### FR-05: Hoàn tiền
- FR-05.1: Hoàn 50% tổng tiền nếu yêu cầu trong 72 giờ đầu và chưa tập buổi nào.
- FR-05.2: Hoàn 20% giá trị còn lại nếu đã tập từ 1 buổi đến ≤ 50% tổng số buổi.
- FR-05.3: Không hoàn tiền nếu đã tập trên 50% tổng số buổi.
- FR-05.4: Không áp dụng hoàn tiền khi thẻ đang trong trạng thái bảo lưu.

### FR-06: Tìm kiếm chi nhánh
- FR-06.1: Tra cứu danh sách chi nhánh theo tỉnh/thành phố.
- FR-06.2: Tìm chi nhánh gần nhất dựa trên tọa độ GPS của học viên.

### FR-07: Quản lý nhân sự
- FR-07.1: Tạo và quản lý hồ sơ nhân viên với đa vai trò (Manager, Trainer, Receptionist).
- FR-07.2: Gán nhân viên vào một hoặc nhiều chi nhánh.
- FR-07.3: Kiểm tra cấp độ chuyên môn của Trainer trước khi phân công lớp học.

### FR-08: Báo cáo & Thống kê (Manager)
- FR-08.1: Xem thống kê check-in theo ngày/tuần/tháng của từng chi nhánh.
- FR-08.2: Xem danh sách học viên sắp hết buổi (còn ≤ 3 buổi).
- FR-08.3: Tra cứu lịch sử hoàn tiền và bảo lưu.

---

## IV. YÊU CẦU PHI CHỨC NĂNG (NON-FUNCTIONAL REQUIREMENTS)

### Hiệu năng
- API phản hồi dưới 300ms với các thao tác thường ngày (check-in, tra cứu thẻ).
- Truy vấn tìm chi nhánh gần nhất (Geo-query) dưới 100ms nhờ chỉ mục `2dsphere`.

### Bảo mật
- Xác thực bằng JWT (Access Token 15 phút, Refresh Token 7 ngày).
- Mã hóa dữ liệu nhạy cảm (CCCD, thông tin thanh toán) khi lưu trữ.
- Phân quyền chặt chẽ theo role và phạm vi chi nhánh.
- Ghi Audit Log cho mọi thao tác tài chính quan trọng.

### Khả năng mở rộng
- Kiến trúc có thể thêm chi nhánh mới mà không cần thay đổi codebase.
- MongoDB cho phép scale horizontal khi dữ liệu tăng trưởng.

### Độ tin cậy
- Tính toán tài chính dùng kiểu `Int64` để triệt tiêu sai số dấu phẩy động.
- Cron Job chạy độc lập, không ảnh hưởng đến luồng xử lý chính.

### Khả năng sử dụng
- Giao diện lễ tân tối ưu cho thao tác nhanh (check-in dưới 3 bước).
- Học viên có thể tự xem lịch sử và trạng thái thẻ qua ứng dụng.

---

## V. QUY TRÌNH NGHIỆP VỤ CỐT LÕI

### Quy trình Đăng ký & Kích hoạt thẻ

```
Học viên cung cấp CCCD
    → Hệ thống kiểm tra CCCD chưa tồn tại
    → Lễ tân tạo hồ sơ Member
    → Chọn gói tập & tính tiền (áp dụng ưu đãi nếu isRegistered = true)
    → Xác nhận thanh toán
    → Tạo Subscription (status = Active, ghi số buổi & unitPrice)
    → Gửi email xác nhận cho học viên
```

### Quy trình Check-in

```
Học viên xuất trình thẻ/CCCD/khuôn mặt
    → Xác định Subscription Active
    → Kiểm tra quota tuần (count < sessionPerWeek)
    → Kiểm tra Geo-fencing (nếu bật)
    → Ghi Attendance (status = Attended)
    → Trừ remainingSessions - 1
    → Nếu remainingSessions = 0 → Chuyển status = Expired
```

### Quy trình Báo nghỉ & Tập bù

```
Báo nghỉ:
    Học viên gửi yêu cầu báo nghỉ
        → Kiểm tra: trong 30 ngày gần nhất có bản ghi Reported_Missed chưa?
        → Chưa có → Tạo Attendance (status = Reported_Missed), KHÔNG trừ buổi
        → Đã có → Từ chối, thông báo lý do

Tập bù:
    Học viên check-in buổi tập bù
        → Kiểm tra: có bản ghi Reported_Missed trong 7 ngày qua không?
        → Có → Tạo Attendance (status = Makeup, isMakeupFor = ngày nghỉ gốc), trừ buổi
        → Không → Từ chối
```
