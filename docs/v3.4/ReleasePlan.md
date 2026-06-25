# GoDriveLog v3.4 release plan

Status: dashboard CLI tail planned
Owner: gauge package implementor

## Purpose

v3.4 defines the next gauge package model for the active Ebiten dashboard path.

The release direction is gauge/display package cleanup and expansion. It does not introduce a platform packaging track.

The original behaviour implementation slices are complete through v3.4.5. The v3.4.6 through v3.4.9 tail adds generated example dashboards that prove the completed gauge types with repeatable assets. The v3.4.10 through v3.4.12 tail now splits dashboard CLI work into full command-tree routing, compact overview output, and gauge-aware harness sweep behaviour.

## Release goal

Create a clear gauge type model where type names describe renderer behaviour and visual identity comes from assets.

The planned gauge types are:

| Type | Purpose |
|---|---|
| `numeric` | Formatted value rendered through image assets per character slot. Replaces `seven_segment`. |
| `radial` | Transform gauge using value-to-angle mapping for a needle or arc. Already exists. |
| `odometer` | Transform gauge using rolling digit/wheel strip movement. |
| `indicator` | Two-state image-selection gauge: off/on. |
| `bar` | Transform gauge where normalized value clips, reveals, fills, or moves an asset layer. |
| `segmented` | Discrete stepped image-selection gauge using sparse percent-threshold images. |

## Final release principles

- Gauge type describes behaviour, not visual style.
- Visual style is asset-only.
- Do not add a `style` field.
- Rename `seven_segment` to `numeric` with no compatibility alias.
- Active code, examples, package YAML, and validation must use `numeric`; historical docs and changelog entries may still mention `seven_segment`.
- Keep `radial` as the existing radial gauge type.
- Treat `radial`, `odometer`, and `bar` as transform gauges.
- Treat `numeric`, `indicator`, and `segmented` as image-selection/composition gauges.
- Keep `bar` and `segmented` separate.
- Keep each implementation slice small.
- Do not add all future gauge families in one monster PR. That beast would need feeding and possibly a small helmet.

## Gauge type boundaries

| Type | Behaviour |
|---|---|
| `numeric` | Format value, split into character slots, draw matching character assets. |
| `radial` | Normalize value and rotate a needle around a pivot. |
| `odometer` | Convert value into rolling wheel strip offsets and draw clipped strip assets. |
| `indicator` | Select one of two state layers. |
| `bar` | Continuously reveal, clip, fill, or move an active layer from normalized value. |
| `segmented` | Select the highest discovered percent-threshold image reached by the current value, with threshold-gap hysteresis. |

## Odometer movement rule

`odometer` supports two movement modes:

```yaml
odometer:
  movement: smooth
```

`movement` defaults to `smooth`.

```text
smooth = continuous strip offset between digit positions
click  = stepped mechanical movement that snaps to digit positions
```

Do not add easing, inertia, gear backlash, curved depth, or rear-wheel wraparound in the first odometer slice. Keep the tiny clockwork goblin unemployed.

## Segmented percent layer rule

`segmented` uses a percent placeholder in its value layer:

```yaml
layers:
  segments: levels/rpm_{percent:03}.png
```

Files such as these are valid:

```text
levels/rpm_000.png
levels/rpm_010.png
levels/rpm_030.png
levels/rpm_040.png
levels/rpm_100.png
```

The renderer discovers matching files, extracts percent thresholds, sorts them, and draws the highest threshold reached by the current normalized percent.

Runtime values are always normalized and clamped to `0..100`; this is not configurable.

Discovery must count filenames only. It must not decode every image. Runtime should lazy-load selected images and cache recent images.

Segmented threshold rules:

- Sparse thresholds are valid.
- Missing `000` is valid.
- If no `000` image exists, values below the first discovered threshold display no segmented value layer.
- A single valid threshold file acts as a value-driven overlay: nothing below the threshold, the image at or above the threshold.
- Files that do not match the `{percent}` pattern are ignored.
- Files above `100` are ignored with a warning.
- Hysteresis defaults to `25`.
- `hysteresis` is a percentage of the adjacent threshold gap, not an absolute percentage of the full `0..100` range.

