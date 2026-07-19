# Custom radial movement policy quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.movement_policy` |
| New Gauge group | `radial_pointer` |
| Paired custom gauge design | `docs/Designs/Gauge/radial_pointer/gauges/custom_radial.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Naming note

This document uses `movement_policy` because that is the current audited runtime surface for the `radial` gauge.

Do not collapse this into the odometer `movement` model. Odometer movement and radial/bar movement policy are related ideas, but they are not the same public configuration shape in current GoDriveLog.

## Design intent

This quirk controls whether the displayed pointer angle jumps immediately or moves toward the target through a simple finite transition.

For the current GoDriveLog `radial` gauge, the behaviour applies to displayed state only. It must not alter the input sensor value, configured range, exported values, or logs.

## Configuration contract

Radial movement policy is configured under the existing realism surface:

```yaml
realism:
  movement_policy: immediate
```

The intended radial policy values are:

| Value | Meaning | Current implementation note |
|---|---|---|
| `immediate` | Display the pointer at the target position without interpolation. | Implemented and used as the normalized default. |
| `linear` | Move from the previous displayed pointer position to the target at constant progress. | Implemented. |
| `ease_out` | Move quickly at first, then settle gently into the target. | Implemented. |
| `bell` | Move with a slow start, faster middle, and slow end. | Desired, but not yet implemented for radial `movement_policy`. |

`bell` belongs in this existing `realism.movement_policy` contract. It must not introduce a separate scalar radial `movement` key.

## Expected visible behaviour

The expected visible effect is bounded value-change motion rather than an unconditional instant redraw. Current documented policy names are simple movement choices, not a nested physics model.

When implemented, every non-immediate movement policy must be:

- deterministic;
- finite;
- display-only;
- bounded in duration;
- settled exactly on the target displayed pointer position when complete.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `radial` renderer and is documented under the `radial_pointer` Gauge group.

It is not a generic definition of every moving radial gauge mechanism, and it is not a promise that all gauge types share the same movement configuration.

## Current implementation note

Implementation status for this quirk is **partially implemented**.

Current code supports `immediate`, `linear`, and `ease_out` for radial `realism.movement_policy`. The desired `bell` policy is not yet accepted or applied for radial movement policy.

Authoritative implementation status belongs in `docs/Status.md`.

## Constraints

Movement policy should remain deterministic, finite, and display-only. It should settle to a stable displayed value and should not introduce perpetual idle motion.

Adding `bell` should extend the existing movement policy contract only. It should not:

- create a new radial `movement` key;
- reuse the odometer movement configuration surface;
- mutate source sensor values;
- alter logging or exports;
- change configured value ranges;
- introduce hidden default movement;
- become a general animation engine.

## Non-goals

This is not damping, stiction, overshoot, peg bounce, backlash, power-on sweep, random vibration, or a general animation engine. Those behaviours remain separate custom quirks or future candidates.

## Documentation boundary

This file documents the current GoDriveLog custom quirk design only.

It does not:
- rename the runtime gauge type;
- change package YAML outside the existing `realism.movement_policy` contract;
- claim generic catalogue coverage;
- replace `docs/Status.md` as the authoritative implementation-status record;
- describe future gauge behaviour as current behaviour.

## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Designs/RealismBehaviour/radial-movement-options.md`
- `docs/Status.md`
