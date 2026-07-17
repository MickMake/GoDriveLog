# Custom radial calibration offset quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.calibration_offset` |
| New Gauge group | `radial_pointer` |
| Paired custom gauge design | `docs/Designs/Gauge/radial_pointer/gauges/custom_radial.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk adds a small display-only offset between the true mapped reading and the visible radial pointer position.

It simulates a gauge whose needle or dial is slightly misaligned, without changing what the sensor value actually is.

## Expected visible behaviour

The expected visible effect is a consistent shifted pointer reading. The source value and configured min/max range stay untouched; only the displayed needle position is offset.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `radial` renderer and is documented under the `radial_pointer` Gauge group.

It is not a generic definition of every calibration or compensation mechanism. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Calibration offset is radial-only in the current model. It must remain display-only, deterministic, and bounded.

## Non-goals

This is not sensor calibration, input correction, logging correction, range remapping, or ECU compensation.

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
