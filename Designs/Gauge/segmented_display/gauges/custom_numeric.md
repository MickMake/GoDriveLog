# Custom numeric gauge

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `numeric` |
| Previous historical name | `seven_segment` |
| New Gauge group | `segmented_display` |
| Documentation role | Custom current GoDriveLog gauge design |
| Runtime code impact | None |

## Design intent

The GoDriveLog `numeric` gauge displays formatted sensor values using image assets per character slot.

The old `seven_segment` name was hard-renamed to `numeric` because the current renderer is not limited to one physical seven-segment style. In the new Gauge taxonomy it is documented under `segmented_display`, which is broad enough to contain seven-segment, segmented, and dot-matrix-like character display families.

## Behaviour model

A `numeric` gauge:

- receives a sensor value;
- formats the value into display text;
- assigns characters to fixed slots;
- draws matching image assets for each visible character;
- treats decimal points as overlays rather than separate consumed slots;
- uses the supplied asset set to define the visual style.

## Asset model

The numeric gauge is asset-composition based. Character appearance comes from the configured image assets, not from a built-in font or physical segment renderer.

## Current design boundaries

The current GoDriveLog `numeric` type is an image-slot formatted character display. It is not yet an individual-segment composition engine.

The mapping is:

```text
Old GoDriveLog type: numeric
New Gauge group:    segmented_display
```

## Future candidate, not current behaviour

Individual segment composition may become the future implementation model. That is not current behaviour and must not be described in implementation records until code exists and is verified.

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

