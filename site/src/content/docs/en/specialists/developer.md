---
title: Developer
description: Turns specs and designs into working code — a plan you approve, then the smallest change that fully solves it, with tests where the logic earns them — the specialist to bring in once the shape of the solution is settled and it's time to build it.
order: 22
locale: en
---

# Developer (`/asdt-developer`)

> Turns specs and designs into working code — a plan you approve, then the smallest change that fully solves it, with tests where the logic earns them — the specialist to bring in once the shape of the solution is settled and it's time to build it.

## What it does

The Developer works spec-first. After loading the project's conventions and design surface (the inline `platform-analysis` prelude), it reads the affected code before anything else — every caller, test, and config the change must reach, and what the project already has that it can reuse — writes a spec — scope in and out, numbered acceptance criteria, the approach, and the exact files it is allowed to touch, each tied to the criteria it serves — and only then, with your approval, writes code.

**The smallest change that fully solves it.** The approach is the first option that works, in this order: not building what nobody asked for, then something the codebase already has, the standard library or platform, a dependency already installed, and only then new code — never a new dependency for what a few lines do. Simplicity trims structure, never coverage: validation at trust boundaries, error handling that prevents data loss, security, accessibility, and anything you asked for are never cut. A hard-to-reverse surface — a public API, a data schema — still gets its rigor up front.

It judges its own chain from what you ask for:

| You ask for | It runs |
|---|---|
| a plan — "how would you do this?" | `explore → spec` — the plan becomes its summary and the run stops; no code is written |
| the change itself | `explore → spec → approve → implement → verify` |
| "implement the plan we approved" | `approve → implement → verify` — the plan in the session scratchpad is shown again instead of re-derived (or rebuilt, when the scratchpad is gone); a plan already built but never verified resumes at `verify` alone |

When a request is ambiguous between a plan and a build, it produces the plan, and building it later is a resume, re-approved at the gate. A resume starts from the scratchpad, never from an earlier approval: once the change is confirmed (when several could match, one question picks), it reads its own summary. A plan goes back to the approval gate; a build never verified goes to `verify`; anything else, it says what the summary already holds and stops. When the scratchpad is gone, which is normal in a new session, it rebuilds the plan from git, the recalled notes, and your words. Since the upstream summaries are gone too, it asks which specialists shaped the plan and offers to re-run them first; it also asks you to restate the plan if needed. The gate then shows you that rebuilt plan to approve. An approval given in another session never carries over.

**The approval gate.** Before any file in your repo is written, it shows you the plan in plain prose: what is in scope and what is explicitly out, the acceptance criteria, the visual direction when there is one, the names of the screens and components UX/UI designed, the approach in one line — what it reuses, the larger option it rejected, any new dependency (with the install you must run, since it never installs anything) and any shortcut it accepts with its limit — the exact files it will create and modify, which of them are tests (or that there are none), and anything it assumed or flagged — such as a UI file in a project with no visual surface. If the change was already delivered once, it also tells you that this build lands on top of that earlier delivery, which stays part of the change. You answer approve, adjust, or stop. Adjust re-runs the spec once with your words — that is how you pick the larger option instead, or drop the tests; stop ends the run at the plan, ready to resume. When no human can answer, nothing gets written.

**Write scope.** `implement` only writes inside the files the spec declared. If a needed edit falls outside them, it stops and reports the path instead of freelancing the write. Nothing under `.asdt/` is ever a Developer edit target: a change wanted there is reported as an open item pointing you to `/asdt-init`, which owns that directory. A spec that declares no files is a plan-only run: the code comes back as snippets in its summary and nothing touches the repo.

**Tests.** There is no separate test step. The spec lists the test files among the files to touch: a small test for each piece of new logic that can break — a branch, a parser, money, security, a data write, a bug fix — after a QA no-go, each of its test cases that the project has a framework for (one it lacks is flagged, not planned), the tests you ask for, and, when `strict_tdd: true` is set in `.asdt/config.yaml`, tests for every unit the change adds or alters. Approving the plan approves them, and `implement` writes exactly those, in the same pass and under the same file scope as the code. If you drop them at the gate under strict TDD, the plan records that you waived them. During a fix round, a test this run created or modified as planned may be corrected to its criterion; any other existing test is never deleted, skipped, or loosened to turn green.

