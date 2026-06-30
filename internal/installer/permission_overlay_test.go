package installer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// decodeSettings unmarshals merged bytes into a generic map for assertions.
func decodeSettings(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("merged output is not valid JSON: %v\n%s", err, data)
	}
	return m
}

func denyOf(t *testing.T, m map[string]any) []string {
	t.Helper()
	perms, ok := m[keyPermissions].(map[string]any)
	if !ok {
		t.Fatalf("missing permissions block: %v", m)
	}
	return toStringSlice(perms[keyDeny])
}

func denyContains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

// testOverlay returns a fresh decoded copy of the embedded overlay asset so
// tests never mutate shared state.
func testOverlay(t *testing.T) map[string]any {
	t.Helper()
	o, err := loadOverlay()
	if err != nil {
		t.Fatalf("loadOverlay: %v", err)
	}
	return o
}

// TestMergePermissionOverlay_FreshFile: an empty existing file produces a clean
// write with both activation keys set and the full deny array.
func TestMergePermissionOverlay_FreshFile(t *testing.T) {
	out, err := mergePermissionOverlay(nil, testOverlay(t), AgentModeAppend, nil)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	m := decodeSettings(t, out)
	perms := m[keyPermissions].(map[string]any)
	if perms[keyDefaultMode] != bypassMode {
		t.Errorf("defaultMode = %v, want %q", perms[keyDefaultMode], bypassMode)
	}
	if skip, _ := m[keySkipDangerous].(bool); !skip {
		t.Errorf("%s = %v, want true", keySkipDangerous, m[keySkipDangerous])
	}
	if len(denyOf(t, m)) == 0 {
		t.Error("deny array is empty")
	}
}

// TestMergePermissionOverlay_PreservesUserKeys covers AC-1: every existing user
// key round-trips, both overlay keys are present, and the output is valid JSON.
func TestMergePermissionOverlay_PreservesUserKeys(t *testing.T) {
	existing := []byte(`{
  "model": "sonnet",
  "env": {"FOO": "bar"},
  "permissions": {
    "allow": ["Read(**/*.md)"],
    "ask": ["Bash(git push*)"],
    "deny": ["Write(**/custom-secret)"]
  },
  "enabledPlugins": ["x"]
}`)

	out, err := mergePermissionOverlay(existing, testOverlay(t), AgentModeAppend, nil)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	m := decodeSettings(t, out)

	if m["model"] != "sonnet" {
		t.Errorf("model not preserved: %v", m["model"])
	}
	if env, ok := m["env"].(map[string]any); !ok || env["FOO"] != "bar" {
		t.Errorf("env not preserved: %v", m["env"])
	}
	if _, ok := m["enabledPlugins"]; !ok {
		t.Error("enabledPlugins not preserved")
	}
	perms := m[keyPermissions].(map[string]any)
	if _, ok := perms["allow"]; !ok {
		t.Error("permissions.allow not preserved")
	}
	if _, ok := perms["ask"]; !ok {
		t.Error("permissions.ask not preserved")
	}

	// Both overlay keys present (set because they were absent — Append).
	if perms[keyDefaultMode] != bypassMode {
		t.Errorf("defaultMode = %v, want %q", perms[keyDefaultMode], bypassMode)
	}
	if skip, _ := m[keySkipDangerous].(bool); !skip {
		t.Errorf("%s not set", keySkipDangerous)
	}

	// User's pre-existing deny entry survives the union.
	if !denyContains(denyOf(t, m), "Write(**/custom-secret)") {
		t.Error("user deny entry not preserved in union")
	}
}

// TestMergePermissionOverlay_AppendPreservesScalars: in Append mode a user's
// pre-existing defaultMode and skip values are NOT overwritten.
func TestMergePermissionOverlay_AppendPreservesScalars(t *testing.T) {
	existing := []byte(`{"permissions":{"defaultMode":"default","deny":[]},"skipDangerousModePermissionPrompt":false}`)
	out, err := mergePermissionOverlay(existing, testOverlay(t), AgentModeAppend, nil)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	m := decodeSettings(t, out)
	perms := m[keyPermissions].(map[string]any)
	if perms[keyDefaultMode] != "default" {
		t.Errorf("Append clobbered defaultMode: got %v, want preserved \"default\"", perms[keyDefaultMode])
	}
	if skip, _ := m[keySkipDangerous].(bool); skip {
		t.Error("Append clobbered skipDangerousModePermissionPrompt")
	}
	// Deny is still unioned regardless of mode.
	if len(denyOf(t, m)) == 0 {
		t.Error("deny not unioned under Append")
	}
}

