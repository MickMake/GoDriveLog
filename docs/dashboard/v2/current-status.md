# Dashboard v2 Current Status

## Active stage
v2.5.x - Scene primitives

## Completed
- v2.0.x - New config schema
- v2.1.x - Dashboard config validation
- v2.2.x - Sensor state boundary
- v2.3.x - Decoder engine
- v2.4.x - Asset registry

## Current branch
- feature/dashboard-v2-scene-primitives

## Decisions
- No legacy dashboard compatibility.
- Sensors are separate from dashboard visuals.
- Dashboard asset registry loads and caches local assets only.
- Asset loading remains independent of rendering.
- Scene primitives evaluate configured dashboard blocks and layers into scene elements without doing Fyne rendering.

## Next prompt
Use v2.5.x prompt from prompts.md until merged. After merge, use v2.6.x prompt.

## Known risks
- Scene conditions are runtime-only in this stage and are supplied to the scene evaluator rather than being added to YAML config.
- Renderer-specific image decoding belongs in the renderer stage.
- Generated frame pattern support is intentionally small: index markers only.
