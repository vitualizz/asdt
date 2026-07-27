# Spec — Developer Specialist

## Purpose
Define exactly what needs to be built: in-scope, out-of-scope, and acceptance criteria.

## Inputs
- Request: the original feature description
- `developer/dev-exploration`: files to understand, patterns, open questions
- `pm/backlog-entry` (optional — acceptance-criteria authority): the canonical acceptance criteria for this change

Extract from dev-exploration: `open_questions` (answer them here), `patterns_to_follow`.
Extract from pm/backlog-entry: ONLY `acceptance_criteria[]` (not the full artifact — bounds context cost).

**DEGRADATION — `pm/backlog-entry` is optional (only produced when a PM backlog entry exists for this change)**: when it arrives as `### INPUT {project}/{change}/pm/backlog-entry: UNRESOLVED`, author the acceptance criteria from dev-exploration context instead; append "pm/backlog-entry absent — acceptance criteria authored from dev-exploration context" to open_items. Never block on this input.

## Context budget
Request + dev-exploration summary: max 1,500 tokens.

## Processing
1. Answer each `open_question` from the exploration step.
2. Define the scope boundary: what IS included and what is explicitly NOT included.
3. Write acceptance criteria (Given/When/Then format, max 5 criteria). When pm/backlog-entry was
   read, REFINE its `acceptance_criteria[]` into Given/When/Then — pm/backlog-entry is the AC
   AUTHORITY. Do NOT re-derive an independent AC set; preserve the intent of each PM AC. Only when
   pm/backlog-entry is absent do you author ACs from dev-exploration context (note this in `open_items`).
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
