# `uneven_brightness`

Applies to: numeric.

Status: **candidate / not implemented**.

## Purpose

Simulate the slight brightness differences commonly seen between digit positions on real electronic displays.

The effect is display-only. It never alters the source value, formatting, calculations, logging, exports, or active digit selection.

## Expected visual behaviour

Each digit slot may have a different, but stable, brightness.

Example:

```text
Slot 0: 0.96
Slot 1: 1.00
Slot 2: 0.91
Slot 3: 0.98
```

The same slot always uses the same brightness multiplier regardless of which digit is displayed.

The display should appear naturally aged or imperfect while remaining easy to read.

## What it simulates

Real displays rarely age uniformly.

Brightness differences may result from:

- LED manufacturing tolerances;
- LED ageing;
- VFD phosphor wear;
- LCD contrast variation;
- driver electronics;
- dirty contacts;
- lens tint.

This feature recreates those stable physical differences.

It does **not** simulate:

- random flicker;
- brightness ripple;
- power supply sag;
- ghosting;
- digit bleed;
- response lag.

## Candidate configuration

Initial implementation should support:

```yaml
realism:
  uneven_brightness:
    slots: [0.96, 1.00, 0.91, 0.98]
```

The array is ordered in visual left-to-right slot order.

Each value represents a brightness multiplier for that slot.

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

## Slot behaviour

Each digit slot owns a fixed brightness multiplier.

The multiplier remains constant regardless of:

- displayed digit;
- response lag;
- ghosting;
- leading-zero behaviour.

Brightness belongs to the slot, not the glyph.

## Decimal point behaviour

The decimal-point overlay inherits the brightness multiplier of the digit slot to which it belongs.

Although implemented as a separate overlay, it behaves as though it is part of the original numeric digit.

## Layer behaviour

The slot brightness multiplier applies to:

- active digit;
- ghost glyph;
- decimal-point overlay.

It does **not** apply to:

- digit bleed;
- background;
- housing;
- bezel;
- glass;
- unrelated overlays.

Digit bleed represents inactive physical structure rather than illuminated content.

## Configuration rules

Brightness values must be within:

```text
0.0 – 1.0
```

Values outside this range fail validation.

If fewer slot values are supplied than required:

- remaining slots default to `1.0`.

If additional slot values are supplied:

- extra values are ignored;
- a warning should be emitted;
- rendering continues normally.

## Behaviour and appearance boundary

Code is responsible for:

- assigning slot brightness;
- validating configuration;
- applying brightness multipliers;
- ensuring deterministic behaviour.

Images are responsible for:

- digit appearance;
- colour;
- glow;
- styling;
- visual identity.

Code must not modify artwork beyond applying the configured brightness multiplier.

## Implementation requirements

If confirmed missing:

- add configuration parsing;
- validate numeric-only usage;
- validate brightness ranges;
- support per-slot brightness multipliers;
- support decimal-point inheritance;
- ignore excess configuration values with a warning;
- default unspecified slots to `1.0`;
- preserve existing behaviour when absent;
- keep behaviour deterministic;
- keep implementation local unless an existing helper cleanly fits;
- do not redesign the renderer.

## Tests

Add or update tests for:

- configuration parsing;
- configuration validation;
- numeric-only scope;
- slot multiplier application;
- decimal-point inheritance;
- missing slot defaults;
- excess slot warnings;
- invalid value rejection;
- deterministic behaviour;
- disabled behaviour.

## Preview and documentation

Add previews demonstrating:

- uniform brightness;
- uneven brightness;
- decimal-point inheritance;
- missing slot defaults;
- readable aged displays.

Update relevant documentation accordingly.

## Constraints

- No source-value mutation.
- No numeric formatting changes.
- No segmented-gauge support in this slice.
- No runtime randomness.
- No automatic brightness generation.
- No renderer redesign.

## Good result

The display appears naturally imperfect through subtle, stable brightness variation while remaining fully readable.

## Bad result

Brightness changes randomly, affects numeric meaning, causes flicker, alters inactive layers, or makes the display difficult to read.

## Non-goals

This design does not define:

- segmented-gauge brightness;
- brightness ripple;
- flicker;
- power supply sag;
- ghosting;
- digit bleed;
- response lag;
- renderer architecture.