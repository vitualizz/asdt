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
// prompts speak (skill/asdt-core/protocol.md §0). Prompts never name a
// provider's tools; the installer binds each verb to the selected provider's
// tool at install time.
type MemoryVerb string

// The memory interface. save upserts a record by stable key; search finds
// records by key prefix or free text; get returns one full record by id.
const (
	MemorySave   MemoryVerb = "save"
	MemorySearch MemoryVerb = "search"
	MemoryGet    MemoryVerb = "get"
)

// memoryVerbs is the canonical verb order: the order the binding lists them and
// the order each namespace's tools are granted to the executor agents.
var memoryVerbs = []MemoryVerb{MemorySave, MemorySearch, MemoryGet}

// MemoryTool is one provider tool bound to a memory verb.
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
	// .asdt/config.yaml, and what every specialist's Prerequisites check reads.
	ConfigValue string
	// Tools binds every verb in memoryVerbs to the provider's tool.
	Tools map[MemoryVerb]MemoryTool
	// ClaudeToolPrefixes are the MCP namespaces Claude Code may expose the
	// tools under (plugin install and plain MCP install). The executor agents'
	// Claude Code allowlist is every prefix × every verb's tool, so at least
	// one is required: none would install agents that cannot reach memory.
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
		Tools: map[MemoryVerb]MemoryTool{
			MemorySave: {
				Name:  "mem_save",
				Usage: "pass the key as `topic_key`, the key without its leading `{project}/` as `title` (a hand-off's title is `{change}/{role}/handoff`), `type: \"decision\"`, and the record as `content`",
			},
			MemorySearch: {
				Name:  "mem_search",
				Usage: "pass the key prefix or the free text as `query`. Each match shows the record's key as its `title`, without the leading `{project}/` (key `{project}/{change}/developer/handoff` → title `{change}/developer/handoff`): match keys and prefixes against `title` with that segment removed",
			},
			MemoryGet: {
				Name:  "mem_get_observation",
				Usage: "pass the id a search result carries as `id`",
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
// every skill file that reaches memory (asdt-core/protocol.md, and through it
// every routed SKILL.md; asdt-core/executor-header.md; asdt-init/SKILL.md; the
// root SKILL.md). skill/embedded_test.go keeps a hand-copy of these two
// literals (package skill cannot import internal/installer) — change both
// copies together.
const (
	memoryBindingBeginMarker = "<!-- ASDT:GENERATED:memory-binding -->"
	memoryBindingEndMarker   = "<!-- /ASDT:GENERATED:memory-binding -->"
)

// claudeMemoryTools returns the Claude Code tool names the executor agents are
// granted for this provider: for each prefix, every verb's tool in
// memoryVerbs order. Derived from Tools, so the grant can never drift from the
// binding the prompts read.
func (p ProviderDescriptor) claudeMemoryTools() []string {
	tools := make([]string, 0, len(p.ClaudeToolPrefixes)*len(memoryVerbs))
	for _, prefix := range p.ClaudeToolPrefixes {
		for _, verb := range memoryVerbs {
			tools = append(tools, prefix+p.Tools[verb].Name)
		}
	}
	return tools
}

// Validate reports the first gap that would install a broken binding: no
// Name or ConfigValue, a verb with no tool Name or Usage, no Claude Code
// prefix (the executor agents would be granted no memory tool), or a probe
// Binary with nowhere to point a user who lacks it.
func (p ProviderDescriptor) Validate() error {
	if p.Name == "" || p.ConfigValue == "" {
		return fmt.Errorf("memory provider %q has no Name or ConfigValue", p.ID)
	}
	for _, verb := range memoryVerbs {
		tool, ok := p.Tools[verb]
		if !ok || tool.Name == "" {
			return fmt.Errorf("memory provider %q binds no tool to verb %q", p.ID, verb)
		}
		if tool.Usage == "" {
			return fmt.Errorf("memory provider %q gives verb %q no Usage", p.ID, verb)
		}
	}
	if len(p.ClaudeToolPrefixes) == 0 {
		return fmt.Errorf("memory provider %q has no ClaudeToolPrefixes", p.ID)
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
	fmt.Fprintf(&b, "\n**Memory provider: %s.** `.asdt/config.yaml` names it as `memory.provider: %s`. Each verb is bound to exactly one tool:\n\n", p.Name, p.ConfigValue)
	for _, verb := range memoryVerbs {
		tool := p.Tools[verb]
		fmt.Fprintf(&b, "- **%s** → `%s` — %s.\n", verb, tool.Name, tool.Usage)
	}

	example := p.Tools[MemorySave].Name
	exposed := make([]string, len(p.ClaudeToolPrefixes))
	for i, prefix := range p.ClaudeToolPrefixes {
		exposed[i] = "`" + prefix + example + "`"
	}
	fmt.Fprintf(&b, "\nA host may prefix these names: Claude Code exposes `%s` as %s, and another host may use a different prefix or none. Match the name after the prefix.\n", example, strings.Join(exposed, " or "))
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
