# UX Handoff — UX/UI Specialist

## Purpose
Consolidate all UX work into two final artifacts: ux-brief (for Developer context) and
component-spec (for implementation). Apply the report shared skill.

## Inputs
- `ux-ui/feature-brief`: actor, problem, success criteria, design_intent, jtbd
- `ux-ui/ia`: sections, navigation, content_intent
- `ux-ui/flows`: interaction sequences
- `ux-ui/components`: component inventory + absorbed responsive (`responsive`, `breakpoint_strategy`, `hidden_on_mobile`) + state_matrix
- `ux-ui/design-tokens`: derived token set
- `ux-ui/design-critique`: critique annotations + needs_review (complex-only → may arrive UNRESOLVED on lower tiers; handle gracefully)

Apply the extraction rules in the report shared skill to each: keep only fields relevant to implementation handoff.

## Context budget
All inputs context-extracted to max 200 tokens each = max 1,000 tokens total.

## Processing
Apply the `report` shared skill:
1. From feature-brief: extract actor + success_criteria + design_intent + jtbd.
2. From ia: extract navigation.entry_point + primary_actions + content_intent.
3. From flows: extract happy-path steps and decision_points (not edge cases — those go in QA).
4. From components: full component inventory. Responsive specs are sourced from the components
   artifact's absorbed `responsive` / `breakpoint_strategy` / `hidden_on_mobile` fields (plus each
   component's `state_matrix`) — components is now authoritative for responsive.
5. From design-tokens: carry the token set into `design_tokens_ref`.
6. From design-critique (if resolved): carry `critique_annotations` and `needs_review`. design-critique
   is complex-only — if it arrives UNRESOLVED (lower tier), omit `critique_annotations`/`needs_review`
   and add `"design critique deferred — tier did not include design-critique"` to open_items.
7. Consolidate open_items from ALL inputs into a deduplicated list.

## Output
Produces: `ux-brief` (final) and `component-spec` (final)

This step produces TWO final artifacts — a genuine dual-artifact shape (two fully
separate schema blocks below), structurally identical to architect's
`technical-handoff` (`architectural-decision` + `system-design`) and security's
`hardening-checklist` (`security-findings` + `hardening-checklist`); NOT a
single-cohesive-payload shape like qa's `quality-report`. Confirmed by reading
this section directly, not inferred from the compound-looking step/artifact names
(per the explicit caution forwarded from PR3/PR4: similarly-shaped compound names
have landed on opposite answers — qa's was single-artifact, security's was dual).

Persist `ux-brief` via `mem_save` under this step's `output_topic_key` in
`workflow.yaml` (`{project}/{change}/ux-ui/ux-brief`); persist the second final
artifact `component-spec` under its own distinct per-type topic_key
`{project}/{change}/ux-ui/component-spec` (see the inline YAML comment on this
step's `workflow.yaml` entry — no suffix needed, this name collides with neither
the primary key nor any intermediate artifact produced earlier in this
specialist's chain: `feature-brief`, `design-tokens`, `ia`, `flows`, `components`, `design-critique`).
Return an envelope covering both persisted keys.

<!-- ASDT:GENERATED:schema-ux-brief -->
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
  open_items: []
```
<!-- /ASDT:GENERATED:schema-ux-brief -->

<!-- ASDT:GENERATED:schema-component-spec -->
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
      state_matrix:
        default:
          applicable: false
          behavior: ""
        hover:
          applicable: false
          behavior: ""
        focus:
          applicable: false
          behavior: ""
        active:
          applicable: false
          behavior: ""
        disabled:
          applicable: false
          behavior: ""
        empty:
          applicable: false
          behavior: ""
        loading:
          applicable: false
          behavior: ""
        error:
          applicable: false
          behavior: ""
        success:
          applicable: false
          behavior: ""
      responsive:
        mobile: ""
        tablet: ""
        desktop: ""
        touch_target_compliant: false
  critique_annotations:
    - target: ""
      issue: ""
      severity: "low" # one of: low|medium|high
      token_ref: ""
  needs_review: false
  open_items: []
```
<!-- /ASDT:GENERATED:schema-component-spec -->
