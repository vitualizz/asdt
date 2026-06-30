package installer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// The permission overlay configures Claude Code to run in bypass mode behind a
// self-protecting deny-net. It is a STRUCTURAL JSON merge into the user's global
// ~/.claude/settings.json — a distinct concern from the marker-text CLAUDE.md
// merge in agent_config_adapters.go, hence its own file and test.
//
// Two keys together activate bypass without a per-launch prompt:
//   - permissions.defaultMode = "bypassPermissions"
//   - skipDangerousModePermissionPrompt = true (pre-accepts Claude Code's
//     dangerous-mode disclaimer on the user's behalf)
//
// The deny array is the SOLE real control under bypass (an allow list is inert),
// so it must protect its own definition, match nested secrets, and curb
// session-escaping persistence. The shipped array lives in the reviewed,
// embedded assets/permission-overlay.json; absolute self-protection entries for
// the real user's ~/.claude are injected at write time so self-defense never
// depends on unverified "~"/glob expansion.

const (
	permissionOverlayAsset = "assets/permission-overlay.json"
	overlayLogFileName     = ".asdt-permission-overlay.json"
	overlayVersion         = "1"

	claudeDirName    = ".claude"
	settingsFileName = "settings.json"

	keyPermissions   = "permissions"
	keyDefaultMode   = "defaultMode"
	keyDeny          = "deny"
	keySkipDangerous = "skipDangerousModePermissionPrompt"

	bypassMode = "bypassPermissions"
)

// PermissionOverlayResult is the outcome of writing the permission overlay.
// It mirrors AgentConfigResult so the TUI can report it the same way.
type PermissionOverlayResult struct {
	Written    []string // settings.json paths written (empty when skipped/failed)
	BackupPath string   // pre-write backup path (empty when none kept, or after rollback restore)
	Skipped    bool     // true when the user declined or chose AgentModeSkip
	Err        error    // non-nil on refusal, merge failure, or rollback
}

// OverlayLog is the audit record written alongside the overlay at
// ~/.claude/.asdt-permission-overlay.json. It is best-effort: a log write
// failure never fails the overlay itself.
type OverlayLog struct {
	ConsentedAt        time.Time `json:"consented_at"`
	OverlayVersion     string    `json:"overlay_version"`
	DenySetHash        string    `json:"deny_set_hash"`
	ManagedDeny        []string  `json:"managed_deny"` // exact deny rules ASDT wrote; pruned and re-placed on the next run
	SettingsBackupPath string    `json:"settings_backup_path,omitempty"`
	Accepted           bool      `json:"accepted"`
}

// readOverlayLog loads the prior audit record (best-effort). On a first run or
// any read/parse error it returns a zero-value log (ManagedDeny nil), so the
// reconcile step simply prunes nothing and behaves like a clean first apply.
func readOverlayLog(home string) OverlayLog {
	var log OverlayLog
	data, err := os.ReadFile(filepath.Join(home, claudeDirName, overlayLogFileName))
	if err != nil {
		return log
	}
	_ = json.Unmarshal(data, &log)
	return log
}

// realUserHome resolves the home directory the overlay must target.
//
// Under sudo (euid 0 with SUDO_USER set) it targets the REAL invoking user, not
// root: Claude Code ignores bypass mode for root, so writing root's home would
// be a silent no-op. Running as root WITHOUT SUDO_USER is refused outright
// rather than silently writing the wrong (or useless) file.
func realUserHome() (string, error) {
	if os.Geteuid() == 0 {
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser == "" {
			return "", errors.New("refusing to write the permission overlay as root: re-run as your normal user (Claude Code ignores bypass mode for root, so this would have no effect)")
		}
		u, err := user.Lookup(sudoUser)
		if err != nil {
			return "", fmt.Errorf("cannot resolve home directory for SUDO_USER %q: %w", sudoUser, err)
		}
		return u.HomeDir, nil
	}
	return os.UserHomeDir()
}

