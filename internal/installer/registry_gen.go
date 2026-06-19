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
// machine-generated §9.2 inline-steps list in the root SKILL.md. They are the
// single source both the install-time generator and any test reference, so the
// region can be located and replaced idempotently.
const (
	inlineStepsBeginMarker = "<!-- ASDT:GENERATED:9.2-inline-steps -->"
	inlineStepsEndMarker   = "<!-- /ASDT:GENERATED:9.2-inline-steps -->"
)

// registryRenderOrder is the display order of specialists in the §9.2
// inline-steps list. It is intentionally NOT the alphabetical directory walk
// order: the hand-authored list reads PM → Architect → QA → Security → UX/UI →
// Developer → Researcher, and the generator must reproduce that exact sequence
// so its first run against the existing SKILL.md is a byte-identical no-op.
var registryRenderOrder = []string{
	"asdt-pm",
	"asdt-architect",
	"asdt-qa",
	"asdt-security",
	"asdt-ux-ui",
	"asdt-developer",
	"asdt-researcher",
}

// inlineStepsDisplayNames maps a specialist directory to the display label used
// in the §9.2 inline-steps list (e.g. "asdt-pm" → "PM").
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

// renderInlineStepsRegion emits the §9.2 inline-steps bullet block body — the
// bytes that sit strictly between the markers — from the registry. It produces
// one "- {Display}: `step`, `step`" line per specialist that has inline steps,
// in registryRenderOrder. The output is leading-newline / trailing-newline
// padded so it slots cleanly between markers that each occupy their own line.
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

// GenerateInlineSteps returns skillMD with its §9.2 inline-steps marker region
// regenerated from the specialist workflow.yaml files in skillsFS. It is the
// install-time entry point: parseRegistry → renderInlineStepsRegion →
// replaceMarkerRegion. The returned bytes are unchanged when the region already
// matches the derived list.
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
