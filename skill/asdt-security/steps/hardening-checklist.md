# Hardening Checklist — Security Specialist

## Purpose
Produce actionable recommendations from all security findings.
Apply the report shared skill. Every finding becomes a concrete action item.

## Inputs
Both inputs arrive ALREADY INJECTED as `### INPUT {topic_key}` blocks; never self-fetch them.

- `security/stride-threats` — STRIDE threats with severity. Extract: Critical + High severity threats only.
- `security/owasp-findings` — OWASP findings with verdicts. Extract: all findings sorted by severity, plus `status`, `reason`, and `evidence` on each.

**DEGRADATION — `security/owasp-findings` is optional (produced only at `risk_surface: high`, where the `owasp-analysis` step runs)**: when it arrives as `### INPUT security/owasp-findings: UNRESOLVED`, build the entire checklist from `security/stride-threats` alone, set `source: "stride-threats"` on every item, and leave the `not_applicable` roll-up empty; append "owasp-findings absent — OWASP Top 10 coverage not assessed at this risk surface, checklist derived from STRIDE threats only" to open_items. Never block on this input.

## Context budget
stride-threats (Critical/High) + owasp-findings: max 1,500 tokens.

## Processing
Apply the `report` reference skill:
1. Deduplicate: merge overlapping findings from STRIDE and OWASP. When a merged item comes from both, record `source: "owasp-findings"` and keep the STRIDE threat id in `finding_ref`.
2. Prioritize: Critical first, then High, then Medium, then Low. OWASP findings with `status: mitigated` produce no action item; findings with `status: at_risk` always do.
3. For each finding, write ONE concrete action item:
   - Not "fix authentication" — write "Add rate limiting to /login endpoint: max 5
     attempts per IP per 15 minutes, return HTTP 429 with Retry-After header"
4. Carry the backing observation forward: each item's `evidence` is the `evidence` of the OWASP finding or the affected component of the STRIDE threat it came from. An action item with no evidence is a guess, not a finding.
5. Group by implementation effort: Quick wins (< 1h), Medium (1h-1 day), Significant (> 1 day).
6. Roll up the conscious skips: every OWASP finding with `status: not_applicable` becomes one `not_applicable[]` entry carrying its category and its `reason`, so a deliberate skip survives into the final artifact instead of looking like an omission.
7. Write a security posture summary: what is the overall risk level? What must be fixed
   before launch vs. what can be deferred?

## Output
Produces: `security/security-findings` (primary) and `security/hardening-checklist` (secondary)

This step produces TWO final artifacts. Persist `security-findings` under this step's output_topic_key ({project}/{change}/security/security-findings); persist `hardening-checklist` under its own distinct per-type topic_key {project}/{change}/security/hardening-checklist, noted on this step's workflow.yaml entry. Return one payload covering both persisted keys.

security-findings schema:
```yaml
payload:
  findings:
    - id: "" # Finding identifier (e.g. F-001).
      owasp_category: "" # OWASP Top 10 2021 category (e.g. A01:2021 - Broken Access Control).
      title: "" # Short title naming the vulnerability.
      description: "" # What the vulnerability is and where it manifests.
      severity: "Low" # one of: Low|Medium|High|Critical
      cwe: "" # Optional CWE identifier (e.g. CWE-79).
      recommendation: "" # Concrete remediation action for the Developer specialist.
  summary: "" # One-sentence security posture summary for the organizational knowledge record.
  overall_risk: "Low" # Overall risk rating for this feature before hardening. | one of: Low|Medium|High|Critical
  open_items: []
```

hardening-checklist schema:
```yaml
payload:
  quick_wins:            # < 1h
    - item: ""           # one concrete action, specific enough to implement without asking a question
      priority: ""       # Critical|High|Medium|Low
      finding_ref: ""    # id of the STRIDE threat or OWASP finding this came from
      evidence: ""       # the observation behind it, carried from the source finding
      source: "stride-threats|owasp-findings"
  medium_effort:         # 1h - 1 day
    - item: ""
      priority: ""
      finding_ref: ""
      evidence: ""
      source: "stride-threats|owasp-findings"
  significant:           # > 1 day
    - item: ""
      priority: ""
      finding_ref: ""
      evidence: ""
      source: "stride-threats|owasp-findings"
  not_applicable:        # conscious skips rolled up from owasp-findings; empty when that input was UNRESOLVED
    - owasp_category: ""
      reason: ""
  must_fix_before_launch: []
  can_defer: []
  summary: ""
  open_items: []
```
