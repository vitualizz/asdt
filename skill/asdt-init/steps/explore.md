# Explore — Init Specialist

## Purpose
Detect the project stack and structural context deterministically, and surface
every low-confidence or genuinely ambiguous field as a question for the clarify
step. This is init's first sub-agent step — it READS the project, it never
writes. The same project must produce the same `stack-detection` no matter which
model or session runs it.

## Inputs
- `inputs: []` — there are no upstream artifacts; the raw project tree is the
  only source.
- **Engram presence arrives as an established fact.** The orchestrator passed
  the gate (`knowledge-gate`) before launching you. Do NOT re-verify Engram's
  tool list — that is the orchestrator's job and it already did it.

## Context budget
Raw request / launch context: max ~500 tokens. Everything else this step needs
it discovers by running bounded shell commands against the project tree — it
does not pull large inputs into context.

## Processing

The analyst is **read-only**: it runs detection commands and returns a
structured result. It writes NO files here — the `write` step owns all
filesystem writes.

**explore NEVER guesses.** No markers matched → `detected_stack=[]`. A probe
that errors is non-fatal: record `value="unknown"`, `source=inferred`,
`confidence=low` for that field and continue. Never infer a stack from
incidental files.

### Step 1 — Detect project stack

Run ONE bounded command scanning for stack marker files down to depth 3 — do not
eyeball a directory listing or infer from visible files. The result must be
identical no matter which model or session runs it. The exclusion list and the
result cap are part of the contract: they keep the output size independent of
repo size.

```
fd -d 3 -t f -H '^(go\.mod|package\.json|Cargo\.toml|pyproject\.toml|requirements\.txt|Gemfile)$' . \
  -E node_modules -E .git -E vendor -E dist -E build -E .venv -E target \
  | awk -F/ '{printf "%d %s\n", NF, $0}' | LC_ALL=C sort -k1,1n -k2 | cut -d' ' -f2- | head -20
```

If `fd` is not available, run the equivalent `find` fallback — same exclusions, same ordering pipeline:

```
find . -maxdepth 3 \( -name node_modules -o -name .git -o -name vendor -o -name dist -o -name build -o -name .venv -o -name target \) -prune \
  -o -type f \( -name go.mod -o -name package.json -o -name Cargo.toml -o -name pyproject.toml -o -name requirements.txt -o -name Gemfile \) -print \
  | awk -F/ '{printf "%d %s\n", NF, $0}' | LC_ALL=C sort -k1,1n -k2 | cut -d' ' -f2- | head -20
```

The `awk | sort | cut` pipeline orders results deterministically: shallowest path first, then lexicographic (`LC_ALL=C`). The cap is the first 20 marker paths.

Map each marker to a canonical language id — no judgment calls:

| Marker file | Language id |
|---|---|
| `go.mod` | `go` |
| `package.json` | `node` |
| `Cargo.toml` | `rust` |
| `pyproject.toml` or `requirements.txt` | `python` |
| `Gemfile` | `ruby` |

From the ordered marker list, derive — mechanically, first occurrence wins:

- **`detected_stack`** — the language ids in scan order, deduplicated. A project can match more than one stack (e.g. a Python `backend/` with a Node `frontend/`).
- **Primary language** — the first entry of `detected_stack`.
- **Language root `{lang_root}`** — for each language, the directory containing its first marker. The §4 probes run against these roots, not blindly against the repo root — a marker in `backend/` means that language's evidence lives in `backend/`.

If nothing matches, record an empty stack — do not guess a stack from other files.

### Step 4 — Detect project context

Every probe has the same shape: **one bounded command, one exact mapping table,
first matching row wins, no model judgment**. The tables are uniform across
languages — supporting a new language means adding rows, never adding logic. If
no row matches, the field is `value="unknown"`, `source=inferred`,
`confidence=low`. A probe error is non-fatal — record the fallback for that
field and continue.

Inputs from Step 1: `detected_stack`, the primary language, and each language's `{lang_root}`. Probes that inspect language evidence run against the primary language's `{lang_root}` — not blindly against the repo root.

**Probe: `is_monorepo`**

Check workspace markers at the repo root with one compound command:

```
ls go.work pnpm-workspace.yaml 2>/dev/null; grep -l '^\[workspace\]' Cargo.toml 2>/dev/null
```

