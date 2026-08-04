---
name: asdt-init
description: "Explicit-invocation-only project setup — invoke it by name to detect the stack, collect the configuration decisions, and WRITE .asdt/config.yaml plus the knowledge files; never fire it from a vague 'set this project up', because it changes files on disk."
user-invocable: true
specialist-id: asdt-init
metadata:
  author: "Lee Palacios (vitualizz)"
  version: "1.0"
---

# ASDT Init Specialist

## Role
Initialize ASDT for the current project: detect the stack, collect the configuration decisions with the human, and write `.asdt/config.yaml` and the `.asdt/knowledge/` files every other specialist reads. It does NOT plan, design, or implement anything — it runs once before the other specialists, and again only on a deliberate recalibration.

## Orchestration Plan

**Setup-class flow — no tiers, no routing.** A user invokes init by name to scaffold a project; no feature request ever routes to setup, so init is deliberately absent from the routing tables in `skill/SKILL.md` and that absence is intentional, not a gap. There are no complexity tiers and no tier→step table: the same five steps always run in order. The gates stay inline with the orchestrator because they need its own tool list and can pause for the human; the `write` sub-agent owns every file write.

| Step | File | Execution | Reads | Writes |
|------|------|-----------|-------|--------|
| knowledge-gate | *(inline — no step file)* | inline | *(orchestrator's own tool list)* | *(no artifact — Engram presence gate)* |
| explore | steps/explore.md | subagent | *(raw project tree — `inputs: []`)* | `init/stack-detection` |
| enrichment | *(inline — no step file)* | inline | *(codegraph MCP survey — no artifact)* | *(no artifact — injects `nuance.*` ambiguities into clarify)* |
| clarify | *(inline — no step file)* | inline | `init/stack-detection.ambiguities[]` + enrichment's `nuance.*` | *(no artifact — injects `answers{}` into write's prompt)* |
| write | steps/write.md | subagent | `init/stack-detection` + `### CLARIFY ANSWERS` | `init/write-summary` |

Setup-class flow: this table stays, unlike the routed specialists', because three of these steps have NO step file — the sections below are their only contract. `workflow.yaml` stays authoritative if the two ever disagree. The detection and file-writing mechanics belong to `steps/explore.md` and `steps/write.md` and are not restated here.

## knowledge-gate — Engram presence (inline)

Resolve Engram presence yourself, first, before delegating anything. "Does THIS session have Engram's memory tools" is a question about your own tool list — the one every other specialist will depend on. A sub-agent carries a different tool list, so asking it risks a false "absent"; checking yourself costs nothing, since you are inspecting your own tools, not running commands. Look for `mem_save`, `mem_search`, `mem_context` (Claude Code exposes them as `mcp__plugin_engram_engram__mem_*`; other hosts may use a different prefix or none).

- **Present** → say so and continue, passing "Engram confirmed present" into explore's prompt as an established fact.
- **Absent** → tell the user Engram is required for ASDT's cross-session memory and is not reachable in this session, explain how to connect it, and STOP. NEVER write `.asdt/config.yaml` with `provider: engram` when the provider is not actually there — that silently points every future specialist at a memory backend that does not exist.

**Recalibration gate, in the same breath.** Before launching explore, check for an existing `.asdt/config.yaml`. If it exists, this project was already initialized: settle recalibrate-vs-leave with the user before any detection runs. "Leave as-is" → stop without launching explore. Only a chosen recalibration, or a fresh project with no config, proceeds.

## enrichment — surface the non-obvious (inline)

Enrichment surfaces the handful of symbols a newcomer would misread — structurally central to the codebase yet non-obvious from the name alone — and turns each into one skippable `nuance.*` question. It runs inline because codegraph is an MCP tool on your own tool list, and its questions join clarify's. It is positive-evidence-only: nothing structurally interesting surfaces, no questions. It never fails and never blocks.

Resolve the capability ladder against your OWN tool list, top rung first:

1. **codegraph present** (any `mcp__codegraph__*` tool). Call `codegraph_explore` with a survey question — "which symbols are the most central to this codebase yet the least self-explanatory from their name?" — and pick ≤ 3 chunks by highest combined caller+callee degree, skipping trivially-named getters, setters, and plain DTOs. Tie-break by file path, then symbol name. A candidate qualifies only at **architecture altitude**: its non-obviousness spans ≥ 2 files, a module boundary, or a cross-cutting invariant. Reject single-function gotchas — a single-function detail is deliberately not asked. Classify each survivor as exactly one of `architectural`, `repo_practice`, or `inconsistency_to_review`.
2. **codegraph absent, `tree-sitter --version` answers.** Syntax is not centrality — a parse tree cannot tell you which symbol matters most. Surface nothing from this rung and note it in `open_items`, e.g. `enrichment: codegraph absent, tree-sitter only — centrality unavailable, skipped surfacing`.
3. **Neither present.** Skip enrichment, note `enrichment: no code-intelligence tooling — skipped`, and proceed.

