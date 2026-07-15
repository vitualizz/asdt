# Write — Init Specialist

## Purpose
Write the four `.asdt` files from the detected stack plus the human's clarify
answers. This is init's only filesystem-writing step. It runs as a `builder`
sub-agent: it cannot pause to ask questions — every question was already asked by
the inline `clarify` step and the answers were injected into this prompt.

## Inputs
- `init/stack-detection` (injected) — the explore step's `detected_stack`,
  `lang_roots`, `fields`, and `ambiguities[]`. Consume the injected block
  directly; do NOT re-fetch it.
- **`### CLARIFY ANSWERS` block (injected)** — the orchestrator composes this
  from the inline clarify step and prepends it to your prompt. Consume it
  directly, never re-fetch it. Its shape:

  ```
  ### CLARIFY ANSWERS
  answers: { field: value, ... }     # may be {} when there was nothing to ask
  skipped: true|false                # true → non-interactive harness, defaults applied
  blocking_open_items: []            # non-empty → HALT (see Halt contract)
  ```

## Context budget
`stack-detection`: max 2,000 tokens. The CLARIFY ANSWERS block is small by
construction. Do not pull anything else into context.

## Recalibration contract
The Engram gate already passed PRE-EXPLORE (the orchestrator checked its own tool
list before launching explore). You do not re-run it.

**Preserve `source: manual` fields. NEVER silently overwrite them.** When
`.asdt/knowledge/project-context.yaml` already exists, any field whose existing
`source` is `manual` was set by a human in a prior recalibration review. Carry it
forward unchanged unless a clarify answer for that exact field explicitly
overrides it. Surface every preserved field in `settings_preserved[]` so the
outcome is auditable.

The same rule extends to `design_fingerprint.<concern>` entries in
`.asdt/knowledge/platform.yaml`: a concern whose existing `source` is `manual` is
carried forward unchanged unless a clarify answer for that exact concern
overrides it, and it appears in `settings_preserved[]` too — a re-init NEVER
silently overwrites a manually-set design_fingerprint concern.

## Processing

`.asdt/` holds static reference data — bootstrapped once and refreshed only on a
deliberate recalibration, never per-change.

1. **`.asdt/config.yaml`** — write `memory.provider: engram` plus any preserved
   settings carried forward from an existing config:

   ```yaml
   memory:
     provider: engram
   ```

2. **`.asdt/knowledge/platform.yaml`** — populate only what the bounded scan
   determined deterministically. `conventions.file_structure` is the one-line
   sentence derived from top-level directory matches. Populate
   `design_fingerprint` from `stack-detection.fields.design_fingerprint`: emit
   one `FieldValue` entry per FIRED pack, in canonical order (`i18n`,
   `css_approach`, `orm`, `state_management`, `ci_cd`, `lint`,
   `code_intelligence`). Omit any pack that did not fire — no key. Write
   `design_fingerprint: {}` only when nothing was emitted at all. A clarify
   answer for a `design_fingerprint.<concern>` field makes that entry
   `source: manual`; a default applied non-interactively keeps the detected
   `source`/`confidence` and is recorded with `origin: default` — the write
   NEVER halts on a design_fingerprint default:

   ```yaml
   schema_version: "1"
   scanned_at: {current UTC timestamp, ISO 8601}
   detected_stack: {stack-detection.detected_stack}
   conventions:
     file_structure: {one-line description}
   design_fingerprint:        # one FieldValue per FIRED pack; {} only when none fired/emitted
     css_approach: { value: "…", source: "…", confidence: "…" }
     # … one entry per fired pack, in canonical order
   ```

3. **`.asdt/knowledge/platform-summary.yaml`** — derived FROM `platform.yaml`,
   never re-analyzed from scratch. Flatten `design_fingerprint` to
   `{concern: value}` scalars — drop the `source`/`confidence` annotations, and
   omit any concern whose value is `unknown`, `none`, or absent:

   ```yaml
   schema_version: "1"
   stack: {platform.yaml detected_stack}
   file_structure: {platform.yaml conventions.file_structure}
   design_fingerprint:        # flat {concern: value} scalars; omit unknown/none/absent
     {concern}: {value}
   ```

