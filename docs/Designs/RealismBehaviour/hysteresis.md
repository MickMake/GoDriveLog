# `hysteresis`

Applies to: radial, bar.

Status: **implemented**.

## What it does

The displayed position may rest slightly differently depending on whether the value approached from below or above.

For radial gauges, this means a small direction-dependent angle offset. For bar gauges, this means a small direction-dependent fill/reveal extent offset.

## What it simulates in real life

Hysteresis appears when the output of a mechanism depends partly on its recent history, not only the current input. In real gauges this can come from friction, linkage play, magnetic effects, spring behaviour, or general mechanical reluctance.

A rising value and a falling value can therefore indicate slightly different positions even when the underlying source value is the same.

## Good result

A rising value and falling value can settle with a tiny direction-dependent display difference, while still clearly representing the same source value.

## Bad result

The displayed offset is large enough to look wrong, changes the source value, or accumulates error over time.