// claudeSettingsPath returns <realHome>/.claude/settings.json, the global write
// target. It is the root/sudo-aware sibling of claudeAgentDir.
func claudeSettingsPath() (string, error) {
	home, err := realUserHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, claudeDirName, settingsFileName), nil
}

// loadOverlay decodes the embedded, reviewed overlay asset into a generic map so
// it can be merged by known key path without a typed schema.
func loadOverlay() (map[string]any, error) {
	data, err := assetsFS.ReadFile(permissionOverlayAsset)
	if err != nil {
		return nil, err
	}
	var overlay map[string]any
	if err := json.Unmarshal(data, &overlay); err != nil {
		return nil, fmt.Errorf("embedded permission overlay asset is not valid JSON: %w", err)
	}
	return overlay, nil
}

// selfProtectionEntries returns the HOME-absolute deny rules that guard the
// overlay's own settings file and the rest of ~/.claude. These are injected at
// write time so the deny-net's self-defense does not rely on "~"/glob expansion
// behaving as expected inside Claude Code rule paths.
func selfProtectionEntries(home string) []string {
	claudeDir := filepath.Join(home, claudeDirName)
	targets := []string{
		filepath.Join(claudeDir, settingsFileName),
		claudeDir + "/**",
	}
	verbs := []string{"Edit", "Write", "NotebookEdit"}
	out := make([]string, 0, len(targets)*len(verbs))
	for _, t := range targets {
		for _, v := range verbs {
			out = append(out, fmt.Sprintf("%s(%s)", v, t))
		}
	}
	return out
}

// injectSelfProtection unions the HOME-absolute self-protection rules into the
// overlay's deny array (keep-first dedup) before the structural merge runs.
func injectSelfProtection(overlay map[string]any, home string) {
	perms, _ := overlay[keyPermissions].(map[string]any)
	if perms == nil {
		perms = map[string]any{}
		overlay[keyPermissions] = perms
	}
	perms[keyDeny] = unionDedup(toStringSlice(perms[keyDeny]), selfProtectionEntries(home))
}

// mergePermissionOverlay is the PURE, structural parse-merge-reserialize core.
//
// It decodes existing settings into a generic map so ALL unknown user keys
// round-trip untouched, then applies the overlay by known key path:
//   - deny: the previously-managed overlay rules (prevManaged) are PRUNED first,
//     then the current managed set is unioned in (keep-first dedup). Anything not
//     in prevManaged — the user's own deny rules — always survives. This makes the
//     overlay self-healing: stale entries (e.g. a verb no longer shipped) are
//     cleaned and the current rules re-placed on every run.
//   - defaultMode / skipDangerousModePermissionPrompt: forced under Overwrite;
//     set only-if-absent under Append (preserving the user's existing value).
//
// The effective mode is Overwrite when the file is fresh or has no permissions
// block (a clean write), otherwise the caller's mode (default Append). An
// existing file that is not valid JSON is REFUSED — never clobbered — so the
// user can fix it first.
func mergePermissionOverlay(existing []byte, overlay map[string]any, mode AgentWriteMode, prevManaged []string) ([]byte, error) {
	settings := map[string]any{}
	fresh := len(bytes.TrimSpace(existing)) == 0
	if !fresh {
		if err := json.Unmarshal(existing, &settings); err != nil {
			return nil, fmt.Errorf("existing %s is not valid JSON; refusing to overwrite it — fix or remove the file, then retry: %w", settingsFileName, err)
		}
	}

	overlayPerms, _ := overlay[keyPermissions].(map[string]any)

	_, hasPerms := settings[keyPermissions]
	effMode := mode
	if fresh || !hasPerms {
		effMode = AgentModeOverwrite
	}

	perms, _ := settings[keyPermissions].(map[string]any)
	if perms == nil {
		perms = map[string]any{}
	}

	if overlayDefault, ok := overlayPerms[keyDefaultMode]; ok {
		if _, exists := perms[keyDefaultMode]; effMode == AgentModeOverwrite || !exists {
			perms[keyDefaultMode] = overlayDefault
		}
	}

	userDeny := subtract(toStringSlice(perms[keyDeny]), prevManaged)
	perms[keyDeny] = unionDedup(userDeny, toStringSlice(overlayPerms[keyDeny]))
	settings[keyPermissions] = perms

	if overlaySkip, ok := overlay[keySkipDangerous]; ok {
		if _, exists := settings[keySkipDangerous]; effMode == AgentModeOverwrite || !exists {
			settings[keySkipDangerous] = overlaySkip
		}
	}

	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, err
	}
	out = append(out, '\n')

	if err := validateOverlay(out); err != nil {
		return nil, err
	}
	return out, nil
}

