# GoDriveLog v3.4 release plan

Status: planning docs in progress
Owner: gauge package implementor

## Purpose

v3.4 defines the next gauge package model for the active Ebiten dashboard path.

The release direction is gauge/display package cleanup and expansion. It does not introduce a platform packaging track.

## Release goal

Create a clear gauge type model where type names describe renderer behaviour and visual identity comes from assets.

The planned gauge types are:

| Type | Purpose |
|---|---|
| `numeric` | Formatted value rendered through image assets per character slot. Replaces `seven_segment`. |
| `radial` | Needle/arc gauge using value-to-angle mapping. Already exists. |
| `odometer` | Rolling wheel gauge with mechanical-style digit movement. |
| `indicator` | Two-state gauge: off/on. |
| `bar` | Continuous value display: fill, reveal, or pointer movement. |
| `segmented` | Discrete stepped value display using pre-rendered percent-threshold images. |

## Final release principles

- Gauge type describes behaviour, not visual style.
- Visual style is asset-only.
- Do not add a `style` field.
- Rename `seven_segment` to `numeric` with no compatibility alias.
- Keep `radial` as the existing radial gauge type.
- Keep `bar` and `segmented` separate.
- Keep each implementation slice small.
- Do not add all future gauge families in one monster PR. That beast would need feeding and possibly a small helmet.

## Gauge type boundaries

| Type | Behaviour |
|---|---|
| `numeric` | Format value, split into character slots, draw matching character assets. |
| `radial` | Normalize value and rotate a needle around a pivot. |
| `odometer` | Convert value into rolling wheel strip positions. |
| `indicator` | Select one of two state layers. |
| `bar` | Continuously reveal, clip, fill, or move an active layer. |
| `segmented` | Select the highest discovered percent-threshold image reached by the current value. |

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

The renderer discovers matching files, extracts percent thresholds, sorts them, and draws the highest threshold less than or equal to the current normalized percent.

Discovery must count filenames only. It must not decode every image. Runtime should lazy-load selected images and cache recent images.

## Planned implementation slices

| Version | Slice | Result |
|---|---|---|
| v3.4.0 | gauge type cleanup docs and naming | Create v3.4 docs, rename plan, and prompt set. |
| v3.4.1 | numeric rename | Replace `seven_segment` with `numeric` in code/examples, no compatibility alias. |
| v3.4.2 | odometer planning/scene model | Add odometer config and scene model without overbuilding wheel depth. |
| v3.4.3 | indicator gauge | Add off/on state rendering. |
| v3.4.4 | bar gauge | Add first continuous level/reveal behaviour. |
| v3.4.5 | segmented gauge | Add `{percent}` discovery and percent-threshold image selection. |

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
