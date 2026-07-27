# Divergent Ideation — Researcher Specialist

## Purpose
Frame the problem and generate divergent candidate directions — WITHOUT
converging. This step opens the solution space wide so `feasibility-scan` and
`discovery-brief` have real options to assess and narrow. Deliberately
generative, never selective.

## Inputs
- `inputs: []` — there are no upstream artifacts; the raw problem statement in
  the prompt is the only source. Declared inputs, when a step has them, arrive
  ALREADY INJECTED as `### INPUT {topic_key}` blocks per `parallel-retrieval.md`;
  this step self-fetches nothing.
- Memory context from `context-recall` (injected inline — prior discovery,
  similar problem framings, and decisions enrich the framing but never constrain
  it).

## Context budget
Raw problem statement + injected recall context: max ~600 tokens. Keep the
framing tight so the generative work has room.

## Output budget
`problem_framing`: max 2 sentences. `ideas[]`: max 7 entries; each `what` and
`why` is ONE sentence. Exceeding the budget is a defect: trim, do not spill.

## Processing
1. Restate the problem as a `problem_framing` — one or two sentences capturing
   the pain or opportunity in neutral terms, NOT a solution.
2. Generate freely, then KEEP the strongest three to seven `Idea` entries. Each
   idea is a distinct direction, not a variation in wording. The ceiling is a
   downstream contract, not a style preference: `feasibility-scan` must emit
   exactly one feasibility per idea inside its own input budget, so an unbounded
   idea set makes that contract unsatisfiable.
3. If fewer than three distinct ideas survive, keep the ones that do and record
   the shortfall in `open_items` — never pad the list to reach the floor.
4. Do NOT converge. Do NOT rank, score, or pick a favourite — that is the job of
   `feasibility-scan` and `discovery-brief`. Mixing convergence in here collapses
   the space prematurely.
5. If the problem description is empty or too thin to frame, proceed best-effort
   on the request as stated and record the thinness of the framing in
   `open_items`. This step runs as a sub-agent with no channel to the user, so
   there is no one to ask — but never fabricate a problem to ideate against.
6. Self-check BEFORE persisting: every `Idea.id` is present and non-empty,
   snake_case, and unique within `ideas[]`; the count is between three and seven
   (or the shortfall is already recorded in `open_items`). Fix any violation
   before saving — a duplicate or missing id breaks both downstream steps.

Each `Idea.id` is a **stable, slug-cased (snake_case) identifier** that names the
direction (e.g. `idea_inline_clarify`) and is unique within this run.
`feasibility-scan` and `discovery-brief` reference ideas by this id as a foreign
key, so it must not change or repeat.

## Output
Produces: `researcher/ideation`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  problem_framing: ""         # neutral one/two-sentence statement of the pain or opportunity
  ideas:                      # 3–7 entries; fewer only with a shortfall noted in open_items
    - id: ""                  # stable, unique, snake_case slug — the FK used by feasibility-scan and discovery-brief
      what: ""                # the candidate direction, one sentence
      why: ""                 # the rationale linking it to problem_framing
      theme: ""               # optional — a grouping label when several ideas share a theme
  open_items: []
```
