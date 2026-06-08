# Dashboard v2 Current Status

## Active stage
v2.6.x - Fyne scene renderer

## Completed
- v2.0.x - New config schema
- v2.1.x - Dashboard config validation
- v2.2.x - Sensor state boundary
- v2.3.x - Decoder engine
- v2.4.x - Asset registry
- v2.5.x - Scene primitives

## Current branch
- feature/dashboard-v2-fyne-scene-renderer

## Decisions
- No legacy dashboard compatibility.
- Sensors are separate from dashboard visuals.
- Dashboard asset registry loads and caches local assets only.
- Fyne rendering consumes resolved scene elements and does not perform decoder logic directly.
- Dashboard UI refreshes from `StateStore` snapshots and re-evaluates scene state for rendering.

## Next prompt
Use v2.6.x prompt from prompts.md until merged. After merge, use v2.7.x prompt.

## Known risks
- Scene conditions are runtime-only in this stage and are supplied to the scene evaluator rather than being added to YAML config.
- Sprite text currently distributes glyphs evenly across the configured text geometry.
- Renderer tests cover object construction and layout, but full visual verification remains manual in mock mode.
