# Test — branch nearby

## Status

- Status: passed after bugfix retest
- Feature: branch nearby
- Plan file: `CHAT_CONTEXT/backend_skills/plans/02_branch_nearby.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/02_branch_nearby.md`
- Review file: `CHAT_CONTEXT/backend_skills/reviews/02_branch_nearby.md`
- First tested at: 2026-05-20 22:22 +07
- Re-tested at: 2026-05-20 22:32 +07

## Commands

- `go build ./...` — pass
  - Output capture failed, but command returned successful execution.
- `go test ./...` — pass

```txt
?   	github.com/SonBestCodeVien5/gym-management-system/cmd/server	[no test files]
?   	github.com/SonBestCodeVien5/gym-management-system/internal/handlers	[no test files]
?   	github.com/SonBestCodeVien5/gym-management-system/internal/models	[no test files]
?   	github.com/SonBestCodeVien5/gym-management-system/internal/repository	[no test files]
ok  	github.com/SonBestCodeVien5/gym-management-system/internal/service	(cached)
?   	github.com/SonBestCodeVien5/gym-management-system/pkg/database	[no test files]
```

## Manual API environment

- MongoDB container: `gym_mongodb Up 34 minutes`
- Server command:
  - `PORT=18081 MONGODB_URI='mongodb://admin:password123@localhost:27017/?authSource=admin' go run ./cmd/server`
- Temporary server stopped after manual checks.
- Test data inserted into local `gym_management.branches`, then cleaned up after user approval.

## Manual API tests

### Happy path

- [x] Create branch with GeoJSON Point `[106.7000, 10.8000]` — `201`
- [x] `GET /api/v1/branches/nearby?lng=106.7001&lat=10.8001&max_distance=10000&limit=10` — `200`
- [x] `GET /api/v1/branches/nearby?lng=106.7001&lat=10.8001` — `200`, omitted `max_distance` defaults to `5000`
- [x] Server log confirms nearby route handled by `BranchHandler.Nearby`.
- [x] `GET /api/v1/branches/:id` for created branch — `200`; route order ok.

### Invalid input

- [x] Missing `lng`: `GET /api/v1/branches/nearby?lat=10.8001` — `400`
- [x] Out-of-range `lat`: `GET /api/v1/branches/nearby?lng=106.7001&lat=91` — `400`
- [x] Explicit zero `max_distance`: `GET /api/v1/branches/nearby?lng=106.7001&lat=10.8001&max_distance=0` — `400`
- [x] Negative `max_distance`: `GET /api/v1/branches/nearby?lng=106.7001&lat=10.8001&max_distance=-1` — `400`
- [x] `limit=101`: `GET /api/v1/branches/nearby?lng=106.7001&lat=10.8001&limit=101` — `400`
- [x] `limit=0`: `GET /api/v1/branches/nearby?lng=106.7001&lat=10.8001&limit=0` — `200`, defaulted to 10.

### Not found

- [x] Not applicable for nearby list: empty result should return `200` by contract.
- [x] Existing branch `GET /api/v1/branches/:id` still resolves after route reorder.

### Conflict/business rule

- [x] Not applicable: feature is read-only geo query.
- [x] Previous issue fixed: explicit `max_distance=0` now returns `400`.

## DB state verification

- [x] MongoDB connection successful.
- [x] Branch repository initialized and created/used `location_2dsphere` index without startup error.
- [x] Test branch inserted with valid `location.type = "Point"` and `[lng, lat]` coordinates.
- [x] Nearby query executed through MongoDB `$geoNear` without DB/index error.
- [x] DB cleanup performed after user approval; manual test branch records removed.

## Issues found

### Issue 1 — explicit `max_distance=0` accepted

- Severity: medium
- First test actual:
  - `GET /api/v1/branches/nearby?lng=106.7001&lat=10.8001&max_distance=0` returned `200`
- Expected by plan/review:
  - `400`
- Root cause:
  - Handler parsed omitted `max_distance` and explicit `max_distance=0` into same zero value.
  - Service treated `maxDistance == 0` as omitted/default `5000`.
- Fix:
  - Handler now distinguishes explicit `max_distance <= 0` from omitted `max_distance`.
  - Explicit invalid value is converted to invalid sentinel so service returns `ErrInvalidBranchInput`.
- Re-test actual:
  - `max_distance=0` — `400`
  - `max_distance=-1` — `400`
  - omitted `max_distance` — `200`
- Status: fixed and verified.

### Issue 2 — manual test data left in local DB

- Severity: low
- Local `gym_management.branches` contained manual test records from first test and retest.
- Cleanup command after user approval:
  - `db.branches.deleteMany({branch_code: {$regex: "^(TNEAR|TFAR|RNEAR)" }})`
- Cleanup result:
  - `deletedCount: 3`
- Verification:
  - matching records query returned `[]`
- Status: fixed.

## Final result

- Result: pass after bugfix retest
- Build/test: pass
- Manual happy path: pass
- Manual invalid inputs: pass
- Previous blocker: fixed
- Ready to update docs/context: yes, proceed to `/backend-complete 02`.
- Manual DB cleanup: done.
