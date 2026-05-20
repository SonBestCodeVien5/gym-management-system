# Code Review — branch nearby

## Status

- Status: re-reviewed after bugfix
- Feature: branch nearby
- Plan file: `CHAT_CONTEXT/backend_skills/plans/02_branch_nearby.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/02_branch_nearby.md`
- Reviewed at: 2026-05-20 22:13 +07
- Re-reviewed at: 2026-05-20 22:29 +07

## Review summary

- Result: pass after bugfix review
- Build status: pass per implementation log (`go build ./...` executed successfully)
- Test status: not run in review phase per `/backend-review` boundary
- Re-review scope: `max_distance=0` bugfix in `internal/handlers/branch_handler.go`

## Checklist

- [x] Code compiles.
- [x] Handler only handles HTTP parse/response.
- [x] Service owns business rules.
- [x] Repository only handles DB.
- [x] Model tags match API/DB contract.
- [x] Errors map to correct HTTP status.
- [x] Atomic updates used where needed.
- [x] Routes have correct order.
- [ ] Docs/API samples match behavior.

## Passed

- `models.BranchNearbyResult` matches planned response fields and BSON/JSON names.
- `BranchRepository.Nearby` uses MongoDB `$geoNear`, `distanceField: "distance_meters"`, `[lng, lat]` coordinate order, `maxDistance`, and `$limit`.
- `NewBranchRepository` creates `location_2dsphere` index and returns startup error to `cmd/server/main.go`.
- `BranchService.NearbyBranches` owns validation/default rules:
  - lng range `[-180, 180]`.
  - lat range `[-90, 90]`.
  - default `maxDistance = 5000`.
  - default `limit = 10`.
  - limit range `1..100`.
- Create/update branch validation hardened to require GeoJSON `Point` and valid coordinate ranges.
- `BranchHandler.Nearby` only parses query params, delegates to service, and maps invalid input to `400`.
- Route order is correct: `/branches/nearby` registered before `/branches/:id`.
- Subscription service test stub updated to satisfy expanded `BranchRepository` interface.
- No money/status/role/computed field trust issue introduced.
- No write race/atomic update concern in this read-only geo query feature.

## Issues found

- Low: `docs/api_contract.md` and `api_test.http` are not updated yet. This is recorded in implementation notes as complete/API-contract phase work, not a code blocker.
- Low: Existing invalid branch documents can make 2dsphere index creation fail or disappear from geo results. Already documented as known limitation.

## Fixes applied during review

- None. Review phase only updated this review report; no implementation code changed.

## Re-review after bugfix

### Bugfix reviewed

- File: `internal/handlers/branch_handler.go`
- Issue from test: explicit `max_distance=0` returned `200`, but plan/test expected `400`.
- Current behavior:
  - omitted `max_distance` stays `0`, service defaults to `5000`.
  - explicit `max_distance <= 0` becomes `-1`.
  - service rejects `maxDistance < 1` with `ErrInvalidBranchInput`.
  - handler maps `ErrInvalidBranchInput` to `400`.
- Result: pass by code inspection.

### Notes

- Fix preserves plan behavior:
  - omitted `max_distance` defaults to `5000`.
  - explicit `max_distance=0` returns `400`.
  - explicit negative `max_distance` returns `400`.
- Handler still only parses query and maps errors; validation remains finalized in service.
- No new DB/repository risk introduced.

## Remaining risks

- MongoDB index creation may fail at startup if existing `branches.location` data is malformed.
- Manual API behavior needs re-test after `max_distance=0` fix.
- Full `go test ./...` not run in review phase.
- Docs/API samples remain pending for complete phase.

## Handoff to test

- Run `go build ./...`.
- Run `go test ./...`.
- Manually verify:
  - `GET /api/v1/branches/nearby?lng=106.7&lat=10.8` returns `200`.
  - Missing `lng` or `lat` returns `400`.
  - Out-of-range `lng`/`lat` returns `400`.
  - `max_distance=0` or negative returns `400`.
  - `limit=0` uses default `10`; `limit=101` returns `400`.
  - `/api/v1/branches/:id` still resolves real ObjectID after route reorder.
- Ensure test DB has valid GeoJSON:
  - `location.type = "Point"`
  - `location.coordinates = [lng, lat]`