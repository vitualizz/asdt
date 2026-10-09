package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// workflowSpecialistDirs lists the specialist directories whose workflow.yaml
// declares per-step agent types. asdt-init has its own workflow.yaml with
// `agent:` fields, but it is a non-routable setup-class command and is
// deliberately excluded from this per-specialist agent-type audit; asdt-core
// is the fragment library and has no workflow.yaml at all.
var workflowSpecialistDirs = []string{
	"asdt-architect",
	"asdt-developer",
	"asdt-pm",
	"asdt-qa",
	"asdt-security",
	"asdt-ux-ui",
	"asdt-researcher",
}

type workflowStep struct {
	Name           string   `yaml:"name"`
	Skill          string   `yaml:"skill"`
	Execution      string   `yaml:"execution"`
	Agent          string   `yaml:"agent"`
	OutputTopicKey string   `yaml:"output_topic_key"`
	ContextInputs  []string `yaml:"context_inputs"`
	HostSkills     []string `yaml:"host_skills"`
}

// readWorkflowFile reads and parses dir/workflow.yaml under the skill root.
func readWorkflowFile(t *testing.T, root, dir string) (string, workflowFile) {
	t.Helper()
	path := filepath.Join(root, dir, "workflow.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var wf workflowFile
	if err := yaml.Unmarshal(data, &wf); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return path, wf
}

type workflowFile struct {
	Specialist string         `yaml:"specialist"`
	Steps      []workflowStep `yaml:"steps"`
}

// skillDir resolves the repository's skill/ directory relative to this
// package directory (internal/installer), where `go test` runs.
func skillDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join("..", "..", "skill")
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("cannot locate skill/ relative to internal/installer: %v", err)
	}
	return dir
}

// TestWorkflowSubagentStepsDeclareKnownAgentTypes asserts that every
// `execution: subagent` step in every specialist workflow.yaml declares an
// `agent:` value drawn from AgentTypeNames, and that the agent-type split is
// exactly the contract: the developer's implement and test steps are the only
// builder steps (2), everything else is analyst (39).
func TestWorkflowSubagentStepsDeclareKnownAgentTypes(t *testing.T) {
	root := skillDir(t)
	known := make(map[string]bool, len(AgentTypeNames))
	for _, name := range AgentTypeNames {
		known[name] = true
	}

	analystCount := 0
	builderCount := 0
	var builderSteps []string

	for _, dir := range workflowSpecialistDirs {
		path, wf := readWorkflowFile(t, root, dir)

		for _, step := range wf.Steps {
			if step.Execution != "subagent" {
				if step.Agent != "" {
					t.Errorf("%s: step %q is %q but declares agent %q; only subagent steps carry an agent type", path, step.Name, step.Execution, step.Agent)
				}
				continue
			}

			if !known[step.Agent] {
				t.Errorf("%s: subagent step %q declares agent %q, want one of %v", path, step.Name, step.Agent, AgentTypeNames)
				continue
			}

			switch step.Agent {
			case "analyst":
				analystCount++
			case "builder":
				builderCount++
				builderSteps = append(builderSteps, dir+"/"+step.Name)
			}
		}
	}

	// Counts over the 7 routed specialists (asdt-init is excluded above).
	// 17 subagent steps across the tree, 15 of them routed: architect 2,
	// developer 4, qa 2, ux-ui 2, pm 2, security 2, researcher 1. Five of
	// those are the `review` study steps. One of the routed 15 is a builder
	// (developer/implement), leaving 14 analysts.
	if analystCount != 14 {
		t.Errorf("analyst subagent steps = %d, want 14", analystCount)
	}
	if builderCount != 1 {
		t.Errorf("builder subagent steps = %d, want 1 (got %v)", builderCount, builderSteps)
	}
	// Only implement writes host files. The former `test` step was absorbed
	// into implement, which now writes tests under the same mode and roots.
	wantBuilder := map[string]bool{
		"asdt-developer/implement": true,
	}
	for _, step := range builderSteps {
		if !wantBuilder[step] {
			t.Errorf("unexpected builder step %q; only asdt-developer implement and test may be builder", step)
		}
	}
}

