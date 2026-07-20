# Custom radial calibration offset quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.calibration_offset` |
| New Gauge group | `radial_pointer` |
| Paired custom quirk design | `docs/Designs/Gauge/radial_pointer/quirks/custom_calibration_offset.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/radial_pointer/gauges/custom_radial.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats `realism.calibration_offset` as a radial-only display offset.

The behaviour affects displayed radial pointer position only. It must not change source sensor values, persisted log output, exported values, or configured range semantics.

## Configuration boundary

The old GoDriveLog realism key remains `realism.calibration_offset`.

This document does not rename that key and does not introduce a new Gauge-tree runtime configuration name.

## Current limitations and exclusions

Calibration offset is not input calibration. Do not use this quirk to document source-value correction, range remapping, or future sensor-compensation work.

## Documentation boundary

This file records current GoDriveLog custom quirk implementation behaviour only.

It does not:
- record implementation status;
- describe future gauge work as implemented;
- rename runtime package types;
- replace or migrate existing documentation.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Status.md`
