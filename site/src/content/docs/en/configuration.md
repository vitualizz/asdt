---
title: Configuration
description: Reference for .asdt/config.yaml, the one file that holds your project's facts.
order: 2
locale: en
---

# Configuration

## Project config — `.asdt/config.yaml`

`.asdt/config.yaml` is the one place your project's facts live. Running `/asdt-init` (inside your assistant) detects them with bounded scans of the repo, settles the ambiguous ones with you, and writes this file. Every specialist reads it: the inline `platform-analysis` step turns it into a short platform block (stack, conventions, design fingerprint, design surface), so no specialist re-detects the stack.

```yaml
memory:
  provider: engram
strict_tdd: false              # optional — never added by /asdt-init
code_intelligence: codegraph   # only when detected
primary_design_surface: mobile # only when answered
stack: [typescript]
file_structure: "src/ for application code, tests/ alongside"
is_monorepo: { value: "false", source: detected, confidence: high }
test_runner: { value: "vitest", source: detected, confidence: high }
naming_style: { value: "kebab-case", source: detected, confidence: medium }
architectural_style: { value: "layered", source: manual, confidence: high }
design_fingerprint:
  css_approach: { value: "tailwind", source: detected, confidence: high }
  ci_cd: { value: "github-actions", source: detected, confidence: high }
```

Treat the detected fields as owned by `/asdt-init`: re-run it to recalibrate rather than editing them by hand. A recalibration shows you a table of every value that would change, plus every value already marked `source: manual`, and asks once whether to accept all or review key by key. A value you set that differs from the one offered is recorded as `source: manual` and never silently overwritten. Re-confirming a value never marks it manual: picking the detected value records the detected entry, and accepting a value you had set by hand keeps it as it is. For a key you set by hand earlier, that value, not the detected one, is the default `/asdt-init` offers. Keys you add by hand that `/asdt-init` doesn't recognize are carried forward untouched, and a re-run that changes nothing leaves the file byte-identical.

### `memory.provider`

The memory provider your install is bound to. Only `engram` is supported today. `/asdt-init` writes the value of the provider your install is bound to, even over an existing value. When no provider is bound, it leaves any existing `memory` block as it is. If the project names a different provider than your install, specialists treat memory as unavailable: they still run, working from the session scratchpad, and nothing is recalled or remembered until you re-run `/asdt-init`. See [Scratchpad & Memory](/asdt/docs/memory-and-engram).

### `strict_tdd`

When `true`, the Developer's `implement` step writes tests in the same pass as the code, inside the same approved file list. There is no separate test step, and `implement` never runs the tests: the `verify` gate that follows offers you the check commands (never one that would rewrite files) and runs them only if you say yes. You can also ask for tests in the request without setting the flag.

Default: absent, which means `false`. `/asdt-init` carries it forward but never adds it.

### `code_intelligence`

Written by `/asdt-init` only when it detects a code-intelligence index (today, `codegraph`). Specialists then prefer it over grep-and-read loops. Never written as `none` — when nothing is detected, the key is removed.

### `primary_design_surface`

The surface UX/UI designs for first: `mobile`, `tablet`, `desktop`, or `none`. `/asdt-init` asks for it once (default `mobile`); a re-init lets you change it.

- **Absent** — the question was never asked (for example, a non-interactive run). The specialists that design or build UI (UX/UI and the Developer) treat it as `mobile` and record one `ASSUMED:` entry saying so.
- **`none`** — the project has no visual surface at all (a CLI, a library, a backend service). UX/UI then skips its design steps and writes a short summary saying there is nothing to design, so downstream specialists read an explicit "no UI" instead of a silent gap. The Developer adds no responsive layer.

### Detected facts

`stack`, `file_structure`, the four decision fields (`is_monorepo`, `test_runner`, `naming_style`, `architectural_style`), and the `design_fingerprint` concerns (`i18n`, `css_approach`, `orm`, `state_management`, `ci_cd`, `lint`) are what `/asdt-init` detects. The decision fields and fingerprint concerns are stored as `{value, source, confidence}`; specialists read only the `value`. A fingerprint concern probed and found absent is stored as `none`, so a later run can tell "probed, absent" from "never probed".

## What is not in the config

**Things you say about the project** — a convention, a gotcha a newcomer would miss — are not configuration. They go to memory as notes, never into a file: `/asdt-init` asks a few optional questions about non-obvious code and remembers your answers, and any specialist remembers what you tell it ("remember that…"). Notes reach a step as separate, labelled project notes, never mixed into the detected platform block. See [Scratchpad & Memory](/asdt/docs/memory-and-engram).

**Working files** — each specialist's summary — live in the session scratchpad, outside your repository, never under `.asdt/`.

### Migrating from `.asdt/knowledge/`

Earlier versions of `/asdt-init` also wrote `.asdt/knowledge/knowledge.yaml` and `.asdt/knowledge/provenance.yaml`. Both are retired. A project is **migrated** once `.asdt/config.yaml` has a top-level `stack` key, which only the current `/asdt-init` writes. Until then, nothing is lost: each specialist says in one line that it found the legacy store and reads it as a fallback — the facts `.asdt/config.yaml` lacks (the config wins where both have a key), and the notes you wrote there, carried as project notes. Once migrated, nothing but `/asdt-init` reads the legacy files: a key the config omits stays absent, and a legacy note counts only once memory holds it. If the files are still there, each specialist says instead: "Legacy ASDT files kept in `.asdt/knowledge/`; nothing reads them now. Delete them via `/asdt-init`." When `/asdt-init` finds them, it treats the project as already initialized and offers a recalibration. On a yes, it folds their detected values (including any you had set by hand) into `.asdt/config.yaml`. If `knowledge.yaml` holds notes you wrote earlier and memory is connected, it asks which of them to remember. It then asks whether to delete the legacy files once the config is written. The default is no, and it only offers the deletion when it can't lose a note you wanted kept. On a yes, and only once the config is written and checks out, `/asdt-init`'s `write` step removes exactly `.asdt/knowledge/knowledge.yaml` and `.asdt/knowledge/provenance.yaml` (`rm`), then the emptied `.asdt/knowledge/` directory (`rmdir`, which refuses a directory that still holds anything else). It runs no other command that changes anything. Without memory connected, init skips the note questions and leaves the notes in their file. Once the config is written nothing reads them, until a later `/asdt-init` with memory connected offers to remember them.

## Environment variables

ASDT reads no environment variables directly. The binary uses the AI assistant's environment (Claude Code or OpenCode) for all runtime configuration.

## Multiple projects

Each project has its own `.asdt/` directory. Switch projects by opening your AI assistant in a different directory — ASDT reads configuration from the nearest `.asdt/config.yaml` walking up from the current working directory. Memory is kept apart per project too: notes are scoped by the project's directory name plus a hash of its git remote (see [Scratchpad & Memory](/asdt/docs/memory-and-engram#the-scratchpad-working-summaries)).
