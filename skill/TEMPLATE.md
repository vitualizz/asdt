# TEMPLATE — Specialist Authoring Contract

This file is repo-only authoring guidance. It is NOT installed into user projects — `skill/embedded.go` embeds `SKILL.md` and `asdt-*` only, so nothing at this level ships. It is the single normative source for the shape of an ASDT specialist: every routed specialist directory (`skill/asdt-<x>/`) MUST conform to it, whether newly added or being normalized.

Authority: each specialist's workflow.yaml owns step identity, execution mode, agent type, and model; each specialist's ## Orchestration Plan in its SKILL.md owns the tier→step lists; the root skill/SKILL.md `Tailored Workflow Generation` per-specialist table is a derived cache of both and never overrides them.

## Frontmatter

Frontmatter is the FIRST content in `SKILL.md` — nothing may precede the opening `---`.

Required keys:

| Key | Rule |
|---|---|
| `name` | `asdt-{id}`, matching the directory name `asdt-<x>` — that directory name is the identity key for the installed command filename, so the two must agree |
| `description` | A SINGLE LINE, double-quoted — the installer's frontmatter reader only sees top-level single-line scalars, so folded or multi-line values are invisible to it |
| `user-invocable` | `true` |
| `specialist-id` | The bare id (`developer`, `security`, …) |
| `metadata` | Map with `author` and `version` |

Optional: `trigger_phrases` — a host-facing discoverability list.

FORBIDDEN: `shared-skills`. The key is retired — no loader ever resolved it. Do not add it.

## Routed specialist layout

The SKILL.md body follows this fixed order.

(1) The FIRST ACTION blockquote, VERBATIM — copy byte-exact from `skill/asdt-developer/SKILL.md`; em-dashes, backticks, and casing are load-bearing because an automated check greps for the literal `FIRST ACTION — self-load the header`:

```markdown
> **FIRST ACTION — self-load the header**: The specialist header is spliced into this file
> immediately below — read it there. Then read `./workflow.yaml` NOW, before acting on
> anything below. Re-read both whenever you can no longer recall their content (e.g. after
> a context compaction).
```

(2) The specialist-header generated region, VERBATIM — committed with an EMPTY body (begin marker immediately followed by end marker). The installer splices `asdt-shared/skills/specialist-header.md` between the markers at install time, so the orchestrator reads the header inline instead of chasing a separate file. Never hand-edit inside the markers, and never place them at or above the closing frontmatter `---` — a marker inside the frontmatter breaks the frontmatter parser and the OpenCode command wrapper. Dropping the region is silent: with both markers absent the splice is a no-op and the specialist installs with no header, which is why `skill/embedded_test.go` greps for both literals.

```markdown
<!-- GENERATED REGION — do not hand-edit; the shared specialist header is spliced in at install time from asdt-shared/skills/specialist-header.md by registry_gen.go. Edits here are overwritten. -->
<!-- ASDT:GENERATED:specialist-header -->
<!-- /ASDT:GENERATED:specialist-header -->
```

(3) The ORCHESTRATOR GATE blockquote, VERBATIM — same source, same byte-exactness; the check greps for the literal `SOLE orchestrator`:

```markdown
> **ORCHESTRATOR GATE (inline copy — full version in specialist-header.md)**: You, the
> calling assistant, are the SOLE orchestrator of this plan. Launch every `subagent` step
> via your native delegation primitive (Agent/Task) — never run subagent steps inline; run
> `inline` steps in your own context. Sub-agents are bound by the executor header injected
> into their prompts, not by this gate.
```

(4) `# {Name} Specialist` H1.

(5) `## Role` — one sentence for what the specialist does, one for what it does NOT do.

(6) `## Orchestration Plan` — see the contract below.

(7) `## Final Output` — the topic_key of the specialist's final artifact and who consumes it.

(8) `## Artifact Persistence` — all artifacts are saved via `mem_save`, never as files under `.asdt/artifacts/` or any local path. For each artifact:

- `title`: `"{change-name}/{specialist}/{artifact-type}"`
- `topic_key`: `"{project}/{change}/{specialist}/{artifact-type}"` — ONE topic_key per artifact type
- `type`: `"architecture"` for design artifacts, `"decision"` for implementation choices
- `content`: structured `What` / `Why` / `Where` (optionally `Learned`)

