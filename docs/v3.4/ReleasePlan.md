# GoDriveLog v3.4 release plan

Status: example dashboard tail planned
Owner: gauge package implementor

## Purpose

v3.4 defines the next gauge package model for the active Ebiten dashboard path.

The release direction is gauge/display package cleanup and expansion. It does not introduce a platform packaging track.

The original behaviour implementation slices are complete through v3.4.5. The remaining v3.4 tail adds generated example dashboards that prove the completed gauge types with repeatable assets.

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
