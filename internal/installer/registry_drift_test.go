package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestRegistryDrift is the structural defense against the `Tailored Workflow
// Generation` "Parity check" prose drifting out of sync (which it already
// did). The canonical roster is DERIVED from the 8 workflow.yaml via
// parseRegistry; each of the 3 mirror sites is parsed for its specialist rows,
// normalized to canonical directory names, and asserted against the roster the
// site is expected to carry. It also asserts the `Tailored Workflow
// Generation` inline-steps list is byte-equal to renderInlineStepsRegion's
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
			name:        "Specialist Registry",
			path:        filepath.Join(root, "SKILL.md"),
			excludeDirs: map[string]bool{},
			rows:        registrySectionRows,
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

	// The inline-steps parity assertion was dropped with the router rewrite:
	// the `Tailored Workflow Generation` region it compared against no longer
	// exists in skill/SKILL.md. renderInlineStepsRegion still exists and is
	// still exercised by its own unit test; it simply has no committed mirror
	// left to drift from.
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

// registrySectionRows extracts the canonical dirs of the `## Registry` table.
func registrySectionRows(content string) []string {
	return commandRowsBetween(content, "## Registry", "| Specialist | Command |")
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

// TestInitEnrichmentInlineStep asserts the asdt-init workflow classifies its
// steps as the enrichment redesign requires: enrichment is a NEW inline step
// sitting between knowledge-gate and clarify, explore and write remain the only
// subagent steps, and because asdt-init is non-routable it never leaks into
// the `Tailored Workflow Generation` inline-steps render.
func TestInitEnrichmentInlineStep(t *testing.T) {
	root := skillDir(t)

	regs, err := parseRegistry(os.DirFS(root))
	if err != nil {
		t.Fatalf("parseRegistry: %v", err)
	}

	var initReg specialistRegistry
	found := false
	for _, r := range regs {
		if r.Dir == "asdt-init" {
			initReg = r
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("asdt-init not found in parsed registry")
	}

	wantInline := []string{"knowledge-gate", "enrichment", "clarify"}
	if !reflect.DeepEqual(initReg.InlineSteps, wantInline) {
		t.Errorf("asdt-init InlineSteps = %v, want %v", initReg.InlineSteps, wantInline)
	}

	wantSubagent := []string{"explore", "write"}
	if !reflect.DeepEqual(initReg.SubagentSteps, wantSubagent) {
		t.Errorf("asdt-init SubagentSteps = %v, want %v", initReg.SubagentSteps, wantSubagent)
	}

	// asdt-init is non-routable, so its inline steps must never render into
	// the `Tailored Workflow Generation` inline-steps region — enrichment in
	// particular must not leak.
	region := renderInlineStepsRegion(regs)
	if strings.Contains(region, "enrichment") {
		t.Errorf("Tailored Workflow Generation inline-steps region leaked non-routable asdt-init step %q:\n%s", "enrichment", region)
	}
}

// TestInitPlanTableMatchesWorkflow asserts the asdt-init SKILL.md Orchestration
// Plan table Step column is byte-identical, in order, to the workflow.yaml step
// sequence — so the hand-authored plan cannot drift from the executable flow.
func TestInitPlanTableMatchesWorkflow(t *testing.T) {
	root := skillDir(t)

	wfData, err := os.ReadFile(filepath.Join(root, "asdt-init", "workflow.yaml"))
	if err != nil {
		t.Fatalf("read asdt-init/workflow.yaml: %v", err)
	}
	var wf registryFile
	if err := yaml.Unmarshal(wfData, &wf); err != nil {
		t.Fatalf("parse asdt-init/workflow.yaml: %v", err)
	}
	var wfSteps []string
	for _, s := range wf.Steps {
		wfSteps = append(wfSteps, s.Name)
	}

	mdData, err := os.ReadFile(filepath.Join(root, "asdt-init", "SKILL.md"))
	if err != nil {
		t.Fatalf("read asdt-init/SKILL.md: %v", err)
	}
	tableSteps := initPlanTableSteps(string(mdData))

	want := []string{"knowledge-gate", "explore", "enrichment", "clarify", "write"}
	if !reflect.DeepEqual(wfSteps, want) {
		t.Errorf("workflow.yaml step order = %v, want %v", wfSteps, want)
	}
	if !reflect.DeepEqual(tableSteps, want) {
		t.Errorf("SKILL.md plan-table Step column = %v, want %v", tableSteps, want)
	}
	if !reflect.DeepEqual(wfSteps, tableSteps) {
		t.Errorf("plan-table drifted from workflow.yaml:\n--- workflow ---\n%v\n--- table ---\n%v", wfSteps, tableSteps)
	}
}

// TestNuanceIsolatedFromProvenance defends the deliberate separation of the
// human_nuance region: it belongs to knowledge.yaml (write step 2) ONLY,
// never to provenance.yaml (write step 3), and platform-context.md's
// auto-injection paths (Reuse Guard, What to Inject) must never pull it in — a
// user-authored note is read directly, never auto-injected.
func TestNuanceIsolatedFromProvenance(t *testing.T) {
	root := skillDir(t)

	writeMD, err := os.ReadFile(filepath.Join(root, "asdt-init", "steps", "write.md"))
	if err != nil {
		t.Fatalf("read asdt-init/steps/write.md: %v", err)
	}
	content := string(writeMD)

	step2 := lineSliceBetween(content, "**`.asdt/knowledge/knowledge.yaml`**", "**`.asdt/knowledge/provenance.yaml`**")
	if step2 == "" {
		t.Fatalf("could not isolate write.md step-2 knowledge section")
	}
	if !strings.Contains(step2, "human_nuance") {
		t.Errorf("write.md step-2 (knowledge) must reference human_nuance")
	}

	step3 := lineSliceBetween(content, "**`.asdt/knowledge/provenance.yaml`**", "## Halt contract")
	if step3 == "" {
		t.Fatalf("could not isolate write.md step-3 provenance section")
	}
	for _, forbidden := range []string{"human_nuance", "NUANCE"} {
		if strings.Contains(step3, forbidden) {
			t.Errorf("write.md step-3 (provenance) must not reference %q", forbidden)
		}
	}

	pcMD, err := os.ReadFile(filepath.Join(root, "asdt-core", "references", "platform-context.md"))
	if err != nil {
		t.Fatalf("read asdt-core/references/platform-context.md: %v", err)
	}
	pc := string(pcMD)

	// Section anchors follow the post-refactor platform-context.md, which
	// collapsed "How to Find" + "What to Inject" + the injection format into a
	// single "## Injection Format" section. The invariant is unchanged:
	// human_nuance is user-authored and carries no confidence rating, so it
	// must never appear in the automatic injection path.
	reuseGuard := lineSliceBetween(pc, "## Reuse Guard", "## Injection Format")
	if reuseGuard == "" {
		t.Fatalf("could not isolate platform-context.md Reuse Guard section")
	}
	if strings.Contains(reuseGuard, "human_nuance") {
		t.Errorf("platform-context.md Reuse Guard must not reference human_nuance")
	}

	whatToInject := lineSliceBetween(pc, "## Injection Format", "## Degradation")
	if whatToInject == "" {
		t.Fatalf("could not isolate platform-context.md Injection Format section")
	}
	if strings.Contains(whatToInject, "human_nuance") {
		t.Errorf("platform-context.md Injection Format must not reference human_nuance")
	}
}

// initPlanTableSteps extracts the ordered Step-column values of the asdt-init
// SKILL.md Orchestration Plan table (the first column of each data row).
func initPlanTableSteps(content string) []string {
	lines := strings.Split(content, "\n")
	var steps []string
	inSection := false
	inTable := false
	for _, line := range lines {
		if !inSection {
			if strings.Contains(line, "## Orchestration Plan") {
				inSection = true
			}
			continue
		}
		if !inTable {
			if strings.Contains(line, "| Step | File | Execution | Reads | Writes |") {
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
		cells := strings.Split(trimmed, "|")
		if len(cells) < 2 {
			continue
		}
		if name := strings.TrimSpace(cells[1]); name != "" {
			steps = append(steps, name)
		}
	}
	return steps
}

// lineSliceBetween returns the lines of content from the first line containing
// startHint (inclusive) up to but excluding the first line containing endHint
// after it. When endHint is "" or is not found, the slice runs to EOF. Returns
// "" when startHint itself is not found.
func lineSliceBetween(content, startHint, endHint string) string {
	lines := strings.Split(content, "\n")
	start := -1
	for i, line := range lines {
		if strings.Contains(line, startHint) {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	end := len(lines)
	if endHint != "" {
		for i := start + 1; i < len(lines); i++ {
			if strings.Contains(lines[i], endHint) {
				end = i
				break
			}
		}
	}
	return strings.Join(lines[start:end], "\n")
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
