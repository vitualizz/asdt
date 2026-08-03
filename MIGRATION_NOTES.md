# ASDT Refactor — Migration Notes

Working document for the 5-phase refactor that replaces the current 12 shared skills / 41 sub-agent steps / persist-everything design with a smaller core. Phases 1, 2, and 3 are complete — see §7 for status.

**Anything a phase defers goes to §6, not into a partial fix.** That backlog is worked in one pass after Phase 5.

## 1. `skill/asdt-core/` is the new core; `skill/asdt-shared/` is deprecated but intact

`skill/asdt-core/protocol.md` is the single mandatory shared skill of the new system. It carries, in one file, what four separate fragments carry today:

| New location | Replaces |
|---|---|
| `asdt-core/protocol.md` §1 Engram contract | `asdt-shared/skills/artifact-loading.md`, `decision-preservation.md` (the dual `mem_save` is gone — one hand-off save, plus one optional journal line) |
| `asdt-core/protocol.md` §2 Intake contract | `asdt-shared/skills/intake-contract.md` (kept near-verbatim; it was already right) |
| `asdt-core/protocol.md` §3 Executor rules | `asdt-shared/skills/executor-header.md` (trimmed; the falsifiability requirement now applies only to steps that read the codebase) |
| `asdt-core/protocol.md` §4 Injection format | `asdt-shared/skills/parallel-retrieval.md` (the cache-ledger prose collapses to "resolve once, inject once") |
| `asdt-core/protocol.md` §5 Hand-off schema | the per-specialist payload schemas scattered across specialist skills |

`skill/asdt-core/references/` holds the OPTIONAL references, consolidated from eight specialist files:

- `scope-definition.md` ← `asdt-shared/skills/scope-definition.md`
- `platform-context.md` ← `asdt-shared/skills/platform-context.md` (Reuse Guard, compact injection format, and degradation kept; the `{{ }}` template section and the line-by-line casuistry dropped)
- `api-design.md` ← `asdt-architect/skills/{api-design,scalability-analysis,architecture-review}.md`
- `owasp.md` ← `asdt-security/skills/{owasp-review,threat-modeling}.md`
- `accessibility.md` ← `asdt-ux-ui/skills/accessibility.md`
- `conventions.md` ← `asdt-developer/skills/{code-generation,test-generation}.md`

**Phase 1 deleted and edited nothing.** The originals remain the live path until Phase 4 purges them; the two trees coexist and `asdt-shared/` stays authoritative for anything still pointing at it. (Phase 2 has since collapsed PM and Researcher — see §4.)

Note on `go:embed`: `skill/embedded.go` embeds `SKILL.md asdt-*`, so `asdt-core/` ships automatically, and `installer.SiblingDestName` maps top-level entries verbatim, so it installs to `{SkillsDir}/asdt-core/` with no Go change. Two consequences worth knowing: a file added under `asdt-core/` whose name starts with `_` or `.` will silently NOT ship (that exclusion is what `TestEmbeddedSharedSkillsMatchDisk` guards for the shared tree), and `asdt-core/` has no `workflow.yaml`, so every registry walk skips it exactly as it skips `asdt-shared/`.

## 2. topic_key rename: `{role}/{step}` → `{role}/handoff`

Today each step persists its own artifact under `{project}/{change}/{role}/{step-artifact}` — `architect/system-design`, `developer/dev-spec`, `security/stride-threats`, and so on. Under the new contract there is exactly **one key per role per change**:

```
{project}/{change}/{role}/handoff      title: "{change}/{role}/handoff"   type: "decision"
```

Intermediate artifacts stop being persisted entirely; they live in the orchestrator's context for the duration of the run. The only other write is the optional one-line append to `{project}/journal`.

This rename lands with the specialist rewrites in Phases 2–3, not before. When it lands, every `inputs:` / `output_topic_key:` in the eight `workflow.yaml` files changes with it, and prose in each specialist `SKILL.md` that names an artifact key changes too. Old keys already in Engram are not migrated — a specialist that looks for a hand-off and finds none records `ASSUMED:` and proceeds, which is exactly the degradation path this design already requires.

## 3. Go-side impacts to expect in Phases 2–5

The Go code is untouched by Phase 1 and `make test` stays green. These are the places that WILL need attention as the later phases move files:

