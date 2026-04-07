# Q&A - Quyết định Kiến trúc & Triết lý Golang (ADR)

Tài liệu này lưu trữ các quyết định thiết kế kiến trúc hệ thống (Architecture Decision Records) và giải thích triết lý lập trình trong dự án Gym Management System.

---

## PHẦN 1: TRIẾT LÝ GOLANG & CẤU TRÚC REPOSITORY (CLEAN ARCHITECTURE)

### 1. Tại sao lại dùng Interface mong đợi đầu ra, nhưng lại trả về một Struct Private? (Nguyên lý Duck Typing)
**Vấn đề:** Trong `member_repo.go`, hệ thống định nghĩa một Interface `MemberRepository` (chỉ chứa tên hàm), một Struct viết thường `memberRepoImpl` (ẩn bên trong), và một hàm `New...` để trả Struct đó ra ngoài dưới danh nghĩa Interface.
**Trả lời:** Đây là cách Go thực thi tính **Đa hình (Polymorphism)** và **Bảo mật (Encapsulation)**, gọi là "Thuyết con vịt" (Duck Typing).
* **Bảo mật dữ liệu:** Bằng cách giấu Struct thành private (`memberRepoImpl`), các tầng khác (như API/Service) không thể chọc thẳng vào biến `collection` để gọi các lệnh nguy hiểm (như `Drop()` xóa bảng DB). Chúng chỉ được phép dùng các hàm mà Interface cấp quyền (như `Create`, `GetByID`).
* **Tính linh hoạt (Loosely Coupled):** Các file khác chỉ giao tiếp qua "Bản hợp đồng" (Interface). Nếu sau này dự án đổi từ MongoDB sang PostgreSQL, ta chỉ cần viết một Struct mới thỏa mãn các hàm trong Interface, tráo đổi ở hàm khởi tạo là xong. Toàn bộ code logic của hệ thống không cần sửa đổi dù chỉ một dòng.

### 2. Sự khác biệt giữa `type` và `func` trong Go là gì?
**Trả lời:** * **`type` (Danh từ):** Dùng để định nghĩa một hình hài, cấu trúc dữ liệu (`Struct`) hoặc một bản hợp đồng (`Interface`). Nó không chứa logic thực thi.
* **`func` (Động từ):** Dùng để định nghĩa hành động (Hàm/Method).
Go không có `class` như Java/C++. Khái niệm và Hành động bị tách rời. Go dùng cú pháp Receiver (ví dụ: `func (r *Struct) DoSomething()`) để gắn hành động vào một cấu trúc dữ liệu.

### 3. Tại sao hàm `NewMemberRepository` lại không có `(r *Struct)` ở trước?
**Trả lời:** Cú pháp có `(r *Struct)` đứng trước tên hàm là **Method** (Hành động của một đối tượng đã tồn tại). 
Hàm `NewMemberRepository` mang nhiệm vụ **khởi tạo** đối tượng đó, nên nó phải là một **Hàm độc lập (Standalone Function)**. Trong Go không có từ khóa `new` cho Class, nên quy ước chung (convention) là sử dụng Design Pattern **Constructor Function** với tiền tố `New...` để khởi tạo, cấp phát bộ nhớ và trả về đối tượng sẵn sàng hoạt động.

---

## PHẦN 2: THIẾT KẾ CƠ SỞ DỮ LIỆU (MONGODB vs SQL)

### 4. Tại sao lại chọn NoSQL (MongoDB) thay vì SQL cho dự án này?
**Trả lời:** SQL (như MySQL) dùng cấu trúc dạng Bảng (Table), cực kỳ chặt chẽ nhưng kém linh hoạt khi dữ liệu lồng nhau phức tạp. MongoDB được chọn vì 3 bản chất:
1. **Lưu trữ dạng Tài liệu (Document-based):** Cho phép gom nhóm dữ liệu (như lịch sử bảo lưu) thành dạng Mảng (Array) hoặc Object lồng nhau, phù hợp tuyệt đối với cấu trúc của hệ thống Gym.
2. **Lược đồ động (Dynamic Schema):** Khi cần thêm tính năng mới (ví dụ thêm mảng `Role` cho nhân viên, hay `FaceVector` cho học viên), DB tự động tiếp nhận mà không cần khóa bảng (ALTER TABLE) làm gián đoạn hệ thống.
3. **Mở rộng ngang (Horizontal Scaling):** Sẵn sàng cho việc phân mảnh dữ liệu (Sharding) khi chuỗi phòng Gym mở rộng ra nhiều tỉnh thành.

