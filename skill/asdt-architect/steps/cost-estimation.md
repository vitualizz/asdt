# Cost Estimation — Architect Specialist

## Purpose
Estimate the cost profile of the system design per relevant NFR dimension and judge
each estimate against the PM's budget. The Architect MEASURES (estimates) and emits a
`within-budget` | `over-budget` verdict — it does not redefine the targets.

## Inputs
- `architect/system-design`: extract service_boundaries, data_model, api_surface, scalability notes.
- `pm/nfr-targets`: the NFR targets to estimate against (HARD cross-specialist InputRef — produced by PM `success-metrics`).

## Context budget
system-design boundaries + API surface + the nfr-targets list: max 1,200 tokens.

## Processing
Apply the `nfr-budget` shared skill (`../asdt-shared/skills/nfr-budget.md`).
1. For each design approach/component, identify which NFR dimensions it drives
   (e.g. a queue → throughput + cost; a cache → latency + memory_footprint).
2. Estimate the cost/value for each relevant dimension, naming the basis of the
   estimate (`measurement_method`).
3. Compare each estimate to that target's `budget`; emit verdict `within-budget`
   when the estimate is at or under budget, `over-budget` when it exceeds it.
4. **DEGRADATION** — if `pm/nfr-targets` arrived as an `### INPUT
   {project}/{change}/pm/nfr-targets: UNRESOLVED` block (PM `success-metrics` did not
   run): estimate the cost profile anyway WITHOUT a budget. Do NOT assume a verdict —
   record each estimate with the note "no budget to compare against", record the gap in
   `open_items`, and PROCEED. Never invent a budget and never silently claim within-budget.

## Output
Produces: `{project}/{change}/architect/cost-estimate`

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Schema (each profile re-states the NFR-target value-object per
`../asdt-shared/skills/nfr-budget.md`, `verdict` narrowed to the Architect owner subset):
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
