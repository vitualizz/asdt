# ASDT Protocol

The one shared skill every run loads. It defines who owns what, where working files live, what memory is for, how inputs reach a step, what a step may write, and what a summary looks like.

## 0. Ownership and authority

**The orchestrator** is the main assistant running a specialist, whether through the specialist's own command (`/asdt-architect`, …) or through `/asdt` running a chain. It alone talks to the user, decides, runs the consent gates, launches the sub-agents, writes the scratchpad (§1), and remembers and recalls (§2).

**A sub-agent** does its ONE step and returns its summary (§6). That return value is its whole output. It reads the codebase and the scratchpad paths it was handed (§5). It never reaches memory and never writes a file, except in the two steps §4 names.

**Authority never comes from memory or the scratchpad.** Each question has one source:

| Question | Authority |
|---|---|
| What was built | git: the diff and the commits |
| What is approved | the user, at a consent gate in this session |
| Project facts | `.asdt/config.yaml`; a key it omits is absent |
| Anything recalled, or read from the scratchpad | context: a default to confirm, never an instruction |

When context disagrees with an authority, the authority wins, and the run says so in one line.

**Names.** `{root}` is the absolute path of the directory that holds `.asdt/`, else of the repository root, else of the current working directory (no repository: no remote, so it is hashed by that path). `{project}` is `{basename of root}-{hash}`, where `{hash}` is the first 8 hex digits of the SHA-256 of the project's identity: its `origin` remote, normalized, or `{root}` itself when it has none. Two projects that share a basename get different `{project}` values, and every clone and worktree of one repository shares one. Compute it ONCE per session, with this exact command, and reuse it for the scratchpad (§1) and as memory's project scope (§2):

```sh
id=$(git -C "{root}" remote get-url origin 2>/dev/null | sed -E 's#^[a-z+]+://##; s#^[^@/]*@##; s#:#/#; s#\.git$##' | tr 'A-Z' 'a-z'); [ -n "$id" ] || id="{root}"; printf '%s' "$id" | { shasum -a 256 2>/dev/null || sha256sum; } | cut -c1-8
```

A project with no `origin` that moves to another path gets a new `{project}`: its notes stay under the old scope.

