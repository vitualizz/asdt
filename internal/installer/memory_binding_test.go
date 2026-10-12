package installer

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vitualizz/asdt/skill"
)

// fakeMemoryProvider is a test-only second provider. Installing with it proves
// a new provider is a Go-only change: the embedded prompts bind to it with no
// edit, and nothing of Engram's survives into the installed tree.
var fakeMemoryProvider = ProviderDescriptor{
	ID:          "acme",
	Name:        "Acme Memory",
	Description: "Test-only provider.",
	ConfigValue: "acme",
	Tools: map[MemoryVerb]MemoryTool{
		MemorySave:   {Name: "acme_put", Usage: "pass the key as `slot`"},
		MemorySearch: {Name: "acme_find", Usage: "pass the query as `q`; each match carries its key as `slot`"},
		MemoryGet:    {Name: "acme_fetch", Usage: "pass the id as `ref`"},
	},
	ClaudeToolPrefixes: []string{"mcp__acme__"},
	// No Detect.Binary: a hosted provider with nothing to probe locally.
}

// boundFiles are the installed files (SkillsDir-relative) that must carry a
// bound memory-binding region: the four that own one, plus every routed
// SKILL.md, which receives protocol.md's copy through the specialist-header
// splice.
var boundFiles = []string{
	"asdt-core/protocol.md",
	"asdt-core/executor-header.md",
	"asdt-init/SKILL.md",
	"asdt/SKILL.md",
	"asdt-pm/SKILL.md",
	"asdt-architect/SKILL.md",
	"asdt-qa/SKILL.md",
	"asdt-security/SKILL.md",
	"asdt-ux-ui/SKILL.md",
	"asdt-developer/SKILL.md",
	"asdt-researcher/SKILL.md",
}

// installWithProvider installs the real embedded skill tree for one throwaway
// assistant bound to provider and returns its SkillsDir.
func installWithProvider(t *testing.T, provider ProviderDescriptor) string {
	t.Helper()

	skillsDir := filepath.Join(t.TempDir(), "skills")
	a := AssistantDescriptor{ID: "memory-binding-test", Name: "Memory Binding Test", BinaryName: "sh", SkillsDir: skillsDir}
	results := InstallWithModels([]AssistantDescriptor{a}, provider, skill.FS(), "", nil, InstallOptions{})
	if results[0].Err != nil {
		t.Fatalf("install with provider %q: %v", provider.ID, results[0].Err)
	}
	return skillsDir
}

// bindingRegion returns the body strictly between the memory-binding markers
// of an installed file, failing when the region is absent.
func bindingRegion(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	content := string(data)
	begin := strings.Index(content, memoryBindingBeginMarker)
	end := strings.Index(content, memoryBindingEndMarker)
	if begin < 0 || end < begin {
		t.Fatalf("%s has no memory-binding region", path)
	}
	return content[begin+len(memoryBindingBeginMarker) : end]
}

func TestInstall_MemoryBindingReachesInstalledFiles(t *testing.T) {
	cases := []struct {
		name     string
		provider ProviderDescriptor
		want     []string // substrings every bound region must carry
		neutral  bool     // no installed file may name a shipped provider, even inside a region
	}{
		{
			name:     "engram binds today's tool names",
			provider: Providers[0],
			want: []string{
				"**Memory provider: Engram.**",
				"`memory.provider: engram`",
				"- **save** → `mem_save` — pass the key as `topic_key`",
				"`type: \"decision\"`",
				"- **search** → `mem_search`",
				"- **get** → `mem_get_observation`",
				"`mcp__plugin_engram_engram__mem_save` or `mcp__engram__mem_save`",
			},
		},
		{
			name:     "a second provider binds with no prompt edits",
			provider: fakeMemoryProvider,
			want: []string{
				"**Memory provider: Acme Memory.**",
				"`memory.provider: acme`",
				"- **save** → `acme_put` — pass the key as `slot`.",
				"- **search** → `acme_find` — pass the query as `q`; each match carries its key as `slot`.",
				"- **get** → `acme_fetch` — pass the id as `ref`.",
				"`mcp__acme__acme_put`",
			},
			neutral: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			skillsDir := installWithProvider(t, c.provider)

			for _, rel := range boundFiles {
				region := bindingRegion(t, filepath.Join(skillsDir, filepath.FromSlash(rel)))
				for _, w := range c.want {
					if !strings.Contains(region, w) {
						t.Errorf("%s binding region missing %q\nregion:%s", rel, w, region)
					}
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
				for _, hit := range providerNameHits(string(data), shippedProviderNames(), c.provider.ClaudeToolPrefixes) {
					t.Errorf("installed %s carries %q — a prompt still hardcodes a shipped provider's binding", path, hit)
				}
				return nil
			})
			if walkErr != nil {
				t.Fatalf("walk %s: %v", skillsDir, walkErr)
			}
		})
	}
}

