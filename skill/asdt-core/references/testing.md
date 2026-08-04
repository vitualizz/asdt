# Testing — Reference

Criteria for judging acceptance criteria, finding edge cases, and choosing a test level.
Optional reference for QA, and for anyone writing a test plan.

## Acceptance criteria

Given (precondition) / When (trigger) / Then (observable outcome). An AC that cannot be
observed by a test is not an AC.

- **Atomic.** One behavior each. An "and" in the Then joining two different behaviors means
  two ACs — unless the outcomes are always true together ("record saved AND audit entry
  written" is one concern).
- **Measurable.** "Loads quickly" is not a criterion; "renders within 200ms at p95 under 100
  concurrent users" is. "Displays correctly on mobile" is not; "renders without horizontal
  scroll at 320px" is. Always a threshold, a count, or an observable state change.
- **Negative case per functional AC.** What happens when the precondition is not met, the
  input is invalid, or a dependency is down — each as its own Given/When/Then.
- **Non-goals stated.** What the AC set explicitly does not cover, so scope cannot creep in
  during implementation.

Gap types worth naming when an AC falls short: `untestable` (nothing can observe it),
`ambiguous` (more than one valid reading), `incomplete` (missing precondition, action, or
outcome) — these three block. `missing-negative` and `missing-nonfunctional` do not block;
they are things to add.

## Edge cases

Where the interesting failures live. Work the categories the change actually touches.

- **Boundaries.** Minimum, maximum, one below, one above, zero or empty, negative where the
  domain implies non-negative, and the field omitted entirely.
- **Null / empty / missing.** `null` vs empty string vs whitespace-only vs absent key vs
  empty collection — these are four distinct states, and code that conflates them is where
  the bug is. Say what each one should do; never assume.
- **Equivalence classes.** One representative per class of identical behavior: one per valid
  shape, one per distinct rejection reason (malformed and unsupported are different classes
  even when both return 400), plus the class the domain forgot — syntactically valid but
  semantically impossible (a future birth date, a negative refund). If two values in a class
  behave differently, the class was wrong.
- **State transitions.** Draw the state table first; the empty cells ARE the edge cases. Test
  each valid transition, each invalid one (especially backward, and skipping a step), the
  same transition applied twice, and any transition out of a terminal state.
- **Concurrency.** Double submit, two writers on one record, two creations racing for one
  unique key, a read during a partial write. These need integration or e2e — a unit test
  cannot catch them.
- **Permissions.** Minimum sufficient role, one role below, unauthenticated, own resource vs
  another user's, admin acting on behalf of someone, and a token revoked mid-session.
- **Dependency failure.** Timeout, 500, connection dropped mid-request, retry against a
  non-idempotent operation. Say which way it fails: open, closed, or retry with backoff.
- **Scale and locale.** Zero items, exactly one (singular vs plural), exactly one page, one
  past the page, no search results, thousands of them. Timezone conversion, DST transitions,
  number formats (`1,000.00` vs `1.000,00`), RTL layout.

## Test level

Most tests are unit tests. Integration covers boundaries, not what unit tests already cover.
E2E covers only critical journeys.

- **Unit** when there is conditional logic, data transformation, validation, or a domain
  rule. NOT for framework wiring, pure configuration, or a pure delegation.
- **Integration** when the code crosses a process boundary, touches a real persistence layer,
  or when mocking would hide the real contract of a dependency.
- **E2E** when the flow spans services or pages and its regression would escape the other two
  levels. Keep each under 30 seconds.

Environment per level: unit runs with no external process (a unit test that opens a socket is
misclassified); integration runs the one real dependency it exists to exercise and stubs
everything past it; e2e runs the full stack against a disposable environment seeded from
fixtures, never a shared one. At every level each test builds and tears down its own state
and never depends on execution order.

Fixtures are the minimum valid state, built by a factory, named for what they are
(`expiredToken`, `emptyCart` — never `user1`). Flakiness budget is zero for unit and e2e,
under 1% for integration; design out wall-clock time, `sleep()` as synchronization, unseeded
random data, and polling without a deadline. A flaky test is worse than no test.

Coverage is a proxy, not a goal — 80% meaningful beats 100% trivial. Rough targets: domain
logic 90%+, application services 80%+, handlers and adapters 70%+, integration judged on
contract coverage rather than percentage.
