---
title: Commands
description: Complete reference for all ASDT CLI commands and slash commands.
order: 4
locale: en
---

# Commands

## CLI

### `asdt-tui`

The only CLI tool. An interactive terminal app — run it with no arguments:

```bash
asdt-tui
```

It checks your setup, lets you pick which assistant(s) to target (Claude Code, OpenCode, or both), and installs or updates the ASDT specialist skills. Everything is driven by the on-screen menu — there are no subcommands or flags.

> Initializing a project (writing `.asdt/config.yaml`) is **not** a CLI command — it's the `/asdt-init` slash command, run inside your assistant. See below.

## Slash Commands (Claude Code / OpenCode)

### `/asdt-init`

Initialize ASDT in the current project. Detects your stack and creates `.asdt/config.yaml` and `.asdt/knowledge/platform.yaml`. Run it inside your assistant, in your project folder.

```
/asdt-init
```

### `/asdt`

Pipeline routing suggestion. Describe what you want to build — ASDT analyzes the request, recommends which specialists to involve and in what order, and waits for your confirmation before listing the commands to run.

```
/asdt Add email verification to the signup flow
```

### `/asdt-pm`

Runs the Product Manager specialist only.

```
/asdt-pm Redesign the notification settings page
```

### `/asdt-architect`

Runs the Architect specialist only.

```
/asdt-architect Design the real-time sync architecture
```

### `/asdt-developer`

Runs the Developer specialist only.

```
/asdt-developer Implement the sidebar component
```

### `/asdt-qa`

Runs the QA specialist only.

```
/asdt-qa Review the checkout flow for edge cases
```

### `/asdt-security`

Runs the Security specialist only.

```
/asdt-security Review the OAuth integration
```

### `/asdt-ux-ui`

Runs the UX/UI specialist only.

```
/asdt-ux-ui Design the onboarding flow
```
