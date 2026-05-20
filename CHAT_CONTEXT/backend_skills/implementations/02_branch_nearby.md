# Implementation — branch nearby

## Status

- Status: implemented
- Feature: branch nearby
- Plan file: `CHAT_CONTEXT/backend_skills/plans/02_branch_nearby.md`
- Started at: 2026-05-20 21:57 +07
- Finished at: 2026-05-20 22:08 +07

## Scope implemented

- [x] Model changes
- [x] Repository changes
- [x] Service changes
- [x] Handler changes
- [x] Route changes
- [ ] Docs/API sample changes

## Files changed

- `internal/models/branch.go`
- `internal/repository/branch_repo.go`
- `internal/service/branch_service.go`
- `internal/handlers/branch_handler.go`
- `cmd/server/main.go`
- `internal/service/subscription_service_test.go`

## Key decisions

- Added `models.BranchNearbyResult` for geo query output with `distance_meters`.
- Added `BranchRepository.Nearby` using MongoDB aggregation with `$geoNear`.
- Added `2dsphere` index creation on `branches.location` in `NewBranchRepository`.
- Changed `NewBranchRepository` to return `(BranchRepository, error)` and wired startup failure handling in `cmd/server/main.go`.
- Kept coordinate validation in service layer:
  - GeoJSON type must be `Point`.
  - coordinates order is `[lng, lat]`.
  - lng range is `[-180, 180]`.
  - lat range is `[-90, 90]`.
- Handler only parses query params and maps errors; service owns defaults and business validation.
- Registered `/branches/nearby` before `/branches/:id`.

## Implementation notes

- Nearby defaults:
  - `max_distance=5000` meters when omitted.
  - `limit=10` when omitted.
- Nearby constraints:
  - `max_distance > 0`.
  - `limit` range `1..100`.
- Create/update branch location validation hardened to reject non-Point GeoJSON and out-of-range coordinates.
- Existing subscription service test stub for `BranchRepository` now includes no-op `Nearby` because repository interface expanded.
- Docs/API sample changes left for complete/API-contract phase per phase boundary.

## Commands run

- `gofmt -w internal/models/branch.go internal/repository/branch_repo.go internal/service/branch_service.go internal/handlers/branch_handler.go cmd/server/main.go internal/service/subscription_service_test.go`
  - Result: executed successfully.
- `go build ./...`
  - Result: executed successfully.

## Known limitations

- Existing branch documents with invalid or non-Point `location` may fail geo index creation or not appear in nearby results.
- No status/opening-hours/capacity filtering because current branch model has no such fields.
- No manual API test run in implementation phase.
- `go test ./...` not run because `/backend-implement` requires build, not full test phase.

## Handoff to review

- Review route order and Gin matching for `/branches/nearby`.
- Review MongoDB `$geoNear` pipeline and index creation behavior.
- Review service validation/defaults for nearby query and branch create/update hardening.
- Complete phase should update `docs/api_contract.md`, `api_test.http`, `CHAT_CONTEXT/README.md`, and `worklog.md`.