Example:

```yaml
segmented:
  hysteresis: 25
```

With thresholds `25` and `50`, the downward hysteresis gap from `50` is `(50 - 25) * 25% = 6.25`. The selected `050` image remains active until the value drops below `43.75`.

## Generated example dashboard tail

After the behaviour slices are complete, v3.4 adds generated example dashboards. These examples must use deterministic local scripts and committed source/config/docs only. They must not depend on remote image generation, downloaded art, or hand-edited opaque assets.

The example dashboard tail has four slices:

| Version | Slice | Result |
|---|---|---|
| v3.4.6 | example asset generation framework | Add deterministic procedural asset generation structure, conventions, docs, and a minimal smoke-test theme. |
| v3.4.7 | ornate timber dashboard | Add a master-carpenter style dashboard using multiple timber treatments, timber needles, timber ticks, and carved/inlaid visual language. |
| v3.4.8 | neon-grid dashboard | Add a dark retro-tech dashboard with Tron-like neon blue glow, grid/circuit accents, and luminous gauge assets. |
| v3.4.9 | steam-scrap dashboard | Add a steampunk/scrapyard dashboard with brass/copper/iron plates, extra pipes, wires, rivets, lamps, and intentionally overbuilt decoration. |

Example dashboard rules:

- Runtime gauge config still uses behaviour types only: `numeric`, `radial`, `odometer`, `indicator`, `bar`, and `segmented`.
- Visual identity belongs in generated image assets and dashboard/package layout.
- Generator-internal theme settings are allowed, but they must not create runtime `style` fields.
- Source asset dimensions are authoritative; if a display needs to fit a different window size, use dashboard/widget config `scale`.
- Generated digit sets may choose their own cell dimensions, but slot-positioned assets within a set must stay internally consistent.
- Decimal points are overlays on the current or preceding digit cell and do not consume separate slots.
- Do not infer one digit set's dimensions from another digit set.
- Decorative elements such as timber panels, pipes, grid lines, screw heads, or wires are background/overlay assets, not new renderer behaviour.
- Generated PNGs should be reproducible from scripts with stable seed/config values.
- Each themed dashboard should cover as many completed gauge types as practical.

## Dashboard CLI tail

v3.4.10 through v3.4.12 reshape the existing flat flag dashboard entry points into dashboard-scoped commands for the active dashboard tooling.

The completed command tree target is:

```text
GoDriveLog dashboard [--config <config-file>]
GoDriveLog dashboard run [vehicle-id] [--config <config-file>] [--renderer ebiten]
GoDriveLog dashboard harness [vehicle-id] [--config <config-file>] [--pattern sweep] [--interval 50ms] [--duration 60s] [--renderer ebiten]
GoDriveLog dashboard examples --output <directory> [--config <config-file>] [--vehicle <vehicle-id>] [--theme framework-smoke] [--force]
GoDriveLog dashboard validate [config-file]
GoDriveLog dashboard validate [--config <config-file>]
```

The CLI tail routing work must not create replacement runtime, renderer, harness, config, validation, or example-generation systems. The existing `cmd/GoDriveLog/main_ebiten.go` flat flag behaviours already reach the required backend paths; the CLI work should move those switches into named command drawers.

Command routing target:

