package installer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// executorHeaderPath is the embedded path of the shared executor header that
// every generated agent definition bakes into its body. Its presence in the
// production embedded FS is guarded by skill/embedded_test.go.
const executorHeaderPath = "asdt-core/executor-header.md"

// agentTypeSpec describes one ASDT executor agent type. The same spec is
// rendered into each assistant's native agent-definition format.
type agentTypeSpec struct {
	ID          string   // "analyst" | "builder"; agent file name is "asdt-"+ID
	Description string   // one-line summary emitted into the frontmatter
	ClaudeTools []string // Claude Code tools allowlist before codegraph grants
	Constraints string   // English constraint prose appended after the executor header
	Guidance    string   // gated section body appended after constraints; empty by
	// default, set only by withCodegraphGuidance when CodegraphFound.
}

// AgentTypeNames lists the canonical ASDT agent type IDs in render order.
// workflow.yaml `agent:` fields must use one of these values.
var AgentTypeNames = []string{"analyst", "builder"}

const analystConstraints = `## asdt-analyst constraints

You are a single-step ANALYSIS executor: you read, inspect, and reason, then
return one summary to the orchestrator. You never write a file, anywhere.

- Bash is INSPECTION-ONLY. Allowed command families: git log, git diff,
  git status, git show, git blame, ls, find, grep, wc, cat.
- Forbidden command families: anything that mutates state — creating, editing,
  moving, or deleting files; git add / git commit / git push; package installs;
  build or codegen output; redirections that write to disk.
- If the step you were given appears to require creating or editing files,
  STOP: do not attempt the write. Record the blocked work and its target paths
  in Open items — that work belongs to asdt-builder.
- You never delegate: do not call Agent/Task and do not run other steps.
`

const builderConstraints = `## asdt-builder constraints

You are a single-step IMPLEMENTATION executor: you create and edit files to
complete the one step you were given, then return one summary to the
orchestrator.

- Create or edit files ONLY within your step's targets: the Edit targets your
  spec declares (or, for asdt-init/write, .asdt/ and the legacy-store
  exception below). If a needed change falls outside them, STOP that edit, do
  not write it, and record the path in Open items.
- Bash is for inspection only — never run build, lint, test, install, or
  codegen. Verification belongs to the orchestrator's verify gate. Its one
  mutation is the exception the executor header names: asdt-init/write
  removing the legacy .asdt/knowledge/ files on recorded consent.
- You never delegate: do not call Agent/Task and do not run other steps.
`

// agentTypeSpecs holds the canonical specs in the same order as AgentTypeNames.
var agentTypeSpecs = []agentTypeSpec{
	{
		ID:          "analyst",
		Description: "Read-only ASDT step executor for analysis steps — inspects the repository with inspection-only Bash and returns one summary, never writing files.",
		ClaudeTools: []string{"Read", "Glob", "Grep", "Bash"},
		Constraints: analystConstraints,
	},
	{
		ID:          "builder",
		Description: "Write-capable ASDT step executor for implementation steps — creates and edits files within declared targets and returns one summary, never delegating.",
		ClaudeTools: []string{"Read", "Glob", "Grep", "Bash", "Edit", "Write"},
		Constraints: builderConstraints,
	},
}

// codegraphClaudeTools are the read-only codegraph MCP tools appended to executor
// agent grants when the installer detects codegraph. Single namespace (codegraph
// ships as a binary MCP, not a plugin), all read-only.
var codegraphClaudeTools = []string{
	"mcp__codegraph__codegraph_explore",
	"mcp__codegraph__codegraph_search",
	"mcp__codegraph__codegraph_callers",
	"mcp__codegraph__codegraph_callees",
	"mcp__codegraph__codegraph_impact",
	"mcp__codegraph__codegraph_node",
	"mcp__codegraph__codegraph_files",
	"mcp__codegraph__codegraph_status",
}

