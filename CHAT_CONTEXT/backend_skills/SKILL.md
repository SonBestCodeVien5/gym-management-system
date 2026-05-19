# Backend Feature Delivery Skill

## Skill name

Backend Feature Delivery — Go + Gin + MongoDB Clean Architecture.

## How to invoke

Use this skill by prompt convention:

```txt
Use backend skill.

/backend-implement 01_refund_pricing
```

Slash-style commands are not native CLI commands. They are project-level prompt aliases defined in this file. Agent should resolve each command by reading this `SKILL.md`, then lazy-load only the files required by that command.

## Trigger

Use this skill when working on backend features in this project, especially when user asks to:

- plan backend work
- implement backend feature
- review backend changes
- test backend/API flow
- resume interrupted backend task
- update project context after feature completion

## Project architecture

Stack:
- Go
- Gin
- MongoDB
- Clean Architecture style

Layer rules:
- `internal/models`: DB/JSON structs.
- `internal/repository`: MongoDB access only.
- `internal/service`: business rules and orchestration.
- `internal/handlers`: HTTP parsing/response/error mapping.
- `cmd/server/main.go`: dependency wiring and route registration.
- `docs/api_contract.md`: FE/BE API contract.
- `api_test.http`: manual API samples.

## Required first read

At start of any backend feature session, read:

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/SKILL.md`
4. matching agent role skill:
   - plan/design: `agent_skills/backend_architect.md`
   - implement: `agent_skills/backend_implementer.md`
   - review: `agent_skills/backend_reviewer.md`
   - test: `agent_skills/backend_tester.md`
   - API/docs alignment: `agent_skills/api_contract_keeper.md`
   - context update/resume: `agent_skills/context_maintainer.md`
5. current phase skill:
   - plan: `01_plan.md`
   - implement: `02_implement.md`
   - review: `03_code_review.md`
   - test: `04_test.md`
6. current feature file in relevant phase folder:
   - `plans/<feature>.md`
   - `implementations/<feature>.md`
   - `reviews/<feature>.md`
   - `tests/<feature>.md`

Then read only relevant source files.

## Phase workflow

Each feature must move through 4 phases:

1. Plan
2. Implement
3. Code review
4. Test

For feature `01_refund_pricing`, phase files are:

```txt
plans/01_refund_pricing.md
implementations/01_refund_pricing.md
reviews/01_refund_pricing.md
tests/01_refund_pricing.md
```

## Operating rules

- Do not implement before plan is clear.
- Do not read all feature files unless needed.
- Do not put business rules in handlers.
- Do not put DB logic in services.
- Do not trust client input for money, status, role, or computed counts.
- Use atomic DB updates for double-submit/race-risk flows.
- Keep context files concise; summarize decisions, do not paste full source code.
- Update phase file while working.
- Update `worklog.md` with short status only.
- Update `CHAT_CONTEXT/README.md` after completed feature.

## Context budget rules

Avoid context bloat:

- Read global context once.
- Read only current feature files.
- Read only current phase skill.
- Read source files by layer as needed.
- Use search/read targeted files, not whole project dumps.
- Prefer summaries in context files over full code copies.

## Output obligations by phase

### Plan phase

Update:
- `plans/<feature>.md`
- `01_plan.md` roadmap status if needed
- `worklog.md` short status

Must include:
- API contract
- business rules
- data changes
- repo/service/handler/route steps
- docs/test plan
- risks

### Implement phase

Update:
- `implementations/<feature>.md`

Must include:
- files changed
- key decisions
- commands run
- limitations
- handoff notes for review

### Review phase

Update:
- `reviews/<feature>.md`

Must include:
- checklist result
- issues found
- fixes applied
- remaining risks
- handoff notes for test

### Test phase

Update:
- `tests/<feature>.md`

Must include:
- build/test command results
- manual API results
- DB verification
- issues found
- final result

## Completion Definition of Done

Feature complete only when:

- Plan exists and matches implementation.
- Code implemented in correct layers.
- Code reviewed with checklist.
- `go build ./...` passes.
- `go test ./...` passes or skip reason recorded.
- `api_test.http` updated.
- `docs/api_contract.md` updated.
- `worklog.md` updated.
- `CHAT_CONTEXT/README.md` updated.
- Phase files reflect final state.

## Slash-style commands

These are project-level prompt aliases. They are not native CLI commands. Use them in chat like:

```txt
Use backend skill.

/backend-implement 01_refund_pricing
```

### `/backend-plan <feature>`

Purpose:
- Plan new or existing backend feature.
- Update `plans/<feature>.md`.
- Keep `01_plan.md` as roadmap only.
- Do not implement.

Required reading:
1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/SKILL.md`
4. `CHAT_CONTEXT/backend_skills/context_loading.md`
5. `CHAT_CONTEXT/backend_skills/agent_skills/backend_architect.md`
6. `CHAT_CONTEXT/backend_skills/01_plan.md`
7. `docs/api_contract.md`
8. `CHAT_CONTEXT/backend_skills/plans/<feature>.md` if exists

