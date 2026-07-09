# `pointer_markers`

Applies to: radial, bar.

Status: **implemented**.

Config key: `realism.pointer_markers`.

## What it does

Pointer markers record and render reference positions for the gauge's final rendered indicator position.

Supported marker enables:

```yaml
realism:
  pointer_markers:
    min: true
    max: true
    average: true
    window: 10m
```

Supported marker keys:

| Key | Gauge families | Meaning |
|---|---|---|
| `pointer_markers.max` | radial, bar | Tracks the highest/furthest final rendered indicator position in the active min/max history mode. |
| `pointer_markers.min` | radial, bar | Tracks the lowest/least final rendered indicator position in the active min/max history mode. |
| `pointer_markers.average` | radial, bar | Shows an old-style highly damped pointer marker. It is not a mathematical average. |
| `pointer_markers.window` | radial, bar | Optional positive finite rolling duration for min/max history only. |

## What it simulates in real life

Pointer markers simulate the small secondary markers found on some physical instruments:

- a **maximum pointer** showing the highest reached needle/bar position;
- a **minimum pointer** showing the lowest reached needle/bar position;
- a **damped follower pointer** that trails the main needle/bar movement, like a slow mechanical witness pointer or highly damped secondary indicator.

For radial gauges, these are extra needle-like markers. For bar gauges, these are marker ticks or overlays placed along the bar travel.

## Source of truth

Pointer markers are instrument-realism features, not statistical overlays.

They observe the gauge's final rendered indicator position:

- radial gauges observe the final rendered needle position;
- bar gauges observe the final rendered bar fill/indicator position.

If no realism effects are enabled, that rendered indicator path is equivalent to the mapped source value. If realism effects are enabled, pointer markers follow the realistic rendered behaviour. For example, a radial max marker may capture visible overshoot if the live pointer overshoots.

## Min and max marker behaviour

### Visual intent

`min` and `max` markers show the extremes reached by the rendered indicator.

### Real-world analogue

This simulates tell-tale, witness, or peak-hold markers on physical instruments. These markers are common on pressure gauges, temperature gauges, tachometers, and other instruments where the operator cares about the highest or lowest reached reading.

### History mode

- Without `window`, min/max markers use daily local reset mode.
- With `window`, min/max markers use rolling-window history.
- Rolling-window history should retain data/update ticks or meaningful final rendered position changes, not every unchanged render frame.

### Good result

The min and max markers clearly show the reached extremes without obscuring the live indicator.

### Bad result

The markers track raw source values instead of rendered position, include unrelated log/export data, reset unpredictably, or turn into statistical chart overlays.

## Average / damped marker behaviour

### Visual intent

`average` is an old-style highly damped pointer marker. It follows the live indicator slowly and calmly.

It is not a mathematical average.

### Real-world analogue

This simulates a secondary damped needle or follower pointer with much more inertia than the live indicator. It gives the impression of an older mechanical instrument where a slow pointer shows the general trend while the main pointer shows the current value.

In v3.6, this marker uses a fixed 10 second time constant.

### Good result

The average marker trails the live indicator smoothly and feels like a physical damped follower.

### Bad result

The marker is described as a true average, samples source/log data instead of rendered position, jitters with every raw input tick, or behaves like a graph/statistics overlay.

## Rendering expectations

Pointer markers render above the live needle/bar and below overlay/glass/bezel/frame layers.

Expected marker assets:

| Gauge family | Marker assets | Notes |
|---|---|---|
| radial | `needle_min.png`, `needle_max.png`, `needle_average.png` | Same pivot/rotation model as the live radial needle. |
| bar | `marker_min.png`, `marker_max.png`, `marker_average.png` | Placed along the bar axis and respects the bar direction/origin. |

## Good result

The markers feel like part of the physical instrument: visible, useful, and quiet.

## Bad result

The markers look like UI analytics, chart annotations, database-backed statistics, or a separate dashboard layer rather than part of the gauge itself.
