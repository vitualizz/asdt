# Design — Developer Specialist

## Purpose
Choose the technical approach and define the data model and API shape.

## Inputs
- `developer/dev-spec`: scope, acceptance criteria, technical requirements

ONLY READ dev-spec. Do NOT load exploration or any other artifact.

## Context budget
dev-spec: max 2,000 tokens.

## Processing
1. Propose the technical approach (1-2 paragraphs, compare 2 options if non-obvious).
   Default to the simplest, most direct approach that satisfies the acceptance criteria, and
   do not add layers, patterns, indirection, or extensibility the requirement does not demand.
   Keep a reversible, two-way-door choice simple; invest day-1 rigor only where the choice is
   hard to reverse or externally observable — where others depend on the surface and changing
   it later breaks them (for example a public API shape, data schema, wire format, or auth
   model; illustrative, not exhaustive). Judge which choices are one-way doors from the
   dev-spec alone. When two options both satisfy the ACs, choose the one with fewer moving
   parts. Justify not only an abstraction you add, but equally a deliberate choice to leave a
   hard-to-reverse or externally-observable surface simple.
2. Define the data model: entities, fields, relationships.
3. Define the API surface: endpoints/functions/interfaces with signatures.
4. Identify any migration or backward-compat concerns.
5. List key technical constraints for implementation (e.g. "use existing AuthMiddleware").

Do NOT write implementation code. Only define the technical structure.

## Output
Produces: `developer/dev-design`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  approach: ""
  data_model:
    - name: ""
      fields: []
  api_surface:
    - name: ""
      signature: ""
      purpose: ""
  migration_notes: []
  key_constraints: []
  open_items: []
```
