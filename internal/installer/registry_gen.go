package installer

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Inline-steps generated-region markers. These literal HTML comments bound the
// machine-generated `Tailored Workflow Generation` inline-steps list in the
// root SKILL.md. They are the single source both the install-time generator
// and any test reference, so the region can be located and replaced
// idempotently.
const (
	inlineStepsBeginMarker = "<!-- ASDT:GENERATED:9.2-inline-steps -->"
	inlineStepsEndMarker   = "<!-- /ASDT:GENERATED:9.2-inline-steps -->"
)

// Specialist-header generated-region markers. These literal HTML comments bound
// the machine-spliced shared specialist header in each routed specialist's
// SKILL.md. Same contract as the inline-steps pair: one literal source for the
// install-time generator and for any test, so the region is located and
// replaced idempotently. skill/embedded_test.go keeps a hand-copy of these two
// literals (package skill cannot import internal/installer) — change both
// copies together.
const (
	specialistHeaderBeginMarker = "<!-- ASDT:GENERATED:specialist-header -->"
	specialistHeaderEndMarker   = "<!-- /ASDT:GENERATED:specialist-header -->"
)

// registryRenderOrder is the display order of specialists in the `Tailored
// Workflow Generation` inline-steps list. It is intentionally NOT the
// alphabetical directory walk order: the hand-authored list reads
// PM → Architect → QA → Security → UX/UI → Developer → Researcher, and the
// generator must reproduce that exact sequence so its first run against the
// existing SKILL.md is a byte-identical no-op.
var registryRenderOrder = []string{
	"asdt-pm",
	"asdt-architect",
	"asdt-qa",
	"asdt-security",
	"asdt-ux-ui",
	"asdt-developer",
	"asdt-researcher",
}

// specialistHeaderFragments lists the shared fragments folded, in this exact
// sequence, into every routed specialist's specialist-header region.
//
// ORDER IS LOAD-BEARING — it is the order a specialist reads the file top-down:
//   - specialist-header.md first: the Prerequisites gate and the ORCHESTRATOR
//     GATE must be the very first thing a specialist reads, before it can act
//     on anything else in the document.
//   - parallel-retrieval.md second: the Cache Ledger Rule and the Injection
//     Format are needed at launch time — the orchestrator must already hold
//     them before it runs the first inline step.
//   - intake-contract.md third: the declared-vs-present input check, the single
//     batched clarification turn, and the harden-always ASSUMED: degradation
//     govern how every step treats its inputs, so they must be settled before
//     any step content arrives.
//   - knowledge-recall.md last: it is the content of that first inline step, so
//     it only has to be in context once the launch contract above is settled.
var specialistHeaderFragments = []string{
	"asdt-shared/skills/specialist-header.md",
	"asdt-shared/skills/parallel-retrieval.md",
	"asdt-shared/skills/intake-contract.md",
	"asdt-shared/skills/knowledge-recall.md",
}

// inlineStepsDisplayNames maps a specialist directory to the display label
// used in the `Tailored Workflow Generation` inline-steps list
// (e.g. "asdt-pm" → "PM").
var inlineStepsDisplayNames = map[string]string{
	"asdt-pm":         "PM",
	"asdt-architect":  "Architect",
	"asdt-qa":         "QA",
	"asdt-security":   "Security",
	"asdt-ux-ui":      "UX/UI",
	"asdt-developer":  "Developer",
	"asdt-researcher": "Researcher",
}

// specialistRegistry is the canonical, derived view of one specialist's
// workflow.yaml. It is the single shape both the install-time generator and the
// CI-time drift guard build on.
type specialistRegistry struct {
	Dir           string   // directory name, e.g. "asdt-developer" (the walk key)
	Specialist    string   // the literal `specialist:` value
	Routable      bool     // the top-level `routable:` key
	SubagentSteps []string // ordered names of every `execution: subagent` step
	InlineSteps   []string // ordered names of every `execution: inline` step
}

// registryStep is the minimal step shape needed to classify steps by execution.
type registryStep struct {
	Name      string `yaml:"name"`
	Execution string `yaml:"execution"`
}

