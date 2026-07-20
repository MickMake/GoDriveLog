# Custom indicator thermal fade quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `indicator` |
| Old realism key | `realism.thermal_fade` |
| New Gauge group | `indicator_lamp` |
| Paired custom quirk design | `docs/Designs/Gauge/indicator_lamp/quirks/custom_thermal_fade.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/indicator_lamp/gauges/custom_indicator.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats `realism.thermal_fade` as an indicator realism option that affects displayed lamp brightness over state changes.

The behaviour applies to rendered indicator state only. It must not change source sensor values, persisted log output, exported values, or configured threshold semantics.

## Configuration boundary

The old GoDriveLog realism key remains `realism.thermal_fade`.

This document does not rename that key and does not introduce a new Gauge-tree runtime configuration name.

## Current limitations and exclusions

Thermal fade is not a power lifecycle model. Do not treat random flicker, brownout dip, lamp failure, or dashboard illumination mode as part of this current custom quirk.

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
