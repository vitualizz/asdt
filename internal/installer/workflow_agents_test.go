package installer

import (
	"os"
	"path/filepath"
	"regexp"
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
	Name       string   `yaml:"name"`
	Skill      string   `yaml:"skill"`
	Execution  string   `yaml:"execution"`
	Agent      string   `yaml:"agent"`
	Inputs     []string `yaml:"inputs"`
	HostSkills []string `yaml:"host_skills"`
	// Output is `summary` on a step whose return the orchestrator writes as
	// the role summary (protocol.md §1); absent otherwise.
	Output          string   `yaml:"output"`
	ReferenceSkills []string `yaml:"reference_skills"`
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
// exactly the contract: the developer's implement step is the only builder
// step (1), everything else is analyst (15).
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
	// 18 subagent steps across the tree, 16 of them routed: architect 2,
	// developer 4, qa 2, ux-ui 3, pm 2, security 2, researcher 1. Five of
	// those are the `review` study steps. One of the routed 16 is a builder
	// (developer/implement), leaving 15 analysts.
	if analystCount != 15 {
		t.Errorf("analyst subagent steps = %d, want 15", analystCount)
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

// inputName is the shape of a workflow.yaml `inputs:` entry (TEMPLATE.md §3):
// `{role}`, another role's summary, or `{role}/{step}`, a step payload of the
// specialist's own run. Both resolve to a scratchpad path (protocol.md §1).
var inputName = regexp.MustCompile(`^[a-z][a-z0-9-]*(/[a-z][a-z0-9-]*)?$`)

// TestStepInputsResolve guards TEMPLATE.md §3 and protocol.md §1, §5: every
// `inputs:` entry names a file the orchestrator can resolve. A `{role}` entry
// names a known specialist; a `{role}/{step}` entry names a step of the same
// workflow, under the workflow's own `specialist:` id, and the consumer must
// mention the `### INPUT {role}/{step}` line it receives in its step file or,
// failing that, in SKILL.md. Only subagent steps declare inputs; `feedback`
// and `previous` are reserved for the orchestrator's re-launches; `output:`
// is `summary` or absent; and every routed specialist writes a role summary
// from at least one step. A renamed producer otherwise leaves the consumer
// reading a path nothing writes.
func TestStepInputsResolve(t *testing.T) {
	root := skillDir(t)
	dirs := append(append([]string{}, workflowSpecialistDirs...), "asdt-init")

	roles := map[string]bool{}
	files := map[string]workflowFile{}
	for _, dir := range dirs {
		_, wf := readWorkflowFile(t, root, dir)
		files[dir] = wf
		roles[wf.Specialist] = true
	}

	for _, dir := range dirs {
		wf := files[dir]
		path := filepath.Join(root, dir, "workflow.yaml")
		skillMD := readText(t, filepath.Join(root, dir, "SKILL.md"))

		stepNames := map[string]bool{}
		summaries := 0
		for _, s := range wf.Steps {
			stepNames[s.Name] = true
			switch s.Output {
			case "":
			case "summary":
				summaries++
			default:
				t.Errorf("%s: step %q declares output %q, want `summary` or nothing", path, s.Name, s.Output)
			}
		}
		if dir != "asdt-init" && summaries == 0 {
			t.Errorf("%s: no step declares `output: summary`; the specialist hands nothing on", path)
		}

		for _, s := range wf.Steps {
			if len(s.Inputs) > 0 && s.Execution != "subagent" {
				t.Errorf("%s: %s step %q declares inputs %v; only subagent steps receive injected inputs", path, s.Execution, s.Name, s.Inputs)
			}
			var stepText string
			if strings.HasPrefix(s.Skill, "steps/") {
				stepText = readText(t, filepath.Join(root, dir, s.Skill))
			}
			for _, name := range s.Inputs {
				if !inputName.MatchString(name) {
					t.Errorf("%s: step %q input %q is neither `{role}` nor `{role}/{step}`", path, s.Name, name)
					continue
				}
				role, step, own := strings.Cut(name, "/")
				if name == "feedback" || name == "previous" {
					t.Errorf("%s: step %q declares the reserved input %q; the orchestrator injects it on a re-launch", path, s.Name, name)
					continue
				}
				if !own {
					if !roles[role] {
						t.Errorf("%s: step %q input %q names no known specialist", path, s.Name, name)
					}
					continue
				}
				if role != wf.Specialist {
					t.Errorf("%s: step %q input %q reads another specialist's step payload; read its role summary `%s` instead", path, s.Name, name, role)
				}
				if !stepNames[step] || step == s.Name {
					t.Errorf("%s: step %q input %q names no other step of this workflow", path, s.Name, name)
				}
				heading := "### INPUT " + name
				if !strings.Contains(stepText, heading) && !strings.Contains(skillMD, heading) {
					t.Errorf("%s: step %q lists input %q, but neither its step file nor SKILL.md mentions the injected `%s` line", path, s.Name, name, heading)
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

// TestWorkflowPathsResolve guards TEMPLATE.md §3: every `skill:` and
// `reference_skills:` path in every workflow.yaml under skill/ resolves, from
// the specialist's own directory, to an existing file inside skill/. A renamed
// or moved reference otherwise ships a step that tells its sub-agent to read a
// file that is not there, with no other test noticing.
func TestWorkflowPathsResolve(t *testing.T) {
	root := skillDir(t)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("resolve %s: %v", root, err)
	}
	matches, err := filepath.Glob(filepath.Join(root, "*", "workflow.yaml"))
	if err != nil {
		t.Fatalf("glob workflow.yaml under %s: %v", root, err)
	}
	if len(matches) == 0 {
		t.Fatalf("no workflow.yaml found under %s", root)
	}

	for _, match := range matches {
		dir := filepath.Base(filepath.Dir(match))
		path, wf := readWorkflowFile(t, root, dir)
		for _, s := range wf.Steps {
			refs := append([]string{}, s.ReferenceSkills...)
			if s.Skill != "" {
				refs = append(refs, s.Skill)
			}
			for _, ref := range refs {
				target, err := filepath.Abs(filepath.Join(root, dir, ref))
				if err != nil {
					t.Errorf("%s: step %q path %q: %v", path, s.Name, ref, err)
					continue
				}
				if rel, err := filepath.Rel(absRoot, target); err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					t.Errorf("%s: step %q path %q resolves to %s, outside skill/", path, s.Name, ref, target)
					continue
				}
				info, err := os.Stat(target)
				if err != nil {
					t.Errorf("%s: step %q path %q does not resolve to a file: %v", path, s.Name, ref, err)
					continue
				}
				if !info.Mode().IsRegular() {
					t.Errorf("%s: step %q path %q resolves to %s, which is not a regular file", path, s.Name, ref, target)
				}
			}
		}
	}
}
