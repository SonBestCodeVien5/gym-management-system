# Backend Skills Context

Folder này chứa context dạng "skill" để tiếp tục hoàn thiện backend theo quy trình lặp: plan → implement → code review → test.

## Cách dùng nhanh

1. Đọc `CHAT_CONTEXT/README.md` để nắm trạng thái dự án.
2. Đọc agent skill protocol:
   - `SKILL.md`
   - `context_loading.md`
   - `resume_prompts.md`
3. Đọc phase skill tương ứng:
   - `01_plan.md`
   - `02_implement.md`
   - `03_code_review.md`
   - `04_test.md`
4. Với mỗi feature backend, cập nhật đủ 4 phase:
   - `plans/<feature>.md`
   - `implementations/<feature>.md`
   - `reviews/<feature>.md`
   - `tests/<feature>.md`
5. Ghi tổng hợp ngắn trong `worklog.md`.
6. Sau khi xong feature, cập nhật:
   - `docs/api_contract.md`
   - `api_test.http`
   - `CHAT_CONTEXT/README.md`

## Folder structure

- `SKILL.md`: agent skill protocol chính + command registry.
- `context_loading.md`: rule đọc context để tránh hụt/ngợp context.
- `resume_prompts.md`: prompt mẫu để resume phiên mới.
- `agent_skills/`: role-specific agent skills.
  - `backend_architect.md`: plan/design backend feature.
  - `backend_implementer.md`: implement backend feature.
  - `backend_reviewer.md`: review backend feature.
  - `backend_tester.md`: test backend feature.
  - `api_contract_keeper.md`: giữ docs/API samples khớp behavior.
  - `context_maintainer.md`: giữ project memory gọn, dễ resume.
- `plans/`: plan chi tiết cho từng feature/cycle.
- `implementations/`: log implement cho từng feature/cycle.
- `reviews/`: code review checklist + kết quả cho từng feature/cycle.
- `tests/`: test report cho từng feature/cycle.
- `01_plan.md`, `02_implement.md`, `03_code_review.md`, `04_test.md`: phase skill template/hướng dẫn thao tác.
- `worklog.md`: tổng hợp trạng thái ngắn.

## Slash-style quick start

Các command này là prompt alias, không phải CLI/native command.

Cách gọi:

```txt
Use backend skill.

/backend-implement 01_refund_pricing
```

Commands map sang agent role skills:

```txt
/backend-plan <feature>       -> agent_skills/backend_architect.md
/backend-implement <feature>  -> agent_skills/backend_implementer.md
/backend-review <feature>     -> agent_skills/backend_reviewer.md
/backend-test <feature>       -> agent_skills/backend_tester.md
/backend-complete <feature>   -> agent_skills/api_contract_keeper.md + context_maintainer.md
/backend-resume <feature>     -> agent_skills/context_maintainer.md + relevant role skill
/backend-status               -> summary files only
```

Ví dụ bắt đầu feature đầu tiên:

```txt
Use backend skill.

/backend-implement 01_refund_pricing
```

Xem chi tiết trong:
- `SKILL.md`
- `resume_prompts.md`

## Git tracking policy

Do not gitignore this folder as a whole.

Commit these files because they are project process/context:
- `SKILL.md`
- `README.md`
- `context_loading.md`
- `resume_prompts.md`
- `agent_skills/*.md`
- `01_plan.md`
- `02_implement.md`
- `03_code_review.md`
- `04_test.md`
- `plans/*.md`
- `implementations/_template.md`
- `reviews/_template.md`
- `tests/_template.md`

For solo/project-study workflow, also commit:
- `implementations/*.md`
- `reviews/*.md`
- `tests/*.md`
- `worklog.md`

Only local scratch/session noise should be ignored:
- `.scratch/`
- `tmp/`
- `session_*.md`
- `local_notes.md`

## Feature priority hiện tại

1. Refund flow & pricing rules
2. Branch nearby geo query
3. Attendance report/makeup endpoints nếu route còn thiếu
4. Auth/login + role guard
5. Validation hardening & error consistency
6. Indexes and data integrity
7. Integration tests & fixtures

## Definition of Done cho mỗi feature

- API contract rõ request/response/status code.
- Handler parse input + map error đúng.
- Service chứa business rule, không đặt rule trong handler.
- Repository thao tác MongoDB rõ, atomic khi có race/double-submit risk.
- Model có BSON/JSON tags đúng.
- Route wired trong `cmd/server/main.go`.
- `api_test.http` có request mẫu.
- `go build ./...` pass.
- Nếu có test: `go test ./...` pass.