// codeIntelligenceGuidance is the shared "## Code intelligence" body appended
// to an executor agent when codegraph is detected. Both agent types receive it
// verbatim — it teaches the codegraph read-only tools and where Bash still wins.
const codeIntelligenceGuidance = `## Code intelligence

This repository is indexed by codegraph, and you hold its read-only tools.
Query the index for anything symbol- or structure-shaped instead of
scanning files by hand:

- codegraph_explore FIRST for symbol, architecture, or "how does this
  work" questions — one call returns the relevant source grouped by file.
- codegraph_callers and codegraph_callees for call-graph questions ("what
  calls this", "what does this call").
- codegraph_impact for "what breaks if I change this".
- codegraph_search to locate a symbol, codegraph_node to read one
  definition in full; the remaining codegraph_* tools cover files and
  index status.

Bash grep and read stay the right tool for text the index does not model:
schemas like Rails db/schema.rb, configs, locales, generated DSL. Reach for
them there without hesitation.
`

// withCodegraphTools returns spec extended with the codegraph tool grants when
// codegraph was detected, and spec unchanged otherwise. The extended copy gets
// a freshly allocated ClaudeTools slice: a range copy of an agentTypeSpec still
// aliases the global agentTypeSpecs backing array, so appending in place would
// mutate the shared specs.
func withCodegraphTools(spec agentTypeSpec, opts InstallOptions) agentTypeSpec {
	if !opts.CodegraphFound {
		return spec
	}

	tools := make([]string, 0, len(spec.ClaudeTools)+len(codegraphClaudeTools))
	tools = append(tools, spec.ClaudeTools...)
	tools = append(tools, codegraphClaudeTools...)
	spec.ClaudeTools = tools
	return spec
}

// withCodegraphGuidance returns spec extended with the "## Code intelligence"
// section body when codegraph was detected, and spec unchanged otherwise.
// Unlike withCodegraphTools, Guidance is a plain string field: setting it on
// the value copy carries no aliasing hazard (no shared backing array to mutate).
func withCodegraphGuidance(spec agentTypeSpec, opts InstallOptions) agentTypeSpec {
	if !opts.CodegraphFound {
		return spec
	}

	spec.Guidance = codeIntelligenceGuidance
	return spec
}

// analystBashAllowlist is the OpenCode bash permission whitelist for the
// read-only analyst agent, in emission order. The catch-all "*": deny is
// emitted last by writeOpenCodePermissions.
var analystBashAllowlist = []string{
	"git log*", "git diff*", "git status", "git show*", "git blame*",
	"ls*", "find*", "grep *", "wc *", "cat *",
}

// renderClaudeAgent renders one agent type as a Claude Code agent definition
// (~/.claude/agents/asdt-<id>.md). The tools line is ALWAYS emitted so the
// allowlist is explicit — omission would grant the harness default toolset.
func renderClaudeAgent(spec agentTypeSpec, executorHeader string) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString("name: asdt-")
	b.WriteString(spec.ID)
	b.WriteString("\n")
	b.WriteString("description: ")
	b.WriteString(spec.Description)
	b.WriteString("\n")
	b.WriteString("tools: ")
	b.WriteString(strings.Join(spec.ClaudeTools, ", "))
	b.WriteString("\n")
	b.WriteString("---\n\n")
	writeAgentBody(&b, executorHeader, spec.Constraints, spec.Guidance)

	return b.String()
}

// renderOpenCodeAgent renders one agent type as an OpenCode subagent
// definition with a structural permission block. It ignores ClaudeTools by
// design: the codegraph grant is Claude-only; OpenCode subagents receive MCP
// tools implicitly, except the memory tools writeOpenCodePermissions denies.
func renderOpenCodeAgent(spec agentTypeSpec, executorHeader string) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString("description: ")
	b.WriteString(spec.Description)
	b.WriteString("\n")
	b.WriteString("mode: subagent\n")
	b.WriteString("permission:\n")
	writeOpenCodePermissions(&b, spec.ID)
	b.WriteString("---\n\n")
	writeAgentBody(&b, executorHeader, spec.Constraints, spec.Guidance)

	return b.String()
}

