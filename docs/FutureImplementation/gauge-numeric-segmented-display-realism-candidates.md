# Numeric and Segmented Display Realism Candidates

Origin: `docs/v3.7/PlannedFeatures.md`

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

## Current-load brightness sag

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

## Digit-slot uneven brightness

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

## Decimal point complication

For ordinary seven-segment digits, an inactive-segment bleed mask could be represented by a faint `8` rendered underneath the active digit.

Decimal points make that less clean. In the current numeric display model, DP is a special overlay rather than part of the normal digit character. A convincing bleed mask for a decimal-capable slot may therefore need:

- a faint `8` mask under the active digit;
- a separate faint decimal-point overlay;
- clear rules for when DP bleed appears;
- a config model that users can understand without knowing the internal renderer layering.

Do not promote bleed/ghosting candidates until the display-mask abstraction and config naming are clear.
