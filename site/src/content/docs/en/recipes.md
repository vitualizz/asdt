---
title: Recipes
description: 'Copy-paste ASDT command recipes for common scenarios: add a form, build an API, run a security review, start mid-pipeline, check test coverage, design a new screen.'
order: 12
locale: en
---

# Recipes

Copy-paste command sequences for the most common ASDT scenarios. Each recipe shows the exact commands to run in Claude Code or OpenCode.

## Full pipeline recipes

### 1. Ship a new user-facing feature

Start here when a feature request is in vague language and needs the full PM → Architect → Developer sequence.

```text
/asdt "Add a contact form with email notification"
# Review the suggested specialist sequence, then run each in order:
/asdt-pm "Add a contact form with email notification"
/asdt-architect
/asdt-developer
```

### 2. Build a new REST API endpoint

For backend-only changes where the API contract is the primary design decision.

```text
/asdt "Add a POST /api/v1/subscriptions endpoint"
/asdt-pm "Add a POST /api/v1/subscriptions endpoint"
/asdt-architect
/asdt-developer
```

### 3. New screen with UX design first

For user-visible features where UI design should precede implementation.

```text
/asdt "Add an onboarding wizard for new users"
/asdt-pm "Add an onboarding wizard for new users"
/asdt-ux-ui
/asdt-architect
/asdt-developer
/asdt-qa
```

### 4. Feature with security review before shipping

For anything touching auth, payments, PII, or external integrations.

```text
/asdt "Add OAuth login with GitHub"
/asdt-pm "Add OAuth login with GitHub"
/asdt-architect
/asdt-security
/asdt-developer
/asdt-qa
```

### 5. Explore before planning (fuzzy problem)

When the problem is unclear and you need discovery before requirements.

```text
/asdt-researcher "We're losing users at the signup step — what could we do?"
# Researcher produces a discovery brief with a recommended direction.
# Then hand the brief to PM:
/asdt-pm "Based on the discovery brief: add progressive disclosure to the signup flow"
```

## Single-specialist recipes

### 6. Lock scope and write user stories

When you have a clear feature idea and just need structured requirements.

```text
/asdt-pm "Add dark mode toggle to user settings"
```

### 7. Document an architecture decision

When the technical approach needs a formal ADR and system design.

```text
/asdt-architect "Document the decision to use PostgreSQL row-level security for multi-tenancy"
```

### 8. Write production code from a settled spec

When scope and architecture are already locked in the knowledge base.

```text
/asdt-developer "Implement the multi-tenancy RLS policy from the Architect's ADR"
```

### 9. Security audit on an existing feature

Run at any point — no prior specialist run required.

```text
/asdt-security "Review the session management in the auth module"
```

### 10. Design a new UI component

When you need a component spec before the developer starts coding.

```text
/asdt-ux-ui "Design a data table component with sorting, filtering, and pagination"
```

### 11. Validate test coverage before shipping

Run after Developer — QA reads the implementation artifact automatically.

```text
/asdt-qa
```

## Mid-pipeline recipes

### 12. Pick up at Developer after an existing ADR

When PM and Architect ran in a previous session. Developer loads prior artifacts automatically.

```text
/asdt-developer
# No need to re-run PM or Architect — artifacts are in the knowledge base.
```

### 13. Add a security review to an in-flight pipeline

Run Security at any point without restarting the pipeline.

```text
/asdt-security
# Security reads system-design and dev-implementation from the knowledge base.
```

### 14. QA a completed feature without a full pipeline

When a feature was built without ASDT — run QA against the existing code.

```text
/asdt-pm "Add acceptance criteria for the existing checkout flow"
/asdt-qa
```
