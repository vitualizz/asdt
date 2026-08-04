# Implement — Developer Specialist

## Purpose
Write the implementation — and its tests when TDD is on — respecting existing conventions and
never leaving the declared edit roots.

## Inputs
- `dev-spec` — injected from the orchestrator's context as `### INPUT dev-spec`. Extract:
  `files_to_create`, `files_to_modify`, `acceptance_criteria[]`, `approach`, `key_constraints`
- `{project}/{change}/architect/handoff` (optional) — the architectural decision, from Engram

**DEGRADATION**: if `dev-spec` is UNRESOLVED, do NOT proceed to writing — there are no declared edit roots, so PLAN-ONLY mode applies; note `ASSUMED:` in open_items. If `architect/handoff` is UNRESOLVED, take the design authority from `dev-spec` and note `ASSUMED:` in open_items.

## Mode resolution (do this FIRST)
1. Resolve `allowedEditRoots` = union of `files_to_create` + `files_to_modify` declared in
   `dev-spec` (cross-check `architect/handoff` for additional declared targets).
2. If `allowedEditRoots` is EMPTY → PLAN-ONLY MODE: emit the snippet-based artifact (unchanged
   schema below, `code_snippets[]`). Write NO host files.
3. If `allowedEditRoots` is NON-EMPTY → WRITING MODE: write the real file(s) to disk, but ONLY
   to paths within `allowedEditRoots`. Before each write, confirm the target path is under a
   declared root; if not, STOP and report the unsafe path in `open_items` — do not write.
4. Match existing conventions: read the current content of any `files_to_modify` target first
   (per `../asdt-core/references/conventions.md`) before editing it.

## Processing

### Plan-only mode
1. Generate the implementation code respecting `key_constraints` from the spec, the
   platform-summary conventions, early-return, no global state, small focused functions.
2. Emit the code as inline `code_snippets[]` — write nothing to the host filesystem.

### Writing mode
1. Confirm each declared path is within `allowedEditRoots`; if any path falls outside, STOP
   before writing it, do not expand scope, and record the unsafe path in `open_items`.
2. Read existing file content for any `files_to_modify` target before editing (match
   conventions, avoid clobbering unrelated code).
3. Generate the implementation respecting `key_constraints`, the platform-summary conventions,
   early-return, no global state, small functions.
4. Write the real file to disk via the filesystem write tool, within the validated root.
5. Record the written path, action (`created`|`modified`), and rationale in `files_changed[]`.
6. Compose — never execute — the verification the USER may run: the build/lint/test commands in
   `suggested_verification.commands`, what a healthy run looks like in `.expected`. This step
   writes code and NEVER runs build, lint, or tests; when to verify is always the user's call.

### Tests
Generate tests HERE, in this same step, when `strict_tdd: true` in `.asdt/config.yaml` or the
user asked for them. Otherwise skip this section entirely. Tests obey the SAME mode and the
SAME `allowedEditRoots` as the code above — in plan-only mode they are `code_snippets[]`
entries like any other file; in writing mode they are real files, and a test path outside the
declared roots STOPS exactly as a source path does.

Per unit under test: one happy-path test for the acceptance criterion, and one edge-case test
for the most likely failure mode. Follow the project's existing test framework; use
table-driven cases where the framework offers them; test behavior, never internals.

This step WRITES tests; it NEVER runs them. Running build, lint, or a test suite is always the
user's call, so `suggested_verification.commands` is an offer to the user, not a step this
specialist performs.

### Coverage
Every acceptance criterion should be answerable with "which code addresses this". After
generating the code, walk `dev-spec.acceptance_criteria[]`: for each one no file addresses,
append a single line to `open_items` — `AC not covered: {ac text}`. Warnings, never a halt.

## Output
Produces: `developer/handoff` — persist via `mem_save` under this step's `output_topic_key`,
using the canonical hand-off schema from `asdt-core/protocol.md`. Set `mode` to the resolved
value; `files_changed` carries real paths in writing mode, `code_snippets` the code in plan-only.

```yaml
payload:
  what: ""                    # what was implemented, one sentence
  mode: "writing | plan-only"
  allowedEditRoots: []        # resolved list, recorded verbatim (writing mode)
  files_changed:              # writing mode
    - path: ""
      action: "created|modified"
      rationale: ""
  code_snippets:              # plan-only mode
    - file: ""
      language: ""
      content: ""
  unsafe_skipped: []          # paths STOPPED on, and what triggered them
  suggested_verification:     # commands the USER may run; this step never runs them
    commands: []
    expected: ""
  decisions: []               # implementation choices worth carrying forward
  risks: []                   # {risk, mitigation}
  open_items: []              # includes one "AC not covered: ..." line per uncovered AC
```
