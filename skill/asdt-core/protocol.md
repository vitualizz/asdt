# ASDT Protocol

The one shared skill every run loads. It defines what gets persisted, how inputs arrive, what a step executor may do, and what a hand-off looks like.

## 1. Engram contract

**Persist hand-offs only.** One key per role per change: `{project}/{change}/{role}/handoff`. Written with `mem_save`, title `"{change}/{role}/handoff"`, type `"decision"`. Everything a run produces on the way there — exploration, drafts, intermediate analysis — lives in the orchestrator's context and dies with the run. If it does not cross a specialist boundary, it is not saved.

**Two intents, one contract.** When the run DELIVERS a change, the key is `{project}/{change}/{role}/handoff`. When it EXAMINES what already exists — an audit, a review, an assessment with nothing to deliver — the key is `{project}/study/{topic}/{role}`, where `{topic}` is derived from the request in short, stable kebab-case ("audit the payments module" → `payments-module`). The specialist judges which one this is from the invocation; it never asks the user to pick, and genuine ambiguity means a change.

Everything else is identical: same schema, same load rules, same degradation. In a study, `decisions[]` carries the judgments and `risks[]` what was found — the schema does not grow a study variant. A past study is organizational memory: later runs meet it through the `knowledge-recall` prelude, never as a declared input.

A step whose `workflow.yaml` entry declares `output: context` instead of `output_topic_key` persists NOTHING: its payload stays in the orchestrator's context and is injected into the next step as an `### INPUT {step-name}` block. A step may declare `context_inputs:` — payloads of earlier `output: context` steps that the orchestrator injects as `### INPUT {name}` blocks.

**Load at start.** ONE `mem_search("{project}/{change}")` to list what exists, then `mem_get_observation(id)` for the `*/handoff` records this specialist declares it consumes — nothing else. Budget: 2–3 MCP calls per run.

**Organizational memory.** When a run closes, and only if it made a decision that is not obvious from the code, append ONE line to topic_key `{project}/journal`:

```
{role}@{change}: {what} — {why}
```

One line, one append, no envelope, no second save. Nothing else goes to the journal.

**Terrain or history.** Will this still be true in three months, whatever the current change? That is TERRAIN, and it belongs in the `human_nuance:` list of `.asdt/knowledge/knowledge.yaml`. Is it something decided or found while working on this change or study? That is HISTORY, and it belongs in the journal line above. Convenience never decides this — the question does.

**Consent.** Nothing enters `human_nuance` without the user's instruction or confirmation. An explicit instruction — "remember that…", "save this", "for the future:" — IS the consent. A durable fact you noticed on your own is PROPOSED in one line and written only on their yes.

**Bounded write.** Once consent exists, the ORCHESTRATOR edits the `human_nuance:` list and nothing else in that file — never a sub-agent, one entry per line, plain language. The rest of `knowledge.yaml` belongs to `/asdt-init`. "Forget the thing about X" removes the matching entry, on the same confirmation. Before adding, read what is already there: a note that contradicts an existing one UPDATES it instead of stacking beside it. Never a secret, a token, or a credential — if asked, decline in one line and say where that belongs instead.

**Degradation.** An expected hand-off that does not exist is recorded in `open_items` with the literal prefix `ASSUMED:` and the run proceeds. A missing input never blocks and never fails a run.

## 2. Intake contract

Declared inputs arrive ALREADY INJECTED in the sub-agent prompt as `### INPUT {topic_key}` blocks, or as `### INPUT {topic_key}: UNRESOLVED`. Every declared input either arrived as a block or it did not — there is no third state, and a sub-agent NEVER fetches its own declared inputs. That work already happened, against a store the sub-agent may not even be able to see.

**One batched clarification turn.** A run gets AT MOST ONE. If gaps are genuinely blocking — no defensible hand-off is possible without an answer — collect every such question across the whole run, ask them TOGETHER as one numbered list, and stop once. Never one round trip per question, never a second turn. If the invocation already carries the router's proposal or otherwise answers your doubts, do not re-ask what is settled. A sharpened invocation answers what the user settled at routing time, not that nothing else is left to ask — raise what the router could not see, once, batched with everything else this run needs. When in doubt between asking and assuming: assume, mark it `ASSUMED:`, keep moving.

**Harden always.** Every non-blocking gap degrades into an `open_items` entry prefixed `ASSUMED:` — what was assumed, and what would confirm or refute it — and the run continues. A stalled run returns nothing; a hardened run returns a hand-off whose weak spots are named and checkable. When in doubt between asking and assuming: assume, mark it, keep moving.

## 3. Step execution rules

> Whoever executes this step — a launched sub-agent, or the orchestrator running it inline — is bound by everything in this section. Do the work of this ONE step and return. Do NOT delegate, do NOT run other steps. Do NOT fetch your inputs — they arrive injected. An `UNRESOLVED` input means record the gap and proceed, never abort.

**Write boundary.** Exactly two steps in ASDT write files, and the step's identity decides it, never the identity of whoever runs it: `developer/implement` writes host source inside the edit roots its spec declares, and `asdt-init/write` writes ASDT's own state under `.asdt/`. Every other step writes NOTHING, anywhere — its only output is `mem_save`. If you are running any other step and reach for Edit or Write, STOP before the write — you have left the plan — and recover inside this same step: name the step that must be delegated instead, record the blocked work in `open_items`, and finish this step normally with a hand-off. STOP scopes to the write, never to the run.

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
