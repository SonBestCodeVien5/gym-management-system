# PHASE 3: THIẾT KẾ GIAO DIỆN & LUỒNG NGƯỜI DÙNG

---

> Ghi chú: Tài liệu này mô tả thiết kế UI/UX mục tiêu cho báo cáo. Đối chiếu phạm vi backend
> hiện tại và các phần planned/future nằm ở
> [07_current_implementation_evidence.md](07_current_implementation_evidence.md).

## I. TỔNG QUAN GIAO DIỆN

Hệ thống có **2 giao diện riêng biệt** phục vụ 2 nhóm người dùng:

| Giao diện | Đối tượng | Nền tảng |
|-----------|-----------|----------|
| **Staff Portal** | Lễ tân, Huấn luyện viên, Quản lý | Web (Desktop-first) |
| **Member App** | Học viên | Web responsive / Mobile |

---

## II. STAFF PORTAL

### 1. Màn hình Đăng nhập
- Form nhập `employeeId` + mật khẩu.
- Sau khi đăng nhập, hệ thống tự động xác định role và điều hướng đến dashboard phù hợp.

### 2. Dashboard (Trang tổng quan - Manager)

```
┌─────────────────────────────────────────────────────────┐
│  🏋️ GYM SYSTEM   [Chi nhánh: HCM-01 ▼]    [Admin ▼]   │
├────────────┬────────────────────────────────────────────┤
│            │  Hôm nay: 03/05/2026                       │
│ 📊 Tổng    │  ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│    quan    │  │Check-in  │ │Thẻ Active│ │Sắp hết   │   │
│            │  │  47      │ │  312     │ │buổi: 18  │   │
│ 👥 Học     │  └──────────┘ └──────────┘ └──────────┘   │
│    viên    │                                             │
│            │  Biểu đồ check-in 7 ngày gần nhất          │
│ 🎫 Thẻ    │  ▁▃▅▇▆▄▅                                   │
│    tập     │  T2 T3 T4 T5 T6 T7 CN                      │
│            │                                             │
│ 📋 Điểm   │  Danh sách check-in hôm nay                 │
│    danh    │  [Nguyễn Văn A] [10:32] [Basic] ✓          │
│            │  [Trần Thị B]   [10:45] [Adv]   ✓          │
│ 👔 Nhân   │                                             │
│    sự      │                                             │
│            │                                             │
│ ⚙️ Cài    │                                             │
│    đặt     │                                             │
└────────────┴────────────────────────────────────────────┘
```

### 3. Màn hình Check-in (Lễ tân)

