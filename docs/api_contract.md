# API Contract (Current)

Cap nhat: 2026-05-05

Muc tieu: chot ten endpoint + request/response co ban de FE va BE dung chung.

---

## API theo collection

### Members
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/members | Tạo hồ sơ học viên mới, check `ccid` unique. | Implemented |
| GET | /api/v1/members/:id | Xem chi tiết hồ sơ học viên. | Implemented |
| PATCH | /api/v1/members/:id/activate | Confirm thanh toán offline và kích hoạt member. | Implemented |
| GET | /api/v1/members/:id/subscriptions | Xem toàn bộ thẻ tập của member. | Planned |

### Subscriptions
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/subscriptions | Tạo thẻ tập mới ở trạng thái `pending`, snapshot quyền từ course. | Implemented |
| GET | /api/v1/subscriptions/:id | Xem thông tin thẻ tập và số buổi còn lại. | Implemented |
| PATCH | /api/v1/subscriptions/:id/suspend | Bảo lưu thẻ tập theo khoảng thời gian. | Implemented |
| PATCH | /api/v1/subscriptions/:id/unsuspend | Kích hoạt lại thẻ sau bảo lưu. | Implemented |
| PATCH | /api/v1/subscriptions/:id/expire | Hết hạn thẻ tập thủ công. | Implemented |
| POST | /api/v1/subscriptions/:id/refund | Tính và xử lý hoàn tiền. | Planned |

### Courses
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/courses | Tạo gói tập mẫu với giá, số buổi và danh sách tag được phép. | Implemented |
| GET | /api/v1/courses | Danh sách gói tập để chọn khi tạo subscription. | Implemented |
| GET | /api/v1/courses/:id | Xem chi tiết một gói tập. | Implemented |
| PATCH | /api/v1/courses/:id | Cập nhật thông tin gói tập và danh sách tag được phép. | Implemented |
| DELETE | /api/v1/courses/:id | Xóa gói tập không còn dùng. | Implemented |

### Branches
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/branches | Tạo chi nhánh mới với location GeoJSON. | Implemented |
| GET | /api/v1/branches | Danh sách chi nhánh để chọn home branch hoặc roaming. | Implemented |
| GET | /api/v1/branches/:id | Xem chi tiết một chi nhánh. | Implemented |
| PATCH | /api/v1/branches/:id | Cập nhật thông tin chi nhánh. | Implemented |
| DELETE | /api/v1/branches/:id | Xóa chi nhánh. | Implemented |
| GET | /api/v1/branches/nearby | Tìm chi nhánh gần vị trí hiện tại. | Planned |

### Attendance
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/attendance/checkin | Ghi nhận check-in tự do hoặc theo session. | Implemented |
| GET | /api/v1/subscriptions/:id/attendance | Xem lịch sử attendance của một subscription. | Implemented |
| POST | /api/v1/attendance/report | Báo nghỉ hợp lệ để mở cửa sổ tập bù. | Planned |
| POST | /api/v1/attendance/makeup | Tạo attendance tập bù từ report đã được duyệt. | Planned |

### Sessions
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/sessions | Tạo lịch lớp do trainer/manager phụ trách. | Implemented |
| GET | /api/v1/sessions | Tìm và lọc lịch lớp theo branch/level/date. | Implemented |
| GET | /api/v1/sessions/:id | Xem chi tiết một session. | Implemented |
| POST | /api/v1/sessions/:id/enroll | Học viên đăng ký chỗ trong session. | Implemented |
| POST | /api/v1/sessions/:id/checkin | Check-in theo session đã đăng ký. | Implemented |

### Auth
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/auth/login | Đăng nhập. | Planned |
| POST | /api/v1/auth/refresh | Cấp lại access token. | Planned |
| POST | /api/v1/auth/logout | Hủy refresh token. | Planned |

---

## Status code mac dinh

- 200: OK
- 201: Created
- 400: Bad request (invalid input)
- 404: Not found
- 409: Conflict (status/logic conflict)
- 500: Internal server error

---

## Notes

| Ghi chú | Nội dung |
|---|---|
| Offline payment | Confirm qua PATCH /members/:id/activate + `subscription_id`. |
| Subscription | Tạo mới ở trạng thái `pending`, sau đó activate khi confirm payment. |
| Sessions MVP | Đã có create/list/get/enroll/checkin; enrollment lưu trên session và check-in tạo attendance có `session_id`. |
| Course tags | `allowed_tags` của course là tập tag được phép dùng để ràng buộc session. |
