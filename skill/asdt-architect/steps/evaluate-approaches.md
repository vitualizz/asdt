# Evaluate Approaches — Architect Specialist

## Purpose
Compare 2-3 viable technical approaches for the key architectural decision.
Choose one with explicit reasoning. Document why alternatives were rejected.

## Inputs
- `architect/constraints-analysis` — arrives as an `### INPUT {topic_key}` block; produced by the `load-constraints` step. Extract: `hard_constraints` (limits the approach space), `soft_constraints`, `opportunities` (could favor one approach).

## Context budget
architect/constraints-analysis: max 1,200 tokens.

## Processing
1. IDENTIFY the central architectural question (the one decision that constrains everything else).
2. GENERATE 2-3 candidate approaches that respect hard constraints.
3. FOR EACH approach, evaluate across these dimensions:
   - Complexity: how hard to implement and maintain?
   - Performance: how does it behave under load?
   - Coupling: how tightly does it bind to existing components?
   - Reversibility: how hard to change later?
   - Familiarity: does the team already use this pattern?
4. SCORE each dimension as Low/Medium/High impact.
5. CHOOSE the approach with the best overall tradeoff — not necessarily the "best" in one dimension.
6. STATE the rejected alternatives with one-line reasons (these become the alternatives recorded by the `decision-record` step).

## Output
Produces: `architect/approaches`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  central_question: ""
  approaches:
    - name: ""
      description: ""
      complexity: "Low|Medium|High"
      performance: "Low|Medium|High"
      coupling: "Low|Medium|High"
      reversibility: "Low|Medium|High"
      familiarity: "Low|Medium|High"
      pros: []
      cons: []
  chosen: ""
  chosen_rationale: ""
  rejected:
    - name: ""
      reason: ""
  open_items: []
```
