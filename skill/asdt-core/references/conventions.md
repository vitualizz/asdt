# Code & Test Conventions — Reference

Language-neutral defaults for writing code and tests. Optional reference for Developer.

**Precedence**: the project wins. Where `platform-context` or the existing code has already established a convention, follow it and ignore the default here. These rules apply only where the project is silent. The no-comments rule is the one exception — only an explicit, enforced convention (a lint or CI rule that fails the build, or a written style guide) overrides it.

## Match what is already there

- Use the identifier casing already dominant in the project. Never introduce a new casing style without a reason.
- Place new files beside their peers, following the existing layout. A feature does not belong in a generic utilities folder.
- Match the project's module, import, and export style.
- Prefer libraries already in the dependency manifest before adding one.

## Simplicity first

Favor the simplest, most direct solution to the problem in front of you.

- Simplicity trims STRUCTURE, never COVERAGE. Handle every contract input, error path, and boundary condition first, then simplify the shape.
- Prefer the direct fix over the clever one. A few explicit lines beat an abstraction that hides what is happening.
- Prefer deleting to adding. If a line, a branch, or a parameter is not earning its place, remove it.
- Apply the rule of three: do not extract shared behavior on the first or second occurrence. Premature abstraction couples code that only looked alike by coincidence.
- Invest rigor up front only where the surface is hard to reverse or externally observable — a public API shape, data schema, wire format, auth model. Leaving one of those simple is a decision that needs justifying, exactly like adding an abstraction does.

## Naming and comments

- Intent-revealing names over comments. A reader should understand what an identifier holds or does from its name alone.
- **Generated code carries no comments by default.** Write one only when the user explicitly asks or an enforced project convention requires it — never at your own discretion. If code seems to need a comment, rewrite it: better names, smaller functions, clearer structure.
- When comments are requested, write them for a teammate reading the repo cold, with no ASDT installed. Never reference the authoring process, an assistant, a turn, an ADR, or a ticket number.
- Name every meaningful constant. Unnamed numbers and inline strings obscure intent and make code unsearchable.

## Structure

- **Early return.** Validate preconditions at the top and return on failure, one invalid case per return, so the happy path stays unindented and linear.
- **Small, focused functions.** One thing at one level of abstraction. If you cannot describe it in a single "and"-free sentence, split it.
- **Composition over inheritance.** Assemble behavior from collaborating parts rather than extending a base type.
- **Explicit boundaries.** Declare each dependency's expected surface at the consumer site, not the implementor, and keep it to the operations the consumer actually calls. Every such boundary is a testing seam.
- **No global state.** Inject dependencies — stores, clients, config, clocks, loggers — through constructor parameters or arguments. Package-level mutable state hides coupling and defeats substitution.
- **Errors carry context.** Enough to locate the failure without a debugger. Add locating context as an error propagates. Never swallow one silently.

## Tests

- **One behavior per case.** Multiple assertions about the same behavior are fine; unrelated behaviors in one case are not. "returns error when token is expired" is a test; "creates token, validates it, sends email, marks it used" is a scenario wearing a test's clothes.
- **Table-driven where it applies.** A list of case records — name, input, expected outcome — iterated over, instead of a copied test body per case. Every framework offers this shape under some name. Give each case its own named subtest so failures name the case that broke.
- **Fixtures, not scattered literals.** A small factory returns a valid object; each case mutates only the one field it is about. That mutation is the case's statement of intent.
- **Mock at the boundary.** Substitute the abstraction the code depends on — never a concrete type. In a language with no interface construct, the abstraction is still the set of methods the consumer calls; keep it small and substitute something that satisfies it.
- **Test behavior, not internals.** Assert the returned value, the surfaced error, or the state a collaborator was left in — not that a private method was called or a counter incremented. If renaming a private variable breaks a test, the test is testing implementation.
- **Assertion messages locate the failure**: what was expected, what arrived, in what context. Where the framework already prints a rich diff, let it.
- **No real filesystem or network in unit tests.** Temp directories via the project's helper, in-memory filesystems where available, HTTP test doubles instead of real calls, and never the user's home directory.
- **Minimum per unit of work**: one happy path, plus one failure test per distinct error condition the implementation documents.
