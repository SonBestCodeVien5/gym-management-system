# Context Loading Protocol

Mục tiêu: đủ context để làm đúng, không nhồi quá nhiều làm hụt context.

## Core rule

Không đọc toàn bộ folder nếu không cần.  
Không đọc tất cả feature.  
Không đọc cả 4 phase nếu chỉ đang làm 1 phase.

## Always read first

Mỗi phiên backend feature đọc:

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/SKILL.md`

## Then read by agent role and phase

### Planning phase

Read:
1. `CHAT_CONTEXT/backend_skills/agent_skills/backend_architect.md`
2. `CHAT_CONTEXT/backend_skills/01_plan.md`
3. `docs/api_contract.md`
4. target `plans/<feature>.md` if exists

Then inspect only relevant code:
- model
- handler
- service
- repository
- route wiring

### Implementation phase

Read:
1. `CHAT_CONTEXT/backend_skills/agent_skills/backend_implementer.md`
2. `CHAT_CONTEXT/backend_skills/02_implement.md`
3. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
4. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`

Then inspect source files from plan.

### Review phase

Read:
1. `CHAT_CONTEXT/backend_skills/agent_skills/backend_reviewer.md`
2. `CHAT_CONTEXT/backend_skills/03_code_review.md`
3. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
4. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`
5. `CHAT_CONTEXT/backend_skills/reviews/<feature>.md`

Then inspect changed files.

### Test phase

Read:
1. `CHAT_CONTEXT/backend_skills/agent_skills/backend_tester.md`
2. `CHAT_CONTEXT/backend_skills/04_test.md`
3. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
4. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`
5. `CHAT_CONTEXT/backend_skills/reviews/<feature>.md`
6. `CHAT_CONTEXT/backend_skills/tests/<feature>.md`

Then inspect API contract and test files.

### Complete/context update phase

Read:
1. `CHAT_CONTEXT/backend_skills/agent_skills/api_contract_keeper.md`
2. `CHAT_CONTEXT/backend_skills/agent_skills/context_maintainer.md`
3. `CHAT_CONTEXT/backend_skills/worklog.md`
4. relevant phase files for `<feature>`

## Source file loading rule

Read source by layer order:

1. models
2. repository
3. service
4. handler
5. route wiring
6. docs/tests

Do not read unrelated modules.

Example for refund/pricing:
- `internal/models/subscription.go`
- `internal/service/subscription_service.go`
- `internal/repository/subscription_repo.go`
- `internal/handlers/subscription_handler.go`
- `cmd/server/main.go`
- `docs/api_contract.md`
- `api_test.http`

## Context file writing rule

Context files should contain:
- decisions
- status
- risks
- file list
- command results
- API examples

Context files should not contain:
- full source files
- huge diffs
- repeated logs
- stack traces longer than needed

## Resume after interruption

When interrupted:
1. Read `CHAT_CONTEXT/README.md`.
2. Read `SKILL.md`.
3. Read `agent_skills/context_maintainer.md`.
4. Read role skill for current phase.
5. Read current phase file for feature.
6. Run `git status --short` if in Act mode.
7. Inspect changed files only.
8. Continue from last unchecked item.

## When context feels too large

Summarize into current phase file:
- current status
- files changed
- remaining steps
- blockers

Then continue using that summary instead of rereading everything.

## Feature key map

| Feature | Plan | Implementation | Review | Test |
|---|---|---|---|---|
| Refund/pricing | `plans/01_refund_pricing.md` | `implementations/01_refund_pricing.md` | `reviews/01_refund_pricing.md` | `tests/01_refund_pricing.md` |
| Branch nearby | `plans/02_branch_nearby.md` | `implementations/02_branch_nearby.md` | `reviews/02_branch_nearby.md` | `tests/02_branch_nearby.md` |
| Attendance report/makeup | `plans/03_attendance_report_makeup.md` | `implementations/03_attendance_report_makeup.md` | `reviews/03_attendance_report_makeup.md` | `tests/03_attendance_report_makeup.md` |
| Auth role guard | `plans/04_auth_role_guard.md` | `implementations/04_auth_role_guard.md` | `reviews/04_auth_role_guard.md` | `tests/04_auth_role_guard.md` |
| Validation/errors | `plans/05_validation_error_consistency.md` | `implementations/05_validation_error_consistency.md` | `reviews/05_validation_error_consistency.md` | `tests/05_validation_error_consistency.md` |
| Index/data integrity | `plans/06_indexes_data_integrity.md` | `implementations/06_indexes_data_integrity.md` | `reviews/06_indexes_data_integrity.md` | `tests/06_indexes_data_integrity.md` |
| Integration tests | `plans/07_integration_tests_fixtures.md` | `implementations/07_integration_tests_fixtures.md` | `reviews/07_integration_tests_fixtures.md` | `tests/07_integration_tests_fixtures.md` |