// TestSubagentStepCeiling guards TEMPLATE.md §3: a specialist declares at most
// four `execution: subagent` steps. Inline preludes and gates are no sub-agent
// of their own and do not count. asdt-init is held to the same ceiling.
func TestSubagentStepCeiling(t *testing.T) {
	const ceiling = 4

	root := skillDir(t)
	for _, dir := range append(append([]string{}, workflowSpecialistDirs...), "asdt-init") {
		path, wf := readWorkflowFile(t, root, dir)
		var subagents []string
		for _, s := range wf.Steps {
			if s.Execution == "subagent" {
				subagents = append(subagents, s.Name)
			}
		}
		if len(subagents) > ceiling {
			t.Errorf("%s: %d subagent steps %v, want at most %d; merge two or split the specialist", path, len(subagents), subagents, ceiling)
		}
	}
}

// TestContextInputsMatchProduces guards protocol.md §1 "Intra-run payloads":
// every name a step lists in `context_inputs:` is the exact string a producer
// declares on a `Produces:` line — searched only where producers declare it: a
// step file's `## Output` section, or the SKILL.md section of an inline gate
// with no step file. The consumer must also name the `### INPUT {name}` block
// it receives, in its step file or, failing that, in SKILL.md. A renamed
// producer otherwise leaves the consumer waiting for a payload nothing emits.
func TestContextInputsMatchProduces(t *testing.T) {
	root := skillDir(t)
	for _, dir := range append(append([]string{}, workflowSpecialistDirs...), "asdt-init") {
		path, wf := readWorkflowFile(t, root, dir)
		skillPath := filepath.Join(root, dir, "SKILL.md")
		skillMD := readText(t, skillPath)

		var producers []string
		stepText := make(map[string]string, len(wf.Steps))
		for _, s := range wf.Steps {
			switch {
			case s.Skill == "" && s.Execution == "inline":
				if section, ok := markdownSection(skillMD, s.Name); ok {
					producers = append(producers, section)
				}
			case strings.HasPrefix(s.Skill, "steps/"):
				stepPath := filepath.Join(root, dir, s.Skill)
				data := readText(t, stepPath)
				stepText[s.Name] = data
				_, output, found := strings.Cut(data, "\n## Output")
				if !found {
					t.Errorf("%s: no `## Output` section", stepPath)
					continue
				}
				producers = append(producers, output)
			}
		}

		for _, s := range wf.Steps {
			for _, name := range s.ContextInputs {
				want := "Produces: `" + name + "`"
				produced := false
				for _, src := range producers {
					if strings.Contains(src, want) {
						produced = true
						break
					}
				}
				if !produced {
					t.Errorf("%s: step %q lists context input %q, but no step file `## Output` or inline-gate SKILL.md section declares %s", path, s.Name, name, want)
				}

				heading := "### INPUT " + name
				if !strings.Contains(stepText[s.Name], heading) && !strings.Contains(skillMD, heading) {
					t.Errorf("%s: step %q lists context input %q, but neither its step file nor %s mentions the injected `%s` block", path, s.Name, name, skillPath, heading)
				}
			}
		}
	}
}

// readText reads a file the test cannot proceed without.
func readText(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// markdownSection returns the body of the level-2 section whose heading is
// `## {name}` — alone, or followed by a space (`## approve — the plan gate`) —
// up to the next level-2 heading.
func markdownSection(doc, name string) (string, bool) {
	lines := strings.Split(doc, "\n")
	for i, line := range lines {
		rest, ok := strings.CutPrefix(line, "## "+name)
		if !ok || (rest != "" && !strings.HasPrefix(rest, " ")) {
			continue
		}
		end := len(lines)
		for j := i + 1; j < len(lines); j++ {
			if strings.HasPrefix(lines[j], "## ") {
				end = j
				break
			}
		}
		return strings.Join(lines[i+1:end], "\n"), true
	}
	return "", false
}
