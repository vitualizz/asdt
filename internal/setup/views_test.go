package setup_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/vitualizz/asdt/internal/installer"
	"github.com/vitualizz/asdt/internal/setup"
	"github.com/vitualizz/asdt/internal/tui/panels"
)

// TestMain forces the English locale for the entire test binary so that
// view string assertions are deterministic regardless of the developer's
// system locale.
func TestMain(m *testing.M) {
	_ = os.Setenv("ASDT_LANG", "en")
	os.Exit(m.Run())
}

func TestView_MainMenuContainsInstallOption(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = sendEngramFound(t, m)
	view := m.View()
	if !strings.Contains(view, "Install / Update Skills") {
		t.Errorf("main menu view missing 'Install / Update Skills', got:\n%s", view)
	}
}

func TestView_MainMenuShowsAsdtTuiHeader(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = sendEngramFound(t, m)
	view := m.View()
	if !strings.Contains(view, "asdt-tui") {
		t.Errorf("main menu view missing 'asdt-tui' header, got:\n%s", view)
	}
}

func TestView_MainMenuShowsUpdateBanner(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "v0.2.0")
	m = sendEngramFound(t, m)

	next, _ := m.Update(setup.UpdateCheckMsg{Current: "v0.2.0", Latest: "v0.3.0"})
	m2 := next.(setup.Model)

	view := m2.View()
	if !strings.Contains(view, "/releases") {
		t.Errorf("main menu view missing update banner releases URL, got:\n%s", view)
	}
	if !strings.Contains(view, "v0.3.0") {
		t.Errorf("main menu view missing latest version v0.3.0, got:\n%s", view)
	}
}

func TestView_SelectAssistantsShowsBothNames(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = sendEngramFound(t, m)         // no-op, already at MainMenu
	m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	next, _ := m.Update(envCheckMsg(true))
	m = next.(setup.Model)
	m2 := updateKey(t, m, tea.KeyEnter) // preflightDone → StateSelectAssistants
	view := m2.View()
	for _, d := range installer.Descriptors {
		if !strings.Contains(view, d.Name) {
			t.Errorf("select assistants missing %q, view:\n%s", d.Name, view)
		}
	}
}

func TestView_SelectAssistantsSelectedItemHasCursor(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = sendEngramFound(t, m)         // no-op, already at MainMenu
	m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	next, _ := m.Update(envCheckMsg(true))
	m = next.(setup.Model)
	m2 := updateKey(t, m, tea.KeyEnter) // preflightDone → StateSelectAssistants
	view := m2.View()
	if !strings.Contains(view, "►") {
		t.Errorf("select assistants missing cursor ►, view:\n%s", view)
	}
}

func TestView_SelectAssistantsUsesBadgeForStatus(t *testing.T) {
	// Force color so panels.Badge styling is distinguishable from plain-text fallback.
	prev := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	defer lipgloss.SetColorProfile(prev)

	m := setup.New(fstest.MapFS{}, "dev")
	m = sendEngramFound(t, m)         // no-op, already at MainMenu
	m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	next, _ := m.Update(envCheckMsg(true))
	m = next.(setup.Model)
	m2 := updateKey(t, m, tea.KeyEnter) // preflightDone → StateSelectAssistants
	view := m2.View()

	wantPresent := panels.NewBadge("present", panels.ColorSuccess).Render()
	wantMissing := panels.NewBadge("missing", panels.ColorError).Render()

	if !strings.Contains(view, wantPresent) && !strings.Contains(view, wantMissing) {
		t.Errorf("expected select assistants to render status via panels.Badge ([present]/[missing] in tone color), got: %q", view)
	}
}

func TestView_SelectAssistantsShowsCheckboxes(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = sendEngramFound(t, m)
	m = updateKey(t, m, tea.KeyEnter) // → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	next, _ := m.Update(envCheckMsg(true))
	m = next.(setup.Model)
	m2 := updateKey(t, m, tea.KeyEnter) // → StateSelectAssistants
	view := m2.View()
	if !strings.Contains(view, "[x]") {
		t.Errorf("select assistants missing pre-selected checkbox [x], view:\n%s", view)
	}
}

