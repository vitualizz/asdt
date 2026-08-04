# Spec — Developer Specialist

## Purpose
Define what gets built and how: scope, acceptance criteria, technical approach, and the files
the implementation is allowed to touch.

## Inputs
All of these arrive ALREADY INJECTED — never self-fetch them.

- Request: the original feature description
- `dev-exploration` — injected from the orchestrator's context as `### INPUT dev-exploration`.
  Extract: `open_questions` (answer them here), `patterns_to_follow`
- `{project}/{change}/pm/handoff` (optional — AC authority). Extract: ONLY `acceptance_criteria[]`
- `{project}/{change}/architect/handoff` (optional). Extract: `decisions`, `constraints`,
  `data_model`, `api_surface` — when it arrived the approach is ALREADY DECIDED, do not re-open it

**DEGRADATION**: if `pm/handoff` is UNRESOLVED, author the acceptance criteria from dev-exploration context and note `ASSUMED:` in open_items. If `architect/handoff` is UNRESOLVED, decide the approach here per step 4.

## Processing

1. Answer each `open_question` from the exploration step.
2. Define the scope boundary: what IS included and what is explicitly NOT included.
3. Write acceptance criteria (Given/When/Then format, max 5 criteria). When pm/handoff was
   read, REFINE its `acceptance_criteria[]` into Given/When/Then — pm/handoff is the AC
   AUTHORITY. Do NOT re-derive an independent AC set; preserve the intent of each PM AC. Only when
   pm/handoff is absent do you author ACs from dev-exploration context (note this in `open_items`).
4. **Technical approach — only as deep as this change earns.** For a change whose shape is obvious
   from the exploration, one paragraph naming the approach is the whole of this section. Go further —
   data model entities and fields, API signatures, migration notes — only when the change introduces
   or reshapes them. You judge that here; there is no separate design step gating the call.

   Default to the simplest, most direct approach that satisfies the acceptance criteria, and do not
   add layers, patterns, indirection, or extensibility the requirement does not demand. Keep a
   reversible, two-way-door choice simple; invest day-1 rigor only where the choice is hard to
   reverse or externally observable — where others depend on the surface and changing it later breaks
   them (for example a public API shape, data schema, wire format, or auth model; illustrative, not
   exhaustive). When two options both satisfy the ACs, choose the one with fewer moving parts. Justify
   not only an abstraction you add, but equally a deliberate choice to leave a hard-to-reverse or
   externally-observable surface simple. When `architect/handoff` arrived, this section RESTATES its
   decision at implementation granularity — a genuine conflict is an `open_items` entry.
5. **Declare the edit targets.** List `files_to_create` and `files_to_modify` — real paths,
   relative to the project root. This list is load-bearing: `implement` resolves its
   `allowedEditRoots` from it, and anything outside it will not be written. An empty list is a
   deliberate choice meaning "plan only, write nothing", not an oversight.

Do NOT write implementation code here.

## Output
Produces: `dev-spec` — retained in the orchestrator's context, NOT persisted. This step declares
`output: context`: do NOT call `mem_save`. The orchestrator injects the payload below into `implement`.

```yaml
payload:
  scope:
    in: []
    out: []
  acceptance_criteria:
    - given: ""
      when: ""
      then: ""
  approach: ""                 # at the depth this change earns
  data_model: []               # only when the change introduces or reshapes one
  api_surface: []              # only when the change introduces or reshapes one
  migration_notes: []
  key_constraints: []          # what implementation must respect
  files_to_create: []          # real paths — feeds implement's allowedEditRoots
  files_to_modify: []          # real paths — feeds implement's allowedEditRoots
  open_questions_answered: {}  # question → answer
  open_items: []
```
