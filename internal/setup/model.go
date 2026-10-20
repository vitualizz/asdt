package setup

import (
	"io/fs"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vitualizz/asdt/internal/i18n"
	"github.com/vitualizz/asdt/internal/installer"
	"github.com/vitualizz/asdt/internal/setup/components"
	"github.com/vitualizz/asdt/internal/tui/panels"
)

// ViewState represents the current screen shown by the TUI.
type ViewState int

const (
	// StateMainMenu is the zero-value state; the TUI starts here showing the
	// top-level menu. The user chooses between installing/updating skills,
	// viewing the dashboard, or quitting.
	StateMainMenu ViewState = iota
	// StateEnvironmentCheck is shown after the user selects "Install / Update
	// Skills". It fans out environment probes (OS, Shell, every memory
	// provider in installer.Providers, every code-intelligence provider in
	// installer.CodeIntelProviders) concurrently; Continue is gated only until
	// all probes resolve. A missing provider is a warning, never a block:
	// memory and code intelligence are both optional.
	StateEnvironmentCheck
	// StateDashboard shows the dashboard overview with per-assistant status.
	StateDashboard
	// StateSelectAssistants allows the user to choose assistants to install.
	// Shows each assistant's state (ASDT installed / detected / not in PATH)
	// as a badge alongside its checkbox.
	StateSelectAssistants
	// StateSelectProvider allows the user to choose a memory provider (or "No
	// memory") and a code-intelligence provider (or "No code intelligence"):
	// two radio groups on one screen. It never gates Continue: an undetected
	// provider shows its install hint as a warning.
	StateSelectProvider
	// StateAgentSetup shows the persona selection step (optional).
	StateAgentSetup
	// StateEmojiPref asks whether the assistants should use emojis. Only
	// reached on the non-skip path out of StateAgentSetup.
	StateEmojiPref
	// StateAgentWriteMode is shown when a selected assistant's global config
	// file already exists; the user toggles per file (space) between updating
	// ASDT's block and skipping the file.
	StateAgentWriteMode
	// StateReview shows a summary of all selections before installation begins.
	StateReview
	// StateInstalling shows installation progress.
	StateInstalling
	// StateDone is the terminal state shown after installation completes.
	StateDone
	// StateLanguageSelect asks which language the installer (and the installed
	// experience) should use. Entered from MainMenu's Install action, before
	// the environment check. Appended at the END of the iota block so existing
	// state values never shift.
	StateLanguageSelect
	// StateModelSetup lets the user choose which AI model each specialist
	// subagent step runs on, cycling through models of detected providers.
	// Reached only via StateModelGate's customize choice. Appended at the
	// END of the iota block so existing state values never shift.
	StateModelSetup
	// StateModelGate offers five model presets plus a "customize per step"
	// choice. Sits between SelectProvider and AgentSetup; the preset choices
	// (0-4, default Chameleon) bypass StateModelSetup entirely, while choice 5
	// opens it. Appended at the END of the iota block so existing state values
	// never shift.
	StateModelGate
)

// preflightState holds state for the StateEnvironmentCheck screen.
type preflightState struct {
	sections        []components.SectionGroup
	done            bool
	assistantsFound map[installer.AssistantID]bool // see EnvironmentCheckMsg
	providersFound  map[installer.ProviderID]bool  // see EnvironmentCheckMsg
	codeIntelFound  map[installer.ProviderID]bool  // see EnvironmentCheckMsg
}

// wizardState holds install-wizard state across SelectAssistants → Installing → Done.
type wizardState struct {
	selected map[int]bool
	// provider indexes installer.Providers; noMemoryChoice() (one past the
	// last provider) is "No memory". The environment check sets it to the
	// first provider found, else "No memory".
	provider int
	// codeIntel indexes installer.CodeIntelProviders; noCodeIntelChoice() is
	// "No code intelligence". Defaulted like provider.
	codeIntel int
	// providerChosen is set once the user reaches the provider screen; a
	// late EnvironmentCheckMsg then leaves provider and codeIntel alone.
	providerChosen  bool
	results         []installer.InstallResult
	agentResults    []installer.AgentConfigResult
	installExpected int  // number of per-assistant install Cmds fired
	agentDone       bool // true once AgentInstallDoneMsg arrives (or skip)

	// Model selection state. modelGateChoice is the StateModelGate radio
	// (0 = Chameleon, 1 = Sprinter, 2 = Craftsman, 3 = Strategist,
	// 4 = Mastermind, 5 = customize per step), defaulting to Chameleon.
	// modelSteps lists every subagent step from the embedded workflows;
	// selectedModels tracks the current choice per step key
	// ("{specialist}/{step}") — prefilled with the source defaults on the
	// customize path, written from a preset's tier mapping by applyPreset, or
	// left empty for Chameleon (which strips every model field at install).
	// modelOptions is the flattened cycle list from detected providers;
	// expandedGroups tracks which specialist accordion groups are open. A
	// selection map identical to the source defaults installs workflow.yaml
	// files verbatim (skip path) — see countCustomizedModels.
	modelGateChoice int
	modelSteps      []installer.WorkflowStepModel
	selectedModels  map[string]string
	detectedAI      []installer.AIProvider
	modelOptions    []installer.AIModel
	expandedGroups  map[string]bool
}

