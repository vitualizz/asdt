# UI Design — UX/UI Specialist

## Purpose
Design the UI the flows need, at the level a strong product designer hands a developer: every
screen's layout, hierarchy, typography, spacing, and states, every component the project lacks
designed part by part, and — without a design system — the visual foundation those designs
stand on. Written as text in the ONE role summary, which carries the flows forward.

## Inputs
- `### INPUT ux-ui/ux-spec` — the flows: read every section — `Design system`, `Brief`,
  `Surface`, `Entry point`, `Flows` (every step's screen and copy), `Components` (every gap),
  `Accessibility`, `Decisions`, `Open items`
- `### HOST SKILL frontend-design` — OPTIONAL. Visual-design craft from the host, injected only
  when it has one

The flows arrive as a path to read — never look for anything else. The platform block is in
this prompt: extract its `CSS:` line, the styling tool the designs are expressed in.

**DEGRADATION**: `ux-ui/ux-spec` UNRESOLVED → derive the flows from the raw request in its shape first — the surface from the platform block's `Design surface:` bare value, `mobile` when it carries none — and record `ASSUMED: no flows from ux-spec — flows derived in ui-design` in `Open items`.
**DEGRADATION**: no host skill → design on `visual-design.md` alone; no entry.

## Processing

Every craft criterion lives in `asdt-core/references/visual-design.md`; this step applies all of
them and adds only the procedure below.

1. **No UI — the explicit no.** When `ux-ui/ux-spec` has no `Flows` section — it found
   no user-facing step — return only the heading, the one-sentence what, and its `Decisions` line
   `Build no UI for this change (no user-facing step)`, and skip steps 2–8.

2. **Visual foundation — decided by `ux-ui/ux-spec`'s `Design system` verdict.**
   - **With a design system, it IS the foundation.** Read the token or theme files and the
     components its evidence points at, and design with their real token names,
     variants, and props. Never change an existing token or variant: a new variant ADDS an
     option in the system's language. A value the system lacks becomes a proposed token — name
     it in the system's convention with its value (`--space-7: 28px`), use that name in the
     design, and repeat it in one `Open items` entry, `Proposed token: --space-7: 28px in
     {token file}`, so `implement` adds it. Never a literal in its place.
   - **Without one, propose the foundation as `Visual direction`** — every field of the schema
     below. Every item is a PROPOSAL the developer builds as the project's first tokens, and
     every design below consumes it by role, level, and step, never by literal. Record it in
     `Decisions` as ONE line — its intent phrase, with the direction you rejected in
     parentheses — the values live in `Visual direction`, never twice.

   Either way, express values in the detected styling tool's terms when there is one (Tailwind
   classes and theme keys, CSS custom properties).

3. **Screens — one entry per distinct screen the flows name.** Fill every field of the
   schema. A component is either an existing one by its real name — variant `as-is` when it is
   reused unchanged — or a `Component designs` name. `States` covers every state the flows name
   on that screen, plus `success` and `disabled` where they occur. `Responsive` starts with the
   base layout for the primary surface, then one entry per further breakpoint, saying only what
   changes. Copy stays on the flow steps — never repeat it on the screen.

4. **Component designs — every gap, plus every new variant.** Design each gap under the name
   `Components` gave it (`new`), and an existing component only when the change needs a
   new variant of it (`variant of {Existing}` — `Variants` lists only the added option). A
   component reused as-is gets no entry. Its accessibility lives in its `Accessibility` row —
   add one when the flows had none, so every designed component has exactly one.

5. **Accessibility — resolve every pair.** Carry the `Accessibility` rows forward, and resolve each
   pair to its values and the ratio it reaches — the token values you read, or the
   `Visual direction` roles. A pair that misses `accessibility.md` changes the design, not the
   requirement; a pair you cannot compute stays advisory in `Open items`.

