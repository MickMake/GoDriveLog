# Custom bar movement policy quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.movement_policy` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom quirk design | `docs/Designs/Gauge/bar_or_wedge_display/quirks/custom_movement_policy.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Runtime code impact | None |

## Naming note

This document uses `movement_policy` because current code uses `Realism.MovementPolicy` for `bar` movement behaviour.

Do not document this as odometer `movement`. Odometer movement has a separate current configuration surface and a separate wheel movement model.

## Current implementation model

Current code treats `bar` movement as a realism policy applied to the displayed fill or reveal extent.

The current audited policy values are `immediate`, `linear`, and `ease_out`, with `immediate` as the normalized default.

The behaviour applies to rendered state only. It must not change source sensor values, persisted log output, exported values, or configured range semantics.

## Configuration boundary

The old GoDriveLog realism key remains `realism.movement_policy`.

This document does not rename that key and does not introduce a new Gauge-tree runtime configuration name.

## Current limitations and exclusions

This is only the current movement-policy surface. It does not imply a nested movement phase model, physics simulation, continuous idle animation, or support for odometer-only movement values.

## Documentation boundary

This file records current GoDriveLog custom quirk implementation behaviour only.

It does not:
- record implementation status;
- describe future gauge work as implemented;
- rename runtime package types;
- replace or migrate existing documentation.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
