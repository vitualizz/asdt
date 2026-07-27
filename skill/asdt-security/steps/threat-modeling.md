# Threat Modeling — Security Specialist

## Purpose
Identify threats using the STRIDE methodology. STRIDE is systematic —
it catches threats that intuition misses.

## Inputs
This step declares no artifact inputs (`inputs: []` in workflow.yaml) — it runs on whatever
context the orchestrator injects. Whatever DOES arrive comes ALREADY INJECTED as an
`### INPUT {topic_key}` block; never self-fetch.

- `platform-summary`: tech stack, component library (to know attack surface shape). Extract: stack entries, framework names, storage engines.
- Any upstream artifact the orchestrator happened to inject (`architect/system-design`, `architect/architectural-decision`, `developer/dev-implementation`). Extract only: API surface entries, service boundaries, data model entities.

If no upstream artifact arrives: derive the model from the platform summary and the request,
and append "upstream design artifacts absent — threat model derived from platform summary and
request text" to open_items.

## Context budget
platform-summary + upstream summaries: max 1,500 tokens.

## Processing
The STRIDE catalogue — the six categories with their detection patterns, key questions,
trust-boundary list, and the text notations for threat trees and data flows — lives in
`skills/threat-modeling.md`, this step's reference skill. Work from it; it is not restated here.

Walk all six categories (S, T, R, I, D, E) against the feature under analysis. For each one:

1. Identify the threats specific to THIS feature — not generic category descriptions — and name the concrete components each threat lands on in `affected_components`.
2. Rate each threat with a CVSS-lite severity: Critical, High, Medium, or Low.
3. In `description`, state the attack path in one or two sentences, using the reference skill's threat-tree or data-flow notation when the path has several hops.
4. When a category yields nothing, emit one entry for it with `title: "No applicable threats identified"` and a one-sentence description of why — a category is never silently skipped.

## Output
Produces: `security/stride-threats`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  scope: ""    # what was analyzed (feature name + available context)
  threats:
    - id: "T-S-001"  # T-{stride-letter}-{number}
      stride_category: "S|T|R|I|D|E"
      title: ""
      description: ""
      affected_components: []
      severity: "Critical|High|Medium|Low"
  threat_count: 0
  open_items: []
```