// validateOverlay is the structural post-condition guard. It rejects truncated
// or corrupted output (which would otherwise break Claude Code at launch): the
// bytes must be valid JSON carrying a permissions block with a non-empty deny
// array and both overlay keys present. It deliberately checks PRESENCE, not
// equality, so Append-mode preservation of a user's pre-existing scalar values
// is not flagged as invalid.
func validateOverlay(data []byte) error {
	if !json.Valid(data) {
		return errors.New("merged settings.json is not valid JSON")
	}
	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}
	perms, ok := settings[keyPermissions].(map[string]any)
	if !ok {
		return errors.New("merged settings.json is missing the permissions block")
	}
	if len(toStringSlice(perms[keyDeny])) == 0 {
		return errors.New("merged settings.json has an empty permissions.deny array")
	}
	if _, ok := perms[keyDefaultMode]; !ok {
		return fmt.Errorf("merged settings.json is missing permissions.%s", keyDefaultMode)
	}
	if _, ok := settings[keySkipDangerous]; !ok {
		return fmt.Errorf("merged settings.json is missing %s", keySkipDangerous)
	}
	return nil
}

// WritePermissionOverlay merges the overlay into the user's global settings.json
// with full safety: refuse-on-root, refuse-on-invalid-existing-JSON, pre-write
// backup, atomic temp+fsync+rename, post-write validation, and rollback to the
// backup if the written file fails validation (so Claude Code always launches).
// It is idempotent: re-running unions to an identical deny set with no dupes.
func WritePermissionOverlay(mode AgentWriteMode) (PermissionOverlayResult, error) {
	var result PermissionOverlayResult

	if mode == AgentModeSkip {
		result.Skipped = true
		return result, nil
	}

	home, err := realUserHome()
	if err != nil {
		result.Err = err
		return result, err
	}
	path := filepath.Join(home, claudeDirName, settingsFileName)

	overlay, err := loadOverlay()
	if err != nil {
		result.Err = err
		return result, err
	}
	injectSelfProtection(overlay, home)

	// The managed set is exactly what we are about to write; the previous run's
	// managed set (from the audit log) is pruned first so stale rules are cleaned.
	overlayPerms, _ := overlay[keyPermissions].(map[string]any)
	managedDeny := toStringSlice(overlayPerms[keyDeny])
	prevManaged := readOverlayLog(home).ManagedDeny

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		result.Err = err
		return result, err
	}

	merged, err := mergePermissionOverlay(existing, overlay, mode, prevManaged)
	if err != nil {
		result.Err = err
		return result, err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		result.Err = err
		return result, err
	}

	var backupPath string
	if len(existing) > 0 {
		backupPath = fmt.Sprintf("%s.asdt-backup-%s", path, time.Now().UTC().Format("20060102T150405Z"))
		if err := os.WriteFile(backupPath, existing, 0o600); err != nil {
			result.Err = err
			return result, err
		}
		result.BackupPath = backupPath
	}

	if err := atomicWriteFile(path, merged, 0o644); err != nil {
		result.Err = err
		return result, err
	}

	readback, err := os.ReadFile(path)
	if err == nil {
		err = validateOverlay(readback)
	}
	if err != nil {
		if backupPath != "" {
			if rollbackErr := os.Rename(backupPath, path); rollbackErr == nil {
				result.BackupPath = ""
			}
		}
		wrapped := fmt.Errorf("permission overlay failed post-write validation and was rolled back: %w", err)
		result.Err = wrapped
		return result, wrapped
	}

	result.Written = append(result.Written, path)
	writeOverlayLog(home, OverlayLog{
		ConsentedAt:        time.Now().UTC(),
		OverlayVersion:     overlayVersion,
		DenySetHash:        denySetHash(currentDeny(merged)),
		ManagedDeny:        managedDeny,
		SettingsBackupPath: result.BackupPath,
		Accepted:           true,
	})
	return result, nil
}

