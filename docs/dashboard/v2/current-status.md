# Dashboard v2 Current Status

## Active stage
v2.4.x - Asset registry

## Completed
- v2.0.x - New config schema
- v2.1.x - Dashboard config validation
- v2.2.x - Sensor state boundary
- v2.3.x - Decoder engine

## Current branch
- feature/dashboard-v2-asset-registry

## Decisions
- No legacy dashboard compatibility.
- Sensors are separate from dashboard visuals.
- Dashboard asset registry loads and caches local assets only.
- Asset loading remains independent of scene primitives and rendering.

## Next prompt
Use v2.4.x prompt from prompts.md until merged. After merge, use v2.5.x prompt.

## Known risks
- Asset registry currently caches file bytes and metadata only; renderer-specific image decoding belongs in the renderer stage.
- Generated frame pattern support is intentionally small: `{index}` and zero-padded `{index:03}` style markers only.
