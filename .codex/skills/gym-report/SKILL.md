---
name: gym-report
description: Prepare formal report source material for this gym management system. Use when Codex is asked to draft, refine, organize, or reconcile report chapters and report materials, or respond to an explicit `$gym-report` or `/gym-report` request.
---

# Gym Report

## Read First

1. Read `docs/report-materials/README.md`.
2. Read the report material files relevant to the requested chapter.
3. Read `docs/system_analysis_design_guide.md` for report structure and traceability guidance.
4. Read `docs/faq_why.md`, code, and `docs/api_contract.md` when a report claim needs technical
   justification or current-implementation evidence.

## Focus

- Turn project evidence into report-ready analysis, design, implementation, testing, and future-work
  material.
- Preserve the difference between requirement target, planned behavior, and implemented behavior.
- Keep report drafts in `docs/report-materials/`.
- Use durable docs and code as evidence; do not treat chat context as final report prose.

## Output Rules

- Update the relevant report material file and its index when a new chapter source is added.
- Label current versus planned behavior when the distinction matters.
- Avoid copying large API tables or raw code into the report unless the user explicitly needs them.
- Record report-specific improvements in report materials, not backend phase memory.
