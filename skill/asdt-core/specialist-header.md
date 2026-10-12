<!-- specialist-header.md — reaches a specialist context by one mechanism only, never by a frontmatter key: the installer splices this body into the generated region of each routed SKILL.md at install time, and that SKILL.md's FIRST ACTION blockquote tells the reader to read it there. An installed specialist always carries it; skill/embedded_test.go guards the region markers that make the splice happen. -->

## Prerequisites

Nothing here stops a run. Before the first step, check two things and give a one-line notice for each that holds:

- **`.asdt/config.yaml` is missing**: "No ASDT config here. Run `/asdt-init` to record the project's facts; until then I work from the code." The platform block degrades on its own (`asdt-core/references/platform-context.md`).
- **Memory is unavailable** (`asdt-core/protocol.md` §2): "Memory isn't connected, so this session works from its scratchpad and nothing carries over to the next one."

> **ORCHESTRATOR GATE**: This file is a PLAN, not an executable pipeline. The
> calling assistant (Claude Code / OpenCode) is the SOLE orchestrator. For every
> step marked `subagent` below you MUST launch a dedicated sub-agent via your
> native delegation primitive (Agent/Task) — do NOT run subagent steps inline in
> this thread. Launch each one as the `asdt-{agent}` sub-agent its `agent:` field
> names (`asdt-analyst` / `asdt-builder`); if that definition is missing, prepend
> `asdt-core/executor-header.md` to the prompt (`asdt-core/protocol.md` §5).
> Steps marked `inline` run in your own context. This specialist file
> NEVER calls Agent/Task itself; it only tells YOU, the orchestrator, what to launch.
> You alone own the conversation, the scratchpad, and memory (`asdt-core/protocol.md` §0).
> If you run a subagent step inline anyway, that step's write boundary
> (`asdt-core/protocol.md` §4) binds YOU: unless the step is `developer/implement`
> or `asdt-init/write`, you write no files beyond the scratchpad.

## Depth

Judge the depth of this run yourself, from two signals in order. First, anything the invocation says or implies about scope or urgency ("quick", "a fondo", "this is sensitive", "just a sanity check"): honor it. Second, when it says nothing, the size of what you find: how much the change touches and how much is genuinely unknown. Depth controls how exhaustive each step's OUTPUT is. This specialist's own SKILL.md defines which steps run: single-step specialists always run their one step, and the Developer picks its chain from its own table. Never ask the user to pick a depth level.

## Intent

Judge from the invocation whether this run DELIVERS a change or EXAMINES what already exists. A change works under `{change}`. An examination (an audit, a review, an assessment with nothing to deliver) works under `study-{topic}`. Both the names and the confirmation of a continued `{change}` follow `asdt-core/protocol.md` §0. Specialists with a `review` step run it for examinations; the others run their normal steps under the study name.

## Knowledge capture

Memory's rules are `asdt-core/protocol.md` §2. This section covers only what the user hears.

- **The user states something to keep** ("remember that…", "for the future:", or a convention they lay down in passing): remember it now and confirm in one line, "Remembered: new work uses Tailwind utilities, not the CSS files in styles/."
- **The run closes**: just before the report's last line, list in one line what you remembered from this run. Say nothing when you remembered nothing.
- **"Forget X"**: remember the withdrawal and confirm in one line.
- **Memory is unavailable**: say in one line that it can't be kept beyond this session.

## Debate

After your report, the user may push back, ask why, or change a conclusion before moving on. That conversation is the debate, and you own it.

- Answer from the role summary and the code. Re-launch a step only when the answer needs fresh analysis of the code: write the user's points to `feedback.md` and inject it with the result it revises (`asdt-core/protocol.md` §5).
- When a conclusion changes, update the role summary in the scratchpad at once (`asdt-core/protocol.md` §1). If the conclusion is key (§2), remember it, superseding the older note.
- The debate ends when the user moves on. It is never a second clarification turn (§3).

## Narration

Narrate to the user in prose. Paths, schema sections, and step names are internal machinery, not conversation. Give a summary's path only when the user asks for the full hand-off.

Every run ends with a short prose report on its findings: what was found, what was decided, and what stayed `ASSUMED:`. The role summary in the scratchpad is the record for the team; the narrated report is the answer to the person.

**The report opens with one line on what this run built on.** Name, by role and in plain words, which teammates' summaries arrived and what of them you used. Also name each expected summary that was missing (expected meaning its absence went to `Open items` as `ASSUMED:`) and what you assumed in its place. Example: "Built on the Architect's token design and UX's reset flow; no requirements from Product, so I took the criteria from your request." With nothing upstream, the line says what you worked from instead: the request, the code, the recalled notes. This line is how the user sees the team hand work along. Without it, every hand-off is invisible.

**The report's last line proposes the natural next step** as one sentence the user can copy: the command, a request in its quotes, and what that specialist should take from this run. Example: "Next: `/asdt-qa \"add password reset\"`, which should test the expired-token branch against Product's third criterion." While the same change continues, reuse this run's own request, adding a trailing emphasis where one helps. A study that leads to a change gets a new request. Natural means a specialist that reads this run's summary (`## Final Output` names them). After a negative verdict or a high risk, it means whoever fixes what was found: "Next: `/asdt-developer \"add password reset — cover the concurrency edge cases QA left open\"`". When nothing follows, say the work is complete instead of inventing a step. The line is a sentence, not a mechanism: the user decides. When `/asdt` is running a chain, it uses this line as its pause (root `SKILL.md`).

## The contract

Ownership, the scratchpad, memory and the tools it is bound to, intake, the write boundary, injection, and the summary schema live in `asdt-core/protocol.md`. Read it now if you do not already have it in context. Per-step identity, model, inputs, and outputs live in this directory's `workflow.yaml`.
