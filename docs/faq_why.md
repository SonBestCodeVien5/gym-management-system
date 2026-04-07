# Q&A - Tại sao lại thiết kế như thế này?

Tài liệu này lưu trữ các quyết định thiết kế kiến trúc và nghiệp vụ cốt lõi (Architecture Decision Records).

### 1. Tại sao lại dùng `Int64` cho mọi biến liên quan đến tiền tệ (`UnitPrice`, `BasePrice`) thay vì `Float64`?
**Trả lời:** Để triệt tiêu hoàn toàn sai số dấu phẩy động (Floating-point precision error) của máy tính. Trong logic hoàn tiền lẻ buổi (Ví dụ: Hoàn 20% của số buổi còn lại), việc tính toán trên số nguyên (đơn vị VNĐ) đảm bảo tính chính xác tuyệt đối về mặt kế toán.

### 2. Tại sao lại tách file `go.mod` và chạy MongoDB qua Docker Volume?
**Trả lời:** * `go.mod` dùng để khóa cứng phiên bản thư viện, tránh việc code chạy được trên máy này nhưng lỗi trên máy khác/server.
* Docker Volume đóng vai trò là "bảo hiểm dữ liệu". Khi Container MongoDB bị tắt hoặc xóa (`docker compose down`), dữ liệu học viên vẫn được giữ lại an toàn trên máy chủ Linux.

### 3. Tại sao luật bảo lưu lại lưu thành mảng `SuspensionHistory` thay vì chỉ cập nhật ngày hết hạn?
**Trả lời:** Để kiểm soát chặt chẽ luật "Bảo lưu 365 ngày". Lưu mảng giúp truy vết được:
1. Học viên đã nghỉ bao nhiêu lần.
2. Tổng số ngày đã dùng.
3. Số buổi bị "đóng băng" tại mỗi thời điểm để phục vụ logic cộng dồn hoặc hoàn tiền nếu xảy ra tranh chấp.

### 4. Tại sao điểm danh lại tách một trường riêng là `IsMakeupFor` (Con trỏ thời gian)?
**Trả lời:** Phục vụ luật "Tập bù trong 7 ngày". 
Khi học viên báo nghỉ hợp lệ (Reported_Missed), hệ thống ghi nhận. Khi họ đi tập bù, một bản ghi điểm danh mới (`Status: Makeup`) được tạo ra và `IsMakeupFor` sẽ trỏ về đúng ngày họ đã báo vắng. Điều này giúp không làm rối lệch bộ đếm `RemainingSessions` (vì buổi vắng đã trừ buổi, buổi bù không trừ thêm).

### 5. Tại sao không cho phép báo nghỉ/tập bù khi thẻ đang `Suspended`?
**Trả lời:** Trạng thái `Suspended` (Đóng băng/Bảo lưu) nghĩa là toàn bộ hợp đồng tạm dừng. Việc phát sinh bản ghi điểm danh (Attendance) trong khoảng thời gian này vi phạm quy tắc vận hành và gây lỗi cho bộ đếm 365 ngày.