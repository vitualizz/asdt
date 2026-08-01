<!-- specialist-header.md — reaches a specialist context by two mechanisms only, never by a frontmatter key: (1) the installer splices this body directly into each routed SKILL.md at install time, so an installed specialist already carries it; (2) the explicit FIRST ACTION Read instruction at the top of every specialist SKILL.md body, which loads it when the body was not spliced. -->

## Prerequisites

Before starting any step, verify:
1. `.asdt/config.yaml` exists with `memory.provider` set
2. The memory provider is reachable (Engram MCP server is running)

If either condition is not met, output this exact message and STOP:

> Memory provider not configured. Run `/asdt-init` and set `memory.provider` in `.asdt/config.yaml` before running any specialist.

> **ORCHESTRATOR GATE**: This file is a PLAN, not an executable pipeline. The
> calling assistant (Claude Code / OpenCode) is the SOLE orchestrator. For every
> step marked `subagent` below you MUST launch a dedicated sub-agent via your
> native delegation primitive (Agent/Task) — do NOT run subagent steps inline in
> this thread. Steps marked `inline` run in your own context. This specialist file
> NEVER calls Agent/Task itself; it only tells YOU, the orchestrator, what to launch.

> **Plain-Language Narration**: everything you say to the user is prose they can
> read without a decoder. Never surface an internal code as-is — step names,
> `topic_key` paths, artifact IDs, tier keywords, and schema fields are machinery,
> not conversation. Say what happened and why it matters; if a code must appear,
> give its meaning in the same breath. This governs narration only: persisted
> artifacts and YAML payloads keep their exact machine wording.

> **Before driving**: read `workflow.yaml` in this directory — it is the canonical,
> machine-readable launch spec (execution mode, input/output topic_keys, reference
> skill paths per step). The table below is a human-readable summary.

> **Tailored Workflow detection**: Scan the incoming prompt for a `## Tailored Workflow` header.
> - If ABSENT: run the full default workflow defined in the step table below.
> - If PRESENT: parse the `steps:` list. Execute ONLY those steps in the order specified.
> - Steps NOT in the tailored list → skip entirely (log annotation that the step was skipped by workflow tailoring).
> - The tailored list overrides the default ordering.
> - A PRESENT block also means the router's clarifying gates already ran — never re-ask a question the routing plan already settled; record any remaining gap in `open_items` instead.
> - The block MAY carry a `depth: quick|standard|deep` field controlling per-step OUTPUT VERBOSITY
>   (orthogonal to WHICH steps run — it never adds or removes steps):
>   - `depth` omitted ⇒ treat as `standard` (no behavior change).
>   - `quick` ⇒ minimum viable output: collapse enumerations to the most important item(s), skip
>     OPTIONAL fields, keep rationale terse. NEVER omit hard schema-required fields.
>   - `standard` ⇒ the current default output volume.
>   - `deep` ⇒ exhaustive: maximize alternatives, edge-cases, and rationale within the existing
>     per-step context budget.

<!-- co-location note: The specialist-specific complexity/tier paragraph and the Artifact Persistence block remain inline in each specialist SKILL.md. Do NOT move them here. The closing Tailored Workflow sentence also stays per-specialist. -->

**Execution policy (the rule, not just the list)**: a step that produces its OWN
persisted artifact (generative / decision-producing) is `subagent`; a step that
produces no artifact of its own and only injects context for the next step
(recall / wrapper) is `inline`. If steps change later, re-apply this rule.

### How to launch a `subagent` step

> Canonical protocol: the **Orchestrator Fetch-Once Cache & Injected-Input Contract** sections
> of this same document — **Cache Ledger Rule**, **Injection Format**, and the UNRESOLVED case in
> **Partial-Context Degradation**. Follow them as written there; they are not restated here.

> **Agent type**: launch each `subagent` step with the agent type its
> `workflow.yaml` entry declares — `agent: analyst` maps to the installed
> `asdt-analyst` agent, `agent: builder` to `asdt-builder`. Fallback is
> MANDATORY: when the named type is not available in your harness, launch with
> the harness default (general-purpose) agent AND prepend the injected
> executor header per **Executor Header Injection** in this same document.

`inline` steps fold into your own orchestrator context — no launch.

### How to resolve declared skill files

`workflow.yaml` declares skill files in two places. YOU, the orchestrator, resolve
both — a declaration nobody reads is a declaration that does nothing.

**`reference_skills:` on a `subagent` step.** Before launching that step, read each
listed file yourself and inject its content into the sub-agent's prompt. The
sub-agent NEVER fetches these itself: sub-agents run from a different working
directory and cannot resolve these paths.

Inject one block per file, resolved:

```
### REFERENCE SKILL {path}
{file content}
```

Or, when the read fails:

```
### REFERENCE SKILL {path}: UNRESOLVED
(orchestrator could not read this skill file — record it in open_items and proceed)
```

**`skill:` on an `inline` step.** Read that file and follow it in your own context.
Nothing is injected anywhere, because inline steps run in the orchestrator. (On a
`subagent` step, `skill:` names the step file that becomes the sub-agent's prompt —
not something you follow yourself.)

Rules for both:

- Paths resolve from the SPECIALIST's own directory. Both
  `../asdt-shared/skills/x.md` and `skills/x.md` are valid and resolve correctly.
- Reads go through the same fetch-once ledger described under **Cache Ledger Rule** in
  this same document, keyed `skill:{path}` so a file entry can never collide with a
  `topic_key` entry. Follow the rule as written there; it is not restated here.
- A declared skill path whose content is already spliced into this header — the sections
  that follow below — is ALREADY resolved: mark its `skill:{path}` ledger entry satisfied
  and do NOT read the file again.
- An inline step with NO `skill:` key at all (`knowledge-gate`, `enrichment`, and
  `clarify` in `asdt-init`) has nothing to read. That is normal — it degrades
  silently and is NOT an UNRESOLVED case.
