# Custom bar gauge

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| New Gauge group | `bar_or_wedge_display` |
| Documentation role | Custom current GoDriveLog gauge design |
| Runtime code impact | None |

## Design intent

The GoDriveLog `bar` gauge is a transform gauge. It displays a numeric sensor value by revealing, clipping, filling, or moving an active visual layer according to a normalised value.

The type name describes the renderer behaviour. The visual form may look like a column, fuel bar, level window, progress strip, or themed dashboard element depending on the supplied assets.

## Behaviour model

A `bar` gauge:

- receives a numeric sensor value;
- maps the value through a configured range;
- normalises the result;
- applies bar-specific geometry such as axis, origin, and bounds;
- draws the active level layer clipped to the current displayed extent.

## Asset model

The bar gauge uses assets for underlay, active level/fill artwork, and overlay decoration. The renderer should not contain visual style presets.

## Current design boundaries

The current GoDriveLog `bar` type is continuous bar-style reveal behaviour. It is separate from the old GoDriveLog `segmented` type even though both now sit under the broader `bar_or_wedge_display` Gauge group.

## Not current behaviour

This file does not merge `bar` and `segmented` into one runtime type. It documents only the existing `bar` renderer model.

## Documentation boundary

This file documents the current GoDriveLog custom gauge design only.

It does not:
- rename the runtime gauge type;
- change package YAML;
- claim generic catalogue coverage;
- record implementation status;
- describe future renderer work as current behaviour.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.4/ReleasePlan.md`
- `docs/v3.4/ImplementationState.md`
- `docs/v3.4/prompts/v3.4.0-gauge-type-docs.md`