// PermissionOverlayExists reports whether the overlay is already active: both
// activation keys set and every shipped overlay deny entry present in the user's
// settings. Used for idempotency and the dashboard probe.
func PermissionOverlayExists() bool {
	path, err := claudeSettingsPath()
	if err != nil {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var settings map[string]any
	if json.Unmarshal(data, &settings) != nil {
		return false
	}
	perms, ok := settings[keyPermissions].(map[string]any)
	if !ok || perms[keyDefaultMode] != bypassMode {
		return false
	}
	if skip, _ := settings[keySkipDangerous].(bool); !skip {
		return false
	}
	overlay, err := loadOverlay()
	if err != nil {
		return false
	}
	overlayPerms, _ := overlay[keyPermissions].(map[string]any)
	have := make(map[string]bool)
	for _, e := range toStringSlice(perms[keyDeny]) {
		have[e] = true
	}
	for _, want := range toStringSlice(overlayPerms[keyDeny]) {
		if !have[want] {
			return false
		}
	}
	return true
}

// atomicWriteFile writes data to a sibling temp file, fsyncs it, then renames
// over the target so a concurrent reader never sees a half-written file.
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp := filepath.Join(dir, fmt.Sprintf("%s.asdt-tmp-%d", settingsFileName, os.Getpid()))
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// writeOverlayLog persists the audit record. Best-effort: errors are swallowed
// so a log failure never fails an otherwise-successful overlay write.
func writeOverlayLog(home string, log OverlayLog) {
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(home, claudeDirName, overlayLogFileName), data, 0o644)
}

// currentDeny extracts the deny array from already-merged settings bytes.
func currentDeny(merged []byte) []string {
	var settings map[string]any
	if json.Unmarshal(merged, &settings) != nil {
		return nil
	}
	perms, _ := settings[keyPermissions].(map[string]any)
	return toStringSlice(perms[keyDeny])
}

// denySetHash is the order-independent fingerprint of a deny set (sha256 over
// the sorted entries), recorded in the audit log.
func denySetHash(deny []string) string {
	sorted := append([]string(nil), deny...)
	sort.Strings(sorted)
	sum := sha256.Sum256([]byte(strings.Join(sorted, "\n")))
	return hex.EncodeToString(sum[:])
}

// toStringSlice coerces a JSON-decoded deny value ([]any of strings) or a native
// []string into []string, dropping any non-string elements.
func toStringSlice(v any) []string {
	switch s := v.(type) {
	case []string:
		return s
	case []any:
		out := make([]string, 0, len(s))
		for _, e := range s {
			if str, ok := e.(string); ok {
				out = append(out, str)
			}
		}
		return out
	}
	return nil
}

// subtract returns the entries of a not present in remove, preserving a's order.
// It strips the previously-written overlay-managed deny rules before the current
// set is re-applied, so a rule that is no longer shipped (e.g. a removed verb) is
// cleaned on the next run while the user's own rules — never in remove — survive.
func subtract(a, remove []string) []string {
	if len(remove) == 0 {
		return a
	}
	drop := make(map[string]bool, len(remove))
	for _, s := range remove {
		drop[s] = true
	}
	out := make([]string, 0, len(a))
	for _, s := range a {
		if !drop[s] {
			out = append(out, s)
		}
	}
	return out
}

// unionDedup returns a ++ b with duplicates removed, preserving first-seen order
// (existing entries first). This makes the merge idempotent and order-stable.
func unionDedup(a, b []string) []string {
	seen := make(map[string]bool, len(a)+len(b))
	out := make([]string, 0, len(a)+len(b))
	for _, list := range [][]string{a, b} {
		for _, s := range list {
			if !seen[s] {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return out
}
