package installer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// orchestrationRowKind classifies one tier cell of a specialist's
// `## Orchestration Plan` table.
type orchestrationRowKind int

const (
	kindStepList orchestrationRowKind = iota
	kindNoRun
	kindUnrecognized
)

// orchestrationRow is one resolved row of a specialist's tier table: its tier
// label, the cell's classification, the raw cell text, and — for kindStepList
// rows only — the resolved, ordered step names.
type orchestrationRow struct {
	Dir   string
	Tier  string
	Kind  orchestrationRowKind
	Cell  string
	Steps []string
}

// inputRef is one entry of a step's `inputs:` list: the topic_key it declares
// and whether that entry's line carries a trailing `# optional` marker.
type inputRef struct {
	TopicKey string
	Optional bool
}

// stepInfo is the minimal shape of one workflow.yaml step needed to build the
// dependency graph: identity, execution mode, what it produces, and what it
// declares as inputs.
type stepInfo struct {
	Name           string
	Execution      string
	OutputTopicKey string
	Inputs         []inputRef
}

// specialistGraph is one specialist's workflow.yaml, parsed independently of
// parseRegistry so every input scalar keeps a usable yaml.Node.Line.
type specialistGraph struct {
	Dir        string
	Specialist string
	Steps      []stepInfo
}

func (g specialistGraph) allStepNames() map[string]bool {
	names := make(map[string]bool, len(g.Steps))
	for _, s := range g.Steps {
		names[s.Name] = true
	}
	return names
}

func (g specialistGraph) subagentNamesInOrder() []string {
	var names []string
	for _, s := range g.Steps {
		if s.Execution == "subagent" {
			names = append(names, s.Name)
		}
	}
	return names
}

func (g specialistGraph) byName(name string) (stepInfo, bool) {
	for _, s := range g.Steps {
		if s.Name == name {
			return s, true
		}
	}
	return stepInfo{}, false
}

func (g specialistGraph) producerOf(topicKey string) (stepInfo, bool) {
	for _, s := range g.Steps {
		if s.OutputTopicKey == topicKey {
			return s, true
		}
	}
	return stepInfo{}, false
}

// parseSpecialistGraph parses a workflow.yaml into a specialistGraph. It
// reads the document into a yaml.Node (never a struct) so every input
// scalar's Node.Line addresses the raw source line it came from, and it
// splits the source once up front so lineHasOptionalMarker can index into it
// directly.
func parseSpecialistGraph(dir string, data []byte) (specialistGraph, error) {
	lines := strings.Split(string(data), "\n")

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return specialistGraph{}, fmt.Errorf("parse workflow.yaml: %w", err)
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return specialistGraph{}, fmt.Errorf("empty workflow.yaml document")
	}
	doc := root.Content[0]

	graph := specialistGraph{Dir: dir}
	if s := mappingValue(doc, "specialist"); s != nil {
		graph.Specialist = s.Value
	}

	stepsNode := mappingNode(doc, "steps")
	if stepsNode == nil || stepsNode.Kind != yaml.SequenceNode {
		return specialistGraph{}, fmt.Errorf("workflow.yaml has no steps sequence")
	}

	for _, stepNode := range stepsNode.Content {
		if stepNode.Kind != yaml.MappingNode {
			continue
		}
		info := stepInfo{}
		if n := mappingValue(stepNode, "name"); n != nil {
			info.Name = n.Value
		}
		if n := mappingValue(stepNode, "execution"); n != nil {
			info.Execution = n.Value
		}
		if n := mappingValue(stepNode, "output_topic_key"); n != nil {
			info.OutputTopicKey = n.Value
		}
		if inputsNode := mappingValue(stepNode, "inputs"); inputsNode != nil && inputsNode.Kind == yaml.SequenceNode {
			for _, item := range inputsNode.Content {
				if item.Kind != yaml.ScalarNode {
					continue
				}
				info.Inputs = append(info.Inputs, inputRef{
					TopicKey: item.Value,
					Optional: lineHasOptionalMarker(lines, item.Line),
				})
			}
		}
		graph.Steps = append(graph.Steps, info)
	}

	return graph, nil
}

