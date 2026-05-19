# Resume Prompts

Dùng file này để bắt đầu phiên mới không hụt context.

## Rule chung

Không yêu cầu agent đọc toàn bộ project. Chỉ đọc:

1. global context
2. skill protocol
3. phase skill hiện tại
4. feature phase files liên quan
5. source files cần thiết

## Slash-style commands

Các command này là prompt alias cấp project, không phải CLI/native command.

Cách gọi chuẩn:

```txt
Use backend skill.

/backend-implement 01_refund_pricing
```

Danh sách command:

```txt
/backend-plan <feature>
/backend-implement <feature>
/backend-review <feature>
/backend-test <feature>
/backend-complete <feature>
/backend-resume <feature>
/backend-status
```

### `/backend-plan <feature>`

```txt
Use backend skill.

/backend-plan <feature>
```

Meaning:
- Read global context + skill protocol + planning skill.
- Create/update `plans/<feature>.md`.
- Keep `01_plan.md` as roadmap only.
- Do not implement.

### `/backend-implement <feature>`

```txt
Use backend skill.

/backend-implement <feature>
```

Meaning:
- Read global context + skill protocol + implement skill.
- Read `plans/<feature>.md`.
- Read/update `implementations/<feature>.md`.
- Implement feature by Clean Architecture layers.
- Run gofmt and `go build ./...`.
- Do not review/test unless asked.

### `/backend-review <feature>`

```txt
Use backend skill.

/backend-review <feature>
```

Meaning:
- Read global context + skill protocol + review skill.
- Read plan + implementation + review files.
- Inspect changed source files.
- Update `reviews/<feature>.md`.
- Do not test unless asked.

### `/backend-test <feature>`

```txt
Use backend skill.

/backend-test <feature>
```

Meaning:
- Read global context + skill protocol + test skill.
- Read plan + implementation + review + test files.
- Run build/test commands.
- Update `tests/<feature>.md`.

### `/backend-complete <feature>`

```txt
Use backend skill.

/backend-complete <feature>
```

Meaning:
- Finalize feature.
- Update project docs/context:
  - `CHAT_CONTEXT/README.md`
  - `worklog.md`
  - `docs/api_contract.md`
  - `api_test.http`

### `/backend-resume <feature>`

```txt
Use backend skill.

/backend-resume <feature>
```

Meaning:
- Resume interrupted work.
- Read relevant phase files.
- In Act mode, inspect `git status --short`.
- Continue from last known state.

### `/backend-status`

```txt
Use backend skill.

/backend-status
```

Meaning:
- Summarize roadmap and current statuses.
- Read summary files only.

## Implement feature prompt

```txt
Read these first:
1. CHAT_CONTEXT/README.md
2. CHAT_CONTEXT/backend_skills/README.md
3. CHAT_CONTEXT/backend_skills/SKILL.md
4. CHAT_CONTEXT/backend_skills/02_implement.md
5. CHAT_CONTEXT/backend_skills/plans/<feature>.md
6. CHAT_CONTEXT/backend_skills/implementations/<feature>.md

Then implement <feature> following Clean Architecture.
Read only relevant source files.
Update implementations/<feature>.md as you work.
After implementation, run gofmt and go build ./...
Do not start review/test phase unless asked.
```

Example for Cycle 01:

```txt
Read these first:
1. CHAT_CONTEXT/README.md
2. CHAT_CONTEXT/backend_skills/README.md
3. CHAT_CONTEXT/backend_skills/SKILL.md
4. CHAT_CONTEXT/backend_skills/02_implement.md
5. CHAT_CONTEXT/backend_skills/plans/01_refund_pricing.md
6. CHAT_CONTEXT/backend_skills/implementations/01_refund_pricing.md

Then implement Cycle 01 Refund Flow & Pricing Rules following Clean Architecture.
Read only relevant source files.
Update implementations/01_refund_pricing.md as you work.
After implementation, run gofmt and go build ./...
Do not start review/test phase unless asked.
```

## Review feature prompt

```txt
Read these first:
1. CHAT_CONTEXT/README.md
2. CHAT_CONTEXT/backend_skills/README.md
3. CHAT_CONTEXT/backend_skills/SKILL.md
4. CHAT_CONTEXT/backend_skills/03_code_review.md
5. CHAT_CONTEXT/backend_skills/plans/<feature>.md
6. CHAT_CONTEXT/backend_skills/implementations/<feature>.md
7. CHAT_CONTEXT/backend_skills/reviews/<feature>.md

Then review <feature>.
Check changed source files and API contract.
Update reviews/<feature>.md with issues, fixes, risks.
Do not start test phase unless asked.
```

## Test feature prompt

```txt
Read these first:
1. CHAT_CONTEXT/README.md
2. CHAT_CONTEXT/backend_skills/README.md
3. CHAT_CONTEXT/backend_skills/SKILL.md
4. CHAT_CONTEXT/backend_skills/04_test.md
5. CHAT_CONTEXT/backend_skills/plans/<feature>.md
6. CHAT_CONTEXT/backend_skills/implementations/<feature>.md
7. CHAT_CONTEXT/backend_skills/reviews/<feature>.md
8. CHAT_CONTEXT/backend_skills/tests/<feature>.md

Then test <feature>.
Run go build ./... and go test ./... if available.
Update tests/<feature>.md with command results and manual API checklist.
```

## Plan new feature prompt

```txt
Read these first:
1. CHAT_CONTEXT/README.md
2. CHAT_CONTEXT/backend_skills/README.md
3. CHAT_CONTEXT/backend_skills/SKILL.md
4. CHAT_CONTEXT/backend_skills/01_plan.md
5. docs/api_contract.md

Then plan next backend feature.
Create/update plans/<feature>.md.
Keep 01_plan.md as roadmap only.
Do not implement.
```

## Resume interrupted implementation prompt

```txt
Resume interrupted backend implementation.

Read:
1. CHAT_CONTEXT/README.md
2. CHAT_CONTEXT/backend_skills/README.md
3. CHAT_CONTEXT/backend_skills/SKILL.md
4. CHAT_CONTEXT/backend_skills/02_implement.md
5. CHAT_CONTEXT/backend_skills/plans/<feature>.md
6. CHAT_CONTEXT/backend_skills/implementations/<feature>.md

Then inspect git diff/status and relevant source files.
Continue from last completed step.
Update implementations/<feature>.md before completion.
```

## Complete feature context update prompt

```txt
Feature <feature> is implemented/reviewed/tested.

Read:
1. CHAT_CONTEXT/README.md
2. CHAT_CONTEXT/backend_skills/SKILL.md
3. CHAT_CONTEXT/backend_skills/worklog.md
4. CHAT_CONTEXT/backend_skills/plans/<feature>.md
5. CHAT_CONTEXT/backend_skills/implementations/<feature>.md
6. CHAT_CONTEXT/backend_skills/reviews/<feature>.md
7. CHAT_CONTEXT/backend_skills/tests/<feature>.md

Then update:
- CHAT_CONTEXT/README.md
- CHAT_CONTEXT/backend_skills/worklog.md
- docs/api_contract.md
- api_test.http if needed
```
