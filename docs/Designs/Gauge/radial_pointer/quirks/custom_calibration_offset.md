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

For the current GoDriveLog `radial` gauge, the behaviour applies to displayed pointer angle only. It must not alter the input sensor value, configured range, exported values, or logs.

## Physical mechanism being imitated

Real analogue gauges are not always perfectly calibrated. The needle might be installed slightly off, the mechanism may have a small fixed bias, or the dial artwork and pointer may not line up perfectly.

## Expected visible behaviour

The expected visible effect is a consistent shifted pointer reading. The source value and configured min/max range stay untouched; only the displayed needle position is offset.

## Good result

The gauge looks slightly imperfect while still clearly representing the configured value range.

## Bad result

The offset changes source values, pushes the needle outside sensible visual bounds, or makes the gauge look broken.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `radial` renderer and is documented under the `radial_pointer` Gauge group.

It is not a generic definition of every calibration or compensation mechanism. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Calibration offset is radial-only in the current model. It must remain display-only, deterministic, fixed for a given configuration, and bounded.

## Non-goals

This is not sensor calibration, input correction, logging correction, range remapping, ECU compensation, or changing the gauge's configured min/max scale.

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
- `docs/Status.md`