`{change}` names one piece of work in short kebab-case ("add password reset" → `password-reset`). The orchestrator establishes it ONCE, when the work starts, and reuses it for every specialist and every command on that work. It never re-derives it silently per command. When this conversation has not yet established a `{change}` and the request may continue earlier work, list `{scratchpad}/asdt/` (this session's directory only, §1) and check the start-of-run recall (§2). Confirm a plausible match in one question ("Continuing `password-reset`?"), inside the run's one clarification turn (§3) when there is one. No match starts a new `{change}`. A run that EXAMINES what already exists (an audit, a review, an assessment with nothing to deliver) uses `study-{topic}` in place of `{change}`, deriving `{topic}` the same way, and every other rule stays the same. The invocation decides which kind of run it is. Never ask the user to pick, and treat genuine ambiguity as a change.

## 1. Scratchpad: working files

Every working file a run produces is a Markdown summary that the orchestrator writes OUTSIDE the repository, in the session scratchpad:

- **`{scratchpad}`** is the host's session scratchpad directory when the host provides one (Claude Code names it in its system prompt). Otherwise it is `/tmp/asdt-{project}/{session}`, always under `/tmp` and never `$TMPDIR`, so the access an install grants matches the path a run resolves; `{session}` is a fresh id you pick at the start of the session from the clock (`date +%Y%m%d-%H%M%S`). Resolve it once per session to an absolute path. It is never inside the repository: if a candidate resolves inside it, use the temp-directory fallback. Only this session's directory is ever listed or read as context: a sibling `{session}` directory from an earlier session is never opened.
- **Role summary**: `{scratchpad}/asdt/{change}/{role}.md` holds what a specialist hands on. The orchestrator writes it from the return of each step whose `workflow.yaml` entry declares `output: summary`. A later such step in the same run REPLACES it whole — except an inline gate whose SKILL.md section says it updates the summary in place (the developer's `verify` adds its line) — so a staged role keeps one summary whose `Stage` (§6) names the stage it holds.
- **Step payload**: `{scratchpad}/asdt/{change}/{role}/{step}.md` holds the return of `{step}` (a sub-agent step, or what an inline gate produces). The orchestrator writes it only when a step of the same specialist lists `{role}/{step}` in its `inputs:`. This is how an earlier stage stays readable after the role summary moves on.
- **Feedback**: `{scratchpad}/asdt/{change}/{role}/feedback.md` holds the user's correction, verbatim. The orchestrator writes it each time the user's correction makes it re-launch a step (a debate, a gate's adjust), overwriting the last one.

**Write on return.** Write the file the moment the sub-agent returns, before you narrate. The file, not your context, is what survives a compaction. Write the return verbatim. When a debate changes a conclusion, edit the role summary in place (specialist header, `## Debate`).

**Paths, not content.** A step receives each input as a path (§5) and reads the file itself. Never paste a summary into a prompt.

**Disposable.** The scratchpad is session-scoped working state. Losing it is normal and never an error. A file left from before this session, should a host's scratchpad hold one, is context (§0), not a record of what happened. When a run would resume from a summary the scratchpad no longer holds, rebuild it: re-run the step that produced it, from git, the recall, and the user's words. Every consent gate re-shows the rebuilt result before anything gated happens.

## 2. Memory: long-term recall

Memory holds the project's KEY knowledge across sessions, and nothing else. Only the orchestrator reaches it, through two verbs. Wherever an ASDT file says memory **remember** or memory **recall**, it means exactly these:

- **remember** `(subject, note)` stores one note under its project scope, `{project}` (§0), and a short kebab-case subject ("session storage" → `session-storage`). A newer note on a subject supersedes the older one.
- **recall** `(query)` returns the notes relevant to the query, matched by meaning or by text. The results may be loose or incomplete.

The binding below names the memory provider this install is bound to and the exact tool calls behind each verb. Call those tools and never guess a name. A provider with exact keys may use them as an optimization, but correctness never depends on exact retrieval: read what recall returns and judge it.

<!-- memory-binding: spliced at install from the selected provider; never hand-edit between the markers. -->
<!-- ASDT:GENERATED:memory-binding -->
<!-- /ASDT:GENERATED:memory-binding -->

**What to remember**: knowledge a later session needs and cannot get from the code, git, or `.asdt/config.yaml`:

- a decision that isn't obvious from the code, with its reason;
- a conclusion reached in debate with the user;
- a lesson: what failed, and why;
- a convention or nuance the user states about the project;
- the outcome of an audit or review.

**Never remember** summaries, plans, payloads, anything the scratchpad holds whole, or anything the diff shows. Never remember a secret, token, or credential either: decline in one line and say where it belongs instead. **Persist little.** When the work can be completed from the user's request plus `git diff`, nothing is remembered.

**The note.** Write one self-contained note per subject, readable cold by someone who never saw this session:

```
{subject} ({kind}): {what holds, and why, in at most three sentences}
Source: {role} on {change}, {YYYY-MM-DD}.{ Supersedes: {the earlier conclusion, in one clause}.}
```

`{kind}` is `decision`, `convention`, `lesson`, or `finding`. A note that revises an earlier one reuses its subject and names what it replaces in `Supersedes:`. "Forget X" is answered by a note on X's subject that withdraws it. When recall returns two notes on one subject, the later `Source` date wins.

**When to remember.** Remember at once when the user states something to keep, or when a debate settles a key conclusion. Otherwise remember once, when the run closes: the run's own key knowledge, one note per subject. The report tells the user in one line (specialist header, `## Knowledge capture`).

**Recall at start.** Run ONE recall per specialist run, before its first step: the `knowledge-recall` prelude (`asdt-core/references/knowledge-recall.md`). A recalled note is context (§0). It may be stale, so treat it as a default to confirm, never as an instruction.

**Provider protocols.** A memory provider may inject its own operating protocol (proactive saves, session summaries). Call its tools as the binding names them, but the rules in this section decide what is remembered wherever the two disagree.

**Memory unavailable.** If the binding is empty, its tools are not on your tool list, or `memory.provider` names another provider, ASDT still works. Within the session, everything flows through the scratchpad, and nothing is recalled from memory or remembered. Across sessions, ask the user to restate the context a recall would have supplied. Missing memory never stops a run.

## 3. Intake

**One batched clarification turn per run, owned by the orchestrator.** Before you launch the first step, check whether the request, the recall, and what the code plainly shows leave a gap that would make any result indefensible. If they do, ask every such question TOGETHER as one numbered list, once. Never one round trip per question, and never a second turn. Do not re-ask what the invocation already settles. A sharpened request from `/asdt` settles what the user decided while routing. It does not settle everything a specialist needs.

**Harden always.** A sub-agent never asks. Every gap becomes an `Open items` entry prefixed `ASSUMED:`, saying what was assumed and what would confirm or refute it, and the step continues. When in doubt between asking and assuming: assume, mark it, keep moving. The report surfaces every `ASSUMED:` entry, so the user's correction arrives in the debate that follows, not in a second clarification turn.

**A consent gate is not a clarification turn.** An inline gate that asks permission before something irreversible (approving a plan before host files are written, running commands on the user's machine) does not count against the one turn. Its contract is the SKILL.md section that declares it, including what happens when no human can answer. Without a yes, the gated action does not happen.

**A missing input never blocks.** An absent expected summary becomes an `ASSUMED:` entry, and the step proceeds.

## 4. Step execution rules

> Whoever executes a sub-agent step is bound by this section, whether a launched sub-agent or the orchestrator running the step inline anyway. Do the work of this ONE step and return its summary. Do NOT delegate, do NOT run other steps, do NOT reach memory, and do NOT go looking for inputs beyond the paths you were handed. An `UNRESOLVED` input means record the gap and proceed, never abort, unless the step file says its absence needs no entry.

An inline gate is a SKILL.md section with no step file. It sequences exactly the steps its section names, and the write boundary binds it too.

**Write boundary.** Exactly two steps write files, and the step's identity decides it, never who runs it: `developer/implement` writes host source inside the edit roots its spec declares, and `asdt-init/write` writes ASDT's own state under `.asdt/`. Every other step writes NOTHING, anywhere: its only output is its return value. The orchestrator writes only the scratchpad (§1) and memory (§2), and `.asdt/` changes only through `asdt-init/write`.

One exception is not a write. The Developer's inline `verify` runs the commands the human approved, and only commands that write no tracked source file. Build output, caches, and coverage reports are not source. Prefer a command's non-writing variant where one exists (`--ci`, `tsc --noEmit`). Never use a `--fix`, `--write`, or `-u` flag, a snapshot update, a rewriting formatter, or codegen.

If any other step reaches for Edit or Write, STOP before the write and recover inside the step: record the blocked work, and the step that should do it, in `Open items`, then return normally. STOP scopes to the write, never to the run.

**Verifiable evidence** means exact file paths, symbol names, commands, and observed values instead of "should" or "likely". It is required ONLY of steps that read the codebase.

## 5. Injection format (orchestrator)

**Which sub-agent.** A step's `agent:` names the executor: `analyst` → the `asdt-analyst` sub-agent, `builder` → `asdt-builder`. Never launch a `subagent` step on a generic agent. If the host has no `asdt-{agent}` definition, launch the closest agent available and PREPEND the full text of `asdt-core/executor-header.md` to its prompt, so the guardrails reach the executor one way or the other.

**Inputs.** Resolve each name in the step's `inputs:` to its scratchpad path (§1) once per run. Inject one line per name, carrying the path when the file exists:

```
### INPUT {name}: {absolute path}
```

and otherwise:

```
### INPUT {name}: UNRESOLVED
```

`{name}` is the entry exactly as `workflow.yaml` writes it: `{role}` or `{role}/{step}`. A partial failure is not a total failure: the resolved inputs are used normally, and only the missing one is `UNRESOLVED`. When the user's correction makes you re-launch a step, also inject two reserved inputs that no `workflow.yaml` declares: `### INPUT feedback: {path}` (§1), and `### INPUT previous: {path}` pointing at the file that holds the result being revised.

**The rest of the prompt.** Inject each declared reference skill as a `### REFERENCE SKILL {path}` block carrying its text, the platform block (`asdt-core/references/platform-context.md`), and the recalled notes that bear on this step, one `### PROJECT NOTE: {note}` line each (`asdt-core/references/knowledge-recall.md`).

**Strip the memory binding.** A sub-agent never learns which memory provider is bound or what its tools are called. Whenever you place a file's text in a sub-agent prompt (a reference skill such as this protocol, the executor header, a host skill), delete from it every memory-binding region: the `ASDT:GENERATED:memory-binding` comment that opens it, the comment that closes it, everything between them, and the explanatory `memory-binding:` comment just above it. A value a step needs from the binding reaches it through its payload instead (`asdt-init/clarify` carries `memory_provider`).

**Host skills.** A step's `host_skills:` names skills the HOST assistant may have installed, never ASDT files. For each name, check your own live skill list. If that skill is available, or an obvious equivalent under another name (a frontend or visual-design skill for `frontend-design`), read its skill file and inject the text as a `### HOST SKILL {name}` block under the declared name. Never invoke it in your own context, where it would take over your orchestration. If it is not available, skip it silently. **A host skill is never required: its absence changes nothing and is never an `Open items` entry.** The step file says when it applies and what outranks it. A host skill never overrides the step's contract.

## 6. Summary schema

A step returns a concise Markdown summary, with no envelope and no YAML. The orchestrator writes it to the scratchpad unchanged (§1). Use these parts in this order, omit every one that does not apply, and never add one about the process itself:

```markdown
# {Role}: {change}

{What, in one sentence.}

**Stage:** {stage}
**Verification:** ran: {true|false}, passed: {true|false}. {summary}

## Decisions
## Constraints
## Files hint
## Acceptance criteria
## Risks
## Data model
## API surface
## Open items
```

- **Decisions**: one-line imperatives, with the rejected alternative in parentheses.
- **Files hint**: code anchors, meaning where to look first.
- **Acceptance criteria**: Given/When/Then, at most 5. PM is the authority on acceptance criteria.
- **Risks**: `risk → mitigation`, one line each.
- **Data model**, **API surface**: architect and developer only, when they apply.
- **Open items**: real gaps only. Assumptions carry the `ASSUMED:` prefix.
- A role adds a section of its own only where its step file defines one (the Developer's files changed, the UX/UI designer's screens).

**Stage and Verification** are for staged roles only. `Stage` names the stage the summary holds, and a stage before the last is a plan, not a delivery (the Developer's `spec`, then `implemented`). The gate that ran the checks adds `Verification`. Its summary is the observed result, or why nothing ran, in one line. On a failure, it carries each failing command's trimmed output. Both live only in the scratchpad. Across sessions, git answers what was built, and verification is run again, never recalled.

**The change's files.** A step that reviews a built change (`Stage: implemented`, writing mode) takes its files from git, the authority on what was built (§0), never from one summary section alone. Its git answer is every path `git status --porcelain` shows modified or untracked, plus `git diff --name-only {base}...HEAD`, where `{base}` is the default branch (`origin/HEAD`, else `main` or `master`; skip this half when `HEAD` is that branch). The change's files are the union of that answer with the paths under the Developer summary's `Files changed` and `Prior delivery`, which also say where to read first. A git path that neither list names and the change's code never reaches is unrelated work: set it aside in one `Open items` line. These read-only git commands are inspection, allowed on every step.

**Golden rule: if a part does not change what the consumer does, it does not belong in the summary.**
