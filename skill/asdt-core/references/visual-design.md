# Visual Design — Reference

The design craft ASDT ships. Reference for UX/UI's `ui-design`, and for anyone designing a screen or a component. Every criterion below is checkable on the design itself. A design that fails one is not finished.

## Intent

- Commit to ONE aesthetic intent, named in a phrase that only this product could carry: derive it from who uses the product, where, and what they need to feel ("calm ledger for tired accountants", not "clean and modern"). Every choice below serves it; a choice that serves no intent is a default.
- Make the intent visible in 2–3 signature moves — one in type, one in color or surface, one in composition — each a concrete value and the default it replaces (`display Fraunces 32/38 600, replacing a sans title`; `surface #F6F1E7 warm paper, replacing white`; `left-aligned column on an 8-col grid, replacing a centered card`). A move with no value, or no replaced default, is not a move.
- Never ship template defaults: no system-font-plus-blue-button, no evenly grey card grid, no gradient for its own sake. If a choice would look identical in any other product, make it again.
- Inside an existing design system, the intent is the system's and it brings no signature moves of yours. Consistency outranks novelty: extend in its language — its tokens, radius, density, and motion — so a new component looks like it was always there.

## Typography

- A scale of 4–6 levels, each with a job: display, title, heading, body, caption (plus label for controls). More levels means the hierarchy is unclear.
- Adjacent levels differ visibly — a size ratio of at least 1.2, or a weight change of two steps. Two levels a reader cannot tell apart are one level.
- Write each level as `face size/line-height weight` (`Fraunces 32/38 600`). Line-height: 1.4–1.6 for body, 1.1–1.3 for display and titles.
- At most two faces: one for display, one for everything else — often the same family. Name each with its source (Google Fonts, self-hosted, or system) and a fallback stack (`Fraunces, Georgia, serif`). A default sans — Inter, Roboto, system-ui — carries a one-line reason tied to the intent; without one it is a template default.
- Body text is at least 16px on touch surfaces, form inputs included (smaller inputs make mobile browsers zoom); 14px body only in compact desktop data UIs.
- Measure: body text runs 45–75 characters per line; cap the text column, never the viewport.
- Numbers in tables and amounts use tabular figures.

## Spacing and layout

- One spacing scale, on a 4px base (4, 8, 12, 16, 24, 32, 48, 64). Every gap, padding, and margin is a scale step; an off-scale value is a bug.
- Proximity is grouping: space inside a group is smaller than space between groups — at least one scale step apart.
- Every screen states its structure: max content width, grid (columns and gutter), alignment, and its regions — each with its size and its contents. Align to few edges; each new edge is noise.
- Breakpoints are widths in the styling tool's terms (Tailwind `md` = 768px, or the system's own). Per breakpoint, say what each region does — stacks, hides, becomes a drawer, changes column count. "Adapts to smaller screens" is not a design.

## Hierarchy

- ONE primary action per screen, styled as the only primary. Secondary actions are visibly quieter; destructive ones are distinct and never the default.
- Order the elements in the sequence the user scans them: what this is, what it says, what to do. The hierarchy list IS that sequence.
- Emphasis is a budget: size, weight, color, and position each spend it. If everything is emphasized, nothing is.

## Color

- Define color by ROLE, never by hue: `surface`, `surface-raised`, `text`, `text-muted`, `border`, `accent`, `danger`, `success`, plus `warning` or `info` only when a flow uses them.
- An interactive role is a family: `accent`, `accent-hover`, `accent-pressed`, `accent-subtle` (tinted background), and `on-accent` (the text and icons on it); the same family for `danger`. ONE `focus` role draws every focus ring.
- Every role that carries text or UI is declared with the role it sits on and its ratio, and meets `accessibility.md` — 4.5:1 for body text, 3:1 for large text and UI parts. A pair you cannot compute is unverified, never a pass.
- `accent` marks the primary action, nothing decorative. Meaning never travels by color alone — pair it with an icon, text, or shape.
- Dark mode, when in scope, gives every role its own dark value with its pairs re-checked — never an inversion.

## Density, radius, elevation, icons

- Density fits the surface and the task: compact for data-dense desktop tools, regular for forms, roomy for first-run and marketing. Touch targets are at least 44×44px on touch surfaces.
- One radius family (e.g. 4 / 8 / 16).
- At most three elevation levels — resting, raised, overlay — each a concrete shadow value (`0 1px 2px rgb(0 0 0 / 0.08)`); a flat style states each level as a border plus `surface-raised` instead.
- One icon set, named with its source; sizes from 16 / 20 / 24, each paired with the type level it sits beside; one stroke width. Two sets on one screen is a bug.

## States

States are designed artifacts, specified as exactly what changes from the default:

- **Empty** — an invitation: what goes here, why it matters, and the one action that fills it. When they differ, design first use, no results, and cleared separately.
- **Loading** — preserves the layout: skeletons or placeholders sized like the content; no layout shift when data lands. Over ~1s, say what is happening.
- **Error** — what happened, in the user's words, and what to do next, placed where the problem is (inline at the field, not only a toast). Input is never lost.
- **Success** — confirms the outcome and points to the next step; transient unless the user must act on it.
- **Disabled** — explains itself when the reason is not obvious; prefer an enabled control with a message.
- **Extremes** — every screen that shows data names how it handles one item, very many, the longest realistic text, and a missing optional field.

## Motion

- Motion has a purpose: show where something came from, what changed, or that the system heard the input. Decorative motion is cut.
- Micro-interactions 100–200ms; transitions between views 200–300ms; ease-out on enter, ease-in on exit.
- Respect `prefers-reduced-motion`: every movement has a static or fade-only fallback.

## Component design

A designed component specifies every one of these — a missing one is a question the developer will answer by guessing:

- **Anatomy** — its named parts (container, label, icon, helper text, …).
- **Props** — its API: each prop's name in the project's convention, its type, and its default. A variant or a size is a prop value, never a separate component.
- **Variants** — by intent (primary, secondary, danger) and only the ones a flow uses. A new variant of an existing component ADDS an option in the system's language; every existing variant stays unchanged.
- **Sizes** — the scale steps it comes in, each with its height, padding, and type level.
- **States** — default, hover, focus-visible, active, disabled, loading, error: what each changes, part by part.
- **Focus-visible** — a ring in the `focus` role, its width and offset stated (2px, 2px offset), at least 3:1 against the surface it sits on, never a color change alone, and never the hover style.
- **Tokens** — the token each part consumes, per state. No literal values.
- **Behavior** — what triggers it, what it emits, how it handles overflow and long text, and what it shows when data is missing.
