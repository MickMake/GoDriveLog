# GoDriveLog dashboard examples 0.1

These are draft GoDriveLog v3 dashboard configuration examples for sleeping on, arguing with tomorrow, and hopefully not summoning a YAML demon.

They are **inspired by** retro digital dashboards from the Hagerty UK article "12 of our favourite digital dashboards", especially the Nissan 300ZX Z31 and Honda S2000 examples. They are not exact copies.

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

- `nissan_300zx_z31_inspired.yaml`
- `honda_s2000_inspired.yaml`

## Notes

- `background` is normal/default for photoreal assets, but not mandatory.
- `foreground` is optional and intended for glass, bezel, scratches, reflections, or lens effects.
- Widgets use `position`, not `rect`, for native-size image assets.
- Asset packs own visual geometry such as spacing and frame counts.
- Widget config owns binding: sensor, position, format, min/max mapping, and asset reference.
- No dashboard scripting, formulas, inheritance, or clever little config goblins.
