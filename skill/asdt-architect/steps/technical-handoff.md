# Technical Handoff — Architect Specialist

## Purpose
Consolidate all architectural work into final artifacts for Developer and QA specialists.
Apply the report shared skill. Surface key constraints the Developer MUST respect.

## Inputs
- `architect/adr`: the decision and its consequences
- `architect/system-design`: data model, API surface
- `architect/risks`: top risks and mitigations

Apply the extraction rules in the report shared skill: keep what Developer needs, discard architect-internal reasoning.

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
Produces: `architectural-decision` (final) and `system-design` (final)

Persist `architectural-decision` via mem_save under this step's `output_topic_key` in workflow.yaml; persist the second final artifact `system-design` under its own distinct per-type topic_key (see the NOTE on this step's workflow.yaml entry — do not collide with the intermediate `architect/system-design` produced earlier by the system-design step); return envelope covering both persisted keys.

<!-- ASDT:GENERATED:schema-architectural-decision -->
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
<!-- /ASDT:GENERATED:schema-architectural-decision -->

<!-- ASDT:GENERATED:schema-system-design -->
system-design schema:
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
<!-- /ASDT:GENERATED:schema-system-design -->
