# Design Critique — UX/UI Specialist

## Purpose
Single-pass annotation of the component inventory against the derived design tokens for
consistency and accessibility. Derive a deterministic `needs_review` signal. This is a
soft-gate ANNOTATE step — it never hard-blocks the workflow.

## Inputs
- `ux-ui/components`: component inventory (reused, extended, new) + carried `reuse_ratio`
- `ux-ui/design-tokens`: the derived token set

Both arrive injected per the parallel-retrieval contract — consume them directly. Do not fetch
either yourself.

Extract from components: the full inventory and the `reuse_ratio` object (carry it verbatim).
Extract from design-tokens: token names/roles to reference in annotations.

## Context budget
Max 2,400 tokens of input material: up to 400 tokens for the token set, plus up to 8 components at
250 tokens each — a component's `state_matrix`, `responsive`, and `accessibility` blocks dominate its
cost, so 250 tokens is the realistic per-component floor. When the inventory exceeds 8 components,
critique `new` first, then `extended`, then `reused`, and record the uncritiqued remainder in
`open_items`.

## Processing
SINGLE PASS ONLY — never iterate, never re-critique, never hard-block.
1. For each component, annotate consistency and accessibility issues against the token set. Each
   annotation references the relevant `token_ref` and carries a severity (`low|medium|high`).
2. CARRY `reuse_ratio` from the components artifact verbatim. NEVER recompute it — component-mapping
   is the sole authority for reuse_ratio.
3. Findings NOT verifiable from text alone (e.g. exact color-contrast ratios) MUST be marked advisory
   / `requires-render-verification` in `open_items` — do NOT assert pass/fail on them.
4. Derive `needs_review` DETERMINISTICALLY: it is `true` iff
   (any annotation `severity == "high"`) OR (carried `reuse_ratio.verdict == "flag"`).
   Otherwise `false`. No judgment beyond this rule.

## Output
Produces: `ux-ui/design-critique`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  critique_annotations: []   # [{target, issue, severity: low|medium|high, token_ref}]
  reuse_ratio: {reused: 0, net_new: 0, ratio: "", verdict: "pass|flag"}  # CARRIED from components, not recomputed
  needs_review: false        # DERIVED deterministically (see rule)
  open_items: []
```
