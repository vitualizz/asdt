# Feature Brief — UX/UI Specialist

## Purpose
Extract the core user problem, primary actor, and success criteria from the feature request.
Establish what "done" looks like from the user's perspective before designing anything.

## Inputs
- Request: the feature description from the user
- `platform-summary`: existing design system, component library, CSS approach

Note: The request and inline-injected `platform-summary` are provided directly by the orchestrator; this is the first generative step and reads no upstream specialist artifact.

Extract from platform-summary: component_library, css_approach.

## Context budget
Request + platform-summary: max 1,100 tokens.

## Processing
1. Identify the PRIMARY ACTOR: who is the main user of this feature?
2. Define the CORE PROBLEM: what pain does this solve? (not the solution — the problem)
3. Establish SUCCESS CRITERIA: 3-5 observable outcomes that mean the feature worked.
4. Note CONSTRAINTS from platform-summary: which design system rules apply?
5. Identify ADJACENT FEATURES: what existing features does this interact with?
6. Capture DESIGN INTENT: the experiential `tone`, guiding `principles`, and a single `north_star`
   statement that anchors later design decisions.
7. Capture JOBS-TO-BE-DONE (`jtbd`): for each, the `actor`, their `motivation`, and the desired
   `outcome`.

Do NOT jump to solutions. Do NOT sketch layouts. Understand the problem first. This is structured
extraction — keep it concise.

## Output
Produces: `ux-ui/feature-brief`

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Schema:
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
  open_items: []
```
