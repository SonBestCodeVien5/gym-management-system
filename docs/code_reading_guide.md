# Huong dan doc hieu code (chi tiet)

Muc tieu: giup ban moi vao du an doc hieu code nhanh, biet cho nao lam gi, tai sao co doan nay, va luong di nhu the nao.

---

## 1) Ban do tong quan (top-down)

Bat dau tu entrypoint va routes de thay he thong co gi.

### 1.1 Entry point
File: cmd/server/main.go

Doc theo thu tu:
1. Tai env va ket noi DB.
2. Tao repo -> tao service -> tao handler.
3. Gan routes vao Gin.

Ket qua: ban se biet toan bo API hien co va dependency wiring.

### 1.2 Phan lop kien truc
- handlers: nhan HTTP request, parse input, goi service, map loi sang HTTP.
- service: kiem tra nghiep vu, chuyen trang thai, goi repository.
- repository: thao tac MongoDB.
- models: struct du lieu.

---

## 2) Luong chinh can hieu (member + subscription)

### 2.1 Dang ky member
File: internal/handlers/member_handler.go

Doc ham Register theo thu tu:
1. Parse JSON body.
2. Map sang models.Member.
3. Goi MemberService.RegisterMember.
4. Map error -> HTTP code.

File: internal/service/member_service.go

Doc ham RegisterMember:
1. Validate input (ccid, full_name).
2. Check CCID unique qua repository.
3. Gan default fields (is_registered, timestamps).
4. Goi repo Create.

File: internal/repository/member_repo.go

Doc Create va GetByCCID de thay no lam gi voi MongoDB.

### 2.2 Tao subscription
File: internal/handlers/subscription_handler.go

Doc ham Create:
1. Parse JSON.
2. Validate ObjectID cho member/course/branch.
3. Parse start_date/end_date RFC3339.
4. Tao models.Subscription.
5. Goi SubscriptionService.CreateSubscription.

File: internal/service/subscription_service.go

Doc ham CreateSubscription:
1. Validate input (IDs, dates, session_per_week).
2. Kiem tra member/course/branch ton tai.
3. Set status pending.
4. Copy gia va so buoi tu course.
5. Goi repo Create.

### 2.3 Activate offline payment
File: internal/handlers/member_handler.go

Doc ham Activate:
1. Validate member id.
2. Parse subscription_id tu body.
3. Goi SubscriptionService.ConfirmSubscriptionPayment.
4. Goi MemberService.ActivateMember.

File: internal/service/subscription_service.go

Doc ham ConfirmSubscriptionPayment:
1. Load subscription.
2. Kiem tra subscription thuoc member.
3. Kiem tra status = pending.
4. Update status = active + set payment_date.

---

## 3) Subscription lifecycle

File: internal/handlers/subscription_handler.go

- Suspend: nhan thong tin suspension (start/end/frozen/reason), goi service.
- Resume (unsuspend): goi service de clear suspension va set active.
- Expire: set status expired.

File: internal/service/subscription_service.go

Doc cac ham:
- SuspendSubscription: validate suspension, status active, chua expired.
- ResumeSubscription: chi cho suspended, neu het han thi loi.
- ExpireSubscription: set status expired.

File: internal/repository/subscription_repo.go

Doc cac ham update status/suspension/remaining_sessions de biet cap nhat DB ra sao.

---

## 4) Course CRUD

Files:
- internal/handlers/course_handler.go
- internal/service/course_service.go
- internal/repository/course_repo.go

Doc theo thu tu: handler -> service -> repo. Trong service co validate input (title, level, base_price, session_count).

---

## 5) Branch CRUD

Files:
- internal/handlers/branch_handler.go
- internal/service/branch_service.go
- internal/repository/branch_repo.go

Doc theo thu tu: handler -> service -> repo. Chu y field location (GeoJSON) va manager_id co the rong.

---

## 6) Attendance check-in

Files:
- internal/handlers/attendance_handler.go
- internal/service/attendance_service.go
- internal/repository/attendance_repo.go

Doc ham CheckIn:
1. Parse JSON, validate subscription_id va branch_id.
2. Parse date va is_makeup_for neu co.
3. Goi service CheckIn.

Doc service CheckIn:
1. Validate input.
2. Load subscription, check status active.
3. Neu status = `attended` hoac `makeup`, kiem tra quota theo tuan (`sessionPerWeek`).
4. Neu `reported_missed`, kiem tra sliding window 30 ngay.
5. Neu `makeup`, kiem tra report goc trong 7 ngay va chua dung.
6. Tao attendance record.
7. Neu attended/makeup: tru remaining_sessions, cong total_sessions_attended.
8. Neu remaining = 0 -> set expired.

Luu y hien tai:
- `attended` va `makeup` deu tinh vao quota tuan va lam giam buoi con lai.
- `absent` chi luu record, khong tru buoi va khong tang count.
- `reported_missed` khong tru buoi nhung bi gioi han 1 lan / 30 ngay.
- `makeup` phai tham chieu den `reported_missed` hop le trong 7 ngay va khong duoc dung lai report do.

Doc ham ListBySubscriptionID de thay cach lay history.

---

## 7) Cac tai lieu ho tro doc hieu

- README.md: tong quan tinh nang va routes.
- docs/local_dev_guide.md: cach chay local va test nhanh.
- docs/development_journal.md: tien do va van de da gap.
- docs/update_2026-05-03.md: tom tat nhung gi vua lam.

