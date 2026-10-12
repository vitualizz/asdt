---
title: Contributing
description: How to add a new specialist, improve prompts, write shared skills, and submit a PR to ASDT.
order: 8
locale: en
---

# Contributing

The most impactful contributions are specialist `SKILL.md` files and workflow step definitions — you don't need Go expertise. The skill layer IS the product. If you can describe a specialist's role, its workflow steps, and the summary it hands on, you can ship a new specialist.

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

- **`inputs:`** are scratchpad names, never paths or memory keys. The orchestrator resolves each one to a file in the session scratchpad and injects it as `### INPUT {name}: {path}`, or `UNRESOLVED` when the file is missing. An entry is either `{role}` (another specialist's summary for this change, `{scratchpad}/asdt/{change}/{role}.md`, where `{role}` is that specialist's `specialist:` id) or `{name}/{step}` (the return of an earlier step of THIS specialist, `{scratchpad}/asdt/{change}/{name}/{step}.md`; listing it is what makes the orchestrator write that file). Never name another specialist's step: read its summary. `feedback` and `previous` are reserved for the orchestrator, which injects them when your correction makes it re-launch a step.
- **`output: summary`** means the step's return is written as the role summary, `{scratchpad}/asdt/{change}/{name}.md`. A later `output: summary` step of the same run replaces it, which is how a staged role (plan, then build) keeps ONE summary with a `Stage`. A step without it hands its return only to the steps that list it.
- An `inline` gate step with no prompt file declares only `name`, `description`, and `execution: inline`, plus `output: summary` when the orchestrator itself rewrites the summary (the Developer's `verify`). No inline step declares `inputs:`.
- **Ceiling: four `subagent` steps.** Inline preludes and gates don't count. Needing a fifth means the design is wrong — merge two, or split the specialist.
- Every cross-specialist input is optional and degrades — mark it with a trailing `# optional — {what to extract}` comment. No specialist may require another's summary.
- No step reaches memory. Durable project knowledge is remembered by the orchestrator alone, and no specialist invents a second store for conventions.

`host_skills:` names skills the **host assistant** may have installed — a craft ASDT does not ship, such as `frontend-design`. They are names, never paths, and never required: when the host has the skill (or an obvious equivalent), the orchestrator reads the skill's file and injects its text into the sub-agent prompt as a `### HOST SKILL {name}` block — never activating it in its own context; when it doesn't, the skill is skipped silently and nothing is recorded. Declare one only on a `subagent` step whose step file says when it applies and what outranks it — a host skill sharpens execution inside the step's contract, never replaces it. Today UX/UI's `ui-design` and the Developer's `implement` declare `frontend-design`.

### 4. Write step files

Create one `.md` per `subagent` step in `skill/{name}/steps/{step}.md`. Each file has four sections — `Purpose`, `Inputs` (one bullet per `### INPUT {name}` line, each a path to read, with one degradation line per optional input), `Processing`, and `Output` (the summary it returns, in the shared protocol's schema). A step never reaches memory, never writes a file (the Developer's `implement` and `/asdt-init`'s `write` are the only exceptions), and never asks the user: every gap becomes an `ASSUMED:` entry. Never write a memory provider's name or tool names into a prompt — `TestPromptsAreProviderNeutral` fails on one.

### 5. Register the specialist

The embed needs nothing from you — `//go:embed SKILL.md asdt-*` already ships your directory. What is **not** automatic is registration: a routable specialist has to be mirrored by hand in a few places.

1. `skill/SKILL.md` — add the row to the `## Registry` table (command, discipline, when to involve), and place it in the dependency prose under that table.
2. `internal/installer/assets/agents-template.md` — add the row to the `## ASDT Specialists` table.
3. `internal/installer/registry_gen.go` — add the directory to `registryRenderOrder` and `inlineStepsDisplayNames`. `TestRegistryDrift` asserts these mirrors agree with the `workflow.yaml` files.
4. `skill/embedded_test.go` — the routed-invariant test keeps a hardcoded specialist list that a maintainer must update.

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

**2. Inline step.** A `workflow.yaml` step with `execution: inline` whose `skill:` names a shared file — `knowledge-recall.md`, `platform-context.md` (declared as the `platform-analysis` step). The orchestrator reads that file and follows it in its own context; no sub-agent is launched. What the prelude produces — the platform block, the recalled project notes — reaches each later step as lines in its prompt.

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
