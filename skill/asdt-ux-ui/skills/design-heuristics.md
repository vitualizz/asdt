# Design Heuristics

## Purpose

A check against generic defaults. Every step that loads this fragment examines the design it just
produced for the archetypes below and RECORDS the outcome as a written verdict.

## Generic-Default Archetypes

Six shapes a design falls into when nothing decided otherwise. Each maps onto one token category.
Read each line as a self-check.

**Flat-Neutral Palette** — you have this when one neutral role is reused for background, text and
border, and color never carries meaning or state.

**Single-Weight Type Scale** — you have this when one size and weight pair carries every level, so
heading, body and caption differ only by position.

**Uniform Elevation** — you have this when one shadow value, or none at all, sits on every surface
regardless of how far that surface floats above the page.

**Zero-or-Uniform Radius** — you have this when the same corner radius applies everywhere, so a
container, a control and an input all share one edge language.

**Linear Motion Sameness** — you have this when one duration and one easing curve drive every
transition, whatever that transition means.

**Undifferentiated Spacing Rhythm** — you have this when a single unit separates everything, with no
1x / 2x / 4x progression grouping related elements or parting unrelated ones.

## Differentiation Obligation

Matching an archetype is not a defect. Leaving it unexamined is. For the archetypes your work
touches, do one of exactly two things: name the deliberate differentiation the stated design intent
justifies, or admit the default was kept and say why keeping it is right here. Both are legitimate
outcomes. Silence — landing on a default without anyone noticing — is the only failing one.

## Signature Choice

Do not differentiate everything. Pick ONE dimension to carry the signature and let the rest stay
conventional; differentiation is a scalpel, not a coat of paint. Tie the choice to the design intent
you were given — its tone, its principles, its north star. Breaking all six archetypes at once reads
as noise, not as character.

## Copy as Design Material

Where a content inventory is among your inputs, treat text as a design surface rather than
decoration. Length drives sizing, tone drives type and color, density drives spacing and hierarchy.
Real copy at real length changes the layout it sits in. Where no such input exists, this section is
inert — skip it without comment.

## Recording the Verdict

Record exactly one line, in this form:

```
design-heuristics: <differentiated|default-kept> — <one-line reason naming an archetype>
```

The enum has those two values and no others. The reason MUST name at least one archetype; a reason
naming none is not a verdict.

```
weak:   design-heuristics: differentiated — the design feels distinctive
strong: design-heuristics: differentiated — broke Single-Weight Type Scale so headings carry the weight contrast a dense reading flow needs
```

The strong line names the archetype, the change and the reason it serves. The weak one could
describe any design at all.

Anything you cannot determine from the text inputs alone goes to `open_items` as advisory, never
asserted into a field. `default-kept` is an honest verdict, not a failure — an unexamined default is
the only failing outcome.

## Calibration

On a platform with an established design system and a strong fingerprint, lean on it and expect
`default-kept`. On greenfield work, or where the design intent is explicitly bold, expect a named
signature. Calibrated 2026-08-01; revisit when the platform fingerprint changes.
