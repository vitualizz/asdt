# Content Design — UX/UI Specialist

## Purpose
Catalog the text touchpoints surfaced by the user flows, elaborating the IA's
`content_intent` into illustrative `TextTouchpoint` entries. This produces a content
inventory — NOT final copy. Representative copy is illustrative only; final wording is
decided downstream.

Model note: `sonnet`; tier `moderate|complex`.

## Inputs
- `ux-ui/flows`: interaction sequences, state changes (where text surfaces appear)
- `ux-ui/ia`: information architecture incl. the thin `content_intent` (`copy_direction` + `empty_state`/`error_state` intent hints)

Both arrive injected per the parallel-retrieval contract — consume them directly. Do not fetch
either yourself.

Extract from flows: every state/interaction that surfaces text (CTAs, labels, errors, empty &
confirmation states).
Extract from ia: the `content_intent` thin intent to back-reference via `elaborates_intent_ref`.

## Context budget
flows + ia: ~800–1,200 tokens.

## Processing
1. SCAN the flows for every text surface: call-to-action buttons, field/section labels, error
   messages, empty states, and confirmation states.
2. For each surface, EMIT a `TextTouchpoint` with: a `surface_type` enum
   (`CTA|confirmation|error|empty_state|label`), a `tone_register` enum
   (`formal|neutral|friendly|playful|urgent`), a `criticality` enum
   (`blocking|important|informational`), an ILLUSTRATIVE `representative_copy` string (≤ 120 chars),
   a `length_range` (`{min, max, unit: chars|words}`), a nullable `i18n_note`, and a MANDATORY
   non-null `elaborates_intent_ref` back to an IA `content_intent` intent.
3. DEGRADE any touchpoint whose `surface_type` is not in the enum into the `microcopy[]` bucket —
   never error, serialization always succeeds.
4. Keep `representative_copy` ILLUSTRATIVE — it demonstrates tone and length, it is NOT the final
   copy. Final copy is produced downstream by implementation.

## Output
Produces: `ux-ui/content-inventory`

Persist via mem_save under the output_topic_key in workflow.yaml; return envelope.

Schema:
```yaml
payload:
  touchpoints:
    - id: "TT-001"                # stable id
      surface_type: "CTA"         # enum: CTA|confirmation|error|empty_state|label (UNRECOGNIZED → microcopy[] bucket, NEVER errors)
      tone_register: "neutral"    # enum: formal|neutral|friendly|playful|urgent
      criticality: "important"    # enum: blocking|important|informational
      representative_copy: ""     # string ≤ 120 chars, ILLUSTRATIVE ONLY (not final copy)
      length_range: {min: 0, max: 0, unit: "chars"}   # unit enum: chars|words
      i18n_note: null             # nullable string
      elaborates_intent_ref: ""   # NON-NULL back-ref to an IA content_intent intent (MANDATORY)
  microcopy: []                   # degradation bucket for unrecognized surface_type
  open_items: []
```
