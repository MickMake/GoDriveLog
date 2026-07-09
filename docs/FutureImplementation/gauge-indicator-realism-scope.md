# Indicator Realism Scope

Origin: `docs/v3.7/PlannedFeatures.md`

Status: implementation planning note

## Canonical behaviour definition

The indicator `thermal_fade` behaviour definition lives in [`../RealismBehaviourGuide/thermal-fade.md`](../RealismBehaviourGuide/thermal-fade.md).

Do not redefine indicator fade behaviour here. Use this file only as backlog/planning context for future indicator realism work.

## Implementation planning notes

- Indicator gauges are image-state driven.
- The `off` and `on` image layers define the static lamp appearance.
- Runtime realism should stay transition-focused unless a later design explicitly adds more indicator display states or display-layer effects.
- `thermal_fade` already supports separate rise/fall timing with `rise_ms` for off-to-on warm-up and `fall_ms` for on-to-off cool-down.
- Weak bulb, tint, ageing, bloom, dirty lens, and uneven illumination usually belong in the supplied images rather than runtime behaviour.

## Suggested future implementation tickets

- Only add more indicator runtime realism if a specific behaviour is promoted from [`../RealismBehaviourGuide/imperfections.md`](../RealismBehaviourGuide/imperfections.md).
- Keep static lamp appearance in image assets unless a future rendering design deliberately introduces display-layer effects.