// dashboardState holds what the dashboard shows, loaded when the user enters
// StateDashboard so the render path never performs I/O. hasBlock marks the
// assistants whose global config file carries ASDT's block; confirmRemove is
// set while the remove action awaits confirmation; removeResults holds the
// outcome of the last removal.
type dashboardState struct {
	meta          map[installer.AssistantID]installer.InstallMeta
	hasBlock      map[installer.AssistantID]bool
	confirmRemove bool
	removeResults []installer.AgentConfigResult
}

// languageState holds state for the StateLanguageSelect screen.
type languageState struct {
	selected int  // index into languageOptions
	touched  bool // true once the user moved the selection; guards late LanguagePrefMsg
}

// languageOptions lists the selectable installer locales, sourced from the
// single installer.SupportedLocales registry so the TUI selector, the persisted
// InstallMeta.Language code, and the injected {{language_directive}} all agree.
// Labels are the locales' NATIVE names — deliberately not catalog strings, so
// each option is readable regardless of the currently active language.
var languageOptions = installer.SupportedLocales

// languageIndex returns the languageOptions index for code, or -1 when the
// code has no option.
func languageIndex(code string) int {
	for i, opt := range languageOptions {
		if opt.Code == code {
			return i
		}
	}
	return -1
}

// agentConfigState holds per-assistant agent config state across AgentSetup,
// EmojiPref, and AgentWriteMode.
type agentConfigState struct {
	selectedPersona int
	skip            bool
	useEmojis       bool
	writeModes      map[string]installer.AgentWriteMode
	conflicts       []string
}

// Model is the root Bubbletea model for the installer TUI.
type Model struct {
	state   ViewState
	cursor  int
	width   int
	spinner spinner.Model

	skillsFS        fs.FS
	catalog         i18n.Catalog
	currentVersion  string
	latestVersion   string
	updateAvailable bool

	preflight   preflightState
	wizard      wizardState
	agentConfig agentConfigState
	dashboard   dashboardState
	language    languageState
}

// New constructs an initial Model with the running binary version. Init()
// fires UpdateCheckCmd to passively detect newer releases and starts the
// spinner tick. EnvironmentCheckCmd is NOT fired in Init() — it is triggered
// on demand when the user selects "Install / Update Skills" from StateMainMenu.
//
// The spinner is built via panels.NewSpinner — exported from the panels
// package specifically so this installer TUI and the dashboard TUI share one
// indeterminate-spinner construction (Dot frames, ColorSecondary tint)
// instead of each hand-rolling its own and risking drift.
func New(skillsFS fs.FS, version string) Model {
	str := i18n.Active()
	// ActiveCode() now returns the full canonical code (en/es-419/es-ES), which
	// matches a SupportedLocales row directly — no base-to-full mapping needed.
	active := i18n.ActiveCode()
	selected := languageIndex(active)
	if selected < 0 {
		selected = 0
	}
	return Model{
		wizard:         wizardState{selected: make(map[int]bool)},
		skillsFS:       skillsFS,
		currentVersion: version,
		width:          80,
		spinner:        panels.NewSpinner(),
		preflight:      preflightState{sections: initialPreflightSections(str.Installer)},
		catalog:        str,
		language:       languageState{selected: selected},
	}
}

// State returns the current ViewState. Exported for tests.
func (m Model) State() ViewState { return m.state }

// Init implements tea.Model. Fires the update check and starts the spinner tick.
// EnvironmentCheckCmd is NOT fired here — it is triggered on demand when the
// user selects "Install / Update Skills" from the main menu.
func (m Model) Init() tea.Cmd {
	return tea.Batch(UpdateCheckCmd(m.currentVersion), m.spinner.Tick)
}