**`internal/installer/registry_gen.go`**
- `specialistHeaderFragments` (≈ line 68) hardcodes four paths under `asdt-shared/skills/` and folds them, in that exact order, into every routed specialist's header region. When `specialist-header.md`, `parallel-retrieval.md`, `intake-contract.md`, and `knowledge-recall.md` are replaced by `asdt-core/protocol.md`, this slice and the load-bearing ordering comment above it both change. A missing path here is an install-time error, so it fails loud — good.
- `specialistHeaderBeginMarker` / `specialistHeaderEndMarker` are hand-copied into `skill/embedded_test.go`; the comment says to change both copies together, and that still holds.
- `registryRenderOrder` and `inlineStepsDisplayNames` enumerate the seven routed specialist directories. They only need editing if a specialist is added, removed, or renamed — not merely rewritten.
- `parseRegistry` derives everything else from `*/workflow.yaml`, so it absorbs step-list changes for free. Note `GenerateSpecialistHeader` returns content unchanged when both markers are absent — a specialist that loses its marker region installs with no header and no error.

**`skill/embedded_test.go`**
- `TestEmbeddedSkillTree` asserts `asdt-shared/skills/executor-header.md` is present by literal path, because `agent_adapters.go` (`executorHeaderPath`, line 14) bakes it into generated agent definitions and silently generates nothing when it is absent. Both the constant and this assertion move together when the executor header folds into `protocol.md`.
- `TestEmbeddedSharedSkillsMatchDisk` walks `asdt-shared/skills/` on disk and asserts a both-directions match with the embedded FS. When Phase 4 empties or deletes that directory, this test needs to be repointed at `asdt-core/` or dropped — `os.ReadDir` on a missing directory is a `t.Fatalf`.
- `TestRoutedSpecialistInvariants` greps each routed `SKILL.md` for the literals `"SOLE orchestrator"`, `"FIRST ACTION — self-load the header"`, and both markers. Rewriting specialist headers in Phase 2 breaks this test unless the replacement prose keeps those exact phrases, or the literals are updated in step.

**`internal/installer/registry_drift_test.go`**
- `TestRegistryDrift` derives the canonical roster from the `workflow.yaml` files and asserts it against three mirror sites: the `Specialist Registry` table in `skill/SKILL.md`, the `Tailored Workflow Generation` trivial table in the same file, and `internal/installer/assets/agents-template.md` (which deliberately omits `asdt-pm`; every site omits the non-routable `asdt-init`). It also asserts the inline-steps region in `skill/SKILL.md` is byte-equal to `renderInlineStepsRegion`'s output.
- Consequence: any Phase 2–3 change to a specialist's `execution: inline` vs `subagent` classification, or to the set of routable specialists, must land in the `workflow.yaml` **and** in all three mirrors in the same commit, or this test fails. It checks set membership, routable booleans, and the inline-steps list — never the descriptive cell prose, so wording in those tables is free to change.

**`internal/installer/workflow_models.go`** walks `*/workflow.yaml` the same way and skips directories without one; step renames flow through it without edits.

## 4. Phase 2 — PM and Researcher collapsed

PM went from 6 sub-agent steps to 1, Researcher from 3 to 1, and both now persist a single hand-off. Nine step files were deleted and two created:

| Specialist | Deleted steps | New step |
|---|---|---|
| PM | `feature-intake`, `user-stories`, `success-metrics`, `scope-analysis`, `prioritization`, `backlog-entry` | `steps/backlog.md` |
| Researcher | `divergent-ideation`, `feasibility-scan`, `discovery-brief` | `steps/discovery.md` |

Both dropped the inline `decision-preservation` step: `protocol.md` §1 replaces it with the optional one-line `{project}/journal` append.

### Keys renamed

| Old topic_key | New topic_key | Where the content went |
|---|---|---|
| `pm/feature-intake` | *(gone)* | intermediate — lives in the run's context only |
| `pm/user-stories` | `pm/handoff` | `decisions[]`, in delivery order (the order IS the priority) |
| `pm/scope-analysis` | `pm/handoff` | `constraints[]` (scope in/out) |
| `pm/prioritization` | *(gone)* | replaced by story ORDER — no MoSCoW block, no separate ranking |
| `pm/nfr-targets` | `pm/handoff` | `constraints[]` (the measurable NFRs) |
| `pm/backlog-entry` | `pm/handoff` | the canonical hand-off itself |
| `researcher/ideation` | `researcher/handoff` | `decisions[]` as `rejected: {direction} — {why}` |
| `researcher/feasibility` | `researcher/handoff` | evidence folded into `decisions[]`; ungrounded verdicts become `ASSUMED:` `open_items` |
| `researcher/discovery-brief` | `researcher/handoff` | the canonical hand-off itself |

