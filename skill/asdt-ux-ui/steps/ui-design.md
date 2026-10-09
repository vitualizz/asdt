# UI Design — UX/UI Specialist

## Purpose
Design the UI the flows need, at the level a strong product designer hands a developer: every
screen's layout, hierarchy, typography, spacing, and states, every component the project lacks
designed part by part, and — without a design system — the visual foundation those designs
stand on. Written as text in the ONE hand-off, which carries the flows forward.

## Inputs
- `ux-flows` — injected from the orchestrator's context as `### INPUT ux-flows`. Extract
  everything: `design_system`, `brief`, `surface`, `entry_point`, `flows` (and every step's
  `screen` and `copy`), `components` (and every `gap`), `a11y_requirements`, `decisions`,
  `open_items`
- Platform summary — injected inline. Extract the `CSS:` line: the styling tool the designs are
  expressed in
- `### HOST SKILL frontend-design` — OPTIONAL. Visual-design craft from the host, injected only
  when it has one

All arrive ALREADY INJECTED — never self-fetch. If `ux-flows` did not arrive, derive the flows
from the raw request in its shape first — the surface from the platform summary's `Design
surface:` line's bare value (strip any parenthesised marker), `mobile` when it carries none — and record `ASSUMED: no ux-flows — flows
derived in ui-design` in `open_items`. If no host skill arrived, design on `visual-design.md` alone — no entry.

## Processing

Every craft criterion lives in `asdt-core/references/visual-design.md`; this step applies all of
them and adds only the procedure below.

1. **No UI — the explicit no.** When `ux-flows.flows` is empty — `ux-spec` found no user-facing
   step — return `screens: []` and `component_designs: []`, omit `visual_direction`, carry
   `ux-flows`' `decisions` line `Build no UI for this change (no user-facing step)`, and skip
   steps 2–8.

2. **Visual foundation — decided by `ux-flows.design_system`.**
   - **With a design system, it IS the foundation.** Read the token or theme files and the
     components `design_system.evidence` points at, and design with their real token names,
     variants, and props. Never change an existing token or variant: a new variant ADDS an
     option in the system's language. A value the system lacks becomes a proposed token — name
     it in the system's convention with its value (`--space-7: 28px`), use that name in the
     design, and repeat it in one `open_items` entry, `Proposed token: --space-7: 28px in
     {token file}`, so `implement` adds it. Never a literal in its place.
   - **Without one, propose the foundation as `visual_direction`** — every field of the schema
     below. Every item is a PROPOSAL the developer builds as the project's first tokens, and
     every design below consumes it by role, level, and step, never by literal. Record it in
     `decisions` as ONE line — its intent phrase, with the direction you rejected in
     parentheses — the values live in `visual_direction`, never twice.

   Either way, express values in the detected styling tool's terms when there is one (Tailwind
   classes and theme keys, CSS custom properties).

3. **Screens — one entry per distinct `screen` the flows name.** Fill every field of the
   schema. A component is either an existing one by its real name — `variant: as-is` when it is
   reused unchanged — or a `component_designs` name. `states` covers every state the flows name
   on that screen, plus `success` and `disabled` where they occur. `responsive` starts with the
   base layout for `surface.primary`, then one entry per further breakpoint, saying only what
   changes. Copy stays on the flow steps — never repeat it on the screen.

4. **Component designs — every `gap`, plus every new variant.** Design each gap under the name
   `components` gave it (`kind: new`), and an existing component only when the change needs a
   new variant of it (`kind: variant` — `variants` lists only the added option). A component
   reused as-is gets no entry. Its accessibility lives in its `a11y_requirements` entry — add one
   when `ux-flows` had none, so every designed component has exactly one.

5. **Accessibility — resolve every pair.** Carry `a11y_requirements` forward, and resolve each
   pair to its values and the ratio it reaches — the token values you read, or the
   `visual_direction` roles. A pair that misses `accessibility.md` changes the design, not the
   requirement; a pair you cannot compute stays advisory in `open_items`.

6. **Host skill.** When `### HOST SKILL frontend-design` arrived, use it to raise the intent's
   distinctiveness and the composition's coherence. Its instructions to write code or files do
   not apply here — translate what they would build into this schema's text fields. It never
   overrides the design system, the flows and their copy, or `accessibility.md`; on taste it
   outranks `visual-design.md`, on a checkable criterion it does not.

7. **Depth** changes how many screens and variants get full detail — a quick pass fills every
   field of every screen and gives full detail to the happy path's screens — never which
   fields exist. Every gap is designed at any depth.

