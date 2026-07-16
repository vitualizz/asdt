package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vitualizz/asdt/skill"
)

// TestInstall_DoesNotShipInstallerAssets asserts that a real install of the
// production skill tree never writes the relocated installer-only assets
// (the persona markdown + agents-template.md) under the assistant's
// .claude/skills/asdt-init/ directory, while the asdt-init SKILL.md — genuine
// conversational skill content — IS still installed. After ADR #2115 the
// assets physically left skill/, so this is structurally guaranteed; this test
// is the regression guard that proves it.
func TestInstall_DoesNotShipInstallerAssets(t *testing.T) {
	skillsDir := filepath.Join(t.TempDir(), "skills")
	a := AssistantDescriptor{ID: "copywalk-test", Name: "Copywalk Test", BinaryName: "sh", SkillsDir: skillsDir}

	results := InstallWithModels([]AssistantDescriptor{a}, Providers[0], skill.FS(), "", nil, InstallOptions{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("install error: %v", results[0].Err)
	}

	initDir := filepath.Join(skillsDir, "asdt-init")

	// SKILL.md is conversational skill content — it MUST be installed.
	if _, err := os.Stat(filepath.Join(initDir, "SKILL.md")); err != nil {
		t.Fatalf("expected asdt-init/SKILL.md to be installed, got: %v", err)
	}

	// The relocated installer assets MUST NOT ship under .claude/skills/.
	shouldBeAbsent := []string{
		filepath.Join(initDir, "agents-template.md"),
		filepath.Join(initDir, "personas"),
		filepath.Join(initDir, "personas", "sky.md"),
	}
	for _, p := range shouldBeAbsent {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("expected %q to be absent after install, but os.Stat err = %v", p, err)
		}
	}
}
