---
title: UX/UI Designer
description: Shapes how people actually experience the product — user flows with their states and copy, designed for the project's primary surface and mapped to its existing components, with the accessibility each one owes — the specialist to bring in before a single screen gets built.
order: 25
locale: en
---

# UX/UI Designer (`/asdt-ux-ui`)

> Shapes how people actually experience the product — user flows with their states and copy, designed for the project's primary surface and mapped to its existing components, with the accessibility each one owes — the specialist to bring in before a single screen gets built.

## What it does

The UX/UI Specialist turns a requirement into flows a developer can build without guessing. It runs one step, `ux-spec`, and hands back one artifact, `ux-ui/handoff`, with these sections:

- **Brief** — four lines: the actor, their problem, what success looks like for them, and the one quality the experience should feel like.
- **Information architecture** — the entry point as a full path from the app's front door, the content hierarchy, and the primary actions.
- **User flows** — the deliverable. Numbered steps with every branch, the empty, loading, and error states named as steps, and the exact copy written inline wherever the wording carries the interaction.
- **Component mapping** — every flow step mapped to a component that already exists in the project, by its real name. Where nothing fits, it says so: that gap is a decision for the Developer, never a component quietly invented.
- **Accessibility** — per component: focus, keyboard, labelling, and the contrast pair it owes. What can't be verified from the values at hand is recorded as advisory, never asserted as a pass.

**Design system first.** When the project already has one — its own tokens, theme, or component library — that is the source of every token and component, and the specialist never invents a palette, type scale, or spacing unit against it. When it has none (a styling tool with no project-specific theme doesn't count), it adds a **visual direction**: a type pairing and scale, palette roles (`surface`, `text`, `accent`, `danger`, …) each with its contrast pair, a spacing scale, a density, and a motion stance. Every item is a proposal, which the Developer builds as the project's first tokens.

**Designed for your surface.** `/asdt-init` asks for the project's primary design surface — `mobile`, `tablet`, `desktop`, or `none`. Flows are designed for that surface first, with one line on what changes on the others; if it was never asked, mobile is assumed. When the answer is `none` — a CLI, a library, a backend service — there are no screens to specify: the specialist skips `ux-spec` and saves a short hand-off that says so, so downstream roles read an explicit "no UI" instead of a silent gap.

If the host assistant has a `frontend-design` skill installed, it is used to raise the quality of the visual direction and the layout choices. It never overrides the project's design system, and its absence changes nothing.

## When to invoke it

- A new screen, dialog, or feature-level UI needs to be designed
- User flows need to be mapped before architecture or implementation begins
- You need to know which existing components cover a change and where the gaps are
- The project has no design system yet and the first screens need a coherent direction
- Accessibility requirements need to be specified explicitly
- You want the Developer to receive a spec rather than infer the UX from the requirements

## On its own

Point it at a screen or a flow that already exists:

```
/asdt-ux-ui "review the accessibility of checkout"
/asdt-ux-ui "which design-system components are we not using in onboarding?"
```

That runs `review` instead: friction, missing states, design-system drift, and accessibility in what already ships.

## Pipeline position

Works best **after PM** (reads `pm/handoff` for the requirement and its acceptance criteria) and **before Architect and Developer**. The Architect reads the flows to shape the API surface it has to serve; the Developer builds them step by step. Running it after a screen is already built means the spec arrives too late to guide it.

## What it produces

`ux-ui/handoff` — brief, surface, information architecture, flows, visual direction (only when there is no design system), component mapping, and accessibility, as sections of ONE artifact.

Consumed by: **Architect** (the flows and component gaps the design has to serve), **Developer** (implements the flows, fills the gaps, builds the visual direction as first tokens), **QA** (every flow branch and every empty, loading, and error state becomes a candidate edge case, and the inline copy becomes exact-wording assertions).

## Common patterns

```
/asdt-ux-ui Design the onboarding flow for new users
# → New multi-step UI — IA and flows before any component work
```

```
/asdt-ux-ui Map the notification preferences screen
# → Existing UI to extend — the mapping names what can be reused and what is missing
```

```
/asdt-ux-ui Design the first screens of the admin panel
# → No design system yet — the hand-off carries a proposed visual direction
```

## Limits — what it does NOT do

- Does not write implementation code or any file — only the hand-off
- Does not produce architecture decisions or test plans
- Never invents tokens or components against an existing design system — an unmet need is a named gap
- Never writes a second set of flows per surface — one design for the primary surface, one adaptation line for the rest
- Never asserts an accessibility pass it can't verify from the tokens or the proposed palette
