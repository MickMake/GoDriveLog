# GoDriveLog v3.4 implementation state

Status: v3.4.2 odometer planning/scene model implemented
Current target: v3.4 gauge/display package cleanup and expansion
Current branch: codex/v3.4.2-odometer-gauge

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

## Pending slices

| Version | Status | Next action |
|---|---|---|
| v3.4.3 | not started | Add indicator gauge behaviour. |
| v3.4.4 | not started | Add first bar gauge transform behaviour. |
| v3.4.5 | not started | Add segmented percent-threshold image discovery, threshold-gap hysteresis, and rendering. |

## Update rule

Every v3.4 implementation PR must update this file with:

- completed version;
- current branch;
- next target;
- any changed decisions or deferrals.
