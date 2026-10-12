package installer

import (
	"fmt"
	"strings"
)

// ProviderCodegraph identifies the Codegraph code-intelligence provider.
const ProviderCodegraph ProviderID = "codegraph"

// CodeIntelVerb names one operation of the provider-neutral code-intelligence
// interface the executor header speaks (skill/asdt-core/executor-header.md,
// "Code intelligence"). Prompts never name a provider's tools; the installer
// binds each verb to the selected provider's tool at install time.
type CodeIntelVerb string

// The code-intelligence interface. explore answers a question with the
// relevant source grouped by file; search locates a symbol; read returns one
// definition in full; callers and callees walk the call graph; impact names
// what breaks if a symbol changes.
const (
	CodeIntelExplore CodeIntelVerb = "explore"
	CodeIntelSearch  CodeIntelVerb = "search"
	CodeIntelRead    CodeIntelVerb = "read"
	CodeIntelCallers CodeIntelVerb = "callers"
	CodeIntelCallees CodeIntelVerb = "callees"
	CodeIntelImpact  CodeIntelVerb = "impact"
)

// codeIntelVerbs is the canonical verb order: the order the binding lists them
// and the Claude Code grant appends their tools.
var codeIntelVerbs = []CodeIntelVerb{CodeIntelExplore, CodeIntelSearch, CodeIntelRead, CodeIntelCallers, CodeIntelCallees, CodeIntelImpact}

// CodeIntelDescriptor describes a known code-intelligence provider: everything
// the installer needs to bind the prompts' code-intelligence interface to it,
// grant its tools, and let /asdt-init detect an indexed project. Adding a
// provider is a new entry in CodeIntelProviders — no prompt edits.
//
// The zero CodeIntelDescriptor is the "no code intelligence" selection:
// installing with it binds no provider (see IsNone), and every step answers
// symbol and structure questions with grep and read.
type CodeIntelDescriptor struct {
	ID          ProviderID
	Name        string
	Description string
	// ConfigValue is what /asdt-init writes to `code_intelligence` in
	// .asdt/config.yaml when it finds IndexMarker; the platform block's
	// `Tooling:` line carries it to every step. Never CodeIntelNone.
	ConfigValue string
	// IndexMarker is the repository-root path whose presence proves the
	// provider has indexed the project (a directory or file name, no glob).
	IndexMarker string
	// Tools binds verbs to the provider's tools. explore is required: the
	// executor header starts every structure question with it and
	// /asdt-init's enrichment calls it. A provider may leave any other verb
	// unbound; prompts call only the verbs the binding lists.
	Tools map[CodeIntelVerb]ProviderTool
	// ClaudeToolPrefixes are the MCP namespaces Claude Code exposes the tools
	// under. Each bound tool is granted to the executor agents under every
	// prefix; a provider without one gets no Claude Code grant.
	ClaudeToolPrefixes []string
	// Detect is how the setup TUI's environment check looks for the provider.
	Detect ProviderDetect
}

// CodeIntelProviders lists all known code-intelligence providers.
var CodeIntelProviders = []CodeIntelDescriptor{
	{
		ID:          ProviderCodegraph,
		Name:        "Codegraph",
		Description: "Semantic code index served as read-only MCP tools.",
		ConfigValue: "codegraph",
		IndexMarker: ".codegraph",
		Tools: map[CodeIntelVerb]ProviderTool{
			CodeIntelExplore: {Name: "codegraph_explore", Usage: "ask the question in plain words; one call returns the relevant source grouped by file"},
			CodeIntelSearch:  {Name: "codegraph_search", Usage: "pass a symbol name, or part of one; returns where each match is defined"},
			CodeIntelRead:    {Name: "codegraph_node", Usage: "pass a symbol name; returns its full definition"},
			CodeIntelCallers: {Name: "codegraph_callers", Usage: "pass a symbol name; returns what calls it"},
			CodeIntelCallees: {Name: "codegraph_callees", Usage: "pass a symbol name; returns what it calls"},
			CodeIntelImpact:  {Name: "codegraph_impact", Usage: "pass a symbol name; returns what depends on it, directly and transitively"},
		},
		ClaudeToolPrefixes: []string{"mcp__codegraph__"},
		Detect: ProviderDetect{
			Binary:     "codegraph",
			InstallURL: "https://github.com/Gentleman-Programming/codegraph",
		},
	},
}

// CodeIntelNone is what install metadata records for an install with no
// code-intelligence provider bound.
const CodeIntelNone = providerNone

// IsNone reports whether p is the "no code intelligence" selection: the zero
// CodeIntelDescriptor. A descriptor with any identity or tool set is a
// provider, and must Validate.
func (p CodeIntelDescriptor) IsNone() bool {
	return p.ID == "" && p.Name == "" && p.ConfigValue == "" && p.IndexMarker == "" && len(p.Tools) == 0
}

