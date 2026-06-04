# Backend Plan - 11 Final Project Package

Status: Planned

Created: 2026-06-04

## Goal

Finish the project as a runnable, demonstrable, and report-ready submission package.

This phase should not add new business-domain features. It should make the existing backend and
frontend easy to run with realistic demo data, complete Docker packaging, align root/user-facing
documentation, and reconcile formal report material with the implementation that now exists.

## Current Baseline

- Backend business surfaces are implemented for members, courses, branches, subscriptions,
  attendance, sessions, auth, employees, dashboard/report aggregates, validation/error consistency,
  indexes, and integration tests.
- Frontend is implemented under `frontend/` as a Vite/React staff portal consuming the live API.
- `docker-compose.yml` currently starts MongoDB only. There is no backend Dockerfile, frontend
  Dockerfile, full-stack compose flow, or Docker health/smoke guidance.
- `.env.example` exists for local backend env, but the server currently always opens database
  `gym_management`; the commented `DB_NAME` value is not consumed by `cmd/server/main.go`.
- There is no seed/demo data command. Test fixtures exist under `internal/testutil`, but they are
  for isolated integration tests, not local demo setup.
- Root `README.md` exists but is backend-centered and does not yet cover the final full-stack
  architecture, Docker flow, seed/demo data, report structure, frontend run flow, or submission
  checklist.
- `docs/report-materials/07_current_implementation_evidence.md` is stale in places: it still marks
  employee management and dashboard/report endpoints as planned/partial even though they are now
  implemented.

## API Contract

No new public HTTP endpoint is planned for this phase.

Existing HTTP behavior should remain unchanged:

- Keep `/api/v1/*` auth, role guards, success response shape, and shared error contract exactly as
  documented in `docs/api_contract.md`.
- Keep dashboard endpoints admin/manager-only.
- Keep employee management admin-only.
- Keep browser auth through bearer access tokens and refresh-token rotation.

Operational additions should be command/documentation/docker surfaces:

- `cmd/seed` or equivalent local seed command for demo data.
- Docker image build/run instructions.
- Full-stack compose service names and ports.
- README/report material updates.

If implementation discovers a small HTTP-support gap that blocks Docker or demo usage, treat it as a
separate explicit change and update `docs/api_contract.md` plus `api_test.http`.

## Business Rules

### Seed/demo data

- Seed data must be deterministic enough for demos and README instructions.
- Seed data must be idempotent by stable unique keys where possible:
  - employees by normalized email or employee ID
  - branches by branch code
  - members by CCID
  - courses by stable course name/code convention if no code field exists
- Seed should create a realistic cross-section:
  - admin, manager, receptionist, and trainer demo employees
  - at least three branches with GeoJSON coordinates suitable for nearby search
  - beginner/intermediate/advanced courses with prices, session counts, tags, and levels
  - members across registered/unregistered states
  - pending, active, suspended, expired, and refunded subscription examples where safe
  - attendance records including attended, reported missed, and makeup examples
  - sessions for today/this week so dashboard and sessions pages are not empty
- Seed command must not silently destroy user data. If a reset mode is needed, make it explicit
  through a flag or env value and document it clearly.
- Demo passwords must be non-secret examples only and clearly marked for local demo use.

### Docker/runtime

- Docker local stack should run MongoDB, backend API, and frontend with one documented command.
- Backend container should use env config, wait for MongoDB health, ensure indexes on startup, and
  bootstrap the configured admin account.
- Frontend container should serve the production Vite build and point browser requests at the
  documented backend base URL.
- Compose should preserve MongoDB data in a named volume and document the reset command separately.
- Secrets in compose/examples are local-development placeholders, not production guidance.

### Report/README

- Report material must distinguish implemented behavior from future work.
- Root README should become the first-stop project overview for an evaluator:
  project purpose, tech stack, architecture, features, screenshots/assets if available, quickstart,
  Docker, seed data, API docs, test commands, and report links.
- Durable docs should point to root README for final run instructions and continue to point to
  `docs/api_contract.md` for exact API behavior.

## Data And Configuration Changes

Planned data/config changes:

- Add configurable database name support, likely `DB_NAME` with fallback `gym_management`, so Docker,
  local dev, tests, and seed command can target the intended database explicitly.
- Add seed/demo command source under `cmd/seed` or a similarly clear path.
- Add/update example env files:
  - root `.env.example` for backend/API and Docker defaults
  - optional `frontend/.env.example` for `VITE_API_BASE_URL`
- Add Docker artifacts:
  - root backend `Dockerfile`
  - `frontend/Dockerfile`
  - `.dockerignore`
  - full-stack `docker-compose.yml` or a clearly named additional compose file if preserving a
    Mongo-only dev compose is preferable
