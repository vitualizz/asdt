# Artifact Loading — Shared Skill

## Purpose

Guide a specialist through retrieving existing upstream artifacts from Engram by topic_key, extracting the relevant fields from each artifact type, and recording absent artifacts in `open_items[]`.

This skill applies ONLY to steps declared with `inputs: []` in `workflow.yaml` — the first artifact-consuming step of a specialist workflow (Developer, QA, Architect, and Security all use it there). A step with declared inputs receives them ALREADY INJECTED as `### INPUT {topic_key}` blocks per `parallel-retrieval.md`, and must never self-fetch them.

---

## Retrieving Artifacts

Artifacts live in Engram, addressed by topic_key in the form `{project}/{change}/{specialist}/{artifact-type}`. One key holds exactly one artifact type, so every lookup resolves unambiguously.

1. List what exists for this change: `mem_search(query: "{project}/{change}", project: "{project}")`.
2. For each result whose topic_key matches an artifact type you need (see the extraction rules below), call `mem_get_observation(id: {id})` to retrieve the full content. When you already know the exact key you want, search for it directly: `mem_search(query: "{project}/{change}/pm/backlog-entry", project: "{project}")`.
3. Apply the extraction rules by artifact type — pull only the listed fields, never the whole payload.
4. If `mem_search` returns nothing, record in `open_items[]`:
   ```
   "No artifacts found in Engram for change '{change}' — proceeding with feature description only"
   ```

**Memory-provider outage fallback**: only if the Engram MCP server is unreachable, scan `.asdt/artifacts/{change}/` for `.yaml` files and treat anything found there as a stale best-effort substitute, recording `"Engram unavailable — loaded artifacts from filesystem fallback at .asdt/artifacts/{change}/; data may be stale"` in `open_items[]`. This is never a primary lookup path: Engram is authoritative, and specialists do not write these files.

---

## Extraction Rules by Artifact Type

### `pm/backlog-entry` (from the PM specialist)

The canonical requirements artifact for a change.

| Field | Where to find it | What to do with it |
|---|---|---|
| `acceptance_criteria` | `payload.acceptance_criteria[]` | Authoritative ACs — refine them, never re-derive an independent set |
| `user_stories` | `payload.user_stories[]` | List each story ID and summary; use as `story_ref` in downstream steps |
| `scope.in` | `payload.scope.in[]` | Constrain the plan to this list |
| `scope.out` | `payload.scope.out[]` | Record any overlap as an `open_items[]` entry |
| `nfrs` | `payload.nfrs[]` | Surface relevant NFRs (performance, security) in the plan |
| `open_questions` | `payload.open_questions[]` | Carry unresolved questions forward into your own `open_items[]` |

### `architect/system-design` (from the Architect specialist)

| Field | Where to find it | What to do with it |
|---|---|---|
| `decisions` | `payload.decisions[]` | Treat as architectural constraints; do not contradict them |
| `components` | `payload.components[]` | Map plan entries to declared components |
| `api_contracts` | `payload.api_contracts[]` | Produce work that satisfies the declared API surface |
| `risk_items` | `payload.risk_items[]` | Surface as `open_items[]` entries when they affect the plan |

### `ux-ui/ux-brief` (from the UX/UI specialist)

| Field | Where to find it | What to do with it |
|---|---|---|
| `user_flows` | `payload.user_flows[]` | Map each flow step to a plan entry |
| `component_specs` | `payload.component_specs[]` | Produce one plan entry per declared component |
| `interaction_notes` | `payload.interaction_notes[]` | Use as constraints for UI work |

---

## Absent Artifact Protocol

When an expected artifact is not found, do NOT stop or error. Follow this protocol:

1. Add a note to `open_items[]` describing what was absent and what was assumed:
   ```yaml
   open_items:
     - "pm/backlog-entry absent — proceeding with inferred scope from feature description"
     - "architect/system-design absent — no architectural constraints applied; flag complex decisions as open_items"
   ```

2. Continue with whatever context is available (`.asdt/knowledge/knowledge.yaml`, visible code, the feature description provided at invocation).

3. Mark any plan entry that depends on an absent artifact with a note in its `rationale`:
   ```
   rationale: "Inferred from feature description — no pm/backlog-entry present to confirm story coverage"
   ```

---

## Summary Guideline

After loading everything you found, produce an internal summary (not written to the artifact) before proceeding:

```
Loaded:
  - pm/backlog-entry: {N} user stories, scope {in/out counts}
  - knowledge.yaml: stack={stack}, conventions={summary}

Missing:
  - architect/system-design
  - ux-ui/ux-brief

open_items to carry forward:
  - "architect/system-design absent — ..."
  - "ux-ui/ux-brief absent — ..."
```

This summary becomes the grounding context for subsequent steps.