func TestView_DoneScreenShowsBothAssistants(t *testing.T) {
	successResult := installer.InstallResult{AssistantID: installer.AssistantClaudeCode, Err: nil}
	failResult := installer.InstallResult{AssistantID: installer.AssistantOpenCode, Err: fmt.Errorf("fail")}
	m := setup.New(fstest.MapFS{}, "dev")
	next, _ := m.Update(setup.InstallDoneMsg{Results: []installer.InstallResult{successResult, failResult}})
	m2 := next.(setup.Model)
	view := m2.View()
	if !strings.Contains(view, "Claude Code") {
		t.Errorf("done screen missing 'Claude Code', view:\n%s", view)
	}
	if !strings.Contains(view, "OpenCode") {
		t.Errorf("done screen missing 'OpenCode', view:\n%s", view)
	}
}

// stateView drives the model from StateMainMenu to the requested ViewState by
// walking the same Update transitions the real TUI uses, and returns the
// rendered View() output for that state.
func stateView(t *testing.T, target string) string {
	t.Helper()

	m := setup.New(fstest.MapFS{}, "dev")

	if target == "PreflightCheck" {
		// Trigger preflight by pressing Enter at cursor-0 (Install), then
		// Enter again through the language screen.
		m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
		m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
		return m.View()
	}

	m = sendEngramFound(t, m) // no-op, already at MainMenu
	if target == "MainMenu" {
		return m.View()
	}

	// For any state that passes through SelectProvider → AgentSetup, use a
	// clean HOME so detectAgentConflicts returns empty (no conflicts).
	// AgentWriteMode intentionally sets up its own HOME with a conflict.
	if target != "AgentWriteMode" {
		t.Setenv("HOME", t.TempDir())
		t.Setenv("XDG_CONFIG_HOME", "")
	}

	// Install path: cursor-0 Enter → StateEnvironmentCheck → EnvironmentCheckMsg → Enter → StateSelectAssistants.
	m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	next, _ := m.Update(envCheckMsg(true))
	m = next.(setup.Model)
	m = updateKey(t, m, tea.KeyEnter) // preflightDone → StateSelectAssistants
	if target == "SelectAssistants" {
		return m.View()
	}

	m = updateKey(t, m, tea.KeyEnter) // SelectAssistants → SelectProvider
	if target == "SelectProvider" {
		return m.View()
	}

	m = updateKey(t, m, tea.KeyEnter) // SelectProvider → ModelGate
	if target == "ModelGate" {
		return m.View()
	}

	if target == "ModelSetup" {
		for i := 0; i < 5; i++ {
			m = updateKey(t, m, tea.KeyDown) // move to "customize per step" (index 5)
		}
		m = updateKey(t, m, tea.KeyEnter) // ModelGate → ModelSetup
		return m.View()
	}

	m = updateKey(t, m, tea.KeyEnter) // ModelGate (Chameleon default) → AgentSetup
	if target == "AgentSetup" {
		return m.View()
	}

	m = updateKey(t, m, tea.KeyEnter) // AgentSetup → EmojiPref
	if target == "EmojiPref" {
		return m.View()
	}

	if target == "AgentWriteMode" {
		// Need a conflict to enter StateAgentWriteMode. Set up a CLAUDE.md with
		// the asdt block marker so detectAgentConflicts returns a hit.
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		t.Setenv("XDG_CONFIG_HOME", "")
		claudeDir := tmpHome + "/.claude"
		if err := os.MkdirAll(claudeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		marker := "<!-- asdt:agent-config -->"
		if err := os.WriteFile(claudeDir+"/CLAUDE.md", []byte(marker+"\n<!-- /asdt:agent-config -->\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		// Re-drive from scratch with the new HOME so detectAgentConflicts fires.
		m2 := setup.New(fstest.MapFS{}, "dev")
		m2 = updateKey(t, m2, tea.KeyEnter) // MainMenu → LanguageSelect
		m2 = updateKey(t, m2, tea.KeyEnter) // LanguageSelect → EnvironmentCheck
		next2, _ := m2.Update(envCheckMsg(true))
		m2 = next2.(setup.Model)
		m2 = updateKey(t, m2, tea.KeyEnter) // EnvironmentCheck → SelectAssistants
		m2 = updateKey(t, m2, tea.KeyEnter) // SelectAssistants → SelectProvider
		m2 = updateKey(t, m2, tea.KeyEnter) // SelectProvider → ModelGate
		m2 = updateKey(t, m2, tea.KeyEnter) // ModelGate (recommended) → AgentSetup (conflict detected)
		m2 = updateKey(t, m2, tea.KeyEnter) // AgentSetup preset → EmojiPref
		m2 = updateKey(t, m2, tea.KeyEnter) // EmojiPref → AgentWriteMode (conflict)
		return m2.View()
	}

	m = updateKey(t, m, tea.KeyEnter) // EmojiPref → Review
	if target == "Review" {
		return m.View()
	}

	m = updateKey(t, m, tea.KeyEnter) // Review → Installing
	if target == "Installing" {
		return m.View()
	}

	t.Fatalf("stateView: unknown target %q", target)
	return ""
}

// TestView_AllStatesHaveBorder proves every one of the view states wraps
// its body in the rounded-border Box style (spec: "every screen's rendered
// output contains a bordered box").
func TestView_AllStatesHaveBorder(t *testing.T) {
	states := []string{
		"PreflightCheck",
		"MainMenu",
		"SelectAssistants",
		"SelectProvider",
		"AgentSetup",
		"EmojiPref",
		"Review",
		"AgentWriteMode",
		"Installing",
	}

	for _, state := range states {
		t.Run(state, func(t *testing.T) {
			view := stateView(t, state)
			if !strings.ContainsAny(view, "╭╮╰╯│") {
				t.Errorf("%s view missing rounded-border runes, got:\n%s", state, view)
			}
		})
	}

	t.Run("Dashboard", func(t *testing.T) {
		m := setup.New(fstest.MapFS{}, "dev")
		m = sendEngramFound(t, m)
		m = updateKey(t, m, tea.KeyDown)  // cursor → 1 (Dashboard)
		m = updateKey(t, m, tea.KeyEnter) // → StateDashboard
		view := m.View()
		if !strings.ContainsAny(view, "╭╮╰╯│") {
			t.Errorf("Dashboard view missing rounded-border runes, got:\n%s", view)
		}
	})

	t.Run("Done", func(t *testing.T) {
		successResult := installer.InstallResult{AssistantID: installer.AssistantClaudeCode, Err: nil}
		m := setup.New(fstest.MapFS{}, "dev")
		next, _ := m.Update(setup.InstallDoneMsg{Results: []installer.InstallResult{successResult}})
		m2 := next.(setup.Model)
		view := m2.View()
		if !strings.ContainsAny(view, "╭╮╰╯│") {
			t.Errorf("Done view missing rounded-border runes, got:\n%s", view)
		}
	})
}

// TestView_FooterRendersHintText proves the per-screen key-hint footer text
// is rendered inside the bordered box on every screen.
func TestView_FooterRendersHintText(t *testing.T) {
	cases := []struct {
		state string
		hint  string
	}{
		{"PreflightCheck", "checking"},
		{"MainMenu", "↑↓"},
		{"SelectAssistants", "space"},
		{"SelectProvider", "esc"},
		{"EmojiPref", "↑↓"},
		{"Review", "esc"},
		{"Installing", "q"},
	}

	for _, tc := range cases {
		t.Run(tc.state, func(t *testing.T) {
			view := stateView(t, tc.state)
			if !strings.Contains(view, tc.hint) {
				t.Errorf("%s view missing footer hint %q, got:\n%s", tc.state, tc.hint, view)
			}
		})
	}

	t.Run("Dashboard", func(t *testing.T) {
		m := setup.New(fstest.MapFS{}, "dev")
		m = sendEngramFound(t, m)
		m = updateKey(t, m, tea.KeyDown)  // cursor → 1 (Dashboard)
		m = updateKey(t, m, tea.KeyEnter) // → StateDashboard
		view := m.View()
		if !strings.Contains(view, "esc") {
			t.Errorf("Dashboard view missing footer hint %q, got:\n%s", "esc", view)
		}
	})

	t.Run("Done", func(t *testing.T) {
		successResult := installer.InstallResult{AssistantID: installer.AssistantClaudeCode, Err: nil}
		m := setup.New(fstest.MapFS{}, "dev")
		next, _ := m.Update(setup.InstallDoneMsg{Results: []installer.InstallResult{successResult}})
		m2 := next.(setup.Model)
		view := m2.View()
		if !strings.Contains(view, "enter/esc") {
			t.Errorf("Done view missing footer hint %q, got:\n%s", "enter/esc", view)
		}
	})
}

func TestView_DashboardShowsAssistantNames(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = sendEngramFound(t, m)
	m = updateKey(t, m, tea.KeyDown)  // cursor → 1 (Dashboard)
	m = updateKey(t, m, tea.KeyEnter) // → StateDashboard
	view := m.View()
	for _, d := range installer.Descriptors {
		if !strings.Contains(view, d.Name) {
			t.Errorf("dashboard view missing assistant name %q, got:\n%s", d.Name, view)
		}
	}
}

// TestView_InstallingShowsSpinner verifies that the StateInstalling screen
// renders an indeterminate spinner.Dot frame glyph alongside the existing
// progress indication — the visual cue that installation is actively running
// (T-030..T-039).
func TestView_InstallingShowsSpinner(t *testing.T) {
	view := stateView(t, "Installing")

	found := false
	for _, frame := range spinner.Dot.Frames {
		glyph := strings.TrimRight(frame, " ")
		if glyph == "" {
			continue
		}
		if strings.Contains(view, glyph) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected Installing view to contain a spinner.Dot frame glyph, got:\n%s", view)
	}
}

func TestView_PreflightCheckShowsTitle(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	view := m.View()
	if !strings.Contains(view, "Pre-flight Check") {
		t.Errorf("preflight view should contain 'Pre-flight Check', got:\n%s", view)
	}
}

func TestView_PreflightCheckShowsSections(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	view := m.View()
	if !strings.Contains(view, "Memory Provider") {
		t.Errorf("preflight view missing 'Memory Provider' section, got:\n%s", view)
	}
}

func TestView_PreflightCheckShowsShellRow(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = updateKey(t, m, tea.KeyEnter) // cursor-0 (Install) → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	view := m.View()
	if !strings.Contains(view, "Shell") {
		t.Errorf("preflight view missing 'Shell' row, got:\n%s", view)
	}
}

func TestView_ReviewShowsAssistantsAndProvider(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "Review")
	for _, d := range installer.Descriptors {
		if !strings.Contains(view, d.Name) {
			t.Errorf("review view missing assistant name %q, got:\n%s", d.Name, view)
		}
	}
	if !strings.Contains(view, "Assistants") {
		t.Errorf("review view missing 'Assistants' label, got:\n%s", view)
	}
}

func TestView_SelectProviderShowsRadioIndicator(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "SelectProvider")
	if !strings.Contains(view, "(•)") {
		t.Errorf("SelectProvider missing selected radio indicator '(•)', got:\n%s", view)
	}
	// Unselected radio only appears when there is more than one provider.
	if len(installer.Providers) > 1 && !strings.Contains(view, "( )") {
		t.Errorf("SelectProvider missing unselected radio indicator '( )', got:\n%s", view)
	}
}

func TestView_EmojiPrefShowsRadioOptionsAndSubtitle(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "EmojiPref")
	if !strings.Contains(view, "Emoji Preference") {
		t.Errorf("EmojiPref missing title 'Emoji Preference', got:\n%s", view)
	}
	if !strings.Contains(view, "Should your assistants use emojis in their responses?") {
		t.Errorf("EmojiPref missing subtitle question, got:\n%s", view)
	}
	if !strings.Contains(view, "Yes — use emojis") {
		t.Errorf("EmojiPref missing Yes option, got:\n%s", view)
	}
	if !strings.Contains(view, "No — plain text only") {
		t.Errorf("EmojiPref missing No option, got:\n%s", view)
	}
	if !strings.Contains(view, "(•)") {
		t.Errorf("EmojiPref missing selected radio indicator '(•)', got:\n%s", view)
	}
	if !strings.Contains(view, "( )") {
		t.Errorf("EmojiPref missing unselected radio indicator '( )', got:\n%s", view)
	}
}

