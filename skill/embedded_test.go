package skill

import (
	"io/fs"
	"os"
	"strings"
	"testing"
)

// TestEmbeddedSkillTree asserts that every asdt-* directory on disk under
// skill/ is present in the embedded FS with its SKILL.md, and that the
// top-level SKILL.md is embedded. This guards against the embed directive
// silently missing a specialist directory (e.g. a newly added asdt-<name>/).
func TestEmbeddedSkillTree(t *testing.T) {
	if _, err := fs.Stat(FS(), "SKILL.md"); err != nil {
		t.Errorf("top-level SKILL.md missing from embedded FS: %v", err)
	}

	// The installer bakes this header into generated agent definitions
	// (internal/installer/agent_adapters.go) and silently generates nothing
	// when it is absent — so its presence here is load-bearing.
	if _, err := fs.Stat(FS(), "asdt-shared/skills/executor-header.md"); err != nil {
		t.Errorf("asdt-shared/skills/executor-header.md missing from embedded FS: %v", err)
	}

	diskEntries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("reading skill/ on disk: %v", err)
	}

	for _, entry := range diskEntries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "asdt-") {
			continue
		}
		name := entry.Name()

		if _, err := fs.Stat(FS(), name); err != nil {
			t.Errorf("directory %s exists on disk but is missing from embedded FS: %v", name, err)
			continue
		}
		// asdt-shared has no SKILL.md; only require it when it exists on disk.
		if _, err := os.Stat(name + "/SKILL.md"); err == nil {
			if _, err := fs.Stat(FS(), name+"/SKILL.md"); err != nil {
				t.Errorf("%s/SKILL.md missing from embedded FS: %v", name, err)
			}
		}
	}
}

// TestEmbeddedSharedSkillsMatchDisk asserts that the embedded FS view of
// skill/asdt-shared/skills/ matches disk in BOTH directions: every regular
// fragment on disk ships, and every underscore- or dot-prefixed entry does not.
//
// This is the only guard against a genuinely silent failure mode: an embed glob
// that is wrong produces no error at all — a fragment simply stops shipping (so
// the installer's header fold fails at install time, or worse, ships short), or
// an excluded draft file suddenly starts shipping. Neither shows up as a build
// or install error, only as a missing/extra file inside the embedded tree.
func TestEmbeddedSharedSkillsMatchDisk(t *testing.T) {
	const sharedSkillsDir = "asdt-shared/skills"

	diskEntries, err := os.ReadDir(sharedSkillsDir)
	if err != nil {
		t.Fatalf("reading %s on disk: %v", sharedSkillsDir, err)
	}

	for _, entry := range diskEntries {
		name := entry.Name()
		embeddedPath := sharedSkillsDir + "/" + name
		excluded := strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".")

		_, statErr := fs.Stat(FS(), embeddedPath)
		if excluded {
			if statErr == nil {
				t.Errorf("%s starts with %q or %q and must NOT be embedded, but it is present in the embedded FS", embeddedPath, "_", ".")
			}
			continue
		}
		if statErr != nil {
			t.Errorf("%s exists on disk but is missing from embedded FS: %v", embeddedPath, statErr)
		}
	}
}

// TestRoutedSpecialistInvariants asserts that every routed specialist's
// embedded SKILL.md body carries the load-bearing header invariants: the SOLE
// orchestrator gate, the FIRST ACTION self-load instruction, and both
// specialist-header generated-region markers. These phrases drive orchestration
// behavior, so dropping them silently would break the launch contract — this
// guards against that.
//
// The marker assertions are the ONLY guard against a real silent failure:
// replaceMarkerRegion returns content unchanged when BOTH markers are absent,
// so a specialist whose region was dropped would install with no header at all
// and no error at install time.
func TestRoutedSpecialistInvariants(t *testing.T) {
	routed := []string{
		"asdt-pm",
		"asdt-ux-ui",
		"asdt-architect",
		"asdt-developer",
		"asdt-qa",
		"asdt-security",
		"asdt-researcher",
	}

	const (
		soleOrchestrator = "SOLE orchestrator"
		firstAction      = "FIRST ACTION — self-load the header"

		// Package skill cannot import internal/installer (the installer imports
		// this package), so these two literals are a hand-copy of
		// specialistHeaderBeginMarker / specialistHeaderEndMarker in
		// internal/installer/registry_gen.go. Change both copies together.
		headerBeginMarker = "<!-- ASDT:GENERATED:specialist-header -->"
		headerEndMarker   = "<!-- /ASDT:GENERATED:specialist-header -->"
	)

	for _, dir := range routed {
		data, err := fs.ReadFile(FS(), dir+"/SKILL.md")
		if err != nil {
			t.Errorf("%s/SKILL.md missing from embedded FS: %v", dir, err)
			continue
		}
		body := string(data)
		if !strings.Contains(body, soleOrchestrator) {
			t.Errorf("%s/SKILL.md missing required phrase %q", dir, soleOrchestrator)
		}
		if !strings.Contains(body, firstAction) {
			t.Errorf("%s/SKILL.md missing required phrase %q", dir, firstAction)
		}
		if !strings.Contains(body, headerBeginMarker) {
			t.Errorf("%s/SKILL.md missing required phrase %q", dir, headerBeginMarker)
		}
		if !strings.Contains(body, headerEndMarker) {
			t.Errorf("%s/SKILL.md missing required phrase %q", dir, headerEndMarker)
		}
	}
}
