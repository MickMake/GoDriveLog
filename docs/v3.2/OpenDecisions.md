# GoDriveLog v3.2 open decisions

Status: planning

This file records decisions that are open, decided, deferred, or explicitly rejected for v3.2.

## Decision states

- Open: still needs a decision.
- Decided: settled for v3.2.
- Deferred: intentionally not in the current slice or series.
- Rejected: explicitly out of scope.

## Decided

### Gauge package discovery

State: Decided

Gauge packages live under:

```text
assets/gauges/**/gauge.yaml
```

The only required filename is `gauge.yaml`.

Anything under `assets/gauges/` may be arbitrarily named. Directory names do not imply renderer type, sensor type, gauge type, style, family, or inheritance.

The dashboard references the gauge package directory using an asset-root relative path, for example:

```yaml
gauge: gauges/classic/rpm
```

The loader resolves that to:

```text
assets/gauges/classic/rpm/gauge.yaml
```

### Gauge type location

State: Decided

Gauge type is declared inside `gauge.yaml`:

```yaml
type: radial
```

It is not inferred from the directory path.

### Sensor binding

State: Decided

For v3.2, the gauge package owns the sensor binding:

```yaml
sensor: engine_rpm
```

A dashboard `type: gauge` widget does not bind or override the sensor.

Reason: the layout and range of a gauge are usually designed with a specific sensor in mind.

### Widget role

State: Decided

A dashboard `type: gauge` widget places a complete gauge package on the dashboard.

The widget owns:

- dashboard-local id;
- `type: gauge`;
- gauge package path;
- position;
- scale.

The widget does not own:

- sensor binding;
- value range;
- pivots;
- layers;
- renderer type;
- shared image paths.

### File-based reuse

State: Decided

Gauge packages may share image files through relative layer paths, for example:

```yaml
layers:
  background: ../images/bezel.png
  face: ../images/face_dark.png
  needle: ../images/needle_red.png
```

This is not code inheritance. It is file reference reuse.

Layer paths are resolved relative to the `gauge.yaml` directory.

### First renderer target

State: Decided

The first gauge type is:

```yaml
type: radial
```

Radial gauges use layered images and a rotating needle.

### Radial gauge pivot model

State: Decided

Radial gauges use two normalised pivots:

- `pivot.face`: where the needle pivot should land on the gauge face;
- `pivot.needle`: the rotation point inside the needle image.

The renderer rotates the needle around `pivot.needle` and places that rotated pivot at `pivot.face`.

### Non-ok sensor states

State: Decided

Gauge widgets must follow the existing v3 dashboard status semantics.

For non-`ok` states, including missing, unsupported, timeout, parse_error, error, stale, and unknown, the dashboard must not render fake live numeric values or pretend the gauge value is valid.

Do not map non-`ok` to zero, min, midpoint, or stale-looking live values unless a later explicit unavailable-state design is added.

## Deferred

### Gauge inheritance

State: Deferred

No code inheritance, style inheritance, base gauge extension, or YAML merge model in v3.2.

Reuse is done by copying gauge packages and/or using relative file references to shared image files.

### Widget-level sensor override

State: Deferred

No widget-level sensor override in v3.2.

If this is needed later, consider an explicit name such as `sensor_override` rather than overloading `sensor`.

### Clusters

State: Deferred

No separate `clusters:` config layer in v3.2.

For now, dashboard `widgets:` is effectively the inline cluster. Add named clusters only if duplication across dashboards starts hurting.

### Procedural gauge artwork

State: Deferred

No generated ticks, labels, arcs, zones, or procedural gauge drawing in the first v3.2 series.

The first version is image-layer based.

### Multiple gauge types

State: Deferred

Only `type: radial` is targeted first.

Future gauge types can be added later without changing the package discovery rule.

## Open

### Exact scene part structure for radial gauges

State: Open

Need to decide whether radial gauge scene data should:

1. extend the existing generic `Part` structure with radial-specific optional fields; or
2. introduce a more explicit radial gauge scene part representation.

Guidance: prefer the smallest clean change that does not turn `Part` into a junk drawer.

Target slice: v3.2.3

### Fyne rotation implementation details

State: Open

Need to confirm the exact Fyne implementation for rotating and placing a needle image around a normalised pivot.

Target slice: v3.2.4

### Example gauge asset source

State: Open

Need to decide which first example gauge package to add and where its image assets come from.

Target slice: v3.2.5
