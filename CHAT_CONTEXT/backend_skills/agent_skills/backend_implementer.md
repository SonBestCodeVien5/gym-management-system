# Agent Skill — Backend Implementer

## Role

Implement backend features using Go + Gin + MongoDB Clean Architecture.

## Use when

- `/backend-implement <feature>`
- plan is ready
- user asks to code backend feature
- interrupted implementation needs resume

## Must read

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/SKILL.md`
3. `CHAT_CONTEXT/backend_skills/context_loading.md`
4. `CHAT_CONTEXT/backend_skills/02_implement.md`
5. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
6. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`

## Responsibilities

- Implement model changes.
- Implement repository DB operations.
- Implement service business rules.
- Implement handler request/response/error mapping.
- Wire routes in `cmd/server/main.go`.
- Update API samples/docs when needed.
- Run formatting/build checks.

## Layer rules

- Model: structs, BSON/JSON tags, constants.
- Repository: MongoDB queries and atomic updates only.
- Service: validation, business rules, orchestration.
- Handler: bind input, call service, map errors.
- Route wiring: dependency injection and endpoints only.

## Safety rules

- Do not trust client money/status/session count/role fields.
- Do not duplicate business rules across handler and service.
- Use atomic update for double-submit/race-prone flows.
- Preserve existing API behavior unless plan says otherwise.
- Do not modify tests to hide broken implementation.
- Do not start review/test phase unless user asks.

## Output

Update:
- `implementations/<feature>.md`

Record:
- files changed
- key decisions
- commands run
- known limitations
- handoff notes for reviewer

## Required commands

Run when relevant:

```bash
gofmt -w <changed-go-files>
go build ./...
```

Run tests if available or if feature has tests:

```bash
go test ./...
```

## Checklist

- [ ] Code follows planned API.
- [ ] Business rules are in service.
- [ ] DB logic is in repository.
- [ ] Handler maps errors consistently.
- [ ] Routes wired correctly.
- [ ] Docs/API samples updated if needed.
- [ ] `gofmt` run.
- [ ] `go build ./...` pass or failure recorded.