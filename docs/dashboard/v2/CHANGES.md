# CHANGES

## 0.6 - 2026-06-08

- Implemented dashboard v2.5.x scene primitive evaluation under `internal/dashboard/scene`.
- Added renderer-independent scene elements for image, sprite_frame, sprite_text, and group blocks.
- Added z-order layer sorting for evaluated scene elements.
- Added runtime condition evaluation using supplied sensor or decoder values.
- Added sprite_frame resolution from decoder frame indexes to frame_set assets.
- Added sprite_text resolution from decoder text/digits to charset glyphs.
- Added tests for z-order, conditions, frame resolution, text/glyph resolution, groups, and missing glyph errors.
- Kept this stage independent of Fyne rendering, gauges, and legacy dashboard widgets.

## 0.5 - 2026-06-08

- Implemented dashboard v2.4.x asset registry under `internal/dashboard/assets`.
- Added cached loading for image, frame_set, and charset assets.
- Added dashboard `asset_root` support so relative asset paths can resolve from the config file directory.
- Added generated frame-set support with `{index}` and zero-padded `{index:03}` patterns, while keeping explicit frame lists supported.
- Added clear load errors for missing image, frame, and glyph files.
- Added tests using small fixture assets for load, cache, missing asset, generated frames, explicit frames, charset glyphs, and remote-path rejection.
- Kept this stage independent of scene primitives and rendering.

## 0.4 - 2026-06-08

- Implemented dashboard v2.3.x decoder engine under `internal/dashboard/decoders`.
- Added ordered decoder execution with support for sensor inputs and earlier decoder outputs.
- Added decoder implementations for normalize, threshold, frame_index, format_number, digits, and boolean.
- Added decoder tests covering all decoder types and common error cases.
- Tightened dashboard decoder config validation so each decoder must define exactly one input source and decoder-to-decoder references must point backwards to an earlier decoder.
- Kept decoder execution independent of Fyne, assets, scene primitives, and rendering.

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