// shippedProviderNames is every name a shipped provider is known by — its
// Name, probe binary, and tool names — plus `mem_`, Engram's whole tool family
// (mem_update, mem_context…), none of which the prompts may reach for. Derived
// from Providers, so a new provider's names are forbidden in the same commit
// that adds its descriptor.
func shippedProviderNames() []string {
	names := []string{"mem_"}
	seen := map[string]bool{"mem_": true}
	add := func(name string) {
		if name = strings.ToLower(name); name != "" && !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	for _, p := range Providers {
		add(p.Name)
		add(p.Detect.Binary)
		for _, verb := range memoryVerbs {
			add(p.Tools[verb].Name)
		}
	}
	return names
}

// codegraphMCPPrefix is the one MCP namespace prompt prose may name: codegraph
// is not a memory provider, and /asdt-init detects it by this prefix.
const codegraphMCPPrefix = "mcp__codegraph__"

// providerNameHits returns what content names outside its memory-binding
// regions that only a binding may name: a forbidden substring (lowercase,
// matched case-insensitively), or an MCP tool name (`mcp__…`) whose namespace is
// neither codegraph's nor in allowedMCP.
func providerNameHits(content string, forbidden, allowedMCP []string) []string {
	for {
		begin := strings.Index(content, memoryBindingBeginMarker)
		end := strings.Index(content, memoryBindingEndMarker)
		if begin < 0 || end < begin {
			break
		}
		content = content[:begin] + content[end+len(memoryBindingEndMarker):]
	}
	lower := strings.ToLower(content)

	var hits []string
	for _, f := range forbidden {
		if strings.Contains(lower, f) {
			hits = append(hits, f)
		}
	}
	allowed := append([]string{codegraphMCPPrefix}, allowedMCP...)
	for rest := lower; ; {
		i := strings.Index(rest, "mcp__")
		if i < 0 {
			break
		}
		rest = rest[i:]
		ok := false
		for _, a := range allowed {
			if strings.HasPrefix(rest, strings.ToLower(a)) {
				ok = true
				break
			}
		}
		if !ok {
			name, _, _ := strings.Cut(rest, "`")
			hits = append(hits, strings.Fields(name)[0])
		}
		rest = rest[len("mcp__"):]
	}
	return hits
}

// TestPromptsAreProviderNeutral asserts no prompt text names a memory provider
// or its tools outside a memory-binding region. Prompts speak the memory
// interface (asdt-core/protocol.md §0); the installer binds it to the selected
// provider, so a provider is added in Go alone. It covers every source of
// prompt text: the embedded skill tree, the executor agents generated for a
// non-shipped provider (Claude Code and OpenCode, codegraph guidance on), the
// Go prompt-prose constants, and the AGENTS.md template and personas.
func TestPromptsAreProviderNeutral(t *testing.T) {
	type source struct {
		name, content string
		allowedMCP    []string // MCP namespaces this source may name: its agents' grants
	}
	var sources []source

	readAll := func(fsys fs.FS, root, label string, allowedMCP []string) {
		walkErr := fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			data, readErr := fs.ReadFile(fsys, p)
			if readErr != nil {
				return readErr
			}
			sources = append(sources, source{label + p, string(data), allowedMCP})
			return nil
		})
		if walkErr != nil {
			t.Fatalf("walk %s: %v", label, walkErr)
		}
	}
	readAll(skill.FS(), ".", "skill/", nil)
	readAll(assetsFS, "assets/personas", "", nil)
	readAll(assetsFS, "assets/agents-template.md", "", nil)

	for name, generate := range map[string]func(fs.FS, string, ProviderDescriptor, InstallOptions) ([]string, error){
		"claude":   generateClaudeAgents,
		"opencode": generateOpenCodeAgents,
	} {
		agentRoot := filepath.Join(t.TempDir(), name)
		if _, err := generate(skill.FS(), agentRoot, fakeMemoryProvider, InstallOptions{CodegraphFound: true}); err != nil {
			t.Fatalf("generate %s agents: %v", name, err)
		}
		readAll(os.DirFS(agentRoot), ".", name+" agent ", fakeMemoryProvider.ClaudeToolPrefixes)
	}

	for name, prose := range map[string]string{
		"analystConstraints":       analystConstraints,
		"builderConstraints":       builderConstraints,
		"codeIntelligenceGuidance": codeIntelligenceGuidance,
	} {
		sources = append(sources, source{"const " + name, prose, nil})
	}

	for _, src := range sources {
		for _, hit := range providerNameHits(src.content, shippedProviderNames(), src.allowedMCP) {
			t.Errorf("%s names %q outside a memory-binding region — speak the memory interface (memory **save** / **search** / **get**) instead", src.name, hit)
		}
	}
}

