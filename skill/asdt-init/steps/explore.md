# Explore — Init Specialist

## Purpose
Detect the project's stack, context fields, and design fingerprint with bounded, deterministic probes, and turn every low-confidence or genuinely ambiguous field into a question for `clarify`. The same project must yield the same detection in any session, with any model.

## Inputs
None. The project tree is the only source, and no `### INPUT` line arrives.

## Processing

Read-only: run the probes below and return what they found. Every probe is one bounded shell command with an exact mapping table: first matching row wins, no model judgment. A probe that errors is non-fatal: record that field as `unknown` / `inferred` / `low` and continue.

**Never guess.** No marker matched means an empty stack, never one inferred from incidental files.

### Detect project stack

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
- **Language root `{lang_root}`** — for each language, the directory containing its first marker. The context probes below run against these roots, not blindly against the repo root — a marker in `backend/` means that language's evidence lives in `backend/`.

If nothing matches, record an empty stack — do not guess a stack from other files.

### Detect project context

The tables are uniform across languages: a new language adds rows, never logic. No row matches → `unknown` / `inferred` / `low`.

Inputs from the stack scan: `detected_stack`, the primary language, and each language's `{lang_root}`. Probes that inspect language evidence run against the primary language's `{lang_root}` — not blindly against the repo root.

**Probe: `is_monorepo`**

Check workspace markers at the repo root with one compound command:

```
ls go.work pnpm-workspace.yaml 2>/dev/null; grep -l '^\[workspace\]' Cargo.toml 2>/dev/null
```

| Evidence (first match wins) | Value | Source | Confidence |
|---|---|---|---|
| Any workspace marker found (`go.work`, `pnpm-workspace.yaml`, `[workspace]` in root `Cargo.toml`) | `true` | detected | high |
| The stack scan found markers for ≥ 2 distinct languages in ≥ 2 distinct directories at depth ≤ 2 (reuse its output — never rescan) | `true` | detected | medium |
| Neither | `false` | detected | medium *(negative evidence — see Confidence rules)* |

**Probe: `test_runner`** — each check is one `test -f` or one `grep -q` against the primary language's `{lang_root}`; for `scripts.test`, read `{lang_root}/package.json` and take the value of `.scripts.test`:

