# OWASP Analysis — Security Specialist

## Purpose
Check the attack surface against the OWASP Top 10 (2021).
Systematic coverage prevents the "I forgot to check X" failure mode.

## Inputs
- `security/attack-surface` — arrives ALREADY INJECTED as an `### INPUT security/attack-surface` block; never self-fetch it. Extract: `entry_points[]`, `data_flows[].vulnerabilities_noted`.

## Context budget
security/attack-surface (entry points + data flows): max 1,200 tokens.

## Processing
The A01–A10 catalogue for 2021 — every category with its detection patterns, stack-specific
concerns, and remediation guidance — lives in `skills/owasp-review.md`, this step's reference
skill. Work from it; it is not restated here.

Walk all ten categories in order against the attack surface and emit at least one `findings[]`
entry per category. Never skip a category and never fold two categories into one entry; when a
single category holds several distinct problems, emit one entry per problem with the same
`owasp_category` and different ids. For every entry produce:

1. `status` — `not_applicable` when nothing in the attack surface can trigger the category at all; otherwise `mitigated` when an existing control already covers it, or `at_risk` when no control does or the control is insufficient.
2. `reason` — REQUIRED when `status: not_applicable` (why the category cannot apply here, so the skip is a conscious one that survives into the final artifact); recommended when `status: at_risk` (why the existing controls fall short).
3. `evidence` — the concrete observation the verdict rests on: a file and symbol when you inspected the repository (`internal/auth/session.go:validateToken`), otherwise the assumption you are reasoning from, stated plainly as an assumption. You reason and inspect; you never run scanners, dependency audits, or any other command, so evidence is never tool output.
4. `severity`, `cwe` where one applies, and a concrete `recommendation` — REQUIRED when `status: at_risk`, since the hardening step turns each one into an action item.

## Output
Produces: `security/owasp-findings`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

```yaml
payload:
  findings:
    - id: "OF-001"
      owasp_category: "A01|A02|...|A10"
      status: "mitigated|at_risk|not_applicable"
      title: ""
      description: ""
      reason: ""       # required when status is not_applicable; recommended when at_risk
      evidence: ""     # file + symbol inspected, or the stated assumption — never scanner output
      severity: "Critical|High|Medium|Low"   # required when status is at_risk
      cwe: ""          # CWE reference number if applicable
      recommendation: ""                     # required when status is at_risk
  open_items: []
```
