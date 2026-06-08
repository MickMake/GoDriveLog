# CHANGES

## 2.9.0 - 2026-06-08

- Added reusable dashboard block type names for `seven_segment_number`, `percent_frame_bar`, `state_lamp`, `glowing_number_box`, `labelled_sensor_value`, `warning_overlay`, `stale_overlay`, and `static_panel`.
- Resolved reusable block names through existing scene primitives so the renderer continues to draw only image, sprite-frame, sprite-text, and group elements.
- Added scene regression tests proving reusable blocks are sensor-agnostic aliases over configured assets and decoder inputs.
- Updated `config.example.yaml` to use reusable block names for the first real dashboard slice.

## 2.8.0 - 2026-06-08

- Removed the old standalone dashboard widget package tree.
- Kept application startup on the dashboard v2 scene path.
- Confirmed config examples and README describe the v2 scene dashboard shape.
- Updated dashboard v2 status, overview, and decision notes for the old-widget removal stage.

## v2.7.1 PR22 stale-review fix

- Added `StaleAfter` to `sensors.SensorDefinition` and `sensors.SensorState`.
- Derived `StaleAfter` from configured sensor refresh interval using a 2x refresh grace window.
- Updated `StateStore` stale evaluation to use each sensor state's own threshold.
- Removed dashboard-level global stale timing.
- Added a regression test proving fast sensors can become stale while slow sensors remain healthy.

## 2.7.0 - 2026-06-08

- Implemented dashboard v2.7.x first real dashboard vertical slice in `config.example.yaml`.
- Added small local SVG fixture assets for a static background, yellow RPM digits, throttle frame bar, redline glow overlay, and status badges.
- Added YAML-backed dashboard block conditions for sensor/decoder status and value checks.
- Added configured scene-condition tests for decoder-driven redline overlays and sensor-status indicators.
- Preserved sensor status/error metadata through decoder outputs so dashboard status elements can respond to stale/error state.
- Removed the standalone old `widget` command path from app startup.
- Retired the unsafe `ui.NewDashboard` constructor so normal dashboard construction keeps the config path for asset resolution.

## 2.6.0 - 2026-06-08

- Implemented dashboard v2.6.x Fyne scene renderer under `internal/dashboard/renderer/fyne`.
- Added rendering support for image, sprite_frame, sprite_text, and group scene elements.
- Wired the app dashboard through asset loading, decoder execution, scene evaluation, and renderer updates.
- Added periodic dashboard refresh from `StateStore` snapshots so scene output follows live sensor state.
- Added renderer tests for visible element rendering, sprite text glyph layout, and group rendering.
- Kept decoder logic outside the renderer and kept polling/logging ownership unchanged.

## 2.5.0 - 2026-06-08

- Implemented dashboard v2.5.x scene primitive evaluation under `internal/dashboard/scene`.
- Added renderer-independent scene elements for image, sprite_frame, sprite_text, and group blocks.
