# Platform Context — Reference

Grounds a run in the project's actual conventions instead of generic defaults. The orchestrator runs it inline as the `platform-analysis` prelude and injects one compact block into each step's prompt (`asdt-core/protocol.md` §5).

## Reuse Guard

Project facts come from `.asdt/config.yaml` alone (protocol §0). Find `.asdt/` by walking up from the working directory to the nearest ancestor that holds it. **When the config exists, do NOT re-analyze the project.** `/asdt-init` detected its facts with bounded probes. Re-deriving them with an LLM costs tokens and gives a different answer every time.

Read these keys, values only, and nothing else in the file. A key stored as `{value, source, confidence}` contributes its `value` alone, and that value counts as absent when it is `unknown` or `none`:

| Key | Feeds |
|---|---|
| `stack` | `Stack` |
| `naming_style`, `file_structure` | `Conventions` |
| `architectural_style` | `Architecture` |
| `is_monorepo` | `Monorepo` |
| `test_runner` | `Tests` |
| `design_fingerprint.{i18n, css_approach, state_management, orm, ci_cd, lint}` | one line each |
| `code_intelligence` | `Tooling` |
| `primary_design_surface` | `Design surface` |

Never inject a value's source or confidence, scan timestamps, raw config, or full file listings. Keep the whole block under ~500 tokens. If it runs longer, cut each line to its single most important fact.

## Injection Format

Build the block from the keys actually present, and **omit any line whose value is empty or absent**: a label with nothing after it conveys nothing and still costs tokens.

```
Stack: {stack, comma-separated}
Conventions: {naming_style}{ | file_structure}
Architecture: {architectural_style}
Monorepo: yes
Tests: {test_runner}
i18n: {design_fingerprint.i18n}
CSS: {design_fingerprint.css_approach}
State: {design_fingerprint.state_management}
ORM: {design_fingerprint.orm}
CI/CD: {design_fingerprint.ci_cd}
Lint: {design_fingerprint.lint}
Tooling: {code_intelligence} index available — prefer code intelligence over grep/read loops
Design surface: {primary_design_surface}
```

`Conventions` joins its two parts with ` | ` only when both are present. With one part, emit it alone; with neither, drop the line. `Monorepo: yes` is emitted only when `is_monorepo` is true. The `Tooling` line is emitted with exactly this wording, and only when `code_intelligence` is present: it is how a step learns this project is indexed (`asdt-core/executor-header.md`, Code intelligence).

**`Design surface` is the one line that is never dropped.** `mobile`, `tablet`, or `desktop` is emitted as-is. An absent key (the question was never asked) emits `Design surface: mobile (default — never asked)`. `none` emits `Design surface: none — no visual surface`. What a step does with this line lives in that step's own file.

A Go-only repository looks like this:

```
Stack: Go
Conventions: cmd/ for binaries, internal/ for private packages
Tests: go test
CI/CD: github-actions
Lint: golangci-lint
Design surface: none — no visual surface
```

Treat detected values as authoritative. `Design surface` is the one exception: it is the human's own answer to `/asdt-init`, not a scan. A note a person wrote about the project never joins this block (see **Project notes** below).

## Degradation

If `.asdt/config.yaml` does not exist, do not halt. Record one `Open items` entry:

```
ASSUMED: .asdt/config.yaml absent — conventions inferred from visible code patterns
```

Then proceed with the conventions visible in the code at hand: file naming, import style, directory layout. When the config exists but lacks a key, inject the lines that are present and say nothing about the rest. Either way, the `Design surface` line is still emitted, with its default when the key is absent.

## Project notes

What the user has said about the project (the thing a newcomer would misread and no scan can detect) does not live in this file. It comes from memory, through the `knowledge-recall` prelude, and reaches a step as its own `### PROJECT NOTE:` lines (`asdt-core/references/knowledge-recall.md`). Never fold a note into the block above, and never rewrite one in the voice of a detected fact: placing a person's note among detected values makes it look like something the scan found.
