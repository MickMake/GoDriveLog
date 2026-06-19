# GoDriveLog v3.2 carry-forward notes

Status: planning

This file records constraints, lessons, and decisions from earlier v3 work that still matter for v3.2.

## From v3.1

### Keep config as data

Carry forward.

Gauge packages should be declarative configuration plus image files. Do not add scriptable config, dynamic expressions, YAML inheritance, templating, or hidden code paths in the first pass.

### Keep dashboard code below the sensor/event boundary

Carry forward.

Gauge rendering should consume sensor state and dashboard events. It should not read OBD endpoints, poll sensors, or know about transport details.

### Preserve explicit sensor status semantics

Carry forward.

Known statuses include `unknown`, `ok`, `missing`, `unsupported`, `timeout`, `parse_error`, `error`, and `stale`.

Gauge widgets must not render fake live values for non-`ok` sensor states.

### Existing widget types must keep working

Carry forward.

Do not break or redesign existing image, digit display, bar display, frame gauge, or indicator widgets while adding `type: gauge`.

### Existing digit sets remain useful

Carry forward.

The current `digit_sets` model is good reusable raw glyph artwork. v3.2 should not throw it away.

The new seven-segment gauge package should use digit sets to build complete mounted displays with panel/bezel, glass, digit count, digit positions, `format`, and sensor binding in the package.

### Harness path matters

Carry forward.

Gauge widgets should be testable through the existing v3 dashboard harness path, using fake sensor events through the real runtime path.

### Scene model must remain display-adapter neutral

Carry forward.

Dashboard runtime should build scene data. Fyne-specific image rotation and drawing should remain in the display adapter.

### Display adapter remains below dashboard runtime

Carry forward.

The gauge scene should provide enough information for Fyne to draw the instrument, but not directly call Fyne from dashboard runtime code.

### v3.1.7 and v3.1.8 are deferred, not abandoned

Carry forward.

v3.1.7 dashboard event efficiency and v3.1.8 retirement readiness should be reconsidered after v3.2 gauge work establishes the new visual path.

## v3.2 architecture to preserve

### Gauge packages are self-contained dashboard instruments

A gauge package owns:

- sensor binding;
- gauge type;
- image layers;
- formatting where applicable;
- value mapping where applicable;
- layout geometry such as digit positions or pivots;
- default presentation details if needed.

### Widgets place gauges

A dashboard `type: gauge` widget owns:

- dashboard-local id;
- gauge package path;
- position;
- scale.

It does not own sensor binding in v3.2. If `sensor` appears on a gauge widget, validation rejects it.

### Directory names are not type names

Everything under `assets/gauges/` is user-defined structure.

This is valid:

```text
assets/gauges/7Seg/amber/4_digit_rpm/gauge.yaml
```

This is also valid:

```text
assets/gauges/my_weird_dash/left_big_spinny_thing/gauge.yaml
```

The loader should care that `gauge.yaml` exists. It should not infer meaning from path names.

### Relative asset reuse is allowed

Gauge packages may reference shared images using relative paths from the `gauge.yaml` directory.

Example:

```yaml
layers:
  panel: ../../7Seg4Digits.png
  glass: ../../Glass.png
```

The implementation must prevent paths escaping the asset tree.

Relative paths such as `../` and `../../` are fine when they stay readable and do not go up and then back down through several unrelated folders.

## Things not to carry into v3.2 yet

- Gauge inheritance.
- Named clusters.
- Widget-level sensor override.
- Procedural gauge drawing.
- Multi-skin or theme systems.
- Animation systems.
- Complex layout engines.
- Renderer-specific config outside gauge type handling.

Small first. Useful first. Fancy later, if fancy earns its lunch.