(9) `## Invariants` — write-scope rules, artifact-prefix scoping, declared-inputs-only reads, missing-input behavior.

## Orchestration Plan contract

`## Orchestration Plan` contains, IN ORDER:

- the gating paragraph — declares the gating axis (see § Gating axis) and how it filters step depth
- the tier table: `| Level | Steps |`
- a **Trivial eligible** line
- an **Inline steps** line (context injection only — never required as explicit list entries)
- a **Conditional** line (which steps are conditional, which are irrenunciable)
- the Tailored-Workflow-precedence sentence: "When a Tailored Workflow block is present in the prompt, its `steps:` list takes precedence over the complexity-based defaults above."
- the step table: `| Step | File | Execution | Reads | Writes |`
- the per-specialist authority sentence: "This section is the authoritative tier→step mapping for this specialist; workflow.yaml owns step identity, execution, and model; skill/SKILL.md `Tailored Workflow Generation` holds a derived cache row — update it when steps change."

## Gating axis

The gating axis is a template slot. Each specialist declares its axis in the gating paragraph of its `## Orchestration Plan`.

- Default slot value: complexity — `{trivial | simple | moderate | complex}`.
- Security fills the slot with risk_surface — `{none | moderate | high}`. Its Tailored Workflow block carries `risk_surface:` INSTEAD OF `complexity:`, never both.

## Step file layout

Each `steps/{step}.md` follows this section order:

1. `## Purpose` — the canonical purpose prose; the `description:` on the step's workflow.yaml entry is a one-line derived label.
2. `## Inputs` — declared topic_keys, each with `Extract:` rules naming the payload fields to use. Injected-input rule: declared inputs arrive ALREADY INJECTED as `### INPUT {topic_key}` blocks per `parallel-retrieval.md` — a step NEVER self-fetches its own declared inputs.
3. `## Context budget` — mandatory on every step: an explicit token cap.
4. `## Output budget` (optional) — caps the emitted payload. `Exceeding the budget is a defect: trim, do not spill.`
5. `## Processing` — numbered steps.
6. `## Output` — a `Produces:` line naming the topic_key, then the canonical closing sentence, then a fenced `payload:` schema whose last key is `open_items: []`.

The canonical closing sentence, used verbatim in every step's `## Output`:

```markdown
Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.
```

It is deliberately self-contained: the artifact contract lives in this repo-only file, so a shipped step file must never lean on a term it cannot define for its reader.

## Artifact contract

The canonical rules, stated once:

- `payload:` is the ROOT of the persisted artifact — no wrapper fields around it.
- `open_items: []` is always present as the trailing key — `[]` when nothing degraded. Entry format: `"{artifact or resource} absent — {what was assumed}"`.
- The final generative step additionally carries `summary: ""` (≤ 150 tokens), read by the decision-preservation shared skill.
- When a declared input arrives as `### INPUT {topic_key}: UNRESOLVED`, proceed best-effort and record the gap in `open_items` — never fail the step.
- Artifacts are addressed exclusively by Engram topic_key `{project}/{change}/{specialist}/{artifact-type}`, one key per artifact type, so a declared input resolves to exactly one artifact through a single `mem_search` / `mem_get_observation` pair. Artifacts are never written to local files.

Retired fields — a specialist MUST NOT emit `schema_version`, `agent`, `change_id`, `created_at`, `prompt_version`, or `input_refs`. Nothing reads them.

## Optional inputs and degradation

An input that only some tiers produce is marked with an end-of-line `# optional` comment on its entry in workflow.yaml AND paired with this paragraph in the consuming step file:

```markdown
**DEGRADATION — `{topic_key}` is optional ({tier condition})**: when it arrives as `### INPUT {topic_key}: UNRESOLVED`, {one-sentence fallback behavior}; append "{artifact} absent — {assumption}" to open_items. Never block on this input.
```

Both halves are required — a `# optional` marker without the paragraph, or the reverse, is an authoring defect. Tier-conditioning NEVER adds or removes steps; it only makes inputs optional.

## Dual-artifact steps

Some steps produce two final artifacts. Their `## Output` uses:

```markdown
This step produces TWO final artifacts. Persist `{primary}` under this step's output_topic_key ({project}/{change}/{specialist}/{primary}); persist `{secondary}` under its own distinct per-type topic_key {project}/{change}/{specialist}/{secondary}, noted on this step's workflow.yaml entry. Return one payload covering both persisted keys.
```

