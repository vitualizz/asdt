# Test Case Generation — QA Specialist

## Purpose
Write structured test cases for the happy path, validated ACs, and critical edge cases.

## Inputs
- `qa/ac-list` — arrives as an `### INPUT {topic_key}` block. Extract: `acceptance_criteria[]` (id, given/when/then, testable).
- `qa/test-strategy` — arrives as an `### INPUT {topic_key}` block. Extract: `unit.what`, `integration.what`, `e2e.what`.
- `qa/edge-cases` — arrives as an `### INPUT {topic_key}` block. Extract: `edge_cases[]` where `priority` is `critical` or `high`.

**DEGRADATION — `qa/test-strategy` is optional (produced only at `moderate` and above)**: when it arrives as `### INPUT {project}/{change}/qa/test-strategy: UNRESOLVED`, assign each test case a level directly from the AC it covers — unit for pure logic and validation rules, integration for anything crossing a process, persistence, or service boundary, e2e for a full user-visible flow; append "qa/test-strategy absent — test levels inferred per acceptance criterion" to open_items. Never block on this input.

**DEGRADATION — `qa/edge-cases` is optional (produced only at `moderate` and above)**: when it arrives as `### INPUT {project}/{change}/qa/edge-cases: UNRESOLVED`, cover only the happy path and the explicit negative path of each acceptance criterion in `qa/ac-list`, and do not invent an edge-case inventory; append "qa/edge-cases absent — coverage limited to AC happy and negative paths" to open_items. Never block on this input.

## Context budget
ac-list + test-strategy summary + critical/high edge cases: max 2,000 tokens.

## Processing
1. For each acceptance criterion in `qa/ac-list` marked `testable`, write at least one
   test case in Given/When/Then format — this is the floor, and it holds whether or not
   the optional inputs arrived.
2. For each item in test-strategy (`unit.what`, `integration.what`, `e2e.what`), write a
   test case and record its level (unit/integration/e2e). When test-strategy is absent,
   derive the level per the degradation rule above.
3. For each critical/high edge case, write a test case covering the edge scenario and
   state precisely what "expected behavior" means (status code, error message, data state).
4. Include the test data setup needed and note any mocks or fixtures required.
5. Cross-reference every test case to the AC it covers via `ac_ref`.

Do NOT write code — write structured test specifications that a developer can implement.

## Output
Produces: `qa/test-cases`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  test_cases:
    - id: "TC-001"
      title: ""
      level: "unit|integration|e2e"
      ac_ref: ""
      given: ""
      when: ""
      then: ""
      setup_required: ""
      mocks_required: []
      priority: "critical|high|medium|low"
  total_count: 0
  critical_count: 0
  open_items: []
```
