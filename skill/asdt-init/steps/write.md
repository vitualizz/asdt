# Write — Init Specialist

## Purpose
Write `.asdt/config.yaml` from the detection and the human's answers, never silently overwriting a manual value. This is one of the two steps that write files (protocol §4), and it writes only this one.

## Inputs
- `### INPUT asdt-init/explore`: the detection. Read its `## Detection` section (stack, lang roots, file structure, the key table) and nothing else.
- `### INPUT asdt-init/clarify`: the answers. Read `skipped`, `memory_provider`, `## Answers`, and `## Blocking open items`. Ignore `## Notes to remember`: notes go to memory, through the orchestrator, never into a file.

Each arrives as a path to read. Both are required: see Processing step 1.

## Processing

1. **Halt with zero writes** when any of these holds, and return the reason as your summary:
   - (a) either input is `UNRESOLVED` or unreadable. A missing answers file means the orchestrator failed to hand it over, not that nothing was asked;
   - (b) `## Blocking open items` lists anything: a non-skippable question went unanswered, and writing over it would bake in a wrong value;
   - (c) an existing `.asdt/config.yaml` does not parse as a YAML mapping. A damaged config is reported, never written over.

2. **Prior values.** Resolve each config key's prior value, its source included where it has one, from the existing `.asdt/config.yaml`; a key it omits, or a config that does not exist, has none.

3. **Resolve each key.** The four decision fields (`is_monorepo`, `test_runner`, `naming_style`, `architectural_style`) and each `design_fingerprint.{concern}` take the form `{value, source, confidence}`, decided in this order:
   - an answer with origin `user` → `{value: answer, source: manual, confidence: high}`;
   - an answer with origin `accepted` → the detected entry, unchanged;
   - no answer or a `default` one, and the prior source is `manual` → the prior entry, unchanged, listed under `## Settings preserved`. This holds for a fingerprint concern whose pack did not fire this time too: a default never overwrites a human's value;
   - otherwise → the detected entry. A concern with no detected entry and no `manual` prior is dropped.

   The plain keys:
   - `stack`, `file_structure`: the answer when there is one, else the detection. `stack` is always written, `[]` when empty. `file_structure` is omitted when neither has a value.
   - `code_intelligence`: the detected value when the detection carries the row, otherwise the key is removed. Never `none`, never `unknown`. Unlike `memory`, nothing is kept here: the key states a fact about this project's index, so a run that does not detect one removes it.
   - `primary_design_surface`: the answer (`mobile`, `tablet`, `desktop`, or `none`), else the prior value, else omitted. An omitted key means "never asked" and consumers default it to `mobile`; `none` is the human's answer that there is no visual surface. Never write one in place of the other.
   - `memory.provider`: the payload's `memory_provider`, even over an existing value. When it is `none` or absent, keep the existing `memory` block as-is, or omit the block when there is none: `none` is never written.
   - `strict_tdd`: carried forward byte-for-byte when present, never added.
   - Every key this step does not recognize: carried forward byte-for-byte.

4. **Render** the file in this canonical key order, then every unrecognized key in its original relative order:

   ```yaml
   memory:
     provider: {memory_provider}
   strict_tdd: false
   code_intelligence: {code_intelligence}
   primary_design_surface: desktop
   stack: [go, node]
   file_structure: "cmd/ for binaries, internal/ for private packages"
   is_monorepo: { value: "true", source: detected, confidence: medium }
   test_runner: { value: "go test ./...", source: detected, confidence: high }
   naming_style: { value: "MixedCaps identifiers (PascalCase exported, camelCase unexported)", source: detected, confidence: high }
   architectural_style: { value: "go-standard-layout", source: manual, confidence: high }
   design_fingerprint:
     i18n: { value: "custom (locale files)", source: detected, confidence: medium }
     css_approach: { value: "tailwind", source: detected, confidence: high }
     ci_cd: { value: "github-actions", source: detected, confidence: high }
     lint: { value: "golangci-lint", source: detected, confidence: high }
   ```

   The formatting is fixed so that an unchanged re-run is byte-identical: two-space indent; `stack` as a flow list of bare ids (`[]` when empty); every `value` and `file_structure` in double quotes; `source`, `confidence`, `memory.provider`, `code_intelligence`, and `primary_design_surface` bare; fingerprint concerns in pack order (`i18n`, `css_approach`, `orm`, `state_management`, `ci_cd`, `lint`), a `none` value included so a later run can tell "probed, absent" from "never probed"; `design_fingerprint` omitted when it has no concern; no timestamp and no comment written by this step; one trailing newline.

5. **Write.** When the rendered bytes equal the existing file, write nothing and report the config unchanged. Otherwise write `.asdt/config.yaml`, creating `.asdt/` when needed.

6. **Self-check.** Re-read the file: it parses, and its keys follow the canonical order. Re-reading what this step wrote is a read, and it is required. Record any mismatch in `Open items` and report the file as written anyway: the bytes on disk are authoritative, so never halt or roll back here.

## Output
Return this summary (protocol §6), with `Open items` and these role sections:

```markdown
# Init: init

{Config written, unchanged, or halted and why, in one sentence.}

## Files changed
- .asdt/config.yaml: {written | unchanged | not written}

## Settings preserved
- {key}: {value} (manual)

## Applied answers
- {field}: {value} ({user | accepted | default})

## Open items
```
