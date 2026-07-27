# Backlog Entry — PM Specialist

## Purpose
Consolidate all PM artifacts into the final structured backlog entry.
This is the canonical requirements artifact consumed by Architect, Developer, and QA.

## Inputs
- `pm/user-stories` — Extract: `user_stories` (all fields), `total_count`, `must_count`, `open_items`. Always present at every tier that runs this step.
- `pm/scope-analysis` — Extract: `in_scope`, `out_of_scope`, `integration_points`, `risk_flags`.
- `pm/prioritization` — Extract: `priority_order`, `deferred`.
- `pm/nfr-targets` — Extract: each target's `dimension`, `budget`, `target`, and `measurement_method`. Mirror these key names verbatim — the producer's vocabulary is canonical.

All declared inputs arrive ALREADY INJECTED as `### INPUT {topic_key}` blocks — consume them
directly and never self-fetch one.

**DEGRADATION — `pm/scope-analysis` is optional (the `scope-analysis` step runs at `moderate` and above, so it is absent at `simple`)**: when it arrives as `### INPUT pm/scope-analysis: UNRESOLVED`, derive the scope block from the stories alone — set `in_scope` to the capabilities the stories cover, write `out_of_scope` by naming the adjacent capabilities the stories deliberately do NOT cover (this list is mandatory and is never left empty), and leave `integration_points` and `risk_flags` empty rather than guessing at systems you cannot see; append "pm/scope-analysis absent — scope boundaries inferred from user stories only; integration points and risk flags not assessed" to open_items. Never block on this input.

**DEGRADATION — `pm/prioritization` is optional (the `prioritization` step runs at `complex` only, so it is absent at `simple` and `moderate`)**: when it arrives as `### INPUT pm/prioritization: UNRESOLVED`, derive `priority_order` from the stories' own MoSCoW `priority` plus `depends_on` ordering (must → should → could, dependencies first) and put every `wont` story in `deferred` with its MoSCoW rating as the reason; append "pm/prioritization absent — delivery order derived from MoSCoW and dependencies, without value/effort/risk rationale" to open_items. Never block on this input.

**DEGRADATION — `pm/nfr-targets` is optional (the `success-metrics` step runs at `moderate` and above, so it is absent at `simple`)**: when it arrives as `### INPUT pm/nfr-targets: UNRESOLVED`, emit `nfr_targets: []` and never invent budgets or measurement methods from the stories; append "pm/nfr-targets absent — no non-functional budgets set for this change" to open_items. Never block on this input.

## Context budget
Load every resolved input artifact in full. Max 2,000 tokens total across all four.

## Processing
1. Merge user stories with their position from the prioritization artifact.
2. Attach the full scope block (in/out of scope, integration points, risk flags).
3. Mirror the non-functional targets into `nfr_targets`, carrying `dimension`, `budget`, `target`, and `measurement_method` verbatim under those exact key names — never rename them, so a reader moving between this artifact and `pm/nfr-targets` sees one vocabulary. Do NOT re-derive, re-budget, or re-judge them — `pm/nfr-targets` remains the authoritative artifact for Architect and QA; this block is a mirror so a reader of the backlog entry alone sees the targets.
4. Write an `executive_summary`: one paragraph — what this feature is, why it matters, and what it explicitly excludes. This is the human-readable entry point for any specialist picking up this artifact.
5. Carry forward any unresolved `open_items` from prior steps, plus every degradation entry recorded above — downstream specialists should address these.
6. Write a `summary` field (≤ 150 tokens) for decision-preservation.

## Output
Produces: `pm/backlog-entry`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  feature_name: ""
  summary: ""              # ≤ 150 tokens — consumed by decision-preservation
  executive_summary: ""    # 1 paragraph: what, why, what it explicitly excludes
  user_stories:
    - id: ""
      role: ""
      action: ""
      benefit: ""
      priority: "must | should | could | wont"
      size: "small | medium | large"
      acceptance_criteria: []   # high-level plain English — QA formalizes these
      depends_on: []
  scope:
    in_scope: []
    out_of_scope: []
    integration_points:
      - system: ""
        nature: ""
    risk_flags: []
  nfr_targets:               # mirrored from pm/nfr-targets — [] when success-metrics did not run
    - dimension: ""          # the NFR dimension, e.g. latency_p95 | cost_per_month | bundle_size
      budget: ""             # the ceiling/allowance
      target: ""             # the goal, at or below budget
      measurement_method: "" # HOW it is measured
  priority_order: []       # ordered list of US IDs (from prioritization)
  deferred: []             # US IDs explicitly deferred with reasons
  open_items: []           # unresolved questions and degradations for downstream specialists
```

## Downstream consumption
- **Architect**: reads `executive_summary`, `scope` (integration_points, risk_flags), `nfr_targets`
- **Developer**: reads `user_stories`, `priority_order`, `acceptance_criteria`
- **QA**: reads `user_stories` + `acceptance_criteria` as the primary requirements source — this replaces the raw request fallback in `load-requirements`