// registryFile is the minimal workflow.yaml shape the registry walk parses.
type registryFile struct {
	Specialist string         `yaml:"specialist"`
	Routable   bool           `yaml:"routable"`
	Steps      []registryStep `yaml:"steps"`
}

// parseRegistry walks skillsFS for */workflow.yaml files and returns the
// canonical registry for each specialist, ordered by directory name. It reuses
// the WorkflowModelSteps walk idiom (fs.ReadDir on ".", sorted dir names,
// fs.ReadFile dir/workflow.yaml, skip dirs without one). Steps are classified
// into SubagentSteps / InlineSteps by their `execution:` field. The fs.FS
// parameter lets the same function serve the embedded skillsFS at install time
// and an os.DirFS(skillDir) tree under test.
func parseRegistry(skillsFS fs.FS) ([]specialistRegistry, error) {
	entries, err := fs.ReadDir(skillsFS, ".")
	if err != nil {
		return nil, fmt.Errorf("read skill root: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	var regs []specialistRegistry
	for _, dir := range names {
		wfPath := path.Join(dir, "workflow.yaml")
		data, readErr := fs.ReadFile(skillsFS, wfPath)
		if readErr != nil {
			continue // directory without workflow.yaml (e.g. asdt-shared)
		}

		var wf registryFile
		if err := yaml.Unmarshal(data, &wf); err != nil {
			return nil, fmt.Errorf("parse %s: %w", wfPath, err)
		}

		reg := specialistRegistry{
			Dir:        dir,
			Specialist: wf.Specialist,
			Routable:   wf.Routable,
		}
		for _, s := range wf.Steps {
			switch s.Execution {
			case "subagent":
				reg.SubagentSteps = append(reg.SubagentSteps, s.Name)
			case "inline":
				reg.InlineSteps = append(reg.InlineSteps, s.Name)
			}
		}
		regs = append(regs, reg)
	}

	return regs, nil
}

// renderInlineStepsRegion emits the `Tailored Workflow Generation`
// inline-steps bullet block body — the bytes that sit strictly between the
// markers — from the registry. It produces one "- {Display}: `step`, `step`"
// line per specialist that has inline steps, in registryRenderOrder. The
// output is leading-newline / trailing-newline padded so it slots cleanly
// between markers that each occupy their own line.
// Pure function; no I/O.
func renderInlineStepsRegion(regs []specialistRegistry) string {
	byDir := make(map[string]specialistRegistry, len(regs))
	for _, r := range regs {
		byDir[r.Dir] = r
	}

	var lines []string
	for _, dir := range registryRenderOrder {
		reg, ok := byDir[dir]
		if !ok || len(reg.InlineSteps) == 0 {
			continue
		}
		display := inlineStepsDisplayNames[dir]
		quoted := make([]string, len(reg.InlineSteps))
		for i, step := range reg.InlineSteps {
			quoted[i] = "`" + step + "`"
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", display, strings.Join(quoted, ", ")))
	}

	// The region body is surrounded by newlines so the begin marker, the list,
	// and the end marker each occupy their own line after splicing.
	return "\n" + strings.Join(lines, "\n") + "\n"
}

// replaceMarkerRegion splices newBody between beginMarker and endMarker in
// content, replacing whatever currently sits strictly between them. It is
// idempotent: when the existing region already equals newBody the original
// content is returned unchanged (mirroring InjectModels' no-op-when-unchanged
// contract).
//
// When NEITHER marker is present the content has no generated region to
// maintain and is returned unchanged — this is the only silent skip, and it
// keeps marker-free documents (e.g. minimal test fixtures, per-specialist
// SKILL.md) from being treated as corrupt. Any PARTIAL or malformed marker
// state — one marker present without the other, a marker appearing more than
// once, or the end marker preceding the begin marker — fails loudly with an
// error and never half-writes, so a damaged generated region in the real root
// SKILL.md is caught rather than silently skipped.
func replaceMarkerRegion(content []byte, beginMarker, endMarker, newBody string) ([]byte, error) {
	begin := []byte(beginMarker)
	end := []byte(endMarker)

	beginCount := bytes.Count(content, begin)
	endCount := bytes.Count(content, end)
	if beginCount == 0 && endCount == 0 {
		return content, nil // no generated region present — nothing to maintain
	}
	if beginCount != 1 {
		return nil, fmt.Errorf("begin marker %q appears %d times, want exactly 1", beginMarker, beginCount)
	}
	if endCount != 1 {
		return nil, fmt.Errorf("end marker %q appears %d times, want exactly 1", endMarker, endCount)
	}

	beginIdx := bytes.Index(content, begin)
	endIdx := bytes.Index(content, end)
	regionStart := beginIdx + len(begin)
	if endIdx < regionStart {
		return nil, fmt.Errorf("end marker %q precedes begin marker %q", endMarker, beginMarker)
	}

	if string(content[regionStart:endIdx]) == newBody {
		return content, nil // region already current — no-op
	}

	var buf bytes.Buffer
	buf.Grow(len(content) - (endIdx - regionStart) + len(newBody))
	buf.Write(content[:regionStart])
	buf.WriteString(newBody)
	buf.Write(content[endIdx:])
	return buf.Bytes(), nil
}

// GenerateInlineSteps returns skillMD with its `Tailored Workflow Generation`
// inline-steps marker region regenerated from the specialist workflow.yaml
// files in skillsFS. It is the install-time entry point: parseRegistry →
// renderInlineStepsRegion → replaceMarkerRegion. The returned bytes are
// unchanged when the region already matches the derived list.
func GenerateInlineSteps(skillsFS fs.FS, skillMD []byte) ([]byte, error) {
	regs, err := parseRegistry(skillsFS)
	if err != nil {
		return nil, fmt.Errorf("parse specialist registry: %w", err)
	}

	body := renderInlineStepsRegion(regs)
	updated, err := replaceMarkerRegion(skillMD, inlineStepsBeginMarker, inlineStepsEndMarker, body)
	if err != nil {
		return nil, fmt.Errorf("regenerate inline-steps region: %w", err)
	}
	return updated, nil
}

// GenerateSpecialistHeader returns skillMD with its specialist-header marker
// region replaced by the shared header fragments from skillsFS, folded together
// in specialistHeaderFragments order. It exists so the header is already
// spliced into each routed specialist's installed SKILL.md: the orchestrator
// then has the header in context the moment it reads the specialist, instead of
// depending on separate reads of the asdt-shared/skills/*.md fragments that may
// never happen. Files without the markers pass through unchanged
// (replaceMarkerRegion's no-region skip), so this is safe to run over every
// SKILL.md the installer writes.
//
// The marker presence check runs BEFORE the fragment reads on purpose: a
// document with no markers declares no region, so requiring the shared
// fragments for it would turn a marker-free tree (a minimal test fixture FS,
// asdt-init) into an install failure — that early return is why marker-free
// trees pass through. A file carrying even one marker needs the fragments, so
// ANY read failure fails the whole call: a partial header is worse than a loud
// install failure. A partial/duplicated marker state still fails in
// replaceMarkerRegion.
func GenerateSpecialistHeader(skillsFS fs.FS, skillMD []byte) ([]byte, error) {
	if !bytes.Contains(skillMD, []byte(specialistHeaderBeginMarker)) &&
		!bytes.Contains(skillMD, []byte(specialistHeaderEndMarker)) {
		return skillMD, nil
	}

	parts := make([]string, 0, len(specialistHeaderFragments))
	for _, p := range specialistHeaderFragments {
		fragment, err := fs.ReadFile(skillsFS, p)
		if err != nil {
			return nil, fmt.Errorf("read specialist header fragment %s: %w", p, err)
		}
		parts = append(parts, strings.TrimRight(string(fragment), "\n"))
	}

	// Same padding convention as renderInlineStepsRegion: leading and trailing
	// newlines so the begin marker, the header, and the end marker each occupy
	// their own line after splicing. Fragments are joined by a blank line so
	// each boundary lands on its own line too.
	body := "\n" + strings.Join(parts, "\n\n") + "\n"
	updated, err := replaceMarkerRegion(skillMD, specialistHeaderBeginMarker, specialistHeaderEndMarker, body)
	if err != nil {
		return nil, fmt.Errorf("regenerate specialist-header region: %w", err)
	}
	return updated, nil
}
