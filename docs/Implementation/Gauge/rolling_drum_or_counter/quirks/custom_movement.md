# Custom odometer movement quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| Old configuration field | `odometer.movement` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired custom quirk design | `docs/Designs/Gauge/rolling_drum_or_counter/quirks/custom_movement.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Runtime code impact | None |

## Naming note

This document uses `odometer.movement` because current code treats that field as the odometer movement source of truth.

Do not document odometer movement as `realism.movement_policy`. That policy is obsolete for odometer movement in the current v3.5 direction.

## Current implementation model

Current code accepts odometer movement values `instant`, `linear`, `ease_out`, `bell`, `smooth`, and `click`.

`instant` jumps to the target value. `linear`, `ease_out`, and `bell` describe concrete transition shapes. `smooth` and `click` are recognised but fall back to `instant` in the current audited implementation.

The behaviour applies to rendered wheel state only. It must not change source sensor values, persisted log output, exported values, or configured range semantics.

## Configuration boundary

The old GoDriveLog configuration field remains `odometer.movement`.

This document does not rename that field and does not introduce a new Gauge-tree runtime configuration name.

## Current limitations and exclusions

The internal odometer movement composition may be described as `route -> lead_in -> travel -> settle -> rest`, but that phase model is not public YAML shape in the current documentation boundary.

This is not backlash, drum slop, carry drag, wraparound, or snap settle. Those are separate current custom quirks or future candidates.

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
