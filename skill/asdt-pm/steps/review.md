# Review — PM Specialist

## Purpose
This step compares what EXISTS against what it was supposed to do; `backlog` formalizes
requirements for a change that has not been built. The product is a gap analysis.

## Inputs
- Platform summary — the project's structure and conventions, injected inline
- The area named in the request

No declared input from another specialist: this step runs standalone, and any prior
requirements for the area were already surfaced by the `knowledge-recall` prelude.
Everything arrives ALREADY INJECTED — never self-fetch.

## Processing

1. **Reconstruct the requirements in force.** Three sources, in this order of authority: the
   hand-offs recall surfaced, the repository's own docs, and the behavior observable in the
   code. Mark every requirement you INFERRED from behavior with `ASSUMED:` — code tells you
   what the product does, never what it was meant to do, and conflating the two is how a bug
   becomes a specification.

2. **Find the gaps**, applying the criteria in `asdt-core/references/scope-definition.md`:
   - **Requirements never implemented** — stated somewhere, absent from the code
   - **Behavior with no requirement** — scope that arrived without anyone deciding it. Name
     it; whether it stays is the owner's call, not yours
   - **Acceptance criteria that cannot be tested** — ambiguous, unobservable, or so broad
     that any outcome satisfies them
   - **NFRs declared with no way to measure them** — "fast", "scalable", a number with no
     method attached

3. **Findings.** At most 7, each with a one-word severity (`high`, `medium`, `low`), the
   evidence that grounds it — a path, a doc, or the hand-off it came from — and in one line
   what would close it. Plus two or three strengths with the same standard: which
   requirements are clearly met, and where the scope held.

Do NOT write new requirements for the gaps. Naming a gap is this step's job; deciding what to
build about it is a change, and a change starts at `backlog`.

## Output
Produces: `{project}/study/{topic}/pm`, with `{topic}` derived from the request in short
kebab-case.

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`:

- `what` — how well this area matches its requirements, in one sentence
- `decisions` — the judgments, strongest first, with the strengths among them
- `constraints` — the requirements you found to be in force, inferred ones marked
- `risks` — `{risk, mitigation}` per gap, highest severity first
- `files_hint` — the code and docs a reader should open first
- `open_items` — every requirement inferred rather than found stated, `ASSUMED:` prefix
