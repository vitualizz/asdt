# ADR-017 — Trigger Phrases, Workflow Depth, and a Cost-First Phased Rollout

Date: 2026-06-15
Status: Accepted

---

## Context

A backlog of eight discoverability/quality enhancements (#1–#8) accumulated for
the specialist layer. Implementing all eight at once carried both cost and risk:
several depend on upstream primitives ASDT does not yet own (statusline hooks,
behavior-benchmark harness, grader tooling), and bundling them would couple
shippable, self-contained changes to blocked ones.

A codebase audit of the eight items found a clean split:

- Three items (**#5 workflow depth**, **#8 trigger phrases**, **#1 CI phrase
  check**) merely EXTEND mechanisms ASDT already has. `## Tailored Workflow`
  block parsing already exists in `specialist-header.md` and is documented in
  `skill/SKILL.md` §9.2 — adding a `depth:` verbosity dial is a sibling field,
  not a new system. Specialist frontmatter already carries structured YAML keys
  (`shared-skills:`) — `trigger_phrases:` is another metadata list. The embedded
  FS already has an invariants test (`TestEmbeddedSkillTree`) — a routed-phrase
  assertion is a second test in the same file with the same imports.

- Five items (**#2 comment convention**, **#6 portability doc**, **#7
  statusline**, **#3 behavior benchmark**, **#4 grader tests**) require new
  infrastructure or external integration points that ASDT does not yet control.

This split motivates a **Cost-First Phased** strategy: ship the three
zero-new-infrastructure items now, defer the five blocked ones with their
blockers recorded so they are never silently dropped.

## Decision

Adopt a two-phase rollout. **Phase 1 is implemented in this change. Phase 2 is
deferred** with explicit upstream blockers.

### Phase 1 — implemented now

**#5 — Workflow depth dial.** The `## Tailored Workflow` block MAY carry an
OPTIONAL `depth: quick|standard|deep` field controlling per-step OUTPUT
VERBOSITY only. It is orthogonal to `complexity` and `risk_surface` (which gate
WHICH steps run) and to the `steps:` list itself.

- `depth` omitted ⇒ `standard` (no behavior change — fully backward compatible).
- `quick` ⇒ minimum viable output: collapse enumerations to the most important
  item(s), skip OPTIONAL fields, terse rationale.
- `standard` ⇒ the current default output volume.
- `deep` ⇒ exhaustive alternatives/edge-cases/rationale within the existing
  per-step context budget.
- **Invariant**: `depth` NEVER overrides hard schema-required fields — those are
  always emitted regardless of depth.

Documented in `skill/asdt-shared/skills/specialist-header.md` (the Tailored
Workflow detection block) and `skill/SKILL.md` (§9.2 format spec, with the field
shown in the §6 and §8 emitted templates).

**#8 — Trigger phrases.** Each of the 7 routed specialist `SKILL.md` files gains
a top-level `trigger_phrases:` frontmatter key — five lowercase, punctuation-free
natural-language phrases that signal when the specialist applies. These give the
meta-orchestrator (and humans scanning the registry) phrase-level routing hints
without changing any execution contract. The key is metadata, mirroring the
existing `shared-skills:` convention; no loader resolves it.

**#1 — CI phrase check.** A new `TestRoutedSpecialistInvariants` in
`skill/embedded_test.go` iterates the 7 routed specialist directories, reads each
`SKILL.md` from the embedded FS, and asserts both load-bearing header invariants
are present: `SOLE orchestrator` and `FIRST ACTION — self-load the header`
(em-dash U+2014). This guards the launch contract against a silent drop during
future edits. It reuses the existing test file's imports — no new dependency.

### Phase 2 — deferred (with upstream blockers)

- **#2 — Comment convention.** Deferred. Needs a ratified inline-annotation
  syntax across SKILL bodies before tooling can rely on it.
- **#6 — Portability doc.** Deferred. Blocked on a stable host-harness capability
  matrix; documenting portability before the matrix freezes would invite churn.
- **#7 — Statusline.** Deferred. Requires a host statusline/hook integration
  point ASDT does not yet own.
- **#3 — Behavior benchmark.** Deferred. Requires a behavior-benchmark harness
  (scenario corpus + scoring) that does not exist yet.
- **#4 — Grader tests.** Deferred. Depends on #3's harness as its substrate;
  cannot land before the benchmark exists.

## Alternatives Considered

**Big-bang all eight at once** — rejected. Couples shippable items to blocked
ones, inflating cost and stalling Phase 1 behind missing infrastructure.

**`depth` as a new top-level orchestration axis** — rejected. Depth is a
verbosity dial, not a step-gating axis; modeling it alongside `complexity`/
`risk_surface` would imply it changes which steps run, which it must never do.

**`trigger_phrases` as a separate registry file** — rejected. Frontmatter
co-locates the phrases with the specialist they describe, reusing the existing
`shared-skills:` metadata pattern; a separate file adds a fourth sync hazard.

## Consequences

**Positive:**

- Three self-contained improvements ship immediately with zero new infrastructure.
- `depth` is fully backward compatible (omitted ⇒ standard).
- Trigger phrases improve routing legibility for orchestrator and humans alike.
- The CI check turns a silent-drop class of regression into a hard test failure.
- Phase 2 blockers are recorded, so deferred items are never silently lost.

**Negative:**

- `trigger_phrases:` adds a 7-file frontmatter surface to keep curated.
- `depth` adds a per-step authoring decision (mitigated by the safe default).
- Phase 2 remains outstanding until its upstream primitives exist.

**Debt:**

- Trigger-phrase wording is unaudited prose — no tooling validates phrase quality.
- `depth` semantics are specified in prose, not enforced by a schema.

## Related

- ADR-011 — Specialist Pipelines as Orchestration Plans: the orchestration model
  the `depth` dial and Tailored Workflow block obey.
- ADR-014 — SKILL.md Self-Load Inline Gate: defines the `FIRST ACTION` and
  `SOLE orchestrator` invariants that #1's CI check now enforces.
- ADR-016 — Init Empowerment and the Researcher Specialist: established the
  seventh routed specialist now covered by the trigger-phrase and CI work.
