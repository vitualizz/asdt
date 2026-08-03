package installer

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"gopkg.in/yaml.v3"

	"github.com/vitualizz/asdt/skill"
)

const testExecutorHeader = `# Executor Header — injected by the orchestrator into every sub-agent prompt

> **EXECUTOR**: You are the sub-agent assigned this single step. Do the work
> described here yourself and return.
`

// agentFixtureFS carries the one file agent generation depends on: the shared
// executor header that gets baked into every agent definition body.
var agentFixtureFS = fstest.MapFS{
	"asdt-core/executor-header.md": &fstest.MapFile{Data: []byte(testExecutorHeader)},
}

// headerlessFS models a partial fixture without the executor header: agent
// generation must be a silent no-op, not an error (production presence is
// guarded by skill/embedded_test.go).
var headerlessFS = fstest.MapFS{
	"asdt-developer/SKILL.md": &fstest.MapFile{Data: []byte("# Developer")},
}

// codegraphToolsSuffix is the exact suffix appended to every Claude tools line
// when codegraph is detected. Hardcoded (NOT derived from codegraphClaudeTools)
// so any drift in the production var is caught here.
const codegraphToolsSuffix = ", mcp__codegraph__codegraph_explore, mcp__codegraph__codegraph_search, mcp__codegraph__codegraph_callers, mcp__codegraph__codegraph_callees, mcp__codegraph__codegraph_impact, mcp__codegraph__codegraph_node, mcp__codegraph__codegraph_files, mcp__codegraph__codegraph_status"

const analystBaseToolsLine = "tools: Read, Glob, Grep, Bash, mcp__plugin_engram_engram__mem_save, mcp__plugin_engram_engram__mem_search, mcp__plugin_engram_engram__mem_get_observation, mcp__engram__mem_save, mcp__engram__mem_search, mcp__engram__mem_get_observation"

const builderBaseToolsLine = "tools: Read, Glob, Grep, Bash, Edit, Write, mcp__plugin_engram_engram__mem_save, mcp__plugin_engram_engram__mem_search, mcp__plugin_engram_engram__mem_get_observation, mcp__plugin_engram_engram__mem_update, mcp__engram__mem_save, mcp__engram__mem_search, mcp__engram__mem_get_observation, mcp__engram__mem_update"

// codeIntelligenceHeading is the H2 that marks the gated codegraph guidance
// section. Hardcoded (NOT derived from codeIntelligenceGuidance) so any drift
// in the production const is caught here.
const codeIntelligenceHeading = "## Code intelligence"

// codeIntelligenceSubstrings are load-bearing markers of the gated guidance
// body, hardcoded independently of the production const: each must appear
// verbatim when codegraph is detected.
var codeIntelligenceSubstrings = []string{
	"## Code intelligence",
	"codegraph_explore",
	"codegraph_callers",
	"codegraph_impact",
	"schemas like Rails db/schema.rb, configs, locales, generated DSL",
}

