# ASDT Protocol

The one shared skill every run loads. It defines what gets persisted, how inputs arrive, what a step executor may do, and what a hand-off looks like.

## 1. Engram contract

**Persist hand-offs only.** One key per role per change: `{project}/{change}/{role}/handoff`, where `{change}` is derived from the request in short, stable kebab-case ("add password reset" → `password-reset`), exactly as `{topic}` is below. Written with `mem_save`, title `"{change}/{role}/handoff"`, type `"decision"`. Everything a run produces on the way there — exploration, drafts, intermediate analysis — lives in the orchestrator's context and dies with the run. If it does not cross a specialist boundary or need to survive into a later run, it is not saved.

**Staged hand-offs.** A role that delivers in stages — a plan the human approves, then the change built from it — writes every stage to its ONE key; it never opens a second key per stage. The record carries `stage` naming which one it holds, and each stage's `mem_save` replaces the previous record whole, so a later stage re-emits whatever of the earlier one its consumers still need. A plan persisted this way survives an interrupted or plan-only run and can be resumed. A staged record with no `stage` — written before stages existed — reads as the last stage — the Developer's `implemented` — and is never a resume candidate.

**Two intents, one contract.** When the run DELIVERS a change, the key is `{project}/{change}/{role}/handoff`. When it EXAMINES what already exists — an audit, a review, an assessment with nothing to deliver — the key is `{project}/study/{topic}/{role}`, where `{topic}` is derived from the request in short, stable kebab-case ("audit the payments module" → `payments-module`). The specialist judges which one this is from the invocation; it never asks the user to pick, and genuine ambiguity means a change.

Everything else is identical: same schema, same load rules, same degradation. In a study, `decisions[]` carries the judgments and `risks[]` what was found — the schema does not grow a study variant. A past study is organizational memory: later runs meet it through the `knowledge-recall` prelude, never as a declared input.

**Intra-run payloads.** A payload reaches a later step of the SAME run through the orchestrator's context, never through Engram. It has a name — the one its producer declares under `Produces:` (`dev-spec`, `security-assessment`): the step file's, or, for an inline step with no step file, the SKILL.md section that is its contract — and the consuming step lists that name in `context_inputs:`; the orchestrator injects each one as an `### INPUT {name}` block. The name in `context_inputs:`, the `Produces:` line, and the injected heading are always the same string. A step whose `workflow.yaml` entry declares `output: context` instead of `output_topic_key` persists NOTHING: its payload dies with the run. A step that declares `output_topic_key` AND is named in a later step's `context_inputs:` does both — it persists, and the orchestrator retains the same payload and injects it from context, never re-fetching what this run just wrote.

**One record, one injection.** A later run that resumes from such a persisted record declares its key as an ordinary optional input. When a step can receive the same record both ways, the orchestrator injects exactly ONE: the in-context payload when this run produced it, the persisted record otherwise — never both, because a stale record from an earlier run would contradict the fresh one.

**Load at start.** ONE `mem_search("{project}/{change}")` to list what exists, then `mem_get_observation(id)` for the `*/handoff` records this specialist declares it consumes — nothing else. A run that resumes a staged record searches `{project}` instead, because the request may not derive the slug the record was saved under; the specialist's SKILL.md says what it looks for. Once a record is chosen, its slug is `{change}`, and the ordinary `mem_search("{project}/{change}")` then loads the rest of the run's declared inputs. Budget: 2–3 MCP calls per run; a resume's `{project}` search and its reads of at most the top 5 candidate records do not count against it.

**Organizational memory.** When a run closes, and only if it made a decision that is not obvious from the code, append ONE line to topic_key `{project}/journal`:

```
{role}@{change}: {what} — {why}
```

One line, one append, no envelope, no second save. Nothing else goes to the journal.

**Terrain or history.** Will this still be true in three months, whatever the current change? That is TERRAIN, and it belongs in the `human_nuance:` list of `.asdt/knowledge/knowledge.yaml`. Is it something decided or found while working on this change or study? That is HISTORY, and it belongs in the journal line above. Convenience never decides this — the question does.

**Consent.** Nothing enters `human_nuance` without the user's instruction or confirmation. An explicit instruction — "remember that…", "save this", "for the future:" — IS the consent. A durable fact you noticed on your own is PROPOSED in one line and written only on their yes.

**Bounded write.** Once consent exists, the ORCHESTRATOR edits the `human_nuance:` list and nothing else in that file — never a sub-agent, one entry per fact, its `note` in plain language — entry shape, `type` values, and the `# ASDT:NUANCE:BEGIN/END` markers per `asdt-init/steps/write.md`. The rest of `knowledge.yaml` belongs to `/asdt-init`. "Forget the thing about X" removes the matching entry, on the same confirmation. Before adding, read what is already there: a note that contradicts an existing one UPDATES it instead of stacking beside it. Never a secret, a token, or a credential — if asked, decline in one line and say where that belongs instead.

**Degradation.** An expected hand-off that does not exist is recorded in `open_items` with the literal prefix `ASSUMED:` and the run proceeds. A missing input never blocks and never fails a run.

## 2. Intake contract

