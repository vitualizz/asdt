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

Any directory named `asdt-{name}` ships automatically on the next build — there is no list to append to and no way to forget it. What *does* need doing by hand is registration: the §5 and §9.2 rows in `skill/SKILL.md`, the ASDT Specialists row in `internal/installer/assets/agents-template.md`, and the hardcoded specialist list in `skill/embedded_test.go`. Miss those and the skill still ships, it just never gets routed.

After adding the directory:

1. Run the sandbox flow to confirm the skill appears under `/tmp/asdt-sandbox/.claude/skills/{name}/`
2. Run `go test ./skill/...` — the embedded registry test verifies the skill is present

## Verifying prompt edits

Editing an existing specialist's `SKILL.md` or `steps/*.md` requires no wiring changes at all. Just:

```sh
go test ./skill/...
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
  asdt-shared/      # shared skill fragments
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
