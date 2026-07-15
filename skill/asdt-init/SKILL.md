---
name: asdt-init
description: "Sets up the ground ASDT stands on — initializes .asdt/config.yaml and wires the memory provider so every other specialist has somewhere to read from and write to."
user-invocable: true
specialist-id: asdt-init
metadata:
  author: "Lee Palacios (vitualizz)"
  version: "1.0"
---

# ASDT Init

## Role
Initialize ASDT for the current project. Detect the project stack, collect configuration, and write `.asdt/config.yaml`.

## Prerequisites
None — this is the setup step. Run this before any other ASDT specialist.

## Orchestration Plan

**asdt-init is STANDALONE.** It is a user-invocable setup command, deliberately
NOT registered in `skill/SKILL.md` §9.2 routing — no feature request ever routes
to setup, so the meta-orchestrator has nothing to route here (per ADR-016). Do
not "fix" this omission. init has NO complexity tiers and a fixed 4-step flow;
there is no tier→step table.

| Step | File | Execution | Reads | Writes |
|------|------|-----------|-------|--------|
| knowledge-gate | *(inline — no step file)* | inline | *(orchestrator's own tool list)* | *(no artifact — Engram presence gate)* |
| explore | steps/explore.md | subagent | *(raw project tree — `inputs: []`)* | `init/stack-detection` |
| enrichment | *(inline — no step file)* | inline | *(codegraph MCP survey — no artifact)* | *(no artifact — injects `nuance.*` ambiguities into clarify)* |
| clarify | *(inline — no step file)* | inline | `init/stack-detection.ambiguities[]` + enrichment's `nuance.*` | *(no artifact — injects `answers{}` into write's prompt)* |
| write | steps/write.md | subagent | `init/stack-detection` + `### CLARIFY ANSWERS` | `init/write-summary` |

Step names byte-match `workflow.yaml`; when they differ, `workflow.yaml` is
authoritative.

## Orchestration

This is a light flow — but it has one gate only the orchestrator can pass correctly, and everything downstream depends on it.

