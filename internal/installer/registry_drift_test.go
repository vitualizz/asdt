package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestRegistryDrift is the structural defense against the §9.2 "Parity check"
// prose drifting out of sync (which it already did). The canonical roster is
// DERIVED from the 8 workflow.yaml via parseRegistry; each of the 3 mirror
// sites is parsed for its specialist rows, normalized to canonical directory
// names, and asserted against the roster the site is expected to carry. It also
// asserts the §9.2 inline-steps list is byte-equal to renderInlineStepsRegion's
// output, so the install-time generator and the committed file agree.
//
// The guard asserts SET membership + routable booleans + the inline-steps list
// ONLY — never the Discipline / When-to-involve / Trivial-eligible cell prose.
func TestRegistryDrift(t *testing.T) {
	root := skillDir(t)

	regs, err := parseRegistry(os.DirFS(root))
	if err != nil {
		t.Fatalf("parseRegistry: %v", err)
	}

	regByDir := make(map[string]specialistRegistry, len(regs))
	for _, r := range regs {
		regByDir[r.Dir] = r
	}

	// Canonical roster = every routable specialist. Non-routable specialists
	// (asdt-init) appear in NO mirror site.
	routable := map[string]bool{}
	for _, r := range regs {
		if r.Routable {
			routable[r.Dir] = true
		}
	}

	// Per-site expected membership = routable roster MINUS the deliberate
	// per-site omissions. agents-template.md has no Product Manager row; all
	// sites omit asdt-init (non-routable / setup-class). The row list is never
	// hardcoded — it is derived, then the explicit policy is subtracted.
	sites := []struct {
		name        string
		path        string
		excludeDirs map[string]bool
		rows        func(string) []string
	}{
		{
			name:        "§5 Specialist Registry",
			path:        filepath.Join(root, "SKILL.md"),
			excludeDirs: map[string]bool{},
			rows:        section5Rows,
		},
		{
			name:        "§9.2 trivial table",
			path:        filepath.Join(root, "SKILL.md"),
			excludeDirs: map[string]bool{},
			rows:        trivialTableRows,
		},
		{
			name: "agents-template",
			// agents-template.md lives in this package's assets/ dir; tests run
			// from internal/installer, so the path is package-relative.
			path:        filepath.Join("assets", "agents-template.md"),
			excludeDirs: map[string]bool{"asdt-pm": true},
			rows:        agentsTemplateRows,
		},
	}

	for _, site := range sites {
		data, err := os.ReadFile(site.path)
		if err != nil {
			t.Fatalf("%s: read %s: %v", site.name, site.path, err)
		}

		parsed := site.rows(string(data))
		seen := map[string]bool{}
		for _, dir := range parsed {
			seen[dir] = true
			// Phantom: a row whose specialist is not in the canonical roster,
			// or one the site is supposed to omit.
			if !routable[dir] {
				if _, known := regByDir[dir]; !known {
					t.Errorf("%s: phantom row %q has no canonical workflow.yaml", site.name, dir)
				} else {
					t.Errorf("%s: phantom row %q is non-routable and must not appear in any mirror site", site.name, dir)
				}
				continue
			}
			if site.excludeDirs[dir] {
				t.Errorf("%s: phantom row %q must NOT appear in this site (deliberate omission)", site.name, dir)
			}
		}

		// Missing: a routable specialist the site should carry but doesn't.
		for dir := range routable {
			if site.excludeDirs[dir] {
				continue
			}
			if !seen[dir] {
				t.Errorf("%s: missing routable specialist %q", site.name, dir)
			}
		}
	}

	// Inline-steps list parity: the committed §9.2 region must equal what the
	// generator would render, so install-time regeneration is a no-op.
	skillMD, err := os.ReadFile(filepath.Join(root, "SKILL.md"))
	if err != nil {
		t.Fatalf("read SKILL.md: %v", err)
	}
	wantBody := renderInlineStepsRegion(regs)
	gotBody, err := extractMarkerRegion(skillMD, inlineStepsBeginMarker, inlineStepsEndMarker)
	if err != nil {
		t.Fatalf("extract §9.2 inline-steps region: %v", err)
	}
	if gotBody != wantBody {
		t.Errorf("§9.2 inline-steps region drifted from workflow.yaml derivation:\n--- committed ---\n%q\n--- derived ---\n%q", gotBody, wantBody)
	}
}