| New command form | Existing flat flag path or machinery being remapped | Existing code path to reuse |
|---|---|---|
| `GoDriveLog dashboard run [vehicle-id]` | Default run path when `--harness=false`; uses existing `--config`, `--vehicle`, `--renderer`, and current runtime duration handling if preserved. | Existing `runV3EbitenCommand(configPath, vehicleID, duration)` path. |
| `GoDriveLog dashboard harness [vehicle-id]` | Existing `--harness=true` path plus `--config`, `--vehicle`, `--pattern`, `--interval`, `--duration`, and `--renderer`. | Existing `runV3EbitenHarnessCommand(configPath, vehicleID, pattern, interval, duration)` path. |
| `GoDriveLog dashboard validate [config-file]` | Existing config load and validation behaviour, reached through a command instead of flat flags. | Existing config parsing/validation helpers; do not create a replacement validator. |
| `GoDriveLog dashboard validate --config <config-file>` | Existing `--config` file selection plus existing config validation behaviour. | Existing config parsing/validation helpers; do not create a replacement validator. |
| `GoDriveLog dashboard [--config <config-file>]` | Existing config load structures, rendered as a compact overview. | Existing config parsing structures; do not invent a new config model. |
| `GoDriveLog dashboard examples --output <directory>` | Existing generated example asset machinery/scripts, plus existing `--config` and `--vehicle` concepts where relevant. | Existing generated-example helpers/scripts; do not build a duplicate generator. |

Flag redistribution:

| Existing flat flag | New command usage |
|---|---|
| `--config` | Used by `dashboard`, `dashboard run`, `dashboard harness`, `dashboard examples`, and `dashboard validate`. |
| `--vehicle` | Becomes positional `[vehicle-id]` for `dashboard run` and `dashboard harness`; remains `--vehicle` for `dashboard examples`. |
| `--renderer` | Used by `dashboard run` and `dashboard harness`. |
| `--duration` | Used by `dashboard harness`; may stay on `dashboard run` if preserving the existing runtime duration behaviour. |
| `--harness` | Replaced by the command name `dashboard harness`. |
| `--pattern` | Used by `dashboard harness`. |
| `--interval` | Used by `dashboard harness`. |
| `--v3` | Removed; the `dashboard` command tree implies the active v3 dashboard path. |

The CLI tail is split into focused slices:

| Version | Slice | Result |
|---|---|---|
| v3.4.10 | dashboard CLI command tree | Remap `run`, `harness`, `examples`, and `validate`; add deterministic config discovery and help-output coverage for implemented commands. |
| v3.4.11 | dashboard overview | Add bare `dashboard` compact config overview using existing config structures. |
| v3.4.12 | gauge-aware harness sweep | Make `dashboard harness --pattern sweep` drive each gauge according to gauge behaviour. |

### Config discovery rules

When a positional config file or `--config` is supplied, load exactly that file and bypass config discovery.

When no config file is supplied, search config files from:

```text
current working directory
/etc/godrivelog recursively
```

Directory traversal must be deterministic:

- sort each directory's entries alphabetically before evaluating them;
- search the current working directory non-recursively;
- search `/etc/godrivelog` recursively;
- evaluate only candidate config filenames.

Candidate config filenames are:

```text
godrivelog.yaml
godrivelog.yml
dashboard.yaml
dashboard.yml
config.yaml
config.yml
```

If no vehicle ID is supplied, use the first valid config in search order that defines exactly one vehicle. If the first valid config defines multiple vehicles, stop and return an error requiring a vehicle ID.

If a vehicle ID is supplied, use the first valid single-vehicle or multi-vehicle config in search order that defines the matching vehicle. Keep searching until a match is found. If no config contains the requested vehicle, return an error that lists searched config files and the vehicles found in each valid config.

Do not default to `config.example.yaml`.

### Dashboard command rules

