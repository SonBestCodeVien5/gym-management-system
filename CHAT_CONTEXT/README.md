# Project Context Snapshot

Read this file when a new chat needs a short project handoff.

## Purpose

- This folder is project memory for resuming work, not a report draft store.
- Durable project docs start at [../docs/README.md](../docs/README.md).
- Formal report source material lives in
  [../docs/report-materials/README.md](../docs/report-materials/README.md).
- Backend phase plans, implementation notes, reviews, and test notes stay in
  [backend_skills/README.md](backend_skills/README.md).
- Frontend phase plans, implementation notes, reviews, and test notes stay in
  [frontend_skills/README.md](frontend_skills/README.md).

## Read Order

1. Read this snapshot.
2. Read [../docs/README.md](../docs/README.md) for the documentation map.
3. For backend feature delivery, use the focused skill: `$gym-plan`, `$gym-implement`,
   `$gym-review`, `$gym-test`, `$gym-complete`, `$gym-resume`, or `$gym-status`.
4. For frontend feature delivery, use `$gym-fe-plan`, `$gym-fe-implement`, `$gym-fe-review`,
   `$gym-fe-test`, or `$gym-fe-complete`.
5. Read source code and current API docs for behavior that must be exact.

## Current State

Snapshot date: 2026-06-02.

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
| Auth | Login, current employee, refresh rotation, logout revoke, access-token middleware, role guard |
| Employees | Admin-only create, list, get, update, password reset, deactivate |
| Dashboard | Admin/manager summary KPIs, revenue buckets, plan distribution, recent members, today's sessions |
| Error handling | Shared HTTP error contract with stable `error.code`, sanitized `message`, and object `details` |
| Data integrity | Central MongoDB index bootstrap with unique/query/partial unique/TTL indexes |
| Frontend readiness | Allow-list CORS for browser FE dev and current employee restore endpoint |
| Integration tests | `internal/app` shared router wiring, `internal/testutil` fixtures, and MongoDB-backed integration tests |

Planned next surfaces:
- CI automation and expanded integration coverage for sessions/not-found cases remain follow-ups.
- FE02 dashboard reference is complete functionally. FE02.1 dashboard responsive repair is a
  temporary containment pass, with final responsive design deferred to FE12.
- FE03 app routing/API foundation is implemented, reviewed, tested, and complete.
- FE04 brand asset integration is implemented, reviewed, tested, and complete. The frontend now uses
  selected official Iron Forge runtime assets for favicon/metadata, login/sidebar/status branding,
  loading state, and app not-found state. No backend API contract changed. Live backend
  login/restore/logout was not re-smoked during FE04 because no credentials or seeded local session
  were available.
- FE06-FE10 interfaces are implemented, reviewed, review-fixed, build-verified, and completed with
  explicit skipped live verification notes. A follow-up mocked browser pass also covered FE05,
  FE06, FE08, and FE10 route rendering, plus FE07 subscription and FE09 session mutation-success /
  background-refresh-failure alerts. The batch covers courses/branches, subscriptions, attendance,
  sessions, and admin employee management. Live backend CRUD/API smokes remain pending because no
  seeded backend credentials/session data were available.
- Backend dashboard/report aggregate APIs are implemented and tested. FE11 can now wire the live
  dashboard against `GET /api/v1/dashboard/summary`, `/revenue`, `/plans`, `/members/recent`, and
  `/sessions/today`.
- Frontend delivery now has dedicated `$gym-fe-*` skills and `frontend_skills/` memory.

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
- Browser FE dev can use `CORS_ALLOWED_ORIGINS` and `GET /api/v1/auth/me` to restore current staff
  context from a bearer access token.
- Dashboard/report endpoints are admin/manager-only. Net revenue is subscription payments minus
  processed refunds; recent members are not branch-scoped because members have no branch field.
- Employee management is admin-only, stores bcrypt password hashes, never returns `password_hash` or
  `normalized_email`, and revokes active refresh tokens on password reset/deactivation.
- Backend error responses use `{"error":{"code":"...","message":"...","details":{}}}` while success
  responses keep the existing `message`/`data` shape.
- Startup runs `pkg/database.EnsureIndexes` before repository construction. Unique indexes enforce
  member CCID, branch code, employee email/ID, refresh-token hash, refund subscription, duplicate
  session check-in, and duplicate makeup reuse. Refresh-token TTL cleanup is eventual.
- `internal/app.NewRouter` owns production/test dependency wiring and route registration. Integration
  tests use isolated `gym_test_*` MongoDB databases and skip when MongoDB is not reachable.

## Where To Update

| When | Update |
|---|---|
| API behavior changes | `docs/api_contract.md`, `api_test.http`, and this snapshot if the project surface changed |
| Backend phase advances | Relevant files under `backend_skills/` plus `backend_skills/worklog.md` |
| Report draft changes | `docs/report-materials/` |
| Documentation structure changes | `docs/README.md` and relevant `$gym-*` skill if loading rules changed |

## Resume Point

FE07/FE09 refresh-alert fixes and FE05-FE10 mocked browser test notes are current.
Start from:

1. `$gym-fe-implement` for FE11 live dashboard APIs using
   `CHAT_CONTEXT/frontend_skills/plans/11_live_dashboard_apis.md`
2. `$gym-fe-test` if a seeded backend/auth session is available for live browser/API smokes
3. `$gym-git` to review/commit/push the current backend dashboard and frontend planning/test notes
