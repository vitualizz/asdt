// Package installer also embeds the installer-only data assets (persona
// markdown + the AGENTS.md template) that renderAgentConfig consumes to
// generate each assistant's agent config. These are NOT conversational
// skill content and never ship into .claude/skills/; they live under
// internal/installer/assets/ separate from the skill/ tree.
package installer

import "embed"

//go:embed assets/personas/*.md assets/agents-template.md
var assetsFS embed.FS
