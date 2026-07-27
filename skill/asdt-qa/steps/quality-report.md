# Quality Report — QA Specialist

## Purpose
Produce the final test-plan artifact. Apply the report shared skill to consolidate
test cases and AC validation into a coherent quality document.

## Inputs
- `qa/test-cases` — arrives as an `### INPUT {topic_key}` block. Extract: `total_count`, `critical_count`, and the critical test cases only.
- `qa/ac-gaps` — arrives as an `### INPUT {topic_key}` block. Extract: `gap_count`, `validated_criteria[]` entries whose `status` is not `valid`, and `open_items[]`.

Apply the extraction rules in the report shared skill: from test-cases keep counts + critical cases only.
From ac-gaps keep gap_count, the non-valid criteria, and open_items only.

## Context budget
qa/test-cases summary + qa/ac-gaps summary: max 1,000 tokens.

## Processing
Apply the `report` shared skill:
1. Check: does every testable AC have at least one test case? If not → open_item.
2. Check: are all critical edge cases covered? If not → open_item.
3. Compute coverage: testable ACs covered / total testable ACs.
4. Summarize test distribution across levels (unit/integration/e2e counts).
5. List any AC gaps that need upstream specialist input before testing can proceed.
6. Write a quality verdict: READY / READY WITH CAVEATS / BLOCKED.

## Output
Produces: `qa/test-plan`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  test_summary:
    total_cases: 0
    unit: 0
    integration: 0
    e2e: 0
  ac_coverage:
    total_testable: 0
    covered: 0
    coverage_percent: ""
  ac_gaps: []
  quality_verdict: "READY" # one of: READY|READY_WITH_CAVEATS|BLOCKED
  verdict_rationale: ""
  test_cases:
    - id: ""
      title: ""
      given: ""
      when: ""
      then: ""
      type: "unit" # one of: unit|integration|e2e|acceptance
      story_id: ""
  open_items: []
```
