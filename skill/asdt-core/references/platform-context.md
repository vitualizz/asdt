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
Design surface: {primary_design_surface}
```

`Conventions` joins its two parts with ` | ` only when both are present; with one, emit it alone; with neither, drop the line. The `Tooling` line is emitted with exactly this wording and only when `code_intelligence` is present in `.asdt/config.yaml`.

**`Design surface` is the one line never dropped.** It reads `primary_design_surface` from `.asdt/config.yaml`, not from `knowledge.yaml`: `mobile`, `tablet`, or `desktop` is emitted as-is; an absent key — the question was never asked — emits `Design surface: mobile (default — never asked)`; `none` emits `Design surface: none — no visual surface`. What a step does with it lives in that step's own file.

Go-only repo, where the Node packs never fire:

```
Stack: Go
Conventions: cmd/ for binaries, internal/ for private packages
CI/CD: github-actions
Lint: golangci-lint
Tooling: codegraph index available — prefer codegraph over grep/read loops
Design surface: none — no visual surface
```

Treat detected conventions as authoritative and user-declared ones as untouchable without explicit approval. This block carries DETECTED values only, with one exception: `Design surface` is the human's own answer to `/asdt-init`, never a scan. Any other note a person wrote about the project never joins it, and travels as its own labelled line instead. See **Human nuance** below.

## Degradation

If `knowledge.yaml` is absent, do not halt. Record one `open_items` entry —

```
ASSUMED: knowledge.yaml absent — conventions inferred from visible code patterns
```

— and proceed using the conventions visible in the code at hand (file naming, import style, directory layout). If the file exists but is partially populated, inject the fields that are present and say nothing about the missing ones. Either way the `Design surface` line is still emitted — it never depended on this file.

## Human nuance

If `knowledge.yaml` carries a `human_nuance:` list, read those entries directly, one by one, as user-authored notes about the topic each one names — the thing a newcomer would misread and no scan can detect.

They are **intentionally NOT auto-injected**. They carry no confidence rating and no provenance, so they never join the compact block above: folding a person's note in among detected values makes it look like something the scan found.

When the orchestrator prepares a step's prompt, the entries RELEVANT to that step's area travel as their own labelled lines — `### PROJECT NOTE (human): legacy CSS in styles/; new work uses Tailwind utilities` — one per line, never merged into the block above and never rewritten into the voice of a detected fact.

Relevance is a judgment: a note about CSS reaches an implement step touching styles, and does not reach a PM backlog. Entries that bear on nothing in this step are simply not carried.