func TestView_EmojiPrefSharesStepFive(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "EmojiPref")
	if !strings.Contains(view, "step 5 of 6") {
		t.Errorf("EmojiPref should share step 5 of 6 with AgentSetup, got:\n%s", view)
	}
}

func TestView_ModelGateShowsPresetRowsWithChameleonDefault(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "ModelGate")
	if !strings.Contains(view, "step 4 of 6") {
		t.Errorf("ModelGate should be step 4 of 6, got:\n%s", view)
	}
	if !strings.Contains(view, "AI Models") {
		t.Errorf("ModelGate view missing title, got:\n%s", view)
	}
	// All six preset rows must be present.
	for _, name := range []string{"Chameleon", "Sprinter", "Craftsman", "Strategist", "Mastermind", "Customize per step"} {
		if !strings.Contains(view, name) {
			t.Errorf("ModelGate view missing preset row %q, got:\n%s", name, view)
		}
	}
	// Chameleon (the default) must be the selected radio.
	chameleonLine := ""
	for _, line := range strings.Split(view, "\n") {
		if strings.Contains(line, "Chameleon") {
			chameleonLine = line
			break
		}
	}
	if !strings.Contains(chameleonLine, "(•)") {
		t.Errorf("Chameleon should be the preselected radio, got line:\n%s", chameleonLine)
	}
}

