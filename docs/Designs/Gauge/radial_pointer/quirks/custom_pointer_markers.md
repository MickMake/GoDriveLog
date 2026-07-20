# Custom radial pointer markers quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.pointer_markers` |
| New Gauge group | `radial_pointer` |
| Paired custom gauge design | `docs/Designs/Gauge/radial_pointer/gauges/custom_radial.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Naming note

This documentation uses `pointer_markers` as the current GoDriveLog realism key.

The same behaviour is also referred to as **witness markers** in older realism/design notes. Within this custom Gauge documentation set, **pointer markers** and **witness markers** are interchangeable names for the same current behaviour unless a document explicitly says otherwise.

## Design intent

This quirk displays retained marker state associated with the gauge reading, such as a minimum, maximum, follower, tell-tale, or damped secondary pointer position.

For the current GoDriveLog `radial` gauge, the behaviour applies to displayed pointer position or angle only. It must not alter the input sensor value, configured ranges, exported values, or logs.

Pointer markers are instrument-realism features, not statistical overlays.

## Configuration contract

Pointer markers are configured under the existing realism surface:

```yaml
realism:
  pointer_markers:
    min: true
    max: true
    average: true
    window: 10m
```

Supported marker keys:

| Key | Meaning |
|---|---|
| `pointer_markers.max` | Tracks the highest/furthest final rendered radial pointer position in the active min/max history mode. |
| `pointer_markers.min` | Tracks the lowest/least final rendered radial pointer position in the active min/max history mode. |
| `pointer_markers.average` | Shows an old-style highly damped pointer marker. It is not a mathematical average. |
| `pointer_markers.window` | Optional positive finite rolling duration for min/max history only. |

## Physical mechanism being imitated

Pointer markers simulate the small secondary markers found on some physical instruments:

- a **maximum pointer** showing the highest reached needle position;
- a **minimum pointer** showing the lowest reached needle position;
- a **damped follower pointer** that trails the main needle movement, like a slow mechanical witness pointer or highly damped secondary indicator.

For radial gauges, these are extra needle-like markers that share the gauge's pointer movement space.

## Source of truth

Pointer markers observe the gauge's final rendered indicator position.

For the current GoDriveLog `radial` gauge, that means the final rendered needle angle after value mapping and any enabled display-only realism effects.

If no realism effects are enabled, that rendered indicator path is equivalent to the mapped source value. If realism effects are enabled, pointer markers follow the realistic rendered behaviour. For example, a radial max marker may capture visible overshoot if the live pointer overshoots.

Pointer markers must not sample raw source values, log/export values, configured ranges, or unrelated dashboard statistics.

## Min and max marker behaviour

### Visual intent

`min` and `max` markers show the extremes reached by the rendered radial pointer.

### Real-world analogue

This simulates tell-tale, witness, or peak-hold markers on physical instruments. These markers are common on pressure gauges, temperature gauges, tachometers, and other instruments where the operator cares about the highest or lowest reached reading.

### History mode

- Without `window`, min/max markers use daily local reset mode.
- With `window`, min/max markers use rolling-window history.
- Rolling-window history should retain data/update ticks or meaningful final rendered position changes, not every unchanged render frame.

### Good result

The min and max markers clearly show the reached extremes without obscuring the live pointer.

### Bad result

The markers track raw source values instead of rendered position, include unrelated log/export data, reset unpredictably, or turn into statistical chart overlays.

## Average / damped marker behaviour

### Visual intent

`average` is an old-style highly damped pointer marker. It follows the live pointer slowly and calmly.

It is not a mathematical average.

### Real-world analogue

This simulates a secondary damped needle or follower pointer with much more inertia than the live indicator. It gives the impression of an older mechanical instrument where a slow pointer shows the general trend while the main pointer shows the current value.

In v3.6, this marker uses a fixed 10 second time constant.

### Good result

The average marker trails the live pointer smoothly and feels like a physical damped follower.

### Bad result

The marker is described as a true average, samples source/log data instead of rendered position, jitters with every raw input tick, or behaves like a graph/statistics overlay.

## Expected visible behaviour

The expected visible effect is one or more marker elements showing remembered positions alongside the live displayed value.

Pointer markers render as part of the physical instrument, not as a separate analytics layer.

## Rendering expectations

Pointer markers render above the live needle and below overlay, glass, bezel, and frame layers.

Expected radial marker assets:

| Marker | Asset | Notes |
|---|---|---|
| Minimum | `needle_min.png` | Same pivot and rotation model as the live radial needle. |
| Maximum | `needle_max.png` | Same pivot and rotation model as the live radial needle. |
| Average / damped follower | `needle_average.png` | Same pivot and rotation model as the live radial needle. |

## Good result

The markers feel like part of the physical instrument: visible, useful, quiet, and clearly tied to the live radial pointer.

## Bad result

The markers look like UI analytics, chart annotations, database-backed statistics, or a separate dashboard layer rather than part of the gauge itself.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `radial` renderer and is documented under the `radial_pointer` Gauge group.

It is not a generic definition of every mechanical witness pointer, tell-tale, min/max register, or statistical marker. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Pointer markers should remain deterministic and should operate on displayed state. They must not mutate source readings or replace the main value mapping.

Min/max marker history should be bounded by the active history mode. Rolling-window history must not become an unbounded frame log.

## Non-goals

This is not `stat_markers`, automatic statistical analysis, logging summary output, a chart overlay, a database-backed report, or a future generic marker subsystem.

## Documentation boundary

This file documents the current GoDriveLog custom quirk design only.

It does not:

- rename the runtime gauge type;
- change package YAML;
- claim generic catalogue coverage;
- record implementation status;
- describe future gauge behaviour as current behaviour.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/Status.md`
