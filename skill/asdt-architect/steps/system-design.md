# System Design — Architect Specialist

## Purpose
Define the concrete technical structure: data model, API surface, and service boundaries.

## Inputs
- `architect/adr`: chosen approach, key constraints, consequences

Extract: decision (the chosen approach), consequences.negative (constraints to design around).

## Context budget
architect/adr: max 1,200 tokens.

## Processing
Apply this step's reference skills while designing: `api-design.md` governs the API surface
(resource naming, verb semantics, error shape, pagination), `scalability-analysis.md` governs
bottlenecks and the scaling/caching strategy, `platform-context.md` supplies the naming and
stack conventions, and `scope-definition.md` bounds what this design covers.

1. DATA MODEL: define entities, fields, types, and relationships.
   - Use the naming conventions from platform-summary (loaded via platform-context).
   - Note which fields are indexed, which are nullable, which have constraints.
2. API SURFACE: for each operation, define:
   - Method/endpoint or function signature
   - Input parameters with types
   - Success response shape
   - Error cases (what HTTP codes or error types, and when)
3. SERVICE BOUNDARIES: which existing services/modules does this touch?
   - What new interfaces (if any) need to be defined?
   - What existing interfaces are being extended?
4. SEQUENCE: one key interaction sequence (the happy path) showing how components collaborate.

Quality gate: a system design MUST carry BOTH a data model and an API surface. Before persisting,
verify `data_model` and `api_surface` are each non-empty. If the change genuinely has no persistent
state or no callable surface, say so explicitly in `open_items` ("no data model — change is
stateless") rather than shipping a silently empty array. An empty array with no stated reason is a
defect, not a design.

## Output
Produces: `architect/system-design`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  data_model:
    - entity: ""
      fields:
        - name: ""
          type: ""
          constraints: ""
      relationships: []
  api_surface:
    - operation: ""
      method: ""
      input: {}
      success_response: {}
      error_cases: []
  service_boundaries:
    touched_modules: []
    new_interfaces: []
    extended_interfaces: []
  key_sequence: []    # numbered steps of happy-path interaction
  open_items: []
```
