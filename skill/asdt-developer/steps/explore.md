# Explore — Developer Specialist

## Purpose
Understand the area of the codebase that will change before a single line is specified.

## Inputs
No `### INPUT` line: this step works from the raw request, the codebase, the platform block
(stack, naming conventions, key patterns), and any `### PROJECT NOTE` lines. Look for nothing
else.

**DEGRADATION**: no platform block → take the conventions from the code around the change, and record `ASSUMED: no platform block — conventions read from the code` in `Open items`.

## Processing
1. Identify which existing files and modules this change is likely to touch.
2. Identify the naming patterns and conventions, from the platform block and the code, that
   this change has to follow.
3. List the known risks and constraints ("this touches the auth layer, which has rate
   limiting").
4. List the open questions `spec` has to answer ("does this need a migration?").

Do NOT design the solution. Do NOT write code. Only explore and understand.

## Output
Return this summary (`asdt-core/protocol.md` §6):

```markdown
# Developer: {change}

{What this change touches, in one sentence.}

## Constraints
- {Naming or structural convention to follow} (`{where it shows}`)

## Files hint
- `{path}` — {why it matters to this change}

## Risks
- {Known constraint or risk} → {how the spec should handle it}

## Open questions
- {Question spec must answer}

## Open items
- ASSUMED: {what was assumed} — {what would confirm or refute it}
```
