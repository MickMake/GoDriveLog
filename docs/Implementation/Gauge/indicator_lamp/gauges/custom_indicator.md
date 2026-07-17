# Custom indicator gauge implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `indicator` |
| New Gauge group | `indicator_lamp` |
| Paired design | `docs/Designs/Gauge/indicator_lamp/gauges/custom_indicator.md` |
| Runtime code impact | None |

## Current implementation model

The current GoDriveLog `indicator` implementation is a binary image-selection gauge.

It selects an on or off visual state using the sensor value and sensor validity.

## Configuration shape

The current package type remains `indicator`.

The model requires an `on` state layer. An `off` state layer may also be supplied. Underlay and overlay layers remain normal gauge package artwork.

## State selection behaviour

The current state rule is:

- sensor must be valid before the on state is rendered;
- boolean `true` renders on;
- non-zero numeric values render on;
- false, zero, or non-valid values render off if an off layer exists;
- if no off layer exists, no state layer is drawn between underlay and overlay for the off/non-valid state.

## Rendering approach

The renderer draws underlay layers, then the selected state layer when present, then overlay layers.

## Current limitations and boundaries

This file documents a two-state indicator. It does not implement multi-state annunciators, blinking policies, thermal fade, or power lifecycle behaviour as base indicator behaviour.


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

