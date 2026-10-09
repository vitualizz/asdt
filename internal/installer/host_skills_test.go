package installer

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// kebabCaseName is the shape of a host skill name: lowercase words joined by
// single hyphens, the form skill directories and Skill-tool names take.
var kebabCaseName = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// TestHostSkillsDeclaredOnSubagentSteps guards the `host_skills:` field
// (asdt-core/protocol.md §4). A host skill is injected by the orchestrator into
// a launched sub-agent's prompt, so it is only meaningful on a `subagent` step;
// on an inline step there is no prompt to inject it into. Every entry is a
// skill NAME looked up in the host's live skill list — never a path — so it
// must be non-empty kebab-case. The step file must name the injected
// `### HOST SKILL` block, or the sub-agent receives a skill its contract never
// says when to apply or what outranks it. The two UI steps must keep declaring
// frontend-design: dropping it silently removes the visual-design craft from
// greenfield UI work, with no other test noticing.
func TestHostSkillsDeclaredOnSubagentSteps(t *testing.T) {
	root := skillDir(t)
	declared := map[string][]string{}

	for _, dir := range append(append([]string{}, workflowSpecialistDirs...), "asdt-init") {
		path, wf := readWorkflowFile(t, root, dir)

		for _, s := range wf.Steps {
			if len(s.HostSkills) == 0 {
				continue
			}
			if s.Execution != "subagent" {
				t.Errorf("%s: step %q is %q but declares host_skills %v; only subagent steps receive injected host skills", path, s.Name, s.Execution, s.HostSkills)
			}
			for _, name := range s.HostSkills {
				if !kebabCaseName.MatchString(name) {
					t.Errorf("%s: step %q host_skills entry %q is not a non-empty kebab-case skill name", path, s.Name, name)
				}
			}
			stepPath := filepath.Join(root, dir, s.Skill)
			data, err := os.ReadFile(stepPath)
			if err != nil {
				t.Errorf("%s: step %q declares host_skills but its step file is unreadable: %v", path, s.Name, err)
			} else if !strings.Contains(string(data), "### HOST SKILL") {
				t.Errorf("%s: step %q declares host_skills %v, but %s never mentions the injected `### HOST SKILL` block", path, s.Name, s.HostSkills, stepPath)
			}
			declared[dir+"/"+s.Name] = s.HostSkills
		}
	}

	for _, want := range []string{"asdt-ux-ui/ui-design", "asdt-developer/implement"} {
		found := false
		for _, name := range declared[want] {
			if name == "frontend-design" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s: host_skills = %v, want it to include %q", want, declared[want], "frontend-design")
		}
	}
}

// TestUXUIDesignCraftIsBuiltIn guards the split of asdt-ux-ui into ux-spec →
// ui-design. ux-spec hands its flows forward in context and persists nothing;
// ui-design receives them as `ux-flows`, persists the specialist's ONE
// hand-off, and declares the shipped visual-design reference. The reference
// is what makes the design craft independent of any host skill: dropping it
// would leave ui-design designing on frontend-design alone, which a host may
// not have, with no other test noticing.
func TestUXUIDesignCraftIsBuiltIn(t *testing.T) {
	const (
		handoffKey = "{project}/{change}/ux-ui/handoff"
		craftRef   = "../asdt-core/references/visual-design.md"
	)

	root := skillDir(t)
	path, wf := readWorkflowFile(t, root, "asdt-ux-ui")
	steps := make(map[string]workflowStep, len(wf.Steps))
	for _, s := range wf.Steps {
		steps[s.Name] = s
	}

	spec, ok := steps["ux-spec"]
	if !ok {
		t.Fatalf("%s: step %q is missing", path, "ux-spec")
	}
	if spec.Output != "context" || spec.OutputTopicKey != "" {
		t.Errorf("%s: ux-spec output = %q, output_topic_key = %q; want output: context and no key, ui-design persists the one hand-off", path, spec.Output, spec.OutputTopicKey)
	}

	design, ok := steps["ui-design"]
	if !ok {
		t.Fatalf("%s: step %q is missing", path, "ui-design")
	}
	if design.OutputTopicKey != handoffKey {
		t.Errorf("%s: ui-design output_topic_key = %q, want %q", path, design.OutputTopicKey, handoffKey)
	}
	if !slices.Contains(design.ContextInputs, "ux-flows") {
		t.Errorf("%s: ui-design context_inputs = %v, want it to include %q", path, design.ContextInputs, "ux-flows")
	}
	if !slices.Contains(design.ReferenceSkills, craftRef) {
		t.Errorf("%s: ui-design reference_skills = %v, want it to include %q", path, design.ReferenceSkills, craftRef)
	}
}
