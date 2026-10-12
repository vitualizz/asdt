---
name: asdt-init
description: "Explicit-invocation-only project setup — invoke it by name to detect the stack, settle the configuration decisions with you, and WRITE .asdt/config.yaml; never fire it from a vague 'set this project up', because it changes files on disk."
user-invocable: true
specialist-id: asdt-init
metadata:
  author: "Lee Palacios (vitualizz)"
  version: "1.0"
---

# ASDT Init Specialist

You are the orchestrator of this run. Your rules (ownership, scratchpad, memory, write boundary, injection, summary schema) are `asdt-core/protocol.md`: read it now if it is not in your context. Read `./workflow.yaml` before acting.

## Role
Detect the project's facts, settle the configuration decisions with the human, and write `.asdt/config.yaml`, the project-facts authority every specialist reads (protocol §0). What the human says about the code goes to memory, never to a file. Init plans, designs, and implements nothing. It runs once before the other specialists, and again only on a deliberate recalibration.

## Orchestration Plan

Setup-class: invoked by name only, never routed (`skill/SKILL.md`), and always the same five steps in order. Init works under `{change}` = `init`, so its scratchpad files live in `{scratchpad}/asdt/init/` (protocol §1).

| Step | Execution | Does |
|------|-----------|------|
| memory-check | inline | Checks this session for the bound memory tools; settles fresh vs recalibration |
| explore | subagent | Detects the stack, context fields, and design fingerprint; returns them with the questions to ask |
| enrichment | inline | Surfaces up to 3 central, non-obvious symbols as skippable note questions |
| clarify | inline | Asks the questions, previews the write, takes the consent; produces the answers write applies |
| write | subagent | Writes `.asdt/config.yaml`; deletes the legacy knowledge files on consent |

Step identity, model, inputs, and outputs: workflow.yaml. Detection mechanics belong to `steps/explore.md`, file mechanics to `steps/write.md`. The three inline steps have no step file: the sections below are their whole contract.

## memory-check — memory and prior state (inline)

Run first, before launching anything. Both checks read only your own context and the filesystem.

**Memory.** The binding below names the bound provider and the tools behind **remember** and **recall** (protocol §2). Look for both on your OWN tool list, under whatever prefix your host adds. A sub-agent carries a different tool list, so never delegate this check.

<!-- ASDT:GENERATED:memory-binding -->
<!-- /ASDT:GENERATED:memory-binding -->

- **Available**: both tools are on your list. Continue. Init recalls nothing: it only remembers what the human tells it.
- **Unavailable** (protocol §2): say in one line that memory isn't connected, so init still detects and writes the config, but it skips the note questions and leaves any legacy notes in their file. Once the config is written, nothing reads them (protocol §0) until a later `/asdt-init` with memory connected offers to remember them. Continue. Missing memory never stops init.

