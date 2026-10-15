# AI Software Delivery Team (ASDT)

**A full software team of AI specialists — architect, developer, QA, security, UX — for when you're building alone.**

Ask one of them to look at what you already have, or hand the whole job to the team. Either way your project remembers: every specialist builds on the work of the ones before it, and the key decisions are still there days later.

📖 **[Read the docs →](https://vitualizz.github.io/asdt)**

---

## What is ASDT?

A plain chat with an AI assistant improvises. It forgets what you decided yesterday, reinvents its process on every prompt, and leaves no trail behind.

ASDT turns that chat into a **team**. It installs a set of AI specialists into [Claude Code](https://claude.com/claude-code) or [OpenCode](https://opencode.ai), each one owning a discipline:

- **Any of them works on its own.** Point the security engineer at your payments module, ask the architect whether a structure will hold, ask QA what your tests are missing. No setup, no pipeline, no ceremony — one question, one answer.
- **Or they work as a team.** Describe what you want built and they pass the work between them: the architect decides, the developer builds on that decision, QA tests it, security reviews it. You say who runs and in what order — the team proposes, you confirm, and it pauses after each specialist so you can push back.
- **Your project remembers.** What each specialist worked out is kept, so the next one starts from it instead of from scratch. You can also just ask what the team decided last week.

> **AI-powered. Human-directed.** You lead, the team executes.

## Install

One command. No Go, no build, no config:

```sh
curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh | bash
```

It downloads the pre-built binary for your platform (Linux/macOS · x86_64/arm64) into `~/.local/bin/`.

> If `~/.local/bin` isn't on your `PATH`, the installer tells you exactly what line to add.

## Get started

**1. Install the specialists into your assistant**

```sh
asdt-tui
```

An interactive menu checks your setup, lets you pick which assistant(s) to target (Claude Code, OpenCode, or both), and copies the ASDT skills in. If you pick a persona, it goes into your global `CLAUDE.md` / `AGENTS.md` as one marked block: ASDT edits only that block, never the rest of the file, and the menu's Dashboard removes it.

**2. Initialize your project**

Open your AI assistant in your project folder and run:

```
/asdt-init
```

It detects your stack, asks a couple of questions, and writes `.asdt/config.yaml`.

**3. Ask a specialist**

The fastest way in. Point one at something you already have — no plan, no setup:

```
/asdt-security "audit the payments module"
```

It reads the code, tells you what it found, and remembers the key findings for later. Every specialist works this way: `/asdt-architect "does this structure scale?"`, `/asdt-qa "what don't our auth tests cover?"`.

**4. Hand work to the team**

When there's something to build and you're not sure who you need, ask for the whole thing:

```
/asdt "add password reset"
```

The team reads the request and proposes who should work on it and in what order — say, the architect, then the developer, then security because passwords are involved. Say yes and it runs them one at a time, pausing after each so you can debate the result — or run the commands it gives you yourself. Each one picks up what the previous left behind.

**5. Ask how it's going**

```
/asdt "what did we decide about the password hashing?"
```

That one is answered from the session's working summaries, git, and what the team remembered — no work is run.

## How it works

Three layers, and each one only does its own job:

```mermaid
flowchart TD
    req([Feature request]) --> asdt
    asdt["/asdt\nreads the request · recommends a route\nruns it on your yes, pausing after each"]
    asdt -->|you confirm| specialists
    specialists["Specialist\n/asdt-architect · /asdt-developer · …\nyour assistant orchestrates its steps"]
    specialists -->|launches, isolated| steps
    steps["Step sub-agent\ndoes one thing · never delegates"]
    steps -.->|returns one summary| specialists
    specialists -->|writes summaries · passes paths| scratch[(Session scratchpad\noutside the repo)]
    scratch -->|the next specialist reads it| specialists
    specialists <-->|remember / recall\nkey knowledge only| memory[(Memory\ndefault: Engram)]

    ask([Question about what exists]) -->|straight to one specialist| specialists
```

**The router recommends.** `/asdt` reads your request, judges how much it touches and how much risk it carries, and proposes a chain of specialists. Nothing runs without your yes; on it, `/asdt` runs the chain one specialist at a time and pauses after each, so you can debate the result before the next one starts.

**Your assistant orchestrates; each specialist finishes with exactly one summary.** The main assistant owns everything that matters: it talks to you, decides, runs the approval gates, and writes the working files. Each step returns a concise Markdown summary, which the assistant writes to the session scratchpad — outside your repository — and hands to the next specialist as a path. Whatever a step worked through on the way stays inside it.

**Memory keeps only what matters.** Across sessions, a memory provider (default: [Engram](https://github.com/Gentleman-Programming/engram)) keeps the project's key knowledge — non-obvious decisions, conclusions reached in debate, conventions, lessons, audit findings — never summaries or plans. It never authorizes anything: git says what was built, you approve what gets built, and `.asdt/config.yaml` holds the project's facts. A plan resumed another day is rebuilt from git, memory, and you, and approved again. Without memory, ASDT still works within a session.

**Steps run isolated.** A step that reads the codebase gets its own sub-agent so it never pollutes the others' context. It never delegates further, never reaches memory, and never writes the scratchpad: only the orchestrator does.

Under the hood, every specialist produces a structured summary before any code is written — what the industry calls spec-driven. You never have to think about it; it is why the next specialist can start from a decision instead of a guess.

The same skill tree installs natively into either host — one embedded source renders to both Claude Code and OpenCode.

## The specialists

Invoke any of these from inside your AI assistant:

| Command | What it does | Produces |
|---|---|---|
| `/asdt` | Route advisor — analyzes your request, recommends which specialists to run and in what order, and runs them on your yes | Routing suggestion, then the chain |
| `/asdt-init` | Initialize ASDT for the project — detects stack, writes config | `.asdt/config.yaml` |
| `/asdt-pm` | Turns loose ideas into user stories with clear scope and criteria | Backlog entry |
| `/asdt-architect` | Architecture decisions, system design, risk analysis | ADR + system design |
| `/asdt-developer` | Implementation plan with production code and tests | Step-by-step plan |
| `/asdt-qa` | Test plan and acceptance criteria | Test cases, quality report |
| `/asdt-security` | Threat model and hardening checklist | Security findings |
| `/asdt-ux-ui` | User flows, screen and component design, accessibility | UX/UI summary: flows, screens, component designs |
| `/asdt-researcher` | Discovery and feasibility before requirements exist | Discovery brief |

**Unsure which one to use? Start with `/asdt`** — it points you to the right specialists for the job.

## Requirements

- [Claude Code](https://claude.com/claude-code) or [OpenCode](https://opencode.ai), installed and authenticated
- Optional, recommended: a memory provider, so key project knowledge (decisions, conventions, lessons) carries across sessions. The installer offers [Engram](https://github.com/Gentleman-Programming/engram) or **No memory**; without a provider, ASDT still works within each session
- Optional: a code-intelligence provider, so specialists answer symbol and structure questions from an index instead of grep-and-read loops in the projects it has indexed. The installer offers [Codegraph](https://github.com/Gentleman-Programming/codegraph) or **No code intelligence**; both kinds of provider plug in through Go descriptors, with no prompt edits
- A terminal (bash or zsh)

> Want to build from source instead of the installer? You'll need Go 1.22+. The one-line installer needs nothing.

## Learn more

- **[Documentation](https://vitualizz.github.io/asdt)** — getting started, the specialist model, recipes, and more

## License

[MIT](LICENSE) — built in the open by Lee Palacios ([vitualizz](https://github.com/vitualizz)).
