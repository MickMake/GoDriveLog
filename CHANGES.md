# CHANGES

## Unreleased

- Added v3 inverse implementation audit documentation for old/current behaviours not yet fully rebuilt as v3.
- Updated v3 migration state for the v3.0.12 inverse implementation audit slice.
- Added v3 retirement audit documentation for old/current paths that may be reviewed for later removal or archiving.
- Added v3 richer dashboard widget rendering for `bar_display` and `frame_gauge`.
- Added dashboard tests for bar fill mapping, reverse fill direction, zones, frame clamping, sensor status handling, and unchanged frame output handling.
- Added v3 richer asset registry support for bar and frame asset families.
- Added reusable decoded bar cell and frame asset structs for later dashboard widgets.
- Added richer asset registry tests for bar cells, frame ranges, decoded frame assets, and related image dimension validation.
- Added v3 smallest selected-dashboard scene runtime for image, digit display, and indicator widgets.
- Added selected-dashboard scene tests for RuntimePlan dashboard selection, digit formatting, decimal point overlays, indicator status mapping, and unchanged formatted output handling.
- Added v3 minimal asset registry for image, digit, and indicator asset families.
- Added reusable decoded image asset structs so future widgets can avoid hot-path asset loading.
- Added tests for repository-root asset path resolution, missing asset errors, decoded digit assets, and required indicator states.
- Updated v3 migration state for the v3.0.11 retirement audit slice.
- Updated v3 migration state for the v3.0.10 implementation slice.
- Updated v3 migration state for the v3.0.9 implementation slice.
- Updated v3 migration state for the v3.0.8 implementation slice.
- Updated v3 migration state for the v3.0.7 implementation slice.

## 0.1 - 2026-06-08

- Created PR-tail package for GoDriveLog dashboard v2.7 throttle fixture completion.
- Restored throttle frame count from 3 to 11 in `config.example.yaml`.
- Added placeholder SVG throttle frames 003 through 010 for 30% through 100%.
- Kept changes limited to the v2.7 example/dashboard fixture assets.
