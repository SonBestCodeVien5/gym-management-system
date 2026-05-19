# Skill 02 — Implement Backend Feature

Dùng skill này khi đã có plan rõ.

## Thứ tự code chuẩn

1. Model
2. Repository interface + implementation
3. Service interface + implementation
4. Handler request/response DTO
5. Route wiring
6. Docs/API samples
7. Build/test

## Nguyên tắc Clean Architecture trong project này

- `internal/models`: struct lưu DB và JSON response.
- `internal/repository`: MongoDB CRUD/query, không chứa business rule.
- `internal/service`: business rule, validate trạng thái, phối hợp repo.
- `internal/handlers`: parse HTTP input, gọi service, map error sang status code.
- `cmd/server/main.go`: dependency wiring + route registration.

## Checklist implement

- [ ] Thêm/sửa model với `bson` + `json` tags.
- [ ] Thêm repo method vào interface.
- [ ] Implement repo method bằng MongoDB query đúng.
- [ ] Repo trả `repository.ErrNotFound` khi không có document.
- [ ] Service thêm error domain rõ.
- [ ] Service validate input và business status.
- [ ] Service không tin client với field server tự tính.
- [ ] Handler DTO không expose field nguy hiểm.
- [ ] Handler parse ObjectID/date/query param đúng.
- [ ] Handler map errors: 400/404/409/500.
- [ ] Route thêm đúng order.
- [ ] Cập nhật `api_test.http`.
- [ ] Cập nhật `docs/api_contract.md`.
- [ ] Chạy `go build ./...`.

## Error mapping mặc định

| Error type | HTTP |
|---|---|
| invalid input/date/ObjectID/query | 400 |
| not found | 404 |
| duplicate/conflict/status conflict/business rule | 409 |
| unknown storage/server error | 500 |

## Pattern service error

```go
var (
    ErrFeatureInvalidInput = errors.New("invalid feature input")
    ErrFeatureNotFound = errors.New("feature not found")
    ErrFeatureConflict = errors.New("feature conflict")
)
```

## Pattern repo not found

```go
if errors.Is(err, mongo.ErrNoDocuments) {
    return nil, ErrNotFound
}
```

## Pattern atomic update

Dùng khi action không được chạy 2 lần: refund, enroll, payment confirm, check-in duplicate prevention.

```go
result, err := collection.UpdateOne(
    ctx,
    bson.M{
        "_id": objID,
        "status": bson.M{"$in": []string{"active", "suspended"}},
    },
    bson.M{"$set": update},
)
if result.MatchedCount == 0 {
    return ErrNotFoundOrConflict
}
```

Sau đó service có thể fetch trước để phân biệt 404 và 409 nếu cần response đẹp.

## Pattern route order

Route cụ thể phải đứng trước route param.

```go
api.GET("/branches/nearby", branchHandler.Nearby)
api.GET("/branches/:id", branchHandler.GetByID)
```

## Sau khi code

Chạy:

```bash
gofmt -w <changed-go-files>
go build ./...
go test ./...
```

Nếu test chưa có hoặc cần DB, vẫn phải chạy `go build ./...`.