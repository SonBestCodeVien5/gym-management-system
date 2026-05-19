# Skill 04 — Test Backend Feature

Dùng skill này sau code review.

## Mục tiêu test

- Xác nhận build pass.
- Xác nhận API chạy đúng happy path.
- Xác nhận error cases đúng status code.
- Xác nhận business rules không bị bypass.
- Xác nhận docs/API samples đủ để test tay.

## Test levels

### 1. Build test

Bắt buộc:

```bash
go build ./...
```

Nếu project có test không phụ thuộc DB hoặc test DB đã sẵn:

```bash
go test ./...
```

### 2. API manual test bằng `api_test.http`

Mỗi feature cần thêm request mẫu:

- Happy path.
- Invalid input.
- Not found.
- Conflict/business rule.

### 3. Integration test sau này

Khi có test harness:
- Setup Mongo test DB riêng.
- Clean collections trước/sau test.
- Seed branch/course/member.
- Gọi router trực tiếp bằng `httptest`.
- Assert status code + response body + DB state.

## Checklist test chung

- [ ] `go build ./...` pass.
- [ ] `go test ./...` pass hoặc ghi rõ lý do chưa chạy được.
- [ ] `api_test.http` có request mẫu.
- [ ] Happy path trả đúng status.
- [ ] Invalid input trả 400.
- [ ] Not found trả 404.
- [ ] Conflict trả 409.
- [ ] DB state sau action đúng.
- [ ] Docs cập nhật trạng thái endpoint.

## Test riêng theo feature

### Refund + pricing

Happy path:
1. Create branch.
2. Create course.
3. Create member.
4. Create subscription với discount nếu có.
5. Activate member/subscription.
6. Check-in để giảm remaining sessions.
7. Refund subscription.
8. Verify:
   - response có `refund_amount`
   - subscription status `refunded`
   - `remaining_sessions = 0`

Error cases:
- refund invalid ObjectID → 400
- refund subscription không tồn tại → 404
- refund pending → 409
- refund expired/refunded → 409
- refund remaining_sessions = 0 → 409
- discount percent > 100 → 400
- discount fixed > subtotal → 400

### Branch nearby

Happy path:
1. Create branches có GeoJSON Point `[lng, lat]`.
2. Call `/api/v1/branches/nearby?lng=...&lat=...`.
3. Verify list sorted near-first nếu dùng `$geoNear`.
4. Verify distance nếu response include.

Error cases:
- missing lat/lng → 400
- lat > 90 hoặc < -90 → 400
- lng > 180 hoặc < -180 → 400
- max_distance <= 0 → 400
- limit quá lớn thì clamp hoặc 400 theo rule đã chọn

### Auth

Happy path:
1. Seed employee password hash.
2. Login.
3. Call protected route với access token.
4. Refresh token.
5. Logout.
6. Refresh lại sau logout phải fail.

Error cases:
- wrong password → 401
- missing token → 401
- role không đủ quyền → 403
- expired/invalid token → 401

## Test report template

Ghi vào `worklog.md`:

```md
## Test — <feature> — <date>

### Commands
- `go build ./...` — pass/fail
- `go test ./...` — pass/fail/skipped reason

### Manual API
- ...

### Results
- ...

### Issues
- ...

### Fixed
- ...
```

## Khi test fail

1. Đọc error đầy đủ.
2. Xác định fail do compile, logic, DB, route, hay test data.
3. Sửa code, không sửa rule để né test.
4. Chạy lại command liên quan.
5. Ghi lỗi và fix vào `worklog.md`.