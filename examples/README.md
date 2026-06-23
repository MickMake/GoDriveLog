# GoDriveLog dashboard examples

These are active GoDriveLog dashboard examples, not versioned planning confetti. Every runnable fixture in this directory should validate against the current v3 config/runtime path.

Asset paths are repository-root relative.

## Files

- `baseline-dashboard.yaml` - reusable baseline dashboard used for renderer comparison.
- `simple_speed_warning.yaml` - deliberately small first-slice example: image + digit displays + indicator.
- `nissan_300zx_z31_inspired.yaml` - retro-inspired richer example using digit, bar, frame, indicator, and image assets.
- `honda_s2000_inspired.yaml` - retro-inspired richer example using digit, bar, frame, indicator, and image assets.
- `dashboards/framework-smoke.yaml` - generated v3.4 smoke-test dashboard proving the deterministic example-asset pipeline.

## Reusable assets

Reusable active assets live under the repository-root `assets/` tree.

Generated v3.4 example dashboard assets live under `examples/assets/v3.4/`, with matching configs under `examples/dashboards/`.

Versioned docs should reference examples and shared assets instead of carrying active runnable copies. Docs explain a slice; examples are the fixtures. This keeps the project from repeatedly moving the same cheese while the mouse files a complaint.

## Notes

- `background` is normal/default for photoreal assets, but not mandatory.
- `foreground` is optional and intended for glass, bezel, scratches, reflections, or lens effects.
- Widgets use `position`, not `rect`, for native-size image assets.
- Asset packs own visual geometry such as spacing and frame counts.
- Widget config owns binding for classic widgets.
- Gauge packages own sensor binding for `type: gauge` widgets.
- Use `characters`, not `digits`, for digit set character maps.
- Keep sensor examples to documented sensor types.
- No dashboard scripting, formulas, inheritance, or clever little config goblins.
