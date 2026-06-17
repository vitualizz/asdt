# Review — QA Specialist

## Purpose
Go/no-go shipping verdict. Evaluate the test-plan and AC gaps holistically and
produce a final release-readiness decision before knowledge is preserved.

This step is READ-ONLY — it never writes host source files.

## Inputs
- `qa/test-plan`: consolidated test plan with quality_verdict, coverage stats, AC gaps
- `qa/ac-gaps`: AC validation results (gaps, missing criteria, open items)

Extract from test-plan: `quality_verdict`, `ac_coverage`, `ac_gaps[]`, `open_items[]`.
Extract from ac-gaps: `gaps[]`, `severity` per gap.

## Context budget
test-plan summary + ac-gaps: max 1,500 tokens.

## Processing

### 1. Coverage gate
Does `ac_coverage.coverage_percent` meet the minimum threshold?
- ≥ 80%: pass
- 50–79%: pass-with-notes
- < 50%: fail

### 2. Blocking open items
Scan `test-plan.open_items[]` and `ac-gaps.open_items[]`.
Flag any item that:
- Requires upstream specialist input before testing can proceed
- Indicates a testable AC with no test case

### 3. AC gap severity
Scan `ac-gaps.gaps[]`. If any gap is marked `critical` or `blocking`:
verdict → fail. Surface the gap explicitly.

### 4. Quality verdict alignment
If `test-plan.quality_verdict` is `BLOCKED` → verdict is fail regardless of
coverage gate. If `READY` and coverage ≥ 80% with no blockers → pass.
Intermediate cases → pass-with-notes with explicit rationale.

### 5. Go/no-go decision
Assign:
- **go** — coverage ≥ 80%, no blocking gaps, no blocking open items
- **go-with-conditions** — minor gaps or caveats; list the conditions that must close before shipping
- **no-go** — one or more blocking items, critical AC gaps, or coverage < 50%; name each blocker explicitly

Do NOT invent issues. Only surface defects evident from the artifact content.

## Output
Produces: `qa/qa-review`

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Include a `summary` field (≤ 150 tokens) — decision-preservation reads this field.

Schema:
```yaml
payload:
  verdict: "go|go-with-conditions|no-go"
  summary: ""     # ≤ 150 tokens — consumed by decision-preservation
  coverage_gate: "pass|pass-with-notes|fail"
  blocking_items:
    - source: "test-plan|ac-gaps"
      item: ""
  conditions:     # populated only when verdict is go-with-conditions
    - ""
  open_items: []
```
