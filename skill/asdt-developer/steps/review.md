# Review — Developer Specialist

## Purpose
This step EXAMINES the quality of code that already exists; `explore` orients you to specify
a change, and `implement` delivers one. The product is judgments with fixes, not a change.

## Inputs
No `### INPUT` line: this step runs standalone, on the area the request names, the codebase,
the platform block (the project's conventions and structure), and any `### PROJECT NOTE`
lines. Look for nothing else.

## Processing

1. **Map the area.** The files, modules, and entry points the request names, plus who calls
   into them. Anchor every observation to a real path and symbol: this step reads the
   codebase, so the evidence rule applies in full.

2. **Judge against the project's own conventions** — the platform block first, then
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
proposal in the summary; turning any of them into code is a separate, explicit request.

## Output
Return this summary (`asdt-core/protocol.md` §6):

```markdown
# Developer: {change}

{The state of this code, in one sentence.}

## Decisions
- {Judgment, strongest first — strengths among them}

## Files hint
- `{path}` — {why a reader should look there first}

## Risks
- [{high|medium|low}] {finding} (`{path}`, `{symbol}`) → {fix}

## Open items
- ASSUMED: {what you could not ground in evidence}
```

`Risks` holds one line per finding, highest severity first.