// lineHasOptionalMarker reports whether the raw source line addressed by an
// input scalar's Node.Line (1-based) carries a trailing `# optional` marker.
// Topic keys never contain `#`, so the first `#` on the line is always the
// start of a comment, never inside the quoted string.
func lineHasOptionalMarker(lines []string, nodeLine int) bool {
	idx := nodeLine - 1
	if idx < 0 || idx >= len(lines) {
		return false
	}
	line := lines[idx]
	hashIdx := strings.Index(line, "#")
	if hashIdx < 0 {
		return false
	}
	return strings.HasPrefix(line[hashIdx:], "# optional")
}

// trailingParenRe strips a trailing parenthetical from a resolved step name,
// e.g. developer's `test (if TDD)` -> `test`.
var trailingParenRe = regexp.MustCompile(`\s*\([^)]*\)\s*$`)

// sameStepsAsRe matches the `Full workflow (same steps as <tier>)` variant and
// captures the referenced tier name.
var sameStepsAsRe = regexp.MustCompile(`same steps as ([A-Za-z][A-Za-z-]*)`)

// extractBacktickChain takes the FIRST backtick span in cell, splits it on
// U+2192, trims each part, and strips a trailing parenthetical. Any span past
// the first is ignored on purpose — PM's trivial cell carries two more
// backtick-quoted topic_keys after the step name, which must never leak into
// the resolved step list.
func extractBacktickChain(cell string) ([]string, error) {
	first := strings.Index(cell, "`")
	if first < 0 {
		return nil, fmt.Errorf("no backtick span found in cell %q", cell)
	}
	rest := cell[first+1:]
	second := strings.Index(rest, "`")
	if second < 0 {
		return nil, fmt.Errorf("unterminated backtick span in cell %q", cell)
	}
	span := rest[:second]

	parts := strings.Split(span, "→")
	steps := make([]string, 0, len(parts))
	for _, p := range parts {
		name := strings.TrimSpace(trailingParenRe.ReplaceAllString(strings.TrimSpace(p), ""))
		if name != "" {
			steps = append(steps, name)
		}
	}
	if len(steps) == 0 {
		return nil, fmt.Errorf("backtick span in cell %q resolved to no step names", cell)
	}
	return steps, nil
}

// classifyFullWorkflowCell resolves the three `Full workflow (...)` variants:
// an explicit backtick chain, `all steps` (every subagent step in workflow.yaml
// declaration order), and `same steps as <tier>` (that tier's already-resolved
// list). Any other wording is unrecognized.
func classifyFullWorkflowCell(trimmed string, graph specialistGraph, resolvedByTier map[string][]string) (orchestrationRowKind, []string, error) {
	if strings.Contains(trimmed, "`") {
		steps, err := extractBacktickChain(trimmed)
		if err != nil {
			return kindUnrecognized, nil, nil
		}
		return kindStepList, steps, nil
	}
	if strings.Contains(trimmed, "all steps") {
		return kindStepList, append([]string(nil), graph.subagentNamesInOrder()...), nil
	}
	if m := sameStepsAsRe.FindStringSubmatch(trimmed); m != nil {
		steps, ok := resolvedByTier[m[1]]
		if !ok {
			return kindUnrecognized, nil, fmt.Errorf("Full workflow references unresolved tier %q", m[1])
		}
		return kindStepList, append([]string(nil), steps...), nil
	}
	return kindUnrecognized, nil, nil
}

// classifyOrchestrationCell classifies one tier cell and, for a stepList
// cell, resolves its step names. noRun is checked by prefix on the trimmed
// cell BEFORE any backtick extraction — qa/trivial's cell carries a `simple`
// backtick span inside its "Not eligible" prose, and extracting backticks
// first would misread it as a phantom stepList.
func classifyOrchestrationCell(cell string, graph specialistGraph, resolvedByTier map[string][]string) (orchestrationRowKind, []string, error) {
	trimmed := strings.TrimSpace(cell)
	switch {
	case strings.HasPrefix(trimmed, "Not called"),
		strings.HasPrefix(trimmed, "Not eligible"),
		strings.HasPrefix(trimmed, "—"):
		return kindNoRun, nil, nil
	case strings.HasPrefix(trimmed, "Full workflow"):
		return classifyFullWorkflowCell(trimmed, graph, resolvedByTier)
	case strings.Contains(trimmed, "`"):
		steps, err := extractBacktickChain(trimmed)
		if err != nil {
			return kindUnrecognized, nil, nil
		}
		return kindStepList, steps, nil
	default:
		return kindUnrecognized, nil, nil
	}
}

