# Test Plan — QA Specialist

## Purpose
Find what the acceptance criteria missed, turn it into concrete test cases, and give a
go/no-go verdict. One step, one artifact.

## Inputs
- `{project}/{change}/pm/handoff` — OPTIONAL. The acceptance criteria and NFR targets
- `{project}/{change}/developer/handoff` — OPTIONAL. Read its `stage` first. `spec` is a plan,
  nothing built: plan cases against its `files_to_create`/`files_to_modify` and
  `acceptance_criteria` (a `files_changed` it carries is a PRIOR delivery, not this plan). `implemented`: `files_changed` (or
  `code_snippets` when it ran plan-only), the carried `acceptance_criteria`, `open_items`
  (including any `AC not covered:` lines), and `verification` — never claim a pass it did not
  record
- `{project}/{change}/architect/handoff` — OPTIONAL. The design and its declared risks
- `{project}/{change}/ux-ui/handoff` — OPTIONAL. Extract: `flows` — every step's `branches`,
  the empty, loading, and error states the flows name as steps, and the step's `copy` wherever
  an assertion should check exact wording; the `states` and `extremes` of `screens`; and only
  the `disabled`, `loading`, and `error` states of `component_designs`
- `{project}/{change}/security/handoff` — OPTIONAL. Extract: `risks` — each finding's
  `severity` and `mitigation`, and `files_hint` for where the mitigations land

All arrive ALREADY INJECTED — never self-fetch. QA can run with NONE of them: work from the
raw request and the codebase, and note `ASSUMED: no upstream hand-off — criteria read from
the raw request` in `open_items`. Never block on a missing input.

**DEGRADATION**: if `ux-ui/handoff` is UNRESOLVED, find the state cases in the code alone —
UX not running is the common case and needs no `open_items` entry.
**DEGRADATION**: if `security/handoff` is UNRESOLVED, there are no findings to prove and the
verdict carries no Security condition — Security not running is the common case and needs no
`open_items` entry.

## Processing

1. **AC gaps.** Take the inherited acceptance criteria and judge them against
   `asdt-core/references/testing.md`: is each one atomic, measurable, and paired with a
   negative case? List every gap with its type. A criterion no test could observe is a
   blocking gap; a missing negative case is not. If the developer hand-off carries
   `AC not covered:` lines, those are gaps too — carry them in.

2. **Edge cases — this is the job.** The ACs describe what someone already thought of. Your
   value is everything they did not. Work the categories in the reference that this change
   actually touches, and group what you find:
   - **Input** — boundaries, null vs empty vs absent, the semantically impossible value
   - **State** — invalid transitions, repeated transitions, terminal states
   - **Concurrency** — double submit, racing writers, read during partial write
   - **Dependency failure** — timeout, 500, connection dropped mid-operation
   When `ux-ui/handoff` exists, every flow branch and every empty, loading, and error state
   it names is a candidate case — file each under the category it exercises — and so is every
   screen state or extreme, and every extracted component state, the flows do not already cover. A branch the designer drew and
   no AC mentions is exactly the case this step exists to find.
   Spend your effort here. A plan that only re-states the ACs as tests has added nothing.

3. **Strategy — three lines.** The unit / integration / e2e split for this change, and a
   coverage target WITH the reason it is that number. Concurrency and dependency-failure
   cases cannot be unit tests; say which level each group lands on.

4. **Test cases.** Given/When/Then for each AC and for each edge case worth a test — not
   every edge case earns one, and saying which ones you dropped is part of the plan. When
   `developer/handoff` exists, reference the path each case exercises — from `files_changed`, or
   the planned files at `stage: spec` — so a reader can go from case to code without guessing.
   When `security/handoff` exists, every finding gets a case that proves its mitigation holds —
   the attack attempted, and refused — with `mitigates` naming the finding's `risk`. A
   mitigation no test can observe (a rotated secret, a policy) goes to `measurement_offered`
   instead, as the check the USER can run.

5. **Verdict.** `go` or `no-go`, with two lines of why. A blocking AC gap, an uncovered
   critical path, or — when `developer/handoff` exists — a Security `high` whose mitigation
   leaves no trace in its files (`files_changed` at `implemented`, the planned files at `spec`)
   is a `no-go`: a `high` reaches real data or identity, and a fix that is nowhere in the
   change ships the threat with it. Everything else is a shipping condition, not a block. At
   `stage: spec`, or `implemented` with `mode: plan-only`, nothing is built, so the verdict is
   `no-go — not built yet`, and
   `verdict_why` carries your judgment of whether the plan is ready to build.

**Never emit a pass or fail on something you did not run.** This step executes nothing. If
`pm/handoff` carries NFR targets, list the command the USER can run to measure each one and
what a healthy result looks like — that is an offer, not a result. A performance claim with
no measurement behind it is worse than no claim.

## Output
Produces: `qa/handoff`

Persist via memory **save** under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`.

```yaml
payload:
  what: ""                    # the quality posture of this change, one sentence
  gaps:                       # AC completeness gaps
    - criterion: ""
      gap: "untestable | ambiguous | incomplete | missing-negative | missing-nonfunctional"
      blocking: true          # the first three block; the rest do not
  edge_cases:
    - category: "input | state | concurrency | dependency-failure"
      case: ""                # what is not covered by any AC
      risk: ""                # what breaks if it goes untested
  strategy:
    split: ""                 # unit / integration / e2e for this change
    coverage_target: ""       # the number AND why it is that number
  test_cases:
    - id: ""
      given: ""
      when: ""
      then: ""
      level: "unit | integration | e2e"
      exercises: ""           # the built or planned path, when developer/handoff exists
      mitigates: ""           # the Security finding this case proves closed, when security/handoff exists
  measurement_offered:        # NEVER a verdict — commands the USER may run
    - target: ""              # the NFR from pm/handoff, or the Security mitigation no test can observe
      command: ""
      healthy: ""
  verdict: "go | no-go | no-go — not built yet"
  verdict_why: ""             # two lines
  open_items: []              # ASSUMED: prefix for anything unverified
```
