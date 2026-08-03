# Review — Developer Specialist

## Purpose
This step EXAMINES the quality of code that already exists; `explore` orients you to specify
a change, and `implement` delivers one. The product is judgments with fixes, not a change.

## Inputs
- Platform summary — the project's conventions and structure, injected inline
- The area named in the request

No declared input from another specialist: this step runs standalone, and anything prior
about the area was already surfaced by the `knowledge-recall` prelude. Everything arrives
ALREADY INJECTED — never self-fetch.

## Processing

1. **Map the area.** The files, modules, and entry points the request names, plus who calls
   into them. Anchor every observation to a real path and symbol: this step reads the
   codebase, so the evidence rule applies in full.

2. **Judge against the project's own conventions** — the platform summary first, then
   `asdt-core/references/conventions.md` where the project is silent. Look for:
   - **Complexity hotspots** — the function or module nobody wants to touch, and why
   - **Duplication that has earned extraction** — the third occurrence, not the second
   - **Error handling** — swallowed errors, lost context, failures that cannot be located
   - **Convention violations** — casing, layout, or import style that fights the codebase
   - **Units that outgrew their purpose** — the helper that became a service, the module
     that accumulated three unrelated responsibilities

   Judge against what this project does, not against a style you prefer. A convention the
   codebase follows consistently is the convention, even when you would have chosen another.

3. **Findings.** At most 7, each with a one-word severity (`high`, `medium`, `low`), the
   path and symbol that grounds it, and the fix in one line. Seven is a ceiling, not a
   target.

4. **Say what is right.** Two or three strengths with the same evidence standard. A review
   that lists only defects miscalibrates the reader: they cannot tell whether the area is
   solid with rough edges or genuinely in trouble.

**This step writes NOT ONE BYTE.** It is `agent: analyst`, not builder. Every fix is a
proposal in the hand-off; turning any of them into code is a separate, explicit request.

## Output
Produces: `{project}/study/{topic}/developer`, with `{topic}` derived from the request in
short kebab-case.

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`:

- `what` — the state of this code in one sentence
- `decisions` — the judgments, strongest first, with the strengths among them
- `risks` — `{risk, mitigation}` per finding, highest severity first
- `files_hint` — where a reader should look first
- `open_items` — anything you could not ground in evidence, `ASSUMED:` prefix
