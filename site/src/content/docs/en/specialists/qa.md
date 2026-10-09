---
title: QA Engineer
description: Builds the safety net before code ships — test plans, acceptance criteria validation, edge case analysis, and quality reports — the specialist to bring in when "it works on my machine" isn't good enough.
order: 23
locale: en
---

# QA Engineer (`/asdt-qa`)

> Builds the safety net before code ships — test plans, acceptance criteria validation, edge case analysis, and quality reports — the specialist to bring in when "it works on my machine" isn't good enough.

## What it does

The QA Specialist finds what the acceptance criteria missed and turns it into a test plan with a go/no-go verdict. It runs one step, `test-plan`, in this order:

1. **AC gaps** — each inherited criterion judged for atomicity, measurability, and a negative case. A criterion no test could observe is a blocking gap.
2. **Edge cases** — the real job: input, state, concurrency, and dependency-failure cases the criteria never mentioned. When UX/UI ran, every flow branch and every empty, loading, and error state it named is a candidate case.
3. **Strategy** — the unit / integration / e2e split for this change, in three lines.
4. **Test cases** — Given/When/Then; when there is a Developer hand-off, each points at the built or planned file it exercises. When Security ran, every finding gets a case proving its mitigation holds; a mitigation no test can observe becomes a check you can run instead.
5. **Verdict** — `go` or `no-go`, with two lines of why. A blocking AC gap, an uncovered critical path, or — when there is a Developer hand-off — a Security `high` whose mitigation leaves no trace in its files (the changed files once built, the planned files on a plan) is a `no-go`. On a plan nothing is built yet, so the verdict is `no-go — not built yet`, and the reasoning says whether the plan is ready to build.

QA executes nothing. It never reports a pass or fail on something that wasn't run: an NFR target becomes a command you can run to measure it.

## When to invoke it

- Code is ready for review and you need a quality gate before it ships
- Acceptance criteria exist but haven't been formally validated (atomicity, measurability, independence)
- You want systematic edge case coverage, not just happy-path tests
- You need a structured test plan that a developer can implement without guessing

## On its own

It doesn't need someone to have just finished coding. Point it at what is already there:

```
/asdt-qa "what don't our auth tests cover?"
/asdt-qa "review the cart's coverage and give me a verdict"
```

It works from whatever it finds — prior hand-offs if any, the codebase if not — and still closes with go/no-go.

## Pipeline position

Typically runs **after Developer** and is the final sign-off before code merges. It reads every hand-off that exists: `pm/handoff` (acceptance criteria and NFR targets), `developer/handoff` (the plan, or what was built and whether its checks passed), `architect/handoff` (the design and its declared risks), `ux-ui/handoff` (flow branches and states), and `security/handoff` (the mitigations to prove). Can run earlier — against the PM hand-off or a saved Developer plan — to catch AC quality issues before implementation starts. That early pass saves far more time than finding gaps after the code is written. Every input is optional: with none of them, it works from the request and the codebase.

## What it produces

`qa/handoff` — AC gaps, edge cases, strategy, test cases, the checks offered to you, and the verdict, as ONE artifact. Auditing an existing suite runs `review` instead and saves `{project}/study/{topic}/qa`.

No specialist declares it as an input: it is the sign-off record. After a `no-go`, the report closes by proposing whoever fixes what was found — usually the Developer.

## Common patterns

```
/asdt-qa Review the checkout flow for edge cases
# → Happy-path is tested but boundary conditions and error paths need coverage
```

```
/asdt-qa Validate acceptance criteria before implementation starts
# → Run QA against pm/handoff to catch AC quality issues early
```

```
/asdt-qa Build a test plan for the authentication module
# → Full test pyramid strategy for security-sensitive code
```

## Limits — what it does NOT do

- Does not write implementation code
- Does not write architecture decisions or UX specs
- Never asserts a pass or fail on something that was not run — it executes nothing
- A plan that only restates the acceptance criteria as tests has added nothing — edge cases are the deliverable
- Test cases are specifications (Given/When/Then) — not executable code
