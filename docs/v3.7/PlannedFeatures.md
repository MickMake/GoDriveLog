# v3.7 Planned Features

This file captures planned or candidate gauge realism work that should not be smuggled into active implementation slices.

It is a holding area for ideas that are useful, plausible, or already documented elsewhere, but which need their own release planning before implementation.

## Status legend

| Marker | Meaning |
| --- | --- |
| ✅ | Implemented or expected to be usable. |
| ❌ | Not planned for that gauge family. |
| 🟡 | Partial, legacy, parse-only, or supported indirectly. |
| ⚠️ | Planned, questionable, in-progress, or needs audit before trusting. |
| 🍺 | Potential candidate / needs beer thought before design or implementation. |

## Gauge realism map

This map is a planning aid only. Do not treat it as implementation truth without checking the current code and completed release docs.

| Realism option | Numeric | Radial | Odometer | Indicator | Bar | Segmented |
| --- | --- | --- | --- | --- | --- | --- |
| `movement` | 🟡 parse only | 🟡 legacy `movement_policy` | ✅ `odometer.movement` | ❌ | 🟡 via damping only | ❌ |
| `wraparound` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `drum_slop` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `carry_drag` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `snap_settle` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `backlash` | ❌ | ❌ | ⚠️ marked implemented / no code | ❌ | ❌ | ❌ |
| `hysteresis` | ❌ | ✅ | ❌ | ❌ | ⚠️ planned / not yet | ❌ |
| `stiction` | ❌ | ✅ | ❌ | ❌ | ⚠️ planned / not yet | ❌ |
| `damping` | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `overshoot` | ❌ | ✅ | ❌ | ❌ | ⚠️ PR open / in progress | ❌ |
| `peg_bounce` | ❌ | ✅ | ❌ | ❌ | ⚠️ planned / not yet | ❌ |
| `thermal_fade` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ✅ | ❌ | 🍺 potential candidate / needs beer thought |
| `per_digit_response_lag` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `leading_zero_behaviour` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `decimal_point_behaviour` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `needle_shadow` | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| `calibration_offset` | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| `segment_bleed` / `digit_bleed` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `ghosting` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `uneven_brightness` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `load_sag` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `stepped_fill` | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought | 🍺 potential candidate / needs beer thought |
| `quantized_fill` | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought | 🍺 potential candidate / needs beer thought |

## Candidate: odometer backlash

`backlash` appears in earlier odometer realism planning and may need a focused implementation/audit slice.

Before implementation, verify the current code state. Do not rely only on old checklists or prompt files.

Treat odometer `backlash` as:

```text
candidate requiring audit before implementation
```

### Candidate behaviour

`backlash` would model direction-change slack for odometer wheels.

Existing odometer realism can create general mechanical feel, but may not fully create direction-change backlash:

| Existing option | Why it is not backlash |
| --- | --- |
| `drum_slop` | Static wheel alignment imperfection; does not care about direction changes. |
| `carry_drag` | Rollover coupling between wheels; not reverse-direction slack. |
| `snap_settle` | Landing/settle effect; not slack when reversing. |
| `movement: linear`, `ease_out`, `bell` | Movement curves; not mechanical play. |
| `wraparound` | Route choice across digit boundaries; not slack. |

If implemented, `backlash` should be its own odometer-only feature.

## Odometer movement cleanup candidates

The reserved odometer movement values should be clarified only inside a focused odometer slice.

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
- `backlash` if implemented.

## Indicator realism scope

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

## Numeric and segmented display realism candidates

Numeric and segmented display realism is plausible, especially for seven-segment-style gauges, but the user-facing model must stay simple enough to explain.

Numeric gauge rendering is image-map driven. Prefer future realism behaviours that operate at the digit-slot or displayed-character level rather than trying to interpret image internals.

Good numeric candidates:

- `thermal_fade`-style character or digit-slot fade;
- `per_digit_response_lag` where digit slots update with a small stagger;
- `leading_zero_behaviour` for blank, dim, or formatted leading-zero slots;
- `decimal_point_behaviour` because DP is overlay-based and should be handled deliberately;
- `uneven_brightness` as a per-digit-slot brightness multiplier;
- `load_sag` where values drawing more lit segments dim the whole display or affected slots.

Candidates needing more design:

- `segment_bleed` / `digit_bleed` using inactive segment masks;
- `ghosting` of previous displayed characters.

### Current-load brightness sag

`load_sag` models an electrical or driver limitation where displays with more lit segments draw more current and therefore appear slightly dimmer.

For a seven-segment display, the visual idea is:

| Displayed value | Approximate lit segment load | Expected brightness |
| --- | --- | --- |
| `111` | low | brighter |
| `777` | medium | normal-ish |
| `888` | high | dimmer |
| `888.8` | very high | dimmest |

This should not require inspecting image internals. A future implementation can use a configured or inferred character load table and apply a brightness multiplier to the whole numeric display, to each digit slot, or to a simple hybrid of both.

Prefer starting with a display-level model because it best matches a shared supply or driver sag: the whole readout dims when total current draw rises.

Possible future config shape:

```yaml
realism:
  load_sag:
    enabled: true
    strength: 0.08
    mode: display
```

Potential modes:

| Mode | Meaning |
| --- | --- |
| `display` | Whole display dims based on total lit-character load. |
| `slot` | Each digit slot dims based on that character's own load. |
| `hybrid` | Whole display sag plus slight per-slot sag. |

### Digit-slot uneven brightness

`uneven_brightness` is still image-map-safe when it is defined at digit-slot level.

It should not mean per-segment brightness variation unless a later display-mask abstraction exists.

A future implementation can apply a stable brightness multiplier per slot:

```text
slot 0: 0.96
slot 1: 1.00
slot 2: 0.91
slot 3: 0.98
```

This works regardless of which glyph image is rendered in the slot.

### Decimal point complication

For ordinary seven-segment digits, an inactive-segment bleed mask could be represented by a faint `8` rendered underneath the active digit.

Decimal points make that less clean. In the current numeric display model, DP is a special overlay rather than part of the normal digit character. A convincing bleed mask for a decimal-capable slot may therefore need:

- a faint `8` mask under the active digit;
- a separate faint decimal-point overlay;
- clear rules for when DP bleed appears;
- a config model that users can understand without knowing the internal renderer layering.

Do not promote bleed/ghosting candidates until the display-mask abstraction and config naming are clear.

## Bar realism scope

Bar gauges are linear fill/reveal gauges. Runtime realism should focus on the displayed fill edge moving toward the target, not on repainting the gauge artwork.

Before planning any bar realism beyond pointer markers, audit the current code and the completed v3.5 docs. Do not contradict completed v3.5 state from backlog notes alone.

Possible future bar candidates:

- `stepped_fill` for block-style bars;
- `quantized_fill` where the bar only visibly changes after the value crosses a display-resolution step;
- focused audits/fixes for already-documented bar realism options if code support is missing.

Both `stepped_fill` and `quantized_fill` need a clear config model before promotion.

## Planning rule

Do not implement anything from this file as part of an unrelated slice.

Promotion should be explicit:

1. choose one small candidate;
2. define its user-facing config;
3. define which gauge families support it;
4. add docs and prompt slice(s);
5. then implement it in a dedicated branch/PR.
