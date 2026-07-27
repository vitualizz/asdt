# ASDT Shared Skills

Cross-specialist utility files. They reach specialist contexts via three mechanisms:

- **Install-time splicing** — the installer writes the body of `specialist-header.md` directly into each routed specialist `SKILL.md`, so an installed specialist already carries it. The **FIRST ACTION Read** at the top of every specialist `SKILL.md` body is the fallback that loads it (and its `workflow.yaml`) when the body was not spliced
- `reference_skills:` in a `subagent` step entry in `workflow.yaml` — the orchestrator reads each listed file and injects it into that step's sub-agent prompt only. The sub-agent never fetches it itself
- `skill:` on an `inline` step entry in `workflow.yaml` — the orchestrator reads that file and follows it in its own context; nothing is injected anywhere

These are **reference text injected into the active context**, not independently executable units. They have no `## Inputs` / `## Output` structure of their own.

The artifact contract itself — the `payload:` root, the trailing `open_items: []`, the retired wrapper fields, and the per-artifact-type topic_key rule — lives in `../../TEMPLATE.md`, the repo-only authoring contract for specialists. It is not a shared skill and is never injected into a step.

## Index

### Structural — required in every specialist

| File | Purpose |
|---|---|
| `specialist-header.md` | Loaded via the FIRST ACTION Read in each specialist `SKILL.md` body. Contains the ORCHESTRATOR GATE and prerequisite logic. |
| `executor-header.md` | Injected into every `subagent` step prompt. Instructs the executor: do the single assigned step, do NOT orchestrate or delegate. |

### Artifact contracts

| File | Purpose |
|---|---|
| `artifact-loading.md` | How a specialist's first artifact-consuming step (declared `inputs: []`) retrieves upstream artifacts from Engram (`mem_search` → `mem_get_observation`), extracts relevant fields, and records missing artifacts in `open_items[]`. |
| `parallel-retrieval.md` | The canonical orchestrator fetch-once cache pattern — prevents duplicate Engram lookups when multiple steps need the same artifact. |

### Workflow utilities

| File | Purpose |
|---|---|
| `knowledge-recall.md` | Queries organizational memory for prior decisions relevant to the current change. Used as the first inline step in most specialists. |
| `platform-context.md` | Injects the project's detected platform knowledge (stack, conventions, design fingerprint) into a specialist's context. Includes the project-level reuse guard (reads `.asdt/knowledge/knowledge.yaml`); also backs the inline `platform-analysis` workflow step. |
| `decision-preservation.md` | Saves a permanent organizational knowledge record after a significant decision is produced. Used as the final inline step in most specialists. |
| `scope-definition.md` | Guidelines for defining explicit, unambiguous project scope. Used by Architect and Developer. |
| `report.md` | Generates a structured handoff document from multiple intermediate artifacts; includes the extraction rules (200-token budget per artifact). Used as the consolidation step in UX/UI, Architect, QA, Security, and PM. |
| `nfr-budget.md` | Shared NFR-target value-object shape + the measure-against-budget protocol; reference-only, produces no artifact. Referenced by PM `success-metrics`, Architect `cost-estimation`, QA `performance-validation`. |

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

There is no frontmatter key that loads shared skills. `specialist-header.md` reaches a specialist through install-time splicing, with the FIRST ACTION Read as fallback; `reference_skills:` is the only declarative way to pull one of these files into a `subagent` step, resolved and injected by the orchestrator; `skill:` on an `inline` step is read by the orchestrator directly.
