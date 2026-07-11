# `leading_zero_behaviour`

Applies to: numeric, segmented.

Status: **candidate / not implemented**.

## What it would do

Leading zero slots could be shown, hidden, dimmed, blanked, or otherwise treated deliberately instead of being an accidental formatting side effect.

## What it simulates in real life

Many real numeric displays handle leading zeroes differently depending on their technology and purpose:

- mechanical counters may show every wheel, including leading zeroes;
- automotive digital displays may blank leading zeroes for readability;
- seven-segment displays may leave unused high-order digits dark;
- worn or low-cost displays may show dim inactive slots.

This option simulates that deliberate display behaviour.

## Candidate visual model

Possible modes:

| Mode | Meaning |
|---|---|
| `show` | Leading zeroes are displayed normally. |
| `blank` | Leading zero slots are empty/off. |
| `dim` | Leading zeroes remain visible but subdued. |
| `placeholder` | Leading slots show a configured placeholder or faint slot state. |

## Good result

Leading zero handling feels intentional and consistent with the display technology.

## Bad result

The display changes numeric meaning, breaks alignment unexpectedly, hides significant zeroes, or treats ordinary formatting as realism without a clear model.

## Design notes

Keep this display-only. The source value and formatted/exported value should remain unchanged unless a separate formatting layer explicitly says otherwise.
