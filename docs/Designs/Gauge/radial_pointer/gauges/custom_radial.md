# Custom radial gauge

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| New Gauge group | `radial_pointer` |
| Documentation role | Custom current GoDriveLog gauge design |
| Runtime code impact | None |

## Design intent

The GoDriveLog `radial` gauge is a transform gauge. It displays a numeric sensor value by mapping the value onto an angular position and rotating a visible needle or arc-like asset around a configured pivot.

The type name describes renderer behaviour, not visual style. Timber dials, brass needles, neon art, warning labels, and similar styling belong in assets and dashboard layout, not in a separate gauge type.

## Behaviour model

A `radial` gauge:

- receives a sensor value;
- maps that value through the configured value range;
- normalises the result for display;
- converts the normalised value to an angle;
- draws the configured visual layers with the moving pointer element transformed around its pivot.

## Asset model

The radial gauge depends on image assets for its visual identity. The renderer should not infer style from the gauge type. Decorative panels, bezels, tick marks, labels, shadows, and needles are supplied as assets.

## Current design boundaries

The current `radial` design is not a generic catalogue definition of every radial instrument. It is the GoDriveLog renderer model for a value-to-angle pointer gauge.

## Not current behaviour

This design file does not claim support for every physical radial-pointer quirk. Quirks such as damping, stiction, hysteresis, overshoot, peg bounce, pointer markers, needle shadow, and calibration offset are documented separately as behaviour/quirk records when they apply.

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

