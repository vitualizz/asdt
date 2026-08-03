# Spec — Developer Specialist

## Purpose
Define exactly what needs to be built: in-scope, out-of-scope, and acceptance criteria.

## Inputs
- Request: the original feature description
- `developer/dev-exploration`: files to understand, patterns, open questions
- `pm/handoff` (optional — acceptance-criteria authority): the canonical acceptance criteria for this change

Extract from dev-exploration: `open_questions` (answer them here), `patterns_to_follow`.
Extract from pm/handoff: ONLY `acceptance_criteria[]` (not the full hand-off — bounds context cost).

**DEGRADATION**: if `pm/handoff` is UNRESOLVED, author the acceptance criteria from dev-exploration context and note `ASSUMED:` in open_items.

## Context budget
Request + dev-exploration summary: max 1,500 tokens.

## Processing
1. Answer each `open_question` from the exploration step.
2. Define the scope boundary: what IS included and what is explicitly NOT included.
3. Write acceptance criteria (Given/When/Then format, max 5 criteria). When pm/handoff was
   read, REFINE its `acceptance_criteria[]` into Given/When/Then — pm/handoff is the AC
   AUTHORITY. Do NOT re-derive an independent AC set; preserve the intent of each PM AC. Only when
   pm/handoff is absent do you author ACs from dev-exploration context (note this in `open_items`).
4. List technical requirements (NFRs: performance targets, error handling expectations).

Do NOT design the technical approach. Do NOT write code. Only define what to build.

## Output
Produces: `developer/dev-spec`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  scope:
    in: []
    out: []
  acceptance_criteria:
    - given: ""
      when: ""
      then: ""
  technical_requirements: []
  open_questions_answered: {}  # map of question → answer
  open_items: []
```
