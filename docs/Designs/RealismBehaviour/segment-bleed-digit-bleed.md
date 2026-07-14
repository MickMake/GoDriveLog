# `digit_bleed`

Applies to: numeric.

Status: **candidate / not implemented**.

## Purpose

Simulate the faint visibility of inactive LED segments beneath an active numeric digit.

The effect is display-only. It does not alter the source value, formatting, logging, exports, calculations, or active digit selection.

## Expected visual behaviour

Each numeric slot renders a faint full-segment digit underneath the active digit.

For a seven-segment display:

- the faint underlay uses the `8` digit image from the same asset set;
- the underlay remains visible beneath the active digit;
- the active digit remains clear and dominant;
- the effect is subtle and must not make the displayed value ambiguous.

The result should resemble the visible inactive segment structure of a physical LED display.

## What it simulates

Real LED displays may reveal faint inactive segment shapes through the lens or surrounding material even when those segments are not illuminated.

This feature recreates that physical display characteristic.

It does not simulate:

- light leaking between neighbouring digits;
- persistence or ghosting;
- brightness variation;
- electrical load effects;
- mechanical behaviour.

## Candidate configuration

Initial implementation should support:

```yaml
realism:
  digit_bleed: true
```

If omitted or false, existing display behaviour is preserved.

No user-configurable alpha value is required for the initial implementation.

## Required behaviour

The implementation must be:

- numeric-only;
- display-only;
- deterministic;
- bounded;
- active only when configured;
- independent of source values.

It must never modify:

- source values;
- calculations;
- logging;
- exports;
- numeric formatting;
- active digit selection.

## Asset behaviour

The bleed layer uses the `8` digit image from the same numeric asset set as the active digit.

The implementation must:

- use the existing `8` image as the complete inactive-segment structure;
- apply a fixed low alpha to the underlay;
- use the same slot geometry and alignment as the active digit;
- preserve the original asset proportions;
- fail validation if the required `8` image is unavailable.

The implementation must not:

- generate segment masks;
- inspect or reverse-engineer image pixels;
- require individual segment assets;
- substitute an asset from another digit set.

## Layer behaviour

The bleed layer must sit immediately beneath the active LED digit layer.

Expected order:

```text
background / housing
digit bleed layer
active LED digit layer
foreground / glass / bezel
```

The bleed layer:

- must not replace the active digit layer;
- must not replace background, housing, glass, bezel, or foreground layers;
- must not reorder unrelated layers;
- must remain confined to the digit slot.

## Behaviour and appearance boundary

Code is responsible for:

- enabling or disabling digit bleed;
- selecting the `8` image from the active digit asset set;
- applying the fixed low alpha;
- placing the underlay directly beneath the active digit.

Images are responsible for:

- segment shape;
- colour;
- glow;
- neighbour glow;
- styling;
- visual identity.

Code must not recreate or reinterpret the segment artwork.

## Decimal points and symbols

Digit bleed applies only to digit slots.

Decimal points, signs, separators, units, and other symbols are excluded from the initial implementation unless they are already part of the digit asset itself.

## Implementation requirements

If confirmed missing:

- add configuration parsing for `realism.digit_bleed`;
- validate that the option is accepted only for numeric gauges;
- validate that the active digit asset set contains an `8` image;
- render the low-alpha `8` underlay beneath each active digit;
- preserve existing behaviour when absent or false;
- keep behaviour deterministic;
- keep implementation local unless an existing helper cleanly fits;
- do not redesign the renderer.

## Tests

Add or update tests for:

- configuration parsing;
- configuration validation;
- numeric-only scope;
- rejection for unsupported gauge types;
- missing `8` asset validation;
- disabled and absent behaviour;
- bleed-layer placement beneath the active digit;
- slot alignment;
- deterministic rendering;
- decimal points and non-digit symbols remaining unaffected.

## Preview and documentation

Add a preview demonstrating:

- digit bleed disabled;
- digit bleed enabled;
- several active digits rendered over the faint `8` underlay;
- correct layer placement;
- unchanged decimal points and symbols.

Update relevant documentation accordingly.

## Constraints

- No source-value mutation.
- No numeric formatting changes.
- No segmented-gauge support in this slice.
- No individual segment assets.
- No generated masks.
- No pixel inspection.
- No configurable alpha in the initial implementation.
- No renderer redesign.
- No layer replacement.
- No unrelated layer reordering.

## Good result

Inactive segment structure is faintly visible beneath the active digit, making the display feel physical while keeping the current value clear and readable.

## Bad result

The bleed layer obscures the active digit, appears above the LED layer, replaces another layer, uses the wrong asset set, misaligns with the active digit, or makes the displayed value ambiguous.

## Non-goals

This design does not define:

- segmented-gauge bleed;
- per-segment rendering;
- ghosting;
- persistence;
- uneven brightness;
- load sag;
- neighbour glow generation;
- configurable bleed intensity;
- renderer architecture.
