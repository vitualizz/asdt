# Platform Context — Shared Skill

## Purpose

Inject the project's detected platform knowledge into any specialist's context. This skill is consumed by every specialist at the start of their workflow to ground code generation, component suggestions, and design decisions in the project's actual conventions. It also serves the inline `platform-analysis` workflow step declared by several specialists — that step points here.

---

## Reuse Guard (project-level)

Before any analysis, check if `.asdt/knowledge/platform-summary.yaml` exists.

If it does, **do not re-analyze**. Read that file and inject its contents as context, as-is. The `/asdt-init` command produces this file deterministically via bounded scan probes; re-deriving it with the LLM wastes tokens and yields non-deterministic stack interpretations.

Only if `.asdt/knowledge/platform-summary.yaml` is **absent**, fall back to finding and reading `.asdt/knowledge/platform.yaml` (next section) and injecting the extracted fields described below.

---

## How to Find platform.yaml

> Fallback path — used only when the Reuse Guard above found no `platform-summary.yaml`.

Walk up from CWD until you find `.asdt/knowledge/platform.yaml`. This is the same nearest-ancestor search used for `.asdt/` itself.

Path: `.asdt/knowledge/platform.yaml` relative to the resolved ASDT root.

---

## What to Inject

When `platform.yaml` is found, extract and inject the following fields into the specialist's context (summarized to under 500 tokens):

| Field | What to inject |
|---|---|
| `detected_stack` | Languages, frameworks, runtimes detected (e.g. "Go 1.22, no frontend framework") |
| `conventions.naming` | Naming conventions in use (e.g. "snake_case for files, PascalCase for exported Go types") |
| `conventions.file_structure` | Directory layout pattern (e.g. "internal/ for packages, cmd/ for binaries") |
| `design_fingerprint` | Per-concern tooling map — inject each present concern via the lines below; omit concerns whose value is empty |
| `design_fingerprint.i18n` | Internationalization library in use — if present |
| `design_fingerprint.css_approach` | CSS approach in use — if present |
| `design_fingerprint.state_management` | State-management library in use — if present |
| `design_fingerprint.orm` | ORM / data-access layer in use — if present |
| `design_fingerprint.ci_cd` | CI/CD platform in use — if present |
| `design_fingerprint.lint` | Linter / formatter in use — if present |
| `design_fingerprint.code_intelligence` | Code-intelligence index present — drives the Tooling line below; if present |

Discard: full file listings, raw config, `layout_patterns`, `scanned_at`, `schema_version`. If the extracted content exceeds 500 tokens, summarize each field to its single most important fact.

Do not inject the entire `platform.yaml` verbatim. Summarize only the fields relevant to the specialist's current step.

---

## Graceful Degradation

If `platform.yaml` does not exist:

1. Do NOT halt the specialist workflow.
2. Record the absence in the artifact's `open_items[]`:
   ```
   "platform.yaml absent — conventions inferred from visible code patterns"
   ```
3. Proceed using conventions inferred from the code visible in the current context (file naming, import style, directory layout).

If `platform.yaml` exists but is partially populated (some fields empty or missing), inject only the fields that are present. Do not halt or record an error for missing optional fields.

---

## Injection Format

Build the injection from only the fields that are actually present — omit a line entirely when its source field is empty or missing. Never emit a label with nothing after it (`Architecture: ` on its own conveys nothing and still costs tokens):

```
Stack: {detected_stack values, comma-separated}
Conventions: {naming style, if present}{ | file structure note, if present}
i18n: {design_fingerprint.i18n, if present}
CSS: {design_fingerprint.css_approach, if present}
State: {design_fingerprint.state_management, if present}
ORM: {design_fingerprint.orm, if present}
CI/CD: {design_fingerprint.ci_cd, if present}
Lint: {design_fingerprint.lint, if present}
Tooling: codegraph index available — prefer codegraph over grep/read loops
```

`Conventions` joins its two parts with ` | ` only when BOTH are present. If only one is present, emit that one alone with no separator. If neither is present, omit the `Conventions` line too. Each `design_fingerprint` concern line is emitted only when its value is present — omit the line entirely otherwise. The `Tooling` line is emitted with EXACTLY this wording, and ONLY when `design_fingerprint.code_intelligence` is present; omit it when there is no code-intelligence index.

Fully-populated example:
```
Stack: Node (TypeScript), React
Conventions: PascalCase exported symbols | src/features/ modular layout
i18n: i18next
CSS: tailwind
State: redux
ORM: prisma
CI/CD: github-actions
Lint: eslint
Tooling: codegraph index available — prefer codegraph over grep/read loops
```

Partially-populated example — e.g. a `platform.yaml` from `/asdt-init` for a Go-only repo, where the node-only packs (`i18n`, `css_approach`, `state_management`) never fire and only the always-on and Go packs contribute:
```
Stack: Go
Conventions: cmd/ for binaries, internal/ for private packages
CI/CD: github-actions
Lint: golangci-lint
Tooling: codegraph index available — prefer codegraph over grep/read loops
```

---

## Conditional: Project Context

> Load this section only if `.asdt/knowledge/project-context.yaml` exists.
> If the file is absent, skip this entire section silently.
> If `schema_version` != `"1"`, skip and note `project-context.yaml: schema_version mismatch, skipped` in open_items.

The following fields describe the structural and stylistic context of this project,
as detected by `/asdt-init`. Each field carries a `source` (detected | inferred | manual)
and a `confidence` (high | medium | low). Fields with an empty `value` are omitted.

**Monorepo**: {{ is_monorepo.value }}  *({{ is_monorepo.source }}, {{ is_monorepo.confidence }})*
**Test runner**: {{ test_runner.value }}  *({{ test_runner.source }}, {{ test_runner.confidence }})*
**Naming style**: {{ naming_style.value }}  *({{ naming_style.source }}, {{ naming_style.confidence }})*
**Architectural style**: {{ architectural_style.value }}  *({{ architectural_style.source }}, {{ architectural_style.confidence }})*

When writing code or tests, treat `detected/high` fields as authoritative conventions.
Treat `inferred/medium` fields as likely conventions — confirm before diverging.
Treat `manual` fields as user-declared — never override without explicit user approval.

If a `human_nuance:` list is present, read each entry directly here as a user-authored note about that topic — it is intentionally NOT auto-injected (source: manual, origin: user, no confidence rating).

---

## Usage Note

This skill is referenced by specialists; it is injected per step via `reference_skills` in each specialist's `workflow.yaml`. The specialist does not need to duplicate this logic — it calls this skill at its Platform Analysis step.
