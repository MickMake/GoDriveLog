# Radial Movement Options

Index: 2

Status: desired

Area: `gauge/radial`, movement policy, runtime animation

Effort: 3-5 Codex hours

Radial gauges should eventually support the scalar `movement` options that already exist for gauge movement selection, while preserving current behaviour as the compatibility default.

## Proposed movement meanings for radial gauges

- `instant`: current radial behaviour; immediately render the needle at the target angle with no interpolation.
- `linear`: interpolate the displayed needle angle from the previous displayed angle to the target angle at constant progress.
- `bell`: interpolate with a slow start, faster middle, and slow end.

## Rules

- `instant` must preserve existing radial semantics.
- Movement must be display-only.
- Movement must animate displayed angle/position only; it must not mutate source values, logs, exported values, configured ranges, or input data.
- Do not pre-render or cache unbounded intermediate needle images.
- Prefer small per-gauge transition state such as previous angle, target angle, elapsed time, duration, movement mode, and active/inactive state.
- Keep needle geometry and image assets reusable; rotate or transform at render time rather than generating a frame cache.
- Do not combine this with damping, stiction, overshoot, peg bounce, needle trail, or peak hold unless a later slice explicitly defines composition.

## Possible future slice

```text
v3.5.x radial movement options
```
