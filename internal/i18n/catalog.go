// Package i18n provides a simple locale-based string catalog for the TUI.
// Add new feature areas as nested structs inside Catalog when they appear.
package i18n

// Catalog holds all translatable strings, organized by feature area so it
// scales naturally as new TUI surfaces (dashboard, settings) are added.
type Catalog struct {
	Installer InstallerStrings
	Dashboard DashboardStrings
	Personas  PersonaStrings
}

// PersonaStrings holds the localized one-line descriptions for the built-in
// persona presets shown in StateAgentSetup and StateReview. Persona names are
// proper nouns and never translated — only their descriptions are.
type PersonaStrings struct {
	Sky         string
	Toffy       string
	Atreus      string
	Babi        string
	LeePalacios string
}

// PersonaDescription returns the localized description for the given persona
// preset ID. An empty field falls back to the English catalog's same field;
// an unknown ID returns "".
func (c Catalog) PersonaDescription(presetID string) string {
	pick := func(p PersonaStrings) string {
		switch presetID {
		case "sky":
			return p.Sky
		case "toffy":
			return p.Toffy
		case "atreus":
			return p.Atreus
		case "babi":
			return p.Babi
		case "lee-palacios":
			return p.LeePalacios
		}
		return ""
	}
	if s := pick(c.Personas); s != "" {
		return s
	}
	return pick(English.Personas)
}