// Update handles all messages and key events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case EnvironmentCheckProgressMsg:
		progress := components.CheckResult{
			Label:    msg.RowLabel,
			Status:   msg.Status,
			Detail:   msg.Detail,
			SoftWarn: msg.SoftWarn,
		}
		m.preflight.sections = components.UpdateRow(m.preflight.sections, msg.RowLabel, progress)
		if allProbesDone(m.preflight.sections) {
			assistantsFound := make(map[installer.AssistantID]bool, len(installer.Descriptors))
			for _, d := range installer.Descriptors {
				assistantsFound[d.ID] = !rowHasStatus(m.preflight.sections, d.Name, components.CheckStatusWarning)
			}
			providersFound := make(map[installer.ProviderID]bool, len(installer.Providers))
			for _, p := range installer.Providers {
				providersFound[p.ID] = !rowHasStatus(m.preflight.sections, p.Name, components.CheckStatusWarning)
			}
			codeIntelFound := make(map[installer.ProviderID]bool, len(installer.CodeIntelProviders))
			for _, p := range installer.CodeIntelProviders {
				codeIntelFound[p.ID] = !rowHasStatus(m.preflight.sections, p.Name, components.CheckStatusWarning)
			}
			return m, func() tea.Msg {
				return EnvironmentCheckMsg{AssistantsFound: assistantsFound, ProvidersFound: providersFound, CodeIntelFound: codeIntelFound}
			}
		}
		return m, nil

	case EnvironmentCheckMsg:
		m.preflight.done = true
		m.preflight.assistantsFound = msg.AssistantsFound
		m.preflight.providersFound = msg.ProvidersFound
		m.preflight.codeIntelFound = msg.CodeIntelFound
		if !m.wizard.providerChosen {
			m.wizard.provider = m.defaultProviderChoice()
			m.wizard.codeIntel = m.defaultCodeIntelChoice()
		}
		return m, nil

	case LanguagePrefMsg:
		// Apply the persisted language only while the user hasn't moved the
		// selection — a late-arriving message must never override an explicit
		// choice.
		if m.language.touched {
			return m, nil
		}
		// msg.Code may be a bare "es"; map it to a full locale row. The full
		// code feeds i18n.ForCode directly — normalize() handles non-canonical
		// shapes inside the i18n package.
		loc := installer.LocaleByCode(msg.Code)
		if idx := languageIndex(loc.Code); idx >= 0 {
			m.language.selected = idx
			m.catalog = i18n.ForCode(loc.Code)
			if m.state == StateLanguageSelect {
				m.cursor = idx
			}
		}
		return m, nil

	case UpdateCheckMsg:
		if msg.Err != nil {
			return m, nil // silent: no banner on any error/timeout/rate-limit
		}
		if newerAvailable(msg.Current, msg.Latest) {
			m.latestVersion = msg.Latest
			m.updateAvailable = true
		}
		return m, nil

	case AssistantInstallProgressMsg:
		m.wizard.results = append(m.wizard.results, msg.Result)
		if m.allInstallsDone() {
			m.state = StateDone
		}
		return m, nil

	case InstallDoneMsg:
		// Kept for test compatibility: direct sends of InstallDoneMsg still
		// transition to StateDone without going through per-assistant tracking.
		m.wizard.results = msg.Results
		m.state = StateDone
		return m, nil

	case AgentInstallDoneMsg:
		m.wizard.agentResults = msg.Results
		m.wizard.agentDone = true
		if m.allInstallsDone() {
			m.state = StateDone
		}
		return m, nil

	case AgentConfigRemoveDoneMsg:
		m.dashboard = loadDashboardState()
		m.dashboard.removeResults = msg.Results
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case spinner.TickMsg:
		// Gated to StateInstalling and StateEnvironmentCheck: the spinner has
		// no visual representation in any other screen.
		if m.state != StateInstalling && m.state != StateEnvironmentCheck {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.state == StateEnvironmentCheck {
			// Advance per-row spinner frames for still-pending rows.
			frame := m.spinner.View()
			for i := range m.preflight.sections {
				for j := range m.preflight.sections[i].Rows {
					m.preflight.sections[i].Rows[j].TickSpinner(frame)
				}
			}
		}
		return m, cmd

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// StateEnvironmentCheck handles its own keys exclusively.
	if m.state == StateEnvironmentCheck {
		return m.handleEnvironmentCheck(msg)
	}

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyRunes:
		if string(msg.Runes) == "q" {
			return m, tea.Quit
		}
	}

	switch m.state {
	case StateDashboard:
		return m.handleDashboard(msg)
	case StateMainMenu:
		return m.handleMainMenu(msg)
	case StateLanguageSelect:
		return m.handleLanguageSelect(msg)
	case StateSelectAssistants:
		return m.handleSelectAssistants(msg)
	case StateSelectProvider:
		return m.handleSelectProvider(msg)
	case StateModelGate:
		return m.handleModelGate(msg)
	case StateModelSetup:
		return m.handleModelSetup(msg)
	case StateAgentSetup:
		return m.handleAgentSetup(msg)
	case StateEmojiPref:
		return m.handleEmojiPref(msg)
	case StateAgentWriteMode:
		return m.handleAgentWriteMode(msg)
	case StateReview:
		return m.handleReview(msg)
	case StateDone:
		return m.handleDone(msg)
	}
	return m, nil
}

// handleEnvironmentCheck handles keys for the pre-flight check screen.
// Enter proceeds to StateSelectAssistants once every probe has resolved,
// whatever it found: a missing memory provider never blocks.
// On first entry the assistants found on PATH are pre-selected — every
// assistant when none was found — so the default is visible.
// Esc returns to StateMainMenu at any point.
// q and ctrl+c always quit.
func (m Model) handleEnvironmentCheck(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		if m.preflight.done {
			m.state = StateSelectAssistants
			if len(m.wizard.selected) == 0 {
				m.wizard.selected = m.defaultAssistantSelection()
			}
			m.cursor = 0
			return m, nil
		}
	case tea.KeyEsc:
		m.state = StateMainMenu
		m.cursor = 0
		return m, nil
	case tea.KeyRunes:
		if string(msg.Runes) == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) handleMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < 2 {
			m.cursor++
		}
	case tea.KeyEnter:
		switch m.cursor {
		case 0: // Install / Update Skills → choose language first
			m.state = StateLanguageSelect
			m.cursor = m.language.selected
			return m, LanguagePrefCmd()
		case 1: // Dashboard
			m.state = StateDashboard
			m.cursor = 0
			m.dashboard = loadDashboardState()
		default: // Quit
			return m, tea.Quit
		}
	case tea.KeyEsc:
		// No-op at MainMenu.
	}
	return m, nil
}

