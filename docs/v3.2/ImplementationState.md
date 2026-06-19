# GoDriveLog v3.2 implementation state

Status: v3.2.3 seven-segment gauge scene model implemented, pending review
Current target: v3.2.4 Fyne seven-segment renderer
Current branch: v3.2.3-seven-segment-gauge-scene-model

## Purpose

This file records the current implementation state for v3.2. Update it in every v3.2 slice PR.

## Current direction

v3.2 adds self-contained gauge packages while keeping the existing dashboard and widget model.

A gauge package lives at:

```text
assets/gauges/**/gauge.yaml
```

The directory names under `assets/gauges/` are arbitrary. They do not imply gauge type, renderer type, sensor type, or style semantics.

The only required filename is `gauge.yaml`.

## Architecture summary

```text
dashboard
  widgets[]
    type: gauge
    gauge: assets/gauges/.../<package-dir>
    position: [x, y]
    scale: n

gauge package
  gauge.yaml
  image files
  optional shared image files via relative paths
```

For v3.2:

- the dashboard widget places the gauge;
- the gauge package owns the sensor binding;
- the gauge package owns value mapping and/or formatting;
- the gauge package owns visual layers;
- the gauge package owns layout geometry such as digit positions or pivots;
- `sensor` on a `type: gauge` widget is rejected;
- no widget-level sensor override is planned;
- no code inheritance is planned;
- no cluster layer is planned.

## Seven-segment package direction

The first concrete gauge package type is `seven_segment`.

This lets a complete 2, 3, 4, or 5 digit seven-segment display be packaged with its panel/bezel, glass, digit count, digit positions, format, and sensor binding.

The current example directory for seven-segment gauge packages is `assets/gauges/7Seg/`.

Existing `digit_sets` remain useful as reusable raw glyph artwork. A seven-segment gauge package turns those glyphs into a complete mounted dashboard instrument.

Example dashboard widget:

```yaml
widgets:
  - id: rpm
    type: gauge
    gauge: assets/gauges/7Seg/amber/4_digit_rpm
    position: [780, 40]
    scale: 1.0
```

Example package:

```text
assets/
  gauges/
    7Seg/
      7Seg4Digits.png
      Glass.png
      7SegBack.png
      amber/
        4_digit_rpm/
          gauge.yaml
        7Seg0.png
        7Seg1.png
```

Example `gauge.yaml`:

```yaml
id: amber_4_digit_rpm
type: seven_segment
sensor: rpm
format: "%04.0f"

size:
  width: 398
  height: 150

layers:
  panel: ../../7Seg4Digits.png
  glass: ../../Glass.png

digit_set:
  background: ../../7SegBack.png
  characters:
    "0": ../7Seg0.png
    "1": ../7Seg1.png
  spacing: 4

digits:
  count: 4
  positions:
    - [35, 35]
    - [117, 35]
    - [199, 35]
    - [281, 35]
```

## Radial package direction

Radial gauges remain in scope, but they follow after the seven-segment package path proves the package model.

Example `gauge.yaml`:

```yaml
id: simple_radial_rpm
type: radial
sensor: rpm

size:
  width: 512
  height: 512

layers:
  background: ../shared/radial/simple_rpm/bezel.png
  face: ../shared/radial/simple_rpm/face.png
  ticks: ../shared/radial/simple_rpm/ticks.png
  needle: ../shared/radial/simple_rpm/needle.png
  overlay: ../shared/radial/simple_rpm/glass.png

pivot:
  face: { x: 0.5, y: 0.55 }
  needle: { x: 0.5, y: 0.9 }

value_map:
  min: 0
  max: 7000
  start_angle: -135
  end_angle: 135
  clamp: true
```

## v3.2.1 implementation notes

- Added `internal/dashboard/gauges` as the loader package for self-contained gauge packages.
- `LoadPackage` accepts a package directory under `assets/gauges/**` and resolves `gauge.yaml` from that directory.
- The loader parses `seven_segment` and `radial` package types with package-owned fields including `sensor`, `format`, `digit_set`, `digits`, `layers`, `pivot`, and `value_map`.
- Layer and digit image paths resolve relative to the package `gauge.yaml` directory.
- Relative image paths such as `../7Seg0.png` and `../../7Seg4Digits.png` are allowed when the resolved path remains inside the inferred `assets` tree.
- Paths escaping the asset tree are rejected.
- This slice intentionally does not add dashboard widget config, scene model, renderer, Fyne, inheritance, cluster, or example asset behaviour.

## v3.2.2 implementation notes

- Added `type: gauge` to dashboard widget config parsing and validation.
- Gauge widgets use `gauge` for the package directory, `position` for placement, and `scale` for sizing.
- Gauge widgets reject widget-level `sensor` because sensor ownership stays inside the gauge package.
- Gauge widgets also reject widget-level `asset`; existing asset families remain for legacy image, digit, bar, frame, and indicator widgets.
- Gauge paths must be repository-root relative package directories under `assets/gauges/`.
- Removed the duplicate post-`Validate` gauge widget ownership pass after the ownership rule moved into normal widget validation.
- This slice intentionally does not add scene model, renderer, Fyne, package loading from dashboard runtime, inheritance, cluster, or example asset behaviour.

## v3.2.3 implementation notes

- Added a seven-segment scene model for already-loaded gauge packages.
- Scene generation combines a loaded `seven_segment` package, dashboard placement, and current sensor state.
- Scene data includes package id/path/type, package-owned sensor id, widget position, scale, package size, status/error, formatted text, static layer parts, digit asset parts, and digit positions.
- Static package layers are included even for non-`ok` sensor states.
- Live digit text, digit backgrounds, character parts, decimal-point parts, and foreground parts are emitted only for `ok` sensor states.
- Scene signatures include package, placement, size, status, text, layer, asset, character, slot, and digit position data so formatted-output and layout changes are detectable.
- This slice intentionally does not add Fyne drawing, radial gauge scene model, dashboard-runtime package loading, example assets, inheritance, cluster, or sensor override behaviour.

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.2.0 | completed | Planning docs, prompts, repo hygiene, active example/assets normalisation, and v3.0 doc archiving. |
| v3.2.1 | completed | Gauge package loader and tests for `assets/gauges/**/gauge.yaml`. |
| v3.2.2 | completed | Dashboard gauge widget config fields and validation. |
| v3.2.3 | implemented, pending review | Seven-segment gauge scene model and tests. |

## Pending slices

| Version | Status | Next action |
|---|---|---|
| v3.2.4 | not started | Add Fyne seven-segment renderer. |
| v3.2.5 | not started | Add radial gauge scene model. |
| v3.2.6 | not started | Add Fyne radial renderer. |
| v3.2.7 | not started | Add example gauge packages. |
| v3.2.8 | not started | Add harness verification. |
| v3.2.9 | not started | Checkpoint next direction. |

## Deferred v3.1 work

- v3.1.7 dashboard event efficiency is deferred.
- v3.1.8 retirement readiness is deferred.

They are not cancelled. They should be reconsidered at the v3.2.9 checkpoint.

## Update rule

Every v3.2 implementation PR must update this file with:

- completed version;
- current branch;
- current state;
- next target;
- any important implementation notes;
