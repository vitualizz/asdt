---
title: Architect
description: Makes architecture decisions and produces ADRs, system design, and API design artifacts — the specialist to bring in when a choice will shape service boundaries, data models, or scalability for the long haul.
order: 21
locale: en
---

# Architect (`/asdt-architect`)

> Makes architecture decisions and produces ADRs, system design, and API design artifacts — the specialist to bring in when a choice will shape service boundaries, data models, or scalability for the long haul.

## What it does

The Architect Specialist makes the technical decisions that everything else is built on. It evaluates competing approaches, documents the chosen path as an Architecture Decision Record (ADR), and produces a concrete system design with data models, API surfaces, and service boundaries — all before a single line of implementation code is written.

Every decision names two or three viable approaches and records, in the same block, why the chosen one won and why each alternative lost — that block is the decision record; there is no separate ADR artifact. It defaults to the simplest approach that satisfies the acceptance criteria, and spends extra rigor only where a choice is hard to reverse or visible to others. This forces honest trade-off analysis instead of post-hoc justification.

The Architect Specialist never writes implementation code, UX specs, or test plans. Its one job is to make the structural decision that the Developer can build against without ambiguity.

## When to invoke it

- A decision will shape service boundaries, data models, or scalability beyond the current feature
- The technical approach is non-obvious and has meaningful trade-offs between at least two viable options
- A cross-cutting concern (caching strategy, auth model, event bus) needs a documented decision
- You want a formal ADR to explain to future engineers why the code is the way it is

## On its own

No change in flight required. Point it at what already exists and it judges it instead of redesigning it — prioritized findings with evidence, and the strengths too:

```
/asdt-architect "does this structure scale if traffic triples?"
/asdt-architect "audit the boundaries of the payments module"
/asdt-architect "which decisions here are already expensive to reverse?"
```

What it finds is kept, so the next run over that area starts already knowing it.

## Pipeline position

Typically runs **after PM** and **before Developer** (Developer reads `architect/handoff`). It reads whatever exists upstream: `pm/handoff` (the requirements the design satisfies), `ux-ui/handoff` (the flows and component gaps the API surface has to serve), `security/handoff` (findings that reshape a boundary become design constraints), and `researcher/handoff` — which frames the problem only when PM was skipped; when PM ran, PM wins. All are optional. On simple changes it is not called at all — the Developer handles those directly. When it does run, it runs one step, `design`, and how deep that step goes is its own call.

## What it produces

`architect/handoff` — the decision and the system design that follows from it, in ONE hand-off:

- **The decision** — the chosen approach first, then each rejected alternative with why it lost
- **The design** — the data model and the API surface (or why the change has neither), the constraints the implementation has to respect, where in the codebase it lands, and the risks with their mitigations

Judging what already exists runs `review` instead and saves `{project}/study/{topic}/architect`.

Consumed by: **Developer** (the decision is settled — its spec restates it at implementation granularity), **QA** (the design and its declared risks), **Security** (the API surface and trust boundaries).

NFR budgets, when PM set any, arrive inside `pm/handoff.constraints` and the design has to live within them. When PM never ran, the design proceeds and records the gap rather than inventing a budget.

## Common patterns

```
/asdt-architect Design the rate-limiting strategy for the public API
# → Cross-cutting concern that will affect every endpoint
```

```
/asdt-architect Choose the event sourcing approach for the order pipeline
# → Non-reversible structural decision with meaningful trade-offs
```

```
/asdt-architect ADR for switching from REST to GraphQL on the mobile client
# → External contract change that needs documented rationale
```

## Limits — what it does NOT do

- Does not write implementation code
- Does not write UX specs or wireframes
- Does not produce test plans or acceptance criteria
- Never skips alternatives — every decision record requires them
- Does not design in isolation — always accounts for existing platform constraints
- System design is always incomplete without both a data model AND an API surface
