# Custom odometer drum slop quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| Old realism key | `realism.drum_slop` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired custom gauge design | `docs/Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk represents small mechanical misalignment between odometer drums so digits do not sit perfectly level and centred at all times.

For the current GoDriveLog `odometer` gauge, the behaviour applies to displayed wheel or strip state only. It must not alter input sensor values, configured ranges, exported values, or logs.

## Physical mechanism being imitated

Real mechanical odometer drums are not always perfectly indexed. Wear, manufacturing tolerance, gear lash, and imperfect assembly can leave each wheel sitting a fraction high or low in the viewing window.

This option simulates static per-wheel alignment imperfection. It is not movement, bounce, drift, or direction-change slack.

## Expected visible behaviour

The expected visible effect is subtle per-wheel offset or imperfect settling that makes the odometer look like a physical drum assembly rather than a perfectly aligned digital overlay.

## Good result

Digits look slightly mechanical and imperfect while still being readable.

## Bad result

Digits become hard to read, offsets change between runs, or wheels drift while idle.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `odometer` renderer and is documented under the `rolling_drum_or_counter` Gauge group.

It is not a generic definition of every rolling-drum mechanism. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Drum slop is odometer-specific. It should be subtle, bounded, deterministic, and display-only.

## Non-goals

This is not carry drag, backlash, random shake, drift, movement, bounce, or a general wheel-physics simulation.

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
- `docs/Designs/RealismBehaviour/odometer-drum-slop.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
