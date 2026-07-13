---
title: Troubleshooting
description: 'Fix common ASDT errors: command not found, PATH issues, Engram not connected, Claude Code auth failure, specialist not appearing, wrong specialist chosen.'
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

ASDT stores specialist artifacts in Engram. If Engram isn't configured, specialists cannot save or load artifacts between steps.

**Fix:**
1. Confirm Engram is wired into your assistant's MCP (Model Context Protocol — the connector protocol that lets your assistant talk to Engram; see [Knowledge Base & Memory](/asdt/docs/memory-and-engram) for how it fits together) configuration. Claude Code and OpenCode each configure MCP servers differently, so follow the MCP setup for the assistant you use (and the [Engram setup guide](https://github.com/Gentleman-Programming/engram)). Both connect to the same Engram server; only the wiring location differs.
2. Restart your assistant after editing the MCP config — both Claude Code and OpenCode start the configured MCP servers on launch.
3. Run `/asdt-init` inside your assistant — it writes the memory provider setting to `.asdt/config.yaml`.
4. Check `.asdt/config.yaml` and confirm `memory.provider: engram` is present.

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

If ASDT reports that the configured model is unavailable, your `.asdt/config.yaml` may reference an outdated model ID.

**Fix:** Open `.asdt/config.yaml` and update the `model` field to a model your assistant can serve. The available model IDs depend on your assistant and its configured provider — for Claude Code, supported models are listed in the [Claude documentation](https://docs.anthropic.com/en/docs/about-claude/models); for OpenCode, use a model ID exposed by your configured provider. You can also pick the **Chameleon** preset during install to strip the `model:` field entirely and let each assistant use its own default.

## Specialist issues

### Wrong specialist chosen — how to recover

Stop the current run. Nothing is lost — prior artifacts remain in the knowledge base. Invoke the correct specialist directly.

**Example:** If you ran `/asdt-developer` before producing an ADR, run `/asdt-architect` to create the decision record. Then re-run `/asdt-developer` — it reads the ADR automatically.

See [Specialist Comparison](/asdt/docs/specialist-comparison) to choose the right specialist for your scenario.

### Specialist command not appearing in your assistant

If typing `/asdt-pm` (or any specialist) shows no autocomplete, the skill files are not installed.

**Fix:**
1. Run `asdt-tui` and (re)install the skills for your assistant. This installs into `~/.claude/skills` for Claude Code, or `~/.config/opencode/skills` (plus command wrappers in `~/.config/opencode/commands/`) for OpenCode.
2. Restart your assistant — both Claude Code and OpenCode reload skill (and command) definitions on startup.
3. If the problem persists, run `asdt-tui` again and reinstall the skill files.

### Artifacts not loading in the next session

Each specialist reads prior artifacts from Engram. If a new session can't find artifacts from the last one, Engram was not running during the previous session when artifacts were saved.

**Fix:** Ensure Engram (via MCP) is running before invoking any specialist. Check the MCP server status in your assistant. Artifacts are only persisted if Engram is active at the time the specialist saves them.

## Known limitations

- **Memory provider required.** ASDT requires a running Engram instance (via MCP) to persist artifacts between specialist runs. There is no fallback storage — if Engram is not connected, artifacts are not saved and the next specialist in the pipeline will not find its inputs.
- **Claude Code or OpenCode required.** ASDT specialists are slash commands invoked inside a supported assistant (Claude Code or OpenCode). They do not run in a standard chat interface or via a model API directly.
- **macOS and Linux only.** The install script targets macOS and Linux (x86_64 and arm64). Windows via WSL2 is untested and unsupported in this release.
- **One active pipeline session at a time.** Running two specialist pipelines simultaneously in the same project directory may cause artifact key collisions in Engram.
