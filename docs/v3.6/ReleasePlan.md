# v3.6 Release Plan

v3.6 follows the completed v3.5 gauge realism pass.

v3.6 is the pointer marker release. It adds optional physical-style min, max, and average pointer markers for radial and bar gauges without reopening v3.5 realism work or turning pointer markers into statistical overlays.

## Theme

v3.6 pointer markers are instrument-realism features.

They simulate physical front-mounted gauge markers, not data-analysis overlays. The live rendered indicator pushes marker state, and the marker is drawn as an additional visual pointer layer.

Pointer markers observe the final rendered indicator position:

- radial gauges observe the final rendered needle position;
- bar gauges observe the final rendered bar fill/indicator position;
- if no realism is enabled, the final rendered indicator path naturally reflects the mapped source value;
- if realism is enabled, the marker follows the realistic rendered movement, including overshoot, damping lag, bounce, stiction, or other visible behaviour.

Existing gauge mapping and realism remain responsible for producing the final live indicator position. Pointer markers must not replace, bypass, simplify, or become the only radial/bar realism mechanism.

The shared marker engine consumes a gauge-family-neutral normalised rendered position where `0.0` is the rendered visual minimum and `1.0` is the rendered visual maximum after source mapping, clamping, and realism effects. Radial and bar renderers convert that marker position into angle or bar-axis geometry. If the existing renderer already has a better equivalent abstraction, use that abstraction while preserving the same meaning.

## Config shape

Use the `realism.pointer_markers` key.

The documented v3.6 config shape is simple boolean marker flags plus an optional rolling window:

```yaml
realism:
  pointer_markers:
    max: true
    min: true
    average: true
    window: 5m
```

Rules:

- `max`, `min`, and `average` are optional boolean marker enables.
- `window` is optional.
- Unknown keys under `pointer_markers` must fail config loading clearly.
- Long-form object config is not supported in v3.6.
- A top-level `pointer_markers: true` value is invalid because it does not name marker types.

Invalid examples:

```yaml
realism:
  pointer_markers: true
```

```yaml
realism:
  pointer_markers:
    max:
      enabled: true
```

## Pointer marker semantics

### Max marker

The max marker records the furthest/highest final rendered indicator position reached by the live gauge indicator within the active marker history mode.

For radial gauges, this means the relevant final rendered needle position using the existing radial geometry model.

For bar gauges, this means the relevant final rendered bar fill/indicator position using the existing bar orientation/direction model.

If the rendered indicator visibly overshoots or bounces to a higher/further position, the max marker may be pushed to that position.

### Min marker

The min marker records the lowest/least final rendered indicator position reached by the live gauge indicator within the active marker history mode.

If the rendered indicator visibly overshoots or bounces to a lower/lesser position, the min marker may be pushed to that position.

### Average marker

`average: true` enables an old-style average pointer marker.

This is a physical/visual highly damped follower. It is not an arithmetic mean, rolling mean, statistical average, or source-data average.

The average marker follows the final rendered indicator position using frame-rate-independent exponential smoothing with a fixed v3.6 time constant of 10 seconds:

```text
alpha = 1 - exp(-dt / 10s)
average_position = average_position + alpha * (rendered_indicator_position - average_position)
```

The 10 second time constant is intentionally not configurable in v3.6.

## Time and reset behaviour

Pointer marker state is runtime-only in v3.6. Do not add database persistence.

There are two history modes.

### Daily local reset mode

When `window` is absent, min/max markers track the current local day and reset at local midnight.

Rules:

- use the host system local timezone;
- for the Raspberry Pi deployment, this means whatever timezone the Pi is configured to use;
- v3.6 does not add dashboard-level or per-gauge timezone configuration;
- reset at host-local midnight;
- do not replay earlier history if the app starts partway through the day;
- marker state starts unset until the first valid final rendered indicator position;
- no database persistence.

### Rolling window mode

When `window` is present, min/max markers use rolling-window history.

Rules:

