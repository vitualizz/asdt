# Test Generation — Developer Skill

## Purpose

This skill provides guidelines for generating high-quality, maintainable test cases. Apply these guidelines in the `test` step of the Developer workflow, whether the step emits test snippets in its payload or writes real test files.

The goal is tests that document intent, catch regressions, and are easy to extend — not tests that merely achieve coverage numbers.

These guidelines are language-neutral. The `platform-context` shared skill is the first authority on the project's test framework, file naming, assertion style, and mocking idiom; defer every concrete syntax choice to it and apply the principles below within whatever idiom the project already uses.

---

## One Test Case Per Behavior

Each test case must verify exactly one observable behavior. Multiple assertions testing the same behavior in one case are fine; testing unrelated behaviors in one case is not.

- **Good**: "returns error when token is expired" verifies one failure mode.
- **Avoid**: "creates token, validates it, sends email, and marks it as used" — this is a scenario test disguised as a unit test.

---

## Table-Driven Tests

Group related cases into a single data-driven table — a list of case records that the test iterates over — rather than copying a whole test body per case. This reduces boilerplate, makes the intent of each case explicit, and makes it trivial to add a new edge case. Every mainstream test framework offers this shape under some name (parameterized tests, shared examples, `each`-style loops, or a plain loop over a list of cases).

Each case record carries: a descriptive name, the input, and the expected outcome.

```
cases = [
  { name: "valid token passes",          token: validToken(),                   expectError: false },
  { name: "expired token returns error", token: validToken(expiresAt: past),    expectError: true, message: "token expired" },
  { name: "missing id returns error",    token: validToken(id: none),           expectError: true, message: "token ID is required" },
]

for case in cases:
    run subtest named case.name:
        result = validateResetToken(case.token)
        if case.expectError: assert result failed and its message contains case.message
        else:                assert result succeeded
```

Give each case its own named subtest so a failure report names the exact case that broke.

---

## Test Fixtures, Not Hardcoded Values

Extract repeated setup into a small factory that returns a valid object, then have each case mutate only the one field it cares about. Hardcoded values scattered across cases make refactoring painful, and a case that mutates one field states its intent plainly.

```
validToken(overrides) -> ResetToken with sane defaults, overridden per case

token = validToken()
token.expiresAt = past   # this case is about expiry, and nothing else
```

Keep fixtures close to the tests that use them and register any needed cleanup through the framework's teardown hook.

---

## Mock at the Boundary

Never mock a concrete type. Substitute the abstraction the code under test depends on — the interface, protocol, port, or injected collaborator named in its constructor or arguments. Every injected dependency from the `code-generation` skill is a seam for testing.

```
# The code under test depends on an abstraction with a minimal surface
TokenStore: save(token) -> result, findById(id) -> token | notFound

# The double implements that same surface, with per-test behavior injected
fakeStore = TokenStore double where:
    save     -> records the call and returns success
    findById -> returns the token this case set up
```

In a language without an explicit interface construct, the abstraction is still the set of methods the consumer calls — keep that set small and substitute an object that satisfies it. Prefer a hand-written double or the project's established mocking library over inventing a new one.

---

## Test Behavior, Not Implementation

Tests should assert the OUTCOME observable from the public interface, not the internal mechanism:

- **Good**: assert the value the function returned, the error it surfaced, or the state a collaborator was left in.
- **Avoid**: assert that a specific private method was called, or that an internal counter was incremented.

If a test breaks when you rename a private variable without changing observable behaviour, the test is testing implementation.

---

## Meaningful Assertion Messages

Every failed assertion must produce a message that identifies what was expected, what was received, and in what context.

```
# Good
assert got.status == "active", "status: got {got.status}, want \"active\""

# Avoid
assert got.status == "active", "wrong status"
```

When the framework already prints a rich diff of expected versus actual, let it — add a message only where the diff alone would not tell the reader which case or which field failed.

---

## No Real Filesystem or Network in Unit Tests

- Use the project's temp-directory helper for any test that writes files, and let the framework clean it up.
- For filesystem-backed logic, prefer the project's in-memory filesystem abstraction, or a temp directory scoped to the test, over touching real project paths.
- Never make real HTTP calls in unit tests. Use the project's HTTP test double — a local stub server or a substituted client abstraction.
- Never read from the user's home directory or other system paths in tests — they produce non-reproducible results in CI.

---

## Coverage Requirement per Task

For every implementation task, generate at minimum:

1. One happy-path test (valid input → expected output).
2. One failure test per distinct error condition the implementation documents.

When the task produces a structured output payload, also assert that its required fields are populated and that `open_items` is present even when empty.

---

## Test Snippet Format

When producing test snippets in the `test` step's payload:

- Set `file` to the test file path the project's convention implies for the code under test.
- Include the complete test unit, not a fragment.
- Use the table-driven shape for any task with more than one case.
- Reference the implementation's function and type names directly so the relationship is traceable.
