# Review — QA Specialist

## Purpose
This step EXAMINES the test suite that already exists for an area; `test-plan` plans tests
for a change. The product is a verdict on the suite's health plus the gaps behind it.

## Inputs
- Platform summary — the project's test framework and conventions, injected inline
- The area named in the request

No declared input from another specialist: this step runs standalone, and anything prior
about the area was already surfaced by the `knowledge-recall` prelude. Everything arrives
ALREADY INJECTED — never self-fetch.

## Processing

1. **Map what is tested.** Find the tests covering the area and read what each one actually
   asserts — not what its name claims. Anchor to real paths and test names: this step reads
   the codebase, so the evidence rule applies in full. A behavior with a test whose
   assertions never reach it is UNTESTED, and saying so is the point of this step.

2. **Find the gaps**, using the catalogue in `asdt-core/references/testing.md`:
   - **Behaviors with no test at all** — start from what the code does, not from what the
     test directory suggests
   - **Edge cases uncovered** — work the categories the area actually touches: boundaries,
     null versus empty versus absent, invalid state transitions, concurrency, dependency
     failure
   - **Weak or tautological assertions** — the test that asserts a mock was called, the one
     that re-computes the expected value with the same code under test
   - **Brittle tests coupled to internals** — the ones that break on a rename and pass
     through a real regression
   - **Misclassified levels** — the "unit" test that opens a socket, the e2e that only
     needed a function call

3. **Verdict on the suite, two lines.** Does this area's suite catch what it is meant to
   catch? Answer plainly, and say what would change the answer.

4. **Findings.** At most 7, each with a one-word severity (`high`, `medium`, `low`),
   evidence, and the fix in one line. Plus two or three strengths with the same evidence
   standard — a coverage audit that only lists holes tells the reader nothing about what
   they can already trust.

**Never a pass or fail on something you did not run.** This step reads tests; it does not
execute them. If measuring would settle a question, name the command the USER can run.

## Output
Produces: `{project}/study/{topic}/qa`, with `{topic}` derived from the request in short
kebab-case.

Persist via `mem_save` under this step's `output_topic_key`, using the canonical hand-off
schema from `asdt-core/protocol.md`:

- `what` — the health of this area's suite in one sentence
- `decisions` — the verdict and the judgments behind it, strengths among them
- `risks` — `{risk, mitigation}` per gap, highest severity first
- `files_hint` — the test files and the code they should be covering
- `open_items` — anything you could not ground in evidence, `ASSUMED:` prefix
