# Accessibility — Reference

WCAG 2.1 AA baseline. Optional reference for UX/UI and for anyone building an interface. Target AA on every new component; treat AAA as aspirational and record AAA-relevant decisions in `open_items`.

## Semantics

- Use landmarks: `<header>`, `<nav>`, `<main>` (one per page), `<aside>`, `<footer>`.
- Headings `h1`–`h6` form a logical outline — never skip levels.
- `<button>` for actions, `<a>` for navigation. A `<div>` is never an interactive element.
- `<ul>`/`<ol>` for item lists, `<dl>` for term/description pairs, `<table>` only for tabular data with `<thead>`, `<tbody>`, and `scope`.
- Every form control is associated with a `<label>` via `for`/`id` or `aria-labelledby`.

## Keyboard

- Every interactive element is reachable and operable by keyboard alone.
- Tab order matches visual reading order.
- No keyboard trap: `Tab` or `Escape` always gets the user out.
- Components with no native keyboard behavior (custom dropdown, date picker) implement the ARIA APG pattern for their role.
- Never override browser or OS shortcuts.

## ARIA

- `aria-label` when no visible label exists; `aria-labelledby` to point at an existing visible heading; `aria-describedby` for hints and error messages.
- `role` only when no native element provides the right semantics. Do not add roles elements already have (`<button role="button">` is noise).
- `aria-live="polite"` for non-urgent updates, `assertive` only for critical alerts.

## Contrast

| Text type | Minimum ratio |
|---|---|
| Normal text (< 18pt / < 14pt bold) | 4.5:1 |
| Large text (≥ 18pt / ≥ 14pt bold) | 3:1 |
| UI components and graphical objects | 3:1 |
| Disabled states | No requirement — note it in `open_items` |

Any color pair you cannot verify against the project's tokens is an `open_items` entry, never an asserted pass.

## Focus

- A visible focus indicator is required everywhere. Custom indicators must meet 3:1 contrast; the default browser outline is acceptable.
- Modals: move focus inside on open, return it to the trigger on close.
- Route changes and in-page navigation: move focus to the new heading or a skip-link target.
- Toasts never steal focus — announce them through a live region.

## Images and media

- Decorative: `alt=""`. Informative: `alt` conveys the same information as the image.
- Complex images (charts, diagrams): long description via `aria-describedby` pointing at a text block.
- Video needs captions; audio description when visuals carry meaning the audio does not.
- Icon-only buttons need an `aria-label` naming the action.

## Per-component checklist

- [ ] Operates with keyboard alone
- [ ] Contrast passes for every text and UI pair
- [ ] All images carry appropriate `alt`
- [ ] Focus order is logical
- [ ] Inputs have associated labels
- [ ] State changes are announced to screen readers
- [ ] No meaning conveyed by color alone

Record what you verified — `aria_role`, `keyboard_interaction`, `focus_management`, `contrast_token_ref` — on the component itself. Anything you could not verify from text alone goes to `open_items` as advisory, never into a field as an assertion.
