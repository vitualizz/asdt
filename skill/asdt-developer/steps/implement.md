# Implement — Developer Specialist

## Purpose
Generate implementation code for each task, respecting existing conventions.

## Inputs
- `developer/dev-tasks` (optional — tier-gated: the `tasks` step runs at complex): ordered task list with files and dependencies
- `developer/dev-design` (optional — tier-gated: the `design` step runs at moderate and above): technical approach, key constraints
- `developer/dev-spec` (optional — spec→code traceability): acceptance criteria for this change

Extract from dev-tasks: `tasks` list.
Extract from dev-design: `key_constraints`, `data_model` field shapes.
Extract from dev-spec: `acceptance_criteria[]` (for the traceability_report — see Processing).

**Soft input (no `dev-spec`)**: if dev-spec is ABSENT, do NOT fail — note the absence in
`open_items`, emit an empty `traceability_report[]`, and proceed with available context.

**Fallback (no `dev-tasks`)**: in simple/moderate tier the `tasks` step does not run, so
`developer/dev-tasks` is ABSENT. Do NOT fail and do NOT STOP on the missing artifact alone —
fall back to `developer/dev-spec` + `developer/dev-design` as primary inputs and derive the
task list inline from those two artifacts, then proceed. The Mode resolution below (PLAN-ONLY
vs WRITING, gated on declared edit roots) applies unchanged to the inline-derived tasks.

**DEGRADATION — `dev-design` is optional (the `design` step runs at `moderate` and above, so it is
absent at `simple`)**: when it arrives as `### INPUT {project}/{change}/developer/dev-design: UNRESOLVED`,
fall back to `developer/dev-spec` as the design authority — take `key_constraints` from the spec's
stated constraints and non-functional requirements, derive data-model field shapes from the
acceptance criteria, and mark every shape the spec does not settle as unresolved rather than
inventing it; append "developer/dev-design absent — implementation derived from dev-spec at spec
granularity" to open_items. Never block on this input. At `simple` BOTH `dev-tasks` and `dev-design`
are absent, so `developer/dev-spec` is the sole primary input and the inline-derived task list comes
from it alone.

## Context budget
dev-tasks + dev-design summary: max 3,000 tokens. Generate code for tasks in batches
if the task list is large.

## Mode resolution (do this FIRST)
1. Resolve `allowedEditRoots` = union of `files_to_create` + `files_to_modify` across all tasks
   in `dev-tasks` (cross-check `dev-design` for additional declared targets).
2. If `allowedEditRoots` is EMPTY → PLAN-ONLY MODE: emit the snippet-based artifact (unchanged
   schema below, `code_snippets[]`). Write NO host files.
3. If `allowedEditRoots` is NON-EMPTY → WRITING MODE: for each task, write the real file(s) to
   disk, but ONLY to paths within `allowedEditRoots`. Before each write, confirm the target path
   is under a declared root; if not, STOP and report the unsafe path in `open_items` — do not write.
4. Match existing conventions: read the current content of any `files_to_modify` target first
   (per `skills/code-generation.md`) before editing it.

## Processing

### Plan-only mode
For each task in dev-tasks:
1. Generate the implementation code respecting `key_constraints` from dev-design.
2. Follow naming conventions from the platform-summary (loaded by platform-context shared skill).
3. Apply early return pattern, no global state, small functions.
4. Emit the code as inline `code_snippets[]` — write nothing to the host filesystem.

### Writing mode
For each task in dev-tasks:
1. Confirm the task's declared `files_to_create`/`files_to_modify` paths are within `allowedEditRoots`;
   if any path falls outside, STOP before writing it, do not expand scope, and record the unsafe
   path plus the triggering task in `open_items`.
2. Read existing file content for any `files_to_modify` target before editing (match conventions,
   avoid clobbering unrelated code).
3. Generate the implementation respecting `key_constraints` from dev-design, naming conventions
   from platform-summary, early-return pattern, no global state, small functions.
4. Write the real file to disk via the filesystem write tool, within the validated root.
5. Record the written path, action (`created`|`modified`), and rationale in the `files_changed[]`
   manifest entry for that task.
6. Compose — never execute — the verification the USER may run: put the build/lint/test commands
   in `suggested_verification.commands` and what a healthy run looks like in
   `suggested_verification.expected`. This step writes code and NEVER runs build, lint, or tests;
   deciding when to verify is always the user's call.

### Traceability (both modes)
Every acceptance criterion must be answerable with "which code addresses this". After generating
code, build a top-level `traceability_report[]`:
1. Read `dev-spec.acceptance_criteria[]`. If dev-spec was absent, emit `traceability_report: []`.
2. For each AC, identify the task(s) / file(s) that address it. Set `coverage_status: covered`
   and list the addressing `task_id`(s) in `addressed_by`.
3. If no task addresses an AC, set `coverage_status: unaddressed` and add a NON-BLOCKING warning
   to `open_items`. Unaddressed ACs are WARNINGS — do NOT halt the step.

## Dual mem_save semantics
This step persists TWICE, and the two saves are DISTINCT — never merge them:
- **PRIMARY (canonical)**: the artifact save under the `output_topic_key` from `workflow.yaml`
  (`developer/dev-implementation`). This is the artifact sub-agents retrieve via their declared `inputs:`.
- **SECONDARY (knowledge record)**: the `decision-preservation` save under its own title pattern
  (`"{specialist-role}: {change-name}"`). This is the permanent organizational memory entry.

See `../asdt-shared/skills/decision-preservation.md` for the shared definition.

## Output
Produces: `developer/dev-implementation`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

The output schema is mode-dependent — set `mode` to the resolved value and emit the matching shape.

### Plan-only mode schema
```yaml
payload:
  mode: "plan-only"
  steps:
    - task_id: "T-001"
      title: ""
      files_to_create: []
      files_to_modify: []
      rationale: ""
      code_snippets:
        - file: ""
          language: ""
          content: ""
  traceability_report:        # top-level — maps each dev-spec AC to the code addressing it
    - ac_id: ""
      ac_text: ""
      addressed_by: []        # [task_id]
      coverage_status: "covered|unaddressed"
  summary: ""                 # ≤ 150 tokens — consumed by decision-preservation
  open_items: []
```

### Writing mode schema
```yaml
payload:
  mode: "writing"
  allowedEditRoots: []        # resolved list, recorded verbatim for traceability
  files_changed:
    - path: ""
      action: "created|modified"
      task_id: "T-001"
      rationale: ""
  unsafe_skipped: []          # paths STOPPED on, with the triggering task_id
  traceability_report:        # top-level — maps each dev-spec AC to the code addressing it
    - ac_id: ""
      ac_text: ""
      addressed_by: []        # [task_id]
      coverage_status: "covered|unaddressed"
  suggested_verification:     # commands the USER may run; this step never runs them
    commands: []              # e.g. the project's build, lint, and test commands
    expected: ""              # what a healthy run looks like
  summary: ""                 # ≤ 150 tokens — consumed by decision-preservation
  open_items: []
```