// parseOrchestrationPlan locates a specialist's tier table by scanning
// forward from the `## Orchestration Plan` heading to the first `| Level |
// Steps |` header row — never a hardcoded line offset, since researcher's
// table sits at a different offset than its siblings — then reads rows until
// the first non-`|` line.
func parseOrchestrationPlan(md []byte, graph specialistGraph) ([]orchestrationRow, error) {
	lines := strings.Split(string(md), "\n")
	inSection := false
	inTable := false
	resolvedByTier := map[string][]string{}
	var rows []orchestrationRow

	for _, line := range lines {
		if !inSection {
			if strings.Contains(line, "## Orchestration Plan") {
				inSection = true
			}
			continue
		}
		if !inTable {
			if strings.Contains(line, "| Level | Steps |") {
				inTable = true
			}
			continue
		}

		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			break
		}
		if strings.HasPrefix(trimmed, "|---") || strings.HasPrefix(trimmed, "|--") {
			continue
		}

		cells := strings.Split(trimmed, "|")
		if len(cells) < 3 {
			continue
		}
		tier := strings.Trim(strings.TrimSpace(cells[1]), "*")
		cell := strings.TrimSpace(cells[2])

		kind, steps, err := classifyOrchestrationCell(cell, graph, resolvedByTier)
		if err != nil {
			return nil, fmt.Errorf("%s/%s: %w", graph.Dir, tier, err)
		}
		if kind == kindStepList {
			resolvedByTier[tier] = steps
		}
		rows = append(rows, orchestrationRow{Dir: graph.Dir, Tier: tier, Kind: kind, Cell: cell, Steps: steps})
	}

	return rows, nil
}

// applyPass1 checks every candidate step name against the full set of names
// declared in the specialist's workflow.yaml — inline steps included.
func applyPass1(steps []string, graph specialistGraph) error {
	known := graph.allStepNames()
	for _, name := range steps {
		if !known[name] {
			return fmt.Errorf("phantom step name %q not declared in %s/workflow.yaml", name, graph.Dir)
		}
	}
	return nil
}

// applyPass2 runs the dependency-completion fixpoint: for each step T (front
// to back), each of T's non-optional inputs is resolved to its producer P; if
// P is a subagent step not already in the list, P is inserted immediately
// before T, and the sweep restarts. Optional inputs, inline producers, and
// inputs with no producer in this specialist's own graph (a cross-specialist
// input) are all skipped independently. The hard sweep cap turns a dependency
// cycle into an error instead of a hang; the caller is expected to t.Fatalf on
// it.
func applyPass2(steps []string, graph specialistGraph) ([]string, bool, error) {
	const maxSweeps = 64

	result := append([]string(nil), steps...)
	insertedAny := false

	for sweep := 0; ; sweep++ {
		if sweep >= maxSweeps {
			return nil, false, fmt.Errorf("pass2 fixpoint did not converge after %d sweeps for %s (possible dependency cycle)", maxSweeps, graph.Dir)
		}

		mutated := false
		for i := 0; i < len(result); i++ {
			step, ok := graph.byName(result[i])
			if !ok {
				continue
			}
			for _, in := range step.Inputs {
				if in.Optional {
					continue
				}
				producer, found := graph.producerOf(in.TopicKey)
				if !found {
					continue
				}
				if producer.Execution != "subagent" {
					continue
				}
				if stepListContains(result, producer.Name) {
					continue
				}

				grown := make([]string, 0, len(result)+1)
				grown = append(grown, result[:i]...)
				grown = append(grown, producer.Name)
				grown = append(grown, result[i:]...)
				result = grown

				mutated = true
				insertedAny = true
				break
			}
			if mutated {
				break
			}
		}
		if !mutated {
			break
		}
	}

	return result, insertedAny, nil
}

func stepListContains(steps []string, name string) bool {
	for _, s := range steps {
		if s == name {
			return true
		}
	}
	return false
}

