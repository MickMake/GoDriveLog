# GoDriveLog v3.2 implementation state

Status: v3.2.6 Fyne radial rendering in progress
Current target: v3.2.6 radial layer rendering, prepared-frame performance fix, and non-blocking display scene submission
Current branch: v3.2.6-fyne-radial-rendering

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

Radial gauges are routed through the dashboard scene runtime. Fyne radial rendering uses package-owned normalised `pivot.face` and `pivot.needle` coordinates to align a rotated needle frame to the face pivot.

The v3.2.6 Fyne adapter prepares a deterministic 1-degree needle frame set for each unique needle asset and pivot pair. Live updates then select an already-prepared frame and update the existing keyed `canvas.Image` resource instead of decoding, rotating, PNG-encoding, and creating a new resource during normal sweep updates.

The v3.2 display path uses latest-only scene coalescing without blocking harness or sensor event generation on Fyne rendering. Slow display rendering can reduce visible frame rate, but it must not throttle producer cadence.

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
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.5 }

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

- Added a seven-segment scene model for loaded gauge packages.
- Wired dashboard runtime `type: gauge` widgets to load the configured `widget.gauge` package.
- Gauge widgets use the package-owned sensor id to read current sensor state.
- Dashboard scene widgets now carry gauge package identity, package path, scale, status/error, formatted text, static layer parts, digit asset parts, digit positions, and scene-signature data.
- Fyne adapter positioning now honours package-owned part coordinates and widget scale while preserving slot-based positioning for legacy digit widgets.
- Static package layers are included even for non-`ok` sensor states.
- Live digit text, digit backgrounds, character parts, decimal-point parts, and foreground parts are emitted only for `ok` sensor states.
- Scene signatures include package, placement, size, status, text, layer, asset, character, slot, and digit-position data so formatted-output and layout changes are detectable.
- This slice intentionally does not add radial gauge scene model, example assets, inheritance, cluster, or sensor override behaviour.

## v3.2.4 implementation notes

- Hardened the Fyne seven-segment rendering path rather than replacing the existing scene work from v3.2.3.
- Reused Fyne image objects by stable scene/widget/part keys so digit-only resource changes do not rebuild the canvas object tree.
- Kept glass/overlay rendering as an ordered scene part and added coverage that glass remains the last rendered object over live digits.
- Added deterministic speed-quality coverage through `BenchmarkSevenSegmentAdapterUpdate`, which measures repeated seven-segment digit updates with allocations.
- The benchmark is intentionally comparative rather than a hard time threshold because CI and Raspberry Pi hardware vary.

## v3.2.5 implementation notes

- Added radial gauge scene model support for loaded `radial` gauge packages.
- Routed radial gauge packages through dashboard runtime `type: gauge` widgets.
- Dashboard `Widget` scene data now preserves radial face pivot, needle pivot, and calculated angle.
- Dashboard `Part` scene data now preserves radial needle angle, face pivot, and needle pivot for the later renderer.
- Value mapping uses package-owned `value_map` min/max/start/end angles and honours `clamp`.
- Static radial package layers are included even for non-`ok` sensor states.
- Live radial needle parts are emitted only for `ok` sensor states.
- Dashboard scene signatures include radial angle and pivot data so angle and geometry changes are detectable.
- Added runtime coverage for radial gauge widget package loading, angle/pivot preservation, non-ok needle suppression, and angle-based redraw detection.
- This slice intentionally does not add Fyne radial drawing/rotation, example assets, inheritance, cluster, or sensor override behaviour.

## v3.2.6 implementation notes

- Added Fyne adapter support for radial gauge scene parts.
- Static radial layer parts render through the existing ordered dashboard scene part list.
- Needle PNG frames are prepared as deterministic 1-degree rotated frame sets keyed by source asset and normalised needle pivot.
- Normal live radial updates select an already-prepared frame and update the existing keyed `canvas.Image` resource.
- Added non-blocking `LatestSink.SubmitLatest` for display/harness paths so Fyne rendering no longer backpressures the harness or sensor event cadence.
- Display sink stats now expose submitted, rendered, superseded, and render-duration values for diagnostics.
- Needle placement aligns the selected frame pivot with the gauge face pivot using normalised package-owned pivots.
- Existing seven-segment and non-gauge widget rendering paths are preserved.
- Added adapter coverage for radial layer ordering, live needle rendering, omitted non-ok needle behaviour, keyed object reuse, prepared frame-set reuse, and radial needle update benchmarks.
- Added scene-sink coverage for non-blocking latest-only submission, error visibility, render timing stats, and no-backpressure producer benchmarks.
- This slice intentionally does not add gauge package loading changes, dashboard config changes, example gauge packages, sensor overrides, inheritance, clusters, animation, or procedural drawing.

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.2.0 | completed | Planning docs, prompts, repo hygiene, active example/assets normalisation, and v3.0 doc archiving. |
| v3.2.1 | completed | Gauge package loader and tests for `assets/gauges/**/gauge.yaml`. |
| v3.2.2 | completed | Dashboard gauge widget config fields and validation. |
| v3.2.3 | completed | Seven-segment gauge scene model, dashboard runtime package loading, and adapter positioning. |
| v3.2.4 | completed | Fyne seven-segment glass overlay verification, keyed object reuse, and deterministic benchmark coverage. |
| v3.2.5 | completed | Radial gauge scene model, dashboard runtime routing, angle/pivot preservation, and runtime coverage. |
| v3.2.6 | in progress | Fyne radial layer rendering, prepared 1-degree needle frames, non-blocking display scene submission, pivot placement, object reuse, and benchmark coverage. |

## Pending slices

| Version | Status | Next action |
|---|---|---|
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
- next target;
- any changed decisions or deferrals.
