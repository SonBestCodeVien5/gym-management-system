# Agent Skill — Backend Tester

## Role

Verify backend feature behavior through build/test commands, manual API checks, and DB state validation.

## Use when

- `/backend-test <feature>`
- review phase is done
- user asks to verify backend feature
- feature behavior must be validated before completion

## Must read

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/SKILL.md`
3. `CHAT_CONTEXT/backend_skills/context_loading.md`
4. `CHAT_CONTEXT/backend_skills/04_test.md`
5. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
6. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`
7. `CHAT_CONTEXT/backend_skills/reviews/<feature>.md`
8. `CHAT_CONTEXT/backend_skills/tests/<feature>.md`

## Responsibilities

- Run build/test commands.
- Verify happy path API flow.
- Verify invalid input cases.
- Verify not found cases.
- Verify conflict/business rule cases.
- Verify DB state changes.
- Record pass/fail and issues.

## Test rules

- Prefer existing test suite over custom-only reproduction.
- If tests fail, treat implementation as suspicious.
- Do not weaken assertions to make tests pass.
- Record environment constraints, e.g. MongoDB not running.
- Keep test report concise and reproducible.
- Update `api_test.http` when manual API sample changes.

## Output

Update:
- `tests/<feature>.md`

Record:
- commands run
- command results
- manual API requests/results
- DB verification
- issues found
- final test result

## Required commands

```bash
go build ./...
go test ./...
```

If MongoDB is required and unavailable, record skip/blocker clearly.

## Checklist

- [ ] `go build ./...` run.
- [ ] `go test ./...` run or skip reason recorded.
- [ ] Happy path tested.
- [ ] Invalid input tested.
- [ ] Not found case tested.
- [ ] Business rule/conflict case tested.
- [ ] DB state verified if relevant.
- [ ] `api_test.http` updated if needed.
- [ ] Test report final result recorded.