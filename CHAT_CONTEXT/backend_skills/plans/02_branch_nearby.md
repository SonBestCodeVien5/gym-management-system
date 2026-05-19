# Cycle 02 — Branch Nearby Geo Query

## Status

- Status: planned
- Priority: high
- Depends on: branch CRUD, GeoJSON branch location
- Endpoint:
  - `GET /api/v1/branches/nearby`

## Goal

Thêm API tìm chi nhánh gần vị trí hiện tại bằng MongoDB geo query, trả danh sách branch theo khoảng cách tăng dần.

## Current source findings

- `models.Branch` đã có `Location GeoLocation` với BSON field `location`.
- `GeoLocation` hiện cho phép `Type` bất kỳ và `Coordinates []float64`.
- `BranchService.CreateBranch` / `UpdateBranch` hiện chỉ validate:
  - required text fields
  - `Location.Type != ""`
  - `len(Location.Coordinates) == 2`
- Chưa validate:
  - `Location.Type == "Point"`
  - lng range `[-180, 180]`
  - lat range `[-90, 90]`
- `BranchRepository` hiện chỉ CRUD, chưa có geo query và chưa tạo `2dsphere` index.
- `cmd/server/main.go` hiện route `/branches/:id` đứng trước route nearby chưa có. Khi thêm nearby phải đặt trước `/branches/:id`.

## API contract

```http
GET /api/v1/branches/nearby?lng=106.7&lat=10.8&max_distance=5000&limit=10
```

Query params:
- `lng` required, float, range `[-180, 180]`
- `lat` required, float, range `[-90, 90]`
- `max_distance` optional, integer meters, default `5000`, must be `> 0`
- `limit` optional, integer, default `10`, range `1..100`

Success response: `200 OK`

```json
{
  "message": "nearby branches fetched successfully",
  "data": [
    {
      "id": "ObjectID",
      "branch_code": "HCM01",
      "name": "Gym HCM 01",
      "address": "...",
      "province": "HCM",
      "location": {
        "type": "Point",
        "coordinates": [106.7, 10.8]
      },
      "manager_id": "ObjectID",
      "distance_meters": 350.5
    }
  ]
}
```

Status codes:
- `200`: nearby branch list returned, may be empty.
- `400`: missing/invalid query params or invalid coordinate range.
- `500`: DB/index/internal error.

## Business rules

### Domain business rules

- Nearby branch means geographically nearest branch from client-provided current location.
- Main use cases:
  - suggest closest branch for member.
  - help choose home branch.
  - help find potential roaming branch.
- This cycle does not validate whether member subscription can train at returned branch.
- This cycle does not filter branch active/inactive because current `Branch` model has no `status` field.
- This cycle does not filter by opening hours/capacity because current `Branch` model has no such fields.
- Future rule: when branch status exists, nearby should return only active branches by default.

### Technical validation rules

- Coordinates order is `[lng, lat]`.
- Branch location must be GeoJSON Point:
  - `type == "Point"`
  - `coordinates` length exactly 2.
- lng must be `[-180, 180]`.
- lat must be `[-90, 90]`.
- Missing/invalid nearby query returns `400`.
- Nearby route must be registered before `/branches/:id`.
- Distance is computed by MongoDB, not by handler.
- No auth requirement in this cycle because auth/role guard is planned later.

## Data changes

### Model

Add nearby result model in `internal/models/branch.go`:

```go
type BranchNearbyResult struct {
    ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    BranchCode     string             `json:"branch_code" bson:"branch_code"`
    Name           string             `json:"name" bson:"name"`
    Address        string             `json:"address" bson:"address"`
    Province       string             `json:"province" bson:"province"`
    Location       GeoLocation        `json:"location" bson:"location"`
    ManagerID      primitive.ObjectID `json:"manager_id" bson:"manager_id"`
    DistanceMeters float64            `json:"distance_meters" bson:"distance_meters"`
}
```

### Collection

- `branches`

### Index

Create `2dsphere` index on `branches.location`.

Preferred place:
- update `NewBranchRepository(db *mongo.Database)` to create index and return `(BranchRepository, error)`

Required follow-up:
- update `cmd/server/main.go` to handle branch repo init error, similar `memberRepo`.

Index spec:

```go
mongo.IndexModel{
    Keys: bson.D{{Key: "location", Value: "2dsphere"}},
}
```

## Repository plan

Update `BranchRepository`:

