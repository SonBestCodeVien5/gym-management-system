# KẾT LUẬN & HƯỚNG PHÁT TRIỂN

---

## I. KẾT QUẢ DỰ KIẾN

Sau khi hoàn thành, hệ thống sẽ giải quyết được các vấn đề cốt lõi đặt ra ban đầu:

| Vấn đề ban đầu | Giải pháp hệ thống |
|----------------|-------------------|
| Học viên không tập được ở chi nhánh khác | Tính năng Roaming với Geo-query tìm chi nhánh gần nhất |
| Gian lận tài khoản ảo để hưởng ưu đãi | Định danh duy nhất bằng CCCD + nhận diện khuôn mặt |
| Hoàn tiền/bảo lưu thiếu minh bạch | Công thức cố định, ghi AuditLog mọi thao tác tài chính |
| Dữ liệu rời rạc giữa các chi nhánh | MongoDB tập trung, đồng bộ real-time toàn hệ thống |
| Quản lý thủ công, dễ sai sót | Tự động hóa qua Cron Jobs và validation chặt chẽ ở API |

### Các tính năng hoàn thành (MVP)

- ✅ Quản lý học viên, nhân sự, chi nhánh, gói tập.
- ✅ Tạo và quản lý thẻ tập với tính năng bảo lưu, hoàn tiền.
- ✅ Điểm danh đa chi nhánh với kiểm soát quota tuần.
- ✅ Báo nghỉ & tập bù theo quy tắc rolling window.
- ✅ Tìm chi nhánh gần nhất bằng GPS.
- ✅ Phân quyền đa vai trò (Member, Receptionist, Trainer, Manager).
- ✅ Tự động hóa vòng đời thẻ tập qua Cron Jobs.
- ✅ Giao diện Staff Portal và Member App.

---

## II. HẠN CHẾ HIỆN TẠI

- **Nhận diện khuôn mặt chưa tích hợp đầy đủ:** MVP sử dụng CCCD là phương thức định danh chính. Nhận diện khuôn mặt được đề xuất nhưng cần tích hợp thêm thư viện/dịch vụ bên thứ ba (ví dụ: FaceIO, AWS Rekognition).
- **Chưa có hệ thống thanh toán trực tuyến:** Hiện tại ghi nhận thanh toán thủ công. Chưa tích hợp cổng thanh toán (VNPay, MoMo...).
- **Chưa có thông báo đẩy (Push Notification):** Cron Job nhắc nhở mới ở mức ghi log, chưa gửi được thông báo thực sự đến thiết bị người dùng.
- **Báo cáo còn cơ bản:** Dashboard Manager chưa có phân tích nâng cao (doanh thu theo gói, tỷ lệ gia hạn, retention rate...).
- **Chưa có tính năng đặt lịch theo slot:** Học viên check-in tự do, chưa hỗ trợ đặt trước buổi tập theo khung giờ cụ thể.

---

## III. HƯỚNG PHÁT TRIỂN TƯƠNG LAI

### Ngắn hạn (3–6 tháng)
- **Tích hợp cổng thanh toán:** Kết nối VNPay hoặc MoMo để học viên thanh toán online, tự động kích hoạt thẻ sau khi giao dịch thành công.
- **Push Notification:** Tích hợp Firebase Cloud Messaging (FCM) để gửi nhắc nhở hết buổi, nhắc lịch tập.
- **Nhận diện khuôn mặt:** Tích hợp dịch vụ Face Recognition bên thứ ba cho kiosk check-in tự động tại cửa vào.

### Trung hạn (6–12 tháng)
- **Đặt lịch theo slot:** Thêm collection `Sessions`, cho phép học viên đặt trước buổi tập theo khung giờ và HLV cụ thể.
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

Đây là nền tảng đủ vững chắc để phát triển thành một sản phẩm thực tế, với lộ trình mở rộng rõ ràng từ MVP đến hệ thống quy mô lớn hơn trong tương lai.

---

*Tài liệu được thực hiện bởi: Nguyễn Khánh Sơn — MSSV: 20235212*
