# Implement — Developer Specialist

## Purpose
Write the implementation — and its tests when TDD is on — respecting existing conventions and
never leaving the declared edit roots.

## Inputs
Everything arrives ALREADY INJECTED — never self-fetch.

- `dev-spec` — this run's spec, injected from the orchestrator's context as `### INPUT dev-spec`.
  Extract: `files_to_create`, `files_to_modify`, `scope`, `acceptance_criteria[]`, `approach`,
  `key_constraints`, and — when present — `ux_flows`, `screens`, `component_designs`,
  `a11y_requirements`, `visual_direction`
- `{project}/{change}/developer/handoff` (OPTIONAL) — the persisted spec, injected INSTEAD of
  `dev-spec` when this run resumes a plan saved in an earlier session. It is a spec only when
  it carries `stage: spec`; then extract the same fields. Below, `dev-spec` means whichever of the
  two arrived.
- `{project}/{change}/architect/handoff` (OPTIONAL) — the architectural decision
- `verification-failures` (OPTIONAL) — injected as `### INPUT verification-failures` only on a
  fix round. Extract: `round`, `failures[]` (each failing command and its output), and
  `files_changed` (everything earlier rounds of this run wrote)
- Platform summary — extract the conventions and the `Design surface:` line. Use its bare value
  (`mobile (default — never asked)` → `mobile`); that default marker, or a summary or line that
  never arrived, means `mobile` — plus ONE `ASSUMED:` entry saying so, only when this step
  produces a UI file
- `### HOST SKILL frontend-design` (OPTIONAL) — visual-design craft from the host, injected only
  when it has one

**DEGRADATION**: if neither `dev-spec` nor a `stage: spec` `developer/handoff` arrived, do NOT proceed to writing — there are no declared edit roots, so PLAN-ONLY mode applies; note `ASSUMED:` in open_items. If `developer/handoff` is absent while `dev-spec` arrived, that is the normal fresh run — no entry. If `architect/handoff` is UNRESOLVED, take the design authority from `dev-spec` and note `ASSUMED:` in open_items. If `verification-failures` is absent, this is the first pass — proceed normally. If no host skill arrived, build the UI on the conventions alone — no entry.

## Processing

### Mode (do this FIRST)
1. `allowedEditRoots` = the union of `files_to_create` + `files_to_modify` in `dev-spec` — those
   two lists and nothing else; the human approved exactly them.
2. EMPTY → PLAN-ONLY MODE: emit every file as a `code_snippets[]` entry; write NO host files.
3. NON-EMPTY → WRITING MODE: write real files, ONLY to paths within `allowedEditRoots`.

### Writing the code
Generate every file respecting `key_constraints`, the platform-summary conventions, early-return,
no global state, and small focused functions. In writing mode:
1. Before each write, confirm the target is under a declared root. If not, STOP before writing
   it, never expand scope, and record the path in `unsafe_skipped` and `open_items`.
2. Read the current content of any `files_to_modify` target before editing it (per
   `../asdt-core/references/conventions.md`) — match its conventions, never clobber unrelated code.
3. Write the file via the filesystem write tool, and record its path, action
   (`created`|`modified`), and rationale in `files_changed[]`.
4. Compose `suggested_verification`: the build/lint/test commands the human may run in
   `.commands`, what a healthy run looks like in `.expected` — only commands the `verify`
   exception of `asdt-core/protocol.md` §3 allows, in their non-writing variant where one exists.

### UI files
Applies to every UI file this step produces — a view, a component, a stylesheet or token file —
whether written or emitted as a snippet; other files skip this section. Lay each one out for the
`Design surface` first, then adapt to the others (as the spec's surface constraint says, when
UX/UI ran). Under `Design surface: none`, add no responsive layer and record in `open_items`
that a UI file contradicts the declared surface. When the spec carries a UX/UI design, every UI
file consumes tokens, never literal values: a `visual_direction` in `dev-spec` is built as the project's first tokens; with a design
system, each `Proposed token:` in `key_constraints` is added under its name and value to the
token file it names — within `allowedEditRoots`, like any write — and no other token is
invented. The `screens` and `component_designs` ARE the UI's design: build each screen to its
layout, hierarchy and its type levels, spacing, components, states, extremes, and breakpoints,
and each designed component to its anatomy, props, variants, sizes, states, tokens, and
behavior — never re-design what they specify. A `kind: variant` design adds its option to the
existing component and leaves every existing variant unchanged.

