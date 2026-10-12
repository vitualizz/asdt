package installer

import (
	"path/filepath"
	"slices"
	"testing"
)

// TestDeveloperSDDShape guards the Spec-Driven Development loop of
// asdt-developer/workflow.yaml: spec writes the developer's role summary (so a
// plan-only run hands a plan on), a human approval gate sits between spec and
// implement (so no host file is written without consent), and an inline verify
// gate follows implement (so the commands run only on the human's yes, from
// the orchestrator — never from the builder sub-agent). Both gates must stay
// inline: as subagent steps they could not pause for the human, and they would
// count against the four-step ceiling. Neither carries a `skill:` — its
// contract is a SKILL.md section, which must exist. verify declares
// `output: summary` because the orchestrator rewrites the role summary with
// the verification outcome added. implement reads the approved spec as its
// step payload — the role summary moves on to `implemented` — and verify's
// failures on a fix round.
func TestDeveloperSDDShape(t *testing.T) {
	root := skillDir(t)
	path, wf := readWorkflowFile(t, root, "asdt-developer")
	skillPath := filepath.Join(root, "asdt-developer", "SKILL.md")
	skillMD := readText(t, skillPath)

	index := make(map[string]int, len(wf.Steps))
	for i, s := range wf.Steps {
		index[s.Name] = i
	}
	for _, name := range []string{"spec", "approve", "implement", "verify"} {
		if _, ok := index[name]; !ok {
			t.Fatalf("%s: step %q is missing; the SDD loop needs spec → approve → implement → verify", path, name)
		}
	}

	cases := []struct {
		step            string
		wantExecution   string
		wantOutput      string   // "" means the step writes no role summary
		wantGateSection bool     // an inline gate: no skill:, its contract is a SKILL.md section
		wantInputs      []string // must all appear in the step's inputs
		after           string   // the step must come after this one ("" = no constraint)
		before          string   // the step must come before this one ("" = no constraint)
	}{
		{step: "spec", wantExecution: "subagent", wantOutput: "summary", wantInputs: []string{"developer/explore"}, before: "approve"},
		{step: "approve", wantExecution: "inline", wantGateSection: true, after: "spec", before: "implement"},
		{step: "implement", wantExecution: "subagent", wantOutput: "summary", wantInputs: []string{"developer/spec", "developer/verify"}, after: "approve", before: "verify"},
		{step: "verify", wantExecution: "inline", wantOutput: "summary", wantGateSection: true, after: "implement"},
	}

	for _, tc := range cases {
		t.Run(tc.step, func(t *testing.T) {
			step := wf.Steps[index[tc.step]]
			if step.Execution != tc.wantExecution {
				t.Errorf("%s: step %q execution = %q, want %q", path, tc.step, step.Execution, tc.wantExecution)
			}
			if step.Output != tc.wantOutput {
				t.Errorf("%s: step %q output = %q, want %q", path, tc.step, step.Output, tc.wantOutput)
			}
			if tc.wantGateSection {
				if step.Skill != "" {
					t.Errorf("%s: step %q declares skill %q, want none; an inline gate's contract is its SKILL.md section", path, tc.step, step.Skill)
				}
				if _, ok := markdownSection(skillMD, tc.step); !ok {
					t.Errorf("%s: no `## %s` section; it is the inline gate's whole contract", skillPath, tc.step)
				}
			}
			for _, name := range tc.wantInputs {
				if !slices.Contains(step.Inputs, name) {
					t.Errorf("%s: step %q inputs = %v, want it to include %q", path, tc.step, step.Inputs, name)
				}
			}
			if tc.after != "" && index[tc.step] <= index[tc.after] {
				t.Errorf("%s: step %q is at position %d, want it after %q (position %d)", path, tc.step, index[tc.step], tc.after, index[tc.after])
			}
			if tc.before != "" && index[tc.step] >= index[tc.before] {
				t.Errorf("%s: step %q is at position %d, want it before %q (position %d)", path, tc.step, index[tc.step], tc.before, index[tc.before])
			}
		})
	}
}
