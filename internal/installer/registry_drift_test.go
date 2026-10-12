package installer

import (
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
}

// asdtCommandRe matches an /asdt-<name> slash command token.
var asdtCommandRe = regexp.MustCompile(`/asdt-[a-z-]+`)

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
// sitting between memory-check and clarify, explore and write remain the only
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

	wantInline := []string{"memory-check", "enrichment", "clarify"}
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

	want := []string{"memory-check", "explore", "enrichment", "clarify", "write"}
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

// TestHumanNotesRoutedToMemory defends where a person's words about the
// project go: to memory, through the orchestrator's remember, never into a
// file. write.md is the only step that writes `.asdt/`, so it must never route
// an enrichment `nuance.*` answer into config.yaml; SKILL.md must remember
// each answer instead; and platform-context.md's injection path (Reuse Guard,
// Injection Format) must never pull a note in among the detected facts, where
// it would read as something the scan found.
func TestHumanNotesRoutedToMemory(t *testing.T) {
	root := skillDir(t)

	writeMD := readText(t, filepath.Join(root, "asdt-init", "steps", "write.md"))
	processing, ok := markdownSection(writeMD, "Processing")
	if !ok {
		t.Fatalf("asdt-init/steps/write.md has no ## Processing section")
	}
	for _, forbidden := range []string{"nuance.", "Notes to remember"} {
		if strings.Contains(processing, forbidden) {
			t.Errorf("write.md Processing references %q; human notes never reach a file", forbidden)
		}
	}
	if !strings.Contains(writeMD, "Ignore `## Notes to remember`") {
		t.Errorf("write.md must tell the writer to ignore the clarify payload's `## Notes to remember`")
	}

	skillMD := readText(t, filepath.Join(root, "asdt-init", "SKILL.md"))
	after, ok := markdownSection(skillMD, "After write")
	if !ok {
		t.Fatalf("asdt-init/SKILL.md has no ## After write section")
	}
	if !strings.Contains(after, "Notes to remember") || !strings.Contains(after, "memory **remember**") {
		t.Errorf("asdt-init/SKILL.md ## After write must remember each `Notes to remember` entry through memory **remember**")
	}

	pc := readText(t, filepath.Join(root, "asdt-core", "references", "platform-context.md"))
	for _, section := range []string{"Reuse Guard", "Injection Format"} {
		body, ok := markdownSection(pc, section)
		if !ok {
			t.Fatalf("platform-context.md has no ## %s section", section)
		}
		for _, forbidden := range []string{"PROJECT NOTE", "Notes to remember", "memory **recall**"} {
			if strings.Contains(body, forbidden) {
				t.Errorf("platform-context.md ## %s references %q; notes never join the detected-facts block", section, forbidden)
			}
		}
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
			if strings.Contains(line, "| Step | Execution | Does |") {
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
