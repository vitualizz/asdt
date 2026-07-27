# Performance Validation — QA Specialist

## Purpose
Validate planned/measured performance against the PM's NFR targets and produce a
go/no-go gate verdict. This step is a GATE and is READ-ONLY with respect to upstream
artifacts — it does NOT overwrite `qa/qa-review`; `review` remains the holistic
shipping verdict and may read this perf verdict.

## Inputs
- `qa/test-plan` — arrives as an `### INPUT {topic_key}` block. Extract: `ac_coverage`, `quality_verdict`, and any performance-relevant test cases (produced by `quality-report`).
- `pm/nfr-targets` — arrives as an `### INPUT {topic_key}` block. Extract: each target's `dimension`, `target`, and `measurement_method` (cross-specialist input, produced by PM `success-metrics`).

**DEGRADATION — `pm/nfr-targets` is optional (the PM `success-metrics` step runs at `moderate` and above, so no target exists below that tier)**: when it arrives as `### INPUT {project}/{change}/pm/nfr-targets: UNRESOLVED`, set `gate_verdict` to `no-target` with the explicit note "no target to validate against" and emit no per-dimension `pass` — NEVER emit a silent `go` when there is no target to validate against; append "pm/nfr-targets absent — performance gate returned no-target, nothing to validate against" to open_items. Never block on this input.

## Context budget
test-plan summary + the nfr-targets list: max 1,200 tokens.

## Processing
Apply the `nfr-budget` shared skill (`../asdt-shared/skills/nfr-budget.md`).
1. For each target, compare the planned or measured performance (from the test-plan)
   to that target's `target` value, using its `measurement_method`.
2. Emit verdict `pass` when the value meets the target, `fail` when it misses.
3. Roll up the gate: `gate_verdict` is `go` ONLY if every target passes; otherwise `no-go`.
   When `pm/nfr-targets` is UNRESOLVED the roll-up is `no-target`, never `go` — see the
   DEGRADATION paragraph in Inputs.

## Output
Produces: `qa/perf-validation`

Include a `summary` field (≤ 150 tokens) describing the gate outcome. Each validation
re-states the NFR-target value-object per `../asdt-shared/skills/nfr-budget.md`, with
`verdict` narrowed to the QA owner subset.

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  validations:
    - dimension: ""              # the NFR axis being validated
      target: ""                 # the desired value from pm/nfr-targets
      measured_or_planned: ""    # the planned/measured value from the test-plan
      measurement_method: ""     # HOW it is measured
      verdict: "pass | fail"
      note: ""
  gate_verdict: "go | no-go | no-target"   # no-target (NOT a silent go) when pm/nfr-targets UNRESOLVED
  summary: ""                    # ≤ 150 tokens — gate outcome
  open_items: []                 # records the missing target when UNRESOLVED
```
