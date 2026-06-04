# Gym Management System

Full-stack multi-branch gym management system built with Go, Gin, MongoDB, and a React/Vite staff
portal.

The project supports member registration, subscription lifecycle management, offline payment
activation, attendance and makeup sessions, class sessions, staff authentication, employee
management, branch/course administration, and management dashboard metrics.

## Tech Stack

| Layer | Technology |
|---|---|
| Backend API | Go, Gin |
| Database | MongoDB |
| Frontend | React 18, Vite |
| Auth | Employee login, access tokens, refresh-token rotation |
| Packaging | Docker Compose with MongoDB, API, frontend, and seed profile |

## Architecture

Backend requests follow a layered flow:

```text
HTTP request
  -> Gin route
  -> Handler: parse request and map HTTP responses
  -> Service: validate business rules and orchestrate workflows
  -> Repository: read/write MongoDB
  -> MongoDB
```

HTTP errors use the shared response contract:

```json
{"error":{"code":"INVALID_INPUT","message":"invalid request body","details":{}}}
```

Successful responses keep the existing `message` and `data` shape.

## Implemented Features

- Auth: login, current employee restore, refresh rotation, logout revoke, role guards.
- Employees: admin-only create/list/get/update/password reset/deactivate.
- Members: register, get by ID, activate offline payment, list member subscriptions.
- Courses: create/list/get/update/delete training packages.
- Branches: create/list/get/update/delete and nearby search with GeoJSON `2dsphere`.
- Subscriptions: create pending subscriptions, activate, suspend, unsuspend, expire, refund.
- Attendance: check-in, report missed sessions, makeup sessions, subscription attendance history.
- Sessions: create/list/get/enroll/check-in class sessions.
- Dashboard: summary KPIs, revenue buckets, plan distribution, recent members, today's sessions.
- Data integrity: central MongoDB index bootstrap for unique, query, partial unique, and TTL indexes.
- Frontend staff portal: dashboard and operational modules backed by the live API.
- Demo seed data: deterministic local dataset for demos and evaluator walkthroughs.

## Quickstart With Docker

Run the full stack:

```bash
docker compose up -d --build
```

If Docker reports a missing `docker-buildx` plugin on this machine, build with the legacy builder
first:

```bash
DOCKER_BUILDKIT=0 docker compose build
docker compose up -d
```

Load deterministic demo data:

```bash
docker compose --profile seed run --rm seed
```

On a machine missing `docker-buildx`, use the same legacy-builder prefix for the seed profile:

```bash
DOCKER_BUILDKIT=0 docker compose --profile seed run --rm seed
```

Open:

- Frontend: `http://localhost:5173`
- API health: `http://localhost:8080/ping`
- MongoDB: `localhost:27017`

Demo accounts:

| Role | Email | Password |
|---|---|---|
| Admin | `admin@gym.test` | `demo123456` |
| Manager | `manager@gym.test` | `demo123456` |
| Receptionist | `receptionist@gym.test` | `demo123456` |
| Trainer | `trainer@gym.test` | `demo123456` |

Stop the stack:

```bash
docker compose down
```

Reset local Docker data only when you intentionally want a fresh database:

```bash
docker compose down -v
```

Use that reset path if MongoDB refuses to start because an old local Docker volume was created by a
newer MongoDB feature compatibility version.

## Local Development

Copy or adjust environment files:

```bash
cp .env.example .env
cp frontend/.env.example frontend/.env
```

Start MongoDB:

```bash
docker compose up -d mongodb
```

Run the backend:

```bash
go run ./cmd/server
```

Seed demo data into the configured `DB_NAME`:

```bash
go run ./cmd/seed
```

Run the frontend:

```bash
npm --prefix frontend install
npm --prefix frontend run dev
```

The frontend reads `VITE_API_BASE_URL`; by default it points at `http://localhost:8080`.

## Configuration

Important backend environment variables:

| Variable | Purpose |
|---|---|
| `MONGODB_URI` | MongoDB connection string |
| `DB_NAME` | Application database name, defaults to `gym_management` |
| `PORT` | Backend API port, defaults to `8080` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated browser origins for local FE dev |
| `JWT_ACCESS_SECRET` | Access-token signing secret |
| `JWT_REFRESH_SECRET` | Refresh-token signing secret |
| `BOOTSTRAP_ADMIN_*` | Optional first admin account created on startup |

The example secrets and demo passwords are local-development placeholders.

## Verification

Backend:

```bash
env GOCACHE=/tmp/gocache go build ./...
env GOCACHE=/tmp/gocache go test ./...
```

Frontend:

```bash
npm --prefix frontend run build
```

Docker:

```bash
docker compose config
docker compose build
```

Use `DOCKER_BUILDKIT=0 docker compose build` on machines without the `docker-buildx` plugin.

Integration tests use isolated `gym_test_*` MongoDB databases and skip cleanly when MongoDB is not
reachable.

## Documentation

- Documentation hub: [docs/README.md](docs/README.md)
- API contract: [docs/api_contract.md](docs/api_contract.md)
- Local development guide: [docs/local_dev_guide.md](docs/local_dev_guide.md)
- Code reading guide: [docs/code_reading_guide.md](docs/code_reading_guide.md)
- Architecture rationale: [docs/faq_why.md](docs/faq_why.md)
- Formal report source material: [docs/report-materials/README.md](docs/report-materials/README.md)
- Current implementation evidence: [docs/report-materials/07_current_implementation_evidence.md](docs/report-materials/07_current_implementation_evidence.md)

Project continuity for Codex/chat handoff starts at [CHAT_CONTEXT/README.md](CHAT_CONTEXT/README.md).
