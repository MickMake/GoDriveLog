# Numeric and Segmented Display Realism Candidates

Origin: `docs/v3.7/PlannedFeatures.md`

Status: historical planning / implementation backlog stub

## Canonical behaviour definitions

Numeric and segmented display realism behaviour definitions now live in `docs/RealismBehaviourGuide/`:

- [`per_digit_response_lag`](../RealismBehaviourGuide/per-digit-response-lag.md)
- [`leading_zero_behaviour`](../RealismBehaviourGuide/leading-zero-behaviour.md)
- [`segment_bleed` / `digit_bleed`](../RealismBehaviourGuide/segment-bleed-digit-bleed.md)
- [`ghosting`](../RealismBehaviourGuide/ghosting.md)
- [`uneven_brightness`](../RealismBehaviourGuide/uneven-brightness.md)
- [`load_sag`](../RealismBehaviourGuide/load-sag.md)
- [`thermal_fade`](../RealismBehaviourGuide/thermal-fade.md)

Do not redefine those behaviours here. Use this file only as backlog/planning context for future implementation work.

## Implementation planning notes

- Numeric gauge rendering is image-map driven.
- Prefer future realism behaviours that operate at the digit-slot or displayed-character level rather than trying to interpret image internals.
- `load_sag` should probably start with a display-level model because it best matches shared supply or driver sag.
- `uneven_brightness` should start as stable per-digit-slot brightness variation, not per-segment brightness.
- `segment_bleed` / `digit_bleed` and `ghosting` still need more design before promotion.
- Decimal points complicate bleed/ghosting because the current numeric display model treats DP as a special overlay rather than part of the normal digit character.

## Suggested future implementation tickets

- Implement numeric/segmented `load_sag`.
- Implement numeric/segmented `uneven_brightness`.
- Specify and implement numeric/segmented `per_digit_response_lag`.
- Defer `segment_bleed` / `digit_bleed` and `ghosting` until the display-mask abstraction and config naming are clear.