- Avoid schema migrations unless a seed-only stable identifier truly requires one. Prefer existing
  unique indexes and current schema.

## Layer Plan

### Backend config/startup

- Update `cmd/server/main.go` to read `DB_NAME` with fallback `gym_management`.
- Keep index bootstrap before router construction.
- Keep auth secret validation in `service.NewAuthService`.
- Keep bootstrap admin behavior unchanged except for documentation/examples.

### Seed command

- Add `cmd/seed/main.go`.
- Load `.env` when present, then read the same MongoDB/DB env values as the server.
- Connect to MongoDB, ensure indexes, and write demo data using existing repositories/services where
  practical so business rules stay consistent.
- Use direct MongoDB upsert only for seed-specific setup that services do not support safely.
- Print a concise summary of created/found records and demo login accounts.
- Return non-zero on missing required env, invalid seed config, or storage/index errors.

### Docker

- Add a multi-stage backend image:
  - build Go binary
  - copy only runtime binary/config needs
  - run as a non-root user when practical
  - expose `8080`
- Add a frontend image:
  - run `npm ci` and `npm run build`
  - serve `dist/` with nginx or another small static server
  - expose a documented port, preferably mapped to `5173` or another evaluator-friendly port
- Update compose:
  - services: `mongodb`, `api`, `frontend`
  - backend depends on MongoDB health
  - frontend depends on API service start
  - documented ports for API and frontend
  - named MongoDB volume
- Decide whether seed runs as:
  - a one-shot compose profile/service, or
  - a documented `docker compose run --rm api ./seed` style command if the image contains both
    binaries.

### Frontend docs/config

- Document `VITE_API_BASE_URL` for local and Docker builds.
- Keep frontend API base behavior aligned with `frontend/src/lib/api.js`.
- Run `npm run build` after Docker/frontend doc changes.

### Documentation and report

- Rewrite/expand root `README.md` as the final project guide.
- Update `docs/local_dev_guide.md` for `DB_NAME`, seed command, frontend local run, and full-stack
  Docker.
- Update `docs/README.md` if the docs map needs a final report or Docker/readme entry.
- Reconcile `docs/report-materials/07_current_implementation_evidence.md` with current implemented
  employee/dashboard/frontend status.
- Create or update a final report assembly file under `docs/report-materials/` if needed, using the
  existing chapter inputs instead of duplicating all source material into chat context.

## Docs And Test Plan

Documentation updates:

- `README.md`
- `.env.example`
- optional `frontend/.env.example`
- `docs/README.md`
- `docs/local_dev_guide.md`
- `docs/report-materials/07_current_implementation_evidence.md`
- optional final report assembly file in `docs/report-materials/`
- `CHAT_CONTEXT/README.md` and `CHAT_CONTEXT/backend_skills/worklog.md` when the phase is complete

Automated verification:

```sh
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
npm --prefix frontend run build
docker compose config
docker compose build
```

Manual/local verification when Docker is available:

```sh
docker compose up -d mongodb api frontend
curl -s http://localhost:8080/ping
go run ./cmd/seed
```

Then smoke:

- Login as demo admin.
- Call `GET /api/v1/auth/me`.
- Call dashboard summary/revenue/plans/recent-members/today-sessions.
- Call course, branch, member, subscription, attendance, session, and employee list/detail flows
  enough to confirm seed data is usable.
- Open the frontend route locally and confirm login/dashboard/modules load against the seeded API.

If Docker cannot run in the current environment, record the limitation and at least verify compose
syntax/buildable files where possible.

## Risks And Boundaries

- A seed command can accidentally mutate real data if pointed at the wrong database. Mitigate with
  explicit DB naming, clear env docs, idempotent upserts, and explicit reset flags only.
- Demo data should avoid asserting financial/accounting accuracy beyond current MVP rules.
- Compose networking can be confusing because the browser needs a host-reachable API URL while
  containers use service names internally. Document the browser-facing `VITE_API_BASE_URL` clearly.
- Building frontend with a fixed API URL means changing the URL requires rebuilding the static image
  unless runtime config is added. Keep this as a documented MVP tradeoff unless runtime config is
  necessary.
- MongoDB transactions are still not part of the current architecture; do not expand final packaging
  into large data-consistency refactors.
- Report material must not claim future features as implemented. Reconcile stale report text before
  final handoff.
- Do not commit real secrets, generated database dumps, local volumes, or build outputs.

## Next Action

Use `$gym-implement` with `CHAT_CONTEXT/backend_skills/plans/11_final_project_package.md`.