| Lang | Check (in table order) | Value | Confidence |
|---|---|---|---|
| go | `{lang_root}/Makefile` contains `go test` | `make test` | medium |
| go | `{lang_root}/go.mod` exists | `go test ./...` | high |
| node | `{lang_root}/package.json` `.scripts.test` is non-empty AND does not contain `no test specified` (npm's placeholder) | that script string | high |
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
| go | `.go` | `^(func\|type\|var\|const) [A-Za-z]` | `^(func\|type\|var\|const) [a-z][A-Za-z0-9]*_` | `MixedCaps identifiers (PascalCase exported, camelCase unexported)` |
| node | `.ts` `.tsx` | `^export (function\|class\|const\|interface\|type) [A-Za-z]` | `^export ((function\|const\|let) [a-z][a-z0-9]*_\|(class\|interface\|type) [a-z])` | `camelCase functions and values, PascalCase classes and types` |
| python | `.py` | `^(def [a-z_]\|class [A-Z])` | `^def [A-Z]` | `snake_case functions, PascalCase classes` |
| ruby | `.rb` | `^( *def [a-z_]\|class [A-Z]\|module [A-Z])` | `^ *def [A-Z]` | `snake_case methods, PascalCase classes` |
| rust | `.rs` | `^(pub )?(fn [a-z_]\|struct [A-Z]\|enum [A-Z])` | `^(pub )?fn [A-Z]` | `snake_case functions, PascalCase types` |

No source files sampled → `unknown`/inferred/low.

The violation regexes flag what breaks the language's OWN idiom — an underscore in a
Go or TS identifier, a capitalised Python `def` — never the lowercase names the idiom
itself prescribes. A Go `func helper()` or a TS `export function parseDate()` conforms.

**Probe: `architectural_style`** — list top-level directories with one command (`fd -d 1 -t d . {dir}` or `find {dir} -maxdepth 1 -type d`), first at the repo root; if no row matches there and the primary `{lang_root}` differs from the root, evaluate the same table once more at `{lang_root}`:

| Layout evidence (first match wins) | Value | Source | Confidence | `file_structure` |
|---|---|---|---|---|
| `cmd/` AND `internal/` present | `go-standard-layout` | detected | high | `cmd/ for binaries, internal/ for private packages` |
| `src/` containing `controllers/`, `models/`, `views/` | `mvc` | detected | high | `src/ split into controllers/, models/, views/` |
| `src/` containing `features/` or `modules/` | `modular` | detected | medium | `src/ organized by feature modules` |
| `src/` (no sub-pattern matched) | `layered` | detected | medium | `source under src/` |
| `lib/` present (no `src/`) | `layered` | detected | medium | `source under lib/` |
| No match at root nor at `{lang_root}` | `unknown` | inferred | low | *(none — omit the field)* |

The matching row also fixes `file_structure`, a plain value with no source or confidence, prefixed with `{lang_root}/` when the row matched at `{lang_root}` rather than the root. It never produces an ambiguity.

`go-standard-layout` names the `cmd/` + `internal/` directory convention only — it says
nothing about hexagonal, layered, or any other architecture inside those packages.

When the command output spans libraries from ≥ 2 distinct rows above, emit the first-matching row's value but cap `confidence` at `medium` (conflicting architectural styles). Here "rows" are the five layout predicates in table order — count how many DISTINCT predicates the same directory listing satisfies; the count is a mechanical row-id count over the one command's output, never a judgment call.

### Detect design fingerprint

The fingerprint records _how_ the codebase is built: its i18n, CSS, ORM, state management, CI/CD, lint, and code-intelligence tooling. Each concern is a **pack** with the same discipline as the probes above, and each fired pack emits one value with its source and confidence.

**Fire gate.** Decide which packs run from the stack scan's `detected_stack`
and, for `code_intelligence`, the code-intelligence binding in your executor
header — a pure lookup, no file reads. A pack that does not fire emits **NO key
and NO ambiguity** (it is silent, never `none`). `{node_root}`, `{go_root}`, and
`{py_root}` are the stack scan's `{lang_root}`s for `node`, `go`, and `python`.

| Pack | Fires when | Scope |
|---|---|---|
| `i18n` | `node` ∈ `detected_stack` | `{node_root}` |
| `css_approach` | `node` ∈ `detected_stack` | `{node_root}` |
| `orm` | `node`, `go`, or `python` ∈ `detected_stack` | first such lang in `detected_stack` order (first-lang-wins) |
| `state_management` | `node` ∈ `detected_stack` | `{node_root}` |
| `ci_cd` | always | repo root |
| `lint` | `node`, `go`, or `python` ∈ `detected_stack` | first such lang in `detected_stack` order (first-lang-wins) |
| `code_intelligence` | the binding names a provider | repo root |

**Firing outcome rules** (uniform across packs):

- Fired, command matches a row → that row's value / `detected` / that row's confidence.
- Fired, command runs clean with no match → the table's terminal `none` / `detected` / `medium` row.
- Fired, command errors → `unknown` / `inferred` / `low` (non-fatal — record and continue).
- `code_intelligence` is positive-evidence-only: absent → emit **no key** (never `none`, never `unknown`, never an ambiguity).
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

When the command output spans libraries from ≥ 2 distinct rows above, emit the first-matching row's value but cap `confidence` at `medium` (conflicting css approaches). CSS-modules probe: `fd -d 3 -g '*.module.css' -E node_modules {node_root} | head -1` (find fallback: `find {node_root} -maxdepth 3 -name '*.module.css' -not -path '*/node_modules/*' | head -1`).

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

When the command output spans libraries from ≥ 2 distinct rows above, emit the first-matching row's value but cap `confidence` at `medium` (conflicting ORMs). Distinct rows are counted WITHIN the winning language's rows only — only the winning language's command ran, so rows of a losing language never count toward distinctness.

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

**Pack: `code_intelligence`** *(fires: the binding names a provider — repo root; modeled on `is_monorepo`, positive-evidence-only)*. Probe the index marker the binding names:

```
ls -d {index marker} 2>/dev/null | head -1
```

| Evidence | Value | Source | Confidence |
|---|---|---|---|
| the index marker present | the config value the binding names | detected | high |
| absent | *(emit NO key — never `none`, never `unknown`, never an ambiguity)* |

### Confidence and source rules

| Source | When to assign |
|---|---|
| `detected` | Value determined by a bounded command with direct file evidence |
| `inferred` | Fallback when no row matched or the probe errored |

(`manual` marks a value the human set; only `write` assigns it, never this step.)

| Confidence | Meaning |
|---|---|
| `high` | Strong signal — treat as authoritative convention |
| `medium` | Likely match — confirm before diverging |
| `low` | Weak signal — best-effort guess |

**Negative-evidence rule**: a value concluded from the *absence* of evidence (e.g. `is_monorepo: "false"` because no workspace marker was found) caps at `confidence=medium`, never `high`. Absence proves the probe found nothing — not that nothing exists. `high` is reserved for direct positive file evidence.

### Emit ambiguities

Emit ONE ambiguity per field that is `medium` or `low` confidence, or genuinely ambiguous (two stacks detected and the primary unclear). Each is a question `clarify` asks the human; never resolve one here. Its default is the detected value. It is skippable unless no defensible default exists.

**Fingerprint ambiguities.** Every FIRED pack with a `medium` or `low` result, a clean `none` included, emits one ambiguity named `design_fingerprint.{concern}` (e.g. `design_fingerprint.css_approach`), always skippable. A pack that did not fire emits none, and `code_intelligence` never does.


## Output
Return this summary (`asdt-core/references/summary-schema.md`), filling `Open items` and the two sections below, which `write` and `clarify` read, so keep their shape exact:

```markdown
# Init detection: init

{The stack in one sentence.}

## Detection
stack: {language ids in scan order, comma-separated; empty when no marker matched}
lang_roots: {lang}={root}, ...
file_structure: {the architectural_style row's file_structure; line omitted when none}

| key | value | source | confidence |
|---|---|---|---|
| is_monorepo | ... | ... | ... |
| test_runner | ... | ... | ... |
| naming_style | ... | ... | ... |
| architectural_style | ... | ... | ... |
| design_fingerprint.{concern} | ... | ... | ... |
| code_intelligence | {config value} | detected | high |

## Ambiguities
1. `{field}`: {the question} | options: {a / b, or "free-form"} | default: {value} | confidence: {medium|low} | {skippable|blocking}

## Open items
```

The table carries one `design_fingerprint.{concern}` row per FIRED pack in pack order (`i18n`, `css_approach`, `orm`, `state_management`, `ci_cd`, `lint`), a `none` result included, and the `code_intelligence` row only when the index marker exists. Write `none` under `## Ambiguities` when there are none.