- `window` must be a positive finite duration;
- zero, negative, or unparseable windows must fail config loading clearly;
- retained samples are timestamped;
- rolling-window min/max history should retain samples from gauge update ticks or meaningful final rendered position changes;
- implementations should avoid retaining duplicate or unchanged frame samples and should coalesce unchanged or effectively identical rendered positions where practical;
- samples older than the configured window are discarded;
- retained history must stay bounded by the configured window;
- min/max markers are recalculated from the remaining samples;
- if no valid samples remain, the affected marker becomes unset until the next valid sample.

This is marker display state only, not a general-purpose history store, database, or time-series subsystem.

`window` applies to min/max marker history. The average marker remains a live 10 second damped follower and must not become a statistical rolling average.

## Rendering model

Pointer marker behaviour is shared across radial and bar gauges. Rendering is family-specific.

Pointer markers render above the live needle/bar and below the final overlay/glass/bezel/frame layer.

This matches the physical gauge metaphor: old min/max drag markers were often attached to, or visually associated with, the glass/bezel assembly at the front of the gauge. The live indicator pushes the marker, but the marker remains a front-mounted visual element.

### Radial assets

Radial gauges use explicit marker needle PNG assets where provided:

```text
needle_min.png
needle_max.png
needle_average.png
```

Rules:

- use the same pivot/rotation geometry model as the live radial needle;
- rotate marker assets to the marker position;
- render above the live needle;
- render below overlay/glass/bezel layers;
- keep marker assets visually distinct from the live needle.

### Bar assets

Bar gauges use explicit marker PNG assets where provided:

```text
marker_min.png
marker_max.png
marker_average.png
```

Rules:

- place marker assets along the same bar axis as the live fill/indicator;
- respect horizontal, vertical, reversed, and origin/direction configuration;
- render above the live bar/fill/indicator;
- render below overlay/frame/glass layers;
- keep marker assets visually distinct from the live fill/indicator.

Do not replace explicit marker assets with arbitrary procedural shapes unless a gauge family already has a documented renderer convention for that asset role.

## Slice plan

| Slice | Name | Intent |
| --- | --- | --- |
| v3.6.0 | Pointer marker planning docs | Rewrite the v3.6 plan around the locked physical pointer marker spec. |
| v3.6.1 | Shared pointer marker config/state | Parse simple boolean config, reject unsupported shapes, and establish shared marker state. |
| v3.6.2 | Shared min/max marker engine | Track min/max from final rendered indicator positions, including daily reset and optional rolling window. |
| v3.6.3 | Radial pointer marker rendering | Render radial min/max markers with marker needle PNG assets above the live needle. |
| v3.6.4 | Bar pointer marker rendering | Render bar min/max markers with marker PNG assets above the live bar/fill/indicator. |
| v3.6.5 | Average pointer marker engine | Add the shared 10 second damped average pointer marker state. |
| v3.6.6 | Average pointer marker rendering | Render average pointer markers for radial and bar gauges using explicit marker assets. |
| v3.6.7 | Tests, previews, docs checkpoint | Verify config, marker state, reset/window behaviour, average damping, radial rendering, bar rendering, and previews. |

## Deferred to v3.7+

The following are intentionally outside v3.6:

- odometer backlash cleanup;
- broad v3.5 realism implementation audit;
- numeric/segmented display realism candidates;
- indicator realism expansion beyond existing thermal fade behaviour;
- bar realism backlog beyond pointer markers;
- persistent marker state;
- true mathematical/statistical overlays;
- configurable average damping time;
- marker labels/tooltips beyond default marker assets.

## Non-goals

- Do not mutate source values, logs, exports, configured ranges, or input data.
- Do not calculate a true statistical average under `pointer_markers.average`.
- Do not implement persistence in v3.6.
- Do not add a `source` switch.
- Do not ignore rendered overshoot/bounce/stiction when those effects visibly move the live indicator.
- Do not add odometer backlash or other non-marker realism work in v3.6.
- Do not let v3.6 become a random wishlist drawer. The drawer has been moved to v3.7, where at least it can be labelled properly.
