package installer

import (
	"path/filepath"
	"slices"
	"strings"
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

// TestDeveloperSmallLane guards the Developer's small lane: a build whose one
// bounded read shows it understood, contained (shared code included), and
// resumable gets a plan the orchestrator writes itself, so neither explore nor
// spec launches, while approve, implement, and verify run with the same
// consent gates as every build. The lane lives in prose only, so the test pins
// the few phrases that carry it, each scoped to its section, and keeps explore
// and spec as subagent steps for every other build. It also pins the router
// side: UX/UI's registry row excludes a cosmetic change, and a trivial request
// with no risk surface is a chain of one.
func TestDeveloperSmallLane(t *testing.T) {
	const smallLaneLead = "**The small lane.**"

	root := skillDir(t)
	workflowPath, wf := readWorkflowFile(t, root, "asdt-developer")
	skillPath := filepath.Join(root, "asdt-developer", "SKILL.md")
	skillMD := readText(t, skillPath)
	routerPath := filepath.Join(root, "SKILL.md")
	routerMD := readText(t, routerPath)

	section := func(path, doc, name string) string {
		t.Helper()
		body, ok := markdownSection(doc, name)
		if !ok {
			t.Fatalf("%s: no `## %s` section, want one", path, name)
		}
		return body
	}
	tableRows := func(body string) []string {
		var rows []string
		for _, line := range strings.Split(body, "\n") {
			if strings.HasPrefix(line, "|") {
				rows = append(rows, line)
			}
		}
		return rows
	}

	plan := section(skillPath, skillMD, "Orchestration Plan")
	_, lane, ok := strings.Cut(plan, smallLaneLead)
	if !ok {
		t.Fatalf("%s: `## Orchestration Plan` has no block led by %s, want one", skillPath, smallLaneLead)
	}
	planRows := tableRows(plan)
	approve := section(skillPath, skillMD, "approve")
	invariants := section(skillPath, skillMD, "Invariants")
	assess := section(routerPath, routerMD, "How to assess")
	var uxUIRows []string
	for _, row := range tableRows(section(routerPath, routerMD, "Registry")) {
		if firstCommandDir(row) == "asdt-ux-ui" {
			uxUIRows = append(uxUIRows, row)
		}
	}

	cases := []struct {
		name  string
		path  string
		where string   // names the candidates in a failure message
		in    []string // one candidate must carry every phrase in want
		want  []string
	}{
		{name: "small build row", path: skillPath, where: "`## Orchestration Plan` table row", in: planRows, want: []string{"small", "`approve → implement → verify`"}},
		{name: "full build row", path: skillPath, where: "`## Orchestration Plan` table row", in: planRows, want: []string{"`explore → spec → approve → implement → verify`"}},
		{name: "three criteria", path: skillPath, where: "small lane block", in: []string{lane}, want: []string{"understood", "contained", "resumable"}},
		{name: "shared code is not contained", path: skillPath, where: "small lane block", in: []string{lane}, want: []string{"shared"}},
		{name: "no explore or spec launch", path: skillPath, where: "small lane block", in: []string{lane}, want: []string{"launch neither `explore` nor `spec`"}},
		{name: "doubt runs the full chain", path: skillPath, where: "small lane block", in: []string{lane}, want: []string{"full chain"}},
		{name: "approve shows the lane", path: skillPath, where: "`## approve` section", in: []string{approve}, want: []string{"Lane: small"}},
		{name: "approve writes the lane payload", path: skillPath, where: "`## approve` section", in: []string{approve}, want: []string{"on a resume or in the small lane"}},
		{name: "lane changes only the plan", path: skillPath, where: "`## Invariants` section", in: []string{invariants}, want: []string{"The small lane changes only how the plan is made"}},
		{name: "approval invariant survives", path: skillPath, where: "`## Invariants` section", in: []string{invariants}, want: []string{"No host file is written without approval"}},
		{name: "ux-ui row excludes a cosmetic change", path: routerPath, where: "`## Registry` row led by the UX/UI command", in: uxUIRows, want: []string{"cosmetic"}},
		{name: "trivial with no risk is a chain of one", path: routerPath, where: "`## How to assess` section", in: []string{assess}, want: []string{"chain of one"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			carriesAll := func(text string) bool {
				for _, phrase := range tc.want {
					if !strings.Contains(text, phrase) {
						return false
					}
				}
				return true
			}
			if !slices.ContainsFunc(tc.in, carriesAll) {
				t.Errorf("%s: no %s contains all of %q, want one that does", tc.path, tc.where, tc.want)
			}
		})
	}

	for _, name := range []string{"explore", "spec"} {
		t.Run("workflow keeps "+name, func(t *testing.T) {
			i := slices.IndexFunc(wf.Steps, func(s workflowStep) bool { return s.Name == name })
			if i < 0 {
				t.Fatalf("%s: step %q is missing, want it kept for every build outside the small lane", workflowPath, name)
			}
			if got := wf.Steps[i].Execution; got != "subagent" {
				t.Errorf("%s: step %q execution = %q, want %q", workflowPath, name, got, "subagent")
			}
		})
	}
}
