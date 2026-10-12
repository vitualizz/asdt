# Knowledge Recall — Shared Skill

## Purpose

Before any step runs, recall what the project already knows, so this run builds on past decisions, conventions, and lessons instead of contradicting them. The orchestrator runs this prelude inline. Memory's rules are `asdt-core/protocol.md` §2.

## Protocol

1. **Legacy notes.** Only on a project not yet migrated (protocol §0), when `.asdt/knowledge/knowledge.yaml` exists, each of its `human_nuance` entries is a note a person wrote: take it as `{topic} ({type}): {note}`, even when memory is unavailable. A file that does not parse contributes nothing. These notes join the recalled ones from step 4 on, and a recalled note on the same subject wins.
2. **Skip the recall when memory is unavailable** (protocol §2). The Prerequisites notice already told the user, so proceed with the legacy notes alone, or none.
3. **Query once.** Run ONE memory **recall**, with 3 to 6 words drawn from the request and this role (e.g. `"password reset token expiry"`, `"payments module security audit"`). A chain run by `/asdt` may reuse the notes it already recalled for the same `{change}` instead of querying again.
4. **Keep what bears on this run**, up to 5 notes and about 300 tokens. When two notes share a subject, keep the one with the later `Source` date.
5. **Check `{change}`.** A note whose `Source` names earlier work that this request may continue feeds the confirmation in protocol §0.
6. **Carry the notes to the steps.** Each step's prompt gets the notes that bear on its area, one line each: `### PROJECT NOTE: {note}`. Relevance is a judgment: a note about CSS reaches a step that touches styles, not a backlog. A note that bears on nothing in a step is not carried.

## How a note binds

- **A default, not an instruction.** Follow it unless the current request or code gives a reason not to, and say what that reason is.
- **A conflict is surfaced.** When the run's direction contradicts a note, the step records `conflicts with recalled note: {subject}` in `Open items`, and the report names it, so the user can confirm which one holds.
- **Authority wins.** A note that git or `.asdt/config.yaml` contradicts is stale (protocol §0). Say so in one line. If the note is key, remember the correction.

## Output

This prelude produces no file. It only enriches the prompts of the steps after it.