Consumers updated in the same pass: `asdt-developer/{workflow.yaml,steps/spec.md,steps/explore.md}`, `asdt-qa/{workflow.yaml,SKILL.md,steps/load-requirements.md,steps/performance-validation.md}`, `asdt-architect/{workflow.yaml,SKILL.md,steps/cost-estimation.md,steps/load-constraints.md}`, `asdt-shared/skills/{artifact-loading,knowledge-recall,report,nfr-budget,parallel-retrieval,_README}.md`, and `skill/SKILL.md`. Every DEGRADATION paragraph that named a renamed key was cut to one line — `protocol.md` §1 already norms the general case. The rule "PM is the AC authority — refine, never re-derive" survives verbatim in `asdt-developer/steps/spec.md`.

Two NFR consumers now read a FIELD rather than a dedicated artifact: `architect/cost-estimation` and `qa/performance-validation` extract the measurable NFRs out of `pm/handoff.constraints`. Their no-budget / `no-target` semantics are unchanged — a missing budget still never becomes a silent `within-budget` or `go`.

### Go-side status after Phase 2

`skill/SKILL.md`'s generated inline-steps region WAS regenerated in this phase (PM and Researcher each lost `decision-preservation`), so **`TestRegistryDrift` is green**. Four tests in `internal/installer` are RED, all of them pinning the pre-collapse shape:

| Test | Assertion | What it needs |
|---|---|---|
| `TestOrchestrationPlanCellClassification` | `wantRowCount` = 4 tier rows for `asdt-pm` / `asdt-researcher`, `totalRows` = 27 | Both SKILL.md files now carry a single-step plan with NO 4-row tier table. Decide whether a single-step specialist is exempt or declares a degenerate table — that decision belongs with the Phase 5 router pass. |
| `TestOptionalMarkerReadsRawSourceLine` | `asdt-pm` optional markers = 3, tree total = 18 | PM now declares ONE optional input (`researcher/handoff`); actual counts are 1 and 16. |
| `TestCollapseOnlyRunsAfterInsertion` | needs a `researcher/complex` tier row | Same root cause as the first row: no tier table to parse. |
| `TestWorkflowSubagentStepsDeclareKnownAgentTypes` | `analystCount` = 40 | 7 sub-agent steps were removed (PM −5, Researcher −2); actual is 33. `builderCount` = 2 is unaffected. |

The `Tailored Workflow Generation` per-specialist table in `skill/SKILL.md` had its PM and Researcher **step-name cells** refreshed (`feature-intake` → `backlog`, `divergent-ideation` → `discovery`) so no row names a step that no longer exists. The **trivial-eligible verdict prose** in those two rows was left conservative and still needs the Phase 5 router pass — with one step at every tier, "trivial eligible" no longer means what it meant when the tier bought a shorter chain.

Everything this phase deliberately left untouched is recorded in §6, not fixed inline.

## 5. Phase 3 — Architect, Developer, and Security collapsed

Architect went 7 sub-agent steps → 1, Developer 6 → 3, Security 4 → 2. Fourteen step files deleted, five written.

| Specialist | Deleted steps | New steps |
|---|---|---|
| Architect | `load-constraints`, `evaluate-approaches`, `decision-record`, `system-design`, `cost-estimation`, `risk-analysis`, `technical-handoff` | `steps/design.md` |
| Developer | `design`, `tasks`, `test` | `steps/explore.md` (kept), `steps/spec.md` (absorbed `design`), `steps/implement.md` (absorbed `test`) |
| Security | `threat-modeling`, `attack-surface`, `owasp-analysis`, `hardening-checklist` | `steps/assess.md`, `steps/harden.md` |

All three dropped the inline `decision-preservation` step. Both dual-artifact steps are gone with it — Architect's `architectural-decision` + `system-design-final` and Security's `security-findings` + `hardening-checklist` are now sections of ONE hand-off each, and the workaround comments that explained the second `output_topic_key` were deleted along with them.

Two orderings were corrected rather than carried over. Security's chain ran STRIDE *before* mapping the attack surface; `assess.md` maps the surface first, runs STRIDE over it, then cross-checks only the applicable OWASP categories. Developer's tier table gated the technical design behind `moderate`, so a `simple` change got no design at all; `spec.md` now judges the needed depth itself.

### New in this phase: `output: context`

A step whose `workflow.yaml` entry declares `output: context` instead of `output_topic_key` persists NOTHING. Its payload stays in the orchestrator's context and is injected into the next step as an `### INPUT {step-name}` block. Documented in one line in `asdt-core/protocol.md` §1.

Three steps use it: `developer/explore`, `developer/spec`, and `security/assess`. Their SKILL.md files carry the matching orchestrator instruction, because retaining and injecting those payloads is the orchestrator's job, not the sub-agent's.