An injected host skill shapes the visual execution only, and only where the spec leaves room:
with a design system, layout and hierarchy within its tokens and components; without one, also
type, palette, composition, and motion. **Precedence, highest first: `allowedEditRoots` and
`scope`; the project's existing design system, or else the dev-spec's `visual_direction`; the
spec's `screens`, `component_designs`, `ux_flows`, and `a11y_requirements`; then the host
skill.** Where the skill pulls against anything above it, the skill yields.

### Tests
Generate tests HERE, in this same step, when `strict_tdd: true` in `.asdt/config.yaml` or the
user asked for them. Otherwise skip this section entirely. Tests obey the SAME mode and the
SAME `allowedEditRoots` as the code above — in plan-only mode they are `code_snippets[]`
entries like any other file; in writing mode they are real files, and a test path outside the
declared roots STOPS exactly as a source path does.

Per unit under test: one happy-path test for the acceptance criterion, and one edge-case test
for the most likely failure mode. Follow the project's existing test framework; use
table-driven cases where the framework offers them; test behavior, never internals.

This step NEVER runs a command that builds, lints, tests, installs, or generates — not even on a
fix round. Running checks is the human's call, at the orchestrator's `verify` gate, which offers
`suggested_verification.commands`.

### Fix round
Only when `verification-failures` arrived. Fix what the failures name and nothing else, in the
SAME mode and within the SAME `allowedEditRoots` resolved from the same `dev-spec` — a fix round
never widens them. Read each failure's output and the files it points at before editing. A fix
that needs a path outside the roots STOPS exactly as in writing mode — report it, never work
around it. Never delete, skip, or loosen a test to turn it green; a test that contradicts an
acceptance criterion is an `open_items` entry, not an edit. `files_changed` covers the whole run: carry forward every entry from
`verification-failures.files_changed`, then add or update your own.

### Coverage
Every acceptance criterion should be answerable with "which code addresses this". After
generating the code, walk `dev-spec.acceptance_criteria[]`: for each one no file addresses,
append a single line to `open_items` — `AC not covered: {ac text}`. Warnings, never a halt.

## Output
Produces: `developer/handoff` — persist via memory **save** under this step's `output_topic_key`,
using the canonical hand-off schema from `asdt-core/protocol.md`, with `stage: implemented`. This
record REPLACES the persisted spec under the same key (`asdt-core/protocol.md` §1), so carry
`scope` and `acceptance_criteria` forward verbatim from `dev-spec` — QA tests against them. Set
`mode` to the resolved value; `files_changed` carries this run's real paths in writing mode —
never a prior delivery `dev-spec` carried, which approving this build replaced —
`code_snippets` the code in plan-only. Return the same payload: the orchestrator's `verify` step
builds on it.

```yaml
payload:
  what: ""                    # what was implemented, one sentence
  stage: implemented
  scope: {in: [], out: []}    # verbatim from dev-spec
  acceptance_criteria: []     # verbatim from dev-spec
  mode: "writing | plan-only"
  allowedEditRoots: []        # resolved list, recorded verbatim (writing mode)
  files_changed:              # writing mode
    - path: ""
      action: "created|modified"
      rationale: ""
  code_snippets:              # plan-only mode
    - file: ""
      language: ""
      content: ""
  unsafe_skipped: []          # paths STOPPED on, and what triggered them
  suggested_verification:     # offered to the human at the verify gate
    commands: []
    expected: ""
  decisions: []               # implementation choices worth carrying forward
  risks: []                   # {risk, mitigation}
  open_items: []              # includes one "AC not covered: ..." line per uncovered AC
```
