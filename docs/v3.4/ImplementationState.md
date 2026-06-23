# GoDriveLog v3.4 implementation state

Status: v3.4.5 segmented gauge implemented
Current target: complete
Current branch: v3.4.5-segmented-gauge

## Purpose

This file records the current implementation state for v3.4. Update it in every v3.4 slice PR.

## Gauge type decision

The gauge type direction is:

```text
numeric    = formatted value rendered through image character slots
radial     = transform gauge: value-to-angle needle/arc movement
odometer   = transform gauge: rolling digit/wheel strip movement
indicator  = image-selection gauge: off/on state
bar        = transform gauge: continuous clip/reveal/fill/movement
segmented  = image-selection gauge: stepped percent-threshold image selection
```

Visual identity belongs to assets. Code should model behaviour only.

Transform gauges currently mean `radial`, `odometer`, and `bar`. They use normalized or formatted values to calculate renderer geometry such as rotation, clipping, reveal bounds, or strip offsets.

Image-selection/composition gauges currently mean `numeric`, `indicator`, and `segmented`. They choose or compose image assets without becoming general geometry-transform systems.

## Non-goals

- No `style` field.
- No `seven_segment` compatibility alias.
- No dot-matrix font/text renderer in this line.
- No merged `bar`/`segmented` supertype.
- No eager loading of all `segmented` percent images.

## Numeric rename

`seven_segment` has been hard-renamed to `numeric`.

The rename is intentionally a hard rename. This project does not need a compatibility layer for old local gauge YAML. If something breaks, it is cheaper to fix the package than to keep a small museum of aliases.

Active code, examples, package YAML, and validation must use `numeric`. Historical docs and changelog entries may still mention `seven_segment`.

The active numeric renderer keeps the existing formatted-value behaviour: character slots, decimal point handling, digit backgrounds/foregrounds, non-ok suppression, and image asset composition.

## Odometer movement model

`odometer` supports `movement: smooth` by default and `movement: click` as the simple mechanical stepped option.

```text
smooth = continuous strip offset between digit positions
click  = stepped movement that snaps to digit positions
```

Do not expand the first odometer slice into easing, inertia, gear backlash, curved depth, or rear-wheel wraparound.

The v3.4.2 implementation adds a flat strip scene model:

- `type: odometer` package validation is active.
- Package config uses `odometer.wheels`, where each wheel declares a strip asset, window position, window size, optional source alignment offset, and optional role.
- `movement: smooth` keeps fractional strip offsets in scene data.
- `movement: click` snaps strip offsets to digit positions.
- `role: sub_unit` maps a wheel to tenths without adding arbitrary decimal formatting.
- Ebiten renders each wheel as a clipped strip subimage through the normal dashboard scene path.

## Indicator state model

`indicator` uses simple image selection with a required `on` package layer and an optional `off` package layer:

```yaml
layers:
  on: on.png
```

Packages may also define an explicit off image:

```yaml
layers:
  off: off.png
  on: on.png
```

The first truth rule is handled from sensor state: a sensor must be `ok` to render `on`. Boolean typed values use their boolean value; otherwise a non-zero numeric value renders `on`. Off and non-ok states render `off` when the optional layer exists; otherwise they draw no state layer between underlay and overlay layers.

The v3.4.3 implementation adds:

- `type: indicator` package validation.
- Required `layers.on` asset validation with optional `layers.off`.
- Scene support that draws underlay layers, the selected state layer when present, then overlay layers.
- Dashboard `type: gauge` routing for indicator packages through the active Ebiten scene path.

## Bar transform model

`bar` uses `value_map` normalization to reveal a level layer from a fixed rectangle:

```yaml
bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [40, 20, 24, 180]
```

The first slice keeps the bar shape deliberately narrow:

- `type: bar` package validation is active.
- `value_map.max` must be greater than `value_map.min`.
- `layers.level` is required.
- `mode`, `axis`, and `origin` are fixed to the first level-reveal configuration.
- `bounds` defines the reveal rectangle within the package artwork using package-space coordinates `[x, y, width, height]`.
- Raw sensor values are normalized through `value_map`; `clamp: true` clamps before normalization, while rendered geometry is always clipped to drawable bounds.
- The scene and Ebiten path clip the level layer from the bottom up as value increases.

Do not add extra bar modes, horizontal fills, or style knobs in this slice.

## Segmented percent model

`segmented` value layers use `{percent}`:

```yaml
layers:
  segments: levels/rpm_{percent:03}.png
```

The renderer discovers files such as:

```text
rpm_000.png
rpm_010.png
rpm_030.png
```

Those files are valid sparse percent thresholds. The renderer selects the highest discovered percent reached by the current normalized value, subject to hysteresis.

Runtime values are normalized and clamped to `0..100`; this is not configurable.

Discovery counts filenames only. Image decoding must stay lazy.

Segmented rules:

- Missing `000` is valid.
- If no `000` image exists, values below the first discovered threshold display no segmented value layer.
- A single threshold file acts as a value-driven overlay: hidden below the threshold, visible at or above it.
- Non-matching files are ignored.
- Files above `100` are ignored with a warning.
- `hysteresis` defaults to `25`.
- `hysteresis` is a percentage of the adjacent threshold gap, not a percentage of the full `0..100` value range.

The v3.4.5 implementation discovers sparse threshold files from the filename pattern, keeps the selected image stable through threshold-gap hysteresis, and routes the result through the v3 dashboard runtime without adding a compatibility alias or a generalized transform mode.

## Baseline dashboard

The v3.4 baseline remains conceptually based on the reusable baseline config:

```text
examples/baseline-dashboard.yaml
```

The current baseline workload remains useful because it exercises numeric displays and radial RPM through the active Ebiten path.

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.4.0 | completed | Planning docs and prompt set for gauge type cleanup and expansion. |
| v3.4.1 | completed | Hard-renamed active `seven_segment` package type to `numeric` in code, validation, dashboard routing, tests, and runnable example package YAML. No compatibility alias was added. |
| v3.4.2 | completed | Added `odometer` package validation, flat wheel-strip scene parts, `smooth` and `click` movement modes, sub-unit wheel support, dashboard routing, Ebiten clipped strip rendering, and focused tests. |
| v3.4.3 | completed | Added `indicator` package validation, required `on` layer with optional `off` layer, two-state scene selection, dashboard gauge routing, and focused tests. |
| v3.4.4 | completed | Added `bar` package validation, required `value_map` normalization, package-space bottom-up clipping, dashboard routing, Ebiten source-rect clipping, and focused tests. |
| v3.4.5 | completed | Added segmented percent-threshold discovery, threshold-gap hysteresis, dashboard routing, and focused package/runtime tests. |

## Pending slices

None. v3.4 implementation work is complete.

## Update rule

Every v3.4 implementation PR must update this file with:

- completed version;
- current branch;
- next target;
- any changed decisions or deferrals.