// stepListPresets derives a specialist's Collapse ladder and preset table
// directly from its own parsed tier rows — the gating axis is never
// hardcoded, so this reads `trivial -> simple -> moderate -> complex` for most
// specialists and `none -> moderate -> high` for security without either
// axis name ever appearing in this function.
func stepListPresets(rows []orchestrationRow) ([]string, map[string][]string) {
	var ladder []string
	presets := make(map[string][]string, len(rows))
	for _, row := range rows {
		if row.Kind != kindStepList {
			continue
		}
		ladder = append(ladder, row.Tier)
		presets[row.Tier] = row.Steps
	}
	return ladder, presets
}

// applyCollapse is a no-op unless Pass 2 inserted at least one step. Once
// gated open, it returns the SMALLEST preset whose step set is a superset of
// the grown list, matching the Collapse fallback's wording in skill/SKILL.md.
// The containment direction is load-bearing: taking the largest preset
// CONTAINED IN the grown list instead would inflate the label whenever two
// tiers share one step set (researcher's three tiers, qa's moderate and
// complex), labelling a list that merely grew to match `simple` as `complex`.
// When no preset covers the grown list, the original label stands.
func applyCollapse(originalTier string, ladder []string, presets map[string][]string, grownList []string, insertedAny bool) string {
	if !insertedAny {
		return originalTier
	}

	grownSet := make(map[string]bool, len(grownList))
	for _, s := range grownList {
		grownSet[s] = true
	}

	for _, tier := range ladder {
		preset, ok := presets[tier]
		if !ok {
			continue
		}
		if isSupersetOf(preset, grownSet) {
			return tier
		}
	}
	return originalTier
}

// isSupersetOf reports whether every member of set appears in preset.
func isSupersetOf(preset []string, set map[string]bool) bool {
	presetSet := make(map[string]bool, len(preset))
	for _, s := range preset {
		presetSet[s] = true
	}
	for s := range set {
		if !presetSet[s] {
			return false
		}
	}
	return true
}

// TestOrchestrationPlanTiersAreFixpoints asserts that every step-list tier row
// of every routed specialist is already dependency-complete: Pass 1 rejects
// nothing, Pass 2 inserts nothing, and Collapse leaves the tier label
// byte-equal to the row's own key.
func TestOrchestrationPlanTiersAreFixpoints(t *testing.T) {
	root := skillDir(t)

	for _, dir := range workflowSpecialistDirs {
		wfPath := filepath.Join(root, dir, "workflow.yaml")
		wfData, err := os.ReadFile(wfPath)
		if err != nil {
			t.Fatalf("%s: read workflow.yaml: %v", dir, err)
		}
		graph, err := parseSpecialistGraph(dir, wfData)
		if err != nil {
			t.Fatalf("%s: parse workflow.yaml: %v", dir, err)
		}

		mdPath := filepath.Join(root, dir, "SKILL.md")
		mdData, err := os.ReadFile(mdPath)
		if err != nil {
			t.Fatalf("%s: read SKILL.md: %v", dir, err)
		}
		rows, err := parseOrchestrationPlan(mdData, graph)
		if err != nil {
			t.Fatalf("%s: parse Orchestration Plan: %v", dir, err)
		}

		ladder, presets := stepListPresets(rows)

		for _, row := range rows {
			if row.Kind != kindStepList {
				continue
			}
			row := row
			t.Run(dir+"/"+row.Tier, func(t *testing.T) {
				if err := applyPass1(row.Steps, graph); err != nil {
					t.Fatalf("Pass 1: %v", err)
				}

				grown, insertedAny, err := applyPass2(row.Steps, graph)
				if err != nil {
					t.Fatalf("Pass 2: %v", err)
				}
				if insertedAny {
					t.Errorf("Pass 2 inserted a step; want the row already dependency-complete, got %v", grown)
				}
				if !reflect.DeepEqual(grown, row.Steps) {
					t.Errorf("Pass 2 result = %v, want unchanged %v", grown, row.Steps)
				}

				label := applyCollapse(row.Tier, ladder, presets, grown, insertedAny)
				if label != row.Tier {
					t.Errorf("Collapse label = %q, want byte-equal %q", label, row.Tier)
				}
			})
		}
	}
}

