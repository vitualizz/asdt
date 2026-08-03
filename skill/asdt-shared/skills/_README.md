# ASDT Shared Skills

> **Repo-only — this file does not ship.** It is maintainer documentation with no runtime reader. `embedded.go` embeds `asdt-*` as a directory match, and Go's `embed` package recursively excludes every name beginning with `_` or `.`. The leading underscore is what keeps this doc out of the binary and out of every installed project, while keeping it in the same directory as the fragments it documents.

Cross-specialist utility files. They reach specialist contexts via three mechanisms:

- **Install-time splicing** — the installer writes the bodies of `specialist-header.md`, `parallel-retrieval.md`, `intake-contract.md`, and `knowledge-recall.md`, in that order, directly into each routed specialist `SKILL.md`, so an installed specialist already carries all four as one continuous document. Every fragment stays its own file on disk; the join happens only in the generator. The **FIRST ACTION Read** at the top of every specialist `SKILL.md` body is the fallback that loads the header (and its `workflow.yaml`) when the body was not spliced
- `reference_skills:` in a `subagent` step entry in `workflow.yaml` — the orchestrator reads each listed file and injects it into that step's sub-agent prompt only. The sub-agent never fetches it itself
- `skill:` on an `inline` step entry in `workflow.yaml` — the orchestrator reads that file and follows it in its own context; nothing is injected anywhere

These are **reference text injected into the active context**, not independently executable units. They have no `## Inputs` / `## Output` structure of their own.

The artifact contract itself — the `payload:` root, the trailing `open_items: []`, the retired wrapper fields, and the per-artifact-type topic_key rule — lives in `../../TEMPLATE.md`, the repo-only authoring contract for specialists. It is not a shared skill and is never injected into a step.

## Index

`Routing` says how the fragment reaches a specialist: **folded** means the installer splices its body into every routed `SKILL.md`, so it is already in context before the first step runs; **runtime** means the orchestrator reads it on demand, via `reference_skills:` on a `subagent` step or `skill:` on an `inline` step.

### Structural — required in every specialist

| File | Routing | Purpose |
|---|---|---|
| `specialist-header.md` | **Folded** (1st) | The spliced header itself: ORCHESTRATOR GATE, prerequisite logic, launch and skill-resolution rules. The FIRST ACTION Read in each specialist `SKILL.md` body is the fallback when the body was not spliced. |
| `executor-header.md` | Runtime | Injected into every `subagent` step prompt — or already baked into the `asdt-analyst` / `asdt-builder` agent definitions. Instructs the executor: do the single assigned step, do NOT orchestrate or delegate. |

### Artifact contracts

| File | Routing | Purpose |
|---|---|---|
| `parallel-retrieval.md` | **Folded** (2nd) | The canonical orchestrator fetch-once cache pattern — Cache Ledger Rule, Injection Format, executor-header injection, UNRESOLVED degradation. Folded because the header depends on it from its very first launch instruction. |
| `intake-contract.md` | **Folded** (3rd) | How every step treats its inputs: the declared-vs-present input check, the single batched clarification turn (fully suppressed under a `## Tailored Workflow` block), and the harden-always `ASSUMED:` degradation into `open_items`. Folded because these rules must be settled before any step content arrives. |
| `artifact-loading.md` | Runtime | How a specialist's first artifact-consuming step (declared `inputs: []`) retrieves upstream artifacts from Engram (`mem_search` → `mem_get_observation`), extracts relevant fields, and records missing artifacts in `open_items[]`. |

### Workflow utilities

| File | Routing | Purpose |
|---|---|---|
| `knowledge-recall.md` | **Folded** (4th) | Queries organizational memory for prior decisions relevant to the current change. It is the first inline step of nearly every specialist and runs at every tier, so folding it costs nothing and saves a read. |
| `platform-context.md` | Runtime — deliberately not folded | Injects the project's detected platform knowledge (stack, conventions, design fingerprint) into a specialist's context. Includes the project-level reuse guard (reads `.asdt/knowledge/knowledge.yaml`); also backs the inline `platform-analysis` workflow step. |
| `decision-preservation.md` | Runtime — deliberately not folded | Saves a permanent organizational knowledge record after a significant decision is produced. Used as the final inline step in most specialists. |
| `scope-definition.md` | Runtime | Guidelines for defining explicit, unambiguous project scope. Used by Architect and Developer. |
| `report.md` | Runtime | Generates a structured handoff document from multiple intermediate artifacts; includes the extraction rules (200-token budget per artifact). Used as the consolidation step in UX/UI, Architect, QA, Security, and PM. |
| `nfr-budget.md` | Runtime | Shared NFR-target value-object shape + the measure-against-budget protocol; reference-only, produces no artifact. Referenced by PM `backlog`, Architect `cost-estimation`, QA `performance-validation`. |

### Why two fragments stay runtime on purpose

These two decisions are written down so they stop being re-litigated:

- **`decision-preservation.md` is tier-gated.** It runs last and only at complexity `moderate` and above. Folding it would charge every `trivial` and `simple` run the full cost of a step those runs never reach — the fragment would sit in context from the first token and be paid for even when it is never used.
- **`platform-context.md` is project-specific.** Its content is re-derived per install from the detected stack, so a copy baked into the shared header would be stale by construction — wrong for every project except the one it was generated against.

## How to Reference

In `workflow.yaml` (step-specific):

```yaml
- name: system-design
  execution: subagent
  reference_skills:
    - ../asdt-shared/skills/platform-context.md
    - ../asdt-shared/skills/scope-definition.md
```

On an `inline` step, `skill:` names the single file the orchestrator reads and follows in its own context:

```yaml
- name: knowledge-recall
  execution: inline
  skill: ../asdt-shared/skills/knowledge-recall.md
```

There is no frontmatter key that loads shared skills. The four **folded** fragments reach a specialist through install-time splicing, with the FIRST ACTION Read as fallback; `reference_skills:` is the only declarative way to pull one of these files into a `subagent` step, resolved and injected by the orchestrator; `skill:` on an `inline` step is read by the orchestrator directly. A `skill:` path naming a folded fragment stays declared — `workflow.yaml` remains the machine-readable truth — but the orchestrator treats it as already resolved instead of reading it again.