// handleLanguageSelect handles keys for the language selection screen.
// Up/Down move the radio selection and re-resolve the catalog live so the
// whole TUI switches language immediately. Enter triggers the preflight
// bootstrap; Esc returns to MainMenu.
func (m Model) handleLanguageSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m = m.syncLanguageSelection()
		}
	case tea.KeyDown:
		if m.cursor < len(languageOptions)-1 {
			m.cursor++
			m = m.syncLanguageSelection()
		}
	case tea.KeyEnter:
		m.preflight.sections = initialPreflightSections(m.catalog.Installer)
		m.preflight.done = false
		m.preflight.providersFound = nil
		m.preflight.codeIntelFound = nil
		m.state = StateEnvironmentCheck
		m.cursor = 0
		return m, tea.Batch(EnvironmentCheckCmd(m.catalog.Installer), m.spinner.Tick)
	case tea.KeyEsc:
		m.state = StateMainMenu
		m.cursor = 0
	}
	return m, nil
}

// syncLanguageSelection returns a Model with the language selection following
// the cursor, the touched guard set, and the catalog re-resolved for the
// newly selected language.
func (m Model) syncLanguageSelection() Model {
	m.language.selected = m.cursor
	m.language.touched = true
	m.catalog = i18n.ForCode(languageOptions[m.language.selected].Code)
	return m
}

// handleDashboard handles the dashboard keys: r asks to remove ASDT's block
// from every global config file that carries one, Enter confirms, and Esc
// cancels a pending confirmation or returns to the main menu.
func (m Model) handleDashboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.dashboard.confirmRemove {
		switch msg.Type {
		case tea.KeyEnter:
			m.dashboard.confirmRemove = false
			return m, AgentConfigRemoveCmd(m.dashboard.blockAssistants())
		case tea.KeyEsc:
			m.dashboard.confirmRemove = false
		}
		return m, nil
	}
	switch msg.Type {
	case tea.KeyRunes:
		if string(msg.Runes) == "r" && len(m.dashboard.blockAssistants()) > 0 {
			m.dashboard.confirmRemove = true
			m.dashboard.removeResults = nil
		}
	case tea.KeyEsc:
		m.state = StateMainMenu
		m.cursor = 0
	}
	return m, nil
}

// blockAssistants returns the assistants whose global config file carries
// ASDT's block, in installer.Descriptors order.
func (d dashboardState) blockAssistants() []installer.AssistantDescriptor {
	var out []installer.AssistantDescriptor
	for _, a := range installer.Descriptors {
		if d.hasBlock[a.ID] {
			out = append(out, a)
		}
	}
	return out
}

func (m Model) handleSelectAssistants(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	count := len(installer.Descriptors)
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < count { // confirm button sits at index count
			m.cursor++
		}
	case tea.KeySpace:
		if m.cursor < count { // guard: confirm button has no selectable index
			m.wizard.selected[m.cursor] = !m.wizard.selected[m.cursor]
		}
	case tea.KeyRunes:
		if string(msg.Runes) == "a" {
			allSelected := true
			for i := range installer.Descriptors {
				if !m.wizard.selected[i] {
					allSelected = false
					break
				}
			}
			for i := range installer.Descriptors {
				m.wizard.selected[i] = !allSelected
			}
		}
	case tea.KeyEnter:
		m.state = StateSelectProvider
		m.wizard.providerChosen = true
		m.cursor = m.wizard.provider
	case tea.KeyEsc:
		m.state = StateEnvironmentCheck
		m.cursor = 0
	}
	return m, nil
}

// handleSelectProvider handles keys for the provider screen: the memory group
// (one radio row per installer.Providers entry, then "No memory"), the
// code-intelligence group (one row per installer.CodeIntelProviders entry,
// then "No code intelligence"), then the confirm button. Up/Down move only the
// cursor — the selection changes on Space or Enter on a row (SelectAssistants
// / ModelSetup idiom), so walking past a "No …" row to the button never
// silently replaces a chosen provider. Space or Enter on a row selects it in
// its own group; Enter then continues, and Enter on the button continues with
// both current selections. An undetected provider is a warning, never a block.
func (m Model) handleSelectProvider(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	count := providerButtonRow() // every row of both groups
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < count { // confirm button sits at index count
			m.cursor++
		}
	case tea.KeySpace:
		m = m.selectProviderRow()
	case tea.KeyEnter:
		return m.selectProviderRow().enterModelGate(), nil
	case tea.KeyEsc:
		m.state = StateSelectAssistants
		m.cursor = 0
	}
	return m, nil
}

// selectProviderRow selects the provider-screen row under the cursor in its
// own group; the button selects nothing.
func (m Model) selectProviderRow() Model {
	switch {
	case m.cursor <= noMemoryChoice():
		m.wizard.provider = m.cursor
	case m.cursor < providerButtonRow():
		m.wizard.codeIntel = m.cursor - codeIntelFirstRow()
	}
	return m
}

// enterModelGate transitions into StateModelGate, detecting available AI
// providers for the subtitle. The gate choice persists across re-entry.
func (m Model) enterModelGate() Model {
	m.wizard.detectedAI = installer.DetectAIProviders()
	m.state = StateModelGate
	m.cursor = m.wizard.modelGateChoice
	return m
}

// modelGateCustomizeChoice is the gate index of the "customize per step" row —
// the only choice that opens StateModelSetup. The lower indices (0-4) are
// presets applied in place.
const modelGateCustomizeChoice = 5