// asdtCommandRe matches an /asdt-<name> slash command token.
var asdtCommandRe = regexp.MustCompile(`/asdt-[a-z-]+`)

// asdtPathRe matches an asdt-<name> directory token inside a skill/ path.
var asdtPathRe = regexp.MustCompile(`asdt-[a-z-]+`)

// firstCommandDir returns the canonical dir for the FIRST /asdt-<name> command
// in a table row (the Command column precedes any prose that references other
// specialists), or "" if the row carries no command.
func firstCommandDir(row string) string {
	if m := asdtCommandRe.FindString(row); m != "" {
		return strings.TrimPrefix(m, "/")
	}
	return ""
}

// section5Rows extracts the canonical dirs of the §5 Specialist Registry table.
func section5Rows(content string) []string {
	return commandRowsBetween(content, "## 5. Specialist Registry", "| Specialist | Command |")
}

// agentsTemplateRows extracts the canonical dirs of the agents-template.md
// ASDT Specialists table.
func agentsTemplateRows(content string) []string {
	return commandRowsBetween(content, "ASDT Specialists", "| Specialist")
}

// commandRowsBetween walks the table that begins after the first header line
// containing headerHint (itself appearing after sectionHint) and returns the
// canonical dir of each data row, identified by its first /asdt-<name> command.
func commandRowsBetween(content, sectionHint, headerHint string) []string {
	lines := strings.Split(content, "\n")
	var dirs []string
	inSection := sectionHint == ""
	inTable := false
	for _, line := range lines {
		if !inSection {
			if strings.Contains(line, sectionHint) {
				inSection = true
			}
			continue
		}
		if !inTable {
			if strings.Contains(line, headerHint) {
				inTable = true
			}
			continue
		}
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			if trimmed == "" {
				break // blank line ends the table
			}
			continue
		}
		if strings.HasPrefix(trimmed, "|---") || strings.HasPrefix(trimmed, "|--") {
			continue // separator row
		}
		if dir := firstCommandDir(trimmed); dir != "" {
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

// trivialTableRows extracts the canonical dirs of the §9.2 per-specialist
// trivial-step table, whose rows reference `skill/asdt-<name>/SKILL.md` paths
// (no slash command).
func trivialTableRows(content string) []string {
	lines := strings.Split(content, "\n")
	var dirs []string
	inTable := false
	for _, line := range lines {
		if !inTable {
			if strings.Contains(line, "| Specialist | File | Trivial step | Trivial eligible? |") {
				inTable = true
			}
			continue
		}
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			if trimmed == "" {
				break
			}
			continue
		}
		if strings.HasPrefix(trimmed, "|---") || strings.HasPrefix(trimmed, "|--") {
			continue
		}
		if m := asdtPathRe.FindString(trimmed); m != "" {
			dirs = append(dirs, m)
		}
	}
	return dirs
}

// extractMarkerRegion returns the bytes strictly between beginMarker and
// endMarker as a string, erroring when either marker is absent or the end
// precedes the begin.
func extractMarkerRegion(content []byte, beginMarker, endMarker string) (string, error) {
	s := string(content)
	begin := strings.Index(s, beginMarker)
	if begin < 0 {
		return "", fmt.Errorf("begin marker %q not found", beginMarker)
	}
	end := strings.Index(s, endMarker)
	if end < 0 {
		return "", fmt.Errorf("end marker %q not found", endMarker)
	}
	regionStart := begin + len(beginMarker)
	if end < regionStart {
		return "", fmt.Errorf("end marker %q precedes begin marker %q", endMarker, beginMarker)
	}
	return s[regionStart:end], nil
}
