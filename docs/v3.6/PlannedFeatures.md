# v3.6 Planned Features

This file captures planned or candidate gauge realism work that should not be smuggled into the active implementation slices.

It is a holding area for ideas that are useful, plausible, or already documented elsewhere, but which need their own release planning before implementation.

## Status legend

| Marker | Meaning |
| --- | --- |
| `implemented` | Code support exists and is expected to be usable. |
| `partial` | Some parsing, validation, or runtime behaviour exists, but the feature is not complete or not the preferred long-term model. |
| `planned / not yet` | Intended by current planning docs, but not implemented yet. |
| `marked implemented / no code` | Documentation or checklist says the feature is complete, but code support could not be found. |
| `potential candidate / needs beer thought` | Plausible future feature, but config model, rendering semantics, or user-facing explanation need more design. |
| `not planned` | Not currently considered suitable for that gauge family. |

## Gauge realism map

| Realism option | Numeric | Radial | Odometer | Indicator | Bar | Segmented |
| --- | --- | --- | --- | --- | --- | --- |
| `movement` | partial: parse only | partial: existing `realism.movement_policy` finite-movement selector | implemented: `odometer.movement` | not planned | partial: finite movement currently tied to damping/related behaviour | not planned |
| `wraparound` | not planned | not planned | implemented | not planned | not planned | not planned |
| `drum_slop` | not planned | not planned | implemented | not planned | not planned | not planned |
| `carry_drag` | not planned | not planned | implemented | not planned | not planned | not planned |
| `snap_settle` | not planned | not planned | implemented | not planned | not planned | not planned |
| `backlash` | not planned | not planned | marked implemented / no code; promoted to v3.6.9 | not planned | not planned | not planned |
| `hysteresis` | not planned | implemented | not planned | not planned | planned / not yet | not planned |
| `stiction` | not planned | implemented | not planned | not planned | planned / not yet | not planned |
| `damping` | not planned | implemented | not planned | not planned | implemented | not planned |
| `overshoot` | not planned | implemented | not planned | not planned | planned / in progress | not planned |
| `peg_bounce` | not planned | implemented | not planned | not planned | planned / not yet | not planned |
| `thermal_fade` | potential candidate / needs beer thought | not planned | not planned | implemented | not planned | potential candidate / needs beer thought |
| `needle_shadow` | not planned | implemented | not planned | not planned | not planned | not planned |
| `calibration_offset` | not planned | implemented | not planned | not planned | not planned | not planned |
| `segment_bleed` / `digit_bleed` | potential candidate / needs beer thought | not planned | not planned | not planned | not planned | potential candidate / needs beer thought |
| `ghosting` | potential candidate / needs beer thought | not planned | not planned | not planned | not planned | potential candidate / needs beer thought |
| `uneven_brightness` | potential candidate / needs beer thought | not planned | not planned | not planned | not planned | potential candidate / needs beer thought |

## Promoted tail slice: odometer backlash

`backlash` is listed as a planned/approved odometer realism behaviour in the v3.5 documentation, and the v3.5 implementation checklist previously marked the slice complete.

However, repository archaeology found only docs/prompt commits for v3.5.15 backlash. No implementation PR, branch, parser field, validator support, runtime behaviour, tests, or preview fixture were found.

Treat odometer `backlash` as:

```text
required missing implementation
```

This is now promoted into the v3.6 tail as:

```text
docs/v3.6/prompts/v3.6.9-implement-odometer-backlash-cleanup.md
```

### Required behaviour

`backlash` must model direction-change slack for odometer wheels.

Existing odometer realism can create general mechanical feel, but it cannot fully create direction-change backlash:

| Existing option | Why it is not backlash |
| --- | --- |
| `drum_slop` | Static wheel alignment imperfection; does not care about direction changes. |
| `carry_drag` | Rollover coupling between wheels; not reverse-direction slack. |
| `snap_settle` | Landing/settle effect; not slack when reversing. |
| `movement: linear`, `ease_out`, `bell` | Movement curves; not mechanical play. |
| `wraparound` | Route choice across digit boundaries; not slack. |

Therefore `backlash` should be implemented as its own odometer-only feature.

## Odometer movement cleanup decisions

The reserved odometer movement values should be cleaned up while implementing `backlash`.

### `smooth`

Do not implement `movement: smooth` as a separate future movement mode unless a later design gives it a meaning that is genuinely different from existing movement curves.

Current smooth odometer movement is already covered by:

- `movement: linear` — continuous constant roll;
- `movement: ease_out` — continuous roll slowing into target;
- `movement: bell` — continuous slow-fast-slow roll.

### `click`

Do not implement `movement: click` as a separate movement mode unless a later slice defines distinct stepped-wheel behaviour.

Most click-like mechanical feel should come from combinations of existing/required realism options:

- `movement: instant`;
- `drum_slop`;
- `carry_drag`;
- `snap_settle`;
- `backlash` once implemented.

## Numeric and segmented display realism candidates

Numeric and segmented display realism is plausible, especially for seven-segment-style gauges, but the user-facing model must stay simple enough to explain.

Candidate ideas:

- `thermal_fade`-style character or segment fade;
- `segment_bleed` / `digit_bleed` using inactive segment masks;
- `ghosting` of previous displayed characters;
- `uneven_brightness` or static per-slot display imperfection.

Current status:

```text
potential candidate / needs beer thought
```

### Decimal point complication

For ordinary seven-segment digits, an inactive-segment bleed mask could be represented by a faint `8` rendered underneath the active digit.

Decimal points make that less clean. In the current numeric display model, DP is a special overlay rather than part of the normal digit character. A convincing bleed mask for a decimal-capable slot may therefore need:

- a faint `8` mask under the active digit;
- a separate faint decimal-point overlay;
- clear rules for when DP bleed appears;
- a config model that users can understand without knowing the internal renderer layering.

Do not promote these numeric/segmented realism candidates until the display-mask abstraction and config naming are clear.

## Planning rule

Do not implement anything from this file as part of an unrelated slice.

Promotion should be explicit:

1. choose one small candidate;
2. define its user-facing config;
3. define which gauge families support it;
4. add docs and prompt slice(s);
5. then implement it in a dedicated branch/PR.
