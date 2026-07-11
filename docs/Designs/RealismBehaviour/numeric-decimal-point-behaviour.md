# `decimal_point_behaviour`

Applies to: numeric, segmented.

Status: **candidate / not implemented**.

## What it would do

Decimal points would be treated as deliberate display elements with their own visibility, fade, bleed, ghosting, and positioning rules.

## What it simulates in real life

On real seven-segment and segmented displays, the decimal point is often a separate LED, lamp, mask, or overlay. It may not behave exactly like the main digit segments. It can be brighter, dimmer, slightly offset, slower to fade, or visibly separate from the digit body.

This option simulates decimal-point-specific display behaviour.

## Candidate visual model

Possible behaviours:

- decimal point turns on/off independently from the digit glyph;
- decimal point has its own fade timing;
- inactive decimal point may have faint bleed;
- decimal point may remain visible during digit ghosting;
- decimal point rules may differ for fixed-point and floating display formats.

## Good result

Decimal points look like real display elements rather than incidental punctuation painted into a glyph.

## Bad result

The decimal point disappears unpredictably, changes numeric meaning, bleeds in a confusing way, or forces the renderer to know too much about image internals.

## Design notes

In the current numeric display model, decimal point handling is overlay-based rather than just part of the digit character. A future implementation should preserve that explicit layering instead of burying decimal point behaviour inside glyph images.
