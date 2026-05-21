# Code Review - attendance report makeup

## Status

- Status: re-reviewed
- Feature: attendance report makeup
- Plan file: `CHAT_CONTEXT/backend_skills/plans/03_attendance_report_makeup.md`
- Implementation file: `CHAT_CONTEXT/backend_skills/implementations/03_attendance_report_makeup.md`
- Test file: `CHAT_CONTEXT/backend_skills/tests/03_attendance_report_makeup.md`
- Reviewed at: 2026-05-21

## Review Summary

- Result: pass for planned minimal scope
- Recommendation: proceed to `/backend-complete 03_attendance_report_makeup`
- `go build ./...`: pass in this review run
- `go test ./...`: pass in this review run
- Test phase evidence: manual API and direct Mongo verification recorded as pass in `tests/03_attendance_report_makeup.md`

## Findings

- No blocking correctness or contract finding was found in this re-review.

## Checklist

- [x] Code compiles or build issue recorded.
- [x] Handler only handles HTTP parse/response/error mapping.
- [x] Service owns business rules.
- [x] Repository only handles DB access.
- [x] Model contract is sufficient for planned minimal scope.
- [x] Errors map to planned HTTP status.
- [ ] Atomic updates used for all race-prone attendance side effects.
- [x] Routes are wired without path-order conflict.
- [x] `docs/api_contract.md` marks report/makeup endpoints implemented.
- [x] `api_test.http` includes report/makeup request samples.
- [x] No secret or local-only config change is present in the Cycle 03 diff.

## Passed

- Dedicated endpoints are exposed in route wiring:
  - `POST /api/v1/attendance/report`
  - `POST /api/v1/attendance/makeup`
- New dedicated handlers keep `status` server-controlled:
  - report forces `reported_missed`
  - makeup forces `makeup`
- New handlers parse ObjectIDs and RFC3339 input, then delegate rule evaluation to `AttendanceService.CheckIn`.
- Service-side guards cover the Cycle 03 behavior:
  - active subscription required
  - subscription expiry checked against attendance date
  - `reported_missed` 30-day window enforced
  - makeup source report must exist by exact `is_makeup_for` date
  - makeup must be within 7 days
  - reused makeup reference rejected by service checks
  - attended/makeup records reject no remaining sessions before insert
  - weekly session limit applies to makeup
- Docs and manual samples are aligned with the new endpoints:
  - `docs/api_contract.md` marks both endpoints implemented
  - `api_test.http` includes request samples
- Test report covers happy path, invalid input, not found, conflicts, overdue makeup window, and direct DB state verification.

## Remaining Risks

- [medium] Duplicate makeup prevention is still service-level. Concurrent makeup requests for the same `is_makeup_for` can still pass the read checks before one insert is visible to the other.
- [medium] Attendance creation, subscription session decrement, and member attended-count increment are not a single atomic unit. A write failure after attendance insert can leave partial state, and concurrent session-consuming requests still depend on non-atomic remaining-session writes.
- [low] Makeup source lookup depends on exact RFC3339 instant equality for `is_makeup_for`; a stable source report ID would be less fragile.
- [low] No feature-specific automated tests exist yet for report/makeup endpoints; current Cycle 03 evidence is build/test plus manual API/DB verification.

## Fixes Applied During Review

- None. Review phase updated the review artifact only.

## Completion Recommendation

- Complete Cycle 03 for the minimal endpoint-exposure scope already planned and tested.
- Carry the concurrency/data-integrity risks into Cycle 07 index/data-integrity work:
  - DB-enforced duplicate makeup protection or a stable report reference strategy
  - atomic remaining-session consumption and side-effect consistency
- Carry endpoint integration coverage into Cycle 08 integration tests/fixtures.
