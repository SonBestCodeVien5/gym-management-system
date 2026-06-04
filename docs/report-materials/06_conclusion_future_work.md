# KẾT LUẬN & HƯỚNG PHÁT TRIỂN

---

## I. KẾT QUẢ HIỆN TẠI VÀ MỤC TIÊU

Phần backend hiện tại đã giải quyết được nhiều vấn đề cốt lõi ở mức MVP, đồng thời còn một số mục
tiêu sản phẩm cần phát triển tiếp:

| Vấn đề ban đầu | Giải pháp hệ thống |
|----------------|-------------------|
| Học viên không tập được ở chi nhánh khác | Backend đã có branch, attendance roaming và geo-query tìm chi nhánh gần nhất |
| Gian lận tài khoản ảo để hưởng ưu đãi | Backend đã có định danh duy nhất bằng CCCD; nhận diện khuôn mặt là future work |
| Hoàn tiền/bảo lưu thiếu minh bạch | Backend đã có rule hoàn tiền MVP và bảo lưu; audit log đầy đủ là future work |
| Dữ liệu rời rạc giữa các chi nhánh | MongoDB tập trung, đồng bộ real-time toàn hệ thống |
| Quản lý thủ công, dễ sai sót | Backend đã có validation, role guard, employee management và dashboard MVP; Cron Jobs là future work |

### Các tính năng backend hoàn thành (MVP)

- ✅ Quản lý học viên, chi nhánh, gói tập.
- ✅ Tạo và quản lý thẻ tập với tính năng bảo lưu, hoàn tiền.
- ✅ Điểm danh đa chi nhánh với kiểm soát quota tuần.
- ✅ Báo nghỉ & tập bù theo quy tắc rolling window.
- ✅ Tìm chi nhánh gần nhất bằng GPS.
- ✅ Đăng nhập nhân viên, refresh/logout token, phân quyền role cho route nghiệp vụ.
- ✅ Quản lý lịch lớp/session ở backend.
- ✅ Quản lý nhân viên admin-only: tạo, liệt kê, xem chi tiết, cập nhật, reset password/deactivate.
- ✅ Dashboard/report aggregate APIs ở mức MVP read-only.
- ✅ Staff Portal React/Vite tích hợp live API cho nhân viên.
- ✅ Full-stack Docker Compose và seed demo data phục vụ demo/chấm bài.

### Các phần chưa hoàn thành trong code hiện tại

- Branch-scope authorization chưa áp dụng.
- Cron Jobs, notification, online payment, face recognition và audit log nâng cao chưa triển khai.
- Member App riêng cho học viên chưa triển khai.
- Report export PDF/CSV và dashboard analytics nâng cao chưa triển khai.

---

## II. HẠN CHẾ HIỆN TẠI

- **Nhận diện khuôn mặt chưa tích hợp đầy đủ:** MVP sử dụng CCCD là phương thức định danh chính. Nhận diện khuôn mặt được đề xuất nhưng cần tích hợp thêm thư viện/dịch vụ bên thứ ba (ví dụ: FaceIO, AWS Rekognition).
- **Chưa có hệ thống thanh toán trực tuyến:** Hiện tại ghi nhận thanh toán thủ công. Chưa tích hợp cổng thanh toán (VNPay, MoMo...).
- **Chưa có thông báo đẩy (Push Notification):** Cron/notification là future work, chưa gửi được
  thông báo thực sự đến thiết bị người dùng.
- **Báo cáo còn cơ bản:** Dashboard/report hiện là MVP read-only, chưa có export PDF/CSV, scheduled
  reports, doanh thu nâng cao, tỷ lệ gia hạn hoặc retention rate.
- **Employee management còn giới hạn theo scope MVP:** Endpoint CRUD nhân sự đã có nhưng đang
  admin-only; chưa có branch-scope permission chi tiết theo từng chi nhánh.
- **Branch-scope authorization chưa có:** Role guard đã có, nhưng chưa ràng buộc quyền theo từng
  chi nhánh.

---

## III. HƯỚNG PHÁT TRIỂN TƯƠNG LAI

### Ngắn hạn (3–6 tháng)
- **Tích hợp cổng thanh toán:** Kết nối VNPay hoặc MoMo để học viên thanh toán online, tự động kích hoạt thẻ sau khi giao dịch thành công.
- **Push Notification:** Tích hợp Firebase Cloud Messaging (FCM) để gửi nhắc nhở hết buổi, nhắc lịch tập.
- **Nhận diện khuôn mặt:** Tích hợp dịch vụ Face Recognition bên thứ ba cho kiosk check-in tự động tại cửa vào.
- **Report export:** Thêm export PDF/CSV và scheduled reports cho dashboard hiện tại.
- **Frontend regression suite:** Bổ sung Playwright suite lâu dài cho các luồng Staff Portal chính.

### Trung hạn (6–12 tháng)
- **Nâng cấp đặt lịch theo slot:** Mở rộng module `Sessions` hiện có với waitlist, hủy đặt chỗ,
  conflict lịch trainer và capacity nâng cao.
- **Analytics nâng cao:** Dashboard báo cáo doanh thu, tỷ lệ gia hạn, học viên mới/rời đi theo tháng, hiệu suất từng chi nhánh.
- **Ứng dụng mobile native:** Phát triển app iOS/Android bằng React Native để có trải nghiệm tốt hơn web responsive.
- **Loyalty Program:** Hệ thống tích điểm và đổi quà cho học viên thân thiết.

### Dài hạn (1 năm+)
- **Chuyển sang Microservices:** Tách các module (Attendance, Finance, Notification) thành service độc lập khi quy mô hệ thống tăng trưởng.
- **API mở cho đối tác:** Cung cấp API để tích hợp với các ứng dụng sức khỏe (Apple Health, Google Fit) hoặc thiết bị wearable.
- **AI/ML:** Gợi ý gói tập phù hợp dựa trên lịch sử tập luyện, dự báo tỷ lệ churn (bỏ gym) để can thiệp sớm.

---

## IV. KẾT LUẬN

Dự án hệ thống quản lý gym đa chi nhánh được xây dựng với trọng tâm là giải quyết các bài toán nghiệp vụ thực tế: định danh học viên chống gian lận, linh hoạt tập luyện đa chi nhánh, và minh bạch hóa các giao dịch tài chính.

Việc lựa chọn **Go + MongoDB** phù hợp với yêu cầu hiệu năng cao và tính linh hoạt của dữ liệu. Thiết kế hệ thống từ đầu đã tính đến khả năng mở rộng — từ cấu trúc module rõ ràng ở backend, đến việc sử dụng index địa lý và Int64 cho tài chính.

Phiên bản hiện tại cũng đã có Staff Portal, Docker Compose full stack và seed demo data, giúp hệ
thống có thể chạy, demo và chụp minh chứng báo cáo bằng một quy trình rõ ràng. Đây là nền tảng đủ
vững chắc để phát triển thành một sản phẩm thực tế, với lộ trình mở rộng rõ ràng từ MVP đến hệ thống
hoàn chỉnh hơn trong tương lai.

---

*Tài liệu được thực hiện bởi: Nguyễn Khánh Sơn — MSSV: 20235212*
