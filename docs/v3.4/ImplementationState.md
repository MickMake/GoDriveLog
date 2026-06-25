# GoDriveLog v3.4 implementation state

Status: v3.4.10 dashboard CLI planned
Current target: v3.4.10 dashboard CLI
Current branch: docs/v3.4.10-dashboard-cli

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
- No remote image generation, stock-art downloads, or non-reproducible manual PNG editing for the generated example dashboard tail.

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

The v3.4.5 implementation discovers sparse threshold files from the filename pattern, normalizes raw sensor values to percent before selection, keeps the selected image stable through threshold-gap hysteresis, and routes the result through the v3 dashboard runtime without adding a compatibility alias or a generalized transform mode.

## Generated example dashboard tail

v3.4.6 through v3.4.9 extend v3.4 with generated example dashboards after the gauge behaviour slices are complete.

The example tail is documentation/example work, not a new renderer model:

- v3.4.6 establishes the deterministic procedural asset generation framework.
- v3.4.7 adds the ornate timber / master carpenter dashboard.
- v3.4.8 adds the neon-grid / Tron-style dashboard.
- v3.4.9 adds the steam-scrap / steampunk dashboard.

The v3.4.6 implementation adds:

- `go run ./scripts/generate-example-assets -theme framework-smoke` as the deterministic asset-generation entry point;
- a small reusable drawing helper package under `internal/assets/examplegen` using only the Go standard library;
- committed `framework-smoke` assets under `examples/framework-smoke/assets/`;
- a runnable smoke dashboard config at `examples/framework-smoke/dashboard.yaml`;
- harness coverage proving the generated asset path loads through the active dashboard runtime.

The v3.4.7 implementation adds:

- `go run ./scripts/generate-example-assets -theme ornate-timber` as the ornate timber regeneration command;
- committed `ornate-timber` dashboard-local artwork under `examples/ornate-timber/assets/`;
- runnable ornate timber gauge packages with `gauge.yaml` beside the local assets under `examples/ornate-timber/assets/gauges/`;
- a runnable ornate timber dashboard config at `examples/ornate-timber/dashboard.yaml`;
- gauge-package example coverage for `numeric`, `radial`, `odometer`, `indicator`, `bar`, and `segmented` inside one themed dashboard;
- harness coverage proving the ornate timber example loads through the active dashboard runtime without new renderer behaviour.

The v3.4.7.1 cleanup adds:

- self-contained generated example dashboard directories under `examples/<dashboard_name>/`;
- dashboard configs moved from `examples/dashboards/` to `examples/<dashboard_name>/dashboard.yaml`;
- dashboard-local assets moved from `examples/assets/v3.4/` to `examples/<dashboard_name>/assets/`;
- runtime gauge packages moved from `assets/gauges/v3.4/` to `examples/<dashboard_name>/assets/gauges/`;
- a movement manifest at `docs/v3.4/ExampleLayoutMoves.md`.

The v3.4.8 implementation adds:

- `go run ./scripts/generate-example-assets -theme neon-grid` as the neon-grid regeneration command;
- committed `neon-grid` dashboard-local artwork under `examples/neon-grid/assets/`;
- runnable neon-grid gauge packages with `gauge.yaml` beside the local assets under `examples/neon-grid/assets/gauges/`;
- a runnable neon-grid dashboard config at `examples/neon-grid/dashboard.yaml`;
- gauge-package example coverage for `numeric`, `radial`, `odometer`, `indicator`, `bar`, and `segmented` inside one dark retro-tech themed dashboard;
- harness coverage proving the neon-grid example loads through the active dashboard runtime without new renderer behaviour.

The v3.4.9 implementation adds:

- `go run ./scripts/generate-example-assets -theme steam-scrap` as the steam-scrap regeneration command;
- committed `steam-scrap` dashboard-local artwork under `examples/steam-scrap/assets/`;
- runnable steam-scrap gauge packages with `gauge.yaml` beside the local assets under `examples/steam-scrap/assets/gauges/`;
- a runnable steam-scrap dashboard config at `examples/steam-scrap/dashboard.yaml`;
- gauge-package example coverage for `numeric`, `radial`, `odometer`, `indicator`, `bar`, and `segmented` inside one brass, iron, and salvaged-hardware themed dashboard;
- harness coverage proving the steam-scrap example loads through the active dashboard runtime without new renderer behaviour.

Example dashboard rules:

- Use local deterministic scripts and stable seed/config values.
- Do not use remote image generation.
- Do not download stock art.
- Do not hand-edit opaque generated PNGs as the source of truth.
- Treat source asset dimensions as authoritative.
- Use dashboard/widget config `scale` when a rendered display needs to fit a smaller or larger window.
- Generated digit sets may choose their own cell size, but slot-positioned assets within that set must share it.
- Decimal points are overlays on the current or preceding digit cell and do not consume a separate slot.
- Do not infer one digit set's dimensions from another digit set.
- Do not add runtime `style` fields.
- Keep visual identity in generated assets and dashboard/package layout.
- Keep decorative timber, glow, pipes, rivets, wires, screws, and panels as assets, not renderer features.
- Prefer small, reviewable slices over one giant asset PR with a top hat and a boiler whistle.

