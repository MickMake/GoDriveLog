# Custom bar overshoot quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.overshoot` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom gauge design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk lets a displayed element pass the target slightly and then settle back.

For the current GoDriveLog `bar` gauge, the behaviour applies to the displayed level only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Physical mechanism being imitated

Overshoot simulates momentum in a moving indicator. A needle, linkage, or damped mechanism can carry a little past the final reading before returning to rest. In a bar display, it simulates the displayed fill edge carrying slightly past the target extent before settling. This should look like a small mechanical or damped response, not a cartoon spring.

## Expected visual behaviour

The fill or reveal extent may pass the target level and return to the settled level.

The effect should remain finite, bounded, deterministic, and readable. It should settle rather than create perpetual background motion.

## Good result

The movement gives a small sense of momentum without stealing attention.

## Bad result

The display swings too far, oscillates repeatedly, overshoots during tiny changes where it looks silly, or renders outside sensible visual bounds.

## Applicable current custom gauge

- `bar` under `bar_or_wedge_display`.

Other gauge types may have related conceptual behaviour, but this file only documents the current custom `bar` design.

## Non-goals

- continuous oscillation;
- random bounce;
- changing the source sensor value;
- simulation of a full physics engine;

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
