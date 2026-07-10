# Odometer Movement Cleanup Candidates

Origin: `docs/v3.7/PlannedFeatures.md`

The reserved odometer movement values should be clarified only inside a focused odometer slice.

## `smooth`

Do not implement `movement: smooth` as a separate future movement mode unless a later design gives it a meaning that is genuinely different from existing movement curves.

Current smooth odometer movement is already covered by:

- `movement: linear` — continuous constant roll;
- `movement: ease_out` — continuous roll slowing into target;
- `movement: bell` — continuous slow-fast-slow roll.

## `click`

Do not implement `movement: click` as a separate movement mode unless a later slice defines distinct stepped-wheel behaviour.

Most click-like mechanical feel should come from combinations of existing/required realism options:

- `movement: instant`;
- `drum_slop`;
- `carry_drag`;
- `snap_settle`;
- `backlash` if implemented.
