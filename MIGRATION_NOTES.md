# ASDT Refactor — Migration Notes

Working document for the 5-phase refactor that replaces the current 12 shared skills / 41 sub-agent steps / persist-everything design with a smaller core. All five phases are complete, plus the post-refactor closing (§10) and the Go pass (§11) — see §12 for status.

**Anything a phase defers goes to §9, not into a partial fix.** That backlog is worked in one pass after Phase 5.

## 1. `skill/asdt-core/` is the new core (`asdt-shared/` was purged in Phase 4)

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

## 3. Go-side impacts (written in Phase 1; see §6 for the Phase 4 state)

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

Everything this phase deliberately left untouched is recorded in §9, not fixed inline.

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

`TestRegistryDrift` is still GREEN — the inline-steps region was regenerated again (Architect, Developer, and Security each lost `decision-preservation`). The same four tests from Phase 2 are red, with larger deltas; the numbers in §9 D1 are updated to the Phase 3 values.

## 6. Phase 4 — QA and UX/UI collapsed, and the purge

QA went 8 sub-agent steps → 1, UX/UI 8 → 1. Sixteen step files deleted, two written.

| Specialist | Deleted steps | New step |
|---|---|---|
| QA | `load-requirements`, `ac-validation`, `edge-case-analysis`, `test-strategy`, `test-case-generation`, `quality-report`, `performance-validation`, `review` | `steps/test-plan.md` |
| UX/UI | `feature-brief`, `design-tokens`, `information-architecture`, `user-flows`, `content-design`, `component-mapping`, `design-critique`, `ux-handoff` | `steps/ux-spec.md` |

QA's three reference skills were merged into a new `asdt-core/references/testing.md` (acceptance-criteria discipline, the edge-case catalogue, and the test-level decision). Two UX/UI behaviors were deliberately dropped rather than migrated: `design-critique` (self-assessment with no user in front of it produces confident prose and no signal) and `content-inventory` (microcopy is now written inline on the flow step that carries it). `asdt-ux-ui/skills/design-heuristics.md` went with the critique it fed; the actionable half of `asdt-ux-ui/skills/information-architecture.md` — entry path, 5–7 top-level items, progressive disclosure, destructive-last — was folded into `ux-spec.md` step 2 rather than kept as a file.

Both specialists lost their dual artifact: QA's `test-plan` + `qa-review` and UX/UI's `ux-brief` + `component-spec` are each one hand-off now.

### Final key table

Every specialist persists exactly ONE key, and that is the complete list of what ASDT writes to Engram:

| Role | Key | Written by |
|---|---|---|
| Researcher | `{project}/{change}/researcher/handoff` | `discovery` |
| PM | `{project}/{change}/pm/handoff` | `backlog` |
| UX/UI | `{project}/{change}/ux-ui/handoff` | `ux-spec` |
| Architect | `{project}/{change}/architect/handoff` | `design` |
| Developer | `{project}/{change}/developer/handoff` | `implement` |
| Security | `{project}/{change}/security/handoff` | `harden` |
| QA | `{project}/{change}/qa/handoff` | `test-plan` |
| *(organizational)* | `{project}/journal` | one line per run, when a non-obvious decision was made |

Phase 4 renames: `qa/ac-list`, `qa/ac-gaps`, `qa/edge-cases`, `qa/test-strategy`, `qa/test-cases`, `qa/test-plan`, `qa/perf-validation`, `qa/qa-review` → `qa/handoff`. `ux-ui/feature-brief`, `design-tokens`, `ia`, `flows`, `content-inventory`, `components`, `design-critique`, `ux-brief`, `component-spec` → `ux-ui/handoff`.

### What was deleted

