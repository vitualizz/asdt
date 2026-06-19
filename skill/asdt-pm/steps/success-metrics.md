# Success Metrics — PM Specialist

## Purpose
Derive measurable non-functional-requirement (NFR) targets from the user stories.
PM only SETS targets — it names the dimensions worth budgeting and the values to aim
for. It does NOT judge whether any design or implementation meets them; that is the
Architect's and QA's job downstream.

## Inputs
- `pm/user-stories`

## Context budget
Extract: user_stories (action + benefit), priority — focus on stories whose benefit
implies a non-functional need (speed, cost, capacity, footprint, accessibility budget).
Max 400 tokens.

## Processing
Apply the `nfr-budget` shared skill (`../asdt-shared/skills/nfr-budget.md`).
1. Read the user stories and surface their non-functional needs — speed, cost,
   throughput, memory/footprint, bundle size, availability. Not every story implies one.
2. For each distinct NFR need, define one target using the NFR-target value-object:
   pick a `dimension`, set its `budget` (the ceiling) and `target` (the goal, at or
   below budget), and name the `measurement_method`.
3. Leave `verdict` as `n/a` for every target — PM sets, it does not judge.
4. If a story implies a non-functional need but no defensible number exists yet,
   record it in `open_items` rather than inventing a budget.

## Output
Produces: `{project}/{change}/pm/nfr-targets`

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Schema (NFR-target value-object per `../asdt-shared/skills/nfr-budget.md`, `verdict`
narrowed to the PM owner subset):
```yaml
payload:
  nfr_targets:
    - dimension: ""            # e.g. latency_p95 | throughput | cost_per_month | memory_footprint | bundle_size
      budget: ""               # the ceiling/allowance, e.g. "200ms", "$50/mo", "512MB"
      target: ""               # the desired value, at or below budget, e.g. "120ms"
      measurement_method: ""   # HOW it is measured, e.g. "p95 over 1k synthetic requests"
      verdict: "n/a"          # PM only SETS targets — always n/a here
  total_count: 0
  open_items: []               # stories with a non-functional need but no defensible budget yet
```
