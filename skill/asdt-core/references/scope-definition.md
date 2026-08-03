# Scope Definition — Reference

Guidelines for defining explicit, unambiguous scope. Optional reference, loaded by any role that must establish or validate boundaries for a change — typically Architect and Developer.

Apply it when producing the `constraints`, `decisions`, and `open_items` of a hand-off.

---

## Explicit In/Out Lists

Scope creep begins with ambiguity. Every scope statement carries both sides:

- **in**: features, interactions, states, and behaviors explicitly within this change. Each item must be concrete enough that someone can decide whether a given task belongs here.
- **out**: features adjacent to this change but explicitly excluded. Listing something as "out" prevents implicit assumption during development.

**Rule**: If a feature is not listed as in, it is not in scope — even if it seems obviously related. If there is ambiguity about whether a feature is included, list it as out explicitly, or move it to an open item.

**Example — password reset feature**:

```yaml
scope:
  in:
    - Requesting a password reset via registered email address
    - Receiving a time-limited reset link by email
    - Setting a new password via the reset link
    - Invalidating the reset link after use or expiry
    - Error handling for unregistered email addresses
  out:
    - Two-factor authentication during reset
    - SMS-based password reset
    - Admin-initiated password reset on behalf of a user
    - Password strength enforcement (tracked separately)
    - Social login account recovery
```

---

## Dependency Identification

Before finalizing scope, scan for dependencies the feature implicitly relies on:

- **Infrastructure**: does this need an email sending service, a background job processor, a new database table?
- **Upstream features**: does this require another feature to exist first (e.g. user registration before password reset)?
- **External systems**: does this integrate with a third party (payment gateway, identity provider)?

Unresolved dependencies go to `open_items`; settled ones go to the in/out lists.

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

They belong in `constraints` when settled, in `acceptance_criteria` when the consumer must verify them.

---

## Open Questions

Open questions are unresolved decisions that materially affect the design. Surface them rather than assume them.

**When something is an open question, not an assumption**:
- The idea implies multiple valid approaches with different tradeoffs.
- A business rule is unclear (e.g. "how long should reset tokens be valid?").
- An actor's permission level is ambiguous (e.g. "can guests trigger this flow?").
- A dependency is uncertain (e.g. "which email provider will be used?").

**Format**: write each as a full interrogative sentence someone could answer directly.

```yaml
open_items:
  - "ASSUMED: reset tokens expire in 15 min — confirm the intended duration (15 min, 1 hour, 24 hours?)"
  - "ASSUMED: the response never reveals whether the email is registered — confirm with product"
  - "ASSUMED: rate-limiting on reset requests is out of scope for this iteration"
```

**Do not** include questions already answered by the original request or by obvious convention. Do not add open items just to appear thorough — only genuinely unresolved ones.

---

## Applying it by role

- **Architect**: the in list covers architectural boundaries, service contracts, and API surface area; the out list excludes implementation detail owned by Developer. NFRs here are scalability, reliability, and security contracts.
- **Developer**: the in list names files to create or modify and the steps included in this change; the out list names adjacent refactors explicitly deferred. NFRs here are performance budgets, error-handling contracts, and test coverage.
