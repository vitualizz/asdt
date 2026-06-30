package setup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vitualizz/asdt/internal/setup/styles"
	"github.com/vitualizz/asdt/internal/tui/panels"
)

// The consent title, copy, and option labels are deliberately ENGLISH
// constants, not catalog strings: the permission overlay is a security contract
// whose exact wording is reviewed and must not drift by locale. Only the
// keyboard-hint labels below come from the active catalog.
const permissionConsentTitle = "Permission Overlay (opt-in)"

// permissionConsentCopy discloses the residual risks the user accepts by opting
// in, names the skipDangerousModePermissionPrompt pre-acceptance explicitly,
// states what the deny-net still blocks, and ends with the trust caveat.
const permissionConsentCopy = `ASDT can configure Claude Code to run without per-action approval prompts
(bypass mode), behind a safety net that blocks known-dangerous paths. This is
opt-in and you can decline — declining changes nothing on your machine.

To make bypass take effect without a launch prompt, ASDT also pre-accepts Claude
Code's dangerous-mode disclaimer on your behalf (it sets
skipDangerousModePermissionPrompt: true in your global settings).

Even with the safety net, the following risks REMAIN and you accept them by
continuing:
 • Prompt injection — hostile text in a file, pull request, web page, or tool
   result can be treated as your instruction and acted on with no prompt.
 • Data exfiltration — the agent could read accessible data and send it out
   (e.g. via curl or git); the send channel itself cannot be blocked.
 • Force-push / history rewrite — git push --force and similar run without
   confirmation.
 • Destructive deletes — in-session file or branch deletion in your projects
   runs without confirmation.
 • Dependency supply-chain — installing packages can run their scripts; this is
   not blocked.
 • No per-action audit — bypass removes the approval prompt, so there is no
   per-action record of what was authorized.

The safety net DOES block: editing the agent's own settings, reading or editing
secret files (.env, secrets, *.pem, *.key), and writing shell/startup dotfiles
that would persist outside the session.

Continue only on machines and projects you trust.`

const (
	permissionConsentAcceptLabel  = "Accept — write the permission overlay"
	permissionConsentDeclineLabel = "Decline — change nothing"
)

// renderPermissionConsent renders the opt-in consent gate: the English contract
// copy followed by two radio rows (Accept / Decline). Decline is the safe
// default, so the model enters this screen with the cursor on Decline.
func renderPermissionConsent(m Model) string {
	s := m.catalog.Installer

	warn := styles.Default.Warning.Render("Review carefully — this changes how Claude Code runs on this machine.")
	copyBlock := styles.Default.Dim.Render(permissionConsentCopy)

	var rows strings.Builder
	options := []string{permissionConsentAcceptLabel, permissionConsentDeclineLabel}
	for i, opt := range options {
		focused := i == m.cursor
		selected := (i == 0) == m.permission.consent

		cursor := "  "
		if focused {
			cursor = cursorChar + " "
		}

		var radioStr string
		if selected {
			radioStr = styles.Default.Cursor.Render("(•)")
		} else {
			radioStr = styles.Default.Dim.Render("( )")
		}

		var nameStr string
		if focused {
			nameStr = styles.Default.Cursor.Render(opt)
		} else {
			nameStr = styles.Default.Dim.Render(opt)
		}

		fmt.Fprintf(&rows, "  %s%s %s\n", cursor, radioStr, nameStr)
	}

	body := lipgloss.JoinVertical(lipgloss.Left, warn, "", copyBlock, "", strings.TrimRight(rows.String(), "\n"))

	footer := panels.RenderKeyboardFooter([]panels.HintGroup{
		{Label: s.HintGroupActions, Hints: []panels.Hint{
			{Key: "↑↓", Description: s.HintNavigate},
			{Key: "enter", Description: s.HintContinue},
			{Key: "esc", Description: s.HintBack},
			{Key: "q", Description: s.HintQuit},
		}},
	}, m.width)
	return frame(permissionConsentTitle, body, footer, true)
}
