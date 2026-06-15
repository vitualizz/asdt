---
title: Specialist Comparison
description: Compare all ASDT specialists — PM, Architect, Developer, QA, Security, UX/UI, Researcher — by trigger condition, artifact output, and when NOT to invoke them.
order: 27
locale: en
---

# Specialist Comparison

Use this table to pick the right specialist in under 60 seconds.

| Specialist | Command | Invoke when… | Produces | Do NOT use when… |
|---|---|---|---|---|
| **Product Manager** | `/asdt-pm` | The request is vague or user-facing, scope isn't locked, or user stories don't exist yet | `pm/backlog-entry` — feature summary, ordered user stories with AC, in/out scope, risk flags | You already have a clear backlog entry — re-running PM regenerates stories from scratch |
| **Architect** | `/asdt-architect` | The solution touches service boundaries, data models, or API contracts, or has two viable technical approaches worth documenting | `architectural-decision` (ADR) + `system-design` (data model, API surface, service boundaries) | You need implementation code, test plans, or UX specs — Architect produces decisions, not code |
| **Developer** | `/asdt-developer` | The shape of the solution is settled and you need an ordered implementation plan or production code written to the repo | `developer/dev-implementation` — ordered file manifest and code plan | You haven't locked scope or architecture yet — Developer will implement against ambiguous requirements |
| **QA Engineer** | `/asdt-qa` | Code is ready for review, AC exists but hasn't been validated, or you need systematic edge case and boundary coverage | `test-plan` — AC coverage %, uncovered gaps, full Given/When/Then test case list, quality verdict | You want executable test code written — QA produces test specifications, not runnable code |
| **Security** | `/asdt-security` | The feature touches auth, sessions, PII, external integrations, webhooks, or new public API endpoints | `security-findings` (severity-rated, CWE-referenced) + `hardening-checklist` (must-fix vs can-defer) | You want implementation code or architectural decisions — Security produces findings and checklists only |
| **UX/UI Designer** | `/asdt-ux-ui` | A new screen or feature-level UI needs design before implementation begins, or user flows need mapping | `ux-brief` (flows, IA, success criteria) + `component-spec` (inventory of reused/extended/new components) | The screen has already been built — a UX spec delivered after implementation is too late to shape it |
| **Researcher** | `/asdt-researcher` | The problem is fuzzy or open-ended — you need discovery and framing before requirements can be written | `researcher/discovery-brief` — problem framing, divergent idea set, feasibility scan, single recommended direction | You already have a defined problem statement — Researcher explores; it does not produce user stories or ADRs |

### Recovery: chose the wrong specialist

Stop the current run. Invoke the correct specialist directly — each specialist reads prior artifacts from the shared knowledge base automatically. Nothing from the interrupted run is lost.

Example: if you ran `/asdt-developer` before creating an ADR, run `/asdt-architect` to produce the decision record, then re-invoke `/asdt-developer`. The Developer will read the ADR automatically on its next run.

See [Troubleshooting](/asdt/docs/troubleshooting) for step-by-step recovery procedures.
