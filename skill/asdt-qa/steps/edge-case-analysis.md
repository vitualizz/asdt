# Edge Case Analysis — QA Specialist

## Purpose
Systematically discover edge cases using structured techniques.
Edge cases are not random — they can be derived methodically.

## Inputs
- `qa/ac-list` — arrives as an `### INPUT {topic_key}` block. Extract: `acceptance_criteria[].given/when/then` (to derive boundaries).

## Context budget
qa/ac-list (AC text only, no metadata): max 1,200 tokens.

## Output budget
Max 1,500 tokens. Exceeding the budget is a defect: trim, do not spill.

## Processing
The technique catalogue lives in `skills/edge-case-analysis.md` — load it and apply the
techniques from there. This step decides WHICH techniques apply and how the result is shaped.

1. Walk each AC and select only the techniques its Given/When/Then actually triggers:
   - `boundary` — the AC names a numeric, length, or date range
   - `equivalence` — the AC accepts a class of inputs rather than a single value
   - `state` — the AC moves the system between named states
   - `concurrent` — the AC touches data more than one actor can write
   - `permission` — the AC is gated by a role, ownership, or authentication check
   - `failure` — the AC depends on persistence, an external call, or the network
2. Do NOT emit techniques the AC does not trigger — an empty technique is noise, not coverage.
3. Write one entry per edge case: the technique that produced it, the AC it tests the
   edges of, the scenario, and the precise expected behavior (status code, error message,
   resulting data state) — never "should fail".
4. Assign a priority. Reserve `critical` for edge cases that corrupt data, leak access,
   or lose money; those are the ones `test-case-generation` is guaranteed to cover.

## Output
Produces: `qa/edge-cases`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  edge_cases:
    - id: "EC-001"
      technique: "boundary|equivalence|state|concurrent|permission|failure"
      ac_ref: "AC-001"   # which AC this tests the edges of
      scenario: ""
      expected_behavior: ""
      priority: "critical|high|medium|low"
  critical_count: 0
  open_items: []
```
