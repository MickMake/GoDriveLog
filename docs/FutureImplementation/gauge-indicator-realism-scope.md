# Indicator Realism Scope

Origin: `docs/v3.7/PlannedFeatures.md`

Indicator gauges are image-state driven. The `off` and `on` image layers define the static lamp appearance.

Runtime realism should stay transition-focused. `thermal_fade` already supports separate rise and fall timing:

```yaml
realism:
  thermal_fade:
    rise_ms: 120
    fall_ms: 240
```

Use `rise_ms` for off-to-on warm-up and `fall_ms` for on-to-off cool-down.

Do not add separate planned runtime features for weak bulb, tint, ageing, bloom, dirty lens, or uneven illumination unless a later design explicitly introduces additional indicator image states or display-layer effects. Those qualities belong in the supplied images.
