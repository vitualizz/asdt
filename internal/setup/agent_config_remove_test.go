package setup_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vitualizz/asdt/internal/installer"
	"github.com/vitualizz/asdt/internal/setup"
)

const (
	testUserContent = "# Mine\n\nKeep me.\n"
	testASDTBlock   = "<!-- asdt:agent-config -->\n# Agent\n<!-- /asdt:agent-config -->\n"
)

// agentConfigPaths points HOME and XDG_CONFIG_HOME at a temp tree and returns
// the Claude Code and OpenCode global config paths inside it.
func agentConfigPaths(t *testing.T) (claude, opencode string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	return filepath.Join(home, ".claude", "CLAUDE.md"), filepath.Join(home, ".config", "opencode", "AGENTS.md")
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func openDashboard(t *testing.T) setup.Model {
	t.Helper()
	next, _ := setup.New(fstest.MapFS{}, "dev").Update(tea.WindowSizeMsg{Width: 200})
	m := next.(setup.Model)           // wide enough for the footer to show hint descriptions
	m = updateKey(t, m, tea.KeyDown)  // cursor → Dashboard
	m = updateKey(t, m, tea.KeyEnter) // → StateDashboard
	if m.State() != setup.StateDashboard {
		t.Fatalf("state = %v, want StateDashboard", m.State())
	}
	return m
}

// pressRunes sends a rune key and returns the model and the command it produced.
func pressRunes(t *testing.T, m setup.Model, r string) (setup.Model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(r)})
	return next.(setup.Model), cmd
}

// TestDashboard_RemoveAgentConfig: r asks for confirmation, Esc cancels with
// nothing written, Enter removes ASDT's block from every file that has one and
// leaves the user's text byte-identical.
func TestDashboard_RemoveAgentConfig(t *testing.T) {
	cases := []struct {
		name          string
		seedOpenCode  bool
		wantRemovedOC bool
	}{
		{name: "both assistants", seedOpenCode: true, wantRemovedOC: true},
		{name: "claude only", seedOpenCode: false, wantRemovedOC: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			claudePath, ocPath := agentConfigPaths(t)
			writeTestFile(t, claudePath, testUserContent+"\n"+testASDTBlock)
			if c.seedOpenCode {
				writeTestFile(t, ocPath, testUserContent+"\n"+testASDTBlock)
			}

			m := openDashboard(t)
			view := m.View()
			for _, want := range []string{"agent config", "remove ASDT block"} {
				if !strings.Contains(view, want) {
					t.Errorf("dashboard missing %q, got:\n%s", want, view)
				}
			}

			m, _ = pressRunes(t, m, "r")
			if view := m.View(); !strings.Contains(view, "Remove ASDT's block from") || !strings.Contains(view, "CLAUDE.md") {
				t.Errorf("confirmation prompt missing, got:\n%s", view)
			}
			m = updateKey(t, m, tea.KeyEsc) // cancel
			if m.State() != setup.StateDashboard || strings.Contains(m.View(), "Remove ASDT's block from") {
				t.Errorf("Esc did not cancel the confirmation: state %v, view:\n%s", m.State(), m.View())
			}
			if got := readTestFile(t, claudePath); got != testUserContent+"\n"+testASDTBlock {
				t.Errorf("cancel changed CLAUDE.md: %q", got)
			}

			m, _ = pressRunes(t, m, "r")
			next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			m = next.(setup.Model)
			if cmd == nil {
				t.Fatal("Enter on the confirmation returned no command")
			}
			next, _ = m.Update(cmd())
			m = next.(setup.Model)

			if got := readTestFile(t, claudePath); got != testUserContent {
				t.Errorf("CLAUDE.md after remove:\n got %q\nwant %q", got, testUserContent)
			}
			if c.wantRemovedOC {
				if got := readTestFile(t, ocPath); got != testUserContent {
					t.Errorf("AGENTS.md after remove:\n got %q\nwant %q", got, testUserContent)
				}
			} else if _, err := os.Stat(ocPath); !os.IsNotExist(err) {
				t.Errorf("remove created AGENTS.md, stat err = %v", err)
			}

			view = m.View()
			if !strings.Contains(view, "Claude Code: ASDT block removed") {
				t.Errorf("dashboard missing the removal result, got:\n%s", view)
			}
			if strings.Contains(view, "remove ASDT block") {
				t.Errorf("remove hint still offered with no block left, got:\n%s", view)
			}
		})
	}
}

// TestDashboard_RemoveAgentConfig_NoBlockNoAction: without a block anywhere, r
// does nothing and the hint is not offered.
func TestDashboard_RemoveAgentConfig_NoBlockNoAction(t *testing.T) {
	claudePath, _ := agentConfigPaths(t)
	writeTestFile(t, claudePath, testUserContent)

	m := openDashboard(t)
	if strings.Contains(m.View(), "remove ASDT block") {
		t.Errorf("remove hint offered with no block, got:\n%s", m.View())
	}
	m, cmd := pressRunes(t, m, "r")
	if cmd != nil || strings.Contains(m.View(), "Remove ASDT's block from") {
		t.Errorf("r with no block asked to remove, view:\n%s", m.View())
	}
}

// TestAgentWriteMode_UserContentIsAConflict: a global config with only the
// user's own text still counts as existing, so the TUI asks before adding
// ASDT's block to it, and the prompt offers exactly update or skip.
func TestAgentWriteMode_UserContentIsAConflict(t *testing.T) {
	claudePath, _ := agentConfigPaths(t)
	writeTestFile(t, claudePath, testUserContent)

	m := setup.New(fstest.MapFS{}, "dev")
	m = advanceToSelectProvider(t, m)
	m = updateKey(t, m, tea.KeyEnter) // → ModelGate
	m = updateKey(t, m, tea.KeyEnter) // → AgentSetup
	m = updateKey(t, m, tea.KeyEnter) // → EmojiPref
	m = updateKey(t, m, tea.KeyEnter) // → AgentWriteMode
	if m.State() != setup.StateAgentWriteMode {
		t.Fatalf("state = %v, want StateAgentWriteMode", m.State())
	}
	modes := m.AgentWriteModes()
	if len(modes) != 1 || modes[string(installer.AssistantClaudeCode)] != installer.AgentModeUpdate {
		t.Errorf("write modes = %v, want only Claude Code at AgentModeUpdate", modes)
	}
	m = updateKeyMsg(t, m, tea.KeyMsg{Type: tea.KeySpace})
	if !strings.Contains(m.View(), "[Skip]") {
		t.Errorf("Space did not switch to Skip, got:\n%s", m.View())
	}
}