// handleModelGate handles keys for the preset gate. Six radio rows track the
// cursor live (EmojiPref pattern). Enter on a preset (0-4) applies it and goes
// straight to AgentSetup — the customization screen is never rendered; Enter on
// "customize per step" (5) opens StateModelSetup.
func (m Model) handleModelGate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m.wizard.modelGateChoice = m.cursor
		}
	case tea.KeyDown:
		if m.cursor < modelGateCustomizeChoice {
			m.cursor++
			m.wizard.modelGateChoice = m.cursor
		}
	case tea.KeyEnter:
		if m.wizard.modelGateChoice == modelGateCustomizeChoice {
			return m.enterModelSetup(), nil
		}
		m = m.applyPreset(m.wizard.modelGateChoice)
		m.state = StateAgentSetup
		m.cursor = 0
		m.agentConfig.conflicts = m.detectAgentConflicts()
		return m, nil
	case tea.KeyEsc:
		m.state = StateSelectProvider
		m.cursor = m.wizard.provider
	}
	return m, nil
}

// applyPreset classifies every subagent step by its source-default model and
// writes selectedModels from the chosen preset's tier mapping. Choice 0
// (Chameleon) leaves selectedModels empty — install strips every model field so
// each step inherits the assistant's own default. Steps with no source default
// are skipped, falling through to the global model. The step list self-loads
// when the customization accordion was never entered (same err-swallow idiom as
// enterModelSetup).
func (m Model) applyPreset(choice int) Model {
	if m.wizard.modelSteps == nil {
		if steps, err := installer.WorkflowModelSteps(m.skillsFS); err == nil {
			m.wizard.modelSteps = steps
		}
	}
	m.wizard.selectedModels = make(map[string]string, len(m.wizard.modelSteps))
	if choice == installer.PresetChameleon {
		return m
	}
	tiers := installer.PresetModels(choice)
	if tiers == nil {
		return m
	}
	for _, s := range m.wizard.modelSteps {
		if s.Model == "" {
			continue
		}
		m.wizard.selectedModels[s.Key()] = tiers[installer.Classify(s.Model)]
	}
	return m
}

// enterModelSetup transitions into StateModelSetup, loading the step list
// from the embedded skill FS and the model cycle options from detected
// providers. Selections and expanded groups persist across re-entry within
// the session; on first entry selections are prefilled with each step's
// source-default model and all accordion groups start collapsed.
func (m Model) enterModelSetup() Model {
	if m.wizard.modelSteps == nil {
		steps, err := installer.WorkflowModelSteps(m.skillsFS)
		if err == nil {
			m.wizard.modelSteps = steps
		}
		m.wizard.selectedModels = make(map[string]string, len(m.wizard.modelSteps))
		for _, s := range m.wizard.modelSteps {
			if s.Model != "" {
				m.wizard.selectedModels[s.Key()] = s.Model
			}
		}
		m.wizard.expandedGroups = make(map[string]bool)
	}
	m.wizard.modelOptions = installer.FlattenModels(m.wizard.detectedAI)
	m.state = StateModelSetup
	m.cursor = 0
	return m
}

// modelRow is one visible row of the model accordion: a specialist group
// header, or a step belonging to an expanded group.
type modelRow struct {
	header     bool
	specialist string
	stepIdx    int // index into wizard.modelSteps; -1 for headers
}

// modelRows flattens the accordion into its currently visible rows —
// one header per specialist plus the steps of expanded groups, in
// workflow order. The confirm button sits at index len(rows).
func (m Model) modelRows() []modelRow {
	var rows []modelRow
	last := ""
	for i, s := range m.wizard.modelSteps {
		if s.Specialist != last {
			rows = append(rows, modelRow{header: true, specialist: s.Specialist, stepIdx: -1})
			last = s.Specialist
		}
		if m.wizard.expandedGroups[s.Specialist] {
			rows = append(rows, modelRow{specialist: s.Specialist, stepIdx: i})
		}
	}
	return rows
}

// handleModelSetup handles keys for the per-step model accordion.
// Up/Down move through visible rows (confirm button sits past the last row);
// Space toggles a group open/closed on headers and cycles forward on steps;
// Left/Right cycle models on steps and collapse/expand on headers; r resets
// the focused step to its shipped default. Enter advances to AgentSetup from
// any cursor position — same idiom as the sibling wizard screens; Esc
// returns to the gate.
func (m Model) handleModelSetup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	rows := m.modelRows()
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(rows) { // confirm button sits at index len(rows)
			m.cursor++
		}
	case tea.KeyLeft:
		if row, ok := focusedRow(rows, m.cursor); ok {
			if row.header {
				m.wizard.expandedGroups[row.specialist] = false
			} else {
				m = m.cycleModel(row.stepIdx, -1)
			}
		}
	case tea.KeyRight:
		if row, ok := focusedRow(rows, m.cursor); ok {
			if row.header {
				m.wizard.expandedGroups[row.specialist] = true
			} else {
				m = m.cycleModel(row.stepIdx, 1)
			}
		}
	case tea.KeySpace:
		if row, ok := focusedRow(rows, m.cursor); ok {
			if row.header {
				m.wizard.expandedGroups[row.specialist] = !m.wizard.expandedGroups[row.specialist]
			} else {
				m = m.cycleModel(row.stepIdx, 1)
			}
		}
	case tea.KeyRunes:
		if string(msg.Runes) == "r" {
			if row, ok := focusedRow(rows, m.cursor); ok && !row.header {
				step := m.wizard.modelSteps[row.stepIdx]
				if step.Model != "" {
					m.wizard.selectedModels[step.Key()] = step.Model
				} else {
					delete(m.wizard.selectedModels, step.Key())
				}
			}
		}
	case tea.KeyEnter:
		m.state = StateAgentSetup
		m.cursor = 0
		m.agentConfig.conflicts = m.detectAgentConflicts()
		return m, nil
	case tea.KeyEsc:
		m.state = StateModelGate
		m.cursor = m.wizard.modelGateChoice
	}
	return m, nil
}

