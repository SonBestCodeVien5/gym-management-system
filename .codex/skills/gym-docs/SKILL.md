---
name: gym-docs
description: Maintain durable project documentation for this gym management system. Use when Codex is asked to refactor docs, align API/local-dev/code-reading documentation with code, clean doc links or duplication, or respond to an explicit `$gym-docs` or `/gym-docs` request.
---

# Gym Docs

## Read First

1. Read `docs/README.md`.
2. Read `README.md` when root guidance or feature summary is touched.
3. Read `docs/api_contract.md` and code when documenting exact backend behavior.
4. Read `CHAT_CONTEXT/README.md` only when docs/context boundaries or resume state are touched.

## Focus

- Keep durable docs in `docs/` and implementation truth anchored to code and API contract.
- Remove duplication when one index, contract, or guide already owns the information.
- Keep links, document boundaries, and update rules clear.
- Avoid moving report drafts or phase memory back into durable docs.

## Output Rules

- Update `docs/README.md` when doc entrypoints or boundaries change.
- Update `README.md` only for repo-level onboarding and feature summary needs.
- Preserve concise docs; point to source-of-truth docs instead of copying long endpoint tables.
- Search for stale paths after moving or deleting docs.
