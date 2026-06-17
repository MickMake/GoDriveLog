# CHANGES

## Unreleased

- Added `v3.1.2` dashboard harness behind `--v3 --harness`, feeding fake sensor events through the real v3 dashboard scene path and Fyne display adapter.
- Removed the temporary `--repo-root` flag and added ordered relative asset search paths: config directory + vehicle ID, current working directory + vehicle ID, config directory, then current working directory.
- Added explicit harness patterns: `sweep`, `heartbeat`, and `fixed`, including tuned sweep and heartbeat timing tests.
- Documented dashboard harness command usage, asset search path order, pattern semantics, and 50ms/100ms cadence options.
- Added `v3.1.1` Fyne display adapter for v3 dashboard scene output, keeping display code below the dashboard runtime boundary.
- Wired the `--v3` command path to show selected v3 dashboard scenes in a Fyne window while retaining the existing runtime as the default path.
- Added adapter tests covering repo-relative asset rendering and rejection of escaping asset paths.
- Added `v3.1.0` runnable v3 command path behind `--v3`, including selected vehicle resolution, endpoint connection, sensor polling runtime startup, selected JSONL event subscribers, dashboard scene boundary logging, and signal-based clean shutdown.
- Added `internal/runtime/v3runtime` orchestration tests covering v3 config load, selected vehicle connection, sensor polling, selected JSONL event output, and reader cleanup.
- Split v3.1 implementation prompts into per-slice files under `docs/v3.1/prompts/`.
- Reframed PR 51 as docs-only planning setup, with implementation starting at `v3.1.0`.
- Added blocking and impact metadata to v3.1 open decisions.
- Clarified v3.1 slice docs-update rules for `MigrationState.md` and `OpenDecisions.md`.
- Revised the v3.1 release plan around the runnable command path, display adapter, dashboard/gauge test harness, and dashboard update cadence targets.
- Expanded v3.1 carry-forward and release plan docs with v3.0 implementation details and per-slice checkpoints.
- Added v3.1 release planning stubs under `docs/v3.1/` for the next implementation phase.
- Added v3 inverse implementation audit documentation for old/current behaviours not yet fully rebuilt as v3.
- Updated v3 migration state for the v3.0.12 inverse implementation audit slice.
- Added v3 retirement audit documentation for old/current paths that may be reviewed for later removal or archiving.
- Added v3 richer dashboard widget rendering for `bar_display` and `frame_gauge`.
- Added dashboard tests for bar fill mapping, reverse fill direction, zones, frame clamping, sensor status handling, and unchanged frame output handling.
- Added v3 richer asset registry support for bar and frame asset families.
- Added reusable decoded bar cell and frame asset structs for later dashboard widgets.
- Added v3 smallest selected-dashboard scene runtime for image, digit display, and indicator widgets.
- Added selected-dashboard scene tests for RuntimePlan dashboard selection, digit formatting, decimal point overlays, indicator status mapping, and unchanged formatted output handling.
- Added v3 minimal asset registry for image, digit, and indicator asset families.
- Added reusable decoded image asset structs so future widgets can avoid hot-path asset loading.
- Added tests for repository-root asset path resolution, missing asset errors, decoded digit assets, and required indicator states.
- Updated v3 migration state for the v3.0.11 retirement audit slice.
- Updated v3 migration state for the v3.0.10 implementation slice.
- Updated v3 migration state for the v3.0.9 implementation slice.
- Updated v3 migration state for the v3.0.7 implementation slice.

## 0.1 - 2026-06-08

- Created PR-tail package for GoDriveLog dashboard v2.7 throttle fixture completion.
- Restored throttle frame count from 3 to 11 in `config.example.yaml`.
- Added placeholder SVG throttle frames 003 through 010 for 30% through 100%.
- Kept changes limited to the v2.7 example/dashboard fixture assets.
