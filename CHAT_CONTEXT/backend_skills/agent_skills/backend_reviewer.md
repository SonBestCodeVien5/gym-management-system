# Agent Skill — Backend Reviewer

## Role

Review backend implementation for correctness, architecture, API consistency, and data safety.

## Use when

- `/backend-review <feature>`
- implementation phase is done
- user asks for code review
- feature needs risk/security/data-integrity check

## Must read

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/SKILL.md`
3. `CHAT_CONTEXT/backend_skills/context_loading.md`
4. `CHAT_CONTEXT/backend_skills/03_code_review.md`
5. `CHAT_CONTEXT/backend_skills/plans/<feature>.md`
6. `CHAT_CONTEXT/backend_skills/implementations/<feature>.md`
7. `CHAT_CONTEXT/backend_skills/reviews/<feature>.md`

## Responsibilities

- Check implementation against plan.
- Check Clean Architecture boundaries.
- Check API contract and docs alignment.
- Check error mapping.
- Check validation and business rules.
- Check MongoDB query/update correctness.
- Check race/double-submit risks.
- Record issues and fixes.

## Review rules

- Existing tests failing after change usually means implementation is wrong.
- Do not change tests to match broken behavior unless user explicitly asks.
- Do not accept handler-level business rules.
- Do not accept client-controlled money/status/role/computed fields.
- Prefer atomic update or unique index for race-prone flows.
- Mark severity for issues: low, medium, high, critical.

## Output

Update:
- `reviews/<feature>.md`

Record:
- review summary
- checklist result
- issues found
- fixes applied
- remaining risks
- handoff notes for tester

## Checklist

- [ ] Code compiles or build issue recorded.
- [ ] Handler only parses/responds.
- [ ] Service owns business rules.
- [ ] Repository owns MongoDB access.
- [ ] Models match BSON/JSON contract.
- [ ] Errors map to correct HTTP status.
- [ ] Atomic updates used where needed.
- [ ] API contract updated if behavior changed.
- [ ] API samples updated if endpoint changed.
- [ ] No secrets or local-only data committed.