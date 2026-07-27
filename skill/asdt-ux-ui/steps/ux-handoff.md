# UX Handoff — UX/UI Specialist

## Purpose
Consolidate all UX work into two final artifacts: ux-brief (for Developer context) and
component-spec (for implementation). Apply the report shared skill.

## Inputs
- `ux-ui/feature-brief`: actor, problem, success criteria, design_intent, jtbd
- `ux-ui/ia`: sections, navigation, content_intent
- `ux-ui/flows`: interaction sequences
- `ux-ui/components`: component inventory + responsive (`responsive`, `breakpoint_strategy`, `hidden_on_mobile`) + `state_matrix` + `accessibility`
- `ux-ui/design-tokens`: derived token set
- `ux-ui/design-critique`: critique annotations + needs_review
- `ux-ui/content-inventory`: text touchpoints (TextTouchpoint list)

All declared inputs arrive ALREADY INJECTED as `### INPUT {topic_key}` blocks per the
parallel-retrieval contract — consume them directly, never self-fetch them.

Apply the extraction rules in the report shared skill to each: keep only fields relevant to implementation handoff.

**DEGRADATION — `ux-ui/design-critique` is optional (produced only at the complex tier)**: when it
arrives as `### INPUT ux-ui/design-critique: UNRESOLVED`, omit `critique_annotations` and emit
`needs_review: "not-evaluated"`; append "design critique absent — component spec assembled without
critique annotations and never evaluated for review" to open_items. Never block on this input.

**DEGRADATION — `ux-ui/content-inventory` is optional (produced only at the moderate and complex
tiers)**: when it arrives as `### INPUT ux-ui/content-inventory: UNRESOLVED`, assemble
`content_intent` from the IA thin intent alone (BRANCH B in Processing step 2); append "content
inventory absent — content_intent degraded to the IA thin intent, microcopy left empty" to
open_items. Never block on this input.

## Context budget
Each declared input is context-extracted to max 200 tokens; with seven declared inputs that is max
1,400 tokens of input material total.

## Processing
Apply the `report` shared skill:
1. From feature-brief: extract actor + success_criteria + design_intent + jtbd.
2. From ia + content-inventory: extract IA `navigation.entry_point` + `primary_actions`, then
   ASSEMBLE `content_intent` via this SOURCE-SELECTION algorithm. The emitted `content_intent`
   block is UNCONDITIONALLY the 4 keys `{copy_direction, microcopy, empty_state, error_state}`
   (R-004 — feed VALUES only; never add or remove a key). `surface_type → content_intent` key map:
   `CTA → microcopy[]`, `label → microcopy[]`, `empty_state → empty_state`, `error → error_state`;
   `copy_direction` is always sourced from the IA thin intent.
   - BRANCH A (content-inventory PRESENT): build `microcopy[]` from touchpoints whose `surface_type`
     ∈ {CTA, label} (use `representative_copy`); `empty_state` from the empty_state touchpoint(s);
     `error_state` from the error touchpoint(s); `copy_direction` from IA.
   - BRANCH B (content-inventory UNRESOLVED/absent): degrade to the IA thin intent — `copy_direction`
     from IA; `microcopy: []`; `empty_state`/`error_state` from the IA intent hints; APPEND the
     open_item named in the content-inventory DEGRADATION paragraph above.
   - BOTH ABSENT (no content-inventory AND no IA content_intent): `copy_direction: ""`, empty values,
     and APPEND the same open_item.
3. From flows: extract happy-path steps and decision_points (not edge cases — those go in QA).
4. From components: full component inventory. Responsive specs are sourced from the components
   artifact's `responsive` / `breakpoint_strategy` / `hidden_on_mobile` fields (plus each
   component's `state_matrix` and `accessibility` block) — components is authoritative for
   responsive, state, and accessibility. Carry those blocks through verbatim.
5. From design-tokens: carry the token set into `design_tokens_ref`.
6. From design-critique: carry `critique_annotations` verbatim, and set `needs_review` to the STRING
   form of the critique's boolean — `"true"` or `"false"`. When design-critique is UNRESOLVED, follow
   its DEGRADATION paragraph above and emit `needs_review: "not-evaluated"`. The tri-state exists so a
   consumer can tell "critiqued and clean" (`"false"`) from "never critiqued" (`"not-evaluated"`);
   never collapse the two.
7. Consolidate open_items from ALL inputs into a deduplicated list.

## Output
Produces: `ux-brief` (final) and `component-spec` (final)

This step produces TWO final artifacts. Persist `ux-brief` under this step's output_topic_key ({project}/{change}/ux-ui/ux-brief); persist `component-spec` under its own distinct per-type topic_key {project}/{change}/ux-ui/component-spec, noted on this step's workflow.yaml entry. Return one payload covering both persisted keys.

ux-brief schema:
```yaml
payload:
  feature_summary: ""
  primary_actor: ""
  success_criteria: []
  design_intent:
    tone: ""
    principles: []
    north_star: ""
  jtbd: []
  content_intent:
    copy_direction: ""
    microcopy: []
    empty_state: ""
    error_state: ""
  user_flows:
    - id: ""
      name: ""
      steps: []
      decision_points: []
  information_architecture:
    sections: []
    navigation_path: ""
  summary: ""    # ≤ 150 tokens — read by the decision-preservation inline step
  open_items: []
```

component-spec schema:
```yaml
payload:
  reused_components: []
  extended_components: []
  new_components:
    - name: ""
      purpose: ""
      reason_existing_insufficient: ""
      props: []
      events: []
      responsive_behavior: ""
      design_tokens_ref: []
      state_matrix: {}     # carried verbatim from ux-ui/components — the 9-state shape defined by the component-mapping step
      accessibility: {}    # carried verbatim from ux-ui/components — {aria_role, keyboard_interaction, focus_management, contrast_token_ref}
      responsive:
        mobile: ""
        tablet: ""
        desktop: ""
        touch_target_compliant: false
  critique_annotations:    # omitted entirely when design-critique did not run
    - target: ""
      issue: ""
      severity: "low" # one of: low|medium|high
      token_ref: ""
  needs_review: "not-evaluated"   # tri-state STRING: "true" | "false" | "not-evaluated" (critique never ran)
  open_items: []
```
