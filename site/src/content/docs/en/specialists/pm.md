---
title: Product Manager
description: Turns a raw request into requirements the rest of the team can build from without guessing — stories in delivery order, explicit scope, and acceptance criteria.
order: 20
locale: en
---

# Product Manager (`/asdt-pm`)

> Turns a raw request into requirements the rest of the team can build from without guessing — stories in delivery order, explicit scope, and acceptance criteria.

## What it does

PM runs one step and hands back one artifact. It takes what you asked for and turns it into 3 to 6 user stories, each carrying one or two acceptance criteria in plain language.

**The order of the stories is the priority.** The first one ships first. No MoSCoW ratings, no separate priority field, no matrix to keep in sync: to change the priority, change the order.

**Out-of-scope is mandatory.** A hand-off that does not say explicitly what is excluded is incomplete — not "there was nothing to exclude". Scope ambiguity is where most scope creep comes from, and naming the adjacent thing is what stops it.

NFRs only make the cut when the feature genuinely implies them and they are measurable. "p95 under 200ms, measured with k6" is an NFR; "it should be fast" is nothing.

## When to bring it in

- The request is in user language and needs pinning down before design or code
- Scope has to be agreed before work starts, so it cannot grow mid-flight
- Several needs have to be reconciled and someone has to decide the order

Don't call it for a refactor, a cosmetic change, or a request that already arrives technically scoped. Go straight to Developer for those.

## How to invoke it

You talk to it. There are no flags: if you want more or less depth, say so inside the request and the specialist adjusts.

```
/asdt-pm "add email and password authentication"
```

```
/asdt-pm "redesign notifications — several teams are involved, take your time"
```

```
/asdt-pm "CSV export on the reports panel, keep it tight"
```

## What it produces

A single hand-off at `{project}/{change}/pm/handoff`:

| Field | What it carries |
|---|---|
| `what` | The change in one sentence |
| `decisions` | The stories, in delivery order |
| `constraints` | Scope in, scope out, and the measurable NFRs |
| `acceptance_criteria` | Given/When/Then, max 5 |
| `risks` | `{risk, mitigation}`, one line each |
| `open_items` | Real gaps, prefixed `ASSUMED:` where PM answered its own question |

**PM is the authority on acceptance criteria.** Everyone downstream refines them — Developer takes them to implementation granularity, QA hunts for what they missed — but nobody rewrites them from scratch.

Consumed by **Architect**, **Developer**, **QA**, and **UX/UI**. All of them read it as an optional input: if PM never ran, each works from the raw request and records that.

## What it consumes

The raw request, the project's detected conventions, and — when Researcher ran first — its `researcher/handoff`: the recommended direction becomes the starting point, and the rejected directions seed out-of-scope.

With none of that available, it works from the request alone and says so.

## Boundaries

- Doesn't write architecture decisions or ADRs
- Doesn't write code or technical designs
- Doesn't write UX flows or component specs
- Never produces a hand-off without explicit out-of-scope items
