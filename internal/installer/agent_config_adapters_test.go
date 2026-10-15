package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testBlockBody = "# Agent\nbody"

// testBlock is the exact block upsertAgentBlock writes for testBlockBody.
var testBlock = asdtBlockStart + "\n" + testBlockBody + "\n" + asdtBlockEnd

// malformedBlockCases lists every marker state that must fail loudly.
var malformedBlockCases = []struct {
	name    string
	content string
}{
	{"start only", "# Mine\n\n" + asdtBlockStart + "\nold\n"},
	{"end only", "# Mine\n\nold\n" + asdtBlockEnd + "\n"},
	{"duplicate blocks", testBlock + "\n\n" + testBlock + "\n"},
	{"duplicate start", asdtBlockStart + "\n" + testBlock + "\n"},
	{"reversed", asdtBlockEnd + "\nold\n" + asdtBlockStart + "\n"},
}

// inlineMarkerProse mentions both markers inside lines of prose. Only a marker
// alone on its own line is ASDT's, so this is plain user text.
var inlineMarkerProse = "# Mine\n\nASDT wraps its block in `" + asdtBlockStart + "` ... `" + asdtBlockEnd + "`.\n"

// agentHome points HOME at a fresh temp dir and clears XDG_CONFIG_HOME, so
// every adapter path resolves under the temp HOME, never the real one.
func agentHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	return home
}

// seedAgentFile writes content to path, creating its directory.
func seedAgentFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readAgentFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func TestUpsertAgentBlock(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{"empty file", "", testBlock + "\n"},
		{"user content", "# Mine\n", "# Mine\n\n" + testBlock + "\n"},
		{"user content without final newline", "# Mine", "# Mine\n\n" + testBlock + "\n"},
		{"markers quoted inline are user text", inlineMarkerProse, inlineMarkerProse + "\n" + testBlock + "\n"},
		{
			"existing block replaced in place",
			"# Mine\n\n" + asdtBlockStart + "\nold\nlines\n" + asdtBlockEnd + "\n\n# After\n",
			"# Mine\n\n" + testBlock + "\n\n# After\n",
		},
		{
			"existing block without final newline",
			"# Mine\n\n" + asdtBlockStart + "\nold\n" + asdtBlockEnd,
			"# Mine\n\n" + testBlock,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := upsertAgentBlock(c.content, testBlockBody+"\n")
			if err != nil {
				t.Fatalf("upsertAgentBlock: %v", err)
			}
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
			again, err := upsertAgentBlock(got, testBlockBody+"\n")
			if err != nil || again != got {
				t.Errorf("second upsert not idempotent: err=%v\n got %q\nwant %q", err, again, got)
			}
		})
	}
	for _, c := range malformedBlockCases {
		t.Run("malformed "+c.name, func(t *testing.T) {
			if _, err := upsertAgentBlock(c.content, testBlockBody); err == nil {
				t.Errorf("upsertAgentBlock accepted malformed markers in %q", c.content)
			}
		})
	}
}

func TestRemoveAgentBlock(t *testing.T) {
	cases := []struct {
		name      string
		content   string
		want      string
		wantFound bool
	}{
		{"no block", "# Mine\n", "# Mine\n", false},
		{"only the block", testBlock + "\n", "", true},
		{"block appended after user content", "# Mine\n\n" + testBlock + "\n", "# Mine\n", true},
		{"block between user content", "# Mine\n\n" + testBlock + "\n\n# After\n", "# Mine\n\n# After\n", true},
		{"block opening the file", testBlock + "\n\n# After\n", "# After\n", true},
		{"markers quoted inline are not a block", inlineMarkerProse, inlineMarkerProse, false},
		{"inline mention beside a real block", inlineMarkerProse + "\n" + testBlock + "\n", inlineMarkerProse, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, found, err := removeAgentBlock(c.content)
			if err != nil {
				t.Fatalf("removeAgentBlock: %v", err)
			}
			if found != c.wantFound || got != c.want {
				t.Errorf("got (%q, %v), want (%q, %v)", got, found, c.want, c.wantFound)
			}
		})
	}
	for _, c := range malformedBlockCases {
		t.Run("malformed "+c.name, func(t *testing.T) {
			if _, _, err := removeAgentBlock(c.content); err == nil {
				t.Errorf("removeAgentBlock accepted malformed markers in %q", c.content)
			}
		})
	}
}

