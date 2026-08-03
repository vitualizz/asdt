package installer

import (
	"os"
	"path/filepath"
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
	Name      string `yaml:"name"`
	Execution string `yaml:"execution"`
	Agent     string `yaml:"agent"`
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
		path := filepath.Join(root, dir, "workflow.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}

		var wf workflowFile
		if err := yaml.Unmarshal(data, &wf); err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

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