func TestView_ModelSetupShowsStepsAndModels(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "ModelSetup")
	if !strings.Contains(view, "step 4 of 6") {
		t.Errorf("ModelSetup should be step 4 of 6, got:\n%s", view)
	}
	if !strings.Contains(view, "Choose Models per Step") {
		t.Errorf("ModelSetup view missing title, got:\n%s", view)
	}
}

func TestView_ReviewShowsModelsRow(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "Review")
	if !strings.Contains(view, "Models:") {
		t.Errorf("review view missing 'Models:' row, got:\n%s", view)
	}
	// The default gate choice is Chameleon, so Review reports the preset name.
	if !strings.Contains(view, "Chameleon preset") {
		t.Errorf("review view should show 'Chameleon preset' on the default path, got:\n%s", view)
	}
	if strings.Contains(view, "recommended defaults") {
		t.Errorf("review view must not show a 'recommended defaults' string, got:\n%s", view)
	}
}

func TestView_ReviewShowsEmojiRowWhenNotSkipped(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "Review")
	if !strings.Contains(view, "Emojis:") {
		t.Errorf("review view missing 'Emojis:' label, got:\n%s", view)
	}
	if !strings.Contains(view, "Yes — use emojis") {
		t.Errorf("review view missing emoji value 'Yes — use emojis' (default), got:\n%s", view)
	}
}

