# Dashboard v2 Current Status

## Active stage
v2.7.x - First real dashboard

## Completed
- v2.0.x - New config schema
- v2.1.x - Dashboard config validation
- v2.2.x - Sensor state boundary
- v2.3.x - Decoder engine
- v2.4.x - Asset registry
- v2.5.x - Scene primitives
- v2.6.x - Fyne scene renderer

## Current branch
- feature/dashboard-v2-real-dashboard

## Decisions
- No legacy dashboard compatibility.
- Sensors are separate from dashboard visuals.
- Dashboard asset registry loads and caches local assets only.
- Fyne rendering consumes resolved scene elements and does not perform decoder logic directly.
- Dashboard UI refreshes from `StateStore` snapshots and re-evaluates scene state for rendering.
- First real dashboard uses configured assets, decoders, scene blocks, and YAML-backed scene conditions.
- Decoder outputs preserve source sensor status/error metadata so configured status indicators can render.

## Next prompt
Use v2.7.x prompt from prompts.md until merged. After merge, use v2.8.x prompt.

## Known risks
- Example dashboard assets are intentionally small SVG placeholders, not final artwork.
- Sprite text currently distributes glyphs evenly across the configured text geometry.
- Full visual verification remains manual in mock mode.
- Broader old-widget package deletion remains a v2.8 cleanup task.
