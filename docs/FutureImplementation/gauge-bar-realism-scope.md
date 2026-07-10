# Bar Realism Scope

Origin: `docs/v3.7/PlannedFeatures.md`

Status: implementation backlog stub

## Canonical behaviour definitions

Bar gauge realism behaviour definitions live in `docs/RealismBehaviourGuide/`:

- [`damping`](../RealismBehaviourGuide/damping.md)
- [`hysteresis`](../RealismBehaviourGuide/hysteresis.md)
- [`stiction`](../RealismBehaviourGuide/stiction.md)
- [`overshoot`](../RealismBehaviourGuide/overshoot.md)
- [`peg_bounce`](../RealismBehaviourGuide/peg-bounce.md)
- [`pointer_markers`](../RealismBehaviourGuide/pointer-markers.md)
- [`stepped_fill`](../RealismBehaviourGuide/stepped-fill.md)
- [`quantized_fill`](../RealismBehaviourGuide/quantized-fill.md)

Do not redefine those behaviours here. Use this file only as backlog/planning context for future bar realism implementation work.

## Implementation planning notes

- Bar gauges are linear fill/reveal gauges.
- Runtime realism should focus on the displayed fill edge moving toward the target, not on repainting the gauge artwork.
- Before planning any bar realism beyond the already implemented options, audit the current code and completed v3.5/v3.6 docs.
- Do not contradict completed implementation state from backlog notes alone.
- `stepped_fill` and `quantized_fill` both need a clear config model before promotion.

## Suggested future implementation tickets

- Specify and implement `stepped_fill` for block-style bar/segmented displays.
- Specify and implement `quantized_fill` where the bar only visibly changes after crossing a display-resolution step.
- Audit any already-documented bar realism option only if current `main` behaviour is suspected to be missing or wrong.
