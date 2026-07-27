# Technical Handoff — Architect Specialist

## Purpose
Consolidate all architectural work into final artifacts for Developer and QA specialists.
Apply the report shared skill. Surface key constraints the Developer MUST respect.

## Inputs
- `architect/adr` — arrives as an `### INPUT {topic_key}` block. Extract: `decision`, `consequences.negative`, `alternatives_considered[]`.
- `architect/system-design` — arrives as an `### INPUT {topic_key}` block. Extract: `data_model`, `api_surface`, `service_boundaries`, `key_sequence`.
- `architect/risks` — arrives as an `### INPUT {topic_key}` block. Extract: `risks[]` (top 3 by likelihood × impact), `top_risk`.

Apply the extraction rules in the report shared skill: keep what Developer needs, discard architect-internal reasoning.

**DEGRADATION — `architect/system-design` is optional (the `system-design` step runs at `complex` only)**: when it arrives as `### INPUT {project}/{change}/architect/system-design: UNRESOLVED`, derive the data model and API surface at the granularity the ADR already fixes and mark every field the ADR does not settle as unresolved rather than inventing it; append "architect/system-design absent — final design derived from the ADR at ADR granularity" to open_items. Never block on this input.

**DEGRADATION — `architect/risks` is optional (the `risk-analysis` step runs at `complex` only)**: when it arrives as `### INPUT {project}/{change}/architect/risks: UNRESOLVED`, carry the ADR's `consequences.negative` and `technical_debt` entries into the handoff in place of rated risks and do NOT fabricate likelihood/impact ratings; append "architect/risks absent — ADR negative consequences carried in place of rated risks" to open_items. Never block on this input.

## Context budget
All inputs context-extracted to max 300 tokens each = max 900 tokens total.

## Processing
Apply the `report` shared skill:
1. From adr: extract decision + consequences.negative (Developer must know these).
2. From system-design: extract full data_model + api_surface + key_sequence.
3. From risks: extract top 3 risks with their mitigations.
4. MANDATORY: write a "Key Constraints" section — explicit rules Developer must not violate.
5. Consolidate open_items from all inputs.

## Output
Produces: `architect/architectural-decision` (final) and `architect/system-design-final` (final)

This step produces TWO final artifacts. Persist `architectural-decision` under this step's output_topic_key ({project}/{change}/architect/architectural-decision); persist `system-design-final` under its own distinct per-type topic_key {project}/{change}/architect/system-design-final, noted on this step's workflow.yaml entry. Return one payload covering both persisted keys.

`architect/system-design-final` is deliberately a different key from the intermediate `architect/system-design` written by the `system-design` step — never overwrite the intermediate.

architectural-decision schema:
```yaml
payload:
  decision_title: ""
  status: "accepted" # one of: accepted|proposed|superseded|deprecated
  context: ""
  decision: ""
  alternatives_considered:
    - name: ""
      reason_rejected: ""
  consequences:
    positive: []
    negative: []
    technical_debt: []
  key_constraints_for_developer: []
  summary: ""
  open_items: []
```

system-design-final schema:
```yaml
payload:
  data_model:
    - entity: ""
      fields:
        - name: ""
          type: ""
          constraints: ""
      relationships: []
  api_surface:
    - operation: ""
      method: "GET" # one of: GET|POST|PUT|PATCH|DELETE|SUBSCRIBE|PUBLISH
      input: {}
      success_response: {}
      error_cases: []
  service_boundaries:
    touched_modules: []
    new_interfaces: []
    extended_interfaces: []
  key_sequence: []
  open_items: []
```
