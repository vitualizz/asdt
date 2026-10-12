# TEMPLATE — Specialist Authoring Contract

How to add or normalize a specialist. Repo-only guidance — `skill/embedded.go` bundles `SKILL.md` and `asdt-*` only, so nothing at this level ships to a user project. The shape is deliberately small: one specialist, a few steps, ONE persisted hand-off. If a design needs a consolidation step, a self-review step, or a second output key, the design is wrong — go back and cut it.

## 1. Directory layout

```
skill/asdt-{name}/
├── SKILL.md          ← the orchestration plan (what runs, in what order)
├── workflow.yaml     ← the machine-readable registry (the Go binary reads this)
└── steps/
    └── {step}.md     ← one prompt file per subagent step
```

No `skills/` directory. Shared criteria live in `asdt-core/references/` and are declared per step via `reference_skills:`.

## 2. SKILL.md

**Frontmatter** — first content in the file, nothing before the opening `---`. `name` must match the directory name: it is the identity key for the installed command filename.

```yaml
---
name: asdt-{name}
description: "One sentence: what it does and when to bring it in."
user-invocable: true
specialist-id: {name}
trigger_phrases:
  - ...
metadata:
  author: "..."
  version: "1.0"
---
```

**Three verbatim blocks follow, in this order.** The Go installer depends on all three; copy them from any existing specialist without editing:

1. The `FIRST ACTION — self-load the header` blockquote.
2. The generated region — the HTML comment, then `<!-- ASDT:GENERATED:specialist-header -->` immediately followed by `<!-- /ASDT:GENERATED:specialist-header -->`, committed with an EMPTY body. The installer splices `asdt-core/specialist-header.md` between the markers. Never hand-edit inside them, and never place them at or above the closing frontmatter `---` (a marker inside the frontmatter breaks the parser and the OpenCode wrapper). Dropping the region is SILENT — with both markers absent the splice is a no-op and the specialist installs with no header, which is why `skill/embedded_test.go` greps for both literals.
3. The `ORCHESTRATOR GATE` blockquote.

**Body**, in this order:

- `## Role` — two or three sentences: what it does, and what it explicitly does not.
- `## Orchestration Plan` — for a single-step specialist, one sentence naming the step and what depth changes; for a multi-step one, a table mapping what the request asks for to the chain that serves it. Then the pointer: `Step identity, model, inputs, and outputs: workflow.yaml.` Never restate a step's file, model, inputs, or outputs here — that duplication is what `workflow.yaml` exists to prevent.
- One `## {step} — …` section per `inline` step that has no prompt file, and only for those — that section IS the step's contract, so it says when the step runs, what it shows or asks, and what happens on every answer, including when no human can answer.
- `## Final Output` — the one key this specialist persists, and who consumes it.
- `## Invariants` — a handful of lines: write scope, prefix, degradation, and whatever is load-bearing for THIS specialist.

## 3. workflow.yaml

```yaml
specialist: {name}
routable: true
steps:
  - name: {step}
    skill: steps/{step}.md
    description: One line — what this step produces.
    execution: subagent          # or: inline
    model: haiku | sonnet | opus
    agent: analyst | builder     # builder ONLY when the step writes host files
    inputs:
      - "{project}/{change}/{role}/handoff"  # optional — say so, and degrade
    output_topic_key: "{project}/{change}/{name}/handoff"
    reference_skills:
      - ../asdt-core/references/{x}.md
      - ../asdt-core/protocol.md
    host_skills:                 # optional — subagent steps only
      - {kebab-case-skill-name}
```

`host_skills:` names skills the host assistant MAY have installed — a craft ASDT does not ship, such as `frontend-design`. They are never paths and never required: the orchestrator injects each one it finds and skips the rest silently (`asdt-core/protocol.md` §4). Declare one only on a `subagent` step whose step file says when it applies and what outranks it — a host skill sharpens execution inside the step's contract, never replaces it.

A step whose payload feeds the NEXT step of the same run and matters to no one after it declares `output: context` instead of `output_topic_key`: it persists nothing, and the orchestrator retains the payload and injects it into the next step. A step whose payload must ALSO outlive the run — a plan a later session resumes — keeps `output_topic_key` on the role's one key and is still injected from context within the run. Both rules, and the staged record they produce, live in `asdt-core/protocol.md` §1; a second key per stage is never the answer.

An `inline` step with no prompt file — a gate where the orchestrator pauses for the human — declares only `name`, `description`, and `execution: inline`, plus `output_topic_key` when the orchestrator itself persists the step's outcome (a gate that records a result on the role's one key). Its contract is a section of SKILL.md (§2).

**Ceiling: four sub-agent steps.** A specialist declares AT MOST four `execution: subagent` steps in its `workflow.yaml`. Inline steps — the orchestrator's own preludes and gates — do not count: they are no sub-agent of their own, and a gate re-launches only steps already declared. Needing a fifth sub-agent step means the design is wrong — merge two, or split the specialist.

**Every specialist works standalone.** No specialist may require another's hand-off: every cross-specialist input is optional and degrades. Whether a run delivers a change or studies what already exists is judged from the invocation per `asdt-core/protocol.md` §1, never from a flag or a mode.

**Permanent project knowledge lives in `human_nuance`** inside `.asdt/knowledge/knowledge.yaml`, and only ever enters with the user's consent. No specialist writes it through a sub-agent, and none invents a second store for conventions.

**The study step is called `review`** and persists under `{project}/study/{topic}/{role}`. A specialist whose normal work ALREADY examines what exists — an audit, a discovery — does not duplicate its chain into a review; it runs the same steps under the study key.

Paths in `skill:` and `reference_skills:` resolve from the specialist's own directory.

## 4. Step file

Four sections. Nothing else.

```markdown
# {Step} — {Specialist} Specialist
## Purpose        One or two sentences.
## Inputs         One bullet per input, marking the optional ones OPTIONAL, and what to
                  extract from each. Then: everything arrives ALREADY INJECTED — never
                  self-fetch. Then one DEGRADATION line per optional input.
## Processing     Numbered steps. Reasoning and judgment, not bureaucracy.
## Output         Produces `{role}/handoff`, with the canonical schema from protocol.md.
```

Reach memory only through the verbs of `asdt-core/protocol.md` §0 — "persist via memory **save** under this step's `output_topic_key`" — and never write a provider's name or tool names into a prompt: the installer binds the verbs to the selected provider, and `TestPromptsAreProviderNeutral` (`internal/installer/memory_binding_test.go`) fails on a hardcoded one.

No numeric context budgets. No per-input degradation paragraph — one line each. No summary envelope: the return value IS the payload. A step file NEVER contains the EXECUTOR block; those guardrails come from the agent definition (`agent: analyst` / `agent: builder` bake in `asdt-core/executor-header.md`) or from the orchestrator prepending that header.

## 5. Artifact contract

One key per role per change: `{project}/{change}/{role}/handoff`. The payload schema is defined once, in `asdt-core/protocol.md` §5 — reference it, never restate it. A role omits keys that do not apply, and never adds keys about the process itself.

## 6. Registering a new specialist

1. Create the directory per §1.
2. Add one row to the `## Registry` table in `skill/SKILL.md`, and place it in the dependency prose under that table.
3. The Go side keeps its own mirrors: `registryRenderOrder` and `inlineStepsDisplayNames` in `internal/installer/registry_gen.go` enumerate the routed directories, and `TestRegistryDrift` asserts those mirrors agree with the `workflow.yaml` files. Update them in the same commit, or that test fails.
