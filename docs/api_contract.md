# API Contract (Current)

Cap nhat: 2026-05-05

Muc tieu: chot ten endpoint + request/response co ban de FE va BE dung chung.

---

## 1) Implemented

### Members
- POST /api/v1/members
- GET /api/v1/members/:id
- PATCH /api/v1/members/:id/activate

### Subscriptions
- POST /api/v1/subscriptions
- GET /api/v1/subscriptions/:id
- PATCH /api/v1/subscriptions/:id/suspend
- PATCH /api/v1/subscriptions/:id/unsuspend
- PATCH /api/v1/subscriptions/:id/expire

### Courses
- POST /api/v1/courses
- GET /api/v1/courses
- GET /api/v1/courses/:id
- PATCH /api/v1/courses/:id
- DELETE /api/v1/courses/:id

### Branches
- POST /api/v1/branches
- GET /api/v1/branches
- GET /api/v1/branches/:id
- PATCH /api/v1/branches/:id
- DELETE /api/v1/branches/:id

### Attendance
- POST /api/v1/attendance/checkin
- GET /api/v1/subscriptions/:id/attendance

---

## 2) Planned (chua co trong code)

- POST /api/v1/attendance/report
- POST /api/v1/attendance/makeup
- GET /api/v1/branches/nearby
- GET /api/v1/members/:id/subscriptions
- POST /api/v1/subscriptions/:id/refund
- Auth endpoints: /api/v1/auth/login, /api/v1/auth/refresh, /api/v1/auth/logout

---

## 3) Status code mac dinh

- 200: OK
- 201: Created
- 400: Bad request (invalid input)
- 404: Not found
- 409: Conflict (status/logic conflict)
- 500: Internal server error

---

## 4) Notes

- Offline payment: confirm qua PATCH /members/:id/activate + subscription_id.
- Subscription tao moi o status pending, activate sau khi payment confirm.
