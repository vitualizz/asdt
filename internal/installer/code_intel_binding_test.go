package installer

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/vitualizz/asdt/skill"
)

// fakeCodeIntelProvider is a test-only second code-intelligence provider that
// binds only two verbs under two Claude Code prefixes. Installing with it
// proves a new provider is a Go-only change: the embedded prompts bind to it
// with no edit, the binding lists only the verbs it serves, and nothing of
// Codegraph's survives into the installed tree.
var fakeCodeIntelProvider = CodeIntelDescriptor{
	ID:          "acmeidx",
	Name:        "Acme Index",
	Description: "Test-only code-intelligence provider.",
	ConfigValue: "acme-index",
	IndexMarker: ".acme-index",
	Tools: map[CodeIntelVerb]ProviderTool{
		CodeIntelExplore: {Name: "acme_explore", Usage: "pass the question as `q`"},
		CodeIntelCallers: {Name: "acme_refs", Usage: "pass the symbol as `sym`"},
	},
	ClaudeToolPrefixes: []string{"mcp__acmeidx__", "mcp__plugin_acme_idx__"},
	// No Detect.Binary: a hosted provider with nothing to probe locally.
}

// codeIntelBoundFiles are the installed files (SkillsDir-relative) that carry
// a code-intel-binding region: the executor header every sub-agent holds and
// asdt-init/SKILL.md, whose enrichment calls explore.
var codeIntelBoundFiles = []string{
	"asdt-core/executor-header.md",
	"asdt-init/SKILL.md",
}

// codeIntelRegion returns the body strictly between the code-intel-binding
// markers of content, failing when the region is absent.
func codeIntelRegion(t *testing.T, label, content string) string {
	t.Helper()
	begin := strings.Index(content, codeIntelBindingBeginMarker)
	end := strings.Index(content, codeIntelBindingEndMarker)
	if begin < 0 || end < begin {
		t.Fatalf("%s has no code-intel-binding region", label)
	}
	return content[begin+len(codeIntelBindingBeginMarker) : end]
}

func TestCodeIntelProviders_Validate(t *testing.T) {
	for _, p := range append([]CodeIntelDescriptor{fakeCodeIntelProvider}, CodeIntelProviders...) {
		if err := p.Validate(); err != nil {
			t.Errorf("code-intelligence provider %q does not validate: %v", p.ID, err)
		}
	}
}

// TestProviders_UniqueAcrossKinds guards the keys the TUI looks providers up
// by: a preflight row is labeled with the provider's Name whatever its kind,
// and each registry is keyed by ID, so neither may repeat — and a
// code-intelligence ConfigValue or IndexMarker may not repeat either.
func TestProviders_UniqueAcrossKinds(t *testing.T) {
	names := map[string]string{}
	for _, p := range Providers {
		names[p.Name] = "memory " + p.ID
	}
	ids, configs, markers := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, p := range CodeIntelProviders {
		if prev, dup := names[p.Name]; dup {
			t.Errorf("code-intelligence provider %q shares Name %q with %s", p.ID, p.Name, prev)
		}
		names[p.Name] = "code intelligence " + p.ID
		for _, f := range []struct {
			value string
			seen  map[string]bool
		}{{p.ID, ids}, {p.ConfigValue, configs}, {p.IndexMarker, markers}} {
			if f.seen[f.value] {
				t.Errorf("code-intelligence provider %q repeats %q", p.ID, f.value)
			}
			f.seen[f.value] = true
		}
	}
}

