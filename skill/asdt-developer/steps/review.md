# Review — Developer Specialist

## Purpose
EXAMINE the quality of code that already exists — an area, or a built change. The product is
judgments with fixes, never a change.

## Inputs
No `### INPUT` line: this step runs standalone, on the area the request names, the codebase,
the platform block (the project's conventions and structure), and any `### PROJECT NOTE`
lines. Look for nothing else.

## Processing

1. **Map the area.** The files, modules, and entry points the request names, plus who calls
   into them. A built change, branch, or uncommitted diff → its files from git
   (`asdt-core/protocol.md` §6) plus every caller of a symbol they change, judged against the
   expected load the repo shows (README, deploy config), recorded as `ASSUMED: expected load —
   {what}`. "What you built" with HEAD on the default branch and a clean tree → the last
   commit (`HEAD~1..HEAD`), recorded as `ASSUMED: reviewed HEAD~1..HEAD`. Anchor every
   observation to a real path and symbol.

2. **Judge against the project's own conventions** — the platform block first, then
   `asdt-core/references/conventions.md` where the project is silent — in this order:
   - **Bug** — wrong result, crash, missed edge case (empty, zero, last item, rounding, time
     zone), a caller the change broke, a fix in one caller while the shared code stays broken
   - **Error handling** — swallowed errors, lost context, writes that can lose data
   - **Missing test** — logic that earns a test (`conventions.md`, Tests) and has none
   - **Lean** — what a higher rung of the `conventions.md` ladder makes unnecessary: dead code,
     unused options, a one-implementation abstraction that is not a testing seam, a helper or
     dependency duplicating one the project has (name its path), duplication on its third
     occurrence
   - **Structure** — the function nobody wants to touch, the unit with unrelated jobs (split
     by job, never by line count), casing, layout, or imports that fight the codebase

   Judge against what this project does, not a style you prefer. Security depth is Security's;
   test plans are QA's. A `shortcut:` a `### PROJECT NOTE` records is a decision, not a
   finding, unless the expected load crosses its limit.

3. **Confirm before reporting.** A Bug or Error-handling finding needs a concrete case — this
   input or situation → this wrong result; any other needs its evidence — the grep of the unused
   symbol across the whole tree, the duplicate's path, the untested branch. Neither, no finding.
   Re-read the lines: the caller exists, the value can really be empty. Fix = the smallest that
   works, a deletion first. No taste, no "consider".

4. **Findings.** At most 7 — a ceiling, not a target — numbered in one sequence, highest
   first: `high` must fix (bug, data loss, breaks at the expected load), `medium` should fix
   (risky logic untested, real slowness, duplication, code that should not exist), `low` nice
   to have. On a change's review, the one-sentence line under the heading is `Verdict: ship.` or
   `Verdict: fix {n, m} first.`; when the ceiling cut findings, that line (or, on an area review,
   the state sentence) ends with `{k} lower findings not listed.`

5. **Say what is right.** Two or three strengths, with the same evidence, so the reader can
   tell an area solid with rough edges from one in trouble.

**This step writes NOT ONE BYTE.** Every fix is a
proposal in the summary; turning any of them into code is a separate, explicit request.

## Output
Return this summary (`asdt-core/protocol.md` §6):

```markdown
# Developer: {change}

{The state of this code, in one sentence.}

## Decisions
- Keep: {strength} (`{path}`, `{symbol}`) — {why it must survive the fixes}

## Files hint
- `{path}` — {why a reader should look there first}

## Risks
- {n}. [{high|medium|low}] {finding} — {case | evidence}: {input or situation → wrong result | what proves it} (`{path}`, `{symbol}`) → {fix}

## Open items
- ASSUMED: {what you could not ground in evidence}
```
