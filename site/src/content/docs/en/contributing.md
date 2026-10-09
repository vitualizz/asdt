---
title: Contributing
description: How to add a new specialist, improve prompts, write shared skills, and submit a PR to ASDT.
order: 8
locale: en
---

# Contributing

The most impactful contributions are specialist `SKILL.md` files and workflow step definitions — you don't need Go expertise. The skill layer IS the product. If you can describe a specialist's role, its workflow steps, and its artifact contracts, you can ship a new specialist.

## Adding a new specialist

### 1. Create the directory structure

```
skill/asdt-{name}/
  SKILL.md          # specialist definition and workflow
  workflow.yaml     # step sequence and metadata
  steps/            # one .md per subagent step
```

There is no `skills/` directory. Shared criteria live in `skill/asdt-core/references/` and are declared per step via `reference_skills:`.

The directory name **must** start with `asdt-`. The binary embeds the skill tree via `//go:embed SKILL.md asdt-*` in `skill/embedded.go` — any directory matching `asdt-*` ships automatically on the next build.

### 2. Write SKILL.md

```markdown
---
name: asdt-{name}
description: "One sentence: what this specialist produces."
user-invocable: true
specialist-id: {name}
metadata:
  author: "Your Name"
  version: "1.0"
---

# {Name} Specialist

## Role
...

## Orchestration Plan
...

## Final Output
...

## Invariants
...
```

Between the frontmatter and `## Role` go three verbatim blocks the installer depends on — the `FIRST ACTION` blockquote, the empty generated region (`<!-- ASDT:GENERATED:specialist-header -->` … `<!-- /ASDT:GENERATED:specialist-header -->`), and the `ORCHESTRATOR GATE` blockquote. Copy them from any existing specialist without editing. An `inline` step with no prompt file — a gate where the orchestrator pauses for the human, like the Developer's `approve` and `verify` — gets its own `## {step} — …` section before `## Final Output`; that section is the step's whole contract.

`metadata` (`author` + `version`) is required in every `SKILL.md`. `trigger_phrases` is the one optional key — a host-facing discoverability list.

`shared-skills` is **forbidden**. The key is retired: no loader ever resolved it, so declaring it documents nothing and loads nothing. See [How shared skills actually load](#how-shared-skills-actually-load) for the three real mechanisms.

### 3. Write workflow.yaml

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

- A step whose payload only feeds the next step of the same run declares `output: context` instead of `output_topic_key` — it persists nothing.
- An `inline` gate step with no prompt file declares only `name`, `description`, and `execution: inline`.
- **Ceiling: four `subagent` steps.** Inline preludes and gates don't count. Needing a fifth means the design is wrong — merge two, or split the specialist.
- Every cross-specialist input is optional and degrades — no specialist may require another's hand-off.

`host_skills:` names skills the **host assistant** may have installed — a craft ASDT does not ship, such as `frontend-design`. They are names, never paths, and never required: when the host has the skill (or an obvious equivalent), the orchestrator reads the skill's file and injects its text into the sub-agent prompt as a `### HOST SKILL {name}` block — never activating it in its own context; when it doesn't, the skill is skipped silently and nothing is recorded. Declare one only on a `subagent` step whose step file says when it applies and what outranks it — a host skill sharpens execution inside the step's contract, never replaces it. Today UX/UI's `ui-design` and the Developer's `implement` declare `frontend-design`.

### 4. Write step files

Create one `.md` per `subagent` step in `skill/{name}/steps/{step}.md`. Each file contains the LLM instructions for that step — what to read, what to produce, what format the artifact should take.

### 5. Register the specialist

The embed needs nothing from you — `//go:embed SKILL.md asdt-*` already ships your directory. What is **not** automatic is registration: a routable specialist has to be mirrored by hand in two places, plus one test fixture.

1. `skill/SKILL.md` — add the row to the `## Registry` table (command, discipline, when to involve).
2. `internal/installer/assets/agents-template.md` — add the row to the `## ASDT Specialists` table.
3. `skill/embedded_test.go` — the routed-invariant test keeps a hardcoded specialist list that a maintainer must update.

**Do not skip these.** The directory ships either way, so nothing fails at build time — the specialist simply never appears in routing, in the installed agents file, or in the invariant test's coverage.

`/asdt-init` is the exception: it is a setup-class specialist, deliberately not routable and deliberately absent from the routing tables. Do not "fix" that omission.

### 6. Verify with the sandbox

```sh
mkdir -p /tmp/asdt-sandbox
HOME=/tmp/asdt-sandbox go run ./cmd/asdt-tui
```

Installs into a throwaway directory. Confirm your specialist appears as its own top-level sibling under `/tmp/asdt-sandbox/.claude/skills/{name}/`.

### 7. Run the embed tests

```sh
go test ./skill/...
```

`skill/embedded_test.go` verifies every `asdt-*` directory on disk is present in the embedded FS and carries a `SKILL.md`. Fails loudly if your specialist is missing.

## Improving a specialist prompt

1. Edit `skill/{specialist}/SKILL.md` or any file under `skill/{specialist}/steps/`.
2. Run `go test ./skill/...` to confirm the embed registry picks up the changes.
3. Open a PR. Prompt-only PRs are first-class contributions.

## Adding a shared skill

Shared skills are capability fragments reused across multiple specialists — platform context detection, knowledge recall, scope definition.

1. Create `skill/asdt-core/references/{name}.md` with the capability instructions.
2. Wire it through one of the three loading mechanisms below. A shared skill that nothing declares is never read — there is no implicit, ambient loading.
3. Open a PR.

### How shared skills actually load

Three mechanisms, and only three. Paths always resolve from the specialist's own directory.

**1. Install-time splice.** The installer splices `asdt-core/specialist-header.md` into a generated region of every routed `SKILL.md`, so the orchestrator reads the header inline instead of chasing a separate file. This applies to that one file only. The FIRST ACTION blockquote no longer instructs reading `specialist-header.md` — the only file it sends you to is `./workflow.yaml`. Never hand-edit between the region markers; the splice overwrites whatever is there.

**2. Inline step.** A `workflow.yaml` step with `execution: inline` whose `skill:` names a shared file — `knowledge-recall.md`, `platform-context.md` (declared as the `platform-analysis` step). The orchestrator reads that file and follows it in its own context. Nothing is injected anywhere and no sub-agent is launched.

**3. `reference_skills:` on a `subagent` step.** Before launching the step, the orchestrator reads each listed file and injects its content into the sub-agent's prompt as a `### REFERENCE SKILL {path}` block. The sub-agent never fetches them itself — sub-agents run from a different working directory and cannot resolve these paths. When a read fails, the block arrives as `### REFERENCE SKILL {path}: UNRESOLVED` and the step proceeds best-effort.

**Not a shared skill: `host_skills:`.** A skill the host assistant has installed (e.g. `frontend-design`) reaches a sub-agent the same way — the orchestrator injects it as a `### HOST SKILL {name}` block — but it lives outside ASDT, is optional, and its absence changes nothing.

## Code standards

- Early return: `if err != nil { return err }` — validate inputs first.
- No global state — constructor injection throughout.
- Interfaces defined close to consumers, not in the implementing package.
- No `utils/`, `helpers/`, `common/`, or `misc/` packages — domain nouns only.
- Table-driven tests for any logic with more than two cases.

## PR process

- One logical change per PR.
- `go test ./...` must pass.
