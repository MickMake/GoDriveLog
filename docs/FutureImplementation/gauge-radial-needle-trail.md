# Radial Needle Trail

Index: 3

Status: desired

Area: `gauge/radial`, renderer, animation history

Effort: 4-7 Codex hours

Add optional radial-only `realism.needle_trail` support.

Needle trail renders a bounded history of previous displayed needle positions as fading ghost needles. It is a visual afterimage effect, not a movement curve.

## Proposed config shape

```yaml
realism:
  needle_trail:
    length: 12
    decay_ms: 500
```

## Options

- `length`: maximum number of historical displayed needle positions retained. Default: `12`.
- `decay_ms`: time in milliseconds for retained trail samples to fade out. Default: `500`.

## Rules

- Radial-only.
- Disabled by default.
- Display-only.
- Must not mutate source values, logs, exported values, configured ranges, or input data.
- Store only a bounded history of displayed needle angles/positions and timestamps.
- Trail samples should fade and be discarded deterministically.
- Do not store an unbounded render history.
- Do not place this under `movement`; `movement` selects the travel curve, while `needle_trail` is a render-history effect.

## Possible future slice

```text
v3.5.19 radial needle trail
```
