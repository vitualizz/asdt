---
title: Developer
description: Turns specs and designs into working code — implementation plans, production code, and test suites — the specialist to bring in once the shape of the solution is settled and it's time to build it.
order: 22
locale: en
---

# Developer (`/asdt-developer`)

> Turns specs and designs into working code — implementation plans, production code, and test suites — the specialist to bring in once the shape of the solution is settled and it's time to build it.

## What it does

The Developer works spec-first. After loading the project's conventions and design surface (the inline `platform-analysis` prelude), it reads the affected code before anything else, writes a spec — scope in and out, acceptance criteria, the technical approach, and the exact files it is allowed to touch — and only then, with your approval, writes code.

It judges its own chain from what you ask for:

| You ask for | It runs |
|---|---|
| a question or a sanity check | `explore` |
| a plan — "how would you do this?" | `explore → spec` — the plan is saved and the run stops; no code is written |
| the change itself | `explore → spec → approve → implement → verify` |
| "implement the plan we approved" | `approve → implement → verify` — the saved plan is loaded instead of re-derived; a plan already built but never verified resumes at `verify` alone |

When a request is ambiguous between a plan and a build, it produces the plan. Because the plan is saved, building it later is a resume, not a redo: it searches memory for open plans, names the one it found — or lists them if there are several — and asks before going on. Finding none, it says so and stops.

**The approval gate.** Before any file in your repo is written, it shows you the plan in plain prose: what is in scope and what is explicitly out, the acceptance criteria, the visual direction when there is one, the names of the screens and components UX/UI designed, the exact files it will create and modify, and anything it assumed or flagged — such as a UI file in a project with no visual surface. If the change was already delivered once, it also tells you that approving this plan replaces that delivered record. You answer approve, adjust, or stop. Adjust re-runs the spec once with your words; stop ends the run with the plan saved, ready to resume. When no human can answer, nothing gets written.

**Write scope.** `implement` only writes inside the files the spec declared. If a needed edit falls outside them, it stops and reports the path instead of freelancing the write. A spec that declares no files is a plan-only run: the code comes back as snippets in the hand-off and nothing touches the repo.

**Tests.** There is no separate test step. When `strict_tdd: true` is set in `.asdt/config.yaml`, or you ask for tests, `implement` writes them in the same pass and under the same file scope as the code.

**The verify gate.** `implement` writes code and tests but never runs anything. After it, the Developer shows you the check commands it suggests — build, lint, test — and what a healthy run looks like, and asks whether to run them. A command that would write your source files — a `--fix` or `--write` flag, a snapshot update, a rewriting formatter, codegen — is never offered; build output and coverage reports don't count, and a non-writing variant (`--ci`, `tsc --noEmit`) is preferred where one exists. On your yes it runs exactly those commands and nothing else; without a yes, the outcome is recorded as not run, never as a pass. If something fails, it gets at most two fix rounds, inside the same files the spec declared, re-running the same commands; still failing after that, it stops and records what is failing. A resume that starts at verify has no spec in hand, so a failure there is recorded with no fix round.

**UI files.** When a change writes views, components, or styles, it lays them out for the project's design surface first. When UX/UI ran, it builds each screen and component to the design in its hand-off — layout, hierarchy, typography, every state and breakpoint. With a design system, it adds only the tokens UX/UI proposed by name, in the system's token file, and a new variant leaves the existing ones untouched. An existing design system — or the visual direction UX/UI proposed when there isn't one — outranks everything but the approved scope. If the host assistant has a `frontend-design` skill installed, it is used to shape the visual execution, and it yields to all of the above.

## When to invoke it

- The shape of the solution is settled (requirements, architecture, or UX are defined)
- You want a plan with file-level targets you can approve before anything is written
- You want production code written to the codebase, inside a scope you approved
- You're picking up a plan you approved in an earlier session

## On its own

Point it at code that already exists and it answers without touching it:

```
/asdt-developer "how is the login flow put together today?"
/asdt-developer "how would you migrate this to the new API?"
/asdt-developer "code review the billing module"
```

A question stops at exploration; asking for a plan reaches the spec; only asking it to build writes files — and only after you approve.

## Pipeline position

Typically runs **after Architect**, and it reads every upstream hand-off that exists: `pm/handoff` (the acceptance criteria authority), `architect/handoff` (the design decision — not re-opened), `ux-ui/handoff` (the flows, the screen and component designs, and any visual direction), and `security/handoff` (the mitigations this change has to carry). All of them are optional — with none, it explores and specs the problem itself and records what it assumed. On simple changes the Architect isn't called at all.

## What it produces

`developer/handoff` — one key, written in stages. `spec` saves the plan (`stage: spec`); `implement` replaces it with what was built (`stage: implemented`) — the files changed, or the code snippets in a plan-only run; `verify` adds whether the checks ran and passed.

Consumed by: **QA** (tests against the plan or what was built, and never claims a pass the record doesn't show), **Security** (the planned surface or the code that changed), and the **Developer** itself when a later run resumes or iterates on it.

## Common patterns

```
/asdt-developer "how would you add CSV export to the reporting dashboard?"
# → Plan only — saved, nothing written
```

```
/asdt-developer "implement the plan we approved"
# → Loads the saved plan, asks for approval, builds it, offers the checks
```

```
/asdt-developer "add CSV export to the reporting dashboard"
# → Standalone build — explores, specs, asks you to approve, implements
```

## Limits — what it does NOT do

- Does not produce architecture decisions or ADRs
- Does not write UX specs or test plans — tests, when on, are code
- Never writes a file in your repo without your approval in the same run
- Never writes outside the files the approved spec declared — stops and reports instead, including during fix rounds
- The `implement` step never runs build, lint, or tests; commands run only at the verify gate, only on your yes, and only the ones it showed you