### Keys renamed

| Old topic_key | New topic_key | Where the content went |
|---|---|---|
| `architect/constraints-analysis` | *(gone)* | intermediate — folded into `design`, step 1 |
| `architect/approaches` | `architect/handoff` | `decisions[]` as `rejected: {approach} — {why}` |
| `architect/adr` | `architect/handoff` | the decision IS the hand-off — no separate ADR artifact |
| `architect/system-design` | `architect/handoff` | `data_model[]` + `api_surface[]` |
| `architect/cost-estimate` | *(gone)* | the step is gone; NFR budgets stay in `pm/handoff.constraints` |
| `architect/risks` | `architect/handoff` | `risks[]` |
| `architect/architectural-decision` | `architect/handoff` | the canonical hand-off itself |
| `architect/system-design-final` | `architect/handoff` | merged — the dual-artifact split is gone |
| `developer/dev-exploration` | *(context)* | `output: context`, injected into `spec` |
| `developer/dev-spec` | *(context)* | `output: context`, injected into `implement` |
| `developer/dev-design` | *(gone)* | absorbed into `spec.md` step 4 |
| `developer/dev-tasks` | *(gone)* | S/M/L estimates and the dependency graph dropped — no consumer used them |
| `developer/dev-implementation` | `developer/handoff` | the canonical hand-off itself |
| `developer/dev-tests` | `developer/handoff` | tests are written inside `implement`, same mode and roots |
| `security/stride-threats` | *(context)* | `assess` output, injected into `harden` |
| `security/attack-surface` | *(context)* | same |
| `security/owasp-findings` | *(context)* | same |
| `security/security-findings` | `security/handoff` | `risks[]` with one-word severity |
| `security/hardening-checklist` | `security/handoff` | `constraints[]` — a section, not a second key |

Also dropped: the structured `traceability_report[]` (now one `AC not covered: {text}` line in `open_items` per uncovered AC), CVSS-lite severity (now one word: `high`/`medium`/`low`), and every numeric context budget in the rewritten steps.

Consumers updated: `asdt-qa/steps/load-requirements.md`, `asdt-shared/skills/{artifact-loading,report,decision-preservation}.md`, `skill/README.md`, and `skill/SKILL.md` (dependency list, the Architect trivial-step cell, and the generated inline-steps region). QA and UX/UI were touched ONLY on those lines — their own collapse is Phase 4.

**`implement`'s write-scope semantics are unchanged.** Mode resolution (plan-only vs writing), `allowedEditRoots`, STOP-on-out-of-scope, and the `.asdt/`-only rule for ASDT's own state survive verbatim in substance. The only change is where the roots come from: `dev-spec.files_to_create` + `files_to_modify` instead of `dev-tasks`/`dev-design`. "Writes tests, NEVER runs them" survives too — `suggested_verification.commands` remains an offer to the user.

### Go-side status after Phase 3

`TestRegistryDrift` is still GREEN — the inline-steps region was regenerated again (Architect, Developer, and Security each lost `decision-preservation`). The same four tests from Phase 2 are red, with larger deltas; the numbers in §6 D1 are updated to the Phase 3 values.

## 6. Deferred — the post-refactor backlog

**Standing rule for this refactor: anything a phase leaves behind gets an entry here instead of a partial fix.** These items are settled AFTER Phase 5, in one pass, for one reason — every one of them tracks a structure that phases 3–5 are still moving. Fixing them per phase means fixing them three times and reviewing churn that says nothing about whether the refactor is correct.

Each entry states what it is, exactly where, and what has to be true before it can be closed.

### D1 — Four red tests in `internal/installer`

Red since Phase 2. Two are pure constants; two need a design call.

Values below are the Phase 3 state, re-measured after Architect, Developer, and Security collapsed.

| Test | Assertion vs actual | Closing move |
|---|---|---|
| `TestWorkflowSubagentStepsDeclareKnownAgentTypes` | `analystCount` = 40, actual **23** | Recount once, at the end. Every phase changes this number, so any value set before Phase 5 is wrong by the next commit. `builderCount` = 2 is stable (only `developer/implement` writes). |
| `TestOptionalMarkerReadsRawSourceLine` | per-specialist `# optional` counts and tree total 18, actual **14** — architect 1 (want 3), developer 3 (want 4), pm 1 (want 3), security 2 (want 1) | Same — recount at the end. Note security went UP: `assess` declares two optional inputs where the old chain declared one. |
| `TestOrchestrationPlanCellClassification` | 4 tier rows per specialist, `totalRows` = 27, actual **12** — only developer, QA, and UX/UI still carry a tier table; architect, pm, researcher, and security parse 0 rows. Also `asdt-security/none` and `asdt-architect/simple` noRun rows not found | Needs a decision, not a number: a collapsed specialist has ONE step at every tier and therefore no 4-row tier table. Decide whether single-step specialists are exempt from the parse or declare a degenerate table, then teach `parseOrchestrationPlan` that shape — including how a noRun row is expressed when there is no table to put it in. |
| `TestCollapseOnlyRunsAfterInsertion` | expects a `researcher/complex` tier row | Same root cause as the row above; closes with it. |

