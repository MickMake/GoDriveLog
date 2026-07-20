# Custom radial needle shadow quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.needle_shadow` |
| New Gauge group | `radial_pointer` |
| Paired custom quirk design | `docs/Designs/Gauge/radial_pointer/quirks/custom_needle_shadow.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/radial_pointer/gauges/custom_radial.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats `realism.needle_shadow` as a radial-only renderer refinement.

The behaviour affects displayed radial artwork only. It must not change source sensor values, persisted log output, exported values, or configured range semantics.

## Configuration boundary

The old GoDriveLog realism key remains `realism.needle_shadow`.

This document does not rename that key and does not introduce a new Gauge-tree runtime configuration name.

## Current limitations and exclusions

Needle shadow is static display refinement. Do not treat dynamic parallax, lighting mode, gyro input, or generated visual-diff work as part of this current custom quirk.

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