func TestGenerateClaudeAgents_ExactFrontmatterAndBody(t *testing.T) {
	branches := []struct {
		name             string
		opts             InstallOptions
		analystToolsLine string
		builderToolsLine string
		wantGuidance     bool
	}{
		{
			name:             "codegraph absent keeps prior tools lines byte-identical",
			opts:             InstallOptions{CodegraphFound: false},
			analystToolsLine: analystBaseToolsLine,
			builderToolsLine: builderBaseToolsLine,
			wantGuidance:     false,
		},
		{
			name:             "codegraph found appends the read-only codegraph grants",
			opts:             InstallOptions{CodegraphFound: true},
			analystToolsLine: analystBaseToolsLine + codegraphToolsSuffix,
			builderToolsLine: builderBaseToolsLine + codegraphToolsSuffix,
			wantGuidance:     true,
		},
	}

	for _, branch := range branches {
		t.Run(branch.name, func(t *testing.T) {
			cases := []struct {
				name          string
				file          string
				wantNameLine  string
				wantToolsLine string
			}{
				{
					name:          "analyst",
					file:          "asdt-analyst.md",
					wantNameLine:  "name: asdt-analyst",
					wantToolsLine: branch.analystToolsLine,
				},
				{
					name:          "builder",
					file:          "asdt-builder.md",
					wantNameLine:  "name: asdt-builder",
					wantToolsLine: branch.builderToolsLine,
				},
			}

			agentRoot := filepath.Join(t.TempDir(), "agents")
			written, err := generateClaudeAgents(agentFixtureFS, agentRoot, branch.opts)
			if err != nil {
				t.Fatalf("generateClaudeAgents returned error: %v", err)
			}
			if len(written) != len(cases) {
				t.Fatalf("expected %d written agent files, got %d: %v", len(cases), len(written), written)
			}

			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					data, readErr := os.ReadFile(filepath.Join(agentRoot, c.file))
					if readErr != nil {
						t.Fatalf("expected agent file %q to exist: %v", c.file, readErr)
					}
					content := string(data)

					if !strings.Contains(content, c.wantNameLine+"\n") {
						t.Errorf("agent %q missing exact name line %q\ngot:\n%s", c.file, c.wantNameLine, content)
					}
					if !strings.Contains(content, c.wantToolsLine+"\n") {
						t.Errorf("agent %q missing EXACT tools line\nwant: %s\ngot:\n%s", c.file, c.wantToolsLine, content)
					}
					assertAgentBody(t, c.file, content, branch.wantGuidance)
					assertNoDelegationTools(t, c.file, c.wantToolsLine)
				})
			}
		})
	}
}