// writeOpenCodePermissions emits an agent's permission block. Beyond the
// edit/task/bash rules, every agent gets read access to the temp scratchpad
// root (external_directory) — on OpenCode the host has no session scratchpad,
// so step inputs live under the protocol.md §1 temp fallback, outside the
// worktree — and a deny for every known memory provider's tools, since OpenCode
// hands MCP tools to subagents implicitly and executors never reach memory.
// The builder's edit rule denies the scratchpad root: the orchestrator alone
// writes it.
func writeOpenCodePermissions(b *strings.Builder, id string) {
	scratch := openCodeScratchpadGlobs

	if id == "builder" {
		b.WriteString("  edit:\n")
		b.WriteString("    \"*\": allow\n")
		for _, glob := range scratch {
			fmt.Fprintf(b, "    %q: deny\n", glob)
		}
		b.WriteString("  task: deny\n")
		b.WriteString("  bash: allow\n")
	} else {
		// analyst (and any future read-only type): deny everything except the
		// inspection-only bash whitelist.
		b.WriteString("  edit: deny\n")
		b.WriteString("  task: deny\n")
		b.WriteString("  bash:\n")
		for _, glob := range analystBashAllowlist {
			fmt.Fprintf(b, "    %q: allow\n", glob)
		}
		b.WriteString("    \"*\": deny\n")
	}

	b.WriteString("  external_directory:\n")
	for _, glob := range scratch {
		fmt.Fprintf(b, "    %q: allow\n", glob)
	}
	for _, pattern := range openCodeMemoryToolPatterns() {
		fmt.Fprintf(b, "  %q: deny\n", pattern)
	}
}

// openCodeScratchpadGlobs are the external_directory patterns covering the
// protocol.md §1 temp scratchpad fallback `/tmp/asdt-{project}/{session}`. The
// fallback is always /tmp, never $TMPDIR, so a grant baked in at install time
// matches the path a run resolves later; /private/tmp covers macOS, where /tmp
// is a symlink to it.
var openCodeScratchpadGlobs = []string{"/tmp/asdt-*/**", "/private/tmp/asdt-*/**"}

// openCodeMemoryToolPatterns returns one OpenCode permission key per tool of
// every known provider, in Providers and verb order, without duplicates.
// OpenCode registers an MCP tool as `{server}_{tool}` and matches permission
// keys as wildcards against that name; the server name is the user's choice,
// so the pattern wildcards it.
func openCodeMemoryToolPatterns() []string {
	var patterns []string
	seen := map[string]bool{}
	for _, p := range Providers {
		for _, verb := range memoryVerbs {
			for _, tool := range p.Tools[verb] {
				pattern := "*_" + tool.Name
				if !seen[pattern] {
					seen[pattern] = true
					patterns = append(patterns, pattern)
				}
			}
		}
	}
	return patterns
}

