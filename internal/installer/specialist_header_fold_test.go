package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vitualizz/asdt/skill"
)

// TestInstall_SpecialistHeaderFoldReachesInstalledSKILL asserts that the folded
// specialist header actually lands in the SKILL.md a user gets on disk, for
// every routed specialist. Reading the source tree proves nothing here: the
// fragments are spliced at install time, so only the INSTALLED artifact shows
// whether the fold ran, ran in the right order, and folded the right set.
//
// The sentinels below are prose lifted from the fragments, so rewording a
// fragment heading breaks this test. That is the intended tradeoff: a loud test
// failure that points at the reword beats a silent fold regression that ships a
// specialist with half a header and no error anywhere.
func TestInstall_SpecialistHeaderFoldReachesInstalledSKILL(t *testing.T) {
	skillsDir := filepath.Join(t.TempDir(), "skills")
	a := AssistantDescriptor{ID: "header-fold-test", Name: "Header Fold Test", BinaryName: "sh", SkillsDir: skillsDir}

	results := InstallWithModels([]AssistantDescriptor{a}, Providers[0], skill.FS(), "", nil, InstallOptions{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("install error: %v", results[0].Err)
	}

	// Ordered exactly like specialistHeaderFragments: the assertions below check
	// both presence and relative position, so this slice IS the expected order.
	folded := []struct {
		name     string
		sentinel string
	}{
		{name: "specialist-header", sentinel: "ORCHESTRATOR GATE"},
		{name: "protocol", sentinel: "# ASDT Protocol"},
	}

	// The expected-order slice above must cover the production fragment list
	// one-to-one: a fragment added to specialistHeaderFragments without a
	// matching sentinel here would fold silently, unasserted.
	if len(specialistHeaderFragments) != len(folded) {
		t.Fatalf("specialistHeaderFragments has %d entries but this test's folded slice has %d — keep registry_gen.go's specialistHeaderFragments and the folded expectations in lockstep", len(specialistHeaderFragments), len(folded))
	}

	// knowledge-recall is deliberately NOT folded — it is a per-specialist inline
	// step, not a header fragment. Its presence would mean the list grew by accident.
	const notFolded = "# Knowledge Recall — Shared Skill"

	tests := []struct {
		name string
		dir  string
	}{
		{name: "pm", dir: "asdt-pm"},
		{name: "architect", dir: "asdt-architect"},
		{name: "qa", dir: "asdt-qa"},
		{name: "security", dir: "asdt-security"},
		{name: "ux-ui", dir: "asdt-ux-ui"},
		{name: "developer", dir: "asdt-developer"},
		{name: "researcher", dir: "asdt-researcher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installed := filepath.Join(skillsDir, tt.dir, "SKILL.md")
			data, err := os.ReadFile(installed)
			if err != nil {
				t.Fatalf("reading installed %s: %v", installed, err)
			}
			body := string(data)

			prev := -1
			prevName := ""
			for _, f := range folded {
				idx := strings.Index(body, f.sentinel)
				if idx < 0 {
					t.Errorf("installed %s/SKILL.md missing %s sentinel %q — fragment not folded", tt.dir, f.name, f.sentinel)
					continue
				}
				if idx < prev {
					t.Errorf("installed %s/SKILL.md has %s at %d, before %s at %d — folded out of order", tt.dir, f.name, idx, prevName, prev)
				}
				prev = idx
				prevName = f.name
			}

			if strings.Contains(body, notFolded) {
				t.Errorf("installed %s/SKILL.md contains %q, but decision-preservation must NOT be folded into the header", tt.dir, notFolded)
			}
		})
	}
}
