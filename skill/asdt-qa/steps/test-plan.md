# Test Plan — QA Specialist

## Purpose
Find what the acceptance criteria missed, turn it into concrete test cases, and give a
go/no-go verdict. One step, one artifact.

## Inputs
- `{project}/{change}/pm/handoff` — OPTIONAL. The acceptance criteria and NFR targets
- `{project}/{change}/developer/handoff` — OPTIONAL. What was actually built:
  `files_changed`, `open_items` (including any `AC not covered:` lines)
- `{project}/{change}/architect/handoff` — OPTIONAL. The design and its declared risks

All arrive ALREADY INJECTED — never self-fetch. QA can run with NONE of them: work from the
raw request and the codebase, and note `ASSUMED: no upstream hand-off — criteria read from
the raw request` in `open_items`. Never block on a missing input.

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
   Spend your effort here. A plan that only re-states the ACs as tests has added nothing.

3. **Strategy — three lines.** The unit / integration / e2e split for this change, and a
   coverage target WITH the reason it is that number. Concurrency and dependency-failure
   cases cannot be unit tests; say which level each group lands on.

4. **Test cases.** Given/When/Then for each AC and for each edge case worth a test — not
   every edge case earns one, and saying which ones you dropped is part of the plan. When
   `developer/handoff` exists, reference the `files_changed` path each case exercises, so a
   reader can go from case to code without guessing.

5. **Verdict.** `go` or `no-go`, with two lines of why. A blocking AC gap or an uncovered
   critical path is a `no-go`; everything else is a shipping condition, not a block.

**Never emit a pass or fail on something you did not run.** This step executes nothing. If
`pm/handoff` carries NFR targets, list the command the USER can run to measure each one and
what a healthy result looks like — that is an offer, not a result. A performance claim with
no measurement behind it is worse than no claim.

## Output
Produces: `qa/handoff`

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
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
      exercises: ""           # files_changed path, when developer/handoff exists
  measurement_offered:        # NEVER a verdict — commands the USER may run
    - target: ""              # the NFR from pm/handoff
      command: ""
      healthy: ""
  verdict: "go | no-go"
  verdict_why: ""             # two lines
  open_items: []              # ASSUMED: prefix for anything unverified
```
