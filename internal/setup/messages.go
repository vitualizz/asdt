// Package setup implements the Bubbletea-based installer TUI for ASDT.
package setup

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vitualizz/asdt/internal/i18n"
	"github.com/vitualizz/asdt/internal/installer"
	"github.com/vitualizz/asdt/internal/setup/components"
)

// AssistantInstallProgressMsg is sent per-assistant as each install completes,
// enabling live progress display during StateInstalling.
type AssistantInstallProgressMsg struct {
	Result installer.InstallResult
}

// InstallDoneMsg is kept for test compatibility — handlers that send this
// directly still transition to StateDone. Normal installs use per-assistant
// AssistantInstallProgressMsg messages instead.
type InstallDoneMsg struct {
	Results []installer.InstallResult
}

// AgentInstallDoneMsg is sent when agent config installation completes.
type AgentInstallDoneMsg struct {
	Results []installer.AgentConfigResult
}

// AgentConfigRemoveDoneMsg is sent when removing ASDT's block from the global
// config files completes.
type AgentConfigRemoveDoneMsg struct {
	Results []installer.AgentConfigResult
}

// AgentConfigRemoveCmd runs installer.RemoveAgentConfig in a goroutine and
// sends AgentConfigRemoveDoneMsg when it completes.
func AgentConfigRemoveCmd(assistants []installer.AssistantDescriptor) tea.Cmd {
	return func() tea.Msg {
		return AgentConfigRemoveDoneMsg{Results: installer.RemoveAgentConfig(assistants)}
	}
}

// EnvironmentCheckProgressMsg is sent by each probe in EnvironmentCheckCmd
// as soon as it resolves. One message per row — carries the row label so
// Update() can locate and update the correct CheckRow in preflightSections.
type EnvironmentCheckProgressMsg struct {
	RowLabel string
	Status   components.CheckStatus
	Detail   string
	SoftWarn bool
}

// EnvironmentCheckMsg is the terminal message sent after ALL probes in
// EnvironmentCheckCmd have resolved. Continue is enabled once it arrives.
type EnvironmentCheckMsg struct {
	// ProvidersFound maps each installer.Providers ID to whether that memory
	// provider was found: its Detect.Binary is on PATH, or it has no probe. A
	// missing ID reads as not found. It picks the provider screen's default
	// and its warnings; it never gates Continue.
	ProvidersFound map[installer.ProviderID]bool
	// CodeIntelFound is ProvidersFound for each installer.CodeIntelProviders
	// ID: it picks the code-intelligence default and its warnings.
	CodeIntelFound map[installer.ProviderID]bool
}

// UpdateCheckMsg carries the result of the GitHub latest-release lookup.
type UpdateCheckMsg struct {
	Latest  string // tag_name from GitHub, e.g. "v0.3.0" (empty on error)
	Current string // version injected at build time, echoed back for comparison
	Err     error  // non-nil on any network/decode/non-200; banner suppressed
}

const (
	latestReleaseURL   = "https://api.github.com/repos/vitualizz/asdt/releases/latest"
	updateCheckTimeout = 3 * time.Second
)

// UpdateCheckCmd returns a tea.Cmd that fetches the latest GitHub release tag
// and reports it via UpdateCheckMsg. All errors are carried in the Msg; it
// never panics and never blocks beyond updateCheckTimeout.
func UpdateCheckCmd(current string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), updateCheckTimeout)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestReleaseURL, nil)
		if err != nil {
			return UpdateCheckMsg{Current: current, Err: err}
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return UpdateCheckMsg{Current: current, Err: err}
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			return UpdateCheckMsg{Current: current, Err: fmt.Errorf("github status %d", resp.StatusCode)}
		}
		var payload struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return UpdateCheckMsg{Current: current, Err: err}
		}
		return UpdateCheckMsg{Latest: payload.TagName, Current: current}
	}
}

// newerAvailable reports whether latest is a different, non-empty release than
// current. Fails closed on dev/empty builds: never returns true for an
// unstamped binary. v1 uses normalized inequality, not semver ordering.
func newerAvailable(current, latest string) bool {
	if current == "dev" || current == "" {
		return false // MANDATORY dev-build guard, before any comparison
	}
	c := strings.TrimPrefix(strings.TrimSpace(current), "v")
	l := strings.TrimPrefix(strings.TrimSpace(latest), "v")
	if l == "" {
		return false // fail closed: no banner on empty/malformed remote
	}
	return l != c
}

// LanguagePrefMsg carries the language preference persisted in install
// metadata. An empty Code means no preference was found.
type LanguagePrefMsg struct {
	Code string
}

// LanguagePrefCmd returns a tea.Cmd that scans the known assistants' install
// metadata and reports the first non-empty persisted language preference via
// LanguagePrefMsg. All metadata I/O stays inside the Cmd — never in Update.
func LanguagePrefCmd() tea.Cmd {
	return func() tea.Msg {
		for _, d := range installer.Descriptors {
			if meta, err := installer.ReadInstallMeta(d); err == nil && meta.Language != "" {
				return LanguagePrefMsg{Code: meta.Language}
			}
		}
		return LanguagePrefMsg{}
	}
}

