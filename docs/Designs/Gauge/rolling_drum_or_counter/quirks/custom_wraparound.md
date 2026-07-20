# Custom odometer wraparound quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| Old realism key | `realism.wraparound` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired custom gauge design | `docs/Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk represents continuous rolling number drums crossing digit-strip boundaries instead of treating each digit as an isolated static image.

For the current GoDriveLog `odometer` gauge, the behaviour applies to displayed wheel or strip state only. It must not alter input sensor values, configured ranges, exported values, or logs.

## Physical mechanism being imitated

A mechanical odometer drum is a continuous wheel, not ten disconnected images. When it passes from `9` to `0`, or from `0` back to `9`, the strip continues around the drum rather than jumping across a flat image list.

This option simulates that continuous cylindrical path.

## Expected visible behaviour

The expected visible effect is a wheel that can roll through the end of a strip and continue from the other side without a visual discontinuity.

## Good result

A rollover looks like one continuous drum motion through the nearest boundary.

## Bad result

The wheel jumps, reverses unexpectedly, rolls the long way around, or briefly shows an impossible digit position.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `odometer` renderer and is documented under the `rolling_drum_or_counter` Gauge group.

It is not a generic definition of every rolling-drum mechanism. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Wraparound is odometer-specific. It should affect only the displayed wheel-strip addressing and must not alter source values, logs, exported values, or configured ranges.

## Non-goals

This is not odometer backlash, gear play, carry drag, route planning, value remapping, or choosing a long-way animation path.

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
- `docs/Status.md`
