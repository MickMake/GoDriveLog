# Custom bar stiction quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.stiction` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom gauge design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk makes small displayed changes resist movement until enough change accumulates to break free.

For the current GoDriveLog `bar` gauge, the behaviour applies to the displayed level only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Physical mechanism being imitated

Stiction is static friction: the extra resistance that must be overcome before a resting mechanism starts moving. A sticky needle, linkage, bearing, seal, or sliding display can ignore tiny input changes until enough force builds up to break it free. This option simulates that thresholded release.

## Expected visual behaviour

The fill or reveal extent may hold through small changes, then release to a new displayed level.

The effect should remain finite, bounded, deterministic, and readable. It should settle rather than create perpetual background motion.

## Good result

Tiny value changes may not move the display immediately. When it does move, it makes a small catch-up movement and settles.

## Bad result

The display sticks during large changes, jumps violently, behaves unpredictably, or changes source values.

## Applicable current custom gauge

- `bar` under `bar_or_wedge_display`.

Other gauge types may have related conceptual behaviour, but this file only documents the current custom `bar` design.

## Non-goals

- random sticking;
- permanent jam simulation;
- sensor fault simulation;
- changing the source sensor value;

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
- `docs/Designs/RealismBehaviour/radial-stiction.md and bar-era equivalents`
- `docs/Status.md`
- Existing `docs/Designs/RealismBehaviour/*` and `docs/Implementation/RealismBehaviour/*` records where present
