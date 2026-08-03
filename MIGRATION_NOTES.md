# ASDT Refactor — Migration Notes

Working document for the 5-phase refactor that replaces the current 12 shared skills / 41 sub-agent steps / persist-everything design with a smaller core. **Phase 1 is complete: it added the new core and changed nothing else.**

## 1. `skill/asdt-core/` is the new core; `skill/asdt-shared/` is deprecated but intact

`skill/asdt-core/protocol.md` is the single mandatory shared skill of the new system. It carries, in one file, what four separate fragments carry today:

| New location | Replaces |
|---|---|
| `asdt-core/protocol.md` §1 Engram contract | `asdt-shared/skills/artifact-loading.md`, `decision-preservation.md` (the dual `mem_save` is gone — one hand-off save, plus one optional journal line) |
| `asdt-core/protocol.md` §2 Intake contract | `asdt-shared/skills/intake-contract.md` (kept near-verbatim; it was already right) |
| `asdt-core/protocol.md` §3 Executor rules | `asdt-shared/skills/executor-header.md` (trimmed; the falsifiability requirement now applies only to steps that read the codebase) |
| `asdt-core/protocol.md` §4 Injection format | `asdt-shared/skills/parallel-retrieval.md` (the cache-ledger prose collapses to "resolve once, inject once") |
| `asdt-core/protocol.md` §5 Hand-off schema | the per-specialist payload schemas scattered across specialist skills |

`skill/asdt-core/references/` holds the OPTIONAL references, consolidated from eight specialist files:

- `scope-definition.md` ← `asdt-shared/skills/scope-definition.md`
- `platform-context.md` ← `asdt-shared/skills/platform-context.md` (Reuse Guard, compact injection format, and degradation kept; the `{{ }}` template section and the line-by-line casuistry dropped)
- `api-design.md` ← `asdt-architect/skills/{api-design,scalability-analysis,architecture-review}.md`
- `owasp.md` ← `asdt-security/skills/{owasp-review,threat-modeling}.md`
- `accessibility.md` ← `asdt-ux-ui/skills/accessibility.md`
- `conventions.md` ← `asdt-developer/skills/{code-generation,test-generation}.md`

**Nothing under `asdt-shared/` or any specialist was deleted, edited, or unwired.** The originals remain the live path until Phase 4 purges them. Until then the two trees coexist and `asdt-shared/` stays authoritative for anything still pointing at it.

Note on `go:embed`: `skill/embedded.go` embeds `SKILL.md asdt-*`, so `asdt-core/` ships automatically, and `installer.SiblingDestName` maps top-level entries verbatim, so it installs to `{SkillsDir}/asdt-core/` with no Go change. Two consequences worth knowing: a file added under `asdt-core/` whose name starts with `_` or `.` will silently NOT ship (that exclusion is what `TestEmbeddedSharedSkillsMatchDisk` guards for the shared tree), and `asdt-core/` has no `workflow.yaml`, so every registry walk skips it exactly as it skips `asdt-shared/`.

## 2. topic_key rename: `{role}/{step}` → `{role}/handoff`

Today each step persists its own artifact under `{project}/{change}/{role}/{step-artifact}` — `architect/system-design`, `developer/dev-spec`, `security/stride-threats`, and so on. Under the new contract there is exactly **one key per role per change**:

```
{project}/{change}/{role}/handoff      title: "{change}/{role}/handoff"   type: "decision"
```

Intermediate artifacts stop being persisted entirely; they live in the orchestrator's context for the duration of the run. The only other write is the optional one-line append to `{project}/journal`.

This rename lands with the specialist rewrites in Phases 2–3, not before. When it lands, every `inputs:` / `output_topic_key:` in the eight `workflow.yaml` files changes with it, and prose in each specialist `SKILL.md` that names an artifact key changes too. Old keys already in Engram are not migrated — a specialist that looks for a hand-off and finds none records `ASSUMED:` and proceeds, which is exactly the degradation path this design already requires.

## 3. Go-side impacts to expect in Phases 2–5

The Go code is untouched by Phase 1 and `make test` stays green. These are the places that WILL need attention as the later phases move files:

**`internal/installer/registry_gen.go`**
- `specialistHeaderFragments` (≈ line 68) hardcodes four paths under `asdt-shared/skills/` and folds them, in that exact order, into every routed specialist's header region. When `specialist-header.md`, `parallel-retrieval.md`, `intake-contract.md`, and `knowledge-recall.md` are replaced by `asdt-core/protocol.md`, this slice and the load-bearing ordering comment above it both change. A missing path here is an install-time error, so it fails loud — good.
- `specialistHeaderBeginMarker` / `specialistHeaderEndMarker` are hand-copied into `skill/embedded_test.go`; the comment says to change both copies together, and that still holds.
- `registryRenderOrder` and `inlineStepsDisplayNames` enumerate the seven routed specialist directories. They only need editing if a specialist is added, removed, or renamed — not merely rewritten.
- `parseRegistry` derives everything else from `*/workflow.yaml`, so it absorbs step-list changes for free. Note `GenerateSpecialistHeader` returns content unchanged when both markers are absent — a specialist that loses its marker region installs with no header and no error.

**`skill/embedded_test.go`**
- `TestEmbeddedSkillTree` asserts `asdt-shared/skills/executor-header.md` is present by literal path, because `agent_adapters.go` (`executorHeaderPath`, line 14) bakes it into generated agent definitions and silently generates nothing when it is absent. Both the constant and this assertion move together when the executor header folds into `protocol.md`.
- `TestEmbeddedSharedSkillsMatchDisk` walks `asdt-shared/skills/` on disk and asserts a both-directions match with the embedded FS. When Phase 4 empties or deletes that directory, this test needs to be repointed at `asdt-core/` or dropped — `os.ReadDir` on a missing directory is a `t.Fatalf`.
- `TestRoutedSpecialistInvariants` greps each routed `SKILL.md` for the literals `"SOLE orchestrator"`, `"FIRST ACTION — self-load the header"`, and both markers. Rewriting specialist headers in Phase 2 breaks this test unless the replacement prose keeps those exact phrases, or the literals are updated in step.

**`internal/installer/registry_drift_test.go`**
- `TestRegistryDrift` derives the canonical roster from the `workflow.yaml` files and asserts it against three mirror sites: the `Specialist Registry` table in `skill/SKILL.md`, the `Tailored Workflow Generation` trivial table in the same file, and `internal/installer/assets/agents-template.md` (which deliberately omits `asdt-pm`; every site omits the non-routable `asdt-init`). It also asserts the inline-steps region in `skill/SKILL.md` is byte-equal to `renderInlineStepsRegion`'s output.
- Consequence: any Phase 2–3 change to a specialist's `execution: inline` vs `subagent` classification, or to the set of routable specialists, must land in the `workflow.yaml` **and** in all three mirrors in the same commit, or this test fails. It checks set membership, routable booleans, and the inline-steps list — never the descriptive cell prose, so wording in those tables is free to change.

**`internal/installer/workflow_models.go`** walks `*/workflow.yaml` the same way and skips directories without one; step renames flow through it without edits.

## 4. Phase status

- **Phase 1 — done.** `asdt-core/protocol.md`, six files in `asdt-core/references/`, this document. Zero existing files modified.
- **Phase 2–3** — rewrite the specialists as declarative orchestrators against `protocol.md`; collapse the step lists; apply the topic_key rename.
- **Phase 4** — purge `asdt-shared/` and the superseded specialist skill files; repoint the Go references listed above.
- **Phase 5** — final pass over the root `SKILL.md` and the registry mirrors.
