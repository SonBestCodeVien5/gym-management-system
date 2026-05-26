# gym-management-system

Multi-branch gym system management backend built with Go, Gin, and MongoDB.

## Overview

This project follows a simple layered structure:

- `handlers`: HTTP layer (request parsing, response mapping).
- `service`: business rules and workflow orchestration.
- `repository`: data access to MongoDB.
- `models`: domain data models.

Documentation map: see [docs/README.md](docs/README.md).

Current API contract: see [docs/api_contract.md](docs/api_contract.md).

Project continuity for Codex and chat handoff starts at
[CHAT_CONTEXT/README.md](CHAT_CONTEXT/README.md). Repo-scoped Codex skills live under
[.codex/skills](.codex/skills); start with `$gym-plan`, `$gym-implement`, `$gym-review`,
`$gym-test`, `$gym-complete`, `$gym-resume`, `$gym-status`, `$gym-docs`, `$gym-report`, or
`$gym-git` for focused work. Prompt and phase workflow guide:
[.codex/GYM_SKILLS_WORKFLOW.md](.codex/GYM_SKILLS_WORKFLOW.md).

## Implemented Features

### 1) Member Registration and Offline Activation

- `POST /api/v1/members`
- `GET /api/v1/members/:id`
- `GET /api/v1/members/:id/subscriptions`
- `PATCH /api/v1/members/:id/activate`

`activate` requires body `subscription_id` and performs:

1. Confirm subscription payment (`pending -> active`, set `payment_date`).
2. Mark member as registered (`is_registered = true`).

### 2) Subscription Management

- `POST /api/v1/subscriptions`
- `GET /api/v1/subscriptions/:id`
- `PATCH /api/v1/subscriptions/:id/suspend`
- `PATCH /api/v1/subscriptions/:id/unsuspend`
- `PATCH /api/v1/subscriptions/:id/expire`
- `POST /api/v1/subscriptions/:id/refund`

Behavior:

- New subscription starts with `status = pending`.
- Offline payment confirmation is tied to member activate flow.
- Suspension stores details in `suspension` field and sets `status = suspended`.

### 3) Course CRUD

- `POST /api/v1/courses`
- `GET /api/v1/courses`
- `GET /api/v1/courses/:id`
- `PATCH /api/v1/courses/:id`
- `DELETE /api/v1/courses/:id`

### 4) Branch CRUD

- `POST /api/v1/branches`
- `GET /api/v1/branches`
- `GET /api/v1/branches/:id`
- `GET /api/v1/branches/nearby`
- `PATCH /api/v1/branches/:id`
- `DELETE /api/v1/branches/:id`

### 5) Attendance

- `POST /api/v1/attendance/checkin`
- `POST /api/v1/attendance/report`
- `POST /api/v1/attendance/makeup`
- `GET /api/v1/subscriptions/:id/attendance`

Check-in effects (for `attended` or `makeup` status):

1. Create attendance record.
2. Decrease `remaining_sessions` of subscription.
3. Increase `total_sessions_attended` of member.
4. Auto-expire subscription if remaining sessions reach 0.

### 6) Sessions

- `POST /api/v1/sessions`
- `GET /api/v1/sessions`
- `GET /api/v1/sessions/:id`
- `POST /api/v1/sessions/:id/enroll`
- `POST /api/v1/sessions/:id/checkin`

### 7) Auth and Role Guard

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

Behavior:

- Business routes under `/api/v1/*` require `Authorization: Bearer <access_token>`.
- Access tokens are short-lived; refresh tokens are rotated and stored only as hashes.
- Role guard currently protects member, subscription, course, branch, attendance, and session routes.
- First admin can be bootstrapped from `BOOTSTRAP_ADMIN_*` environment variables.

### 8) Employee Management

- `POST /api/v1/employees`
- `GET /api/v1/employees`
- `GET /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id`
- `PATCH /api/v1/employees/:id/password`

Behavior:

- Employee management routes are admin-only.
- Staff passwords are stored only as bcrypt hashes.
- Responses never expose `password_hash` or `normalized_email`.
- Password reset and employee deactivation revoke active refresh tokens.
- Offboarding uses `status = inactive`; no hard delete endpoint is exposed.

## Run Locally

1. Set environment variables from `.env.example`, especially `MONGODB_URI`,
   `JWT_ACCESS_SECRET`, and `JWT_REFRESH_SECRET`.
2. Run server:

```bash
go run cmd/server/main.go
```

3. Login with the bootstrap admin, paste the returned tokens into `api_test.http`, then run the
   protected API samples.
