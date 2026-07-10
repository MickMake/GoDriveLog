# v3.7 Implementation State

Status: planned slice list complete; implementation not started

Current target: `v3.7.0 release planning docs`

Current branch: `feature/v3.7-gauge-realism`

## Scope

v3.7 is the odometer and image-based numeric display realism pass after the completed v3.6 pointer marker release.

The release adds:

- odometer `backlash`;
- numeric/seven-segment `per_digit_response_lag`;
- numeric/seven-segment `leading_zero_behaviour`;
- numeric/seven-segment `segment_bleed` and `digit_bleed`;
- numeric/seven-segment `ghosting`;
- numeric/seven-segment `uneven_brightness`;
- numeric/seven-segment `load_sag`.

v3.7 extends the existing renderers. It does not replace or redesign them.

## Current decisions

- Follow the KISS principle.
- Keep each realism option small and locally owned.
- Do not create a generic realism engine.
- Do not create a shared numeric-display runtime framework.
- Share a helper only when doing so is plainly simpler than keeping the code local.
- Keep every new option optional, display-only, deterministic, bounded, and subtle.
- Never mutate source values, logs, exports, configured ranges, or input data.
- Keep the existing numeric renderer responsible for formatting, slot assignment, decimal placement, asset lookup, and layer composition.
- Keep the existing whole-image `segmented` renderer separate from image-based numeric/seven-segment digit rendering.
- Do not move existing segmented hysteresis into a new shared display abstraction.
- Keep decimal-point handling explicit and separate from character assets.
- Use explicit image assets and existing scene parts where practical.
- Do not analyse image pixels to infer segment geometry or segment load.
- Implement one slice at a time.

## Scope boundaries

Allowed in v3.7:

- release docs and prompts;
- odometer backlash;
- the listed numeric/seven-segment realism options;
- small feature-local runtime state for lag or ghosting;
- small obvious helpers where reuse is already demonstrated;
- config validation;
- tests and preview packages;
- updates to the Realism Behaviour Guide;
- regression tests for existing odometer, numeric, and segmented gauges.

Not allowed in v3.7:

- generic realism engines;
- new renderer architecture;
- a shared numeric-display state framework;
- merging numeric and segmented gauge implementations;
- unrelated radial, bar, indicator, or pointer-marker work;
- persistent realism state;
- random frame-to-frame variation;
- procedural seven-segment replacement rendering;
- implementing multiple later slices during an earlier slice;
- opportunistic refactors unrelated to the active slice.

## Config ownership

| Key | Gauge families | Notes |
| --- | --- | --- |
| `realism.backlash` | odometer | Bounded slack after direction reversal; must settle exactly on the target. |
| `realism.per_digit_response_lag` | numeric | Small deterministic slot update delays. |
| `realism.leading_zero_behaviour` | numeric | Presentation-only handling of leading digit slots. |
| `realism.segment_bleed` | numeric | Faint inactive seven-segment imagery using explicit assets/layers. |
| `realism.digit_bleed` | numeric | Faint inactive digit-slot imagery using explicit assets/layers. |
| `realism.ghosting` | numeric | Finite fade of the previous character after a slot changes. |
| `realism.uneven_brightness` | numeric | Stable deterministic per-slot brightness variation. |
| `realism.load_sag` | numeric | Subtle brightness reduction based on known displayed segment load. |

Final config shapes belong to their individual slices and must be documented before implementation in that slice.

## Checklist

- [ ] v3.7.0 release planning docs
- [ ] v3.7.1 odometer backlash
- [ ] v3.7.2 per-digit response lag
- [ ] v3.7.3 leading-zero behaviour
- [ ] v3.7.4 segment and digit bleed
- [ ] v3.7.5 ghosting
- [ ] v3.7.6 uneven brightness
- [ ] v3.7.7 load sag
- [ ] v3.7.8 tests, previews, docs checkpoint

## Slice workflow

1. Read this file.
2. Use `Current target` only if it is explicitly set and its checklist item is still unchecked; otherwise find the first unchecked allowed slice.
3. Read `docs/v3.7/ReleasePlan.md`.
4. Read the matching prompt in `docs/v3.7/prompts/`.
5. Inspect the current implementation before designing the slice.
6. Define and document the slice's exact config shape before changing runtime behaviour.
7. Make only that slice's changes.
8. Keep the implementation local to the existing renderer unless a tiny shared helper is demonstrably simpler.
9. Add or update focused tests and previews.
10. Update this checklist, `Current target`, and relevant realism-guide docs.
11. Run the relevant local tests/checks.
12. Commit the completed slice with a clear message.
13. Push the branch to GitHub.
14. Raise a pull request against `main`.

Then enter the review-fix loop:

15. Inspect Codex GitHub review feedback and CI results.
16. If CI fails, review requests changes, or unresolved review comments require code changes:
    - make the smallest safe fixes only;
    - do not refactor unrelated code;
    - rerun relevant tests/checks;
    - commit and push the fixes.
17. Repeat the review-fix loop at most 3 times.

Stop when either:

- CI passes and review is clear enough to merge; or
- the review-fix loop hits the limit and needs human direction.

## Completion rule

When a slice is complete:

- mark only that slice complete;
- advance `Current target` to the next intended unchecked slice;
- do not implement the next slice in the same branch/PR unless explicitly instructed;
- preserve this file as the operational record of the release.