// TestGenerateAgents_MemoryGrantsFollowProvider asserts the executor agents'
// memory grants and baked-in binding come from the selected provider, against
// the REAL executor header (agentFixtureFS carries no binding region).
func TestGenerateAgents_MemoryGrantsFollowProvider(t *testing.T) {
	cases := []struct {
		name        string
		provider    ProviderDescriptor
		analystLine string
		builderLine string
		bodyWant    string
	}{
		{
			name:        "engram",
			provider:    Providers[0],
			analystLine: analystBaseToolsLine,
			builderLine: builderBaseToolsLine,
			bodyWant:    "- **save** → `mem_save`",
		},
		{
			name:        "fake provider",
			provider:    fakeMemoryProvider,
			analystLine: "tools: Read, Glob, Grep, Bash, mcp__acme__acme_put, mcp__acme__acme_find, mcp__acme__acme_fetch",
			builderLine: "tools: Read, Glob, Grep, Bash, Edit, Write, mcp__acme__acme_put, mcp__acme__acme_find, mcp__acme__acme_fetch",
			bodyWant:    "- **save** → `acme_put`",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			agentRoot := filepath.Join(t.TempDir(), "agents")
			if _, err := generateClaudeAgents(skill.FS(), agentRoot, c.provider, InstallOptions{}); err != nil {
				t.Fatalf("generateClaudeAgents: %v", err)
			}
			for file, line := range map[string]string{"asdt-analyst.md": c.analystLine, "asdt-builder.md": c.builderLine} {
				data, err := os.ReadFile(filepath.Join(agentRoot, file))
				if err != nil {
					t.Fatalf("read %s: %v", file, err)
				}
				content := string(data)
				if !strings.Contains(content, line+"\n") {
					t.Errorf("%s missing EXACT tools line\nwant: %s\ngot:\n%s", file, line, content)
				}
				if !strings.Contains(bindingRegion(t, filepath.Join(agentRoot, file)), c.bodyWant) {
					t.Errorf("%s body binding region missing %q", file, c.bodyWant)
				}
			}
		})
	}
}

func TestBindMemory(t *testing.T) {
	const begin, end = memoryBindingBeginMarker, memoryBindingEndMarker
	incomplete := fakeMemoryProvider
	incomplete.Tools = map[MemoryVerb]MemoryTool{MemorySave: fakeMemoryProvider.Tools[MemorySave]}
	noUsage := fakeMemoryProvider
	noUsage.Tools = map[MemoryVerb]MemoryTool{
		MemorySave:   fakeMemoryProvider.Tools[MemorySave],
		MemorySearch: {Name: "acme_find"},
		MemoryGet:    fakeMemoryProvider.Tools[MemoryGet],
	}
	noPrefixes := fakeMemoryProvider
	noPrefixes.ClaudeToolPrefixes = nil
	emptyPrefix := fakeMemoryProvider
	emptyPrefix.ClaudeToolPrefixes = []string{""}
	probeWithoutURL := fakeMemoryProvider
	probeWithoutURL.Detect = ProviderDetect{Binary: "acme"}

	cases := []struct {
		name     string
		content  string
		provider ProviderDescriptor
		wantErr  string // substring; "" = success
		wantSame bool   // content must pass through unchanged
	}{
		{name: "marker-free content passes through, even for an empty provider", content: "# plain\n", provider: ProviderDescriptor{}, wantSame: true},
		{name: "empty region is bound", content: "a\n" + begin + "\n" + end + "\nb\n", provider: fakeMemoryProvider},
		{name: "begin marker without end fails", content: begin + "\n", provider: fakeMemoryProvider, wantErr: "end marker"},
		{name: "duplicated region fails", content: begin + end + begin + end, provider: fakeMemoryProvider, wantErr: "appears 2 times"},
		{name: "provider missing a verb fails", content: begin + end, provider: incomplete, wantErr: `binds no tool to verb "search"`},
		{name: "verb without Usage fails", content: begin + end, provider: noUsage, wantErr: `gives verb "search" no Usage`},
		{name: "provider without Claude prefixes fails", content: begin + end, provider: noPrefixes, wantErr: "has no ClaudeToolPrefixes"},
		{name: "empty Claude prefix fails", content: begin + end, provider: emptyPrefix, wantErr: "empty ClaudeToolPrefixes entry"},
		{name: "probe without install URL fails", content: begin + end, provider: probeWithoutURL, wantErr: "no Detect.InstallURL"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := bindMemory([]byte(c.content), c.provider)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("bindMemory error = %v, want one containing %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("bindMemory: %v", err)
			}
			if c.wantSame {
				if string(got) != c.content {
					t.Errorf("bindMemory changed marker-free content: %q", got)
				}
				return
			}
			again, err := bindMemory(got, c.provider)
			if err != nil || string(again) != string(got) {
				t.Errorf("bindMemory is not idempotent: err=%v\nfirst:  %q\nsecond: %q", err, got, again)
			}
			if !strings.Contains(string(got), "`acme_put`") {
				t.Errorf("bound content missing the save tool: %q", got)
			}
		})
	}
}