// InstallCmd returns a tea.Batch of per-assistant Cmds, each emitting one
// AssistantInstallProgressMsg when that assistant's install completes.
// Running concurrently lets the TUI update row-by-row as each finishes.
// lang is the language code chosen in the wizard, recorded in install metadata.
// models carries the per-step model selections from StateModelSetup, injected
// into each workflow.yaml as it is written; nil installs files unmodified.
func InstallCmd(assistants []installer.AssistantDescriptor, provider installer.ProviderDescriptor, skillsFS fs.FS, lang string, models map[string]string, opts installer.InstallOptions) tea.Cmd {
	cmds := make([]tea.Cmd, len(assistants))
	for i, a := range assistants {
		a := a // capture loop variable
		cmds[i] = func() tea.Msg {
			results := installer.InstallWithModels([]installer.AssistantDescriptor{a}, provider, skillsFS, lang, models, opts)
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

// AgentInstallCmd runs installer.InstallAgentConfig in a goroutine and sends
// AgentInstallDoneMsg when it completes. useEmojis selects the rendered
// {{emoji_preference}} bullet. localeCode is the full BCP-47 locale whose
// Directive is injected as {{language_directive}}. modes maps each AssistantID
// to its write mode; a nil or empty map defaults all assistants to
// AgentModeUpdate.
func AgentInstallCmd(assistants []installer.AssistantDescriptor, preset installer.PersonaPreset, useEmojis bool, localeCode string, modes map[string]installer.AgentWriteMode) tea.Cmd {
	return func() tea.Msg {
		results := installer.InstallAgentConfig(assistants, preset, useEmojis, localeCode, modes)
		return AgentInstallDoneMsg{Results: results}
	}
}

const environmentCheckTimeout = 5 * time.Second

// EnvironmentCheckCmd fans out environment probes concurrently.
// Each probe sends an EnvironmentCheckProgressMsg immediately when it
// resolves. Uses context.WithTimeout(5s) for all async probes to prevent
// TUI hang on slow PATH resolution. s localizes the probe row details.
func EnvironmentCheckCmd(s i18n.InstallerStrings) tea.Cmd {
	cmds := []tea.Cmd{osProbeCmd(), shellProbeCmd()}
	for _, p := range installer.Providers {
		cmds = append(cmds, providerProbeCmd(p.Name, p.Detect, s.DetailProviderNotFound))
	}
	for _, p := range installer.CodeIntelProviders {
		cmds = append(cmds, providerProbeCmd(p.Name, p.Detect, s.DetailCodeIntelNotFound))
	}
	return tea.Batch(cmds...)
}

func osProbeCmd() tea.Cmd {
	return func() tea.Msg {
		detail := runtime.GOOS + "/" + runtime.GOARCH
		return EnvironmentCheckProgressMsg{
			RowLabel: "OS / Arch",
			Status:   components.CheckStatusOK,
			Detail:   detail,
		}
	}
}

func shellProbeCmd() tea.Cmd {
	return func() tea.Msg {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "unknown"
		}
		return EnvironmentCheckProgressMsg{
			RowLabel: "Shell",
			Status:   components.CheckStatusOK,
			Detail:   shell,
		}
	}
}

// providerProbeCmd resolves a provider's preflight row, labeled label: a soft
// warning carrying notFound (a printf format taking detect.InstallURL) when
// detect.Binary is not on PATH, OK with its path when it is, and OK with an
// informational detail when the provider has no local probe. No outcome
// blocks: every provider — memory or code intelligence — is optional.
func providerProbeCmd(label string, detect installer.ProviderDetect, notFound string) tea.Cmd {
	return func() tea.Msg {
		if detect.Binary == "" {
			return EnvironmentCheckProgressMsg{
				RowLabel: label,
				Status:   components.CheckStatusOK,
				Detail:   "no local check — /asdt-init verifies it is connected",
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), environmentCheckTimeout)
		defer cancel()
		path, err := lookPathCtx(ctx, detect.Binary)
		if err != nil {
			return EnvironmentCheckProgressMsg{
				RowLabel: label,
				Status:   components.CheckStatusWarning,
				Detail:   fmt.Sprintf(notFound, detect.InstallURL),
				SoftWarn: true,
			}
		}
		return EnvironmentCheckProgressMsg{
			RowLabel: label,
			Status:   components.CheckStatusOK,
			Detail:   path,
		}
	}
}

// lookPathCtx wraps exec.LookPath with context cancellation. exec.LookPath
// itself is synchronous — this goroutine+channel pattern enforces the timeout.
func lookPathCtx(ctx context.Context, name string) (string, error) {
	type result struct {
		path string
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		p, err := exec.LookPath(name)
		ch <- result{p, err}
	}()
	select {
	case r := <-ch:
		return r.path, r.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
