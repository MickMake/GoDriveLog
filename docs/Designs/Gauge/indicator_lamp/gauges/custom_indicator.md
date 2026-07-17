# Custom indicator gauge

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `indicator` |
| New Gauge group | `indicator_lamp` |
| Documentation role | Custom current GoDriveLog gauge design |
| Runtime code impact | None |

## Design intent

The GoDriveLog `indicator` gauge is a two-state image-selection gauge. It displays an off or on state according to the current sensor value and sensor validity.

It maps best to the `indicator_lamp` Gauge group because the common physical expression is a lamp, illuminated legend, warning tell-tale, or simple state indicator. The implementation itself remains asset-driven and does not require the artwork to literally be a lamp.

## Behaviour model

An `indicator` gauge:

- receives a sensor state;
- requires the sensor to be valid before rendering the on state;
- treats boolean `true` as on;
- treats non-zero numeric values as on;
- renders off or no state layer when the value is false, zero, or not valid, depending on configured assets.

## Asset model

The on layer is the essential state asset. An explicit off layer may also be supplied. Underlay and overlay artwork remain normal asset layers.

## Current design boundaries

The current indicator design is binary. It does not define multi-state legends, colour ramps, flashing patterns, fault priority logic, or physical bulb simulation as part of the base gauge type.

## Not current behaviour

Thermal fade, uneven brightness, warm-up, power lifecycle, blinking, and lamp ageing are quirk/behaviour topics, not base `indicator` gauge identity.

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

