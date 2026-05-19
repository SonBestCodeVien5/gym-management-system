# Agent Skill — Context Maintainer

## Role

Maintain concise project memory so future sessions can resume without context loss.

## Use when

- feature phase finishes
- `/backend-complete <feature>`
- user asks to update context
- task is interrupted and needs resumable state
- docs/logs are stale

## Must read

1. `CHAT_CONTEXT/README.md`
2. `CHAT_CONTEXT/backend_skills/SKILL.md`
3. `CHAT_CONTEXT/backend_skills/context_loading.md`
4. `CHAT_CONTEXT/backend_skills/worklog.md`
5. relevant phase files for current feature

## Responsibilities

- Update high-level project state.
- Update feature phase status.
- Record decisions and remaining risks.
- Keep worklog concise.
- Avoid copying full source or huge diffs.
- Keep resume prompts usable.

## Rules

- Summarize, do not dump.
- Keep `CHAT_CONTEXT/README.md` current but short.
- Keep `worklog.md` chronological and concise.
- Keep phase files source-of-truth for detailed status.
- Do not store secrets.
- Do not store local machine-specific noise.
- If context grows too much, compress into bullet summaries.

## Output

Update as needed:
- `CHAT_CONTEXT/README.md`
- `CHAT_CONTEXT/backend_skills/worklog.md`
- `plans/<feature>.md`
- `implementations/<feature>.md`
- `reviews/<feature>.md`
- `tests/<feature>.md`
- `resume_prompts.md` if workflow changes
- `context_loading.md` if loading rule changes

## Checklist

- [ ] Current state updated.
- [ ] Next action clear.
- [ ] Completed phase marked.
- [ ] Remaining tasks listed.
- [ ] Risks/blockers recorded.
- [ ] No secret/local noise added.
- [ ] Resume prompt still works.