Blocked on: Phase 5 fixing the final tier/router semantics. Until then the counts keep moving.

### D2 — User-facing docs still describe the old artifacts

Eleven files under `site/src/content/docs/`, in BOTH languages, still name `backlog-entry`, `nfr-targets`, `discovery-brief`, `ideation`, and the old multi-step chains:

```
en/commands.mdx            es/tutorial.mdx
en/tutorial.mdx            es/specialists/{pm,qa,architect,researcher}.md
en/specialists/{pm,qa,architect,researcher}.md
```

`site/src/content/docs/{en,es}/specialists/pm.md` and `researcher.md` need the most work — they document the six-step and three-step chains as the product's behavior, not just the artifact names. Blocked on: the last specialist collapsing (Phase 3), so the pages are rewritten once against the final shape.

### D3 — The authoring contract still teaches the old shape

- **`skill/TEMPLATE.md`** — the normative contract for adding a specialist. Still prescribes one artifact per step (`:107`, `:139`), the `summary: ""` field read by decision-preservation (`:118`), `decision-preservation` as a standard inline step (`:205`), and a worked example built from a `discover → brief` chain with `researcher/discovery` / `researcher/feasibility-brief` keys (`:212`, `:213`, `:225`, `:235`). Every one of those is false under `protocol.md`.
- **`skill/README.md`** — "one artifact per step" (`:15`, `:24`), and `decision-preservation` listed among the standard inline steps (`:54`).
- **`README.md`** (repo root) — the architecture diagram and `:81` describe "one artifact per step" saved to Engram, which is exactly the behavior this refactor removes.

Blocked on: nothing structural — but rewriting the authoring contract before the last specialist is collapsed would document a shape that is still moving. This is the highest-value item in the list: TEMPLATE.md is what the next contributor copies.

### D4 — Go test fixtures naming dead artifacts

`internal/installer/workflow_models_test.go` (21 hits) and `internal/setup/model_test.go` (9 hits) use `pm/feature-intake` and `pm/backlog-entry` as SYNTHETIC inline YAML fixtures. They are not references to the shipped tree, and they pass as-is — this is naming hygiene, not a defect.

One caveat when it is picked up: `model_test.go` pins `feature-intake` to a model-tier classification (`:478`, `:528`), so renaming the fixture can change what the test asserts, not just what it reads. Rename fixture and expectation together.

### D5 — Router prose left conservative on purpose

In `skill/SKILL.md`, the `Tailored Workflow Generation` per-specialist table had its PM and Researcher step-name cells refreshed in Phase 2, but the **trivial-eligible verdict prose** in those two rows was not re-reasoned. With one step at every tier, "trivial eligible" no longer means what it meant when a lower tier bought a shorter chain. Closes with the Phase 5 router pass, which owns that column's semantics.

### D6 — `asdt-shared/` still live

After Phase 3, only `qa` and `ux-ui` still declare the inline `decision-preservation` step, and the shared fragments remain wired into the installer. Phase 4 removes those two declarations, purges the directory, and repoints the Go references listed in §3 — this entry exists so the dependency is visible from the backlog, not to schedule work outside those phases.

## 7. Phase status

- **Phase 1 — done.** `asdt-core/protocol.md`, six files in `asdt-core/references/`, this document. Zero existing files modified.
- **Phase 2 — done.** PM 6→1 step, Researcher 3→1, both on `{role}/handoff`; all in-tree consumers repointed. See §4.
- **Phase 3 — done.** Architect 7→1, Developer 6→3, Security 4→2; `output: context` introduced for intra-run payloads. See §5.
- **Phase 4** — the same collapse for QA and UX/UI; purge `asdt-shared/` and the superseded specialist `skills/` directories; repoint the Go references listed in §3.
- **Phase 5** — router pass over the root `SKILL.md` and the registry mirrors.
- **Post-refactor** — work the §6 backlog in one pass: the red tests, the `site/` docs, and the authoring contract (`TEMPLATE.md`, the READMEs).

Each phase appends what it defers to §6 rather than fixing it partially.
