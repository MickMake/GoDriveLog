# Custom bar peg bounce quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.peg_bounce` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom gauge design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk adds bounded rebound when the displayed element contacts an end stop.

For the current GoDriveLog `bar` gauge, the behaviour applies to the displayed level only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Physical mechanism being imitated

Many analogue gauges have physical stop pegs or hard limits. When a needle reaches the stop with momentum, it may make a tiny rebound before settling. A bar gauge does not literally hit a peg, but the same idea can apply visually to the fill edge or reveal extent reaching its display limit.

## Expected visual behaviour

The fill or reveal extent can rebound after contacting its minimum or maximum end of travel. The config key remains `realism.peg_bounce` even though bars do not have literal pegs.

The effect should remain finite, bounded, deterministic, and readable. It should settle rather than create perpetual background motion.

## Good result

The bounce is small, quick, deterministic, and only visible at the configured visual stop.

## Bad result

The display bounces during ordinary in-range movement, passes through the stop, keeps bouncing, or changes source values.

## Applicable current custom gauge

- `bar` under `bar_or_wedge_display`.

Other gauge types may have related conceptual behaviour, but this file only documents the current custom `bar` design.

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
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Designs/RealismBehaviour/radial-peg-bounce.md`
- `docs/Status.md`
- Existing `docs/Designs/RealismBehaviour/*` and `docs/Implementation/RealismBehaviour/*` records where present
