package setup

import (
	"fmt"
	"strings"

	"github.com/vitualizz/asdt/internal/i18n"
	"github.com/vitualizz/asdt/internal/installer"
	"github.com/vitualizz/asdt/internal/setup/components"
	"github.com/vitualizz/asdt/internal/setup/styles"
	"github.com/vitualizz/asdt/internal/tui/panels"
)

// initialPreflightSections returns the seed state for the pre-flight check screen.
// Section titles come from the active catalog; row labels are kept as constants
// — a memory provider's row is its Name — because they also serve as lookup
// keys in rowHasStatus and in probe messages.
func initialPreflightSections(s i18n.InstallerStrings) []components.SectionGroup {
	providerRows := make([]components.CheckRow, len(installer.Providers))
	for i, p := range installer.Providers {
		providerRows[i] = components.CheckRow{Label: p.Name, Status: components.CheckStatusPending}
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
			Title: s.SectionAIEnhancements,
			Rows: []components.CheckRow{
				{Label: "Codegraph", Status: components.CheckStatusPending, SoftWarn: true},
			},
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

	missing := m.missingProviders()
	blocked := m.preflight.done && len(missing) == len(installer.Providers)
	if blocked {
		writeProviderRecovery(&b, s, missing)
	}

	var footer string
	switch {
	case !m.preflight.done:
		footer = panels.RenderKeyboardFooter([]panels.HintGroup{
			{Label: s.HintGroupStatus, Hints: []panels.Hint{{Key: s.HintChecking, Description: s.HintEnvironment}}},
		}, m.width)
	case blocked:
		footer = panels.RenderKeyboardFooter([]panels.HintGroup{
			{Label: s.HintGroupRequired, Hints: providerRequiredHints(s, missing)},
		}, m.width)
	default:
		footer = panels.RenderKeyboardFooter([]panels.HintGroup{
			{Label: s.HintGroupActions, Hints: []panels.Hint{
				{Key: "enter", Description: s.HintContinue},
				{Key: "esc", Description: s.HintBack},
			}},
		}, m.width)
	}
	return frame(s.TitlePreflightCheck, strings.TrimRight(b.String(), "\n"), footer, true)
}

// writeProviderRecovery writes the recovery block for each memory provider the
// environment check did not find: why it is required, where to install it,
// and that the TUI must restart to see it.
func writeProviderRecovery(b *strings.Builder, s i18n.InstallerStrings, missing []installer.ProviderDescriptor) {
	for _, p := range missing {
		fmt.Fprintf(b, "\n  %s\n", styles.Default.Warning.Render(fmt.Sprintf(s.PrefProviderRequired, p.Name)))
		fmt.Fprintf(b, "  %s\n", styles.Default.Dim.Render(fmt.Sprintf(s.PrefProviderInstall, p.Detect.InstallURL)))
	}
	fmt.Fprintf(b, "  %s\n", styles.Default.Dim.Render(s.PrefProviderRestart))
}

// providerRequiredHints returns the footer hints of a screen blocked on missing
// memory providers: one per provider, keyed by its probe binary, then back.
func providerRequiredHints(s i18n.InstallerStrings, missing []installer.ProviderDescriptor) []panels.Hint {
	hints := make([]panels.Hint, 0, len(missing)+1)
	for _, p := range missing {
		hints = append(hints, panels.Hint{Key: p.Detect.Binary, Description: s.HintProviderRequired})
	}
	return append(hints, panels.Hint{Key: "esc", Description: s.HintBack})
}