### 5. Tại sao thiết kế lại dùng song song cả Tham chiếu (Reference) và Nhúng (Embedded)?
**Trả lời:** Đây là cách tối ưu hóa hiệu năng đọc của MongoDB (Nguyên lý: "Data that is accessed together, lives together").
* **Nhúng (Embedded) cho `Suspension`:** Lịch sử bảo lưu phụ thuộc hoàn toàn vào thẻ tập (`Subscription`). Việc nhúng thẳng vào trong giúp hệ thống chỉ cần 1 thao tác đọc là lấy được toàn bộ trạng thái hiện tại, triệt tiêu các lệnh JOIN đắt đỏ.
* **Tham chiếu (Reference) cho `MemberID` và `CourseID`:** Hội viên và Gói tập là các thực thể độc lập, dùng chung ở nhiều thẻ tập khác nhau. Việc tham chiếu giúp tránh phình to CSDL và ngăn chặn lỗi bất đồng bộ khi cập nhật thông tin.

### 6. Tại sao ID lại dùng `primitive.ObjectID` thay vì số tự tăng (Auto-increment)?
**Trả lời:** Số tự tăng hoạt động tốt trên 1 máy chủ SQL đơn lẻ, nhưng sẽ xung đột trên hệ thống phân tán nhiều máy chủ của NoSQL. `ObjectID` là chuỗi thập lục phân 12-byte do chính thuật toán của MongoDB sinh ra (bao gồm cả timestamp). Nó đảm bảo tính độc nhất vô nhị toàn cầu và tăng tốc độ truy vấn qua Index mà không cần cơ chế khóa (locking).

---

## PHẦN 3: LOGIC NGHIỆP VỤ & KIỂU DỮ LIỆU (MODELS)

### 7. Tại sao lại dùng kiểu `Int64` cho mọi biến tiền tệ (`UnitPrice`, `BasePrice`)?
**Trả lời:** Để triệt tiêu sai số dấu phẩy động (Floating-point precision error) vốn có của máy tính. Các phép tính hoàn tiền (ví dụ hoàn 20% thẻ tập) được xử lý trên số nguyên gốc (đơn vị VNĐ nhỏ nhất), đảm bảo tính toàn vẹn và chính xác tuyệt đối trong kế toán tài chính.

### 8. Tag `omitempty` trong BSON/JSON có tác dụng gì?
**Trả lời:** Từ khóa này báo cho Go compiler biết: "Nếu trường dữ liệu này rỗng hoặc mang giá trị mặc định (nil, 0, chuỗi rỗng), hãy tự động bỏ qua, KHÔNG lưu nó xuống Database và KHÔNG xuất ra file JSON". Điều này giúp tiết kiệm đáng kể dung lượng lưu trữ (Storage) và băng thông mạng (Bandwidth).

### 9. Tại sao trường `Suspension` (trong thẻ tập) và `IsMakeupFor` (trong điểm danh) lại dùng kiểu con trỏ `*`?
**Trả lời:**
* Nếu không dùng con trỏ, các biến thời gian hoặc Struct luôn chiếm bộ nhớ với giá trị rác mặc định (ví dụ thời gian là `0001-01-01`).
* Dùng con trỏ `*`, nếu học viên không bảo lưu hoặc buổi đó không phải tập bù, biến này sẽ mang giá trị `nil` (null). Kết hợp với `omitempty`, Database sẽ hoàn toàn bỏ qua các trường thừa thãi này, giúp Data cực kỳ "sạch" và code kiểm tra logic `if field != nil` ngắn gọn hơn rất nhiều.

### 10. Tại sao phải lưu trữ mảng `SuspensionHistory` thay vì chỉ thay đổi biến ngày hết hạn?
**Trả lời:** Để kiểm soát chặt chẽ luật "Bảo lưu tối đa 365 ngày" của doanh nghiệp. Việc lưu mảng lịch sử giúp hệ thống có bằng chứng (Audit trail) truy vết số ngày đã nghỉ ở từng đợt, số buổi bị đóng băng, làm cơ sở giải quyết tranh chấp hoàn tiền hoặc cộng dồn chu kỳ bảo lưu tiếp theo.