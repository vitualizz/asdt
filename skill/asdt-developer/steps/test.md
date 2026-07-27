# Test — Developer Specialist

## Purpose
Generate tests for the implementation, covering happy paths and key edge cases.

## Inputs
- `developer/dev-tasks` (optional — tier-gated: the `tasks` step runs at complex): task list with acceptance criteria references
- `developer/dev-implementation`: code snippets per task

Extract from dev-tasks: `tasks[].ac_ref`, `tasks[].id`.
Extract from dev-implementation: `steps[].code_snippets[].content` (signatures only, not full bodies).

**DEGRADATION — `dev-tasks` is optional (the `tasks` step runs at `complex` only, so it is absent at
`simple` and `moderate`)**: when it arrives as `### INPUT {project}/{change}/developer/dev-tasks: UNRESOLVED`,
derive coverage from `developer/dev-implementation` instead of from a task list — treat each entry of its
`files_changed[]` (writing mode) or `steps[]` (plan-only mode) as one unit under test, key each suite by that
entry's `task_id`, and take the behavior to assert from that artifact's `traceability_report[]` acceptance
criteria in place of `tasks[].ac_ref`; append "developer/dev-tasks absent — test coverage derived from
dev-implementation units instead of the task list" to open_items. Never block on this input.

## Context budget
dev-tasks AC list + dev-implementation function signatures: max 2,500 tokens.

## Mode resolution (do this FIRST)
1. Inherit `mode` and `allowedEditRoots` from `dev-implementation`'s payload — do NOT
   re-derive them independently. The `test` step always runs in the SAME mode `implement`
   resolved for this change.
2. If `dev-implementation.payload.mode` is `plan-only` → PLAN-ONLY MODE: emit the
   snippet-based artifact (unchanged schema below, `test_cases[].code_snippet`). Write NO
   host files.
3. If `dev-implementation.payload.mode` is `writing` → WRITING MODE: write real test files
   to disk, but ONLY to paths within the SAME `allowedEditRoots` `implement` used. Before
   each write, confirm the target path is under a declared root; if not, STOP and report the
   unsafe path in `open_items` — do not write.
4. Mode resolution decides only WHERE the resulting test artifact lives — an Engram snippet
   or a real file. It never decides whether the tests execute: in BOTH modes this step writes
   tests and never runs them.

## Processing

### Plan-only mode
For each task in dev-tasks — or, when dev-tasks is UNRESOLVED, each unit derived per DEGRADATION above:
1. Write one happy-path test covering the acceptance criterion (ac_ref) as a `code_snippet`.
2. Write one edge-case test for the most likely failure mode as a `code_snippet`.
3. Follow the existing test framework convention from platform-summary.
4. Use table-driven tests where appropriate.

### Writing mode
For each task in dev-tasks — or, when dev-tasks is UNRESOLVED, each unit derived per DEGRADATION above:
1. Confirm the test file path (from `dev-tasks` test entries / declared paths, or from the
   `dev-implementation` entry the unit was derived from) is within the
   inherited `allowedEditRoots`; if it falls outside, STOP before writing it, do not expand
   scope, and record the unsafe path plus the triggering task in `open_items`.
2. Write one happy-path test covering the acceptance criterion (ac_ref) and one edge-case
   test for the most likely failure mode, as real test files on disk.
3. Follow the existing test framework convention from platform-summary; use table-driven
   tests where appropriate.
4. Compose — never execute — the single command the USER can run to execute these tests, and
   record it in `suggested_run_command`. For each test case, record in `expected_when_run` what
   the user should observe when they run it.

This step WRITES tests; it NEVER runs them. Running build, lint, or a test suite is always the
user's call, so `suggested_run_command` is an offer to the user, not a step this specialist
performs.

Do NOT test implementation internals — test behavior.

## Output
Produces: `developer/dev-tests`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

The output schema is mode-dependent — set `mode` to the inherited value and emit the matching shape.

### Plan-only mode schema
```yaml
payload:
  mode: "plan-only"
  test_suites:
    - task_id: "T-001"
      test_cases:
        - id: "TC-001"
          title: ""
          type: "unit|integration"
          given: ""
          when: ""
          then: ""
          code_snippet:
            file: ""
            language: ""
            content: ""
  open_items: []
```

### Writing mode schema
```yaml
payload:
  mode: "writing"
  allowedEditRoots: []        # inherited verbatim from dev-implementation
  suggested_run_command: ""   # the single command the USER can run to execute these tests
  test_suites:
    - task_id: "T-001"
      test_cases:
        - id: "TC-001"
          title: ""
          type: "unit|integration"
          given: ""
          when: ""
          then: ""
          file: ""            # real path written, within allowedEditRoots
          expected_when_run: "" # what the user should observe when they run it
  unsafe_skipped: []          # paths STOPPED on, with the triggering task_id
  open_items: []
```