// TestWriteAgentBlock_BothAssistants drives the shared write path through each
// adapter's real config path: the file ends up holding the user's text
// byte-identical plus exactly one block, and a second write changes nothing.
func TestWriteAgentBlock_BothAssistants(t *testing.T) {
	scenarios := []struct {
		name     string
		existing *string // nil: no file
		want     string
	}{
		{"no file", nil, testBlock + "\n"},
		{"user content", strPtr("# Mine\n\nKeep me.\n"), "# Mine\n\nKeep me.\n\n" + testBlock + "\n"},
		{
			"existing block",
			strPtr("# Mine\n\n" + asdtBlockStart + "\nold\n" + asdtBlockEnd + "\n# Tail\n"),
			"# Mine\n\n" + testBlock + "\n# Tail\n",
		},
	}
	for _, adapter := range AgentConfigAdapters {
		for _, sc := range scenarios {
			t.Run(string(adapter.AssistantID)+"/"+sc.name, func(t *testing.T) {
				agentHome(t)
				path := adapter.ConfigPath()
				if sc.existing != nil {
					seedAgentFile(t, path, *sc.existing)
				}

				if err := writeAgentBlock(path, testBlockBody); err != nil {
					t.Fatalf("first write: %v", err)
				}
				first := readAgentFile(t, path)
				if first != sc.want {
					t.Errorf("after write:\n got %q\nwant %q", first, sc.want)
				}

				if err := writeAgentBlock(path, testBlockBody); err != nil {
					t.Fatalf("second write: %v", err)
				}
				if second := readAgentFile(t, path); second != first {
					t.Errorf("second write not byte-identical:\n got %q\nwant %q", second, first)
				}
			})
		}
	}
}

// TestWriteAgentBlock_MalformedMarkersWriteNothing: a malformed marker state is
// an error and the file is left byte-identical.
func TestWriteAgentBlock_MalformedMarkersWriteNothing(t *testing.T) {
	for _, adapter := range AgentConfigAdapters {
		for _, c := range malformedBlockCases {
			t.Run(string(adapter.AssistantID)+"/"+c.name, func(t *testing.T) {
				agentHome(t)
				path := adapter.ConfigPath()
				seedAgentFile(t, path, c.content)

				if err := writeAgentBlock(path, testBlockBody); err == nil {
					t.Fatal("writeAgentBlock accepted malformed markers")
				}
				if got := readAgentFile(t, path); got != c.content {
					t.Errorf("file changed despite the error:\n got %q\nwant %q", got, c.content)
				}
			})
		}
	}
}

func TestOpenCodeConfigPath_RespectsXDGConfigHome(t *testing.T) {
	agentHome(t)
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	adapter, _ := AgentConfigAdapterFor(AssistantOpenCode)
	if got, want := adapter.ConfigPath(), filepath.Join(xdg, "opencode", "AGENTS.md"); got != want {
		t.Errorf("ConfigPath = %q, want %q", got, want)
	}
}

func TestAgentConfigAdapter_ConfigExistsAndHasBlock(t *testing.T) {
	cases := []struct {
		name         string
		content      *string
		wantExists   bool
		wantHasBlock bool
	}{
		{"no file", nil, false, false},
		{"empty file", strPtr(""), false, false},
		{"user content", strPtr("# Mine\n"), true, false},
		{"block", strPtr(testBlock + "\n"), true, true},
		{"lone end marker", strPtr(asdtBlockEnd + "\n"), true, true},
		{"markers quoted inline", strPtr(inlineMarkerProse), true, false},
		{"CRLF block", strPtr(asdtBlockStart + "\r\nold\r\n" + asdtBlockEnd + "\r\n"), true, true},
	}
	for _, adapter := range AgentConfigAdapters {
		for _, c := range cases {
			t.Run(string(adapter.AssistantID)+"/"+c.name, func(t *testing.T) {
				agentHome(t)
				if c.content != nil {
					seedAgentFile(t, adapter.ConfigPath(), *c.content)
				}
				if got := adapter.ConfigExists(); got != c.wantExists {
					t.Errorf("ConfigExists = %v, want %v", got, c.wantExists)
				}
				if got := adapter.HasBlock(); got != c.wantHasBlock {
					t.Errorf("HasBlock = %v, want %v", got, c.wantHasBlock)
				}
			})
		}
	}
}

