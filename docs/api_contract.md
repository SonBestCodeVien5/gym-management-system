# API Contract (Current)

Cap nhat: 2026-05-25

Muc tieu: chot ten endpoint + request/response co ban de FE va BE dung chung.

---

## API theo collection

### Members
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/members | Tạo hồ sơ học viên mới, check `ccid` unique. | Implemented |
| GET | /api/v1/members/:id | Xem chi tiết hồ sơ học viên. | Implemented |
| PATCH | /api/v1/members/:id/activate | Confirm thanh toán offline và kích hoạt member. | Implemented |
| GET | /api/v1/members/:id/subscriptions | Xem toàn bộ thẻ tập của member. | Implemented |

### Subscriptions
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/subscriptions | Tạo thẻ tập mới ở trạng thái `pending`, snapshot quyền + pricing/discount từ course. | Implemented |
| GET | /api/v1/subscriptions/:id | Xem thông tin thẻ tập và số buổi còn lại. | Implemented |
| PATCH | /api/v1/subscriptions/:id/suspend | Bảo lưu thẻ tập theo khoảng thời gian. | Implemented |
| PATCH | /api/v1/subscriptions/:id/unsuspend | Kích hoạt lại thẻ sau bảo lưu. | Implemented |
| PATCH | /api/v1/subscriptions/:id/expire | Hết hạn thẻ tập thủ công. | Implemented |
| POST | /api/v1/subscriptions/:id/refund | Tính và xử lý hoàn tiền theo số buổi còn lại. | Implemented |

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
| GET | /api/v1/branches/nearby | Tìm chi nhánh gần vị trí hiện tại bằng GeoJSON 2dsphere. | Implemented |

### Attendance
| Method | Endpoint | Code làm gì | Trạng thái |
|---|---|---|---|
| POST | /api/v1/attendance/checkin | Ghi nhận check-in tự do hoặc theo session. | Implemented |
| GET | /api/v1/subscriptions/:id/attendance | Xem lịch sử attendance của một subscription. | Implemented |
| POST | /api/v1/attendance/report | Báo nghỉ hợp lệ để mở cửa sổ tập bù. | Implemented |
| POST | /api/v1/attendance/makeup | Tạo attendance tập bù từ report đã được duyệt. | Implemented |

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
| POST | /api/v1/auth/login | Đăng nhập employee/staff và cấp access + refresh token. | Implemented |
| POST | /api/v1/auth/refresh | Rotate refresh token và cấp lại access token. | Implemented |
| POST | /api/v1/auth/logout | Hủy refresh token. | Implemented |

---

## Endpoint details

### Authentication

Public routes:

- `GET /ping`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

All other `/api/v1/*` business routes require:

```http
Authorization: Bearer <access_token>
```

Role guard matrix:

| Surface | Roles |
|---|---|
| Course CRUD | `admin`, `manager` |
| Branch CRUD | `admin`, `manager` |
| Member create/get/list-subscriptions/activate | `admin`, `manager`, `receptionist` |
| Subscription create/get/refund/suspend/unsuspend/expire | `admin`, `manager`, `receptionist` |
| Attendance checkin/report/makeup/history | `admin`, `manager`, `receptionist` |
| Session create/list/get/enroll/checkin | `admin`, `manager`, `trainer` |

Missing, malformed, expired, or inactive-employee access tokens return `401`. Authenticated staff
without an allowed role returns `403`.

Startup creates:

- unique sparse index on `employees.normalized_email`
- unique sparse index on `employees.employee_id`
- unique index on `refresh_tokens.token_hash`

If `BOOTSTRAP_ADMIN_EMAIL` and `BOOTSTRAP_ADMIN_PASSWORD` are configured, startup creates the first
admin account only when the normalized email does not already exist.

### Auth login

`POST /api/v1/auth/login`

Body:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `email` | string | yes | Normalized before lookup. |
| `password` | string | yes | Compared against stored bcrypt password hash. |

Response `200`:

```json
{
  "message": "login successful",
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "employee": {
      "id": "69f20c000c4cd4cdf5768500",
      "employee_id": "ADMIN001",
      "email": "admin@gym.test",
      "full_name": "Gym Admin",
      "role": ["admin"],
      "branch_id": []
    }
  }
}
```

Status codes:

- `200`: credentials accepted and token pair issued.
- `400`: invalid body or missing email/password.
- `401`: invalid credentials or inactive employee.
- `500`: token/storage/internal failure.

Notes:

- The API never returns password hash.
- The refresh token is stored only as a SHA-256 hash for revoke/rotation checks.

### Auth refresh

`POST /api/v1/auth/refresh`

Body:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `refresh_token` | string | yes | Must be signed, unexpired, known, and not revoked. |

Response `200`:

```json
{
  "message": "token refreshed successfully",
  "data": {
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

Status codes:

- `200`: refresh token valid; old refresh token revoked and replacement token pair issued.
- `400`: invalid body or missing refresh token.
- `401`: invalid, expired, revoked, unknown token, or inactive employee.
- `500`: token/storage/internal failure.

### Auth logout

`POST /api/v1/auth/logout`

Body:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `refresh_token` | string | yes | Must be a signed refresh token. |

Response `200`:

```json
{
  "message": "logout successful"
}
```

Status codes:

- `200`: token revoked; repeated logout is idempotent.
- `400`: invalid body or missing refresh token.
- `401`: malformed or unverifiable refresh token.
- `500`: token/storage/internal failure.

### Attendance report

`POST /api/v1/attendance/report`

Body:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `subscription_id` | ObjectID string | yes | Must reference an existing active subscription. |
| `branch_id` | ObjectID string | yes | Attendance branch. |
| `date` | RFC3339 string | no | Missed date; defaults to server time when omitted. |

Response `201`:

```json
{
  "message": "attendance report recorded successfully",
  "data": {
    "id": "69f20c000c4cd4cdf5768500",
    "sub_id": "69f20b22f79bb78cac99aa0a",
    "branch_id": "69f20a180c4cd4cdf57684fe",
    "date": "2026-05-12T08:00:00Z",
    "status": "reported_missed",
    "is_makeup_for": null
  }
}
```

Status codes:

- `201`: reported-missed attendance created.
- `400`: invalid body, ObjectID, or RFC3339 date input.
- `404`: subscription not found.
- `409`: subscription state/expiry conflict or reported-missed 30-day limit reached.
- `500`: internal server error.

Notes:

- Client does not send `status`; handler sets `reported_missed`.
- `reported_missed` does not consume a remaining session.

### Attendance makeup

`POST /api/v1/attendance/makeup`

Body:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `subscription_id` | ObjectID string | yes | Must reference an existing active subscription. |
| `branch_id` | ObjectID string | yes | Attendance branch. |
| `date` | RFC3339 string | no | Makeup date; defaults to server time when omitted. |
| `is_makeup_for` | RFC3339 string | yes | Must equal the source `reported_missed` date for the same subscription. |

Response `201`:

```json
{
  "message": "attendance makeup recorded successfully",
  "data": {
    "id": "69f20c010c4cd4cdf5768501",
    "sub_id": "69f20b22f79bb78cac99aa0a",
    "branch_id": "69f20a180c4cd4cdf57684fe",
    "date": "2026-05-14T08:00:00Z",
    "status": "makeup",
    "is_makeup_for": "2026-05-12T08:00:00Z"
  }
}
```

Status codes:

- `201`: makeup attendance created.
- `400`: invalid body, ObjectID, RFC3339 date input, or missing `is_makeup_for`.
- `404`: subscription not found.
- `409`: subscription state/expiry conflict, weekly or remaining-session limit, invalid/not-found makeup source reference, or reused makeup reference.
- `500`: internal server error.

Notes:

- Client does not send `status`; handler sets `makeup`.
- Makeup must reference a reported-missed date within 7 days and consumes one remaining session.

### Branch nearby

`GET /api/v1/branches/nearby`

Query:

| Field | Type | Required | Rule |
|---|---|---:|---|
| `lng` | number | yes | Longitude, range `-180..180`. |
| `lat` | number | yes | Latitude, range `-90..90`. |
| `max_distance` | integer | no | Meter distance. Default `5000`. Explicit `<= 0` returns `400`. |
| `limit` | integer | no | Default `10`. `0` means default. Valid final range `1..100`; `>100` returns `400`. |

Response `200`:

```json
{
  "message": "nearby branches fetched successfully",
  "data": [
    {
      "id": "69f20a180c4cd4cdf57684fe",
      "branch_code": "HCM01",
      "name": "Ho Chi Minh Main Branch",
      "address": "123 Nguyen Hue, District 1",
      "province": "Ho Chi Minh",
      "location": {
        "type": "Point",
        "coordinates": [106.7009, 10.7769]
      },
      "manager_id": "000000000000000000000000",
      "distance_meters": 123.45
    }
  ]
}
```

Status codes:

- `200`: success, including empty `data`.
- `400`: missing/invalid query, out-of-range coordinates, invalid `max_distance`, invalid `limit`.
- `500`: internal server error or geo query/index failure.

Notes:

- Coordinates use GeoJSON order `[lng, lat]`.
- Route is registered before `/api/v1/branches/:id`.
- Branch create/update require `location.type = "Point"` and valid coordinates.
- Startup creates MongoDB `branches.location` 2dsphere index.

---

## Status code mac dinh

- 200: OK
- 201: Created
- 400: Bad request (invalid input)
- 401: Unauthorized (missing/invalid/expired token or invalid credentials)
- 403: Forbidden (authenticated staff does not have an allowed role)
- 404: Not found
- 409: Conflict (status/logic conflict)
- 500: Internal server error

---

## Notes

| Ghi chú | Nội dung |
|---|---|
| Offline payment | Confirm qua PATCH /members/:id/activate + `subscription_id`. |
| Subscription | Tạo mới ở trạng thái `pending`, server tính `subtotal_amount`, `discount_amount`, `total_amount_paid`, sau đó activate khi confirm payment. |
| Refund | `POST /api/v1/subscriptions/:id/refund` chỉ áp dụng cho `active`, reject `pending`/`suspended`/`expired`/`refunded`, tính `refund_amount = total_amount_paid * remaining_sessions / total_sessions`, sau đó set subscription `refunded` và `remaining_sessions = 0`. |
| Sessions MVP | Đã có create/list/get/enroll/checkin; enrollment lưu trên session và check-in tạo attendance có `session_id`. |
| Branch nearby | `GET /api/v1/branches/nearby` dùng `lng`, `lat`, optional `max_distance`, `limit`; trả thêm `distance_meters`. |
| Course tags | `allowed_tags` của course là tập tag được phép dùng để ràng buộc session. |