// focusedRow returns the visible row under the cursor, or ok=false when the
// cursor sits on the confirm button (or the list is empty).
func focusedRow(rows []modelRow, cursor int) (modelRow, bool) {
	if cursor < 0 || cursor >= len(rows) {
		return modelRow{}, false
	}
	return rows[cursor], true
}

// cycleModel advances the given step's model selection by dir through
// modelOptions, wrapping around. A selection not present in the options
// (e.g. a source default from an undetected provider) restarts at index 0.
func (m Model) cycleModel(stepIdx, dir int) Model {
	if stepIdx < 0 || stepIdx >= len(m.wizard.modelSteps) || len(m.wizard.modelOptions) == 0 {
		return m
	}
	key := m.wizard.modelSteps[stepIdx].Key()
	current := m.wizard.selectedModels[key]

	idx := -1
	for i, opt := range m.wizard.modelOptions {
		if opt.ID == current {
			idx = i
			break
		}
	}
	if idx < 0 {
		idx = 0
	} else {
		n := len(m.wizard.modelOptions)
		idx = (idx + dir + n) % n
	}
	m.wizard.selectedModels[key] = m.wizard.modelOptions[idx].ID
	return m
}

// countCustomizedModels returns how many steps currently differ from their
// shipped source default. Zero means the install is byte-identical to the
// skill sources (the skip path), regardless of how the user navigated.
func (m Model) countCustomizedModels() int {
	n := 0
	for _, s := range m.wizard.modelSteps {
		if sel, ok := m.wizard.selectedModels[s.Key()]; ok && sel != s.Model {
			n++
		}
	}
	return n
}

// detectAgentConflicts returns the selected assistants whose global config file
// already exists with content, so the user chooses per file whether ASDT
// updates its block there or leaves the file alone.
func (m Model) detectAgentConflicts() []string {
	var conflicts []string
	assistants := m.selectedAssistants()
	for _, a := range assistants {
		adapter, ok := installer.AgentConfigAdapterFor(a.ID)
		if !ok {
			continue
		}
		if adapter.ConfigExists() {
			conflicts = append(conflicts, a.ID)
		}
	}
	return conflicts
}

// selectedAssistants returns the list of assistants selected by the user (or all
// if none were explicitly selected, mirroring buildInstallCmd).
func (m Model) selectedAssistants() []installer.AssistantDescriptor {
	var assistants []installer.AssistantDescriptor
	for i, d := range installer.Descriptors {
		if m.wizard.selected[i] {
			assistants = append(assistants, d)
		}
	}
	if len(assistants) == 0 {
		assistants = installer.Descriptors
	}
	return assistants
}

func (m Model) handleAgentSetup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 4 navigable presets (0-3); [ Skip → ] button sits at index len(PersonaPresets).
	presets := len(installer.PersonaPresets)
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m = m.syncAgentSetupSelection()
		}
	case tea.KeyDown:
		if m.cursor < presets { // [ Skip → ] button at index presets
			m.cursor++
			m = m.syncAgentSetupSelection()
		}
	case tea.KeyEnter:
		if m.cursor == presets {
			// User pressed Enter on the [ Skip → ] button — bypass EmojiPref and
			// go straight to Review.
			m.agentConfig.skip = true
			m.state = StateReview
			m.cursor = 0
			return m, nil
		}
		// m.agentConfig.selectedPersona reflects the current choice; ask the
		// emoji preference next, defaulting to Yes (cursor 0).
		m.agentConfig.useEmojis = true
		m.state = StateEmojiPref
		m.cursor = 0
		return m, nil
	case tea.KeyEsc:
		// The gate is the stable anchor for backward navigation — never the
		// customization screen, which may not have been visited at all.
		m.state = StateModelGate
		m.cursor = m.wizard.modelGateChoice
	}
	return m, nil
}

// handleEmojiPref handles keys for the emoji preference screen. Two radio rows:
// cursor 0 = Yes, cursor 1 = No; useEmojis stays in sync with the cursor.
// Enter continues to AgentWriteMode when conflicts exist, otherwise to Review.
func (m Model) handleEmojiPref(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m.agentConfig.useEmojis = m.cursor == 0
		}
	case tea.KeyDown:
		if m.cursor < 1 {
			m.cursor++
			m.agentConfig.useEmojis = m.cursor == 0
		}
	case tea.KeyEnter:
		if len(m.agentConfig.conflicts) > 0 {
			m.agentConfig.writeModes = make(map[string]installer.AgentWriteMode)
			for _, id := range m.agentConfig.conflicts {
				m.agentConfig.writeModes[id] = installer.AgentModeUpdate // touches only ASDT's block
			}
			m.state = StateAgentWriteMode
			m.cursor = 0
			return m, nil
		}
		m.state = StateReview
		m.cursor = 0
		return m, nil
	case tea.KeyEsc:
		m.state = StateAgentSetup
		m.cursor = 0
	}
	return m, nil
}