- **`asdt-shared/` — the whole directory.** `artifact-loading.md`, `decision-preservation.md`, `intake-contract.md`, `nfr-budget.md`, `parallel-retrieval.md`, `platform-context.md`, `report.md`, `scope-definition.md`, `_README.md`. Their content either lives in `asdt-core/protocol.md` (intake, injection, persistence, degradation) or in `asdt-core/references/` (platform-context, scope-definition).
- **Every per-specialist `skills/` directory**: `asdt-architect/skills/`, `asdt-developer/skills/`, `asdt-qa/skills/`, `asdt-security/skills/`, `asdt-ux-ui/skills/`. Nine files, consolidated into `asdt-core/references/` across Phases 1 and 4.

Three files survived the purge by moving rather than dying:

| Was | Is now | Why it survived |
|---|---|---|
| `asdt-shared/skills/specialist-header.md` | `asdt-core/specialist-header.md` | the Go installer splices it into every routed SKILL.md |
| `asdt-shared/skills/executor-header.md` | `asdt-core/executor-header.md` | baked into generated agent definitions |
| `asdt-shared/skills/knowledge-recall.md` | `asdt-core/references/knowledge-recall.md` | **a third exception, decided here.** Reading ORGANIZATIONAL memory — prior work on other changes — is not covered by `protocol.md` §1, which only norms loading THIS change's hand-offs. Deleting it would have dropped a real capability, so it was repointed instead. Both header files move as-is; they are rewritten in Phase 5 |

`asdt-init/` logic was not touched. Its one broken path — `write`'s `reference_skills` pointing at the deleted `report.md` — now points at `asdt-core/protocol.md`.

### Installer paths the maintainer must repoint

`internal/installer/registry_gen.go` `specialistHeaderFragments` (≈ line 68) still lists four paths under `asdt-shared/skills/`. After this phase it should be:

```go
var specialistHeaderFragments = []string{
    "asdt-core/specialist-header.md",
    "asdt-core/protocol.md",
}
```

`parallel-retrieval.md` and `intake-contract.md` collapse into `protocol.md` §2/§4; `knowledge-recall.md` is no longer a header fragment — it is a per-specialist inline step pointing at `asdt-core/references/knowledge-recall.md`. `internal/installer/agent_adapters.go` `executorHeaderPath` (line 14) becomes `"asdt-core/executor-header.md"`.

### Go-side status after Phase 4

Eleven tests are RED. `TestRegistryDrift` is GREEN — the inline-steps region was regenerated again. The failures split cleanly in two:

**Broken by the purge — pure path repointing, fixed with the two constants above:**

| Test | What it reads |
|---|---|
| `TestExecutorHeaderCarriesFalsifiabilityCondition` | `asdt-shared/skills/executor-header.md` |
| `TestInstall_SpecialistHeaderFoldReachesInstalledSKILL` | `specialistHeaderFragments` → install fails loud |
| `TestInstall_DoesNotShipInstallerAssets` | same install path |
| `TestNuanceIsolatedFromProvenance` | same install path |
| `TestPass2InsertsWhenOptionalMarkerRemoved` | `asdt-shared/skills/platform-context.md` |
| `TestEmbeddedSkillTree` | asserts `asdt-shared/skills/executor-header.md` is embedded |
| `TestEmbeddedSharedSkillsMatchDisk` | `os.ReadDir("asdt-shared/skills")` — a `t.Fatalf` on a directory that no longer exists; repoint at `asdt-core/` or drop the test |

**The pre-existing shape/count family (§7 D1), now at final values:** `TestWorkflowSubagentStepsDeclareKnownAgentTypes` wants `analystCount` 40, actual **9** (`builderCount` 2 unchanged); `TestOptionalMarkerReadsRawSourceLine` wants 18, actual **11**; `TestOrchestrationPlanCellClassification` wants 27 rows, actual **4** — only `asdt-developer` still carries a tier table; `TestCollapseOnlyRunsAfterInsertion` still wants a `researcher/complex` row. Three noRun rows are now missing rather than two: `asdt-security/none`, `asdt-architect/simple`, `asdt-qa/trivial`.

### Final step inventory