```
┌─────────────────────────────────────────┐
│  CHECK-IN HỌC VIÊN                      │
│                                         │
│  🔍 Tìm kiếm                            │
│  ┌──────────────────────────────────┐   │
│  │ Nhập CCCD / Mã thẻ / Họ tên...  │   │
│  └──────────────────────────────────┘   │
│         [ Hoặc quét khuôn mặt 📷 ]      │
│                                         │
│  Kết quả:                               │
│  ┌──────────────────────────────────┐   │
│  │ 👤 Nguyễn Khánh Sơn             │   │
│  │ Thẻ: SUB-00123 | Advanced       │   │
│  │ Còn lại: 8 buổi                 │   │
│  │ Tuần này: 2/4 buổi              │   │
│  │ Chi nhánh gốc: HCM-01           │   │
│  │                                 │   │
│  │      [ ✅ XÁC NHẬN CHECK-IN ]   │   │
│  └──────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### 4. Màn hình Tạo thẻ tập

```
┌─────────────────────────────────────────────┐
│  TẠO THẺ TẬP MỚI                            │
│                                             │
│  Học viên: [Tìm kiếm theo CCCD...]          │
│  ► Nguyễn Khánh Sơn — CCCD: 0123456789      │
│    isRegistered: ✅ (Giảm 10%)              │
│                                             │
│  Chọn gói tập:                              │
│  ○ Basic    — 20 buổi — 2,000,000đ          │
│  ● Advanced — 20 buổi — 3,500,000đ ★        │
│  ○ Pro      — 20 buổi — 5,000,000đ          │
│                                             │
│  Số buổi/tuần: [ 3 ] [ 4 ] [ 5 ] [●6]      │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │ Tổng: 3,500,000đ                    │   │
│  │ Ưu đãi (10%): - 350,000đ            │   │
│  │ Thực thu: 3,150,000đ                │   │
│  │ Giá/buổi: 157,500đ                  │   │
│  └──────────────────────────────────────┘   │
│                                             │
│  [ HỦY ]              [ ✅ XÁC NHẬN TẠO ]  │
└─────────────────────────────────────────────┘
```

### 5. Màn hình Quản lý học viên

```
┌─────────────────────────────────────────────────────┐
│  DANH SÁCH HỌC VIÊN          [+ Thêm mới]           │
│  🔍 [Tìm kiếm...]   Lọc: [Tất cả ▼] [Active ▼]     │
├──────────────┬──────────┬──────────┬─────────────────┤
│ Họ tên       │ CCCD     │ Level    │ Trạng thái thẻ  │
├──────────────┼──────────┼──────────┼─────────────────┤
│ Nguyễn V. A  │ 012...   │ Advanced │ 🟢 Active (8b)  │
│ Trần Thị B   │ 034...   │ Basic    │ 🟡 Suspended    │
│ Lê Văn C     │ 056...   │ Pro      │ 🔴 Expired      │
└──────────────┴──────────┴──────────┴─────────────────┘
```

---

## III. MEMBER APP (Ứng dụng học viên)

### 1. Màn hình Trang chủ

```
┌─────────────────────────────┐
│  Xin chào, Khánh Sơn 👋     │
│                             │
│  ┌─────────────────────┐    │
│  │  THẺ TẬP CỦA BẠN   │    │
│  │  Advanced Package   │    │
│  │  ████████░░  8/20   │    │
│  │  buổi còn lại       │    │
│  │  HCM-01 | Active 🟢 │    │
│  └─────────────────────┘    │
│                             │
│  Tuần này: 2 / 4 buổi       │
│  ▓▓░░                       │
│                             │
│  [ 📍 Tìm chi nhánh ]       │
│  [ 📅 Lịch sử tập ]         │
│  [ 🔔 Báo nghỉ ]            │
└─────────────────────────────┘
```

### 2. Màn hình Tìm chi nhánh (Roaming)

```
┌─────────────────────────────┐
│  ← TÌM CHI NHÁNH GẦN NHẤT  │
│                             │
│  📍 Vị trí hiện tại:        │
│  Q. Bình Thạnh, TP.HCM      │
│                             │
│  ┌────────────────────────┐ │
│  │    [BẢN ĐỒ]            │ │
│  │   📌HCM-01  📌HCM-02   │ │
│  │         📌HCM-03       │ │
│  └────────────────────────┘ │
│                             │
│  Chi nhánh gần nhất:        │
│  📌 HCM-01 — 0.3 km         │
│     123 Đinh Tiên Hoàng     │
│     [ XEM CHI TIẾT ]        │
│                             │
│  📌 HCM-02 — 1.2 km         │
│     45 Võ Thị Sáu           │
│     [ XEM CHI TIẾT ]        │
└─────────────────────────────┘
```

### 3. Màn hình Lịch sử điểm danh

```
┌─────────────────────────────┐
│  ← LỊCH SỬ TẬP LUYỆN       │
│                             │
│  Tháng 5/2026               │
│  ──────────────────         │
│  ✅ 03/05 — HCM-02 (Adv)   │
│  ✅ 01/05 — HCM-01 (Adv)   │
│                             │
│  Tháng 4/2026               │
│  ──────────────────         │
│  ✅ 29/04 — HCM-01          │
│  🔵 27/04 — Tập bù (Adv)   │
│  📝 25/04 — Báo nghỉ        │
│  ✅ 22/04 — HCM-01          │
└─────────────────────────────┘
```

---

## IV. USER FLOW CHÍNH

### Flow 1: Đăng ký & Mua thẻ tập lần đầu

```
[Học viên đến quầy]
        │
        ▼
[Lễ tân nhập CCCD]
        │
        ├─ CCCD đã tồn tại ──► Thông báo "Tài khoản đã có"
        │
        ▼
[Tạo hồ sơ Member mới]
        │
        ▼
[Chọn gói tập + số buổi/tuần]
        │
        ▼
[Xác nhận & Thanh toán]
        │
        ▼
[Subscription Active — Gửi email xác nhận]
```

### Flow 2: Check-in (Roaming)

```
[Học viên đến chi nhánh khác]
        │
        ▼
[Xuất trình CCCD / khuôn mặt]
        │
        ▼
[Hệ thống tìm Subscription Active]
        │
        ▼
[Kiểm tra quota tuần: count < sessionPerWeek?]
        │
        ├─ Không ──► Từ chối, hiển thị "Đã đủ buổi tuần này"
        │
        ▼
[Geo-fencing: Trong bán kính 500m?] (nếu bật)
        │
        ├─ Không ──► Từ chối, hiển thị vị trí chi nhánh
        │
        ▼
[Ghi Attendance: Attended — Trừ remainingSessions]
        │
        ▼
[Hiển thị xác nhận: còn X buổi]
```

### Flow 3: Hoàn tiền

> Đây là flow nghiệp vụ mục tiêu. Backend hiện tại đã có MVP hoàn tiền theo tỉ lệ số buổi còn lại:
> `refundAmount = totalAmountPaid * remainingSessions / totalSessions`, sau đó chuyển thẻ sang
> `refunded` và đặt `remainingSessions = 0`.

```
[Học viên yêu cầu hoàn tiền]
        │
        ▼
[Manager kiểm tra điều kiện]
        │
        ├─ Trong 72h & chưa tập ──► Hoàn 50% tổng tiền
        │
        ├─ Đã tập 1~50% buổi ──► Hoàn 20% giá trị còn lại
        │                         (remainingSessions × unitPrice × 20%)
        │
        └─ Đã tập > 50% buổi ──► Từ chối hoàn tiền
                │
                ▼
        [Ghi AuditLog — Cập nhật status = Refunded]
```
