# `per_digit_response_lag`

Applies to: numeric, segmented.

Status: **candidate / not implemented**.

## What it would do

Digit slots would not all update at exactly the same instant. Each slot may lag the source value by a tiny, controlled amount so that a multi-digit display appears to settle across the display rather than snapping as one perfectly synchronised image.

## What it simulates in real life

Some digital displays and driver circuits update digits sequentially or with slightly uneven response. Older electronics, multiplexed displays, and mechanical/electromechanical readouts can show small timing differences between slots.

This option simulates that slot-level response lag without changing the source value.

## Candidate visual model

```text
slot 0 updates first
slot 1 follows after a small delay
slot 2 follows after another small delay
```

The delay should be subtle. This is display character, not slot-machine theatrics.

## Good result

The display feels like a real driven display with tiny timing differences, while still becoming readable quickly and settling exactly on the correct value.

## Bad result

Digits lag so much that the reading is confusing, update in random order, produce impossible long-lived values, or mutate logs/exported/source data.

## Design notes

Prefer defining lag per digit slot rather than trying to inspect glyph images. The renderer already knows which slot is being updated; that is the right abstraction layer.
