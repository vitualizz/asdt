# Review — Architect Specialist

## Purpose
Judge the architecture that ALREADY EXISTS in an area. The product is prioritized findings,
not a design — nothing is being built here.

## Inputs
- Platform summary — the stack, conventions, and existing structure, injected inline
- The area named in the request

This step declares no upstream hand-off. It runs in studies, where there is no `{change}` to
read from; if the area has prior work in memory, the `knowledge-recall` prelude already
surfaced it. Everything arrives ALREADY INJECTED — never self-fetch.

## Processing

1. **Map what is there.** Components, boundaries, dependencies, and the contracts between
   them, for the area the request names. Anchor every one to a real path or symbol. This
   step READS the codebase, so the evidence rule applies in full: a component you cannot
   point at is not a finding, it is an impression.

2. **Judge against three axes.**
   - **Coupling and boundaries** — what is cheap to change today that will be expensive
     tomorrow? Follow who imports what, and which contracts cannot move without a migration.
   - **Scalability under the load this system actually expects**, using
     `asdt-core/references/api-design.md` as the criteria. Look for the N+1s, the missing
     indexes, the synchronous chains, and the unbounded work.
   - **One-way doors already crossed** — which decisions are now expensive to reverse, and
     which are still open? Naming an open door while it is still open is the highest-value
     thing this review produces.

3. **Prioritize.** At most 7 findings. Each carries a one-word severity (`high`, `medium`,
   `low`), the evidence that grounds it (path or symbol), and one line of mitigation. Seven
   is a ceiling, not a target — five real findings beat seven padded ones.

4. **Say what is right.** Two or three strengths, each with the same evidence standard. A
   review that lists only defects is incomplete and it miscalibrates the reader: they cannot
   tell whether the area is sound with rough edges or rotten throughout.

Do NOT design a replacement. Judging what exists and proposing its rewrite are different
jobs, and mixing them buries the judgment under the proposal.

## Output
Produces: `{project}/study/{topic}/architect`, with `{topic}` derived from the request in
short kebab-case.

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`:

- `what` — the state of this architecture in one sentence
- `decisions` — the judgments, strongest first, and the strengths among them
- `risks` — `{risk, mitigation}` per finding, highest severity first
- `constraints` — what any future change to this area has to respect
- `files_hint` — where a reader should look first
- `open_items` — anything you could not ground in evidence, `ASSUMED:` prefix
