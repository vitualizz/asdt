# NFR Budget — Shared Skill

## Purpose

Define what a non-functional-requirement (NFR) budget is and give every specialist
that touches one a single shared SHAPE for it. An NFR budget pairs a measurable
quality dimension (latency, cost, bundle size, memory, throughput) with a ceiling
the design and implementation must stay within. A shared shape matters because the
same target is SET by one specialist (PM) and MEASURED against by others (Architect,
QA): if each invented its own structure, the budget could not flow across specialist
boundaries unchanged.

## When to Use

Referenced (via `reference_skills:` in `workflow.yaml`) by every step that defines or
measures against an NFR budget:

- **PM `success-metrics`** — DEFINES targets from user stories (verdict `n/a`; PM only sets).
- **Architect `cost-estimation`** — MEASURES estimated cost vs. budget (verdict `within-budget` | `over-budget`).
- **QA `performance-validation`** — MEASURES planned/measured performance vs. target as a gate (verdict `pass` | `fail`).

## Protocol

The generic "measure against a budget" sequence — each consuming step adapts it to
its owner role:

1. **Identify the dimension(s)** being budgeted (e.g. `latency_p95`, `cost_per_month`, `bundle_size`).
2. **State the budget** (the ceiling/allowance) and the **target** (the desired value, at or below budget).
3. **Name the measurement_method** — HOW the dimension is measured or verified.
4. **Compare** the measured or estimated value to the budget.
5. **Emit a verdict** drawn ONLY from your role's allowed subset (see schema below).
6. **If no budget is available**, record the gap (in `open_items`) and proceed honestly —
   never assume a pass. A missing budget yields a "no target" outcome, not a silent green.

## Reference value-object schema

The NFR-target value-object is the DRY single source of shape. Each consuming step
re-states this block in its own `## Output` schema, narrowing `verdict` to its owner
subset (schema conformance is documentation-trusted — there is no validation tooling).

```yaml
nfr_target:
  dimension: ""            # the NFR axis being budgeted, e.g. latency_p95 | throughput | cost_per_month | memory_footprint | bundle_size
  budget: ""               # the ceiling/allowance the design & implementation must stay within, e.g. "200ms", "$50/mo", "512MB"
  target: ""               # the desired value to aim for, at or below budget, e.g. "120ms"; may equal budget
  measurement_method: ""   # HOW the dimension is measured/verified, e.g. "p95 over 1k synthetic requests", "monthly cloud bill", "lighthouse CI"
  verdict: ""              # owner-dependent; allowed set is the UNION across owners (each step uses ONLY its subset):
                          #   n/a            — PM sets targets only, never judges
                          #   within-budget  — Architect: estimated value is at or under budget
                          #   over-budget    — Architect: estimated value exceeds budget
                          #   pass           — QA: planned/measured value meets the target
                          #   fail           — QA: planned/measured value misses the target
                          #   no-target      — any measurer: no budget available to compare against (honest gap, NOT a pass)
```

**Verdict ownership (pinned)**: PM `success-metrics` uses `n/a`. Architect
`cost-estimation` uses `within-budget` | `over-budget`. QA `performance-validation`
uses `pass` | `fail` (rolling up to a `go` | `no-go` gate). Any measurer falls back to
`no-target` when the budget input is UNRESOLVED — a gate must NEVER be silently green
without a target to validate against.

## Context Budget

No added input budget; reference text only.

## Output

No artifact produced by this skill. It only enriches context for the step it precedes.
