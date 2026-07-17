# Custom radial needle shadow quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.needle_shadow` |
| New Gauge group | `radial_pointer` |
| Paired custom gauge design | `docs/Designs/Gauge/radial_pointer/gauges/custom_radial.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk adds a static visual shadow for a radial gauge needle so the pointer appears to sit above the dial face rather than being painted flat onto it.

For the current GoDriveLog `radial` gauge, the behaviour is display-only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Expected visible behaviour

The expected visible effect is a subtle shadow or offset duplicate of the needle artwork. It should add physical depth without behaving like dynamic parallax, lighting simulation, or gyro-driven movement.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `radial` renderer and is documented under the `radial_pointer` Gauge group.

It is not a generic definition of every radial-pointer shadow effect. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Needle shadow is radial-only in the current model. It should remain deterministic, bounded, and visually subtle.

## Non-goals

This is not dynamic lighting, moving parallax, dashboard illumination state, or a general shadow engine.

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
