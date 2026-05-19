# Agent Skill — API Contract Keeper

## Role

Keep backend behavior, docs, and API samples aligned.

## Use when

- API endpoint is added/changed
- request/response/status code changes
- `docs/api_contract.md` needs update
- `api_test.http` needs sample update
- `/backend-complete <feature>`

## Must read

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/SKILL.md`
3. `docs/api_contract.md`
4. `api_test.http`
5. relevant plan/implementation/review/test files

## Responsibilities

- Ensure API docs match actual handlers.
- Ensure sample requests match required fields.
- Ensure status codes are documented.
- Ensure error responses are consistent.
- Ensure frontend-facing contract is clear.

## Rules

- Do not invent API behavior not implemented.
- Do not document fields that handler does not accept/return.
- Keep examples minimal and runnable.
- Include IDs/placeholders clearly.
- Mark auth requirements when auth exists.
- Preserve current docs style.

## Output

Update as needed:
- `docs/api_contract.md`
- `api_test.http`
- `CHAT_CONTEXT/README.md`
- `worklog.md`

## Checklist

- [ ] New/changed endpoint documented.
- [ ] Request body/query documented.
- [ ] Response documented.
- [ ] Status codes documented.
- [ ] Sample request added/updated.
- [ ] Context updated after feature completion.