// TestOrchestrationPlanCellClassification asserts the tri-state classifier
// over every specialist's tier table: the exact row counts per specialist,
// the exact set of noRun cells, zero unrecognized cells, and that asdt-init
// carries the heading but no tier table at all.
func TestOrchestrationPlanCellClassification(t *testing.T) {
	root := skillDir(t)

	wantRowCount := map[string]int{
		"asdt-architect":  4,
		"asdt-developer":  4,
		"asdt-pm":         4,
		"asdt-qa":         4,
		"asdt-researcher": 4,
		"asdt-security":   3,
		"asdt-ux-ui":      4,
	}
	wantNoRun := map[string]bool{
		"asdt-architect/simple": true,
		"asdt-qa/trivial":       true,
		"asdt-security/none":    true,
	}
	foundNoRun := map[string]bool{}

	totalRows := 0
	unrecognizedCount := 0

	for _, dir := range workflowSpecialistDirs {
		wfData, err := os.ReadFile(filepath.Join(root, dir, "workflow.yaml"))
		if err != nil {
			t.Fatalf("%s: read workflow.yaml: %v", dir, err)
		}
		graph, err := parseSpecialistGraph(dir, wfData)
		if err != nil {
			t.Fatalf("%s: parse workflow.yaml: %v", dir, err)
		}

		mdData, err := os.ReadFile(filepath.Join(root, dir, "SKILL.md"))
		if err != nil {
			t.Fatalf("%s: read SKILL.md: %v", dir, err)
		}
		rows, err := parseOrchestrationPlan(mdData, graph)
		if err != nil {
			t.Fatalf("%s: parse Orchestration Plan: %v", dir, err)
		}

		if got := len(rows); got != wantRowCount[dir] {
			t.Errorf("%s: row count = %d, want %d", dir, got, wantRowCount[dir])
		}
		totalRows += len(rows)

		for _, row := range rows {
			key := dir + "/" + row.Tier
			switch row.Kind {
			case kindNoRun:
				foundNoRun[key] = true
				if !wantNoRun[key] {
					t.Errorf("unexpected noRun row %q", key)
				}
			case kindUnrecognized:
				unrecognizedCount++
				t.Errorf("%s: cell classified unrecognized: %q", key, row.Cell)
			case kindStepList:
				if len(row.Steps) == 0 {
					t.Errorf("%s: stepList row resolved to an empty step list", key)
				}
			}
		}
	}

	if totalRows != 27 {
		t.Errorf("total orchestration rows = %d, want 27", totalRows)
	}
	for key := range wantNoRun {
		if !foundNoRun[key] {
			t.Errorf("expected noRun row %q was not found", key)
		}
	}
	if unrecognizedCount != 0 {
		t.Errorf("unrecognized cell count = %d, want 0", unrecognizedCount)
	}

	initMD, err := os.ReadFile(filepath.Join(root, "asdt-init", "SKILL.md"))
	if err != nil {
		t.Fatalf("read asdt-init/SKILL.md: %v", err)
	}
	initContent := string(initMD)
	if !strings.Contains(initContent, "## Orchestration Plan") {
		t.Errorf("asdt-init/SKILL.md missing ## Orchestration Plan heading")
	}
	if strings.Contains(initContent, "| Level | Steps |") {
		t.Errorf("asdt-init/SKILL.md unexpectedly declares a | Level | Steps | tier table")
	}
}