// TestMergePermissionOverlay_OverwriteForcesScalars: Overwrite forces the
// overlay's activation values over the user's existing ones.
func TestMergePermissionOverlay_OverwriteForcesScalars(t *testing.T) {
	existing := []byte(`{"permissions":{"defaultMode":"default","deny":[]},"skipDangerousModePermissionPrompt":false}`)
	out, err := mergePermissionOverlay(existing, testOverlay(t), AgentModeOverwrite, nil)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	m := decodeSettings(t, out)
	perms := m[keyPermissions].(map[string]any)
	if perms[keyDefaultMode] != bypassMode {
		t.Errorf("Overwrite defaultMode = %v, want %q", perms[keyDefaultMode], bypassMode)
	}
	if skip, _ := m[keySkipDangerous].(bool); !skip {
		t.Error("Overwrite did not force skipDangerousModePermissionPrompt=true")
	}
}

// TestOverlayAsset_WriteVerbs covers AC-4: the shipped deny array carries the
// real write verbs (Edit/Write/NotebookEdit — MultiEdit is not a Claude Code
// permission-rule tool) for every write-protected path and uses **/ globs for
// nested secrets.
func TestOverlayAsset_WriteVerbs(t *testing.T) {
	out, err := mergePermissionOverlay(nil, testOverlay(t), AgentModeOverwrite, nil)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	deny := denyOf(t, decodeSettings(t, out))

	protected := []string{
		"**/.claude/settings.json", "**/.claude/**",
		"**/.env", "**/.env.*", "**/secrets/**", "**/*.pem", "**/*.key",
		"**/.zshrc", "**/.bashrc", "**/.bash_profile", "**/.profile",
		"**/.envrc", "**/.npmrc", "**/.gitconfig", "**/.mcp.json",
	}
	verbs := []string{"Edit", "Write", "NotebookEdit"}
	for _, p := range protected {
		for _, v := range verbs {
			want := v + "(" + p + ")"
			if !denyContains(deny, want) {
				t.Errorf("deny missing %q", want)
			}
		}
	}

	// Nested-secret **/ read globs and tool-substitution denies.
	for _, want := range []string{
		"Read(**/.env)", "Read(**/secrets/**)", "Read(**/*.pem)", "Read(**/*.key)",
		"Bash(cat **/.env*)", "Bash(rm -rf /)", "Bash(rm -rf ~)",
	} {
		if !denyContains(deny, want) {
			t.Errorf("deny missing %q", want)
		}
	}

	// No bare "~" path forms in the embedded asset (self-protection absolutes
	// are injected at write time, not shipped).
	for _, e := range deny {
		if strings.Contains(e, "(~") {
			t.Errorf("embedded asset must not ship bare ~ form: %q", e)
		}
	}
}

