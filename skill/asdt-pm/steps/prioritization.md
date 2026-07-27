# Prioritization — PM Specialist

## Purpose
Order user stories into a delivery sequence that maximizes value while respecting
dependencies and managing risk. This is not just a re-sort of MoSCoW — it produces
a concrete ordered list with rationale.

## Inputs
- `pm/user-stories` — Extract: `user_stories` (`id`, `priority`, `size`, `depends_on`).
- `pm/scope-analysis` — Extract: `risk_flags`.

Both declared inputs arrive ALREADY INJECTED as `### INPUT {topic_key}` blocks — consume them
directly and do NOT self-fetch either one. If one arrives UNRESOLVED, record the gap in
`open_items` and proceed on the other.

## Context budget
Injected user-stories + scope-analysis fields: max 400 tokens combined.

## Processing
1. Start from MoSCoW priority assigned in the user-stories step.
2. Apply dependency ordering: a story with `depends_on` entries must come after all its dependencies.
3. Apply risk adjustment: stories touching `risk_flags` may need earlier scheduling (de-risk first) or deferral (low confidence). State the reasoning explicitly.
4. Produce a final ordered list. Each entry must include a one-line rationale.
5. Move "Won't" stories to `deferred` with an explicit reason. Deferred ≠ deleted — they are backlog candidates for a future iteration.

## Output
Produces: `pm/prioritization`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  priority_order:
    - id: "US-001"
      rationale: ""        # one line: why this position in the sequence
  deferred:
    - id: "US-00X"
      reason: ""           # why deferred, not just "Won't" — what condition would promote it
  open_items: []
```
