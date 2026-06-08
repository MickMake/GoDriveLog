# Dashboard v2 Current Status

## Active stage
v2.8.x - Remove old widgets

## Completed
- v2.0.x - New config schema
- v2.1.x - Dashboard config validation
- v2.2.x - Sensor state boundary
- v2.3.x - Decoder engine
- v2.4.x - Asset registry
- v2.5.x - Scene primitives
- v2.6.x - Fyne scene renderer
- v2.7.x - First real dashboard

## Current branch
- feature/v2-8-remove-widgets

## Decisions
- No legacy dashboard compatibility.
- Sensors are separate from dashboard visuals.
- Dashboard asset registry loads and caches local assets only.
- Fyne rendering consumes resolved scene elements and does not perform decoder logic directly.
- Dashboard UI refreshes from `StateStore` snapshots and re-evaluates scene state for rendering.
- First real dashboard uses configured assets, decoders, scene blocks, and YAML-backed scene conditions.
- Decoder outputs preserve source sensor status/error metadata so configured status indicators can render.
- Sensor stale status is owned by `StateStore` and derived per sensor from configured refresh intervals.
- Dashboard scenes consume already-classified sensor status; dashboard rendering must not apply global stale timing.
- The throttle example uses an 11-frame fixture set for 10% visual steps.
- The old standalone widget package tree has been removed from the v2 application path.

## Next prompt
Use v2.8.x prompt from prompts.md until merged. After merge, use v2.9.x prompt.

## Known risks
- Example dashboard assets are intentionally small SVG placeholders, not final artwork.
- Sprite text currently distributes glyphs evenly across the configured text geometry.
- Full visual verification remains manual in mock mode.
- Broader reusable dashboard blocks remain a v2.9 cleanup/buildout task.