// InstallerStrings holds every user-visible string in the wizard TUI.
type InstallerStrings struct {
	// Frame titles
	TitleMainMenu         string
	TitlePreflightCheck   string
	TitleSelectAssistants string
	TitleSelectProvider   string
	TitleModelGate        string
	TitleModelSetup       string
	TitleAgentSetup       string
	TitleEmojiPref        string
	TitleAgentWriteMode   string
	TitleReview           string
	TitleInstalling       string
	TitleDone             string
	TitleDashboard        string
	TitleLanguageSelect   string

	// Step indicator words ("step 1 of 5" / "paso 1 de 5")
	StepWord   string
	StepOfWord string

	// Main menu items
	MenuInstall   string
	MenuDashboard string
	MenuQuit      string

	// Keyboard hint descriptions (what a key does)
	HintNavigate     string
	HintSelect       string
	HintQuit         string
	HintToggle       string
	HintAllNone      string
	HintBack         string
	HintBackToMenu   string
	HintCycleModel   string
	HintExpand       string
	HintResetDefault string
	HintContinue     string
	HintChecking     string
	HintEnvironment  string

	// Hint group labels
	HintGroupNav     string
	HintGroupActions string
	HintGroupStatus  string

	// Inline confirm buttons (full label including brackets)
	BtnContinue string
	BtnSkip     string
	BtnInstall  string

	// Body / subtitle text
	BodyAgentSetupSubtitle string
	BodyEmojiPrefSubtitle  string
	BodyModelSetupSubtitle string // printf format: %s = detected provider names

	// Model gate options (radio rows) and per-step screen extras. The six
	// OptionModelsPreset rows are the gate choices (0 = Chameleon … 4 =
	// Mastermind, 5 = customize per step); each focused row shows its matching
	// HintPresetAssignment. All copy is provider-agnostic — no model IDs.
	OptionModelsPreset0        string
	OptionModelsPreset1        string
	OptionModelsPreset2        string
	OptionModelsPreset3        string
	OptionModelsPreset4        string
	OptionModelsPreset5        string
	HintPresetAssignment0      string
	HintPresetAssignment1      string
	HintPresetAssignment2      string
	HintPresetAssignment3      string
	HintPresetAssignment4      string
	HintPresetAssignment5      string
	HintRunsAs                 string // printf format: %s = Claude Code model enum
	LabelModels                string
	ValueModelsPreset          string // printf format: %s = preset short name
	ValueModelsCustomized      string // printf format: %d = customized step count
	BodyAgentWriteMode         string
	BodyInstalling             string
	BodyLanguageSelectSubtitle string

	// Emoji preference options (radio rows + Review value)
	OptionEmojiYes string
	OptionEmojiNo  string

	// AgentSetup state messages
	WarnExistingConfig     string
	WarnExistingConfigNote string
	InfoWriteGlobal        string

	// Review screen labels
	LabelAssistants string
	LabelMemory     string
	LabelCodeIntel  string
	LabelPersona    string
	LabelEmojis     string
	LabelConfig     string
	LabelSkipped    string

	// Agent write mode labels for an existing global config file: update
	// ASDT's marked block, or leave the file untouched.
	LabelModeUpdate string
	LabelModeSkip   string

	// Done screen
	LabelAgentConfig      string
	MsgAgentConfigWritten string // printf format: first %s = assistant name
	MsgAgentConfigSkipped string // printf format: first %s = assistant name
	MsgOpenCodeNote       string
	MsgGetStarted         string
	MsgStaleRemoved       string // printf format: %d = count of pruned stale files

	// Preflight section group titles (display only; row labels are lookup keys)
	SectionYourEnvironment  string
	SectionAssistants       string
	SectionMemoryProvider   string
	SectionCodeIntelligence string

	// Assistant detection. The preflight row of an assistant whose binary is
	// not on PATH is a soft warning (DetailAssistantNotFound); on the
	// selection screen each assistant reads one of three states — ASDT
	// already installed, binary detected (new install), or binary not on
	// PATH — as a badge plus a short note, and a selected assistant that is
	// not on PATH gets WarnAssistantMissing under the list. None blocks.
	DetailAssistantNotFound string // printf format: %s = binary name
	DescSelectAssistants    string
	BadgeAssistantInstalled string
	NoteAssistantInstalled  string
	BadgeAssistantDetected  string
	NoteAssistantDetected   string
	BadgeAssistantNotFound  string
	NoteAssistantNotFound   string
	WarnAssistantMissing    string // printf format: %s = assistant Name, %s = binary name

	// Memory provider selection. A provider is optional: the "No memory" row
	// sits after every installer.Providers entry, and a selected provider whose
	// probe binary is not on PATH shows a warning with its install hint but
	// never blocks Continue. The Review row reads ValueMemoryNone for "No
	// memory" and ValueMemoryMissing for an undetected provider.
	OptionNoMemory      string
	DescNoMemory        string
	WarnProviderMissing string // printf format: %s = provider Name
	PrefProviderInstall string // printf format: %s = provider Detect.InstallURL
	// DetailProviderNotFound is the preflight row detail for a provider whose
	// probe binary is not on PATH (a soft warning, never a block).
	DetailProviderNotFound string // printf format: %s = provider Detect.InstallURL
	ValueMemoryNone        string // printf format: %s = provider Names, joined by " / "
	ValueMemoryMissing     string // printf format: %s = provider Name, %s = provider Detect.InstallURL

	// Code-intelligence provider selection, on the same screen as memory. A
	// provider is optional: the "No code intelligence" row sits after every
	// installer.CodeIntelProviders entry, and a selected provider whose probe
	// binary is not on PATH shows a warning with its install hint
	// (PrefProviderInstall) but never blocks Continue.
	OptionNoCodeIntel       string
	DescNoCodeIntel         string
	WarnCodeIntelMissing    string // printf format: %s = provider Name
	DetailCodeIntelNotFound string // printf format: %s = provider Detect.InstallURL
	ValueCodeIntelNone      string
	ValueCodeIntelMissing   string // printf format: %s = provider Name, %s = provider Detect.InstallURL
}

// DashboardStrings holds user-visible strings for the dashboard TUI.
type DashboardStrings struct {
	LabelInstalled string // "Installed" — prefix for the install date
	LabelPersona   string // "Persona" — prefix for the configured persona name
	LabelMemory    string // "Memory" — prefix for the bound memory provider
	ValueNoMemory  string // shown after LabelMemory for a "no memory" install
	// LabelCodeIntel prefixes the bound code-intelligence provider;
	// ValueNoCodeIntel follows it for an install with none.
	LabelCodeIntel   string
	ValueNoCodeIntel string

	// Assistant status badges: whether the assistant's CLI is on PATH and
	// whether ASDT's skills are installed for it.
	BadgeCLIFound      string
	BadgeCLIMissing    string
	BadgeASDTInstalled string
	BadgeASDTMissing   string

	// Agent config block: the badge on an assistant whose global config file
	// carries ASDT's block, and the remove action.
	BadgeAgentConfig         string
	HintRemoveAgentConfig    string
	ConfirmRemoveAgentConfig string // printf format: %s = config file paths, joined by ", "
	HintConfirm              string
	HintCancel               string
	MsgAgentConfigRemoved    string // printf format: %s = assistant name
	MsgAgentConfigNoBlock    string // printf format: %s = assistant name
}
