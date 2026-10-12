# TEMPLATE — Specialist Authoring Contract

How to add or normalize a specialist. Repo-only guidance — `skill/embedded.go` bundles `SKILL.md` and `asdt-*` only, so nothing at this level ships to a user project. The shape is deliberately small: one specialist, a few steps, ONE role summary. If a design needs a consolidation step, a self-review step, or a second summary, the design is wrong — go back and cut it.

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
- `## Final Output` — what this specialist's role summary carries, which step writes it, and who consumes it.
- `## Invariants` — a handful of lines: write scope, prefix, degradation, and whatever is load-bearing for THIS specialist.

## 3. workflow.yaml

```yaml
specialist: {name}               # the role id: names the summary file and own-step inputs
routable: true
steps:
  - name: knowledge-recall       # inline prelude: names the reference it follows
    skill: ../asdt-core/references/knowledge-recall.md
    description: One line.
    execution: inline

  - name: {step}
    skill: steps/{step}.md
    description: One line — what this step returns.
    execution: subagent          # or: inline
    model: haiku | sonnet | opus
    agent: analyst | builder     # builder ONLY when the step writes host files
    inputs:                      # optional — subagent steps only
      - pm # optional — what to extract from Product's summary
      - {name}/{earlier-step} # this specialist's own step payload
    output: summary              # optional — this step's return becomes the role summary
    reference_skills:
      - ../asdt-core/references/{x}.md
      - ../asdt-core/protocol.md
    host_skills:                 # optional — subagent steps only
      - {kebab-case-skill-name}
```

**`inputs:`** — scratchpad names, never paths or memory keys. The orchestrator resolves each one to a file (`asdt-core/protocol.md` §1) and injects it as `### INPUT {name}: {path}`, or `UNRESOLVED` (§5). An entry is one of:

- `{role}` — another specialist's role summary for this change, `{scratchpad}/asdt/{change}/{role}.md`, where `{role}` is that specialist's `specialist:` id. Every cross-specialist input is optional: mark it with a trailing `# optional — {what to extract}` comment, and have the step file degrade when it is missing. `{role}` may also be this specialist's own id, meaning its current summary (a plan being revised).
- `{name}/{step}` — the return of an earlier step of THIS specialist, `{scratchpad}/asdt/{change}/{name}/{step}.md`. Listing it is what makes the orchestrator write that file. The step may be an inline gate that produces a payload (the Developer's `verify` failures). Mark it `# optional` when its producer may not have run. Never name another specialist's step: read its role summary.
- Never `feedback` or `previous`. Those are reserved for the orchestrator, which injects them whenever the user's correction makes it re-launch a step (§5).

**`output: summary`** — the step's return is written as the role summary, `{scratchpad}/asdt/{change}/{name}.md`. It is the only value. A later `output: summary` step of the same run replaces the file, which is how a staged role (plan, then build) keeps ONE summary with a `Stage`. An inline gate declares it when the orchestrator itself rewrites the summary (the Developer's `verify`). A step without it hands its return only to the steps that list it, and to the orchestrator's narration.

`host_skills:` names skills the host assistant MAY have installed — a craft ASDT does not ship, such as `frontend-design`. They are never paths and never required: the orchestrator injects each one it finds and skips the rest silently (`asdt-core/protocol.md` §5). Declare one only on a `subagent` step whose step file says when it applies and what outranks it — a host skill sharpens execution inside the step's contract, never replaces it.

An `inline` step with no prompt file — a gate where the orchestrator pauses for the human — declares only `name`, `description`, and `execution: inline`, plus `output: summary` when it rewrites the role summary. No inline step declares `inputs:`. Its contract, including any payload it produces for a later step, is a section of SKILL.md (§2).

**Ceiling: four sub-agent steps.** A specialist declares AT MOST four `execution: subagent` steps in its `workflow.yaml`. Inline steps — the orchestrator's own preludes and gates — do not count: they are no sub-agent of their own, and a gate re-launches only steps already declared. Needing a fifth sub-agent step means the design is wrong — merge two, or split the specialist.

**Every specialist works standalone.** No specialist may require another's summary: every cross-specialist input is optional and degrades. Whether a run delivers a change or studies what already exists is judged from the invocation per `asdt-core/protocol.md` §0, never from a flag or a mode.

**Durable project knowledge lives in memory**, remembered by the orchestrator alone (`asdt-core/protocol.md` §2). No step reaches memory, and no specialist invents a second store for conventions.

**The study step is called `review`** and declares `output: summary`; a study run writes under `study-{topic}` instead of `{change}`. A specialist whose normal work ALREADY examines what exists — an audit, a discovery — does not duplicate its chain into a review; it runs the same steps under the study name.

Paths in `skill:` and `reference_skills:` resolve from the specialist's own directory.

## 4. Step file

Four sections. Nothing else.

```markdown
# {Step} — {Specialist} Specialist
## Purpose        One or two sentences.
## Inputs         One bullet per `### INPUT {name}` line, marking the optional ones
                  OPTIONAL, and what to extract from that file. Then: each arrives as a
                  path to read — never look for anything else. Then one DEGRADATION line
                  per optional input.
## Processing     Numbered steps. Reasoning and judgment, not bureaucracy.
## Output         "Return this summary:" — the protocol §6 schema, naming the sections
                  this step fills and any role-specific section it adds.
```

A step never reaches memory, never writes a file (the two exceptions are `asdt-core/protocol.md` §4), and never asks the user — it hardens every gap into an `ASSUMED:` entry. Never write a provider's name or tool names — memory or code intelligence — into any prompt: the installer binds the verbs, and `TestPromptsAreProviderNeutral` (`internal/installer/memory_binding_test.go`) fails on a hardcoded one. A step never reads the memory binding: the orchestrator strips it from every text it injects (protocol §5), so a value a step needs from it arrives in the step's payload. The code-intelligence binding is the step's own: it rides in the executor header, so a step speaks the code-intelligence verbs (**explore**, **search**, …) and never a tool name.

No numeric context budgets. No per-input degradation paragraph — one line each. No envelope: the return value IS the summary. A step file NEVER contains the EXECUTOR block; those guardrails come from the agent definition (`agent: analyst` / `agent: builder` bake in `asdt-core/executor-header.md`) or from the orchestrator prepending that header.

## 5. Summary contract

One role summary per role per change: `{scratchpad}/asdt/{change}/{role}.md`. Its schema is defined once, in `asdt-core/protocol.md` §6 — reference it, never restate it. A role omits sections that do not apply, and never adds sections about the process itself.

## 6. Registering a new specialist

1. Create the directory per §1.
2. Add one row to the `## Registry` table in `skill/SKILL.md`, and place it in the dependency prose under that table.
3. The Go side keeps its own mirrors: `registryRenderOrder` and `inlineStepsDisplayNames` in `internal/installer/registry_gen.go` enumerate the routed directories, and `TestRegistryDrift` asserts those mirrors agree with the `workflow.yaml` files. Update them in the same commit, or that test fails.
