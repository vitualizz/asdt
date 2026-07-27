# Review — QA Specialist

## Purpose
Go/no-go shipping verdict. Evaluate the test-plan, the AC validation results, and the
performance gate holistically and produce a final release-readiness decision before
knowledge is preserved.

## Inputs
- `qa/test-plan` — arrives as an `### INPUT {topic_key}` block. Extract: `quality_verdict`, `ac_coverage.coverage_percent`, `ac_gaps[]`, `open_items[]`.
- `qa/ac-gaps` — arrives as an `### INPUT {topic_key}` block. Extract: `validated_criteria[].id`, `validated_criteria[].status` (`valid` | `needs-revision` | `invalid`), `validated_criteria[].issue`, `gap_count`, `open_items[]`.
- `qa/perf-validation` — arrives as an `### INPUT {topic_key}` block. Extract: `gate_verdict` (`go` | `no-go` | `no-target`), `validations[].dimension` and `validations[].verdict` for any failing target.

**DEGRADATION — `qa/perf-validation` is optional (produced only when the tier or Tailored Workflow ran `performance-validation`)**: when it arrives as `### INPUT {project}/{change}/qa/perf-validation: UNRESOLVED`, decide the verdict from coverage, AC status, and blocking open items alone with NO penalty for the missing gate; append "qa/perf-validation absent — performance gate not run, verdict covers functional readiness only" to open_items. Never block on this input.

## Context budget
test-plan summary + ac-gaps `validated_criteria[]` + perf gate verdict: max 1,500 tokens.

## Output budget
Max 600 tokens. Exceeding the budget is a defect: trim, do not spill.

## Processing

### 1. Coverage gate
Does `test-plan.ac_coverage.coverage_percent` meet the minimum threshold?
- ≥ 80%: pass
- 50–79%: pass-with-notes
- < 50%: fail

### 2. Blocking open items
Scan `test-plan.open_items[]` and `ac-gaps.open_items[]`.
Flag any item that:
- Requires upstream specialist input before testing can proceed
- Indicates a testable AC with no test case

### 3. AC status roll-up
Scan `ac-gaps.validated_criteria[]` and bucket by `status`:
- `invalid` — BLOCKING finding. Record it as a blocking item (`source: ac-gaps`, `item:` the AC id plus its `issue`) and force the no-go path in step 5.
- `needs-revision` — MAJOR finding. Record it as a condition (the AC id plus its `issue`) and cap the verdict at go-with-conditions; it never on its own forces no-go.
- `valid` — no action.

If `validated_criteria[]` is empty but `gap_count` is greater than zero, treat the
discrepancy as a blocking item rather than assuming the ACs are clean.

### 4. Performance gate
Read `perf-validation.gate_verdict`:
- `no-go` — forces this review's verdict to no-go. Name each failing `validations[].dimension` in `blocking_items`.
- `no-target` — no target existed to validate against; cap the verdict at go-with-conditions and add the condition "define NFR targets before shipping".
- `go` — no action.

### 5. Quality verdict alignment
If `test-plan.quality_verdict` is `BLOCKED` → verdict is no-go regardless of the
coverage gate. If `READY` and coverage ≥ 80% with no blockers → go. Intermediate
cases → go-with-conditions with explicit rationale.

### 6. Go/no-go decision
Assign:
- **go** — coverage ≥ 80%, no `invalid` AC, no blocking open items, performance gate `go` or absent
- **go-with-conditions** — `needs-revision` ACs, a `no-target` performance gate, or minor caveats; list every condition that must close before shipping
- **no-go** — one or more `invalid` ACs, a `no-go` performance gate, a `BLOCKED` quality verdict, a blocking open item, or coverage < 50%; name each blocker explicitly

Do NOT invent issues. Only surface defects evident from the artifact content.

## Output
Produces: `qa/qa-review`

Include a `summary` field (≤ 150 tokens) — decision-preservation reads this field.

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  verdict: "go|go-with-conditions|no-go"
  summary: ""     # ≤ 150 tokens — consumed by decision-preservation
  coverage_gate: "pass|pass-with-notes|fail"
  performance_gate: "go|no-go|no-target|not-run"
  blocking_items:
    - source: "test-plan|ac-gaps|perf-validation"
      item: ""
  conditions:     # populated only when verdict is go-with-conditions
    - ""
  open_items: []
```
