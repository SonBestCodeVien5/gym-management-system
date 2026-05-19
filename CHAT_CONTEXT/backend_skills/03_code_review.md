# Skill 03 — Code Review Backend Feature

Dùng skill này sau khi implement xong, trước khi test/build final.

## Mục tiêu review

- Bắt lỗi logic.
- Bắt lỗi Clean Architecture.
- Bắt lỗi MongoDB query/index.
- Bắt lỗi route/path conflict.
- Bắt lỗi API contract mismatch.
- Bắt lỗi security/auth nếu có.

## Checklist review tổng quát

- [ ] Code compile về mặt import/package/signature.
- [ ] Handler không chứa business rule phức tạp.
- [ ] Service không phụ thuộc Gin/HTTP.
- [ ] Repository không chứa business rule.
- [ ] Model field name Go đúng style, JSON/BSON đúng contract.
- [ ] Error domain rõ, không trả raw DB error cho client.
- [ ] HTTP status code đúng contract.
- [ ] Route cụ thể đứng trước route param.
- [ ] Không tin input client với tiền, status, số buổi, role.
- [ ] Date/time dùng RFC3339 ở handler.
- [ ] ObjectID validate trước khi query nếu cần.
- [ ] Mongo update atomic khi có nguy cơ double-submit.
- [ ] Index cần thiết đã có plan hoặc đã implement.
- [ ] Docs/API samples khớp code.

## Review theo layer

### Models

Kiểm tra:
- `bson:"_id,omitempty"` cho ID.
- Field money dùng `int64`.
- Field count dùng `int`.
- Optional time dùng `*time.Time`.
- Không đặt tên Go kiểu `Total_Amount_Paid`; dùng `TotalAmountPaid`.

### Repository

Kiểm tra:
- Interface và implementation đồng bộ.
- `primitive.ObjectIDFromHex` error được xử lý.
- `mongo.ErrNoDocuments` map sang `repository.ErrNotFound`.
- Cursor được close.
- Update/Delete check `MatchedCount`/`DeletedCount`.
- Query atomic cho actions nhạy cảm.
- Không swallow error.

### Service

Kiểm tra:
- Input nil/empty/zero handled.
- Reference existence checked nếu feature cần.
- State transition hợp lệ.
- Business rule không bị bypass.
- Tính tiền/số buổi từ snapshot DB, không từ request.
- Conflict trả error riêng.

### Handler

Kiểm tra:
- `ShouldBindJSON`/query parse đúng.
- Date parse RFC3339 nếu input là string.
- Không expose internal-only fields.
- Map service error đúng status:
  - 400 invalid input
  - 404 not found
  - 409 conflict
  - 500 unknown
- Response shape nhất quán với phần còn lại.

### Routes

Kiểm tra:
- Path trong `docs/api_contract.md` khớp `main.go`.
- Route static trước route param:
  - `/branches/nearby` trước `/branches/:id`
- Auth middleware không chặn endpoint public.

## Review riêng cho feature sắp làm

### Refund + pricing

- [ ] Không cho refund pending/expired/refunded.
- [ ] Không refund khi remaining_sessions <= 0.
- [ ] Refund amount tính từ `total_amount_paid * remaining_sessions / total_sessions`.
- [ ] Set subscription `status = refunded`.
- [ ] Set `remaining_sessions = 0`.
- [ ] Chặn double refund bằng atomic update hoặc unique index `refunds.subscription_id`.
- [ ] Refund record có audit fields.
- [ ] Discount không âm, không vượt subtotal.
- [ ] Client không tự set total_amount_paid.

### Branch nearby

- [ ] Validate lat/lng range.
- [ ] Default/max limit hợp lý.
- [ ] Default max_distance hợp lý.
- [ ] Có 2dsphere index cho `location`.
- [ ] GeoJSON coordinates đúng thứ tự `[lng, lat]`.
- [ ] Response có distance nếu dùng `$geoNear`.

### Auth

- [ ] Password hash, không lưu plain text.
- [ ] JWT secret từ env.
- [ ] Refresh token hash nếu lưu DB.
- [ ] Role guard đúng endpoint.
- [ ] Không log token/password.
- [ ] Logout invalidate refresh token.

## Output review

Ghi review vào `worklog.md`:

```md
## Code Review — <feature> — <date>

### Passed
- ...

### Issues found
- ...

### Fixes applied
- ...

### Remaining risks
- ...
```
