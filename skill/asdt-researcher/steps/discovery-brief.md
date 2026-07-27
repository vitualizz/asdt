# Discovery Brief — Researcher Specialist

## Purpose
Converge the explored space to ONE recommended direction for PM handoff. This
step turns the divergent ideas of `divergent-ideation` plus the verdicts of
`feasibility-scan` into a single clear recommendation, plus the won't-do
candidates that seed PM's out-of-scope list. The brief's `summary` is the prose
the orchestrator hands to `/asdt-pm`.

## Inputs
- `researcher/ideation` — arrives ALREADY INJECTED as an `### INPUT {topic_key}`
  block per `parallel-retrieval.md`; never self-fetch it. Extract:
  `problem_framing`, `ideas[]` (`id`, `what`, `why`).
- `researcher/feasibility` — arrives ALREADY INJECTED the same way. Extract:
  `feasibilities[]` (`idea_id`, `verdict`, `evidence`, `evidence_source`,
  `effort`).
- If EITHER input arrives as an `### INPUT ...: UNRESOLVED` block, record the gap
  in `open_items` and degrade gracefully — recommend from whatever survives,
  noting the reduced confidence.

## Context budget
Combined `ideation` + `feasibility`: max 2,500 tokens. Apply the `report.md`
extraction rules — pull `problem_framing`, the surviving ideas, and their
verdicts; discard the rest.

## Output budget
`options[]` and `feasibility_notes[]`: max 7 entries each, one note per option.
`wont_candidates[]`: max 7 entries. Each `rationale` and `note` is at most 2
sentences; `summary` is at most 150 tokens.
Exceeding the budget is a defect: trim, do not spill.

## Processing
1. VERIFICATION GATE — run this BEFORE anything else. If BOTH `researcher/ideation`
   and `researcher/feasibility` arrived UNRESOLVED, do NOT fabricate a
   recommendation: persist a degraded brief with `recommended_direction: ""`,
   name both missing inputs in `open_items`, and say plainly in `summary` that no
   direction could be recommended because both upstream artifacts were missing.
   Stop after persisting that degraded brief.
2. Carry the `context` from `ideation.problem_framing`.
3. Build `options` — the surviving directions worth presenting (drop `red`
   ideas unless none survive). Each option keeps its `idea_id` so the brief stays
   traceable to ideation.
4. Write `feasibility_notes` — a condensed green/yellow/red rationale per
   surviving option. **MANDATORY**: this is what makes the recommendation
   defensible. A brief without feasibility notes is incomplete — the step
   consumes `feasibility` precisely so this section can exist.
5. Choose exactly ONE `recommended_direction` — prose that restates a single
   entry from `options[]`, so the orchestrator can hand it to PM as-is. The brief
   always recommends one — never zero (except under the degraded path in the
   verification gate), never a tie.
6. Collect `wont_candidates` — the directions explicitly NOT recommended. These
   seed PM's `scope.out_of_scope` at handoff.
7. Write a `summary` (≤ 150 tokens) — the prose the orchestrator renders into the
   `/asdt-pm` RAW REQUEST. This is the load-bearing handoff text; make it
   self-contained.
8. Self-check BEFORE persisting — every one of these must hold, or fix the brief
   and re-check:
   - `recommended_direction` restates one of the listed `options[]` — the option
     it names is identifiable, and no direction outside `options[]` is proposed;
   - exactly ONE direction is recommended, never a tie;
   - no `wont_candidates[].idea_id` also appears in `options[]`;
   - every option in `options[]` has a `feasibility_notes[]` entry with the same
     `idea_id`, and every note maps back to a listed option.

## Output
Produces: `researcher/discovery-brief`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  context: ""                 # the problem_framing carried from ideation
  options:                    # surviving directions considered
    - idea_id: ""             # FK to Idea.id from researcher/ideation
      direction: ""           # the option in one sentence
      rationale: ""           # why it is worth presenting
  feasibility_notes:          # MANDATORY — exactly one note per entry in options[]
    - idea_id: ""             # FK — MUST match an idea_id present in options[]
      verdict: "green | yellow | red"
      note: ""                # condensed rationale carried from researcher/feasibility
  recommended_direction: ""   # prose restating exactly ONE entry from options[]; "" only on the degraded path
  wont_candidates:            # directions not recommended — seed PM scope.out_of_scope
    - idea_id: ""             # FK to Idea.id — MUST NOT also appear in options[]
      reason: ""              # why it is set aside
  summary: ""                 # ≤150 tokens — prose handed to /asdt-pm as the raw request
  open_items: []
```

## Downstream consumption
- The orchestrator renders `summary` + `recommended_direction` to prose and
  passes it as the RAW REQUEST to `/asdt-pm`.
- `pm/feature-intake` keeps `inputs: []` — UNCHANGED. The handoff is prose-only;
  the Researcher introduces no declared-input contract on PM.
- The belt is SOFT: a prior `researcher/discovery-brief` MAY be recalled via
  `knowledge-recall` for richer context, but it is never a required input.