41 sub-agent steps at the start of the refactor → **12**: pm 1, researcher 1, architect 1, developer 3, security 2, qa 1, ux-ui 1, init 2.

## 7. Phase 5 — the router, the headers, and the authoring contract

The last rewrite. `skill/SKILL.md` went 355 lines → 61, `specialist-header.md` 110 → 26, `executor-header.md` 25 → 20, `TEMPLATE.md` 246 → 96.

**The router now uses judgment, not classifiers.** Deleted outright, with no replacement: the Complexity and Risk-Surface keyword tables, the 9.1c/9.1d gates, the dagger (✝) provenance notation, the two-pass fixpoint Step List Validation algorithm, the `Tailored Workflow Generation` per-specialist table and its generated regions, the routing-plan persistence at `{project}/{change}/routing/tailored-workflow`, and the two character-for-character output strings. What replaced them is ten lines of criteria and one instruction: judge both axes, and when a tier is arguable take the lower defensible one and say so.

**Nothing about routing is persisted any more.** The tier travels as a `--tier=` argument on the command line the user copies. That removes the last write outside the seven hand-offs and the journal.

**Duplication resolved.** `workflow.yaml` is the single machine-readable source of step identity, model, inputs, and outputs. Every specialist SKILL.md dropped its File/Execution/Reads/Writes table and the "This section is the authoritative…" ritual sentence; what remains is the tier→steps mapping (only where a specialist has more than one step) and a one-line pointer at `workflow.yaml`.

`asdt-init` was left alone except for the two orphan `## Context budget` sections its step files still carried. Its five-row step table was deliberately NOT collapsed: three of its steps are inline with no step file, so that table is their only contract, and `TestInitPlanTableMatchesWorkflow` guards it. It also has no specialist-header marker region — `installer.go` documents that silent pass-through as intended.

## 8. Checklist for the Go maintainer

Everything below is Go-side and was NOT touched by any phase. Work it in one pass.

**1. Repoint two constants.**

```go
// internal/installer/registry_gen.go  (~line 68)
var specialistHeaderFragments = []string{
    "asdt-core/specialist-header.md",
    "asdt-core/protocol.md",
}

// internal/installer/agent_adapters.go  (line 14)
const executorHeaderPath = "asdt-core/executor-header.md"
```

`parallel-retrieval.md` and `intake-contract.md` collapsed into `protocol.md` §2/§4. `knowledge-recall.md` is no longer a header fragment at all — it is a per-specialist inline step pointing at `asdt-core/references/knowledge-recall.md`.

**2. Fix `skill/embedded_test.go`.** `TestEmbeddedSkillTree` asserts `asdt-shared/skills/executor-header.md` is embedded — the literal becomes `asdt-core/executor-header.md`. `TestEmbeddedSharedSkillsMatchDisk` does `os.ReadDir("asdt-shared/skills")` and `t.Fatalf`s on a directory that no longer exists: repoint it at `asdt-core/references/` or drop it. `TestRoutedSpecialistInvariants` still passes — the three verbatim blocks were preserved in every specialist, so `"SOLE orchestrator"`, `"FIRST ACTION — self-load the header"`, and both markers are all still present. The stale comment on `embedded.go:15` names `asdt-shared/skills/platform-context.md` as a path example.

**3. Regenerate `TestRegistryDrift` against the new tree.** It currently passes, because the inline-steps region was regenerated at the end of every phase. But two of its three mirror sites no longer exist in `skill/SKILL.md`: the `Tailored Workflow Generation` trivial table (deleted this phase) and, with it, `trivialTableRows`. The `Specialist Registry` site survives — the table is still there under `## Registry` — but `section5Rows` anchors on the old `## 5. Specialist Registry` heading and now parses zero rows, which is exactly how this test currently fails. Update the anchor, and drop `trivialTableRows` and its site entry entirely.

