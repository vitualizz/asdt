# ADR-018 — QA review step, declarative specialist authoring, and skill versioning for drift detection

**Status**: Accepted  
**Date**: 2026-06-16  
**Decided in**: Lote 3 — Workflow Design Research

---

## Context

The ASDT specialist model (ADR-006, ADR-011) defines per-specialist pipelines as orchestration plans built from subagent steps that exchange artifact-envelope contracts (ADR-003). `SKILL.md` self-loads an inline conformance gate as its `FIRST ACTION Read` (ADR-014).

Three gaps surfaced, all bound by the same hard constraints — HC-1 (no runtime logic), HC-2 (two-level subagent ceiling), HC-5 (SKILL.md must be a real committed file):

1. **Q1 — Missing conformance step in QA**: The QA specialist pipeline lacks a final `review` step that consolidates findings from all QA activities into a go/no-go verdict. The conformance gate belongs in QA — not Developer — because QA is the specialist responsible for release readiness. Developer produces the artifact; QA judges it.

2. **Q2 — Specialist authoring boilerplate**: Adding a new specialist requires hand-authoring boilerplate `SKILL.md` files that drift from one another over time as the shared structure evolves.

3. **Q3 — Unversioned skill references**: `reference_skills[]` entries are referenced by flat path with no version signal, so renames break silently and skill content can drift out of sync with no detection mechanism.

These three concerns share the same constraints and were decided together.

---

## Decision

### Q1 (A2) — Add `review` step to the QA pipeline

Add `review` as the final subagent step of the **QA** specialist pipeline at the **sonnet** tier. The `review` step belongs in QA because QA is the release-readiness specialist — Developer produces the artifact, QA judges it.

- `execution: subagent`, `model: sonnet`, `agent: analyst` (read-only)
- `inputs`: QA findings artifacts produced by prior QA steps
- `output_topic_key: qa/qa-review`
- Positioned immediately **before** `decision-preservation` (inline)
- Intra-specialist (HC-2 respected), follows SC-3 closing convention

### Q2 (B2) — Declarative specialist authoring via template + manifest

Introduce two new authoring artifacts:

- `skill/asdt-shared/skills/SKILL.md.template` — canonical template with `{{variable}}` placeholders for specialist-specific fields
- `skill/asdt-<specialist>/manifest.yaml` — per-specialist values file

`SKILL.md` is **generated once** from template + manifest at authoring time, then **checked in** as a real file. The `FIRST ACTION Read` gate and `ORCHESTRATOR GATE` blocks are **literal text** in the template — generation must reproduce them byte-for-byte so the inline gate survives (HC-5 satisfied). HC-1 is satisfied because generation is authoring-time, not runtime.

### Q3 (C2) — Optional `skill_version` frontmatter + `path@version` pinning

Add two optional, backward-compatible version signals:

- **`skill_version` frontmatter** on shared skill files (semver string, optional — absence means "unversioned")
- **`path@version` pinning** in `reference_skills[]` entries (the `@version` suffix is stripped before file resolution; it is purely additive)

No loader is introduced (HC-1 compatible). On mismatch, the on-disk file is injected anyway and the mismatch is appended to `open_items[]`:

```
"skill_version mismatch: {path} pinned @{expected} but file declares skill_version {actual|none}; proceeded with the on-disk version."
```

Files and entries without version annotations continue to behave exactly as today. Degradation only — never a hard fail.

---

## Alternatives Considered

**Q1:**
- **A1 — Add review to Developer pipeline**: Rejected — the conformance gate belongs to the specialist that judges release readiness (QA), not the one that produces code. Developer self-review creates a conflict of interest and duplicates QA's role.
- **A3 — Opt-in review via Tailored Workflow**: A conformance gate should be default-on, not opt-out.

**Q2:**
- **B1 — Keep current hand-authored SKILL.md**: Acceptable but pays ongoing boilerplate-drift cost; dominated by B2.
- **B3 — Pure YAML, no SKILL.md**: Violates HC-5/ADR-014 and HC-1; architecturally forbidden.

**Q3:**
- **C1 — Keep flat paths, no versioning**: Leaves SC-6/SC-7 drift and rename-breakage unaddressed.
- **C3 — Formal registry index**: Hand-maintained with no loader reproduces the drift it claims to fix; high coupling, violates HC-1 in spirit.

---

## Consequences

### Positive

- QA pipeline closes the conformance gap with a go/no-go verdict step before knowledge is preserved.
- Clear separation: Developer produces, QA judges — no conflict of interest.
- New specialists authored from one template + manifest, reducing boilerplate divergence.
- Skill drift and silent rename-breakage become detectable via version mismatch.
- All three changes stay inside HC-1/HC-2/HC-5 — no runtime logic, no new delegation level, `SKILL.md` stays a real file.

### Negative

- `qa-review` adds a subagent invocation (latency + token cost) to every QA run.
- Review can produce false-negative blocks or rubber-stamp passes; a weak rubric gives false confidence.
- B2's committed `SKILL.md` can silently drift from its template/manifest if edited directly — generation provenance is not enforced at runtime (HC-1), so the template stops being the single source of truth unless authoring discipline holds.
- C2's version signal is `OPTIONAL` and only degrades via `open_items`; nothing forces producers to bump `skill_version` or consumers to pin, so drift can still pass unnoticed.
- Bundling three concerns couples their future revisitation — superseding one requires re-opening a record governing the other two.

### Technical Debt

- `SKILL.md.template` + manifest schema are defined but no generator or regeneration-drift check exists yet.
- `skill_version` tooling to compare/flag versions does not exist yet — only the field and syntax are reserved.
- `skill/asdt-qa/steps/review.md` body (output format + quality gate) must be authored; this ADR specifies the contract, not the body.
- Coordinated edit (`asdt-qa/workflow.yaml` + `asdt-qa/SKILL.md` + `skill/SKILL.md` §9.2) plus the new `review.md` is builder work requiring go-ahead.

---

## Open Items

- [ ] Author `skill/asdt-qa/steps/review.md` (output schema + quality gate for go/no-go verdict)
- [ ] Update `skill/asdt-qa/workflow.yaml` to insert `review` before `decision-preservation`
- [ ] Update `skill/asdt-qa/SKILL.md` complexity table + step table to include `review`
- [ ] Update `skill/SKILL.md` §9.2 cache row (co-located with QA tier map)
- [ ] Author `skill/asdt-shared/skills/SKILL.md.template` with full placeholder set
- [ ] Author `manifest.yaml` for each specialist (developer example is specified in system-design artifact obs 1861)
- [ ] Implement B2 generator tool (authoring-time only, not runtime)
- [ ] Implement C2 comparison-at-injection in orchestrator (no loader — pure string comparison)
