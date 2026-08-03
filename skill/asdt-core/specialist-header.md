<!-- specialist-header.md — reaches a specialist context by two mechanisms only, never by a frontmatter key: (1) the installer splices this body directly into each routed SKILL.md at install time, so an installed specialist already carries it; (2) the explicit FIRST ACTION Read instruction at the top of every specialist SKILL.md body, which loads it when the body was not spliced. -->

## Prerequisites

Before starting any step, verify `.asdt/config.yaml` exists with `memory.provider` set and that the provider is reachable. If either fails, output this message and STOP:

> Memory provider not configured. Run `/asdt-init` and set `memory.provider` in `.asdt/config.yaml` before running any specialist.

> **ORCHESTRATOR GATE**: This file is a PLAN, not an executable pipeline. The
> calling assistant (Claude Code / OpenCode) is the SOLE orchestrator. For every
> step marked `subagent` below you MUST launch a dedicated sub-agent via your
> native delegation primitive (Agent/Task) — do NOT run subagent steps inline in
> this thread. Steps marked `inline` run in your own context. This specialist file
> NEVER calls Agent/Task itself; it only tells YOU, the orchestrator, what to launch.
> If you run a subagent step inline anyway, that step's write boundary
> (`asdt-core/protocol.md` §3) binds YOU: unless the step is `developer/implement`
> or `asdt-init/write`, you write no files at all.

## Depth

Judge the depth of this run yourself, from two signals in order: (1) anything the invocation says or implies about scope or urgency ("quick", "a fondo", "this is sensitive", "just a sanity check") — honor it; (2) absent that, the size of what you find: how much the change touches and how much is genuinely unknown. Depth controls how exhaustive each step's OUTPUT is. Which steps run is defined by this specialist's own SKILL.md — single-step specialists always run their one step; the Developer judges its chain per its own table. Never ask the user to pick a depth level.

## Intent

Judge from the invocation whether this run DELIVERS a change or EXAMINES what already exists. A change persists under `{project}/{change}/{role}/handoff`; an examination — an audit, a review, an assessment with nothing to deliver — persists under `{project}/study/{topic}/{role}`, with `{topic}` derived from the request in short kebab-case. Same schema, same rules. Specialists with a `review` step run it for examinations; the others run their normal steps under the study key. Never ask the user to pick; the phrasing decides, and when genuinely ambiguous, treat it as a change.

## Narration

Narrate to the user in prose. Topic keys, schema fields, and step names are internal machinery, not conversation.

Every run ends by presenting its findings to the user as a short prose report: what was found, what was decided, what stayed `ASSUMED:`. The persisted hand-off is the record for the team; the narrated report is the answer to the person. **The report is never persisted** — there is no report artifact.

When a run closes on a negative verdict or a high risk, the report's last line proposes the next invocation in plain language, as a sentence they can copy — "this would be fixed by `/asdt-developer \"cover the concurrency edge cases qa/handoff left open\"`". It is a sentence, not a mechanism: the user decides.

## The contract

Inputs, Engram persistence, injection format, and degradation live in `asdt-core/protocol.md` — read it now if you do not already have it in context. Per-step identity, model, inputs, and outputs live in this directory's `workflow.yaml`.
