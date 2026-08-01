# Feature Brief — UX/UI Specialist

## Purpose
Extract the core user problem, primary actor, and success criteria from the feature request.
Establish what "done" looks like from the user's perspective before designing anything.

## Inputs
- Request: the feature description from the user
- `platform-summary`: existing design system, component library, CSS approach

The request and the `platform-summary` arrive ALREADY INJECTED by the orchestrator (the latter from
its inline `platform-analysis` step) — consume them directly, never self-fetch them. This is the
first generative step and reads no upstream specialist artifact.

Extract from platform-summary: component_library, css_approach.

## Context budget
Request + platform-summary: max 1,100 tokens.

## Processing
1. Identify the PRIMARY ACTOR: who is the main user of this feature?
2. Define the CORE PROBLEM: what pain does this solve? (not the solution — the problem)
3. Establish SUCCESS CRITERIA: 3-5 observable outcomes that mean the feature worked.
4. Note CONSTRAINTS and ADJACENT FEATURES: which platform-summary design rules apply, and what this
   touches. ENRICH: fill gaps the request leaves silent from platform-summary, never an empty field.
   CHALLENGE: each fill appends one `open_items` entry with the LITERAL prefix `ASSUMED:` (spelled
   only here): `ASSUMED: <field> — <assumption and its platform-summary signal>`.
5. Capture DESIGN INTENT: the experiential `tone`, guiding `principles`, and a single `north_star`
   statement that anchors later design decisions.
6. Capture JOBS-TO-BE-DONE (`jtbd`): for each, the `actor`, their `motivation`, and the desired
   `outcome`.

Do NOT jump to solutions. Do NOT sketch layouts. Understand the problem first. This is structured
extraction — keep it concise.

## Output
Produces: `ux-ui/feature-brief`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  primary_actor: ""
  core_problem: ""
  success_criteria: []
  design_constraints:
    component_library: ""
    css_approach: ""
    existing_adjacent_features: []
  design_intent: {tone: "", principles: [], north_star: ""}
  jtbd: []   # [{actor, motivation, outcome}]
  open_items: []   # ENRICH fills carry the `ASSUMED:` prefix (step 4).
```