**4. Recount the shape/count assertions.** These pin the pre-refactor shape and are red; final values are in §9 D1. `analystCount` 40 → **9**; optional markers 18 → **11**; orchestration rows 27 → **4**. `TestOrchestrationPlanCellClassification` and `TestCollapseOnlyRunsAfterInsertion` need the design call recorded in D1 — only `asdt-developer` and `asdt-security` still carry a tier table, so a single-step specialist has no rows to parse.

**5. Teach the parsers `output: context`.** `workflow_models.go` and `parseRegistry` walk `workflow.yaml` for `output_topic_key`. Three steps — `developer/explore`, `developer/spec`, `security/assess` — declare `output: context` instead: they persist nothing, and their payload is retained by the orchestrator and injected into the next step. Anything that assumes every subagent step has an `output_topic_key` needs to tolerate its absence.

## 9. Deferred — the post-refactor backlog

**Standing rule for this refactor: anything a phase leaves behind gets an entry here instead of a partial fix.** These items are settled AFTER Phase 5, in one pass, for one reason — every one of them tracks a structure that phases 3–5 are still moving. Fixing them per phase means fixing them three times and reviewing churn that says nothing about whether the refactor is correct.

Each entry states what it is, exactly where, and what has to be true before it can be closed.

### D1 — Four red tests in `internal/installer` — **CLOSED by the Go pass**

Recounted against the final tree (`analystCount` 9, `builderCount` 1, optional markers 11) and the design call was made: a single-step specialist has no tier table, and the tier-table machinery was deleted rather than exempted — see §11.

<details><summary>Original entry</summary>


Red since Phase 2. Two are pure constants; two need a design call.

Values below are the Phase 3 state, re-measured after Architect, Developer, and Security collapsed.

| Test | Assertion vs actual | Closing move |
|---|---|---|
| `TestWorkflowSubagentStepsDeclareKnownAgentTypes` | `analystCount` = 40, actual **23** | Recount once, at the end. Every phase changes this number, so any value set before Phase 5 is wrong by the next commit. `builderCount` = 2 is stable (only `developer/implement` writes). |
| `TestOptionalMarkerReadsRawSourceLine` | per-specialist `# optional` counts and tree total 18, actual **14** — architect 1 (want 3), developer 3 (want 4), pm 1 (want 3), security 2 (want 1) | Same — recount at the end. Note security went UP: `assess` declares two optional inputs where the old chain declared one. |
| `TestOrchestrationPlanCellClassification` | 4 tier rows per specialist, `totalRows` = 27, actual **12** — only developer, QA, and UX/UI still carry a tier table; architect, pm, researcher, and security parse 0 rows. Also `asdt-security/none` and `asdt-architect/simple` noRun rows not found | Needs a decision, not a number: a collapsed specialist has ONE step at every tier and therefore no 4-row tier table. Decide whether single-step specialists are exempt from the parse or declare a degenerate table, then teach `parseOrchestrationPlan` that shape — including how a noRun row is expressed when there is no table to put it in. |
| `TestCollapseOnlyRunsAfterInsertion` | expects a `researcher/complex` tier row | Same root cause as the row above; closes with it. |

</details>

Blocked on: Phase 5 fixing the final tier/router semantics. Until then the counts keep moving.

### D2 — User-facing docs still describe the old artifacts — **CLOSED by the docs pass**

All pages rewritten in both languages, plus the site's data layer (`specialist-steps.ts`, `artifact-graph.ts`), the components that render it, and the i18n step catalogue.


Eleven files under `site/src/content/docs/`, in BOTH languages, still name `backlog-entry`, `nfr-targets`, `discovery-brief`, `ideation`, and the old multi-step chains:

```
en/commands.mdx            es/tutorial.mdx
en/tutorial.mdx            es/specialists/{pm,qa,architect,researcher}.md
en/specialists/{pm,qa,architect,researcher}.md
```

