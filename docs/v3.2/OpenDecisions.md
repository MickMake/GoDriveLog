# GoDriveLog v3.2 open decisions

Status: active; v3.2.6 is in progress

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

The current seven-segment examples use the `7Seg` directory name.

The dashboard references the gauge package directory using the configured package path, for example:

```yaml
gauge: assets/gauges/7Seg/amber/4_digit_rpm
```

The loader resolves that to:

```text
assets/gauges/7Seg/amber/4_digit_rpm/gauge.yaml
```

### One gauge per file

State: Decided

A package contains exactly one gauge definition in `gauge.yaml`.

The supported shape is the positive rule: one package directory, one `gauge.yaml`, one instrument definition.

### Gauge type location

State: Decided

Gauge type is declared inside `gauge.yaml`:

```yaml
type: seven_segment
```

or:

```yaml
type: radial
```

It is not inferred from the directory path.

### Gauge packages are complete instrument packages

State: Decided

A gauge package is a complete dashboard instrument package. It can represent a radial gauge, a seven-segment display module, or later another instrument type.

The dashboard widget places the package. The package owns the instrument behaviour.

### First renderer target

State: Decided

The first concrete package type is:

```yaml
type: seven_segment
```

Reason: seven-segment displays reuse known existing digit-display behaviour while proving the new ownership model.

A seven-segment package owns panel/bezel artwork, glass artwork, digit count, digit positions, format, and sensor binding.

### Existing digit sets stay useful

State: Decided

Existing `digit_sets` remain reusable raw glyph artwork.

A seven-segment gauge package is a complete mounted display that can reference a digit set.

Relationship:

```text
digit_sets = reusable raw glyphs
seven_segment gauge = complete display module
widget = placement only
```

### Sensor binding

State: Decided

For v3.2, the gauge package owns the sensor binding:

```yaml
sensor: engine_rpm
```

A dashboard `type: gauge` widget does not bind or override the sensor.

If `sensor` appears on a `type: gauge` widget, validation rejects it.

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
- `format`;
- digit count;
- digit positions;
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
  panel: ../../7Seg4Digits.png
  glass: ../../Glass.png
  background: ../../7SegBack.png
```

This is not code inheritance. It is file reference reuse.

Layer paths are resolved relative to the `gauge.yaml` directory.

Relative paths such as `../` and `../../` are acceptable when they remain inside the asset tree and do not go up and then back down through several unrelated folders.

### Radial gauge pivot model

State: Decided

Radial gauges use two normalised pivots:

- `pivot.face`: where the needle pivot should land on the gauge face;
- `pivot.needle`: the rotation point inside the needle image.

The renderer rotates the needle around `pivot.needle` and places that rotated pivot at `pivot.face`.

### Fyne radial rotation implementation details

State: Decided

The v3.2.6 Fyne adapter prepares deterministic 1-degree rotated needle frame sets for each unique source needle asset and normalised needle pivot pair.

Normal live radial updates must not decode, rotate, PNG-encode, or allocate new rotated resources for every frame. They should select an already-prepared frame and update the existing keyed `canvas.Image` resource only when the quantised frame changes.

The prepared-frame strategy preserves the v3.2.4 keyed canvas object reuse model and avoids returning to full Fyne canvas tree rebuilds.

### Non-ok sensor states

State: Decided

Gauge widgets must follow the existing v3 dashboard status semantics.

For non-`ok` states, including missing, unsupported, timeout, parse_error, error, stale, and unknown, the dashboard must not render fake live numeric values or pretend the gauge value is valid.

Do not map non-`ok` to zero, min, midpoint, or stale-looking live values unless a later explicit unavailable-state design is added.

### Example asset source

State: Decided

Use simple generated/placeholder PNG assets for first example gauge packages.

The assets only need to prove package layout, layer order, digit positions, relative paths, and renderer behaviour. They can be replaced later.

### Seven-segment scene and dashboard runtime structure

State: Decided

The v3.2.3 seven-segment path is generated from:

- a dashboard `type: gauge` widget;
- the widget's `gauge` package path;
- the already-loaded `seven_segment` gauge package;
- dashboard placement (`position`, `scale`);
- the current sensor state for the package-owned sensor.

Dashboard runtime loads the configured gauge package and emits ordinary dashboard scene widgets/parts that the adapter can consume.

Scene data records package identity, package path, gauge type, sensor id, widget placement, scale, package size, status/error, formatted text, static layers, digit parts, digit slots, digit characters, asset paths, and digit positions.

Static package layers are emitted for all sensor states so the mounted instrument can remain visible.

Live digit parts are emitted only for `ok` sensor states. Non-`ok` states must not emit fake numeric characters, digit backgrounds, decimal-point overlays, or foreground digit parts that look like a live value.

Scene signatures include package, placement, size, status, text, layer, asset, character, slot, and digit-position data so output and layout changes are detectable.

The Fyne adapter honours package-owned part coordinates and widget scale. Broader visual polish, examples, and harness work remain later slices.

### Exact scene part structure for radial gauges

State: Decided

Radial gauge scene support extends the existing generic `Part` and `Widget` scene structures with radial-specific optional fields instead of introducing a separate radial scene representation.

The added radial data is limited to package-owned face pivot, needle pivot, and calculated angle. Static layers still use ordinary layer parts. Live needle parts are emitted only for `ok` sensor states.

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

## Open

No open decisions currently block the v3.2.6 Fyne radial renderer slice.
