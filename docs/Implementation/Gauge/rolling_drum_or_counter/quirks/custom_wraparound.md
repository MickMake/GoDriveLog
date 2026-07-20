# Custom odometer wraparound quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| Old realism key | `realism.wraparound` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired custom quirk design | `docs/Designs/Gauge/rolling_drum_or_counter/quirks/custom_wraparound.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats wraparound as an odometer-only realism option, but the wheel-strip renderer is already circular in practice. Do not infer a separate non-wrap rendering path unless the code is audited and changed.

The behaviour applies to rendered odometer wheel or strip state. It must not change the source sensor value, persisted log output, exported values, or configured range semantics.

## Configuration boundary

The old GoDriveLog realism key remains `realism.wraparound`.

This document does not rename that key and does not introduce a new Gauge-tree runtime configuration name.

## Interaction with odometer movement

Odometer movement is controlled by the odometer movement model. This quirk may affect the displayed wheel state during or around movement, but it should not expose hidden internal movement phases as public YAML configuration.

## Current limitations and exclusions

This is not odometer backlash, gear play, carry drag, or route planning. It only concerns crossing digit-strip boundaries cleanly.

Do not treat `realism.backlash`, per-digit response lag, or future v3.7 odometer work as part of this current custom quirk.


## Documentation boundary

This file records current GoDriveLog odometer quirk implementation behaviour only.

It does not:
- record implementation status;
- describe future odometer work as implemented;
- rename runtime package types;
- replace or migrate existing documentation.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Status.md`