// TestGenerateClaudeAgents_CodegraphDoesNotMutateSpecs guards the shared
// agentTypeSpecs global: a codegraph-enabled generation must not leak the
// appended grants into a later codegraph-absent generation.
func TestGenerateClaudeAgents_CodegraphDoesNotMutateSpecs(t *testing.T) {
	withRoot := filepath.Join(t.TempDir(), "agents-with")
	if _, err := generateClaudeAgents(agentFixtureFS, withRoot, InstallOptions{CodegraphFound: true}); err != nil {
		t.Fatalf("codegraph-enabled run: %v", err)
	}

	withoutRoot := filepath.Join(t.TempDir(), "agents-without")
	if _, err := generateClaudeAgents(agentFixtureFS, withoutRoot, InstallOptions{}); err != nil {
		t.Fatalf("codegraph-absent run: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(withoutRoot, "asdt-analyst.md"))
	if err != nil {
		t.Fatalf("read analyst file: %v", err)
	}
	if strings.Contains(string(data), "mcp__codegraph__") {
		t.Error("codegraph-absent generation carries codegraph grants — agentTypeSpecs was mutated by the earlier codegraph-enabled run")
	}
	if strings.Contains(string(data), codeIntelligenceHeading) {
		t.Error("codegraph-absent generation carries the codegraph guidance section — agentTypeSpecs.Guidance was mutated by the earlier codegraph-enabled run")
	}
}

// assertAgentBody verifies the body carries the executor header VERBATIM,
// followed by the constraints prose, in that order. When wantGuidance is true
// it further requires the gated codegraph section to appear AFTER constraints
// and to carry every drift marker; when false it requires the section absent.
func assertAgentBody(t *testing.T, file, content string, wantGuidance bool) {
	t.Helper()

	headerIdx := strings.Index(content, testExecutorHeader)
	if headerIdx == -1 {
		t.Errorf("agent %q body does not contain the executor header verbatim", file)
		return
	}
	constraintsIdx := strings.Index(content, "constraints")
	if constraintsIdx == -1 {
		t.Errorf("agent %q body does not contain constraint prose", file)
		return
	}
	if constraintsIdx < headerIdx {
		t.Errorf("agent %q constraints appear before the executor header; header must come first", file)
	}

	headingIdx := strings.Index(content, codeIntelligenceHeading)
	if !wantGuidance {
		if headingIdx != -1 {
			t.Errorf("agent %q body carries the codegraph guidance section, want it absent", file)
		}
		return
	}

	if headingIdx == -1 {
		t.Errorf("agent %q body does not contain the codegraph guidance heading %q", file, codeIntelligenceHeading)
		return
	}
	if headingIdx < constraintsIdx {
		t.Errorf("agent %q guidance section appears before the constraints; constraints must come first", file)
	}
	for _, want := range codeIntelligenceSubstrings {
		if !strings.Contains(content, want) {
			t.Errorf("agent %q guidance section missing expected substring %q", file, want)
		}
	}
}

// assertNoDelegationTools verifies neither Agent nor Task appears as a tool
// list element — executor agents must never delegate.
func assertNoDelegationTools(t *testing.T, file, toolsLine string) {
	t.Helper()

	list := strings.TrimPrefix(toolsLine, "tools: ")
	for _, tool := range strings.Split(list, ", ") {
		if tool == "Agent" || tool == "Task" {
			t.Errorf("agent %q tools list contains delegation tool %q", file, tool)
		}
	}
}

// openCodeAgentFrontmatter is the structural shape asserted for generated
// OpenCode agent definitions. Bash is `any` because it is a scalar ("allow")
// for builder and a glob->verdict map for analyst.
type openCodeAgentFrontmatter struct {
	Description string `yaml:"description"`
	Mode        string `yaml:"mode"`
	Permission  struct {
		Edit string `yaml:"edit"`
		Task string `yaml:"task"`
		Bash any    `yaml:"bash"`
	} `yaml:"permission"`
}

func parseOpenCodeAgentFile(t *testing.T, path string) (openCodeAgentFrontmatter, string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	content := string(data)

	const delim = "---\n"
	if !strings.HasPrefix(content, delim) {
		t.Fatalf("agent file %q does not start with a frontmatter delimiter", path)
	}
	rest := content[len(delim):]
	end := strings.Index(rest, "\n"+delim)
	if end == -1 {
		t.Fatalf("agent file %q has no closing frontmatter delimiter", path)
	}

	var fm openCodeAgentFrontmatter
	if err := yaml.Unmarshal([]byte(rest[:end]), &fm); err != nil {
		t.Fatalf("agent file %q frontmatter is not valid YAML: %v", path, err)
	}
	body := rest[end+1+len(delim):]
	return fm, body
}

func TestGenerateOpenCodeAgents_StructuralPermissions(t *testing.T) {
	agentRoot := filepath.Join(t.TempDir(), "agents")
	written, err := generateOpenCodeAgents(agentFixtureFS, agentRoot, InstallOptions{})
	if err != nil {
		t.Fatalf("generateOpenCodeAgents returned error: %v", err)
	}
	if len(written) != 2 {
		t.Fatalf("expected 2 written agent files, got %d: %v", len(written), written)
	}

	t.Run("analyst", func(t *testing.T) {
		fm, body := parseOpenCodeAgentFile(t, filepath.Join(agentRoot, "asdt-analyst.md"))
		if fm.Mode != "subagent" {
			t.Errorf("analyst mode = %q, want %q", fm.Mode, "subagent")
		}
		if fm.Description == "" {
			t.Error("analyst description is empty")
		}
		if fm.Permission.Edit != "deny" {
			t.Errorf("analyst permission.edit = %q, want %q", fm.Permission.Edit, "deny")
		}
		if fm.Permission.Task != "deny" {
			t.Errorf("analyst permission.task = %q, want %q", fm.Permission.Task, "deny")
		}

		bashMap, ok := fm.Permission.Bash.(map[string]any)
		if !ok {
			t.Fatalf("analyst permission.bash = %T (%v), want a glob->verdict map", fm.Permission.Bash, fm.Permission.Bash)
		}
		if got := bashMap["*"]; got != "deny" {
			t.Errorf("analyst permission.bash[%q] = %v, want %q", "*", got, "deny")
		}
		for _, glob := range analystBashAllowlist {
			if got := bashMap[glob]; got != "allow" {
				t.Errorf("analyst permission.bash[%q] = %v, want %q", glob, got, "allow")
			}
		}
		if len(bashMap) != len(analystBashAllowlist)+1 {
			t.Errorf("analyst permission.bash has %d entries, want %d (whitelist + catch-all deny)", len(bashMap), len(analystBashAllowlist)+1)
		}
		assertAgentBody(t, "asdt-analyst.md", body, false)
	})

	t.Run("builder", func(t *testing.T) {
		fm, body := parseOpenCodeAgentFile(t, filepath.Join(agentRoot, "asdt-builder.md"))
		if fm.Mode != "subagent" {
			t.Errorf("builder mode = %q, want %q", fm.Mode, "subagent")
		}
		if fm.Permission.Edit != "allow" {
			t.Errorf("builder permission.edit = %q, want %q", fm.Permission.Edit, "allow")
		}
		if fm.Permission.Task != "deny" {
			t.Errorf("builder permission.task = %q, want %q", fm.Permission.Task, "deny")
		}
		if bash, ok := fm.Permission.Bash.(string); !ok || bash != "allow" {
			t.Errorf("builder permission.bash = %v (%T), want scalar %q", fm.Permission.Bash, fm.Permission.Bash, "allow")
		}
		assertAgentBody(t, "asdt-builder.md", body, false)
	})
}

func TestGenerateAgents_IdempotentDoubleGenerate(t *testing.T) {
	generators := []struct {
		name     string
		generate func(skillsFS fstest.MapFS, agentRoot string) ([]string, error)
	}{
		{name: "claude", generate: func(fsys fstest.MapFS, root string) ([]string, error) {
			return generateClaudeAgents(fsys, root, InstallOptions{})
		}},
		{name: "opencode", generate: func(fsys fstest.MapFS, root string) ([]string, error) {
			return generateOpenCodeAgents(fsys, root, InstallOptions{})
		}},
	}

	for _, g := range generators {
		t.Run(g.name, func(t *testing.T) {
			agentRoot := filepath.Join(t.TempDir(), "agents")

			written1, err1 := g.generate(agentFixtureFS, agentRoot)
			if err1 != nil {
				t.Fatalf("first run: %v", err1)
			}
			first := readWrapperContents(t, written1)

			written2, err2 := g.generate(agentFixtureFS, agentRoot)
			if err2 != nil {
				t.Fatalf("second run: %v", err2)
			}
			if len(written1) != len(written2) {
				t.Fatalf("written counts differ across runs: first=%d, second=%d", len(written1), len(written2))
			}

			second := readWrapperContents(t, written2)
			for path, want := range first {
				if got := second[path]; got != want {
					t.Errorf("agent file %q changed across runs (not byte-identical)", path)
				}
			}
		})
	}
}

// TestGenerateOpenCodeAgents_CodegraphStructuralNeutrality asserts codegraph
// detection changes OpenCode output ONLY by appending the gated guidance
// section — never its structure. This deliberately overturns the prior
// byte-identity premise (an ADR-follow-up decision by the project owner): the
// codegraph tool grant stays Claude-only, so the frontmatter must be identical
// across branches, and the absent body must be an exact prefix of the with body
// (the sole difference is the appended "## Code intelligence" section).
func TestGenerateOpenCodeAgents_CodegraphStructuralNeutrality(t *testing.T) {
	withoutRoot := filepath.Join(t.TempDir(), "agents-without")
	if _, err := generateOpenCodeAgents(agentFixtureFS, withoutRoot, InstallOptions{CodegraphFound: false}); err != nil {
		t.Fatalf("codegraph-absent run: %v", err)
	}

	withRoot := filepath.Join(t.TempDir(), "agents-with")
	if _, err := generateOpenCodeAgents(agentFixtureFS, withRoot, InstallOptions{CodegraphFound: true}); err != nil {
		t.Fatalf("codegraph-enabled run: %v", err)
	}

	for _, file := range []string{"asdt-analyst.md", "asdt-builder.md"} {
		t.Run(file, func(t *testing.T) {
			fmWithout, bodyWithout := parseOpenCodeAgentFile(t, filepath.Join(withoutRoot, file))
			fmWith, bodyWith := parseOpenCodeAgentFile(t, filepath.Join(withRoot, file))

			if !reflect.DeepEqual(fmWithout, fmWith) {
				t.Errorf("agent %q frontmatter differs between codegraph branches; the grant is Claude-only, structure must be identical\nwithout: %+v\nwith:    %+v", file, fmWithout, fmWith)
			}

			if strings.Contains(bodyWithout, codeIntelligenceHeading) {
				t.Errorf("agent %q codegraph-absent body carries the guidance heading, want it absent", file)
			}
			if !strings.Contains(bodyWith, codeIntelligenceHeading) {
				t.Errorf("agent %q codegraph-enabled body is missing the guidance heading %q", file, codeIntelligenceHeading)
			}
			for _, want := range codeIntelligenceSubstrings {
				if !strings.Contains(bodyWith, want) {
					t.Errorf("agent %q codegraph-enabled body missing expected substring %q", file, want)
				}
			}

			if !strings.HasPrefix(bodyWith, bodyWithout) {
				t.Errorf("agent %q codegraph-enabled body is not a superset of the absent body; the only difference must be the appended guidance section", file)
			}
		})
	}
}

func TestGenerateAgentFiles_MissingExecutorHeaderIsNoOp(t *testing.T) {
	agentRoot := filepath.Join(t.TempDir(), "agents")

	written, err := generateClaudeAgents(headerlessFS, agentRoot, InstallOptions{})
	if err != nil {
		t.Fatalf("expected nil error for a fixture FS without the executor header, got: %v", err)
	}
	if written != nil {
		t.Errorf("expected no written files, got %v", written)
	}
	checkFileAbsent(t, agentRoot)
}

// TestExecutorHeaderCarriesEvidenceCondition reads the REAL embedded executor
// header — not the trimmed agentFixtureFS copy — because the evidence clause's
// condition is a one-way-door contract: it is what exempts steps that never
// touched the codebase from the verifiable-evidence requirement. The sentence is
// asserted verbatim so any reword that widens or narrows the exemption fails
// loudly here instead of shipping silently baked into both agent definitions.
//
// The condition was rekeyed by the refactor: it used to exempt steps without an
// `output_topic_key`, which stopped being the right axis once steps could
// declare `output: context`. It now turns on whether the step read the codebase.
func TestExecutorHeaderCarriesEvidenceCondition(t *testing.T) {
	data, err := fs.ReadFile(skill.FS(), "asdt-core/executor-header.md")
	if err != nil {
		t.Fatalf("read embedded executor-header.md: %v", err)
	}

	const condition = "if your step read the codebase, anchor every claim to something"
	if !strings.Contains(string(data), condition) {
		t.Errorf("embedded executor-header.md is missing the verbatim evidence condition sentence %q", condition)
	}
}

func TestAgentAdapters_BothAssistantsRegistered(t *testing.T) {
	for _, id := range []AssistantID{AssistantClaudeCode, AssistantOpenCode} {
		adapter, ok := agentAdapterFor(id)
		if !ok {
			t.Errorf("agentAdapterFor(%q) not found; both assistants must carry an agent adapter", id)
			continue
		}
		if adapter.Generate == nil {
			t.Errorf("agentAdapterFor(%q).Generate is nil, want non-nil", id)
		}
	}
}

func TestAgentRootFor(t *testing.T) {
	cases := []struct {
		name           string
		id             AssistantID
		wantNonEmpty   bool
		wantPathSuffix string
	}{
		{name: "claude code resolves to its agents directory", id: AssistantClaudeCode, wantNonEmpty: true, wantPathSuffix: "/.claude/agents"},
		{name: "opencode resolves to its agents directory", id: AssistantOpenCode, wantNonEmpty: true, wantPathSuffix: "/opencode/agents"},
		{name: "unknown assistant has no agent root", id: "a1", wantNonEmpty: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := agentRootFor(c.id)
			if c.wantNonEmpty && got == "" {
				t.Fatalf("agentRootFor(%q) = %q, want a non-empty path", c.id, got)
			}
			if !c.wantNonEmpty && got != "" {
				t.Fatalf("agentRootFor(%q) = %q, want \"\" (no known agent root)", c.id, got)
			}
			if c.wantPathSuffix != "" && !strings.HasSuffix(got, c.wantPathSuffix) {
				t.Errorf("agentRootFor(%q) = %q, want it to end in %q", c.id, got, c.wantPathSuffix)
			}
		})
	}
}

func TestAgentTypeSpecs_MatchCanonicalNames(t *testing.T) {
	if len(agentTypeSpecs) != len(AgentTypeNames) {
		t.Fatalf("agentTypeSpecs has %d entries, AgentTypeNames has %d — they must stay in lockstep", len(agentTypeSpecs), len(AgentTypeNames))
	}
	for i, spec := range agentTypeSpecs {
		if spec.ID != AgentTypeNames[i] {
			t.Errorf("agentTypeSpecs[%d].ID = %q, want %q", i, spec.ID, AgentTypeNames[i])
		}
	}
}
