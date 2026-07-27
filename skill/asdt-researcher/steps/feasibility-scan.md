# Feasibility Scan — Researcher Specialist

## Purpose
Assess each ideated direction for feasibility — a green/yellow/red verdict with
supporting evidence and an effort estimate. This is the convergence-enabling
step: it gives `discovery-brief` the grounded signal it needs to recommend one
direction without guessing.

## Inputs
- `researcher/ideation` — arrives ALREADY INJECTED as an `### INPUT {topic_key}`
  block per `parallel-retrieval.md`; never self-fetch it. Extract:
  `problem_framing`, `ideas[]` (`id`, `what`, `why`).
- If this input arrives as an `### INPUT ...: UNRESOLVED` block, record the gap
  in `open_items` and proceed with whatever ideas are available.

## Context budget
`ideation`: max 1,500 tokens. Extract `problem_framing` and the `ideas[]` ids,
`what`, and `why`. Discard anything heavier.

## Output budget
One `Feasibility` per `Idea.id` — the ideation ceiling of seven ideas caps this
list at seven entries. `evidence`: max 2 sentences.
Exceeding the budget is a defect: trim, do not spill.

## Processing
1. Produce EXACTLY one `Feasibility` per `Idea.id` from the ideation artifact —
   no more, no fewer. The set of `feasibilities[].idea_id` must equal the set of
   `ideas[].id`.
2. Each `Feasibility.idea_id` is a foreign key — it MUST match an existing
   `Idea.id`. If you would emit a feasibility whose `idea_id` has no matching
   idea, HALT and record the dangling id in `open_items` rather than inventing an
   idea to match it.
3. Assign a `verdict`: `green` (clearly feasible), `yellow` (feasible with
   caveats or unknowns), `red` (blocked or impractical as framed).
4. Attach `evidence` — the concrete reason for the verdict (a constraint, a
   dependency, a known limitation), not a restatement of the verdict.
5. Attach `evidence_source` — where that reason came from: `platform-summary`
   (the loaded platform context), `recalled-artifact` (memory recalled for this
   change), or `stated-assumption` (your own reasoning, unbacked by either). A
   reader must be able to tell a grounded verdict from an educated guess; when
   the source is `stated-assumption`, say so rather than dressing it up.
6. Estimate `effort`: `low` | `medium` | `high`.

Do not converge here — every idea gets assessed. Choosing the winner is
`discovery-brief`'s job.

## Output
Produces: `researcher/feasibility`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  feasibilities:              # exactly one entry per Idea.id — no more, no fewer
    - idea_id: ""             # FK — MUST match an existing Idea.id from researcher/ideation
      verdict: "green | yellow | red"
      evidence: ""            # the concrete reason for the verdict, not a restatement
      evidence_source: "platform-summary|recalled-artifact|stated-assumption"
      effort: "low | medium | high"
  open_items: []
```