// TestRemoveAgentConfig_RestoresUserContent: install then remove leaves the
// user's file byte-identical, for both assistants, except that a missing final
// newline is added; the file is never deleted.
func TestRemoveAgentConfig_RestoresUserContent(t *testing.T) {
	originals := []struct {
		name    string
		content *string
		want    string
	}{
		{"user content", strPtr("# Mine\n\nKeep me.\n"), "# Mine\n\nKeep me.\n"},
		{"user content without final newline gains one", strPtr("# Mine\n\nKeep me."), "# Mine\n\nKeep me.\n"},
		{"markers quoted inline", strPtr(inlineMarkerProse), inlineMarkerProse},
		{"no prior file", nil, ""},
	}
	for _, o := range originals {
		t.Run(o.name, func(t *testing.T) {
			home := agentHome(t)
			assistants := testAgentAssistants(t, home)
			for _, a := range AgentConfigAdapters {
				if o.content != nil {
					seedAgentFile(t, a.ConfigPath(), *o.content)
				}
			}

			for _, r := range InstallAgentConfig(assistants, PersonaPresets[0], true, "en", nil) {
				if r.Err != nil || r.Skipped {
					t.Fatalf("install %s: %+v", r.AssistantID, r)
				}
			}
			results := RemoveAgentConfig(assistants)
			for i, a := range AgentConfigAdapters {
				r := results[i]
				if r.Err != nil || r.Skipped || len(r.Written) != 1 {
					t.Errorf("remove %s: %+v", r.AssistantID, r)
				}
				if got := readAgentFile(t, a.ConfigPath()); got != o.want {
					t.Errorf("%s after remove:\n got %q\nwant %q", a.AssistantID, got, o.want)
				}
				meta, err := ReadInstallMeta(assistants[i])
				if err != nil {
					t.Fatalf("read meta: %v", err)
				}
				if meta.Persona != "" || meta.Emojis != "" {
					t.Errorf("%s meta keeps persona after remove: %+v", a.AssistantID, meta)
				}
			}
		})
	}
}

func TestRemoveAgentConfig_NothingToRemove(t *testing.T) {
	cases := []struct {
		name    string
		content *string
		wantErr bool
	}{
		{"no file", nil, false},
		{"no block", strPtr("# Mine\n"), false},
		{"malformed markers", strPtr("# Mine\n\n" + asdtBlockStart + "\n"), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			home := agentHome(t)
			for _, a := range AgentConfigAdapters {
				if c.content != nil {
					seedAgentFile(t, a.ConfigPath(), *c.content)
				}
			}
			results := RemoveAgentConfig(testAgentAssistants(t, home))
			for i, a := range AgentConfigAdapters {
				r := results[i]
				if (r.Err != nil) != c.wantErr || (!c.wantErr && !r.Skipped) || len(r.Written) != 0 {
					t.Errorf("%s: %+v, wantErr %v", a.AssistantID, r, c.wantErr)
				}
				_, statErr := os.Stat(a.ConfigPath())
				if c.content == nil {
					if !os.IsNotExist(statErr) {
						t.Errorf("%s: remove created a file", a.AssistantID)
					}
					continue
				}
				if got := readAgentFile(t, a.ConfigPath()); got != *c.content {
					t.Errorf("%s changed:\n got %q\nwant %q", a.AssistantID, got, *c.content)
				}
			}
		})
	}
}

// testAgentAssistants returns one descriptor per agent config adapter, each with
// a skills dir under home so install metadata reads and writes stay in the
// temp tree.
func testAgentAssistants(t *testing.T, home string) []AssistantDescriptor {
	t.Helper()
	out := make([]AssistantDescriptor, 0, len(AgentConfigAdapters))
	for _, a := range AgentConfigAdapters {
		dir := filepath.Join(home, "skills-"+string(a.AssistantID))
		if err := os.MkdirAll(filepath.Join(dir, "asdt"), 0o755); err != nil {
			t.Fatal(err)
		}
		out = append(out, AssistantDescriptor{ID: a.AssistantID, Name: string(a.AssistantID), SkillsDir: dir})
	}
	return out
}

func strPtr(s string) *string { return &s }

