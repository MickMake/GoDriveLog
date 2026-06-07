# CHANGES

## 0.3 - 2026-06-08

- Implemented dashboard v2.2.x sensor state boundary.
- Added a neutral sensor state store with current value, unit, min, max, status, error text, and update time.
- Added store tests for initial definitions, value updates, error updates, stale status, and sorted snapshots.
- Wired active sensor polling to write latest success/error state while preserving existing JSONL logging behaviour.
- Passed the state store into the dashboard placeholder so future dashboard stages can consume state instead of sensor config.

## 0.2 - 2026-06-08

- Implemented dashboard v2.1.x config validation schema for assets, decoders, blocks, and layers.
- Added validation for duplicate and missing IDs across dashboard config sections.
- Added validation for unsupported dashboard asset, decoder, and block types.
- Added validation for sensor, asset, decoder, block, and layer references.
- Added geometry validation for non-group dashboard blocks.
- Updated `config.example.yaml` and `docs/config.md` to show the v2.1.x schema.

## 0.1 - 2026-06-07

- Created dashboard implementation planning document set for the GoDriveLog v2.x.x dashboard rewrite.
- Added implementation overview mapping stages 1-10 to v2.0.x through v2.9.x.
- Added condensed implementation prompt series with guardrails per stage.
- Added reference/checklist file for schema, decoders, assets, scene primitives, and validation.
- Packaged documents into a zip archive with a top-level directory.