### `/backend-implement <feature>`

Purpose:
- Implement backend feature from existing plan.
- Update `implementations/<feature>.md`.
- Run `gofmt` and `go build ./...`.
- Do not start review/test unless asked.

Required reading:
1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/SKILL.md`
4. `CHAT_CONTEXT/backend_skills/context_loading.md`
5. `CHAT_CONTEXT/backend_skills/agent_skills/backend_implementer.md`
6. `CHAT_CONTEXT/backend_skills/02_implement.md`
7. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
8. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`

### `/backend-review <feature>`

Purpose:
- Review implemented feature.
- Update `reviews/<feature>.md`.
- Check code layers, API contract, error mapping, data integrity.
- Do not start test unless asked.

Required reading:
1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/SKILL.md`
4. `CHAT_CONTEXT/backend_skills/context_loading.md`
5. `CHAT_CONTEXT/backend_skills/agent_skills/backend_reviewer.md`
6. `CHAT_CONTEXT/backend_skills/03_code_review.md`
7. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
8. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`
9. `CHAT_CONTEXT/backend_skills/reviews/<feature>.md`

### `/backend-test <feature>`

Purpose:
- Test feature.
- Update `tests/<feature>.md`.
- Run `go build ./...` and `go test ./...` when possible.
- Record manual API test results.

Required reading:
1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/SKILL.md`
4. `CHAT_CONTEXT/backend_skills/context_loading.md`
5. `CHAT_CONTEXT/backend_skills/agent_skills/backend_tester.md`
6. `CHAT_CONTEXT/backend_skills/04_test.md`
7. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
8. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`
9. `CHAT_CONTEXT/backend_skills/reviews/<feature>.md`
10. `CHAT_CONTEXT/backend_skills/tests/<feature>.md`

### `/backend-complete <feature>`

Purpose:
- Finalize feature after plan/implement/review/test.
- Update project context and docs.
- Mark status in `worklog.md`.

Required reading:
1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/SKILL.md`
3. `CHAT_CONTEXT/backend_skills/agent_skills/api_contract_keeper.md`
4. `CHAT_CONTEXT/backend_skills/agent_skills/context_maintainer.md`
5. `CHAT_CONTEXT/backend_skills/worklog.md`
6. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
7. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`
8. `CHAT_CONTEXT/backend_skills/reviews/<feature>.md`
9. `CHAT_CONTEXT/backend_skills/tests/<feature>.md`

Update:
- `CHAT_CONTEXT/README.md`
- `CHAT_CONTEXT/backend_skills/worklog.md`
- `docs/api_contract.md`
- `api_test.http`

### `/backend-resume <feature>`

Purpose:
- Resume interrupted backend work.
- Determine current phase from phase files and user instruction.
- Inspect changed files only.
- Continue from last known state.

Required reading:
1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/SKILL.md`
4. `CHAT_CONTEXT/backend_skills/context_loading.md`
5. `CHAT_CONTEXT/backend_skills/agent_skills/context_maintainer.md`
6. relevant role skill for current phase
7. relevant phase files for `<feature>`

If in Act mode, also run:
```bash
git status --short
```

### `/backend-status`

Purpose:
- Summarize backend roadmap and current feature states.
- Read only summary files unless user asks for detail.

Required reading:
1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/README.md`
3. `CHAT_CONTEXT/backend_skills/01_plan.md`
4. `CHAT_CONTEXT/backend_skills/worklog.md`

## Lazy-loaded files

Do not load every file every session. Load by command/phase:

- `context_loading.md`: read when deciding what context to load.
- `resume_prompts.md`: read when user asks for prompt examples or resume pattern.
- `agent_skills/backend_architect.md`: read for `/backend-plan`.
- `agent_skills/backend_implementer.md`: read for `/backend-implement`.
- `agent_skills/backend_reviewer.md`: read for `/backend-review`.
- `agent_skills/backend_tester.md`: read for `/backend-test`.
- `agent_skills/api_contract_keeper.md`: read for `/backend-complete` and docs/API alignment.
- `agent_skills/context_maintainer.md`: read for `/backend-resume`, `/backend-complete`, and context updates.
- `01_plan.md`: read for `/backend-plan`.
- `02_implement.md`: read for `/backend-implement`.
- `03_code_review.md`: read for `/backend-review`.
- `04_test.md`: read for `/backend-test`.
- `plans/<feature>.md`: read for target feature plan.
- `implementations/<feature>.md`: read during implementation/review/test.
- `reviews/<feature>.md`: read during review/test/complete.
- `tests/<feature>.md`: read during test/complete.
- `worklog.md`: read for status/complete only.

Core principle:

```txt
SKILL.md = entrypoint + command registry
Other files = lazy-loaded context
```

## Current roadmap

1. Refund flow & pricing rules
2. Branch nearby geo query
3. Attendance report/makeup endpoints
4. Auth/login + role guard
5. Validation hardening & error consistency
6. Indexes and data integrity
7. Integration tests & fixtures