- Keep active dashboard CLI functions under `dashboard`.
- Do not add `dashboard preview` yet.
- `--renderer` is optional and defaults to `ebiten`.
- `ebiten` remains the only active renderer.
- Bare `dashboard` prints a compact overview of the resolved config in v3.4.11.
- The overview prints the configured vehicle OBD source string as-is; it does not infer source-type labels.
- The overview must list gauges/widgets with gauge ID/name, type, source of data, and PID where applicable.
- Do not add a separate PID section.
- `dashboard examples` requires `--output` in v3.4.10.
- `dashboard examples --output <directory>` treats `<directory>` as the generated dashboard root, creates it when missing, writes `dashboard.yaml` and `assets/` directly inside it, and does not add a theme subdirectory.
- Generated example output must use paths relative to the output directory so each generated dashboard is self-contained and movable.
- If the examples output directory exists and is non-empty, interactive terminals may prompt unless `--force` is supplied; non-interactive use must fail unless `--force` is supplied.
- `dashboard harness --pattern sweep` becomes gauge-aware in v3.4.12 while keeping other patterns for later scope decisions.
- Do not add compatibility or migration behaviour for any earlier flat flag shape.

Gauge-aware sweep rules:

| Gauge type | Sweep behaviour |
|---|---|
| `odometer` | Let `n = 0`; increment from `n - 20` to `n + 20` for 5 seconds, then `n + 20` to `n + 30` for 5 seconds. |
| `numeric` | Same as `odometer`. |
| `radial` | Keep the existing sweep style. |
| `indicator` | Flash on/off for 5 seconds with a 1s duty cycle, then flash on/off for 5 seconds with a 250ms duty cycle. |
| `bar` | Heartbeat pulse at 90 bpm. |
| `segmented` | Same input shape as `radial`. |

Do not change gauge package semantics, renderer scene semantics, or generated example artwork rules in the CLI tail.

## Planned implementation slices

| Version | Slice | Result |
|---|---|---|
| v3.4.0 | gauge type cleanup docs and naming | Create v3.4 docs, rename plan, and prompt set. |
| v3.4.1 | numeric rename | Replace `seven_segment` with `numeric` in code/examples, no compatibility alias. |
| v3.4.2 | odometer planning/scene model | Add odometer config and flat strip scene model with `smooth` and `click` movement modes. |
| v3.4.3 | indicator gauge | Add off/on state rendering. |
| v3.4.4 | bar gauge | Add first continuous transform behaviour for level/reveal. |
| v3.4.5 | segmented gauge | Add `{percent}` discovery, threshold-gap hysteresis, and percent-threshold image selection. |
| v3.4.6 | example asset generation framework | Add deterministic procedural asset generation structure and smoke-test output. |
| v3.4.7 | ornate timber dashboard | Add generated ornate timber dashboard assets/config. |
| v3.4.8 | neon-grid dashboard | Add generated Tron-like dark neon dashboard assets/config. |
| v3.4.9 | steam-scrap dashboard | Add generated steampunk/scrapyard dashboard assets/config. |
| v3.4.10 | dashboard CLI command tree | Remap dashboard-scoped `run`, `harness`, `examples`, and `validate`, add deterministic config discovery, and help-output coverage. |
| v3.4.11 | dashboard overview | Add compact config overview for bare `dashboard`. |
| v3.4.12 | gauge-aware harness sweep | Refine `dashboard harness --pattern sweep` so synthetic input matches gauge behaviour. |

## Branch-chat workflow

Each implementation chat should:

1. Read this file.
2. Read `docs/v3.4/ImplementationState.md`.
3. Confirm the previous relevant PR is merged into `main`.
4. Confirm there are no blocking open PRs.
5. Create a branch from latest `main` using the full target version prefix.
6. Implement only that version slice.
7. Update `CHANGES.md` and `docs/v3.4/ImplementationState.md`.
8. Open a PR.
9. Stop.

## Things not to do

- Do not add a `style` field.
- Do not keep `seven_segment` as a compatibility alias.
- Do not make dot-matrix a font/text renderer in this line.
- Do not merge `bar` and `segmented` just because both can look like level displays.
- Do not preload all percent images for `segmented`.
- Do not chase curved odometer wheel depth before flat strip scrolling works.
- Do not add unrelated platform/package work here.
- Do not use image-generation services, downloaded stock art, or non-reproducible manual PNG editing for the v3.4 example dashboard tail.
- Do not add `dashboard preview` until it has a job distinct from harness/runtime.