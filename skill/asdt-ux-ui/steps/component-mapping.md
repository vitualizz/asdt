# Component Mapping — UX/UI Specialist

## Purpose
Identify which existing components to reuse, which to extend, and which must be created.
Maximize reuse. Justify every new component.

## Inputs
- `ux-ui/flows`: interaction sequences, state changes
- `platform-summary`: component_library, existing patterns
- `ux-ui/content-inventory`: text touchpoints (TextTouchpoint list) — informs content-driven component sizing.

`ux-ui/flows` and `ux-ui/content-inventory` arrive ALREADY INJECTED per the parallel-retrieval
contract — consume them directly. `platform-summary` is injected by the orchestrator's inline
`platform-analysis` step. Never self-fetch any of them.

**DEGRADATION — `ux-ui/content-inventory` is optional (produced only at the moderate and complex
tiers)**: when it arrives as `### INPUT ux-ui/content-inventory: UNRESOLVED`, proceed on the flows
alone and skip content-driven component sizing; append "content inventory absent — component sizing
derived from flows only, without text-surface length data" to open_items. Never block on this input.

Extract from flows: all unique UI states and interactions.
Extract from platform-summary: component_library name and approach.

This step owns responsive behavior for every component it maps, defines each component's
`state_matrix` and `accessibility` block, and is the SOLE authority for `reuse_ratio`.

## Context budget
flows summary (UI states only) + platform-summary + content-inventory touchpoints (when present):
max 2,600 tokens total.

## Processing
For each UI state identified in the flows:
1. CHECK if an existing component handles this state. If yes → REUSE.
2. If the existing component needs minor changes → EXTEND (document what changes).
3. Only if no existing component fits → CREATE NEW (document why existing ones don't work).

For each new component:
- Name it following the project's naming convention (from platform-summary).
- Define its props/inputs (data it needs).
- Define its events/outputs (what it emits).
- Note its responsive behavior (how it changes at breakpoints).
- Fill its `accessibility` block: the `aria_role` it presents, the `keyboard_interaction` it must
  support, its `focus_management` rule, and the `contrast_token_ref` its text/UI pairs resolve
  against. This step runs at EVERY tier that maps components, so this block is where the WCAG
  checklist results are recorded even when `design-critique` never runs. State an honest gap rather
  than inventing a value: anything you cannot determine goes to `open_items`, never a guessed field.

### State matrix (per interactive component)
For each interactive component, define the 9-state `state_matrix`: `default`, `hover`, `focus`,
`active`, `disabled`, `empty`, `loading`, `error`, `success`. Each state is
`{applicable: bool, behavior: ""}` — `behavior` is REQUIRED whenever `applicable: true`. Mark
non-applicable states `{applicable: false}` with no behavior.

### Responsive (mobile-first)
Define a top-level `breakpoint_strategy` and a per-component `responsive` spec. Mobile-first: start
at the smallest viewport and expand up.
1. MOBILE (320-767px): layout, what is visible, what collapses or stacks.
2. TABLET (768-1023px): changes from mobile.
3. DESKTOP (1024px+): full desktop layout.
4. TOUCH TARGETS: confirm all interactive elements are ≥ 44×44px on mobile
   (`touch_target_compliant`).
5. CONTENT PRIORITY: if content must be hidden on small screens, add it to `hidden_on_mobile[]`
   WITH justification. Never hide critical actions on mobile — collapse or reorder instead.

### Reuse ratio (this step is the sole authority)
Compute `reuse_ratio` as a structured object: `{reused, net_new, ratio, verdict}`. `net_new` =
extended + new components. `verdict` is `pass` iff `reused / net_new >= 2`, otherwise `flag`.
The ratio of reused to new components should be > 2:1 for features on existing platforms. If you're
creating more than 30% new components, revisit whether existing ones can be extended instead.

### Design heuristics
Read the finished mapping as ONE screen rather than a list of components, and compare that screen
against the generic-default archetypes. Ask what a user actually sees when these components sit
together: whether anything distinguishes one surface from another, whether hierarchy survives, and
whether the set reads as a deliberate whole. Then append EXACTLY ONE entry to `open_items` — no
more, no fewer — carrying the verdict:
`design-heuristics: <differentiated|default-kept> — <one-line reason naming an archetype>`

## Output
Produces: `ux-ui/components`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  breakpoint_strategy: {mobile: "", tablet: "", desktop: ""}
  reused_components:
    - name: ""
      use_case: ""
  extended_components:
    - name: ""
      changes_needed: []
      use_case: ""
  new_components:
    - name: ""
      reason_existing_insufficient: ""
      props: []
      events: []
      responsive_behavior: ""        # one-liner concept
      accessibility:                 # WCAG 2.1 AA checklist result for this component
        aria_role: ""                # native element or explicit role it presents
        keyboard_interaction: ""     # keys it must support and what each does
        focus_management: ""         # where focus goes on open/close/update
        contrast_token_ref: ""       # token pair its text/UI contrast resolves against
      state_matrix:                  # 9 states; behavior required when applicable
        default: {applicable: true, behavior: ""}
        hover: {applicable: false, behavior: ""}
        focus: {applicable: false, behavior: ""}
        active: {applicable: false, behavior: ""}
        disabled: {applicable: false, behavior: ""}
        empty: {applicable: false, behavior: ""}
        loading: {applicable: false, behavior: ""}
        error: {applicable: false, behavior: ""}
        success: {applicable: false, behavior: ""}
      responsive: {mobile: "", tablet: "", desktop: "", touch_target_compliant: true}
  hidden_on_mobile: []   # [{content, justification}]
  reuse_ratio: {reused: 0, net_new: 0, ratio: "", verdict: "pass|flag"}  # verdict=pass iff reused/net_new >= 2
  open_items: []
```
