package installer

import (
	"bytes"
	"fmt"
	"strings"
)

// ProviderID identifies a known provider: memory or code intelligence.
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

// ProviderTool is one provider tool call bound to a verb: a memory verb's
// recipe step, or a code-intelligence verb's tool.
type ProviderTool struct {
	Name  string // bare tool name, without any host prefix (e.g. "mem_save")
	Usage string // how to call it for this verb, in prompt prose; spliced verbatim
}

// ProviderDescriptor describes a known memory provider: everything the
// installer needs to bind the prompts' memory interface to it. Adding a
// provider is a new entry in Providers — no prompt edits.
//
// The zero ProviderDescriptor is the "no memory" selection: installing with it
// binds no provider (see IsNone). A memory provider is optional; ASDT works
// without one (protocol.md §2, "Memory unavailable").
type ProviderDescriptor struct {
	ID          ProviderID
	Name        string
	Description string
	// ConfigValue is what /asdt-init writes to `memory.provider` in
	// .asdt/config.yaml. protocol.md §2 ("Memory unavailable") treats memory as
	// unavailable when that config value names a different provider. It is
	// never MemoryNone: that value means no provider is bound.
	ConfigValue string
	// Tools binds every verb in memoryVerbs to its recipe: one or more tool
	// calls, made in order (recall may need a search, then a full read).
	Tools map[MemoryVerb][]ProviderTool
	// ClaudeToolPrefixes are the MCP namespaces Claude Code may expose the
	// tools under (plugin install and plain MCP install). The binding names
	// them so the orchestrator matches its own tool list; optional — a
	// provider without one gets no prefix sentence.
	ClaudeToolPrefixes []string
	// Detect is how the setup TUI's environment check looks for the provider.
	Detect ProviderDetect
}

