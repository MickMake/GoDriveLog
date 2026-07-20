# Custom odometer movement quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| Old configuration field | `odometer.movement` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired custom gauge design | `docs/Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Naming note

This document uses `movement` for odometers because current GoDriveLog treats `odometer.movement` as the odometer wheel movement surface.

Do not collapse this into `realism.movement_policy`. The v3.5 design explicitly treats `realism.movement_policy` as obsolete for odometer movement.

## Design intent

This quirk controls how odometer number wheels transition from an old displayed value to a new displayed value.

For the current GoDriveLog `odometer` gauge, the behaviour applies to displayed wheel motion only. It must not alter the input sensor value, configured range, exported values, or logs.

## Expected visible behaviour

The expected visible effect is a finite wheel transition such as an immediate jump, constant-speed roll, eased roll, or bell-shaped roll, depending on the accepted movement value.

The movement model remains display-only and should settle to a stable displayed reading.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `odometer` renderer and is documented under the `rolling_drum_or_counter` Gauge group.

It is not a generic definition of every rolling drum, mechanical counter, gear train, detent, backlash, or carry mechanism. Those remain separate quirks or future candidates.

## Constraints

Odometer movement should remain deterministic, finite, and display-only. It should not expose the internal phase model as public YAML unless a later audited design explicitly changes that boundary.

## Non-goals

This is not drum slop, carry drag, wraparound, snap settle, backlash, random vibration, or a general physics engine.

## Documentation boundary

This file documents the current GoDriveLog custom quirk design only.

It does not:
- rename the runtime gauge type;
- change package YAML;
- claim generic catalogue coverage;
- record implementation status;
- describe future gauge behaviour as current behaviour.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Status.md`
