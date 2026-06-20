---
title: Getting Started
description: Install and run your first ASDT pipeline in minutes.
order: 1
locale: en
---

# Getting Started

## Requirements

Before using ASDT, you need:

- **Claude Code** or **OpenCode** — installed and authenticated
- **A memory provider** — required for cross-session persistence (default: [Engram](https://github.com/Gentleman-Programming/engram))
- A terminal (bash or zsh)

> **Building from source?** Go 1.22+ is required. Running the one-line installer downloads a pre-built binary — no compiler needed.

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh | bash
```

Downloads the latest pre-built binary (`asdt-tui`) for your platform and installs it to `~/.local/bin/`. No Go or compiler needed.

> If `~/.local/bin` is not on your `PATH`, the installer prints the exact line to add to your shell profile.

## Install the specialists

Run the interactive installer:

```bash
asdt-tui
```

It checks your setup, lets you pick which assistant(s) to target (Claude Code, OpenCode, or both), and copies the ASDT skills in — each specialist as its own independently invocable skill.

## Initialize your project

Open your AI assistant in your project folder and run:

```
/asdt-init
```

It detects your stack, asks a couple of questions, and writes `.asdt/config.yaml` with sensible defaults.

## Your first pipeline

```
/asdt Add user authentication with email and password
```

ASDT analyzes the request and recommends a specialist sequence — for example: `/asdt-pm` → `/asdt-architect` → `/asdt-developer`. Confirm the plan, then run each command. Each specialist saves its output to the knowledge base so the next one picks up where the previous left off.

## Running individual specialists

```
/asdt-pm Add dark mode to the settings page
/asdt-architect Design the caching strategy
/asdt-developer Implement the user profile component
```