Each chosen chunk becomes one `Ambiguity`: `field` is `nuance.<type>.<slug>` (e.g. `nuance.architectural.replaceMarkerRegion`), `question` asks the human what makes the symbol non-obvious — its role, an invariant, a gotcha a newcomer would miss — with `options: []`, `default: ""`, and `skippable: true` ALWAYS. A `nuance.*` ambiguity is NEVER a blocking open item. Three is the ceiling here; across the whole clarify pass — explore's ambiguities plus these — the ceiling is seven.

## clarify — collect the decisions (inline)

explore returns `init/stack-detection` carrying `ambiguities[]`, one per low/medium-confidence or genuinely ambiguous field; enrichment adds its `nuance.*` entries. Clarify resolves them with the human and hands the answers to write. It runs inline because only an inline step can pause for a question — the `write` sub-agent cannot ask anything, and it, not you, performs every file write.

1. Ask explore's ambiguities first, one question at a time, offering `options` when present and taking a free-form value otherwise. Then ask the ≤ 3 `nuance.*` questions (free-form, all skippable), with a skip-all shortcut so the human can decline the whole block in one answer. Finally ask the ONE standing question — `Which surface is this project's primary design target? (mobile / tablet / desktop / none) [default: mobile]` — and record the reply under the field `primary_design_surface`; accepting the default records `mobile`, and `none` means the project has no visual surface at all (a CLI, a library, a backend service). It is NOT an ambiguity and does not count against any ambiguity ceiling: on a non-interactive run it is never asked and emits no answer at all, so `write` omits the key entirely.
2. On a recalibration, frame the questions as a review of what would change. Present a delta table (field, old value, new value, changed?) covering the four decision fields — `is_monorepo`, `test_runner`, `naming_style`, `architectural_style` — plus every fired `design_fingerprint.<concern>`, and `code_intelligence` and `primary_design_surface` when that key is present on either side; old fingerprint values come from the `provenance.yaml` sidecar. Ask one question: "accept all changes, or review field by field?" On field-by-field, ask accept / reject / set manually, one field at a time. **Human answers always win** — a field whose existing `source` is `manual` is never silently overwritten, so it always appears in the table and requires explicit acceptance. A field the human sets here carries `source: manual` into write.
3. Collect the answers into `answers{}` (field → value). `nuance.*` answers ride along like any other field; write routes them into `knowledge.yaml`'s `human_nuance` region.
4. Compose the `### CLARIFY ANSWERS` block for write's prompt. It is required even when there was nothing to ask:

   ```
   ### CLARIFY ANSWERS
   answers: { field: value, ... }     # {} when nothing was asked
   skipped: true|false
   blocking_open_items: []
   ```

5. **Non-interactive harness** (no way to reach the human): skip the questions. Apply each skippable ambiguity's `default` — write records those as `origin: default` and never halts on one — and add every non-skippable ambiguity to `blocking_open_items[]`. Set `skipped: true`.
6. **User abort**: if the user cancels, do not launch write. No files are written.

Clarify produces no artifact of its own; the block is its entire output.

## Pending writes — preview, then launch write

Before launching write, show the user what is about to land. This is the last moment anything can pause, and only an inline step can pause:

- `.asdt/config.yaml` — the memory provider, plus `code_intelligence` when it was detected and `primary_design_surface` when it was answered.
- `.asdt/knowledge/knowledge.yaml` — the stack, file structure, design fingerprint, the four decision fields, and any human nuance notes.
- `.asdt/knowledge/provenance.yaml` — the write-only fingerprint provenance sidecar.

List the key deltas alongside them — the fresh values, or the changes the human accepted in the recalibration review — and confirm. On confirmation, launch `steps/write.md` (`agent: builder`) via your delegation primitive with `init/stack-detection` and the `### CLARIFY ANSWERS` block injected. On a decline, nothing is written.

When write returns, tell the user which files landed, surface anything it reported in `settings_preserved[]` or `open_items[]`, and point them at `/asdt-architect`, `/asdt-developer`, and the rest.

## Final Output
`init/write-summary` — the write step's record of which files landed, which `source: manual` settings were carried forward, and which answers were applied. It closes the run; no downstream specialist consumes it.

## Artifact Persistence

Both artifacts are saved via `mem_save`, never as files under `.asdt/artifacts/` or any local path:

- `title`: `"{change-name}/init/{artifact-type}"`
- `topic_key`: `"{project}/{change}/init/{artifact-type}"` — one key per artifact type (`stack-detection`, `write-summary`)
- `type`: `"decision"`
- `content`: structured `What` / `Why` / `Where`

The `.asdt` files themselves are project state, not artifacts — that carve-out belongs to init alone and covers only the three paths listed under Pending writes.

## Invariants
- **Write scope**: `.asdt/config.yaml`, `.asdt/knowledge/knowledge.yaml`, and `.asdt/knowledge/provenance.yaml` only, written only by the `write` sub-agent. The inline steps and explore write nothing.
- **explore NEVER guesses**: no markers matched means an empty stack, never one inferred from incidental files. Every probe is a bounded shell command with no dependency on any particular language or on this repo's own tooling, so the same project yields the same detection in any session.
- Each step reads only its declared inputs; a missing input is recorded in `open_items` and the step proceeds best-effort.
- Artifacts are scoped under the `init/` prefix.
