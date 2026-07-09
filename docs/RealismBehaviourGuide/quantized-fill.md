# `quantized_fill`

Applies to: bar, segmented.

Status: **candidate / not implemented**.

## What it would do

A bar or segmented display would only visibly change after the value crosses a display-resolution threshold.

Unlike `stepped_fill`, which implies visibly block-style rendering, `quantized_fill` can still look like a continuous bar but update in discrete increments.

## What it simulates in real life

Real instruments often have limited display resolution. A digital bar graph, LCD level display, or low-resolution driver may only represent a finite number of positions even if the source value changes smoothly.

This option simulates display quantisation: the instrument only has so many visible states.

## Candidate visual model

```text
source value: 42.1 -> rendered as 42
source value: 42.4 -> rendered as 42
source value: 42.6 -> rendered as 43
```

For a bar, this means the displayed fill/reveal extent snaps to the nearest configured display increment.

## Good result

The gauge looks like a real limited-resolution display while remaining stable and predictable.

## Bad result

The indicator jitters rapidly around thresholds, loses too much useful detail, or changes exported/source values instead of only the rendered display.

## Design notes

A future implementation should define whether quantisation happens before or after other realism effects such as damping, overshoot, and pointer markers. Pointer markers should continue to observe the final rendered indicator position.
