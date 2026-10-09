# Spec — Developer Specialist

## Purpose
Define what gets built and how: scope, acceptance criteria, technical approach, and the files
the implementation is allowed to touch.

## Inputs
All of these arrive ALREADY INJECTED — never self-fetch them.

- Request: the original feature description
- Platform summary — injected inline. Extract the `Design surface:` line's bare value (strip
  any parenthesised marker; `none — no visual surface` reads as `none`)
- `dev-exploration` (OPTIONAL) — injected from the orchestrator's context as
  `### INPUT dev-exploration`. Extract: `open_questions` (answer them here), `patterns_to_follow`
- `{project}/{change}/pm/handoff` (optional — AC authority). Extract: ONLY `acceptance_criteria[]`
- `{project}/{change}/architect/handoff` (optional). Extract: `decisions`, `constraints`,
  `data_model`, `api_surface`, `files_hint` — when it arrived the approach is ALREADY DECIDED, do
  not re-open it
- `{project}/{change}/ux-ui/handoff` (optional). Extract: `flows`, `components`, `screens`,
  `component_designs`, `a11y_requirements`, `surface`, `visual_direction` — `flows` (with their
  branches, states, and `copy`), `screens`, `component_designs`, and `a11y_requirements` are
  carried verbatim, because they are the design the UI code must implement; every
  `component_designs` entry — each component `gap` has one — is a file to create or modify, or
  an `open_items` entry when its location is unclear; `surface` becomes a `key_constraints` entry; a
  `visual_direction` is carried verbatim, and the files that implement it as the project's
  first tokens join the edit targets; each `open_items` entry starting `Proposed token:` is
  carried verbatim into `key_constraints`, and the token file it names joins `files_to_modify`
- `{project}/{change}/security/handoff` (optional). Extract: `risks[].mitigation` and
  `constraints` (the hardening checklist) — each mitigation that lands in this change's files is
  a `key_constraints` entry
- `{project}/{change}/developer/handoff` (OPTIONAL) — the record this plan will replace under
  the same key. Extract only `files_changed` and `verification`, and only when it carries them:
  they are a delivery — this record's, or one an earlier plan already carried forward
- `spec-feedback` (OPTIONAL) — injected as `### INPUT spec-feedback` only when the human chose
  adjust at the approval gate. Extract: `spec` (the plan they were shown) and `feedback` (their
  words, verbatim)

**DEGRADATION**: if `pm/handoff` is UNRESOLVED, author the acceptance criteria from dev-exploration context and note `ASSUMED:` in open_items. If `architect/handoff` is UNRESOLVED, decide the approach here per step 4. If `ux-ui/handoff`, `security/handoff`, or `developer/handoff` is UNRESOLVED, proceed without it — no `open_items` entry needed. If `spec-feedback` is absent, this is a first pass — proceed normally. If `dev-exploration` is absent (an adjust on a resumed run), revise the spec inside `spec-feedback` instead of re-deriving one. If the `Design surface:` line is missing, skip the surface check in step 5 — no entry.

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
   relative to the project root. Fold in every path of the Architect's `files_hint` this change
   creates or modifies — the design lands nowhere this list does not name. This list is
   load-bearing: `implement` resolves its `allowedEditRoots` from it and nothing else, and
   anything outside it will not be written. An empty list is a deliberate choice meaning "plan
   only, write nothing", not an oversight. Under `Design surface: none`, each UI target — a view,
   a component, a stylesheet or token file — gets one `open_items` entry, `surface conflict:
   {path} is UI, but the project declares no visual surface`, so the human sees it at approval.
6. **On a revision** (`spec-feedback` arrived), steps 1–5 start from the spec it carries, not
   from scratch: the human's words OVERRIDE it wherever the two disagree, and everything they did
   not touch stays as it was. Step 5 still runs in full — an adjustment to scope moves the edit
   targets with it.

Do NOT write implementation code here.

**A new iteration never loses the delivery.** When `developer/handoff` carried `files_changed` or
`verification`, copy them verbatim into this plan: it replaces that record under the same key,
and until the human approves building over it, that delivery is still what the record reports.

## Output
Produces: `dev-spec` — the plan, persisted AND handed forward (`asdt-core/protocol.md` §1). Persist
it via `mem_save` under this step's `output_topic_key`, with `stage: spec`, so a plan-only or
interrupted run can be resumed in a later session; then return the same payload — the orchestrator
shows it to the human for approval and injects it into `implement` as `### INPUT dev-spec`.

```yaml
payload:
  what: ""                     # the change this plan delivers, one sentence
  stage: spec
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
  ux_flows: []                 # verbatim `flows` from ux-ui/handoff — branches, states, copy
  screens: []                  # verbatim from ux-ui/handoff — layout, hierarchy, type, states per screen
  component_designs: []        # verbatim from ux-ui/handoff — anatomy, variants, states, tokens
  a11y_requirements: []        # verbatim from ux-ui/handoff
  visual_direction: {}         # verbatim from ux-ui/handoff, only when it carried one
  files_to_create: []          # real paths — feeds implement's allowedEditRoots
  files_to_modify: []          # real paths — feeds implement's allowedEditRoots
  files_changed: []            # verbatim from developer/handoff — the prior delivery, never this plan
  verification: {}             # verbatim from developer/handoff
  open_questions_answered: {}  # question → answer
  open_items: []
```
