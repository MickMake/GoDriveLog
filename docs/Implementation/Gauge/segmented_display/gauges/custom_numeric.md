# Custom numeric gauge implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `numeric` |
| Previous historical name | `seven_segment` |
| New Gauge group | `segmented_display` |
| Paired design | `docs/Designs/Gauge/segmented_display/gauges/custom_numeric.md` |
| Runtime code impact | None |

## Current implementation model

The current GoDriveLog `numeric` implementation renders formatted values through image character slots.

The current code name is `numeric`. Historical documentation may still mention `seven_segment`, but active package YAML and examples use `numeric`.

## Configuration shape

The current model is asset-composition based:

- formatted value behaviour;
- fixed character slots;
- digit/character image assets;
- decimal point handling as overlays;
- asset-defined visual identity.

The renderer does not require the assets to be physically seven-segment, even though that was the historical name.

## Rendering approach

The numeric renderer formats a sensor value, maps each visible character to a configured image asset, and draws the resulting character slots. Decimal points do not consume independent character slots; they overlay the current or preceding slot according to the current numeric rendering rules.

## Current limitations and boundaries

This implementation record does not claim individual-segment rendering.

The current mapping is:

```text
Old GoDriveLog type: numeric
New Gauge group:    segmented_display
```

Individual segment composition is a future candidate only. It must not be recorded as current implementation until the code exists and is verified.


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