// reviewPriorState returns the step that precedes Review for the current path,
// so Review's Esc retraces the path the user took.
func reviewPriorState(m Model) ViewState {
	switch {
	case len(m.agentConfig.conflicts) > 0 && !m.agentConfig.skip:
		return StateAgentWriteMode
	case !m.agentConfig.skip:
		return StateEmojiPref
	default:
		return StateAgentSetup
	}
}

// syncAgentSetupSelection returns a Model with selectedPersona / skip
// updated to match the current cursor. Only called for cursor positions within
// the preset list (0 to len(PersonaPresets)-1); the [ Skip → ] button position
// is handled explicitly in the Enter handler.
func (m Model) syncAgentSetupSelection() Model {
	if m.cursor < len(installer.PersonaPresets) {
		m.agentConfig.selectedPersona = m.cursor
		m.agentConfig.skip = false
	}
	return m
}

func (m Model) handleAgentWriteMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	count := len(m.agentConfig.conflicts)
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < count { // confirm button sits at index count
			m.cursor++
		}
	case tea.KeySpace:
		if m.cursor < len(m.agentConfig.conflicts) {
			id := m.agentConfig.conflicts[m.cursor]
			if m.agentConfig.writeModes[id] == installer.AgentModeSkip {
				m.agentConfig.writeModes[id] = installer.AgentModeUpdate
			} else {
				m.agentConfig.writeModes[id] = installer.AgentModeSkip
			}
		}
	case tea.KeyEnter:
		m.state = StateReview
		m.cursor = 0
	case tea.KeyEsc:
		m.state = StateEmojiPref
		m.cursor = 0
	}
	return m, nil
}

func (m Model) handleReview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyDown:
		if m.cursor == 0 {
			m.cursor = 1 // scroll to Install button
		}
	case tea.KeyUp:
		if m.cursor == 1 {
			m.cursor = 0
		}
	case tea.KeyEnter:
		m.state = StateInstalling
		m.cursor = 0
		m.wizard.installExpected = len(m.selectedAssistants())
		if m.agentConfig.skip {
			m.wizard.agentDone = true
		}
		return m, tea.Batch(m.buildInstallCmd(), m.spinner.Tick)
	case tea.KeyEsc:
		m.state = reviewPriorState(m)
		m.cursor = 0
	}
	return m, nil
}

func (m Model) handleDone(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyEnter:
		m.state = StateMainMenu
		m.cursor = 0
		m.wizard.results = nil
		m.wizard.agentResults = nil
		m.wizard.installExpected = 0
		m.wizard.agentDone = false
		m.wizard.selected = make(map[int]bool)
	}
	return m, nil
}

func (m Model) buildInstallCmd() tea.Cmd {
	assistants := m.selectedAssistants()
	provider := m.SelectedProvider()
	lang := languageOptions[m.language.selected].Code
	opts := installer.InstallOptions{CodeIntel: m.SelectedCodeIntel()}

	var installCmd tea.Cmd
	if m.wizard.modelGateChoice == installer.PresetChameleon {
		// Chameleon must be detected BEFORE the nil-guard below: its
		// selectedModels map is empty, which would otherwise fall through to a
		// verbatim install and wrongly keep the source defaults. Instead the
		// remove-models path strips every model field.
		installCmd = RemoveModelsInstallCmd(assistants, provider, m.skillsFS, lang, opts)
	} else {
		// Selections identical to the shipped defaults pass nil so workflow.yaml
		// files install verbatim — a preset that reproduces the defaults
		// (Craftsman), or customizing and reverting, never rewrites what the
		// skill ships.
		models := m.wizard.selectedModels
		if m.countCustomizedModels() == 0 {
			models = nil
		}
		installCmd = InstallCmd(assistants, provider, m.skillsFS, lang, models, opts)
	}

	cmds := []tea.Cmd{installCmd}
	if !m.agentConfig.skip {
		preset := installer.PersonaPresets[m.agentConfig.selectedPersona]
		cmds = append(cmds, AgentInstallCmd(assistants, preset, m.agentConfig.useEmojis, lang, m.agentConfig.writeModes))
	}
	if len(cmds) == 1 {
		return cmds[0]
	}
	return tea.Batch(cmds...)
}

// RemoveModelsInstallCmd is InstallCmd for the Chameleon preset: each
// per-assistant install strips the `model:` field from every subagent step as
// its workflow.yaml is written. It mirrors InstallCmd's per-assistant batching
// so the TUI still updates row-by-row.
func RemoveModelsInstallCmd(assistants []installer.AssistantDescriptor, provider installer.ProviderDescriptor, skillsFS fs.FS, lang string, opts installer.InstallOptions) tea.Cmd {
	cmds := make([]tea.Cmd, len(assistants))
	for i, a := range assistants {
		a := a // capture loop variable
		cmds[i] = func() tea.Msg {
			results := installer.InstallRemovingModels([]installer.AssistantDescriptor{a}, provider, skillsFS, lang, opts)
			if len(results) > 0 {
				return AssistantInstallProgressMsg{Result: results[0]}
			}
			return AssistantInstallProgressMsg{
				Result: installer.InstallResult{AssistantID: a.ID},
			}
		}
	}
	return tea.Batch(cmds...)
}

