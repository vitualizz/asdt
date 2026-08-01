# Intake Contract — Shared Skill

How every specialist run treats its inputs: what arrived, what is missing, and when a question to the human is allowed. This is reference text governing the whole run, not a step of its own.

## Declared vs. present

On step start, compare the `### INPUT {topic_key}` blocks present in the prompt against the step's declared inputs in `workflow.yaml`. Every declared input either arrived as a block or it did not — there is no third state. A declared-but-missing input is recorded in `open_items` and the step proceeds; it is NEVER self-fetched. The orchestrator injects inputs exactly once; a specialist that goes looking for its own inputs is re-doing work that already happened, against a store it may not even be able to see.

## One batched clarification turn

A specialist run gets AT MOST ONE clarification turn. If gaps are genuinely blocking — the step cannot produce a defensible artifact without an answer — collect every such question across the whole run, ask them TOGETHER as one numbered list, and stop once. Never one round trip per question, never a second turn.

This turn is FULLY SUPPRESSED when the prompt carries a `## Tailored Workflow` block: the router already ran its clarifying gates before routing, and its answers are settled. In that case there is nothing left to ask — every remaining gap degrades as described below.

## Harden always

Every non-blocking gap degrades into an `open_items` entry carrying the literal prefix `ASSUMED:` — what was assumed, and what would confirm or refute it — and the run proceeds. A specialist NEVER stalls on a missing input: a stalled run returns nothing, while a hardened run returns an artifact whose weak spots are named and checkable. When in doubt between asking and assuming, assume, mark it `ASSUMED:`, and keep moving.