// TestOptionalMarkerReadsRawSourceLine drives the raw-source-line marker
// reader against four in-memory fixtures — A, B, and C each place the `#
// optional` marker in a different structural position and must be detected;
// D carries a non-`# optional` comment on a sibling `model:` line and must
// not be detected as a marker on any input — then cross-checks the real tree's
// marker distribution against a hardcoded per-specialist count.
func TestOptionalMarkerReadsRawSourceLine(t *testing.T) {
	fixtureA := []byte(`specialist: fixture
steps:
  - name: sample
    execution: subagent
    inputs:
      - "fixture/a"
      - "fixture/b"  # optional
    output_topic_key: "fixture/out"
`)
	fixtureB := []byte(`specialist: fixture
steps:
  - name: sample
    execution: subagent
    inputs:
      - "fixture/a"
      - "fixture/b"  # optional
      # trailing comment block before the next key
    output_topic_key: "fixture/out"
`)
	fixtureC := []byte(`specialist: fixture
steps:
  - name: sample
    execution: subagent
    inputs:
      - "fixture/a"  # optional
      - "fixture/b"
    output_topic_key: "fixture/out"
`)
	fixtureD := []byte(`specialist: fixture
steps:
  - name: sample
    execution: subagent
    model: sonnet  # favor throughput, not an optional marker
    inputs:
      - "fixture/a"
      - "fixture/b"
    output_topic_key: "fixture/out"
`)

	cases := []struct {
		name       string
		data       []byte
		wantMarked map[string]bool
	}{
		{"A", fixtureA, map[string]bool{"fixture/a": false, "fixture/b": true}},
		{"B", fixtureB, map[string]bool{"fixture/a": false, "fixture/b": true}},
		{"C", fixtureC, map[string]bool{"fixture/a": true, "fixture/b": false}},
		{"D", fixtureD, map[string]bool{"fixture/a": false, "fixture/b": false}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			graph, err := parseSpecialistGraph("fixture", tc.data)
			if err != nil {
				t.Fatalf("parse fixture: %v", err)
			}
			step, ok := graph.byName("sample")
			if !ok {
				t.Fatalf("step %q not found in fixture", "sample")
			}
			got := map[string]bool{}
			for _, in := range step.Inputs {
				got[in.TopicKey] = in.Optional
			}
			if !reflect.DeepEqual(got, tc.wantMarked) {
				t.Errorf("optional flags = %v, want %v", got, tc.wantMarked)
			}
		})
	}

	root := skillDir(t)
	wantCounts := map[string]int{
		"asdt-architect": 3,
		"asdt-developer": 4,
		"asdt-pm":        3,
		"asdt-qa":        4,
		"asdt-security":  1,
		"asdt-ux-ui":     3,
	}

	total := 0
	for _, dir := range workflowSpecialistDirs {
		wfData, err := os.ReadFile(filepath.Join(root, dir, "workflow.yaml"))
		if err != nil {
			t.Fatalf("%s: read workflow.yaml: %v", dir, err)
		}
		graph, err := parseSpecialistGraph(dir, wfData)
		if err != nil {
			t.Fatalf("%s: parse workflow.yaml: %v", dir, err)
		}

		count := 0
		for _, step := range graph.Steps {
			for _, in := range step.Inputs {
				if in.Optional {
					count++
				}
			}
		}
		if count != wantCounts[dir] {
			t.Errorf("%s: optional marker count = %d, want %d", dir, count, wantCounts[dir])
		}
		total += count
	}
	if total != 18 {
		t.Errorf("total optional markers across the tree = %d, want 18", total)
	}
}