`site/src/content/docs/{en,es}/specialists/pm.md` and `researcher.md` need the most work — they document the six-step and three-step chains as the product's behavior, not just the artifact names. Blocked on: the last specialist collapsing (Phase 3), so the pages are rewritten once against the final shape.

### D3 — The authoring contract still teaches the old shape — **CLOSED by the docs pass**

`TEMPLATE.md` was already rewritten in Phase 5 and audited clean this pass; both READMEs now describe the hand-off contract instead of one-artifact-per-step.


- **`skill/TEMPLATE.md`** — the normative contract for adding a specialist. Still prescribes one artifact per step (`:107`, `:139`), the `summary: ""` field read by decision-preservation (`:118`), `decision-preservation` as a standard inline step (`:205`), and a worked example built from a `discover → brief` chain with `researcher/discovery` / `researcher/feasibility-brief` keys (`:212`, `:213`, `:225`, `:235`). Every one of those is false under `protocol.md`.
- **`skill/README.md`** — "one artifact per step" (`:15`, `:24`), and `decision-preservation` listed among the standard inline steps (`:54`).
- **`README.md`** (repo root) — the architecture diagram and `:81` describe "one artifact per step" saved to Engram, which is exactly the behavior this refactor removes.

Blocked on: nothing structural — but rewriting the authoring contract before the last specialist is collapsed would document a shape that is still moving. This is the highest-value item in the list: TEMPLATE.md is what the next contributor copies.

### D4 — Go test fixtures naming dead artifacts — **CLOSED by the Go pass**

Renamed to new-world names. The documented caveat held: the model-tier classification keys on the fixture's `model: haiku` value, not on the step name, so renaming the step and its comment together left the assertion intact.


`internal/installer/workflow_models_test.go` (21 hits) and `internal/setup/model_test.go` (9 hits) use `pm/feature-intake` and `pm/backlog-entry` as SYNTHETIC inline YAML fixtures. They are not references to the shipped tree, and they pass as-is — this is naming hygiene, not a defect.

One caveat when it is picked up: `model_test.go` pins `feature-intake` to a model-tier classification (`:478`, `:528`), so renaming the fixture can change what the test asserts, not just what it reads. Rename fixture and expectation together.

### D5 — Router prose left conservative on purpose

In `skill/SKILL.md`, the `Tailored Workflow Generation` per-specialist table had its PM and Researcher step-name cells refreshed in Phase 2, but the **trivial-eligible verdict prose** in those two rows was not re-reasoned. With one step at every tier, "trivial eligible" no longer means what it meant when a lower tier bought a shorter chain. Closes with the Phase 5 router pass, which owns that column's semantics.

### D6 — `asdt-shared/` still live — **CLOSED in Phase 4**

The directory and every per-specialist `skills/` directory are deleted. What survived by moving, and the two installer constants that must be repointed, are recorded in §6.

## 10. Closing — `--tier` removed

An audit found the flag carrying FOUR incompatible meanings at once: the router emitted complexity values AND risk values through it, the specialist header defined it as `quick|standard|deep` verbosity that "never controls which steps run", the Developer's table used it to gate exactly that, and Security called it the risk surface. Four readings, one flag, no way to be right.

**The flag is gone, with no substitute.** Depth now travels in natural language inside the invocation, and in its absence the specialist judges it — the same judgment the router already applies. That is only possible because the refactor left six of eight specialists with a single step: the flag's remaining jobs were output verbosity and the Developer's chain, and both are judgment calls the specialist is better placed to make than a caller typing a keyword.

```
/asdt-architect "add password reset"
/asdt-security "add password reset — touches password hashing, go deep"
/asdt-developer "rename this helper, quick one"
```

The two assessment axes are unchanged: the router still judges complexity and risk surface independently, and still states both. Only the transport changed — prose, not a flag.

