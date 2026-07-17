# Custom bar gauge implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| New Gauge group | `bar_or_wedge_display` |
| Paired design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Runtime code impact | None |

## Current implementation model

The current GoDriveLog `bar` implementation treats bar gauges as transform gauges. A numeric sensor value is normalised and used to reveal or clip an active level layer.

## Configuration shape

The current documented shape includes:

- package type `bar`;
- a required level-style asset layer;
- bar configuration for the reveal geometry;
- a value range used to normalise the sensor value.

The first implemented model is intentionally narrow: a level-reveal bar rather than a general bar/graph rendering framework.

## Rendering approach

The renderer clips the active level artwork to the current displayed extent. The artwork supplies the visual style, while the renderer supplies the continuous reveal behaviour.

## Current limitations and boundaries

This file documents the existing `bar` runtime type only. It does not merge the old `segmented` type into `bar`, even though both are now documented under the broader `bar_or_wedge_display` Gauge group.


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

