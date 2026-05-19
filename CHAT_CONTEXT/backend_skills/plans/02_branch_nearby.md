# Cycle 02 — Branch Nearby Geo Query

## Status

- Status: planned
- Priority: high
- Depends on: branch CRUD, GeoJSON branch location
- Endpoint:
  - `GET /api/v1/branches/nearby`

## Goal

Thêm API tìm chi nhánh gần vị trí hiện tại bằng MongoDB geo query.

## API plan

```http
GET /api/v1/branches/nearby?lng=106.7&lat=10.8&max_distance=5000&limit=10
```

Query params:
- `lng` required, range `[-180, 180]`
- `lat` required, range `[-90, 90]`
- `max_distance` optional, default `5000`, unit meters
- `limit` optional, default `10`, max `100`

Response:
```json
{
  "branches": [
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
      "distance_meters": 350.5
    }
  ]
}
```

## Business rules

- Coordinates order is `[lng, lat]`.
- Invalid/missing query returns 400.
- Nearby route must be registered before `/branches/:id`.
- Branch location must use GeoJSON Point.
- Distance query requires `2dsphere` index.

## Data/index plan

MongoDB index:
- collection: `branches`
- field: `location`
- type: `2dsphere`

## Repository plan

Update `BranchRepository`:
```go
Nearby(ctx context.Context, lng float64, lat float64, maxDistance int64, limit int64) ([]models.BranchNearbyResult, error)
```

Use `$geoNear` aggregation if response includes distance.

## Service plan

Add:
```go
NearbyBranches(ctx context.Context, lng float64, lat float64, maxDistance int64, limit int64) ([]models.BranchNearbyResult, error)
```

Validate:
- lng/lat range
- maxDistance > 0
- limit default/max

## Handler plan

Add:
```go
func (h *BranchHandler) Nearby(c *gin.Context)
```

Parse query params and map errors.

## Route plan

Important order:
```go
api.GET("/branches/nearby", branchHandler.Nearby)
api.GET("/branches/:id", branchHandler.GetByID)
```

## Docs/test plan

Update:
- `docs/api_contract.md`
- `api_test.http`
- `CHAT_CONTEXT/README.md`
- `worklog.md`

Run:
```bash
go build ./...
go test ./...
```

## Risks

- Missing 2dsphere index causes runtime query failure.
- Existing branch create validation only checks coordinate length; should also validate type/range.