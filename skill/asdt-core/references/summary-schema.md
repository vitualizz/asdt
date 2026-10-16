# Summary Schema — Reference

The shape every step returns; every sub-agent step loads it as a reference. The orchestrator reads it alongside `asdt-core/protocol.md`.

A step returns a concise Markdown summary, with no envelope and no YAML. The orchestrator writes it to the scratchpad unchanged. Use these parts in this order, omit every one that does not apply, and never add one about the process itself:

```markdown
# {Role}: {change}

{What, in one sentence.}

**Stage:** {stage}
**Verification:** ran: {true|false}, passed: {true|false}. {summary}

## Decisions
## Constraints
## Files hint
## Acceptance criteria
## Risks
## Data model
## API surface
## Open items
```

- **Decisions**: one-line imperatives, with the rejected alternative in parentheses.
- **Files hint**: code anchors, meaning where to look first.
- **Acceptance criteria**: Given/When/Then, at most 5. PM is the authority on acceptance criteria.
- **Risks**: `risk → mitigation`, one line each.
- **Data model**, **API surface**: architect and developer only, when they apply.
- **Open items**: real gaps only. Assumptions carry the `ASSUMED:` prefix.
- A role adds a section of its own only where its step file defines one (the Developer's files changed, the UX/UI designer's screens).

**Stage and Verification** are for staged roles only. `Stage` names the stage the summary holds, and a stage before the last is a plan, not a delivery (the Developer's `spec`, then `implemented`). The gate that ran the checks adds `Verification`. Its summary is the observed result, or why nothing ran, in one line. On a failure, it carries each failing command's trimmed output. Both live only in the scratchpad. Across sessions, git answers what was built, and verification is run again, never recalled.

**The change's files.** A step that reviews a built change (`Stage: implemented`, writing mode) takes its files from git, the authority on what was built, never from one summary section alone. Its git answer is every path `git status --porcelain` shows modified or untracked, plus `git diff --name-only {base}...HEAD`, where `{base}` is the default branch (`origin/HEAD`, else `main` or `master`; skip this half when `HEAD` is that branch). The change's files are the union of that answer with the paths under the Developer summary's `Files changed` and `Prior delivery`, which also say where to read first. A git path that neither list names and the change's code never reaches is unrelated work: set it aside in one `Open items` line. These read-only git commands are inspection, allowed on every step.

**Golden rule: if a part does not change what the consumer does, it does not belong in the summary.**