Current cases: architect `technical-handoff`, security `hardening-checklist`, ux-ui `ux-handoff`.

## Setup-class variant

`asdt-init` is the one setup-class specialist. Deltas from the routed shape:

- Deliberately NOT routable and deliberately absent from the routing tables — do not "fix" that omission.
- No tier table; a fixed five-step flow: `knowledge-gate → explore → enrichment → clarify → write`.
- The gates stay inline with the orchestrator because they need its tool list and can pause for the human.
- The `write` sub-agent owns every file write, and a post-write self-check is required.
- It follows the universal core of this template and is exempt only from the routed-mandatory blocks (the FIRST ACTION / ORCHESTRATOR GATE literals, the specialist-header generated region, and the routable registration mirrors). Its `SKILL.md` matches the installer's per-specialist predicate but carries no markers, so the header splice is an intended no-op there.

## Cross-references

Reference a section by its backticked NAME, never by its number. A stale number keeps resolving — silently, to the wrong section — while a stale name fails visibly and gets fixed.

- Same file: the name alone — `Gating axis`, `Step file layout`.
- Across files: the file path AND the name — `skill/SKILL.md` `Tailored Workflow Generation`. Never the path alone, never the name alone.
- Use the FULL section name. `Tailored Workflow Generation` is not `Tailored Workflow`, which collides with the literal `## Tailored Workflow` block headings.
- The numbered headings in `skill/SKILL.md` STAY numbered: the numbers are navigation aids, and a drift test matches `## 5. Specialist Registry` as a literal string. Names replace pointers, never targets — renumber nothing.

## Registration ritual

Adding a routable specialist means three manual mirrors:

1. `skill/SKILL.md` — the `Specialist Registry` row.
2. `skill/SKILL.md` — the `Tailored Workflow Generation` per-specialist table row. Never hand-edit inside the generated inline-steps markers; that sub-region is regenerated at install.
3. `internal/installer/assets/agents-template.md` — the ASDT Specialists row.

workflow.yaml `name:` fields stay authoritative for step identity. The repo's routed-invariant test (`skill/embedded_test.go`) also keeps a hardcoded specialist list that a maintainer must update.

## Worked example

Frontmatter for a routed specialist:

```yaml
---
name: asdt-researcher
description: "Turns a vague idea into a validated problem statement — discovery, ideation, and a feasibility brief before any spec is written."
user-invocable: true
specialist-id: researcher
trigger_phrases:
  - explore this idea
  - is this feasible
metadata:
  author: "Lee Palacios (vitualizz)"
  version: "1.0"
---
```

Minimal Orchestration Plan (gating paragraph and lines abbreviated):

```markdown
## Orchestration Plan

**Complexity-based step filtering**: Researcher gates step depth by complexity.

| Level | Steps |
|-------|-------|
| **trivial** | `discover` |
| **simple** | `discover → brief` |

**Trivial eligible**: Yes — `discover` has `inputs: []`.
**Inline steps** (context injection only): `knowledge-recall`, `decision-preservation`
**Conditional**: none — `discover` is irrenunciable.

When a Tailored Workflow block is present in the prompt, its `steps:` list takes precedence over the complexity-based defaults above.

| Step | File | Execution | Reads | Writes |
|------|------|-----------|-------|--------|
| discover | steps/discover.md | subagent | *(request + platform-summary)* | `researcher/discovery` |
| brief | steps/brief.md | subagent | `researcher/discovery` | `researcher/feasibility-brief` |

This section is the authoritative tier→step mapping for this specialist; workflow.yaml owns step identity, execution, and model; skill/SKILL.md `Tailored Workflow Generation` holds a derived cache row — update it when steps change.
```

Step-file skeleton (`steps/brief.md`), following § Step file layout:

````markdown
## Purpose
Consolidate discovery findings into the final feasibility brief.

## Inputs
- `researcher/discovery` — arrives as an `### INPUT {topic_key}` block. Extract: `findings[]`, `open_questions[]`.

## Context budget
Max 1,000 tokens of input material.

## Processing
1. Rank findings by product impact.
2. State a go / no-go recommendation with rationale.

## Output
Produces: `researcher/feasibility-brief`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  recommendation: ""
  findings: []
  summary: ""
  open_items: []
```
````