## Dashboard CLI tail

v3.4.10 is planned as a dashboard-scoped CLI slice.

Target command tree:

```text
GoDriveLog dashboard [--config <config-file>]
GoDriveLog dashboard run [vehicle-id] [--config <config-file>] [--renderer ebiten]
GoDriveLog dashboard harness [vehicle-id] [--config <config-file>] [--pattern sweep] [--interval 50ms] [--duration 60s] [--renderer ebiten]
GoDriveLog dashboard examples --output <directory> [--config <config-file>] [--vehicle <vehicle-id>] [--theme framework-smoke] [--force]
GoDriveLog dashboard validate [config-file]
GoDriveLog dashboard validate [--config <config-file>]
```

Decision state:

- All active dashboard CLI functions live under `dashboard`.
- `dashboard preview` is intentionally not part of v3.4.10.
- `--config` is optional and uses the same directory search ordering concept as assets.
- `--renderer` is optional and defaults to `ebiten`.
- `dashboard` without a subcommand prints a compact config overview.
- The overview must include vehicle data source type: serial, Bluetooth, Wi-Fi/TCP, fake/test, or unknown.
- Gauge overview rows must keep gauge ID/name, type, source of data, and PID where applicable.
- `dashboard examples` requires `--output` and uses `--force` to suppress overwrite prompts.
- `dashboard harness --pattern sweep` should become gauge-aware while keeping other patterns for later scope decisions.

Gauge-aware harness `sweep` target:

| Gauge type | Sweep behaviour |
|---|---|
| `odometer` | Let `n = 0`; increment from `n - 20` to `n + 20` for 5 seconds, then `n + 20` to `n + 30` for 5 seconds. |
| `numeric` | Same as `odometer`. |
| `radial` | Keep the existing sweep style. |
| `indicator` | Flash on/off for 5 seconds with a 1s duty cycle, then flash on/off for 5 seconds with a 250ms duty cycle. |
| `bar` | Heartbeat pulse at 90 bpm. |
| `segmented` | Same input shape as `radial`. |

## Baseline dashboard

The v3.4 baseline remains conceptually based on the reusable baseline config:

```text
examples/baseline-dashboard.yaml
```

The current baseline workload remains useful because it exercises numeric displays and radial RPM through the active Ebiten path.

The generated example dashboard tail should add richer example coverage for the completed gauge types without replacing the baseline dashboard as the simple renderer verification workload.

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.4.0 | completed | Planning docs and prompt set for gauge type cleanup and expansion. |
| v3.4.1 | completed | Hard-renamed active `seven_segment` package type to `numeric` in code, validation, dashboard routing, tests, and runnable example package YAML. No compatibility alias was added. |
| v3.4.2 | completed | Added `odometer` package validation, flat wheel-strip scene parts, `smooth` and `click` movement modes, sub-unit wheel support, dashboard routing, Ebiten clipped strip rendering, and focused tests. |
| v3.4.3 | completed | Added `indicator` package validation, required `on` layer with optional `off` layer, two-state scene selection, dashboard gauge routing, and focused tests. |
| v3.4.4 | completed | Added `bar` package validation, required `value_map` normalization, package-space bottom-up clipping, dashboard routing, Ebiten source-rect clipping, and focused tests. |
| v3.4.5 | completed | Added segmented percent-threshold discovery, raw-value normalization before selection, threshold-gap hysteresis, dashboard routing, and focused package/runtime tests. |
| v3.4.6 | completed | Added the deterministic example-asset generation entry point, standard-library drawing helpers, committed `framework-smoke` output under `examples/framework-smoke/assets/`, a runnable smoke dashboard config, and harness coverage for the generated example path. |
| v3.4.7 | completed | Added the ornate timber generated dashboard, committed generated theme artwork under `examples/ornate-timber/assets/`, runnable gauge packages under `examples/ornate-timber/assets/gauges/`, a runnable ornate dashboard config, and harness coverage for the themed example path. |
| v3.4.7.1 | completed | Rehomed the generated framework-smoke and ornate-timber example dashboards under self-contained `examples/<dashboard_name>/` directories, including dashboard configs, dashboard-local assets, and co-located gauge packages. |
| v3.4.8 | completed | Added the neon-grid generated dashboard, committed generated theme artwork under `examples/neon-grid/assets/`, runnable gauge packages under `examples/neon-grid/assets/gauges/`, a runnable neon-grid dashboard config, and harness coverage for the themed example path. |
| v3.4.9 | completed | Added the steam-scrap generated dashboard, committed generated theme artwork under `examples/steam-scrap/assets/`, runnable gauge packages under `examples/steam-scrap/assets/gauges/`, a runnable steam-scrap dashboard config, and harness coverage for the themed example path. |

## Pending slices

| Version | Target | Notes |
|---|---|---|
| v3.4.10 | dashboard CLI | Add dashboard-scoped commands for run, harness, examples, validation, and compact config overview. |

## Update rule

Every v3.4 implementation PR must update this file with:

- completed version;
- current branch;
- next target;
- any changed decisions or deferrals.
