# Test Strategy Guidelines

## Purpose

Test pyramid principles and when to use each test type. Applied during the `test-strategy` step of the QA workflow.

This file exists to settle three decisions and nothing else: the unit / integration / e2e
split, the fixture and test-data strategy, and the environment each level needs. Anything
that does not change one of those decisions does not belong here.

## Test Pyramid

```
        /\
       /  \
      / E2E \        ← Few; slow; expensive; high confidence on user flows
     /--------\
    /Integration\    ← Some; test service boundaries and external deps
   /--------------\
  /   Unit Tests   \ ← Many; fast; cheap; test business logic in isolation
 /------------------\
```

Rules:
- Most tests should be unit tests.
- Integration tests cover only service boundaries — not re-testing what unit tests already cover.
- E2E tests cover only critical user journeys — not every feature.

## When to Write Unit Tests

Write a unit test when:
- The code contains conditional logic (if/switch).
- The code transforms data (parsing, formatting, calculation).
- The code validates input.
- The code implements a domain rule.
- The behavior must remain stable under refactoring.

Do NOT write unit tests for:
- Framework boilerplate (router registration, dependency wiring).
- Pure configuration (no logic, just values).
- Code that only delegates to another well-tested component.

## When to Write Integration Tests

Write an integration test when:
- The code crosses a process boundary (HTTP, RPC, message queue).
- The code interacts with a real persistence layer (DB query, file system).
- The behavior depends on the combined behavior of two or more components.
- Mocking would hide the real contract of an external dependency.

Keep integration tests focused on the boundary, not on business logic already covered by unit tests.

## When to Write E2E Tests

Write an E2E test when:
- The user flow spans multiple services or pages.
- The feature is critical to revenue or user retention.
- Regression of this flow would not be caught by unit or integration tests.

Limit E2E tests to the minimum set of critical paths. Each E2E test should take less than 30 seconds to run.

## Fixture Design

- Fixtures must be the minimum valid state for the test — no extra data.
- Use builder patterns or factory functions for creating test entities.
- Never share mutable state between tests.
- Name fixtures descriptively: `validUser`, `expiredToken`, `emptyCart` — not `user1`, `data`, `obj`.

## Environment and Isolation

This is the `environment` decision for each level — state it explicitly per level:

- Unit: no external process. Fakes and stubs at every boundary; a unit test that opens a socket is misclassified.
- Integration: the real dependency it exists to exercise (a real DB, a real broker) and nothing else. Everything past that boundary stays stubbed.
- E2E: the full stack against a disposable environment, seeded from fixtures, never a shared one.

Isolation rules that hold at every level: each test builds and tears down its own state,
no test depends on execution order, and database tests run inside a transaction rolled
back afterwards or against a per-test schema.

## Flakiness Tolerance

Set the tolerance per level (unit: 0, integration: < 1%, e2e: 0) and name the sources the
strategy must design out — wall-clock time (use a fixed clock), real network in unit tests,
`sleep()` as synchronization, unseeded random data, and polling without a deadline. A flaky
test is worse than no test; budget for fixing it, not for tolerating it.

## Coverage Targets

| Layer | Target |
|-------|--------|
| Business logic (domain) | 90%+ |
| Application services | 80%+ |
| HTTP handlers / adapters | 70%+ |
| Integration tests | Focus on contract, not percentage |
| E2E | Critical paths only |

Coverage is a proxy, not a goal. 80% with meaningful tests beats 100% with trivial tests.
