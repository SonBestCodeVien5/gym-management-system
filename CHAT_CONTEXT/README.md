# Project Context Snapshot

Read this file when a new chat needs a short project handoff.

## Purpose

- This folder is project memory for resuming work, not a report draft store.
- Durable project docs start at [../docs/README.md](../docs/README.md).
- Formal report source material lives in
  [../docs/report-materials/README.md](../docs/report-materials/README.md).
- Backend phase plans, implementation notes, reviews, and test notes stay in
  [backend_skills/README.md](backend_skills/README.md).

## Read Order

1. Read this snapshot.
2. Read [../docs/README.md](../docs/README.md) for the documentation map.
3. For backend feature delivery, use the focused skill: `$gym-plan`, `$gym-implement`,
   `$gym-review`, `$gym-test`, `$gym-complete`, `$gym-resume`, or `$gym-status`.
4. Read source code and current API docs for behavior that must be exact.

## Current State

Snapshot date: 2026-05-28.

Stack:
- Go + Gin + MongoDB.
- Layered backend flow: handler -> service -> repository -> MongoDB.

Implemented backend surfaces:

| Area | Current surface |
|---|---|
| Members | Register, get by ID, activate offline payment, list member subscriptions |
| Subscriptions | Create pending subscription, get, suspend, unsuspend, expire, refund |
| Courses | CRUD |
| Branches | CRUD and nearby geo search |
| Attendance | Free check-in, report missed, makeup, history by subscription |
| Sessions | Create, list, get, enroll, session check-in |
| Auth | Login, refresh rotation, logout revoke, access-token middleware, role guard |
| Employees | Admin-only create, list, get, update, password reset, deactivate |
| Error handling | Shared HTTP error contract with stable `error.code`, sanitized `message`, and object `details` |
| Data integrity | Central MongoDB index bootstrap with unique/query/partial unique/TTL indexes |

Planned next surfaces:
- Integration tests and fixtures.

## Rules Worth Remembering

- Subscription creation validates member/course/branch references and snapshots course pricing,
  session count, and allowed tags.
- Refund currently applies only to active subscriptions and sets refunded subscriptions to zero
  remaining sessions after computing the refund amount.
- Attendance enforces weekly session limits, a 30-day reported-missed window, and a 7-day makeup
  reference window.
- Branch nearby search depends on GeoJSON coordinates and a MongoDB `2dsphere` index.
- Session enrollment stores subscription IDs on the session and session check-in reuses attendance
  rules.
- Auth requires `JWT_ACCESS_SECRET` and `JWT_REFRESH_SECRET`; bootstrap admin is created from
  `BOOTSTRAP_ADMIN_*` env values only when the normalized email is absent.
- Employee management is admin-only, stores bcrypt password hashes, never returns `password_hash` or
  `normalized_email`, and revokes active refresh tokens on password reset/deactivation.
- Backend error responses use `{"error":{"code":"...","message":"...","details":{}}}` while success
  responses keep the existing `message`/`data` shape.
- Startup runs `pkg/database.EnsureIndexes` before repository construction. Unique indexes enforce
  member CCID, branch code, employee email/ID, refresh-token hash, refund subscription, duplicate
  session check-in, and duplicate makeup reuse. Refresh-token TTL cleanup is eventual.

## Where To Update

| When | Update |
|---|---|
| API behavior changes | `docs/api_contract.md`, `api_test.http`, and this snapshot if the project surface changed |
| Backend phase advances | Relevant files under `backend_skills/` plus `backend_skills/worklog.md` |
| Report draft changes | `docs/report-materials/` |
| Documentation structure changes | `docs/README.md` and relevant `$gym-*` skill if loading rules changed |

## Resume Point

Cycle 07 indexes and data-integrity hardening is complete. The next backend cycle is integration
tests and fixtures.
Start from:

1. `backend_skills/plans/08_integration_tests_fixtures.md` if present, otherwise create it with
   `$gym-plan`
2. the requested phase skill for the next task
3. only the source files needed for that task