// TestMergePermissionOverlay_Idempotent: re-merging the previous output yields
// identical bytes and an identical deny-set hash (no duplicate entries).
func TestMergePermissionOverlay_Idempotent(t *testing.T) {
	first, err := mergePermissionOverlay(nil, testOverlay(t), AgentModeAppend, nil)
	if err != nil {
		t.Fatalf("first merge: %v", err)
	}
	second, err := mergePermissionOverlay(first, testOverlay(t), AgentModeAppend, nil)
	if err != nil {
		t.Fatalf("second merge: %v", err)
	}
	if string(first) != string(second) {
		t.Errorf("merge not idempotent:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if denySetHash(currentDeny(first)) != denySetHash(currentDeny(second)) {
		t.Error("deny-set hash drifted across idempotent re-run")
	}

	// No duplicate deny entries.
	deny := currentDeny(second)
	seen := map[string]bool{}
	for _, e := range deny {
		if seen[e] {
			t.Errorf("duplicate deny entry: %q", e)
		}
		seen[e] = true
	}
}

// TestMergePermissionOverlay_ReconcilesStaleManaged: a previously-managed entry
// that is no longer shipped (e.g. a removed MultiEdit verb) is pruned on re-merge,
// the current managed rules are re-placed, and the user's OWN deny rule survives.
func TestMergePermissionOverlay_ReconcilesStaleManaged(t *testing.T) {
	existing := []byte(`{
  "permissions": {
    "deny": [
      "Write(**/custom-secret)",
      "MultiEdit(**/.env)",
      "Edit(**/.env)"
    ]
  }
}`)
	// What the previous (buggy) run wrote and recorded as managed.
	prevManaged := []string{"MultiEdit(**/.env)", "Edit(**/.env)"}

	out, err := mergePermissionOverlay(existing, testOverlay(t), AgentModeAppend, prevManaged)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	deny := denyOf(t, decodeSettings(t, out))

	if denyContains(deny, "MultiEdit(**/.env)") {
		t.Error("stale managed entry MultiEdit(**/.env) was not pruned")
	}
	if !denyContains(deny, "Write(**/custom-secret)") {
		t.Error("user's own deny rule was lost during reconciliation")
	}
	if !denyContains(deny, "Edit(**/.env)") {
		t.Error("current managed entry Edit(**/.env) missing after re-place")
	}
}

// TestMergePermissionOverlay_RefusesInvalidExistingJSON: a present-but-invalid
// existing file is refused, never clobbered.
func TestMergePermissionOverlay_RefusesInvalidExistingJSON(t *testing.T) {
	_, err := mergePermissionOverlay([]byte("{not valid json"), testOverlay(t), AgentModeAppend, nil)
	if err == nil {
		t.Fatal("expected refusal on invalid existing JSON, got nil error")
	}
	if !strings.Contains(err.Error(), "valid JSON") {
		t.Errorf("error should explain the invalid-JSON refusal, got: %v", err)
	}
}

// TestValidateOverlay_RejectsInvalid covers AC-5's trigger: the post-condition
// guard that drives rollback rejects corrupted or incomplete output.
func TestValidateOverlay_RejectsInvalid(t *testing.T) {
	cases := []struct {
		name string
		data string
	}{
		{"not json", "{broken"},
		{"missing permissions", `{"skipDangerousModePermissionPrompt":true}`},
		{"empty deny", `{"permissions":{"defaultMode":"bypassPermissions","deny":[]},"skipDangerousModePermissionPrompt":true}`},
		{"missing defaultMode", `{"permissions":{"deny":["Write(x)"]},"skipDangerousModePermissionPrompt":true}`},
		{"missing skip key", `{"permissions":{"defaultMode":"bypassPermissions","deny":["Write(x)"]}}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := validateOverlay([]byte(c.data)); err == nil {
				t.Errorf("validateOverlay accepted invalid input %q", c.data)
			}
		})
	}

	// A well-formed merged document validates.
	good, err := mergePermissionOverlay(nil, testOverlay(t), AgentModeOverwrite, nil)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	if err := validateOverlay(good); err != nil {
		t.Errorf("validateOverlay rejected a valid document: %v", err)
	}
}

// TestSelfProtectionEntries: HOME-absolute self-defense rules cover both targets
// across the three write verbs (Edit/Write/NotebookEdit) and contain no bare "~".
func TestSelfProtectionEntries(t *testing.T) {
	entries := selfProtectionEntries("/home/alice")
	if len(entries) != 6 {
		t.Fatalf("want 6 self-protection entries, got %d: %v", len(entries), entries)
	}
	for _, e := range entries {
		if !strings.Contains(e, "/home/alice/.claude") {
			t.Errorf("self-protection entry not HOME-absolute: %q", e)
		}
		if strings.Contains(e, "(~") {
			t.Errorf("self-protection entry must not use ~: %q", e)
		}
	}
}

// TestWritePermissionOverlay_Skip: AgentModeSkip writes nothing (AC-2 backing).
func TestWritePermissionOverlay_Skip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")

	res, err := WritePermissionOverlay(AgentModeSkip)
	if err != nil {
		t.Fatalf("skip returned error: %v", err)
	}
	if !res.Skipped || len(res.Written) != 0 {
		t.Errorf("skip should write nothing: %+v", res)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude", "settings.json")); !os.IsNotExist(err) {
		t.Errorf("settings.json should not exist after skip, stat err = %v", err)
	}
}

// TestWritePermissionOverlay_FreshAndIdempotent: a real write to a temp HOME
// activates the overlay; a second run is idempotent and produces a backup.
func TestWritePermissionOverlay_FreshAndIdempotent(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root: realUserHome refuses without SUDO_USER by design")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	path := filepath.Join(home, ".claude", "settings.json")

	res, err := WritePermissionOverlay(AgentModeAppend)
	if err != nil {
		t.Fatalf("first write: %v (%+v)", err, res)
	}
	if len(res.Written) != 1 {
		t.Fatalf("expected one written path, got %+v", res)
	}
	if res.BackupPath != "" {
		t.Errorf("fresh write should not produce a backup, got %q", res.BackupPath)
	}
	if !PermissionOverlayExists() {
		t.Error("PermissionOverlayExists false after a successful write")
	}
	first, _ := os.ReadFile(path)

	res2, err := WritePermissionOverlay(AgentModeAppend)
	if err != nil {
		t.Fatalf("second write: %v", err)
	}
	if res2.BackupPath == "" {
		t.Error("re-run over an existing file should produce a backup")
	}
	second, _ := os.ReadFile(path)
	if string(first) != string(second) {
		t.Error("re-run was not idempotent on disk")
	}
}

// TestWritePermissionOverlay_RefusesInvalidExisting: an invalid existing file is
// not clobbered, and no backup/write occurs.
func TestWritePermissionOverlay_RefusesInvalidExisting(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root: realUserHome refuses without SUDO_USER by design")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SUDO_USER", "")
	dir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "settings.json")
	original := []byte("{ this is not json")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := WritePermissionOverlay(AgentModeAppend)
	if err == nil {
		t.Fatal("expected refusal on invalid existing settings.json")
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(original) {
		t.Error("invalid existing settings.json was modified despite refusal")
	}
}
