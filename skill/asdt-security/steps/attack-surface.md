# Attack Surface Analysis — Security Specialist

## Purpose
Map every entry point where an attacker can interact with the system.
A smaller attack surface = fewer vulnerabilities.

## Inputs
- `security/stride-threats` — arrives ALREADY INJECTED as an `### INPUT security/stride-threats` block; never self-fetch it. Extract: `threats[].affected_components`, `threats[].stride_category`.

## Context budget
security/stride-threats (components + categories only): max 800 tokens.

## Processing
For each affected component identified in STRIDE threats:
1. ENTRY POINTS: list every place external input enters this component
   (HTTP endpoints, form fields, file uploads, message queue events, webhooks, CLI args).
2. TRUST BOUNDARIES: where does data cross from untrusted to trusted zones?
   (External → internal, user-controlled → system-controlled)
3. DATA FLOWS: trace the data from entry point to storage and back.
   At each step: is the data validated? Is it sanitized? Is it escaped on output?
4. AUTHENTICATION CHECKPOINTS: which entry points require authentication?
   Record the answer twice per entry point: `authentication_required` is what the system
   enforces TODAY, `authentication_expected` is what the data it touches DEMANDS. When the
   two differ, `authentication_gap` carries what is missing and why it matters — that field
   is where "should require it but doesn't" gets written down instead of being lost in prose.
5. AUTHORIZATION CHECKPOINTS: for each entry point that touches data,
   is the current user authorized to access that specific data object?

## Output
Produces: `security/attack-surface`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  entry_points:
    - id: "EP-001"
      type: "http|form|upload|queue|webhook|cli"
      description: ""
      authentication_required: true   # actual: enforced by the system today
      authentication_expected: true   # expected: what the data this entry point touches demands
      authentication_gap: ""          # required when expected differs from actual — what is missing and the impact
      authorization_check: ""
      input_validation: ""
  trust_boundaries:
    - from: ""
      to: ""
      controls: []
  data_flows:
    - name: ""
      steps: []
      vulnerabilities_noted: []
  open_items: []
```
