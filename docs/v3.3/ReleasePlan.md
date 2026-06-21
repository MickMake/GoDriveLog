# GoDriveLog v3.3 release plan

Status: planning
Owner: migration implementor

## Purpose

v3.3 plans and measures an Ebiten renderer spike beside the existing Fyne renderer.

The goal is not to rewrite the dashboard. The goal is to test whether Ebiten is clearly faster, smoother, and simpler for the real GoDriveLog dashboard path on Raspberry Pi hardware.

## Release goal

Add a renderer comparison path that keeps Fyne as the default while allowing an Ebiten backend to be tested against the same baseline dashboard workload and the same upstream runtime path.

The baseline workload is:

- 3 digit seven-segment temperature display using `coolant_temperature`, range `-10..40`, including minus-symbol rendering.
- 3 digit seven-segment speed display using `speed`.
- 4 digit seven-segment RPM display using `rpm`.
- Radial RPM gauge using the same `rpm` sensor as the numeric RPM display.

The canonical baseline config is:

```text
examples/baseline-dashboard.yaml
```

## Release principles

- Keep Fyne as the default until Ebiten clearly wins.
- Add `--renderer fyne|ebiten`; do not remove Fyne.
- Reuse the existing v3 dashboard scene model.
- Do not redesign gauge packages.
- Do not add widget-level sensor overrides.
- Do not add inheritance or clusters.
- Do not use renderer-local fake data for the spike.
- Test through the full path: OBD or harness source, prepared vehicle/sensor data, runtime events, dashboard scene generation, display sink, selected renderer, screen.
- Prefer Raspberry Pi measurements over desktop impressions.
- Decide with numbers, not vibes in a hi-vis vest.

## Branch-chat workflow

Each implementation chat should:

1. Read this file.
2. Read `docs/v3.3/prompts/README.md`.
3. Read `docs/v3.3/ImplementationState.md`.
4. Read the prompt file for the target slice under `docs/v3.3/prompts/`.
5. Confirm the previous relevant PR is merged into `main`.
6. Confirm there are no blocking open PRs.
7. Create a branch from latest `main` using the full target version prefix.
8. Implement only that version slice.
9. Update `CHANGES.md` and `docs/v3.3/ImplementationState.md`.
10. Open a PR.
11. Stop.

Do not redesign the release plan inside a slice chat.

## Planned implementation slices

| Version | Slice | Goal |
|---|---|---|
| v3.3.0 | renderer checkpoint planning | Create v3.3 docs, prompts, and reusable examples structure. |
| v3.3.1 | Ebiten renderer spike | Add an experimental Ebiten backend beside Fyne using the real dashboard path. |
| v3.3.2 | renderer A/B comparison | Run Fyne and Ebiten against the same baseline with fixed duration and comparable stats. |
| v3.3.3 | renderer decision | Decide whether to continue, promote, pause, or abandon Ebiten. |
| v3.3.4 | act on decision | Only needed if v3.3.3 identifies a clear follow-up implementation path. |

## v3.3.0 checkpoints

- `docs/v3.3/` contains the planning documents.
- `docs/v3.3/prompts/` contains one prompt per planned decision slice.
- Reusable baseline config is moved out of versioned docs into `examples/baseline-dashboard.yaml`.
- Versioned docs reference the examples path instead of carrying active runnable config copies.
- Existing v3.2 baseline behaviour is preserved conceptually.
- No renderer implementation code is changed.

## v3.3.1 checkpoints

- Add `--renderer fyne|ebiten`.
- Keep Fyne as the default.
- Add `internal/dashboard/adapter/ebiten` or equivalent narrow backend package.
- Reuse `[]v3dashboard.Scene`.
- Do not use a demo-only renderer loop or renderer-private fake values.
- Support the baseline dashboard enough for comparison: static layers, seven-segment digits, minus glyph, RPM numeric, and radial RPM.
- Start with Ebiten runtime needle rotation.
- Document that prepared radial needle frames may be needed if runtime rotation is too costly on the Pi.
- Preserve comparable display stats.

## v3.3.2 checkpoints

- Add or use a fixed run duration such as `--duration 60s` for comparable renderer runs.
- Run the same config, vehicle, pattern, interval, and duration through Fyne and Ebiten.
- Capture `events`, `display_submitted`, `display_rendered`, `display_superseded`, and render duration stats.
- Capture CPU use and subjective smoothness on Raspberry Pi.
- Record results in `docs/v3.3/BaselineDashboardVerification.md` or a linked results section.

## v3.3.3 checkpoints

- Record a clear decision: continue, promote, pause, or abandon Ebiten.
- Keep Fyne primary unless Ebiten clearly wins on Raspberry Pi.
- Treat a 10-20% Ebiten win as probably not enough.
- Treat a 2x+ win plus smoother behaviour and manageable code complexity as a strong continue signal.
- Identify the next slice only after the decision.

## Things not to do

- Do not remove Fyne.
- Do not redesign gauge packages.
- Do not change sensor ownership rules.
- Do not add widget-level sensor overrides.
- Do not add inheritance.
- Do not add clusters.
- Do not add procedural gauge artwork.
- Do not rebuild the dashboard model around Ebiten.
- Do not chase perfect visual polish before measuring performance.
