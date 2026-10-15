# Scope Definition — Reference

Guidelines for defining explicit, unambiguous scope. Optional reference, loaded by any role that must establish or validate boundaries for a change — typically Architect and Developer.

Apply it when writing a summary's `Constraints`, `Decisions`, `Acceptance criteria`, and `Open items` (`asdt-core/references/summary-schema.md`).

---

## Explicit In/Out Lists

Scope creep begins with ambiguity. Every scope statement carries both sides:

- **In**: features, interactions, states, and behaviors explicitly within this change. Each item must be concrete enough that someone can decide whether a given task belongs here.
- **Out**: features adjacent to this change but explicitly excluded. Listing something as out prevents implicit assumption during development.

**Rule**: If a feature is not listed as In, it is not in scope — even if it seems obviously related. If there is ambiguity about whether a feature is included, list it as Out explicitly, or record it in `Open items`.

**Example — password reset feature**:

```markdown
## Scope
- In: requesting a password reset via a registered email address
- In: receiving a time-limited reset link by email
- In: setting a new password via the reset link
- In: invalidating the reset link after use or expiry
- In: error handling for unregistered email addresses
- Out: two-factor authentication during reset
- Out: SMS-based password reset
- Out: admin-initiated password reset on behalf of a user
- Out: password strength enforcement (tracked separately)
- Out: social login account recovery
```

A role whose summary has no `Scope` section states its boundary in `Constraints` instead (`- Out of scope: SMS-based password reset`).

---

## Dependency Identification

Before finalizing scope, scan for dependencies the feature implicitly relies on:

- **Infrastructure**: does this need an email sending service, a background job processor, a new database table?
- **Upstream features**: does this require another feature to exist first (e.g. user registration before password reset)?
- **External systems**: does this integrate with a third party (payment gateway, identity provider)?

An unresolved dependency is an `Open items` entry; a settled one joins the In or Out list.

---

## NFR Categories

Non-functional requirements are constraints on quality attributes. Check each category for applicability:

| Category | Example |
|---|---|
| **Performance** | "Password reset emails must be delivered within 60 seconds of the request." |
| **Security** | "Reset tokens must be single-use and expire after 15 minutes." |
| **Accessibility** | "The reset form must be operable with keyboard-only navigation and screen readers (WCAG 2.1 AA)." |
| **Internationalisation (i18n)** | "Error messages must be translatable via the existing i18n system." |
| **Reliability** | "The feature must degrade gracefully if the email provider is unavailable." |
| **Data retention** | "Expired reset tokens must be purged within 24 hours." |

**Rules**:
- Only include NFRs directly implied by or highly relevant to the feature being specified.
- Do not invent NFRs for features with no evident quality constraints.
- NFRs must be testable or measurable (not "it should be fast").

A settled NFR is a `Constraints` entry; one the consumer must verify is also an `Acceptance criteria` entry, in Given/When/Then.

---

## Open Questions

Open questions are unresolved decisions that materially affect the design. Make each one visible: never bury an assumption in prose.

**When something is an open question, not settled by convention**:
- The idea implies multiple valid approaches with different tradeoffs.
- A business rule is unclear (e.g. "how long should reset tokens be valid?").
- An actor's permission level is ambiguous (e.g. "can guests trigger this flow?").
- A dependency is uncertain (e.g. "which email provider will be used?").

**Format**: a step never asks (protocol §3), so each open question becomes one `Open items` line: the `ASSUMED:` prefix, the assumption it proceeds on, then what would confirm or refute it.

```markdown
## Open items
- ASSUMED: reset tokens expire in 15 minutes — confirm the intended duration (15 minutes, 1 hour, 24 hours?)
- ASSUMED: the response never reveals whether the email is registered — confirm with Product
- ASSUMED: rate-limiting on reset requests is out of scope for this iteration — confirm it is tracked elsewhere
```

**Do not** include questions already answered by the original request or by obvious convention. Do not add `Open items` entries just to appear thorough — only genuinely unresolved ones.

---

## Applying it by role

- **Architect**: the In list covers architectural boundaries, service contracts, and API surface area; the Out list excludes implementation detail owned by Developer. NFRs here are scalability, reliability, and security contracts.
- **Developer**: the In list names the behavior and the steps included in this change (the files themselves are `Edit targets`); the Out list names adjacent refactors explicitly deferred. NFRs here are performance budgets, error-handling contracts, and test coverage.
