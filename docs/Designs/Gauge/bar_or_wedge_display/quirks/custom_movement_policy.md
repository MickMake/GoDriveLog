# Custom bar movement policy quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.movement_policy` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom gauge design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Naming note

This document uses `movement_policy` because that is the current audited runtime surface for the `bar` gauge.

Do not collapse this into the odometer `movement` model. Odometer movement and radial/bar movement policy are related ideas, but they are not the same public configuration shape in current GoDriveLog.

## Design intent

This quirk controls whether the displayed fill or reveal extent jumps immediately or moves toward the target through a simple finite transition.

For the current GoDriveLog `bar` gauge, the behaviour applies to displayed state only. It must not alter the input sensor value, configured range, exported values, or logs.

## Expected visible behaviour

The expected visible effect is bounded value-change motion rather than an unconditional instant redraw. Current documented policy names are simple movement choices, not a nested physics model.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `bar` renderer and is documented under the `bar_or_wedge_display` Gauge group.

It is not a generic definition of every moving bar gauge mechanism, and it is not a promise that all gauge types share the same movement configuration.

## Constraints

Movement policy should remain deterministic, finite, and display-only. It should settle to a stable displayed value and should not introduce perpetual idle motion.

## Non-goals

This is not damping, stiction, overshoot, peg bounce, backlash, power-on sweep, random vibration, or a general animation engine. Those behaviours remain separate custom quirks or future candidates.

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
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
