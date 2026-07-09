# `damping`

Applies to: radial, bar.

Status: **implemented**.

## What it does

The displayed position lags behind the source value and catches up smoothly.

## What it simulates in real life

Damping simulates resistance that smooths movement: viscous damping in an analogue gauge, mass/inertia in a needle mechanism, electrical filtering, or deliberately damped display electronics.

In practice, the instrument does not instantly follow every input change. It moves calmly toward the new reading.

## Good result

Needles and bars feel weighty and calm during value changes.

## Bad result

The display feels sluggish, never reaches the target, or keeps drifting after the value has settled.
