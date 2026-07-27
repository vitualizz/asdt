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

**Workflow** — an ordered list of steps specific to that discipline, gated by complexity so the depth matches the change. At `complex`, the Developer runs `explore → spec → design → tasks → implement → test (if TDD)`; at `simple` it runs just `explore → spec → implement`. The UX/UI specialist at `complex` runs `feature-brief → design-tokens → information-architecture → user-flows → content-design → component-mapping → design-critique → ux-handoff`. These are not the same pipeline applied to different names — each specialist's workflow reflects how that discipline actually works.

**Skill composition** — shared skills (platform context, knowledge recall, scope definition) plus specialist-scoped skills (threat modeling for Security, code generation for Developer). Nothing loads ambiently: a shared skill is read only where it is declared, either as an `inline` step in `workflow.yaml` or in a step's `reference_skills:` list. Capabilities are mixed in rather than inherited.

**Artifact contract** — what the specialist reads (`inputs`) and writes (`outputs`). Inputs are soft: a missing input degrades to an `open_items[]` note, never an error. Outputs have stable `topic_key` values so other specialists can retrieve them by key.

## Adding a specialist

Adding a new specialist requires exactly two things:

1. One `SpecialistDescriptor` value literal in the registry
2. One `skill/{id}/SKILL.md` tree with step files

Zero new Go packages, zero new switch arms. The `asdt-*` embed glob in `skill/embedded.go` picks up any directory matching the pattern and ships it in the next build. See [Contributing](/asdt/docs/contributing) for the full authoring contract.

## The independence guarantee

Any specialist may run first — there is no required predecessor. If the Developer finds no Architect artifact in Engram, it proceeds with `open_items: ["architect/adr not found"]` and makes reasonable assumptions. The resulting implementation artifact is less precise than if the Architect had run first, but it's valid output.

This design choice prioritizes flexibility over correctness guarantees. You can always run specialists out of order. ASDT trusts you to decide when to involve each discipline.
