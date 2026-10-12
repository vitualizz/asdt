package installer

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

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
	Tools: map[MemoryVerb][]ProviderTool{
		MemoryRemember: {{Name: "acme_put", Usage: "pass the subject as `slot` and the note as `body`"}},
		MemoryRecall:   {{Name: "acme_ask", Usage: "pass the query as `q`; each answer is a full note"}},
	},
	ClaudeToolPrefixes: []string{"mcp__acme__"},
	// No Detect.Binary: a hosted provider with nothing to probe locally.
}

// boundFiles are the installed files (SkillsDir-relative) that must carry a
// bound memory-binding region: the three orchestrator-facing files that own
// one, plus every routed SKILL.md, which receives protocol.md's copy through
// the specialist-header splice. The executor header is not among them: a
// sub-agent never reaches memory.
var boundFiles = []string{
	"asdt-core/protocol.md",
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
// assistant bound to provider (the zero ProviderDescriptor: no memory) and
// returns its SkillsDir.
func installWithProvider(t *testing.T, provider ProviderDescriptor) string {
	t.Helper()

	skillsDir := filepath.Join(t.TempDir(), "skills")
	a := AssistantDescriptor{ID: "memory-binding-test", Name: "Memory Binding Test", BinaryName: "sh", SkillsDir: skillsDir}
	results := InstallWithModels([]AssistantDescriptor{a}, provider, skill.FS(), "", nil, InstallOptions{})
	if results[0].Err != nil {
		t.Fatalf("install with provider %q: %v", provider.ID, results[0].Err)
	}
	meta, err := ReadInstallMeta(a)
	if err != nil {
		t.Fatalf("read install meta: %v", err)
	}
	if want := provider.MetaValue(); meta.Memory != want {
		t.Errorf("install meta Memory = %q, want %q", meta.Memory, want)
	}
	if meta.CodeIntel != CodeIntelNone {
		t.Errorf("install meta CodeIntel = %q, want %q", meta.CodeIntel, CodeIntelNone)
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
		neutral  bool     // no installed file may name a shipped provider outside a region
	}{
		{
			name:     "engram binds today's tool names",
			provider: Providers[0],
			want: []string{
				"**Memory provider: Engram.**",
				"`memory.provider: engram`",
				"- **remember** → `mem_save` — pass `project: {project}`, `topic_key: {project}/notes/{subject}`",
				"- **recall** → `mem_search` — pass the query as `query`",
				"; then `mem_get_observation` — pass a match's `id` as `id`",
				"`mcp__plugin_engram_engram__mem_save` or `mcp__engram__mem_save`",
			},
		},
		{
			name:     "a second provider binds with no prompt edits",
			provider: fakeMemoryProvider,
			want: []string{
				"**Memory provider: Acme Memory.**",
				"`memory.provider: acme`",
				"- **remember** → `acme_put` — pass the subject as `slot` and the note as `body`.",
				"- **recall** → `acme_ask` — pass the query as `q`; each answer is a full note.",
				"`mcp__acme__acme_put`",
			},
			neutral: true,
		},
		{
			// protocol.md §2 and asdt-init/SKILL.md read a binding that names no
			// provider as "Memory unavailable"; init then passes
			// memory_provider: none and writes no memory.provider key.
			name:     "no memory binds no provider",
			provider: ProviderDescriptor{},
			want: []string{
				"**No memory provider bound.**",
				`"Memory unavailable"`,
				"Re-running the ASDT installer with a provider selected connects one.",
			},
			neutral: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			skillsDir := installWithProvider(t, c.provider)

			header, err := os.ReadFile(filepath.Join(skillsDir, "asdt-core", "executor-header.md"))
			if err != nil {
				t.Fatalf("read installed executor header: %v", err)
			}
			if strings.Contains(string(header), memoryBindingBeginMarker) || strings.Contains(string(header), memoryBindingEndMarker) {
				t.Errorf("installed executor-header.md carries a memory-binding region; sub-agents never reach memory")
			}

			for _, rel := range boundFiles {
				region := bindingRegion(t, filepath.Join(skillsDir, filepath.FromSlash(rel)))
				for _, w := range c.want {
					if !strings.Contains(region, w) {
						t.Errorf("%s binding region missing %q\nregion:%s", rel, w, region)
					}
				}
				if c.provider.IsNone() {
					if hits := providerNameHits(region, append(shippedProviderNames(), "memory provider:", "memory.provider:"), nil); len(hits) > 0 {
						t.Errorf("%s no-memory binding names %v; it must bind no provider\nregion:%s", rel, hits, region)
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

// shippedProviderNames is every name a shipped provider — memory or code
// intelligence — is known by: its Name, probe binary, tool names, and (code
// intelligence) index marker, plus `mem_`, Engram's whole tool family
// (mem_update, mem_context…), none of which the prompts may reach for. Derived
// from Providers and CodeIntelProviders, so a new provider's names are
// forbidden in the same commit that adds its descriptor.
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
			for _, tool := range p.Tools[verb] {
				add(tool.Name)
			}
		}
	}
	for _, p := range CodeIntelProviders {
		add(p.Name)
		add(p.Detect.Binary)
		add(p.IndexMarker)
		for _, tool := range p.Tools {
			add(tool.Name)
		}
	}
	return names
}

// providerNameHits returns what content names outside its binding regions
// (memory and code intelligence) that only a binding may name: a forbidden
// substring (lowercase, matched case-insensitively), or an MCP tool name
// (`mcp__…`) whose namespace is not in allowedMCP.
func providerNameHits(content string, forbidden, allowedMCP []string) []string {
	content = stripRegions(content, memoryBindingBeginMarker, memoryBindingEndMarker)
	content = stripRegions(content, codeIntelBindingBeginMarker, codeIntelBindingEndMarker)
	lower := strings.ToLower(content)

	var hits []string
	for _, f := range forbidden {
		if strings.Contains(lower, f) {
			hits = append(hits, f)
		}
	}
	for rest := lower; ; {
		i := strings.Index(rest, "mcp__")
		if i < 0 {
			break
		}
		rest = rest[i:]
		ok := false
		for _, a := range allowedMCP {
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

// stripRegions returns content without any begin…end marker region, markers
// included.
func stripRegions(content, begin, end string) string {
	for {
		b := strings.Index(content, begin)
		e := strings.Index(content, end)
		if b < 0 || e < b {
			return content
		}
		content = content[:b] + content[e+len(end):]
	}
}

// TestPromptsAreProviderNeutral asserts no prompt text names a provider —
// memory or code intelligence — or its tools outside a binding region. Prompts
// speak the memory interface (asdt-core/protocol.md §2) and the
// code-intelligence interface (asdt-core/executor-header.md); the installer
// binds each to the selected provider, so a provider is added in Go alone. It
// covers every source of prompt text: the embedded skill tree, the generated
// executor agents (Claude Code and OpenCode, Codegraph bound — its binding
// region is the one place its names may appear), the Go prompt-prose
// constants, and the AGENTS.md template and personas. An agent's frontmatter
// is configuration, so only its body is scanned.
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

	for name, generate := range map[string]func(fs.FS, string, InstallOptions) ([]string, error){
		"claude":   generateClaudeAgents,
		"opencode": generateOpenCodeAgents,
	} {
		agentRoot := filepath.Join(t.TempDir(), name)
		if _, err := generate(skill.FS(), agentRoot, InstallOptions{CodeIntel: CodeIntelProviders[0]}); err != nil {
			t.Fatalf("generate %s agents: %v", name, err)
		}
		first := len(sources)
		readAll(os.DirFS(agentRoot), ".", name+" agent ", nil)
		// The frontmatter is configuration, not prompt text: OpenCode's
		// permission block must name the memory tools it denies, and
		// TestGenerateAgents_NoMemoryAccess guards it. Only the body is prose.
		for i := first; i < len(sources); i++ {
			sources[i].content = stripFrontmatter(sources[i].content)
		}
	}

	for name, prose := range map[string]string{
		"analystConstraints": analystConstraints,
		"builderConstraints": builderConstraints,
	} {
		sources = append(sources, source{"const " + name, prose, nil})
	}

	for _, src := range sources {
		for _, hit := range providerNameHits(src.content, shippedProviderNames(), src.allowedMCP) {
			t.Errorf("%s names %q outside a binding region — speak the provider-neutral interface (memory **remember** / **recall**, code intelligence **explore** …) instead", src.name, hit)
		}
	}
}

// stripFrontmatter returns content after its leading `---` frontmatter block,
// or content unchanged when it has none.
func stripFrontmatter(content string) string {
	const delim = "---\n"
	if !strings.HasPrefix(content, delim) {
		return content
	}
	if _, body, ok := strings.Cut(content[len(delim):], "\n"+delim); ok {
		return body
	}
	return content
}

// TestGenerateAgents_NoMemoryAccess asserts the executor agents, generated
// from the REAL executor header, hold no memory tool and no binding region:
// the orchestrator alone remembers and recalls (protocol.md §0, §2), so a
// sub-agent that could reach memory would reopen the second writer the
// ownership rule closes.
func TestGenerateAgents_NoMemoryAccess(t *testing.T) {
	agentRoot := filepath.Join(t.TempDir(), "agents")
	if _, err := generateClaudeAgents(skill.FS(), agentRoot, InstallOptions{}); err != nil {
		t.Fatalf("generateClaudeAgents: %v", err)
	}
	for file, line := range map[string]string{"asdt-analyst.md": analystBaseToolsLine, "asdt-builder.md": builderBaseToolsLine} {
		data, err := os.ReadFile(filepath.Join(agentRoot, file))
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		content := string(data)
		if !strings.Contains(content, line+"\n") {
			t.Errorf("%s missing EXACT tools line\nwant: %s\ngot:\n%s", file, line, content)
		}
		if strings.Contains(content, memoryBindingBeginMarker) || strings.Contains(content, memoryBindingEndMarker) {
			t.Errorf("%s carries a memory-binding region; executors never reach memory", file)
		}
	}

	// OpenCode hands MCP tools to subagents implicitly, so the only guard is
	// an explicit deny for every tool of every known provider.
	ocRoot := filepath.Join(t.TempDir(), "opencode-agents")
	if _, err := generateOpenCodeAgents(skill.FS(), ocRoot, InstallOptions{}); err != nil {
		t.Fatalf("generateOpenCodeAgents: %v", err)
	}
	for _, file := range []string{"asdt-analyst.md", "asdt-builder.md"} {
		fm, body := parseOpenCodeAgentFile(t, filepath.Join(ocRoot, file))
		if strings.Contains(body, memoryBindingBeginMarker) || strings.Contains(body, memoryBindingEndMarker) {
			t.Errorf("OpenCode %s carries a memory-binding region; executors never reach memory", file)
		}
		for _, p := range Providers {
			for _, verb := range memoryVerbs {
				for _, tool := range p.Tools[verb] {
					key := "*_" + tool.Name
					if got := fm.Permission.Tools[key]; got != "deny" {
						t.Errorf("OpenCode %s permission[%q] = %v, want deny: provider %q's %s tool stays reachable", file, key, got, p.ID, verb)
					}
				}
			}
		}
	}
}

func TestBindMemory(t *testing.T) {
	const begin, end = memoryBindingBeginMarker, memoryBindingEndMarker
	incomplete := fakeMemoryProvider
	incomplete.Tools = map[MemoryVerb][]ProviderTool{MemoryRemember: fakeMemoryProvider.Tools[MemoryRemember]}
	noName := fakeMemoryProvider
	noName.Tools = map[MemoryVerb][]ProviderTool{
		MemoryRemember: fakeMemoryProvider.Tools[MemoryRemember],
		MemoryRecall:   {fakeMemoryProvider.Tools[MemoryRecall][0], {Usage: "read the note"}},
	}
	noUsage := fakeMemoryProvider
	noUsage.Tools = map[MemoryVerb][]ProviderTool{
		MemoryRemember: fakeMemoryProvider.Tools[MemoryRemember],
		MemoryRecall:   {{Name: "acme_ask"}},
	}
	noPrefixes := fakeMemoryProvider
	noPrefixes.ClaudeToolPrefixes = nil
	emptyPrefix := fakeMemoryProvider
	emptyPrefix.ClaudeToolPrefixes = []string{""}
	noID := fakeMemoryProvider
	noID.ID = ""
	configNone := fakeMemoryProvider
	configNone.ConfigValue = MemoryNone
	probeWithoutURL := fakeMemoryProvider
	probeWithoutURL.Detect = ProviderDetect{Binary: "acme"}

	cases := []struct {
		name     string
		content  string
		provider ProviderDescriptor
		wantErr  string // substring; "" = success
		wantSame bool   // content must pass through unchanged
		wantText string // substring the bound content must carry ("" = the remember tool)
		lackText string // substring the bound content must NOT carry
	}{
		{name: "marker-free content passes through, even with no provider", content: "# plain\n", provider: ProviderDescriptor{}, wantSame: true},
		{name: "no memory binds the no-provider body", content: begin + end, provider: ProviderDescriptor{}, wantText: "**No memory provider bound.**", lackText: "Memory provider:"},
		{name: "provider without ID fails", content: begin + end, provider: noID, wantErr: "no ID, Name or ConfigValue"},
		{name: "ConfigValue none fails", content: begin + end, provider: configNone, wantErr: "which means no provider"},
		{name: "empty region is bound", content: "a\n" + begin + "\n" + end + "\nb\n", provider: fakeMemoryProvider},
		{name: "begin marker without end fails", content: begin + "\n", provider: fakeMemoryProvider, wantErr: "end marker"},
		{name: "duplicated region fails", content: begin + end + begin + end, provider: fakeMemoryProvider, wantErr: "appears 2 times"},
		{name: "provider missing a verb fails", content: begin + end, provider: incomplete, wantErr: `binds no tool to verb "recall"`},
		{name: "recipe tool without Name fails", content: begin + end, provider: noName, wantErr: `binds a tool with no Name to verb "recall"`},
		{name: "tool without Usage fails", content: begin + end, provider: noUsage, wantErr: `gives verb "recall" tool "acme_ask" no Usage`},
		{name: "provider without Claude prefixes binds without the prefix sentence", content: begin + end, provider: noPrefixes, lackText: "A host may prefix"},
		{name: "multi-call recipe binds in order", content: begin + end, provider: Providers[0], wantText: "`mem_search` — pass the query as `query` and `project: {project}`; each match is a truncated preview with its `id` — a preview is never the note; then `mem_get_observation`"},
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
			want := c.wantText
			if want == "" {
				want = "`" + c.provider.Tools[MemoryRemember][0].Name + "`"
			}
			if !strings.Contains(string(got), want) {
				t.Errorf("bound content missing %q: %q", want, got)
			}
			if c.lackText != "" && strings.Contains(string(got), c.lackText) {
				t.Errorf("bound content carries %q: %q", c.lackText, got)
			}
		})
	}
}

// bindingStripRule is the heading of the protocol.md §5 rule that makes the
// orchestrator delete every memory-binding region from text it injects into a
// sub-agent prompt.
const bindingStripRule = "**Strip the memory binding.**"

// stripInjectedBinding applies the protocol.md §5 strip rule as worded: the
// explanatory `memory-binding:` comment above a region, its opening and closing
// comments, and everything between them go.
func stripInjectedBinding(content string) string {
	for {
		begin := strings.Index(content, memoryBindingBeginMarker)
		end := strings.Index(content, memoryBindingEndMarker)
		if begin < 0 || end < begin {
			return content
		}
		start := begin
		if note := strings.LastIndex(content[:begin], "<!-- memory-binding:"); note >= 0 && !strings.Contains(content[note:begin], "-->\n\n") {
			start = note
		}
		content = content[:start] + content[end+len(memoryBindingEndMarker):]
	}
}

// TestSubagentInjectionStripsMemoryBinding guards the one path by which a
// binding could reach a sub-agent: protocol.md is a reference skill of most
// sub-agent steps, and the installed copy carries the bound region. The region
// must stay (the orchestrator reads it), so the guard is the §5 rule: it must
// exist and name the markers whenever any injected file carries a region, and
// applying it as worded must leave no provider name behind. No step file may
// lean on the binding either — a value a step needs from it arrives in its
// payload (asdt-init/clarify's memory_provider). The code-intelligence binding
// is different: the executor header carries it to every sub-agent, so a step
// file may lean on it (asdt-init/explore probes its index marker).
func TestSubagentInjectionStripsMemoryBinding(t *testing.T) {
	skillsDir := installWithProvider(t, Providers[0])

	protocol, err := os.ReadFile(filepath.Join(skillsDir, "asdt-core", "protocol.md"))
	if err != nil {
		t.Fatalf("read installed protocol.md: %v", err)
	}
	sec5Start := strings.Index(string(protocol), "## 5.")
	sec5End := strings.Index(string(protocol), "## 6.")
	if sec5Start < 0 || sec5End < sec5Start {
		t.Fatal("installed protocol.md has no §5 followed by §6")
	}
	sec5 := string(protocol)[sec5Start:sec5End]
	ruleIdx := strings.Index(sec5, bindingStripRule)
	if ruleIdx < 0 {
		t.Fatalf("protocol.md §5 lacks %q; the bound region would reach every sub-agent that loads protocol.md", bindingStripRule)
	}
	rule, _, _ := strings.Cut(sec5[ruleIdx:], "\n\n")
	for _, want := range []string{"ASDT:GENERATED:memory-binding", "`memory-binding:` comment", "sub-agent prompt"} {
		if !strings.Contains(rule, want) {
			t.Errorf("protocol.md §5 strip rule does not name %q:\n%s", want, rule)
		}
	}

	stepFileWord := regexp.MustCompile(`(?i)memory[- ]binding|\bmemory\b[^.\n]*\bbinding\b|\bbinding\b[^.\n]*\bmemory\b`)
	workflows, err := filepath.Glob(filepath.Join(skillsDir, "*", "workflow.yaml"))
	if err != nil || len(workflows) == 0 {
		t.Fatalf("glob installed workflow.yaml: %v (found %d)", err, len(workflows))
	}
	carriers := 0
	for _, wfPath := range workflows {
		dir := filepath.Dir(wfPath)
		data, err := os.ReadFile(wfPath)
		if err != nil {
			t.Fatalf("read %s: %v", wfPath, err)
		}
		var wf workflowFile
		if err := yaml.Unmarshal(data, &wf); err != nil {
			t.Fatalf("parse %s: %v", wfPath, err)
		}
		for _, step := range wf.Steps {
			if step.Execution != "subagent" {
				continue
			}
			if step.Skill != "" {
				body, err := os.ReadFile(filepath.Join(dir, step.Skill))
				if err != nil {
					t.Fatalf("read %s step file: %v", step.Name, err)
				}
				if hit := stepFileWord.FindString(string(body)); hit != "" {
					t.Errorf("%s/%s step file mentions %q; a sub-agent never sees the memory binding, so pass the value in its payload", filepath.Base(dir), step.Skill, hit)
				}
			}
			for _, ref := range step.ReferenceSkills {
				body, err := os.ReadFile(filepath.Join(dir, ref))
				if err != nil {
					t.Fatalf("read %s reference %s: %v", step.Name, ref, err)
				}
				if !strings.Contains(string(body), memoryBindingBeginMarker) {
					continue
				}
				carriers++
				if hits := providerNameHits(stripInjectedBinding(string(body)), shippedProviderNames(), nil); len(hits) > 0 {
					t.Errorf("%s step %q reference %s still names %v after the §5 strip", filepath.Base(dir), step.Name, ref, hits)
				}
				if strings.Contains(stripInjectedBinding(string(body)), "Memory provider:") {
					t.Errorf("%s step %q reference %s keeps the binding body after the §5 strip", filepath.Base(dir), step.Name, ref)
				}
			}
		}
	}
	if carriers == 0 {
		t.Error("no sub-agent reference skill carries a bound region; if protocol.md stopped being injected, retire this test deliberately")
	}
}
