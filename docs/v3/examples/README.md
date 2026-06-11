# GoDriveLog dashboard examples

These are draft GoDriveLog v3 dashboard configuration examples for sleeping on, arguing with tomorrow, and hopefully not summoning a YAML demon.

They are active v3 examples, not decorative confetti. Every file in this directory should validate against the same v3 schema rules as `docs/v3/config.example.yaml` and `docs/v3/config.full.yaml`.

Asset paths are repository-root relative. Use `assets/dashboard/...` paths in active v3 examples.

The key model used here is:

```text
asset background
+ value/state-driven dynamic layer
+ optional foreground/glass/bezel overlay
= photoreal widget
```

The examples assume the v3 direction:

```text
vehicles
sensors
assets
logs
dashboards
```

Sensors own polling. Logs and dashboards subscribe to sensor events.

## Files

- `simple_speed_warning.yaml` — deliberately small first-slice example: image + digit displays + indicator
- `nissan_300zx_z31_inspired.yaml` — retro-inspired richer example using digit, bar, frame, indicator, and image assets
- `honda_s2000_inspired.yaml` — retro-inspired richer example using digit, bar, frame, indicator, and image assets

## Notes

- `background` is normal/default for photoreal assets, but not mandatory.
- `foreground` is optional and intended for glass, bezel, scratches, reflections, or lens effects.
- Widgets use `position`, not `rect`, for native-size image assets.
- Asset packs own visual geometry such as spacing and frame counts.
- Widget config owns binding: sensor, position, format, min/max mapping, and asset reference.
- Use `characters`, not `digits`, for digit set character maps.
- Keep sensor examples to documented sensor types. First-slice v3 starts with `type: obd`.
- No dashboard scripting, formulas, inheritance, or clever little config goblins.