**Resolve Engram presence yourself, first, before delegating anything** (the `knowledge-gate` step / Step 2's detection). "Does THIS session have Engram's memory tools" is a question about the orchestrator's own tool list — the one the user is actually relying on for every other specialist. A sub-agent has its own tool list (narrower for specialized agent types, full for `general-purpose`); asking it risks a false "absent" when Engram is actually present in the session that matters. It costs nothing to check yourself — you're inspecting your own tools, not running commands or reading files. This gate is undelegatable; it stays inline with you.

- **Absent** → stop right here and tell the user, exactly as Step 2 describes. Nothing downstream can run without it — don't launch a sub-agent only to have it discover the same dead end.
- **Present** → run the 4-step flow:
  - **explore** (subagent, `agent: analyst`) — launch `steps/explore.md` via your native delegation primitive, passing "Engram confirmed present" into its prompt as an established fact. It detects the stack and context, flags `ambiguities[]`, and returns `init/stack-detection`. It writes NO files.
  - **enrichment** (inline — your own context) — run the codegraph-backed survey that surfaces ≤3 structurally central yet non-obvious symbols as skippable `nuance.*` ambiguities (see Step 2.4 below). Positive-evidence-only, degrades gracefully, and NEVER fails.
  - **clarify** (inline — your own context) — resolve `stack-detection.ambiguities[]` plus enrichment's `nuance.*` ambiguities with the human, ONE question at a time, then compose the `### CLARIFY ANSWERS` block (see the clarify contract in the Workflow below).
  - **write** (subagent, `agent: builder`) — launch `steps/write.md`, injecting both `init/stack-detection` and the `### CLARIFY ANSWERS` block. It writes the four `.asdt` files and returns `init/write-summary`.

**PRE-EXPLORE recalibration gate.** Before launching explore, check for an
existing `.asdt/config.yaml`. If it exists, this project was already
initialized — resolve recalibrate-vs-leave WITH THE USER first. Fail fast: never
run detection only to discard it. If the user chooses "leave as-is", stop without
launching explore. Only proceed to explore once the user has chosen to
recalibrate (a fresh setup with no existing config proceeds directly).

Routing work through sub-agents keeps the bash output, file reads, and intermediate reasoning out of your main context — which is the whole point of the specialist model.

## Workflow

### Step 1 — Detect project stack *(in `explore`)*

The marker-scan mechanics — the bounded `fd`/`find` pipeline, its exclusions and
deterministic ordering, the result cap, the marker→language mapping table, and
the `detected_stack` / primary-language / `{lang_root}` derivation — live in
**`steps/explore.md`**. The explore sub-agent runs them and returns
`init/stack-detection`. Do not duplicate the pipeline here.

### Step 2 — Detect the memory provider

**Detect Engram yourself — do not ask the user to confirm something you can observe directly.** Check your own current tool list for Engram's memory tools (`mem_save`, `mem_search`, `mem_context`, etc. — Claude Code exposes them prefixed as `mcp__plugin_engram_engram__mem_*`; other host assistants may expose the same tools under a different prefix or none).

- If they're present → Engram is installed and reachable. Tell the user so and continue.
- If they're absent → tell the user Engram is required for ASDT's cross-session memory and is not reachable in this session, explain how to install/connect it, and STOP. Do not write `.asdt/config.yaml` with `provider: engram` when the provider isn't actually present — that would silently point every future specialist at a memory backend that doesn't exist.

### Step 2.4 — Enrichment *(inline — your own context, between explore and clarify)*

Enrichment surfaces the handful of symbols a newcomer would misread — structurally
central to the codebase yet non-obvious from their name alone — and turns each into
one skippable `nuance.*` question for the human. It runs inline because, like
clarify, it depends on the orchestrator's own tool list (codegraph is an MCP tool)
and feeds its output into the same `### CLARIFY ANSWERS` block clarify composes.
It is **positive-evidence-only**: if nothing structurally interesting surfaces, it
emits no questions. It NEVER fails and NEVER blocks — every ambiguity it emits is
skippable.

Resolve the capability LADDER against your OWN tool list, top rung first:

1. **Rung 1 — codegraph present** (any `mcp__codegraph__*` tool is in your tool
   list). Call `mcp__codegraph__codegraph_explore` with a natural-language survey —
   "which symbols are the most central to this codebase yet the least
   self-explanatory from their name?" From the grouped results, pick **≤3 chunks**
   that are structurally central yet non-obvious. Heuristic: prefer the highest
   combined caller+callee degree; skip trivially-named getters, setters, and plain
   DTOs (their name already tells the whole story). Deterministic tie-break: order
   by file path, then by symbol name.
2. **Rung 2 — codegraph absent, tree-sitter CLI present** (probe
   `tree-sitter --version`). Syntax is not centrality — a parse tree cannot tell you
   which symbol matters most — so DO NOT surface anything from this rung. Emit
   **zero** ambiguities and note the degraded rung in `open_items` (e.g.
   `enrichment: codegraph absent, tree-sitter only — centrality unavailable, skipped surfacing`).
3. **Rung 3 — neither present.** Skip enrichment entirely, note it in `open_items`
   (e.g. `enrichment: no code-intelligence tooling — skipped`), and proceed. NEVER
   fail.

**Emission.** For each chunk chosen at Rung 1, emit exactly one `Ambiguity`:

- `field`: `nuance.<slug>` — a short, stable slug derived from the symbol name
  (e.g. `nuance.replaceMarkerRegion`).
- `question`: prose asking the human to note what makes this symbol non-obvious —
  its role, an invariant, a gotcha a newcomer would miss.
- `options`: `[]` (free-form note, never a menu).
- `default`: `""`.
- `skippable`: `true` — ALWAYS. A `nuance.*` ambiguity is NEVER a
  `blocking_open_item`.

Emit no more than three. Nothing structurally interesting → no question at all.

### Step 2.5 — Clarify *(inline — your own context, between explore and write)*

explore returns `init/stack-detection` carrying `ambiguities[]` — one entry per
low/medium-confidence or genuinely ambiguous field. clarify is where you, the
orchestrator, resolve them WITH the human. clarify ALSO consumes the `nuance.*`
ambiguities that Step 2.4 (enrichment) surfaced — they are ordinary skippable
`Ambiguity` entries and flow through this same contract. This step runs inline
because only an inline step can pause for a question; the `write` sub-agent cannot.

The inline contract:

1. Ask explore's `stack-detection.ambiguities[]` FIRST — one question at a time,
   the same one-question-per-field discipline §4.3 uses for recalibration. Offer
   `options` when present; otherwise take a free-form value. THEN ask the ≤3
   `nuance.*` ambiguities from enrichment (all skippable, free-form). Offer a
   skip-all shortcut for the `nuance.*` block — the human may decline the whole set
   in one answer.
2. Collect the answers into `answers{}` (field → value). `nuance.*` answers land in
   `answers{}` like any other field; the `write` step routes them to
   `human_nuance` in `project-context.yaml` (see `steps/write.md`).
3. Compose the `### CLARIFY ANSWERS` block and inject it into write's prompt —
   it is REQUIRED even when there was nothing to ask:

   ```
   ### CLARIFY ANSWERS
   answers: { field: value, ... }     # {} when nothing was asked
   skipped: true|false
   blocking_open_items: []
   ```

4. **Non-interactive harness** (no way to ask the human): SKIP the questions. For
   each ambiguity, if `skippable: true` apply its `default` (recorded later as
   `origin: default`); if `skippable: false` add it to `blocking_open_items[]`.
   Set `skipped: true`.
5. **User abort**: if the user cancels, do NOT launch write — no files are
   written.

The `### CLARIFY ANSWERS` block is the only thing the inline clarify step
produces; it carries no artifact of its own.

### Step 3 — Write configuration files *(in `write`)*

The file-writing mechanics — the `.asdt/config.yaml` (`memory.provider: engram`)
write, the `platform.yaml` scan and `conventions.file_structure` derivation, the
`platform-summary.yaml` derived FROM `platform.yaml`, and the
`project-context.yaml` build from `stack-detection.fields` + applied answers —
live in **`steps/write.md`**, along with the idempotency check, the
`source: manual` preservation rule, and the halt contract. The write sub-agent
owns all four file writes. Do not duplicate the write mechanics here.

### Step 4 — Detect project context

Produce `.asdt/knowledge/project-context.yaml` — a machine-written file that records _how_ the project is structured and coded (monorepo shape, test runner, naming style, architectural pattern). This is separate from `platform.yaml`, which records _what_ is installed.

#### 4.1 Check for existing project-context.yaml

Look for `{root}/knowledge/project-context.yaml`:

- **Absent** → fresh detection path (§4.2).
- **Present** → recalibration path (§4.3).

#### 4.2 Fresh detection *(in `explore`)*

The four probes — `is_monorepo`, `test_runner`, `naming_style`,
`architectural_style` — each with its bounded command and exact mapping table,
plus the "one bounded command, first matching row wins, no model judgment" rule
and the per-field `FieldValue` shape, live in **`steps/explore.md`** (§4 probes).
explore detects every field, attaches `source`/`confidence`, and emits an
`Ambiguity` for any low/medium-confidence field. The `write` sub-agent applies
the clarify answers on top and writes `project-context.yaml`; explore itself
writes nothing. Per ADR-013, all probes run as bounded shell commands with no
dependency on this repo's Go code.

Alongside them, seven **design-fingerprint packs** — `i18n`, `css_approach`,
`orm`, `state_management`, `ci_cd`, `lint`, `code_intelligence` — fire
conditionally on `detected_stack` (the fire gate and the per-pack command +
mapping tables also live in **`steps/explore.md`**, §5). Their `FieldValue`s land
in `platform.yaml`'s `design_fingerprint`, NOT `project-context.yaml`. A pack
that does not fire emits no key; `code_intelligence` is positive-evidence-only
(absent → omitted, never an `Ambiguity`).

#### 4.3 Recalibration (project-context.yaml already exists)

When `project-context.yaml` already exists:

1. Run fresh detection → produce `NewContext` (same rules as §4.2).
2. Compute a delta table:

   | Field | Old value | New value | Changed? |
   |---|---|---|---|
   | is_monorepo | … | … | yes/no |
   | test_runner | … | … | yes/no |
   | naming_style | … | … | yes/no |
   | architectural_style | … | … | yes/no |
   | design_fingerprint.i18n | … | … | yes/no |
   | design_fingerprint.css_approach | … | … | yes/no |
   | design_fingerprint.orm | … | … | yes/no |
   | design_fingerprint.state_management | … | … | yes/no |
   | design_fingerprint.ci_cd | … | … | yes/no |
   | design_fingerprint.lint | … | … | yes/no |
   | design_fingerprint.code_intelligence | … | … | yes/no |

   The `design_fingerprint.*` rows track `platform.yaml` (only for packs that
   fired); include a row only for a concern that is present in either the old or
   the new fingerprint.

3. Present the delta table to the user.
4. Ask ONE question: "Accept all changes, or review field by field?"
5. If "accept all" → overwrite `project-context.yaml` with `NewContext`.
6. If "field by field" → for each changed field, ask the user to accept / reject / set manually. One question per field.
7. **Human answers always win.** Fields where the existing `source=manual` are NEVER silently overwritten — they must appear in the delta table and require explicit user acceptance. This holds for `source=manual` `design_fingerprint.<concern>` entries in `platform.yaml` too.

#### 4.4 Confidence and source rules

| Source | When to assign |
|---|---|
| `detected` | Value determined by a bounded command with direct file evidence |
| `inferred` | Pattern matched without direct file evidence (fallback / best-effort) |
| `manual` | User explicitly set this value during a recalibration review |

| Confidence | Meaning |
|---|---|
| `high` | Strong signal — treat as authoritative convention |
| `medium` | Likely match — confirm before diverging |
| `low` | Weak signal — best-effort guess |

Confidence thresholds are assigned by each probe's algorithm (see §4.2 rules). Do not reassign confidence based on judgment — use the exact rules above.

**Negative-evidence rule**: a value concluded from the *absence* of evidence (e.g. `is_monorepo: "false"` because no workspace marker was found) caps at `confidence=medium`, never `high`. Absence proves the probe found nothing — not that nothing exists. `high` is reserved for direct positive file evidence.

#### 4.5 Output

- `{root}/knowledge/project-context.yaml` written (fresh) or confirmed/updated (recalibration).
- Orchestrator receives a `DetectionSummary` for display to the user.
- Proceed to Step 6.

### Step 6 — Confirm
Tell the user:
- Configuration written to `.asdt/config.yaml`
- Detected stack and platform info written to `.asdt/knowledge/`
- Project context written to `.asdt/knowledge/project-context.yaml`
- They can now use `/asdt-architect`, `/asdt-developer`, etc.
