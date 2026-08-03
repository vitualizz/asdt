---
title: Researcher
description: Takes a fuzzy problem and returns ONE recommended direction with feasibility behind it — the specialist to bring in before requirements exist, when you don't yet know what to build.
order: 26
locale: en
---

# Researcher (`/asdt-researcher`)

> Takes a fuzzy problem and returns ONE recommended direction with feasibility behind it — the specialist to bring in before requirements exist, when you don't yet know what to build.

## What it does

It diverges before PM converges, in a single step:

1. **Frames the problem.** A fuzzy request usually hides two or three different problems. It names which one it is exploring and which ones it is deliberately setting aside.
2. **Diverges.** Generates 3 to 5 genuinely different directions. Different means they fail for different reasons — three variations on one idea is one direction, not three.
3. **Judges feasibility.** Green, yellow, or red per direction, each with **one line of evidence**: a file that already does something similar, a missing dependency, a constraint that rules it out. A verdict with nothing behind it isn't a verdict — it gets marked `ASSUMED:` and stays visible.
4. **Converges.** Recommends exactly one direction and says why it beat the others. Every discarded one carries its reason: a rejected direction with no reason written down is the one that comes back next quarter.

## When to bring it in

- You don't know what to build yet, only what hurts
- Several paths look plausible and you want the tradeoffs before committing
- Someone already proposed a solution and you want to know whether it's the only one

## How to invoke it

Plain language, no flags. If you want it to go deeper or lighter, say so in the request.

```
/asdt-researcher "users drop out of onboarding and we don't know where"
```

```
/asdt-researcher "we want semantic search — explore thoroughly before we commit"
```

## What it produces

A single hand-off at `{project}/{change}/researcher/handoff`:

| Field | What it carries |
|---|---|
| `what` | The recommended direction in one sentence |
| `decisions` | The recommendation first, then every discarded direction as `rejected: {direction} — {why}` |
| `constraints` | What any implementation of that direction has to respect |
| `files_hint` | The code it actually read while judging feasibility |
| `risks` | `{risk, mitigation}` for the recommended direction |
| `open_items` | Verdicts it couldn't ground in evidence, prefixed `ASSUMED:` |

`decisions` is where the exploration survives: not just what it picked, but what it looked at and set down.

Consumed by **PM** as an optional input: the recommended direction becomes its starting point, and the rejected ones seed its out-of-scope list.

## Where it sits

The only **pre-requirements** specialist. It runs before PM and never replaces it: it recommends a direction, PM decides what gets built. It also works standalone, when structured exploration is all you want.

## Boundaries

- Doesn't write requirements, architecture, code, or tests
- Never writes to the filesystem
- Doesn't rank while diverging — ranking early throws away the option you hadn't finished thinking about