// MetaValue is what install metadata records for p: its ID, or CodeIntelNone.
func (p CodeIntelDescriptor) MetaValue() string {
	if p.IsNone() {
		return CodeIntelNone
	}
	return p.ID
}

// Validate reports the first gap that would install a broken binding: the
// identity gaps every provider kind shares (validateProviderIdentity), no
// IndexMarker or one that is a path rather than a root entry name, no explore
// tool, or a bound tool with no Name or Usage.
func (p CodeIntelDescriptor) Validate() error {
	const kind = "code-intelligence provider"
	if err := validateProviderIdentity(kind, p.ID, p.Name, p.ConfigValue, p.ClaudeToolPrefixes, p.Detect); err != nil {
		return err
	}
	if p.IndexMarker == "" || strings.ContainsAny(p.IndexMarker, "/*?[ ") {
		return fmt.Errorf("%s %q has IndexMarker %q, want one repository-root entry name", kind, p.ID, p.IndexMarker)
	}
	if _, ok := p.Tools[CodeIntelExplore]; !ok {
		return fmt.Errorf("%s %q binds no tool to verb %q", kind, p.ID, CodeIntelExplore)
	}
	for verb, tool := range p.Tools {
		if !isCodeIntelVerb(verb) {
			return fmt.Errorf("%s %q binds unknown verb %q", kind, p.ID, verb)
		}
		if err := validateProviderTool(kind, p.ID, string(verb), tool); err != nil {
			return err
		}
	}
	return nil
}

func isCodeIntelVerb(verb CodeIntelVerb) bool {
	for _, v := range codeIntelVerbs {
		if v == verb {
			return true
		}
	}
	return false
}

// claudeTools returns the Claude Code tool names granted to the executor
// agents for p: every bound tool under every ClaudeToolPrefixes entry, in
// verb order. The "no code intelligence" selection grants nothing.
func (p CodeIntelDescriptor) claudeTools() []string {
	var tools []string
	for _, verb := range codeIntelVerbs {
		tool, ok := p.Tools[verb]
		if !ok {
			continue
		}
		for _, prefix := range p.ClaudeToolPrefixes {
			tools = append(tools, prefix+tool.Name)
		}
	}
	return tools
}

// noCodeIntelBinding is the binding region body of a "no code intelligence"
// install. Its opening phrase is what executor-header.md and asdt-init/SKILL.md
// test for: a binding that names no provider.
const noCodeIntelBinding = "\n**No code-intelligence provider bound.** This install has none: answer symbol and structure questions with grep and read, even when a code-intelligence tool is on your tool list. Re-running the ASDT installer with a provider selected connects one.\n"

// Code-intelligence-binding generated-region markers. They bound the provider
// binding in every skill file that calls a code-intelligence verb:
// asdt-core/executor-header.md (every sub-agent, and through it
// asdt-init/explore's index probe) and asdt-init/SKILL.md (enrichment).
// skill/embedded_test.go keeps a hand-copy of these two literals (package skill
// cannot import internal/installer) — change both copies together.
const (
	codeIntelBindingBeginMarker = "<!-- ASDT:GENERATED:code-intel-binding -->"
	codeIntelBindingEndMarker   = "<!-- /ASDT:GENERATED:code-intel-binding -->"
)

// renderCodeIntelBinding emits the binding region body for p — the bytes
// strictly between the code-intel-binding markers. The "no code intelligence"
// selection renders noCodeIntelBinding. Any other p must Validate.
func renderCodeIntelBinding(p CodeIntelDescriptor) (string, error) {
	if p.IsNone() {
		return noCodeIntelBinding, nil
	}
	if err := p.Validate(); err != nil {
		return "", err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "\n**Code-intelligence provider: %s.** Its index marker is `%s` at the repository root. `/asdt-init` records an indexed project in `.asdt/config.yaml` as `code_intelligence: %s`, which reaches a step's platform block as the line `Tooling: %s index available`. Each verb it serves is bound to one tool:\n\n", p.Name, p.IndexMarker, p.ConfigValue, p.ConfigValue)
	for _, verb := range codeIntelVerbs {
		if tool, ok := p.Tools[verb]; ok {
			fmt.Fprintf(&b, "- **%s** → `%s` — %s.\n", verb, tool.Name, tool.Usage)
		}
	}
	b.WriteString(hostPrefixNote(p.Tools[CodeIntelExplore].Name, p.ClaudeToolPrefixes))
	return b.String(), nil
}

// bindCodeIntel splices p's binding into content's code-intel-binding region
// (see bindRegion).
func bindCodeIntel(content []byte, p CodeIntelDescriptor) ([]byte, error) {
	return bindRegion(content, codeIntelBindingBeginMarker, codeIntelBindingEndMarker, "code-intel-binding", func() (string, error) {
		return renderCodeIntelBinding(p)
	})
}
