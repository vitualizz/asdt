# Decision Preservation — Shared Skill

## Purpose

After producing a significant decision or recommendation, ensure a **permanent
organizational knowledge record** is written — distinct from the change-scoped
artifact. This is what transforms ASDT from a per-change tool into a knowledge
system that accumulates domain expertise over time.

## When to Use

Any step that produces a **final decision artifact**:

- Architect: `technical-handoff` → `architectural-decision` (the earlier `decision-record` step produces the intermediate ADR, not this final artifact)
- Developer: `implement` → `dev-implementation`
- Security: `hardening-checklist` → `hardening-checklist`
- QA: `quality-report` → `quality-report`
- Any specialist's **last step** that produces a recommendations artifact

## Protocol

1. **Extract What**: the decision or recommendation in one clear sentence.
2. **Extract Why**: the driving constraint, user request, or force that made this
   decision necessary.
3. **Extract Where**: the change ID and the Engram topic_key for the primary artifact
   (e.g. `topic_key: "{project}/{change}/architect/handoff"`).
4. **Extract Learned** (optional): a gotcha, rejected alternative, or edge case
   worth remembering for future changes.
5. **Shape as Entry**:
   - `Title`: `"{specialist-role}: {change-name}"` (e.g. `"architect: add-auth"`)
   - `Type`: `architecture` for design decisions, `decision` for policy/choice
   - `Content.What`: the extracted What sentence
   - `Content.Why`: the extracted Why
   - `Content.Where`: the extracted Where path
   - `Content.Learned`: optional

6. **Include a `summary` field** in your output artifact payload (≤ 150 tokens):
   ```yaml
   summary: "Chose JWT over sessions for stateless auth — enables horizontal scaling"
   ```
   You MUST call `mem_save` yourself with this summary as the permanent record —
   there is no runner; persistence is your responsibility.

## Dual mem_save semantics

This `mem_save` is the SECONDARY permanent knowledge record — it uses the title pattern
`"{specialist-role}: {change-name}"` and is NOT stored at the step's `output_topic_key`.
The PRIMARY canonical envelope save was already performed by the calling step itself, under
the `output_topic_key` declared in `workflow.yaml` — that is the artifact sub-agents retrieve
via their declared `inputs:`. These two saves are DISTINCT and MUST NOT be merged: collapsing
them would break artifact retrieval at the canonical topic_key. This applies to every specialist
whose final step invokes this skill, not just one.

## Context Budget

No added input budget — this skill operates on the artifact you already produced.
The `summary` field it adds to your payload is the only output overhead.

## Output Schema

Your final artifact payload MUST include:

```yaml
summary: string  # ≤150 tokens — the decision in one sentence
```

Your `mem_save` call persists this as a permanent record (type `architecture` or
`decision`) so future specialists can query it via `knowledge-recall`.
