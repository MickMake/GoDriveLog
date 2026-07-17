# Custom radial gauge implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| New Gauge group | `radial_pointer` |
| Paired design | `docs/Designs/Gauge/radial_pointer/gauges/custom_radial.md` |
| Runtime code impact | None |

## Current implementation model

The current GoDriveLog `radial` implementation treats radial gauges as transform gauges.

The renderer uses the configured value range and gauge geometry to convert a sensor value into a displayed angular position. The moving pointer/needle artwork is then drawn through the active dashboard rendering path.

## Configuration shape

The current model is behaviour-oriented:

- package type remains `radial`;
- visual identity is supplied by assets;
- value mapping determines the displayed position;
- gauge artwork supplies the dial, pointer, ticks, labels, and decoration.

## Rendering approach

The rendering path calculates display state from sensor input, then draws the configured gauge package layers. Pointer rotation is renderer behaviour; decorative style is not.

## Current limitations and boundaries

This implementation record does not claim every radial-pointer quirk is implemented. Quirk behaviours such as damping, hysteresis, stiction, overshoot, peg bounce, pointer markers, needle shadow, and calibration offset are separate implementation records when verified.


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

