---
title: Troubleshooting
description: 'Fix common ASDT errors: command not found, PATH issues, Engram not connected, Claude Code auth failure, specialist not appearing, switching specialists.'
order: 11
locale: en
---

# Troubleshooting

## Installation issues

### `command not found: asdt-tui` after install

The install script places the binary at `~/.local/bin/asdt-tui`. If your shell doesn't find it, your PATH doesn't include that directory.

**Fix:**

```bash
export PATH="$HOME/.local/bin:$PATH"
```

To make this permanent, add the line above to your `~/.bashrc`, `~/.zshrc`, or `~/.profile`, then restart your terminal or run `source ~/.zshrc`.

Verify the install by launching the menu:

```bash
asdt-tui
```

### Install script fails silently

If the install script exits without printing a version, run it with debug output:

```bash
bash -x <(curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh)
```

The `-x` flag prints each command as it runs. Look for the first line that fails. Common causes: `~/.local/bin` doesn't exist (fix: `mkdir -p ~/.local/bin`), or the download was blocked by a firewall.

## Memory provider issues

### Engram is not connected

ASDT remembers the project's key knowledge (decisions, conventions, lessons) in Engram. If Engram isn't connected, specialists still run: each one says in one line that memory isn't connected, works from the session scratchpad, and recalls and remembers nothing, so nothing carries over to the next session.

If you installed with **No memory**, first re-run `asdt-tui` and pick Engram on the memory step.

**Fix:**
1. Confirm Engram is wired into your assistant's MCP (Model Context Protocol — the connector protocol that lets your assistant talk to Engram; see [Scratchpad & Memory](/asdt/docs/memory-and-engram) for how it fits together) configuration. Claude Code and OpenCode each configure MCP servers differently, so follow the MCP setup for the assistant you use (and the [Engram setup guide](https://github.com/Gentleman-Programming/engram)). Both connect to the same Engram server; only the wiring location differs.
2. Restart your assistant after editing the MCP config — both Claude Code and OpenCode start the configured MCP servers on launch.
3. Run `/asdt-init` inside your assistant — it writes the memory provider setting to `.asdt/config.yaml`.
4. Check `.asdt/config.yaml` and confirm `memory.provider: engram` is present. A project that names a different provider than your install is treated as having no memory.

### Connection refused from Engram

If Engram reports a connection error, the MCP server process is not running.

**Fix:** Restart your assistant — both Claude Code and OpenCode start the configured MCP servers on launch. If the error persists, check your assistant's logs for MCP startup errors.

## Assistant issues (Claude Code / OpenCode)

### Auth failure when running a specialist

If a specialist command returns an authentication error, your assistant's session has expired.

**Fix:** Re-authenticate through your assistant's login flow, then restart it:

- **Claude Code** — sign in again with the `claude` CLI (`claude auth login`).
- **OpenCode** — sign in again with the `opencode` CLI (`opencode auth login`).

Follow the prompts to re-authenticate, then restart the assistant.

### Model not available

If your assistant reports that a model is unavailable, the installed specialist steps or executor agents may reference a model ID it can't serve. Models are set at install time, not in `.asdt/config.yaml`.

**Fix:** Re-run `asdt-tui` and pick a different model preset, or set the model per step yourself. The available model IDs depend on your assistant and its configured provider — for Claude Code, supported models are listed in the [Claude documentation](https://docs.anthropic.com/en/docs/about-claude/models); for OpenCode, use a model ID exposed by your configured provider. You can also pick the **Chameleon** preset to strip the `model:` field entirely and let each assistant use its own default.

## Specialist issues

### Switching specialists

There's no wrong specialist to recover from — each one works independently and reads whatever summaries earlier specialists left in the session scratchpad, so switching costs nothing. If a request doesn't fit the specialist you invoked, just stop and run a different one directly; nothing is lost.

**Example:** If you ran `/asdt-developer` before there was a design, run `/asdt-architect` to make the decision. Then re-run `/asdt-developer` on the same request — it reads the Architect's summary as an input.

If you're not sure which specialist fits your request, run `/asdt` — it recommends a more suitable one when your request doesn't fit. See [Specialists](/asdt/docs/specialists) to choose the right specialist for your scenario.

### Specialist command not appearing in your assistant

If typing `/asdt-pm` (or any specialist) shows no autocomplete, the skill files are not installed.

**Fix:**
1. Run `asdt-tui` and (re)install the skills for your assistant. This installs into `~/.claude/skills` for Claude Code, or `~/.config/opencode/skills` (plus command wrappers in `~/.config/opencode/commands/`) for OpenCode.
2. Restart your assistant — both Claude Code and OpenCode reload skill (and command) definitions on startup.
3. If the problem persists, run `asdt-tui` again and reinstall the skill files.

### Earlier work missing in the next session

This is expected for working files: each specialist's summary lives in the session scratchpad, and the scratchpad belongs to the session. A new session rebuilds what it needs from git, from the notes memory recalls, and from you, and every gate (like the Developer's plan approval) asks again.

What should carry over is the key knowledge: decisions, debate conclusions, conventions, lessons, audit findings. If a new session recalls none of it, memory wasn't connected when it would have been remembered.

**Fix:** Ensure Engram (via MCP) is running before invoking a specialist; each specialist warns in one line when it isn't. Notes are only remembered while memory is connected. For anything lost, tell the specialist again, or say "remember that…" to keep it.

## Known limitations

- **Memory is what crosses sessions.** Without a connected memory provider, ASDT works within a session, but no knowledge carries over to the next one: you restate the context a recall would have supplied.
- **Claude Code or OpenCode required.** ASDT specialists are slash commands invoked inside a supported assistant (Claude Code or OpenCode). They do not run in a standard chat interface or via a model API directly.
- **macOS and Linux only.** The install script targets macOS and Linux (x86_64 and arm64). Windows via WSL2 is untested and unsupported in this release.
- **Sessions don't share working files.** Each session has its own scratchpad, so a second session working on the same change doesn't see the first one's summaries: it rebuilds what it needs from git, memory, and you.
