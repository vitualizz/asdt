package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// This file guards ONE property that yaml.v3 cannot express: a step's input is
// marked optional by a trailing `# optional` comment, and comments are dropped
// by the parser. Optionality therefore has to be read off the raw source line
// the node came from — if that ever regresses, every optional input silently
// becomes mandatory.
//
// It is what survived tier_preset_validation_test.go. The rest of that file
// validated the router's two-pass "Step List Validation" fixpoint against the
// per-specialist tier tables; Phase 5 deleted that algorithm from the router
// and the tier tables with it, so the Go reimplementation had nothing left to
// validate.

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

func (g specialistGraph) byName(name string) (stepInfo, bool) {
	for _, s := range g.Steps {
		if s.Name == name {
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
	// Re-measured against the final post-refactor tree. asdt-researcher declares
	// no inputs at all, so it contributes zero and is absent from this map.
	wantCounts := map[string]int{
		"asdt-architect": 1,
		"asdt-developer": 3,
		"asdt-pm":        1,
		"asdt-qa":        3,
		"asdt-security":  2,
		"asdt-ux-ui":     1,
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
	if total != 11 {
		t.Errorf("total optional markers across the tree = %d, want 11", total)
	}
}