func TestView_ReviewOmitsEmojiRowWhenSkipped(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	// Drive the skip path to Review: EmojiPref is bypassed entirely.
	m := setup.New(fstest.MapFS{}, "dev")
	m = updateKey(t, m, tea.KeyEnter) // MainMenu → LanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → EnvironmentCheck
	next, _ := m.Update(envCheckMsg(true))
	m = next.(setup.Model)
	m = updateKey(t, m, tea.KeyEnter) // EnvironmentCheck → SelectAssistants
	m = updateKey(t, m, tea.KeyEnter) // SelectAssistants → SelectProvider
	m = updateKey(t, m, tea.KeyEnter) // SelectProvider → ModelSetup
	m = updateKey(t, m, tea.KeyEnter) // ModelSetup → AgentSetup
	for range installer.PersonaPresets {
		m = updateKey(t, m, tea.KeyDown) // navigate to [ Skip → ]
	}
	m = updateKey(t, m, tea.KeyEnter) // Skip → Review
	if m.State() != setup.StateReview {
		t.Fatalf("expected StateReview, got %v", m.State())
	}
	view := m.View()
	if strings.Contains(view, "Emojis:") {
		t.Errorf("review view must omit 'Emojis:' row when persona was skipped, got:\n%s", view)
	}
}

func TestView_AgentSetupShowsRadioIndicator(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	view := stateView(t, "AgentSetup")
	if !strings.Contains(view, "(•)") {
		t.Errorf("AgentSetup missing selected radio indicator '(•)', got:\n%s", view)
	}
	if !strings.Contains(view, "( )") {
		t.Errorf("AgentSetup missing unselected radio indicator '( )', got:\n%s", view)
	}
}

// TestView_PreflightNeverBlocksOnProviders asserts that with no memory
// provider found the environment check still offers Continue and shows no
// blocking recovery block: memory is optional.
func TestView_PreflightNeverBlocksOnProviders(t *testing.T) {
	m := setup.New(fstest.MapFS{}, "dev")
	m = updateKey(t, m, tea.KeyEnter) // → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	next, _ := m.Update(envCheckMsg(false))
	view := next.(setup.Model).View()
	if !strings.Contains(view, "enter") {
		t.Errorf("preflight with no provider must offer Continue (enter), got:\n%s", view)
	}
	for _, blocked := range []string{"required", "Required", "restart asdt-tui"} {
		if strings.Contains(view, blocked) {
			t.Errorf("preflight with no provider must not read as blocked (%q), got:\n%s", blocked, view)
		}
	}
}