**The verify gate.** `implement` writes code and tests but never runs anything. After it, the Developer shows you the check commands it suggests — build, lint, test — and what a healthy run looks like, and asks whether to run them — when the plan added a dependency, it tells you to run its install before you say yes. A command that would write your source files — a `--fix` or `--write` flag, a snapshot update, a rewriting formatter, codegen — is never offered; build output and coverage reports don't count, and a non-writing variant (`--ci`, `tsc --noEmit`) is preferred where one exists. On your yes it runs exactly those commands and nothing else; without a yes, the outcome is recorded as not run, never as a pass. If something fails, it gets at most two fix rounds, inside the same files the spec declared, re-running the same commands; still failing after that, it stops and records what is failing. A resume that starts at verify has no spec in hand, so a failure there is recorded with no fix round. The report ends by saying what was not checked — checks you declined, criteria no file covers, paths it stopped short of — when anything was.

**Shortcuts.** When the approved approach accepts a known limit ("fine up to a few thousand rows — upgrade when it grows"), the plan states it and the gate shows it. Once built, it is remembered as a project decision when memory is connected, so a later review treats it as a choice, not a finding, until the expected load crosses it.

**UI files.** When a change writes views, components, or styles, it lays them out for the project's design surface first. When UX/UI ran, it builds each screen and component to the design in UX/UI's summary — layout, hierarchy, typography, every state and breakpoint. With a design system, it adds only the tokens UX/UI proposed by name, in the system's token file, and a new variant leaves the existing ones untouched. An existing design system — or the visual direction UX/UI proposed when there isn't one — outranks everything but the approved scope. If the host assistant has a `frontend-design` skill installed, it is used to shape the visual execution, and it yields to all of the above.

## When to invoke it

- The shape of the solution is settled (requirements, architecture, or UX are defined)
- You want a plan with file-level targets you can approve before anything is written
- You want production code written to the codebase, inside a scope you approved
- You're picking up a plan from an earlier session (you approve it again at the gate)

## On its own

Point it at code that already exists and it answers without touching it:

```
/asdt-developer "how is the login flow put together today?"
/asdt-developer "how would you migrate this to the new API?"
/asdt-developer "code review the billing module"
/asdt-developer "review what you just built"
```

A review reads the code and its callers and reports numbered findings, highest severity first — bugs, error handling, missing tests on risky logic, code that should not exist, structure — each with a concrete case that breaks it (or, beyond bugs and error handling, the evidence — the duplicate's path, the unused symbol's search) and the smallest fix. Reviewing a change opens with a verdict: ship, or which findings to fix first.

A question is answered from the code with no step; asking for a plan reaches the spec; only asking it to build writes files — and only after you approve.

## Pipeline position

Typically runs **after Architect**, and it reads every upstream summary that exists: PM's (the acceptance criteria authority), the Architect's (the design decision — not re-opened), UX/UI's (the flows, the screen and component designs, and any visual direction), Security's (the mitigations this change has to carry), and QA's (after a `no-go`, the edge cases and test cases the fix has to close). All of them are optional — with none, it explores and specs the problem itself and records what it assumed. On simple changes the Architect isn't called at all.

## What it produces

The Developer's summary (`developer`) — one summary, written in stages. `spec` writes the plan (`Stage: spec`); `implement` replaces it with what was built (`Stage: implemented`) — the files changed, or the code snippets in a plan-only run; `verify` adds whether the checks ran and passed. An earlier delivery of the same change is never dropped: each new plan and build carries its files forward under `## Prior delivery`. Reviewing code that already exists runs `review` instead and writes its summary under `study-{topic}`.

Consumed by: **QA** (tests against the plan or what was built, and never claims a pass the summary doesn't show), **Security** (the planned surface or the code that changed), and the **Developer** itself when a later run resumes or iterates on it.

## Common patterns

```
/asdt-developer "how would you add CSV export to the reporting dashboard?"
# → Plan only — the plan is the summary, nothing written
```

```
/asdt-developer "implement the plan we approved"
# → Shows the plan again (rebuilt if the session ended), asks for approval, builds it, offers the checks
```

```
/asdt-developer "add CSV export to the reporting dashboard"
# → Standalone build — explores, specs, asks you to approve, implements
```

## Limits — what it does NOT do

- Does not produce architecture decisions or ADRs
- Does not write UX specs or test plans — the tests it writes are code, declared in the plan you approve
- Never writes a file in your repo without your approval in the same run
- Never writes outside the files the approved spec declared — stops and reports instead, including during fix rounds
- Never edits anything under `.asdt/` — that is `/asdt-init`'s to write
- The `implement` step never runs build, lint, or tests; commands run only at the verify gate, only on your yes, and only the ones it showed you
