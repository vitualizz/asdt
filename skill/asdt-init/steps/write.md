# Write — Init Specialist

## Purpose
Write the three `.asdt` files — `config.yaml`, `knowledge.yaml`, and its
write-only provenance sidecar — from the detected stack plus the human's clarify
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
`.asdt/knowledge/knowledge.yaml` already exists, any of its four inline
FieldValue fields (`is_monorepo`, `test_runner`, `naming_style`,
`architectural_style`) whose existing `source` is `manual` was set by a human in
a prior recalibration review. Carry it forward unchanged unless a clarify answer
for that exact field explicitly overrides it. Surface every preserved field in
`settings_preserved[]` so the outcome is auditable.

The `design_fingerprint` values in `knowledge.yaml` are flat scalars and carry
no `source` annotation — their provenance lives in the write-only sidecar
`.asdt/knowledge/provenance.yaml` (Processing step 3). During a recalibration,
read the sidecar: a concern whose sidecar `source` is `manual` is carried
forward unchanged unless a clarify answer for that exact concern overrides it,
and it appears in `settings_preserved[]` too — a re-init NEVER silently
overwrites a manually-set design_fingerprint concern. If the sidecar is absent
or unparseable, re-detect the fingerprint fresh and add the open_item
`"provenance.yaml unreadable — design_fingerprint re-detected"`; the write
proceeds normally.

## Processing

`.asdt/` holds static reference data — bootstrapped once and refreshed only on a
deliberate recalibration, never per-change.

1. **`.asdt/config.yaml`** — write `memory.provider: engram` plus any preserved
   settings carried forward from an existing config. Preserve `memory.provider`,
   `strict_tdd`, and EVERY unknown key byte-wise via the same preservation path
   that carries `strict_tdd` forward today — this step never drops a key it does
   not recognize.

   `code_intelligence` is a top-level scalar, **positive-evidence-only**,
   modeled on `strict_tdd`'s placement:

   - The `code_intelligence` pack fired and matched → write
     `code_intelligence: <value>` (e.g. `codegraph`). NO FieldValue wrapper —
     the value is `detected`/`high` by construction.
   - Nothing detected → the key is REMOVED from the file (never `none`, never
     `unknown`).

   `primary_design_surface` is a top-level scalar too, written with the SAME
   mechanics — but explicitly NOT `code_intelligence`'s detection semantics.
   This key is ASKED, never detected:

   - A `primary_design_surface` clarify answer is present → write
     `primary_design_surface: <answer>`, one of `mobile` | `tablet` | `desktop` |
     `none`. NO FieldValue wrapper.
   - No such answer — the question was never asked, e.g. a non-interactive run →
     the key is REMOVED from the file.

   These two states are DIFFERENT and must never be conflated: an ABSENT key
   means the question was never asked, and every consumer treats that as
   `mobile`. An explicit `none` means the human answered that this project has no
   visual surface at all, and consumers emit no responsive output. Never write
   `none` as a stand-in for "unanswered", and never drop the key when the answer
   was `none`.

   ```yaml
   memory:
     provider: engram
   strict_tdd: false            # preserved byte-wise when present
   code_intelligence: codegraph # positive evidence only; key removed when absent
   primary_design_surface: mobile # asked, not detected; key removed when never answered
   ```

2. **`.asdt/knowledge/knowledge.yaml`** — the ONLY knowledge file specialists
   read. Built from `stack-detection` with each applied clarify answer overlaid.
   Canonical key order is BYTE-LOCKED (an unchanged re-run must produce a
   zero-byte diff): `schema_version`, `scanned_at`, `stack`, `file_structure`,
   `design_fingerprint`, `is_monorepo`, `test_runner`, `naming_style`,
   `architectural_style`, then the fenced nuance region at the file tail.

   - `schema_version: "2"` (string literal, quoted).
   - `scanned_at`: current UTC timestamp, ISO 8601.
   - `stack`: language ids in scan order, deduplicated
     (`stack-detection.detected_stack`).
   - `file_structure`: the one-line sentence derived from top-level directory
     matches.
   - `design_fingerprint`: FLAT `{concern: value}` scalars — NO
     source/confidence annotations. Emit one scalar per FIRED pack in canonical
     order (`i18n`, `css_approach`, `orm`, `state_management`, `ci_cd`,
     `lint`); omit any concern whose value is `unknown` or `none`.
     `code_intelligence` is NEVER written here — it lives in `config.yaml`
     (step 1). Write `design_fingerprint: {}` only when nothing was emitted at
     all.
   - `is_monorepo`, `test_runner`, `naming_style`, `architectural_style`: the
     four decision fields keep the inline FieldValue mapping
     `{value, source, confidence}` — the ONLY mapping-form fields in this file.
     A field answered by the human becomes `source: manual`; a default applied
     non-interactively keeps the detected `source`/`confidence` and is recorded
     with `origin: default` — the write NEVER halts on a default.

   ```yaml
   schema_version: "2"
   scanned_at: {current UTC timestamp, ISO 8601}
   stack: {stack-detection.detected_stack}
   file_structure: {one-line description}
   design_fingerprint:        # flat {concern: value} scalars; omit unknown/none; no code_intelligence
     css_approach: tailwind
     ci_cd: github-actions
   is_monorepo: { value: "…", source: "…", confidence: "…" }
   test_runner: { value: "…", source: "…", confidence: "…" }
   naming_style: { value: "…", source: "…", confidence: "…" }
   architectural_style: { value: "…", source: "…", confidence: "…" }
   # ASDT:NUANCE:BEGIN
   human_nuance:
     - topic: "replaceMarkerRegion"
       type: architectural
       note: "fail-loud on partial markers; any new region writer must mirror its halt semantics"
       source: manual
       origin: user
   # ASDT:NUANCE:END
   ```

   **`human_nuance` — typed entries, routed from clarify answers.** The
   enrichment step (SKILL.md, § enrichment) surfaces typed `nuance.*` ambiguities; the
   human's answers arrive in the `### CLARIFY ANSWERS` block's `answers{}`. For
   each answer key matching
   `^nuance\.(architectural|repo_practice|inconsistency_to_review)\.(.+)$`, the
   first capture is the entry `type`, the second capture is the `topic` slug,
   and the answer value is the `note`:

   - **non-empty value** → append one entry
     `{ topic: <slug>, type: <type>, note: <value>, source: manual, origin: user }`.
   - **empty value** (the human skipped it) → NO entry.

   `type` is one of `architectural`, `repo_practice`, or
   `inconsistency_to_review`.

   Non-`nuance.*` answer keys are unaffected — they overlay the detected fields
   exactly as before. When no `nuance.*` answer produces an entry, render the
   empty form `human_nuance: []` inside the markers.

   **Merge-into-existing algorithm.** This region mirrors the
   `replaceMarkerRegion` fail-loud contract in
   `internal/installer/registry_gen.go` — same marker discipline, same
   halt-on-partial rule. Use the YAML-comment markers `# ASDT:NUANCE:BEGIN` and
   `# ASDT:NUANCE:END`. When `knowledge.yaml` already exists, read its bytes
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
     the `note` (and `type`, when the answer key carries one) in place for any
     `topic` that already exists. Render the body — the non-empty
     `human_nuance:` block, or `human_nuance: []` when the merge leaves it
     empty. Idempotency: if the rendered body equals what already sits between
     the markers, no-op (no write). Otherwise splice
     `content[:regionStart] + body + content[endIdx:]`, exactly as
     `replaceMarkerRegion` does.

   The `source: manual` flat-field preservation rule (Recalibration contract
   above) is UNCHANGED and applies alongside this region — `human_nuance`
   merging does not alter how the decision fields carry manual values forward.

