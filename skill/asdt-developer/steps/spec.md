# Spec — Developer Specialist

## Purpose
Define what gets built and how: scope, acceptance criteria, technical approach, and the files
the implementation is allowed to touch.

## Inputs
- `### INPUT developer/explore` — OPTIONAL. The exploration: `Open questions` (answer them
  here), `Constraints` (the patterns to follow), `Files hint` (each line's reach or `reuse`
  label), `Risks`
- `### INPUT pm` — OPTIONAL, the acceptance-criteria authority. Extract ONLY its `Acceptance
  criteria`
- `### INPUT architect` — OPTIONAL. Extract `Decisions`, `Constraints`, `Data model`, `API
  surface`, `Files hint` — when it arrived the approach is ALREADY DECIDED, do not re-open it
- `### INPUT ux-ui` — OPTIONAL. UX/UI's summary. Its `Surface` becomes a `Constraints` entry.
  Its `Flows` (with their branches, states, and copy), `Screens`, `Component designs`, and
  `Accessibility` are copied verbatim, because they are the design the UI code must implement;
  every `Component designs` entry is a file to create or modify, or an `Open items` entry when
  its location is unclear. A `Visual direction` is copied verbatim, and the files that implement
  it as the project's first tokens join the edit targets. Each `Open items` entry starting
  `Proposed token:` is copied verbatim into `Constraints`, and the token file it names joins the
  files to modify. A summary that says the project declares no visual surface carries no design
- `### INPUT security` — OPTIONAL. Extract each `Risks` mitigation and the `Constraints` (the
  hardening checklist) — each mitigation that lands in this change's files is a `Constraints`
  entry
- `### INPUT qa` — OPTIONAL. QA's summary. Extract its `Verdict`, `Edge cases`, and `Test cases`.
  After a `no-go` (either form), its reason is a `Constraints` entry, each `Edge cases` entry
  this change closes is a `Scope` In item (the rest are Out), and each `Test cases` row that
  exercises an In item is copied verbatim. Any other verdict contributes nothing
- `### INPUT developer` — OPTIONAL. The summary this plan replaces. Extract only a delivery, and
  only when it holds one: at `Stage: implemented`, its `**Verification:**` line, `Files changed`,
  and any `Prior delivery`; at `Stage: spec`, the `Prior delivery` an earlier plan already
  carried forward
- `### INPUT feedback` and `### INPUT previous` — injected only when the human chose adjust at
  the approval gate: `previous` is the plan they were shown, `feedback` their words, verbatim

