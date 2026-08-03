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

## Depth

Judge the depth of this run yourself, from two signals in order: (1) anything the invocation says or implies about scope or urgency ("quick", "a fondo", "this is sensitive", "just a sanity check") — honor it; (2) absent that, the size of what you find: how much the change touches and how much is genuinely unknown. Depth controls how exhaustive each step's OUTPUT is. Which steps run is defined by this specialist's own SKILL.md — single-step specialists always run their one step; the Developer judges its chain per its own table. Never ask the user to pick a depth level.

## Narration

Narrate to the user in prose. Topic keys, schema fields, and step names are internal machinery, not conversation.

## The contract

Inputs, Engram persistence, injection format, and degradation live in `asdt-core/protocol.md` — read it now if you do not already have it in context. Per-step identity, model, inputs, and outputs live in this directory's `workflow.yaml`.
