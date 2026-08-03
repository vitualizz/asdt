# Harden — Security Specialist

## Purpose
Turn the assessment into prioritized findings and an actionable hardening checklist. One
step, one artifact.

## Inputs
- `security-assessment` — injected from the orchestrator's context as
  `### INPUT security-assessment`. Extract: `attack_surface`, `threats`, `open_items`

**DEGRADATION**: if the assessment is UNRESOLVED, produce no findings rather than invented ones — record `ASSUMED: assessment unavailable — no findings produced` in open_items.

## Processing

1. **Prioritize.** Give every threat a severity in ONE word — `high`, `medium`, or `low` —
   and nothing else. No CVSS vectors, no numeric scores, no composite ratings. Severity is a
   judgment about what an attacker gets and how easily; write that judgment as the finding's
   `why`, and let the reader disagree with it.

   Rank by what an attacker actually gains. A `high` is something that reaches real data or
   real identity. Everything cannot be `high` — an inflated severity list gets skimmed, and
   then the real one gets skimmed with it.

2. **Give every finding a mitigation.** Concrete and specific to this codebase: the check to
   add, the path to guard, the default to change. "Validate input" is not a mitigation;
   "reject non-allowlisted schemes in `fetchRemote`" is.

3. **Build the checklist.** The findings say what is wrong; the checklist says what to DO,
   in order, each item small enough to close in one sitting and phrased so someone can tell
   whether it is done. A checklist item nobody can mark complete is a wish.

Carry any `ASSUMED:` entries from the assessment forward — a threat you could not verify
stays visible.

## Output
Produces: `security/handoff`

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`. Findings and checklist are SECTIONS of this one
hand-off — there is no second artifact and no second key.

```yaml
payload:
  what: ""                    # the security posture of this change, one sentence
  risks:                      # the findings, highest severity first
    - risk: ""                # what an attacker can do
      severity: "high | medium | low"
      why: ""                 # what they gain, and how easily
      evidence: ""            # the path, symbol, or absence that grounds it
      mitigation: ""          # the concrete change that closes it
  decisions: []               # hardening choices worth carrying forward
  constraints: []             # the checklist — ordered, each item verifiable as done
  files_hint: []              # where the mitigations land
  open_items: []              # ASSUMED: prefix for anything unverified
```