// TestPass2InsertsWhenOptionalMarkerRemoved proves Pass 2 actually enforces
// the `# optional` skip clause: with the marker removed from qa's
// test-case-generation input on qa/test-strategy — as an in-memory fixture
// ONLY, the marker is never stripped from the working tree — Pass 2 must
// auto-insert test-strategy before test-case-generation.
//
// test-strategy's OWN qa/edge-cases input carries no `# optional` marker
// anywhere in the real qa/workflow.yaml, so the fixpoint recurses on the
// newly-inserted test-strategy and also inserts edge-case-analysis before it.
// The grown list ends up element-for-element equal to qa's moderate preset —
// and, because qa/complex resolves as "same steps as moderate", to the
// complex preset too — so Collapse promotes the label all the way to
// `complex`, not `simple`.
func TestPass2InsertsWhenOptionalMarkerRemoved(t *testing.T) {
	root := skillDir(t)

	wfPath := filepath.Join(root, "asdt-qa", "workflow.yaml")
	wfData, err := os.ReadFile(wfPath)
	if err != nil {
		t.Fatalf("read %s: %v", wfPath, err)
	}

	marker := []byte(`"{project}/{change}/qa/test-strategy"  # optional (tier-gated: test-strategy runs at moderate+)`)
	stripped := []byte(`"{project}/{change}/qa/test-strategy"`)
	fixtureData := bytes.Replace(wfData, marker, stripped, 1)
	if bytes.Equal(fixtureData, wfData) {
		t.Fatalf("in-memory fixture did not find the qa/test-strategy optional marker to strip; the source line may have changed")
	}

	graph, err := parseSpecialistGraph("asdt-qa", fixtureData)
	if err != nil {
		t.Fatalf("parse fixture workflow.yaml: %v", err)
	}

	mdData, err := os.ReadFile(filepath.Join(root, "asdt-qa", "SKILL.md"))
	if err != nil {
		t.Fatalf("read asdt-qa/SKILL.md: %v", err)
	}
	rows, err := parseOrchestrationPlan(mdData, graph)
	if err != nil {
		t.Fatalf("parse Orchestration Plan: %v", err)
	}
	ladder, presets := stepListPresets(rows)

	var simpleRow orchestrationRow
	found := false
	for _, row := range rows {
		if row.Kind == kindStepList && row.Tier == "simple" {
			simpleRow = row
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("qa/simple row not found")
	}

	grown, insertedAny, err := applyPass2(simpleRow.Steps, graph)
	if err != nil {
		t.Fatalf("Pass 2: %v", err)
	}
	if !insertedAny {
		t.Fatalf("Pass 2 did not insert a step after the optional marker was stripped; want test-strategy and edge-case-analysis inserted")
	}

	wantGrown := []string{
		"load-requirements",
		"ac-validation",
		"edge-case-analysis",
		"test-strategy",
		"test-case-generation",
		"quality-report",
		"performance-validation",
		"review",
	}
	if !reflect.DeepEqual(grown, wantGrown) {
		t.Errorf("Pass 2 result = %v, want %v", grown, wantGrown)
	}

	label := applyCollapse(simpleRow.Tier, ladder, presets, grown, insertedAny)
	if label != "moderate" {
		t.Errorf("Collapse label = %q, want %q (the grown list equals qa's moderate preset exactly, and moderate is the SMALLEST preset covering it — qa's complex shares that same step set but is larger on the ladder)", label, "moderate")
	}
}

// TestCollapseOnlyRunsAfterInsertion pins both halves of the Collapse gate on
// researcher, whose simple, moderate and complex tiers all resolve to one
// identical three-step list — the specialist the unguarded rule relabelled
// downward. With zero insertions the label must survive untouched. With a real
// growth (trivial's single step grown to that three-step list) the label must
// land on `simple`, the smallest preset covering it, never on `complex`.
func TestCollapseOnlyRunsAfterInsertion(t *testing.T) {
	root := skillDir(t)
	dir := "asdt-researcher"

	wfData, err := os.ReadFile(filepath.Join(root, dir, "workflow.yaml"))
	if err != nil {
		t.Fatalf("read workflow.yaml: %v", err)
	}
	graph, err := parseSpecialistGraph(dir, wfData)
	if err != nil {
		t.Fatalf("parse workflow.yaml: %v", err)
	}

	mdData, err := os.ReadFile(filepath.Join(root, dir, "SKILL.md"))
	if err != nil {
		t.Fatalf("read SKILL.md: %v", err)
	}
	rows, err := parseOrchestrationPlan(mdData, graph)
	if err != nil {
		t.Fatalf("parse Orchestration Plan: %v", err)
	}
	ladder, presets := stepListPresets(rows)

	byTier := make(map[string]orchestrationRow, len(rows))
	for _, row := range rows {
		if row.Kind == kindStepList {
			byTier[row.Tier] = row
		}
	}
	complexRow, ok := byTier["complex"]
	if !ok {
		t.Fatalf("researcher/complex row not found")
	}
	trivialRow, ok := byTier["trivial"]
	if !ok {
		t.Fatalf("researcher/trivial row not found")
	}

	if label := applyCollapse(complexRow.Tier, ladder, presets, complexRow.Steps, false); label != "complex" {
		t.Errorf("Collapse with zero insertions = %q, want %q unchanged", label, "complex")
	}

	if label := applyCollapse(trivialRow.Tier, ladder, presets, complexRow.Steps, true); label != "simple" {
		t.Errorf("Collapse after real growth = %q, want %q (trivial's single step grew to the three-step list, and simple is the smallest preset covering it)", label, "simple")
	}
}