Declared inputs arrive ALREADY INJECTED in the sub-agent prompt as `### INPUT {topic_key}` blocks, or as `### INPUT {topic_key}: UNRESOLVED`. Every declared input either arrived as a block or it did not — there is no third state, and a sub-agent NEVER fetches its own declared inputs. That work already happened, against a store the sub-agent may not even be able to see.

**One batched clarification turn.** A run gets AT MOST ONE. If gaps are genuinely blocking — no defensible hand-off is possible without an answer — collect every such question across the whole run, ask them TOGETHER as one numbered list, and stop once. Never one round trip per question, never a second turn. If the invocation already carries the router's proposal or otherwise answers your doubts, do not re-ask what is settled. A sharpened invocation answers what the user settled at routing time, not that nothing else is left to ask — raise what the router could not see, once, batched with everything else this run needs. When in doubt between asking and assuming: assume, mark it `ASSUMED:`, keep moving.

**A consent gate is not a clarification turn.** An inline step a specialist declares to get the human's go-ahead before something irreversible — approving a plan before host files are written, running commands on their machine — asks for permission, not for missing facts, so it does not count against the one clarification turn and the clarification rule does not limit it. Each gate's contract, including what happens when no human can answer, lives in the SKILL.md section that declares it; without a yes, the gated action does not happen.

**Harden always.** Every non-blocking gap degrades into an `open_items` entry prefixed `ASSUMED:` — what was assumed, and what would confirm or refute it — and the run continues. A stalled run returns nothing; a hardened run returns a hand-off whose weak spots are named and checkable. When in doubt between asking and assuming: assume, mark it, keep moving.

## 3. Step execution rules

> Whoever executes a sub-agent step or an inline prelude — a launched sub-agent, or the orchestrator running it inline — is bound by everything in this section. Do the work of this ONE step and return. Do NOT delegate, do NOT run other steps. Do NOT fetch your inputs — they arrive injected. An `UNRESOLVED` input means record the gap and proceed, never abort — unless the step file says its absence needs no entry.

An inline gate — a SKILL.md section with no step file, such as the Developer's `approve` and `verify` — sequences exactly the steps its section names, nothing else, and the rules below bind it too.

**Write boundary.** Exactly two steps in ASDT write files, and the step's identity decides it, never the identity of whoever runs it: `developer/implement` writes host source inside the edit roots its spec declares, and `asdt-init/write` writes ASDT's own state under `.asdt/`. One explicit exception is not a write: the Developer's inline `verify` runs the commands the human approved, and only commands that write no tracked source file — build output, caches, and coverage reports are not source. Prefer a command's non-writing variant where one exists (`--ci`, `tsc --noEmit`); never a `--fix`, `--write`, or `-u` flag, a snapshot update, a rewriting formatter, or codegen. Every other step writes NOTHING, anywhere — its only output is `mem_save`. If you are running any other step and reach for Edit or Write, STOP before the write — you have left the plan — and recover inside this same step: name the step that must be delegated instead, record the blocked work in `open_items`, and finish this step normally with a hand-off. STOP scopes to the write, never to the run.

**Verifiable evidence** — exact file paths, symbol names, commands, observed values instead of "should" or "likely" — is required ONLY of steps that read the codebase. Steps that do not touch code do not carry this requirement.

## 4. Injection format (orchestrator)

**Which sub-agent.** A step's `agent:` field names the executor to launch: `agent: analyst` → the `asdt-analyst` sub-agent, `agent: builder` → `asdt-builder`. Never launch a `subagent` step on a generic agent. If the host has no `asdt-{agent}` definition, launch the closest agent available and PREPEND the full text of `asdt-core/executor-header.md` to its prompt — the guardrails must reach the executor one way or the other.

The orchestrator resolves each declared input ONCE per run and injects it. Resolved:

```
### INPUT {topic_key}
{full content}
```

Failed to resolve:

```
### INPUT {topic_key}: UNRESOLVED
(could not be fetched — degrade as your step file says, and proceed)
```

**Partial failure is not total failure**: resolved inputs are injected and used normally; only the failed one becomes an `UNRESOLVED` block. Declared reference skills are read by the orchestrator and injected as `### REFERENCE SKILL {path}` blocks in the same prompt.

**Host skills.** A step's `host_skills:` names skills the HOST assistant may have installed — never ASDT files. For each name, check your own live skill list: when that skill, or an obvious equivalent under another name (a frontend or visual-design skill for `frontend-design`), is available, read its skill file's text and inject it as a `### HOST SKILL {name}` block in the same prompt, under the declared name — never invoke or activate it in your own context, where it would take over your orchestration; when it is not available, skip it silently. **A host skill is never required: its absence changes nothing and is never an `open_items` entry.** Injecting is yours alone — a sub-agent never fetches a skill itself, because it may not hold the tool. The step file says when it applies, how far it reaches, and what outranks it; a host skill never overrides the step's contract.

## 5. Hand-off schema

Roles omit keys that do not apply. Roles never add process keys.

```yaml
payload:
  what: ""                  # one sentence
  stage: ""                 # staged roles only (§1) — a stage before the last is a plan, not a delivery
  verification: {}          # staged roles only — {ran, passed, summary}; summary: the observed result or why nothing ran, in one line — on a failure, each failing command's trimmed output
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
