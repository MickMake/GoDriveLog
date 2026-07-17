# Custom odometer gauge implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| New Gauge group | `rolling_drum_or_counter` |
| Paired design | `docs/Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md` |
| Runtime code impact | None |

## Current implementation model

The current GoDriveLog `odometer` implementation is a rolling wheel-strip display.

Each wheel uses a strip asset and a clipped window. The displayed value is represented by strip offsets rather than by composing ordinary text.

## Configuration shape

The current documented model includes `odometer.wheels`. Each wheel can define:

- strip asset;
- window position;
- window size;
- optional source alignment offset;
- optional role.

A sub-unit role can map a wheel to tenths without turning the odometer into arbitrary decimal text formatting.

## Movement behaviour

The current public movement modes are:

| Movement | Current behaviour |
|---|---|
| `smooth` | Keeps fractional strip offsets between digit positions. |
| `click` | Snaps strip offsets to digit positions. |

## Rendering approach

The renderer clips each wheel strip through its configured window and draws the visible strip section through the normal dashboard scene path.

## Current limitations and boundaries

The current implementation is flat strip/window rendering. It does not implement full mechanical gearing, backlash, curved drum depth, or advanced odometer physics as part of the base gauge type.


## Documentation boundary

This file records current GoDriveLog implementation behaviour only.

It does not:
- record implementation status;
- describe intended future work as implemented;
- rename runtime package types;
- replace or migrate existing documentation.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.4/ReleasePlan.md`
- `docs/v3.4/ImplementationState.md`
- `docs/v3.4/prompts/v3.4.0-gauge-type-docs.md`