func TestBindCodeIntel(t *testing.T) {
	const begin, end = codeIntelBindingBeginMarker, codeIntelBindingEndMarker
	noExplore := fakeCodeIntelProvider
	noExplore.Tools = map[CodeIntelVerb]ProviderTool{CodeIntelCallers: fakeCodeIntelProvider.Tools[CodeIntelCallers]}
	unknownVerb := fakeCodeIntelProvider
	unknownVerb.Tools = map[CodeIntelVerb]ProviderTool{CodeIntelExplore: fakeCodeIntelProvider.Tools[CodeIntelExplore], "rename": {Name: "acme_mv", Usage: "x"}}
	noUsage := fakeCodeIntelProvider
	noUsage.Tools = map[CodeIntelVerb]ProviderTool{CodeIntelExplore: {Name: "acme_explore"}}
	noMarker := fakeCodeIntelProvider
	noMarker.IndexMarker = ""
	pathMarker := fakeCodeIntelProvider
	pathMarker.IndexMarker = "build/.acme-index"
	configNone := fakeCodeIntelProvider
	configNone.ConfigValue = CodeIntelNone
	noPrefixes := fakeCodeIntelProvider
	noPrefixes.ClaudeToolPrefixes = nil

	cases := []struct {
		name     string
		content  string
		provider CodeIntelDescriptor
		wantErr  string   // substring; "" = success
		wantSame bool     // content must pass through unchanged
		want     []string // substrings the bound content must carry
		lack     []string // substrings it must NOT carry
	}{
		{name: "marker-free content passes through, even with no provider", content: "# plain\n", wantSame: true},
		{name: "no provider binds the none body", content: begin + end, want: []string{"**No code-intelligence provider bound.**", "grep and read"}, lack: []string{"Code-intelligence provider:"}},
		{
			name: "codegraph binds every verb in order", content: begin + end, provider: CodeIntelProviders[0],
			want: []string{
				"**Code-intelligence provider: Codegraph.** Its index marker is `.codegraph` at the repository root.",
				"`code_intelligence: codegraph`",
				"`Tooling: codegraph index available`",
				"- **explore** → `codegraph_explore` — ask the question in plain words; one call returns the relevant source grouped by file.\n- **search** → `codegraph_search`",
				"- **read** → `codegraph_node`",
				"- **impact** → `codegraph_impact`",
				"`mcp__codegraph__codegraph_explore`",
			},
		},
		{
			name: "a partial provider lists only the verbs it serves", content: begin + end, provider: fakeCodeIntelProvider,
			want: []string{"- **explore** → `acme_explore` — pass the question as `q`.\n- **callers** → `acme_refs` — pass the symbol as `sym`.\n", "`mcp__acmeidx__acme_explore` or `mcp__plugin_acme_idx__acme_explore`"},
			lack: []string{"**search**", "**impact**"},
		},
		{name: "provider without prefixes binds without the prefix sentence", content: begin + end, provider: noPrefixes, lack: []string{"A host may prefix"}},
		{name: "explore is required", content: begin + end, provider: noExplore, wantErr: `binds no tool to verb "explore"`},
		{name: "unknown verb fails", content: begin + end, provider: unknownVerb, wantErr: `unknown verb "rename"`},
		{name: "tool without Usage fails", content: begin + end, provider: noUsage, wantErr: `gives verb "explore" tool "acme_explore" no Usage`},
		{name: "missing index marker fails", content: begin + end, provider: noMarker, wantErr: "IndexMarker"},
		{name: "index marker path fails", content: begin + end, provider: pathMarker, wantErr: "one repository-root entry name"},
		{name: "ConfigValue none fails", content: begin + end, provider: configNone, wantErr: "which means no provider"},
		{name: "duplicated region fails", content: begin + end + begin + end, provider: fakeCodeIntelProvider, wantErr: "appears 2 times"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := bindCodeIntel([]byte(c.content), c.provider)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("bindCodeIntel error = %v, want one containing %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("bindCodeIntel: %v", err)
			}
			if c.wantSame {
				if string(got) != c.content {
					t.Errorf("bindCodeIntel changed marker-free content: %q", got)
				}
				return
			}
			again, err := bindCodeIntel(got, c.provider)
			if err != nil || string(again) != string(got) {
				t.Errorf("bindCodeIntel is not idempotent: err=%v\nfirst:  %q\nsecond: %q", err, got, again)
			}
			for _, w := range c.want {
				if !strings.Contains(string(got), w) {
					t.Errorf("bound content missing %q:\n%s", w, got)
				}
			}
			for _, l := range c.lack {
				if strings.Contains(string(got), l) {
					t.Errorf("bound content carries %q:\n%s", l, got)
				}
			}
		})
	}
}

