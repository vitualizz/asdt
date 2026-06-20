---
title: 'Tutorial: Your First Pipeline'
description: 'Step-by-step ASDT tutorial: install, init, and run a complete PM → Architect → Developer pipeline for "Add a contact form with validation".'
order: 3
locale: en
---

# Tutorial: Your First Pipeline

In this tutorial you will install ASDT and run a complete PM → Architect → Developer pipeline for a concrete feature: **adding a contact form with validation**.

Estimated time: 15 minutes.

## Prerequisites

Before starting, make sure you have:

- **Claude Code** (or OpenCode) installed and authenticated. Run `claude auth login` if you haven't already.
- **Engram** running as an MCP server in Claude Code. Check your MCP config in Claude Code settings.
- A **bash or zsh** terminal.

## Step 1 — Install ASDT

Run the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh | bash
```

The binary is placed at `~/.local/bin/asdt-tui`. If your shell doesn't find it, add this to your shell profile:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then install the specialist skills into your assistant:

```bash
asdt-tui
```

This interactive menu lets you pick your assistant(s) and copies the skill files in. Then **restart Claude Code** — it loads the skill definitions on startup.

If `asdt-tui` shows `command not found`, see [Troubleshooting](/asdt/docs/troubleshooting).

## Step 2 — Initialize ASDT in your project

Open Claude Code in your project directory and run:

```text
/asdt-init
```

This creates `.asdt/config.yaml` (memory provider settings) and `.asdt/knowledge/platform.yaml` (detected stack), so every specialist tailors its output to your project.

## Step 3 — Ask ASDT for a pipeline recommendation

Open Claude Code in your project and run:

```text
/asdt "Add a contact form with validation"
```

ASDT analyzes your request and recommends which specialists to involve and in what order.

```text
# Example output — your results will vary based on your project

Analyzing: "Add a contact form with validation"

Recommended pipeline:
  1. /asdt-pm        → Define scope and acceptance criteria
  2. /asdt-architect → Design validation approach and API contract
  3. /asdt-developer → Implement the form handler and validation logic

Run each command in order. Each specialist reads the previous one's output automatically.
```

## Step 4 — Run the Product Manager

```text
/asdt-pm "Add a contact form with validation"
```

PM locks scope, writes user stories, and saves a `pm/backlog-entry` artifact to the knowledge base.

```text
# Example output — your results will vary based on your project

Running PM specialist...

Produced: pm/backlog-entry

  Feature: Contact form with validation
  
  User stories:
    US-01: As a visitor, I want to submit a contact form so that I can reach the team.
      AC1: Form fields: name (required), email (required, valid format), message (required, max 500 chars).
      AC2: Submission shows a success message and sends an email notification.
      AC3: Invalid fields show inline error messages without losing existing input.
    
    US-02: As a visitor, I want the form to prevent submission if my email is invalid
      so that I don't have to wait for a bounce to know something went wrong.
      AC1: Email field validates format on blur and on submit.
      AC2: Submit button is disabled while validation errors exist.

  In scope: form UI, client-side validation, form submission handler, email notification.
  Out of scope: CAPTCHA, rate limiting, analytics tracking.

Artifact saved to knowledge base.
```

## Step 5 — Run the Architect

```text
/asdt-architect
```

Architect reads the PM backlog entry automatically, evaluates technical approaches, and produces an Architecture Decision Record (ADR) and system design.

```text
# Example output — your results will vary based on your project

Running Architect specialist...
Loading: pm/backlog-entry ✓

Produced: architectural-decision + system-design

  ADR: Use native HTML5 Constraint Validation API with progressive enhancement
  
  Decision: Use the browser's native Constraint Validation API for field validation,
  augmented with a thin JS layer for custom error messages and blur-event handling.
  
  Alternatives considered:
    - Zod schema validation (rejected: over-engineered for a contact form)
    - React Hook Form (rejected: introduces unnecessary dependency)
  
  System design:
    - POST /api/contact — accepts { name, email, message }, returns 200 or 422
    - Email notification via nodemailer (existing dep)
    - No database write — form data forwarded directly to email

Artifacts saved to knowledge base.
```

## Step 6 — Run the Developer

```text
/asdt-developer
```

Developer reads the ADR and system design, breaks work into ordered tasks, and produces an implementation plan (or writes code in writing mode).

```text
# Example output — your results will vary based on your project

Running Developer specialist...
Loading: architectural-decision ✓
Loading: system-design ✓

Produced: developer/dev-implementation

  Implementation plan:

  T-01  Create src/components/ContactForm.astro
        — Form markup with native constraint validation attributes
        
  T-02  Create src/pages/api/contact.ts
        — POST handler: validate, send email via nodemailer, return 200/422
        
  T-03  Create src/styles/contact-form.css
        — Error state styles for invalid fields
        
  T-04  Update src/pages/contact.astro
        — Import and render ContactForm component

Artifact saved to knowledge base.
```

## What's next

Your pipeline is complete. From here you can:

- Run **[QA](/asdt/docs/specialists/qa)** (`/asdt-qa`) to validate acceptance criteria coverage and generate a test plan.
- Run **[Security](/asdt/docs/specialists/security)** (`/asdt-security`) to review the form handler for OWASP vulnerabilities.
- Browse **[Recipes](/asdt/docs/recipes)** for command sequences for other common scenarios.
- See **[Troubleshooting](/asdt/docs/troubleshooting)** if anything went wrong.
