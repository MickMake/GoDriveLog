# GoDriveLog v3.3 release plan

Status: final housekeeping slice in progress
Owner: migration implementor

## Purpose

v3.3 planned, tested, and completed the renderer migration decision for the active v3 dashboard path.

The result is clear: Ebiten is the default and only active v3.3 dashboard renderer implementation. Fyne support remains available only from the v3.2.x line.

## Release goal

Move the active v3 dashboard command path to Ebiten while preserving the existing runtime, sensor, dashboard scene, and display-sink boundaries.

The baseline workload is:

- 3 digit seven-segment temperature display using `coolant_temperature`, range `-10..40`, including minus-symbol rendering.
- 3 digit seven-segment speed display using `speed`.
- 4 digit seven-segment RPM display using `rpm`.
- Radial RPM gauge using the same `rpm` sensor as the numeric RPM display.

The canonical baseline config is:

```text
examples/baseline-dashboard.yaml
```

## Final release principles

- Ebiten is the active v3.3 renderer implementation.
- The renderer remains a boundary; runtime, sensors, logs, dashboard scenes, and display sinks stay renderer-independent.
- Fyne support ends with v3.2.x.
- v3.2.9 is superseded by the v3.3 renderer decision; do not create a separate v3.2.9 branch.
- Reuse the existing v3 dashboard scene model.
- Do not redesign gauge packages.
- Do not add widget-level sensor overrides.
- Do not add inheritance or clusters.
- Do not use renderer-local fake data for renderer validation.
- Test through the full path: OBD or harness source, prepared vehicle/sensor data, runtime events, dashboard scene generation, display sink, renderer adapter, screen.
- Prefer Raspberry Pi measurements over desktop impressions.
- Keep v3.4/v4.0 planning out of the v3.3.4 housekeeping slice.

## Branch-chat workflow

Each implementation chat should:

1. Read this file.
2. Read `docs/v3.3/ImplementationState.md`.
3. Confirm the previous relevant PR is merged into `main`.
4. Confirm there are no blocking open PRs.
5. Create a branch from latest `main` using the full target version prefix where the connector allows it.
6. Implement only that version slice.
7. Update `CHANGES.md` and `docs/v3.3/ImplementationState.md`.
8. Open a PR.
9. Stop.

## Completed implementation slices

| Version | Slice | Result |
|---|---|---|
| v3.3.0 | renderer checkpoint planning | Created v3.3 docs, prompts, and reusable examples structure. |
| v3.3.1 | Ebiten renderer spike | Added an Ebiten backend using the real dashboard path. |
| v3.3.2 | renderer decision | Promoted Ebiten to the primary v3.3 renderer and retired Fyne from the active runtime path. |
| v3.3.3 | dependency cleanup | Removed legacy Fyne code packages and Fyne module dependencies from the active branch. |
| v3.3.4 | post-Fyne renderer housekeeping | Audits Fyne removal, preserves renderer-boundary wording, records Pi performance evidence, documents example geometry, and closes v3.2.9 as superseded. |

## v3.3.0 checkpoints

- `docs/v3.3/` contains the planning documents.
- `docs/v3.3/prompts/` contains one prompt per planned decision slice.
- Reusable baseline config is moved out of versioned docs into `examples/baseline-dashboard.yaml`.
- Versioned docs reference the examples path instead of carrying active runnable config copies.
- Existing v3.2 baseline behaviour is preserved conceptually.
- No renderer implementation code is changed.

## v3.3.1 checkpoints

- Add Ebiten support using the same v3 runtime/harness/dashboard scene path.
- Do not use a demo-only renderer loop or renderer-private fake values.
- Support the baseline dashboard enough for comparison: static layers, seven-segment digits, minus glyph, RPM numeric, and radial RPM.
- Start with Ebiten runtime needle rotation.
- Document that prepared radial needle frames may be needed if runtime rotation is too costly on the Pi.
- Preserve comparable display stats.

## v3.3.2 checkpoints

- Promote Ebiten when real workload results justify it.
- Keep Fyne support in the v3.2.x line only.
- Keep the dashboard scene model intact.
- Keep mobile/platform packaging as a later follow-up, not part of v3.3.

## v3.3.3 checkpoints

- Remove legacy Fyne adapter, renderer, and UI code from the active branch.
- Remove Fyne module dependencies from `go.mod` and `go.sum`.
- Simplify active renderer selection to Ebiten.
- Update README, `CHANGES.md`, and v3.3 implementation state.

## v3.3.4 checkpoints

- Audit active code and module files for Fyne removal.
- Keep Ebiten as the default and only active v3.3 renderer implementation.
- Preserve the renderer boundary in documentation and wording.
- Record Raspberry Pi performance evidence in `docs/v3.3/PerformanceRuns.md`.
- Update baseline verification docs with the Fyne usability note and Ebiten run numbers.
- Document that example seven-segment digit positions are source-artwork alignment coordinates, not simple logical-width coordinates.
- Mark v3.2.9 as superseded by the v3.3 renderer decision.

## Things not to do

- Do not redesign gauge packages.
- Do not change sensor ownership rules.
- Do not add widget-level sensor overrides.
- Do not add inheritance.
- Do not add clusters.
- Do not add procedural gauge artwork.
- Do not rebuild the dashboard model around Ebiten.
- Do not add v3.4 feature work in v3.3.4.
- Do not start v4.0 product/platform hardening here.
