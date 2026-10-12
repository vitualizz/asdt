---
title: Development
description: How to build, test, and run ASDT locally — including the specialist registration checklist and sandbox testing flow.
order: 9
locale: en
---

# Development

Practical day-to-day workflows for working on ASDT itself: building, linting, adding skills, and exercising the installer TUI without touching your real AI assistant configuration.

## Prerequisites

- Go (version pinned in `go.mod`)
- golangci-lint **v2** — install the v2 module path specifically:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

> The `.golangci.yml` config declares `version: "2"`. Installing the v1 module path gives you a v1 binary that refuses the config.

- lefthook (optional) — `make hooks` installs it and registers the pre-push lint hook

## Day-to-day commands

```sh
make build   # go build ./...
make test    # go test ./...
make lint    # golangci-lint run ./...
make hooks   # install lefthook + register pre-push hook
```

Run the installer TUI from source without a build step:

```sh
go run ./cmd/asdt-tui
```

## Testing in a sandbox

`asdt-tui` writes directly into your AI assistant's real skills directory (`~/.claude/skills`, `~/.config/opencode/skills`). To avoid touching your real config, override `$HOME`:

```sh
mkdir -p /tmp/asdt-sandbox
HOME=/tmp/asdt-sandbox go run ./cmd/asdt-tui
```

This installs into `/tmp/asdt-sandbox/.claude/skills` and `/tmp/asdt-sandbox/.config/opencode/skills`. Your real `~/.claude` is never touched.

Inspect the result:

```sh
eza --tree /tmp/asdt-sandbox/.claude/skills
```

Re-run against the same `$HOME` to exercise the "already installed / update" detection path. Clean up when done:

```sh
rm -rf /tmp/asdt-sandbox
```

## Adding a new specialist — registration checklist

Full authoring details are in [Contributing](/asdt/docs/contributing). The thing worth knowing up front:

**The embed needs no wiring.** `skill/embedded.go` uses a glob:

```go
//go:embed SKILL.md asdt-*
var skillFS embed.FS
```

Any directory named `asdt-{name}` ships automatically on the next build — there is no list to append to and no way to forget it. What *does* need doing by hand is registration: the row in the `## Registry` table of `skill/SKILL.md`, the `## ASDT Specialists` row in `internal/installer/assets/agents-template.md`, the `registryRenderOrder` and `inlineStepsDisplayNames` mirrors in `internal/installer/registry_gen.go` (`TestRegistryDrift` checks them), and the hardcoded specialist list in `skill/embedded_test.go`. Miss those and the skill still ships, it just never gets routed.

After adding the directory:

1. Run the sandbox flow to confirm the skill appears under `/tmp/asdt-sandbox/.claude/skills/{name}/`
2. Run `go test ./skill/...` — the embedded registry test verifies the skill is present

## Adding a memory provider

The prompts never name a memory provider. They speak two verbs — **remember** `(subject, note)`, which stores one self-contained note under the project's scope and a short subject (a newer note on a subject supersedes the older), and **recall** `(query)`, which returns the notes relevant to a query, by meaning or by text — defined once in `skill/asdt-core/protocol.md` §2. Only the orchestrator reaches memory, so the binding lives in the three orchestrator-facing files that do, each carrying a `<!-- ASDT:GENERATED:memory-binding -->` region committed empty: `protocol.md` (spliced, with it, into every routed `SKILL.md`), `asdt-init/SKILL.md`, and the root `SKILL.md`. At install, the installer renders the selected provider's binding into those regions. A **No memory** install (the zero `ProviderDescriptor`) renders a body saying no provider is bound, which protocol §2 reads as memory unavailable. The `asdt-analyst` and `asdt-builder` agents get no memory tools and no binding: the orchestrator strips every binding region from any file it injects into a sub-agent prompt (protocol §5), and the generated OpenCode agents deny every known provider's memory tools, a list built from `Providers`, so a new provider's tools are denied the moment you add it.

A new provider is a Go-only change — no prompt edits:

1. Add a `ProviderDescriptor` to `Providers` in `internal/installer/providers.go`:
   - `ID`, `Name`, `Description` — the TUI's provider screen lists it from these.
   - `ConfigValue` — what `/asdt-init` writes to `memory.provider`. A specialist whose project names a different provider than the install is bound to treats memory as unavailable and keeps working from the scratchpad.
   - `Tools` — each verb's recipe: one or more `MemoryTool` calls, made in order (recall may need a search, then a full read). Each carries the bare tool name and a one-line `Usage` telling the model how to call it, spliced verbatim into the binding. Engram's, for example, has remember save under the topic key `{project}/notes/{subject}`, so a newer note replaces the older, and has recall search, then read each relevant match in full. Exact keys are an optimization a provider may use; correctness never depends on them.
   - `ClaudeToolPrefixes` — optional: the MCP namespaces Claude Code may expose the tools under. The binding names them so the orchestrator matches its own tool list; a provider without one gets no prefix sentence.
   - `Detect` — how the TUI's environment check looks for the provider: `Binary`, looked up on your `PATH`, and `InstallURL`, where a user who lacks it is pointed. Leave `Binary` empty for a provider with nothing to probe locally: its row is informational and never blocks. The check lists every provider; *Continue* needs at least one usable, and the provider screen continues only on a usable one.
2. Run `go test ./...`. `TestProviders_BindEveryVerb` fails on any gap the descriptor's `Validate` finds — an unbound verb, a tool with no name or `Usage`, an empty prefix entry; `TestInstall_MemoryBindingReachesInstalledFiles` shows what a binding renders into; `TestGenerateAgents_NoMemoryAccess` keeps memory out of the generated agents; `TestPromptsAreProviderNeutral` fails if any prompt, generated agent, or persona names a shipped provider or its tools outside a binding region — the forbidden names come from `Providers`, so yours are covered the moment you add it.

## Verifying prompt edits

Editing an existing specialist's `SKILL.md` or `steps/*.md` requires no wiring changes at all. Just:

```sh
go test ./skill/... ./internal/installer/...
```

`go:embed` re-reads files at build time, so `go test` and `go run` always reflect your latest edits — no caching to worry about.

## Project structure

```
cmd/asdt-tui/       # installer TUI entrypoint
internal/
  grader/           # grades artifact payloads against probe sets
  i18n/             # TUI string catalog (en + es)
  installer/        # skill detection, installation, update logic
  setup/            # TUI views and state machine
  tui/panels/       # pure TUI rendering primitives (badge, hero, spinner)
skill/
  SKILL.md          # root orchestrator + specialist registry
  TEMPLATE.md       # specialist authoring contract (repo-only)
  README.md         # skill layer overview
  embedded.go       # go:embed registry
  embedded_test.go  # routed-invariant + embed tests
  asdt-core/        # protocol, headers, and optional references
  asdt-architect/   # Architect specialist
  asdt-developer/   # Developer specialist
  asdt-pm/          # Product Manager specialist
  asdt-qa/          # QA specialist
  asdt-researcher/  # Researcher specialist
  asdt-security/    # Security specialist
  asdt-ux-ui/       # UX/UI specialist
  asdt-init/        # Project initialization specialist (setup-class)
site/               # This documentation site (Astro)
```
