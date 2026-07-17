# Custom odometer snap settle quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| Old realism key | `realism.snap_settle` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired custom gauge design | `docs/Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk represents a moving odometer wheel landing into a stable detent or final display position after motion.

For the current GoDriveLog `odometer` gauge, the behaviour applies to displayed wheel or strip state only. It must not alter input sensor values, configured ranges, exported values, or logs.

## Expected visible behaviour

The expected visible effect is a finite settle into the target position rather than an endless oscillation or a perfectly weightless stop.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `odometer` renderer and is documented under the `rolling_drum_or_counter` Gauge group.

It is not a generic definition of every rolling-drum mechanism. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Snap settle is current odometer realism. It should be finite, bounded, deterministic, and display-only.

## Non-goals

This is not full spring physics, random bounce, odometer backlash, or an exposed multi-phase public movement model.


## Documentation boundary

This file documents the current GoDriveLog custom odometer quirk design only.

It does not:
- rename the runtime gauge type;
- change package YAML;
- claim generic catalogue coverage;
- record implementation status;
- describe future odometer behaviour as current behaviour.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`

