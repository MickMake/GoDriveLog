# Custom radial damping quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.damping` |
| New Gauge group | `radial_pointer` |
| Paired custom gauge design | `docs/Designs/Gauge/radial_pointer/gauges/custom_radial.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk adds visible smoothing and lag so the display does not jump instantly to every new target value.

For the current GoDriveLog `radial` gauge, the behaviour applies to the displayed angle only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Physical mechanism being imitated

This quirk imitates mass, friction, fluid damping, electrical smoothing, or deliberately damped display response.

## Expected visual behaviour

the needle moves toward the target angle with a damped response rather than teleporting to the new angle.

The effect should remain finite, bounded, deterministic, and readable. It should settle rather than create perpetual background motion.

## Applicable current custom gauge

- `radial` under `radial_pointer`.

Other gauge types may have related conceptual behaviour, but this file only documents the current custom `radial` design.

## Non-goals

- random wobble;
- perpetual idle movement;
- changing the source sensor value;
- changing logged/exported values;

## Relationship to generic catalogue quirks

This file is a GoDriveLog-specific `custom_` quirk record. Generic catalogue quirk files in the same Gauge group describe physical display families more broadly and should not be treated as current implementation documentation.


## Documentation boundary

This file documents current GoDriveLog custom quirk design only.

It does not:

- rename runtime gauge types;
- change package YAML;
- claim generic catalogue coverage;
- record implementation status;
- describe future renderer work as current behaviour.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
- Existing `docs/Designs/RealismBehaviour/*` and `docs/Implementation/RealismBehaviour/*` records where present

