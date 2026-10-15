# Review — PM Specialist

## Purpose
This step compares what EXISTS against what it was supposed to do; `backlog` formalizes
requirements for a change that has not been built. The product is a gap analysis.

## Inputs
No `### INPUT` line: this step runs standalone, on the area the request names, the codebase,
the platform block, and any `### PROJECT NOTE` lines (prior requirements for the area). Look
for nothing else.

## Processing

1. **Reconstruct the requirements in force.** Three sources, in this order of authority: the
   project notes recall surfaced, the repository's own docs, and the behavior observable in the
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
   evidence that grounds it — a path, a doc, or the project note it came from — and in one line
   what would close it. Plus two or three strengths with the same standard: which
   requirements are clearly met, and where the scope held.

Do NOT write new requirements for the gaps. Naming a gap is this step's job; deciding what to
build about it is a change, and a change starts at `backlog`.

## Output
Return this summary (`asdt-core/references/summary-schema.md`):

```markdown
# PM: {change}

{How well this area matches its requirements, in one sentence.}

## Decisions
- {Judgment, strongest first — strengths among them}

## Constraints
- {Requirement in force} ({stated in `{doc}` | ASSUMED: inferred from `{path}`})

## Files hint
- `{path or doc}` — {why to open it first}

## Risks
- [{high|medium|low}] {gap} (`{evidence}`) → {what would close it}

## Open items
- ASSUMED: {requirement inferred rather than found stated}
```

`Risks` holds one line per gap, highest severity first.