| Evidence (first match wins) | Value | Source | Confidence |
|---|---|---|---|
| Any workspace marker found (`go.work`, `pnpm-workspace.yaml`, `[workspace]` in root `Cargo.toml`) | `true` | detected | high |
| Step 1 found markers for ≥ 2 distinct languages in ≥ 2 distinct directories at depth ≤ 2 (reuse Step 1's output — never rescan) | `true` | detected | medium |
| Neither | `false` | detected | medium *(negative evidence — see Confidence rules)* |

**Probe: `test_runner`** — each check is one `test -f` or one `grep -q` against the primary language's `{lang_root}`; for `scripts.test`, read `{lang_root}/package.json` and take the value of `.scripts.test`:

| Lang | Check (in table order) | Value | Confidence |
|---|---|---|---|
| go | `{lang_root}/Makefile` contains `go test` | `make test` | medium |
| go | `{lang_root}/go.mod` exists | `go test ./...` | high |
| node | `package.json` `.scripts.test` is non-empty | that script string | high |
| node | `{lang_root}/jest.config.js` exists | `jest` | medium |
| node | `{lang_root}/vitest.config.ts` exists | `vitest` | medium |
| python | `{lang_root}/pytest.ini` exists OR `pyproject.toml` contains `[tool.pytest` | `pytest` | high |
| python | `pytest` appears in `pyproject.toml` or `requirements.txt` | `pytest` | medium |
| ruby | `{lang_root}/Gemfile` contains `rspec` | `bundle exec rspec` | high |
| ruby | `{lang_root}/.rspec` exists | `bundle exec rspec` | medium |
| rust | `{lang_root}/Cargo.toml` exists | `cargo test` | high |

All table matches are `source=detected`.

**Probe: `naming_style`** — sample up to 8 source files (≤ 64 KB each) under the primary language's `{lang_root}`, depth ≤ 3, deterministic order:

```
fd -d 3 -t f -S -64k {extension flags from table} . {lang_root} | LC_ALL=C sort | head -8
```

(`find` fallback: `-maxdepth 3 -type f -size -65536c` with the same `-name` patterns, then the same `sort | head` pipeline.)

A file **conforms** when it matches the positive regex (`grep -qE`) and does not match the violation regex. Conforming ratio maps to confidence: ≥ 75% high, 50–74% medium, < 50% → `unknown`/inferred/low. At ≥ 50% the table value is emitted with `source=detected`.

| Lang | Extensions | Positive regex | Violation regex | Value when dominant |
|---|---|---|---|---|
| go | `.go` | `^(func\|type\|var\|const) [A-Z]` | `^(func\|type\|var\|const) [a-z]` | `snake_case filenames, PascalCase exported symbols` |
| node | `.ts` `.tsx` | `^export (function\|class\|const\|interface) [A-Z]` | `^export (function\|const) [a-z]` | `PascalCase exported symbols` |
| python | `.py` | `^(def [a-z_]\|class [A-Z])` | `^def [A-Z]` | `snake_case functions, PascalCase classes` |
| ruby | `.rb` | `^( *def [a-z_]\|class [A-Z]\|module [A-Z])` | `^ *def [A-Z]` | `snake_case methods, PascalCase classes` |
| rust | `.rs` | `^(pub )?(fn [a-z_]\|struct [A-Z]\|enum [A-Z])` | `^(pub )?fn [A-Z]` | `snake_case functions, PascalCase types` |

No source files sampled → `unknown`/inferred/low.

**Probe: `architectural_style`** — list top-level directories with one command (`fd -d 1 -t d . {dir}` or `find {dir} -maxdepth 1 -type d`), first at the repo root; if no row matches there and the primary `{lang_root}` differs from the root, evaluate the same table once more at `{lang_root}`:

| Layout evidence (first match wins) | Value | Source | Confidence |
|---|---|---|---|
| `cmd/` AND `internal/` present | `hexagonal` | detected | high |
| `src/` containing `controllers/`, `models/`, `views/` | `mvc` | detected | high |
| `src/` containing `features/` or `modules/` | `modular` | detected | medium |
| `src/` (no sub-pattern matched) | `layered` | detected | medium |
| `lib/` present (no `src/`) | `layered` | detected | medium |
| No match at root nor at `{lang_root}` | `unknown` | inferred | low |

### Step 5 — Detect design fingerprint

The `design_fingerprint` records _how_ the codebase is built — its i18n, CSS,
ORM, state-management, CI/CD, lint, and code-intelligence tooling. Each concern
is a **pack** with the same discipline as the §4 probes: one bounded command,
one exact mapping table, first matching row wins, no model judgment. Fired packs
land in `platform.yaml` (not `project-context.yaml`); each is a `FieldValue`.

**Fire gate.** Decide which packs run from Step 1's `detected_stack` alone — a
pure lookup, no file reads. A pack that does not fire emits **NO key and NO
ambiguity** (it is silent, never `none`). `{node_root}`, `{go_root}`, and
`{py_root}` are the `{lang_root}`s from Step 1 for `node`, `go`, and `python`.

| Pack | Fires when | Scope |
|---|---|---|
| `i18n` | `node` ∈ `detected_stack` | `{node_root}` |
| `css_approach` | `node` ∈ `detected_stack` | `{node_root}` |
| `orm` | `node`, `go`, or `python` ∈ `detected_stack` | first such lang in `detected_stack` order (first-lang-wins) |
| `state_management` | `node` ∈ `detected_stack` | `{node_root}` |
| `ci_cd` | always | repo root |
| `lint` | `node`, `go`, or `python` ∈ `detected_stack` | first such lang in `detected_stack` order (first-lang-wins) |
| `code_intelligence` | always | repo root |

**Firing outcome rules** (uniform across packs):

- Fired, command matches a row → that row's value / `detected` / that row's confidence.
- Fired, command runs clean with no match → the table's terminal `none` / `detected` / `medium` row.
- Fired, command errors → `unknown` / `inferred` / `low` (non-fatal — record and continue).
- `code_intelligence` is positive-evidence-only: absent → emit **no key** (never `none`, never `unknown`, never an Ambiguity).
- Did not fire → no key, no ambiguity.

**Pack: `i18n`** *(fires: node)* — grep the node manifest dependency block:

```
grep -oE '"(i18next|react-i18next|next-i18next|vue-i18n|react-intl|@formatjs/[a-z-]+|@lingui/[a-z-]+|svelte-i18n|@nuxtjs/i18n)"' {node_root}/package.json | tr -d '"' | LC_ALL=C sort -u | head -10
```

| Evidence (first match wins) | Value | Source | Confidence |
|---|---|---|---|
| `i18next`, `react-i18next`, or `next-i18next` | `i18next` | detected | high |
| `vue-i18n` or `@nuxtjs/i18n` | `vue-i18n` | detected | high |
| `react-intl` or any `@formatjs/*` | `react-intl (formatjs)` | detected | high |
| any `@lingui/*` | `lingui` | detected | high |
| `svelte-i18n` | `svelte-i18n` | detected | high |
| no dep, but a `locales` directory exists (dual fd/find, depth ≤ 3) | `custom (locale files)` | detected | medium |
| none of the above | `none` | detected | medium |

When the command output spans libraries from ≥ 2 distinct rows above, emit the first-matching row's value but cap `confidence` at `medium` (conflicting i18n libraries). Locale-dir fallback: `fd -d 3 -t d -H '^locales$' -E node_modules {node_root} | head -1` (find fallback: `find {node_root} -maxdepth 3 -type d -name locales -not -path '*/node_modules/*' | head -1`).

**Pack: `css_approach`** *(fires: node)*:

```
grep -oE '"(tailwindcss|styled-components|@emotion/[a-z]+|sass|less|@stitches/[a-z]+|@vanilla-extract/[a-z-]+)"' {node_root}/package.json | tr -d '"' | LC_ALL=C sort -u | head -10
```

| Evidence (first match wins) | Value | Source | Confidence |
|---|---|---|---|
| `tailwindcss` | `tailwind` | detected | high |
| `styled-components` | `styled-components` | detected | high |
| any `@emotion/*` | `emotion` | detected | high |
| any `@vanilla-extract/*` | `vanilla-extract` | detected | high |
| any `@stitches/*` | `stitches` | detected | high |
| `sass` or `less` only (no CSS-in-JS above) | `sass/less preprocessor` | detected | medium |
| no dep, but a `*.module.css` file exists (dual fd/find) | `css-modules` | detected | medium |
| none of the above | `none` | detected | medium |

CSS-modules probe: `fd -d 3 -g '*.module.css' -E node_modules {node_root} | head -1` (find fallback: `find {node_root} -maxdepth 3 -name '*.module.css' -not -path '*/node_modules/*' | head -1`).

**Pack: `orm`** *(fires: node | go | python — evaluate the FIRST of these in `detected_stack` order; first-lang-wins)*:

Run only the winning language's command against its `{lang_root}`; each pipes `| LC_ALL=C sort -u | head -10`.

- node — `grep -oE '"(prisma|@prisma/client|typeorm|drizzle-orm|sequelize|mongoose|@mikro-orm/[a-z-]+|knex)"' {node_root}/package.json | tr -d '"'`
- go — `grep -oE '(gorm\.io/gorm|entgo\.io/ent|github\.com/uptrace/bun|github\.com/jmoiron/sqlx)' {go_root}/go.mod`
- python — `grep -hioE '(sqlalchemy|django|tortoise-orm|peewee)' {py_root}/pyproject.toml {py_root}/requirements.txt 2>/dev/null`

| Lang | Evidence (first match wins) | Value | Confidence |
|---|---|---|---|
| node | `prisma` or `@prisma/client` | `prisma` | high |
| node | `typeorm` | `typeorm` | high |
| node | `drizzle-orm` | `drizzle` | high |
| node | `sequelize` | `sequelize` | high |
| node | `mongoose` | `mongoose` | high |
| node | any `@mikro-orm/*` | `mikro-orm` | high |
| node | `knex` only | `knex (query builder)` | medium |
| go | `gorm.io/gorm` | `gorm` | high |
| go | `entgo.io/ent` | `ent` | high |
| go | `github.com/uptrace/bun` | `bun` | high |
| go | `github.com/jmoiron/sqlx` | `sqlx (query builder)` | medium |
| python | `sqlalchemy` | `sqlalchemy` | high |
| python | `django` | `django orm` | high |
| python | `tortoise-orm` | `tortoise` | high |
| python | `peewee` | `peewee` | high |
| (winning lang) | no match | `none` | medium |

Every row is `source: detected`.

**Pack: `state_management`** *(fires: node)*:

```
grep -oE '"(redux|@reduxjs/toolkit|zustand|jotai|recoil|mobx|valtio|xstate|pinia|@ngrx/store|@tanstack/react-query|react-query|swr)"' {node_root}/package.json | tr -d '"' | LC_ALL=C sort -u | head -10
```

| Evidence (first match wins) | Value | Source | Confidence |
|---|---|---|---|
| `@reduxjs/toolkit` or `redux` | `redux` | detected | high |
| `zustand` | `zustand` | detected | high |
| `jotai` | `jotai` | detected | high |
| `recoil` | `recoil` | detected | high |
| `mobx` | `mobx` | detected | high |
| `valtio` | `valtio` | detected | high |
| `xstate` | `xstate` | detected | high |
| `pinia` | `pinia` | detected | high |
| `@ngrx/store` | `ngrx` | detected | high |
| only `@tanstack/react-query`, `react-query`, or `swr` (no client-state lib above) | `react-query (server state)` | detected | medium |
| none of the above | `none` | detected | medium |

**Pack: `ci_cd`** *(fires: always — repo root)*:

```
{ ls -d .github/workflows 2>/dev/null; ls .gitlab-ci.yml .circleci/config.yml azure-pipelines.yml Jenkinsfile bitbucket-pipelines.yml .drone.yml 2>/dev/null; } | LC_ALL=C sort | head -10
```

| Evidence (first match wins) | Value | Source | Confidence |
|---|---|---|---|
| `.github/workflows` | `github-actions` | detected | high |
| `.gitlab-ci.yml` | `gitlab-ci` | detected | high |
| `.circleci/config.yml` | `circleci` | detected | high |
| `azure-pipelines.yml` | `azure-pipelines` | detected | high |
| `Jenkinsfile` | `jenkins` | detected | high |
| `bitbucket-pipelines.yml` | `bitbucket-pipelines` | detected | high |
| `.drone.yml` | `drone` | detected | high |
| none of the above | `none` | detected | medium |

**Pack: `lint`** *(fires: node | go | python — FIRST of these in `detected_stack` order; first-lang-wins)*:

Config-file checks are fixed-path `ls` at `{lang_root}`; dep checks grep the manifest. Run only the winning language's command; each pipes `| LC_ALL=C sort -u | head -10`.

- node — `{ ls {node_root}/.eslintrc* {node_root}/eslint.config.* {node_root}/biome.json {node_root}/.prettierrc* 2>/dev/null; grep -oE '"(eslint|@biomejs/biome|standard|prettier)"' {node_root}/package.json | tr -d '"'; }`
- go — `ls {go_root}/.golangci.yml {go_root}/.golangci.yaml {go_root}/.golangci.toml 2>/dev/null`
- python — `{ ls {py_root}/ruff.toml {py_root}/.flake8 {py_root}/.pylintrc 2>/dev/null; grep -hioE '(ruff|flake8|pylint)' {py_root}/pyproject.toml {py_root}/requirements.txt 2>/dev/null; }`

| Lang | Evidence (first match wins) | Value | Confidence |
|---|---|---|---|
| node | `.eslintrc*`, `eslint.config.*`, or an `eslint` dep | `eslint` | high |
| node | `@biomejs/biome` dep or `biome.json` | `biome` | high |
| node | `standard` dep | `standard` | high |
| node | only `prettier` / `.prettierrc*` (no linter above) | `prettier (format only)` | medium |
| go | `.golangci.{yml,yaml,toml}` | `golangci-lint` | high |
| go | no golangci config | `none (go vet only)` | medium |
| python | `ruff` dep or `ruff.toml` | `ruff` | high |
| python | `.flake8` or `flake8` dep | `flake8` | high |
| python | `.pylintrc` or `pylint` dep | `pylint` | high |
| (winning lang) | no match | `none` | medium |

Every row is `source: detected`.

**Pack: `code_intelligence`** *(fires: always — repo root; modeled on `is_monorepo`, positive-evidence-only)*:

```
ls -d .codegraph 2>/dev/null | head -1
```

| Evidence | Value | Source | Confidence |
|---|---|---|---|
| `.codegraph` present | `codegraph` | detected | high |
| absent | *(emit NO key — never `none`, never `unknown`, never an Ambiguity)* |

### Confidence and source rules

| Source | When to assign |
|---|---|
| `detected` | Value determined by a bounded command with direct file evidence |
| `inferred` | Pattern matched without direct file evidence (fallback / best-effort) |
| `manual` | User explicitly set this value during a recalibration review (set by clarify/write, never here) |

| Confidence | Meaning |
|---|---|
| `high` | Strong signal — treat as authoritative convention |
| `medium` | Likely match — confirm before diverging |
| `low` | Weak signal — best-effort guess |

**Negative-evidence rule**: a value concluded from the *absence* of evidence (e.g. `is_monorepo: "false"` because no workspace marker was found) caps at `confidence=medium`, never `high`. Absence proves the probe found nothing — not that nothing exists. `high` is reserved for direct positive file evidence.

### Emit ambiguities

Emit ONE `Ambiguity` per field that is low/medium-confidence OR genuinely
ambiguous (e.g. two stacks detected and the primary is unclear). Each ambiguity
is the question the clarify step will ask the human. Do NOT resolve them here —
explore detects and flags; clarify asks; write applies.

**design_fingerprint ambiguities.** For every FIRED pack whose result is
`medium` or `low` confidence — including a clean `none` result — emit exactly one
`Ambiguity` whose `field` is the dotted name `design_fingerprint.<concern>` (e.g.
`design_fingerprint.css_approach`), with `default` set to the detected value and
`skippable: true` ALWAYS. A pack that did not fire emits NO ambiguity, and
`code_intelligence` NEVER produces one (positive-evidence-only). design_fingerprint
ambiguities are never `blocking_open_items`.

## Output
Produces: `init/stack-detection`

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Schema:
```yaml
payload:
  detected_stack: []          # language ids in scan order, deduplicated; [] when no markers matched
  lang_roots:                 # language id -> directory of its first marker
    - lang: ""
      root: ""
  fields:                     # FieldValue per detected context field
    is_monorepo: { value: "", source: "", confidence: "" }
    test_runner: { value: "", source: "", confidence: "" }
    naming_style: { value: "", source: "", confidence: "" }
    architectural_style: { value: "", source: "", confidence: "" }
    design_fingerprint:       # one FieldValue per FIRED pack; OMIT non-firing packs entirely (no key)
      i18n: { value: "", source: "", confidence: "" }
      css_approach: { value: "", source: "", confidence: "" }
      orm: { value: "", source: "", confidence: "" }
      state_management: { value: "", source: "", confidence: "" }
      ci_cd: { value: "", source: "", confidence: "" }
      lint: { value: "", source: "", confidence: "" }
      code_intelligence: { value: "", source: "", confidence: "" }  # present only when .codegraph exists
  ambiguities: []             # one Ambiguity per low/medium-confidence or genuinely ambiguous field; design_fingerprint ambiguities use dotted field names and are always skippable
  open_items: []
```

**FieldValue** (value-object):
```yaml
value: ""        # the detected value, or "unknown"
source: ""       # detected | inferred  (manual is set later by clarify/write, never here)
confidence: ""   # high | medium | low
```

**Ambiguity** (value-object):
```yaml
field: ""        # the field name this question resolves (e.g. "test_runner")
question: ""     # the prose question clarify asks the human, one at a time
options: []      # optional — concrete choices to offer; omit when free-form
default: ""      # the value applied when the harness is non-interactive and skippable=true
skippable: true  # true → SKIP with default when non-interactive; false → open_item when non-interactive
```
