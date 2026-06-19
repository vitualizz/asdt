package installer

import (
	"io/fs"
	"strings"
)

const agentConfigPlaceholder = "detected at project init"

// emojiPrefYes / emojiPrefNo are the rendered {{emoji_preference}} bullets.
// The template prose is always English; the injected {{language_directive}}
// sets the spoken-language variant for the assistant's replies.
const (
	emojiPrefYes = "- **Emojis**: Use emojis naturally to add warmth and expressiveness — they are part of this persona's voice."
	emojiPrefNo  = "- **Emojis**: Never use emojis — keep all output plain text."
)

// AgentConfigAdapterFor returns the AgentConfigAdapterDescriptor for id, if one exists.
func AgentConfigAdapterFor(id AssistantID) (AgentConfigAdapterDescriptor, bool) {
	for _, adapter := range AgentConfigAdapters {
		if adapter.AssistantID == id {
			return adapter, true
		}
	}
	return AgentConfigAdapterDescriptor{}, false
}

// renderAgentConfig reads the AGENTS.md template and persona file from the
// package-private assetsFS (internal/installer/assets/) and substitutes all
// placeholders with the preset values. Only {{agent_name}},
// {{agent_description}}, {{persona_block}}, and {{emoji_preference}} receive
// real values; {{stack}} and {{architectural_style}} receive the sentinel
// string.
func renderAgentConfig(preset PersonaPreset, useEmojis bool, localeCode string) (string, error) {
	tmpl, err := fs.ReadFile(assetsFS, "assets/agents-template.md")
	if err != nil {
		return "", err
	}

	persona, err := fs.ReadFile(assetsFS, preset.File)
	if err != nil {
		return "", err
	}

	out := string(tmpl)
	out = strings.ReplaceAll(out, "{{agent_name}}", preset.Name)
	out = strings.ReplaceAll(out, "{{agent_description}}", preset.Description)
	out = strings.ReplaceAll(out, "{{persona_block}}", strings.TrimSpace(string(persona)))
	emojiPref := emojiPrefNo
	if useEmojis {
		emojiPref = emojiPrefYes
	}
	out = strings.ReplaceAll(out, "{{emoji_preference}}", emojiPref)
	out = strings.ReplaceAll(out, "{{language_directive}}", LocaleByCode(localeCode).Directive)
	out = strings.ReplaceAll(out, "{{stack}}", agentConfigPlaceholder)
	out = strings.ReplaceAll(out, "{{architectural_style}}", agentConfigPlaceholder)

	return out, nil
}

// allAssistantsError builds one AgentConfigResult per assistant, each carrying
// the same render error. Used when renderAgentConfig fails before any write so
// the failure is reported uniformly across all selected assistants.
func allAssistantsError(assistants []AssistantDescriptor, err error) []AgentConfigResult {
	results := make([]AgentConfigResult, len(assistants))
	for i, a := range assistants {
		results[i] = AgentConfigResult{AssistantID: a.ID, Err: err}
	}
	return results
}

// InstallAgentConfig renders the AGENTS.md template for the given preset and
// writes it to the global config location for each selected assistant.
// modes maps each AssistantID to its AgentWriteMode; assistants absent from
// the map default to AgentModeOverwrite (clean write, no prior conflict).
// One result per assistant; per-assistant failure does not abort the others.
func InstallAgentConfig(assistants []AssistantDescriptor, preset PersonaPreset, useEmojis bool, localeCode string, modes map[string]AgentWriteMode) []AgentConfigResult {
	rendered, err := renderAgentConfig(preset, useEmojis, localeCode)
	if err != nil {
		return allAssistantsError(assistants, err)
	}

	results := make([]AgentConfigResult, len(assistants))
	for i, a := range assistants {
		mode, hasMode := modes[a.ID]
		if !hasMode {
			mode = AgentModeOverwrite // non-conflicted assistant: clean write
		}

		if mode == AgentModeSkip {
			results[i] = AgentConfigResult{AssistantID: a.ID, Skipped: true}
			continue
		}

		adapter, found := AgentConfigAdapterFor(a.ID)
		if !found {
			results[i] = AgentConfigResult{AssistantID: a.ID, Skipped: true}
			continue
		}

		r, writeErr := adapter.Write(rendered, mode)
		if writeErr != nil {
			results[i] = AgentConfigResult{AssistantID: a.ID, Err: writeErr}
			continue
		}
		results[i] = r

		// Best-effort: persist persona name and emoji preference so the
		// dashboard can display them.
		if !r.Skipped {
			if meta, merr := ReadInstallMeta(a); merr == nil {
				meta.Persona = preset.Name
				if useEmojis {
					meta.Emojis = "yes"
				} else {
					meta.Emojis = "no"
				}
				_ = WriteInstallMeta(a, meta)
			}
		}
	}

	return results
}
