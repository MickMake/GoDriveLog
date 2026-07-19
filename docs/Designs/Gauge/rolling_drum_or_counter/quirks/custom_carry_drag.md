# Custom odometer carry drag quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| Old realism key | `realism.carry_drag` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired custom gauge design | `docs/Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk represents rollover coupling where a lower odometer drum begins to drag the next higher drum as the lower digit approaches carry.

For the current GoDriveLog `odometer` gauge, the behaviour applies to displayed wheel or strip state only. It must not alter input sensor values, configured ranges, exported values, or logs.

## Physical mechanism being imitated

In a mechanical odometer, the lower drum does not always leave the next drum perfectly untouched until the exact rollover point. The carry mechanism can start to load or nudge the neighbouring drum before the final click.

This option simulates light rollover coupling between adjacent number drums.

## Expected visible behaviour

The expected visible effect is the next wheel beginning to move slightly before or during a rollover, giving the display a mechanical carry interaction.

## Good result

The next wheel looks lightly dragged by the rolling lower wheel, then lands in the correct final digit.

## Bad result

The higher digit moves too early, moves too far, or appears to change value before the rollover is visually justified.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `odometer` renderer and is documented under the `rolling_drum_or_counter` Gauge group.

It is not a generic definition of every rolling-drum mechanism. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Carry drag is odometer-specific. It should be tied to displayed rollover state and must remain deterministic and bounded.

## Non-goals

This is not backlash, free gear play, arbitrary per-digit lag, or independent wheel animation.

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
- `docs/Designs/RealismBehaviour/odometer-carry-drag.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