// writeAgentBody emits the shared agent body in fixed order: the executor
// header VERBATIM, one blank line, the agent type's constraint prose, then —
// when non-empty — one blank line and the gated guidance section.
func writeAgentBody(b *strings.Builder, executorHeader, constraints, guidance string) {
	b.WriteString(executorHeader)
	if !strings.HasSuffix(executorHeader, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(constraints)
	if !strings.HasSuffix(constraints, "\n") {
		b.WriteString("\n")
	}
	if guidance != "" {
		b.WriteString("\n")
		b.WriteString(guidance)
		if !strings.HasSuffix(guidance, "\n") {
			b.WriteString("\n")
		}
	}
}

// AgentAdapterDescriptor describes how to generate executor agent definitions
// for one assistant, on top of the shared skill-tree copy.
type AgentAdapterDescriptor struct {
	AssistantID AssistantID
	Generate    func(skillsFS fs.FS, agentRoot string, opts InstallOptions) ([]string, error)
}

// AgentAdapters lists assistants that receive generated executor agent
// definitions. Unlike CommandAdapters, BOTH assistants carry an entry —
// each has a native agents directory.
var AgentAdapters = []AgentAdapterDescriptor{
	{
		AssistantID: AssistantClaudeCode,
		Generate:    generateClaudeAgents,
	},
	{
		AssistantID: AssistantOpenCode,
		Generate:    generateOpenCodeAgents,
	},
}

func agentAdapterFor(id AssistantID) (AgentAdapterDescriptor, bool) {
	for _, adapter := range AgentAdapters {
		if adapter.AssistantID == id {
			return adapter, true
		}
	}
	return AgentAdapterDescriptor{}, false
}

// agentRootFor resolves the agents directory for one assistant, or "" when
// the assistant has no known agent root.
func agentRootFor(id AssistantID) string {
	switch id {
	case AssistantClaudeCode:
		home, _ := os.UserHomeDir()
		return home + "/.claude/agents"
	case AssistantOpenCode:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return xdg + "/opencode/agents"
		}
		home, _ := os.UserHomeDir()
		return home + "/.config/opencode/agents"
	default:
		return ""
	}
}

func generateClaudeAgents(skillsFS fs.FS, agentRoot string, opts InstallOptions) ([]string, error) {
	return generateAgentFiles(skillsFS, agentRoot, opts, renderClaudeAgent)
}

func generateOpenCodeAgents(skillsFS fs.FS, agentRoot string, opts InstallOptions) ([]string, error) {
	return generateAgentFiles(skillsFS, agentRoot, opts, renderOpenCodeAgent)
}

// generateAgentFiles writes one agent definition file per agent type spec
// under agentRoot, baking the shared executor header into every body. Executors
// never reach memory — the orchestrator alone does — so no binding is spliced,
// the Claude Code allowlist grants no memory tool, and the OpenCode permission
// block denies every known provider's memory tools. A
// skillsFS without the executor header (a partial fixture) generates nothing
// — the embedded production FS is guaranteed to carry it. Per-file failures
// are isolated: the first error is reported, the remaining specs still write.
func generateAgentFiles(skillsFS fs.FS, agentRoot string, opts InstallOptions, render func(agentTypeSpec, string) string) ([]string, error) {
	header, readErr := fs.ReadFile(skillsFS, executorHeaderPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", executorHeaderPath, readErr)
	}

	if mkErr := os.MkdirAll(agentRoot, 0o755); mkErr != nil {
		return nil, fmt.Errorf("mkdir %s: %w", agentRoot, mkErr)
	}

	var written []string
	var firstErr error

	for _, spec := range agentTypeSpecs {
		target := filepath.Join(agentRoot, "asdt-"+spec.ID+".md")
		effSpec := withCodegraphGuidance(withCodegraphTools(spec, opts), opts)
		content := render(effSpec, string(header))

		if writeErr := os.WriteFile(target, []byte(content), 0o644); writeErr != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("write %s: %w", target, writeErr)
			}
			continue
		}
		written = append(written, target)
	}

	return written, firstErr
}

// generateAgents mirrors generateCommands: resolve the adapter and root for
// the assistant (absence of either is the no-op), generate, fold written
// paths into result.WrittenCommands and any error into result.Err.
func generateAgents(assistant AssistantDescriptor, skillsFS fs.FS, opts InstallOptions, result *InstallResult) {
	adapter, ok := agentAdapterFor(assistant.ID)
	if !ok {
		return
	}

	agentRoot := agentRootFor(assistant.ID)
	if agentRoot == "" {
		return
	}

	written, genErr := adapter.Generate(skillsFS, agentRoot, opts)
	result.WrittenCommands = append(result.WrittenCommands, written...)
	if genErr != nil {
		result.Err = genErr
	}
}
