# Explore — Developer Specialist

## Purpose
Understand the area of the codebase that will change before a single line is specified.

## Inputs
No `### INPUT` line: this step works from the raw request, the codebase, the platform block
(stack, naming conventions, key patterns), and any `### PROJECT NOTE` lines. Look for nothing
else.

**DEGRADATION**: no platform block → take the conventions from the code around the change, and record `ASSUMED: no platform block — conventions read from the code` in `Open items`.

## Processing
1. **Reach.** Every place the change must reach: the code it changes, every caller of a symbol
   it touches (on an indexed project, the code-intelligence **callers** and **impact** verbs),
   and the tests, fixtures, config, and exports that depend on it — one `Files hint` line each.
2. **Reuse.** For each thing the change needs, what already provides it — an existing helper,
   component, or pattern; a stdlib or platform feature; a dependency in the manifest — by path
   or package, one `Files hint` line each.
3. The naming patterns and conventions, from the platform block and the code, this change
   has to follow.
4. The known risks ("this touches the auth layer, which has rate limiting") — first what the
   change could break for users: data destroyed or exposed, callers that stop working.
5. The open questions `spec` has to answer ("does this need a migration?").

Do NOT design the solution: name the reuse candidates, `spec` chooses. Do NOT write code.

## Output
Return this summary (`asdt-core/protocol.md` §6):

```markdown
# Developer: {change}

{What this change touches, in one sentence.}

## Constraints
- {Naming or structural convention to follow} (`{where it shows}`)

## Files hint
- `{path}` — {change | caller | test | fixture | config | export}: {why it must change or be checked}
- `{path or package}` — reuse: {what it already does for this change}

## Risks
- {Known constraint or risk} → {how the spec should handle it}

## Open questions
- {Question spec must answer}

## Open items
- ASSUMED: {what was assumed} — {what would confirm or refute it}
```
