---
title: UX/UI Designer
description: Designs how a change looks and works before a single screen gets built — user flows with their states and copy, then every screen and every missing component designed in detail, in the project's design system or a proposed visual foundation, with the accessibility each one owes.
order: 25
locale: en
---

# UX/UI Designer (`/asdt-ux-ui`)

> Designs how a change looks and works before a single screen gets built — user flows with their states and copy, then every screen and every missing component designed in detail, in the project's design system or a proposed visual foundation, with the accessibility each one owes.

## What it does

The UX/UI Specialist does the work of a product designer: it decides how the change flows, then designs the UI it needs, precisely enough that the Developer builds it without guessing. It runs two steps — `ux-spec` writes the flows, `ui-design` designs the screens and components — and hands back one artifact, `ux-ui/handoff`. The design is text inside that hand-off; no file is written.

From `ux-spec`:

- **Brief** — four lines: the actor, their problem, what success looks like for them, and the one quality the experience should feel like.
- **Information architecture** — the entry point as a full path from the app's front door, the content hierarchy, and the primary actions.
- **User flows** — numbered steps with every branch, the screen each one happens on, the empty, loading, and error states named as steps, and the exact copy written inline wherever the wording carries the interaction.
- **Component mapping** — every flow step mapped to a component that already exists in the project, by its real name. Where nothing fits, the gap gets a name — and `ui-design` designs it.
- **Accessibility** — per component: focus, keyboard, labelling, and the contrast pair it owes. What can't be verified from the values at hand is recorded as advisory, never asserted as a pass.

From `ui-design`:

- **Screens** — one per screen the flows touch: layout (max width, grid, alignment, and each region with its size and contents), hierarchy in scanning order with each element's type level and one primary action, spacing as scale steps, the components it composes, every state and exactly what changes in it, how it handles data extremes (one item, very many, the longest text, a missing field), what each region does at each breakpoint, and motion only where it differs from the foundation. The copy stays on the flow steps. For example:

  ```
  Screen: Request reset
    Layout: centered column, max 420px
    Hierarchy: H1 → text → input → CTA
    Typography: H1 Fraunces 28/34 600 (Google Fonts)
    States: empty / sending / error
  ```

- **Component designs** — every component the flows need and the project lacks, plus any new variant of an existing component, designed part by part: anatomy, props, variants, sizes, states (default, hover, focus-visible, active, disabled, loading, error), the token each part consumes, and behavior. A component reused as-is gets no design — the screen that uses it marks it as-is.

The design craft is built in: a shared visual-design reference — one aesthetic intent made visible in two or three signature moves, a type scale on named faces with their source, a spacing scale and breakpoints, hierarchy, color by role with interactive families and contrast pairs, density, elevation, and icons, states — empty variants and data extremes included — as designed artifacts, motion with reduced-motion respected, visible focus rings, and component anatomy and props — so the quality does not depend on any plugin. Before handing off, `ui-design` checks its design against that reference and fixes what fails.

**Design system first.** When the project already has one — its own tokens, theme, or component library — that is the foundation: its tokens and components are used as they are, never restyled: a new variant adds an option in the system's language and leaves the existing ones unchanged, a value the system lacks is proposed as a named token (`--space-7: 28px`) for the Developer to add, and a component it lacks is designed in its language. When it has none (a styling tool with no project-specific theme doesn't count), `ui-design` proposes a **visual direction**: the intent and its signature moves, the faces with their source and a type scale, palette roles (`surface`, `text`, `accent` with its hover, pressed, and on-accent values, `focus`, `danger`, …) each with its contrast pair, a spacing scale, radius, elevation, icons, a density, and a motion stance. Every item is a proposal, which the Developer builds as the project's first tokens.

**Designed for your surface.** `/asdt-init` asks for the project's primary design surface — `mobile`, `tablet`, `desktop`, or `none`. Flows and screens are designed for that surface first, with one line on what changes on the others; if it was never asked, mobile is assumed. When the answer is `none` — a CLI, a library, a backend service — there are no screens to design: the specialist skips both steps and saves a short hand-off that says so, so downstream roles read an explicit "no UI" instead of a silent gap. A change with no user-facing step — a background job, an API-only change — gets the same explicit no on any surface: both steps run and return no flows and no screens.

If the host assistant has a `frontend-design` skill installed, `ui-design` uses it to raise the craft further. It never overrides the project's design system, the flows, or accessibility, and its absence changes nothing.

## When to invoke it

- A new screen, dialog, or feature-level UI needs to be designed
- User flows need to be mapped before architecture or implementation begins
- You need to know which existing components cover a change, and the missing ones designed
- You want the screens designed — layout, hierarchy, typography, states — before anyone builds them
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

`ux-ui/handoff` — brief, surface, information architecture, flows, component mapping, screens, component designs, visual direction (only when there is no design system), and accessibility, as sections of ONE artifact.

Consumed by: **Architect** (the flows and component gaps the design has to serve), **Developer** (builds each screen and component to its design, implements the flows, builds the visual direction as first tokens, adds the tokens it proposed), **QA** (every flow branch, every screen state and data extreme, and every designed component's disabled, loading, and error state becomes a candidate edge case, and the inline copy becomes exact-wording assertions).

## Common patterns

```
/asdt-ux-ui Design the onboarding flow for new users
# → New multi-step UI — IA and flows before any component work
```

```
/asdt-ux-ui Map the notification preferences screen
# → Existing UI to extend — the mapping names what is reused as-is, and what is missing gets designed
```

```
/asdt-ux-ui Design the first screens of the admin panel
# → No design system yet — the screens are designed on a proposed visual direction
```

## Limits — what it does NOT do

- Does not write implementation code or any file — the design lives in the hand-off
- Does not produce architecture decisions or test plans
- Never restyles an existing design system or invents a token against it — a missing component is designed in its language, a new variant leaves the existing ones untouched, a missing token is proposed by name and value as an open item
- Never writes a second design per surface — one for the primary surface, one adaptation line for the rest
- Never asserts an accessibility pass it can't verify from the tokens or the proposed palette
