# Assess — Security Specialist

## Purpose
Map what an attacker can reach, work out what they could do with it, and check that against
the OWASP Top 10. Analysis only — `harden` turns it into findings.

## Inputs
- `### INPUT developer` — OPTIONAL. Developer's summary. Read its `**Stage:**` line first.
  `spec` is a plan, nothing built: map the PLANNED surface from its `Edit targets` and
  `Decisions` (its `Approach:` line) and any `Data model`, `API surface`, and `Migration notes`, and record `ASSUMED: planned surface — nothing built yet` in `Open items`.
  `implemented`: the code that changed, the change's files (`asdt-core/references/summary-schema.md`: git,
  with `Files changed` and `Prior delivery` as hints) — or `Code snippets` when `Edit roots`
  reads `Mode: plan-only`
- `### INPUT architect` — OPTIONAL. The Architect's summary: the design and its trust boundaries

Each arrives as a path to read — never look for anything else. The platform block (stack and
conventions) is in this prompt.

**DEGRADATION**: both UNRESOLVED — common, since Security runs at ANY point → assess the repository directly and record `ASSUMED: no upstream summary — surface mapped from the codebase alone` in `Open items`. Never block.

## Processing

Work in this order. It matters: threats are only real against a surface that exists, and
OWASP categories are only worth checking where the surface actually reaches them.

1. **Map the attack surface.** From the changed code, the design, and the repository, list
   the entry points this change creates or touches: endpoints, inputs, file and network
   reads, auth boundaries, secrets, third-party calls. For each one, name what crosses it,
   in which direction, and what validation stands at the boundary. Cite real paths and
   symbols — a surface you cannot point at is not evidence, it is a guess.

2. **STRIDE over that surface.** For each entry point, ask the six questions: can someone
   claim to be who they are not (Spoofing), change what they should not (Tampering), deny
   what they did (Repudiation), read what they should not (Information disclosure), exhaust
   it (Denial of service), or do more than their role allows (Elevation of privilege)? Keep
   only the threats that this surface actually admits.

3. **Cross-check OWASP.** Take the threats that survived step 2 and check them against the
   applicable Top 10 categories in `asdt-core/references/owasp.md`. Only the applicable ones
   — walking all ten against a change that adds one pure function is theater, and it buries
   the findings that matter.

Reason over the change and inspect the repository for evidence. NEVER run scanners,
dependency audits, or any command beyond the read-only git inspection of `asdt-core/references/summary-schema.md`.

## Output
Return this summary (`asdt-core/references/summary-schema.md`):

```markdown
# Security: {change}

{What this change exposes to an attacker, in one sentence.}

## Attack surface
| Entry point | Crosses | Validation today |
|---|---|---|
| `{path or symbol}` | {what data, in which direction} | {what stands at the boundary} |

## Threats
- {S | T | R | I | D | E} at `{entry point}`: {what the attacker actually does} — OWASP: {category, or —} — evidence: `{path, symbol, or absence}`

## Open items
- ASSUMED: {what was assumed} — {what would confirm or refute it}
```