Each arrives as a path to read — never look for anything else. The platform block is in this
prompt: extract its `Design surface:` line's bare value (strip any parenthesised marker; `none
— no visual surface` reads as `none`), and whether its `Tests:` line carries `strict TDD`.

**DEGRADATION**: `pm` UNRESOLVED → author the acceptance criteria from the exploration and the request, and record `ASSUMED:` in `Open items`.
**DEGRADATION**: `architect` UNRESOLVED → choose the approach here, per step 4.
**DEGRADATION**: `ux-ui`, `security`, `qa`, or `developer` UNRESOLVED → proceed without it; no entry.
**DEGRADATION**: `developer/explore` UNRESOLVED → on a revision, revise `previous` instead of re-deriving the plan; otherwise read the affected code yourself. No entry.
**DEGRADATION**: no `Design surface:` line → skip the surface check in step 5; no entry.
**DEGRADATION — overrides every "no entry" above**: when the request or `feedback` refers to an
earlier plan, design, or finding (an approved plan, "the design", "what QA found") and the input
that would carry it is UNRESOLVED, record `ASSUMED: upstream {role} summary not available this
session — {what you planned without it}` in `Open items`, one per such role.

## Processing

1. Answer each of the exploration's open questions.
2. Define the scope boundary: what IS included and what is explicitly NOT. What the `not
   needed` rung names is an Out item.
3. Write at most 5 acceptance criteria, numbered: `AC1: Given …, when …, then …`. Product's
   summary is the AC AUTHORITY: when it arrived, REFINE its criteria into this form, preserving
   each one's intent — never re-derive an independent set.
4. **Approach — the smallest that fully works.** Walk the `conventions.md` ladder against the
   exploration's `reuse` hints; stop at the first rung that satisfies every AC.
   - `Approach:` is the FIRST `Decisions` line, always: what it reuses and the larger option
     rejected (`rejected: none — no larger option` for a trivial change). With the Architect's
     summary it RESTATES that decision at implementation granularity; a genuine conflict is an
     `Open items` entry.
   - `New dependency:` → the manifest joins the edit targets, and `run {install command} —
     lockfile not updated` joins `Open items`: nothing in this chain installs.
   - `shortcut:` for a known limit the approach accepts; `never-cut: {item} — {why}` in `Open
     items` for a never-cut item (`conventions.md`) it leaves out.
5. **Declare the edit targets** — real paths relative to the project root, one per line: the
   action, the backticked path, then what it serves (`- modify: \`{path}\` — AC1, AC3`) —
   `implement` writes nothing else. Fold in every path of the Architect's `Files hint` this change creates or
   modifies. Close the list both ways:
   - every AC is named by a target, and every target names an AC, a `Test cases` ID, its
     reach label (`reach: caller`), or the decision or input that mandates it (`New
     dependency`, `Proposed token`, `Visual direction`, `Component designs`, Architect `Files
     hint`) — never cut those; a target with none is scope creep: cut it (a plan-only `- none`
     list is exempt);
   - every exploration hint labelled `change`, `caller`, `test`, `fixture`, `config`, or
     `export` is a target, an `Open items` entry, or left untouched in one `Scope` Out line
     with its reason (`Out: callers unchanged — \`{path}\`, \`{path}\`; props unchanged`).

   Mark every test-file target `(test)` (`- create: \`{path}\` — AC2 (test)`), an existing test
   the change breaks included (`reach: test`). Add test targets where new logic earns a test
   (`conventions.md`, Tests), for each `Test cases` row you copied, when the request asks for
   tests, and, under strict TDD, for every unit the change adds or alters — no others. A row
   whose `Level` has no framework in the project (the `Tests:` line, the manifest) is no
   target: `no {Level} framework: {ID} — not tested` in `Open items`. A `feedback` that drops
   the tests under strict TDD wins: record `strict TDD waived by the human — no tests` in
   `Open items`.

   An empty list means plan only: write `- none: plan only, nothing is written`. A path under
   `.asdt/` is never a target, whoever asked — ASDT's own state changes only through
   `/asdt-init`: record `.asdt/ change: {path} — {what was wanted}; run /asdt-init` in `Open
   items`. Under `Design surface: none`, each UI target — a view, a component, a stylesheet or
   token file — gets `surface conflict: {path} is UI, but the project declares no visual
   surface` in `Open items`, so the human sees it at approval.
6. **On a revision** (`feedback` arrived), steps 1–5 start from `previous`, not from scratch:
   the human's words OVERRIDE it wherever the two disagree, and everything they did not touch
   stays as it was. Step 5 still runs in full — an adjustment to scope moves the edit
   targets with it.

Do NOT write implementation code here.

**A new iteration never loses the delivery.** When the summary this plan replaces held one,
copy it verbatim into `Prior delivery`: it stays part of the change.

## Output
Return this summary (`asdt-core/protocol.md` §6). The orchestrator files it as the role
summary and as the payload `implement` builds from. Omit every section that does not apply;
the design sections, `Test cases`, and `Prior delivery` are copied verbatim, never rewritten.

```markdown
# Developer: {change}

{The change this plan delivers, in one sentence.}

**Stage:** spec

## Decisions
- Approach: {what, naming the reused path or package} (rejected: {the larger option})
- New dependency: {package} — {why no smaller rung works}
- shortcut: {the limit} — upgrade when {trigger}
- {Open question} → {answer}

## Constraints
- {What implementation must respect: the surface, each mitigation, each proposed token}

## Acceptance criteria
- AC1: Given {context}, when {action}, then {observable outcome}

## Risks
- {Risk} → {mitigation}

## Data model
## API surface

## Scope
- In: {item}
- Out: {item}

## Approach
{Only for a data-model, API, or migration change: that design, at the depth it earns.}

## Migration notes
- {Phase or step}

## Edit targets
- create: `{path}` — AC1
- modify: `{path}` — AC2, AC3
- modify: `{path}` — reach: {caller | test | fixture | config | export}
- create: `{path}` — AC2 (test)

## Flows
## Screens
## Component designs
## Accessibility
## Visual direction
## Test cases

## Prior delivery
{The replaced summary's `**Verification:**` line, then the file entries of its `Files changed` and its own `Prior delivery`, verbatim.}

## Open items
- ASSUMED: {what was assumed} — {what would confirm or refute it}
```
