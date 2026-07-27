# Cost Estimation — Architect Specialist

## Purpose
Estimate the cost profile of the system design per relevant NFR dimension and judge
each estimate against the PM's budget. The Architect MEASURES (estimates) and emits a
`within-budget` | `over-budget` verdict — it does not redefine the targets.

## Inputs
- `architect/system-design` — arrives as an `### INPUT {topic_key}` block; produced by the `system-design` step. Extract: `service_boundaries`, `data_model`, `api_surface`, and any scalability notes.
- `pm/nfr-targets` — arrives as an `### INPUT {topic_key}` block; produced by PM `success-metrics`. Extract: each target's dimension and `budget`.

**DEGRADATION — `pm/nfr-targets` is optional (PM `success-metrics` runs at `moderate` and above, so no budget exists below that tier)**: when it arrives as `### INPUT {project}/{change}/pm/nfr-targets: UNRESOLVED`, estimate every cost profile anyway WITHOUT a budget — leave `budget` empty and record the note "no budget to compare against" in place of a verdict, never inventing a budget and never silently claiming within-budget; append "pm/nfr-targets absent — cost profiles estimated with no budget to compare against" to open_items. Never block on this input.

## Context budget
system-design boundaries + API surface + the nfr-targets list: max 1,200 tokens.

## Processing
Apply the `nfr-budget` shared skill (`../asdt-shared/skills/nfr-budget.md`).
1. For each design approach/component, identify which NFR dimensions it drives
   (e.g. a queue → throughput + cost; a cache → latency + memory_footprint).
2. Estimate the cost/value for each relevant dimension, naming the basis of the
   estimate (`measurement_method`).
3. Compare each estimate to that target's `budget`; emit verdict `within-budget`
   when the estimate is at or under budget, `over-budget` when it exceeds it. When the
   budget is absent, apply the DEGRADATION rule above instead of assuming a verdict.

## Output
Produces: `architect/cost-estimate`

Each profile re-states the NFR-target value-object per `../asdt-shared/skills/nfr-budget.md`,
with `verdict` narrowed to the Architect owner subset.

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  cost_profiles:
    - approach: ""             # design approach / component being estimated
      nfr_target_ref: ""       # the dimension this profile estimates against
      estimated_cost: ""       # the estimated value for the dimension
      budget: ""               # the budget from pm/nfr-targets ("" when UNRESOLVED)
      measurement_method: ""   # basis of the estimate
      verdict: "within-budget | over-budget"   # when budget UNRESOLVED, note "no budget to compare against" instead
      rationale: ""
  open_items: []               # MUST carry the UNRESOLVED pm/nfr-targets note when the budget input is absent
```