// TestView_ProviderWarningFollowsSelection asserts the provider screen lists
// every provider plus "No memory" with why a provider is recommended, and that
// a selected provider the environment check did not find shows a warning with
// its install hint — while Continue still works — and no other row does.
func TestView_ProviderWarningFollowsSelection(t *testing.T) {
	withHostedProvider(t)
	const missing = "Engram was not found. ASDT installs bound to it and runs without memory until you install and connect it."
	const install = "Install: https://github.com/Gentleman-Programming/engram"

	m := setup.New(fstest.MapFS{}, "dev")
	m = updateKey(t, m, tea.KeyEnter) // → StateLanguageSelect
	m = updateKey(t, m, tea.KeyEnter) // LanguageSelect → StateEnvironmentCheck
	msg := envCheckMsg(true)
	msg.ProvidersFound[installer.ProviderEngram] = false
	next, _ := m.Update(msg)
	m = updateKey(t, next.(setup.Model), tea.KeyEnter) // → SelectAssistants
	m = updateKey(t, m, tea.KeyEnter)                  // → SelectProvider (hosted provider, the first found)

	view := m.View()
	for _, want := range []string{"Engram", hostedProvider.Name, "No memory", "keeps key project knowledge between sessions"} {
		if !strings.Contains(view, want) {
			t.Errorf("provider screen missing %q, got:\n%s", want, view)
		}
	}
	if strings.Contains(view, missing) {
		t.Errorf("provider screen on a found provider must show no warning, got:\n%s", view)
	}

	onEngram := updateKey(t, updateKey(t, m, tea.KeyUp), tea.KeySpace)
	if view := onEngram.View(); !strings.Contains(view, missing) || !strings.Contains(view, install) {
		t.Errorf("provider screen on a missing Engram must show its warning and install hint, got:\n%s", view)
	}
	if got := updateKey(t, onEngram, tea.KeyEnter).State(); got != setup.StateModelGate {
		t.Errorf("Enter on a missing Engram: state = %v, want StateModelGate (a warning, never a block)", got)
	}

	if view := updateKey(t, updateKey(t, m, tea.KeyDown), tea.KeySpace).View(); strings.Contains(view, missing) {
		t.Errorf("provider screen on No memory must show no provider warning, got:\n%s", view)
	}
}

// TestView_ReviewMemoryRow asserts the Review screen's Memory row: the chosen
// provider's name, a not-found flag with its install hint, or "none" with the
// recommended provider.
func TestView_ReviewMemoryRow(t *testing.T) {
	cases := []struct {
		name        string
		engramFound bool
		downs       int // Down presses on the provider screen from the default
		want        string
	}{
		{name: "engram found", engramFound: true, want: "Memory:      Engram"},
		{name: "no memory by default", want: "Memory:      none (recommended: Engram, to keep key project knowledge between sessions)"},
		{name: "engram chosen while missing", downs: -1, want: "Engram (not found — install: https://github.com/Gentleman-Programming/engram)"},
		{name: "no memory chosen beside engram", engramFound: true, downs: 1, want: "none (recommended: Engram"},
		// Walking past "No memory" and the code-intelligence group onto
		// [ Continue → ] keeps the default.
		{name: "engram kept on the button", engramFound: true, downs: providerScreenRows, want: "Memory:      Engram"},
		{name: "no memory kept on the button", downs: providerScreenRows, want: "Memory:      none (recommended: Engram"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())
			t.Setenv("XDG_CONFIG_HOME", "")
			m := setup.New(fstest.MapFS{}, "dev")
			m = updateKey(t, m, tea.KeyEnter) // → LanguageSelect
			m = updateKey(t, m, tea.KeyEnter) // → EnvironmentCheck
			next, _ := m.Update(envCheckMsg(c.engramFound))
			m = updateKey(t, next.(setup.Model), tea.KeyEnter) // → SelectAssistants
			m = updateKey(t, m, tea.KeyEnter)                  // → SelectProvider
			for range c.downs {
				m = updateKey(t, m, tea.KeyDown)
			}
			for range -c.downs {
				m = updateKey(t, m, tea.KeyUp)
			}
			m = updateKey(t, m, tea.KeyEnter) // → ModelGate
			m = updateKey(t, m, tea.KeyEnter) // → AgentSetup
			for range installer.PersonaPresets {
				m = updateKey(t, m, tea.KeyDown) // → [ Skip → ]
			}
			m = updateKey(t, m, tea.KeyEnter) // → Review
			if m.State() != setup.StateReview {
				t.Fatalf("state = %v, want StateReview", m.State())
			}
			if view := m.View(); !strings.Contains(view, c.want) {
				t.Errorf("review missing %q, got:\n%s", c.want, view)
			}
		})
	}
}