Files touched: `skill/SKILL.md` (Output section only), `skill/asdt-core/specialist-header.md` (`## Tier` → `## Depth`), `skill/asdt-core/protocol.md` (§2 suppression clause, which keyed on a flag anyone could type by hand without the router ever having asked), `skill/asdt-developer/SKILL.md` (the tier table became a request→chain judgment table), `skill/asdt-security/SKILL.md`, `skill/asdt-architect/{SKILL.md,steps/design.md}`, `skill/asdt-{pm,researcher,qa,ux-ui}/SKILL.md`, `skill/TEMPLATE.md`.

Two words `tier` survive, both in the router and both correct: they name the complexity and risk-surface assessments, which still exist as the router's criteria.

**Defect D-C closed in the same pass.** Payloads travelling by context were invisible to any machine reader — `workflow.yaml` said `output: context` on the producer but nothing on the consumer. Both consumers now declare it: `developer/implement` carries `context_inputs: [dev-spec]` and `security/harden` carries `context_inputs: [assess]`. The key is defined in one line of `protocol.md` §1. `skill/embedded.go`'s path-example comment, which still named the deleted `asdt-shared/` tree, now names `asdt-core/protocol.md` — comment only, no code touched.

### Two more items for the §8 Go checklist

**6. Audit `internal/installer/preset_tiers.go` and `tier_preset_validation_test.go` as dead code.** They exist to validate tier presets against per-specialist tier tables. Only the Developer still has a chain table, and it is now keyed on what the request asks for rather than on a tier level, so the concept these two files encode may no longer exist. Read them before deleting — three of the currently-red tests live in that file, and the question is whether they should be fixed or removed.

**7. Delete any Go-side parsing of `--tier`.** Nothing in `skill/` emits or reads it any more. If the CLI or the installer parses that argument anywhere, it is now dead — and worse, it would accept a value no prompt will ever honor.

## 11. The Go pass

One pass over `internal/` and `cmd/`, executing the §8 checklist. `go build`, `go vet`, and `go test ./...` are green except for one real defect, reported below.

**Constants repointed.** `specialistHeaderFragments` is now `["asdt-core/specialist-header.md", "asdt-core/protocol.md"]` with its load-bearing order comment rewritten; `executorHeaderPath` is `"asdt-core/executor-header.md"`.

**Numbers were re-measured, not copied.** §8 said `analystCount` 9 and "builderCount 2 unchanged". The first was right; the second was wrong — the only builder step left is `developer/implement`, because `developer/test` was absorbed into it in Phase 3. Actual: analyst **9**, builder **1**. Optional markers: 11 total, matching the notes.

**Three test sites lost their target and were removed, not repointed.** `trivialTableRows` (the `Tailored Workflow Generation` trivial table), the inline-steps parity assertion (its `<!-- ASDT:GENERATED:9.2-inline-steps -->` region), and the tier-table machinery. All three mirrored router sections that Phase 5 deleted. `section5Rows` was repointed (`## 5. Specialist Registry` → `## Registry`) and renamed `registrySectionRows`; the `agents-template.md` mirror still matches the workflow.yaml roster.

**`tier_preset_validation_test.go` DELETED (26 KB).** Four of its five tests validated the router's two-pass "Step List Validation" fixpoint (`applyPass1`/`applyPass2`/`applyCollapse`) against the per-specialist tier tables — an algorithm Phase 5 deleted from the router, and tables that no longer exist. Its parser anchored on the literal `| Level | Steps |` header, so after the `--tier` removal retitled the Developer's table, ZERO specialists parsed. The fifth test survived: `TestOptionalMarkerReadsRawSourceLine` guards a property nothing else does — `yaml.v3` drops comments, so an input's optionality has to be read off the raw source line the node came from. It moved to `internal/installer/workflow_inputs_test.go` with its minimal helpers and recounted expectations.

