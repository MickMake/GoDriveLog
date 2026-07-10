# Value Zones / Warning-Danger Assets

Index: 5

Status: desired

Area: `gauge/assets`, renderer, config validation

Effort: 4-7 Codex hours

Support optional value zones that select warning/danger variants of gauge assets when the source value reaches a configured range.

This should be a separate gauge-display feature, not part of `realism.overshoot`.

## Proposed config shape

```yaml
zones:
  warning:
    min: 6000
    max: 7000
  danger:
    min: 7000
    max: 8000
```

## Proposed asset convention

```text
needle.png
needle_warning.png
needle_danger.png
face.png
face_warning.png
face_danger.png
bar.png
bar_warning.png
bar_danger.png
```

## Rules

- Zone selection should follow the real/source target value, not any temporary animated display value.
- If a zone-specific asset exists for a layer, use it.
- If a zone-specific asset does not exist, fall back to the normal asset.
- Overshoot may visually pass a threshold, but should not change the zone state unless the real/source value is in that zone.
- Avoid surprising behaviour where a temporary animation makes the gauge appear to enter warning or danger falsely.

## Possible future slice

```text
v3.5.x value zones / warning-danger assets
```