// providerScreenRows is the number of radio rows on the provider screen: each
// group's providers plus its "No …" row. That many Down presses from any row
// reach the confirm button.
var providerScreenRows = len(installer.Providers) + len(installer.CodeIntelProviders) + 2

// TestView_CodeIntelOnProviderScreen asserts the provider screen carries the
// code-intelligence group — every provider plus "No code intelligence" — and
// warns, with its install hint, only while a selected provider was not found.
func TestView_CodeIntelOnProviderScreen(t *testing.T) {
	codegraph := installer.CodeIntelProviders[0]
	missing := codegraph.Name + " was not found. ASDT installs bound to it; specialists use grep and read until you install it and index a project."
	for _, found := range []bool{true, false} {
		t.Run(fmt.Sprintf("found=%v", found), func(t *testing.T) {
			m := setup.New(fstest.MapFS{}, "dev")
			m = updateKey(t, m, tea.KeyEnter) // → LanguageSelect
			m = updateKey(t, m, tea.KeyEnter) // → EnvironmentCheck
			msg := envCheckMsg(true)
			msg.CodeIntelFound[codegraph.ID] = found
			next, _ := m.Update(msg)
			m = updateKey(t, next.(setup.Model), tea.KeyEnter) // → SelectAssistants
			m = updateKey(t, m, tea.KeyEnter)                  // → SelectProvider

			view := m.View()
			for _, want := range []string{"Code Intelligence", codegraph.Name, codegraph.Description, "No code intelligence"} {
				if !strings.Contains(view, want) {
					t.Errorf("provider screen missing %q, got:\n%s", want, view)
				}
			}
			if strings.Contains(view, missing) {
				t.Errorf("default selection must show no code-intelligence warning (missing defaults to none), got:\n%s", view)
			}

			// Select Codegraph explicitly: the warning follows found.
			for i := 0; i <= len(installer.Providers); i++ {
				m = updateKey(t, m, tea.KeyDown)
			}
			view = updateKey(t, m, tea.KeySpace).View()
			if got := strings.Contains(view, missing) && strings.Contains(view, "Install: "+codegraph.Detect.InstallURL); got == found {
				t.Errorf("Codegraph selected (found=%v): warning shown = %v, want %v\n%s", found, got, !found, view)
			}
		})
	}
}

// TestView_ReviewCodeIntelRow asserts the Review screen's code-intelligence
// row: the chosen provider's name, a not-found flag with its install hint, or
// "none".
func TestView_ReviewCodeIntelRow(t *testing.T) {
	codegraph := installer.CodeIntelProviders[0]
	cases := []struct {
		name   string
		found  bool
		choose bool // Space on the Codegraph row
		want   string
	}{
		{name: "found by default", found: true, want: "Code intel:  " + codegraph.Name},
		{name: "none by default", want: "Code intel:  none (grep and read)"},
		{name: "chosen while missing", choose: true, want: codegraph.Name + " (not found — install: " + codegraph.Detect.InstallURL + ")"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())
			t.Setenv("XDG_CONFIG_HOME", "")
			m := setup.New(fstest.MapFS{}, "dev")
			m = updateKey(t, m, tea.KeyEnter) // → LanguageSelect
			m = updateKey(t, m, tea.KeyEnter) // → EnvironmentCheck
			msg := envCheckMsg(true)
			msg.CodeIntelFound[codegraph.ID] = c.found
			next, _ := m.Update(msg)
			m = updateKey(t, next.(setup.Model), tea.KeyEnter) // → SelectAssistants
			m = updateKey(t, m, tea.KeyEnter)                  // → SelectProvider
			if c.choose {
				for i := 0; i <= len(installer.Providers); i++ {
					m = updateKey(t, m, tea.KeyDown)
				}
				m = updateKey(t, m, tea.KeySpace)
			}
			for range providerScreenRows {
				m = updateKey(t, m, tea.KeyDown) // → [ Continue → ]
			}
			m = updateKey(t, m, tea.KeyEnter) // → ModelGate
			m = updateKey(t, m, tea.KeyEnter) // → AgentSetup
			for range installer.PersonaPresets {
				m = updateKey(t, m, tea.KeyDown) // → [ Skip → ]
			}
			m = updateKey(t, m, tea.KeyEnter) // → Review
			if m.State() != setup.StateReview {
				t.Fatalf("state = %v, want StateReview", m.State())
			}
			if view := m.View(); !strings.Contains(view, c.want) {
				t.Errorf("review missing %q, got:\n%s", c.want, view)
			}
		})
	}
}

