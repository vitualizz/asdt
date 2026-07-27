# Scope Analysis — PM Specialist

## Purpose
Define explicit boundaries: what IS and IS NOT in scope for this feature.
A backlog entry without explicit out-of-scope items is incomplete — scope ambiguity
is the root cause of most scope creep.

## Inputs
- `pm/user-stories` — Extract: `user_stories` (`id`, `role`, `action`, `priority`, `depends_on`) and the carried `open_items`.

Extract ONLY those fields. Stakeholders are NOT available here — that field lives in
`pm/feature-intake`, which is not a declared input of this step; use each story's `role` as the
actor signal instead, and never invent a stakeholder list.

This declared input arrives ALREADY INJECTED as an `### INPUT {topic_key}` block — consume it
directly and do NOT self-fetch it. If it arrives UNRESOLVED, record the gap in `open_items` and
proceed best-effort.

## Context budget
Injected user-stories fields: max 300 tokens.

## Processing
1. List what IS being built, mapped to user story IDs or capabilities.
2. List what IS NOT being built now — make adjacent capabilities explicit. When in doubt, call it out of scope.
3. Identify integration points: other systems, services, or modules this feature touches, reads from, writes to, or depends on.
4. Flag scope risk items: anything that could cause unplanned expansion (e.g., "if we build X, we will also need Y").

## Output
Produces: `pm/scope-analysis`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  in_scope:
    - ""                   # capabilities or story IDs being built
  out_of_scope:
    - ""                   # explicit list of what is NOT being built in this iteration
  integration_points:
    - system: ""
      nature: ""           # reads-from | writes-to | triggers | depends-on
  risk_flags:
    - ""                   # things that could cause scope expansion if not managed
  open_items: []
```
