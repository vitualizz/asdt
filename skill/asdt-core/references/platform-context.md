# Platform Context — Reference

Grounds a run in the project's actual conventions instead of generic defaults. Optional reference, useful to any role that writes code, suggests components, or makes design decisions.

## Reuse Guard

Before any analysis, check whether `.asdt/knowledge/knowledge.yaml` exists — walking up from CWD, the same nearest-ancestor search used to find `.asdt/` itself.

**If it exists, do NOT re-analyze the project.** Read that file and inject the compact summary below. `/asdt-init` produced it with bounded scan probes; re-deriving it with an LLM costs tokens and yields a different answer every time.

Read values only. Never inject `source`, `confidence`, `scanned_at`, raw config, full file listings, or the `provenance.yaml` sidecar. Keep the whole injection under ~500 tokens; if it runs longer, cut each field to its single most important fact.

## Injection Format

Build the block from the fields actually present and **omit any line whose source value is empty** — a label with nothing after it conveys nothing and still costs tokens.

```
Stack: {stack values, comma-separated}
Conventions: {naming style}{ | file structure note}
i18n: {design_fingerprint.i18n}
CSS: {design_fingerprint.css_approach}
State: {design_fingerprint.state_management}
ORM: {design_fingerprint.orm}
CI/CD: {design_fingerprint.ci_cd}
Lint: {design_fingerprint.lint}
Tooling: codegraph index available — prefer codegraph over grep/read loops
```

`Conventions` joins its two parts with ` | ` only when both are present; with one, emit it alone; with neither, drop the line. The `Tooling` line is emitted with exactly this wording and only when `code_intelligence` is present in `.asdt/config.yaml`.

Go-only repo, where the Node packs never fire:

```
Stack: Go
Conventions: cmd/ for binaries, internal/ for private packages
CI/CD: github-actions
Lint: golangci-lint
Tooling: codegraph index available — prefer codegraph over grep/read loops
```

Treat detected conventions as authoritative and user-declared ones as untouchable without explicit approval. If `knowledge.yaml` carries a `human_nuance:` list, read those entries directly — they are user-authored notes, not detected values.

## Degradation

If `knowledge.yaml` is absent, do not halt. Record one `open_items` entry —

```
ASSUMED: knowledge.yaml absent — conventions inferred from visible code patterns
```

— and proceed using the conventions visible in the code at hand (file naming, import style, directory layout). If the file exists but is partially populated, inject the fields that are present and say nothing about the missing ones.
