# gym-management-system

Multi-branch gym system management backend built with Go, Gin, and MongoDB.

## Overview

This project follows a simple layered structure:

- `handlers`: HTTP layer (request parsing, response mapping).
- `service`: business rules and workflow orchestration.
- `repository`: data access to MongoDB.
- `models`: domain data models.

## Implemented Features

### 1) Member Registration and Offline Activation

- `POST /api/v1/members`
- `GET /api/v1/members/:id`
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
- `PATCH /api/v1/branches/:id`
- `DELETE /api/v1/branches/:id`

### 5) Attendance

- `POST /api/v1/attendance/checkin`
- `GET /api/v1/subscriptions/:id/attendance`

Check-in effects (for `attended` or `makeup` status):

1. Create attendance record.
2. Decrease `remaining_sessions` of subscription.
3. Increase `total_sessions_attended` of member.
4. Auto-expire subscription if remaining sessions reach 0.

## Run Locally

1. Set environment variable `MONGODB_URI`.
2. Run server:

```bash
go run cmd/server/main.go
```

3. Use requests in `api_test.http` for manual API testing.
