# `uneven_brightness`

Applies to: numeric, segmented.

Status: **candidate / not implemented**.

## What it would do

Digit slots or display regions would render with small, stable brightness differences.

For numeric gauges, the simplest useful model is per-slot brightness rather than per-segment brightness.

## What it simulates in real life

Real displays rarely age perfectly evenly. Individual digit positions may be slightly dimmer or brighter due to LED binning, lamp ageing, lens tint, driver variance, dirty contacts, VFD wear, or LCD contrast differences.

This option simulates stable, deterministic brightness variation across the display.

## Candidate visual model

```text
slot 0: 0.96
slot 1: 1.00
slot 2: 0.91
slot 3: 0.98
```

The same slot should keep the same brightness multiplier regardless of which glyph is shown.

## Good result

The display feels slightly aged or imperfect while remaining readable.

## Bad result

Brightness varies randomly every frame, makes digits unreadable, changes source values, or requires image-internal segment analysis before the abstraction exists.

## Design notes

Start with digit-slot brightness. Per-segment brightness can wait until there is a clear display-mask or segment-level abstraction.
