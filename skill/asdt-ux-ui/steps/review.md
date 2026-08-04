# Review — UX/UI Specialist

## Purpose
This step AUDITS the flows that already ship; `ux-spec` specifies flows for a change. And it
audits **the user's product** — it is not the deleted design-critique, which was the model
grading its own output and produced confident prose with no signal.

## Inputs
- Platform summary — the project's detected design system and conventions, injected inline
- The screen, flow, or area named in the request

No declared input from another specialist: this step runs standalone, and anything prior
about the area was already surfaced by the `knowledge-recall` prelude. Everything arrives
ALREADY INJECTED — never self-fetch.

## Processing

1. **Walk the flows that exist.** Read the views, components, and routes for the named area,
   and reconstruct what the user actually goes through step by step. Anchor to real paths,
   component names, and route definitions: this step reads the codebase, so the evidence rule
   applies in full.

2. **Judge the experience.**
   - **Friction** — steps that ask for something the product already knows, decisions the
     user has no basis to make yet, confirmations on reversible actions
   - **Missing states** — empty, loading, and error. A flow that only renders the happy path
     is not finished, and the absence is usually invisible until it ships
   - **Inconsistency with the detected design system** — a component restyled locally, a
     spacing or type choice that fights the tokens the project already has

3. **Audit accessibility** against `asdt-core/references/accessibility.md`: focus behavior
   and visible indicators, keyboard reachability and no traps, labelling for every control,
   and contrast pairs. A requirement you cannot verify from the code and tokens is an
   `open_items` entry, never an asserted pass.

4. **Findings.** At most 7, each with a one-word severity (`high`, `medium`, `low`), the path
   or component that grounds it, and the fix in one line. Plus two or three strengths with
   the same evidence standard — an audit that only lists defects gives the reader no idea
   what already works.

## Output
Produces: `{project}/study/{topic}/ux-ui`, with `{topic}` derived from the request in short
kebab-case.

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`:

- `what` — the state of this experience in one sentence
- `decisions` — the judgments, strongest first, with the strengths among them
- `risks` — `{risk, mitigation}` per finding, highest severity first
- `files_hint` — the views and components a reader should open first
- `open_items` — anything you could not verify from the code and tokens, `ASSUMED:` prefix