```go
type BranchRepository interface {
    Create(ctx context.Context, branch *models.Branch) error
    GetByID(ctx context.Context, id string) (*models.Branch, error)
    List(ctx context.Context) ([]models.Branch, error)
    UpdateByID(ctx context.Context, id string, branch *models.Branch) error
    DeleteByID(ctx context.Context, id string) error
    Nearby(ctx context.Context, lng float64, lat float64, maxDistance int64, limit int64) ([]models.BranchNearbyResult, error)
}
```

Use aggregation with `$geoNear` so response includes distance:

```go
pipeline := mongo.Pipeline{
    {{"$geoNear", bson.D{
        {"near", bson.D{
            {"type", "Point"},
            {"coordinates", bson.A{lng, lat}},
        }},
        {"distanceField", "distance_meters"},
        {"maxDistance", maxDistance},
        {"spherical", true},
    }}},
    {{"$limit", limit}},
}
```

Decode into `[]models.BranchNearbyResult`.

## Service plan

Extend `BranchService`:

```go
NearbyBranches(ctx context.Context, lng float64, lat float64, maxDistance int64, limit int64) ([]models.BranchNearbyResult, error)
```

Validation:
- `lng >= -180 && lng <= 180`
- `lat >= -90 && lat <= 90`
- if `maxDistance == 0`, default to `5000`
- `maxDistance > 0`
- if `limit == 0`, default to `10`
- `limit >= 1 && limit <= 100`
- return `ErrInvalidBranchInput` for invalid values.

Also harden create/update location validation:
- `Location.Type == "Point"`
- coordinate ranges valid.

Keep business rules in service, not handler.

## Handler plan

Add method:

```go
func (h *BranchHandler) Nearby(c *gin.Context)
```

Handler responsibilities:
- parse query params:
  - `lng` / `lat` via `strconv.ParseFloat`
  - `max_distance` / `limit` via `strconv.ParseInt`
- use zero values for omitted optional params, letting service set defaults.
- call `h.branchService.NearbyBranches(...)`.
- map:
  - parse error or `service.ErrInvalidBranchInput` → `400`
  - other errors → `500`

Do not compute distance or apply business rules in handler.

## Route plan

Update route order in `cmd/server/main.go`.

Required order:

```go
api.POST("/branches", branchHandler.Create)
api.GET("/branches", branchHandler.List)
api.GET("/branches/nearby", branchHandler.Nearby)
api.GET("/branches/:id", branchHandler.GetByID)
api.PATCH("/branches/:id", branchHandler.Update)
api.DELETE("/branches/:id", branchHandler.Delete)
```

Reason:
- Gin route `/branches/:id` can capture `nearby` if static route is registered after dynamic route or conflict depending route tree behavior.

## Docs/API sample plan

Update in implementation/complete phases:
- `docs/api_contract.md`
  - mark `GET /api/v1/branches/nearby` as Implemented after code done.
  - add query params/response notes.
- `api_test.http`
  - add nearby request after branch create/list examples.
- `CHAT_CONTEXT/README.md`
  - mark Branch nearby implemented after completion phase.
- `worklog.md`
  - short status update.

## Verification plan for implement phase

Run:
```bash
gofmt -w internal/models/branch.go internal/repository/branch_repo.go internal/service/branch_service.go internal/handlers/branch_handler.go cmd/server/main.go
go build ./...
```

Test phase later:
```bash
go test ./...
```

Manual API checks:
1. create branch with valid GeoJSON Point.
2. call nearby with valid lng/lat.
3. call nearby with missing lng/lat → expect 400.
4. call nearby with invalid lng/lat range → expect 400.
5. call nearby with `limit=101` → expect 400.
6. call `/api/v1/branches/:id` still works after route reorder.

## Risks

- Missing `2dsphere` index causes runtime query failure.
- Existing records with invalid `location.type` or coordinate order may not appear in nearby query.
- Changing `NewBranchRepository` signature requires updating `cmd/server/main.go`.
- Geo query assumes coordinates saved as `[lng, lat]`; existing clients may send `[lat, lng]` by mistake.
- FE may confuse nearby branch with eligible branch; API contract must state this is distance-only in this cycle.
- If route order is wrong, `/branches/nearby` may be treated as ObjectID and return `400 invalid branch id`.

## Architect checklist

- [x] Endpoint/method/path clear.
- [x] Request query params clear.
- [x] Success response clear.
- [x] Error cases clear.
- [x] Model changes clear.
- [x] Mongo query/index needs clear.
- [x] Business rules assigned to service.
- [x] Race/atomic concerns identified: no write race; index init risk only.
- [x] Route order concern identified.
- [x] Docs/test updates listed.