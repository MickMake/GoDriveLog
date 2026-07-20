# Custom radial peg bounce quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.peg_bounce` |
| New Gauge group | `radial_pointer` |
| Paired custom gauge design | `docs/Designs/Gauge/radial_pointer/gauges/custom_radial.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk adds bounded rebound when the displayed element contacts an end stop.

For the current GoDriveLog `radial` gauge, the behaviour applies to the displayed angle only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Physical mechanism being imitated

Many analogue gauges have physical stop pegs or hard limits. When a needle reaches the stop with momentum, it may make a tiny rebound before settling. A bar gauge does not literally hit a peg, but the same idea can apply visually to the fill edge or reveal extent reaching its display limit.

## Expected visual behaviour

The needle can rebound after contacting the configured minimum or maximum stop.

The effect should remain finite, bounded, deterministic, and readable. It should settle rather than create perpetual background motion.

## Good result

The bounce is small, quick, deterministic, and only visible at the configured visual stop.

## Bad result

The display bounces during ordinary in-range movement, passes through the stop, keeps bouncing, or changes source values.

## Applicable current custom gauge

- `radial` under `radial_pointer`.

Other gauge types may have related conceptual behaviour, but this file only documents the current custom `radial` design.

## Non-goals

- bounce away from non-stop positions;
- random jitter;
- unbounded oscillation;
- changing configured min/max values;

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
- `docs/Status.md`
