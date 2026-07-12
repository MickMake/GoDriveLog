# Numeric Display

## Purpose

The Numeric Display renders formatted numeric values using reusable digit artwork.

It owns the display structure for digits, signs, decimal points, formatting and slot placement. Optional realism behaviours may alter how those rendered elements appear or respond, but they do not create separate display types.

## Scope

The Numeric Display is responsible for:

- formatting a source value for display;
- allocating the formatted characters to display slots;
- rendering digit glyphs;
- rendering a minus sign where required;
- rendering the decimal point;
- preserving the numeric meaning of the formatted value;
- defining the rendering boundary used by optional numeric-display realism behaviours.

Leading-zero behaviour, ghosting, load sag, bleed, uneven brightness and per-digit response lag remain separate realism behaviours. They act upon the Numeric Display rather than redefining its basic architecture.

## Digit artwork

Digit artwork is supplied as reusable assets for:

```text
0 1 2 3 4 5 6 7 8 9
-
```

The display should not require duplicate digit assets merely to represent punctuation or state variations that can be rendered independently.

## Formatting and slot allocation

Formatting determines the textual value presented by the Numeric Display.

Slot allocation determines which visible character belongs in each display position.

A decimal point does not consume a digit slot. It is associated with the preceding numeric character and rendered as an overlay within that character's slot.

The renderer must preserve the numeric meaning of the formatted value. It must not:

- move a decimal point to a different numeric position;
- drop a required decimal point;
- show a decimal point when the format does not require one;
- allow visual effects to change the represented value.

## Decimal-point overlay

The decimal point is part of the Numeric Display renderer.

It is not embedded within every digit glyph. Instead, it is rendered as a separate overlay positioned relative to the active digit.

The overlay exists to reduce artwork duplication, not to create an independent rendering pipeline or a separate realism feature.

### Design rationale

Embedding the decimal point into each digit would require duplicate artwork:

```text
0 1 2 3 4 5 6 7 8 9
0. 1. 2. 3. 4. 5. 6. 7. 8. 9.
```

Using an overlay requires only:

```text
Digits:
0 1 2 3 4 5 6 7 8 9
-

Overlay:
decimal point
```

This approach:

- reduces the number of assets that must be created and maintained;
- keeps the decimal-point appearance consistent;
- permits decimal-point placement to be controlled by the renderer;
- supports different numeric formats without duplicate glyph sets;
- allows future visual behaviour without replacing digit artwork.

### Rendering behaviour

The renderer determines:

- whether the formatted value requires a decimal point;
- which digit slot owns the decimal point;
- the decimal point's position within that slot;
- whether the decimal-point asset is available;
- the order in which the decimal point and other slot layers are composed.

The decimal point remains part of the displayed numeric value. It is not a separate gauge, widget or display type.

### Asset requirements

A Numeric Display that uses a format containing a decimal point requires a decimal-point overlay asset.

A format that does not contain a decimal point must not require that asset.

Missing required artwork should fail validation clearly rather than silently changing the displayed value.

## Interaction with realism behaviours

The decimal-point overlay may participate in numeric-display realism behaviours where those behaviours explicitly define the interaction.

Examples include:

- ghosting;
- segment or digit bleed;
- uneven brightness;
- load sag;
- per-digit response lag.

A future behaviour may treat the decimal-point overlay differently from the digit glyph where that difference represents the physical display technology. Examples could include:

- independent fade timing;
- faint inactive illumination;
- brightness mismatch;
- slight positional variation;
- delayed illumination or extinction.

Those effects belong to the relevant realism behaviour. They do not make `decimal_point_behaviour` a separate feature or configuration namespace.

## Constraints

- The decimal point must remain visibly associated with the correct digit.
- Overlay rendering must be deterministic.
- The overlay must not consume a digit slot.
- The overlay must not alter the underlying source value.
- The overlay must not require duplicate digit artwork.
- Optional visual effects must not make the numeric meaning ambiguous.
- The renderer must degrade safely when a decimal point is not required.
- Missing required assets must produce a clear validation error.

## Good result

The display shows the correct formatted value using a compact, reusable asset set. Decimal points are consistently positioned and remain visually part of the number, while optional realism behaviours can affect them without duplicating digit artwork.

## Bad result

The decimal point consumes a digit position, becomes detached from its number, changes the apparent value, requires duplicate glyph sets, disappears despite being required by the format, or evolves into a second rendering system with its own tiny civil service.

## Non-goals

This design does not define:

- leading-zero policy;
- ghosting behaviour;
- load-sag behaviour;
- bleed behaviour;
- uneven-brightness behaviour;
- per-digit timing;
- numeric-display power lifecycle.

Those remain separate behaviour designs and may reference this display architecture where required.
