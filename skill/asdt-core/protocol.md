# ASDT Protocol

The one shared skill every run loads. It defines what gets persisted, how inputs arrive, what a step executor may do, and what a hand-off looks like.

## 1. Engram contract

**Persist hand-offs only.** One key per role per change: `{project}/{change}/{role}/handoff`. Written with `mem_save`, title `"{change}/{role}/handoff"`, type `"decision"`. Everything a run produces on the way there — exploration, drafts, intermediate analysis — lives in the orchestrator's context and dies with the run. If it does not cross a specialist boundary, it is not saved.

**Load at start.** ONE `mem_search("{project}/{change}")` to list what exists, then `mem_get_observation(id)` for the `*/handoff` records this specialist declares it consumes — nothing else. Budget: 2–3 MCP calls per run.

**Organizational memory.** When a run closes, and only if it made a decision that is not obvious from the code, append ONE line to topic_key `{project}/journal`:

```
{role}@{change}: {what} — {why}
```

One line, one append, no envelope, no second save. Nothing else goes to the journal.

**Degradation.** An expected hand-off that does not exist is recorded in `open_items` with the literal prefix `ASSUMED:` and the run proceeds. A missing input never blocks and never fails a run.

## 2. Intake contract

Declared inputs arrive ALREADY INJECTED in the sub-agent prompt as `### INPUT {topic_key}` blocks, or as `### INPUT {topic_key}: UNRESOLVED`. Every declared input either arrived as a block or it did not — there is no third state, and a sub-agent NEVER fetches its own declared inputs. That work already happened, against a store the sub-agent may not even be able to see.

**One batched clarification turn.** A run gets AT MOST ONE. If gaps are genuinely blocking — no defensible hand-off is possible without an answer — collect every such question across the whole run, ask them TOGETHER as one numbered list, and stop once. Never one round trip per question, never a second turn. This turn is FULLY SUPPRESSED when the prompt carries a `## Tailored Workflow` block: the router already ran its clarifying gates and its answers are settled.

**Harden always.** Every non-blocking gap degrades into an `open_items` entry prefixed `ASSUMED:` — what was assumed, and what would confirm or refute it — and the run continues. A stalled run returns nothing; a hardened run returns a hand-off whose weak spots are named and checkable. When in doubt between asking and assuming: assume, mark it, keep moving.

## 3. Executor rules

> You are the sub-agent for this single step. Do the work and return. Do NOT delegate, do NOT orchestrate, do NOT run other steps. Do NOT fetch your inputs — they arrive injected. An `UNRESOLVED` input means record the gap and proceed, never abort.

**Role boundary.** If your step is NOT `developer/implement` (or `developer/test` under strict TDD), you write ZERO files in the user's repo. Your outputs go to Engram via `mem_save`, full stop. ASDT's own state lives only under `.asdt/`.

**Verifiable evidence** — exact file paths, symbol names, commands, observed values instead of "should" or "likely" — is required ONLY of steps that read the codebase. Steps that do not touch code do not carry this requirement.

## 4. Injection format (orchestrator)

The orchestrator resolves each declared input ONCE per run and injects it. Resolved:

```
### INPUT {topic_key}
{full content}
```

Failed to resolve:

```
### INPUT {topic_key}: UNRESOLVED
(could not be fetched — record it in open_items and proceed)
```

**Partial failure is not total failure**: resolved inputs are injected and used normally; only the failed one becomes an `UNRESOLVED` block. Declared reference skills are read by the orchestrator and injected as `### REFERENCE SKILL {path}` blocks in the same prompt.

## 5. Hand-off schema

Roles omit keys that do not apply. Roles never add process keys.

```yaml
payload:
  what: ""                  # one sentence
  decisions: []             # one-line imperatives, rejected alternative in parentheses
  constraints: []
  files_hint: []            # code anchors: where to look first
  acceptance_criteria: []   # Given/When/Then, max 5 — PM is the authority on ACs
  risks: []                 # {risk, mitigation}, one line each
  data_model: []            # architect/developer only, when applicable
  api_surface: []           # architect/developer only, when applicable
  open_items: []            # real gaps only, ASSUMED: prefix
```

**Golden rule: if a field does not change what the consumer types, it does not belong in the hand-off.**
