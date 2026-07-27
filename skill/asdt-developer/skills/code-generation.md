# Code Generation — Developer Skill

## Purpose

Guidelines for generating idiomatic, production-quality code in the `implement` step of the Developer workflow. The `platform-context` shared skill is the first authority on conventions; these guidelines are language-neutral defaults that apply only where platform-context is silent.

---

## Precedence

When platform-context conventions or existing project conventions conflict with these defaults, the platform-context and project consistency take precedence. These guidelines apply only where no project-specific convention has been established. Platform-context injects the project's detected conventions (from `knowledge.yaml`); treat its detected/high-confidence and manual fields as authoritative. Defer every language-specific idiom to it — never apply a default from this file over a convention the project has already established. Exception: the no-comments rule below is never overridden by an inferred or informal project style — only an explicit, enforced convention (a lint or CI rule that fails the build without them, or a written repo style guide) can make comments required.

---

## Match Existing Conventions

Read the injected platform context before writing any code. Then:

- Casing: use the identifier casing already dominant in the project. Never introduce a new casing style without justification.
- File structure: place new files in directories consistent with the existing layout. A feature belongs beside its peers, not in a generic utilities folder.
- Imports and exports: match the project's established module and export style rather than imposing your own.
- Libraries: prefer libraries already declared in the project's dependency manifest before introducing new ones. Treat every manifest equally as the source of truth.

If platform context is absent, infer conventions from the visible code and note the inference in the step's rationale.

---

## Simplicity First (YAGNI)

This is the umbrella principle the rules below serve. Favour the simplest, most direct
solution that solves the problem in front of you.

- Write the simplest, most direct code for decisions that are reversible and internal, where a
  later change stays contained. Invest rigor up front only where the code is hard to reverse or
  externally observable — a public API shape, data schema, wire format, or auth model
  (illustrative, not exhaustive) — because callers depend on that surface and reworking it
  later breaks them.
- Simplicity trims STRUCTURE, never COVERAGE. Handle every contract input, error path, and
  boundary condition first, then simplify the shape of the code once they are all covered.
- Prefer the direct fix over the clever one. A few explicit lines beat an abstraction that
  hides what is happening. Reach for a pattern only when the problem genuinely exhibits its
  shape, not to pre-empt one a reversible design can adopt later at low cost.
- Prefer deleting to adding. The simplest change that satisfies the acceptance criteria is
  usually the right one; if a line, a branch, or a parameter is not earning its place, remove it.
- Every abstraction, indirection, and dependency is a cost paid in surface, maintenance, and
  the next reader's time. It must earn that cost — by resolving duplication that already exists,
  or by hardening a hard-to-reverse, externally-observable surface.
- Justify the decision to leave a sensitive or contract surface simple exactly as you justify
  an abstraction you add: neither an added indirection nor an un-hardened one-way door goes
  unexamined.

The tactical rules that follow — rule of three, dependency minimalism, small focused
functions, composition over inheritance — are specific applications of this principle. When in
doubt, choose less.

---

## Naming and Self-Documentation

Make the code reveal its own intent so it needs little explanation:

- Prefer intent-revealing names over comments. A reader should understand what an identifier holds or does from its name alone.
- Generated code carries no comments by default. Write a comment only when explicitly requested — by the user in the current request, or by an explicit, enforced project convention (a lint/CI rule or a written style guide) — never at your own discretion. If code seems to need a comment to be understood, rewrite it instead: better names, smaller functions, clearer structure.
- When comments were explicitly requested, write every comment for a teammate reading the repository cold, with no knowledge of how the code was produced and no ASDT installed. Never reference the authoring process or tooling-internal artifacts in code — no mention of an assistant, a turn, an ADR, a ticket number, or any internal concept the reader may not have access to. If a rationale depends on such context, restate it in plain terms a reader without that context can act on.
- Name all significant constants. Unnamed numbers and inline strings obscure intent exactly as opaque identifiers do, and make code fragile and unsearchable. Bind each meaningful literal to a named constant.

---

## Prefer Composition Over Inheritance

Favour small, focused units composed together over deep inheritance hierarchies. Assemble behaviour from collaborating parts rather than extending a base type. Defer the concrete composition mechanism to the idiom your language and platform-context prescribe.

---

## Early Return / Early Exit

Validate preconditions at the top of the function and return early on failure. Each early return handles one invalid case and exits, so the happy path stays unindented and linear. Avoid nesting the success path inside conditional branches.

---

## Explicit Boundaries

Declare each injected dependency's expected surface at the consumer site, not at the implementor, and keep that surface small — only the operations the consumer actually calls. A narrow, consumer-owned boundary is easier to satisfy and to substitute. Express it with whatever construct the language provides — an interface, a protocol, an abstract type, or, where the language has none, a documented set of methods the collaborator must answer to. Every such boundary is a testing seam consumed by the `test-generation` skill; keeping it minimal keeps tests focused.

---

## Dependencies

- Inject all dependencies — data stores, clients, configuration, clocks, loggers — through constructor parameters or function arguments. Enforce no global state: never reach for package-level or module-level mutable state, which hides coupling and defeats substitution.
- Practise dependency minimalism. Prefer dependencies already declared in the project manifest; every new dependency is a cost in surface, maintenance, and risk. Add one only when it earns its place.

---

## Small, Focused Functions

A function should do one thing at one level of abstraction. If you cannot describe what it does in a single "and"-free sentence, split it. Small functions stay testable and readable.

Apply the rule of three to abstraction: do not extract shared behaviour on the first or second occurrence. Tolerate the duplication until the third repetition, then extract — premature abstraction couples code that only looked alike by coincidence.

---

## Error Handling

Errors must carry enough context to locate the failure without a debugger. Choose the error mechanism your language idiomatically prescribes and apply it consistently. Never swallow an error silently or discard it without handling. Add locating context as an error propagates outward, and defer the concrete error type and wrapping syntax to platform-context.

---

## Code Snippet Format

When producing code snippets in the `implement` step's payload (`code_snippets[]`):

- Show the complete, relevant code unit (function, type, method) — not a fragment that forces the reader to guess the surrounding context.
- Include whatever module, namespace, or package declaration the language requires when introducing a new file.
- If the snippet is an excerpt from a larger file, add markers to show where it fits within the existing code.
- Set `file` to the relative path from the project root.
- Set `language` to the file's language identifier.
