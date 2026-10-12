package installer

import (
	"bytes"
	"fmt"
	"strings"
)

// ProviderID identifies a known memory provider.
type ProviderID = string

// ProviderEngram identifies the Engram memory provider.
const ProviderEngram ProviderID = "engram"

// MemoryVerb names one operation of the provider-neutral memory interface the
// orchestrator-facing prompts speak (skill/asdt-core/protocol.md §2). Prompts
// never name a provider's tools; the installer binds each verb to the selected
// provider's tools at install time.
type MemoryVerb string

// The memory interface. remember stores one concise, self-contained note with
// its project scope and a short subject (a newer note on a subject supersedes
// the older); recall retrieves the notes relevant to a query, by meaning or by
// text. Correctness never depends on exact-key retrieval.
const (
	MemoryRemember MemoryVerb = "remember"
	MemoryRecall   MemoryVerb = "recall"
)

// memoryVerbs is the canonical verb order: the order the binding lists them.
var memoryVerbs = []MemoryVerb{MemoryRemember, MemoryRecall}

// MemoryTool is one provider tool call in a verb's recipe.
type MemoryTool struct {
	Name  string // bare tool name, without any host prefix (e.g. "mem_save")
	Usage string // how to call it for this verb, in prompt prose; spliced verbatim
}

// ProviderDescriptor describes a known memory provider: everything the
// installer needs to bind the prompts' memory interface to it. Adding a
// provider is a new entry in Providers — no prompt edits.
type ProviderDescriptor struct {
	ID          ProviderID
	Name        string
	Description string
	// ConfigValue is what /asdt-init writes to `memory.provider` in
	// .asdt/config.yaml. protocol.md §2 ("Memory unavailable") treats memory as
	// unavailable when that config value names a different provider.
	ConfigValue string
	// Tools binds every verb in memoryVerbs to its recipe: one or more tool
	// calls, made in order (recall may need a search, then a full read).
	Tools map[MemoryVerb][]MemoryTool
	// ClaudeToolPrefixes are the MCP namespaces Claude Code may expose the
	// tools under (plugin install and plain MCP install). The binding names
	// them so the orchestrator matches its own tool list; optional — a
	// provider without one gets no prefix sentence.
	ClaudeToolPrefixes []string
	// Detect is how the setup TUI's environment check looks for the provider.
	Detect ProviderDetect
}

// ProviderDetect is how the setup TUI's environment check looks for a provider
// before installing.
type ProviderDetect struct {
	// Binary is the executable looked up on PATH. Empty means the provider has
	// nothing to probe locally: its row is informational and never blocks.
	Binary string
	// InstallURL is where a user whose Binary is missing is pointed. Required
	// when Binary is set.
	InstallURL string
}

// Providers lists all known memory providers.
var Providers = []ProviderDescriptor{
	{
		ID:          ProviderEngram,
		Name:        "Engram",
		Description: "Persistent cross-session memory via the Engram MCP server.",
		ConfigValue: "engram",
		Tools: map[MemoryVerb][]MemoryTool{
			MemoryRemember: {{
				Name:  "mem_save",
				Usage: "pass `project: {project}`, `topic_key: {project}/notes/{subject}` (saving to an existing `topic_key` replaces that note, which is how a newer note supersedes the older), `title: {subject}`, `type` (`decision` for a decision, `pattern` for a convention, `learning` for a lesson, `discovery` for a finding), and the note as `content`",
			}},
			MemoryRecall: {
				{
					Name:  "mem_search",
					Usage: "pass the query as `query` and `project: {project}`; each match is a truncated preview with its `id` — a preview is never the note",
				},
				{
					Name:  "mem_get_observation",
					Usage: "pass a match's `id` as `id` to read the full note, for each match whose preview bears on the query",
				},
			},
		},
		ClaudeToolPrefixes: []string{"mcp__plugin_engram_engram__", "mcp__engram__"},
		Detect: ProviderDetect{
			Binary:     "engram",
			InstallURL: "https://github.com/Gentleman-Programming/engram",
		},
	},
}