6. **Host skill.** When `### HOST SKILL frontend-design` arrived, use it to raise the intent's
   distinctiveness and the composition's coherence. Its instructions to write code or files do
   not apply here — translate what they would build into this schema's text. It never
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
Return this summary (`asdt-core/references/summary-schema.md`). Copy `Brief`, `Surface`, `Entry point`,
`Flows`, and `Components` verbatim from the flows; `Decisions` and `Open items` carry their
entries plus yours. `Design system` is not carried — a `Visual direction` present says there
was none. Omit `Visual direction` with a design system. Every `{…}` is filled, and every table
row repeats per item.

```markdown
# UX/UI: {change}

{The experience this change delivers, in one sentence.}

## Decisions
- {The flows' entries; the foundation's one line; each design call}

## Brief
## Surface
## Entry point
## Flows
## Components

## Visual direction
**Intent:** {phrase}
**Signature moves:**
- {type | color | composition}: {value}, replacing {default}

**Type faces:**
| Use | Family | Source | Fallback | Reason |
|---|---|---|---|---|
| {display / text} | {family} | {google-fonts / self-hosted / system} | {fallback stack} | {only for a default sans} |

**Type scale:**
| Level | Style |
|---|---|
| {display / title / heading / body / caption / label} | {face size/line-height weight} |

**Palette:**
| Role | Value | Dark | On | Ratio |
|---|---|---|---|---|
| {surface / surface-raised / text / text-muted / border / accent / accent-hover / accent-pressed / accent-subtle / on-accent / focus / danger / danger-hover / danger-pressed / danger-subtle / on-danger / success / …} | {value} | {only when dark mode is in scope} | {role it sits on; — for a surface} | {against On, per mode: `4.8:1 / 5.1:1 dark`} |

**Spacing:** {scale, px}
**Radius:** {name} {value}, …
**Elevation:** resting {shadow, or border + surface-raised}; raised {…}; overlay {…}
**Icons:** {set} ({source}); sizes {from 16 / 20 / 24}; stroke {width}
**Density:** {compact | regular | roomy}
**Motion:** micro {…}; transition {…}; easing {…}; reduced {…}

## Screens
### {Screen name, as the flows use it}
**Layout:** max width {…}; grid {columns and gutter, in the tool's terms}; alignment {…}

| Region | Size | Contents, in scanning order |
|---|---|---|
| {region} | {scale step or tool value} | {element}, {element} |

**Hierarchy**, in scanning order:
1. {element} — {type level or token; — for a component carrying its own}

**Primary action:** {the one element styled as primary}

**Spacing:**
- {A → B, or inside {region}}: {scale step or token}

**Components:**
| Element | Component | Source | Variant | Props |
|---|---|---|---|---|
| {hierarchy element} | `{real name, or a Component designs name}` | {existing / designed} | {`as-is`, or the variant} | {name=value, …} |

**States:**
- {empty-first-use / empty-no-results / empty-cleared / loading / error / success / disabled / …}: {exactly what differs from the default}

**Extremes** (screens showing data):
- {one item / very many / longest text / missing optional field}: {handling}

**Responsive** (first row: the base layout on the primary surface):
| Breakpoint | Regions |
|---|---|
| {`base`, `md ≥768px`, …} | {region}: {stacks / hides / drawer / columns: n / …} |

**Motion:** {only when it differs from the visual direction's motion}

## Component designs
### {Name} ({new | variant of `{Existing}`})
**Anatomy:** {part}, {part}, …

**Props:**
| Name | Type | Default |
|---|---|---|

**Variants:** {for a variant, only the added option}

**Sizes:**
| Name | Height | Padding | Type |
|---|---|---|---|

**States:**
| State | Changes, part by part |
|---|---|
| {default / hover / focus-visible / active / disabled / loading / error} | {…} |

**Tokens:**
| Part | State | Token |
|---|---|---|

**Behavior:** triggers {…}; emits {…}; overflow {…}; missing data {…}

## Accessibility
| Component | Focus | Keyboard | Labelling | Contrast |
|---|---|---|---|---|
| `{name}` | {behavior} | {interaction} | {label} | {pair, its values, and the ratio it reaches} |

## Open items
- ASSUMED: {anything unverified} — {what would confirm or refute it}
- Proposed token: {name}: {value} in {token file}
```
