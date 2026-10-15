package installer

import (
	"io/fs"
	"regexp"
	"strings"
)

// htmlCommentRE matches a single HTML comment and the newlines after it,
// non-greedy so multiple or multiline comments each strip cleanly without
// leaving a run of blank lines. Used by stripHTMLComments to remove
// source-only authoring comments (e.g. the H2-ownership contract) from the
// rendered agent config before it is wrapped for write.
var htmlCommentRE = regexp.MustCompile(`<!--[\s\S]*?-->\n*`)

// stripHTMLComments removes every HTML comment from s. It runs inside
// renderAgentConfig only (pre-wrap), never on the write path, so the asdt block
// markers — themselves HTML comments injected later — stay intact.
func stripHTMLComments(s string) string { return htmlCommentRE.ReplaceAllString(s, "") }

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

// renderAgentConfig reads the agent config template and persona file from the
// package-private assetsFS (internal/installer/assets/) and substitutes every
// placeholder: {{agent_name}}, {{agent_description}}, {{persona_block}},
// {{emoji_preference}}, and {{language_directive}}. The result is the body of
// ASDT's block in each assistant's global config, loaded in every session, so
// it carries only user-level preferences, never ASDT-operational rules.
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
	out = stripHTMLComments(out)

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

// InstallAgentConfig renders the agent config for the given preset and upserts
// it as ASDT's marker block into the global config file of each selected
// assistant. modes maps each AssistantID to its AgentWriteMode; assistants
// absent from the map default to AgentModeUpdate. One result per assistant; a
// per-assistant failure does not abort the others.
func InstallAgentConfig(assistants []AssistantDescriptor, preset PersonaPreset, useEmojis bool, localeCode string, modes map[string]AgentWriteMode) []AgentConfigResult {
	rendered, err := renderAgentConfig(preset, useEmojis, localeCode)
	if err != nil {
		return allAssistantsError(assistants, err)
	}

	results := make([]AgentConfigResult, len(assistants))
	for i, a := range assistants {
		results[i] = AgentConfigResult{AssistantID: a.ID}
		adapter, found := AgentConfigAdapterFor(a.ID)
		if modes[a.ID] == AgentModeSkip || !found {
			results[i].Skipped = true
			continue
		}

		path := adapter.ConfigPath()
		if err := writeAgentBlock(path, rendered); err != nil {
			results[i].Err = err
			continue
		}
		results[i].Written = []string{path}

		// Best-effort: persist persona name and emoji preference so the
		// dashboard can display them.
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

	return results
}
