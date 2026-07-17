# Custom segmented percent-threshold gauge

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `segmented` |
| New Gauge group | `bar_or_wedge_display` |
| Documentation role | Custom current GoDriveLog gauge design |
| Runtime code impact | None |

## Design intent

The GoDriveLog `segmented` gauge is a discrete image-selection gauge for percent-threshold displays. It chooses from a sparse set of pre-rendered level images based on the current normalised value.

This old type name is potentially confusing. In the new Gauge taxonomy it is documented under `bar_or_wedge_display` because the current behaviour is stepped level/threshold display, not a general segmented character display.

## Behaviour model

A `segmented` gauge:

- receives a numeric sensor value;
- maps the value to a clamped `0..100` percent;
- discovers available threshold images from filenames containing a `{percent}` placeholder;
- selects the highest discovered threshold reached by the current value;
- applies threshold-gap hysteresis so the visible image does not chatter around threshold boundaries.

## Asset model

The gauge expects sparse percent-threshold images. Missing intermediate thresholds are valid. Visual identity belongs entirely to the image assets.

## Current design boundaries

The GoDriveLog `segmented` type remains a distinct old runtime/config type. This document does not rename it to `bar`.

The new Gauge group mapping is:

```text
Old GoDriveLog type: segmented
New Gauge group:    bar_or_wedge_display
```

## Not current behaviour

This is not the same thing as `segmented_display`. It does not document seven-segment digits, dot-matrix characters, LCD character modules, or individual luminous segment composition.

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

