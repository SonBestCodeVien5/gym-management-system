## Cập nhật thiết kế: Thêm Collection Sessions (Lịch học)

### Bối cảnh
Hệ thống hiện tại chưa kết nối được Trainer với Member, và tính năng
Roaming chỉ dừng lại ở mức "check-in được ở chi nhánh khác" mà chưa
có giá trị thực tế cho học viên. Để giải quyết, bổ sung collection
Sessions đóng vai trò là **lịch học của phòng gym**.

### Sessions là gì?
Sessions là thời khóa biểu do phòng gym đăng trước. Mỗi Session đại
diện cho một buổi học cụ thể: diễn ra lúc nào, ở đâu, trainer nào
phụ trách, còn bao nhiêu chỗ. Học viên xem lịch và đăng ký chỗ trước
khi đến tập.

Sessions KHÔNG phải là buổi thứ N trong Course của học viên — đó là
Attendance. Sessions là lịch chạy liên tục của gym, độc lập với số
buổi còn lại của từng học viên.

### Schema mới: Sessions
```
Sessions
├── _id           : ObjectId
├── branchId      : Ref(Branches)       — chi nhánh tổ chức
├── trainerId     : Ref(Employees)      — trainer phụ trách
├── courseLevel   : Enum(Basic/Advanced/Professional)
├── scheduledAt   : Date                — ngày + giờ bắt đầu
├── durationMin   : Int                 — thời lượng (phút)
├── capacity      : Int                 — sĩ số tối đa
├── enrolledCount : Int                 — số học viên đã đăng ký
└── tags          : Array[String]       — ví dụ: "yoga", "HIIT"
```

### Thay đổi schema hiện có
Attendance thêm 1 field nullable:
```
sessionId : Ref(Sessions) | null
```
- Có sessionId  = học viên tập theo lịch lớp (đã đăng ký trước)
- null          = học viên tập tự do (giữ nguyên flow hiện tại)

Subscriptions, Members, Courses, Branches giữ nguyên hoàn toàn.

### API cần bổ sung
- GET  /api/v1/sessions?branchId=&level=&date=   — tìm lịch lớp phù hợp
- POST /api/v1/sessions                           — tạo lịch (Trainer/Manager)
- POST /api/v1/sessions/:id/enroll               — học viên đăng ký chỗ
- POST /api/v1/sessions/:id/checkin              — check-in theo lịch lớp

### Tác dụng
- Trainer có vai trò thực sự trong hệ thống, được gắn với buổi dạy cụ thể.
- Học viên roaming sang chi nhánh khác có thể tìm lớp tương tự
  (cùng level, cùng tag) và đăng ký trước khi đến.
- Phòng gym kiểm soát được sĩ số, tránh quá tải.
- Hỗ trợ cả 2 mô hình song song: tập tự do và tập theo lịch lớp.

### Ảnh hưởng đến thiết kế hiện tại
- `Attendance` vẫn là nguồn sự thật cho việc điểm danh và trừ buổi của member.
- `Sessions` chỉ là lịch lớp, không thay thế `Subscription` và không thay thế `Attendance`.
- Khi `sessionId = null`, hệ thống chạy đúng flow cũ: check-in tự do theo subscription.
- Khi `sessionId` có giá trị, attendance phải gắn với lịch lớp cụ thể để quản lý sĩ số và trainer.
- `sessionPerWeek` vẫn là quota của member theo subscription; `capacity` là quota của lớp theo session.
- `reported_missed` và `makeup` vẫn thuộc nghiệp vụ attendance, không chuyển sang `Sessions`.
- `roaming` giờ có ý nghĩa hơn vì học viên có thể xem/lọc lịch lớp ở chi nhánh khác rồi đăng ký chỗ.

### Tác động cần xử lý khi implement
- Bổ sung `sessionId` nullable vào Attendance schema.
- Thêm collection `Sessions` và index phù hợp để tránh trùng lịch/giờ.
- Tách luồng check-in thành 2 mode:
  - mode tự do: không có `sessionId`.
  - mode theo lớp: có `sessionId`, cần kiểm tra capacity và quan hệ branch/trainer/level.
- Quy tắc enrollment cần chốt rõ:
  - chỉ tăng `enrolledCount` khi chỗ còn trống,
  - không cho đăng ký trùng cùng session,
  - có cơ chế hủy nếu sau này cần hoàn chỗ.
- API mới nên được dựng theo thứ tự: list/search sessions -> create session -> enroll -> check-in.

### Gợi ý phạm vi triển khai tối thiểu
1. Làm `GET /sessions` để học viên xem lịch.
2. Làm `POST /sessions` để trainer/manager tạo lịch.
3. Làm `POST /sessions/:id/enroll` để giữ chỗ.
4. Làm `POST /sessions/:id/checkin` hoặc mở rộng attendance check-in với `sessionId`.
5. Chỉ sau đó mới tối ưu search theo tag/branch/level hoặc geo.

### Quan hệ mới trong ERD
- Branches  (1) → (N) Sessions
- Employees (1) → (N) Sessions   — Trainer phụ trách
- Sessions  (1) → (N) Attendance — các học viên tham dự buổi đó