---

## 8) Cach doc khi gap loi hoac can debug

1. Xem handler co nhan dung body/param khong.
2. Xem service co tra dung error khong.
3. Xem repository query dung field khong.
4. Kiem tra ID co dung ObjectID khong.
5. Kiem tra status/remaining_sessions trong DB.

---

## 9) Checklist tu hoc de nho luong

- Tim route trong cmd/server/main.go.
- Tim handler cua route do.
- Tim service ma handler goi.
- Tim repository ma service goi.
- Kiem tra model lien quan.

---

## 10) Goi y doc theo thu tu (toi uu)

1. cmd/server/main.go
2. internal/handlers/member_handler.go
3. internal/service/member_service.go
4. internal/repository/member_repo.go
5. internal/handlers/subscription_handler.go
6. internal/service/subscription_service.go
7. internal/repository/subscription_repo.go
8. internal/handlers/attendance_handler.go
9. internal/service/attendance_service.go
10. internal/repository/attendance_repo.go
11. internal/handlers/course_handler.go
12. internal/service/course_service.go
13. internal/repository/course_repo.go
14. internal/handlers/branch_handler.go
15. internal/service/branch_service.go
16. internal/repository/branch_repo.go

---

## 11) So do luong (flow diagram) theo API chinh

### 11.1 Member registration
Request -> member_handler.Register -> member_service.RegisterMember -> member_repo.Create -> MongoDB -> Response

### 11.2 Create subscription
Request -> subscription_handler.Create -> subscription_service.CreateSubscription ->
member_repo.GetByID + course_repo.GetByID + branch_repo.GetByID -> subscription_repo.Create -> Response

### 11.3 Activate offline payment
Request -> member_handler.Activate -> subscription_service.ConfirmSubscriptionPayment -> subscription_repo.UpdateStatusAndPaymentDate ->
member_service.ActivateMember -> member_repo.UpdateRegistrationStatus -> Response

### 11.4 Attendance check-in
Request -> attendance_handler.CheckIn -> attendance_service.CheckIn ->
subscription_repo.GetByID -> attendance_repo.Create ->
subscription_repo.UpdateRemainingSessions(AndStatus) + member_repo.IncrementSessionsAttended -> Response

Neu status la `attended`/`makeup`, service se kiem tra them quota theo tuan truoc khi tao record.
Neu status la `reported_missed`, service se chan neu da co report trong vong 30 ngay.
Neu status la `makeup`, service se can `is_makeup_for` va report goc phai nam trong 7 ngay.

---

## 12) Bang input/output cho tung endpoint

### 12.1 Members
- POST /api/v1/members
	- Input: ccid, full_name, email, phone, gender, level
	- Output: member object

- GET /api/v1/members/:id
	- Input: path id
	- Output: member object

- PATCH /api/v1/members/:id/activate
	- Input: subscription_id
	- Output: message only

### 12.2 Subscriptions
- POST /api/v1/subscriptions
	- Input: member_id, course_id, home_branch_id, start_date, end_date, session_per_week
	- Output: subscription object

- GET /api/v1/subscriptions/:id
	- Input: path id
	- Output: subscription object

- PATCH /api/v1/subscriptions/:id/suspend
	- Input: start_date, end_date, frozen_session, reason
	- Output: message only

- PATCH /api/v1/subscriptions/:id/unsuspend
	- Input: none
	- Output: message only

- PATCH /api/v1/subscriptions/:id/expire
	- Input: none
	- Output: message only

### 12.3 Courses
- POST /api/v1/courses
	- Input: title, level, base_price, session_count, description
	- Output: course object

- GET /api/v1/courses
	- Input: none
	- Output: list course

- GET /api/v1/courses/:id
	- Input: path id
	- Output: course object

- PATCH /api/v1/courses/:id
	- Input: title, level, base_price, session_count, description
	- Output: message only

- DELETE /api/v1/courses/:id
	- Input: path id
	- Output: message only

### 12.4 Branches
- POST /api/v1/branches
	- Input: branch_code, name, address, province, location{type, coordinates}, manager_id (optional)
	- Output: branch object

- GET /api/v1/branches
	- Input: none
	- Output: list branch

- GET /api/v1/branches/:id
	- Input: path id
	- Output: branch object

- PATCH /api/v1/branches/:id
	- Input: branch_code, name, address, province, location{type, coordinates}, manager_id (optional)
	- Output: message only

- DELETE /api/v1/branches/:id
	- Input: path id
	- Output: message only

### 12.5 Attendance
- POST /api/v1/attendance/checkin
	- Input: subscription_id, branch_id, date (optional), status, is_makeup_for (optional)
	- Output: attendance object

- GET /api/v1/subscriptions/:id/attendance
	- Input: subscription id
	- Output: list attendance

---

## 13) Checklist doc code theo luong va test nhanh

### 13.1 Doc code theo luong
1. Tim route trong cmd/server/main.go.
2. Mo handler tuong ung.
3. Xem ham service ma handler goi.
4. Xem repo ma service goi.
5. Mo model lien quan de hieu field.

### 13.2 Test nhanh theo thu tu
1. Tao branch va course.
2. Tao member.
3. Tao subscription (pending).
4. Activate member + confirm payment.
5. Check-in attendance.
6. Suspend/unsuspend/expire subscription.

