# Design Tokens — UX/UI Specialist

## Purpose
Derive the design token set (color, typography, spacing, radius, elevation, motion) for the
feature, grounded in the feature brief's `design_intent` and the platform's existing design
system. Prefer reusing established system values; never invent tokens that conflict with the
existing design system.

## Inputs
- `ux-ui/feature-brief`: `design_intent` (tone, principles, north_star), design_constraints
- `platform-summary`: design system fingerprint (existing tokens, css_approach, component_library)

`ux-ui/feature-brief` arrives injected per the parallel-retrieval contract — consume it directly.
`platform-summary` is injected by the orchestrator's inline `platform-analysis` step. Do not fetch
either yourself.

Extract from feature-brief: `design_intent` (drives tone/expressiveness of the token set).
Extract from platform-summary: existing token values and naming conventions (the fingerprint).

## Context budget
feature-brief `design_intent` + platform-summary fingerprint: max 900 tokens.

## Processing
1. GROUND every token in the platform fingerprint. Read the existing design system's values and
   naming convention FIRST.
2. For each token, prefer an EXISTING system value (`source: reused`). If the design_intent needs a
   variant of an existing value → `source: extended` (document the base it extends). Only when no
   existing value fits → `source: new`.
3. Translate `design_intent` (tone, principles, north_star) into token choices — e.g. a calm tone
   favors restrained motion durations and soft elevation; a bold tone favors higher contrast color
   roles.
4. NEVER invent tokens that conflict with the existing design system (e.g. a new primary color when
   one already exists). Conflicts go in `open_items`, not into the token set.
5. COMPARE the assembled token set against the generic-default archetypes — one archetype per token
   category. For every category you filled, either name the differentiation `design_intent`
   justifies, or admit the default was kept and say why keeping it is right for this platform.
   Record the outcome as exactly one verdict line, in this form:
   `design-heuristics: <differentiated|default-kept> — <one-line reason naming an archetype>`
6. Compute the `token_reuse_advisory` counts (reused / extended / new) and ratio. This is ADVISORY
   ONLY — it is never a hard gate and never blocks the workflow. Compose `note` from up to two
   parts: a short reuse-low sentence FIRST, only when reuse is low; then the verdict line from
   step 5, ALWAYS present and ALWAYS last. Join the two with ` | ` only when both are present —
   when the reuse-low sentence is absent, `note` starts at `design-heuristics:`, never with a bare,
   leading or trailing ` | `.

## Output
Produces: `ux-ui/design-tokens`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  tokens:
    color: []        # [{name, value, source: reused|extended|new, role}]
    typography: []   # [{name, size, line_height, weight, source}]
    spacing: []      # [{name, value, source}]
    radius: []       # [{name, value, source}]
    elevation: []    # [{name, value, source}]
    motion: []       # [{name, duration, easing, source}]
  token_reuse_advisory: {reused: 0, extended: 0, new: 0, ratio: "", note: ""}  # ADVISORY only
  # note format: "<reuse-low sentence — only when reuse is low> | design-heuristics: <differentiated|default-kept> — <one-line reason>"
  # the heuristics verdict is ALWAYS present and ALWAYS last; with no reuse-low sentence the note starts at "design-heuristics:" — never a bare " | "
  open_items: []
```
