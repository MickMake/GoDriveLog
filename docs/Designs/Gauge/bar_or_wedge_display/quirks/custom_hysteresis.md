# Custom bar hysteresis quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.hysteresis` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom gauge design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk adds direction-dependent displayed offset so approach direction can affect the visible reading.

For the current GoDriveLog `bar` gauge, the behaviour applies to the displayed level only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Physical mechanism being imitated

This quirk imitates mechanical friction, elastic memory, magnetic lag, or linkage behaviour where the indicated position depends partly on recent direction of travel.

## Expected visual behaviour

the fill or reveal extent can hold or settle differently depending on whether the value is rising or falling.

The effect should remain finite, bounded, deterministic, and readable. It should settle rather than create perpetual background motion.

## Applicable current custom gauge

- `bar` under `bar_or_wedge_display`.

Other gauge types may have related conceptual behaviour, but this file only documents the current custom `bar` design.

## Non-goals

- changing the source sensor value;
- changing thresholds used by non-display logic;
- random drift;
- generic debounce for all widgets;

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

