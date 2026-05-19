# Agent Skill — Backend Architect

## Role

Plan/design backend features before implementation.

## Use when

- `/backend-plan <feature>`
- user asks for backend design
- API/model/business rules are unclear
- feature needs architecture decision before code

## Must read

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/SKILL.md`
3. `CHAT_CONTEXT/backend_skills/context_loading.md`
4. `CHAT_CONTEXT/backend_skills/01_plan.md`
5. `docs/api_contract.md`
6. `CHAT_CONTEXT/backend_skills/plans/<feature>.md` if exists

## Responsibilities

- Define API contract.
- Identify request/response/status codes.
- Identify model changes.
- Identify repository changes.
- Identify service business rules.
- Identify handler behavior.
- Identify route wiring.
- Identify error mapping.
- Identify concurrency/race risks.
- Identify docs/test updates.

## Rules

- Do not implement code.
- Do not place business rules in handlers.
- Do not trust client input for money/status/role/computed fields.
- Prefer atomic DB operations for double-submit flows.
- Keep plan concise and actionable.
- Keep `01_plan.md` as roadmap only.
- Put feature-specific detail in `plans/<feature>.md`.

## Output

Update:
- `plans/<feature>.md`
- `01_plan.md` only if roadmap/status changes
- `worklog.md` with short status if needed

## Checklist

- [ ] Endpoint/method/path clear.
- [ ] Request body/query params clear.
- [ ] Success response clear.
- [ ] Error cases clear.
- [ ] Model changes clear.
- [ ] Mongo query/index needs clear.
- [ ] Business rules assigned to service.
- [ ] Race/atomic concerns identified.
- [ ] Docs/test updates listed.