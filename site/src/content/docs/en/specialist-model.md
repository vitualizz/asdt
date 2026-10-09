---
title: Specialist Model
description: How ASDT models software delivery as a team of independent specialists, each owning a discipline.
order: 6
locale: en
---

# Specialist Model

## Why specialists, not a pipeline

The first version of ASDT modeled software delivery as a fixed four-phase FSM (a finite-state machine — a rigid flow that only moves through a fixed set of steps in a fixed order): `requirements → plan → implement → review`. Adding a new role required a new Go package, a new struct, and a new switch arm — code, not prompt authoring. The FSM hardcoded `requirements` as the only valid entry point, so a security engineer or UX designer had no valid place in the model without restructuring the entire graph.

This is the wrong model. Real software delivery is performed by a team of specialists, each owning an independent discipline. A security engineer doesn't wait for a developer to finish before reviewing auth code. A UX designer doesn't follow a requirements → plan workflow — they follow their own creative process.

So ASDT replaced the FSM with a different unit: a Specialist is a composable, independent unit defined by its identity, its own workflow steps, its artifact contract, and an independence guarantee — any specialist may run first, with no required predecessor.

## What defines a specialist

A specialist has four parts:

**Identity** — a stable `id` (e.g. `developer`), a human name, and a description that the pipeline advisor uses to route requests.

**Workflow** — a short list of steps specific to that discipline, declared in its `workflow.yaml`. The specialist judges which of them the request needs; depth changes how thorough each step's output is, never which steps exist. The Developer picks its chain from what you ask: a question runs `explore`; a plan runs `explore → spec` and saves the plan so it can be resumed later; a build runs `explore → spec → approve → implement → verify`, pausing for your approval before any file is written; and "implement the plan we approved" resumes at `approve → implement → verify` — or at `verify` alone, when the plan was built but never checked. The UX/UI specialist runs a single `ux-spec` — flows for the project's design surface, mapped to its components — or `review` when you ask it to audit what already ships. These are not the same pipeline applied to different names — each specialist's workflow reflects how that discipline actually works.

**Skill composition** — shared references (platform context, knowledge recall, scope definition, OWASP, accessibility) declared per step, plus optional skills the host assistant may already have installed (`host_skills:`, e.g. `frontend-design` for UX/UI and for the Developer's UI files). Nothing loads ambiently: a shared reference is read only where it is declared, either as an `inline` step in `workflow.yaml` or in a step's `reference_skills:` list. Capabilities are mixed in rather than inherited.

**Artifact contract** — which teammates' hand-offs the specialist reads (`inputs`) and the one hand-off it writes, at the stable key `{project}/{change}/{role}/handoff`, so other specialists can retrieve it by key. Inputs are soft: a missing input degrades to an `ASSUMED:` note in `open_items`, never an error.

## Adding a specialist

Adding a new specialist requires two things:

1. One `skill/asdt-{id}/` directory — `SKILL.md`, `workflow.yaml`, and one step file per sub-agent step
2. Its row in the routing tables — the `## Registry` table in `skill/SKILL.md` and the `## ASDT Specialists` table in the installed agents template — plus the routed list in `skill/embedded_test.go`

Zero new Go packages, zero new switch arms. The `asdt-*` embed glob in `skill/embedded.go` picks up any directory matching the pattern and ships it in the next build. See [Contributing](/asdt/docs/contributing) for the full authoring contract.

## The independence guarantee

Any specialist may run first — there is no required predecessor. If the Developer finds no PM hand-off in Engram, it writes the acceptance criteria itself and records `ASSUMED: no PM hand-off — acceptance criteria authored from the exploration` in `open_items`. Its report says so in its first line, which names what it built on and what was missing. The result is less precise than if PM had run first, but it's valid output.

This design choice prioritizes flexibility over correctness guarantees. You can always run specialists out of order. ASDT trusts you to decide when to involve each discipline.
