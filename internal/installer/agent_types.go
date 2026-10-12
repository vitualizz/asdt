package installer

// PersonaPreset is a selectable agent persona shipped with the installer.
type PersonaPreset struct {
	ID          string // matches persona filename stem: sky|toffy|atreus|babi|lee-palacios
	Name        string // "Sky"
	Description string // one-line summary shown in StateAgentSetup
	File        string // embedded asset path: "assets/personas/sky.md"
}

// PersonaPresets lists the built-in agent persona presets in display order.
var PersonaPresets = []PersonaPreset{
	{
		ID:          "sky",
		Name:        "Sky",
		Description: "Sharp and precise. Roasts bad practices with wit, then explains what matters.",
		File:        "assets/personas/sky.md",
	},
	{
		ID:          "toffy",
		Name:        "Toffy",
		Description: "Warm and enthusiastic. Simple answers, accidental sarcasm.",
		File:        "assets/personas/toffy.md",
	},
	{
		ID:          "atreus",
		Name:        "Atreus",
		Description: "Bold and reckless. Ships the unconventional solution first — and it works.",
		File:        "assets/personas/atreus.md",
	},
	{
		ID:          "babi",
		Name:        "Babi",
		Description: "Your biggest fan. Warm, fast, and celebrates everything loudly.",
		File:        "assets/personas/babi.md",
	},
	{
		ID:          "lee-palacios",
		Name:        "Lee Palacios",
		Description: "Cat lover, coder, otaku. Your pair-programming companion for the whole journey.",
		File:        "assets/personas/lee-palacios.md",
	},
}

// AgentWriteMode controls what the installer does with ASDT's block in an
// assistant's global config file.
type AgentWriteMode int

const (
	// AgentModeUpdate writes ASDT's block: it replaces the block when present,
	// else adds it at the end. Text outside the block is never touched.
	AgentModeUpdate AgentWriteMode = iota
	// AgentModeSkip leaves the file untouched.
	AgentModeSkip
)

// AgentConfigResult is the per-assistant outcome of writing or removing the
// agent config block.
type AgentConfigResult struct {
	AssistantID AssistantID
	Written     []string // config files written (CLAUDE.md, AGENTS.md)
	Skipped     bool     // true when nothing was written (skip chosen, no adapter, or no block to remove)
	Err         error
}

// AgentConfigAdapterDescriptor declares where one assistant keeps its global
// config file, the file that carries ASDT's marker block.
type AgentConfigAdapterDescriptor struct {
	AssistantID AssistantID
	ConfigPath  func() string // absolute path, also shown in the TUI
}