// Memory-binding generated-region markers. They bound the provider binding in
// every orchestrator-facing skill file that reaches memory
// (asdt-core/protocol.md, and through it every routed SKILL.md;
// asdt-init/SKILL.md; the root SKILL.md). Sub-agents never reach memory, so
// asdt-core/executor-header.md carries no region. skill/embedded_test.go keeps
// a hand-copy of these two literals (package skill cannot import
// internal/installer) — change both copies together.
const (
	memoryBindingBeginMarker = "<!-- ASDT:GENERATED:memory-binding -->"
	memoryBindingEndMarker   = "<!-- /ASDT:GENERATED:memory-binding -->"
)

// Validate reports the first gap that would install a broken binding: no
// Name or ConfigValue, a verb with no tool or a tool with no Name or Usage, an
// empty Claude Code prefix, or a probe Binary with nowhere to point a user who
// lacks it.
func (p ProviderDescriptor) Validate() error {
	if p.Name == "" || p.ConfigValue == "" {
		return fmt.Errorf("memory provider %q has no Name or ConfigValue", p.ID)
	}
	for _, verb := range memoryVerbs {
		tools := p.Tools[verb]
		if len(tools) == 0 {
			return fmt.Errorf("memory provider %q binds no tool to verb %q", p.ID, verb)
		}
		for _, tool := range tools {
			if tool.Name == "" {
				return fmt.Errorf("memory provider %q binds a tool with no Name to verb %q", p.ID, verb)
			}
			if tool.Usage == "" {
				return fmt.Errorf("memory provider %q gives verb %q tool %q no Usage", p.ID, verb, tool.Name)
			}
		}
	}
	for _, prefix := range p.ClaudeToolPrefixes {
		if prefix == "" {
			return fmt.Errorf("memory provider %q has an empty ClaudeToolPrefixes entry", p.ID)
		}
	}
	if p.Detect.Binary != "" && p.Detect.InstallURL == "" {
		return fmt.Errorf("memory provider %q probes %q but has no Detect.InstallURL", p.ID, p.Detect.Binary)
	}
	return nil
}

// renderMemoryBinding emits the binding region body for p — the bytes strictly
// between the memory-binding markers. It fails when p does not Validate: a
// partial binding would leave a prompt with a verb it cannot call.
func renderMemoryBinding(p ProviderDescriptor) (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "\n**Memory provider: %s.** `.asdt/config.yaml` names it as `memory.provider: %s`. Each verb is bound to the tool calls below, made in order:\n\n", p.Name, p.ConfigValue)
	for _, verb := range memoryVerbs {
		calls := make([]string, len(p.Tools[verb]))
		for i, tool := range p.Tools[verb] {
			calls[i] = fmt.Sprintf("`%s` — %s", tool.Name, tool.Usage)
		}
		fmt.Fprintf(&b, "- **%s** → %s.\n", verb, strings.Join(calls, "; then "))
	}

	if len(p.ClaudeToolPrefixes) > 0 {
		example := p.Tools[MemoryRemember][0].Name
		exposed := make([]string, len(p.ClaudeToolPrefixes))
		for i, prefix := range p.ClaudeToolPrefixes {
			exposed[i] = "`" + prefix + example + "`"
		}
		fmt.Fprintf(&b, "\nA host may prefix these names: Claude Code exposes `%s` as %s, and another host may use a different prefix or none. Match the name after the prefix.\n", example, strings.Join(exposed, " or "))
	}
	return b.String(), nil
}

// bindMemory splices p's binding into content's memory-binding region. Content
// without the markers passes through unchanged, so it is safe to run over every
// file the installer writes; partial or duplicated markers fail loudly in
// replaceMarkerRegion. The marker check runs before rendering so a marker-free
// file never depends on the provider being complete.
func bindMemory(content []byte, p ProviderDescriptor) ([]byte, error) {
	if !bytes.Contains(content, []byte(memoryBindingBeginMarker)) &&
		!bytes.Contains(content, []byte(memoryBindingEndMarker)) {
		return content, nil
	}

	body, err := renderMemoryBinding(p)
	if err != nil {
		return nil, err
	}
	bound, err := replaceMarkerRegion(content, memoryBindingBeginMarker, memoryBindingEndMarker, body)
	if err != nil {
		return nil, fmt.Errorf("regenerate memory-binding region: %w", err)
	}
	return bound, nil
}