// ProviderDetect is how the setup TUI's environment check looks for a provider
// — memory or code intelligence — before installing.
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
		Tools: map[MemoryVerb][]ProviderTool{
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

// providerNone is the value that stands for "no provider" of either kind. No
// provider's ConfigValue may equal it.
const providerNone = "none"

// MemoryNone is the value that stands for "no memory provider": it is what
// /asdt-init passes as `memory_provider` when the binding names none (and then
// writes no `memory.provider` key), and what install metadata records.
const MemoryNone = providerNone

// IsNone reports whether p is the "no memory" selection: the zero
// ProviderDescriptor. A descriptor with any identity or tool set is a provider,
// and must Validate.
func (p ProviderDescriptor) IsNone() bool {
	return p.ID == "" && p.Name == "" && p.ConfigValue == "" && len(p.Tools) == 0
}

// MetaValue is what install metadata records for p: its ID, or MemoryNone.
func (p ProviderDescriptor) MetaValue() string {
	if p.IsNone() {
		return MemoryNone
	}
	return p.ID
}

// noProviderBinding is the binding region body of a "no memory" install. Its
// opening phrase is what protocol.md §2 ("Memory unavailable") and
// asdt-init/SKILL.md test for: a binding that names no provider.
const noProviderBinding = "\n**No memory provider bound.** This install has no memory: nothing is remembered or recalled, even when a memory tool is on your tool list (`asdt-core/protocol.md` §2, \"Memory unavailable\"). Re-running the ASDT installer with a provider selected connects one.\n"

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

// Validate reports the first gap that would install a broken binding: no ID,
// Name or ConfigValue, a ConfigValue of MemoryNone, a verb with no tool or a
// tool with no Name or Usage, an empty Claude Code prefix, or a probe Binary
// with nowhere to point a user who lacks it.
func (p ProviderDescriptor) Validate() error {
	const kind = "memory provider"
	if err := validateProviderIdentity(kind, p.ID, p.Name, p.ConfigValue, p.ClaudeToolPrefixes, p.Detect); err != nil {
		return err
	}
	for _, verb := range memoryVerbs {
		tools := p.Tools[verb]
		if len(tools) == 0 {
			return fmt.Errorf("%s %q binds no tool to verb %q", kind, p.ID, verb)
		}
		for _, tool := range tools {
			if err := validateProviderTool(kind, p.ID, string(verb), tool); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateProviderIdentity holds the checks every provider kind shares: an ID,
// Name and ConfigValue, a ConfigValue that is not the "none" selection, no empty
// Claude Code prefix, and an install URL for any probed binary. kind names the
// provider kind in the error ("memory provider").
func validateProviderIdentity(kind, id, name, configValue string, prefixes []string, detect ProviderDetect) error {
	if id == "" || name == "" || configValue == "" {
		return fmt.Errorf("%s %q has no ID, Name or ConfigValue", kind, id)
	}
	if configValue == providerNone {
		return fmt.Errorf("%s %q has ConfigValue %q, which means no provider", kind, id, providerNone)
	}
	for _, prefix := range prefixes {
		if prefix == "" {
			return fmt.Errorf("%s %q has an empty ClaudeToolPrefixes entry", kind, id)
		}
	}
	if detect.Binary != "" && detect.InstallURL == "" {
		return fmt.Errorf("%s %q probes %q but has no Detect.InstallURL", kind, id, detect.Binary)
	}
	return nil
}

// validateProviderTool reports a tool bound to verb with no Name or Usage.
func validateProviderTool(kind, id, verb string, tool ProviderTool) error {
	if tool.Name == "" {
		return fmt.Errorf("%s %q binds a tool with no Name to verb %q", kind, id, verb)
	}
	if tool.Usage == "" {
		return fmt.Errorf("%s %q gives verb %q tool %q no Usage", kind, id, verb, tool.Name)
	}
	return nil
}

// renderMemoryBinding emits the binding region body for p — the bytes strictly
// between the memory-binding markers. The "no memory" selection renders
// noProviderBinding. Any other p must Validate: a partial binding would leave a
// prompt with a verb it cannot call.
func renderMemoryBinding(p ProviderDescriptor) (string, error) {
	if p.IsNone() {
		return noProviderBinding, nil
	}
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

	b.WriteString(hostPrefixNote(p.Tools[MemoryRemember][0].Name, p.ClaudeToolPrefixes))
	return b.String(), nil
}

// hostPrefixNote is the binding's closing paragraph naming how Claude Code
// exposes example under each of prefixes, so the reader matches its own tool
// list; empty when the provider declares no prefix.
func hostPrefixNote(example string, prefixes []string) string {
	if len(prefixes) == 0 {
		return ""
	}
	exposed := make([]string, len(prefixes))
	for i, prefix := range prefixes {
		exposed[i] = "`" + prefix + example + "`"
	}
	return fmt.Sprintf("\nA host may prefix these names: Claude Code exposes `%s` as %s, and another host may use a different prefix or none. Match the name after the prefix.\n", example, strings.Join(exposed, " or "))
}

// bindMemory splices p's binding into content's memory-binding region (see
// bindRegion).
func bindMemory(content []byte, p ProviderDescriptor) ([]byte, error) {
	return bindRegion(content, memoryBindingBeginMarker, memoryBindingEndMarker, "memory-binding", func() (string, error) {
		return renderMemoryBinding(p)
	})
}

// bindRegion splices render's body into content's begin/end marker region.
// Content without the markers passes through unchanged, so it is safe to run
// over every file the installer writes; partial or duplicated markers fail
// loudly in replaceMarkerRegion. The marker check runs before rendering so a
// marker-free file never depends on the provider being complete. region names
// the region in errors.
func bindRegion(content []byte, begin, end, region string, render func() (string, error)) ([]byte, error) {
	if !bytes.Contains(content, []byte(begin)) && !bytes.Contains(content, []byte(end)) {
		return content, nil
	}

	body, err := render()
	if err != nil {
		return nil, err
	}
	bound, err := replaceMarkerRegion(content, begin, end, body)
	if err != nil {
		return nil, fmt.Errorf("regenerate %s region: %w", region, err)
	}
	return bound, nil
}
