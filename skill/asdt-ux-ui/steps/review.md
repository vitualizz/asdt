# Review — UX/UI Specialist

## Purpose
This step AUDITS the flows that already ship; `ux-spec → ui-design` designs a change. And it
audits **the user's product** — it is not the deleted design-critique, which was the model
grading its own output and produced confident prose with no signal.

## Inputs
No `### INPUT` line: this step runs standalone, on the screen, flow, or area the request names,
the codebase, the platform block (detected design system and conventions), and any `### PROJECT
NOTE` lines. Look for nothing else.

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
   `Open items` entry, never an asserted pass.

4. **Findings.** At most 7, each with a one-word severity (`high`, `medium`, `low`), the path
   or component that grounds it, and the fix in one line. Plus two or three strengths with
   the same evidence standard — an audit that only lists defects gives the reader no idea
   what already works.

## Output
Return this summary (`asdt-core/references/summary-schema.md`):

```markdown
# UX/UI: {change}

{The state of this experience, in one sentence.}

## Decisions
- {Judgment, strongest first — strengths among them}

## Files hint
- `{view or component}` — {why to open it first}

## Risks
- [{high|medium|low}] {finding} (`{path or component}`) → {fix}

## Open items
- ASSUMED: {what you could not verify from the code and tokens}
```

`Risks` holds one line per finding, highest severity first.