8. **Conformance pass — silent.** Before returning, check every screen and component design
   against each criterion of `visual-design.md` and fix every failure in place. Nothing about
   the check enters the output: no score, no checklist, no notes.

No self-critique section: judging your own design without a user in front of it adds
confident prose and no signal. Step 8 fixes the design; it reports nothing.

## Output
Produces: `ux-ui/handoff`

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`. ONE artifact — flows, screens, component designs, visual
foundation, and accessibility are sections of it. `brief`, `surface`, `entry_point`, `flows`,
and `components` are carried verbatim from `ux-flows`; `decisions` and `open_items` carry its
entries plus yours. `design_system` is not carried — a `visual_direction` present says there
was none.

```yaml
payload:
  what: ""                    # the experience this change delivers, one sentence
  brief: {actor: "", problem: "", success: "", design_intent: ""}
  surface: {primary: "", adaptation: ""}
  entry_point: ""
  flows: []                   # verbatim from ux-flows — steps, each with its screen, branches, copy
  components: []              # verbatim from ux-flows — the mapping and its gaps
  visual_direction:           # OMIT the key with a design system, or when screens is empty
    intent:
      phrase: ""
      signature_moves: []     # 2–3 {area: type | color | composition, value: "", replaces: ""}
    type:
      faces: []               # {use: display | text, family: "", source: google-fonts | self-hosted | system, fallback: "", reason: ""}  # reason only for a default sans
      scale: []               # {level: display | title | heading | body | caption | label, style: "face size/line-height weight"}
    palette:
      - role: ""              # surface | surface-raised | text | text-muted | border | accent | accent-hover | accent-pressed | accent-subtle | on-accent | focus | danger | danger-hover | danger-pressed | danger-subtle | on-danger | success | …
        value: ""
        value_dark: ""        # only when the request, the brief, or the project's existing theme includes dark mode
        on: ""                # the role it sits on; "" for a surface role
        ratio: ""             # against `on`, per mode (`4.8:1 / 5.1:1 dark`)
    spacing: []               # the scale, px
    radius: []                # {name, value}
    elevation: []             # {level: resting | raised | overlay, value: shadow, or border + surface-raised}
    icons: {set: "", sizes: [], stroke: ""}   # set with its source; sizes from 16 / 20 / 24
    density: ""               # compact | regular | roomy
    motion: {micro: "", transition: "", easing: "", reduced: ""}
  screens:
    - name: ""                # the `screen` value the flows use
      layout:
        max_width: ""
        grid: ""              # columns and gutter, in the styling tool's terms
        alignment: ""
        regions:
          - name: ""
            size: ""          # height or width, as a scale step or tool value
            contents: []      # element names, in scanning order
      hierarchy:              # every element, in scanning order
        - element: ""
          type: ""            # scale level or token; "" for a component carrying its own
      primary_action: ""      # the one element styled as primary
      spacing:
        - between: ""         # `A → B`, or `inside {region}`
          step: ""            # scale step or token
      components:
        - element: ""         # the hierarchy element it renders
          name: ""            # real name, or a component_designs name
          source: ""          # existing (reused as-is) | designed (a new component, or a `kind: variant` of an existing one)
          variant: ""         # `as-is` when an existing component is reused unchanged
          props: []           # {name, value}
      states:
        - state: ""           # empty-first-use | empty-no-results | empty-cleared | loading | error | success | disabled | …
          changes: ""         # exactly what differs from the default
      extremes: []            # screens showing data: {case: one item | very many | longest text | missing optional field, handling: ""}
      responsive:             # first entry: the base layout on surface.primary
        - breakpoint: ""      # tool name and width (`base`, `md ≥768px`)
          regions: []         # {name, does} — stacks | hides | drawer | columns: n | …
      motion: ""              # OMIT unless it differs from visual_direction.motion
  component_designs:
    - name: ""                # the gap's name in components, or the existing component
      kind: ""                # new | variant
      anatomy: []             # named parts
      props: []               # {name, type, default}
      variants: []            # kind: variant → only the added option
      sizes: []               # {name, height, padding, type}
      states:
        - state: ""           # default | hover | focus-visible | active | disabled | loading | error
          changes: ""         # part by part
      tokens: []              # {part, state, token}
      behavior: {triggers: "", emits: "", overflow: "", missing_data: ""}
  a11y_requirements:
    - component: ""
      focus: ""
      keyboard: ""
      labelling: ""
      contrast: ""            # the token pair, or the resolved role pair and its ratio
  decisions: []               # ux-flows' entries, the foundation's one line, design calls
  open_items: []              # ASSUMED: prefix for anything unverified; `Proposed token:` entries
```
