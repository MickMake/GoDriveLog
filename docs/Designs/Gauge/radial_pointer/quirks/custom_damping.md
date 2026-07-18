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

Damping simulates resistance that smooths movement: viscous damping in an analogue gauge, mass/inertia in a needle mechanism, electrical filtering, or deliberately damped display electronics. In practice, the instrument does not instantly follow every input change. It moves calmly toward the new reading.

## Expected visual behaviour

The needle moves toward the target angle with a damped response rather than teleporting to the new angle.

The effect should remain finite, bounded, deterministic, and readable. It should settle rather than create perpetual background motion.

## Good result

Needles and bars feel weighty and calm during value changes.

## Bad result

The display feels sluggish, never reaches the target, or keeps drifting after the value has settled.

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
- `docs/Designs/RealismBehaviour/radial-damping.md`
- `docs/Status.md`
- Existing `docs/Designs/RealismBehaviour/*` and `docs/Implementation/RealismBehaviour/*` records where present