3. **`.asdt/knowledge/provenance.yaml`** — the write-only provenance sidecar.
   Written by this step and read SOLELY by this step during a later
   recalibration — no specialist ever reads it, and it never enters any
   injection path. Write it only AFTER `knowledge.yaml` (step 2) has been
   written successfully.

   - `schema_version: "2"` (string literal, quoted).
   - `scanned_at`: MUST equal the `scanned_at` written to `knowledge.yaml` —
     both files come from the same detection in the same write.
   - `design_fingerprint`: a map of EVERY fired pack to its full FieldValue
     `{value, source, confidence}` — INCLUDING packs whose value is `none`
     (e.g. `orm: { value: none, source: detected, confidence: medium }`), so a
     later recalibration can distinguish "probed, found absent" from "never
     probed". Canonical concern order; `code_intelligence` excluded (it lives
     in `config.yaml`).

   ```yaml
   schema_version: "2"
   scanned_at: {same timestamp as knowledge.yaml}
   design_fingerprint:        # one FieldValue per FIRED pack, INCLUDING none-valued
     i18n: { value: "custom (locale files)", source: detected, confidence: medium }
     orm: { value: none, source: detected, confidence: medium }
   ```

   A concern whose prior sidecar `source` is `manual` is carried forward
   unchanged unless a clarify answer for that exact concern overrides it, and
   is surfaced in `settings_preserved[]`. Regeneration invariant: the
   fingerprint scalars in `knowledge.yaml` and the FieldValues here are written
   together, in the same step, from the same detection — never regenerate or
   hand-edit one alone. This file carries fingerprint provenance ONLY;
   user-authored notes belong to the fenced tail region of `knowledge.yaml`
   (step 2) and never appear here.

All three files stay bounded — their size grows with the number of detected
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
- **(c) Malformed `human_nuance` markers in an existing `knowledge.yaml`.**
  If the file has a partial or duplicated marker set — `# ASDT:NUANCE:BEGIN`
  count ≠ 1, `# ASDT:NUANCE:END` count ≠ 1, or the end marker preceding the begin
  marker — HALT with an error and ZERO writes. A damaged region must be caught,
  never half-written over. Both markers absent is NOT malformed — that is the
  region-absent path handled in Processing step 2, not a halt.

On halt, write nothing and return the error in the envelope.

## Post-write self-check

The Halt contract gates everything BEFORE the writes; this gates what actually
landed. After the three files are written — and before reporting them as
written — re-read each one from disk and verify:

1. **It parses.** Each file is valid YAML.
2. **Key order matches.** `knowledge.yaml`'s top-level keys appear in the
   canonical order declared in Processing step 2; `provenance.yaml`'s in the
   order declared in Processing step 3.
3. **The tail marker region is intact and balanced.** In `knowledge.yaml`, the
   fenced region described in Processing step 2 has exactly one begin marker and
   exactly one end marker, with the begin preceding the end. A file with no
   region at all is fine — that is the region-absent path, not a defect.

Re-reading files this step just wrote is a READ, not a build or a test run — it
is permitted here and it is required.

On any mismatch, record it in `open_items` (e.g.
`"knowledge.yaml key order drifted from canonical — verify before the next recalibration"`)
and report the files as written anyway. Do NOT halt and do NOT roll back: the
bytes on disk are authoritative, and a self-check finding is an audit signal,
not a reason to leave the project half-configured.

## Output
Produces: `init/write-summary`

Persist via mem_save under this step's output_topic_key in workflow.yaml; return the payload above with open_items populated.

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
