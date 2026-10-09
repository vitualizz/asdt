---
title: Security
description: Hunts for the gaps an attacker would find first — threat models, OWASP reviews, hardening checklists — the specialist to bring in whenever auth, data handling, or external integrations are on the table, at any point in the pipeline.
order: 24
locale: en
---

# Security (`/asdt-security`)

> Hunts for the gaps an attacker would find first — threat models, OWASP reviews, hardening checklists — the specialist to bring in whenever auth, data handling, or external integrations are on the table, at any point in the pipeline.

## What it does

The Security Specialist performs threat modeling and security analysis using STRIDE and the OWASP Top 10. It maps the attack surface, identifies threats systematically, and produces a prioritized hardening checklist where every finding has a concrete, actionable mitigation — not "monitor it" or "add logging."

The critical invariant: **Security has no required predecessor.** It can run at any stage — on a fresh project with no prior artifacts, mid-development, or after launch. If upstream artifacts exist (architecture decisions, implementation), it reads them. If they don't, it works from the platform context and request alone, noting gaps in `open_items` and proceeding.

Depth is gated by the risk surface, not complexity — the specialist judges it from what the change touches: authentication, secrets, data handling, external integrations. This is the only specialist where the question isn't "how complex is the feature?" but "how large is the attack surface?" The risk surface sets how deep the pass goes; both steps, `assess` and `harden`, run every time.

## When to invoke it

- Authentication, session management, or authorization is involved
- The feature handles or stores personally identifiable information
- External integrations, webhooks, or user-controlled URLs are present
- New API endpoints are being exposed publicly
- Any time before shipping to production when security hasn't been reviewed

## On its own

The specialist most often used alone. It requires no change in progress:

```
/asdt-security "audit the payments module"
/asdt-security "review how we handle sessions"
/asdt-security "what does our webhooks endpoint expose?"
```

It maps the surface, judges it, and leaves you prioritized findings with a concrete mitigation each.

## Pipeline position

**No required predecessor** — invoke at any point. For maximum impact, run it after the Architect (it reads `architect/handoff` for the API surface and trust boundaries, and `developer/handoff` for the code that changed, when they exist). Given a Developer plan that isn't built yet, it maps the planned surface from the files the plan declares and records that nothing is built yet. For a quick threat model early in design, run it before architecture is finalized to surface design-level risks before they're baked in.

Its hand-off is read by the Architect, the Developer, and QA.

## What it produces

`security/handoff` — findings and hardening checklist as sections of ONE artifact (`{project}/study/{topic}/security` when the run audits what already exists):

- **Findings** — highest severity first, each with a one-word severity (`high`, `medium`, `low`), what an attacker gains, the evidence that grounds it, and the concrete mitigation that closes it
- **Hardening checklist** — ordered items, each verifiable as done, plus where the mitigations land

Consumed by: **Architect** (a finding that reshapes a boundary is a design constraint), **Developer** (each mitigation that lands in the change's files becomes a constraint on the spec), **QA** (every finding gets a test case that proves its mitigation holds; a `high` whose mitigation leaves no trace in the Developer's changed or planned files is a `no-go`).

## Common patterns

```
/asdt-security Audit the OAuth integration
# → External auth flow with token handling — high risk surface
```

```
/asdt-security Threat model the new payment webhook handler
# → User-controlled input hitting financial logic
```

```
/asdt-security Quick security pass before the v2 launch
# → No prior artifacts needed — runs from platform context alone
```

## Limits — what it does NOT do

- Does not write implementation code
- Does not produce architecture decisions or UX specs
- Does not produce test plans (though its findings inform what QA should cover)
- Every finding must have a concrete mitigation — "add monitoring" is not a mitigation
- Severity is one word — `high`, `medium`, or `low`. No CVSS, no numeric scores
- Never runs scanners, dependency audits, or any other command — it reasons over the change and reads the repository