func TestView_AgentWriteModeShowsPerAssistantRows(t *testing.T) {
	view := stateView(t, "AgentWriteMode")
	// The default label for an existing file is "Update ASDT block".
	if !strings.Contains(view, "Update ASDT block") {
		t.Errorf("agent write mode view missing the update label, got:\n%s", view)
	}
	for _, gone := range []string{"Overwrite", "Append", "Keep"} {
		if strings.Contains(view, gone) {
			t.Errorf("agent write mode view still offers %q, got:\n%s", gone, view)
		}
	}
}

func TestView_DoneScreenShowsNextStepHint(t *testing.T) {
	successResult := installer.InstallResult{AssistantID: installer.AssistantClaudeCode, Err: nil}
	m := setup.New(fstest.MapFS{}, "dev")
	next, _ := m.Update(setup.InstallDoneMsg{Results: []installer.InstallResult{successResult}})
	m2 := next.(setup.Model)
	view := m2.View()
	if !strings.Contains(view, "/asdt") {
		t.Errorf("done screen missing next-step hint '/asdt', got:\n%s", view)
	}
}

// TestView_DoneScreenShowsStaleRemovedLine verifies that a result carrying
// pruned files renders the dim "stale file(s) cleaned up" sub-line under the
// assistant's success row.
func TestView_DoneScreenShowsStaleRemovedLine(t *testing.T) {
	result := installer.InstallResult{
		AssistantID: installer.AssistantClaudeCode,
		Removed:     []string{"asdt/old.md", "asdt-core/references/gone.md"},
	}
	m := setup.New(fstest.MapFS{}, "dev")
	next, _ := m.Update(setup.InstallDoneMsg{Results: []installer.InstallResult{result}})
	m2 := next.(setup.Model)
	view := m2.View()
	if !strings.Contains(view, "2 stale file(s) cleaned up") {
		t.Errorf("done screen missing stale-removed line for 2 pruned files, got:\n%s", view)
	}
}

// TestView_DoneScreenOmitsStaleRemovedLineWhenNothingPruned verifies the
// sub-line is absent when no files were pruned.
func TestView_DoneScreenOmitsStaleRemovedLineWhenNothingPruned(t *testing.T) {
	result := installer.InstallResult{AssistantID: installer.AssistantClaudeCode}
	m := setup.New(fstest.MapFS{}, "dev")
	next, _ := m.Update(setup.InstallDoneMsg{Results: []installer.InstallResult{result}})
	m2 := next.(setup.Model)
	view := m2.View()
	if strings.Contains(view, "stale file(s) cleaned up") {
		t.Errorf("done screen must omit the stale-removed line when Removed is empty, got:\n%s", view)
	}
}

// sendEngramFound is a no-op helper kept for caller compatibility.
// New() already starts at StateMainMenu, so no navigation is needed.
func sendEngramFound(_ *testing.T, m setup.Model) setup.Model {
	return m
}

// TestView_InstallingSpinnerTintedWithColorSecondary verifies that the
// rendered spinner glyph is styled with panels.ColorSecondary — the same
// cyan/sky pastel reserved for in-progress indicators across the dashboard
// (mirrors panels.NewSpinner's tint, see specialists panel wiring).
func TestView_InstallingSpinnerTintedWithColorSecondary(t *testing.T) {
	prev := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	defer lipgloss.SetColorProfile(prev)

	view := stateView(t, "Installing")

	tinted := lipgloss.NewStyle().Foreground(panels.ColorSecondary)
	found := false
	for _, frame := range spinner.Dot.Frames {
		if strings.Contains(view, tinted.Render(frame)) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected Installing view to contain a spinner frame tinted with ColorSecondary, got:\n%s", view)
	}
}