// AgentWriteModes returns the per-assistant write mode map. Exported for tests.
func (m Model) AgentWriteModes() map[string]installer.AgentWriteMode {
	return m.agentConfig.writeModes
}

// SelectedModels returns the per-step model selection map. Exported for tests.
func (m Model) SelectedModels() map[string]string {
	return m.wizard.selectedModels
}

// ModelsTouched reports whether any step currently differs from its shipped
// default. Exported for tests.
func (m Model) ModelsTouched() bool {
	return m.countCustomizedModels() > 0
}

// UseEmojis returns the chosen emoji preference. Exported for tests.
func (m Model) UseEmojis() bool {
	return m.agentConfig.useEmojis
}

// LanguageCode returns the currently selected language code. Exported for tests.
func (m Model) LanguageCode() string {
	return languageOptions[m.language.selected].Code
}

// View renders the current state.
func (m Model) View() string {
	return renderState(m)
}

// loadDashboardState reads install metadata and agent config block presence
// for all known assistants. Called when the user navigates to StateDashboard,
// and after a removal, so the render path never performs I/O.
func loadDashboardState() dashboardState {
	meta := make(map[installer.AssistantID]installer.InstallMeta, len(installer.Descriptors))
	hasBlock := make(map[installer.AssistantID]bool, len(installer.Descriptors))
	for _, d := range installer.Descriptors {
		if m, err := installer.ReadInstallMeta(d); err == nil {
			meta[d.ID] = m
		}
		if adapter, ok := installer.AgentConfigAdapterFor(d.ID); ok {
			hasBlock[d.ID] = adapter.HasBlock()
		}
	}
	return dashboardState{meta: meta, hasBlock: hasBlock}
}

// allInstallsDone reports whether both the per-assistant installs and the agent
// config step have finished, so Update can transition to StateDone.
func (m Model) allInstallsDone() bool {
	if len(m.wizard.results) < m.wizard.installExpected {
		return false
	}
	return m.wizard.agentDone
}

func allProbesDone(sections []components.SectionGroup) bool {
	for _, sg := range sections {
		for _, row := range sg.Rows {
			if row.Status == components.CheckStatusPending {
				return false
			}
		}
	}
	return true
}

// defaultAssistantSelection is the assistant-screen default: every assistant
// the environment check found on PATH, or every assistant when it found none
// (or reported nothing), so the screen never opens with an empty selection.
func (m Model) defaultAssistantSelection() map[int]bool {
	selected := make(map[int]bool, len(installer.Descriptors))
	for i, d := range installer.Descriptors {
		if m.preflight.assistantsFound[d.ID] {
			selected[i] = true
		}
	}
	if len(selected) == 0 {
		for i := range installer.Descriptors {
			selected[i] = true
		}
	}
	return selected
}

// noMemoryChoice is the provider-screen index of the "No memory" row: one past
// the last installer.Providers entry.
func noMemoryChoice() int { return len(installer.Providers) }

// defaultProviderChoice is the provider-screen default: the first provider the
// environment check found (a provider with no local probe counts as found),
// else "No memory".
func (m Model) defaultProviderChoice() int {
	for i, p := range installer.Providers {
		if m.preflight.providersFound[p.ID] {
			return i
		}
	}
	return noMemoryChoice()
}

// SelectedProvider returns the chosen memory provider, or the zero
// ProviderDescriptor — the installer's "no memory" selection — when "No
// memory" is chosen. Exported for tests.
func (m Model) SelectedProvider() installer.ProviderDescriptor {
	if m.wizard.provider < noMemoryChoice() {
		return installer.Providers[m.wizard.provider]
	}
	return installer.ProviderDescriptor{}
}

// codeIntelFirstRow is the provider-screen index of the first
// code-intelligence row: just below the memory group's "No memory".
func codeIntelFirstRow() int { return noMemoryChoice() + 1 }

// noCodeIntelChoice is the code-intelligence group index of "No code
// intelligence": one past the last installer.CodeIntelProviders entry.
func noCodeIntelChoice() int { return len(installer.CodeIntelProviders) }

// providerButtonRow is the provider-screen index of the confirm button: one
// past the code-intelligence group's "No code intelligence".
func providerButtonRow() int { return codeIntelFirstRow() + noCodeIntelChoice() + 1 }

// defaultCodeIntelChoice is the code-intelligence default: the first provider
// the environment check found, else "No code intelligence".
func (m Model) defaultCodeIntelChoice() int {
	for i, p := range installer.CodeIntelProviders {
		if m.preflight.codeIntelFound[p.ID] {
			return i
		}
	}
	return noCodeIntelChoice()
}

// SelectedCodeIntel returns the chosen code-intelligence provider, or the zero
// CodeIntelDescriptor — the installer's "no code intelligence" selection.
// Exported for tests.
func (m Model) SelectedCodeIntel() installer.CodeIntelDescriptor {
	if m.wizard.codeIntel < noCodeIntelChoice() {
		return installer.CodeIntelProviders[m.wizard.codeIntel]
	}
	return installer.CodeIntelDescriptor{}
}

func rowHasStatus(sections []components.SectionGroup, label string, status components.CheckStatus) bool {
	for _, sg := range sections {
		for _, row := range sg.Rows {
			if row.Label == label {
				return row.Status == status
			}
		}
	}
	return false
}
