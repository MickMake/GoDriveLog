# `stepped_fill`

Applies to: bar, segmented.

Status: **candidate / not implemented**.

## What it would do

A bar or segmented fill would advance in visible blocks or steps rather than as a perfectly continuous fill edge.

## What it simulates in real life

Some level indicators are built from discrete segments, blocks, lamps, LCD cells, or LED columns. They cannot show every possible intermediate value; they show the nearest lit step.

This option simulates a block-style or segmented fill display.

## Candidate visual model

A value range is divided into a fixed number of visual steps:

```text
0%   -> no blocks lit
25%  -> first block lit
50%  -> first two blocks lit
75%  -> first three blocks lit
100% -> all blocks lit
```

## Good result

The display clearly looks like a discrete block/segment indicator while still reflecting the source value range sensibly.

## Bad result

The bar jitters between steps, hides meaningful changes for too long, or looks like a rendering bug rather than a deliberate stepped display.

## Design notes

This needs a clear config model before implementation. It may belong in bar rendering, segmented gauge rendering, or a shared display-resolution abstraction depending on final design.
