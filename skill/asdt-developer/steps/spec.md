# Spec — Developer Specialist

## Purpose
Define exactly what needs to be built: in-scope, out-of-scope, and acceptance criteria.

## Inputs
- Request: the original feature description
- `developer/dev-exploration`: files to understand, patterns, open questions
- `pm/backlog-entry` (SOFT — ADR-019 AC authority): the canonical acceptance criteria for this change

Extract from dev-exploration: `open_questions` (answer them here), `patterns_to_follow`.
Extract from pm/backlog-entry: ONLY `acceptance_criteria[]` (not the full artifact — bounds context cost).

**Soft input (no `pm/backlog-entry`)**: if pm/backlog-entry is ABSENT from Engram, do NOT
hard-fail. Note the absence in `open_items`, author ACs from dev-exploration context as before,
and proceed. When pm/backlog-entry IS read, add its topic_key to the dev-spec envelope `input_refs`.

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

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Schema:
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