4. **`.asdt/knowledge/project-context.yaml`** — built from
   `stack-detection.fields` with each applied clarify answer overlaid. A field
   answered by the human becomes `source: manual`. Record every applied
   `Ambiguity` answer with its `origin` (`user` | `default`). The
   `human_nuance` region (below) is rendered at the END of this file, after the
   detected fields:

   ```yaml
   schema_version: "1"
   detected_at: {current UTC timestamp, ISO 8601}
   is_monorepo: { value: "…", source: "…", confidence: "…" }
   test_runner: { value: "…", source: "…", confidence: "…" }
   naming_style: { value: "…", source: "…", confidence: "…" }
   architectural_style: { value: "…", source: "…", confidence: "…" }
   # ASDT:NUANCE:BEGIN
   human_nuance:
     - topic: "replaceMarkerRegion"
       note: "fail-loud on partial markers; any new region writer must mirror its halt semantics"
       source: manual
       origin: user
   # ASDT:NUANCE:END
   ```

   **`human_nuance` — routing from clarify answers.** The enrichment step
   (SKILL.md §2.4) surfaces `nuance.*` ambiguities; the human's answers arrive in
   the `### CLARIFY ANSWERS` block's `answers{}`. For each answer key matching
   `^nuance\.(.+)$`, take the captured suffix as the `topic` slug and the answer
   value as the `note`:

   - **non-empty value** → append one entry
     `{ topic: <slug>, note: <value>, source: manual, origin: user }`.
   - **empty value** (the human skipped it) → NO entry.

   Non-`nuance.*` answer keys are unaffected — they overlay the detected fields
   exactly as before. When no `nuance.*` answer produces an entry, render the
   empty form `human_nuance: []` inside the markers.

   **Merge-into-existing algorithm.** This region mirrors the
   `replaceMarkerRegion` fail-loud contract in
   `internal/installer/registry_gen.go` — same marker discipline, same
   halt-on-partial rule. Use the YAML-comment markers `# ASDT:NUANCE:BEGIN` and
   `# ASDT:NUANCE:END`. When `project-context.yaml` already exists, read its bytes
   and count the begin/end markers:

   - **Both counts 0 (region absent).** If there are non-empty entries to write —
     OR the file is being freshly created — append a fresh region at EOF: a blank
     line, then `# ASDT:NUANCE:BEGIN`, the body, then `# ASDT:NUANCE:END`. If there
     are NO entries AND the file pre-existed, this is a NO-OP — do not append an
     empty region to a file that never had one.
   - **begin count ≠ 1 OR end count ≠ 1** (partial or duplicated markers) → HALT
     with ZERO writes (see Halt contract clause (c)).
   - **end marker precedes begin marker** (reversed) → HALT with ZERO writes.
   - **Region present and valid.** Parse the existing entries and preserve them,
     then merge the new answers: append an entry for each new `topic`, and UPDATE
     the `note` in place for any `topic` that already exists. Render the body —
     the non-empty `human_nuance:` block, or `human_nuance: []` when the merge
     leaves it empty. Idempotency: if the rendered body equals what already sits
     between the markers, no-op (no write). Otherwise splice
     `content[:regionStart] + body + content[endIdx:]`, exactly as
     `replaceMarkerRegion` does.

   The existing `source: manual` flat-field preservation rule (above) is
   UNCHANGED and applies alongside this region — `human_nuance` merging does not
   alter how the detected fields carry manual values forward.

All four files stay bounded — their size grows with the number of detected
stacks and human notes, never with repo size.

## Halt contract

HALT with an error and ZERO writes if any condition holds:

- **(a) The `### CLARIFY ANSWERS` block is absent from the prompt.** Absence
  means the orchestrator failed to pass it — NOT that there were no answers. The
  block is REQUIRED even when empty (`answers: {}`, `skipped: true|false`,
  `blocking_open_items: []`). No block → halt; do not write partial config.
- **(b) `blocking_open_items[]` is non-empty.** A blocking open item means a
  non-skippable ambiguity went unresolved; writing config on top of it would
  bake in a wrong default. Halt and surface the items.
- **(c) Malformed `human_nuance` markers in an existing `project-context.yaml`.**
  If the file has a partial or duplicated marker set — `# ASDT:NUANCE:BEGIN`
  count ≠ 1, `# ASDT:NUANCE:END` count ≠ 1, or the end marker preceding the begin
  marker — HALT with an error and ZERO writes. A damaged region must be caught,
  never half-written over. Both markers absent is NOT malformed — that is the
  region-absent path handled in step 4, not a halt.

On halt, write nothing and return the error in the envelope.

## Output
Produces: `init/write-summary`

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Schema:
```yaml
payload:
  files_written: []           # absolute or repo-relative paths of the .asdt files written
  settings_preserved: []      # source:manual fields carried forward unchanged
  applied_answers:            # every Ambiguity answer applied, with its origin
    - field: ""
      value: ""
      origin: "user | default"
  open_items: []
```

If the `mem_save` call fails, record the failure in `open_items` — the files on
disk are authoritative and the writes already succeeded. Do NOT halt or roll back
on a persistence failure; the config is the durable outcome, the summary is the
audit trail.
