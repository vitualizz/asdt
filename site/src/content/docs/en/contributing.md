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
  steps/            # one .md per workflow step
  skills/           # specialist-scoped skill fragments (optional)
```

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

## Invariants
...
```

`metadata` (`author` + `version`) is required in every `SKILL.md`. `trigger_phrases` is the one optional key — a host-facing discoverability list.

`shared-skills` is **forbidden**. The key is retired: no loader ever resolved it, so declaring it documents nothing and loads nothing. See [How shared skills actually load](#how-shared-skills-actually-load) for the three real mechanisms.

### 3. Write workflow.yaml

```yaml
specialist: {name}
steps:
  - id: step-one
    name: Step One
  - id: step-two
    name: Step Two
```

### 4. Write step files

Create one `.md` per step in `skill/{name}/steps/{step-id}.md`. Each file contains the LLM instructions for that step — what to read, what to produce, what format the artifact should take.

### 5. Register the specialist

The embed needs nothing from you — `//go:embed SKILL.md asdt-*` already ships your directory. What is **not** automatic is registration: a routable specialist has to be mirrored by hand in three places, plus one test fixture.

1. `skill/SKILL.md` — add the `Specialist Registry` row (command, discipline, when to involve).
2. `skill/SKILL.md` — add the `Tailored Workflow Generation` per-specialist table row. Never hand-edit inside the generated inline-steps markers in that table; that sub-region is regenerated at install time and your edits are overwritten.
3. `internal/installer/assets/agents-template.md` — add the ASDT Specialists row.
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

1. Edit `skill/{specialist}/SKILL.md` or any file under `skill/{specialist}/steps/` or `skill/{specialist}/skills/`.
2. Run `go test ./skill/...` to confirm the embed registry picks up the changes.
3. Open a PR. Prompt-only PRs are first-class contributions.

## Adding a shared skill

Shared skills are capability fragments reused across multiple specialists — platform context detection, knowledge recall, scope definition.

1. Create `skill/asdt-shared/skills/{name}.md` with the capability instructions.
2. Wire it through one of the three loading mechanisms below. A shared skill that nothing declares is never read — there is no implicit, ambient loading.
3. Open a PR.

### How shared skills actually load

Three mechanisms, and only three. Paths always resolve from the specialist's own directory.

**1. Install-time splice.** The installer splices `asdt-shared/skills/specialist-header.md` into a generated region of every routed `SKILL.md`, so the orchestrator reads the header inline instead of chasing a separate file. This applies to that one file only. The FIRST ACTION blockquote no longer instructs reading `specialist-header.md` — the only file it sends you to is `./workflow.yaml`. Never hand-edit between the region markers; the splice overwrites whatever is there.

**2. Inline step.** A `workflow.yaml` step with `execution: inline` whose `skill:` names a shared file — `knowledge-recall.md`, `platform-context.md` (declared as the `platform-analysis` step), `decision-preservation.md`. The orchestrator reads that file and follows it in its own context. Nothing is injected anywhere and no sub-agent is launched.

**3. `reference_skills:` on a `subagent` step.** Before launching the step, the orchestrator reads each listed file and injects its content into the sub-agent's prompt as a `### REFERENCE SKILL {path}` block. The sub-agent never fetches them itself — sub-agents run from a different working directory and cannot resolve these paths. When a read fails, the block arrives as `### REFERENCE SKILL {path}: UNRESOLVED` and the step proceeds best-effort.

## Code standards

- Early return: `if err != nil { return err }` — validate inputs first.
- No global state — constructor injection throughout.
- Interfaces defined close to consumers, not in the implementing package.
- No `utils/`, `helpers/`, `common/`, or `misc/` packages — domain nouns only.
- Table-driven tests for any logic with more than two cases.

## PR process

- One logical change per PR.
- `go test ./...` must pass.
