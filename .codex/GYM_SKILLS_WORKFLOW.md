# Gym Codex Skills Workflow

This guide is for using the repo-scoped skills under `.codex/skills/`.

## 1. Basic Invocation

Pick a skill from the `$` picker or mention it directly in the prompt:

```txt
$gym-plan
Plan Cycle 04 auth/login + role guard from the current backend context.
```

Prefer one focused skill per phase. Add the feature key or target files when you know them.

```txt
$gym-implement
Implement `04_auth_role_guard` from its plan. Keep docs and implementation notes aligned.
```

If the request crosses docs, report material, backend memory, and code without a clear phase, use:

```txt
$gym-project-maintainer
Route this task to the right docs/context/backend workflow before editing.
```

## 2. Skill Map

| Skill | Use it for |
|---|---|
| `$gym-plan` | Design a backend feature before coding |
| `$gym-implement` | Implement an existing backend plan |
| `$gym-review` | Review backend changes and risks |
| `$gym-test` | Build, test, manual API verification, and test notes |
| `$gym-complete` | Final docs/API/context synchronization for a feature |
| `$gym-resume` | Resume interrupted backend work from memory and git state |
| `$gym-status` | Summarize roadmap/current feature status without editing |
| `$gym-docs` | Durable docs refactor/update |
| `$gym-report` | Formal report source material |
| `$gym-git` | Git, GitHub, commit, branch, diff, PR, and version-control workflow |
| `$gym-project-maintainer` | Cross-surface routing when one focused skill is not clear |

## 3. Standard Backend Phase Order

Use this order for a normal backend feature:

1. `$gym-status` when you need to know the current roadmap or next cycle.
2. `$gym-plan` to create or refine the feature plan.
3. `$gym-implement` to change code from the plan.
4. `$gym-review` to inspect correctness, regressions, risks, and missing tests.
5. `$gym-test` to record verification evidence.
6. `$gym-complete` to align API docs, API samples, worklog, and chat context.
7. `$gym-git` when you want to inspect diff, split/stage commits, commit, push, or prepare PR work.

Do not skip from plan to complete. If the user intentionally skips a phase, keep the skipped state
visible in the phase notes and final summary.

## 4. Backend Prompt Templates

### Plan

```txt
$gym-plan
Plan feature `<feature-key>`.
Check current code and API behavior first.
Do not implement yet.
```

### Implement

```txt
$gym-implement
Implement `<feature-key>` from the existing plan.
Run the required verification commands and update implementation notes.
```

### Review

```txt
$gym-review
Review `<feature-key>`.
Find bugs, regressions, data-integrity risks, and missing tests first.
```

### Test

```txt
$gym-test
Verify `<feature-key>`.
Run build/tests when feasible and record manual API or skipped prerequisites clearly.
```

### Complete

```txt
$gym-complete
Finalize `<feature-key>`.
Align API contract, API samples, worklog, and chat context with the code that exists.
```

### Resume

```txt
$gym-resume
Resume the current backend task from chat context, backend memory, and git state.
```

### Status

```txt
$gym-status
Summarize current backend status, blockers, verification state, and next step.
```

## 5. When A Phase Gets Stuck

| Problem | Call |
|---|---|
| Unsure what cycle is next | `$gym-status` |
| Work was interrupted or chat context changed | `$gym-resume` |
| Plan is weak, incomplete, or disagrees with code | `$gym-plan` |
| Implementation is diverging from plan | `$gym-plan` first, then `$gym-implement` |
| Code exists but risk is unclear | `$gym-review` |
| Build/test/manual API evidence is missing | `$gym-test` |
| Code works but docs/context are stale | `$gym-complete` |
| Docs are duplicated or confusing | `$gym-docs` |
| Report text mixes planned and implemented behavior | `$gym-report` |
| Diff/branch/commit/PR handling is unclear | `$gym-git` |

## 6. Docs And Report Workflow

Use docs/report skills outside the backend phase loop when the task is primarily documentation:

```txt
$gym-docs
Refactor the local dev and code reading docs without duplicating the API contract.
```

```txt
$gym-report
Prepare the analysis and design report material for attendance and sessions.
Mark planned behavior separately from implemented behavior.
```

Boundaries:

- Durable project docs live in `docs/`.
- Report source material lives in `docs/report-materials/`.
- Chat resume memory lives in `CHAT_CONTEXT/`.
- Backend feature memory lives in `CHAT_CONTEXT/backend_skills/`.

## 7. Git And GitHub Workflow

Use `$gym-git` after code/docs work when version-control work is the main task:

```txt
$gym-git
Inspect the current diff and suggest a clean commit split.
```

```txt
$gym-git
Stage only the skill and docs refactor changes, then show the staged scope before committing.
```

```txt
$gym-git
Prepare a PR summary for the current branch with verification and residual risks.
```

Do not assume commit or push intent. Ask for those operations explicitly when you want them.

## 8. Skill Loading Notes

- Skills are repo-scoped here: `.codex/skills/<skill-name>/SKILL.md`.
- Use the `$` picker when possible. A picked skill is injected into the chat as a skill block.
- If a newly created skill does not appear in the picker, reload or reopen the Codex session for the
  repo and try again.
- `$skill-creator` is for creating or updating skills.
- `$skill-installer` is for installing skills from curated/GitHub sources into Codex skill storage;
  it is not required to use the repo-local gym skills.

## 9. Recommended Feature Handoff

When a feature phase ends, leave the next skill obvious:

| Finished | Next |
|---|---|
| `$gym-plan` | `$gym-implement` |
| `$gym-implement` | `$gym-review` |
| `$gym-review` | `$gym-test` |
| `$gym-test` | `$gym-complete` |
| `$gym-complete` | `$gym-git` if version-control work is requested |

Keep handoff notes in the feature memory files, not in long chat-only instructions.