// TestAgentWriteModes_OnlyUpdateAndSkip pins the write semantics: there is no
// append mode, so a re-install can never stack a second block.
func TestAgentWriteModes_OnlyUpdateAndSkip(t *testing.T) {
	if AgentModeUpdate != 0 || AgentModeSkip != 1 {
		t.Fatalf("AgentModeUpdate=%d AgentModeSkip=%d, want 0 and 1", AgentModeUpdate, AgentModeSkip)
	}
	home := agentHome(t)
	assistants := testAgentAssistants(t, home)
	for range 3 {
		InstallAgentConfig(assistants, PersonaPresets[0], true, "en", map[string]AgentWriteMode{})
	}
	for _, a := range AgentConfigAdapters {
		if n := strings.Count(readAgentFile(t, a.ConfigPath()), asdtBlockStart); n != 1 {
			t.Errorf("%s holds %d blocks after three installs, want 1", a.AssistantID, n)
		}
	}
}

// TestAgentConfigWrites_FollowSymlink: a symlinked config file stays a symlink
// through install and remove; the file it points to gets the content.
func TestAgentConfigWrites_FollowSymlink(t *testing.T) {
	home := agentHome(t)
	assistants := testAgentAssistants(t, home)
	for _, a := range AgentConfigAdapters {
		real := filepath.Join(home, "dotfiles", string(a.AssistantID)+".md")
		seedAgentFile(t, real, "# Mine\n")
		link := a.ConfigPath()
		if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(real, link); err != nil {
			t.Fatal(err)
		}
	}

	assertLink := func(stage string) {
		t.Helper()
		for _, a := range AgentConfigAdapters {
			info, err := os.Lstat(a.ConfigPath())
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode()&os.ModeSymlink == 0 {
				t.Errorf("%s after %s: config is no longer a symlink", a.AssistantID, stage)
			}
		}
	}

	for _, r := range InstallAgentConfig(assistants, PersonaPresets[0], true, "en", nil) {
		if r.Err != nil || r.Skipped {
			t.Fatalf("install %s: %+v", r.AssistantID, r)
		}
	}
	assertLink("install")
	for _, a := range AgentConfigAdapters {
		real := filepath.Join(home, "dotfiles", string(a.AssistantID)+".md")
		if got := readAgentFile(t, real); !strings.HasPrefix(got, "# Mine\n\n"+asdtBlockStart+"\n") {
			t.Errorf("%s link target lacks the block: %q", a.AssistantID, got)
		}
	}

	for _, r := range RemoveAgentConfig(assistants) {
		if r.Err != nil || r.Skipped {
			t.Fatalf("remove %s: %+v", r.AssistantID, r)
		}
	}
	assertLink("remove")
	for _, a := range AgentConfigAdapters {
		real := filepath.Join(home, "dotfiles", string(a.AssistantID)+".md")
		if got := readAgentFile(t, real); got != "# Mine\n" {
			t.Errorf("%s link target after remove: %q", a.AssistantID, got)
		}
	}
}

// TestAgentConfigWrites_PreserveMode: install and remove keep an existing
// file's mode; a new file is created 0o644.
func TestAgentConfigWrites_PreserveMode(t *testing.T) {
	home := agentHome(t)
	assistants := testAgentAssistants(t, home)
	claude, _ := AgentConfigAdapterFor(AssistantClaudeCode)
	seedAgentFile(t, claude.ConfigPath(), "# Mine\n")
	if err := os.Chmod(claude.ConfigPath(), 0o600); err != nil {
		t.Fatal(err)
	}
	openCode, _ := AgentConfigAdapterFor(AssistantOpenCode) // no prior file

	assertMode := func(stage string, a AgentConfigAdapterDescriptor, want os.FileMode) {
		t.Helper()
		info, err := os.Stat(a.ConfigPath())
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != want {
			t.Errorf("%s after %s: mode %o, want %o", a.AssistantID, stage, got, want)
		}
	}

	InstallAgentConfig(assistants, PersonaPresets[0], true, "en", nil)
	assertMode("install", claude, 0o600)
	assertMode("install", openCode, 0o644)
	RemoveAgentConfig(assistants)
	assertMode("remove", claude, 0o600)
	assertMode("remove", openCode, 0o644)

	entries, err := os.ReadDir(filepath.Dir(claude.ConfigPath()))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".asdt-tmp-") {
			t.Errorf("temp file left behind: %s", e.Name())
		}
	}
}
