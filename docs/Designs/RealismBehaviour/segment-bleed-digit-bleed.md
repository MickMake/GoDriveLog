# `segment_bleed` / `digit_bleed`

Applies to: numeric, segmented.

Status: **candidate / needs design**.

## What it would do

Inactive segments or digit areas would remain faintly visible underneath the active displayed character.

For ordinary seven-segment digits, this might look like a very faint inactive `8` mask behind the active glyph. For other digit artwork, it may need a dedicated inactive-segment asset.

## What it simulates in real life

Real displays often have visible inactive elements:

- unlit LED segments may still be visible under the lens;
- LCD segments may leave a faint outline;
- gas-discharge or VFD elements may glow or bleed slightly;
- old display masks may show the full segment structure even when inactive.

This option simulates that inactive-element visibility.

## Candidate visual model

Possible approaches:

| Approach | Meaning |
|---|---|
| `digit_bleed` | Render a faint full-digit mask under the active glyph. |
| `segment_bleed` | Render only inactive segment masks when the display has segment-level assets. |
| `slot_bleed` | Render a faint per-slot background or inactive display shape independent of glyph internals. |

## Good result

Inactive display structure is faintly visible and helps the display feel physical without hurting readability.

## Bad result

Bleed makes the active digit ambiguous, requires the renderer to reverse-engineer image internals, or creates artefacts that look like wrong numbers.

## Design notes

Do not promote this until the display-mask abstraction and config naming are clear. Decimal points complicate the model because they may need separate bleed handling from the main digit body.
