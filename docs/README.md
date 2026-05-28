# Documentation Hub

This folder is the durable documentation surface for the project.

## Source Of Truth

| Need | Start here | Notes |
|---|---|---|
| Current API behavior | [api_contract.md](api_contract.md) | Reconcile with handlers, services, routes, and `api_test.http` before changing API docs. |
| Local development | [local_dev_guide.md](local_dev_guide.md) | Run and debug the Go + MongoDB backend locally. |
| Code navigation | [code_reading_guide.md](code_reading_guide.md) | Read the backend by feature flow and layer. |
| Architecture reasoning | [faq_why.md](faq_why.md) | ADR-style explanations and report-ready rationale. |
| Development history | [development_journal.md](development_journal.md) | Timeline, issues, and learning notes. |
| System analysis/report guide | [system_analysis_design_guide.md](system_analysis_design_guide.md) | Guidance for assembling the formal report. |
| Report source material | [report-materials/README.md](report-materials/README.md) | Requirement, design, UI, rollout, and conclusion drafts. |

## Boundaries

- Keep implementation-facing docs in `docs/`.
- Keep report drafts and chapter source material in `docs/report-materials/`.
- Keep resumable chat/project memory in `CHAT_CONTEXT/`.
- Keep backend phase logs in `CHAT_CONTEXT/backend_skills/` until they are summarized into durable docs.
- Treat code and the current API contract as implementation truth. Report material may describe target scope or planned work and must be reconciled before submission.

## Update Matrix

| Change | Update |
|---|---|
| Endpoint, request, response, status, API sample | `docs/api_contract.md`, `api_test.http`, related context summary |
| Local setup or debugging flow | `docs/local_dev_guide.md` |
| Architecture rationale or report explanation | `docs/faq_why.md`, `docs/system_analysis_design_guide.md`, or report material |
| Completed backend feature | Relevant backend phase files, `CHAT_CONTEXT/README.md`, durable docs touched by the behavior |
| Report-only draft work | Files under `docs/report-materials/`; do not store it in `CHAT_CONTEXT/` |

## Codex Context

Use [../CHAT_CONTEXT/README.md](../CHAT_CONTEXT/README.md) for a short project snapshot.
Use the repo-scoped skills under [../.codex/skills](../.codex/skills):

| Task | Skill |
|---|---|
| Backend plan | `$gym-plan` |
| Backend implementation | `$gym-implement` |
| Backend review | `$gym-review` |
| Backend verification | `$gym-test` |
| Feature completion/docs-context sync | `$gym-complete` |
| Frontend plan | `$gym-fe-plan` |
| Frontend implementation | `$gym-fe-implement` |
| Frontend review | `$gym-fe-review` |
| Frontend verification | `$gym-fe-test` |
| Frontend completion/docs-context sync | `$gym-fe-complete` |
| Resume backend memory | `$gym-resume` |
| Backend status | `$gym-status` |
| Durable docs | `$gym-docs` |
| Report material | `$gym-report` |
| Git/GitHub and version control | `$gym-git` |

Use `$gym-project-maintainer` only when a task crosses these surfaces and the right focused skill is
not already clear.

Prompt templates and phase protocol: [../.codex/GYM_SKILLS_WORKFLOW.md](../.codex/GYM_SKILLS_WORKFLOW.md).
