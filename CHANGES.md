# CHANGES

## Unreleased

- Added v3 minimal asset registry for image, digit, and indicator asset families.
- Added reusable decoded image asset structs so future widgets can avoid hot-path asset loading.
- Added tests for repository-root asset path resolution, missing asset errors, decoded digit assets, and required indicator states.
- Updated v3 migration state for the v3.0.7 implementation slice.

## 0.1 - 2026-06-08

- Created PR-tail package for GoDriveLog dashboard v2.7 throttle fixture completion.
- Restored throttle frame count from 3 to 11 in `config.example.yaml`.
- Added placeholder SVG throttle frames 003 through 010 for 30% through 100%.
- Kept changes limited to the v2.7 example/dashboard fixture assets.
