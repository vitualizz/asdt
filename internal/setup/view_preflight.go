package setup

import (
	"fmt"
	"strings"

	"github.com/vitualizz/asdt/internal/i18n"
	"github.com/vitualizz/asdt/internal/installer"
	"github.com/vitualizz/asdt/internal/setup/components"
	"github.com/vitualizz/asdt/internal/tui/panels"
)

// initialPreflightSections returns the seed state for the pre-flight check screen.
// Every provider row — memory or code intelligence — is a soft warning: a
// provider is optional, so a missing one never blocks Continue.
// Section titles come from the active catalog; row labels are kept as constants
// — a provider's row is its Name — because they also serve as lookup keys in
// rowHasStatus and in probe messages.
func initialPreflightSections(s i18n.InstallerStrings) []components.SectionGroup {
	providerRows := make([]components.CheckRow, len(installer.Providers))
	for i, p := range installer.Providers {
		providerRows[i] = components.CheckRow{Label: p.Name, Status: components.CheckStatusPending, SoftWarn: true}
	}
	codeIntelRows := make([]components.CheckRow, len(installer.CodeIntelProviders))
	for i, p := range installer.CodeIntelProviders {
		codeIntelRows[i] = components.CheckRow{Label: p.Name, Status: components.CheckStatusPending, SoftWarn: true}
	}
	return []components.SectionGroup{
		{
			Title: s.SectionYourEnvironment,
			Rows: []components.CheckRow{
				{Label: "OS / Arch", Status: components.CheckStatusPending},
				{Label: "Shell", Status: components.CheckStatusPending},
			},
		},
		{
			Title: s.SectionMemoryProvider,
			Rows:  providerRows,
		},
		{
			Title: s.SectionCodeIntelligence,
			Rows:  codeIntelRows,
		},
	}
}

func renderPreflightCheck(m Model) string {
	s := m.catalog.Installer
	var b strings.Builder

	fmt.Fprintf(&b, "  %s\n\n", stepLine(s, 1, 6))

	for _, sg := range m.preflight.sections {
		b.WriteString(sg.Render(m.width))
		b.WriteString("\n")
	}

	var footer string
	if m.preflight.done {
		footer = panels.RenderKeyboardFooter([]panels.HintGroup{
			{Label: s.HintGroupActions, Hints: []panels.Hint{
				{Key: "enter", Description: s.HintContinue},
				{Key: "esc", Description: s.HintBack},
			}},
		}, m.width)
	} else {
		footer = panels.RenderKeyboardFooter([]panels.HintGroup{
			{Label: s.HintGroupStatus, Hints: []panels.Hint{{Key: s.HintChecking, Description: s.HintEnvironment}}},
		}, m.width)
	}
	return frame(s.TitlePreflightCheck, strings.TrimRight(b.String(), "\n"), footer, true)
}