**`preset_tiers.go` KEPT — it is not what §10 assumed.** It has nothing to do with the deleted `--tier` flag: it classifies a step's *source-default model* (haiku/sonnet/opus) into cost tiers for the setup TUI's model-preset gate, and it is consumed by `internal/setup/model.go:620,628`. `preset_tiers_test.go` goes with it. The filename collision is what made it look dead. No Go code parses a `--tier` argument anywhere; nothing to delete on that front.

**`output: context` needed no parser change.** `workflow_models.go` and `parseRegistry` only read `name`/`execution`/`model`, so they already tolerated a missing `output_topic_key`. `context_inputs:` is likewise ignored by both and can never be confused with `inputs:`. The "fail loud when a subagent step declares neither" guard belongs in an authoring test rather than in the installer, and is noted as optional below.

**Dead-tree references swept.** `grep -rn "asdt-shared" internal/ cmd/ skill/*.go` → 0. That included synthetic fixture directory names in `adapters_test.go`, `prune_test.go`, and `installer_test.go`, which used `asdt-shared` as a stand-in for "a directory with no SKILL.md" — now `asdt-core`, which is that same shape for real.

### One real defect found — FIXED

`TestNuanceIsolatedFromProvenance` was red, and it was right. The Phase 1 trim of `platform-context.md` dropped the original's explicit clarifier — "it is intentionally NOT auto-injected" — and left the `human_nuance` sentence sitting inside `## Injection Format`, the section that defines what goes into the automatic ≤500-token block. Those entries are user-authored and carry no confidence rating, so folding them in among detected values makes a person's note look like something the scan found.

Fixed by giving it its own `## Human nuance` section at the tail of the file, with the clarifier restored. The separation is now structural rather than a matter of wording: the slice the test inspects (`## Injection Format` → `## Degradation`) cannot contain it by construction.

Deleting the sentence instead was considered and rejected. It is the ONLY consumer of `human_nuance` in the tree, and the producer side is fully built in `asdt-init`: the `enrichment` inline step surveys codegraph for structurally central but non-obvious symbols, `clarify` turns them into up to three skippable questions, and `steps/write.md` routes the answers into a marked region of `knowledge.yaml` (including malformed-marker handling). Removing the reader would have left init interviewing the user for notes nothing ever reads — a worse defect than the misplacement, and a silent one.

**`go test ./...` is now green across every package.**

### Optional follow-up, not done

An authoring guard that a `subagent` step declares either `output_topic_key` or `output: context` — never neither. Nothing enforces it today; the three `output: context` steps are correct, but a fourth added without either key would install silently.

## 12. Phase status

- **Phase 1 — done.** `asdt-core/protocol.md`, six files in `asdt-core/references/`, this document. Zero existing files modified.
- **Phase 2 — done.** PM 6→1 step, Researcher 3→1, both on `{role}/handoff`; all in-tree consumers repointed. See §4.
- **Phase 3 — done.** Architect 7→1, Developer 6→3, Security 4→2; `output: context` introduced for intra-run payloads. See §5.
- **Phase 4 — done.** QA 8→1, UX/UI 8→1; `asdt-shared/` and every per-specialist `skills/` directory deleted. 41 sub-agent steps → 12. See §6.
- **Phase 5 — done.** Router rewritten on judgment instead of keyword tables; both headers rewritten; `workflow.yaml` made the single machine-readable source; `TEMPLATE.md` rewritten for the new system. See §7.
- **Go maintainer** — the checklist in §8, now seven items (§10 added two).
- **Closing — done.** `--tier` removed system-wide; depth is judged, not flagged. Defect D-C closed. See §10.
- **Docs pass — done.** Both READMEs, the site in both languages, and the site's data/components/i18n layer rewritten against the final system. D2 and D3 closed.
- **Go pass — done.** §8 checklist executed; D1, D4, D6 closed. `go test ./...` green except one real skill-tree defect awaiting a decision. See §11.
- **Post-refactor** — work the §9 backlog: the `site/` docs and the remaining red tests, alongside the §8 checklist.

Each phase appends what it defers to §9 rather than fixing it partially.