**Prior state.** The project is already initialized when `.asdt/config.yaml` or `.asdt/knowledge/knowledge.yaml` (an earlier init's store, now legacy) exists. Then ask once: recalibrate, or leave it as-is? "Leave" ends the run with nothing launched. With no human to answer, leave it: a recalibration needs a yes. A project with neither file is fresh and proceeds.

## enrichment — surface the non-obvious (inline)

Surface the few symbols a newcomer would misread: structurally central yet non-obvious from the name. Each becomes one skippable question whose answer you remember as a convention note. Skip this step entirely when memory is unavailable: an answer with nowhere to go wastes the human's time. Enrichment is positive-evidence-only, never fails, and never blocks.

Take the first rung your OWN tool list supports:

1. **codegraph** (any `mcp__codegraph__*` tool): call `codegraph_explore` with "which symbols are the most central to this codebase yet the least self-explanatory from their name?". Rank chunks by caller+callee degree, skip trivially named getters, setters, and plain DTOs, and break ties by file path, then symbol name. Keep a candidate only at architecture altitude: its non-obviousness spans 2 or more files, a module boundary, or a cross-cutting invariant. Never a single-function gotcha. Keep at most 3.
2. **No codegraph**: surface nothing. A parse tree shows syntax, not centrality.

Each kept symbol becomes the question `nuance.{symbol}`: "What makes `{symbol}` non-obvious: its role, an invariant, a gotcha a newcomer would miss?", free-form, default empty, always skippable, never blocking.

## clarify — ask, preview, consent (inline)

Clarify resolves explore's ambiguities and enrichment's questions with the human, previews the write, and produces the payload `write` applies. It runs inline because only you can pause for the human.

**Ceiling.** At most 7 questions in total, the standing design-surface question aside. When more are pending, ask in this order and stop at 7: non-skippable ambiguities, skippable `low`, skippable `medium`, then the `nuance.*` questions, ties in the order explore emitted them. A skippable question past the ceiling takes its default, with origin `default`, and the preview lists it so the human can still correct it. A non-skippable one is never cut.

1. **Ambiguities.** Ask explore's ambiguities one at a time, offering its options when it gives them, a free-form value otherwise. The offered default is explore's, except for a key whose prior source is `manual`: then it is that prior value.
2. **Notes.** Ask the `nuance.*` questions, free-form, with a skip-all shortcut for the whole block. Keep only non-empty answers.
3. **Design surface.** Ask the ONE standing question: `Which surface is this project's primary design target? (mobile / tablet / desktop / none) [default: {current value, else mobile}]`, under the field `primary_design_surface`. `none` means no visual surface at all (a CLI, a library, a backend service). It is not an ambiguity and is outside the ceiling.
4. **Recalibration review.** Resolve the prior values (`steps/write.md`, step 2) and show a delta table (key, old value, new value) of every key whose value would change, plus every key whose prior source is `manual`. Ask once: "accept all, or review key by key?" Key by key: accept, reject, or set manually. Accepting over a `manual` prior value is the answer `accepted`; rejecting is the old value as a `user` answer; setting manually is a `user` answer, unless the value equals the detected one (then `accepted`). Human answers always win: a `manual` prior value appears in the table and changes only on an explicit answer.
5. **Legacy notes.** When the legacy `knowledge.yaml` carries `human_nuance` entries and memory is available, list them and ask ONE batched question: remember all, pick which, or none.
6. **Preview and consent.** Show what will land: `.asdt/config.yaml` with its key values and deltas, and every default applied past the ceiling. When the legacy `.asdt/knowledge/` files exist, ask in the same turn whether to delete them once the config is written, defaulting to no. Offer that deletion only when it cannot lose a note the human wanted kept: memory is available, or the file holds no `human_nuance`. On a decline, nothing is written and the run ends.
7. **Remember the legacy notes** the human chose, on that yes and before launching write, one memory **remember** each: subject = the entry's `topic` in kebab-case, note per protocol §2 with kind `convention`. If any of them fails, mark the deletion declined and say so in one line.

**Origins.** Each answer carries one: `user` (the human typed or chose a value that differs from the offered default, which write records as `source: manual`), `accepted` (the human chose the detected value, by typing it, picking it, or accepting it, which write records as the detected entry unchanged, over a `manual` prior value too), or `default` (a skippable default applied without the human, or the human re-confirming a `manual` prior value, which write keeps unchanged). An answer equal to the offered default is never `user`: re-confirming what is already recorded leaves `.asdt/config.yaml` byte-identical.

**The payload** is the `asdt-init/clarify` file you write on the yes, in this shape, every section present even when empty:

```markdown
# Init answers: init

skipped: {true|false}
delete_legacy: {yes|no}
memory_provider: {the `memory.provider` value the binding above names, or none when it is empty}

## Answers
- {field}: {value} ({user|accepted|default})

## Blocking open items
- {field}: {why it could not be answered}

## Notes to remember
- {symbol}: {the human's answer, verbatim}
```

**No human to answer** (a non-interactive run): ask nothing. On a fresh project the invocation by name is the consent to write the config. Apply each skippable ambiguity's default with origin `default`, list each non-skippable one under `Blocking open items`, emit no `primary_design_surface` answer, set `skipped: true` and `delete_legacy: no`, and remember nothing.

**User abort**: write nothing, launch nothing.

## After write

Write returns its summary, which the orchestrator files as the role summary. Then:

1. **Remember** each `Notes to remember` entry, one memory **remember** each: subject = the symbol in kebab-case, kind `convention`, note per protocol §2 naming the symbol. Do this only when write reports the config written or unchanged. On a halt, remember nothing.
2. Report in prose: which files changed, every `manual` value carried forward, every default applied, write's `Open items`, and, in one line, the notes remembered. When legacy files remain beside a written or unchanged config, say in one line why, and that no specialist reads them any more.
3. Close with the next step: `/asdt-architect`, `/asdt-developer`, or `/asdt` with the user's first request.

## Final Output
The role summary is write's return: what landed in `.asdt/config.yaml`, which `manual` values were carried forward, which answers were applied, and whether the legacy files were deleted. It closes the run; no specialist consumes it.

## Invariants
- **Write scope**: `.asdt/config.yaml`, and the legacy `.asdt/knowledge/knowledge.yaml`, `.asdt/knowledge/provenance.yaml`, and the then-empty `.asdt/knowledge/` (removed, only on a yes, through the executor header's one Bash exception). Only `write` touches them (protocol §4).
- **Human notes live in memory.** Init never writes a note into any file, and never deletes a legacy note it could not remember.
- **Never silently overwrite a manual value**, and an unchanged re-run leaves `.asdt/config.yaml` byte-identical.
- **explore never guesses**: no marker matched means an empty stack, never one inferred from incidental files.