// TestInstall_CodeIntelBindingReachesInstalledFiles installs the real skill
// tree bound to Codegraph, to a second provider, and to none, and asserts each
// code-intel-binding region carries exactly that binding, the install metadata
// records it, and — for the second provider and none — no installed file names
// a shipped provider anywhere outside its regions.
func TestInstall_CodeIntelBindingReachesInstalledFiles(t *testing.T) {
	for _, c := range []struct {
		name     string
		provider CodeIntelDescriptor
		neutral  bool
	}{
		{name: "codegraph", provider: CodeIntelProviders[0]},
		{name: "a second provider binds with no prompt edits", provider: fakeCodeIntelProvider, neutral: true},
		{name: "no code intelligence", neutral: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			skillsDir := filepath.Join(t.TempDir(), "skills")
			a := AssistantDescriptor{ID: "code-intel-test", Name: "Code Intel Test", BinaryName: "sh", SkillsDir: skillsDir}
			results := InstallWithModels([]AssistantDescriptor{a}, ProviderDescriptor{}, skill.FS(), "", nil, InstallOptions{CodeIntel: c.provider})
			if results[0].Err != nil {
				t.Fatalf("install: %v", results[0].Err)
			}
			meta, err := ReadInstallMeta(a)
			if err != nil {
				t.Fatalf("read install meta: %v", err)
			}
			if want := c.provider.MetaValue(); meta.CodeIntel != want {
				t.Errorf("install meta CodeIntel = %q, want %q", meta.CodeIntel, want)
			}

			want, err := renderCodeIntelBinding(c.provider)
			if err != nil {
				t.Fatalf("renderCodeIntelBinding: %v", err)
			}
			for _, rel := range codeIntelBoundFiles {
				data, err := os.ReadFile(filepath.Join(skillsDir, filepath.FromSlash(rel)))
				if err != nil {
					t.Fatalf("read %s: %v", rel, err)
				}
				if got := codeIntelRegion(t, rel, string(data)); got != want {
					t.Errorf("%s code-intel-binding region:\n%s\nwant:\n%s", rel, got, want)
				}
			}

			if !c.neutral {
				return
			}
			walkErr := filepath.WalkDir(skillsDir, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				data, readErr := os.ReadFile(path)
				if readErr != nil {
					return readErr
				}
				for _, hit := range providerNameHits(string(data), shippedProviderNames(), nil) {
					t.Errorf("installed %s carries %q — a prompt still hardcodes a shipped provider", path, hit)
				}
				return nil
			})
			if walkErr != nil {
				t.Fatalf("walk %s: %v", skillsDir, walkErr)
			}
		})
	}
}

// TestGenerateAgents_CodeIntelBinding generates both assistants' executor
// agents from the REAL executor header under each provider and asserts the
// provider changes only what it should: the body's code-intel-binding region
// (and, on Claude Code, the tools line). The OpenCode frontmatter is identical
// across providers — OpenCode hands MCP tools to subagents implicitly — and its
// memory-tool denies never catch a code-intelligence tool.
func TestGenerateAgents_CodeIntelBinding(t *testing.T) {
	providers := []CodeIntelDescriptor{{}, CodeIntelProviders[0], fakeCodeIntelProvider}

	for name, generate := range map[string]func(string, InstallOptions) ([]string, error){
		"claude": func(root string, o InstallOptions) ([]string, error) {
			return generateClaudeAgents(skill.FS(), root, o)
		},
		"opencode": func(root string, o InstallOptions) ([]string, error) {
			return generateOpenCodeAgents(skill.FS(), root, o)
		},
	} {
		t.Run(name, func(t *testing.T) {
			outside := map[string]string{} // file → body outside the region, from the first provider
			for _, p := range providers {
				root := filepath.Join(t.TempDir(), "agents")
				if _, err := generate(root, InstallOptions{CodeIntel: p}); err != nil {
					t.Fatalf("generate with %q: %v", p.MetaValue(), err)
				}
				want, err := renderCodeIntelBinding(p)
				if err != nil {
					t.Fatalf("renderCodeIntelBinding: %v", err)
				}
				for _, file := range []string{"asdt-analyst.md", "asdt-builder.md"} {
					data, err := os.ReadFile(filepath.Join(root, file))
					if err != nil {
						t.Fatalf("read %s: %v", file, err)
					}
					content := string(data)
					if got := codeIntelRegion(t, file, content); got != want {
						t.Errorf("%s with %q: region\n%s\nwant\n%s", file, p.MetaValue(), got, want)
					}
					rest := stripRegions(content, codeIntelBindingBeginMarker, codeIntelBindingEndMarker)
					if name == "claude" {
						// The tools line is TestGenerateClaudeAgents_ExactFrontmatterAndBody's.
						rest = stripFrontmatter(rest)
					} else {
						fm, _ := parseOpenCodeAgentFile(t, filepath.Join(root, file))
						for key := range fm.Permission.Tools {
							for _, tool := range p.Tools {
								if key == "*_"+tool.Name {
									t.Errorf("OpenCode %s denies %q, a %s code-intelligence tool", file, key, p.Name)
								}
							}
						}
					}
					if prev, ok := outside[file]; !ok {
						outside[file] = rest
					} else if !reflect.DeepEqual(prev, rest) {
						t.Errorf("%s with %q differs outside the code-intel-binding region; a provider is a Go-only change", file, p.MetaValue())
					}
				}
			}
		})
	}
}
